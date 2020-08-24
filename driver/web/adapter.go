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

package web

import (
	"encoding/json"
	"fmt"
	"health/core"
	"health/core/model"
	"health/driver/web/rest"
	"health/utils"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"

	httpSwagger "github.com/swaggo/http-swagger"
)

//Adapter entity
type Adapter struct {
	host string
	auth *Auth

	apisHandler      rest.ApisHandler
	adminApisHandler rest.AdminApisHandler

	app *core.Application
}

// @title Rokwire Health Building Block API
// @description Rokwire Health Building Block API Documentation.
// @version 1.10.0
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost
// @BasePath /health
// @schemes https

// @securityDefinitions.apikey RokwireAuth
// @in header
// @name ROKWIRE-API-KEY

// @securityDefinitions.apikey AppUserAuth
// @in header (add Bearer prefix to the Authorization value)
// @name Authorization

// @securityDefinitions.apikey ProvidersAuth
// @in header
// @name ROKWIRE-HS-API-KEY

// @securityDefinitions.apikey AdminUserAuth
// @in header (add Bearer prefix to the Authorization value)
// @name Authorization

//Start starts the module
func (we Adapter) Start() {

	//add listener to the application
	we.app.AddListener(&AppListener{&we})

	we.auth.Start()

	router := mux.NewRouter().StrictSlash(true)

	// handle apis
	subrouter := router.PathPrefix("/health").Subrouter()
	subrouter.PathPrefix("/doc/ui").Handler(we.serveDocUI())
	subrouter.HandleFunc("/doc", we.serveDoc)
	subrouter.HandleFunc("/version", we.wrapFunc(we.apisHandler.Version)).Methods("GET")

	// handle covid19 rest apis /////////////
	covid19RestSubrouter := subrouter.PathPrefix("/covid19").Subrouter()

	//app id token auth
	covid19RestSubrouter.HandleFunc("/login", we.loginUser).Methods("POST")
	covid19RestSubrouter.HandleFunc("/user", we.getUser).Methods("GET")
	covid19RestSubrouter.HandleFunc("/user/clear", we.userAuthWrapFunc(we.apisHandler.ClearUserData)).Methods("GET")

	covid19RestSubrouter.HandleFunc("/ctests", we.userAuthWrapFunc(we.apisHandler.GetCTests)).Methods("GET")
	covid19RestSubrouter.HandleFunc("/ctests/{id}", we.userAuthWrapFunc(we.apisHandler.UpdateCTest)).Methods("PUT")
	covid19RestSubrouter.HandleFunc("/ctests", we.userAuthWrapFunc(we.apisHandler.DeleteCTests)).Methods("DELETE")

	covid19RestSubrouter.HandleFunc("/v2/statuses", we.userAuthWrapFunc(we.apisHandler.GetStatusV2Deprecated)).Methods("GET")
	covid19RestSubrouter.HandleFunc("/v2/statuses", we.userAuthWrapFunc(we.apisHandler.CreateOrUpdateStatusV2Deprecated)).Methods("PUT")
	covid19RestSubrouter.HandleFunc("/v2/statuses", we.userAuthWrapFunc(we.apisHandler.DeleteStatusV2Deprecated)).Methods("DELETE")
	covid19RestSubrouter.HandleFunc("/v2/app-version/{app-version}/statuses", we.userAuthWrapFunc(we.apisHandler.GetStatusV2)).Methods("GET")
	covid19RestSubrouter.HandleFunc("/v2/app-version/{app-version}/statuses", we.userAuthWrapFunc(we.apisHandler.CreateOrUpdateStatusV2)).Methods("PUT")
	covid19RestSubrouter.HandleFunc("/v2/app-version/{app-version}/statuses", we.userAuthWrapFunc(we.apisHandler.DeleteStatusV2)).Methods("DELETE")

	covid19RestSubrouter.HandleFunc("/v2/histories", we.userAuthWrapFunc(we.apisHandler.GetHistoriesV2)).Methods("GET")
	covid19RestSubrouter.HandleFunc("/v2/histories", we.userAuthWrapFunc(we.apisHandler.CreateHistoryV2)).Methods("POST")
	covid19RestSubrouter.HandleFunc("/v2/histories/{id}", we.userAuthWrapFunc(we.apisHandler.UpdateHistoryV2)).Methods("PUT")
	covid19RestSubrouter.HandleFunc("/v2/histories", we.userAuthWrapFunc(we.apisHandler.DeleteHistoriesV2)).Methods("DELETE")

	//provider auth
	covid19RestSubrouter.HandleFunc("/users/uin/{uin}", we.providerAuthWrapFunc(we.apisHandler.GetUserByShibbolethUIN)).Methods("GET")
	covid19RestSubrouter.HandleFunc("/users/re-post", we.providerAuthWrapFunc(we.apisHandler.GetUsersForRePost)).Methods("GET")
	covid19RestSubrouter.HandleFunc("/ctests", we.providerAuthWrapFunc(we.apisHandler.CreateExternalCTest)).Methods("POST")

	// api key auth
	covid19RestSubrouter.HandleFunc("/counties", we.authWrapFunc(we.apisHandler.GetCounties)).Methods("GET")
	covid19RestSubrouter.HandleFunc("/counties/{id}", we.authWrapFunc(we.apisHandler.GetCounty)).Methods("GET")

	covid19RestSubrouter.HandleFunc("/rules/county/{county-id}", we.authWrapFunc(we.apisHandler.GetRulesByCounty)).Methods("GET")
	covid19RestSubrouter.HandleFunc("/symptom-rules/county/{county-id}", we.authWrapFunc(we.apisHandler.GetSymptomRuleByCounty)).Methods("GET")
	covid19RestSubrouter.HandleFunc("/access-rules/county/{county-id}", we.authWrapFunc(we.apisHandler.GetAccessRuleByCounty)).Methods("GET")

	covid19RestSubrouter.HandleFunc("/resources", we.authWrapFunc(we.apisHandler.GetResources)).Methods("GET")

	covid19RestSubrouter.HandleFunc("/faq", we.authWrapFunc(we.apisHandler.GetFAQ)).Methods("GET")

	covid19RestSubrouter.HandleFunc("/news", we.authWrapFunc(we.apisHandler.GetNews)).Methods("GET")

	covid19RestSubrouter.HandleFunc("/providers", we.authWrapFunc(we.apisHandler.GetProvidersByCounties)).Methods("GET").Queries("county-ids", "")
	covid19RestSubrouter.HandleFunc("/providers", we.authWrapFunc(we.apisHandler.GetProviders)).Methods("GET")
	covid19RestSubrouter.HandleFunc("/providers/county/{county-id}", we.authWrapFunc(we.apisHandler.GetProvidersByCounty)).Methods("GET")

	covid19RestSubrouter.HandleFunc("/locations", we.authWrapFunc(we.apisHandler.GetLocationsByCountyIDProviderID)).Methods("GET").Queries("county-id", "", "provider-id", "")
	covid19RestSubrouter.HandleFunc("/locations", we.authWrapFunc(we.apisHandler.GetLocationsByCountyID)).Methods("GET").Queries("county-id", "")
	covid19RestSubrouter.HandleFunc("/locations/{id}", we.authWrapFunc(we.apisHandler.GetLocation)).Methods("GET")

	covid19RestSubrouter.HandleFunc("/test-types", we.authWrapFunc(we.apisHandler.GetTestTypesByIDs)).Methods("GET").Queries("ids", "")
	covid19RestSubrouter.HandleFunc("/test-types", we.authWrapFunc(we.apisHandler.GetTestTypes)).Methods("GET")

	covid19RestSubrouter.HandleFunc("/symptom-groups", we.authWrapFunc(we.apisHandler.GetSymptomGroups)).Methods("GET")

	covid19RestSubrouter.HandleFunc("/trace/report", we.authWrapFunc(we.apisHandler.AddTraceReport)).Methods("POST")
	covid19RestSubrouter.HandleFunc("/trace/exposures", we.authWrapFunc(we.apisHandler.GetExposures)).Methods("GET")

	// handle admin rest apis /////////////////
	adminRestSubrouter := router.PathPrefix("/health/admin").Subrouter()

	//admin app id token auth
	adminRestSubrouter.HandleFunc("/covid19-config", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetCovid19Config)).Methods("GET")
	adminRestSubrouter.HandleFunc("/covid19-config", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateCovid19Config)).Methods("PUT")

	adminRestSubrouter.HandleFunc("/news", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetNews)).Methods("GET")
	adminRestSubrouter.HandleFunc("/news", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateNews)).Methods("POST")
	adminRestSubrouter.HandleFunc("/news/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateNews)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/news/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteNews)).Methods("DELETE")

	adminRestSubrouter.HandleFunc("/resources", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetResources)).Methods("GET")
	adminRestSubrouter.HandleFunc("/resources", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateResources)).Methods("POST")
	adminRestSubrouter.HandleFunc("/resources/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateResource)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/resources/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteResource)).Methods("DELETE")
	adminRestSubrouter.HandleFunc("/resources/display-order", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateDisplaOrderResources)).Methods("POST")

	//TODO refactor
	adminRestSubrouter.HandleFunc("/faq", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetFAQs)).Methods("GET")
	adminRestSubrouter.HandleFunc("/faq", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateFAQItem)).Methods("POST")
	adminRestSubrouter.HandleFunc("/faq/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateFAQItem)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/faq/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteFAQItem)).Methods("DELETE")
	adminRestSubrouter.HandleFunc("/faq/section/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateFAQSection)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/faq/section/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteFAQSection)).Methods("DELETE")

	adminRestSubrouter.HandleFunc("/providers", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetProviders)).Methods("GET")
	adminRestSubrouter.HandleFunc("/providers", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateProvider)).Methods("POST")
	adminRestSubrouter.HandleFunc("/providers/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateProvider)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/providers/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteProvider)).Methods("DELETE")

	adminRestSubrouter.HandleFunc("/test-types", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetTestTypes)).Methods("GET")
	adminRestSubrouter.HandleFunc("/test-types", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateTestType)).Methods("POST")
	adminRestSubrouter.HandleFunc("/test-types/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateTestType)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/test-types/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteTestType)).Methods("DELETE")

	adminRestSubrouter.HandleFunc("/test-type-results", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateTestTypeResult)).Methods("POST")
	adminRestSubrouter.HandleFunc("/test-type-results/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateTestTypeResult)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/test-type-results/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteTestTypeResult)).Methods("DELETE")
	adminRestSubrouter.HandleFunc("/test-type-results", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetTestTypeResultsByTestTypeID)).Methods("GET").Queries("test-type-id", "")

	adminRestSubrouter.HandleFunc("/counties", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetCounties)).Methods("GET")
	adminRestSubrouter.HandleFunc("/counties", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateCounty)).Methods("POST")
	adminRestSubrouter.HandleFunc("/counties/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateCounty)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/counties/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteCounty)).Methods("DELETE")

	adminRestSubrouter.HandleFunc("/guidelines", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateGuideline)).Methods("POST")
	adminRestSubrouter.HandleFunc("/guidelines/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateGuideline)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/guidelines/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteGuideline)).Methods("DELETE")
	adminRestSubrouter.HandleFunc("/guidelines", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetGuidelinesByCountyID)).Methods("GET").Queries("county-id", "")

	adminRestSubrouter.HandleFunc("/county-statuses", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateCountyStatus)).Methods("POST")
	adminRestSubrouter.HandleFunc("/county-statuses/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateCountyStatus)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/county-statuses/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteCountyStatus)).Methods("DELETE")
	adminRestSubrouter.HandleFunc("/county-statuses", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetCountyStatusesByCountyID)).Methods("GET").Queries("county-id", "")

	adminRestSubrouter.HandleFunc("/rules", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetRules)).Methods("GET")
	adminRestSubrouter.HandleFunc("/rules", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateRule)).Methods("POST")
	adminRestSubrouter.HandleFunc("/rules/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateRule)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/rules/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteRule)).Methods("DELETE")

	adminRestSubrouter.HandleFunc("/locations", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetLocations)).Methods("GET")
	adminRestSubrouter.HandleFunc("/locations", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateLocation)).Methods("POST")
	adminRestSubrouter.HandleFunc("/locations/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateLocation)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/locations/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteLocation)).Methods("DELETE")

	adminRestSubrouter.HandleFunc("/symptoms", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateSymptom)).Methods("POST")
	adminRestSubrouter.HandleFunc("/symptoms/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateSymptom)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/symptoms/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteSymptom)).Methods("DELETE")

	adminRestSubrouter.HandleFunc("/symptom-groups", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetSymptomGroups)).Methods("GET")

	adminRestSubrouter.HandleFunc("/symptom-rules", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetSymptomRules)).Methods("GET")
	adminRestSubrouter.HandleFunc("/symptom-rules", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateSymptomRule)).Methods("POST")
	adminRestSubrouter.HandleFunc("/symptom-rules/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateSymptomRule)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/symptom-rules/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteSymptomRule)).Methods("DELETE")

	adminRestSubrouter.HandleFunc("/manual-tests", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetManualTestsByCountyID)).Methods("GET").Queries("county-id", "")
	adminRestSubrouter.HandleFunc("/manual-tests/{id}/process", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.ProcessManualTest)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/manual-tests/{id}/image", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetManualTestImage)).Methods("GET")

	adminRestSubrouter.HandleFunc("/access-rules", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetAccessRules)).Methods("GET")
	adminRestSubrouter.HandleFunc("/access-rules", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.CreateAccessRule)).Methods("POST")
	adminRestSubrouter.HandleFunc("/access-rules/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.UpdateAccessRule)).Methods("PUT")
	adminRestSubrouter.HandleFunc("/access-rules/{id}", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.DeleteAccessRule)).Methods("DELETE")

	adminRestSubrouter.HandleFunc("/user", we.adminAppIDTokenAuthWrapFunc(we.adminApisHandler.GetUserByExternalID)).Methods("GET").Queries("external-id", "")

	adminRestSubrouter.HandleFunc("/actions", we.adminAppIDTokenAuthWrapFunc(we.apisHandler.CreateAction)).Methods("POST")

	log.Fatal(http.ListenAndServe(":80", router))
}

