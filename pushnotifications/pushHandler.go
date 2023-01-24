package pushnotifications

import (
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/therecipe/qt/core"
	"mastodon-client/files"
	"mastodon-client/global"
)

type PushHandler struct {
	core.QObject

	_ func() `constructor:"init"`

	_ string `property:"pushNotificationToken"`
	_ bool   `property:"pushSubscriptionRegistered"`

	_ func(string) `slot:"initialize"`
	_ func(string) `slot:"handle"`
	_ func(string) `slot:"handleError"`
}

func (handler *PushHandler) init() {
	handler.ConnectInitialize(handler.initialize)
	handler.ConnectHandle(handler.handle)
	handler.ConnectHandleError(handler.handleError)
}

func (handler *PushHandler) initialize(token string) {
	if files.FileExists(global.ConfigFileDir + configFilename) {
		fmt.Println("Found pushClient config file, do not need to generate keys")
		return
	}
	client := &PushClient{Config: &Config{}}
	err := client.GenerateNewKeys()
	if err != nil {
		fmt.Println(err)
		return
	}
	privateKey := client.ExportPrivateKey()
	sharedSecret := client.ExportSharedSecret()
	publicKeyBytes := elliptic.Marshal(client.PrivateKey.PublicKey.Curve, client.PrivateKey.PublicKey.X, client.PrivateKey.PublicKey.Y)
	publicKey := base64.RawURLEncoding.EncodeToString(publicKeyBytes)
	pushConfig := pushConfigFile{
		PublicKeyString:    publicKey,
		PrivateKeyString:   privateKey,
		SharedSecretString: sharedSecret,
		PushToken:          token,
	}
	fileContent, _ := json.MarshalIndent(&pushConfig, "", " ")
	files.CreateFile(global.ConfigFileDir, configFilename, fileContent)
}
func (handler *PushHandler) handle(message string) {
	println(message)
}
func (handler *PushHandler) handleError(message string) {
	println(message)
}
