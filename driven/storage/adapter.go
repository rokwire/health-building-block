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

package storage

import (
	"context"
	"errors"
	"fmt"
	"health/core"
	"health/core/model"
	"health/utils"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type accessRule struct {
	ID       string                   `bson:"_id"`
	CountyID string                   `bson:"county_id"`
	Rules    []accessRuleCountyStatus `bson:"rules"`

	DateCreated time.Time  `bson:"date_created"`
	DateUpdated *time.Time `bson:"date_updated"`
}

type accessRuleCountyStatus struct {
	CountyStatusID string `bson:"county_status_id"`
	Value          string `bson:"value"` //granted or denied
}

type eManualTest struct {
	ID         string  `bson:"_id"`
	UserID     string  `bson:"user_id"`
	EHistoryID string  `bson:"ehistory_id"`
	LocationID *string `bson:"location_id"`
	CountyID   *string `bson:"county_id"`

	EncryptedKey  string `bson:"encrypted_key"`
	EncryptedBlob string `bson:"encrypted_blob"`

	EncryptedImageKey  string `bson:"encrypted_image_key"`
	EncryptedImageBlob string `bson:"encrypted_image_blob"`

	Status string `bson:"status"` //unverified, verified, rejected

	DateCreated time.Time `bson:"date_created"`
}

type symptomRule struct {
	ID       string `bson:"_id"`
	CountyID string `bson:"county_id"`

	Gr1Count int `bson:"gr1_count"`
	Gr2Count int `bson:"gr2_count"`

	Items []symptomRuleItem `bson:"items"`

	DateCreated time.Time  `bson:"date_created"`
	DateUpdated *time.Time `bson:"date_updated"`
}

type symptomRuleItem struct {
	Gr1            bool   `bson:"gr1"`
	Gr2            bool   `bson:"gr2"`
	CountyStatusID string `bson:"county_status_id"`
	NextStep       string `bson:"next_step"`
}

type symptom struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`

	DateCreated time.Time  `bson:"date_created"`
	DateUpdated *time.Time `bson:"date_updated"`
}

type symptomGroup struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`

	Symptoms []symptom `bson:"symptoms"`
}

type location struct {
	ID string `bson:"_id"`

	Name            string         `bson:"name"`
	Address1        string         `bson:"address_1"`
	Address2        string         `bson:"address_2"`
	City            string         `bson:"city"`
	State           string         `bson:"state"`
	ZIP             string         `bson:"zip"`
	Country         string         `bson:"country"`
	Latitude        float64        `bson:"latitude"`
	Longitude       float64        `bson:"longitude"`
	Timezone        string         `bson:"timezone"`
	Contact         string         `bson:"contact"`
	DaysOfOperation []operationDay `bson:"days_of_operation"`
	URL             string         `bson:"url"`
	Notes           string         `bson:"notes"`
	WaitTimeColor   *string        `bson:"wait_time_color"`

	ProviderID string `bson:"provider_id"`
	CountyID   string `bson:"county_id"`

	AvailableTests []string `bson:"available_tests"`

	DateCreated time.Time  `bson:"date_created"`
	DateUpdated *time.Time `bson:"date_updated"`
}

type operationDay struct {
	Name      string `bson:"name"`
	OpenTime  string `bson:"open_time"`
	CloseTime string `bson:"close_time"`
}

type rule struct {
	ID         string `bson:"_id"`
	CountyID   string `bson:"county_id"`
	TestTypeID string `bson:"test_type_id"`
	Priority   *int   `bson:"priority"`

	ResultsStates []testTypeResultCountyStatus `bson:"results_statuses"`

	DateCreated time.Time  `bson:"date_created"`
	DateUpdated *time.Time `bson:"date_updated"`
}

type testTypeResultCountyStatus struct {
	TestTypeResultID string `bson:"test_type_result_id"`
	CountyStatusID   string `bson:"county_status_id"`
}

type testType struct {
	ID       string `bson:"_id"`
	Name     string `bson:"name"`
	Priority *int   `bson:"priority"`

	Results []testTypeResult `bson:"results"`

	DateCreated time.Time  `bson:"date_created"`
	DateUpdated *time.Time `bson:"date_updated"`
}

type testTypeResult struct {
	ID                  string `bson:"_id"`
	Name                string `bson:"name"`
	NextStep            string `bson:"next_step"`
	NextStepOffset      *int   `bson:"next_step_offset"`
	ResultExpiresOffset *int   `bson:"result_expires_offset"`

	DateCreated time.Time  `bson:"date_created"`
	DateUpdated *time.Time `bson:"date_updated"`
}

type provider struct {
	ID                  string   `bson:"_id"`
	ProviderName        string   `bson:"provider_name"`
	ManualTest          bool     `bson:"manual_test"`
	AvailableMechanisms []string `bson:"available_mechanisms"`

	DateCreated time.Time  `bson:"date_created"`
	DateUpdated *time.Time `bson:"date_updated"`
}

type county struct {
	ID string `bson:"_id"`

	Name          string `bson:"name"`
	StateProvince string `bson:"state_province"`
	Country       string `bson:"country"`

	Guidelines     []guidline     `bson:"guidelines"`
	CountyStatuses []countyStatus `bson:"county_statuses"`

	DateCreated time.Time  `bson:"date_created"`
	DateUpdated *time.Time `bson:"date_updated"`
}

type guidline struct {
	ID          string `bson:"id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`

	Items []guidlineItem `bson:"items"`

	DateCreated time.Time  `bson:"date_created"`
	DateUpdated *time.Time `bson:"date_updated"`
}

type guidlineItem struct {
	Icon        string `bson:"icon"`
	Description string `bson:"description"`
	Type        string `bson:"type"`
}

type countyStatus struct {
	ID          string `bson:"id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`

	DateCreated time.Time  `bson:"date_created"`
	DateUpdated *time.Time `bson:"date_updated"`
}

//Adapter implements the Storage interface
type Adapter struct {
	db *database
}

//Start starts the storage
func (sa *Adapter) Start() error {
	err := sa.db.start()
	return err
}

//SetStorageListener sets listener for the storage
func (sa *Adapter) SetStorageListener(storageListener core.StorageListener) {
	sa.db.listener = storageListener
}

//ReadAllAppVersions reads all the versions. It gives them in a sorted way as the latest version is on position 0
func (sa *Adapter) ReadAllAppVersions() ([]string, error) {
	filter := bson.D{}
	var result []interface{}
	err := sa.db.appversions.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	res := make([]string, len(result))
	for index, current := range result {
		cc := current.(bson.D)
		v := cc.Map()["version"].(string)
		res[index] = v
	}

	//sort the versions list
	utils.SortVersions(res)

	return res, nil
}

//CreateAppVersion preates app version
func (sa *Adapter) CreateAppVersion(version string) error {
	item := bson.D{primitive.E{Key: "version", Value: version}}
	_, err := sa.db.appversions.InsertOne(item)
	if err != nil {
		return err
	}

	return nil
}

