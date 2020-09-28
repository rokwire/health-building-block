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

package core

import (
	"errors"
	"health/core/model"
	"health/utils"
	"log"
	"sync"
	"time"
)

//Application represents the core application code based on hexagonal architecture
type Application struct {
	version string
	build   string

	Services       Services       //expose to the drivers adapters
	Administration Administration //expose to the drivrs adapters

	dataProvider DataProvider
	sender       Sender
	messaging    Messaging
	profileBB    ProfileBuildingBlock
	audit        Audit

	storage Storage

	//cache config data
	cvLock              *sync.RWMutex
	cachedCovid19Config *model.COVID19Config

	listeners []ApplicationListener

	supportedVersions []string
}

//Start starts the core part of the application
func (app *Application) Start() {
	//set storage listener
	storageListener := storageListenerImpl{app: app}
	app.storage.SetStorageListener(&storageListener)

	//cache the configs
	app.loadCovid19Config()

	go app.loadNewsData()
	//Disable the resource data loading as we cannot map the new created data
	//go app.loadResourcesData()

	go app.setupLocationWaitTimeColorTimer()
}

//AddListener adds application listener
func (app *Application) AddListener(listener ApplicationListener) {
	log.Println("Application -> AddListener")

	app.listeners = append(app.listeners, listener)
}

