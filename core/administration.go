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
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (app *Application) getCovid19Config() (*model.COVID19Config, error) {
	config, err := app.storage.ReadCovid19Config()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (app *Application) updateCovid19Config(config *model.COVID19Config) error {
	err := app.storage.SaveCovid19Config(config)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) getAppVersions() ([]string, error) {
	appVersions, err := app.storage.ReadAllAppVersions()
	if err != nil {
		return nil, err
	}
	return appVersions, nil
}

func (app *Application) createAppVersion(current model.User, group string, audit *string, version string) error {
	//First validate the version input. We accept x.x.x or x.x which is the short for x.x.0
	//If the input is 3.5.0 then we will store 3.5 as the system works with the short view when the patch is 0
	var major, minor, patch int
	var err error
	elements := strings.Split(version, ".")
	elementsCount := len(elements)
	if !(elementsCount == 2 || elementsCount == 3) {
		return errors.New("format must be x.x.x or x.x")
	}

	major, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}
	minor, err = strconv.Atoi(elements[1])
	if err != nil {
		return err
	}
	if elementsCount == 2 {
		patch = 0
	} else {
		patch, err = strconv.Atoi(elements[2])
		if err != nil {
			return err
		}
	}

	res := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	//the system work with the shor view when patch is 0
	if patch == 0 {
		res = fmt.Sprintf("%d.%d", major, minor)
	}

	err = app.storage.CreateAppVersion(res)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "version", Value: version}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "app-version", version, lData, audit)

	return nil
}

func (app *Application) getAllNews() ([]*model.News, error) {
	news, err := app.storage.ReadNews(0)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (app *Application) createNews(current model.User, group string, audit *string, date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error) {
	news, err := app.storage.CreateNews(date, title, description, htmlContent, link)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "date", Value: fmt.Sprint(date)}, {Key: "title", Value: title}, {Key: "description", Value: description},
		{Key: "htmlContent", Value: htmlContent}, {Key: "link", Value: utils.GetString(link)}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "news", news.ID, lData, audit)

	return news, nil
}

func (app *Application) updateNews(current model.User, group string, audit *string, ID string, date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error) {
	news, err := app.storage.FindNews(ID)
	if err != nil {
		return nil, err
	}
	if news == nil {
		return nil, errors.New("news is nil for id " + ID)
	}

	//add the new values
	news.Date = date
	news.Title = title
	news.Description = description
	news.HTMLContent = htmlContent

	//save it
	err = app.storage.SaveNews(news)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "date", Value: fmt.Sprint(date)}, {Key: "title", Value: title}, {Key: "description", Value: description},
		{Key: "htmlContent", Value: htmlContent}, {Key: "link", Value: utils.GetString(link)}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "news", ID, lData, audit)

	return news, nil
}

func (app *Application) deleteNews(current model.User, group string, ID string) error {
	err := app.storage.DeleteNews(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "news", ID)
	return nil
}

func (app *Application) getAllResources() ([]*model.Resource, error) {
	news, err := app.storage.ReadAllResources()
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (app *Application) createResource(current model.User, group string, audit *string, title string, link string, displayOrder int) (*model.Resource, error) {
	resource, err := app.storage.CreateResource(title, link, displayOrder)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "title", Value: title}, {Key: "link", Value: link}, {Key: "displayOrder", Value: fmt.Sprint(displayOrder)}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "resource", resource.ID, lData, audit)

	return resource, nil
}

func (app *Application) updateResource(current model.User, group string, audit *string, ID string, title string, link string, displayOrder int) (*model.Resource, error) {
	resource, err := app.storage.FindResource(ID)
	if err != nil {
		return nil, err
	}
	if resource == nil {
		return nil, errors.New("resource is nil for id " + ID)
	}

	//add the new values
	resource.Title = title
	resource.Link = link
	resource.DisplayOrder = displayOrder

	//save it
	err = app.storage.SaveResource(resource)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "title", Value: title}, {Key: "link", Value: link}, {Key: "displayOrder", Value: fmt.Sprint(displayOrder)}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "resource", ID, lData, audit)

	return resource, nil
}

func (app *Application) deleteResource(current model.User, group string, ID string) error {
	err := app.storage.DeleteResource(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "resource", ID)

	return nil
}

func (app *Application) updateResourceDisplayOrder(IDs []string) error {
	if IDs == nil {
		return errors.New("updateResourceDisplayOrder() -> no IDs provided")
	}

	//1. validate if the input items are the same size as the stored items
	inputCount := len(IDs)
	allResources, err := app.storage.ReadAllResources()
	if err != nil {
		return err
	}
	itemsCount := len(allResources)
	if inputCount != itemsCount {
		return errors.New("updateResourceDisplayOrder() -> input count != items count")
	}

	//2. update the display order for the items
	//TODO refactor
	for _, resource := range allResources {
		resource.DisplayOrder = app.findDisplayOrder(resource.ID, IDs)
		app.storage.SaveResource(resource)
	}
	return nil
}

func (app *Application) findDisplayOrder(ID string, IDs []string) int {
	for index, current := range IDs {
		if ID == current {
			return index
		}
	}
	return -1
}

func (app *Application) getFAQs() (*model.FAQ, error) {
	faq, err := app.storage.ReadFAQ()
	if err != nil {
		return nil, err
	}
	if faq == nil {
		return nil, errors.New("FAQs is nil")
	}
	return faq, nil
}

