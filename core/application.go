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
	"fmt"
	"health/core/model"
	"health/utils"
	"log"
	"strings"
	"sync"
	"time"

	idgen "github.com/google/uuid"
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

	//cache app versions
	avLock            *sync.RWMutex
	cachedAppVersions []string

	listeners []ApplicationListener
}

//Start starts the core part of the application
func (app *Application) Start() {
	//set storage listener
	storageListener := storageListenerImpl{app: app}
	app.storage.SetStorageListener(&storageListener)

	//cache the configs
	app.loadCovid19Config()

	//cache the app versions
	app.loadAppVersions()

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

	//determine the first moment of the wait time color check
	//we check it every 30 minutes
	now := time.Now()
	currentMinutes := now.Minute()
	currentSecconds := now.Second()
	var desiredMoment int
	if currentMinutes < 30 {
		log.Println("Application -> setupLocationWaitTimeColorTimer -> desired is 30")
		desiredMoment = 30
	} else {
		log.Println("Application -> setupLocationWaitTimeColorTimer -> desired is 60")
		desiredMoment = 60
	}

	desiredMomentInSec := desiredMoment * 60
	currentMomentInSec := (currentMinutes * 60) + currentSecconds
	//we add 5 seconds which is insignificant from user point of view but we quarantee that the check is in the desired 30 minutes interval
	difference := (desiredMomentInSec - currentMomentInSec) + 5
	duration := time.Second * time.Duration(difference)
	log.Printf("Application -> setupLocationWaitTimeColorTimer -> start after - %s", duration)
	timer := time.NewTimer(duration)
	<-timer.C

	//check it for first time
	go app.checkLocationsWaitTimesColors()

	//check it every 30 minutes
	ticker := time.NewTicker(30 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go app.checkLocationsWaitTimesColors()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	//close the quit channel when want to cancel the ticker.
}

func (app *Application) checkLocationsWaitTimesColors() {
	log.Println("Application -> checkLocationsWaitTimesColors")

	// load locations
	locations, err := app.storage.ReadAllLocations()
	if err != nil {
		log.Printf("error loading locations for wait time color check - %s", err)
	}

	for _, loc := range locations {
		app.checkLocationWaitTimeColor(loc)
	}
}

func (app *Application) checkLocationWaitTimeColor(location *model.Location) {
	log.Printf("Application -> checkLocationWaitTimeColor for %s with timezone %s", location.Name, location.Timezone)

	//find the day of the week and the passed seconds within the day
	timeLocation, err := time.LoadLocation(location.Timezone)
	if err != nil {
		log.Printf("Error getting time location:%s\n", err.Error())
	}
	now := time.Now().In(timeLocation)
	nowWeekDay := now.Weekday()
	nowMomentInSec := (now.Hour() * 60 * 60) + (now.Minute() * 60) + now.Second()
	log.Printf("... -> now week day - %s, now moment in secs - %d\n", nowWeekDay, nowMomentInSec)

	isLocationOpen := app.isLocationOpen(location, nowWeekDay.String(), nowMomentInSec)
	if isLocationOpen {
		log.Printf("... -> %s is OPEN, set it to green only if nil or grey\n", location.Name)
		if location.WaitTimeColor == nil || *location.WaitTimeColor == "grey" {
			log.Println("... -> setting it to green because the current wait time color is nil or grey")

			waitTimeColor := "green"
			location.WaitTimeColor = &waitTimeColor
			err = app.storage.SaveLocation(location)
			if err != nil {
				log.Printf("error saving a location after setting green wait time color - %s", err)
			} else {
				log.Printf("... -> Successfully set green wait time color for %s\n", location.Name)
			}
		} else {
			log.Printf("... -> nothing to set because the current wait time color is %s\n", *location.WaitTimeColor)
		}
	} else {
		log.Printf("... -> %s is CLOSED, set it to grey if not grey\n", location.Name)
		if location.WaitTimeColor == nil || *location.WaitTimeColor != "grey" {
			log.Println("... -> setting it to gray because the current wait time color is nil or not grey")

			waitTimeColor := "grey"
			location.WaitTimeColor = &waitTimeColor
			err = app.storage.SaveLocation(location)
			if err != nil {
				log.Printf("error saving a location after setting grey wait time color - %s", err)
			} else {
				log.Printf("... -> Successfully set grey wait time color for %s\n", location.Name)
			}
		} else {
			log.Println("... -> nothing to set because the current wait time color is grey")
		}
	}

}

func (app *Application) isLocationOpen(location *model.Location, day string, passedSeconds int) bool {
	daysOfOperations := location.DaysOfOperation
	if len(daysOfOperations) == 0 {
		return false
	}

	//check if the location is open this day
	var operationDay *model.OperationDay
	for _, current := range daysOfOperations {
		if day == current.Name {
			operationDay = &current
			break
		}
	}
	if operationDay == nil {
		return false
	}

	//check if the location is open this moment in the day
	openTime, err := time.Parse("03:04pm", operationDay.OpenTime)
	if err != nil {
		log.Printf("error parsing open time - %s", openTime)
	}
	openTimeInSec := (openTime.Hour() * 60 * 60) + (openTime.Minute() * 60) + openTime.Second()
	log.Printf("... open moment is secs - %d", openTimeInSec)

	closeTime, err := time.Parse("03:04pm", operationDay.CloseTime)
	if err != nil {
		log.Printf("error parsing close time - %s", closeTime)
	}
	closeTimeInSec := (closeTime.Hour() * 60 * 60) + (closeTime.Minute() * 60) + closeTime.Second()
	log.Printf("... close moment is secs - %d", closeTimeInSec)
	if !(passedSeconds >= openTimeInSec && passedSeconds < closeTimeInSec) {
		return false
	}

	//it is open
	return true
}

func (app *Application) notifyListeners(message string, data interface{}) {
	go func() {
		for _, listener := range app.listeners {
			if message == "onClearUserData" {
				listener.OnClearUserData(data.(model.User))
			} else if message == "onUserUpdated" {
				listener.OnUserUpdated(data.(model.User))
			} else if message == "onRostersUpdated" {
				listener.OnRostersUpdated()
			}
		}
	}()
}

func (app *Application) checkAppVersion(v *string) (*string, error) {
	//use the latest version if not provided
	if v == nil {
		latest := app.getLatestVersion()
		return &latest, nil
	}

	//check if it matches
	matches, version := app.isVersionSupported(*v)
	if matches {
		return version, nil
	}

	//if it does not match then use the latest one which is less that the desired one
	//the versions are sorted as the latest one is on possition 0
	for _, current := range app.getCachedAppVersions() {
		if utils.IsVersionLess(current, *v) {
			return &current, nil
		}
	}

	return nil, errors.New("Not supported version")
}

func (app *Application) getLatestVersion() string {
	//the versions list is sorted, the first element is the latest one
	return app.getCachedAppVersions()[0]
}

func (app *Application) isVersionSupported(v string) (bool, *string) {
	//if the input is 2.8.0 then we search for 2.8 because the system works with the short view when patch is 0
	forSearch := v
	elements := strings.Split(v, ".")
	elementsCount := len(elements)
	if !(elementsCount == 2 || elementsCount == 3) {
		return false, nil
	}
	lastElement := elements[elementsCount-1]
	if elementsCount == 3 && lastElement == "0" {
		forSearch = fmt.Sprintf("%s.%s", elements[0], elements[1])
	}

	//search for it
	for _, current := range app.getCachedAppVersions() {
		if current == forSearch {
			return true, &current
		}
	}
	return false, nil
}

func (app *Application) loadAppVersions() {
	log.Println("Load App versions")

	versions, err := app.storage.ReadAllAppVersions()
	if err != nil {
		log.Printf("Error reading the app versions %s", err)
	}
	app.setCachedAppVersions(versions)
}

func (app *Application) setCachedAppVersions(versions []string) {
	app.avLock.RLock()
	app.cachedAppVersions = versions
	app.avLock.RUnlock()
}

func (app *Application) getCachedAppVersions() []string {
	app.avLock.RLock()
	defer app.avLock.RUnlock()

	return app.cachedAppVersions
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

	//add default account
	accountID, _ := idgen.NewUUID()
	defaultAccount := &model.Account{ID: accountID.String(), ExternalID: externalID, Default: true}

	user, err := app.storage.CreateUser(nil, externalID, uuid, publicKey, consent, exposureNotification, rePost, encryptedKey, encryptedBlob, defaultAccount)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//CreateAdminAppUser creates an admin app user
func (app *Application) CreateAdminAppUser(shibboAuth *model.ShibbolethAuth) (*model.User, error) {
	externalID := "a_" + shibboAuth.Uin //TODO
	user, err := app.storage.CreateUser(shibboAuth, externalID, "", "", false, false, false, nil, nil, nil)
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

//LoadAllRosters loads all rosters
func (app *Application) LoadAllRosters() ([]map[string]string, error) {
	rosters, err := app.storage.ReadAllRosters()
	if err != nil {
		return nil, err
	}
	return rosters, nil
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
	avLock := &sync.RWMutex{}
	listeners := []ApplicationListener{}

	application := Application{version: version, build: build, dataProvider: dataProvider, sender: sender, messaging: messaging,
		profileBB: profileBB, storage: storage, audit: audit, cvLock: cvLock, avLock: avLock, listeners: listeners}

	//add the drivers ports/interfaces
	application.Services = &servicesImpl{app: &application}
	application.Administration = &administrationImpl{app: &application}

	return &application
}
