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
	"bytes"
	"fmt"
	"health/core/model"
	"health/utils"
	"time"
)

//Services exposes APIs for the driver adapters
type Services interface {
	GetVersion() string

	ClearUserData(current model.User) error

	GetUserByShibbolethUIN(shibbolethUIN string) (*model.User, error)
	GetUsersForRePost() ([]*model.User, error)

	GetResources() ([]*model.Resource, error)

	GetFAQ() (*model.FAQ, error)

	GetNews(limit int64) ([]*model.News, error)

	GetEStatusByUserID(userID string, appVersion *string) (*model.EStatus, error)
	CreateOrUpdateEStatus(userID string, appVersion *string, date *time.Time, encryptedKey string, encryptedBlob string) (*model.EStatus, error)
	DeleteEStatus(userID string, appVersion *string) error

	GetEHistoriesByUserID(userID string) ([]*model.EHistory, error)
	CreateЕHistory(userID string, date time.Time, eType string, encryptedKey string, encryptedBlob string) (*model.EHistory, error)
	CreateManualЕHistory(userID string, date time.Time, encryptedKey string, encryptedBlob string, encryptedImageKey *string, encryptedImageBlob *string,
		countyID *string, locationID *string) (*model.EHistory, error)
	DeleteEHitories(userID string) (int64, error)
	UpdateEHistory(userID string, ID string, date *time.Time, encryptedKey *string, encryptedBlob *string) (*model.EHistory, error)

	GetCTests(urrent model.User, processed bool) ([]*model.CTest, []*model.Provider, error)
	CreateExternalCTest(providerID string, uin string, encryptedKey string, encryptedBlob string) error
	DeleteCTests(userID string) (int64, error)
	UpdateCTest(current model.User, ID string, processed bool) (*model.CTest, error)

	GetProviders() ([]*model.Provider, error)

	FindCounties(f *utils.Filter) ([]*model.County, error)
	GetCounty(ID string) (*model.County, error)

	GetRulesByCounty(countyID string) ([]*model.Rule, []*model.CountyStatus, []*model.TestType, error)

	GetLocation(ID string) (*model.Location, error)
	GetLocationsByProviderIDCountyID(providerID string, countyID string) ([]*model.Location, error)
	GetLocationsByCountyID(countyID string) ([]*model.Location, error)
	GetLocationsByCounties(countyIDs []string) ([]*model.Location, error)

	GetAllTestTypes() ([]*model.TestType, error)
	GetTestTypesByIDs(ids []string) ([]*model.TestType, error)

	GetSymptomGroups() ([]*model.SymptomGroup, error)

	GetSymptomRuleByCounty(countyID string) (*model.SymptomRule, []*model.CountyStatus, error)
	GetAccessRuleByCounty(countyID string) (*model.AccessRule, []*model.CountyStatus, error)

	AddTraceReport(items []model.TraceExposure) (int, error)
	GetExposures(timestamp *int64, dateAdded *int64) ([]model.TraceExposure, error)
}

type servicesImpl struct {
	app *Application
}

func (s *servicesImpl) GetVersion() string {
	return s.app.getVersion()
}

func (s *servicesImpl) ClearUserData(current model.User) error {
	return s.app.clearUserData(current)
}

func (s *servicesImpl) GetUserByShibbolethUIN(shibbolethUIN string) (*model.User, error) {
	return s.app.getUserByShibbolethUIN(shibbolethUIN)
}

func (s *servicesImpl) GetUsersForRePost() ([]*model.User, error) {
	return s.app.getUsersForRePost()
}

func (s *servicesImpl) GetResources() ([]*model.Resource, error) {
	return s.app.getResources()
}

func (s *servicesImpl) GetFAQ() (*model.FAQ, error) {
	return s.app.getFAQ()
}

func (s *servicesImpl) GetNews(limit int64) ([]*model.News, error) {
	return s.app.getNews(limit)
}

func (s *servicesImpl) GetEStatusByUserID(userID string, appVersion *string) (*model.EStatus, error) {
	return s.app.getEStatusByUserID(userID, appVersion)
}

func (s *servicesImpl) CreateOrUpdateEStatus(userID string, appVersion *string, date *time.Time, encryptedKey string, encryptedBlob string) (*model.EStatus, error) {
	return s.app.createOrUpdateEStatus(userID, appVersion, date, encryptedKey, encryptedBlob)
}

