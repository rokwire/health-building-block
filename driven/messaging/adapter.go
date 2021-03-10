/*
 *   Copyright (c) 2020 Board of Trustees of the University of Illinois.
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package messaging

import (
	"context"
	"log"

	fire "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	firemessaging "firebase.google.com/go/messaging"
)

//FirebaseAdapter implements Firebase messaging
type FirebaseAdapter struct {
	app    *fire.App
	client *firemessaging.Client
}

//SendNotificationMessage send a notification message
func (fa *FirebaseAdapter) SendNotificationMessage(tokens []string, title string, body string, data map[string]string) {
	if len(tokens) <= 0 {
		log.Println("SendNotificationMessage -> cannot send messages without tokens")
	}

	for _, token := range tokens {
		go func(token string, title string, body string, data map[string]string) {
			ntf := firemessaging.Notification{Title: title, Body: body}
			androidConfig := fa.getAndroidConfig()
			apnsConfig := fa.getAPNSConfig()
			msg := messaging.Message{Token: token, Notification: &ntf, Android: androidConfig, APNS: apnsConfig, Data: data}

			response, err := fa.client.Send(context.Background(), &msg)
			if err != nil {
				log.Printf("Error sending notification message - %s\n", err)
				return
			}
			log.Printf("Successfully sent notification message - %s\n", response)
		}(token, title, body, data)
	}
}

func (fa *FirebaseAdapter) getAndroidConfig() *firemessaging.AndroidConfig {
	notification := firemessaging.AndroidNotification{DefaultSound: true}
	androidConfig := firemessaging.AndroidConfig{Priority: "high", Notification: &notification}
	return &androidConfig
}

func (fa *FirebaseAdapter) getAPNSConfig() *firemessaging.APNSConfig {
	aps := firemessaging.Aps{Sound: "default", ContentAvailable: true}
	payload := firemessaging.APNSPayload{Aps: &aps}
	apnsConfig := firemessaging.APNSConfig{Payload: &payload}
	return &apnsConfig
}

//NewFirebaseAdapter creates a new firebase adapter instance
func NewFirebaseAdapter(authFile string, projectID string) *FirebaseAdapter {
	/*	conf, err := google.JWTConfigFromJSON([]byte(authFile),
			"https://www.googleapis.com/auth/firebase",
			"https://www.googleapis.com/auth/cloud-platform")
		if err != nil {
			log.Fatal(err.Error())
			return nil
		}

		tokenSource := conf.TokenSource(context.Background())
		creds := google.Credentials{ProjectID: projectID, TokenSource: tokenSource}
		opt := option.WithCredentials(&creds)
		app, err := fire.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Fatalf("error creating Firebase app: %v\n", err)
			return nil
		}
		client, err := app.Messaging(context.Background())
		if err != nil {
			log.Fatalf("error getting Messaging client: %v\n", err)
			return nil
		}

		return &FirebaseAdapter{app: app, client: client}*/
	return nil
}
