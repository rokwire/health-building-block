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
	GetUINsByOrderNumbers(orderNumbers []string) (map[string]*string, error)
	GetCTestsByExternalUserIDs(externalUserIDs []string) (map[string][]*model.CTest, error)

	GetResources() ([]*model.Resource, error)

	GetFAQ() (*model.FAQ, error)

	GetNews(limit int64) ([]*model.News, error)

	GetEStatusByAccountID(accountID string, appVersion *string) (*model.EStatus, error)
	CreateOrUpdateEStatus(accountID string, appVersion *string, date *time.Time, encryptedKey string, encryptedBlob string) (*model.EStatus, error)
	DeleteEStatus(accountID string, appVersion *string) error

	GetEHistoriesByAccountID(accountID string) ([]*model.EHistory, error)
	CreateЕHistory(accountID string, date time.Time, eType string, encryptedKey string, encryptedBlob string) (*model.EHistory, error)
	CreateManualЕHistory(accountID string, date time.Time, encryptedKey string, encryptedBlob string, encryptedImageKey *string, encryptedImageBlob *string,
		countyID *string, locationID *string) (*model.EHistory, error)
	DeleteEHitories(accountID string) (int64, error)
	UpdateEHistory(accountID string, ID string, date *time.Time, encryptedKey *string, encryptedBlob *string) (*model.EHistory, error)

	GetCTests(account model.Account, processed bool) ([]*model.CTest, []*model.Provider, error)
	CreateExternalCTest(providerID string, uin string, encryptedKey string, encryptedBlob string, orderNumber *string) error
	DeleteCTests(accountID string) (int64, error)
	UpdateCTest(account model.Account, ID string, processed bool) (*model.CTest, error)

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
	GetSymptoms(appVersion *string) (*model.Symptoms, error)

	GetSymptomRuleByCounty(countyID string) (*model.SymptomRule, []*model.CountyStatus, error)
	GetCRulesByCounty(appVersion *string, countyID string) (*model.CRules, error)
	GetAccessRuleByCounty(countyID string) (*model.AccessRule, []*model.CountyStatus, error)

	AddTraceReport(items []model.TraceExposure) (int, error)
	GetExposures(timestamp *int64, dateAdded *int64) ([]model.TraceExposure, error)

	GetUINOverride(account model.Account) (*model.UINOverride, error)
	CreateOrUpdateUINOverride(account model.Account, interval int, category *string, expiration *time.Time) error

	GetExtUINOverrides(uin *string, sort *string) ([]*model.UINOverride, error)
	CreateExtUINOverride(uin string, interval int, category *string, expiration *time.Time) (*model.UINOverride, error)
	UpdateExtUINOverride(uin string, interval int, category *string, expiration *time.Time) (*string, error)
	DeleteExtUINOverride(uin string) error

	SetUINBuildingAccess(account model.Account, date time.Time, access string) error
	GetExtUINBuildingAccess(uin string) (*model.UINBuildingAccess, error)

	GetRosterByPhone(phone string) (map[string]string, error)

	GetExtJoinExternalApproval(account model.Account) ([]RokmetroJoinGroupExtApprovement, error)
	UpdateExtJoinExternalApprovement(jeaID string, status string) error
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

func (s *servicesImpl) GetUINsByOrderNumbers(orderNumbers []string) (map[string]*string, error) {
	return s.app.getUINsByOrderNumbers(orderNumbers)
}

