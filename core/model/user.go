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

package model

import (
	"time"
)

//User represents user entity
type User struct {
	ID string `json:"id" bson:"_id"`

	ShibbolethAuth *ShibbolethAuth `json:"shibboleth_auth" bson:"shibboleth_auth"`

	ExternalID string `json:"external_id" bson:"external_id"`

	UUID                 string  `json:"uuid" bson:"uuid"`
	PublicKey            string  `json:"public_key" bson:"public_key"`
	Consent              bool    `json:"consent" bson:"consent"`
	ExposureNotification bool    `json:"exposure_notification" bson:"exposure_notification"`
	RePost               bool    `json:"re_post" bson:"re_post"`
	EncryptedKey         *string `json:"encrypted_key" bson:"encrypted_key"`
	EncryptedBlob        *string `json:"encrypted_blob" bson:"encrypted_blob"`

	Accounts []Account `json:"accounts" bson:"accounts"`

	DateCreated time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time `json:"date_updated" bson:"date_updated"`
}

//GetAccount gives the user account for the provided account id
func (user User) GetAccount(accountID string) *Account {
	if len(user.Accounts) == 0 {
		return nil
	}

	for _, current := range user.Accounts {
		if current.ID == accountID {
			return &current
		}
	}
	return nil
}

//GetDefaultAccount gives the default user account
func (user User) GetDefaultAccount() *Account {
	if len(user.Accounts) == 0 {
		return nil
	}

	for _, current := range user.Accounts {
		if current.Default {
			return &current
		}
	}
	return nil
}

//HasDefaultAccount says if the user has a default account
func (user User) HasDefaultAccount() bool {
	defaultAccount := user.GetDefaultAccount()
	if defaultAccount != nil {
		return true
	}
	return false
}

//IsAdmin says if the user is admin
func (user User) IsAdmin() bool {
	if user.ShibbolethAuth == nil {
		return false
	}
	isMemberOfList := user.ShibbolethAuth.IsMemberOf
	if isMemberOfList == nil {
		return false
	}
	for _, group := range *isMemberOfList {
		if group == "urn:mace:uiuc.edu:urbana:authman:app-rokwire-service-policy-rokwire admin app" {
			return true
		}
	}
	return false
}

//IsPublicHealth says if the user is public health
func (user User) IsPublicHealth() bool {
	if user.ShibbolethAuth == nil {
		return false
	}
	isMemberOfList := user.ShibbolethAuth.IsMemberOf
	if isMemberOfList == nil {
		return false
	}
	for _, group := range *isMemberOfList {
		if group == "urn:mace:uiuc.edu:urbana:authman:app-rokwire-service-policy-rokwire public health" {
			return true
		}
	}
	return false
}

//IsMemberOf says if the user is member of a group
func (user User) IsMemberOf(group string) bool {
	if user.ShibbolethAuth == nil {
		return false
	}
	isMemberOfList := user.ShibbolethAuth.IsMemberOf
	if isMemberOfList == nil {
		return false
	}
	for _, current := range *isMemberOfList {
		if current == group {
			return true
		}
	}
	return false
}

//GetLogData gives the user audit log data
func (user User) GetLogData() (string, string) {
	if user.ShibbolethAuth == nil {
		return "", ""
	}

	userIdentifier := user.ID
	userInfo := user.ShibbolethAuth.Email

	return userIdentifier, userInfo
}

//ShibbolethAuth represents shibboleth auth entity
type ShibbolethAuth struct {
	Uin        string    `json:"uiucedu_uin" bson:"uiucedu_uin"`
	Email      string    `json:"email" bson:"email"`
	IsMemberOf *[]string `json:"uiucedu_is_member_of" bson:"uiucedu_is_member_of"`
}

//Account represents account entity
type Account struct {
	ID         string `json:"id" bson:"id"`
	ExternalID string `json:"external_id" bson:"external_id"`
	Default    bool   `json:"default" bson:"default"`
	Active     bool   `json:"active" bson:"active"`

	FirstName  string `json:"first_name" bson:"first_name"`
	MiddleName string `json:"middle_name" bson:"middle_name"`
	LastName   string `json:"last_name" bson:"last_name"`
	BirthDate  string `json:"birth_date" bson:"birth_date"`
	Gender     string `json:"gender" bson:"gender"`
	Address1   string `json:"address1" bson:"address1"`
	Address2   string `json:"address2" bson:"address2"`
	Address3   string `json:"address3" bson:"address3"`
	City       string `json:"city" bson:"city"`
	State      string `json:"state" bson:"state"`
	ZipCode    string `json:"zip_code" bson:"zip_code"`
	Phone      string `json:"phone" bson:"phone"`
	Email      string `json:"email" bson:"email"`
}

//RawSubAccount represents raw sub account entity
type RawSubAccount struct {
	UIN        string `json:"uin" bson:"uin"`
	FirstName  string `json:"first_name" bson:"first_name"`
	MiddleName string `json:"middle_name" bson:"middle_name"`
	LastName   string `json:"last_name" bson:"last_name"`
	BirthDate  string `json:"birth_date" bson:"birth_date"`
	Gender     string `json:"gender" bson:"gender"`
	Address1   string `json:"address1" bson:"address1"`
	Address2   string `json:"address2" bson:"address2"`
	Address3   string `json:"address3" bson:"address3"`
	City       string `json:"city" bson:"city"`
	State      string `json:"state" bson:"state"`
	ZipCode    string `json:"zip_code" bson:"zip_code"`
	Phone      string `json:"phone" bson:"phone"`
	NetID      string `json:"net_id" bson:"net_id"`
	Email      string `json:"email" bson:"email"`

	PrimaryAccount string `json:"primary_account" bson:"primary_account"`

	AccountID string `json:"account_id" bson:"account_id"`
}
