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
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"health/core"
	"health/core/model"
	"health/utils"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
	"golang.org/x/sync/syncmap"
	"gopkg.in/ericchiang/go-oidc.v2"
)

type cacheUser struct {
	user      *model.User
	lastUsage time.Time
}

//Auth handler
type Auth struct {
	apiKeysAuth   *APIKeysAuth
	userAuth      *UserAuth
	adminAuth     *AdminAuth
	providersAuth *ProvidersAuth
	externalAuth  *ExternalAuth
}

//Start starts the auth module
func (auth *Auth) Start() error {
	auth.adminAuth.start()
	auth.userAuth.start()

	return nil
}

func (auth *Auth) clientIDCheck(w http.ResponseWriter, r *http.Request) (bool, *string) {
	clientID := r.Header.Get("APP")
	if len(clientID) == 0 {
		clientID = "edu.illinois.rokwire"
	}

	log.Println(fmt.Sprintf("400 - Bad Request"))
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad Request"))
	return false, nil
}

func (auth *Auth) apiKeyCheck(w http.ResponseWriter, r *http.Request) (bool, *string) {
	return auth.apiKeysAuth.check(w, r)
}

func (auth *Auth) adminCheck(w http.ResponseWriter, r *http.Request) (bool, *model.User, string, *model.ShibbolethAuth) {
	return auth.adminAuth.check(w, r)
}

func (auth *Auth) createAdminAppUser(shibboAuth *model.ShibbolethAuth) (*model.User, error) {
	return auth.adminAuth.createAdminAppUser(shibboAuth)
}

func (auth *Auth) externalAuthCheck(w http.ResponseWriter, r *http.Request) bool {
	return auth.externalAuth.check(w, r)
}

func (auth *Auth) providersCheck(w http.ResponseWriter, r *http.Request) bool {
	return auth.providersAuth.check(w, r)
}

func (auth *Auth) userCheck(w http.ResponseWriter, r *http.Request) (bool, *model.User, *string, *string, *string) {
	return auth.userAuth.userCheck(w, r)
}

func (auth *Auth) userAccountsCheck(w http.ResponseWriter, r *http.Request) (bool, *model.User, *model.Account) {
	return auth.userAuth.userAccountsCheck(w, r)
}

func (auth *Auth) updateAppUser(user model.User, uuid string, publicKey string, consent bool, consentVaccine bool, exposureNotification bool, rePost *bool, encryptedKey *string, encryptedBlob *string, encryptedPK *string) error {
	return auth.userAuth.updateAppUser(user, uuid, publicKey, consent, consentVaccine, exposureNotification, rePost, encryptedKey, encryptedBlob, encryptedPK)
}

func (auth *Auth) createAppUser(externalID string, uuid string, publicKey string, consent bool, consentVaccine bool, exposureNotification bool, rePost bool, encryptedKey *string, encryptedBlob *string, encryptedPK *string) error {
	return auth.userAuth.createAppUser(externalID, uuid, publicKey, consent, consentVaccine, exposureNotification, rePost, encryptedKey, encryptedBlob, encryptedPK)
}

//NewAuth creates new auth handler
func NewAuth(app *core.Application, appKeys []string, oidcProvider string,
	oidcAppClientID string, appClientID string, webAppClientID string, phoneAuthSecret string,
	authKeys string, authIssuer string, providersAPIKeys []string, externalAPIKeys []string) *Auth {
	apiKeysAuth := newAPIKeysAuth(appKeys)
	userAuth2 := newUserAuth(app, oidcProvider, oidcAppClientID, phoneAuthSecret, authKeys, authIssuer)
	adminAuth := newAdminAuth(app, oidcProvider, appClientID, webAppClientID)
	providersAuth := newProviderAuth(providersAPIKeys)
	externalAuth := newExternalAuth(externalAPIKeys)

	auth := Auth{apiKeysAuth: apiKeysAuth, userAuth: userAuth2, adminAuth: adminAuth, providersAuth: providersAuth, externalAuth: externalAuth}
	return &auth
}

/////////////////////////////////////

//APIKeysAuth entity
type APIKeysAuth struct {
	appKeys []string
}

func (auth *APIKeysAuth) check(w http.ResponseWriter, r *http.Request) (bool, *string) {
	vHeader := r.Header.Get("v")
	var appVersion *string
	if len(vHeader) > 0 {
		appVersion = &vHeader
	}

	apiKey := r.Header.Get("ROKWIRE-API-KEY")
	//check if there is api key in the header
	if len(apiKey) == 0 {
		//no key, so return 400
		log.Println(fmt.Sprintf("400 - Bad Request"))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return false, nil
	}

	//check if the api key is one of the listed
	appKeys := auth.appKeys
	exist := false
	for _, element := range appKeys {
		if element == apiKey {
			exist = true
			break
		}
	}
	if !exist {
		//not exist, so return 401
		log.Println(fmt.Sprintf("401 - Unauthorized for key %s", apiKey))

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return false, nil
	}
	return true, appVersion
}