func (app *Application) createFAQ(current model.User, group string, audit *string, section string, sectionDisplayOrder int, title string, description string, questionDisplayOrder int) error {
	faq, err := app.storage.ReadFAQ()
	if err != nil {
		return err
	}
	sections := faq.Sections
	if sections == nil {
		return errors.New("for some reasons the sections are nil")
	}

	//create the new question
	qID, _ := uuid.NewUUID()
	question := model.Question{ID: qID.String(), Title: title, Description: description, DisplayOrder: questionDisplayOrder}

	//add the new question for the corresponding section
	foundedSection := app.findSection(section, sections)
	if foundedSection != nil {
		log.Printf("there is a section with name %s\n", section)
		sQuestions := *foundedSection.Questions
		sQuestions = append(sQuestions, &question)

		foundedSection.Questions = &sQuestions
		foundedSection.DisplayOrder = sectionDisplayOrder
	} else {
		log.Printf("there is no a section with name %s\n", section)
		//create a new section
		sQuestions := []*model.Question{&question}
		sID, _ := uuid.NewUUID()
		newSection := model.Section{ID: sID.String(), Title: section, Questions: &sQuestions, DisplayOrder: sectionDisplayOrder}

		//TODO refactor

		s := *sections
		s = append(s, &newSection)

		sections = &s
	}

	faq.Sections = sections
	faq.DateUpdate = time.Now()

	//store the update item
	err = app.storage.SaveFAQ(faq)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "section", Value: section}, {Key: "sectionDisplayOrder", Value: fmt.Sprint(sectionDisplayOrder)}, {Key: "title", Value: title},
		{Key: "description", Value: description}, {Key: "questionDisplayOrder", Value: fmt.Sprint(questionDisplayOrder)}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "faq-question", question.ID, lData, audit)

	return nil
}

func (app *Application) findSection(title string, sections *[]*model.Section) *model.Section {
	for _, section := range *sections {
		if title == section.Title {
			return section
		}
	}
	return nil
}

func (app *Application) updateFAQ(current model.User, group string, audit *string, ID string, title string, description string, displayOrder int) error {
	faq, err := app.storage.ReadFAQ()
	if err != nil {
		return err
	}
	sections := faq.Sections
	if sections == nil {
		return errors.New("updateFAQ -> for some reasons the sections are nil")
	}

	updated := false
	for _, section := range *sections {
		questions := section.Questions
		if questions != nil {
			for _, question := range *questions {
				if question.ID == ID {
					question.Title = title
					question.Description = description
					question.DisplayOrder = displayOrder

					updated = true
				}
			}
		}
	}
	if !updated {
		return errors.New("cannot update faq with the provided id " + ID)
	}

	faq.DateUpdate = time.Now()

	//store the updated item
	err = app.storage.SaveFAQ(faq)
	if err != nil {
		return err
	}

	//title string, description string, displayOrder int

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "title", Value: title}, {Key: "description", Value: description}, {Key: "displayOrder", Value: fmt.Sprint(displayOrder)}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "faq-question", ID, lData, audit)

	return nil
}

func (app *Application) deleteFAQ(current model.User, group string, ID string) error {
	faq, err := app.storage.ReadFAQ()
	if err != nil {
		return err
	}
	sections := faq.Sections
	if sections == nil {
		return errors.New("deleteFAQ -> for some reasons the sections are nil")
	}

	sectionIndex := -1
	questionIndex := -1
	for sIndex, section := range *sections {
		questions := section.Questions
		if questions != nil {
			for qIndex, question := range *questions {
				if question.ID == ID {
					sectionIndex = sIndex
					questionIndex = qIndex
				}
			}
		}
	}
	if sectionIndex == -1 || questionIndex == -1 {
		return errors.New("cannot remove faq with the provided id " + ID)
	}

	section := (*sections)[sectionIndex]
	questions := section.Questions
	modifiedList := remove(*questions, questionIndex)
	section.Questions = &modifiedList

	faq.DateUpdate = time.Now()

	//store the updated item
	err = app.storage.SaveFAQ(faq)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "faq-question", ID)

	return nil
}

func (app *Application) deleteFAQSection(current model.User, group string, ID string) error {
	err := app.storage.DeleteFAQSection(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "faq-section", ID)

	return nil
}

func remove(questions []*model.Question, s int) []*model.Question {
	return append(questions[:s], questions[s+1:]...)
}

func (app *Application) updateFAQSection(current model.User, group string, audit *string, ID string, title string, displayOrder int) error {
	faq, err := app.storage.ReadFAQ()
	if err != nil {
		return err
	}
	sections := faq.Sections
	if sections == nil {
		return errors.New("updateFAQSection -> for some reasons the sections are nil")
	}

	updated := false
	for _, section := range *sections {
		if section.ID == ID {
			section.Title = title
			section.DisplayOrder = displayOrder

			updated = true
		}
	}
	if !updated {
		return errors.New("cannot update faq section with the provided id " + ID)
	}

	faq.DateUpdate = time.Now()

	//store the updated item
	err = app.storage.SaveFAQ(faq)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "title", Value: title}, {Key: "displayOrder", Value: fmt.Sprint(displayOrder)}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "faq-section", ID, lData, audit)

	return nil
}

func (app *Application) createProvider(current model.User, group string, audit *string, providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error) {
	provider, err := app.storage.CreateProvider(providerName, manualTest, availableMechanisms)
	if err != nil {
		return nil, err
	}

	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "providerName", Value: providerName}, {Key: "manualTest", Value: fmt.Sprint(manualTest)}, {Key: "availableMechanisms", Value: fmt.Sprint(availableMechanisms)}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "provider", provider.ID, lData, audit)

	return provider, nil
}

