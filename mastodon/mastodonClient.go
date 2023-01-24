package mastodon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattn/go-mastodon"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/network"
	"github.com/therecipe/qt/webengine"
	"mastodon-client/files"
	"mastodon-client/global"
	"mastodon-client/pushnotifications"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const clientName string = "ubuntu-touch-mastodon"
const clientScopes string = "read write follow push"

type UserAction int

const (
	IsLoggingIn UserAction = iota
	IsRegistering
)

type ServerResponse struct {
	Domain           string   `json:"domain"`
	Version          string   `json:"version"`
	Description      string   `json:"description"`
	Languages        []string `json:"languages"`
	Region           string   `json:"region"`
	Categories       []string `json:"categories"`
	ProxiedThumbnail string   `json:"proxied_thumbnail"`
	TotalUsers       int      `json:"total_users"`
	LastWeekUsers    int      `json:"last_week_users"`
	ApprovalRequired bool     `json:"approval_required"`
	Language         string   `json:"language"`
	Category         string   `json:"category"`
}

type QServer struct {
	core.QObject
	_ string   `property:"domain"`
	_ string   `property:"description"`
	_ bool     `property:"selected"`
	_ int      `property:"totalUsers"`
	_ string   `property:"language"`
	_ bool     `property:"matchingSearchTerm"`
	_ []string `property:"serverRules"`
	_ bool     `property:"requiresApproval"`
}