//ClearUserData removes all the user data in the storage. It uses a transaction
func (sa *Adapter) ClearUserData(userID string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			fmt.Println(err)
			return err
		}

		//get the user
		uFilter := bson.D{primitive.E{Key: "_id", Value: userID}}
		var usersResult []*model.User
		err = sa.db.users.FindWithContext(sessionContext, uFilter, &usersResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(usersResult) == 0 {
			abortTransaction(sessionContext)
			log.Printf("there is no a user for the provided user id - %s %s", userID, err)
			return err
		}
		user := usersResult[0]

		//remove the user accounts data
		if len(user.Accounts) > 0 {
			for _, account := range user.Accounts {

				// delete the account data
				err = sa.deleteAccountData(sessionContext, account)
				if err != nil {
					abortTransaction(sessionContext)
					return err
				}
			}
		}

		//remove from users
		usersFilter := bson.D{primitive.E{Key: "_id", Value: userID}}
		_, err = sa.db.users.DeleteOneWithContext(sessionContext, usersFilter, nil)
		if err != nil {
			abortTransaction(sessionContext)
			log.Printf("error deleting user record for a user - %s", err)
			return err
		}

		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.Printf("error on commiting a transaction - %s", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//FindUser finds the user for the provided id
func (sa *Adapter) FindUser(ID string) (*model.User, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*model.User
	err := sa.db.users.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	return result[0], nil
}

//FindUserByExternalID finds the user for the provided external id. It does not look in the accounts
func (sa *Adapter) FindUserByExternalID(externalID string) (*model.User, error) {
	filter := bson.D{primitive.E{Key: "external_id", Value: externalID}}
	var result []*model.User
	err := sa.db.users.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	return result[0], nil
}

//FindUserAccountsByExternalID looks in primary user + accounts
func (sa *Adapter) FindUserAccountsByExternalID(externalID string) (*model.User, error) {
	filter := bson.D{primitive.E{Key: "$or", Value: []interface{}{
		bson.D{primitive.E{Key: "external_id", Value: externalID}},
		bson.D{primitive.E{Key: "accounts.external_id", Value: externalID}},
	}}}
	var result []*model.User
	err := sa.db.users.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	return result[0], nil
}

//FindUserByShibbolethID finds the user for the provided shibboleth id
func (sa *Adapter) FindUserByShibbolethID(shibbolethID string) (*model.User, error) {
	filter := bson.D{primitive.E{Key: "shibboleth_auth.uiucedu_uin", Value: shibbolethID}}
	var result []*model.User
	err := sa.db.users.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	return result[0], nil
}

//FindUsersByRePost finds the users filtered by re_post
func (sa *Adapter) FindUsersByRePost(rePost bool) ([]*model.User, error) {
	filter := bson.D{primitive.E{Key: "re_post", Value: rePost}}
	var result []*model.User
	err := sa.db.users.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//CreateAppUser creates an app user
func (sa *Adapter) CreateAppUser(externalID string,
	userUUID string, publicKey string, consent bool, exposureNotification bool, rePost bool, encryptedKey *string, encryptedBlob *string) (*model.User, error) {
	var user *model.User

	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		userID, _ := uuid.NewUUID()

		var accounts []model.Account

		//add default account
		accounts = append(accounts, model.Account{ID: userID.String(), ExternalID: externalID, Default: true, Active: true})

		//check if there are added sub accounts for this user
		rawSubAccountsFilter := bson.D{primitive.E{Key: "primary_account", Value: externalID}}
		var rawSubAccounts []*model.RawSubAccount
		err = sa.db.rawsubaccounts.FindWithContext(sessionContext, rawSubAccountsFilter, &rawSubAccounts, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(rawSubAccounts) > 0 {
			//if there are, then hook them up
			for _, rsa := range rawSubAccounts {
				subAccID, _ := uuid.NewUUID()
				newSubAccount := model.Account{ID: subAccID.String(), ExternalID: rsa.UIN, Default: false, Active: true, FirstName: rsa.FirstName,
					MiddleName: rsa.MiddleName, LastName: rsa.LastName, BirthDate: rsa.BirthDate, Gender: rsa.Gender, Address1: rsa.Address1,
					Address2: rsa.Address2, Address3: rsa.Address3, City: rsa.City, State: rsa.State, ZipCode: rsa.ZipCode, Phone: rsa.Phone,
					Email: rsa.Email}

				accounts = append(accounts, newSubAccount)

				//update the raw sub account with the inserted account id
				rawSubAccountUpdateFilter := bson.D{primitive.E{Key: "uin", Value: rsa.UIN}}
				rawSubAccountUpdate := bson.D{
					primitive.E{Key: "$set", Value: bson.D{
						primitive.E{Key: "account_id", Value: subAccID.String()},
					}},
				}
				_, err := sa.db.rawsubaccounts.UpdateOneWithContext(sessionContext, rawSubAccountUpdateFilter, rawSubAccountUpdate, nil)
				if err != nil {
					abortTransaction(sessionContext)
					return err
				}
			}
		}

		//insert the created user
		user = &model.User{ID: userID.String(), ExternalID: externalID, UUID: userUUID,
			PublicKey: publicKey, Consent: consent, ExposureNotification: exposureNotification, RePost: rePost,
			EncryptedKey: encryptedKey, EncryptedBlob: encryptedBlob, Accounts: accounts, DateCreated: time.Now()}
		_, err = sa.db.users.InsertOneWithContext(sessionContext, user)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

//CreateAdminUser creates an admin user
func (sa *Adapter) CreateAdminUser(shibboAuth *model.ShibbolethAuth, externalID string,
	userUUID string, publicKey string, consent bool, exposureNotification bool, rePost bool, encryptedKey *string, encryptedBlob *string) (*model.User, error) {

	id, _ := uuid.NewUUID()

	dateCreated := time.Now()

	user := model.User{ID: id.String(), ShibbolethAuth: shibboAuth, ExternalID: externalID, UUID: userUUID,
		PublicKey: publicKey, Consent: consent, ExposureNotification: exposureNotification, RePost: rePost,
		EncryptedKey: encryptedKey, EncryptedBlob: encryptedBlob, DateCreated: dateCreated}
	_, err := sa.db.users.InsertOne(&user)
	if err != nil {
		return nil, err
	}

	//return the inserted item
	return &user, nil
}

//CreateDefaultAccount creates a default account for the user
func (sa *Adapter) CreateDefaultAccount(userID string) (*model.User, error) {
	var user *model.User

	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		// find the user
		findFilter := bson.D{primitive.E{Key: "_id", Value: userID}}
		var usersResult []*model.User
		err = sa.db.users.FindWithContext(sessionContext, findFilter, &usersResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if usersResult == nil || len(usersResult) == 0 {
			abortTransaction(sessionContext)
			return errors.New("there is no an user for the provided id " + userID)
		}
		user = usersResult[0]

		// error if already there is a default account
		if user.HasDefaultAccount() {
			abortTransaction(sessionContext)
			return errors.New("there is already a default account for this user " + userID)
		}

		// update it
		newAccounts := user.Accounts
		if newAccounts == nil {
			//there are no any accounts yet
			newAccounts = make([]model.Account, 0)
		}
		newAccounts = append(newAccounts, model.Account{ID: user.ID, ExternalID: user.ExternalID, Default: true, Active: true})
		updateFilter := bson.D{primitive.E{Key: "_id", Value: userID}}
		update := bson.D{
			primitive.E{Key: "$set", Value: bson.D{
				primitive.E{Key: "accounts", Value: newAccounts},
			}},
		}
		updateResult, err := sa.db.users.UpdateOne(updateFilter, update, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if updateResult.ModifiedCount != 1 {
			abortTransaction(sessionContext)
			return errors.New("modifier count != 1")
		}

		//the updated user is returned
		user.Accounts = newAccounts

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

//SaveUser saves the user
func (sa *Adapter) SaveUser(user *model.User) error {
	filter := bson.D{primitive.E{Key: "_id", Value: user.ID}}

	dateUpdated := time.Now()
	user.DateUpdated = &dateUpdated

	err := sa.db.users.ReplaceOne(filter, user, nil)
	if err != nil {
		return err
	}
	return nil
}

//ReadCovid19Config reads the covid19 configuration from the storage
func (sa *Adapter) ReadCovid19Config() (*model.COVID19Config, error) {
	filter := bson.D{primitive.E{Key: "name", Value: "covid19"}}
	var result []*model.COVID19Config
	err := sa.db.configs.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		return nil, errors.New("no covid19 config found")
	}
	if len(result) > 1 {
		return nil, errors.New("more than 1 covid19 configs were found")
	}
	return result[0], nil
}

//SaveCovid19Config saves the covid19 configuration to the storage
func (sa *Adapter) SaveCovid19Config(covid19Config *model.COVID19Config) error {
	filter := bson.D{primitive.E{Key: "name", Value: covid19Config.Name}}
	err := sa.db.configs.ReplaceOne(filter, covid19Config, nil)
	if err != nil {
		return err
	}
	return nil
}

//ReadAllResources reads all covid19 resources
func (sa *Adapter) ReadAllResources() ([]*model.Resource, error) {
	filter := bson.D{}
	var result []*model.Resource
	err := sa.db.resources.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//CreateResource creates a resource item
func (sa *Adapter) CreateResource(title string, link string, displayOrder int) (*model.Resource, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	resource := model.Resource{ID: id.String(), Title: title, Link: link, DisplayOrder: displayOrder}
	_, err = sa.db.resources.InsertOne(&resource)
	if err != nil {
		return nil, err
	}

	//return the inserted item
	return &resource, nil
}

//DeleteResource deletes a resource item
func (sa *Adapter) DeleteResource(ID string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	result, err := sa.db.resources.DeleteOne(filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return errors.New("result is nil for resource item with id " + ID)
	}
	deletedCount := result.DeletedCount
	if deletedCount != 1 {
		return errors.New("error occured while deleting a resource item with id " + ID)
	}
	return nil
}

//FindResource finds resource
func (sa *Adapter) FindResource(ID string) (*model.Resource, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*model.Resource
	err := sa.db.resources.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	return result[0], nil
}

//SaveResource saves resource entity to the storage
func (sa *Adapter) SaveResource(resource *model.Resource) error {
	filter := bson.D{primitive.E{Key: "_id", Value: resource.ID}}
	err := sa.db.resources.ReplaceOne(filter, resource, nil)
	if err != nil {
		return err
	}
	return nil
}

//ReadFAQ reads the covid19 FAQs
func (sa *Adapter) ReadFAQ() (*model.FAQ, error) {
	filter := bson.D{}
	var result []*model.FAQ
	err := sa.db.faq.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if len(result) < 1 {
		return nil, errors.New("no faq data")
	}
	return result[0], nil
}

//SaveFAQ saves faq entity to the storage
func (sa *Adapter) SaveFAQ(faq *model.FAQ) error {
	//It is always 1 item
	err := sa.db.faq.ReplaceOne(bson.D{}, faq, nil)
	if err != nil {
		return err
	}
	return nil
}

//DeleteFAQSection deletes a faq section
func (sa *Adapter) DeleteFAQSection(ID string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. get the faq
		filter := bson.D{}
		var faqRes []*model.FAQ
		err = sa.db.faq.FindWithContext(sessionContext, filter, &faqRes, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(faqRes) < 1 {
			abortTransaction(sessionContext)
			return err
		}
		faq := faqRes[0]

		//2. find the section
		sections := faq.Sections
		if sections == nil {
			abortTransaction(sessionContext)
			return errors.New("for some reasons the sections are nil")
		}
		sectionIndex := -1
		for sIndex, section := range *sections {
			if section.ID == ID {
				sectionIndex = sIndex
			}
		}
		if sectionIndex == -1 {
			abortTransaction(sessionContext)
			return errors.New("cannot remove faq section with the provided id " + ID)
		}
		section := (*sections)[sectionIndex]

		//3. check if there are question for this section
		if section.Questions != nil && len(*section.Questions) > 0 {
			abortTransaction(sessionContext)
			return errors.New("cannot remove faq section because there are associated questions")
		}

		//4. remove the section from the faq
		modifiedList := removeSection(*sections, sectionIndex)
		faq.Sections = &modifiedList

		//5. save the faq
		err = sa.db.faq.ReplaceOneWithContext(sessionContext, bson.D{}, faq, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func removeSection(sections []*model.Section, s int) []*model.Section {
	return append(sections[:s], sections[s+1:]...)
}

//ReadNews reads all covid19 news
func (sa *Adapter) ReadNews(limit int64) ([]*model.News, error) {
	filter := bson.D{}
	var result []*model.News

	options := options.Find()
	options.SetSort(bson.D{primitive.E{Key: "date", Value: -1}}) //sort by "date"

	if limit > 0 {
		options.SetLimit(limit)
	}

	err := sa.db.news.Find(filter, &result, options)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//CreateNews creates a new covid19 news
func (sa *Adapter) CreateNews(date time.Time, title string, description string, htmlContent string, link *string) (*model.News, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	news := model.News{ID: id.String(), Date: date, Title: title, Description: description,
		HTMLContent: htmlContent, Link: link}
	insertedID, err := sa.db.news.InsertOne(&news)
	if err != nil {
		return nil, err
	}

	//return the inserted item
	news.ID = insertedID.(string)
	return &news, nil
}

//DeleteNews deletes a new covid19 news
func (sa *Adapter) DeleteNews(ID string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	result, err := sa.db.news.DeleteOne(filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return errors.New("result is nil for item with id " + ID)
	}
	deletedCount := result.DeletedCount
	if deletedCount != 1 {
		return errors.New("error occured while deleting an item with id " + ID)
	}
	return nil
}

//FindNews finds news
func (sa *Adapter) FindNews(ID string) (*model.News, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*model.News
	err := sa.db.news.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	return result[0], nil
}

//SaveNews saves news entity to the storage
func (sa *Adapter) SaveNews(news *model.News) error {
	filter := bson.D{primitive.E{Key: "_id", Value: news.ID}}
	err := sa.db.news.ReplaceOne(filter, news, nil)
	if err != nil {
		return err
	}
	return nil
}

//CreateEStatus creates a new covid19 passport status
func (sa *Adapter) CreateEStatus(appVersion *string, accountID string, date *time.Time, encryptedKey string, encryptedBlob string) (*model.EStatus, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	status := model.EStatus{ID: id.String(), AppVersion: appVersion, UserID: accountID, Date: date, EncryptedKey: encryptedKey, EncryptedBlob: encryptedBlob}
	_, err = sa.db.estatus.InsertOne(&status)
	if err != nil {
		return nil, err
	}

	//return the inserted item
	return &status, nil
}

//FindEStatusByAccountID finds a status by account id
func (sa *Adapter) FindEStatusByAccountID(appVersion *string, accountID string) (*model.EStatus, error) {
	filter := bson.D{primitive.E{Key: "user_id", Value: accountID},
		primitive.E{Key: "app_version", Value: appVersion}}
	var result []*model.EStatus
	err := sa.db.estatus.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	return result[0], nil
}

//SaveEStatus saves the status
func (sa *Adapter) SaveEStatus(status *model.EStatus) error {
	filter := bson.D{primitive.E{Key: "_id", Value: status.ID}}

	dateUpdated := time.Now()
	status.DateUpdated = &dateUpdated

	err := sa.db.estatus.ReplaceOne(filter, status, nil)
	if err != nil {
		return err
	}
	return nil
}

//DeleteEStatus deletes the status for the user account
func (sa *Adapter) DeleteEStatus(appVersion *string, accountID string) error {
	filter := bson.D{primitive.E{Key: "user_id", Value: accountID},
		primitive.E{Key: "app_version", Value: appVersion}}
	result, err := sa.db.estatus.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("error deleting a estatus - %s", err)
		return err
	}
	if result == nil {
		return errors.New("result is nil forestatus with account id " + accountID)
	}
	deletedCount := result.DeletedCount
	if deletedCount == 0 {
		return errors.New("there is no a estatus for account id " + accountID)
	}
	if deletedCount > 1 {
		return errors.New("deleted more than one records for account id " + accountID)
	}
	return nil
}

//CreateEHistory creates a history
func (sa *Adapter) CreateEHistory(accountID string, date time.Time, eType string, encryptedKey string, encryptedBlob string) (*model.EHistory, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	history := model.EHistory{ID: id.String(), UserID: accountID, Date: date, Type: eType,
		EncryptedKey: encryptedKey, EncryptedBlob: encryptedBlob}
	_, err = sa.db.ehistory.InsertOne(&history)
	if err != nil {
		return nil, err
	}

	//return the inserted item
	return &history, nil
}

//CreateManualЕHistory creates a history
func (sa *Adapter) CreateManualЕHistory(accountID string, date time.Time, encryptedKey string, encryptedBlob string, encryptedImageKey *string, encryptedImageBlob *string,
	countyID *string, locationID *string) (*model.EHistory, error) {
	var history model.EHistory

	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. insert history item
		historyID, _ := uuid.NewUUID()
		history = model.EHistory{ID: historyID.String(), UserID: accountID, Date: date, Type: "unverified_manual_test",
			EncryptedKey: encryptedKey, EncryptedBlob: encryptedBlob}
		insertedID, err := sa.db.ehistory.InsertOneWithContext(sessionContext, &history)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//2. insert manual test item
		manualTestID, _ := uuid.NewUUID()
		manualTest := eManualTest{ID: manualTestID.String(), UserID: accountID, EHistoryID: insertedID.(string),
			LocationID: locationID, CountyID: countyID, EncryptedKey: encryptedKey, EncryptedBlob: encryptedBlob,
			EncryptedImageKey: *encryptedImageKey, EncryptedImageBlob: *encryptedImageBlob, Status: "unverified", DateCreated: time.Now()}
		_, err = sa.db.emanualtests.InsertOneWithContext(sessionContext, &manualTest)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &history, nil
}

//FindEHistories finds all histories for an user account
func (sa *Adapter) FindEHistories(accountID string) ([]*model.EHistory, error) {
	filter := bson.D{primitive.E{Key: "user_id", Value: accountID}}
	var result []*model.EHistory

	options := options.Find()
	options.SetSort(bson.D{primitive.E{Key: "date", Value: -1}}) //sort by "date"

	err := sa.db.ehistory.Find(filter, &result, options)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//DeleteEHistories deletes all histories for an user account
func (sa *Adapter) DeleteEHistories(accountID string) (int64, error) {
	filter := bson.D{primitive.E{Key: "user_id", Value: accountID}}

	result, err := sa.db.ehistory.DeleteMany(filter, nil)
	if err != nil {
		return -1, err
	}
	if result == nil {
		return -1, errors.New("delete result is nil for some reasons")
	}

	//return the inserted item
	return result.DeletedCount, nil
}

//FindEHistory finds a history item
func (sa *Adapter) FindEHistory(ID string) (*model.EHistory, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*model.EHistory
	err := sa.db.ehistory.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	history := result[0]
	return history, nil
}

//SaveEHistory saves a history item
func (sa *Adapter) SaveEHistory(history *model.EHistory) error {
	filter := bson.D{primitive.E{Key: "_id", Value: history.ID}}
	err := sa.db.ehistory.ReplaceOne(filter, history, nil)
	if err != nil {
		return err
	}

	return nil
}

//ReadAllProviders reads all the providers
func (sa *Adapter) ReadAllProviders() ([]*model.Provider, error) {
	filter := bson.D{}
	var result []*provider
	err := sa.db.providers.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	var resultList []*model.Provider
	if result != nil {
		for _, current := range result {
			item := &model.Provider{ID: current.ID, Name: current.ProviderName, ManualTest: current.ManualTest, AvailableMechanisms: current.AvailableMechanisms}
			resultList = append(resultList, item)
		}
	}
	return resultList, nil
}

//CreateProvider creates a provider
func (sa *Adapter) CreateProvider(providerName string, manualTest bool, availableMechanisms []string) (*model.Provider, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	dateCreated := time.Now()

	provider := provider{ID: id.String(), ProviderName: providerName, ManualTest: manualTest, AvailableMechanisms: availableMechanisms, DateCreated: dateCreated}
	_, err = sa.db.providers.InsertOne(&provider)
	if err != nil {
		return nil, err
	}

	//return the inserted item
	result := &model.Provider{ID: provider.ID, Name: provider.ProviderName, ManualTest: provider.ManualTest, AvailableMechanisms: provider.AvailableMechanisms}
	return result, nil
}

//FindProvider finds a provider
func (sa *Adapter) FindProvider(ID string) (*model.Provider, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*provider
	err := sa.db.providers.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	provider := result[0]
	resultEntity := &model.Provider{ID: provider.ID, Name: provider.ProviderName, ManualTest: provider.ManualTest, AvailableMechanisms: provider.AvailableMechanisms}
	return resultEntity, nil
}

//SaveProvider saves a provider
func (sa *Adapter) SaveProvider(entity *model.Provider) error {
	findFilter := bson.D{primitive.E{Key: "_id", Value: entity.ID}}
	var result []*provider
	err := sa.db.providers.Find(findFilter, &result, nil)
	if err != nil {
		return err
	}
	if result == nil || len(result) == 0 {
		//not found
		return errors.New("there is no a provider for the provided id")
	}
	provider := result[0]

	//update the values
	provider.ProviderName = entity.Name
	provider.ManualTest = entity.ManualTest
	provider.AvailableMechanisms = entity.AvailableMechanisms
	dateUpdated := time.Now()
	provider.DateUpdated = &dateUpdated

	filter := bson.D{primitive.E{Key: "_id", Value: provider.ID}}
	err = sa.db.providers.ReplaceOne(filter, provider, nil)
	if err != nil {
		return err
	}

	return nil
}

//DeleteProvider deletes a provider
func (sa *Adapter) DeleteProvider(ID string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. check if there are associated locations for this provider
		locationsFilter := bson.D{primitive.E{Key: "provider_id", Value: ID}}
		var locResult []*location
		err = sa.db.locations.FindWithContext(sessionContext, locationsFilter, &locResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(locResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated locations for this provider")
		}

		//2. check if there are associated ctests for this provider
		ctestsFilter := bson.D{primitive.E{Key: "provider_id", Value: ID}}
		var cTestsResult []*model.CTest
		err = sa.db.ctests.FindWithContext(sessionContext, ctestsFilter, &cTestsResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(cTestsResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated ctests for this provider")
		}

		//3. now we can delete the provider
		deleteFilter := bson.D{primitive.E{Key: "_id", Value: ID}}
		result, err := sa.db.providers.DeleteOneWithContext(sessionContext, deleteFilter, nil)
		if err != nil {
			log.Printf("error deleting a provider - %s", err)
			abortTransaction(sessionContext)
			return err
		}
		if result == nil {
			abortTransaction(sessionContext)
			return errors.New("result is nil for provider with id " + ID)
		}
		deletedCount := result.DeletedCount
		if deletedCount == 0 {
			abortTransaction(sessionContext)
			return errors.New("there is no a provider for id " + ID)
		}
		if deletedCount > 1 {
			abortTransaction(sessionContext)
			return errors.New("deleted more than one records for id " + ID)
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//CreateExternalCTest creates an external ctests record
func (sa *Adapter) CreateExternalCTest(providerID string, uin string, encryptedKey string, encryptedBlob string, processed bool, orderNumber *string) (*model.CTest, *model.User, error) {
	var cTest model.CTest
	var user model.User

	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. check if there is a provider for the provided identifier
		provFilter := bson.D{primitive.E{Key: "_id", Value: providerID}}
		var provResult []*provider
		err = sa.db.providers.FindWithContext(sessionContext, provFilter, &provResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if provResult == nil || len(provResult) == 0 {
			//not found
			abortTransaction(sessionContext)
			return errors.New("there is no a provider for the provided identifier")
		}

		//2. check if there is a user with the provided uin - handle the accounts case
		userFilter := bson.D{primitive.E{Key: "$or", Value: []interface{}{
			bson.D{primitive.E{Key: "external_id", Value: uin}},
			bson.D{primitive.E{Key: "accounts.external_id", Value: uin}},
		}}}
		var userResult []*model.User
		err = sa.db.users.FindWithContext(sessionContext, userFilter, &userResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if userResult == nil || len(userResult) == 0 {
			//not found
			abortTransaction(sessionContext)
			return errors.New("there is no a user for the provided identifier")
		}
		user = *userResult[0]

		//3. find the account id
		var accountID string
		if len(user.Accounts) > 0 {
			for _, acc := range user.Accounts {
				if acc.ExternalID == uin {
					accountID = acc.ID
					break
				}
			}
		}
		if len(accountID) == 0 {
			//we can go in a situation where the user does not have any acccount yet.
			accountID = user.ID
		}

		//4. create a ctest
		id, err := uuid.NewUUID()
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		dateCreated := time.Now()
		cTest = model.CTest{ID: id.String(), ProviderID: providerID, UserID: accountID,
			EncryptedKey: encryptedKey, EncryptedBlob: encryptedBlob, Processed: processed, OrderNumber: orderNumber, DateCreated: dateCreated}
		_, err = sa.db.ctests.InsertOneWithContext(sessionContext, &cTest)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//5. Set the user re-post field as "false"
		sUserfilter := bson.D{primitive.E{Key: "_id", Value: user.ID}}
		dateUpdated := time.Now()
		user.DateUpdated = &dateUpdated
		user.RePost = false
		err = sa.db.users.ReplaceOneWithContext(sessionContext, sUserfilter, user, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return &cTest, &user, nil
}

//CreateAdminCTest creates an admin ctests record
func (sa *Adapter) CreateAdminCTest(providerID string, accountID string, encryptedKey string, encryptedBlob string, processed bool, orderNumber *string) (*model.CTest, *model.User, error) {
	var cTest model.CTest
	var user model.User

	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. check if there is a provider for the provided identifier
		provFilter := bson.D{primitive.E{Key: "_id", Value: providerID}}
		var provResult []*provider
		err = sa.db.providers.FindWithContext(sessionContext, provFilter, &provResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if provResult == nil || len(provResult) == 0 {
			//not found
			abortTransaction(sessionContext)
			return errors.New("there is no a provider for the provided identifier")
		}

		//2. check if there is a user with the provided identifier - handle the accounts case
		userFilter := bson.D{primitive.E{Key: "$or", Value: []interface{}{
			bson.D{primitive.E{Key: "_id", Value: accountID}},
			bson.D{primitive.E{Key: "accounts.id", Value: accountID}},
		}}}
		var userResult []*model.User
		err = sa.db.users.FindWithContext(sessionContext, userFilter, &userResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if userResult == nil || len(userResult) == 0 {
			//not found
			abortTransaction(sessionContext)
			return errors.New("there is no a user for the provided identifier")
		}

		//3. create a ctest
		id, err := uuid.NewUUID()
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		dateCreated := time.Now()
		cTest = model.CTest{ID: id.String(), ProviderID: providerID, UserID: accountID,
			EncryptedKey: encryptedKey, EncryptedBlob: encryptedBlob, Processed: processed, OrderNumber: orderNumber, DateCreated: dateCreated}
		_, err = sa.db.ctests.InsertOneWithContext(sessionContext, &cTest)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return &cTest, &user, nil
}

//FindCTest finds ctest
func (sa *Adapter) FindCTest(ID string) (*model.CTest, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*model.CTest
	err := sa.db.ctests.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	ctest := result[0]
	return ctest, nil
}

//FindCTests finds ctests for user account and processed
func (sa *Adapter) FindCTests(accountID string, processed bool) ([]*model.CTest, error) {
	filter := bson.D{primitive.E{Key: "user_id", Value: accountID},
		primitive.E{Key: "processed", Value: processed}}

	options := options.Find()
	options.SetSort(bson.D{primitive.E{Key: "date_created", Value: 1}}) //sort by "date_created"

	var result []*model.CTest
	err := sa.db.ctests.Find(filter, &result, options)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//FindCTestsByExternalUserIDs finds ctests lists for the provided external user IDs
func (sa *Adapter) FindCTestsByExternalUserIDs(externalUserIDs []string) (map[string][]*model.CTest, error) {
	//1. get based on the user
	ubList, err := sa.findCTBEIDBaesOnUser(externalUserIDs)
	if err != nil {
		return nil, err
	}

	//2. get based on the user accounts
	uabList, err := sa.findCTBEIDBaesOnUserAccounts(externalUserIDs)
	if err != nil {
		return nil, err
	}

	//3. merge the maps
	for key, value := range uabList {
		ubList[key] = value
	}

	return ubList, nil
}

func (sa *Adapter) findCTBEIDBaesOnUser(externalUserIDs []string) (map[string][]*model.CTest, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}},
		{"$match": bson.M{"user.external_id": bson.M{"$in": externalUserIDs}}},
		{"$unwind": "$user"},
		{"$project": bson.M{
			"_id": 1, "provider_id": 1, "order_number": 1, "encrypted_key": 1, "encrypted_blob": 1,
			"processed": 1, "date_created": 1, "date_updated": 1,
			"user_id": "$user._id", "user_external_id": "$user.external_id",
		}}}

	type ctuJoin struct {
		ID            string     `bson:"_id"`
		ProviderID    string     `bson:"provider_id"`
		OrderNumber   *string    `bson:"order_number"`
		EncryptedKey  string     `bson:"encrypted_key"`
		EncryptedBlob string     `bson:"encrypted_blob"`
		Processed     bool       `bson:"processed"`
		DateCreated   time.Time  `bson:"date_created"`
		DateUpdated   *time.Time `bson:"date_updated"`

		UserID         string `bson:"user_id"`
		UserExternalID string `bson:"user_external_id"`
	}

	var result []*ctuJoin
	err := sa.db.ctests.Aggregate(pipeline, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return make(map[string][]*model.CTest, 0), nil
	}

	//construct the result
	mapData := make(map[string][]*model.CTest, len(externalUserIDs))
	for _, v := range result {
		if v.OrderNumber != nil {
			userExternalID := v.UserExternalID
			list := mapData[userExternalID]
			if list == nil {
				list = []*model.CTest{}
			}
			list = append(list, &model.CTest{ID: v.ID, ProviderID: v.ProviderID, UserID: v.UserID,
				EncryptedKey: v.EncryptedKey, EncryptedBlob: v.EncryptedBlob, OrderNumber: v.OrderNumber, Processed: v.Processed,
				DateCreated: v.DateCreated, DateUpdated: v.DateUpdated})

			mapData[userExternalID] = list
		}
	}
	return mapData, nil
}

func (sa *Adapter) findCTBEIDBaesOnUserAccounts(externalUserIDs []string) (map[string][]*model.CTest, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "ctests",
			"localField":   "accounts.id",
			"foreignField": "user_id",
			"as":           "ctests",
		}},
		{"$unwind": "$accounts"},
		{"$match": bson.M{"accounts.external_id": bson.M{"$in": externalUserIDs}}},
		{"$project": bson.M{
			"_id": 1, "external_id": 1, "accounts": 1,
			"ctests": bson.M{
				"$filter": bson.M{
					"input": "$ctests",
					"as":    "ctest",
					"cond": bson.M{
						"$and": []bson.M{
							{"$eq": []interface{}{"$$ctest.user_id", "$accounts.id"}},
						},
					},
				},
			},
		}},
	}

	type usersAccountsJoin struct {
		ID         string `bson:"_id"`
		ExternalID string `bson:"external_id"`
		Account    struct {
			ID         string `json:"id" bson:"id"`
			ExternalID string `json:"external_id" bson:"external_id"`
		} `bson:"accounts"`

		CTests []struct {
			ID            string     `bson:"_id"`
			ProviderID    string     `bson:"provider_id"`
			OrderNumber   *string    `bson:"order_number"`
			EncryptedKey  string     `bson:"encrypted_key"`
			EncryptedBlob string     `bson:"encrypted_blob"`
			Processed     bool       `bson:"processed"`
			DateCreated   time.Time  `bson:"date_created"`
			DateUpdated   *time.Time `bson:"date_updated"`
			UserID        string     `bson:"user_id"`
		} `bson:"ctests"`
	}

	var result []usersAccountsJoin
	err := sa.db.users.Aggregate(pipeline, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return make(map[string][]*model.CTest, 0), nil
	}

	//construct the result
	mapData := make(map[string][]*model.CTest, len(result))
	for _, acc := range result {
		accountExternalID := acc.Account.ExternalID
		accountID := acc.Account.ID
		list := mapData[accountExternalID]
		if list == nil {
			list = []*model.CTest{}
		}

		if acc.CTests != nil {
			for _, current := range acc.CTests {
				if current.OrderNumber != nil {
					list = append(list, &model.CTest{ID: current.ID, ProviderID: current.ProviderID, UserID: accountID,
						EncryptedKey: current.EncryptedKey, EncryptedBlob: current.EncryptedBlob, OrderNumber: current.OrderNumber, Processed: current.Processed,
						DateCreated: current.DateCreated, DateUpdated: current.DateUpdated})
				}
			}
		}

		mapData[accountExternalID] = list
	}
	return mapData, nil
}

//DeleteCTests deletes all ctest for a user account
func (sa *Adapter) DeleteCTests(accountID string) (int64, error) {
	filter := bson.D{primitive.E{Key: "user_id", Value: accountID}}

	result, err := sa.db.ctests.DeleteMany(filter, nil)
	if err != nil {
		return -1, err
	}
	if result == nil {
		return -1, errors.New("delete result is nil for some reasons")
	}

	//return the inserted item
	return result.DeletedCount, nil
}

//SaveCTest saves the ctest
func (sa *Adapter) SaveCTest(entity *model.CTest) error {
	findFilter := bson.D{primitive.E{Key: "_id", Value: entity.ID}}
	var result []*model.CTest
	err := sa.db.ctests.Find(findFilter, &result, nil)
	if err != nil {
		return err
	}
	if result == nil || len(result) == 0 {
		//not found
		return errors.New("there is no a ctest for the provided id")
	}
	ctest := result[0]

	//update the values
	ctest.Processed = entity.Processed
	dateUpdated := time.Now()
	ctest.DateUpdated = &dateUpdated

	filter := bson.D{primitive.E{Key: "_id", Value: ctest.ID}}
	err = sa.db.ctests.ReplaceOne(filter, ctest, nil)
	if err != nil {
		return err
	}

	return nil
}

//CreateCounty creates a county
func (sa *Adapter) CreateCounty(name string, stateProvince string, country string) (*model.County, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	dateCreated := time.Now()

	county := county{ID: id.String(), Name: name, StateProvince: stateProvince, Country: country,
		DateCreated: dateCreated}
	_, err = sa.db.counties.InsertOne(&county)
	if err != nil {
		return nil, err
	}

	//return the inserted item
	result := &model.County{ID: county.ID, Name: county.Name,
		StateProvince: county.StateProvince, Country: county.Country}
	return result, nil
}

//FindCounty finds a county
func (sa *Adapter) FindCounty(ID string) (*model.County, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*county
	err := sa.db.counties.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	county := result[0]

	//guidelines
	var guidelines []model.Guideline
	if county.Guidelines != nil {
		for _, gl := range county.Guidelines {
			var glItems []model.GuidelineItem
			if gl.Items != nil {
				for _, inner := range gl.Items {
					glItemType := model.GuidelineItemType{Value: inner.Type}
					innerItem := model.GuidelineItem{Icon: inner.Icon,
						Description: inner.Description, Type: glItemType}
					glItems = append(glItems, innerItem)
				}
			}

			item := model.Guideline{ID: gl.ID, Name: gl.Name, Description: gl.Description, Items: glItems}
			guidelines = append(guidelines, item)
		}
	}

	//county statuses
	var countyStatuses []model.CountyStatus
	if county.CountyStatuses != nil {
		for _, cs := range county.CountyStatuses {
			item := model.CountyStatus{ID: cs.ID, Name: cs.Name, Description: cs.Description}
			countyStatuses = append(countyStatuses, item)
		}
	}

	resultEntity := &model.County{ID: county.ID, Name: county.Name, StateProvince: county.StateProvince,
		Country: county.Country, Guidelines: guidelines, CountyStatuses: countyStatuses}
	return resultEntity, nil
}

//SaveCounty saves a county
func (sa *Adapter) SaveCounty(entity *model.County) error {
	findFilter := bson.D{primitive.E{Key: "_id", Value: entity.ID}}
	var result []*county
	err := sa.db.counties.Find(findFilter, &result, nil)
	if err != nil {
		return err
	}
	if result == nil || len(result) == 0 {
		//not found
		return errors.New("there is no a county for the provided id")
	}
	county := result[0]

	//update the values
	county.Name = entity.Name
	county.StateProvince = entity.StateProvince
	county.Country = entity.Country

	filter := bson.D{primitive.E{Key: "_id", Value: county.ID}}

	dateUpdated := time.Now()
	county.DateUpdated = &dateUpdated

	err = sa.db.counties.ReplaceOne(filter, county, nil)
	if err != nil {
		return err
	}

	return nil
}

//FindCounties finds counties
func (sa *Adapter) FindCounties(f *utils.Filter) ([]*model.County, error) {
	//add filter
	var filter interface{}
	if f != nil {
		filter = constructDataFilter(f)
	}

	var result []*county
	err := sa.db.counties.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	var resultList []*model.County
	if result != nil {
		for _, county := range result {
			//guidelines
			var guidelines []model.Guideline
			if county.Guidelines != nil {
				for _, gl := range county.Guidelines {
					var glItems []model.GuidelineItem
					if gl.Items != nil {
						for _, inner := range gl.Items {
							glItemType := model.GuidelineItemType{Value: inner.Type}
							innerItem := model.GuidelineItem{Icon: inner.Icon,
								Description: inner.Description, Type: glItemType}
							glItems = append(glItems, innerItem)
						}
					}

					item := model.Guideline{ID: gl.ID, Name: gl.Name, Description: gl.Description, Items: glItems}
					guidelines = append(guidelines, item)
				}
			}

			//county statuses
			var countyStatuses []model.CountyStatus
			if county.CountyStatuses != nil {
				for _, cs := range county.CountyStatuses {
					item := model.CountyStatus{ID: cs.ID, Name: cs.Name, Description: cs.Description}
					countyStatuses = append(countyStatuses, item)
				}
			}

			entity := &model.County{ID: county.ID, Name: county.Name, StateProvince: county.StateProvince,
				Country: county.Country, Guidelines: guidelines, CountyStatuses: countyStatuses}
			resultList = append(resultList, entity)
		}
	}
	return resultList, nil
}

//DeleteCounty deletes a county
func (sa *Adapter) DeleteCounty(ID string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. check if there are associated locations for this county
		locationsFilter := bson.D{primitive.E{Key: "county_id", Value: ID}}
		var locResult []*location
		err = sa.db.locations.FindWithContext(sessionContext, locationsFilter, &locResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(locResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated locations for this county")
		}

		//2. check if there are associated test types rules for this county
		rulesFilter := bson.D{primitive.E{Key: "county_id", Value: ID}}
		var rulesResult []*rule
		err = sa.db.rules.FindWithContext(sessionContext, rulesFilter, &rulesResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(rulesResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated test types rules for this county")
		}

		//3. check if there are associated symptom rules for this county
		symptomRulesFilter := bson.D{primitive.E{Key: "county_id", Value: ID}}
		var symptomRulesResult []*symptomRule
		err = sa.db.symptomrules.FindWithContext(sessionContext, symptomRulesFilter, &symptomRulesResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(symptomRulesResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated symptom rules for this county")
		}

		//4. check if there are associated access rules for this county
		accessRulesFilter := bson.D{primitive.E{Key: "county_id", Value: ID}}
		var accessRulesResult []*accessRule
		err = sa.db.accessrules.FindWithContext(sessionContext, accessRulesFilter, &accessRulesResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(accessRulesResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated access rules for this county")
		}

		//5. now we can delete the county
		deleteFilter := bson.D{primitive.E{Key: "_id", Value: ID}}
		result, err := sa.db.counties.DeleteOneWithContext(sessionContext, deleteFilter, nil)
		if err != nil {
			log.Printf("error deleting a county - %s", err)

			abortTransaction(sessionContext)

			return err
		}
		if result == nil {
			abortTransaction(sessionContext)
			return errors.New("result is nil for county with id " + ID)
		}
		deletedCount := result.DeletedCount
		if deletedCount == 0 {
			abortTransaction(sessionContext)
			return errors.New("there is no a county for id " + ID)
		}
		if deletedCount > 1 {
			abortTransaction(sessionContext)
			return errors.New("deleted more than one records for id " + ID)
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//CreateGuideline creates a guidline
func (sa *Adapter) CreateGuideline(countyID string, name string, description string, items []model.GuidelineItem) (*model.Guideline, error) {
	//1. find the county
	findFilter := bson.D{primitive.E{Key: "_id", Value: countyID}}
	var result []*county
	err := sa.db.counties.Find(findFilter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, errors.New("there is no a county for the provided id")
	}
	county := result[0]

	//2. create guideline
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	var gdItems []guidlineItem
	for _, v := range items {
		item := guidlineItem{Icon: v.Icon, Description: v.Description, Type: v.Type.Value}
		gdItems = append(gdItems, item)
	}
	dateCreated := time.Now()
	guideline := guidline{ID: id.String(), Name: name, Description: description, Items: gdItems, DateCreated: dateCreated}

	//3. add the guideline to the county
	guidelines := county.Guidelines
	guidelines = append(guidelines, guideline)
	county.Guidelines = guidelines

	//4. save the county
	saveFilter := bson.D{primitive.E{Key: "_id", Value: county.ID}}
	err = sa.db.counties.ReplaceOne(saveFilter, county, nil)
	if err != nil {
		return nil, err
	}

	//5. return the inserted item
	createdItem := &model.Guideline{ID: id.String(), Name: name, Description: description, Items: items}
	return createdItem, nil
}

//FindGuideline finds a guideline
func (sa *Adapter) FindGuideline(ID string) (*model.Guideline, error) {
	//1. first find the county
	filter := bson.D{primitive.E{Key: "guidelines.id", Value: ID}}
	var result []*county
	err := sa.db.counties.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	county := result[0]

	//2. get the guideline from the county
	var guideline guidline
	allGuidelines := county.Guidelines
	if allGuidelines != nil {
		for _, v := range allGuidelines {
			if v.ID == ID {
				guideline = v
				break
			}
		}
	}

	//3. construct the result
	var items []model.GuidelineItem
	for _, current := range guideline.Items {
		itemType := model.GuidelineItemType{Value: current.Type}
		gli := model.GuidelineItem{Icon: current.Icon, Description: current.Description, Type: itemType}
		items = append(items, gli)
	}
	resultItem := &model.Guideline{ID: guideline.ID, Name: guideline.Name, Description: guideline.Description, Items: items}

	return resultItem, nil
}

//FindGuidelineByCountyID finds guidelines for the provided county id
func (sa *Adapter) FindGuidelineByCountyID(countyID string) ([]*model.Guideline, error) {
	//1. first find the county
	filter := bson.D{primitive.E{Key: "_id", Value: countyID}}
	var result []*county
	err := sa.db.counties.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	county := result[0]

	//2. construct the result
	var resultList []*model.Guideline
	allGuidelines := county.Guidelines
	if allGuidelines != nil {
		for _, current := range allGuidelines {
			var items []model.GuidelineItem
			if current.Items != nil {
				for _, c := range current.Items {
					itemType := model.GuidelineItemType{Value: c.Type}
					items = append(items, model.GuidelineItem{Icon: c.Icon, Description: c.Description, Type: itemType})
				}
			}
			item := &model.Guideline{ID: current.ID, Name: current.Name, Description: current.Description, Items: items}
			resultList = append(resultList, item)
		}
	}

	return resultList, nil
}

//SaveGuideline saves a guideline
func (sa *Adapter) SaveGuideline(guideline *model.Guideline) error {
	//1. first find the county
	filter := bson.D{primitive.E{Key: "guidelines.id", Value: guideline.ID}}
	var result []*county
	err := sa.db.counties.Find(filter, &result, nil)
	if err != nil {
		return err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil
	}
	county := result[0]

	//2. update the guideline in the county
	guidelines := county.Guidelines
	var newGuidelines []guidline
	if guidelines != nil {
		for _, v := range guidelines {
			if v.ID == guideline.ID {
				//date updated
				dateUpdated := time.Now()
				v.DateUpdated = &dateUpdated

				//name
				v.Name = guideline.Name
				v.Description = guideline.Description

				//items
				var gdItems []guidlineItem
				if guideline.Items != nil {
					for _, v := range guideline.Items {
						item := guidlineItem{Icon: v.Icon, Description: v.Description, Type: v.Type.Value}
						gdItems = append(gdItems, item)
					}
				}
				v.Items = gdItems
			}
			newGuidelines = append(newGuidelines, v)
		}
	}
	county.Guidelines = newGuidelines

	//3. save the county
	saveFilter := bson.D{primitive.E{Key: "_id", Value: county.ID}}
	err = sa.db.counties.ReplaceOne(saveFilter, county, nil)
	if err != nil {
		return err
	}

	return nil
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

//DeleteGuideline deletes a guideline
func (sa *Adapter) DeleteGuideline(ID string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. first find the county
		filter := bson.D{primitive.E{Key: "guidelines.id", Value: ID}}
		var result []*county
		err = sa.db.counties.FindWithContext(sessionContext, filter, &result, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if result == nil || len(result) == 0 {
			abortTransaction(sessionContext)
			return errors.New("there is no a guideline for id " + ID)
		}
		county := result[0]

		//2. remove the guideline from the county
		indextToDelete := -1
		guidelines := county.Guidelines
		for index, v := range guidelines {
			if v.ID == ID {
				indextToDelete = index
				break
			}
		}
		guidelines = append(guidelines[:indextToDelete], guidelines[indextToDelete+1:]...)
		county.Guidelines = guidelines

		//3. save the county
		saveFilter := bson.D{primitive.E{Key: "_id", Value: county.ID}}
		err = sa.db.counties.ReplaceOneWithContext(sessionContext, saveFilter, county, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//CreateCountyStatus creates county status
func (sa *Adapter) CreateCountyStatus(countyID string, name string, description string) (*model.CountyStatus, error) {
	//1. find the county
	findFilter := bson.D{primitive.E{Key: "_id", Value: countyID}}
	var result []*county
	err := sa.db.counties.Find(findFilter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, errors.New("there is no a county for the provided id")
	}
	county := result[0]

	//2. create county status
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	dateCreated := time.Now()
	countyStatus := countyStatus{ID: id.String(), Name: name, Description: description, DateCreated: dateCreated}

	//3. add the county status to the county
	countyStatuses := county.CountyStatuses
	countyStatuses = append(countyStatuses, countyStatus)
	county.CountyStatuses = countyStatuses

	//4. save the county
	saveFilter := bson.D{primitive.E{Key: "_id", Value: county.ID}}
	err = sa.db.counties.ReplaceOne(saveFilter, county, nil)
	if err != nil {
		return nil, err
	}

	//5. return the inserted item
	createdItem := &model.CountyStatus{ID: id.String(), Name: name, Description: description}
	return createdItem, nil
}

//FindCountyStatus finds county status by ID
func (sa *Adapter) FindCountyStatus(ID string) (*model.CountyStatus, error) {
	//1. first find the county
	filter := bson.D{primitive.E{Key: "county_statuses.id", Value: ID}}
	var result []*county
	err := sa.db.counties.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	county := result[0]

	//2. get the county status from the county
	var countyStatus countyStatus
	allCountyStatuses := county.CountyStatuses
	if allCountyStatuses != nil {
		for _, v := range allCountyStatuses {
			if v.ID == ID {
				countyStatus = v
				break
			}
		}
	}

	//3. construct the result
	resultItem := &model.CountyStatus{ID: countyStatus.ID, Name: countyStatus.Name,
		Description: countyStatus.Description}

	return resultItem, nil
}

//FindCountyStatusesByCountyID finds county statuses by county ID
func (sa *Adapter) FindCountyStatusesByCountyID(countyID string) ([]*model.CountyStatus, error) {
	//1. first find the county
	filter := bson.D{primitive.E{Key: "_id", Value: countyID}}
	var result []*county
	err := sa.db.counties.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	county := result[0]

	//2. construct the result
	var resultList []*model.CountyStatus
	allCountyStatuses := county.CountyStatuses
	if allCountyStatuses != nil {
		for _, current := range allCountyStatuses {
			item := &model.CountyStatus{ID: current.ID, Name: current.Name, Description: current.Description}
			resultList = append(resultList, item)
		}
	}

	return resultList, nil
}

//SaveCountyStatus saves the county status
func (sa *Adapter) SaveCountyStatus(entity *model.CountyStatus) error {
	//1. first find the county
	filter := bson.D{primitive.E{Key: "county_statuses.id", Value: entity.ID}}
	var result []*county
	err := sa.db.counties.Find(filter, &result, nil)
	if err != nil {
		return err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil
	}
	county := result[0]

	//2. update the county status in the county
	countyStatuses := county.CountyStatuses
	var newCountyStatuses []countyStatus
	if countyStatuses != nil {
		for _, v := range countyStatuses {
			if v.ID == entity.ID {
				//date updated
				dateUpdated := time.Now()
				v.DateUpdated = &dateUpdated

				v.Name = entity.Name
				v.Description = entity.Description
			}
			newCountyStatuses = append(newCountyStatuses, v)
		}
	}
	county.CountyStatuses = newCountyStatuses

	//3. save the county
	saveFilter := bson.D{primitive.E{Key: "_id", Value: county.ID}}
	err = sa.db.counties.ReplaceOne(saveFilter, county, nil)
	if err != nil {
		return err
	}

	return nil
}

//DeleteCountyStatus deletes county status
func (sa *Adapter) DeleteCountyStatus(ID string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. check if there are associated rules for this county status
		rulesFilter := bson.D{primitive.E{Key: "results_statuses.county_status_id", Value: ID}}
		var rulesResult []*rule
		err = sa.db.rules.FindWithContext(sessionContext, rulesFilter, &rulesResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(rulesResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated test type rules for this county status")
		}

		//2. check if there are associated symptom rules for this county status
		sRulesFilter := bson.D{primitive.E{Key: "items.county_status_id", Value: ID}}
		var sRulesResult []*symptomRule
		err = sa.db.symptomrules.FindWithContext(sessionContext, sRulesFilter, &sRulesResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(sRulesResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated symptom rules for this county status")
		}

		//3. check if there are associated access rules for this county status
		aRulesFilter := bson.D{primitive.E{Key: "rules.county_status_id", Value: ID}}
		var aRulesResult []*accessRule
		err = sa.db.accessrules.FindWithContext(sessionContext, aRulesFilter, &aRulesResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(aRulesResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated access rules for this county status")
		}

		//4. first find the county
		filter := bson.D{primitive.E{Key: "county_statuses.id", Value: ID}}
		var result []*county
		err = sa.db.counties.FindWithContext(sessionContext, filter, &result, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if result == nil || len(result) == 0 {
			abortTransaction(sessionContext)
			return errors.New("there is no a county status for id " + ID)
		}
		county := result[0]

		//5. remove the county status from the county
		indextToDelete := -1
		countyStatuses := county.CountyStatuses
		for index, v := range countyStatuses {
			if v.ID == ID {
				indextToDelete = index
				break
			}
		}
		countyStatuses = append(countyStatuses[:indextToDelete], countyStatuses[indextToDelete+1:]...)
		county.CountyStatuses = countyStatuses

		//6. save the county
		saveFilter := bson.D{primitive.E{Key: "_id", Value: county.ID}}
		err = sa.db.counties.ReplaceOneWithContext(sessionContext, saveFilter, county, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//ReadAllTestTypes reads all test types
func (sa *Adapter) ReadAllTestTypes() ([]*model.TestType, error) {
	filter := bson.D{}
	var result []*testType
	err := sa.db.testtypes.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	var resultList []*model.TestType
	if result != nil {
		for _, current := range result {
			var ttResults []model.TestTypeResult
			if current.Results != nil {
				for _, inner := range current.Results {
					ttResult := model.TestTypeResult{ID: inner.ID, Name: inner.Name, NextStep: inner.NextStep,
						NextStepOffset: inner.NextStepOffset, ResultExpiresOffset: inner.ResultExpiresOffset}
					ttResults = append(ttResults, ttResult)
				}
			}

			item := &model.TestType{ID: current.ID, Name: current.Name, Priority: current.Priority, Results: ttResults}
			resultList = append(resultList, item)
		}
	}
	return resultList, nil
}

//CreateTestType creates a test type
func (sa *Adapter) CreateTestType(name string, priority *int) (*model.TestType, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	dateCreated := time.Now()

	testType := testType{ID: id.String(), Name: name, Priority: priority, DateCreated: dateCreated}
	_, err = sa.db.testtypes.InsertOne(&testType)
	if err != nil {
		return nil, err
	}

	//return the inserted item
	result := &model.TestType{ID: testType.ID, Name: testType.Name, Priority: testType.Priority}
	return result, nil
}

//FindTestType finds a test type by ID
func (sa *Adapter) FindTestType(ID string) (*model.TestType, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*testType
	err := sa.db.testtypes.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	testType := result[0]

	var ttResults []model.TestTypeResult
	if testType.Results != nil {
		for _, ttr := range testType.Results {
			testTypeResult := model.TestTypeResult{ID: ttr.ID, Name: ttr.Name, NextStep: ttr.NextStep,
				NextStepOffset: ttr.NextStepOffset, ResultExpiresOffset: ttr.ResultExpiresOffset}
			ttResults = append(ttResults, testTypeResult)
		}
	}
	resultEntity := &model.TestType{ID: testType.ID, Name: testType.Name, Priority: testType.Priority,
		Results: ttResults}
	return resultEntity, nil
}

//FindTestTypesByIDs finds the test types for the provided ids
func (sa *Adapter) FindTestTypesByIDs(ids []string) ([]*model.TestType, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}

	var result []*testType
	err := sa.db.testtypes.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}

	var resultList []*model.TestType
	for _, testType := range result {
		var ttResults []model.TestTypeResult
		if testType.Results != nil {
			for _, ttr := range testType.Results {
				testTypeResult := model.TestTypeResult{ID: ttr.ID, Name: ttr.Name, NextStep: ttr.NextStep,
					NextStepOffset: ttr.NextStepOffset, ResultExpiresOffset: ttr.ResultExpiresOffset}
				ttResults = append(ttResults, testTypeResult)
			}
		}
		resultEntity := &model.TestType{ID: testType.ID, Name: testType.Name, Priority: testType.Priority,
			Results: ttResults}
		resultList = append(resultList, resultEntity)
	}
	return resultList, nil
}

//SaveTestType saves the test type
func (sa *Adapter) SaveTestType(entity *model.TestType) error {
	findFilter := bson.D{primitive.E{Key: "_id", Value: entity.ID}}
	var result []*testType
	err := sa.db.testtypes.Find(findFilter, &result, nil)
	if err != nil {
		return err
	}
	if result == nil || len(result) == 0 {
		//not found
		return errors.New("there is no a test type for the provided id")
	}
	testType := result[0]

	//update the values
	testType.Name = entity.Name
	testType.Priority = entity.Priority
	dateUpdated := time.Now()
	testType.DateUpdated = &dateUpdated

	//save
	filter := bson.D{primitive.E{Key: "_id", Value: testType.ID}}
	err = sa.db.testtypes.ReplaceOne(filter, testType, nil)
	if err != nil {
		return err
	}

	return nil
}

//DeleteTestType delete a test type
func (sa *Adapter) DeleteTestType(ID string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. check if there are associated locations for this test type
		locationsFilter := bson.M{"available_tests": ID}
		var locResult []*location
		err = sa.db.locations.FindWithContext(sessionContext, locationsFilter, &locResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(locResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated locations for this test type")
		}

		//2. check if there are associated county rules for this test type
		rulesFilter := bson.D{primitive.E{Key: "test_type_id", Value: ID}}
		var rulesResult []*rule
		err = sa.db.rules.FindWithContext(sessionContext, rulesFilter, &rulesResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(rulesResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated county rules for this test type")
		}

		//3. check if there are associated test type results for this test type
		ttrFilter := bson.M{"_id": ID}
		var ttResult []*testType
		err = sa.db.testtypes.FindWithContext(sessionContext, ttrFilter, &ttResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(ttResult) > 0 {
			testType := ttResult[0]
			if len(testType.Results) > 0 {
				abortTransaction(sessionContext)
				return errors.New("there are associated test type results for this test type")
			}
		}

		//4. now we can delete the provider
		deleteFilter := bson.D{primitive.E{Key: "_id", Value: ID}}
		result, err := sa.db.testtypes.DeleteOneWithContext(sessionContext, deleteFilter, nil)
		if err != nil {
			log.Printf("error deleting a test type - %s", err)
			abortTransaction(sessionContext)
			return err
		}
		if result == nil {
			abortTransaction(sessionContext)
			return errors.New("result is nil for test type with id " + ID)
		}
		deletedCount := result.DeletedCount
		if deletedCount == 0 {
			abortTransaction(sessionContext)
			return errors.New("there is no a test type for id " + ID)
		}
		if deletedCount > 1 {
			abortTransaction(sessionContext)
			return errors.New("deleted more than one records for id " + ID)
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//CreateTestTypeResult creates a test type result
func (sa *Adapter) CreateTestTypeResult(testTypeID string, name string, nextStep string, nextStepOffset *int,
	resultExpiresOffset *int) (*model.TestTypeResult, error) {

	//1. find the test type
	findFilter := bson.D{primitive.E{Key: "_id", Value: testTypeID}}
	var result []*testType
	err := sa.db.testtypes.Find(findFilter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, errors.New("there is no a test type for the provided id")
	}
	testType := result[0]

	//2. create test type result
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	dateCreated := time.Now()
	testTypeResult := testTypeResult{ID: id.String(), Name: name, NextStep: nextStep,
		NextStepOffset: nextStepOffset, ResultExpiresOffset: resultExpiresOffset, DateCreated: dateCreated}

	//3. add the test type result to the test type
	results := testType.Results
	results = append(results, testTypeResult)
	testType.Results = results

	//4. save the test type
	saveFilter := bson.D{primitive.E{Key: "_id", Value: testType.ID}}
	err = sa.db.testtypes.ReplaceOne(saveFilter, testType, nil)
	if err != nil {
		return nil, err
	}

	//5. return the inserted item
	createdItem := &model.TestTypeResult{ID: testTypeResult.ID, Name: testTypeResult.Name, NextStep: testTypeResult.NextStep,
		NextStepOffset: testTypeResult.NextStepOffset, ResultExpiresOffset: testTypeResult.ResultExpiresOffset}
	return createdItem, nil
}

//FindTestTypeResult finds a test type result by id
func (sa *Adapter) FindTestTypeResult(ID string) (*model.TestTypeResult, error) {
	//1. first find the test type
	filter := bson.D{primitive.E{Key: "results._id", Value: ID}}
	var result []*testType
	err := sa.db.testtypes.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	testType := result[0]

	//2. get the test type results from the test type
	var testTypeResult testTypeResult
	allResults := testType.Results
	if allResults != nil {
		for _, v := range allResults {
			if v.ID == ID {
				testTypeResult = v
				break
			}
		}
	}

	//3. construct the result
	resultItem := &model.TestTypeResult{ID: testTypeResult.ID, Name: testTypeResult.Name, NextStep: testTypeResult.NextStep,
		NextStepOffset: testTypeResult.NextStepOffset, ResultExpiresOffset: testTypeResult.ResultExpiresOffset}

	return resultItem, nil
}

//FindTestTypeResultsByTestTypeID finds all test type results for a test type
func (sa *Adapter) FindTestTypeResultsByTestTypeID(testTypeID string) ([]*model.TestTypeResult, error) {
	//1. first find the test type
	filter := bson.D{primitive.E{Key: "_id", Value: testTypeID}}
	var result []*testType
	err := sa.db.testtypes.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	testType := result[0]

	//2. construct the result
	var resultList []*model.TestTypeResult
	allResults := testType.Results
	if allResults != nil {
		for _, current := range allResults {
			item := &model.TestTypeResult{ID: current.ID, Name: current.Name, NextStep: current.NextStep,
				NextStepOffset: current.NextStepOffset, ResultExpiresOffset: current.ResultExpiresOffset}
			resultList = append(resultList, item)
		}
	}

	return resultList, nil
}

//SaveTestTypeResult save the test type result
func (sa *Adapter) SaveTestTypeResult(entity *model.TestTypeResult) error {
	//1. first find the test type
	filter := bson.D{primitive.E{Key: "results._id", Value: entity.ID}}
	var result []*testType
	err := sa.db.testtypes.Find(filter, &result, nil)
	if err != nil {
		return err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil
	}
	testType := result[0]

	//2. update the test type result in the test type
	results := testType.Results
	var newResults []testTypeResult
	if results != nil {
		for _, v := range results {
			if v.ID == entity.ID {
				//date updated
				dateUpdated := time.Now()
				v.DateUpdated = &dateUpdated

				//name
				v.Name = entity.Name
				v.NextStep = entity.NextStep
				v.NextStepOffset = entity.NextStepOffset
				v.ResultExpiresOffset = entity.ResultExpiresOffset
			}
			newResults = append(newResults, v)
		}
	}
	testType.Results = newResults

	//3. save the test type
	saveFilter := bson.D{primitive.E{Key: "_id", Value: testType.ID}}
	err = sa.db.testtypes.ReplaceOne(saveFilter, testType, nil)
	if err != nil {
		return err
	}

	return nil
}

//DeleteTestTypeResult deletes a test type result
func (sa *Adapter) DeleteTestTypeResult(ID string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. check if there are associated county rules for this test type result
		rulesFilter := bson.D{primitive.E{Key: "results_statuses.test_type_result_id", Value: ID}}
		var rulesResult []*rule
		err = sa.db.rules.FindWithContext(sessionContext, rulesFilter, &rulesResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(rulesResult) > 0 {
			abortTransaction(sessionContext)
			return errors.New("there are associated county rules for this test type result")
		}

		//2. find the test type
		filter := bson.D{primitive.E{Key: "results._id", Value: ID}}
		var result []*testType
		err = sa.db.testtypes.FindWithContext(sessionContext, filter, &result, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if result == nil || len(result) == 0 {
			abortTransaction(sessionContext)
			return errors.New("there is no a test type result for id " + ID)
		}
		testType := result[0]

		//3. remove the test type result from the test type
		indextToDelete := -1
		results := testType.Results
		for index, v := range results {
			if v.ID == ID {
				indextToDelete = index
				break
			}
		}
		results = append(results[:indextToDelete], results[indextToDelete+1:]...)
		testType.Results = results

		//4. save the test type
		saveFilter := bson.D{primitive.E{Key: "_id", Value: testType.ID}}
		err = sa.db.testtypes.ReplaceOneWithContext(sessionContext, saveFilter, testType, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//ReadAllRules reads all rules
func (sa *Adapter) ReadAllRules() ([]*model.Rule, error) {
	filter := bson.D{}
	var result []*rule
	err := sa.db.rules.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	var resultList []*model.Rule
	if result != nil {
		for _, rule := range result {
			county := model.County{ID: rule.CountyID}
			testType := model.TestType{ID: rule.TestTypeID}
			var resultsStatuses []model.TestTypeResultCountyStatus
			if rule.ResultsStates != nil {
				for _, current := range rule.ResultsStates {
					item := model.TestTypeResultCountyStatus{TestTypeResultID: current.TestTypeResultID, CountyStatusID: current.CountyStatusID}
					resultsStatuses = append(resultsStatuses, item)
				}
			}
			resultItem := &model.Rule{ID: rule.ID, County: county, TestType: testType,
				Priority: rule.Priority, ResultsStates: resultsStatuses}
			resultList = append(resultList, resultItem)
		}
	}
	return resultList, nil
}

//FindRulesByCountyID finds the rules for a county
func (sa *Adapter) FindRulesByCountyID(countyID string) ([]*model.Rule, error) {
	filter := bson.D{primitive.E{Key: "county_id", Value: countyID}}
	var result []*rule
	err := sa.db.rules.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}

	//construct the result
	var resultList []*model.Rule
	for _, rule := range result {
		county := model.County{ID: rule.CountyID}
		testType := model.TestType{ID: rule.TestTypeID}
		var resultsStatuses []model.TestTypeResultCountyStatus
		if rule.ResultsStates != nil {
			for _, current := range rule.ResultsStates {
				item := model.TestTypeResultCountyStatus{TestTypeResultID: current.TestTypeResultID, CountyStatusID: current.CountyStatusID}
				resultsStatuses = append(resultsStatuses, item)
			}
		}
		resultItem := &model.Rule{ID: rule.ID, County: county, TestType: testType,
			Priority: rule.Priority, ResultsStates: resultsStatuses}
		resultList = append(resultList, resultItem)
	}
	return resultList, nil
}

//FindRule finds a rule
func (sa *Adapter) FindRule(ID string) (*model.Rule, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*rule
	err := sa.db.rules.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	rule := result[0]

	county := model.County{ID: rule.CountyID}
	testType := model.TestType{ID: rule.TestTypeID}
	var resultsStatuses []model.TestTypeResultCountyStatus
	if rule.ResultsStates != nil {
		for _, current := range rule.ResultsStates {
			item := model.TestTypeResultCountyStatus{TestTypeResultID: current.TestTypeResultID, CountyStatusID: current.CountyStatusID}
			resultsStatuses = append(resultsStatuses, item)
		}
	}
	resultItem := &model.Rule{ID: rule.ID, County: county, TestType: testType,
		Priority: rule.Priority, ResultsStates: resultsStatuses}
	return resultItem, nil
}

//FindRuleByCountyIDTestTypeID finds the rule for a county and test type
func (sa *Adapter) FindRuleByCountyIDTestTypeID(countyID string, testTypeID string) (*model.Rule, error) {
	filter := bson.D{primitive.E{Key: "county_id", Value: countyID}, primitive.E{Key: "test_type_id", Value: testTypeID}}
	var result []rule
	err := sa.db.rules.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	rule := result[0]

	county := model.County{ID: rule.CountyID}
	testType := model.TestType{ID: rule.TestTypeID}
	var resultsStatuses []model.TestTypeResultCountyStatus
	if rule.ResultsStates != nil {
		for _, current := range rule.ResultsStates {
			item := model.TestTypeResultCountyStatus{TestTypeResultID: current.TestTypeResultID, CountyStatusID: current.CountyStatusID}
			resultsStatuses = append(resultsStatuses, item)
		}
	}
	resultItem := model.Rule{ID: rule.ID, County: county, TestType: testType,
		Priority: rule.Priority, ResultsStates: resultsStatuses}
	return &resultItem, nil
}

//CreateRule create a rule
func (sa *Adapter) CreateRule(countyID string, testTypeID string, priority *int, resultsStates []model.TestTypeResultCountyStatus) (*model.Rule, error) {

	var resultItem model.Rule

	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. create rule
		id, err := uuid.NewUUID()
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		dateCreated := time.Now()

		var resultsStatuses []testTypeResultCountyStatus
		if resultsStates != nil {
			for _, current := range resultsStates {
				item := testTypeResultCountyStatus{TestTypeResultID: current.TestTypeResultID, CountyStatusID: current.CountyStatusID}
				resultsStatuses = append(resultsStatuses, item)
			}
		}
		countyTestType := rule{ID: id.String(), CountyID: countyID,
			TestTypeID: testTypeID, Priority: priority, ResultsStates: resultsStatuses, DateCreated: dateCreated}
		_, err = sa.db.rules.InsertOneWithContext(sessionContext, &countyTestType)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		county := model.County{ID: countyTestType.CountyID}
		testType := model.TestType{ID: countyTestType.TestTypeID}
		resultItem = model.Rule{ID: countyTestType.ID, County: county, TestType: testType,
			Priority: countyTestType.Priority, ResultsStates: resultsStates}

		//2. delete all statuses of the users
		err = sa.deleteAllStatuses(sessionContext)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &resultItem, nil
}

//SaveRule saves a rule
func (sa *Adapter) SaveRule(entity *model.Rule) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. save
		findFilter := bson.D{primitive.E{Key: "_id", Value: entity.ID}}
		var result []*rule
		err = sa.db.rules.FindWithContext(sessionContext, findFilter, &result, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if result == nil || len(result) == 0 {
			//not found
			abortTransaction(sessionContext)
			return errors.New("there is no a rule for the provided id")
		}
		rule := result[0]

		//update the values
		rule.Priority = entity.Priority

		var resSt []testTypeResultCountyStatus
		if entity.ResultsStates != nil {
			for _, current := range entity.ResultsStates {
				item := testTypeResultCountyStatus{TestTypeResultID: current.TestTypeResultID, CountyStatusID: current.CountyStatusID}
				resSt = append(resSt, item)
			}
		}
		rule.ResultsStates = resSt

		dateUpdated := time.Now()
		rule.DateUpdated = &dateUpdated

		filter := bson.D{primitive.E{Key: "_id", Value: rule.ID}}
		err = sa.db.rules.ReplaceOneWithContext(sessionContext, filter, rule, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//2. delete all statuses of the users
		err = sa.deleteAllStatuses(sessionContext)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//DeleteRule deletes a rule
func (sa *Adapter) DeleteRule(ID string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. delete
		filter := bson.D{primitive.E{Key: "_id", Value: ID}}
		result, err := sa.db.rules.DeleteOneWithContext(sessionContext, filter, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if result == nil {
			abortTransaction(sessionContext)
			return errors.New("result is nil for rule item with id " + ID)
		}
		deletedCount := result.DeletedCount
		if deletedCount == 0 {
			abortTransaction(sessionContext)
			return errors.New("there is no a rule for id " + ID)
		}
		if deletedCount > 1 {
			abortTransaction(sessionContext)
			return errors.New("deleted more than one records for id " + ID)
		}

		//2. delete all statuses of the users
		err = sa.deleteAllStatuses(sessionContext)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//ReadAllLocations reads all the locations
func (sa *Adapter) ReadAllLocations() ([]*model.Location, error) {
	filter := bson.D{}
	var result []*location
	err := sa.db.locations.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	var resultList []*model.Location
	if result != nil {
		for _, location := range result {
			provider := model.Provider{ID: location.ProviderID}
			county := model.County{ID: location.CountyID}
			var avTests []model.TestType
			if location.AvailableTests != nil {
				for _, id := range location.AvailableTests {
					testType := model.TestType{ID: id}
					avTests = append(avTests, testType)
				}
			}
			daysOfOperation := convertToDaysOfOperation(location.DaysOfOperation)
			locationEntity := &model.Location{ID: location.ID, Name: location.Name, Address1: location.Address1,
				Address2: location.Address2, City: location.City, State: location.State, ZIP: location.ZIP, Country: location.Country,
				Latitude: location.Latitude, Longitude: location.Longitude, Contact: location.Contact, Timezone: location.Timezone,
				DaysOfOperation: daysOfOperation, URL: location.URL, Notes: location.Notes, WaitTimeColor: location.WaitTimeColor,
				Provider: provider, County: county, AvailableTests: avTests}
			resultList = append(resultList, locationEntity)
		}
	}
	return resultList, nil
}

//CreateLocation creates a location
func (sa *Adapter) CreateLocation(providerID string, countyID string, name string, address1 string, address2 string, city string,
	state string, zip string, country string, latitude float64, longitude float64, contact string, daysOfOperation []model.OperationDay,
	url string, notes string, waitTimeColor *string, availableTests []string) (*model.Location, error) {

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	dateCreated := time.Now()

	doo := convertFromDaysOfOperation(daysOfOperation)
	location := location{ID: id.String(), Name: name, Address1: address1, Address2: address2, City: city,
		State: state, ZIP: zip, Country: country, Latitude: latitude, Longitude: longitude, Timezone: "America/Chicago", Contact: contact,
		DaysOfOperation: doo, URL: url, Notes: notes, WaitTimeColor: waitTimeColor, ProviderID: providerID, CountyID: countyID,
		AvailableTests: availableTests, DateCreated: dateCreated}
	_, err = sa.db.locations.InsertOne(&location)
	if err != nil {
		return nil, err
	}

	//return the inserted item
	provider := model.Provider{ID: providerID}
	county := model.County{ID: countyID}
	var avTests []model.TestType
	if location.AvailableTests != nil {
		for _, id := range location.AvailableTests {
			testType := model.TestType{ID: id}
			avTests = append(avTests, testType)
		}
	}
	result := &model.Location{ID: id.String(), Name: name, Address1: address1, Address2: address2, City: city,
		State: state, ZIP: zip, Country: country, Latitude: latitude, Longitude: longitude, Contact: contact,
		DaysOfOperation: daysOfOperation, URL: url, Notes: notes, WaitTimeColor: waitTimeColor, Provider: provider, County: county, AvailableTests: avTests}
	return result, nil
}

//FindLocationsByProviderIDCountyID finds the locations for a provider and county
func (sa *Adapter) FindLocationsByProviderIDCountyID(providerID string, countyID string) ([]*model.Location, error) {
	filter := bson.D{primitive.E{Key: "provider_id", Value: providerID},
		primitive.E{Key: "county_id", Value: countyID}}
	var result []*location
	err := sa.db.locations.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}

	var resultList []*model.Location
	for _, location := range result {
		provider := model.Provider{ID: location.ProviderID}
		county := model.County{ID: countyID}
		var avTests []model.TestType
		if location.AvailableTests != nil {
			for _, id := range location.AvailableTests {
				testType := model.TestType{ID: id}
				avTests = append(avTests, testType)
			}
		}
		daysOfOperations := convertToDaysOfOperation(location.DaysOfOperation)
		locationEntity := &model.Location{ID: location.ID, Name: location.Name, Address1: location.Address1,
			Address2: location.Address2, City: location.City, State: location.State, ZIP: location.ZIP,
			Country: location.Country, Latitude: location.Latitude, Longitude: location.Longitude, Timezone: location.Timezone,
			Contact: location.Contact, DaysOfOperation: daysOfOperations, URL: location.URL, Notes: location.Notes,
			WaitTimeColor: location.WaitTimeColor, Provider: provider, County: county, AvailableTests: avTests}
		resultList = append(resultList, locationEntity)
	}
	return resultList, nil
}

type locationProviderJoin struct {
	ID              string         `bson:"_id"`
	Name            string         `bson:"name"`
	Address1        string         `bson:"address_1"`
	Address2        string         `bson:"address_2"`
	City            string         `bson:"city"`
	State           string         `bson:"state"`
	ZIP             string         `bson:"zip"`
	Country         string         `bson:"country"`
	Latitude        float64        `bson:"latitude"`
	Longitude       float64        `bson:"longitude"`
	Timezone        string         `bson:"timezone"`
	Contact         string         `bson:"contact"`
	DaysOfOperation []operationDay `bson:"days_of_operation"`
	URL             string         `bson:"url"`
	Notes           string         `bson:"notes"`
	WaitTimeColor   *string        `bson:"wait_time_color"`
	AvailableTests  []string       `bson:"available_tests"`
	CountyID        string         `bson:"county_id"`

	ProviderID                  string   `bson:"provider_id"`
	ProviderName                string   `bson:"provider_name"`
	ProviderAvailableMechanisms []string `bson:"provider_available_mechanisms"`
}

//FindLocationsByCountyIDDeep finds the locations for a county - deep request!
func (sa *Adapter) FindLocationsByCountyIDDeep(countyID string) ([]*model.Location, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "providers",
			"localField":   "provider_id",
			"foreignField": "_id",
			"as":           "provider",
		}},
		{"$match": bson.M{"county_id": countyID}},
		{"$unwind": "$provider"},
		{"$project": bson.M{
			"_id": 1, "name": 1, "address_1": 1, "address_2": 1, "city": 1, "state": 1, "zip": 1, "country": 1, "latitude": 1, "longitude": 1, "timezone": 1,
			"contact": 1, "days_of_operation": 1, "url": 1, "notes": 1, "wait_time_color": 1, "available_tests": 1, "county_id": 1,
			"provider_id": "$provider._id", "provider_name": "$provider.provider_name", "provider_available_mechanisms": "$provider.available_mechanisms",
		}}}

	var result []*locationProviderJoin
	err := sa.db.locations.Aggregate(pipeline, &result, nil)
	if err != nil {
		return nil, err

	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}

	var resultList []*model.Location
	for _, location := range result {
		provider := model.Provider{ID: location.ProviderID, Name: location.ProviderName,
			AvailableMechanisms: location.ProviderAvailableMechanisms}
		county := model.County{ID: location.CountyID}
		var avTests []model.TestType
		if location.AvailableTests != nil {
			for _, id := range location.AvailableTests {
				testType := model.TestType{ID: id}
				avTests = append(avTests, testType)
			}
		}
		daysOfOperations := convertToDaysOfOperation(location.DaysOfOperation)
		locationEntity := &model.Location{ID: location.ID, Name: location.Name, Address1: location.Address1, Address2: location.Address2,
			City: location.City, State: location.State, ZIP: location.ZIP, Country: location.Country, Latitude: location.Latitude, Timezone: location.Timezone,
			Longitude: location.Longitude, Contact: location.Contact, DaysOfOperation: daysOfOperations, URL: location.URL,
			Notes: location.Notes, WaitTimeColor: location.WaitTimeColor, Provider: provider, County: county, AvailableTests: avTests}
		resultList = append(resultList, locationEntity)
	}
	return resultList, nil
}

//FindLocationsByCountiesDeep finds the locations for a list of county items - deep request!
func (sa *Adapter) FindLocationsByCountiesDeep(countyIDs []string) ([]*model.Location, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "providers",
			"localField":   "provider_id",
			"foreignField": "_id",
			"as":           "provider",
		}},
		{"$match": bson.M{"county_id": bson.M{"$in": countyIDs}}},
		{"$unwind": "$provider"},
		{"$project": bson.M{
			"_id": 1, "name": 1, "address_1": 1, "address_2": 1, "city": 1, "state": 1, "zip": 1, "country": 1, "latitude": 1, "longitude": 1, "timezone": 1,
			"contact": 1, "days_of_operation": 1, "url": 1, "notes": 1, "wait_time_color": 1, "available_tests": 1, "county_id": 1,
			"provider_id": "$provider._id", "provider_name": "$provider.provider_name", "provider_available_mechanisms": "$provider.available_mechanisms",
		}}}

	var result []*locationProviderJoin
	err := sa.db.locations.Aggregate(pipeline, &result, nil)
	if err != nil {
		return nil, err

	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}

	var resultList []*model.Location
	for _, location := range result {
		provider := model.Provider{ID: location.ProviderID, Name: location.ProviderName,
			AvailableMechanisms: location.ProviderAvailableMechanisms}
		county := model.County{ID: location.CountyID}
		var avTests []model.TestType
		if location.AvailableTests != nil {
			for _, id := range location.AvailableTests {
				testType := model.TestType{ID: id}
				avTests = append(avTests, testType)
			}
		}
		daysOfOperations := convertToDaysOfOperation(location.DaysOfOperation)
		locationEntity := &model.Location{ID: location.ID, Name: location.Name, Address1: location.Address1,
			Address2: location.Address2, City: location.City, State: location.State, ZIP: location.ZIP,
			Country: location.Country, Latitude: location.Latitude, Longitude: location.Longitude, Timezone: location.Timezone,
			Contact: location.Contact, DaysOfOperation: daysOfOperations, URL: location.URL, Notes: location.Notes, WaitTimeColor: location.WaitTimeColor,
			Provider: provider, County: county, AvailableTests: avTests}
		resultList = append(resultList, locationEntity)
	}
	return resultList, nil
}

//FindLocation finds a location by id
func (sa *Adapter) FindLocation(ID string) (*model.Location, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*location
	err := sa.db.locations.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	location := result[0]

	provider := model.Provider{ID: location.ProviderID}
	county := model.County{ID: location.CountyID}
	var avTests []model.TestType
	if location.AvailableTests != nil {
		for _, id := range location.AvailableTests {
			testType := model.TestType{ID: id}
			avTests = append(avTests, testType)
		}
	}
	daysOfOperations := convertToDaysOfOperation(location.DaysOfOperation)
	resultEntity := &model.Location{ID: location.ID, Name: location.Name, Address1: location.Address1,
		Address2: location.Address2, City: location.City, State: location.State, ZIP: location.ZIP,
		Country: location.Country, Latitude: location.Latitude, Longitude: location.Longitude, Timezone: location.Timezone,
		Contact: location.Contact, DaysOfOperation: daysOfOperations, URL: location.URL, Notes: location.Notes,
		WaitTimeColor: location.WaitTimeColor, Provider: provider, County: county, AvailableTests: avTests}
	return resultEntity, nil
}

//SaveLocation save a location
func (sa *Adapter) SaveLocation(entity *model.Location) error {
	findFilter := bson.D{primitive.E{Key: "_id", Value: entity.ID}}
	var result []*location
	err := sa.db.locations.Find(findFilter, &result, nil)
	if err != nil {
		return err
	}
	if result == nil || len(result) == 0 {
		//not found
		return errors.New("there is no a location for the provided id")
	}
	location := result[0]

	//update the values
	location.Name = entity.Name
	location.Address1 = entity.Address1
	location.Address2 = entity.Address2
	location.City = entity.City
	location.State = entity.State
	location.ZIP = entity.ZIP
	location.Country = entity.Country
	location.Latitude = entity.Latitude
	location.Longitude = entity.Longitude
	location.Timezone = entity.Timezone
	location.Contact = entity.Contact
	location.DaysOfOperation = convertFromDaysOfOperation(entity.DaysOfOperation)
	location.URL = entity.URL
	location.Notes = entity.Notes
	location.WaitTimeColor = entity.WaitTimeColor
	var avTR []string
	if entity.AvailableTests != nil {
		for _, testType := range entity.AvailableTests {
			avTR = append(avTR, testType.ID)
		}
	}
	location.AvailableTests = avTR

	dateUpdated := time.Now()
	location.DateUpdated = &dateUpdated

	filter := bson.D{primitive.E{Key: "_id", Value: location.ID}}
	err = sa.db.locations.ReplaceOne(filter, location, nil)
	if err != nil {
		return err
	}

	return nil
}

//DeleteLocation deletes a location
func (sa *Adapter) DeleteLocation(ID string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	result, err := sa.db.locations.DeleteOne(filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return errors.New("result is nil for location item with id " + ID)
	}
	deletedCount := result.DeletedCount
	if deletedCount == 0 {
		return errors.New("there is no a location for id " + ID)
	}
	if deletedCount > 1 {
		return errors.New("deleted more than one records for id " + ID)
	}

	//success - count = 1
	return nil
}

//FindSymptom finds a symptom
func (sa *Adapter) FindSymptom(ID string) (*model.Symptom, error) {
	//1. first find the symptom group
	filter := bson.D{primitive.E{Key: "symptoms.id", Value: ID}}
	var result []*symptomGroup
	err := sa.db.symptomgroups.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	symptomGroup := result[0]

	//2. get the symptom from the symptom group
	var symptom symptom
	allSymptoms := symptomGroup.Symptoms
	if allSymptoms != nil {
		for _, v := range allSymptoms {
			if v.ID == ID {
				symptom = v
				break
			}
		}
	}

	//3. construct the result
	resultItem := &model.Symptom{ID: symptom.ID, Name: symptom.Name}

	return resultItem, nil
}

//CreateSymptom creates a symptom
func (sa *Adapter) CreateSymptom(name string, group string) (*model.Symptom, error) {
	//1. find the symptom group
	findFilter := bson.D{primitive.E{Key: "name", Value: group}}
	var result []*symptomGroup
	err := sa.db.symptomgroups.Find(findFilter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, errors.New("there is no a symptom group for the provided name")
	}
	symptomGroup := result[0]

	//2. create a symptom
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	dateCreated := time.Now()
	symptom := symptom{ID: id.String(), Name: name, DateCreated: dateCreated}

	//3. add the symptom to the symtom group
	symptoms := symptomGroup.Symptoms
	symptoms = append(symptoms, symptom)
	symptomGroup.Symptoms = symptoms

	//4. save the county
	saveFilter := bson.D{primitive.E{Key: "_id", Value: symptomGroup.ID}}
	err = sa.db.symptomgroups.ReplaceOne(saveFilter, symptomGroup, nil)
	if err != nil {
		return nil, err
	}

	//5. return the inserted item
	createdItem := &model.Symptom{ID: id.String(), Name: name}
	return createdItem, nil
}

//DeleteSymptom deletes a symptom
func (sa *Adapter) DeleteSymptom(ID string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. find the symptom group
		filter := bson.D{primitive.E{Key: "symptoms.id", Value: ID}}
		var result []*symptomGroup
		err = sa.db.symptomgroups.FindWithContext(sessionContext, filter, &result, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if result == nil || len(result) == 0 {
			abortTransaction(sessionContext)
			return errors.New("there is no a symptom for id " + ID)
		}
		symptomGroup := result[0]

		//2. remove the symptom from the symptom group
		indextToDelete := -1
		symptoms := symptomGroup.Symptoms
		for index, v := range symptoms {
			if v.ID == ID {
				indextToDelete = index
				break
			}
		}
		symptoms = append(symptoms[:indextToDelete], symptoms[indextToDelete+1:]...)
		symptomGroup.Symptoms = symptoms

		//3. save the symptom group
		saveFilter := bson.D{primitive.E{Key: "_id", Value: symptomGroup.ID}}
		err = sa.db.symptomgroups.ReplaceOneWithContext(sessionContext, saveFilter, symptomGroup, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//SaveSymptom saves a symptom
func (sa *Adapter) SaveSymptom(entity *model.Symptom) error {
	//1. first find the symptom group
	filter := bson.D{primitive.E{Key: "symptoms.id", Value: entity.ID}}
	var result []*symptomGroup
	err := sa.db.symptomgroups.Find(filter, &result, nil)
	if err != nil {
		return err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil
	}
	symptomGroup := result[0]

	//2. update the symptom in the symptom group
	symptoms := symptomGroup.Symptoms
	var newSymptoms []symptom
	if symptoms != nil {
		for _, v := range symptoms {
			if v.ID == entity.ID {
				//date updated
				dateUpdated := time.Now()
				v.DateUpdated = &dateUpdated

				//name
				v.Name = entity.Name
			}
			newSymptoms = append(newSymptoms, v)
		}
	}
	symptomGroup.Symptoms = newSymptoms

	//3. save the symptom group
	saveFilter := bson.D{primitive.E{Key: "_id", Value: symptomGroup.ID}}
	err = sa.db.symptomgroups.ReplaceOne(saveFilter, symptomGroup, nil)
	if err != nil {
		return err
	}

	return nil
}

//ReadAllSymptomGroups reads all symptom groups
func (sa *Adapter) ReadAllSymptomGroups() ([]*model.SymptomGroup, error) {
	options := options.Find()
	options.SetSort(bson.D{primitive.E{Key: "name", Value: 1}}) //sort by "name" //gr1 and gr2

	var result []*symptomGroup
	err := sa.db.symptomgroups.Find(nil, &result, options)
	if err != nil {
		return nil, err
	}
	var resultList []*model.SymptomGroup
	if result != nil {
		for _, sg := range result {
			//symptoms
			var symptoms []model.Symptom
			if sg.Symptoms != nil {
				for _, s := range sg.Symptoms {
					item := model.Symptom{ID: s.ID, Name: s.Name}
					symptoms = append(symptoms, item)
				}
			}

			entity := &model.SymptomGroup{ID: sg.ID, Name: sg.Name, Symptoms: symptoms}
			resultList = append(resultList, entity)
		}
	}
	return resultList, nil
}

//ReadSymptoms reads all the symptoms
func (sa *Adapter) ReadSymptoms(appVersion string) (*model.Symptoms, error) {
	filter := bson.D{primitive.E{Key: "app_version", Value: appVersion}}
	var symptoms *model.Symptoms
	err := sa.db.symptoms.FindOne(filter, &symptoms, nil)
	if err != nil {
		return nil, err
	}
	return symptoms, nil
}

//CreateOrUpdateSymptoms creates symptoms for the provided version or update them if already created
func (sa *Adapter) CreateOrUpdateSymptoms(appVersion string, items string) (*bool, error) {
	filter := bson.D{primitive.E{Key: "app_version", Value: appVersion}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "items", Value: items},
		}},
	}

	//insert if not exists
	opt := options.Update()
	upsert := true
	opt.Upsert = &upsert

	updateResult, err := sa.db.symptoms.UpdateOne(filter, update, opt)
	if err != nil {
		return nil, err
	}

	create := true
	if updateResult.ModifiedCount == 1 {
		create = false //modified
	}
	return &create, nil
}

//ReadAllSymptomRules reads all the symptom rules
func (sa *Adapter) ReadAllSymptomRules() ([]*model.SymptomRule, error) {
	filter := bson.D{}
	var result []*symptomRule
	err := sa.db.symptomrules.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	var resultList []*model.SymptomRule
	if result != nil {
		for _, symptomRule := range result {
			county := model.County{ID: symptomRule.CountyID}
			var items []model.SymptomRuleItem
			if symptomRule.Items != nil {
				for _, c := range symptomRule.Items {
					countyStatus := model.CountyStatus{ID: c.CountyStatusID}
					item := model.SymptomRuleItem{Gr1: c.Gr1, Gr2: c.Gr2, CountyStatus: countyStatus, NextStep: c.NextStep}
					items = append(items, item)
				}
			}
			resultItem := &model.SymptomRule{ID: symptomRule.ID, County: county, Gr1Count: symptomRule.Gr1Count,
				Gr2Count: symptomRule.Gr2Count, Items: items}
			resultList = append(resultList, resultItem)
		}
	}
	return resultList, nil
}

//CreateSymptomRule creates a symptom rule
func (sa *Adapter) CreateSymptomRule(countyID string, gr1Count int, gr2Count int, items []model.SymptomRuleItem) (*model.SymptomRule, error) {
	var resultItem model.SymptomRule

	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. create symptom rule
		id, err := uuid.NewUUID()
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		dateCreated := time.Now()
		var dItems []symptomRuleItem
		if items != nil {
			for _, current := range items {
				item := symptomRuleItem{Gr1: current.Gr1, Gr2: current.Gr2,
					CountyStatusID: current.CountyStatus.ID, NextStep: current.NextStep}
				dItems = append(dItems, item)
			}
		}
		symptomRule := symptomRule{ID: id.String(), CountyID: countyID,
			Gr1Count: gr1Count, Gr2Count: gr2Count, Items: dItems, DateCreated: dateCreated}
		_, err = sa.db.symptomrules.InsertOneWithContext(sessionContext, &symptomRule)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		county := model.County{ID: symptomRule.CountyID}
		resultItem = model.SymptomRule{ID: symptomRule.ID, County: county,
			Gr1Count: symptomRule.Gr1Count, Gr2Count: symptomRule.Gr2Count, Items: items}

		//2. delete all statuses of the users
		err = sa.deleteAllStatuses(sessionContext)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &resultItem, nil
}

//FindSymptomRule finds a symptom rule
func (sa *Adapter) FindSymptomRule(ID string) (*model.SymptomRule, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*symptomRule
	err := sa.db.symptomrules.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	symptomRule := result[0]

	var items []model.SymptomRuleItem
	if symptomRule.Items != nil {
		for _, c := range symptomRule.Items {
			countyStatus := model.CountyStatus{ID: c.CountyStatusID}
			item := model.SymptomRuleItem{Gr1: c.Gr1, Gr2: c.Gr2,
				CountyStatus: countyStatus, NextStep: c.NextStep}
			items = append(items, item)
		}
	}
	county := model.County{ID: symptomRule.CountyID}
	resultItem := model.SymptomRule{ID: symptomRule.ID, County: county,
		Gr1Count: symptomRule.Gr1Count, Gr2Count: symptomRule.Gr2Count, Items: items}
	return &resultItem, nil
}

//FindSymptomRuleByCountyID finds a symptom rule by county id
func (sa *Adapter) FindSymptomRuleByCountyID(countyID string) (*model.SymptomRule, error) {
	filter := bson.D{primitive.E{Key: "county_id", Value: countyID}}
	var result []*symptomRule
	err := sa.db.symptomrules.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	symptomRule := result[0]

	var items []model.SymptomRuleItem
	if symptomRule.Items != nil {
		for _, c := range symptomRule.Items {
			countyStatus := model.CountyStatus{ID: c.CountyStatusID}
			item := model.SymptomRuleItem{Gr1: c.Gr1, Gr2: c.Gr2,
				CountyStatus: countyStatus, NextStep: c.NextStep}
			items = append(items, item)
		}
	}
	county := model.County{ID: symptomRule.CountyID}
	resultItem := model.SymptomRule{ID: symptomRule.ID, County: county,
		Gr1Count: symptomRule.Gr1Count, Gr2Count: symptomRule.Gr2Count, Items: items}
	return &resultItem, nil
}

//SaveSymptomRule saves a symptom rule
func (sa *Adapter) SaveSymptomRule(entity *model.SymptomRule) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. save
		findFilter := bson.D{primitive.E{Key: "_id", Value: entity.ID}}
		var result []*symptomRule
		err = sa.db.symptomrules.FindWithContext(sessionContext, findFilter, &result, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if result == nil || len(result) == 0 {
			//not found
			abortTransaction(sessionContext)
			return errors.New("there is no a symptom rule for the provided id")
		}
		symptomRule := result[0]

		//update the values
		symptomRule.CountyID = entity.County.ID
		symptomRule.Gr1Count = entity.Gr1Count
		symptomRule.Gr2Count = entity.Gr2Count

		var items []symptomRuleItem
		if entity.Items != nil {
			for _, c := range entity.Items {
				item := symptomRuleItem{Gr1: c.Gr1, Gr2: c.Gr2, CountyStatusID: c.CountyStatus.ID, NextStep: c.NextStep}
				items = append(items, item)
			}
		}
		symptomRule.Items = items

		dateUpdated := time.Now()
		symptomRule.DateUpdated = &dateUpdated

		filter := bson.D{primitive.E{Key: "_id", Value: symptomRule.ID}}
		err = sa.db.symptomrules.ReplaceOneWithContext(sessionContext, filter, symptomRule, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//2. delete all statuses of the users
		err = sa.deleteAllStatuses(sessionContext)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//DeleteSymptomRule deletes a symptom rule
func (sa *Adapter) DeleteSymptomRule(ID string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. delete
		filter := bson.D{primitive.E{Key: "_id", Value: ID}}
		result, err := sa.db.symptomrules.DeleteOneWithContext(sessionContext, filter, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if result == nil {
			abortTransaction(sessionContext)
			return errors.New("result is nil for symptom rule item with id " + ID)
		}
		deletedCount := result.DeletedCount
		if deletedCount == 0 {
			abortTransaction(sessionContext)
			return errors.New("there is no a symptom rule for id " + ID)
		}
		if deletedCount > 1 {
			abortTransaction(sessionContext)
			return errors.New("deleted more than one records for id " + ID)
		}

		//2. delete all statuses of the users
		err = sa.deleteAllStatuses(sessionContext)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//FindCRulesByCountyID finds the rules for county
func (sa *Adapter) FindCRulesByCountyID(appVersion string, countyID string) (*model.CRules, error) {
	filter := bson.D{primitive.E{Key: "app_version", Value: appVersion},
		primitive.E{Key: "county_id", Value: countyID}}
	var symptomsRules *model.CRules
	err := sa.db.crules.FindOne(filter, &symptomsRules, nil)
	if err != nil {
		return nil, err
	}
	return symptomsRules, nil
}

//CreateOrUpdateCRules creates crule or update it if already created
func (sa *Adapter) CreateOrUpdateCRules(appVersion string, countyID string, data string) (*bool, error) {
	filter := bson.D{primitive.E{Key: "app_version", Value: appVersion}, primitive.E{Key: "county_id", Value: countyID}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "data", Value: data},
		}},
	}

	//insert if not exists
	opt := options.Update()
	upsert := true
	opt.Upsert = &upsert

	updateResult, err := sa.db.crules.UpdateOne(filter, update, opt)
	if err != nil {
		return nil, err
	}

	create := true
	if updateResult.ModifiedCount == 1 {
		create = false //modified
	}
	return &create, nil
}

//CreateTraceReports creates trace reports items
func (sa *Adapter) CreateTraceReports(items []model.TraceExposure) (int, error) {

	//we need create []Interface{}!
	data := make([]interface{}, len(items))
	for i, v := range items {
		data[i] = v
	}

	//insert the items
	result, err := sa.db.traceexposures.InsertMany(data, nil)
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, errors.New("for some reasons the result is nil when create many trace items")
	}

	//return the result
	insertedCount := len(result.InsertedIDs)
	return insertedCount, nil
}

//ReadTraceExposures reads the exposures
func (sa *Adapter) ReadTraceExposures(timestamp *int64, dateAdded *int64) ([]model.TraceExposure, error) {
	filter := bson.M{}

	if timestamp != nil {
		filter["timestamp"] = bson.M{"$gte": timestamp}
	}
	if dateAdded != nil {
		filter["date_added"] = bson.M{"$gte": dateAdded}
	}

	options := options.Find()
	options.SetSort(bson.D{primitive.E{Key: "timestamp", Value: 1}}) //sort by "timestamp"

	var result []model.TraceExposure
	err := sa.db.traceexposures.Find(filter, &result, options)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type manualTestUserJoin struct {
	ID            string    `bson:"_id"`
	HistoryID     string    `bson:"ehistory_id"`
	LocationID    *string   `bson:"location_id"`
	CountyID      *string   `bson:"county_id"`
	EncryptedKey  string    `bson:"encrypted_key"`
	EncryptedBlob string    `bson:"encrypted_blob"`
	Status        string    `bson:"status"`
	DateCreated   time.Time `bson:"date_created"`

	UserID                   string  `bson:"user_id"`
	UserExternalID           string  `bson:"user_external_id"`
	UserUUID                 string  `bson:"user_uuid"`
	UserPublicKey            string  `bson:"user_public_key"`
	UserConsent              bool    `bson:"user_consent"`
	UserExposureNotification bool    `bson:"user_exposure_notification"`
	UserEncryptedKey         *string `bson:"user_encrypted_key"`
	UserEncryptedBlob        *string `bson:"user_encrypted_blob"`
	UserAccounts             []struct {
		ID         string `bson:"id"`
		ExternalID string `bson:"external_id"`
		Default    bool   `bson:"default"`
		Active     bool   `bson:"active"`

		FirstName  string `bson:"first_name"`
		MiddleName string `bson:"middle_name"`
		LastName   string `bson:"last_name"`
		BirthDate  string `bson:"birth_date"`
		Gender     string `bson:"gender"`
		Address1   string `bson:"address1"`
		Address2   string `bson:"address2"`
		Address3   string `bson:"address3"`
		City       string `bson:"city"`
		State      string `bson:"state"`
		ZipCode    string `bson:"zip_code"`
		Phone      string `bson:"phone"`
		Email      string `bson:"email"`
	} `bson:"user_accounts"`
}

//FindManualTestsByCountyIDDeep find the manual test for a county
func (sa *Adapter) FindManualTestsByCountyIDDeep(countyID string, status *string) ([]*model.EManualTest, error) {
	// construct the query
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "users",
		"localField":   "user_id",
		"foreignField": "_id",
		"as":           "user",
	}})
	if countyID != "all" {
		//we need to filter by county

		locsFilter := bson.D{primitive.E{Key: "county_id", Value: countyID}}
		var locsResult []*location
		err := sa.db.locations.Find(locsFilter, &locsResult, nil)
		if err != nil {
			return nil, err
		}
		if len(locsResult) == 0 {
			return nil, errors.New("there is no any location for this county")
		}
		var locationIDs []string
		for _, item := range locsResult {
			locationIDs = append(locationIDs, item.ID)
		}
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$or": []interface{}{bson.M{"county_id": countyID}, bson.M{"location_id": bson.M{"$in": locationIDs}}}}})
	}
	if status != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"status": status}})
	}

	pipeline = append(pipeline, bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{
			"_id": 1, "ehistory_id": 1, "location_id": 1, "county_id": 1, "encrypted_key": 1, "encrypted_blob": 1, "status": 1, "date_created": 1,
			"user_id": "$user._id", "user_external_id": "$user.external_id", "user_uuid": "$user.uuid", "user_public_key": "$user.public_key",
			"user_consent": "$user.consent", "user_exposure_notification": "$user._exposure_notification",
			"user_info": "$user.info", "user_encrypted_key": "$user.encrypted_key", "user_encrypted_blob": "$user.encrypted_blob",
			"user_accounts": "$user.accounts",
		}},
		bson.M{"$sort": bson.D{primitive.E{Key: "date_created", Value: -1}}})

	var result []*manualTestUserJoin
	err := sa.db.emanualtests.Aggregate(pipeline, &result, nil)
	if err != nil {
		return nil, err

	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}

	var resultList []*model.EManualTest
	for _, item := range result {

		accounts := make([]model.Account, len(item.UserAccounts))
		if item.UserAccounts != nil {
			for i, acc := range item.UserAccounts {
				accounts[i] = model.Account{ID: acc.ID, ExternalID: acc.ExternalID, Default: acc.Default, Active: acc.Active, FirstName: acc.FirstName,
					MiddleName: acc.MiddleName, LastName: acc.LastName, BirthDate: acc.BirthDate, Gender: acc.Gender, Address1: acc.Address1,
					Address2: acc.Address2, Address3: acc.Address3, City: acc.City, State: acc.State, ZipCode: acc.ZipCode, Phone: acc.Phone, Email: acc.Email}
			}
		}

		user := model.User{ID: item.UserID, ExternalID: item.UserExternalID, UUID: item.UserUUID, PublicKey: item.UserPublicKey,
			Consent: item.UserConsent, ExposureNotification: item.UserExposureNotification,
			EncryptedKey: item.UserEncryptedKey, EncryptedBlob: item.UserEncryptedBlob, Accounts: accounts}

		mt := model.EManualTest{ID: item.ID, HistoryID: item.HistoryID, LocationID: item.LocationID, CountyID: item.CountyID,
			EncryptedKey: item.EncryptedKey, EncryptedBlob: item.EncryptedBlob,
			Status: item.Status, Date: item.DateCreated, User: user, AccountID: item.UserID}

		resultList = append(resultList, &mt)
	}
	return resultList, nil
}

//FindManualTestImage finds the manual test image
func (sa *Adapter) FindManualTestImage(ID string) (*string, *string, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	var result []*eManualTest
	err := sa.db.emanualtests.Find(filter, &result, nil)
	if err != nil {
		return nil, nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil, nil
	}
	manualTest := result[0]
	return &manualTest.EncryptedImageKey, &manualTest.EncryptedImageBlob, nil
}

//ProcessManualTest processes manual test
func (sa *Adapter) ProcessManualTest(ID string, status string, encryptedKey *string, encryptedBlob *string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			fmt.Println(err)
			return err
		}

		//1. find the manual test
		mtFilter := bson.D{primitive.E{Key: "_id", Value: ID}}
		var mt *eManualTest
		err = sa.db.emanualtests.FindOneWithContext(sessionContext, mtFilter, &mt, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if mt == nil {
			//not found
			abortTransaction(sessionContext)
			return errors.New("there is no a manual test for the provided id")
		}
		historyID := mt.EHistoryID

		//2. update the manual test status or delete the manual test if the status is "verified"
		if status == "verified" {
			//delete the manual test
			_, err = sa.db.emanualtests.DeleteOneWithContext(sessionContext, mtFilter, nil)
			if err != nil {
				abortTransaction(sessionContext)
				return err
			}
		} else {
			//update the manual tests
			mt.Status = status

			//save the manual test
			err = sa.db.emanualtests.ReplaceOneWithContext(sessionContext, mtFilter, mt, nil)
			if err != nil {
				abortTransaction(sessionContext)
				return err
			}
		}

		//3. update the history item if the status is "verified"
		if status == "verified" {
			//3.1 find the history
			historyFilter := bson.D{primitive.E{Key: "_id", Value: historyID}}
			var history *model.EHistory
			err = sa.db.ehistory.FindOneWithContext(sessionContext, historyFilter, &history, nil)
			if err != nil {
				abortTransaction(sessionContext)
				return err
			}
			if history == nil {
				//not found
				abortTransaction(sessionContext)
				return errors.New("there is no a history for the provided manual test")
			}

			//3.2 update the history
			history.Type = "verified_manual_test"
			history.EncryptedKey = *encryptedKey
			history.EncryptedBlob = *encryptedBlob

			//3.3 save the history item
			err = sa.db.ehistory.ReplaceOneWithContext(sessionContext, historyFilter, history, nil)
			if err != nil {
				abortTransaction(sessionContext)
				return err
			}

			//3.4 remove the status of the user
			statusFilter := bson.D{primitive.E{Key: "user_id", Value: history.UserID}}
			//from estatus
			_, err = sa.db.estatus.DeleteManyWithContext(sessionContext, statusFilter, nil)
			if err != nil {
				abortTransaction(sessionContext)
				return err
			}
		}

		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.Printf("error on commiting a transaction - %s", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//ReadAllAccessRules reads all access rules
func (sa *Adapter) ReadAllAccessRules() ([]*model.AccessRule, error) {
	filter := bson.D{}
	var result []*accessRule
	err := sa.db.accessrules.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	var resultList []*model.AccessRule
	if result != nil {
		for _, accessRule := range result {
			county := model.County{ID: accessRule.CountyID}
			var rules []model.AccessRuleCountyStatus
			if accessRule.Rules != nil {
				for _, c := range accessRule.Rules {
					rule := model.AccessRuleCountyStatus{CountyStatusID: c.CountyStatusID, Value: c.Value}
					rules = append(rules, rule)
				}
			}
			resultItem := &model.AccessRule{ID: accessRule.ID, County: county, Rules: rules}
			resultList = append(resultList, resultItem)
		}
	}
	return resultList, nil
}

//CreateAccessRule creates an access rule
func (sa *Adapter) CreateAccessRule(countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error) {
	var resultItem model.AccessRule

	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. validate the input data
		err = sa.validateAccessRuleData(sessionContext, countyID, rules)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//2. create access rule
		id, err := uuid.NewUUID()
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		dateCreated := time.Now()
		var rItems []accessRuleCountyStatus
		if rules != nil {
			for _, current := range rules {
				item := accessRuleCountyStatus{CountyStatusID: current.CountyStatusID, Value: current.Value}
				rItems = append(rItems, item)
			}
		}
		accessRule := accessRule{ID: id.String(), CountyID: countyID, Rules: rItems, DateCreated: dateCreated}
		_, err = sa.db.accessrules.InsertOneWithContext(sessionContext, &accessRule)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		c := model.County{ID: accessRule.CountyID}
		resultItem = model.AccessRule{ID: accessRule.ID, County: c, Rules: rules}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &resultItem, nil
}

//UpdateAccessRule updates an access rule
func (sa *Adapter) UpdateAccessRule(ID string, countyID string, rules []model.AccessRuleCountyStatus) (*model.AccessRule, error) {
	var resultItem model.AccessRule

	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//1. find the access rule item
		arFilter := bson.D{primitive.E{Key: "_id", Value: ID}}
		var arResult []*accessRule
		err = sa.db.accessrules.FindWithContext(sessionContext, arFilter, &arResult, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if arResult == nil || len(arResult) == 0 {
			abortTransaction(sessionContext)
			return errors.New("there is no an access rule for the provided id")
		}
		accessRule := arResult[0]

		//2. validate the input data
		err = sa.validateAccessRuleData(sessionContext, countyID, rules)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		//3. update the access rule
		accessRule.CountyID = countyID
		dateUpdated := time.Now()
		accessRule.DateUpdated = &dateUpdated
		var rItems []accessRuleCountyStatus
		if rules != nil {
			for _, current := range rules {
				item := accessRuleCountyStatus{CountyStatusID: current.CountyStatusID, Value: current.Value}
				rItems = append(rItems, item)
			}
		}
		accessRule.Rules = rItems

		//4. save the access rule
		saveFilter := bson.D{primitive.E{Key: "_id", Value: accessRule.ID}}
		err = sa.db.accessrules.ReplaceOneWithContext(sessionContext, saveFilter, &accessRule, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		c := model.County{ID: accessRule.CountyID}
		resultItem = model.AccessRule{ID: accessRule.ID, County: c, Rules: rules}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &resultItem, nil
}

func (sa *Adapter) validateAccessRuleData(sessionContext mongo.SessionContext, countyID string, rules []model.AccessRuleCountyStatus) error {
	//1. validate the county id
	countyFilter := bson.D{primitive.E{Key: "_id", Value: countyID}}
	var countyResult []*county
	err := sa.db.counties.FindWithContext(sessionContext, countyFilter, &countyResult, nil)
	if err != nil {
		return err
	}
	if countyResult == nil || len(countyResult) == 0 {
		return errors.New("there is no a county for the provided id")
	}
	county := countyResult[0]

	//2. validate the county statuses ids
	countyStatuses := county.CountyStatuses
	if countyStatuses == nil || len(countyStatuses) == 0 {
		return errors.New("there is no county statuses for this county")
	}
	if rules != nil && len(rules) > 0 {
		for _, r := range rules {
			contains := sa.containsCountyStatus(r.CountyStatusID, countyStatuses)
			if !contains {
				return errors.New("there is invalid county status id")
			}
		}
	}
	return nil
}

//FindAccessRuleByCountyID finds the access rule for a specific county
func (sa *Adapter) FindAccessRuleByCountyID(countyID string) (*model.AccessRule, error) {
	filter := bson.D{primitive.E{Key: "county_id", Value: countyID}}
	var result []*accessRule
	err := sa.db.accessrules.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	accessRule := result[0]

	var items []model.AccessRuleCountyStatus
	if accessRule.Rules != nil {
		for _, c := range accessRule.Rules {
			item := model.AccessRuleCountyStatus{CountyStatusID: c.CountyStatusID, Value: c.Value}
			items = append(items, item)
		}
	}
	county := model.County{ID: accessRule.CountyID}
	resultItem := model.AccessRule{ID: accessRule.ID, County: county, Rules: items}
	return &resultItem, nil
}

//DeleteAccessRule deletes an access rule
func (sa *Adapter) DeleteAccessRule(ID string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	result, err := sa.db.accessrules.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("error deleting an access rule - %s", err)
		return err
	}
	if result == nil {
		return errors.New("result is nil for access rule with id " + ID)
	}
	deletedCount := result.DeletedCount
	if deletedCount == 0 {
		return errors.New("there is no an access rule for id " + ID)
	}
	if deletedCount > 1 {
		return errors.New("deleted more than one records for id " + ID)
	}
	return nil
}

//FindExternalUserIDsByTestsOrderNumbers finds the external users ids for the tests orders numbers
func (sa *Adapter) FindExternalUserIDsByTestsOrderNumbers(orderNumbers []string) (map[string]*string, error) {

	//1. get based on the user
	ubList, err := sa.findEUIDsBasedOnUser(orderNumbers)
	if err != nil {
		return nil, err
	}

	//2. get based on the user accounts
	uabList, err := sa.findEUIDsBasedOnUserAccounts(orderNumbers)
	if err != nil {
		return nil, err
	}

	//3. merge the maps
	for key, value := range uabList {
		ubList[key] = value
	}

	return ubList, nil
}

func (sa *Adapter) findEUIDsBasedOnUser(orderNumbers []string) (map[string]*string, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}},
		{"$match": bson.M{"order_number": bson.M{"$in": orderNumbers}}},
		{"$unwind": "$user"},
		{"$project": bson.M{
			"_id": 1, "order_number": 1,
			"user_id": "$user._id", "user_external_id": "$user.external_id",
		}}}

	type ctuJoin struct {
		ID          string `bson:"_id"`
		OrderNumber string `bson:"order_number"`

		UserID         string `bson:"user_id"`
		UserExternalID string `bson:"user_external_id"`
	}
	var result []*ctuJoin
	err := sa.db.ctests.Aggregate(pipeline, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return make(map[string]*string, 0), nil
	}
	mapData := make(map[string]*string, len(result))
	for _, v := range result {
		mapData[v.OrderNumber] = &v.UserExternalID
	}
	return mapData, nil
}

func (sa *Adapter) findEUIDsBasedOnUserAccounts(orderNumbers []string) (map[string]*string, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "accounts.id",
			"as":           "user",
		}},
		{"$match": bson.M{"order_number": bson.M{"$in": orderNumbers}}},
		{"$unwind": "$user"},
		{"$project": bson.M{
			"_id": 1, "order_number": 1, "user_id": 1,
			"user_account": bson.M{
				"$filter": bson.M{
					"input": "$user.accounts",
					"as":    "ac",
					"cond": bson.M{
						"$and": []bson.M{
							{"$eq": []interface{}{"$$ac.id", "$user_id"}},
						},
					},
				},
			},
		}}}

	type ctuJoin struct {
		ID          string `bson:"_id"`
		OrderNumber string `bson:"order_number"`
		UserID      string `bson:"user_id"`

		UserAccount []struct {
			ExternalID string `bson:"external_id"`
		} `bson:"user_account"`
	}

	var result []ctuJoin
	err := sa.db.ctests.Aggregate(pipeline, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return make(map[string]*string, 0), nil
	}
	mapData := make(map[string]*string, len(result))
	for _, v := range result {
		if len(v.UserAccount) > 0 {
			mapData[v.OrderNumber] = &v.UserAccount[0].ExternalID
		}
	}
	return mapData, nil
}

//CreateOrUpdateUINOverride creates a new uin override entity or updates it if already created
func (sa *Adapter) CreateOrUpdateUINOverride(uin string, interval int, category *string, expiration *time.Time) error {
	filter := bson.D{primitive.E{Key: "uin", Value: uin}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "interval", Value: interval},
			primitive.E{Key: "category", Value: category},
			primitive.E{Key: "expiration", Value: expiration},
		}},
	}

	//insert if not exists
	opt := options.Update()
	upsert := true
	opt.Upsert = &upsert

	_, err := sa.db.uinoverrides.UpdateOne(filter, update, opt)
	if err != nil {
		return err
	}

	return nil
}

//FindUINOverride finds the uin override for the provided uin. It makes additional check for the expiration because of the mongoDB TTL delay
func (sa *Adapter) FindUINOverride(uin string) (*model.UINOverride, error) {
	now := time.Now()
	filter := bson.D{
		primitive.E{Key: "uin", Value: uin},
		primitive.E{Key: "$or", Value: []interface{}{
			bson.D{primitive.E{Key: "expiration", Value: bson.M{"$gte": now}}},
			bson.D{primitive.E{Key: "expiration", Value: nil}},
		}},
	}

	var uinOverrides []*model.UINOverride
	err := sa.db.uinoverrides.Find(filter, &uinOverrides, nil)
	if err != nil {
		return nil, err
	}
	if len(uinOverrides) == 0 {
		//not found
		return nil, nil
	}

	return uinOverrides[0], nil
}

//FindUINOverrides finds the uin override for the provided uin. If uin is nil then it gives all
func (sa *Adapter) FindUINOverrides(uin *string, sort *string) ([]*model.UINOverride, error) {
	//filter by uin if provided
	filter := bson.D{}
	if uin != nil {
		filter = bson.D{primitive.E{Key: "uin", Value: *uin}}
	}

	// sort by if provided
	var opt *options.FindOptions
	if sort != nil {
		opt = options.Find()
		opt.SetSort(bson.D{primitive.E{Key: *sort, Value: 1}})
	}

	var result []*model.UINOverride
	err := sa.db.uinoverrides.Find(filter, &result, opt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//CreateUINOverride creates a new uin override entity
func (sa *Adapter) CreateUINOverride(uin string, interval int, category *string, expiration *time.Time) (*model.UINOverride, error) {
	uinOverride := model.UINOverride{UIN: uin, Interval: interval, Category: category, Expiration: expiration}
	_, err := sa.db.uinoverrides.InsertOne(&uinOverride)
	if err != nil {
		return nil, err
	}

	return &uinOverride, nil
}

//UpdateUINOverride updates uin override entity
func (sa *Adapter) UpdateUINOverride(uin string, interval int, category *string, expiration *time.Time) (*string, error) {
	filter := bson.D{primitive.E{Key: "uin", Value: uin}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "interval", Value: interval},
			primitive.E{Key: "category", Value: category},
			primitive.E{Key: "expiration", Value: expiration},
		}},
	}

	result, err := sa.db.uinoverrides.UpdateOne(filter, update, nil)
	if err != nil {
		return nil, err
	}

	res := fmt.Sprintf("%d matched, %d modified", result.MatchedCount, result.ModifiedCount)
	return &res, nil
}

//DeleteUINOverride deletes uin override entity
func (sa *Adapter) DeleteUINOverride(uin string) error {
	filter := bson.D{primitive.E{Key: "uin", Value: uin}}
	result, err := sa.db.uinoverrides.DeleteOne(filter, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return errors.New("result is nil for uin override item with uin " + uin)
	}
	deletedCount := result.DeletedCount
	if deletedCount == 0 {
		return errors.New("there is no a uin override for uin " + uin)
	}
	if deletedCount > 1 {
		return errors.New("deleted more than one records for uin " + uin)
	}

	//success - count = 1
	return nil
}

//FindUINBuildingAccess finds the UIN buinding access for the provided uin
func (sa *Adapter) FindUINBuildingAccess(uin string) (*model.UINBuildingAccess, error) {
	filter := bson.D{primitive.E{Key: "uin", Value: uin}}
	var uinBuildingAccess *model.UINBuildingAccess
	err := sa.db.uinbuildingaccess.FindOne(filter, &uinBuildingAccess, nil)
	if err != nil {
		return nil, err
	}
	return uinBuildingAccess, nil
}

//CreateOrUpdateUINBuildingAccess creates UIN building access or update it if already created
func (sa *Adapter) CreateOrUpdateUINBuildingAccess(uin string, date time.Time, access string) error {
	filter := bson.D{primitive.E{Key: "uin", Value: uin}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "date", Value: date},
			primitive.E{Key: "access", Value: access},
		}},
	}

	//insert if not exists
	opt := options.Update()
	upsert := true
	opt.Upsert = &upsert

	_, err := sa.db.uinbuildingaccess.UpdateOne(filter, update, opt)
	if err != nil {
		return err
	}

	return nil
}

//ReadAllRosters reads all rosters
func (sa *Adapter) ReadAllRosters() ([]map[string]string, error) {
	filter := bson.D{}
	var result []map[string]string
	err := sa.db.rosters.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

//FindRosterByPhone finds the roster for the user with the given phone number
func (sa *Adapter) FindRosterByPhone(phone string) (map[string]string, error) {
	filter := bson.D{primitive.E{Key: "phone", Value: phone}}

	var result []map[string]string
	err := sa.db.rosters.Find(filter, &result, nil)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		//not found
		return nil, nil
	}
	item := result[0]
	return item, nil
}

//FindRosters returns the roster members matching filters, sorted, and paginated
func (sa *Adapter) FindRosters(f *utils.Filter, sortBy string, sortOrder int, limit int, offset int) ([]map[string]interface{}, error) {
	var filter bson.D
	if f != nil {
		filter = constructDataFilter(f).(bson.D)
	}

	options := options.Find()
	options.SetSort(bson.D{primitive.E{Key: sortBy, Value: sortOrder}})
	if limit > 0 {
		options.SetLimit(int64(limit))
	}
	if offset > 0 {
		options.SetSkip(int64(offset))
	}

	projection := bson.D{
		bson.E{Key: "_id", Value: 0},
	}
	options.SetProjection(projection)

	var result []map[string]interface{}
	err := sa.db.rosters.Find(filter, &result, options)
	if err != nil {
		log.Println("GetRoster:", err.Error())
		return []map[string]interface{}{}, err
	} else if len(result) < 1 {
		log.Println("GetRoster: no roster data found")
		return []map[string]interface{}{}, nil
	}

	return result, nil
}

//CreateRoster creates a roster
func (sa *Adapter) CreateRoster(phone string, uin string, firstName string, middleName string, lastName string, birthDate string, gender string,
	address1 string, address2 string, address3 string, city string, state string, zipCode string, email string, badgeType string) error {

	//insert the roster
	item := map[string]string{"phone": phone, "uin": uin, "first_name": firstName, "middle_name": middleName,
		"last_name": lastName, "birth_date": birthDate, "gender": gender, "address1": address1, "address2": address2,
		"address3": address3, "city": city, "state": state, "zip_code": zipCode, "email": email, "badge_type": badgeType}
	_, err := sa.db.rosters.InsertOne(&item)
	if err != nil {
		return err
	}

	return nil
}

//UpdateRoster updates a roster
func (sa *Adapter) UpdateRoster(uin string, firstName string, middleName string, lastName string, birthDate string, gender string,
	address1 string, address2 string, address3 string, city string, state string, zipCode string, email string, badgeType string) error {

	filter := bson.D{primitive.E{Key: "uin", Value: uin}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "first_name", Value: firstName},
			primitive.E{Key: "middle_name", Value: middleName},
			primitive.E{Key: "last_name", Value: lastName},
			primitive.E{Key: "birth_date", Value: birthDate},
			primitive.E{Key: "gender", Value: gender},
			primitive.E{Key: "address1", Value: address1},
			primitive.E{Key: "address2", Value: address2},
			primitive.E{Key: "address3", Value: address3},
			primitive.E{Key: "city", Value: city},
			primitive.E{Key: "state", Value: state},
			primitive.E{Key: "zip_code", Value: zipCode},
			primitive.E{Key: "email", Value: email},
			primitive.E{Key: "badge_type", Value: badgeType},
		}},
	}

	result, err := sa.db.rosters.UpdateOne(filter, update, nil)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("there is no an item for the provided uin")
	}

	return nil
}

//CreateRosterItems creates roster items
func (sa *Adapter) CreateRosterItems(items []map[string]string) error {
	//insert the items
	//need to prepare the input data
	data := make([]interface{}, len(items))
	for i, c := range items {
		data[i] = c
	}
	result, err := sa.db.rosters.InsertMany(data, nil)
	if err != nil {
		log.Printf("error inserting many roster items - %s", err)
		return err
	}
	if result == nil {
		log.Println("for some reasons the result is nil when create many roster items")
		return errors.New("for some reasons the result is nil when create many roster items")
	}

	return nil
}

//DeleteRosterByPhone deletes the roster for the provided phone
func (sa *Adapter) DeleteRosterByPhone(phone string) error {
	deleteFilter := bson.D{primitive.E{Key: "phone", Value: phone}}
	result, err := sa.db.rosters.DeleteOne(deleteFilter, nil)
	if err != nil {
		log.Printf("error deleting a roster - %s", err)
		return err
	}
	if result == nil {
		return errors.New("result is nil for roster with phone " + phone)
	}
	deletedCount := result.DeletedCount
	if deletedCount == 0 {
		return errors.New("there is no a roster for phone " + phone)
	}
	if deletedCount > 1 {
		return errors.New("deleted more than one records for phone " + phone)
	}

	return nil
}

//DeleteRosterByUIN deletes the roster for the provided uin
func (sa *Adapter) DeleteRosterByUIN(uin string) error {
	deleteFilter := bson.D{primitive.E{Key: "uin", Value: uin}}
	result, err := sa.db.rosters.DeleteOne(deleteFilter, nil)
	if err != nil {
		log.Printf("error deleting a roster - %s", err)
		return err
	}
	if result == nil {
		return errors.New("result is nil for roster with uin " + uin)
	}
	deletedCount := result.DeletedCount
	if deletedCount == 0 {
		return errors.New("there is no a roster for uin " + uin)
	}
	if deletedCount > 1 {
		return errors.New("deleted more than one records for uin " + uin)
	}

	return nil
}

//DeleteAllRosters deletes all rosters
func (sa *Adapter) DeleteAllRosters() error {
	//delete all rosters
	deleteRostersFilter := bson.D{}
	_, err := sa.db.rosters.DeleteMany(deleteRostersFilter, nil)
	if err != nil {
		return err
	}

	return nil
}

//CreateRawSubAccountItems creates raw sub account items
func (sa *Adapter) CreateRawSubAccountItems(items []model.RawSubAccount) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		//apply logic for every raw sub account
		for _, current := range items {
			var accountID *string

			// check if there is a primary account added
			uFilter := bson.D{primitive.E{Key: "external_id", Value: current.PrimaryAccount}}
			var usersResult []*model.User
			err = sa.db.users.FindWithContext(sessionContext, uFilter, &usersResult, nil)
			if err != nil {
				abortTransaction(sessionContext)
				return err
			}
			if len(usersResult) > 0 {
				//there is a primary user for this account

				user := usersResult[0]
				subAccount := user.GetAccountByExternalID(current.UIN)

				if subAccount != nil {
					//there is a sub account, so we need to update it

					//we know the account id
					accountID = &subAccount.ID

					updateSubAccountFilter := bson.D{primitive.E{Key: "accounts.id", Value: subAccount.ID}}
					updateSubAccount := bson.D{
						primitive.E{Key: "$set", Value: bson.D{
							primitive.E{Key: "accounts.$.active", Value: true},
							primitive.E{Key: "accounts.$.first_name", Value: current.FirstName},
							primitive.E{Key: "accounts.$.middle_name", Value: current.MiddleName},
							primitive.E{Key: "accounts.$.last_name", Value: current.LastName},
							primitive.E{Key: "accounts.$.birth_date", Value: current.BirthDate},
							primitive.E{Key: "accounts.$.gender", Value: current.Gender},
							primitive.E{Key: "accounts.$.address1", Value: current.Address1},
							primitive.E{Key: "accounts.$.address2", Value: current.Address2},
							primitive.E{Key: "accounts.$.address3", Value: current.Address3},
							primitive.E{Key: "accounts.$.city", Value: current.City},
							primitive.E{Key: "accounts.$.state", Value: current.State},
							primitive.E{Key: "accounts.$.zip_code", Value: current.ZipCode},
							primitive.E{Key: "accounts.$.phone", Value: current.Phone},
							primitive.E{Key: "accounts.$.email", Value: current.Email},
						}},
					}

					_, err := sa.db.users.UpdateOneWithContext(sessionContext, updateSubAccountFilter, updateSubAccount, nil)
					if err != nil {
						abortTransaction(sessionContext)
						return err
					}
				} else {
					//there is no a sub account, so we need to create it

					//generate sub account id
					genSubAccountID, err := uuid.NewUUID()
					if err != nil {
						abortTransaction(sessionContext)
						return err
					}
					subAccID := genSubAccountID.String()
					accountID = &subAccID

					newSubAccount := model.Account{ID: subAccID, ExternalID: current.UIN, Default: false, Active: true, FirstName: current.FirstName,
						MiddleName: current.MiddleName, LastName: current.LastName, BirthDate: current.BirthDate, Gender: current.Gender, Address1: current.Address1,
						Address2: current.Address2, Address3: current.Address3, City: current.City, State: current.State, ZipCode: current.ZipCode, Phone: current.Phone,
						Email: current.Email}

					createSubAccountFilter := bson.D{primitive.E{Key: "_id", Value: user.ID}}
					createSubAccount := bson.D{
						primitive.E{Key: "$push", Value: bson.D{
							primitive.E{Key: "accounts", Value: newSubAccount},
						}},
					}

					_, err = sa.db.users.UpdateOneWithContext(sessionContext, createSubAccountFilter, createSubAccount, nil)
					if err != nil {
						abortTransaction(sessionContext)
						return err
					}
				}
			}

			//save the raw sub account
			current.AccountID = accountID
			_, err = sa.db.rawsubaccounts.InsertOneWithContext(sessionContext, &current)
			if err != nil {
				abortTransaction(sessionContext)
				return err
			}

		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//FindRawSubAccounts finds raw sub accounts
func (sa *Adapter) FindRawSubAccounts(f *utils.Filter, sortBy string, sortOrder int, limit int, offset int) ([]model.RawSubAccount, error) {
	var filter bson.D
	if f != nil {
		filter = constructDataFilter(f).(bson.D)
	}

	options := options.Find()
	options.SetSort(bson.D{primitive.E{Key: sortBy, Value: sortOrder}})
	if limit > 0 {
		options.SetLimit(int64(limit))
	}
	if offset > 0 {
		options.SetSkip(int64(offset))
	}

	var result []model.RawSubAccount
	err := sa.db.rawsubaccounts.Find(filter, &result, options)
	if err != nil {
		return nil, err
	}
	if len(result) < 1 {
		return make([]model.RawSubAccount, 0), nil
	}

	return result, nil
}

//UpdateRawSubAcccount updates raw sub account
func (sa *Adapter) UpdateRawSubAcccount(uin string, firstName string, middleName string, lastName string, birthDate string, gender string,
	address1 string, address2 string, address3 string, city string, state string, zipCode string, netID string, email string) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			fmt.Println(err)
			return err
		}

		//get the raw sub account
		rsaFilter := bson.D{primitive.E{Key: "uin", Value: uin}}
		var rawSubAccounts []*model.RawSubAccount
		err = sa.db.rawsubaccounts.FindWithContext(sessionContext, rsaFilter, &rawSubAccounts, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}
		if len(rawSubAccounts) == 0 {
			abortTransaction(sessionContext)
			log.Printf("there is no a sub account for the provided uin - %s", uin)
			return errors.New("there is no a sub account for the provided uin - " + uin)
		}
		rawSubAccount := rawSubAccounts[0]

		//update the sub account for the user if there is a such
		if rawSubAccount.AccountID != nil {
			updateSubAccountFilter := bson.D{primitive.E{Key: "accounts.id", Value: *rawSubAccount.AccountID}}
			updateSubAccount := bson.D{
				primitive.E{Key: "$set", Value: bson.D{
					primitive.E{Key: "accounts.$.first_name", Value: firstName},
					primitive.E{Key: "accounts.$.middle_name", Value: middleName},
					primitive.E{Key: "accounts.$.last_name", Value: lastName},
					primitive.E{Key: "accounts.$.birth_date", Value: birthDate},
					primitive.E{Key: "accounts.$.gender", Value: gender},
					primitive.E{Key: "accounts.$.address1", Value: address1},
					primitive.E{Key: "accounts.$.address2", Value: address2},
					primitive.E{Key: "accounts.$.address3", Value: address3},
					primitive.E{Key: "accounts.$.city", Value: city},
					primitive.E{Key: "accounts.$.state", Value: state},
					primitive.E{Key: "accounts.$.zip_code", Value: zipCode},
					primitive.E{Key: "accounts.$.email", Value: email},
				}},
			}

			_, err := sa.db.users.UpdateOneWithContext(sessionContext, updateSubAccountFilter, updateSubAccount, nil)
			if err != nil {
				abortTransaction(sessionContext)
				return err
			}
		}

		//update the raw sub account
		rsaSubAccountUpdateFilter := bson.D{primitive.E{Key: "uin", Value: uin}}
		rsaSubAccountUpdate := bson.D{
			primitive.E{Key: "$set", Value: bson.D{
				primitive.E{Key: "first_name", Value: firstName},
				primitive.E{Key: "middle_name", Value: middleName},
				primitive.E{Key: "last_name", Value: lastName},
				primitive.E{Key: "birth_date", Value: birthDate},
				primitive.E{Key: "gender", Value: gender},
				primitive.E{Key: "address1", Value: address1},
				primitive.E{Key: "address2", Value: address2},
				primitive.E{Key: "address3", Value: address3},
				primitive.E{Key: "city", Value: city},
				primitive.E{Key: "state", Value: state},
				primitive.E{Key: "zip_code", Value: zipCode},
				primitive.E{Key: "net_id", Value: netID},
				primitive.E{Key: "email", Value: email},
			}},
		}
		_, err = sa.db.rawsubaccounts.UpdateOneWithContext(sessionContext, rsaSubAccountUpdateFilter, rsaSubAccountUpdate, nil)
		if err != nil {
			abortTransaction(sessionContext)
			return err
		}

		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			log.Printf("error on commiting a transaction - %s", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (sa *Adapter) containsCountyStatus(ID string, list []countyStatus) bool {
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

func (sa *Adapter) deleteAllStatuses(sessionContext mongo.SessionContext) error {
	filter := bson.D{}
	//delete from estatus
	_, err := sa.db.estatus.DeleteManyWithContext(sessionContext, filter, nil)
	if err != nil {
		return err
	}
	return nil
}

func (sa *Adapter) deleteAccountData(sessionContext mongo.SessionContext, account model.Account) error {

	//remove from uinoverrides
	uinOverridesFilter := bson.D{primitive.E{Key: "uin", Value: account.ExternalID}}
	_, err := sa.db.uinoverrides.DeleteOneWithContext(sessionContext, uinOverridesFilter, nil)
	if err != nil {
		log.Printf("error deleting uinoverride record for a user - %s", err)
		return err
	}

	//remove from uinbuildingaccess
	uinBuildingAccessFilter := bson.D{primitive.E{Key: "uin", Value: account.ExternalID}}
	_, err = sa.db.uinbuildingaccess.DeleteOneWithContext(sessionContext, uinBuildingAccessFilter, nil)
	if err != nil {
		log.Printf("error deleting uinbuildingaccess record for a user - %s", err)
		return err
	}

	//remove from ctest
	cTestFilter := bson.D{primitive.E{Key: "user_id", Value: account.ID}}
	_, err = sa.db.ctests.DeleteManyWithContext(sessionContext, cTestFilter, nil)
	if err != nil {
		log.Printf("error deleting ctests for a user - %s", err)
		return err
	}

	//remove from history
	historyFilter := bson.D{primitive.E{Key: "user_id", Value: account.ID}}
	//from ehistory
	_, err = sa.db.ehistory.DeleteManyWithContext(sessionContext, historyFilter, nil)
	if err != nil {
		log.Printf("error deleting ehistories for a user - %s", err)
		return err
	}

	//remove from status
	statusFilter := bson.D{primitive.E{Key: "user_id", Value: account.ID}}
	//from estatus
	_, err = sa.db.estatus.DeleteManyWithContext(sessionContext, statusFilter, nil)
	if err != nil {
		log.Printf("error deleting estatus for a user - %s", err)
		return err
	}

	//remove from manual tests
	mtFilter := bson.D{primitive.E{Key: "user_id", Value: account.ID}}
	_, err = sa.db.emanualtests.DeleteOneWithContext(sessionContext, mtFilter, nil)
	if err != nil {
		log.Printf("error deleting manual tests for a user - %s", err)
		return err
	}

	return nil
}

//NewStorageAdapter creates a new storage adapter instance
func NewStorageAdapter(mongoDBAuth string, mongoDBName string, mongoTimeout string) *Adapter {
	timeout, err := strconv.Atoi(mongoTimeout)
	if err != nil {
		log.Println("Set default timeout - 500")
		timeout = 500
	}
	timeoutMS := time.Millisecond * time.Duration(timeout)

	db := &database{mongoDBAuth: mongoDBAuth, mongoDBName: mongoDBName, mongoTimeout: timeoutMS}
	return &Adapter{db: db}
}

func constructDataFilter(f *utils.Filter) interface{} {
	if f == nil || len(f.Items) == 0 {
		return bson.D{}
	}
	var filter bson.D
	for _, item := range f.Items {
		if len(item.Value) == 1 {
			quoteMeta := regexp.QuoteMeta(item.Value[0])
			filter = append(filter, bson.E{Key: item.Field, Value: bson.D{{"$regex", primitive.Regex{Pattern: quoteMeta, Options: "i"}}}})
		} else if len(item.Value) > 1 {
			var vals []interface{}
			for _, value := range item.Value {
				quoteMeta := regexp.QuoteMeta(value)
				vals = append(vals, bson.D{bson.E{Key: item.Field, Value: bson.D{{"$regex", primitive.Regex{Pattern: quoteMeta, Options: "i"}}}}})
			}
			filter = append(filter, bson.E{Key: "$or", Value: vals})
		}
	}
	return filter
}

func convertToDaysOfOperation(list []operationDay) []model.OperationDay {
	var result []model.OperationDay
	if list != nil {
		for _, d := range list {
			item := model.OperationDay{Name: d.Name, OpenTime: d.OpenTime, CloseTime: d.CloseTime}
			result = append(result, item)
		}
	}
	return result
}

func convertFromDaysOfOperation(list []model.OperationDay) []operationDay {
	var result []operationDay
	if list != nil {
		for _, d := range list {
			item := operationDay{Name: d.Name, OpenTime: d.OpenTime, CloseTime: d.CloseTime}
			result = append(result, item)
		}
	}
	return result
}

func abortTransaction(sessionContext mongo.SessionContext) {
	err := sessionContext.AbortTransaction(sessionContext)
	if err != nil {
		log.Printf("error on aborting a transaction - %s", err)
	}
}
