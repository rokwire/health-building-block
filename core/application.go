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
		log.Printf("... -> %s is OPEN\n", location.Name)
		//TODO
	} else {
		//TODO
		log.Printf("... -> %s is CLOSED\n", location.Name)
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

func (app *Application) test() {
	log.Println("Application -> checkLocationsWaitTimesColors")

	// load locations
	locations, err := app.storage.ReadAllLocations()
	if err != nil {
		log.Printf("error loading locations for wait time color check - %s", err)
	}

	for _, loc := range locations {
		log.Printf("%s - %s", loc.Name, loc.DaysOfOperation)
	}

	//	ref, _ := time.Parse("03:04 PM", "12:00 AM")
	//	t, err := time.Parse("03:04 PM", "11:22 PM")
	//	if err != nil {
	//		panic(err)
	//	}

	//fmt.Println(t.Sub(ref).Seconds())

	t, err := time.Parse("03:04pm", "06:30pm")
	log.Printf("Parsed -> t - hours:%d minutes:%d seconds:%d\n", t.Hour(), t.Minute(), t.Second())

	//TODO
	/*
		location, err := time.LoadLocation("America/Chicago")
		if err != nil {
			log.Printf("Error getting location:%s\n", err.Error())
		}
		now := time.Now().In(location)
		//log.Printf("setupColorChangeTimer -> now - hours:%d minutes:%d seconds:%d\n", now.Hour(), now.Minute(), now.Second())

		//now := time.Now()
		log.Printf("NOW: %s", now.Weekday())

		log.Printf("A de: %s", now.Format("03:04:05PM"))

		log.Printf("NOWWWW -> now - hours:%d minutes:%d seconds:%d\n", now.Hour(), now.Minute(), now.Second())

		currentMomentInSec := (now.Hour() * 60 * 60) + (now.Minute() * 60) + now.Second()

		log.Printf("NOWWWW in sec %d \n", currentMomentInSec)

	*/
	/*	layout1 := "03:04:05PM"
		layout2 := "15:04:05"
		t, err := time.Parse(layout1, "07:05:45PM")
		if err != nil {
		    fmt.Println(err)
		    return
		}
		fmt.Println(t.Format(layout1))
		fmt.Println(t.Format(layout2)) */
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
