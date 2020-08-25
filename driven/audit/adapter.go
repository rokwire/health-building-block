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

package audit

//Adapter implements the Audit interface
type Adapter struct {
}

//Log logs an item
func (a *Adapter) Log(entity AuditEntity) error {

}

//Find finds items
func (a *Adapter) Find() ([]AuditEntity, error) {

}

//NewAuditAdapter creates a new audit adapter instance
func NewAuditAdapter() *Adapter {
	return &Adapter{}
}