func (s *servicesImpl) DeleteEStatus(userID string, appVersion *string) error {
	return s.app.deleteEStatus(userID, appVersion)
}

func (s *servicesImpl) GetEHistoriesByUserID(userID string) ([]*model.EHistory, error) {
	return s.app.getEHistoriesByUserID(userID)
}

func (s *servicesImpl) CreateЕHistory(userID string, date time.Time, eType string, encryptedKey string, encryptedBlob string) (*model.EHistory, error) {
	return s.app.createЕHistory(userID, date, eType, encryptedKey, encryptedBlob)
}

func (s *servicesImpl) CreateManualЕHistory(userID string, date time.Time, encryptedKey string, encryptedBlob string, encryptedImageKey *string, encryptedImageBlob *string,
	countyID *string, locationID *string) (*model.EHistory, error) {
	return s.app.createManualЕHistory(userID, date, encryptedKey, encryptedBlob, encryptedImageKey, encryptedImageBlob, countyID, locationID)
}

func (s *servicesImpl) DeleteEHitories(userID string) (int64, error) {
	return s.app.deleteEHitories(userID)
}

func (s *servicesImpl) UpdateEHistory(userID string, ID string, date *time.Time, encryptedKey *string, encryptedBlob *string) (*model.EHistory, error) {
	return s.app.updateEHistory(userID, ID, date, encryptedKey, encryptedBlob)
}

func (s *servicesImpl) GetCTests(current model.User, processed bool) ([]*model.CTest, []*model.Provider, error) {
	return s.app.getCTests(current, processed)
}

func (s *servicesImpl) CreateExternalCTest(providerID string, uin string, encryptedKey string, encryptedBlob string) error {
	return s.app.createExternalCTest(providerID, uin, encryptedKey, encryptedBlob)
}

func (s *servicesImpl) DeleteCTests(userID string) (int64, error) {
	return s.app.deleteCTests(userID)
}

func (s *servicesImpl) UpdateCTest(current model.User, ID string, processed bool) (*model.CTest, error) {
	return s.app.updateCTest(current, ID, processed)
}

func (s *servicesImpl) GetProviders() ([]*model.Provider, error) {
	return s.app.getProviders()
}

func (s *servicesImpl) FindCounties(f *utils.Filter) ([]*model.County, error) {
	return s.app.findCounties(f)
}

func (s *servicesImpl) GetCounty(ID string) (*model.County, error) {
	return s.app.getCounty(ID)
}

func (s *servicesImpl) GetRulesByCounty(countyID string) ([]*model.Rule, []*model.CountyStatus, []*model.TestType, error) {
	return s.app.getRulesByCounty(countyID)
}

func (s *servicesImpl) GetLocation(ID string) (*model.Location, error) {
	return s.app.getLocation(ID)
}

func (s *servicesImpl) GetLocationsByProviderIDCountyID(providerID string, countyID string) ([]*model.Location, error) {
	return s.app.getLocationsByProviderIDCountyID(providerID, countyID)
}

func (s *servicesImpl) GetLocationsByCountyID(countyID string) ([]*model.Location, error) {
	return s.app.getLocationsByCountyID(countyID)
}

func (s *servicesImpl) GetLocationsByCounties(countyIDs []string) ([]*model.Location, error) {
	return s.app.getLocationsByCounties(countyIDs)
}

func (s *servicesImpl) GetAllTestTypes() ([]*model.TestType, error) {
	return s.app.getAllTestTypes()
}

func (s *servicesImpl) GetTestTypesByIDs(ids []string) ([]*model.TestType, error) {
	return s.app.getTestTypesByIDs(ids)
}

func (s *servicesImpl) GetSymptomGroups() ([]*model.SymptomGroup, error) {
	return s.app.getSymptomGroups()
}

func (s *servicesImpl) GetSymptomRuleByCounty(countyID string) (*model.SymptomRule, []*model.CountyStatus, error) {
	return s.app.getSymptomRuleByCounty(countyID)
}

func (s *servicesImpl) GetAccessRuleByCounty(countyID string) (*model.AccessRule, []*model.CountyStatus, error) {
	return s.app.getAccessRuleByCounty(countyID)
}

func (s *servicesImpl) AddTraceReport(items []model.TraceExposure) (int, error) {
	return s.app.аddTraceReport(items)
}

func (s *servicesImpl) GetExposures(timestamp *int64, dateAdded *int64) ([]model.TraceExposure, error) {
	return s.app.getExposures(timestamp, dateAdded)
}