type QClient struct {
	core.QAbstractListModel

	_ func() `constructor:"init"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*QServer               `property:"servers"`
	_ *QServer                 `property:"selectedServer"`
	_ bool                     `property:"awaitingActivation"`
	_ UserAction               `property:"userAction"`
	_ string                   `property:"webviewUrl"`

	_ func(*QServer)                                                                                 `slot:"addServer"`
	_ func(row int)                                                                                  `slot:"setSelected"`
	_ func(displayName string, username string, email string, password string, reason string) string `slot:"createAccount"`
	_ func()                                                                                         `slot:"setLoginUrl"`
	_ func()                                                                                         `slot:"setIsLoggingIn"`
	_ func()                                                                                         `slot:"setIsRegistering"`
	_ func() bool                                                                                    `slot:"getIsLoggingIn"`
	_ func() bool                                                                                    `slot:"getIsRegistering"`
	_ func(authCode string)                                                                          `slot:"handleAuthCode"`
	_ func(searchTerm string)                                                                        `slot:"filterServers"`
	_ func()                                                                                         `slot:"getServerRules"`
	_ func()                                                                                         `slot:"resendConfirmationEmail"`
	_ func(v0 *webengine.QQuickWebEngineProfile)                                                     `slot:"setProfile"`
	_ func() bool                                                                                    `slot:"shouldSkipSelection"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type CreateAccountResponse struct {
	TokenResponse
	Error string `json:"error"`
}

type PushApplication struct {
	*mastodon.Application
	HasPushSubscription bool `json:"has_push_subscription"`
}

type chosenServer struct {
	Domain string `json:"domain"`
}

const (
	Domain = int(core.Qt__UserRole) + 1<<iota
	Description
	Selected
	TotalUsers
	Language
	MatchingSearchTerm
	ServerRules
	RequiresApproval
)

var listRoles = map[int]*core.QByteArray{
	Domain:             core.NewQByteArray2("domain", -1),
	Description:        core.NewQByteArray2("description", -1),
	Selected:           core.NewQByteArray2("selected", -1),
	TotalUsers:         core.NewQByteArray2("totalUsers", -1),
	Language:           core.NewQByteArray2("language", -1),
	MatchingSearchTerm: core.NewQByteArray2("matchingSearchTerm", -1),
	ServerRules:        core.NewQByteArray2("serverRules", -1),
	RequiresApproval:   core.NewQByteArray2("requiresApproval", -1),
}

var userAccessToken *string

func (m *QClient) init() {
	m.SetRoles(listRoles)
	m.ConnectData(m.data)
	m.ConnectRowCount(m.rowCount)
	m.ConnectRoleNames(m.roleNames)
	m.ConnectAddServer(m.addServer)
	m.ConnectSetSelected(m.setSelected)
	m.ConnectCreateAccount(m.createAccount)
	m.ConnectSetLoginUrl(m.setLoginUrl)
	m.ConnectSetIsLoggingIn(m.setIsLoggingIn)
	m.ConnectSetIsRegistering(m.setIsRegistering)
	m.ConnectGetIsLoggingIn(m.getIsLoggingIn)
	m.ConnectGetIsRegistering(m.getIsRegistering)
	m.ConnectHandleAuthCode(m.handleAuthCode)
	m.ConnectFilterServers(m.filterServers)
	m.ConnectGetServerRules(m.getServerRules)
	m.ConnectResendConfirmationEmail(m.resendConfirmationEmail)
	m.ConnectSetProfile(m.setProfile)
	m.ConnectShouldSkipSelection(m.shouldSkipSelection)
	m.SetAwaitingActivationDefault(true)
}

func (m *QClient) rowCount(parent *core.QModelIndex) int {
	return len(m.Servers())
}

func (m *QClient) shouldSkipSelection() bool {
	fileBytes, err := os.ReadFile(global.ConfigFileDir + "chosenServer")
	if err != nil {
		fmt.Println(err)
		return false
	}
	chosenServer := chosenServer{}
	err = json.Unmarshal(fileBytes, &chosenServer)
	dom := "https://" + chosenServer.Domain
	m.SetWebviewUrl(dom)
	m.WebviewUrlChanged(dom)
	return true
}

func (m *QClient) setProfile(v0 *webengine.QQuickWebEngineProfile) {
	store := v0.CookieStore()
	store.ConnectCookieAdded(func(cookie *network.QNetworkCookie) {
		if cookie.Name().Data() == "_mastodon_session" {
			chosenServer := chosenServer{
				Domain: cookie.Domain(),
			}
			fileContent, _ := json.MarshalIndent(&chosenServer, "", " ")
			files.CreateFile(global.ConfigFileDir, "chosenServer", fileContent)
		}
	})
	store.LoadAllCookies()
}

func (m *QClient) resendConfirmationEmail() {
	fmt.Println("Resending confirmation email for ", *userAccessToken)
	r, err := http.NewRequest(http.MethodPost, "https://"+m.SelectedServer().Domain()+"/api/v1/emails/confirmation", nil)
	r.Header.Add("Authorization", "Bearer "+*userAccessToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.StatusCode)
}

func (m *QClient) getServerRules() {
	r, err := http.NewRequest(http.MethodGet, "https://"+m.SelectedServer().Domain()+"/api/v1/instance", nil)
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	type rule struct {
		Text string `json:"text"`
	}
	type rules struct {
		Rules []rule `json:"rules"`
	}
	var ruleResponse rules
	err = json.NewDecoder(resp.Body).Decode(&ruleResponse)
	var serverRules []string
	for _, ruleResponse := range ruleResponse.Rules {
		serverRules = append(serverRules, ruleResponse.Text)
	}
	for index, server := range m.Servers() {
		if !server.IsSelected() {
			continue
		}
		server.SetServerRules(serverRules)
		server.ServerRulesChanged(serverRules)
		var pIndex = m.Index(index, 0, core.NewQModelIndex())
		m.DataChanged(pIndex, pIndex, []int{ServerRules})
	}
}

func (m *QClient) handleAuthCode(authCode string) {
	domain := m.SelectedServer().Domain()
	app, err := getApp(domain)
	if err != nil {
		return
	}

	accessToken, err := getAuthorizationCodeAccessToken(domain, app.ClientID, app.ClientSecret, authCode)
	if err != nil {
		return
	}
	createPushNotificationSubscription(app, domain, accessToken)
	m.SetWebviewUrl("https://" + m.SelectedServer().Domain())
	m.WebviewUrlChanged("https://" + m.SelectedServer().Domain())
}

func (m *QClient) setIsLoggingIn() {
	m.SetUserAction(IsLoggingIn)
}

func (m *QClient) setIsRegistering() {
	m.SetUserAction(IsRegistering)
}
func (m *QClient) getIsLoggingIn() bool {
	return m.UserAction() == IsLoggingIn
}

func (m *QClient) getIsRegistering() bool {
	return m.UserAction() == IsRegistering
}

func (m *QClient) roleNames() map[int]*core.QByteArray {
	return m.Roles()
}

func (m *QClient) data(index *core.QModelIndex, role int) *core.QVariant {
	if _, ok := listRoles[role]; !ok {
		return core.NewQVariant()
	}
	server := m.Servers()[index.Row()]
	switch role {
	case Domain:
		return core.NewQVariant1(server.Domain())
	case Description:
		return core.NewQVariant1(server.Description())
	case Selected:
		return core.NewQVariant1(server.IsSelected())
	case TotalUsers:
		return core.NewQVariant1(server.TotalUsers())
	case Language:
		return core.NewQVariant1(server.Language())
	case MatchingSearchTerm:
		return core.NewQVariant1(server.IsMatchingSearchTerm())
	case ServerRules:
		return core.NewQVariant1(server.ServerRules())
	case RequiresApproval:
		return core.NewQVariant1(server.IsRequiresApproval())
	default:
		return core.NewQVariant()
	}
}

func (m *QClient) filterServers(searchTerm string) {
	var hasMatchingServer bool
	for index, server := range m.Servers() {
		isMatching := strings.Contains(strings.ToLower(server.Domain()), strings.ToLower(searchTerm))
		server.SetMatchingSearchTerm(isMatching)
		server.MatchingSearchTermChanged(isMatching)
		var pIndex = m.Index(index, 0, core.NewQModelIndex())
		m.DataChanged(pIndex, pIndex, []int{MatchingSearchTerm})
		if isMatching {
			hasMatchingServer = true
		}
	}
	if !hasMatchingServer {
		tryAddServerInstance(m, searchTerm)
	}
}

func tryAddServerInstance(qClient *QClient, searchTerm string) {
	fmt.Println("Trying to add server", searchTerm)
	r, err := http.NewRequest(http.MethodGet, "https://"+searchTerm+"/api/v1/instance", nil)
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	type instanceStats struct {
		UserCount int `json:"user_count"`
	}
	type instance struct {
		ShortDescription string        `json:"short_description"`
		Languages        []string      `json:"languages"`
		Stats            instanceStats `json:"stats"`
	}
	var instanceResponse instance
	err = json.NewDecoder(resp.Body).Decode(&instanceResponse)

	if resp.StatusCode == http.StatusOK {
		mastodonServer := NewQServer(nil)
		mastodonServer.SetDomain(searchTerm)
		mastodonServer.SetDescription(instanceResponse.ShortDescription)
		mastodonServer.SetSelected(false)
		mastodonServer.SetTotalUsers(instanceResponse.Stats.UserCount)
		mastodonServer.SetLanguage(instanceResponse.Languages[0])
		qClient.AddServer(mastodonServer)
		fmt.Println("Added server", searchTerm)
	}
}

func (m *QClient) addServer(p *QServer) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.Servers()), len(m.Servers()))
	m.SetServers(append(m.Servers(), p))
	m.EndInsertRows()
}