func (we Adapter) serveDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("access-control-allow-origin", "*")
	http.ServeFile(w, r, "./docs/swagger.yaml")
}

func (we Adapter) serveDocUI() http.Handler {
	url := fmt.Sprintf("%s/health/doc", we.host)
	return httpSwagger.Handler(httpSwagger.URL(url))
}

func (we Adapter) wrapFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		handler(w, req)
	}
}

func (we Adapter) authWrapFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		authenticated := we.auth.apiKeyCheck(w, req)
		if !authenticated {
			return
		}

		handler(w, req)
	}
}

type authFunc = func(model.User, http.ResponseWriter, *http.Request)

func (we Adapter) adminAppIDTokenAuthWrapFunc(handler authFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		var err error

		ok, user, shibboAuth := we.auth.adminCheck(w, req)
		if !ok {
			return
		}
		if user == nil {
			//it is valid but the user does not exist, so create it first
			user, err = we.auth.createAdminAppUser(shibboAuth)
			if err != nil {
				log.Printf("Error on creating admin app user - %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			log.Println("Admin user created")
		}

		//handle global access control for now
		if !user.IsAdmin() {
			log.Println("Access control error")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		handler(*user, w, req)
	}
}

func (we Adapter) userAuthWrapFunc(handler authFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		ok, user, _, _ := we.auth.userCheck(w, req)
		if !ok {
			return
		}
		if user == nil {
			//it is valid but the user is not logged in - return 200/null
			log.Println("200 - Not logged in")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("null"))
			return
		}

		handler(*user, w, req)
	}
}

