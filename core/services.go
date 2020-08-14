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
	"log"
	"time"
)

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) clearUserData(current model.User) error {
	err := app.storage.ClearUserData(current.ID)
	if err != nil {
		return err
	}

	defer app.notifyListeners("onClearUserData", current)

	return nil
}

func (app *Application) getUserByShibbolethUIN(shibbolethUIN string) (*model.User, error) {
	user, err := app.storage.FindUserByExternalID(shibbolethUIN)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (app *Application) getUsersForRePost() ([]*model.User, error) {
	users, err := app.storage.FindUsersByRePost(true)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (app *Application) getResources() ([]*model.Resource, error) {
	resources, err := app.storage.ReadAllResources()
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (app *Application) getFAQ() (*model.FAQ, error) {
	faq, err := app.storage.ReadFAQ()
	if err != nil {
		return nil, err
	}
	return faq, nil
}

func (app *Application) getNews(limit int64) ([]*model.News, error) {
	if limit < 0 {
		return nil, errors.New("cannot pass limit < 0")
	}

	news, err := app.storage.ReadNews(limit)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (app *Application) getCTests(current model.User, processed bool) ([]*model.CTest, []*model.Provider, error) {
	//We get the data with two requests to the database - NoSQL approach in some cases!

	//1. first get the ctest
	ctests, err := app.storage.FindCTests(current.ID, processed)
	if err != nil {
		return nil, nil, err
	}

	//2. we need to get the providers
	providers, err := app.storage.ReadAllProviders()
	if err != nil {
		return nil, nil, err
	}

	return ctests, providers, nil
}

func (app *Application) createExternalCTest(providerID string, uin string, encryptedKey string, encryptedBlob string) error {
	//1. create a ctest
	_, user, err := app.storage.CreateExternalCTest(providerID, uin, encryptedKey, encryptedBlob, false)
	if err != nil {
		return err
	}

	//2. notify that the user is updated
	defer app.notifyListeners("onUserUpdated", *user)

	//3. send a firebase notification to the user that the ctest is arrived.
	go func(userUUID string) {
		if len(userUUID) <= 0 {
			log.Println("user uuid is empty")
			return
		}
		//1. load the user data, we need the fcm tokens
		userData, err := app.profileBB.LoadUserData(user.UUID)
		if err != nil {
			log.Printf("Error loading user data - %s\n", err)
			return
		}

		//2. send notification message
		//2.1 prepare the data
		data := make(map[string]string)
		data["type"] = "health.covid19.notification"
		data["health.covid19.notification.type"] = "process-pending-tests"
		data["title"] = "COVID-19"
		data["body"] = "You have received a COVID-19 update"
		data["click_action"] = "FLUTTER_NOTIFICATION_CLICK"
		app.messaging.SendNotificationMessage(userData.FCMTokens, "COVID-19", "You have received a COVID-19 update", data)
	}(user.UUID)

	return nil
}

func (app *Application) deleteCTests(userID string) (int64, error) {
	deletedCount, err := app.storage.DeleteCTests(userID)
	if err != nil {
		return -1, err
	}
	return deletedCount, nil
}

func (app *Application) updateCTest(current model.User, ID string, processed bool) (*model.CTest, error) {
	//check if we have a such ctest entity
	ctest, err := app.storage.FindCTest(ID)
	if err != nil {
		return nil, err
	}
	if ctest == nil {
		return nil, errors.New("ctest is nil for id " + ID)
	}

	//check if the user owns the ctest
	if current.ID != ctest.UserID {
		return nil, errors.New("the user does not owns this ctest")
	}

	//add the new values
	ctest.Processed = processed

	//save it
	err = app.storage.SaveCTest(ctest)
	if err != nil {
		return nil, err
	}

	return ctest, nil
}

func (app *Application) getRulesByCounty(countyID string) ([]*model.Rule, []*model.CountyStatus, []*model.TestType, error) {
	//1. first check if we have a county for the provided id
	county, err := app.storage.FindCounty(countyID)
	if err != nil {
		return nil, nil, nil, err
	}
	if county == nil {
		return nil, nil, nil, errors.New("there is no a county for the provided id")
	}

	//2. get the rules
	rules, err := app.storage.FindRulesByCountyID(countyID)
	if err != nil {
		return nil, nil, nil, err
	}
	log.Println(rules)

	//3. we need to get the county statuses and the test types also - NoSQL!
	countyStatuses, err := app.storage.FindCountyStatusesByCountyID(countyID)
	if err != nil {
		return nil, nil, nil, err
	}
	testTypes, err := app.storage.ReadAllTestTypes()
	if err != nil {
		return nil, nil, nil, err
	}

	return rules, countyStatuses, testTypes, nil
}

func (app *Application) createOrUpdateEStatus(userID string, appVersion *string, date *time.Time, encryptedKey string, encryptedBlob string) (*model.EStatus, error) {
	//determine if we need to create or update it
	status, err := app.storage.FindEStatusByUserID(appVersion, userID)
	if err != nil {
		return nil, err
	}

	if status == nil {
		//we need to create it
		status, err = app.storage.CreateEStatus(appVersion, userID, date, encryptedKey, encryptedBlob)
		if err != nil {
			return nil, err
		}
		return status, nil
	}
	//we need to update it

	//add the new values
	status.Date = date
	status.EncryptedKey = encryptedKey
	status.EncryptedBlob = encryptedBlob

	//save it
	err = app.storage.SaveEStatus(status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (app *Application) getEStatusByUserID(userID string, appVersion *string) (*model.EStatus, error) {
	status, err := app.storage.FindEStatusByUserID(appVersion, userID)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (app *Application) deleteEStatus(userID string, appVersion *string) error {
	err := app.storage.DeleteEStatus(appVersion, userID)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) getLocation(ID string) (*model.Location, error) {
	location, err := app.storage.FindLocation(ID)
	if err != nil {
		return nil, err
	}
	return location, nil
}

func (app *Application) getLocationsByProviderIDCountyID(providerID string, countyID string) ([]*model.Location, error) {
	locations, err := app.storage.FindLocationsByProviderIDCountyID(providerID, countyID)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (app *Application) getLocationsByCountyID(countyID string) ([]*model.Location, error) {
	locations, err := app.storage.FindLocationsByCountyIDDeep(countyID)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (app *Application) getLocationsByCounties(countyIDs []string) ([]*model.Location, error) {
	locations, err := app.storage.FindLocationsByCountiesDeep(countyIDs)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (app *Application) getAllTestTypes() ([]*model.TestType, error) {
	testTypes, err := app.storage.ReadAllTestTypes()
	if err != nil {
		return nil, err
	}
	return testTypes, nil
}

func (app *Application) getTestTypesByIDs(ids []string) ([]*model.TestType, error) {
	testTypes, err := app.storage.FindTestTypesByIDs(ids)
	if err != nil {
		return nil, err
	}
	return testTypes, nil
}

func (app *Application) getSymptomRuleByCounty(countyID string) (*model.SymptomRule, []*model.CountyStatus, error) {
	//get the symptom rule
	symptomRule, err := app.storage.FindSymptomRuleByCountyID(countyID)
	if err != nil {
		return nil, nil, err
	}
	if symptomRule == nil {
		//no item
		return nil, nil, nil
	}

	//get the county statuses for the provided county
	countyStatuses, err := app.storage.FindCountyStatusesByCountyID(countyID)
	if err != nil {
		return nil, nil, err
	}

	return symptomRule, countyStatuses, nil
}

func (app *Application) getAccessRuleByCounty(countyID string) (*model.AccessRule, []*model.CountyStatus, error) {
	//get the access rule
	accessRule, err := app.storage.FindAccessRuleByCountyID(countyID)
	if err != nil {
		return nil, nil, err
	}
	if accessRule == nil {
		//no item
		return nil, nil, nil
	}

	//get the county statuses for the provided county
	countyStatuses, err := app.storage.FindCountyStatusesByCountyID(countyID)
	if err != nil {
		return nil, nil, err
	}

	return accessRule, countyStatuses, nil
}

func (app *Application) аddTraceReport(items []model.TraceExposure) (int, error) {
	insertedCount, err := app.storage.CreateTraceReports(items)
	if err != nil {
		return 0, err
	}
	return insertedCount, nil
}

func (app *Application) getExposures(timestamp *int64, dateAdded *int64) ([]model.TraceExposure, error) {
	items, err := app.storage.ReadTraceExposures(timestamp, dateAdded)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (app *Application) deleteEHitories(userID string) (int64, error) {
	deletedCount, err := app.storage.DeleteEHistories(userID)
	if err != nil {
		return -1, err
	}
	return deletedCount, nil
}

func (app *Application) updateEHistory(userID string, ID string, date *time.Time, encryptedKey *string, encryptedBlob *string) (*model.EHistory, error) {
	history, err := app.storage.FindEHistory(ID)
	if err != nil {
		return nil, err
	}
	if history == nil {
		return nil, errors.New("history is nil for id " + ID)
	}
	if history.UserID != userID {
		return nil, errors.New("not allowed to modify history with id " + ID)
	}

	//add the new values
	if date != nil {
		history.Date = *date
	}
	if encryptedKey != nil {
		history.EncryptedKey = *encryptedKey
	}
	if encryptedBlob != nil {
		history.EncryptedBlob = *encryptedBlob
	}

	//save it
	err = app.storage.SaveEHistory(history)
	if err != nil {
		return nil, err
	}

	return history, nil
}