func (app *Application) setupLocationWaitTimeColorTimer() {
	log.Println("Application -> setupLocationWaitTimeColorTimer")

	//TODO

	ticker := time.NewTicker(30 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				app.checkLocationsWaitTimesColors()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (app *Application) checkLocationsWaitTimesColors() {
	log.Println("Application -> checkLocationsWaitTimesColors")

	//TODO
}

func (app *Application) notifyListeners(message string, data interface{}) {
	go func() {
		for _, listener := range app.listeners {
			if message == "onClearUserData" {
				listener.OnClearUserData(data.(model.User))
			} else if message == "onUserUpdated" {
				listener.OnUserUpdated(data.(model.User))
			}
		}
	}()
}

func (app *Application) checkAppVersion(v *string) (*string, error) {
	if v == nil {
		//use the latest version if not provided
		latest := app.getLatestVersion()
		return &latest, nil
	}

	//check if supported
	supported := app.isVersionSupported(*v)
	if !supported {
		return nil, errors.New("the provided version is not supported")
	}

	//the version is ok
	return v, nil
}

func (app *Application) getLatestVersion() string {
	return "2.6"
}

func (app *Application) isVersionSupported(v string) bool {
	for _, current := range app.supportedVersions {
		if current == v {
			return true
		}
	}
	return false
}

func (app *Application) loadCovid19Config() {
	log.Println("Load Covid19 config")

	covid19Config, err := app.storage.ReadCovid19Config()
	if err != nil {
		log.Printf("Error reading the covid19 config %s", err)
	}
	app.setCachedCovid19Config(covid19Config)
}

func (app *Application) setCachedCovid19Config(covid19 *model.COVID19Config) {
	app.cvLock.RLock()
	app.cachedCovid19Config = covid19
	app.cvLock.RUnlock()
}

func (app *Application) getCachedCovid19Config() *model.COVID19Config {
	app.cvLock.RLock()
	defer app.cvLock.RUnlock()

	return app.cachedCovid19Config
}

func (app *Application) loadNewsData() {
	log.Println("loadNewsData() -> load data from the provider")

	//1. load the provider data
	providerData, err := app.dataProvider.LoadNews()
	if err != nil {
		log.Printf("loadNewsData() -> error on loading the provider data %s", err)

		go app.execNewsTimer()
		return
	}

	//2. find the latest news date
	var latestDate *time.Time
	newsList, err := app.storage.ReadNews(0)
	if err != nil {
		log.Printf("loadNewsData() -> error on finding the latest news date %s", err)

		go app.execNewsTimer()
		return
	}
	if newsList != nil && len(newsList) > 0 {
		latestItem := newsList[0]

		log.Printf("loadNewsData() -> latest item is %s with date %s", latestItem.Title, latestItem.Date)

		latestDate = &latestItem.Date
	} else {
		log.Println("loadNewsData() -> the news list is empty")
	}

	//3. find only the new items
	newItems := app.findNewNewsItems(providerData, latestDate)

	var addedNews []*model.News

	//4. store the new items
	newItemsCount := len(newItems)
	if newItems != nil && newItemsCount > 0 {
		log.Printf("loadNewsData() -> there are %d new news items\n", newItemsCount)

		for _, item := range newItems {
			description := utils.ModifyHTMLContent(item.Description)
			htmlContent := utils.ModifyHTMLContent(item.ContentEncoded)
			created, err := app.storage.CreateNews(item.PubDate, item.Title, description, htmlContent, nil)
			if err != nil {
				log.Printf("loadNewsData() -> error on saving news - %s\n", item.Title)
			} else {
				log.Printf("loadNewsData() -> successfully saved news - %s\n", item.Title)
				addedNews = append(addedNews, created)
			}
		}
	} else {
		log.Println("loadNewsData() -> there is no new news items")
	}

	//stop sending emails
	//5. notify the recipients for the added items
	//if len(addedNews) > 0 {
	//	go app.sender.SendForNews(addedNews)
	//}

	go app.execNewsTimer()
}

func (app *Application) execNewsTimer() {
	periodInMinutes := app.getCachedCovid19Config().NewsUpdatePeriod
	nextLoad := time.Minute * time.Duration(periodInMinutes)
	log.Printf("execNewsTimer() -> next exec after %s\n", nextLoad)
	timer := time.NewTimer(nextLoad)
	<-timer.C
	log.Println("execNewsTimer() -> timer expired")
	app.loadNewsData()
}

func (app *Application) loadResourcesData() {
	/*log.Println("loadResourcesData() -> load data from the provider")

	//1. load the provider data
	providerData, err := app.dataProvider.LoadResources()
	if err != nil {
		log.Printf("loadResourcesData() -> error on loading the provider data %s", err)

		go app.execResourcesTimer()
		return
	}

	//2. Load the resoruces from the storage. They are prety small size
	resourceList, err := app.storage.ReadAllResources()
	if err != nil {
		log.Printf("loadResourcesData() -> error on reading all the resources %s", err)

		go app.execResourcesTimer()
		return
	}

	//3. find only the new items
	newResources := app.findNewResourcesItems(resourceList, providerData)

	var addedResources []*model.Resource

	//4. store the new items
	newItemsCount := len(newResources)
	if newResources != nil && newItemsCount > 0 {
		log.Printf("loadResourcesData() -> there are %d new resource items\n", newItemsCount)

		for _, item := range newResources {
			created, err := app.storage.CreateResource(item.Title, item.Link)
			if err != nil {
				log.Printf("loadResourcesData() -> error on saving resource - %s\n", item.Title)
			} else {
				log.Printf("loadResourcesData() -> successfully saved resource - %s\n", item.Title)
				addedResources = append(addedResources, created)
			}
		}
	} else {
		log.Println("loadResourcesData() -> there is no new resource items")
	}

	//5. notify the recipients for the added resource items
	if len(addedResources) > 0 {
		go app.sender.SendForResources(addedResources)
	}

	go app.execResourcesTimer() */

}

func (app *Application) findNewResourcesItems(list []*model.Resource, providerList []ProviderResource) []ProviderResource {
	log.Println("findNewResourcesItems() -> start")

	var result []ProviderResource

	if providerList != nil {
		for _, pItem := range providerList {
			if !app.containsResourceItem(list, pItem.Link) {
				result = append(result, pItem)
			}
		}
	}

	log.Println("findNewResourcesItems() -> end")
	return result
}

func (app *Application) containsResourceItem(list []*model.Resource, link string) bool {
	if list != nil {
		for _, item := range list {
			if item.Link == link {
				return true
			}
		}
	}
	return false
}

func (app *Application) execResourcesTimer() {
	nextLoad := time.Hour * 1
	timer := time.NewTimer(nextLoad)
	<-timer.C
	log.Println("execResourcesTimer() -> timer expired")
	app.loadResourcesData()
}

func (app *Application) findNewNewsItems(list []ProviderNews, latestDate *time.Time) []ProviderNews {
	log.Println("findNewNewsItems() -> start")

	if latestDate == nil {
		//if there is no latest date then it means that the list is empty, so return the full list
		return list
	}

	var result []ProviderNews
	if list != nil {
		for _, item := range list {
			if item.PubDate.Unix() > latestDate.Unix() {
				result = append(result, item)
			}
		}
	}

	log.Println("findNewNewsItems() -> end")
	return result
}

//FindUserByShibbolethID finds an user for the provided shibboleth id
func (app *Application) FindUserByShibbolethID(shibbolethID string) (*model.User, error) {
	user, err := app.storage.FindUserByShibbolethID(shibbolethID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//FindUserByExternalID finds an user for the provided external id
func (app *Application) FindUserByExternalID(externalID string) (*model.User, error) {
	user, err := app.storage.FindUserByExternalID(externalID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//CreateAppUser creates an app user
func (app *Application) CreateAppUser(externalID string, uuid string, publicKey string,
	consent bool, exposureNotification bool, rePost bool, encryptedKey *string, encryptedBlob *string) (*model.User, error) {
	user, err := app.storage.CreateUser(nil, externalID, uuid, publicKey, consent, exposureNotification, rePost, encryptedKey, encryptedBlob)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//CreateAdminAppUser creates an admin app user
func (app *Application) CreateAdminAppUser(shibboAuth *model.ShibbolethAuth) (*model.User, error) {
	externalID := "a_" + shibboAuth.Uin //TODO
	user, err := app.storage.CreateUser(shibboAuth, externalID, "", "", false, false, false, nil, nil)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//UpdateUser updates the user
func (app *Application) UpdateUser(user *model.User) error {
	err := app.storage.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) getEHistoriesByUserID(userID string) ([]*model.EHistory, error) {
	histories, err := app.storage.FindEHistories(userID)
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (app *Application) createЕHistory(userID string, date time.Time, eType string, encryptedKey string, encryptedBlob string) (*model.EHistory, error) {
	history, err := app.storage.CreateEHistory(userID, date, eType, encryptedKey, encryptedBlob)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (app *Application) createManualЕHistory(userID string, date time.Time, encryptedKey string, encryptedBlob string, encryptedImageKey *string, encryptedImageBlob *string,
	countyID *string, locationID *string) (*model.EHistory, error) {
	history, err := app.storage.CreateManualЕHistory(userID, date, encryptedKey, encryptedBlob, encryptedImageKey, encryptedImageBlob, countyID, locationID)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (app *Application) getProviders() ([]*model.Provider, error) {
	providers, err := app.storage.ReadAllProviders()
	if err != nil {
		return nil, err
	}
	return providers, nil
}

func (app *Application) findCounties(f *utils.Filter) ([]*model.County, error) {
	counties, err := app.storage.FindCounties(f)
	if err != nil {
		return nil, err
	}
	return counties, nil
}

func (app *Application) getCounty(ID string) (*model.County, error) {
	county, err := app.storage.FindCounty(ID)
	if err != nil {
		return nil, err
	}
	return county, nil
}

func (app *Application) containsCountyStatusP(ID string, list []*model.CountyStatus) bool {
	if list == nil {
		return false
	}
	for _, item := range list {
		if item.ID == ID {
			return true
		}
	}
	return false
}

func (app *Application) containsCountyStatus(ID string, list []model.CountyStatus) bool {
	if list == nil {
		return false
	}
	for _, item := range list {
		if item.ID == ID {
			return true
		}
	}
	return false
}

func (app *Application) getSymptomGroups() ([]*model.SymptomGroup, error) {
	symptomGroups, err := app.storage.ReadAllSymptomGroups()
	if err != nil {
		return nil, err
	}
	return symptomGroups, nil
}

//NewApplication creates new Application
func NewApplication(version string, build string, dataProvider DataProvider, sender Sender, messaging Messaging, profileBB ProfileBuildingBlock, storage Storage, audit Audit) *Application {
	cvLock := &sync.RWMutex{}
	listeners := []ApplicationListener{}

	supportedVersion := []string{"2.6"}

	application := Application{version: version, build: build, dataProvider: dataProvider, sender: sender, messaging: messaging,
		profileBB: profileBB, storage: storage, audit: audit, cvLock: cvLock, listeners: listeners, supportedVersions: supportedVersion}

	//add the drivers ports/interfaces
	application.Services = &servicesImpl{app: &application}
	application.Administration = &administrationImpl{app: &application}

	return &application
}