func (app *Application) updateProvider(current model.User, group string, audit *string, ID string, providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error) {
	provider, err := app.storage.FindProvider(ID)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, errors.New("provider is nil for id " + ID)
	}

	//add the new values
	provider.Name = providerName
	provider.ManualTest = manualTest
	provider.AvailableMechanisms = availableMechanisms

	//save it
	err = app.storage.SaveProvider(provider)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "providerName", Value: providerName}, {Key: "manualTest", Value: fmt.Sprint(manualTest)}, {Key: "availableMechanisms", Value: fmt.Sprint(availableMechanisms)}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "provider", ID, lData, audit)

	return provider, nil
}

func (app *Application) deleteProvider(current model.User, group string, ID string) error {
	err := app.storage.DeleteProvider(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "provider", ID)

	return nil
}

func (app *Application) createCounty(current model.User, group string, audit *string, name string, stateProvince string, country string) (*model.County, error) {
	county, err := app.storage.CreateCounty(name, stateProvince, country)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "name", Value: name}, {Key: "stateProvince", Value: stateProvince}, {Key: "country", Value: country}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "county", county.ID, lData, audit)

	return county, nil
}

func (app *Application) updateCounty(current model.User, group string, audit *string, ID string, name string, stateProvince string, country string) (*model.County, error) {
	county, err := app.storage.FindCounty(ID)
	if err != nil {
		return nil, err
	}
	if county == nil {
		return nil, errors.New("county is nil for id " + ID)
	}

	//add the new values
	county.Name = name
	county.StateProvince = stateProvince
	county.Country = country

	//save it
	err = app.storage.SaveCounty(county)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "name", Value: name}, {Key: "stateProvince", Value: stateProvince}, {Key: "country", Value: country}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "county", county.ID, lData, audit)

	return county, nil
}

func (app *Application) deleteCounty(current model.User, group string, ID string) error {
	err := app.storage.DeleteCounty(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "county", ID)

	return nil
}

func (app *Application) createGuideline(current model.User, group string, audit *string, countyID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error) {
	//1. find if we have a county for the provided ID
	county, err := app.storage.FindCounty(countyID)
	if err != nil {
		return nil, err
	}
	if county == nil {
		return nil, errors.New("there is no a county for the provided id")
	}

	//2. create the guideline entity
	guideline, err := app.storage.CreateGuideline(countyID, name, description, items)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "name", Value: name}, {Key: "description", Value: description}, {Key: "items", Value: fmt.Sprint(items)}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "guideline", guideline.ID, lData, audit)

	return guideline, nil
}

func (app *Application) updateGuideline(current model.User, group string, audit *string, ID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error) {
	guideline, err := app.storage.FindGuideline(ID)
	if err != nil {
		return nil, err
	}
	if guideline == nil {
		return nil, errors.New("guideline is nil for id " + ID)
	}

	//add the new values
	guideline.Name = name
	guideline.Description = description
	guideline.Items = items

	//save it
	err = app.storage.SaveGuideline(guideline)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "name", Value: name}, {Key: "description", Value: description}, {Key: "items", Value: fmt.Sprint(items)}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "guideline", ID, lData, audit)

	return guideline, nil
}

func (app *Application) deleteGuideline(current model.User, group string, ID string) error {
	err := app.storage.DeleteGuideline(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "guideline", ID)

	return nil
}

func (app *Application) getGuidelinesByCountyID(countyID string) ([]*model.Guideline, error) {
	//1. first check if we have a county for the provided id
	county, err := app.storage.FindCounty(countyID)
	if err != nil {
		return nil, err
	}
	if county == nil {
		return nil, errors.New("there is no a county for the provided id")
	}

	//2. find the guidelines
	guidelines, err := app.storage.FindGuidelineByCountyID(countyID)
	if err != nil {
		return nil, err
	}
	return guidelines, nil
}

func (app *Application) createCountyStatus(current model.User, group string, audit *string, countyID string, name string, description string) (*model.CountyStatus, error) {
	//1. find if we have a county for the provided ID
	county, err := app.storage.FindCounty(countyID)
	if err != nil {
		return nil, err
	}
	if county == nil {
		return nil, errors.New("there is no a county for the provided id")
	}

	//2. create the county status entity
	countyStatus, err := app.storage.CreateCountyStatus(countyID, name, description)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "countyID", Value: countyID}, {Key: "name", Value: name}, {Key: "description", Value: description}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "county-status", countyStatus.ID, lData, audit)

	return countyStatus, nil
}

func (app *Application) updateCountyStatus(current model.User, group string, audit *string, ID string, name string, description string) (*model.CountyStatus, error) {
	countyStatus, err := app.storage.FindCountyStatus(ID)
	if err != nil {
		return nil, err
	}
	if countyStatus == nil {
		return nil, errors.New("county status is nil for id " + ID)
	}

	//add the new values
	countyStatus.Name = name
	countyStatus.Description = description

	//save it
	err = app.storage.SaveCountyStatus(countyStatus)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "name", Value: name}, {Key: "description", Value: description}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "county-status", ID, lData, audit)

	return countyStatus, nil
}

func (app *Application) deleteCountyStatus(current model.User, group string, ID string) error {
	err := app.storage.DeleteCountyStatus(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "county-status", ID)

	return nil
}

