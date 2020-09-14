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

//Symptom represents symptom entity
type Symptom struct {
	ID   string
	Name string

	Group SymptomGroup
}

//SymptomGroup represents symptom group entity
type SymptomGroup struct {
	ID   string
	Name string

	Symptoms []Symptom
}

//Symptoms represents raw symptoms entity for e specific app version
type Symptoms struct {
	AppVersion string `json:"app_version" bson:"app_version"`
	Items      string `json:"items" bson:"items"`
}
