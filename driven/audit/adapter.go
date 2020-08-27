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
	"bytes"
	"fmt"
	"health/core"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

//LogCreateEvent logs a create event item
func (a *Adapter) LogCreateEvent(userIdentifier string, userInfo string, usedGroup string, entity string, entityID string, data []core.AuditDataEntry) {
	go func(userIdentifier string, userInfo string, usedGroup string, entity string, entityID string) {
		dataFormatted := a.prepareData(data)
		auditEntity := core.AuditEntity{UserIdentifier: userIdentifier, UserInfo: userInfo,
			UsedGroup: usedGroup, Entity: entity, EntityID: entityID,
			Operation: "create", Data: dataFormatted, CreatedAt: time.Now()}

		a.log(auditEntity)

	}(userIdentifier, userInfo, usedGroup, entity, entityID)
}

//LogUpdateEvent logs an update event item
func (a *Adapter) LogUpdateEvent(userIdentifier string, userInfo string, usedGroup string, entity string, entityID string, data []core.AuditDataEntry) {
	go func(userIdentifier string, userInfo string, usedGroup string, entity string, entityID string) {
		dataFormatted := a.prepareData(data)
		auditEntity := core.AuditEntity{UserIdentifier: userIdentifier, UserInfo: userInfo,
			UsedGroup: usedGroup, Entity: entity, EntityID: entityID,
			Operation: "update", Data: dataFormatted, CreatedAt: time.Now()}

		a.log(auditEntity)

	}(userIdentifier, userInfo, usedGroup, entity, entityID)
}

func (a *Adapter) prepareData(data []core.AuditDataEntry) *string {
	if len(data) <= 0 {
		return nil
	}

	var b bytes.Buffer
	i := 0
	count := len(data)
	for _, current := range data {
		res := fmt.Sprintf("%s:%s", current.Key, current.Value)
		b.WriteString(res)

		if i < (count - 1) {
			b.WriteString(", ")
		}
		i++
	}
	dataFormatted := b.String()
	return &dataFormatted
}

//LogDeleteEvent logs a delete event item
func (a *Adapter) LogDeleteEvent(userIdentifier string, userInfo string, usedGroup string, entity string, entityID string) {
	go func(userIdentifier string, userInfo string, usedGroup string, entity string, entityID string) {
		auditEntity := core.AuditEntity{UserIdentifier: userIdentifier, UserInfo: userInfo,
			UsedGroup: usedGroup, Entity: entity, EntityID: entityID,
			Operation: "delete", Data: nil, CreatedAt: time.Now()}

		a.log(auditEntity)

	}(userIdentifier, userInfo, usedGroup, entity, entityID)
}

func (a *Adapter) log(entity core.AuditEntity) {
	_, err := a.db.audit.InsertOne(&entity)
	if err != nil {
		log.Printf("error audit logging - %s", err.Error())
	}
}

//Find finds items
func (a *Adapter) Find(sortBy *string, asc *bool) ([]*core.AuditEntity, error) {
	options := options.Find()

	//add sort
	if sortBy != nil && asc != nil {
		ascValue := -1
		if *asc {
			ascValue = 1
		}
		options.SetSort(bson.D{primitive.E{Key: *sortBy, Value: ascValue}})
	}

	var result []*core.AuditEntity
	err := a.db.audit.Find(nil, &result, options)
	if err != nil {
		return nil, err
	}
	return result, nil
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
