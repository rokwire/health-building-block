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

package profilebb

import (
	"encoding/json"
	"errors"
	"fmt"
	"health/core"
	"io/ioutil"
	"log"
	"net/http"
)

//Adapter implements the ProfileBuildingBlock interface
type Adapter struct {
	host   string
	apiKey string
}

//LoadUserData loads the user data by uuid
func (a *Adapter) LoadUserData(uuid string) (*core.ProfileUserData, error) {
	url := fmt.Sprintf("%s/%s", a.host, uuid)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error creating load user data request - %s", err)
		return nil, err
	}
	req.Header.Set("ROKWIRE-API-KEY", a.apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error loading user data - %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("error with response code - %d", resp.StatusCode)
		return nil, errors.New("error with response code != 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading the body data for the loading user data request - %s", err)
		return nil, err
	}

	var result core.ProfileUserData
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("error converting data for the loading user data request - %s", err)
		return nil, err
	}

	return &result, nil
}

//NewProfileBBAdapter creates a new profile building block adapter instance
func NewProfileBBAdapter(profileHost string, profileAPIKey string) *Adapter {
	return &Adapter{host: profileHost, apiKey: profileAPIKey}
}