func (m *QClient) setSelected(row int) {
	for index, server := range m.Servers() {
		if index == row {
			server.SetSelected(true)
			m.SetSelectedServer(server)
		} else {
			server.SetSelected(false)
		}
		var pIndex = m.Index(index, 0, core.NewQModelIndex())
		m.DataChanged(pIndex, pIndex, []int{Selected})
	}
}
func (m *QClient) setLoginUrl() {
	app, err := getCachedAppOrCreate(m.SelectedServer().Domain())

	if err != nil {
		return
	}
	if app.HasPushSubscription {
		dom := "https://" + m.SelectedServer().Domain()
		m.SetWebviewUrl(dom)
		m.WebviewUrlChanged(dom)
		return
	}
	fmt.Println("Setting login url to ", app.AuthURI)
	m.SetWebviewUrl(app.AuthURI)
	m.WebviewUrlChanged(app.AuthURI)
}

func getApp(domain string) (*PushApplication, error) {
	filename := global.ConfigFileDir + "app-" + domain + ".json"
	if !files.FileExists(filename) {
		fmt.Println("Could not find app")
		return nil, nil
	}

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("could not read config from location: " + filename)
	}
	app := &PushApplication{}
	err = json.Unmarshal(fileBytes, app)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return app, nil
}
func getCachedAppOrCreate(domain string) (*PushApplication, error) {
	cachedApp, err := getApp(domain)
	if err != nil || cachedApp != nil {
		return cachedApp, err
	}

	app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:       "https://" + domain,
		ClientName:   clientName,
		Scopes:       clientScopes,
		RedirectURIs: "mastodon://oauth",
	})
	fileContent, _ := json.MarshalIndent(&app, "", " ")
	files.CreateFile(global.ConfigFileDir, "app-"+domain+".json", fileContent)
	fmt.Println("App created")
	return &PushApplication{
		Application:         app,
		HasPushSubscription: false,
	}, nil
}