func (s *servicesImpl) GetCTestsByExternalUserIDs(externalUserIDs []string) (map[string][]*model.CTest, error) {
	return s.app.getCTestsByExternalUserIDs(externalUserIDs)
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

func (s *servicesImpl) GetEStatusByAccountID(accountID string, appVersion *string) (*model.EStatus, error) {
	return s.app.getEStatusByAccountID(accountID, appVersion)
}

func (s *servicesImpl) CreateOrUpdateEStatus(accountID string, appVersion *string, date *time.Time, encryptedKey string, encryptedBlob string) (*model.EStatus, error) {
	return s.app.createOrUpdateEStatus(accountID, appVersion, date, encryptedKey, encryptedBlob)
}

func (s *servicesImpl) DeleteEStatus(accountID string, appVersion *string) error {
	return s.app.deleteEStatus(accountID, appVersion)
}

func (s *servicesImpl) GetEHistoriesByAccountID(accountID string) ([]*model.EHistory, error) {
	return s.app.getEHistoriesByAccountID(accountID)
}

func (s *servicesImpl) CreateЕHistory(accountID string, date time.Time, eType string, encryptedKey string, encryptedBlob string) (*model.EHistory, error) {
	return s.app.createЕHistory(accountID, date, eType, encryptedKey, encryptedBlob)
}

func (s *servicesImpl) CreateManualЕHistory(accountID string, date time.Time, encryptedKey string, encryptedBlob string, encryptedImageKey *string, encryptedImageBlob *string,
	countyID *string, locationID *string) (*model.EHistory, error) {
	return s.app.createManualЕHistory(accountID, date, encryptedKey, encryptedBlob, encryptedImageKey, encryptedImageBlob, countyID, locationID)
}

func (s *servicesImpl) DeleteEHitories(accountID string) (int64, error) {
	return s.app.deleteEHitories(accountID)
}

func (s *servicesImpl) UpdateEHistory(accountID string, ID string, date *time.Time, encryptedKey *string, encryptedBlob *string) (*model.EHistory, error) {
	return s.app.updateEHistory(accountID, ID, date, encryptedKey, encryptedBlob)
}

func (s *servicesImpl) GetCTests(account model.Account, processed bool) ([]*model.CTest, []*model.Provider, error) {
	return s.app.getCTests(account, processed)
}

func (s *servicesImpl) CreateExternalCTest(providerID string, uin string, encryptedKey string, encryptedBlob string, orderNumber *string) error {
	return s.app.createExternalCTest(providerID, uin, encryptedKey, encryptedBlob, orderNumber)
}

func (s *servicesImpl) DeleteCTests(accountID string) (int64, error) {
	return s.app.deleteCTests(accountID)
}

func (s *servicesImpl) UpdateCTest(account model.Account, ID string, processed bool) (*model.CTest, error) {
	return s.app.updateCTest(account, ID, processed)
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

func (s *servicesImpl) GetSymptoms(appVersion *string) (*model.Symptoms, error) {
	return s.app.getSymptoms(appVersion)
}

func (s *servicesImpl) GetSymptomRuleByCounty(countyID string) (*model.SymptomRule, []*model.CountyStatus, error) {
	return s.app.getSymptomRuleByCounty(countyID)
}

func (s *servicesImpl) GetCRulesByCounty(appVersion *string, countyID string) (*model.CRules, error) {
	return s.app.getCRulesByCounty(appVersion, countyID)
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

func (s *servicesImpl) GetUINOverride(account model.Account) (*model.UINOverride, error) {
	return s.app.getUINOverride(account)
}

func (s *servicesImpl) CreateOrUpdateUINOverride(account model.Account, interval int, category *string, expiration *time.Time) error {
	return s.app.createOrUpdateUINOverride(account, interval, category, expiration)
}

func (s *servicesImpl) GetExtUINOverrides(uin *string, sort *string) ([]*model.UINOverride, error) {
	return s.app.getExtUINOverrides(uin, sort)
}

func (s *servicesImpl) CreateExtUINOverride(uin string, interval int, category *string, expiration *time.Time) (*model.UINOverride, error) {
	return s.app.createExtUINOverride(uin, interval, category, expiration)
}

func (s *servicesImpl) UpdateExtUINOverride(uin string, interval int, category *string, expiration *time.Time) (*string, error) {
	return s.app.updateExtUINOverride(uin, interval, category, expiration)
}

func (s *servicesImpl) DeleteExtUINOverride(uin string) error {
	return s.app.deleteExtUINOverride(uin)
}

func (s *servicesImpl) SetUINBuildingAccess(account model.Account, date time.Time, access string) error {
	return s.app.setUINBuildingAccess(account, date, access)
}

func (s *servicesImpl) GetExtUINBuildingAccess(uin string) (*model.UINBuildingAccess, error) {
	return s.app.getExtUINBuildingAccess(uin)
}

func (s *servicesImpl) GetRosterByPhone(phone string) (map[string]string, error) {
	return s.app.getRosterByPhone(phone)
}

func (s *servicesImpl) GetExtJoinExternalApproval(account model.Account) ([]RokmetroJoinGroupExtApprovement, error) {
	return s.app.getExtJoinExternalApproval(account)
}

func (s *servicesImpl) UpdateExtJoinExternalApprovement(jeaID string, status string) error {
	return s.app.updateExtJoinExternalApprovement(jeaID, status)
}

//Administration exposes administration APIs for the driver adapters
type Administration interface {
	GetCovid19Config() (*model.COVID19Config, error)
	UpdateCovid19Config(config *model.COVID19Config) error
	GetCovid19Configs() ([]model.COVID19Config, error)

	GetAppVersions() ([]string, error)
	CreateAppVersion(current model.User, group string, audit *string, version string) error

	GetNews() ([]*model.News, error)
	CreateNews(current model.User, group string, audit *string, date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error)
	UpdateNews(current model.User, group string, audit *string, ID string, date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error)
	DeleteNews(current model.User, group string, ID string) error

	GetResources() ([]*model.Resource, error)
	CreateResource(current model.User, group string, audit *string, title string, link string, displayOrder int) (*model.Resource, error)
	UpdateResource(current model.User, group string, audit *string, ID string, title string, link string, displayOrder int) (*model.Resource, error)
	DeleteResource(current model.User, group string, ID string) error
	UpdateResourceDisplayOrder(IDs []string) error

	GetFAQs() (*model.FAQ, error)
	CreateFAQ(current model.User, group string, audit *string, section string, sectionDisplayOrder int, title string, description string, questionDisplayOrder int) error
	UpdateFAQ(current model.User, group string, audit *string, ID string, title string, description string, displayOrder int) error
	DeleteFAQ(current model.User, group string, ID string) error

	DeleteFAQSection(current model.User, group string, ID string) error
	UpdateFAQSection(current model.User, group string, audit *string, ID string, title string, displayOrder int) error

	GetProviders() ([]*model.Provider, error)
	CreateProvider(current model.User, group string, audit *string, providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error)
	UpdateProvider(current model.User, group string, audit *string, ID string, providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error)
	DeleteProvider(current model.User, group string, ID string) error

	FindCounties(f *utils.Filter) ([]*model.County, error)
	CreateCounty(current model.User, group string, audit *string, name string, stateProvince string, country string) (*model.County, error)
	UpdateCounty(current model.User, group string, audit *string, ID string, name string, stateProvince string, country string) (*model.County, error)
	DeleteCounty(current model.User, group string, ID string) error

	CreateGuideline(current model.User, group string, audit *string, countyID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error)
	UpdateGuideline(current model.User, group string, audit *string, ID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error)
	DeleteGuideline(current model.User, group string, ID string) error
	GetGuidelinesByCountyID(countyID string) ([]*model.Guideline, error)

	CreateCountyStatus(current model.User, group string, audit *string, countyID string, name string, description string) (*model.CountyStatus, error)
	UpdateCountyStatus(current model.User, group string, audit *string, ID string, name string, description string) (*model.CountyStatus, error)
	DeleteCountyStatus(current model.User, group string, ID string) error
	GetCountyStatusByCountyID(countyID string) ([]*model.CountyStatus, error)

	GetTestTypes() ([]*model.TestType, error)
	CreateTestType(current model.User, group string, audit *string, name string, priority *int) (*model.TestType, error)
	UpdateTestType(current model.User, group string, audit *string, ID string, name string, priority *int) (*model.TestType, error)
	DeleteTestType(current model.User, group string, ID string) error

	CreateTestTypeResult(current model.User, group string, audit *string, testTypeID string, name string, nextStep string, nextStepOffset *int, resultExpiresOffset *int) (*model.TestTypeResult, error)
	UpdateTestTypeResult(current model.User, group string, audit *string, ID string, name string, nextStep string, nextStepOffset *int, resultExpiresOffset *int) (*model.TestTypeResult, error)
	DeleteTestTypeResult(current model.User, group string, ID string) error
	GetTestTypeResultsByTestTypeID(testTypeID string) ([]*model.TestTypeResult, error)

	GetRules() ([]*model.Rule, error)
	CreateRule(current model.User, group string, audit *string, countyID string, testTypeID string, priority *int, resultsStates []model.TestTypeResultCountyStatus) (*model.Rule, error)
	UpdateRule(current model.User, group string, audit *string, ID string, priority *int, resultsStates []model.TestTypeResultCountyStatus) (*model.Rule, error)
	DeleteRule(current model.User, group string, ID string) error

	GetLocations() ([]*model.Location, error)
	CreateLocation(current model.User, group string, audit *string, providerID string, countyID string, name string, address1 string, address2 string, city string,
		state string, zip string, country string, latitude float64, longitude float64, contact string,
		daysOfOperation []model.OperationDay, url string, notes string, waitTimeColor *string, availableTests []string) (*model.Location, error)
	UpdateLocation(current model.User, group string, audit *string, ID string, name string, address1 string, address2 string, city string,
		state string, zip string, country string, latitude float64, longitude float64, contact string,
		daysOfOperation []model.OperationDay, url string, notes string, waitTimeColor *string, availableTests []string) (*model.Location, error)
	DeleteLocation(current model.User, group string, ID string) error

	CreateSymptom(current model.User, group string, Name string, SymptomGroup string) (*model.Symptom, error)
	UpdateSymptom(current model.User, group string, ID string, name string) (*model.Symptom, error)
	DeleteSymptom(current model.User, group string, ID string) error

	GetSymptomGroups() ([]*model.SymptomGroup, error)

	GetSymptomRules() ([]*model.SymptomRule, error)
	CreateSymptomRule(current model.User, group string, countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error)
	UpdateSymptomRule(current model.User, group string, ID string, countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error)
	DeleteSymptomRule(current model.User, group string, ID string) error

	GetManualTestByCountyID(countyID string, status *string) ([]*model.EManualTest, error)
	ProcessManualTest(ID string, status string, encryptedKey *string, encryptedBlob *string, date *time.Time) error
	GetManualTestImage(ID string) (*string, *string, error)

	GetAccessRules() ([]*model.AccessRule, error)
	CreateAccessRule(current model.User, group string, audit *string, countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error)
	UpdateAccessRule(current model.User, group string, audit *string, ID string, countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error)
	DeleteAccessRule(current model.User, group string, ID string) error

	GetCRules(countyID string, appVersion string) (*model.CRules, error)
	CreateOrUpdateCRules(current model.User, group string, audit *string, countyID string, appVersion string, data string) error

	GetSymptoms(appVersion string) (*model.Symptoms, error)
	CreateOrUpdateSymptoms(current model.User, group string, audit *string, appVersion string, items string) error

	GetUINOverrides(uin *string, sort *string) ([]*model.UINOverride, error)
	CreateUINOverride(current model.User, group string, audit *string, uin string, interval int, category *string, expiration *time.Time) (*model.UINOverride, error)
	UpdateUINOverride(current model.User, group string, audit *string, uin string, interval int, category *string, expiration *time.Time) (*string, error)
	DeleteUINOverride(current model.User, group string, uin string) error

	CreateRoster(current model.User, group string, audit *string, phone string, uin string, firstName string,
		middleName string, lastName string, birthDate string, gender string, address1 string, address2 string,
		address3 string, city string, state string, zipCode string, email string, badgeType string) error
	UpdateRoster(current model.User, group string, audit *string, uin string, firstName string, middleName string, lastName string, birthDate string, gender string,
		address1 string, address2 string, address3 string, city string, state string, zipCode string, email string, badgeType string) error
	CreateRosterItems(current model.User, group string, audit *string, items []map[string]string) error
	GetRosters(filter *utils.Filter, sortBy string, sortOrder int, limit int, offset int) ([]map[string]interface{}, error)
	DeleteRosterByPhone(current model.User, group string, phone string) error
	DeleteRosterByUIN(current model.User, group string, uin string) error
	DeleteAllRosters(current model.User, group string) error

	CreateRawSubAccountItems(current model.User, group string, audit *string, items []model.RawSubAccount) error
	GetRawSubAccounts(filter *utils.Filter, sortBy string, sortOrder int, limit int, offset int) ([]model.RawSubAccount, error)
	UpdateRawSubAccount(current model.User, group string, audit *string, uin string, firstName string, middleName string, lastName string, birthDate string, gender string,
		address1 string, address2 string, address3 string, city string, state string, zipCode string, phone string, netID string, email string) error
	DeleteRawSubAccountByUIN(current model.User, group string, uin string) error
	DeleteAllRawSubAccounts(current model.User, group string) error

	GetUserByExternalID(externalID string) (*model.User, error)

	CreateAction(current model.User, group string, audit *string, providerID string, accountID string, encryptedKey string, encryptedBlob string) (*model.CTest, error)

	GetAudit(current model.User, group string, userIdentifier *string, entity *string, entityID *string, operation *string, clientData *string,
		createdAt *time.Time, sortBy *string, asc *bool, limit *int64) ([]*AuditEntity, error)
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

func (s *administrationImpl) GetCovid19Configs() ([]model.COVID19Config, error) {
	return s.app.getCovid19Configs()
}

func (s *administrationImpl) GetNews() ([]*model.News, error) {
	return s.app.getAllNews()
}

func (s *administrationImpl) GetAppVersions() ([]string, error) {
	return s.app.getAppVersions()
}

func (s *administrationImpl) CreateAppVersion(current model.User, group string, audit *string, version string) error {
	return s.app.createAppVersion(current, group, audit, version)
}

func (s *administrationImpl) CreateNews(current model.User, group string, audit *string, date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error) {
	return s.app.createNews(current, group, audit, date, title, description, htmlContent, link)
}

func (s *administrationImpl) UpdateNews(current model.User, group string, audit *string, ID string, date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error) {
	return s.app.updateNews(current, group, audit, ID, date, title, description, htmlContent, nil)
}

func (s *administrationImpl) DeleteNews(current model.User, group string, ID string) error {
	return s.app.deleteNews(current, group, ID)
}

func (s *administrationImpl) GetResources() ([]*model.Resource, error) {
	return s.app.getAllResources()
}

func (s *administrationImpl) CreateResource(current model.User, group string, audit *string, title string, link string, displayOrder int) (*model.Resource, error) {
	return s.app.createResource(current, group, audit, title, link, displayOrder)
}

func (s *administrationImpl) UpdateResource(current model.User, group string, audit *string, ID string, title string, link string, displayOrder int) (*model.Resource, error) {
	return s.app.updateResource(current, group, audit, ID, title, link, displayOrder)
}

func (s *administrationImpl) DeleteResource(current model.User, group string, ID string) error {
	return s.app.deleteResource(current, group, ID)
}

func (s *administrationImpl) UpdateResourceDisplayOrder(IDs []string) error {
	return s.app.updateResourceDisplayOrder(IDs)
}

func (s *administrationImpl) GetFAQs() (*model.FAQ, error) {
	return s.app.getFAQs()
}

func (s *administrationImpl) CreateFAQ(current model.User, group string, audit *string, section string, sectionDisplayOrder int, title string, description string, questionDisplayOrder int) error {
	return s.app.createFAQ(current, group, audit, section, sectionDisplayOrder, title, description, questionDisplayOrder)
}

func (s *administrationImpl) UpdateFAQ(current model.User, group string, audit *string, ID string, title string, description string, displayOrder int) error {
	return s.app.updateFAQ(current, group, audit, ID, title, description, displayOrder)
}

func (s *administrationImpl) DeleteFAQ(current model.User, group string, ID string) error {
	return s.app.deleteFAQ(current, group, ID)
}

func (s *administrationImpl) DeleteFAQSection(current model.User, group string, ID string) error {
	return s.app.deleteFAQSection(current, group, ID)
}

func (s *administrationImpl) UpdateFAQSection(current model.User, group string, audit *string, ID string, title string, displayOrder int) error {
	return s.app.updateFAQSection(current, group, audit, ID, title, displayOrder)
}

func (s *administrationImpl) GetProviders() ([]*model.Provider, error) {
	return s.app.getProviders()
}

func (s *administrationImpl) CreateProvider(current model.User, group string, audit *string, providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error) {
	return s.app.createProvider(current, group, audit, providerName, manualTest, availableMechanisms)
}

func (s *administrationImpl) UpdateProvider(current model.User, group string, audit *string, ID string, providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error) {
	return s.app.updateProvider(current, group, audit, ID, providerName, manualTest, availableMechanisms)
}

func (s *administrationImpl) DeleteProvider(current model.User, group string, ID string) error {
	return s.app.deleteProvider(current, group, ID)
}

func (s *administrationImpl) FindCounties(f *utils.Filter) ([]*model.County, error) {
	return s.app.findCounties(f)
}

func (s *administrationImpl) CreateCounty(current model.User, group string, audit *string, name string, stateProvince string, country string) (*model.County, error) {
	return s.app.createCounty(current, group, audit, name, stateProvince, country)
}

func (s *administrationImpl) UpdateCounty(current model.User, group string, audit *string, ID string, name string, stateProvince string, country string) (*model.County, error) {
	return s.app.updateCounty(current, group, audit, ID, name, stateProvince, country)
}

func (s *administrationImpl) DeleteCounty(current model.User, group string, ID string) error {
	return s.app.deleteCounty(current, group, ID)
}

func (s *administrationImpl) CreateGuideline(current model.User, group string, audit *string, countyID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error) {
	return s.app.createGuideline(current, group, audit, countyID, name, description, items)
}

func (s *administrationImpl) UpdateGuideline(current model.User, group string, audit *string, ID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error) {
	return s.app.updateGuideline(current, group, audit, ID, name, description, items)
}

func (s *administrationImpl) DeleteGuideline(current model.User, group string, ID string) error {
	return s.app.deleteGuideline(current, group, ID)
}

func (s *administrationImpl) GetGuidelinesByCountyID(countyID string) ([]*model.Guideline, error) {
	return s.app.getGuidelinesByCountyID(countyID)
}

func (s *administrationImpl) CreateCountyStatus(current model.User, group string, audit *string, countyID string, name string, description string) (*model.CountyStatus, error) {
	return s.app.createCountyStatus(current, group, audit, countyID, name, description)
}

func (s *administrationImpl) UpdateCountyStatus(current model.User, group string, audit *string, ID string, name string, description string) (*model.CountyStatus, error) {
	return s.app.updateCountyStatus(current, group, audit, ID, name, description)
}

func (s *administrationImpl) DeleteCountyStatus(current model.User, group string, ID string) error {
	return s.app.deleteCountyStatus(current, group, ID)
}

func (s *administrationImpl) GetCountyStatusByCountyID(countyID string) ([]*model.CountyStatus, error) {
	return s.app.getCountyStatusByCountyID(countyID)
}

func (s *administrationImpl) GetTestTypes() ([]*model.TestType, error) {
	return s.app.getTestTypes()
}

func (s *administrationImpl) CreateTestType(current model.User, group string, audit *string, name string, priority *int) (*model.TestType, error) {
	return s.app.createTestType(current, group, audit, name, priority)
}

func (s *administrationImpl) UpdateTestType(current model.User, group string, audit *string, ID string, name string, priority *int) (*model.TestType, error) {
	return s.app.updateTestType(current, group, audit, ID, name, priority)
}

func (s *administrationImpl) DeleteTestType(current model.User, group string, ID string) error {
	return s.app.deleteTestType(current, group, ID)
}

func (s *administrationImpl) CreateTestTypeResult(current model.User, group string, audit *string, testTypeID string, name string, nextStep string, nextStepOffset *int, resultExpiresOffset *int) (*model.TestTypeResult, error) {
	return s.app.createTestTypeResult(current, group, audit, testTypeID, name, nextStep, nextStepOffset, resultExpiresOffset)
}

func (s *administrationImpl) UpdateTestTypeResult(current model.User, group string, audit *string, ID string, name string, nextStep string, nextStepOffset *int, resultExpiresOffset *int) (*model.TestTypeResult, error) {
	return s.app.updateTestTypeResult(current, group, audit, ID, name, nextStep, nextStepOffset, resultExpiresOffset)
}

func (s *administrationImpl) DeleteTestTypeResult(current model.User, group string, ID string) error {
	return s.app.deleteTestTypeResult(current, group, ID)
}

func (s *administrationImpl) GetTestTypeResultsByTestTypeID(testTypeID string) ([]*model.TestTypeResult, error) {
	return s.app.getTestTypeResultsByTestTypeID(testTypeID)
}

func (s *administrationImpl) GetRules() ([]*model.Rule, error) {
	return s.app.getRules()
}

func (s *administrationImpl) CreateRule(current model.User, group string, audit *string, countyID string, testTypeID string, priority *int, resultsStatuses []model.TestTypeResultCountyStatus) (*model.Rule, error) {
	return s.app.createRule(current, group, audit, countyID, testTypeID, priority, resultsStatuses)
}

func (s *administrationImpl) UpdateRule(current model.User, group string, audit *string, ID string, priority *int, resultsStates []model.TestTypeResultCountyStatus) (*model.Rule, error) {
	return s.app.updateRule(current, group, audit, ID, priority, resultsStates)
}

func (s *administrationImpl) DeleteRule(current model.User, group string, ID string) error {
	return s.app.deleteRule(current, group, ID)
}

func (s *administrationImpl) GetLocations() ([]*model.Location, error) {
	return s.app.getLocations()
}

func (s *administrationImpl) CreateLocation(current model.User, group string, audit *string, providerID string, countyID string, name string, address1 string, address2 string, city string,
	state string, zip string, country string, latitude float64, longitude float64, contact string,
	daysOfOperation []model.OperationDay, url string, notes string, waitTimeColor *string, availableTests []string) (*model.Location, error) {
	return s.app.createLocation(current, group, audit, providerID, countyID, name, address1, address2, city, state, zip, country,
		latitude, longitude, contact, daysOfOperation, url, notes, waitTimeColor, availableTests)
}

func (s *administrationImpl) UpdateLocation(current model.User, group string, audit *string, ID string, name string, address1 string, address2 string, city string,
	state string, zip string, country string, latitude float64, longitude float64, contact string,
	daysOfOperation []model.OperationDay, url string, notes string, waitTimeColor *string, availableTests []string) (*model.Location, error) {
	return s.app.updateLocation(current, group, audit, ID, name, address1, address2, city, state, zip, country,
		latitude, longitude, contact, daysOfOperation, url, notes, waitTimeColor, availableTests)
}

func (s *administrationImpl) DeleteLocation(current model.User, group string, ID string) error {
	return s.app.deleteLocation(current, group, ID)
}

func (s *administrationImpl) CreateSymptom(current model.User, group string, name string, symptomGroup string) (*model.Symptom, error) {
	return s.app.createSymptom(current, group, name, symptomGroup)
}

func (s *administrationImpl) UpdateSymptom(current model.User, group string, ID string, name string) (*model.Symptom, error) {
	return s.app.updateSymptom(current, group, ID, name)
}

func (s *administrationImpl) DeleteSymptom(current model.User, group string, ID string) error {
	return s.app.deleteSymptom(current, group, ID)
}

func (s *administrationImpl) GetSymptomGroups() ([]*model.SymptomGroup, error) {
	return s.app.getSymptomGroups()
}

func (s *administrationImpl) GetSymptomRules() ([]*model.SymptomRule, error) {
	return s.app.getSymptomRules()
}

func (s *administrationImpl) CreateSymptomRule(current model.User, group string, countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error) {
	return s.app.createSymptomRule(current, group, countyID, gr1Count, gr2Count, items)
}

func (s *administrationImpl) UpdateSymptomRule(current model.User, group string, ID string, countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error) {
	return s.app.updateSymptomRule(current, group, ID, countyID, gr1Count, gr2Count, items)
}

func (s *administrationImpl) DeleteSymptomRule(current model.User, group string, ID string) error {
	return s.app.deleteSymptomRule(current, group, ID)
}

func (s *administrationImpl) GetManualTestByCountyID(countyID string, status *string) ([]*model.EManualTest, error) {
	return s.app.getManualTestByCountyID(countyID, status)
}

func (s *administrationImpl) ProcessManualTest(ID string, status string, encryptedKey *string, encryptedBlob *string, date *time.Time) error {
	return s.app.processManualTest(ID, status, encryptedKey, encryptedBlob, date)
}

func (s *administrationImpl) GetManualTestImage(ID string) (*string, *string, error) {
	return s.app.getManualTestImage(ID)
}

func (s *administrationImpl) GetAccessRules() ([]*model.AccessRule, error) {
	return s.app.getAccessRules()
}

func (s *administrationImpl) CreateAccessRule(current model.User, group string, audit *string, countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error) {
	return s.app.createAccessRule(current, group, audit, countyID, rules)
}

func (s *administrationImpl) UpdateAccessRule(current model.User, group string, audit *string, ID string, countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error) {
	return s.app.updateAccessRule(current, group, audit, ID, countyID, rules)
}

func (s *administrationImpl) DeleteAccessRule(current model.User, group string, ID string) error {
	return s.app.deleteAccessRule(current, group, ID)
}

func (s *administrationImpl) GetCRules(countyID string, appVersion string) (*model.CRules, error) {
	return s.app.getCRules(countyID, appVersion)
}

func (s *administrationImpl) CreateOrUpdateCRules(current model.User, group string, audit *string, countyID string, appVersion string, data string) error {
	return s.app.createOrUpdateCRules(current, group, audit, countyID, appVersion, data)
}

func (s *administrationImpl) GetSymptoms(appVersion string) (*model.Symptoms, error) {
	return s.app.getASymptoms(appVersion)
}

func (s *administrationImpl) CreateOrUpdateSymptoms(current model.User, group string, audit *string, appVersion string, items string) error {
	return s.app.createOrUpdateSymptoms(current, group, audit, appVersion, items)
}

func (s *administrationImpl) GetUINOverrides(uin *string, sort *string) ([]*model.UINOverride, error) {
	return s.app.getUINOverrides(uin, sort)
}

func (s *administrationImpl) CreateUINOverride(current model.User, group string, audit *string, uin string, interval int, category *string, expiration *time.Time) (*model.UINOverride, error) {
	return s.app.createUINOverride(current, group, audit, uin, interval, category, expiration)
}

func (s *administrationImpl) UpdateUINOverride(current model.User, group string, audit *string, uin string, interval int, category *string, expiration *time.Time) (*string, error) {
	return s.app.updateUINOverride(current, group, audit, uin, interval, category, expiration)
}

func (s *administrationImpl) DeleteUINOverride(current model.User, group string, uin string) error {
	return s.app.deleteUINOverride(current, group, uin)
}

func (s *administrationImpl) GetUserByExternalID(externalID string) (*model.User, error) {
	return s.app.getUserByExternalID(externalID)
}

func (s *administrationImpl) CreateRoster(current model.User, group string, audit *string, phone string, uin string, firstName string,
	middleName string, lastName string, birthDate string, gender string, address1 string, address2 string,
	address3 string, city string, state string, zipCode string, email string, badgeType string) error {
	return s.app.createRoster(current, group, audit, phone, uin, firstName, middleName, lastName, birthDate, gender,
		address1, address2, address3, city, state, zipCode, email, badgeType)
}

func (s *administrationImpl) UpdateRoster(current model.User, group string, audit *string, uin string, firstName string, middleName string, lastName string, birthDate string, gender string,
	address1 string, address2 string, address3 string, city string, state string, zipCode string, email string, badgeType string) error {
	return s.app.updateRoster(current, group, audit, uin, firstName, middleName, lastName, birthDate, gender,
		address1, address2, address3, city, state, zipCode, email, badgeType)
}

func (s *administrationImpl) CreateRosterItems(current model.User, group string, audit *string, items []map[string]string) error {
	return s.app.createRosterItems(current, group, audit, items)
}

func (s *administrationImpl) GetRosters(filter *utils.Filter, sortBy string, sortOrder int, limit int, offset int) ([]map[string]interface{}, error) {
	return s.app.getRosters(filter, sortBy, sortOrder, limit, offset)
}

func (s *administrationImpl) DeleteRosterByPhone(current model.User, group string, phone string) error {
	return s.app.deleteRosterByPhone(current, group, phone)
}

func (s *administrationImpl) DeleteRosterByUIN(current model.User, group string, uin string) error {
	return s.app.deleteRosterByUIN(current, group, uin)
}

func (s *administrationImpl) DeleteAllRosters(current model.User, group string) error {
	return s.app.deleteAllRosters(current, group)
}

func (s *administrationImpl) CreateRawSubAccountItems(current model.User, group string, audit *string, items []model.RawSubAccount) error {
	return s.app.createRawSubAccountItems(current, group, audit, items)
}

func (s *administrationImpl) GetRawSubAccounts(filter *utils.Filter, sortBy string, sortOrder int, limit int, offset int) ([]model.RawSubAccount, error) {
	return s.app.getRawSubAccounts(filter, sortBy, sortOrder, limit, offset)
}

func (s *administrationImpl) UpdateRawSubAccount(current model.User, group string, audit *string, uin string, firstName string, middleName string, lastName string, birthDate string, gender string,
	address1 string, address2 string, address3 string, city string, state string, zipCode string, phone string, netID string, email string) error {
	return s.app.updateRawSubAccount(current, group, audit, uin, firstName, middleName, lastName, birthDate, gender, address1, address2,
		address3, city, state, zipCode, phone, netID, email)
}

func (s *administrationImpl) DeleteRawSubAccountByUIN(current model.User, group string, uin string) error {
	return s.app.deleteRawSubAccountByUIN(current, group, uin)
}

func (s *administrationImpl) DeleteAllRawSubAccounts(current model.User, group string) error {
	return s.app.deleteAllRawSubAccounts(current, group)
}

func (s *administrationImpl) CreateAction(current model.User, group string, audit *string, providerID string, accountID string, encryptedKey string, encryptedBlob string) (*model.CTest, error) {
	return s.app.createAction(current, group, audit, providerID, accountID, encryptedKey, encryptedBlob)
}

func (s *administrationImpl) GetAudit(current model.User, group string, userIdentifier *string, entity *string, entityID *string, operation *string,
	clientData *string, createdAt *time.Time, sortBy *string, asc *bool, limit *int64) ([]*AuditEntity, error) {
	return s.app.getAudit(current, group, userIdentifier, entity, entityID, operation, clientData, createdAt, sortBy, asc, limit)
}

//Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	SetStorageListener(storageListener StorageListener)

	ReadAllAppVersions() ([]string, error)
	CreateAppVersion(version string) error

	ClearUserData(userID string) error
	FindUser(userID string) (*model.User, error)
	//it does not look the accounts, it looks just the primary user
	FindUserByExternalID(externalID string) (*model.User, error)
	//it looks primary user + accounts
	FindUserAccountsByExternalID(externalID string) (*model.User, error)
	FindUserByShibbolethID(shibbolethID string) (*model.User, error)
	FindUsersByRePost(rePost bool) ([]*model.User, error)
	CreateAppUser(externalID string,
		uuid string, publicKey string, consent bool, exposureNotification bool, rePost bool, encryptedKey *string, encryptedBlob *string, encryptedPK *string) (*model.User, error)
	CreateAdminUser(shibboAuth *model.ShibbolethAuth, externalID string,
		uuid string, publicKey string, consent bool, exposureNotification bool, rePost bool, encryptedKey *string, encryptedBlob *string) (*model.User, error)
	CreateDefaultAccount(userID string) (*model.User, error)
	SaveUser(user *model.User) error

	GetCovid19Configs() ([]model.COVID19Config, error)
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

	CreateEStatus(appVersion *string, accountID string, date *time.Time, encryptedKey string, encryptedBlob string) (*model.EStatus, error)
	FindEStatusByAccountID(appVersion *string, accountID string) (*model.EStatus, error)
	SaveEStatus(status *model.EStatus) error
	DeleteEStatus(appVersion *string, accountID string) error

	CreateEHistory(accountID string, date time.Time, eType string, encryptedKey string, encryptedBlob string) (*model.EHistory, error)
	CreateManualЕHistory(accountID string, date time.Time, encryptedKey string, encryptedBlob string, encryptedImageKey *string, encryptedImageBlob *string,
		countyID *string, locationID *string) (*model.EHistory, error)
	FindEHistories(accountID string) ([]*model.EHistory, error)
	DeleteEHistories(accountD string) (int64, error)
	FindEHistory(ID string) (*model.EHistory, error)
	SaveEHistory(history *model.EHistory) error

	ReadAllProviders() ([]*model.Provider, error)
	CreateProvider(providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error)
	FindProvider(ID string) (*model.Provider, error)
	SaveProvider(provider *model.Provider) error
	DeleteProvider(ID string) error

	CreateExternalCTest(providerID string, uin string, encryptedKey string, encryptedBlob string, processed bool, orderNumber *string) (*model.CTest, *model.User, error)
	CreateAdminCTest(providerID string, accountID string, encryptedKey string, encryptedBlob string, processed bool, orderNumber *string) (*model.CTest, *model.User, error)
	FindCTest(ID string) (*model.CTest, error)
	FindCTests(accountID string, processed bool) ([]*model.CTest, error)
	FindCTestsByExternalUserIDs(externalUserIDs []string) (map[string][]*model.CTest, error)
	DeleteCTests(accountID string) (int64, error)
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
		daysOfOperation []model.OperationDay, url string, notes string, waitTimeColor *string, availableTests []string) (*model.Location, error)
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

	ReadSymptoms(appVersion string) (*model.Symptoms, error)
	CreateOrUpdateSymptoms(appVersion string, items string) (*bool, error)

	ReadAllSymptomRules() ([]*model.SymptomRule, error)
	CreateSymptomRule(countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error)
	FindSymptomRule(ID string) (*model.SymptomRule, error)
	FindSymptomRuleByCountyID(countyID string) (*model.SymptomRule, error)
	SaveSymptomRule(symptomRule *model.SymptomRule) error
	DeleteSymptomRule(ID string) error

	FindCRulesByCountyID(appVersion string, countyID string) (*model.CRules, error)
	CreateOrUpdateCRules(appVersion string, countyID string, data string) (*bool, error)

	CreateTraceReports(items []model.TraceExposure) (int, error)
	ReadTraceExposures(timestamp *int64, dateAdded *int64) ([]model.TraceExposure, error)

	FindManualTestsByCountyIDDeep(countyID string, status *string) ([]*model.EManualTest, error)
	FindManualTestImage(ID string) (*string, *string, error)
	ProcessManualTest(ID string, status string, encryptedKey *string, encryptedBlob *string, date *time.Time) error

	ReadAllAccessRules() ([]*model.AccessRule, error)
	CreateAccessRule(countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error)
	UpdateAccessRule(ID string, countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error)
	FindAccessRuleByCountyID(countyID string) (*model.AccessRule, error)
	DeleteAccessRule(ID string) error

	FindExternalUserIDsByTestsOrderNumbers(orderNumbers []string) (map[string]*string, error)

	CreateOrUpdateUINOverride(uin string, interval int, category *string, expiration *time.Time) error
	//finds the uin override for the provided uin. It makes additional check for the expiration because of the mongoDB TTL delay
	FindUINOverride(uin string) (*model.UINOverride, error)

	//finds the uin override for the provided uin. If uin is nil then it gives all
	FindUINOverrides(uin *string, sort *string) ([]*model.UINOverride, error)
	CreateUINOverride(uin string, interval int, category *string, expiration *time.Time) (*model.UINOverride, error)
	UpdateUINOverride(uin string, interval int, category *string, expiration *time.Time) (*string, error)
	DeleteUINOverride(uin string) error

	FindUINBuildingAccess(uin string) (*model.UINBuildingAccess, error)
	CreateOrUpdateUINBuildingAccess(uin string, date time.Time, access string) error

	ReadAllRosters() ([]map[string]string, error)
	FindRosterByPhone(phone string) (map[string]string, error)
	FindRosters(filter *utils.Filter, sortBy string, sortOrder int, limit int, offset int) ([]map[string]interface{}, error)
	CreateRoster(phone string, uin string, firstName string, middleName string, lastName string, birthDate string, gender string,
		address1 string, address2 string, address3 string, city string, state string, zipCode string, email string, badgeType string) error
	UpdateRoster(uin string, firstName string, middleName string, lastName string, birthDate string, gender string,
		address1 string, address2 string, address3 string, city string, state string, zipCode string, email string, badgeType string) error
	CreateRosterItems(items []map[string]string) error
	DeleteRosterByPhone(phone string) error
	DeleteRosterByUIN(uin string) error
	DeleteAllRosters() error

	CreateRawSubAccountItems(items []model.RawSubAccount) error
	FindRawSubAccounts(filter *utils.Filter, sortBy string, sortOrder int, limit int, offset int) ([]model.RawSubAccount, error)
	UpdateRawSubAcccount(uin string, firstName string, middleName string, lastName string, birthDate string, gender string,
		address1 string, address2 string, address3 string, city string, state string, zipCode string, phone string, netID string, email string) error
	DeleteRawSubAccountByUIN(uin string) error
	DeleteAllSubAccounts() error
}

//StorageListener listenes for change data storage events
type StorageListener interface {
	OnConfigsChanged()
	OnAppVersionsChanged()
	OnRostersChanged()
	OnRawSubAccountsChanged()
}

type storageListenerImpl struct {
	app *Application
}

func (a *storageListenerImpl) OnConfigsChanged() {
	//reload the configs
	a.app.loadCovid19Config()
}

func (a *storageListenerImpl) OnAppVersionsChanged() {
	//reload the app versions
	a.app.loadAppVersions()
}

func (a *storageListenerImpl) OnRostersChanged() {
	//notify that th rosters has been changed
	a.app.notifyListeners("onRostersUpdated", nil)
}

func (a *storageListenerImpl) OnRawSubAccountsChanged() {
	//notify that the raw sub accounts have been changed
	a.app.notifyListeners("onRawSubAccountsUpdated", nil)
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

//Rokmetro is used by core to communicate with the rokmetro ecosystem
type Rokmetro interface {
	GetExtJoinExternalApproval(externalApproverID string) ([]RokmetroJoinGroupExtApprovement, error)
	UpdateExtJoinExternalApprovement(jeaID string, status string) error
}

//RokmetroJoinGroupExtApprovement represent Rokmetro join group external approvement entity
type RokmetroJoinGroupExtApprovement struct {
	ID                       string    `json:"id"`
	GroupName                string    `json:"group_name"`
	FirstName                string    `json:"first_name"`
	LastName                 string    `json:"last_name"`
	Email                    string    `json:"email"`
	Phone                    string    `json:"phone"`
	DateCreated              time.Time `json:"date_created"`
	ExternalApproverID       string    `json:"external_approver_id"`
	ExternalApproverLastName string    `json:"external_approver_last_name"`
	Status                   string    `json:"status"`
}

//Audit is used by core to log history
type Audit interface {
	LogCreateEvent(userIdentifier string, userInfo string, usedGroup string, entity string, entityID string, data []AuditDataEntry, clientData *string)
	LogUpdateEvent(userIdentifier string, userInfo string, usedGroup string, entity string, entityID string, data []AuditDataEntry, clientData *string)
	LogDeleteEvent(userIdentifier string, userInfo string, usedGroup string, entity string, entityID string)

	Find(userIdentifier *string, usedGroup *string, entity *string, entityID *string, operation *string,
		clientData *string, createdAt *time.Time, sortBy *string, asc *bool, limit *int64) ([]*AuditEntity, error)
}

//AuditEntity represents audit module entity
type AuditEntity struct {
	UserIdentifier string    `json:"user_identifier" bson:"user_identifier"`
	UserInfo       string    `json:"user_info" bson:"user_info"`
	UsedGroup      string    `json:"used_group" bson:"used_group"`
	Entity         string    `json:"entity" bson:"entity"`
	EntityID       string    `json:"entity_id" bson:"entity_id"`
	Operation      string    `json:"operation" bson:"operation"`
	Data           *string   `json:"data" bson:"data"`
	ClientData     *string   `json:"client_data" bson:"client_data"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
} // @name AuditEntity

//AuditDataEntry represents audit data entry
type AuditDataEntry struct {
	Key   string
	Value string
}

//ApplicationListener represents application listener
type ApplicationListener interface {
	OnClearUserData(user model.User)
	OnUserUpdated(user model.User)
	OnRostersUpdated()
	OnSubAccountsUpdated()
}
