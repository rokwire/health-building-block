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

import "time"

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

	DateCreated time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time `json:"date_updated" bson:"date_updated"`
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