func (app *Application) getCountyStatusByCountyID(countyID string) ([]*model.CountyStatus, error) {
	//1. first check if we have a county for the provided id
	county, err := app.storage.FindCounty(countyID)
	if err != nil {
		return nil, err
	}
	if county == nil {
		return nil, errors.New("there is no a county for the provided id")
	}

	//2. find the county statuses
	countyStatuses, err := app.storage.FindCountyStatusesByCountyID(countyID)
	if err != nil {
		return nil, err
	}
	return countyStatuses, nil
}

func (app *Application) getTestTypes() ([]*model.TestType, error) {
	testTypes, err := app.storage.ReadAllTestTypes()
	if err != nil {
		return nil, err
	}
	return testTypes, nil
}

func (app *Application) createTestType(current model.User, group string, audit *string, name string, priority *int) (*model.TestType, error) {
	testType, err := app.storage.CreateTestType(name, priority)
	if err != nil {
		return nil, err
	}

	//audit

	priority = nil

	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "name", Value: name}, {Key: "priority", Value: fmt.Sprint(utils.GetInt(priority))}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "test-type", testType.ID, lData, audit)

	return testType, nil
}

func (app *Application) updateTestType(current model.User, group string, audit *string, ID string, name string, priority *int) (*model.TestType, error) {
	testType, err := app.storage.FindTestType(ID)
	if err != nil {
		return nil, err
	}
	if testType == nil {
		return nil, errors.New("test type is nil for id " + ID)
	}

	//add the new values
	testType.Name = name
	testType.Priority = priority

	//save it
	err = app.storage.SaveTestType(testType)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "name", Value: name}, {Key: "priority", Value: fmt.Sprint(utils.GetInt(priority))}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "test-type", ID, lData, audit)

	return testType, nil
}

func (app *Application) deleteTestType(current model.User, group string, ID string) error {
	err := app.storage.DeleteTestType(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "test-type", ID)

	return nil
}

func (app *Application) createTestTypeResult(current model.User, group string, audit *string, testTypeID string, name string, nextStep string, nextStepOffset *int, resultExpiresOffset *int) (*model.TestTypeResult, error) {
	//1. find if we have a test type for the provided ID
	testType, err := app.storage.FindTestType(testTypeID)
	if err != nil {
		return nil, err
	}
	if testType == nil {
		return nil, errors.New("there is no a test type for the provided id")
	}

	//2. create the test type result entity
	testTypeResult, err := app.storage.CreateTestTypeResult(testTypeID, name, nextStep, nextStepOffset, resultExpiresOffset)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "testTypeID", Value: testTypeID}, {Key: "name", Value: name}, {Key: "nextStep", Value: nextStep},
		{Key: "nextStepOffset", Value: fmt.Sprint(utils.GetInt(nextStepOffset))}, {Key: "resultExpiresOffset", Value: fmt.Sprint(utils.GetInt(resultExpiresOffset))}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "test-type-result", testTypeResult.ID, lData, audit)

	return testTypeResult, nil
}

func (app *Application) updateTestTypeResult(current model.User, group string, audit *string, ID string, name string, nextStep string, nextStepOffset *int, resultExpiresOffset *int) (*model.TestTypeResult, error) {
	testTypeResult, err := app.storage.FindTestTypeResult(ID)
	if err != nil {
		return nil, err
	}
	if testTypeResult == nil {
		return nil, errors.New("test type result is nil for id " + ID)
	}

	//add the new values
	testTypeResult.Name = name
	testTypeResult.NextStep = nextStep
	testTypeResult.NextStepOffset = nextStepOffset
	testTypeResult.ResultExpiresOffset = resultExpiresOffset

	//save it
	err = app.storage.SaveTestTypeResult(testTypeResult)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "name", Value: name}, {Key: "nextStep", Value: nextStep},
		{Key: "nextStepOffset", Value: fmt.Sprint(utils.GetInt(nextStepOffset))}, {Key: "resultExpiresOffset", Value: fmt.Sprint(utils.GetInt(resultExpiresOffset))}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "test-type-result", ID, lData, audit)

	return testTypeResult, nil
}

func (app *Application) deleteTestTypeResult(current model.User, group string, ID string) error {
	err := app.storage.DeleteTestTypeResult(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "test-type-result", ID)

	return nil
}

func (app *Application) getTestTypeResultsByTestTypeID(testTypeID string) ([]*model.TestTypeResult, error) {
	//1. first check if we have a test type for the provided id
	testType, err := app.storage.FindTestType(testTypeID)
	if err != nil {
		return nil, err
	}
	if testType == nil {
		return nil, errors.New("there is no a test type for the provided id")
	}

	//2. find the test type results
	testTypeResults, err := app.storage.FindTestTypeResultsByTestTypeID(testTypeID)
	if err != nil {
		return nil, err
	}
	return testTypeResults, nil
}

func (app *Application) getRules() ([]*model.Rule, error) {
	rules, err := app.storage.ReadAllRules()
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func (app *Application) getCRules(countyID string, appVersion string) (*model.CRules, error) {
	supported, v := app.isVersionSupported(appVersion)
	if !supported {
		return nil, errors.New("app version is not supported")
	}

	cRules, err := app.storage.FindCRulesByCountyID(*v, countyID)
	if err != nil {
		return nil, err
	}
	return cRules, nil
}

func (app *Application) createOrUpdateCRules(current model.User, group string, audit *string, countyID string, appVersion string, data string) error {
	supported, v := app.isVersionSupported(appVersion)
	if !supported {
		return errors.New("app version is not supported")
	}

	create, err := app.storage.CreateOrUpdateCRules(*v, countyID, data)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "countyID", Value: countyID}, {Key: "appVersion", Value: appVersion}, {Key: "data", Value: data}}
	if *create {
		//create
		defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "crules", "", lData, audit)
	} else {
		//update
		defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "crules", "", lData, audit)
	}

	return nil
}