func newAPIKeysAuth(appKeys []string) *APIKeysAuth {
	auth := APIKeysAuth{appKeys}
	return &auth
}

////////////////////////////////////

//ExternalAuth entity
type ExternalAuth struct {
	appKeys []string
}

func (auth *ExternalAuth) check(w http.ResponseWriter, r *http.Request) bool {
	apiKey := r.Header.Get("ROKWIRE-EXT-HS-API-KEY")
	//check if there is api key in the header
	if len(apiKey) == 0 {
		//no key, so return 400
		log.Println(fmt.Sprintf("400 - Bad Request"))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return false
	}

	//check if the api key is one of the listed
	appKeys := auth.appKeys
	exist := false
	for _, element := range appKeys {
		if element == apiKey {
			exist = true
			break
		}
	}
	if !exist {
		//not exist, so return 401
		log.Println(fmt.Sprintf("401 - Unauthorized for key %s", apiKey))

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return false
	}
	return true
}

//newExternalAuth creates new internal auth
func newExternalAuth(apiKeys []string) *ExternalAuth {
	auth := ExternalAuth{appKeys: apiKeys}
	return &auth
}

////////////////////////////////////

type userData struct {
	UIuceduUIN        *string   `json:"uiucedu_uin"`
	Sub               *string   `json:"sub"`
	Email             *string   `json:"email"`
	UIuceduIsMemberOf *[]string `json:"uiucedu_is_member_of"`
}

//AdminAuth entity
type AdminAuth struct {
	app *core.Application

	appVerifier    *oidc.IDTokenVerifier
	appClientID    string
	webAppVerifier *oidc.IDTokenVerifier
	webAppClientID string

	cachedUsers     *syncmap.Map //cache users while active - 5 minutes timeout
	cachedUsersLock *sync.RWMutex
}

func (auth *AdminAuth) start() {
	go auth.cleanCacheUser()
}

//cleanChacheUser cleans all users from the cache with no activity > 5 minutes
func (auth *AdminAuth) cleanCacheUser() {
	log.Println("AdminAuth -> cleanCacheUser -> start")

	toRemove := []string{}

	//find all users to remove - more than 5 minutes period from their last usage
	now := time.Now().Unix()
	auth.cachedUsers.Range(func(key, value interface{}) bool {
		cacheUser, ok := value.(*cacheUser)
		if !ok {
			return false //break the iteration
		}
		externalID, ok := key.(string)
		if !ok {
			return false //break the iteration
		}

		difference := now - cacheUser.lastUsage.Unix()
		//5 minutes
		if difference > 300 {
			toRemove = append(toRemove, externalID)
		}

		// this will continue iterating
		return true
	})

	//remove the selected ones
	count := len(toRemove)
	if count > 0 {
		log.Printf("AdminAuth -> cleanCacheUser -> %d items to remove\n", count)

		for _, key := range toRemove {
			auth.deleteCacheUser(key)
		}
	} else {
		log.Println("AdminAuth -> cleanCacheUser -> nothing to remove")
	}

	nextLoad := time.Minute * 5
	log.Printf("AdminAuth -> cleanCacheUser() -> next exec after %s\n", nextLoad)
	timer := time.NewTimer(nextLoad)
	<-timer.C
	log.Println("AdminAuth -> cleanCacheUser() -> timer expired")

	auth.cleanCacheUser()
}

func (auth *AdminAuth) check(w http.ResponseWriter, r *http.Request) (bool, *model.User, string, *model.ShibbolethAuth) {
	//1. Get the token from the request
	rawIDToken, tokenType, err := auth.getIDToken(r)
	if err != nil {
		auth.responseBadRequest(w)
		return false, nil, "", nil
	}

	//2. Get the group from the request
	group := r.Header.Get("GROUP")
	if len(group) <= 0 {
		auth.responseBadRequest(w)
		return false, nil, "", nil
	}

	//3. Validate the token
	idToken, err := auth.verify(*rawIDToken, *tokenType)
	if err != nil {
		log.Printf("error validating token - %s\n", err)

		auth.responseUnauthorized(*rawIDToken, w)
		return false, nil, "", nil
	}

	//4. Get the user data from the token
	var userData userData
	if err := idToken.Claims(&userData); err != nil {
		log.Printf("error getting user data from token - %s\n", err)

		auth.responseUnauthorized(*rawIDToken, w)
		return false, nil, "", nil
	}
	//we must have UIuceduUIN
	if userData.UIuceduUIN == nil {
		log.Printf("error - missing uiuceuin data in the token - %s\n", err)

		auth.responseUnauthorized(*rawIDToken, w)
		return false, nil, "", nil
	}

	//5. Get the user for the provided external id.
	user, err := auth.getUser(*userData.UIuceduUIN)
	if err != nil {
		log.Printf("error getting an user for external id - %s\n", err)

		auth.responseInternalServerError(w)
		return false, nil, "", nil
	}

	shibboAuth := &model.ShibbolethAuth{Uin: *userData.UIuceduUIN, Email: *userData.Email,
		IsMemberOf: userData.UIuceduIsMemberOf}

	//6.
	if user == nil {
		//we do not have a such user yet but the ID token is valid so return ok
		return true, nil, "", shibboAuth
	}
	//we have a such user, check if need to update the shibbo data before to return it
	user, err = auth.updateShiboDataIfNeeded(*user, userData)
	if err != nil {
		log.Printf("error updating an user for external id - %s\n", err)

		auth.responseInternalServerError(w)
		return false, nil, "", nil
	}

	//7. Check if the user is member of the group
	if !user.IsMemberOf(group) {
		auth.responseForbbiden(fmt.Sprintf("Security - %s is trying to access not allowed resource", *userData.Email), w)
		return false, nil, "", nil
	}

	return true, user, group, shibboAuth
}