type loginAppUserRequest struct {
	UUID                 string  `json:"uuid" validate:"required"`
	PublicKey            string  `json:"public_key" validate:"required"`
	Consent              *bool   `json:"consent" validate:"required"`
	ExposureNotification *bool   `json:"exposure_notification" validate:"required"`
	RePost               *bool   `json:"re_post"`
	EncryptedKey         *string `json:"encrypted_key"`
	EncryptedBlob        *string `json:"encrypted_blob"`
} //@name loginUserRequest

//all manipulating of the user must happen via the auth module. We cache the users in the auth module
// @Description Creates a user, updates it if already created
// @Tags Covid19
// @ID loginUser
// @Accept json
// @Produce json
// @Param data body web.loginAppUserRequest true "body data"
// @Success 200 {object} string "Successfully created or Successfully updated"
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AppUserAuth
// @Router /covid19/login [post]
func (we Adapter) loginUser(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)

	ok, user, externalID, authType := we.auth.userCheck(w, r)
	if !ok {
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal login user - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData loginAppUserRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the login user request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate data
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating login data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuid := requestData.UUID
	publicKey := requestData.PublicKey
	consent := requestData.Consent
	exposureNotification := requestData.ExposureNotification
	rePost := requestData.RePost
	encryptedKey := requestData.EncryptedKey
	encryptedBlob := requestData.EncryptedBlob

	if user == nil {
		//we need to create

		rePostValue := false
		if authType != nil && *authType == "shibboleth" {
			rePostValue = true
		}
		err = we.auth.createAppUser(*externalID, uuid, publicKey, *consent, *exposureNotification, rePostValue, encryptedKey, encryptedBlob)
		if err != nil {
			log.Println("Error on creating user")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully created"))
	} else {
		//we need to update

		err = we.auth.updateAppUser(*user, uuid, publicKey, *consent, *exposureNotification, rePost, encryptedKey, encryptedBlob)
		if err != nil {
			log.Println("Error on updating user")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully updated"))
	}
}

// @Description Gives the current user.
// @Tags Covid19
// @ID getUser
// @Accept json
// @Success 200 {object} rest.AppUserResponse
// @Security AppUserAuth
// @Router /covid19/user [get]
func (we Adapter) getUser(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)

	ok, user, _, _ := we.auth.userCheck(w, r)
	if !ok {
		return
	}

	if user == nil {
		//it is valid but the user is not logged in - return 200/null
		log.Println("200 - Not logged in")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("null"))
		return

	}

	response := rest.AppUserResponse{UUID: user.UUID, PublicKey: user.PublicKey,
		Consent: user.Consent, ExposureNotification: user.ExposureNotification, RePost: user.RePost,
		EncryptedKey: user.EncryptedKey, EncryptedBlob: user.EncryptedBlob}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the user")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charser=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (we Adapter) providerAuthWrapFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		authenticated := we.auth.providersCheck(w, req)
		if !authenticated {
			return
		}

		handler(w, req)
	}
}