func (app *Application) getASymptoms(appVersion string) (*model.Symptoms, error) {
	supported, v := app.isVersionSupported(appVersion)
	if !supported {
		return nil, errors.New("app version is not supported")
	}

	symptoms, err := app.storage.ReadSymptoms(*v)
	if err != nil {
		return nil, err
	}
	return symptoms, nil
}

func (app *Application) createOrUpdateSymptoms(current model.User, group string, audit *string, appVersion string, items string) error {
	supported, v := app.isVersionSupported(appVersion)
	if !supported {
		return errors.New("app version is not supported")
	}

	create, err := app.storage.CreateOrUpdateSymptoms(*v, items)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "appVersion", Value: appVersion}, {Key: "items", Value: items}}
	if *create {
		//create
		defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "symptoms", "", lData, audit)
	} else {
		//update
		defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "symptoms", "", lData, audit)
	}

	return nil
}

func (app *Application) getUINOverrides(uin *string, sort *string) ([]*model.UINOverride, error) {
	uinOverrides, err := app.storage.FindUINOverrides(uin, sort)
	if err != nil {
		return nil, err
	}
	return uinOverrides, nil
}

func (app *Application) createUINOverride(current model.User, group string, audit *string, uin string, interval int, category *string, expiration *time.Time) (*model.UINOverride, error) {
	uinOverride, err := app.storage.CreateUINOverride(uin, interval, category, expiration)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "uin", Value: uin}, {Key: "interval", Value: fmt.Sprint(interval)}, {Key: "category", Value: utils.GetString(category)}, {Key: "expiration", Value: utils.GetTime(expiration)}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "uin-override", uin, lData, audit)

	return uinOverride, nil
}

func (app *Application) updateUINOverride(current model.User, group string, audit *string, uin string, interval int, category *string, expiration *time.Time) (*string, error) {
	result, err := app.storage.UpdateUINOverride(uin, interval, category, expiration)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "uin", Value: uin}, {Key: "interval", Value: fmt.Sprint(interval)}, {Key: "category", Value: utils.GetString(category)}, {Key: "expiration", Value: utils.GetTime(expiration)}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "uin-override", uin, lData, audit)

	return result, nil
}

func (app *Application) deleteUINOverride(current model.User, group string, uin string) error {
	err := app.storage.DeleteUINOverride(uin)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "uin-override", uin)

	return nil
}

func (app *Application) createRule(current model.User, group string, audit *string, countyID string, testTypeID string, priority *int,
	resultsStatuses []model.TestTypeResultCountyStatus) (*model.Rule, error) {

	//TODO - transactions, consistency!!!

	//First validate
	//1. Check if we have a county with the provided ID
	county, err := app.storage.FindCounty(countyID)
	if err != nil {
		return nil, err
	}
	if county == nil {
		return nil, errors.New("there is no a county for the provided id")
	}

	//2. Check if we have a test type with the provided ID
	testType, err := app.storage.FindTestType(testTypeID)
	if err != nil {
		return nil, err
	}
	if testType == nil {
		return nil, errors.New("there is no a test type for the provided id")
	}

	//3. Check if we already have a rule for the provided county id and test type id.
	founded, err := app.storage.FindRuleByCountyIDTestTypeID(countyID, testTypeID)
	if err != nil {
		return nil, err
	}
	if founded != nil {
		return nil, errors.New("there is already a rule for this county and this test type")
	}

	//4. Check if the rule data is vaid
	valid, err := app.isRuleDataValid(countyID, testTypeID, resultsStatuses)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, err
	}

	//5. Now create it
	rule, err := app.storage.CreateRule(countyID, testTypeID, priority, resultsStatuses)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "countyID", Value: countyID}, {Key: "testTypeID", Value: testTypeID},
		{Key: "priority", Value: fmt.Sprint(utils.GetInt(priority))}, {Key: "resultsStatuses", Value: fmt.Sprint(resultsStatuses)}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "rule", rule.ID, lData, audit)

	return rule, nil
}

func (app *Application) updateRule(current model.User, group string, audit *string, ID string, priority *int, resultsStatuses []model.TestTypeResultCountyStatus) (*model.Rule, error) {
	//1. find the rule
	rule, err := app.storage.FindRule(ID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, errors.New("rule is nil for id " + ID)
	}

	//2. Check if the rule data is valid
	valid, err := app.isRuleDataValid(rule.County.ID, rule.TestType.ID, resultsStatuses)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, err
	}

	//3. Add the new values
	rule.Priority = priority
	rule.ResultsStates = resultsStatuses

	//4. Save it
	err = app.storage.SaveRule(rule)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "priority", Value: fmt.Sprint(utils.GetInt(priority))}, {Key: "resultsStatuses", Value: fmt.Sprint(resultsStatuses)}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "rule", ID, lData, audit)

	return rule, nil
}

func (app *Application) isRuleDataValid(countyID string, testTypeID string, resultsStatuses []model.TestTypeResultCountyStatus) (bool, error) {
	//1. Check if the provided test type results are part from this test type
	testTypeResultsValid, err := app.areTestTypeResultsValid(testTypeID, resultsStatuses)
	if err != nil {
		return false, err
	}
	if !testTypeResultsValid {
		return false, errors.New("the provided test types results are not valid for this test type")
	}

	//2. Check if the provided county statuses are part from this county
	countyStatusesValid, err := app.areCountyStatusesValid(countyID, resultsStatuses)
	if err != nil {
		return false, err
	}
	if !countyStatusesValid {
		return false, errors.New("the provided county statuses are not valid for this county")
	}
	return true, nil
}