//gets the token from the request - as cookie or as Authorization header.
//returns the id token and its type - mobile or web. If the token is taken by the cookie it is web otherwise it is mobile
func (auth *AdminAuth) getIDToken(r *http.Request) (*string, *string, error) {
	var tokenType string

	//1. Check if there is a cookie
	cookie, err := r.Cookie("rwa-at-data")
	if err == nil && cookie != nil && len(cookie.Value) > 0 {
		//there is a cookie
		tokenType = "web"
		return &cookie.Value, &tokenType, nil
	}

	//2. Check if there is a token in the Authorization header
	authorizationHeader := r.Header.Get("Authorization")
	if len(authorizationHeader) <= 0 {
		return nil, nil, errors.New("error getting Authorization header")
	}
	splitAuthorization := strings.Fields(authorizationHeader)
	if len(splitAuthorization) != 2 {
		return nil, nil, errors.New("error processing the Authorization header")
	}
	// expected - Bearer 1234
	if splitAuthorization[0] != "Bearer" {
		return nil, nil, errors.New("error processing the Authorization header")
	}
	rawIDToken := splitAuthorization[1]
	tokenType = "mobile"
	return &rawIDToken, &tokenType, nil
}

func (auth *AdminAuth) verify(rawIDToken string, tokenType string) (*oidc.IDToken, error) {
	switch tokenType {
	case "mobile":
		log.Println("AdminAuth -> mobile app client token")
		return auth.appVerifier.Verify(context.Background(), rawIDToken)
	case "web":
		log.Println("AdminAuth -> web app client token")
		return auth.webAppVerifier.Verify(context.Background(), rawIDToken)
	default:
		return nil, errors.New("AdminAuth -> there is an issue with the audience")
	}
}

