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

//EHistory represents status history entity which keeps encrypted data
type EHistory struct {
	ID            string    `json:"id" bson:"_id"`
	UserID        string    `json:"user_id" bson:"user_id"`
	Date          time.Time `json:"date" bson:"date"`
	Type          string    `json:"type" bson:"type"`
	EncryptedKey  string    `json:"encrypted_key" bson:"encrypted_key"`
	EncryptedBlob string    `json:"encrypted_blob" bson:"encrypted_blob"`
} // @name History
