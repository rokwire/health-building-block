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

package rest

import (
	"health/core/model"
	"time"
)

type guidelinesResponse struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Items       []guidelinesItemsResponse `json:"items"`
} // @name Guideline

type guidelinesItemsResponse struct {
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Type        string `json:"type"`
} // @name GuidelineItem

//AppUserResponse represents user response entity
type AppUserResponse struct {
	ID                   string  `json:"id"`
	UUID                 string  `json:"uuid"`
	PublicKey            string  `json:"public_key"`
	Consent              bool    `json:"consent"`
	ExposureNotification bool    `json:"exposure_notification"`
	RePost               bool    `json:"re_post"`
	EncryptedKey         *string `json:"encrypted_key"`
	EncryptedBlob        *string `json:"encrypted_blob"`

	Accounts []AppUserAccountResponse `json:"accounts"`
} //@name User

//AppUserAccountResponse represents user account response entity
type AppUserAccountResponse struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Default    bool   `json:"default"`
	Active     bool   `json:"active"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	BirthDate  string `json:"birth_date"`
	Gender     string `json:"gender"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	Address3   string `json:"address3"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zip_code"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
} //@name Account

type manualTestImageResponse struct {
	EncryptedImageKey  string `json:"encrypted_image_key"`
	EncryptedImageBlob string `json:"encrypted_image_blob"`
} // @name ManualTestImage

type manualTestResponse struct {
	ID   string          `json:"id"`
	User AppUserResponse `json:"user"`

	//history id
	HistoryID string `json:"history_id"`
	//+ all history fields
	ProviderID *string `json:"provider_id"`
	Provider   *string `json:"provider"`
	LocationID *string `json:"location_id"`
	Location   *string `json:"location"`
	TestType   *string `json:"test_type"`
	Result     *string `json:"result"`
	//+ county id
	CountyID *string `json:"county_id"`
	//deprecated
	Verified bool      `json:"verified"`
	Status   string    `json:"status"`
	Date     time.Time `json:"date"`
}

type eManualTestResponse struct {
	ID            string          `json:"id"`
	User          AppUserResponse `json:"user"`
	HistoryID     string          `json:"history_id"`
	LocationID    *string         `json:"location_id"`
	CountyID      *string         `json:"county_id"`
	EncryptedKey  string          `json:"encrypted_key"`
	EncryptedBlob string          `json:"encrypted_blob"`
	Status        string          `json:"status"`
	Date          time.Time       `json:"date"`
} // @name ManualTest

type locationResponse struct {
	ID              string                         `json:"id"`
	Name            string                         `json:"name"`
	Address1        string                         `json:"address_1"`
	Address2        string                         `json:"address_2"`
	City            string                         `json:"city"`
	State           string                         `json:"state"`
	ZIP             string                         `json:"zip"`
	Country         string                         `json:"country"`
	Latitude        float64                        `json:"latitude"`
	Longitude       float64                        `json:"longitude"`
	Timezone        string                         `json:"timezone"`
	Contact         string                         `json:"contact"`
	DaysOfOperation []locationOperationDayResponse `json:"days_of_operation"`
	URL             string                         `json:"url"`
	Notes           string                         `json:"notes"`
	WaitTimeColor   *string                        `json:"wait_time_color"`

	ProviderID string `json:"provider_id"`
	CountyID   string `json:"county_id"`

	AvailableTests []string `json:"available_tests"`
} // @name Location

type locationOperationDayResponse struct {
	Name      string `json:"name"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
} // @name OperationDay

type providerResponse struct {
	ID                  string   `json:"id"`
	ProviderName        string   `json:"provider_name"`
	ManualTest          bool     `json:"manual_test"`
	AvailableMechanisms []string `json:"available_mechanisms"`
} // @name Provider

type symptomResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
} // @name Symptom

//PUserResponse represents user response entity used by the providers
type PUserResponse struct {
	UIN       string `json:"uin"`
	PublicKey string `json:"public_key"`
	Consent   bool   `json:"consent"`
} //@name PUser

type historyResponse struct {
	ID            string    `json:"id"`
	AccountID     string    `json:"account_id"`
	Date          time.Time `json:"date"`
	Type          string    `json:"type"`
	EncryptedKey  string    `json:"encrypted_key"`
	EncryptedBlob string    `json:"encrypted_blob"`
} // @name History

type statusResponse struct {
	ID            string     `json:"id"`
	AccountID     string     `json:"account_id"`
	Date          *time.Time `json:"date"`
	EncryptedKey  string     `json:"encrypted_key"`
	EncryptedBlob string     `json:"encrypted_blob"`
	DateUpdated   *time.Time `json:"date_updated"`
	AppVersion    *string    `json:"app_version"`
} // @name Status

func convertToDaysOfOperations(list []locationOperationDayRequest) []model.OperationDay {
	var doo []model.OperationDay
	if list != nil {
		for _, d := range list {
			item := model.OperationDay{Name: d.Name, OpenTime: d.OpenTime, CloseTime: d.CloseTime}
			doo = append(doo, item)
		}
	}
	return doo
}

func convertFromDaysOfOperations(list []model.OperationDay) []locationOperationDayResponse {
	var doo []locationOperationDayResponse
	if list != nil {
		for _, d := range list {
			item := locationOperationDayResponse{Name: d.Name, OpenTime: d.OpenTime, CloseTime: d.CloseTime}
			doo = append(doo, item)
		}
	}
	return doo
}