func (auth *AdminAuth) createAdminAppUser(shibboAuth *model.ShibbolethAuth) (*model.User, error) {
	user, err := auth.app.CreateAdminAppUser(shibboAuth)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (auth *AdminAuth) updateShiboDataIfNeeded(current model.User, userData userData) (*model.User, error) {
	currentList := current.ShibbolethAuth.IsMemberOf
	newList := userData.UIuceduIsMemberOf

	isEqual := utils.EqualPointers(currentList, newList)
	if !isEqual {
		log.Println("updateUserIfNeeded -> need to update user")

		//1. remove it from the cache
		auth.deleteCacheUser(current.ShibbolethAuth.Uin)

		//2. update it
		current.ShibbolethAuth.IsMemberOf = userData.UIuceduIsMemberOf
		err := auth.app.UpdateUser(&current)
		if err != nil {
			return nil, err
		}
	}

	return &current, nil
}

func (auth *AdminAuth) getCachedUser(externalID string) *cacheUser {
	auth.cachedUsersLock.RLock()
	defer auth.cachedUsersLock.RUnlock()

	var cachedUser *cacheUser //to return

	item, _ := auth.cachedUsers.Load(externalID)
	if item != nil {
		cachedUser = item.(*cacheUser)
	}

	//keep the last get time
	if cachedUser != nil {
		cachedUser.lastUsage = time.Now()
		auth.cachedUsers.Store(externalID, cachedUser)
	}

	return cachedUser
}

func (auth *AdminAuth) cacheUser(externalID string, user *model.User) {
	auth.cachedUsersLock.RLock()

	cacheUser := &cacheUser{user: user, lastUsage: time.Now()}
	auth.cachedUsers.Store(externalID, cacheUser)

	auth.cachedUsersLock.RUnlock()
}

func (auth *AdminAuth) deleteCacheUser(externalID string) {
	auth.cachedUsersLock.RLock()

	auth.cachedUsers.Delete(externalID)

	auth.cachedUsersLock.RUnlock()
}

func (auth *AdminAuth) getUser(uiUceduUIN string) (*model.User, error) {
	var err error

	//1. First check if cached
	cachedUser := auth.getCachedUser(uiUceduUIN)
	if cachedUser != nil {
		return cachedUser.user, nil
	}

	//2. Check if we have a such user in the application
	user, err := auth.app.FindUserByShibbolethID(uiUceduUIN)
	if err != nil {
		log.Printf("error finding an for external id - %s\n", err)
		return nil, err
	}
	if user != nil {
		//cache it
		auth.cacheUser(uiUceduUIN, user)
		return user, nil
	}
	//there is no a such user
	return nil, nil
}

func (auth *AdminAuth) responseBadRequest(w http.ResponseWriter) {
	log.Println("AdminAuth -> 400 - Bad Request")

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad Request"))
}

func (auth *AdminAuth) responseUnauthorized(token string, w http.ResponseWriter) {
	log.Printf("AdminAuth -> 401 - Unauthorized for token %s", token)

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Unauthorized"))
}

func (auth *AdminAuth) responseForbbiden(info string, w http.ResponseWriter) {
	log.Printf("AdminAuth -> 403 - Forbidden - %s", info)

	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Forbidden"))
}

func (auth *AdminAuth) responseInternalServerError(w http.ResponseWriter) {
	log.Println("AdminAuth -> 500 - Internal Server Error")

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}

func newAdminAuth(app *core.Application, oidcProvider string, appClientID string, webAppClientID string) *AdminAuth {
	provider, err := oidc.NewProvider(context.Background(), oidcProvider)
	if err != nil {
		log.Fatalln(err)
	}

	appVerifier := provider.Verifier(&oidc.Config{ClientID: appClientID})
	webAppVerifier := provider.Verifier(&oidc.Config{ClientID: webAppClientID})

	cacheUsers := &syncmap.Map{}
	lock := &sync.RWMutex{}

	auth := AdminAuth{app: app, appVerifier: appVerifier, appClientID: appClientID,
		webAppVerifier: webAppVerifier, webAppClientID: webAppClientID,
		cachedUsers: cacheUsers, cachedUsersLock: lock}
	return &auth
}

/////////////////////////////////////

//ProvidersAuth entity
type ProvidersAuth struct {
	appKeys []string
}

func (auth *ProvidersAuth) check(w http.ResponseWriter, r *http.Request) bool {
	apiKey := r.Header.Get("ROKWIRE-HS-API-KEY")
	//check if there is api key in the header
	if len(apiKey) == 0 {
		//no key, so return 400
		log.Println(fmt.Sprintf("400 - Bad Request"))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return false
	}

	//check if the api key is one of the listed
	appKeys := auth.appKeys
	exist := false
	for _, element := range appKeys {
		if element == apiKey {
			exist = true
			break
		}
	}
	if !exist {
		//not exist, so return 401
		log.Println(fmt.Sprintf("401 - Unauthorized for key %s", apiKey))

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return false
	}
	return true
}

func newProviderAuth(appKeys []string) *ProvidersAuth {
	auth := ProvidersAuth{appKeys}
	return &auth
}

type shData struct {
	UIuceduUIN *string `json:"uiucedu_uin"`
}

type tokenData struct {
	UID      string
	Name     string
	Email    string
	Phone    string
	ClientID string
	Groups   string
	Auth     string
	Type     string
	ISS      string
}

//UserAuth entity
type UserAuth struct {
	app *core.Application

	//shibboleth - keep for back compatability
	appIDTokenVerifier *oidc.IDTokenVerifier

	//phone - keep for back compatability
	phoneAuthSecret string

	//auth service
	Keys   *jwk.Set
	Issuer string

	cachedUsers     *syncmap.Map //cache users while active - 5 minutes timeout
	cachedUsersLock *sync.RWMutex

	rosters     []map[string]string //cache rosters
	rostersLock *sync.RWMutex
}

func (auth *UserAuth) start() {
	auth.loadRosters()
	go auth.cleanCacheUser()
}

//loadRosters loads all rosters
func (auth *UserAuth) loadRosters() {
	log.Println("UserAuth -> load Rosters")

	rosters, err := auth.app.LoadAllRosters()
	if err != nil {
		log.Println("cannot load rosters")
	}

	count := len(rosters)
	newValues := make([]map[string]string, count)
	if count > 0 {
		for index, item := range rosters {
			newValues[index] = item
		}
	}
	auth.setRosters(newValues)
}

//cleanChacheUser cleans all users from the cache with no activity > 5 minutes
func (auth *UserAuth) cleanCacheUser() {
	log.Println("UserAuth -> cleanCacheUser -> start")

	toRemove := []string{}

	//find all users to remove - more than 5 minutes period from their last usage
	now := time.Now().Unix()
	auth.cachedUsers.Range(func(key, value interface{}) bool {
		cacheUser, ok := value.(*cacheUser)
		if !ok {
			return false //break the iteration
		}
		externalID, ok := key.(string)
		if !ok {
			return false //break the iteration
		}

		difference := now - cacheUser.lastUsage.Unix()
		//5 minutes
		if difference > 300 {
			toRemove = append(toRemove, externalID)
		}

		// this will continue iterating
		return true
	})

	//remove the selected ones
	count := len(toRemove)
	if count > 0 {
		log.Printf("UserAuth -> cleanCacheUser -> %d items to remove\n", count)

		for _, key := range toRemove {
			auth.deleteCacheUser(key)
		}
	} else {
		log.Println("UserAuth -> cleanCacheUser -> nothing to remove")
	}

	nextLoad := time.Minute * 5
	log.Printf("UserAuth -> cleanCacheUser() -> next exec after %s\n", nextLoad)
	timer := time.NewTimer(nextLoad)
	<-timer.C
	log.Println("UserAuth -> cleanCacheUser() -> timer expired")

	auth.cleanCacheUser()
}

func (auth *UserAuth) mainCheck(w http.ResponseWriter, r *http.Request) (bool, *model.User, *string, *string, *string) {
	vHeader := r.Header.Get("v")
	var appVersion *string
	if len(vHeader) > 0 {
		appVersion = &vHeader
	}

	//get the tokens
	token, tokenSourceType, csrfToken, err := auth.getTokens(r)
	if err != nil {
		log.Printf("error gettings tokens - %s", err)

		auth.responseInternalServerError(w)
		return false, nil, nil, nil, nil
	}

	//check if all input data is available
	if token == nil || len(*token) == 0 {
		auth.responseBadRequest(w)
		return false, nil, nil, nil, nil
	}
	rawToken := *token //we have token
	if *tokenSourceType == "cookie" && (csrfToken == nil || len(*csrfToken) == 0) {
		//if the token is sent via cookie then we must have csrf token as well
		auth.responseBadRequest(w)
		return false, nil, nil, nil, nil
	}

	// determine the token type: 1 for shibboleth, 2 for phone, 3 for auth access token
	// 1 & 2 are deprecated but we support them for back compatability
	tokenType, err := auth.getTokenType(rawToken)
	if err != nil {
		auth.responseUnauthorized(err.Error(), w)
		return false, nil, nil, nil, nil
	}
	if !(*tokenType == 1 || *tokenType == 2 || *tokenType == 3) {
		auth.responseUnauthorized("not supported token type", w)
		return false, nil, nil, nil, nil
	}

	// process the token - validate it, extract the user identifier
	var externalID string
	var authType string

	switch *tokenType {
	case 1:
		//support this for back compatability
		uin, err := auth.processShibbolethToken(rawToken)
		if err != nil {
			auth.responseUnauthorized(err.Error(), w)
			return false, nil, nil, nil, nil
		}
		externalID = *uin
		authType = "shibboleth"
	case 2:
		//support this for back compatability
		phone, err := auth.processPhoneToken(rawToken)
		if err != nil {
			auth.responseUnauthorized(err.Error(), w)
			return false, nil, nil, nil, nil
		}
		externalID = *phone
		authType = "phone"
	case 3:
		//mobile app sends just token, the browser sends token + csrf token

		csrfCheck := false
		if *tokenSourceType == "cookie" {
			csrfCheck = true
		}

		tokenData, err := auth.processAccessToken(rawToken, csrfCheck, csrfToken)
		if err != nil {
			auth.responseUnauthorized(err.Error(), w)
			return false, nil, nil, nil, nil
		}

		tokenAuth := tokenData.Auth
		if tokenAuth == "oidc" {
			externalID = tokenData.UID
			authType = "shibboleth"
		} else if tokenAuth == "rokwire_phone" {
			externalID = tokenData.UID
			authType = "phone"
		} else {
			auth.responseUnauthorized("not supported token auth type", w)
			return false, nil, nil, nil, nil
		}
	}

	//TODO - refactor!!!
	// if phone token then treat it as shibboleth
	if authType == "phone" {
		foundedUIN := auth.findUINByPhone(externalID)
		if foundedUIN == nil {
			//not found, it means that this phone is not added, so return unauthorized
			auth.responseUnauthorized(fmt.Sprintf("%s phone is not added in the system", externalID), w)
			return false, nil, nil, nil, nil
		}
		//it was found
		externalID = *foundedUIN
		authType = "shibboleth"
	}

	// get the user for the provided external id.
	user, err := auth.getUser(externalID)
	if err != nil {
		log.Printf("error getting an user for external id - %s\n", err)

		auth.responseInternalServerError(w)
		return false, nil, nil, nil, nil
	}

	// we do not have a such user yet but the ID token is valid so return ok
	if user == nil {
		return true, nil, &externalID, &authType, nil
	}

	// once we have the user we must check if we need to create a default account, every user must have at least one default account
	user, err = auth.createDefaultAccountIfNeeded(*user)
	if err != nil {
		log.Printf("error creating a default account for user - %s - %s\n", utils.GetLogUUIDValue(user.ID), err)

		auth.responseInternalServerError(w)
		return false, nil, nil, nil, nil
	}

	return true, user, &externalID, &authType, appVersion
}

//token source type - cookie and header
func (auth *UserAuth) getTokens(r *http.Request) (*string, *string, *string, error) {
	//1. Check if there is a cookie
	cookie, err := r.Cookie("rokwire-access")
	if err == nil && cookie != nil && len(cookie.Value) > 0 {
		//there is a cookie
		tokenSourceType := "cookie"
		csrfToken := r.Header.Get("CSRF")

		return &cookie.Value, &tokenSourceType, &csrfToken, nil
	}

	//2. Check if there is a token in the Authorization header
	authorizationHeader := r.Header.Get("Authorization")
	if len(authorizationHeader) <= 0 {
		//no authorization
		return nil, nil, nil, nil
	}
	splitAuthorization := strings.Fields(authorizationHeader)
	if len(splitAuthorization) != 2 {
		//bad authorization
		return nil, nil, nil, nil
	}
	// expected - Bearer 1234
	if splitAuthorization[0] != "Bearer" {
		//bad authorization
		return nil, nil, nil, nil
	}
	token := splitAuthorization[1]
	tokenSourceType := "header"
	return &token, &tokenSourceType, nil, nil
}

func (auth *UserAuth) userCheck(w http.ResponseWriter, r *http.Request) (bool, *model.User, *string, *string, *string) {
	//apply main check
	ok, user, externalID, authType, appVersion := auth.mainCheck(w, r)
	if !ok {
		return false, nil, nil, nil, nil
	}

	return true, user, externalID, authType, appVersion
}

func (auth *UserAuth) userAccountsCheck(w http.ResponseWriter, r *http.Request) (bool, *model.User, *model.Account) {
	//apply main check
	ok, user, _, _, _ := auth.mainCheck(w, r)
	if !ok {
		return false, nil, nil
	}
	if user == nil {
		//it is ok but there is no user
		return true, nil, nil
	}

	// determine the used user account
	/// get the account id from the header
	passedAccountID := r.Header.Get("ROKWIRE-ACC-ID")
	/// if account id not passed then we use the default one, we support the old client version
	if len(passedAccountID) == 0 {
		defAccount := user.GetDefaultAccount()
		if defAccount == nil {
			log.Printf("error getting default account - %s\n", utils.GetLogUUIDValue(user.ID))

			auth.responseInternalServerError(w)
			return false, nil, nil
		}
		passedAccountID = defAccount.ID
	}
	/// now get the user account for the provided account id
	account := user.GetAccount(passedAccountID)
	if account == nil {
		auth.responseForbbiden(fmt.Sprintf("Security - %s is trying to use account %s", utils.GetLogUUIDValue(user.ID), passedAccountID), w)
		return false, nil, nil
	}

	return true, user, account
}

func (auth *UserAuth) createDefaultAccountIfNeeded(current model.User) (*model.User, error) {
	if current.HasDefaultAccount() {
		//we have a default account, so do nothing
		return &current, nil
	}

	log.Printf("createDefaultAccountIfNeeded -> we need to create default account!")

	//1. remove it from the cache
	auth.deleteCacheUser(current.ExternalID)

	//2. create default account
	user, err := auth.app.CreateDefaultAccount(current.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//mobile app sends just token, the browser sends token + csrf token
func (auth *UserAuth) processAccessToken(token string, csrfCheck bool, csrfToken *string) (*tokenData, error) {

	//1. apply csrf check
	if csrfCheck {

		if csrfToken == nil || len(*csrfToken) == 0 {
			return nil, errors.New("missing csrf token")
		}

		crsfTokenData, err := auth.validateToken(*csrfToken, "csrf")
		if err != nil {
			log.Printf("error trying to validate csrf token - %s", err)
			return nil, err
		}

		if crsfTokenData == nil {
			log.Printf("not valid csrf token - %s", *csrfToken)
			return nil, errors.New("not valid csrf token")
		}
	}

	//2. apply access token check
	accessTokenData, err := auth.validateToken(token, "access")
	if err != nil {
		log.Printf("error trying to validate access token - %s", err)
		return nil, err
	}

	if accessTokenData == nil {
		log.Printf("not valid access token - %s", token)
		return nil, errors.New("not valid access token")
	}

	return accessTokenData, nil
}

//token type - access or csrf
func (auth *UserAuth) validateToken(token string, tokenType string) (*tokenData, error) {
	//extract the data - header and payload
	tokenSegments := strings.Split(token, ".")
	if len(tokenSegments) != 3 {
		return nil, errors.New("token segments count is != 3")
	}
	//header data
	headerData, err := jwt.DecodeSegment(tokenSegments[0])
	if err != nil {
		log.Printf("error decoding the header segment - %s", err)
		return nil, err
	}
	headerMap := make(map[string]string)
	err = json.Unmarshal(headerData, &headerMap)
	if err != nil {
		log.Println("error unmarshaling the header data" + err.Error())
		return nil, err
	}

	//payload
	payloadData, err := jwt.DecodeSegment(tokenSegments[1])
	if err != nil {
		log.Printf("error decoding the payload segment - %s", err)
		return nil, err
	}
	var tokenData *tokenData
	err = json.Unmarshal(payloadData, &tokenData)
	if err != nil {
		log.Println("error unmarshaling the payload data" + err.Error())
		return nil, err
	}

	//check issuer
	if tokenData.ISS != auth.Issuer {
		log.Printf("issuer does not match: - %s", tokenData.ISS)
		return nil, errors.New("issuer does not match:" + tokenData.ISS)
	}

	//check keys
	kid := headerMap["kid"]
	if len(kid) == 0 {
		log.Println("kid header is missing")
		return nil, errors.New("kid header is missing")
	}
	jwkMatch := auth.Keys.LookupKeyID(kid)
	if len(jwkMatch) == 0 {
		log.Printf("no matching kid found")
		return nil, errors.New("no matching kid found")
	} else if len(jwkMatch) > 1 {
		log.Printf("multiple matching kids found")
		return nil, errors.New("multiple matching kids found")
	}
	publicKey := jwkMatch[0].(jwk.RSAPublicKey)

	//validate
	jwk := rsa.PublicKey{}
	parsedToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if err := publicKey.Raw(&jwk); err != nil {
			log.Println("failed to create public key:", err)
			return nil, err
		}
		return &jwk, nil
	})
	if err != nil {
		log.Printf("error parse/validate token - %s", err)
		return nil, err
	}
	if !parsedToken.Valid {
		log.Printf("not valid token - %s", token)
		return nil, errors.New("not valid token:" + token)
	}

	//check token type
	if tokenData.Type != tokenType {
		log.Printf("invalid type %s", tokenData.Type)
		return nil, errors.New("invalid type - " + token)
	}

	return tokenData, nil
}

func (auth *UserAuth) processShibbolethToken(token string) (*string, error) {
	// Validate the token
	idToken, err := auth.appIDTokenVerifier.Verify(context.Background(), token)
	if err != nil {
		log.Printf("error validating token - %s\n", err)
		return nil, err
	}

	// Get the user data from the token
	var userData shData
	if err := idToken.Claims(&userData); err != nil {
		log.Printf("error getting user data from token - %s\n", err)
		return nil, err
	}
	//we must have UIuceduUIN
	if userData.UIuceduUIN == nil {
		log.Printf("missing uiuceuin data in the token - %s\n", token)
		return nil, errors.New("missing uiuceuin data in the token")
	}
	return userData.UIuceduUIN, nil
}

func (auth *UserAuth) findUINByPhone(phone string) *string {
	rosters := auth.getRosters()
	if len(rosters) == 0 {
		return nil
	}

	for _, item := range rosters {
		cPhone := item["phone"]
		if cPhone == phone {
			uin := item["uin"]
			return &uin
		}
	}
	return nil
}

func (auth *UserAuth) processPhoneToken(token string) (*string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.phoneAuthSecret), nil
	})
	if err != nil {
		return nil, err
	}

	for key, val := range claims {
		if key == "phoneNumber" {
			phoneValue := val.(string)
			return &phoneValue, nil
		}
	}
	return nil, errors.New("there is no phoneNumber claim in the phone token")
}

