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

//Provider represents provider entity
type Provider struct {
	ID         string
	Name       string
	ManualTest bool

	Locations           []Location
	AvailableMechanisms []string
}

//Location represents provider location entity
type Location struct {
	ID string

	Name            string
	Address1        string
	Address2        string
	City            string
	State           string
	ZIP             string
	Country         string
	Latitude        float64
	Longitude       float64
	Contact         string //phone
	DaysOfOperation []OperationDay
	URL             string
	Notes           string

	Provider Provider
	County   County

	AvailableTests []TestType
}

//OperationDay represents a day from the week saying the operation hours
type OperationDay struct {
	Name      string
	OpenTime  string
	CloseTime string
}
