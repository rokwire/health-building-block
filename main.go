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

package main

import (
	"health/core"
	audit "health/driven/audit"
	dataprovider "health/driven/dataprovider"
	messaging "health/driven/messaging"
	profilebb "health/driven/profilebb"
	rokmetro "health/driven/rokmetro"
	sender "health/driven/sender"
	storage "health/driven/storage"
	driver "health/driver/web"
	"log"
	"os"
	"strings"
)

var (
	// Version : version of this executable
	Version string
	// Build : build date of this executable
	Build string
)

func main() {
	if len(Version) == 0 {
		Version = "dev"
	}

	//mongoDB adapter
	mongoDBAuth := getEnvKey("HEALTH_MONGO_AUTH", true)
	mongoDBName := getEnvKey("HEALTH_MONGO_DATABASE", true)
	mongoTimeout := getEnvKey("HEALTH_MONGO_TIMEOUT", false)
	storageAdapter := storage.NewStorageAdapter(mongoDBAuth, mongoDBName, mongoTimeout)
	err := storageAdapter.Start()
	if err != nil {
		log.Fatal("Cannot start the mongoDB adapter - " + err.Error())
	}

	//audit adapter
	auditAdapter := audit.NewAuditAdapter(mongoDBAuth, mongoDBName, mongoTimeout)
	err = auditAdapter.Start()
	if err != nil {
		log.Fatal("Cannot start the audit adapter - " + err.Error())
	}

	//data provider adapter
	newsRSSURL := getEnvKey("HEALTH_NEWS_RSS_URL", true)
	resourcesURL := getEnvKey("HEALTH_RESOURCES_URL", true)
	dataProvider := dataprovider.NewDataProviderAdapter(newsRSSURL, resourcesURL)

	//sender adapter
	smtpHost := getEnvKey("HEALTH_SMTP_HOST", true)
	smtpPort := getEnvKey("HEALTH_SMTP_PORT", true)
	user := getEnvKey("HEALTH_SMTP_USER", true)
	password := getEnvKey("HEALTH_SMTP_PASSWORD", true)
	from := getEnvKey("HEALTH_EMAIL_FROM", true)
	to := getEmailsRecepients()
	sender := sender.NewSenderAdapter(smtpHost, smtpPort, user, password, from, to)

	firebaseAuth := getEnvKey("HEALTH_FIREBASE_AUTH", true)
	firebaseProjectID := getEnvKey("HEALTH_FIREBASE_PROJECT_ID", true)
	messaging := messaging.NewFirebaseAdapter(firebaseAuth, firebaseProjectID)

	//profile bb adapter
	profileHost := getEnvKey("HEALTH_PROFILE_HOST", true)
	profileAPIKey := getEnvKey("HEALTH_PROFILE_API_KEY", true)
	profileBBAdapter := profilebb.NewProfileBBAdapter(profileHost, profileAPIKey)

	//rokmetro adapter
	rokmetroGroupsHost := getEnvKey("HEALTH_ROKMETRO_GROUPS_HOST", true)
	rokmetroGroupsAPIKey := getEnvKey("HEALTH_ROKMETRO_GROUPS_API_KEY", true)
	rokmetroAdapter := rokmetro.NewRokmetroAdapter(rokmetroGroupsHost, rokmetroGroupsAPIKey)

	//application
	application := core.NewApplication(Version, Build, dataProvider, sender, messaging, profileBBAdapter, rokmetroAdapter, storageAdapter, auditAdapter)
	application.Start()

	//web adapter
	apiKeys := getAPIKeys()
	//TODO - get ROKWIRE-EXT-HS-API-KEYS from the environment
	externalApiKey := getExternalApiKey()
	host := getEnvKey("HEALTH_HOST", true)
	oidcProvider := getEnvKey("HEALTH_OIDC_PROVIDER", true)
	oidcAppClientID := getEnvKey("HEALTH_OIDC_APP_CLIENT_ID", true)
	adminAppClientID := getEnvKey("HEALTH_OIDC_ADMIN_CLIENT_ID", true)
	adminWebAppClientID := getEnvKey("HEALTH_OIDC_ADMIN_WEB_CLIENT_ID", true)
	phoneSecret := getEnvKey("HEALTH_PHONE_SECRET", true)
	providersKeys := getHSAPIKeys()
	authKeys := getEnvKey("HEALTH_AUTH_KEYS", true)
	authIssuer := getEnvKey("HEALTH_AUTH_ISSUER", true)

	webAdapter := driver.NewWebAdapter(host, application, apiKeys, oidcProvider, oidcAppClientID, adminAppClientID, adminWebAppClientID,
		phoneSecret, authKeys, authIssuer, providersKeys, externalApiKey)

	webAdapter.Start()
}

func getEmailsRecepients() []string {
	//get from the environment
	emails, exist := os.LookupEnv("HEALTH_EMAIL_TO")
	if !exist {
		return nil
	}
	if len(emails) <= 0 {
		return nil
	}
	printEnvVar("HEALTH_EMAIL_TO", emails)

	//it is comma separated format
	emailsList := strings.Split(emails, ",")
	if len(emailsList) <= 0 {
		log.Fatal("For some reasons the emails list is empty")
	}

	return emailsList
}

func getAPIKeys() []string {
	//get from the environment
	rokwireAPIKeys := getEnvKey("ROKWIRE_API_KEYS", true)

	//it is comma separated format
	rokwireAPIKeysList := strings.Split(rokwireAPIKeys, ",")
	if len(rokwireAPIKeysList) <= 0 {
		log.Fatal("For some reasons the apis keys list is empty")
	}

	return rokwireAPIKeysList
}
func getExternalApiKey() []string {
	//get from the environment
	rokwireExternalApiKey := getExternalApiKeys("ROKWIRE-EXT-HS-API-KEYS", true)

	//it is comma separated format
	externalKeyList := strings.Split(rokwireExternalApiKey, ",")
	if len(externalKeyList) <= 0 {
		log.Fatal("The external key list is empty")
	}

	return externalKeyList
}

func getHSAPIKeys() []string {
	//get from the environment
	providersAPIKeys := getEnvKey("HEALTH_PROVIDERS_KEY", true)

	//it is comma separated format
	providersAPIKeysList := strings.Split(providersAPIKeys, ",")
	if len(providersAPIKeysList) <= 0 {
		log.Fatal("For some reasons the providers apis keys list is empty")
	}

	return providersAPIKeysList
}

func getEnvKey(key string, required bool) string {
	//get from the environment
	value, exist := os.LookupEnv(key)
	if !exist {
		if required {
			log.Fatal("No provided environment variable for " + key)
		} else {
			log.Printf("No provided environment variable for " + key)
		}
	}
	printEnvVar(key, value)
	return value
}
func getExternalApiKeys(key string, required bool) string {
	//get from the environment
	value, exist := os.LookupEnv(key)
	if !exist {
		if required {
			log.Fatal("No provided environment variable for " + key)
		} else {
			log.Printf("No provided environment variable for " + key)
		}
	}
	printEnvVar(key, value)
	return value
}

func printEnvVar(name string, value string) {
	if Version == "dev" {
		log.Printf("%s=%s", name, value)
	}
}