// type: 1 for shibboleth, 2 for phone, 3 for auth access token
// 1 & 2 are deprecated but we support them for back compatability
func (auth *UserAuth) getTokenType(token string) (*int, error) {
	parser := new(jwt.Parser)
	claims := jwt.MapClaims{}
	_, _, err := parser.ParseUnverified(token, claims)
	if err != nil {
		return nil, err
	}

	for key := range claims {
		if key == "uiucedu_uin" {
			tokenType := 1
			return &tokenType, nil
		}
		if key == "phoneNumber" {
			tokenType := 2
			return &tokenType, nil
		}
		if key == "uid" {
			tokenType := 3
			return &tokenType, nil
		}
	}
	return nil, errors.New("not supported token type")
}

func (auth *UserAuth) createAppUser(externalID string, uuid string, publicKey string,
	consent bool, consentVaccine bool, exposureNotification bool, rePost bool, encryptedKey *string, encryptedBlob *string, encryptedPK *string) error {

	_, err := auth.app.CreateAppUser(externalID, uuid, publicKey, consent, consentVaccine, exposureNotification, rePost, encryptedKey, encryptedBlob, encryptedPK)
	if err != nil {
		return err
	}

	return nil
}

func (auth *UserAuth) updateAppUser(user model.User, uuid string, publicKey string, consent bool, consentVaccine bool, exposureNotification bool, rePost *bool,
	encryptedKey *string, encryptedBlob *string, encryptedPK *string) error {

	//1. remove it from the cache
	auth.deleteCacheUser(user.ExternalID)

	//2. Set the new values
	user.UUID = uuid
	user.PublicKey = publicKey
	user.Consent = consent
	user.ConsentVaccine = consentVaccine
	user.ExposureNotification = exposureNotification
	if rePost != nil {
		user.RePost = *rePost
	}
	user.EncryptedKey = encryptedKey
	user.EncryptedBlob = encryptedBlob
	user.EncryptedPK = encryptedPK

	//3. Update the user
	err := auth.app.UpdateUser(&user)
	if err != nil {
		return err
	}

	return nil
}

