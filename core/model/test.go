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

//CTest represents encrypted provider test
type CTest struct {
	ID string `json:"id" bson:"_id"`

	ProviderID string `json:"provider_id" bson:"provider_id"`
	UserID     string `json:"user_id" bson:"user_id"`

	EncryptedKey  string `json:"encrypted_key" bson:"encrypted_key"`
	EncryptedBlob string `json:"encrypted_blob" bson:"encrypted_blob"`

	OrderNumber *string `json:"order_number" bson:"order_number"`

	Processed bool `json:"processed" bson:"processed"`

	DateCreated time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time `json:"date_updated" bson:"date_updated"`
}

//EManualTest represents manual test
type EManualTest struct {
	ID        string
	User      User
	AccountID string
	HistoryID string

	LocationID *string
	CountyID   *string

	EncryptedKey  string
	EncryptedBlob string

	EncryptedImageKey  string
	EncryptedImageBlob string

	Image  string
	Status string //unverified, verified, rejected
	Date   time.Time
}
