package pushnotifications

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"mastodon-client/global"
	"os"
	"time"
)

const configFilename string = "pushConfigFile.json"

type Config struct {
	PrivateKey   ecdsa.PrivateKey
	ServerKey    ecdsa.PublicKey
	SharedSecret []byte
	PushToken    string
}

type PushClient struct {
	*Config
}

type PushMessage struct {
	Message      `json:"message"`
	Notification `json:"notification"`
}

type Notification struct {
	Card    `json:"card"`
	Vibrate bool `json:"vibrate"`
	Sound   bool `json:"sound"`
}
type Card struct {
	Icon    string   `json:"icon"`
	Summary string   `json:"summary"`
	Body    string   `json:"body"`
	Popup   bool     `json:"popup"`
	Persist bool     `json:"persist"`
	Actions []string `json:"actions"`
}
type Message struct {
	PublicKey        string           `json:"PublicKey"`
	Salt             string           `json:"Salt"`
	Expiration       time.Time        `json:"Expiration"`
	Priority         string           `json:"Priority"`
	Topic            string           `json:"Topic"`
	Urgency          string           `json:"Urgency"`
	Payload          string           `json:"Payload"`
	DecryptedPayload DecryptedPayload `json:"DecryptedPayload"`
}

type DecryptedPayload struct {
	AccessToken      string `json:"access_token"`
	PreferredLocale  string `json:"preferred_locale"`
	NotificationID   int    `json:"notification_id"`
	NotificationType string `json:"notification_type"`
	Icon             string `json:"icon"`
	Title            string `json:"title"`
	Body             string `json:"body"`
}

type pushConfigFile struct {
	PublicKeyString    string
	PrivateKeyString   string
	SharedSecretString string
	PushToken          string
}

func GetPushClient() *PushClient {
	config, err := getConfig()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	client := &PushClient{Config: config}

	return client
}

func getConfig() (*Config, error) {
	fileBytes, err := os.ReadFile(global.ConfigFileDir + configFilename)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("could not read config from location: " + global.ConfigFileDir + configFilename)
	}
	pushConfigFile := pushConfigFile{}
	err = json.Unmarshal(fileBytes, &pushConfigFile)
	if err != nil {
		return nil, err
	}
	config := &Config{
		PushToken: pushConfigFile.PushToken,
	}
	err = config.ImportSharedSecret(pushConfigFile.SharedSecretString)
	err = config.ImportPrivateKey(pushConfigFile.PrivateKeyString)
	if err != nil {
		return nil, err
	}
	return config, nil
}