func (auth *UserAuth) getCachedUser(externalID string) *cacheUser {
	auth.cachedUsersLock.RLock()
	defer auth.cachedUsersLock.RUnlock()

	var cachedUser *cacheUser //to return

	item, _ := auth.cachedUsers.Load(externalID)
	if item != nil {
		cachedUser = item.(*cacheUser)
	}

	//keep the last get time
	if cachedUser != nil {
		cachedUser.lastUsage = time.Now()
		auth.cachedUsers.Store(externalID, cachedUser)
	}

	return cachedUser
}

func (auth *UserAuth) cacheUser(externalID string, user *model.User) {
	auth.cachedUsersLock.RLock()

	cacheUser := &cacheUser{user: user, lastUsage: time.Now()}
	auth.cachedUsers.Store(externalID, cacheUser)

	auth.cachedUsersLock.RUnlock()
}

func (auth *UserAuth) deleteCacheUser(externalID string) {
	auth.cachedUsersLock.RLock()

	auth.cachedUsers.Delete(externalID)

	auth.cachedUsersLock.RUnlock()
}

func (auth *UserAuth) clearCacheUsers() {
	log.Println("UserAuth -> clearCacheUsers")

	auth.cachedUsersLock.RLock()

	auth.cachedUsers = &syncmap.Map{}

	auth.cachedUsersLock.RUnlock()
}