//Administration exposes administration APIs for the driver adapters
type Administration interface {
	GetCovid19Config() (*model.COVID19Config, error)
	UpdateCovid19Config(config *model.COVID19Config) error

	GetNews() ([]*model.News, error)
	CreateNews(date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error)
	UpdateNews(ID string, date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error)
	DeleteNews(ID string) error

	GetResources() ([]*model.Resource, error)
	CreateResource(title string, link string, displayOrder int) (*model.Resource, error)
	UpdateResource(ID string, title string, link string, displayOrder int) (*model.Resource, error)
	DeleteResource(ID string) error
	UpdateResourceDisplayOrder(IDs []string) error

	GetFAQs() (*model.FAQ, error)
	CreateFAQ(section string, sectionDisplayOrder int, title string, description string, questionDisplayOrder int) error
	UpdateFAQ(ID string, title string, description string, displayOrder int) error
	DeleteFAQ(ID string) error
	DeleteFAQSection(ID string) error

	UpdateFAQSection(ID string, title string, displayOrder int) error

	GetProviders() ([]*model.Provider, error)
	CreateProvider(providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error)
	UpdateProvider(ID string, providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error)
	DeleteProvider(ID string) error

	FindCounties(f *utils.Filter) ([]*model.County, error)
	CreateCounty(current model.User, name string, stateProvince string, country string) (*model.County, error)
	UpdateCounty(ID string, name string, stateProvince string, country string) (*model.County, error)
	DeleteCounty(current model.User, ID string) error

	CreateGuideline(countyID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error)
	UpdateGuideline(current model.User, ID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error)
	DeleteGuideline(ID string) error
	GetGuidelinesByCountyID(countyID string) ([]*model.Guideline, error)

	CreateCountyStatus(countyID string, name string, description string) (*model.CountyStatus, error)
	UpdateCountyStatus(ID string, name string, description string) (*model.CountyStatus, error)
	DeleteCountyStatus(ID string) error
	GetCountyStatusByCountyID(countyID string) ([]*model.CountyStatus, error)

	GetTestTypes() ([]*model.TestType, error)
	CreateTestType(name string, priority *int) (*model.TestType, error)
	UpdateTestType(ID string, name string, priority *int) (*model.TestType, error)
	DeleteTestType(ID string) error

	CreateTestTypeResult(testTypeID string, name string, nextStep string, nextStepOffset *int, resultExpiresOffset *int) (*model.TestTypeResult, error)
	UpdateTestTypeResult(ID string, name string, nextStep string, nextStepOffset *int, resultExpiresOffset *int) (*model.TestTypeResult, error)
	DeleteTestTypeResult(ID string) error
	GetTestTypeResultsByTestTypeID(testTypeID string) ([]*model.TestTypeResult, error)

	GetRules() ([]*model.Rule, error)
	CreateRule(countyID string, testTypeID string, priority *int, resultsStates []model.TestTypeResultCountyStatus) (*model.Rule, error)
	UpdateRule(ID string, priority *int, resultsStates []model.TestTypeResultCountyStatus) (*model.Rule, error)
	DeleteRule(ID string) error

	GetLocations() ([]*model.Location, error)
	CreateLocation(providerID string, countyID string, name string, address1 string, address2 string, city string,
		state string, zip string, country string, latitude float64, longitude float64, contact string,
		daysOfOperation []model.OperationDay, url string, notes string, availableTests []string) (*model.Location, error)
	UpdateLocation(ID string, name string, address1 string, address2 string, city string,
		state string, zip string, country string, latitude float64, longitude float64, contact string,
		daysOfOperation []model.OperationDay, url string, notes string, availableTests []string) (*model.Location, error)
	DeleteLocation(ID string) error

	CreateSymptom(Name string, SymptomGroup string) (*model.Symptom, error)
	UpdateSymptom(ID string, name string) (*model.Symptom, error)
	DeleteSymptom(ID string) error

	GetSymptomGroups() ([]*model.SymptomGroup, error)

	GetSymptomRules() ([]*model.SymptomRule, error)
	CreateSymptomRule(countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error)
	UpdateSymptomRule(ID string, countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error)
	DeleteSymptomRule(ID string) error

	GetManualTestByCountyID(countyID string, status *string) ([]*model.EManualTest, error)
	ProcessManualTest(ID string, status string, encryptedKey *string, encryptedBlob *string) error
	GetManualTestImage(ID string) (*string, *string, error)

	GetAccessRules() ([]*model.AccessRule, error)
	CreateAccessRule(countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error)
	UpdateAccessRule(ID string, countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error)
	DeleteAccessRule(ID string) error

	GetUserByExternalID(externalID string) (*model.User, error)

	CreateAction(providerID string, userID string, encryptedKey string, encryptedBlob string) (*model.CTest, error)
}

