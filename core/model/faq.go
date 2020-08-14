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
	"sort"
	"time"
)

//FAQ represents FAQ entity
type FAQ struct {
	DateUpdate time.Time   `json:"dateUpdated" bson:"dateUpdated"`
	General    *[]*General `json:"general" bson:"general"`
	Sections   *[]*Section `json:"sections" bson:"sections"`
} // @name FAQ

//Sort sorts the sections and the question within the sections based on the display order field
func (faq FAQ) Sort() {
	sections := faq.Sections
	if sections == nil {
		return
	}

	//sort the questions within a section
	for _, s := range *sections {
		questions := s.Questions
		if questions != nil {
			sort.Slice(*questions, func(i, j int) bool {
				return (*questions)[i].DisplayOrder < (*questions)[j].DisplayOrder
			})
		}
	}
	//sort the sections
	sort.Slice(*sections, func(i, j int) bool {
		return (*sections)[i].DisplayOrder < (*sections)[j].DisplayOrder
	})

}

//General represents general section entity
type General struct {
	Title       string  `json:"title" bson:"title"`
	Description string  `json:"description" bson:"description"`
	Link        *string `json:"link" bson:"link"`
} // @name FAQGeneral

//Section represents section entity
type Section struct {
	ID           string       `json:"id" bson:"id"`
	Title        string       `json:"title" bson:"title"`
	DisplayOrder int          `json:"display_order" bson:"display_order"`
	Questions    *[]*Question `json:"questions" bson:"questions"`
} // @name FAQSection

//Question represents section question entity
type Question struct {
	ID           string  `json:"id" bson:"id"`
	Title        string  `json:"title" bson:"title"`
	Description  string  `json:"description" bson:"description"`
	Link         *string `json:"link" bson:"link"`
	DisplayOrder int     `json:"display_order" bson:"display_order"`
} // @name FAQQuestion
