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
	"health/core/model"
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

func (a *Adapter) LogCreateEvent(user model.User, entity string, entityID string) {
	go func(e *core.AuditEntity) {
		if e == nil {
			log.Printf("cannot log for nil entity")
			return
		}
		_, err := a.db.audit.InsertOne(e)
		if err != nil {
			log.Printf("error audit logging - %s", err.Error())
		}
	}(entity)

	/*if user.ShibbolethAuth == nil {
		return nil
	}

	//TODO groups
	var groups []string
	return &AuditEntity{UserIdentifier: user.ID, UserInfo: user.ShibbolethAuth.Email,
		UserGroups: groups, Entity: "county", EntityID: entityID,
		Operation: "create", Change: nil, CreatedAt: time.Now()} */
}

//Log logs an item
func (a *Adapter) log(entity *core.AuditEntity) {
	if e == nil {
		log.Printf("cannot log for nil entity")
		return
	}
	_, err := a.db.audit.InsertOne(e)
	if err != nil {
		log.Printf("error audit logging - %s", err.Error())
	}
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
