package main

import (
	"encoding/json"
	"fmt"
	"github.com/nofeaturesonlybugs/z85"
	"mastodon-client/pushnotifications"
	"os"
)

func main() {
	fmt.Println("Handling push notification")
	pushClient := pushnotifications.GetPushClient()

	args := os.Args[1:3]

	firstFileBytes, openErr := os.ReadFile(args[0])
	if openErr != nil {
		fmt.Println("Could not open first file", openErr)
		return
	}

	pushMessage := &pushnotifications.PushMessage{}
	unmarshalErr := json.Unmarshal(firstFileBytes, pushMessage)
	fmt.Println("UNMARSHAL:", string(firstFileBytes))
	if unmarshalErr != nil {
		fmt.Println("Could not unmarshal firstFile:", unmarshalErr)
		return
	}
	if pushMessage.Message.PublicKey == "" {
		writeErr := os.WriteFile(args[1], firstFileBytes, os.ModeDevice)
		if writeErr != nil {
			fmt.Println("Could not write file", writeErr)
		}
		return
	}

	decodedPubKeyBytes, _ := z85.PaddedDecode(pushMessage.Message.PublicKey)
	decodedCypherText, _ := z85.PaddedDecode(pushMessage.Message.Payload)
	decodedSalt, _ := z85.PaddedDecode(pushMessage.Message.Salt)

	decryptedPayloadBytes, decryptErr := pushClient.Decrypt(decodedPubKeyBytes, decodedSalt, decodedCypherText)
	if decryptErr != nil {
		fmt.Println("Could not decrypt payload", openErr)
	}
	fmt.Println("Decrypted payload:", string(decryptedPayloadBytes))
	decryptedPayload := &pushnotifications.DecryptedPayload{}
	unmarshalErr = json.Unmarshal(decryptedPayloadBytes, decryptedPayload)
	if unmarshalErr != nil {
		return
	}
	pushMessage.Card.Body = decryptedPayload.Body
	pushMessage.Card.Summary = decryptedPayload.Title
	pushMessage.DecryptedPayload = *decryptedPayload
	pushMessage.Icon = decryptedPayload.Icon

	updatedFirstFileBytes, updateMessageErr := json.Marshal(pushMessage)
	if updateMessageErr != nil {
		fmt.Println("Could not update first message bytes", openErr)
		return
	}
	writeErr := os.WriteFile(args[1], updatedFirstFileBytes, os.ModeDevice)
	if writeErr != nil {
		fmt.Println("Could not write file", writeErr)
	}
}