//NewWebAdapter creates new WebAdapter instance
func NewWebAdapter(host string, app *core.Application, appKeys []string, oidcProvider string,
	oidcAppClientID string, oidcAdminClientID string, phoneAuthSecret string, providersKeys []string) Adapter {
	auth := NewAuth(app, appKeys, oidcProvider, oidcAppClientID, oidcAdminClientID, phoneAuthSecret, providersKeys)

	apisHandler := rest.NewApisHandler(app)
	adminApisHandler := rest.NewAdminApisHandler(app)
	return Adapter{host: host, auth: auth, apisHandler: apisHandler, adminApisHandler: adminApisHandler, app: app}
}

//AppListener implements core.ApplicationListener interface
type AppListener struct {
	adapter *Adapter
}

//OnClearUserData notifies that a user has been removed
func (al *AppListener) OnClearUserData(user model.User) {
	log.Println("AppListener -> OnClearUserData -> " + user.ID)

	//take out the removed user from the cached users
	al.adapter.auth.userAuth.deleteCacheUser(user.ExternalID)
}

//OnUserUpdated notifies that a user has been updated
func (al *AppListener) OnUserUpdated(user model.User) {
	log.Println("AppListener -> OnUserUpdated -> " + user.ID)

	//take out the updated user from the cached users
	al.adapter.auth.userAuth.deleteCacheUser(user.ExternalID)
}
