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

	"github.com/dgrijalva/jwt-go"
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
}

//Start starts the auth module
func (auth *Auth) Start() error {
	auth.adminAuth.start()
	auth.userAuth.start()

	return nil
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

func (auth *Auth) providersCheck(w http.ResponseWriter, r *http.Request) bool {
	return auth.providersAuth.check(w, r)
}

func (auth *Auth) userCheck(w http.ResponseWriter, r *http.Request) (bool, *model.User, *string, *string) {
	return auth.userAuth.check(w, r)
}

func (auth *Auth) updateAppUser(user model.User, uuid string, publicKey string, consent bool, exposureNotification bool, rePost *bool, encryptedKey *string, encryptedBlob *string) error {
	return auth.userAuth.updateAppUser(user, uuid, publicKey, consent, exposureNotification, rePost, encryptedKey, encryptedBlob)
}

func (auth *Auth) createAppUser(externalID string, uuid string, publicKey string, consent bool, exposureNotification bool, rePost bool, encryptedKey *string, encryptedBlob *string) error {
	return auth.userAuth.createAppUser(externalID, uuid, publicKey, consent, exposureNotification, rePost, encryptedKey, encryptedBlob)
}

//NewAuth creates new auth handler
func NewAuth(app *core.Application, appKeys []string, oidcProvider string,
	oidcAppClientID string, appClientID string, webAppClientID string, phoneAuthSecret string, providersAPIKeys []string) *Auth {
	apiKeysAuth := newAPIKeysAuth(appKeys)
	userAuth2 := newUserAuth(app, oidcProvider, oidcAppClientID, phoneAuthSecret)
	adminAuth := newAdminAuth(app, oidcProvider, appClientID, webAppClientID)
	providersAuth := newProviderAuth(providersAPIKeys)

	auth := Auth{apiKeysAuth: apiKeysAuth, userAuth: userAuth2, adminAuth: adminAuth, providersAuth: providersAuth}
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
	//1. Get the token and the group from the request
	authorizationHeader := r.Header.Get("Authorization")
	if len(authorizationHeader) <= 0 {
		auth.responseBadRequest(w)
		return false, nil, "", nil
	}
	group := r.Header.Get("GROUP")
	if len(group) <= 0 {
		auth.responseBadRequest(w)
		return false, nil, "", nil
	}

	splitAuthorization := strings.Fields(authorizationHeader)
	if len(splitAuthorization) != 2 {
		auth.responseBadRequest(w)
		return false, nil, "", nil
	}
	// expected - Bearer 1234
	if splitAuthorization[0] != "Bearer" {
		auth.responseBadRequest(w)
		return false, nil, "", nil
	}
	rawIDToken := splitAuthorization[1]

	//2. Validate the token
	idToken, err := auth.verify(rawIDToken)
	if err != nil {
		log.Printf("error validating token - %s\n", err)

		auth.responseUnauthorized(rawIDToken, w)
		return false, nil, "", nil
	}

	//3. Get the user data from the token
	var userData userData
	if err := idToken.Claims(&userData); err != nil {
		log.Printf("error getting user data from token - %s\n", err)

		auth.responseUnauthorized(rawIDToken, w)
		return false, nil, "", nil
	}
	//we must have UIuceduUIN
	if userData.UIuceduUIN == nil {
		log.Printf("error - missing uiuceuin data in the token - %s\n", err)

		auth.responseUnauthorized(rawIDToken, w)
		return false, nil, "", nil
	}

	//4. Get the user for the provided external id.
	user, err := auth.getUser(*userData.UIuceduUIN)
	if err != nil {
		log.Printf("error getting an user for external id - %s\n", err)

		auth.responseInternalServerError(w)
		return false, nil, "", nil
	}

	shibboAuth := &model.ShibbolethAuth{Uin: *userData.UIuceduUIN, Email: *userData.Email,
		IsMemberOf: userData.UIuceduIsMemberOf}

	//5.
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

	//6. Check if the user is member of the group
	if !user.IsMemberOf(group) {
		auth.responseForbbiden(fmt.Sprintf("Security - %s is trying to access not allowed resource", *userData.Email), w)
		return false, nil, "", nil
	}

	return true, user, group, shibboAuth
}

func (auth *AdminAuth) verify(rawIDToken string) (*oidc.IDToken, error) {
	parser := new(jwt.Parser)
	claims := jwt.MapClaims{}
	_, _, errr := parser.ParseUnverified(rawIDToken, claims)
	if errr != nil {
		return nil, errr
	}

	for key := range claims {
		if key == "aud" {
			audience := claims[key]
			switch audience {
			case auth.appClientID:
				log.Println("AuthAdmin -> app client token")
				return auth.appVerifier.Verify(context.Background(), rawIDToken)
			case auth.webAppClientID:
				log.Println("AuthAdmin -> web app client token")
				return auth.webAppVerifier.Verify(context.Background(), rawIDToken)
			default:
				return nil, errors.New("there is an issue with the audience")
			}

		}

	}
	return nil, errors.New("there is an issue with the audience")
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

//UserAuth entity
type UserAuth struct {
	app *core.Application

	//shibboleth
	appIDTokenVerifier *oidc.IDTokenVerifier

	//phone
	phoneAuthSecret string

	cachedUsers     *syncmap.Map //cache users while active - 5 minutes timeout
	cachedUsersLock *sync.RWMutex
}

func (auth *UserAuth) start() {
	go auth.cleanCacheUser()
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

func (auth *UserAuth) check(w http.ResponseWriter, r *http.Request) (bool, *model.User, *string, *string) {
	authorizationHeader := r.Header.Get("Authorization")
	if len(authorizationHeader) <= 0 {
		auth.responseBadRequest(w)
		return false, nil, nil, nil
	}
	splitAuthorization := strings.Fields(authorizationHeader)
	if len(splitAuthorization) != 2 {
		auth.responseBadRequest(w)
		return false, nil, nil, nil
	}
	// expected - Bearer 1234
	if splitAuthorization[0] != "Bearer" {
		auth.responseBadRequest(w)
		return false, nil, nil, nil
	}
	rawIDToken := splitAuthorization[1]

	// determine the token type - 1 for shibboleth, 2 for phone
	tokenType, err := auth.getTokenType(rawIDToken)
	if err != nil {
		auth.responseUnauthorized(err.Error(), w)
		return false, nil, nil, nil
	}
	if !(*tokenType == 1 || *tokenType == 2) {
		auth.responseUnauthorized("not supported token type", w)
		return false, nil, nil, nil
	}

	// process the token - validate it, extract the user identifier
	var externalID string
	var authType string
	if *tokenType == 1 {
		uin, err := auth.processShibbolethToken(rawIDToken)
		if err != nil {
			auth.responseUnauthorized(err.Error(), w)
			return false, nil, nil, nil
		}
		externalID = *uin
		authType = "shibboleth"
	} else if *tokenType == 2 {
		phone, err := auth.processPhoneToken(rawIDToken)
		if err != nil {
			auth.responseUnauthorized(err.Error(), w)
			return false, nil, nil, nil
		}
		externalID = *phone
		authType = "phone"
	}

	// get the user for the provided external id.
	user, err := auth.getUser(externalID)
	if err != nil {
		log.Printf("error getting an user for external id - %s\n", err)

		auth.responseInternalServerError(w)
		return false, nil, nil, nil
	}
	return true, user, &externalID, &authType
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

// type: 1 for shibboleth, 2 for phone
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
	}
	return nil, errors.New("not supported token type")
}

func (auth *UserAuth) createAppUser(externalID string, uuid string, publicKey string,
	consent bool, exposureNotification bool, rePost bool, encryptedKey *string, encryptedBlob *string) error {

	_, err := auth.app.CreateAppUser(externalID, uuid, publicKey, consent, exposureNotification, rePost, encryptedKey, encryptedBlob)
	if err != nil {
		return err
	}

	return nil
}

func (auth *UserAuth) updateAppUser(user model.User, uuid string, publicKey string, consent bool, exposureNotification bool, rePost *bool,
	encryptedKey *string, encryptedBlob *string) error {

	//1. remove it from the cache
	auth.deleteCacheUser(user.ExternalID)

	//2. Set the new values
	user.UUID = uuid
	user.PublicKey = publicKey
	user.Consent = consent
	user.ExposureNotification = exposureNotification
	if rePost != nil {
		user.RePost = *rePost
	}
	user.EncryptedKey = encryptedKey
	user.EncryptedBlob = encryptedBlob

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

func newUserAuth(app *core.Application, oidcProvider string, oidcAppClientID string,
	phoneAuthSecret string) *UserAuth {
	provider, err := oidc.NewProvider(context.Background(), oidcProvider)
	if err != nil {
		log.Fatalln(err)
	}

	appIDTokenVerifier := provider.Verifier(&oidc.Config{ClientID: oidcAppClientID})

	cacheUsers := &syncmap.Map{}
	lock := &sync.RWMutex{}

	auth := UserAuth{app: app, appIDTokenVerifier: appIDTokenVerifier,
		phoneAuthSecret: phoneAuthSecret, cachedUsers: cacheUsers, cachedUsersLock: lock}
	return &auth
}