func (app *Application) areTestTypeResultsValid(testTypeID string, resultsStates []model.TestTypeResultCountyStatus) (bool, error) {
	testTypeResults, err := app.storage.FindTestTypeResultsByTestTypeID(testTypeID)
	if err != nil {
		return false, err
	}
	if resultsStates == nil || len(resultsStates) == 0 {
		//nothing to check
		return true, nil
	}

	for _, rs := range resultsStates {
		candidateID := rs.TestTypeResultID
		contains := app.containsTestTypeResultP(candidateID, testTypeResults)
		if !contains {
			return false, nil
		}
	}
	return true, nil
}

func (app *Application) areCountyStatusesValid(countyID string, resultsStates []model.TestTypeResultCountyStatus) (bool, error) {
	countyStatuses, err := app.storage.FindCountyStatusesByCountyID(countyID)
	if err != nil {
		return false, err
	}
	if countyStatuses == nil || len(countyStatuses) == 0 {
		//nothing to check
		return true, nil
	}

	for _, rs := range resultsStates {
		candidateID := rs.CountyStatusID
		contains := app.containsCountyStatusP(candidateID, countyStatuses)
		if !contains {
			return false, nil
		}
	}
	return true, nil
}

func (app *Application) containsTestTypeResultP(ID string, list []*model.TestTypeResult) bool {
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

func (app *Application) containsTestTypeResult(ID string, list []model.TestTypeResult) bool {
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

func (app *Application) deleteRule(current model.User, group string, ID string) error {
	err := app.storage.DeleteRule(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "rule", ID)

	return nil
}

func (app *Application) getLocations() ([]*model.Location, error) {
	locations, err := app.storage.ReadAllLocations()
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (app *Application) createLocation(current model.User, group string, audit *string, providerID string, countyID string, name string, address1 string, address2 string, city string,
	state string, zip string, country string, latitude float64, longitude float64, contact string,
	daysOfOperation []model.OperationDay, url string, notes string, waitTimeColor *string, availableTests []string) (*model.Location, error) {
	//1. check if the location data is valid
	err := app.isLocationDataValid(providerID, countyID, availableTests)
	if err != nil {
		return nil, err
	}

	//2. create the entity
	location, err := app.storage.CreateLocation(providerID, countyID, name, address1, address2, city,
		state, zip, country, latitude, longitude, contact, daysOfOperation, url, notes, waitTimeColor, availableTests)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "providerID", Value: providerID}, {Key: "countyID", Value: countyID}, {Key: "name", Value: name}, {Key: "address1", Value: address1},
		{Key: "address2", Value: address2}, {Key: "city", Value: city}, {Key: "state", Value: state}, {Key: "zip", Value: zip}, {Key: "country", Value: country},
		{Key: "latitude", Value: fmt.Sprint(latitude)}, {Key: "longitude", Value: fmt.Sprint(longitude)}, {Key: "contact", Value: contact},
		{Key: "daysOfOperation", Value: fmt.Sprint(daysOfOperation)}, {Key: "url", Value: url}, {Key: "notes", Value: notes}, {Key: "waitTimeColor", Value: utils.GetString(waitTimeColor)},
		{Key: "availableTests", Value: fmt.Sprint(availableTests)}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "location", location.ID, lData, audit)

	return location, nil
}

func (app *Application) updateLocation(current model.User, group string, audit *string, ID string, name string, address1 string, address2 string, city string,
	state string, zip string, country string, latitude float64, longitude float64, contact string,
	daysOfOperation []model.OperationDay, url string, notes string, waitTimeColor *string, availableTests []string) (*model.Location, error) {

	// find if we have a location for the provided id
	location, err := app.storage.FindLocation(ID)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, errors.New("location is nil for id " + ID)
	}

	// check if the provided test types ids are valid
	areTestTypesValid, err := app.areTestTypesValid(availableTests)
	if err != nil {
		return nil, err
	}
	if !areTestTypesValid {
		return nil, errors.New("the provided test types are not valid")
	}

	// add the new values
	location.Name = name
	location.Address1 = address1
	location.Address2 = address2
	location.City = city
	location.State = state
	location.ZIP = zip
	location.Country = country
	location.Latitude = latitude
	location.Longitude = longitude
	location.Timezone = "America/Chicago"
	location.Contact = contact
	location.DaysOfOperation = daysOfOperation
	location.URL = url
	location.Notes = notes
	location.WaitTimeColor = waitTimeColor
	var avTR []model.TestType
	if availableTests != nil {
		for _, id := range availableTests {
			testType := model.TestType{ID: id}
			avTR = append(avTR, testType)
		}
	}
	location.AvailableTests = avTR

	// save it
	err = app.storage.SaveLocation(location)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "name", Value: name}, {Key: "address1", Value: address1}, {Key: "address2", Value: address2}, {Key: "city", Value: city},
		{Key: "state", Value: state}, {Key: "zip", Value: zip}, {Key: "country", Value: country}, {Key: "latitude", Value: fmt.Sprint(latitude)},
		{Key: "longitude", Value: fmt.Sprint(longitude)}, {Key: "contact", Value: contact}, {Key: "daysOfOperation", Value: fmt.Sprint(daysOfOperation)},
		{Key: "url", Value: url}, {Key: "notes", Value: notes}, {Key: "waitTimeColor", Value: utils.GetString(waitTimeColor)},
		{Key: "availableTests", Value: fmt.Sprint(availableTests)}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "location", ID, lData, audit)

	return location, nil
}