func (auth *UserAuth) setRosters(rosters []map[string]string) {
	auth.rostersLock.RLock()

	auth.rosters = rosters

	auth.rostersLock.RUnlock()
}

func (auth *UserAuth) getRosters() []map[string]string {
	auth.rostersLock.RLock()
	defer auth.rostersLock.RUnlock()

	return auth.rosters
}

func (auth *UserAuth) getUser(externalID string) (*model.User, error) {
	var err error

	//1. First check if cached
	cachedUser := auth.getCachedUser(externalID)
	if cachedUser != nil {
		return cachedUser.user, nil
	}

	//2. Check if we have a such user in the application
	user, err := auth.app.FindUserByExternalID(externalID)
	if err != nil {
		log.Printf("error finding an user for external id - %s\n", err)
		return nil, err
	}
	if user != nil {
		//cache it
		auth.cacheUser(externalID, user)
		return user, nil
	}
	//there is no a such user
	return nil, nil
}

func (auth *UserAuth) responseBadRequest(w http.ResponseWriter) {
	log.Println(fmt.Sprintf("400 - Bad Request"))

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad Request"))
}

func (auth *UserAuth) responseUnauthorized(logInfo string, w http.ResponseWriter) {
	log.Println(fmt.Sprintf("401 - Unauthorized - %s", logInfo))

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Unauthorized"))
}