type administrationImpl struct {
	app *Application
}

func (s *administrationImpl) GetCovid19Config() (*model.COVID19Config, error) {
	return s.app.getCovid19Config()
}

func (s *administrationImpl) UpdateCovid19Config(config *model.COVID19Config) error {
	return s.app.updateCovid19Config(config)
}

func (s *administrationImpl) GetNews() ([]*model.News, error) {
	return s.app.getAllNews()
}

func (s *administrationImpl) CreateNews(date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error) {
	return s.app.createNews(date, title, description, htmlContent, link)
}

func (s *administrationImpl) UpdateNews(ID string, date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error) {
	return s.app.updateNews(ID, date, title, description, htmlContent, nil)
}

func (s *administrationImpl) DeleteNews(ID string) error {
	return s.app.deleteNews(ID)
}

func (s *administrationImpl) GetResources() ([]*model.Resource, error) {
	return s.app.getAllResources()
}

func (s *administrationImpl) CreateResource(title string, link string, displayOrder int) (*model.Resource, error) {
	return s.app.createResource(title, link, displayOrder)
}

func (s *administrationImpl) UpdateResource(ID string, title string, link string, displayOrder int) (*model.Resource, error) {
	return s.app.updateResource(ID, title, link, displayOrder)
}

func (s *administrationImpl) DeleteResource(ID string) error {
	return s.app.deleteResource(ID)
}

func (s *administrationImpl) UpdateResourceDisplayOrder(IDs []string) error {
	return s.app.updateResourceDisplayOrder(IDs)
}

func (s *administrationImpl) GetFAQs() (*model.FAQ, error) {
	return s.app.getFAQs()
}

func (s *administrationImpl) CreateFAQ(section string, sectionDisplayOrder int, title string, description string, questionDisplayOrder int) error {
	return s.app.createFAQ(section, sectionDisplayOrder, title, description, questionDisplayOrder)
}

func (s *administrationImpl) UpdateFAQ(ID string, title string, description string, displayOrder int) error {
	return s.app.updateFAQ(ID, title, description, displayOrder)
}

func (s *administrationImpl) DeleteFAQ(ID string) error {
	return s.app.deleteFAQ(ID)
}

func (s *administrationImpl) DeleteFAQSection(ID string) error {
	return s.app.deleteFAQSection(ID)
}

func (s *administrationImpl) UpdateFAQSection(ID string, title string, displayOrder int) error {
	return s.app.updateFAQSection(ID, title, displayOrder)
}

func (s *administrationImpl) GetProviders() ([]*model.Provider, error) {
	return s.app.getProviders()
}

func (s *administrationImpl) CreateProvider(providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error) {
	return s.app.createProvider(providerName, manualTest, availableMechanisms)
}

func (s *administrationImpl) UpdateProvider(ID string, providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error) {
	return s.app.updateProvider(ID, providerName, manualTest, availableMechanisms)
}

func (s *administrationImpl) DeleteProvider(ID string) error {
	return s.app.deleteProvider(ID)
}

func (s *administrationImpl) FindCounties(f *utils.Filter) ([]*model.County, error) {
	return s.app.findCounties(f)
}

func (s *administrationImpl) CreateCounty(current model.User, name string, stateProvince string, country string) (*model.County, error) {
	return s.app.createCounty(current, name, stateProvince, country)
}

func (s *administrationImpl) UpdateCounty(ID string, name string, stateProvince string, country string) (*model.County, error) {
	return s.app.updateCounty(ID, name, stateProvince, country)
}

func (s *administrationImpl) DeleteCounty(current model.User, ID string) error {
	return s.app.deleteCounty(current, ID)
}

func (s *administrationImpl) CreateGuideline(countyID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error) {
	return s.app.createGuideline(countyID, name, description, items)
}

func (s *administrationImpl) UpdateGuideline(current model.User, ID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error) {
	return s.app.updateGuideline(current, ID, name, description, items)
}