func (app *Application) isLocationDataValid(providerID string, countyID string, availableTests []string) error {
	//1. check if we have a provider for the provided id
	provider, err := app.storage.FindProvider(providerID)
	if err != nil {
		return err
	}
	if provider == nil {
		return errors.New("there is no a provider for the provided id")
	}

	//2. check if we have a county for the provided id
	county, err := app.storage.FindCounty(countyID)
	if err != nil {
		return err
	}
	if county == nil {
		return errors.New("there is no a county for the provided id")
	}

	//3. check if the test types ids are valid
	areTestTypesValid, err := app.areTestTypesValid(availableTests)
	if err != nil {
		return err
	}
	if !areTestTypesValid {
		return errors.New("the provided test types are not valid")
	}

	return nil
}

func (app *Application) areTestTypesValid(availableTests []string) (bool, error) {
	if availableTests == nil {
		return true, nil
	}

	allTestTypes, err := app.storage.ReadAllTestTypes()
	if err != nil {
		return false, err
	}

	for _, item := range availableTests {
		contains := app.containsTestType(item, allTestTypes)
		if !contains {
			return false, errors.New("there is no a test type for " + item)
		}
	}
	return true, nil
}

func (app *Application) containsTestType(ID string, list []*model.TestType) bool {
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

func (app *Application) deleteLocation(current model.User, group string, ID string) error {
	err := app.storage.DeleteLocation(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "location", ID)

	return nil
}

func (app *Application) createSymptom(current model.User, group string, name string, symptomGroup string) (*model.Symptom, error) {
	symptom, err := app.storage.CreateSymptom(name, symptomGroup)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "name", Value: name}, {Key: "symptomGroup", Value: symptomGroup}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "symptom", symptom.ID, lData, nil)

	return symptom, nil
}

func (app *Application) updateSymptom(current model.User, group string, ID string, name string) (*model.Symptom, error) {
	symptom, err := app.storage.FindSymptom(ID)
	if err != nil {
		return nil, err
	}
	if symptom == nil {
		return nil, errors.New("symptom nil for id " + ID)
	}

	//add the new values
	symptom.Name = name

	//save it
	err = app.storage.SaveSymptom(symptom)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "name", Value: name}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "symptom", ID, lData, nil)

	return symptom, nil
}

func (app *Application) deleteSymptom(current model.User, group string, ID string) error {
	err := app.storage.DeleteSymptom(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "symptom", ID)
	return nil
}

func (app *Application) getSymptomRules() ([]*model.SymptomRule, error) {
	symptomRules, err := app.storage.ReadAllSymptomRules()
	if err != nil {
		return nil, err
	}
	return symptomRules, nil
}

func (app *Application) createSymptomRule(current model.User, group string, countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error) {
	// validate the data
	err := app.validateSymptomRuleData(countyID, gr1Count, gr2Count, items)
	if err != nil {
		return nil, err
	}

	//create it
	symptomRule, err := app.storage.CreateSymptomRule(countyID, gr1Count, gr2Count, items)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "countyID", Value: countyID}, {Key: "gr1Count", Value: fmt.Sprint(gr1Count)},
		{Key: "gr2Count", Value: fmt.Sprint(gr2Count)}, {Key: "items", Value: fmt.Sprint(items)}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "symptom-rule", symptomRule.ID, lData, nil)

	return symptomRule, nil
}

func (app *Application) updateSymptomRule(current model.User, group string, ID string, countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error) {
	//1. find the symptom rule
	symptomRule, err := app.storage.FindSymptomRule(ID)
	if err != nil {
		return nil, err
	}
	if symptomRule == nil {
		return nil, errors.New("symptom rule is nil for id " + ID)
	}

	// validate the data
	err = app.validateSymptomRuleData(countyID, gr1Count, gr2Count, items)
	if err != nil {
		return nil, err
	}

	//3. Add the new values
	symptomRule.County = model.County{ID: countyID}
	symptomRule.Gr1Count = gr1Count
	symptomRule.Gr2Count = gr2Count
	symptomRule.Items = items

	//4. Save it
	err = app.storage.SaveSymptomRule(symptomRule)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "countyID", Value: countyID}, {Key: "gr1Count", Value: fmt.Sprint(gr1Count)},
		{Key: "gr2Count", Value: fmt.Sprint(gr2Count)}, {Key: "items", Value: fmt.Sprint(items)}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "symptom-rule", ID, lData, nil)

	return symptomRule, nil
}

func (app *Application) validateSymptomRuleData(countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) error {
	//check if there is a county for the provided id
	county, err := app.storage.FindCounty(countyID)
	if err != nil {
		return err
	}
	if county == nil {
		return errors.New("there is no a county for the provided id")
	}

	//check if the county statuses are valid
	countyStatuses := county.CountyStatuses
	for _, item := range items {
		contains := app.containsCountyStatus(item.CountyStatus.ID, countyStatuses)
		if !contains {
			return errors.New("there is invalid county status id")
		}
	}

	//check if we have a full combination
	valid := app.areGr1Gr2Valid(true, true, items)
	if !valid {
		return errors.New("invalid gr1 and gr2 items")
	}
	valid = app.areGr1Gr2Valid(true, false, items)
	if !valid {
		return errors.New("invalid gr1 and gr2 items")
	}
	valid = app.areGr1Gr2Valid(false, true, items)
	if !valid {
		return errors.New("invalid gr1 and gr2 items")
	}
	valid = app.areGr1Gr2Valid(false, false, items)
	if !valid {
		return errors.New("invalid gr1 and gr2 items")
	}

	return nil
}