func getClientCredentialsAccessToken(domain string, clientId string, clientSecret string) (string, error) {
	return getAccessToken(domain, clientId, clientSecret, "client_credentials", "")
}
func getAuthorizationCodeAccessToken(domain string, clientId string, clientSecret string, code string) (string, error) {
	return getAccessToken(domain, clientId, clientSecret, "authorization_code", code)
}
func getAccessToken(domain string, clientId string, clientSecret string, grantType string, code string) (string, error) {
	tokenFormData := url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"grant_type":    {grantType},
		"redirect_uri":  {"mastodon://oauth"},
		"scope":         {clientScopes},
	}
	if code != "" {
		tokenFormData.Add("code", code)
	}
	r, err := http.NewRequest(http.MethodPost, "https://"+domain+"/oauth/token", strings.NewReader(tokenFormData.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	tokenResponse := TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	println("App authenticated", tokenResponse.AccessToken)
	return tokenResponse.AccessToken, nil
}

func (m *QClient) createAccount(displayName string, username string, email string, password string, reason string) string {
	if len(password) < 8 {
		return "Password must be 8 characters or more."
	}
	for _, qServer := range m.Servers() {
		if qServer.IsSelected() {
			app, err := getCachedAppOrCreate(qServer.Domain())
			if err != nil {
				fmt.Println(err)
				return "Could not create application"
			}
			client := &http.Client{}
			accessToken, err := getClientCredentialsAccessToken(qServer.Domain(), app.ClientID, app.ClientSecret)
			if err != nil {
				return "Could not get an access token"
			}
			accountFormData := url.Values{
				"username":  {username},
				"email":     {email},
				"password":  {password},
				"agreement": {"1"},
				"locale":    {"en"},
			}
			if reason != "" {
				accountFormData.Add("reason", reason)
			}
			r, err := http.NewRequest(http.MethodPost, "https://"+qServer.Domain()+"/api/v1/accounts", strings.NewReader(accountFormData.Encode()))
			r.Header.Add("Authorization", "Bearer "+accessToken)
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, _ := client.Do(r)
			if err != nil {
				fmt.Println(err)
				return "Could reach the server"
			}

			tokenResponse := CreateAccountResponse{}
			err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
			if err != nil {
				fmt.Println(err)
				return "Could not parse the response"
			}
			if tokenResponse.Error != "" {
				return tokenResponse.Error
			}
			userAccessToken = &tokenResponse.AccessToken
			awaitActivation(m, qServer, app, displayName, tokenResponse.AccessToken)
			return ""
		}
	}
	return "Something went wrong"
}

func awaitActivation(qClient *QClient, qServer *QServer, app *PushApplication, displayName string, accessToken string) {
	qClient.AwaitingActivationChanged(true)

	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				{
					credentialsResponse := verifyCredentials(qServer, accessToken)
					if credentialsResponse.StatusCode == http.StatusOK {
						close(quit)
						qClient.SetAwaitingActivation(false)
						qClient.AwaitingActivationChanged(false)
						fmt.Println("Activated account")
						loginUrl := "https://" + qServer.Domain() + "/auth/sign_in"
						fmt.Println("Setting webview url to", loginUrl)
						qClient.SetWebviewUrl(loginUrl)
						qClient.WebviewUrlChanged(loginUrl)
						go createPushNotificationSubscription(app, qServer.Domain(), accessToken)
						go updateUserCredentials(qServer.Domain(), displayName, accessToken)
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func updateUserCredentials(domain string, displayName string, accessToken string) {
	client := &http.Client{}
	updateCredentialsFormData := url.Values{
		"display_name": {displayName},
	}
	r, err := http.NewRequest(http.MethodPost, "https://"+domain+"/api/v1/accounts/update_credentials", strings.NewReader(updateCredentialsFormData.Encode()))
	r.Header.Add("Authorization", "Bearer "+accessToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		fmt.Println("Could not update credentials")
	}
	_, err = client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Updated credentials for", accessToken)
}
func createPushNotificationSubscription(app *PushApplication, domain string, accessToken string) {
	client := &http.Client{}
	pushClient := pushnotifications.GetPushClient()
	if pushClient == nil {
		return
	}
	publicKey := pushClient.ExportPublicKey()
	sharedSecret := pushClient.ExportSharedSecret()
	pushFormData := url.Values{
		"subscription[endpoint]":     {"https://mastodon-relay.joinubuntutouch.dev/" + pushClient.PushToken},
		"subscription[keys][p256dh]": {publicKey},
		"subscription[keys][auth]":   {sharedSecret},
		"data[alerts][follow]":       {"true"},
		"data[alerts][favourite]":    {"true"},
		"data[alerts][reblog]":       {"true"},
		"data[alerts][mention]":      {"true"},
	}
	r, err := http.NewRequest(http.MethodPost, "https://"+domain+"/api/v1/push/subscription", strings.NewReader(pushFormData.Encode()))
	r.Header.Add("Authorization", "Bearer "+accessToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		fmt.Println("Could not create subscription")
	}
	_, err = client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Created push subscription for", accessToken)
	app.HasPushSubscription = true
	fileContent, _ := json.MarshalIndent(&app, "", " ")
	files.CreateFile(global.ConfigFileDir, "app-"+domain+".json", fileContent)
	fmt.Println("App updated")
}

func verifyCredentials(qServer *QServer, accessToken string) *http.Response {
	client := &http.Client{}
	fmt.Println("Checking status")
	r, err := http.NewRequest(http.MethodGet, "https://"+qServer.Domain()+"/api/v1/accounts/verify_credentials", nil)
	r.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	return resp
}

func GetQClient() (*QClient, error) {
	client := &http.Client{}
	result, err := client.Get("https://api.joinmastodon.org/servers")
	if err != nil {
		fmt.Println("Couldn't get a response when getting servers")
		return nil, err
	}
	var serverResponse []ServerResponse

	err = json.NewDecoder(result.Body).Decode(&serverResponse)
	if err != nil {
		fmt.Println("Couldn't parse the response when getting servers")
		return nil, err
	}
	qClient := NewQClient(nil)
	for _, mastodonServerResponse := range serverResponse {
		mastodonServer := NewQServer(nil)
		mastodonServer.SetDomain(mastodonServerResponse.Domain)
		mastodonServer.SetDescription(mastodonServerResponse.Description)
		mastodonServer.SetSelected(false)
		mastodonServer.SetTotalUsers(mastodonServerResponse.TotalUsers)
		mastodonServer.SetLanguage(mastodonServerResponse.Language)
		mastodonServer.SetRequiresApproval(mastodonServerResponse.ApprovalRequired)
		qClient.AddServer(mastodonServer)
	}

	return qClient, nil
}