func (s *administrationImpl) DeleteGuideline(ID string) error {
	return s.app.deleteGuideline(ID)
}

func (s *administrationImpl) GetGuidelinesByCountyID(countyID string) ([]*model.Guideline, error) {
	return s.app.getGuidelinesByCountyID(countyID)
}

func (s *administrationImpl) CreateCountyStatus(countyID string, name string, description string) (*model.CountyStatus, error) {
	return s.app.createCountyStatus(countyID, name, description)
}

func (s *administrationImpl) UpdateCountyStatus(ID string, name string, description string) (*model.CountyStatus, error) {
	return s.app.updateCountyStatus(ID, name, description)
}

func (s *administrationImpl) DeleteCountyStatus(ID string) error {
	return s.app.deleteCountyStatus(ID)
}

func (s *administrationImpl) GetCountyStatusByCountyID(countyID string) ([]*model.CountyStatus, error) {
	return s.app.getCountyStatusByCountyID(countyID)
}

func (s *administrationImpl) GetTestTypes() ([]*model.TestType, error) {
	return s.app.getTestTypes()
}

func (s *administrationImpl) CreateTestType(name string, priority *int) (*model.TestType, error) {
	return s.app.createTestType(name, priority)
}

func (s *administrationImpl) UpdateTestType(ID string, name string, priority *int) (*model.TestType, error) {
	return s.app.updateTestType(ID, name, priority)
}

func (s *administrationImpl) DeleteTestType(ID string) error {
	return s.app.deleteTestType(ID)
}

func (s *administrationImpl) CreateTestTypeResult(testTypeID string, name string, nextStep string, nextStepOffset *int, resultExpiresOffset *int) (*model.TestTypeResult, error) {
	return s.app.createTestTypeResult(testTypeID, name, nextStep, nextStepOffset, resultExpiresOffset)
}

func (s *administrationImpl) UpdateTestTypeResult(ID string, name string, nextStep string, nextStepOffset *int, resultExpiresOffset *int) (*model.TestTypeResult, error) {
	return s.app.updateTestTypeResult(ID, name, nextStep, nextStepOffset, resultExpiresOffset)
}

func (s *administrationImpl) DeleteTestTypeResult(ID string) error {
	return s.app.deleteTestTypeResult(ID)
}

func (s *administrationImpl) GetTestTypeResultsByTestTypeID(testTypeID string) ([]*model.TestTypeResult, error) {
	return s.app.getTestTypeResultsByTestTypeID(testTypeID)
}

func (s *administrationImpl) GetRules() ([]*model.Rule, error) {
	return s.app.getRules()
}

func (s *administrationImpl) CreateRule(countyID string, testTypeID string, priority *int, resultsStatuses []model.TestTypeResultCountyStatus) (*model.Rule, error) {
	return s.app.createRule(countyID, testTypeID, priority, resultsStatuses)
}

func (s *administrationImpl) UpdateRule(ID string, priority *int, resultsStates []model.TestTypeResultCountyStatus) (*model.Rule, error) {
	return s.app.updateRule(ID, priority, resultsStates)
}

func (s *administrationImpl) DeleteRule(ID string) error {
	return s.app.deleteRule(ID)
}

func (s *administrationImpl) GetLocations() ([]*model.Location, error) {
	return s.app.getLocations()
}

func (s *administrationImpl) CreateLocation(providerID string, countyID string, name string, address1 string, address2 string, city string,
	state string, zip string, country string, latitude float64, longitude float64, contact string,
	daysOfOperation []model.OperationDay, url string, notes string, availableTests []string) (*model.Location, error) {
	return s.app.createLocation(providerID, countyID, name, address1, address2, city, state, zip, country,
		latitude, longitude, contact, daysOfOperation, url, notes, availableTests)
}

func (s *administrationImpl) UpdateLocation(ID string, name string, address1 string, address2 string, city string,
	state string, zip string, country string, latitude float64, longitude float64, contact string,
	daysOfOperation []model.OperationDay, url string, notes string, availableTests []string) (*model.Location, error) {
	return s.app.updateLocation(ID, name, address1, address2, city, state, zip, country,
		latitude, longitude, contact, daysOfOperation, url, notes, availableTests)
}

func (s *administrationImpl) DeleteLocation(ID string) error {
	return s.app.deleteLocation(ID)
}

func (s *administrationImpl) CreateSymptom(name string, symptomGroup string) (*model.Symptom, error) {
	return s.app.createSymptom(name, symptomGroup)
}

