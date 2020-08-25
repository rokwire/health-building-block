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

import (
	"health/core"
	"log"
	"strconv"
	"time"
)

//Adapter implements the Audit interface
type Adapter struct {
	db *database
}

//Start starts the audit
func (sa *Adapter) Start() error {
	err := sa.db.start()
	return err
}

//Log logs an item
func (a *Adapter) Log(entity core.AuditEntity) {
	go func(e core.AuditEntity) {
		_, err := a.db.audit.InsertOne(&e)
		if err != nil {
			log.Printf("error audit logging - %s", err.Error())
		}
	}(entity)
}

//Find finds items
func (a *Adapter) Find() ([]core.AuditEntity, error) {
	//TODO
	return nil, nil
}

//NewAuditAdapter creates a new audit adapter instance
func NewAuditAdapter(mongoDBAuth string, mongoDBName string, mongoTimeout string) *Adapter {
	timeout, err := strconv.Atoi(mongoTimeout)
	if err != nil {
		log.Println("Audit - Set default timeout - 500")
		timeout = 500
	}
	timeoutMS := time.Millisecond * time.Duration(timeout)

	db := &database{mongoDBAuth: mongoDBAuth, mongoDBName: mongoDBName, mongoTimeout: timeoutMS}
	return &Adapter{db: db}
}