func (auth *UserAuth) responseInternalServerError(w http.ResponseWriter) {
	log.Println(fmt.Sprintf("500 - Internal Server Error"))

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}

func (auth *UserAuth) responseForbbiden(info string, w http.ResponseWriter) {
	log.Printf("403 - Forbidden - %s", info)

	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Forbidden"))
}

func newUserAuth(app *core.Application, oidcProvider string, oidcAppClientID string,
	phoneAuthSecret string, keys string, issuer string) *UserAuth {

	provider, err := oidc.NewProvider(context.Background(), oidcProvider)
	if err != nil {
		log.Fatalln(err)
	}
	appIDTokenVerifier := provider.Verifier(&oidc.Config{ClientID: oidcAppClientID})

	keysSet, err := jwk.ParseString(keys)
	if err != nil {
		log.Fatalln(err)
	}

	cacheUsers := &syncmap.Map{}
	lock := &sync.RWMutex{}

	cacheRosters := []map[string]string{}
	rostersLock := &sync.RWMutex{}

	auth := UserAuth{app: app, appIDTokenVerifier: appIDTokenVerifier, phoneAuthSecret: phoneAuthSecret, Keys: keysSet, Issuer: issuer,
		cachedUsers: cacheUsers, cachedUsersLock: lock, rosters: cacheRosters, rostersLock: rostersLock}
	return &auth
}