func (s *administrationImpl) UpdateSymptom(ID string, name string) (*model.Symptom, error) {
	return s.app.updateSymptom(ID, name)
}

func (s *administrationImpl) DeleteSymptom(ID string) error {
	return s.app.deleteSymptom(ID)
}

func (s *administrationImpl) GetSymptomGroups() ([]*model.SymptomGroup, error) {
	return s.app.getSymptomGroups()
}

func (s *administrationImpl) GetSymptomRules() ([]*model.SymptomRule, error) {
	return s.app.getSymptomRules()
}

func (s *administrationImpl) CreateSymptomRule(countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error) {
	return s.app.createSymptomRule(countyID, gr1Count, gr2Count, items)
}

func (s *administrationImpl) UpdateSymptomRule(ID string, countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error) {
	return s.app.updateSymptomRule(ID, countyID, gr1Count, gr2Count, items)
}

func (s *administrationImpl) DeleteSymptomRule(ID string) error {
	return s.app.deleteSymptomRule(ID)
}

func (s *administrationImpl) GetManualTestByCountyID(countyID string, status *string) ([]*model.EManualTest, error) {
	return s.app.getManualTestByCountyID(countyID, status)
}

func (s *administrationImpl) ProcessManualTest(ID string, status string, encryptedKey *string, encryptedBlob *string) error {
	return s.app.processManualTest(ID, status, encryptedKey, encryptedBlob)
}

func (s *administrationImpl) GetManualTestImage(ID string) (*string, *string, error) {
	return s.app.getManualTestImage(ID)
}

func (s *administrationImpl) GetAccessRules() ([]*model.AccessRule, error) {
	return s.app.getAccessRules()
}

func (s *administrationImpl) CreateAccessRule(countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error) {
	return s.app.createAccessRule(countyID, rules)
}

func (s *administrationImpl) UpdateAccessRule(ID string, countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error) {
	return s.app.updateAccessRule(ID, countyID, rules)
}

func (s *administrationImpl) DeleteAccessRule(ID string) error {
	return s.app.deleteAccessRule(ID)
}

func (s *administrationImpl) GetUserByExternalID(externalID string) (*model.User, error) {
	return s.app.getUserByExternalID(externalID)
}

func (s *administrationImpl) CreateAction(providerID string, userID string, encryptedKey string, encryptedBlob string) (*model.CTest, error) {
	return s.app.createAction(providerID, userID, encryptedKey, encryptedBlob)
}

//Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	SetStorageListener(storageListener StorageListener)

	ClearUserData(userID string) error
	FindUser(userID string) (*model.User, error)
	FindUserByExternalID(externalID string) (*model.User, error)
	FindUserByShibbolethID(shibbolethID string) (*model.User, error)
	FindUsersByRePost(rePost bool) ([]*model.User, error)
	CreateUser(shibboAuth *model.ShibbolethAuth, externalID string,
		uuid string, publicKey string, consent bool, exposureNotification bool, rePost bool, encryptedKey *string, encryptedBlob *string) (*model.User, error)
	SaveUser(user *model.User) error

	ReadCovid19Config() (*model.COVID19Config, error)
	SaveCovid19Config(covid19Config *model.COVID19Config) error

	ReadAllResources() ([]*model.Resource, error)
	CreateResource(title string, link string, displayOrder int) (*model.Resource, error)
	DeleteResource(ID string) error
	FindResource(ID string) (*model.Resource, error)
	SaveResource(resource *model.Resource) error

	ReadFAQ() (*model.FAQ, error)
	SaveFAQ(faq *model.FAQ) error
	DeleteFAQSection(ID string) error

	ReadNews(limit int64) ([]*model.News, error)
	CreateNews(date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error)
	DeleteNews(ID string) error
	FindNews(ID string) (*model.News, error)
	SaveNews(news *model.News) error

	CreateEStatus(appVersion *string, userID string, date *time.Time, encryptedKey string, encryptedBlob string) (*model.EStatus, error)
	FindEStatusByUserID(appVersion *string, userID string) (*model.EStatus, error)
	SaveEStatus(status *model.EStatus) error
	DeleteEStatus(appVersion *string, userID string) error

	CreateEHistory(userID string, date time.Time, eType string, encryptedKey string, encryptedBlob string) (*model.EHistory, error)
	CreateManualЕHistory(userID string, date time.Time, encryptedKey string, encryptedBlob string, encryptedImageKey *string, encryptedImageBlob *string,
		countyID *string, locationID *string) (*model.EHistory, error)
	FindEHistories(userID string) ([]*model.EHistory, error)
	DeleteEHistories(userID string) (int64, error)
	FindEHistory(ID string) (*model.EHistory, error)
	SaveEHistory(history *model.EHistory) error

	ReadAllProviders() ([]*model.Provider, error)
	CreateProvider(providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error)
	FindProvider(ID string) (*model.Provider, error)
	SaveProvider(provider *model.Provider) error
	DeleteProvider(ID string) error

	CreateExternalCTest(providerID string, uin string, encryptedKey string, encryptedBlob string, processed bool) (*model.CTest, *model.User, error)
	CreateAdminCTest(providerID string, userID string, encryptedKey string, encryptedBlob string, processed bool) (*model.CTest, *model.User, error)
	FindCTest(ID string) (*model.CTest, error)
	FindCTests(userID string, processed bool) ([]*model.CTest, error)
	DeleteCTests(userID string) (int64, error)
	SaveCTest(ctest *model.CTest) error

	FindCounties(f *utils.Filter) ([]*model.County, error)
	CreateCounty(name string, stateProvince string, country string) (*model.County, error)
	FindCounty(ID string) (*model.County, error)
	SaveCounty(county *model.County) error
	DeleteCounty(ID string) error

	CreateGuideline(countyID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error)
	FindGuideline(ID string) (*model.Guideline, error)
	FindGuidelineByCountyID(countyID string) ([]*model.Guideline, error)
	SaveGuideline(guideline *model.Guideline) error
	DeleteGuideline(ID string) error

	CreateCountyStatus(countyID string, name string, description string) (*model.CountyStatus, error)
	FindCountyStatus(ID string) (*model.CountyStatus, error)
	FindCountyStatusesByCountyID(countyID string) ([]*model.CountyStatus, error)
	SaveCountyStatus(countyStatus *model.CountyStatus) error
	DeleteCountyStatus(ID string) error

	ReadAllTestTypes() ([]*model.TestType, error)
	CreateTestType(name string, priority *int) (*model.TestType, error)
	FindTestType(ID string) (*model.TestType, error)
	FindTestTypesByIDs(ids []string) ([]*model.TestType, error)
	SaveTestType(testType *model.TestType) error
	DeleteTestType(ID string) error

	CreateTestTypeResult(testTypeID string, name string, nextStep string, nextStepOffset *int, resultExpiresOffset *int) (*model.TestTypeResult, error)
	FindTestTypeResult(ID string) (*model.TestTypeResult, error)
	FindTestTypeResultsByTestTypeID(testTypeID string) ([]*model.TestTypeResult, error)
	SaveTestTypeResult(testTypeResult *model.TestTypeResult) error
	DeleteTestTypeResult(ID string) error

	ReadAllRules() ([]*model.Rule, error)
	FindRulesByCountyID(countyID string) ([]*model.Rule, error)
	FindRule(ID string) (*model.Rule, error)
	FindRuleByCountyIDTestTypeID(countyID string, testTypeID string) (*model.Rule, error)
	CreateRule(countyID string, testTypeID string, priority *int, resultsStates []model.TestTypeResultCountyStatus) (*model.Rule, error)
	SaveRule(rule *model.Rule) error
	DeleteRule(ID string) error

	ReadAllLocations() ([]*model.Location, error)
	CreateLocation(providerID string, countyID string, name string, address1 string, address2 string, city string,
		state string, zip string, country string, latitude float64, longitude float64, contact string,
		daysOfOperation []model.OperationDay, url string, notes string, availableTests []string) (*model.Location, error)
	FindLocationsByProviderIDCountyID(providerID string, countyID string) ([]*model.Location, error)
	FindLocationsByCountyIDDeep(countyID string) ([]*model.Location, error)
	FindLocationsByCountiesDeep(countyIDs []string) ([]*model.Location, error)
	FindLocation(ID string) (*model.Location, error)
	SaveLocation(location *model.Location) error
	DeleteLocation(ID string) error

	FindSymptom(ID string) (*model.Symptom, error)
	CreateSymptom(name string, symptomGroup string) (*model.Symptom, error)
	DeleteSymptom(ID string) error
	SaveSymptom(symptom *model.Symptom) error

	ReadAllSymptomGroups() ([]*model.SymptomGroup, error)

	ReadAllSymptomRules() ([]*model.SymptomRule, error)
	CreateSymptomRule(countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error)
	FindSymptomRule(ID string) (*model.SymptomRule, error)
	FindSymptomRuleByCountyID(countyID string) (*model.SymptomRule, error)
	SaveSymptomRule(symptomRule *model.SymptomRule) error
	DeleteSymptomRule(ID string) error

	CreateTraceReports(items []model.TraceExposure) (int, error)
	ReadTraceExposures(timestamp *int64, dateAdded *int64) ([]model.TraceExposure, error)

	FindManualTestsByCountyIDDeep(countyID string, status *string) ([]*model.EManualTest, error)
	FindManualTestImage(ID string) (*string, *string, error)
	ProcessManualTest(ID string, status string, encryptedKey *string, encryptedBlob *string) error

	ReadAllAccessRules() ([]*model.AccessRule, error)
	CreateAccessRule(countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error)
	UpdateAccessRule(ID string, countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error)
	FindAccessRuleByCountyID(countyID string) (*model.AccessRule, error)
	DeleteAccessRule(ID string) error
}