func (app *Application) areGr1Gr2Valid(gr1 bool, gr2 bool, items []model.SymptomRuleItem) bool {
	for _, item := range items {
		if item.Gr1 == gr1 && item.Gr2 == gr2 {
			return true
		}
	}
	return false
}

func (app *Application) deleteSymptomRule(current model.User, group string, ID string) error {
	err := app.storage.DeleteSymptomRule(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "symptom-rule", ID)

	return nil
}

func (app *Application) getManualTestByCountyID(countyID string, status *string) ([]*model.EManualTest, error) {
	manualTests, err := app.storage.FindManualTestsByCountyIDDeep(countyID, status)
	if err != nil {
		return nil, err
	}
	return manualTests, nil
}

func (app *Application) processManualTest(ID string, status string, encryptedKey *string, encryptedBlob *string) error {
	err := app.storage.ProcessManualTest(ID, status, encryptedKey, encryptedBlob)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) getManualTestImage(ID string) (*string, *string, error) {
	encryptedImageKey, encryptedImageBlob, err := app.storage.FindManualTestImage(ID)
	if err != nil {
		return nil, nil, err
	}
	return encryptedImageKey, encryptedImageBlob, nil
}

func (app *Application) getAccessRules() ([]*model.AccessRule, error) {
	accessRules, err := app.storage.ReadAllAccessRules()
	if err != nil {
		return nil, err
	}
	return accessRules, nil
}

func (app *Application) createAccessRule(current model.User, group string, audit *string, countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error) {
	accessRule, err := app.storage.CreateAccessRule(countyID, rules)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "countyID", Value: countyID}, {Key: "rules", Value: fmt.Sprint(rules)}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "access-rule", accessRule.ID, lData, audit)

	return accessRule, nil
}

func (app *Application) updateAccessRule(current model.User, group string, audit *string, ID string, countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error) {
	accessRule, err := app.storage.UpdateAccessRule(ID, countyID, rules)
	if err != nil {
		return nil, err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "countyID", Value: countyID}, {Key: "rules", Value: fmt.Sprint(rules)}}
	defer app.audit.LogUpdateEvent(userIdentifier, userInfo, group, "access-rule", ID, lData, audit)

	return accessRule, nil
}

func (app *Application) deleteAccessRule(current model.User, group string, ID string) error {
	err := app.storage.DeleteAccessRule(ID)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "access-rule", ID)
	return nil
}

func (app *Application) getUserByExternalID(externalID string) (*model.User, error) {
	user, err := app.storage.FindUserByExternalID(externalID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (app *Application) createRoster(current model.User, group string, audit *string, phone string, uin string) error {
	err := app.storage.CreateRoster(phone, uin)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "phone", Value: phone}, {Key: "uin", Value: uin}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "roster", "", lData, audit)

	return nil
}

func (app *Application) getRosters(filter *utils.Filter, sortBy string, sortOrder int, limit int, offset int) ([]map[string]interface{}, error) {
	rosters, err := app.storage.FindRosters(filter, sortBy, sortOrder, limit, offset)
	if err != nil {
		return nil, err
	}
	return rosters, nil
}

func (app *Application) deleteRosterByPhone(current model.User, group string, phone string) error {
	err := app.storage.DeleteRosterByPhone(phone)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "roster", fmt.Sprintf("phone:%s", phone))

	return nil
}

func (app *Application) deleteRosterByUIN(current model.User, group string, uin string) error {
	err := app.storage.DeleteRosterByUIN(uin)
	if err != nil {
		return err
	}

	//audit
	userIdentifier, userInfo := current.GetLogData()
	defer app.audit.LogDeleteEvent(userIdentifier, userInfo, group, "roster", fmt.Sprintf("uin:%s", uin))

	return nil
}

func (app *Application) createAction(current model.User, group string, audit *string, providerID string, userID string, encryptedKey string, encryptedBlob string) (*model.CTest, error) {
	//1. create a ctest
	item, user, err := app.storage.CreateAdminCTest(providerID, userID, encryptedKey, encryptedBlob, false, nil)
	if err != nil {
		return nil, err
	}

	//2. send a firebase notification to the user.
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

	//audit
	userIdentifier, userInfo := current.GetLogData()
	lData := []AuditDataEntry{{Key: "providerID", Value: providerID}, {Key: "userID", Value: userID},
		{Key: "encryptedKey", Value: encryptedKey}, {Key: "encryptedBlob", Value: encryptedBlob}}
	defer app.audit.LogCreateEvent(userIdentifier, userInfo, group, "action", item.ID, lData, audit)

	return item, nil
}

func (app *Application) getAudit(current model.User, group string, userIdentifier *string, entity *string, entityID *string, operation *string,
	clientData *string, createdAt *time.Time, sortBy *string, asc *bool, limit *int64) ([]*AuditEntity, error) {

	//Admin can look all logs
	var usedGroup *string
	if current.IsAdmin() {
		usedGroup = nil
	} else {
		usedGroup = &group
	}

	items, err := app.audit.Find(userIdentifier, usedGroup, entity, entityID, operation, clientData, createdAt, sortBy, asc, limit)
	if err != nil {
		return nil, err
	}
	return items, nil
}
