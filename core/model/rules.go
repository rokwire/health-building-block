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

//Rule represents county and test types rule item
type Rule struct {
	ID       string
	County   County
	TestType TestType
	Priority *int

	ResultsStates []TestTypeResultCountyStatus
}

//TestTypeResultCountyStatus represents test type result and county status mapping
type TestTypeResultCountyStatus struct {
	TestTypeResultID string
	CountyStatusID   string
}

//SymptomRule represents a symptom rule entity
type SymptomRule struct {
	ID     string
	County County

	Gr1Count int
	Gr2Count int

	Items []SymptomRuleItem
}

//CRules represents all rules in a raw format
type CRules struct {
	AppVersion string `json:"app_version" bson:"app_version"`
	CountyID   string `json:"county_id" bson:"county_id"`
	Data       string `json:"data" bson:"data"`
}

//SymptomRuleItem represents a symptom rule item entity
type SymptomRuleItem struct {
	Gr1          bool
	Gr2          bool
	CountyStatus CountyStatus
	NextStep     string
}

//AccessRule represents an access rule entity
type AccessRule struct {
	ID     string
	County County
	Rules  []AccessRuleCountyStatus
}

//AccessRuleCountyStatus represents mapping for "granted"/"denied" county status
type AccessRuleCountyStatus struct {
	CountyStatusID string
	Value          string //granted or denied
}