//StorageListener listenes for change data storage events
type StorageListener interface {
	OnConfigsChanged()
}

type storageListenerImpl struct {
	app *Application
}

func (a *storageListenerImpl) OnConfigsChanged() {
	//reload the configs
	a.app.loadCovid19Config()
}

//DataProvider is used by core to access needed data
type DataProvider interface {
	LoadNews() ([]ProviderNews, error)
	LoadResources() ([]ProviderResource, error)
}

//ProviderNews represents data provider news entity
type ProviderNews struct {
	PubDate        time.Time
	Title          string
	Description    string
	ContentEncoded string
}

//ProviderResource represents data provider resource entity
type ProviderResource struct {
	Title string
	Link  string
}

//Sender is used by core to send emails
type Sender interface {
	SendForNews(newsList []*model.News)
	SendForResources(resourcesList []*model.Resource)
}

//Messaging is used by core to send user messages
type Messaging interface {
	SendNotificationMessage(tokens []string, title string, body string, data map[string]string)
}

//ProfileBuildingBlock is used by core to communicate with the profile building block.
type ProfileBuildingBlock interface {
	LoadUserData(uuid string) (*ProfileUserData, error)
}

//ProfileUserData represents the profile building block user data entity
type ProfileUserData struct {
	FCMTokens []string `json:"fcmTokens"`
}

//Audit is used by core to log history
type Audit interface {
	LogCreateEvent(userIdentifier string, userInfo string, userGroups []string, entity string, entityID string, data AuditData)
	LogUpdateEvent(userIdentifier string, userInfo string, userGroups []string, entity string, entityID string, data map[string]interface{})
	LogDeleteEvent(userIdentifier string, userInfo string, userGroups []string, entity string, entityID string)
	//TODO add params
	Find() ([]AuditEntity, error)
}

//AuditEntity represents audit module entity
type AuditEntity struct {
	UserIdentifier string    `json:"user_identifier" bson:"user_identifier"`
	UserInfo       string    `json:"user_info" bson:"user_info"`
	UserGroups     []string  `json:"user_groups" bson:"user_groups"`
	Entity         string    `json:"entity" bson:"entity"`
	EntityID       string    `json:"entity_id" bson:"entity_id"`
	Operation      string    `json:"operation" bson:"operation"`
	Data           *string   `json:"data" bson:"data"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
}

//AuditData represents audit data
type AuditData struct {
	data  map[string]interface{}
	order []string
}

//GetData gets the data in string
func (ad AuditData) GetData() *string {
	if len(ad.data) <= 0 || len(ad.order) <= 0 {
		return nil
	}
	if len(ad.data) != len(ad.order) {
		value := "bad data"
		return &value
	}

	var b bytes.Buffer
	i := 0
	count := len(ad.order)
	for _, key := range ad.order {
		value := ad.data[key]
		res := fmt.Sprintf("%s:%s", key, value)
		b.WriteString(res)

		if i < (count - 1) {
			b.WriteString(", ")
		}
		i++
	}
	dataFormatted := b.String()
	return &dataFormatted
}

func NewAuditData(data map[string]interface{}, order []string) AuditData {
	return AuditData{data: data, order: order}
}

//ApplicationListener represents application listener
type ApplicationListener interface {
	OnClearUserData(user model.User)
	OnUserUpdated(user model.User)
}
