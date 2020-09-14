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

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"health/core"
	"health/core/model"
	"health/utils"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

//ApisHandler handles the rest APIs implementation
type ApisHandler struct {
	app *core.Application
}

//Version gives the service version
func (h ApisHandler) Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.app.Services.GetVersion()))
}

//ClearUserData clears everything in the service for that user as if the user has never seen the service before
// @Description Clears everything for that user as if the user has never seen the service before.
// @Tags Covid19
// @ID clearUserData
// @Accept plain
// @Success 200 {object} string "Successfully cleared"
// @Security AppUserAuth
// @Router /covid19/user/clear [get]
func (h ApisHandler) ClearUserData(current model.User, w http.ResponseWriter, r *http.Request) {
	err := h.app.Services.ClearUserData(current)
	if err != nil {
		log.Printf("error on clearing the user data - %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully cleared"))
}

type getUserByShibbolethIDResponse struct {
	PublicKey string `json:"public_key"`
	Consent   bool   `json:"consent"`
} // @name GetUserByShibbolethUINResponse

//GetUserByShibbolethUIN gives the user info needed for the providers
// @Description Gives the user info needed for the providers
// @Tags Providers
// @ID getUserByShibbolethUIN
// @Accept json
// @Param id path string true "User ID"
// @Success 200 {object} getUserByShibbolethIDResponse
// @Security ProvidersAuth
// @Router /covid19/users/uin/{id} [get]
func (h ApisHandler) GetUserByShibbolethUIN(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shibbolethUIN := params["uin"]
	if len(shibbolethUIN) <= 0 {
		log.Println("uin is required")
		http.Error(w, "uin is required", http.StatusBadRequest)
		return
	}

	user, err := h.app.Services.GetUserByShibbolethUIN(shibbolethUIN)
	if err != nil {
		log.Printf("Error on getting user by shibboleth id %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if user == nil {
		//return not found
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	result := getUserByShibbolethIDResponse{PublicKey: user.PublicKey, Consent: user.Consent}

	data, err := json.Marshal(result)
	if err != nil {
		log.Println("Error on marshal getUserByShibbolethIDResponse")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetUsersForRePost gives the users for re-posting the test results
// @Description Gives the users for re-posting the test results
// @Tags Providers
// @ID GetUsersForRePost
// @Accept json
// @Success 200 {array} PUserResponse
// @Security ProvidersAuth
// @Router /covid19/users/re-post [get]
func (h ApisHandler) GetUsersForRePost(w http.ResponseWriter, r *http.Request) {
	users, err := h.app.Services.GetUsersForRePost()
	if err != nil {
		log.Printf("Error on getting users for re-post %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var result []PUserResponse
	if len(users) <= 0 {
		result = make([]PUserResponse, 0)
	} else {
		for _, user := range users {
			pUser := PUserResponse{UIN: user.ExternalID, Consent: user.Consent, PublicKey: user.PublicKey}
			result = append(result, pUser)
		}
	}

	data, err := json.Marshal(result)
	if err != nil {
		log.Println("Error on marshal PUserResponse")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createCTestRequest struct {
	ProviderID    string  `json:"provider_id" validate:"required"`
	UIN           string  `json:"uin" validate:"required"`
	EncryptedKey  string  `json:"encrypted_key" validate:"required"`
	EncryptedBlob string  `json:"encrypted_blob" validate:"required"`
	OrderNumber   *string `json:"order_number"`
} // @name createCTestRequest

//CreateExternalCTest creates CTest
// @Description Creates CTest.
// @Tags Providers
// @ID createCTest
// @Accept json
// @Produce json
// @Param data body createCTestRequest true "body data"
// @Success 200 {object} string "Successfully created"
// @Security ProvidersAuth
// @Router /covid19/ctests [post]
func (h ApisHandler) CreateExternalCTest(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a ctest - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createCTestRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create ctest request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create ctest data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	providerID := requestData.ProviderID
	uin := requestData.UIN
	encryptedKey := requestData.EncryptedKey
	encryptedBlob := requestData.EncryptedBlob
	orderNumber := requestData.OrderNumber

	err = h.app.Services.CreateExternalCTest(providerID, uin, encryptedKey, encryptedBlob, orderNumber)
	if err != nil {
		log.Printf("Error on creating a ctest - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully created"))
}

type getMCountyResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	StateProvince string `json:"state_province"`
	Country       string `json:"country"`

	CountyStatuses []getMCountyCountyStatusResponse `json:"county_statuses"`
	Guidelines     []getMCountyGuidelineResponse    `json:"guidelines"`
} // @name County

type getMCountyCountyStatusResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Descpription string `json:"descpription"`
} // @name CountyStatus

type getMCountyGuidelineResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Items []getMCountyGuidelineItemResponse `json:"items"`
} // @name Guideline

type getMCountyGuidelineItemResponse struct {
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Type        string `json:"type"`
} // @name GuidelineItem

type gubonResponse map[string]*string // @name gubonResponse

//GetUINsByOrderNumbers gives the corresponding UINs for the provided order numbers list
// @Description Gives the corresponding UINs for the provided order numbers list. The list must be comma separated. The response looks like {"ordernumber1":"uin 1","ordernumber2":"uin 2"}
// @Tags Providers
// @ID GetUINsByOrderNumbers
// @Accept json
// @Param order-numbers query string true "Comma separated - ordernumber1,ordernumber2"
// @Success 200 {object} gubonResponse
// @Security ProvidersAuth
// @Router /covid19/track/uins [get]
func (h ApisHandler) GetUINsByOrderNumbers(w http.ResponseWriter, r *http.Request) {
	orderNumbersKeys, ok := r.URL.Query()["order-numbers"]
	if !ok || len(orderNumbersKeys[0]) < 1 {
		log.Println("url param 'order-numbers' is missing")
		return
	}
	orderNumbersKey := orderNumbersKeys[0]
	orderNumbers := strings.Split(orderNumbersKey, ",")
	if len(orderNumbers) == 0 {
		http.Error(w, "order-numbers is required", http.StatusBadRequest)
		return
	}

	var resData gubonResponse
	resData, err := h.app.Services.GetUINsByOrderNumbers(orderNumbers)
	if err != nil {
		log.Printf("Error on getting UINs - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Println("Error on marshal the uins by order numbers")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type ilbuResponseItem struct {
	OrderNumber *string   `json:"order_number"`
	DateCreated time.Time `json:"date_created"`
} // @name ilbuResponseItem

type ilbuResponse map[string][]ilbuResponseItem // @name ilbuResponse

//GetItemsListsByUINs gives the tracks items list for the provided UINs
// @Description Gives the items list for the provided UINs. The list must be comma separated. The response looks like {"”777778":[{"order_number":null,"date_created":"2020-08-12T05:52:47.467Z”},…],”777777":[{"order_number":"9","date_created":"2020-09-10T05:02:14.716Z"}]}
// @Tags Providers
// @ID GetItemsListsByUINs
// @Accept json
// @Param uins query string true "Comma separated - uin1,uin2"
// @Success 200 {object} ilbuResponse
// @Security ProvidersAuth
// @Router /covid19/track/items [get]
func (h ApisHandler) GetItemsListsByUINs(w http.ResponseWriter, r *http.Request) {
	uinsKeys, ok := r.URL.Query()["uins"]
	if !ok || len(uinsKeys[0]) < 1 {
		log.Println("url param 'uins' is missing")
		return
	}
	uinsKey := uinsKeys[0]
	uins := strings.Split(uinsKey, ",")
	if len(uins) == 0 {
		http.Error(w, "uins is required", http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.GetCTestsByExternalUserIDs(uins)
	if err != nil {
		log.Printf("Error on getting track items by external id - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//prepare the response
	responseData := make(ilbuResponse, len(resData))
	for key, currentList := range resData {
		list := resData[key]

		if list == nil {
			continue
		}

		var resList []ilbuResponseItem
		for _, item := range currentList {
			resList = append(resList, ilbuResponseItem{OrderNumber: item.OrderNumber, DateCreated: item.DateCreated})
		}
		responseData[key] = resList
	}

	data, err := json.Marshal(responseData)
	if err != nil {
		log.Println("Error on marshal the track items by uins")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetCounty gets a county
// @Description Gets a county
// @Tags Covid19
// @ID getCounty
// @Accept json
// @Param id path string true "ID"
// @Success 200 {object} getMCountyResponse
// @Security RokwireAuth
// @Router /covid19/counties/{id} [get]
func (h ApisHandler) GetCounty(appVersion *string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("id is required")
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	county, err := h.app.Services.GetCounty(ID)
	if err != nil {
		log.Printf("Error on getting the counties items - %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if county == nil {
		//not found
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	//county statuses
	var countyStatuses []getMCountyCountyStatusResponse
	if county.CountyStatuses != nil {
		for _, cs := range county.CountyStatuses {
			csItem := getMCountyCountyStatusResponse{ID: cs.ID, Name: cs.Name, Descpription: cs.Description}
			countyStatuses = append(countyStatuses, csItem)
		}
	}

	//guidelines
	var guidelines []getMCountyGuidelineResponse
	if county.Guidelines != nil {
		for _, gl := range county.Guidelines {
			var glItems []getMCountyGuidelineItemResponse
			if gl.Items != nil {
				for _, inner := range gl.Items {
					innerItem := getMCountyGuidelineItemResponse{Icon: inner.Icon,
						Description: inner.Description, Type: inner.Type.Value}
					glItems = append(glItems, innerItem)
				}
			}

			item := getMCountyGuidelineResponse{ID: gl.ID, Name: gl.Name, Description: gl.Description, Items: glItems}
			guidelines = append(guidelines, item)
		}
	}

	responseItem := getMCountyResponse{ID: county.ID, Name: county.Name, StateProvince: county.StateProvince,
		Country: county.Country, CountyStatuses: countyStatuses, Guidelines: guidelines}
	data, err := json.Marshal(responseItem)
	if err != nil {
		log.Println("Error on marshal the county items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type getMCountiesResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	StateProvince string `json:"state_province"`
	Country       string `json:"country"`

	CountyStatuses []getMCountiesCountyStatusResponse `json:"county_statuses"`
	Guidelines     []getMCountiesGuidelineResponse    `json:"guidelines"`
}

type getMCountiesCountyStatusResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Descpription string `json:"descpription"`
}

type getMCountiesGuidelineResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Items []getMCountiesGuidelineItemResponse `json:"items"`
}

type getMCountiesGuidelineItemResponse struct {
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

//GetCounties gets all the counties
// @Description Gives the counties. Optionally you can filter the results by one or many of the following three fields - name, state_province and country
// @Tags Covid19
// @ID GetCounties
// @Accept json
// @Param name query string false "name"
// @Param state_province query string false "State province"
// @Param country query string false "Country"
// @Success 200 {array} getMCountyResponse
// @Security RokwireAuth
// @Router /covid19/counties [get]
func (h ApisHandler) GetCounties(appVersion *string, w http.ResponseWriter, r *http.Request) {
	filter := utils.ConstructFilter(r)
	counties, err := h.app.Services.FindCounties(filter)
	if err != nil {
		log.Println("Error on getting the counties items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseList := make([]getMCountiesResponse, len(counties))
	if counties != nil {
		for i, county := range counties {
			//county statuses
			var countyStatuses []getMCountiesCountyStatusResponse
			if county.CountyStatuses != nil {
				for _, cs := range county.CountyStatuses {
					csItem := getMCountiesCountyStatusResponse{ID: cs.ID, Name: cs.Name, Descpription: cs.Description}
					countyStatuses = append(countyStatuses, csItem)
				}
			}

			//guidelines
			var guidelines []getMCountiesGuidelineResponse
			if county.Guidelines != nil {
				for _, gl := range county.Guidelines {
					var glItems []getMCountiesGuidelineItemResponse
					if gl.Items != nil {
						for _, inner := range gl.Items {
							innerItem := getMCountiesGuidelineItemResponse{Icon: inner.Icon,
								Description: inner.Description, Type: inner.Type.Value}
							glItems = append(glItems, innerItem)
						}
					}

					item := getMCountiesGuidelineResponse{ID: gl.ID, Name: gl.Name, Description: gl.Description, Items: glItems}
					guidelines = append(guidelines, item)
				}
			}

			county := getMCountiesResponse{ID: county.ID, Name: county.Name, StateProvince: county.StateProvince,
				Country: county.Country, CountyStatuses: countyStatuses, Guidelines: guidelines}
			responseList[i] = county
		}
	}
	data, err := json.Marshal(responseList)
	if err != nil {
		log.Println("Error on marshal the counties items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type getCTestsResponse struct {
	ID string `json:"id"`

	ProviderID   string `json:"provider_id"`
	ProviderName string `json:"provider"`
	UserID       string `json:"user_id"`

	EncryptedKey  string `json:"encrypted_key"`
	EncryptedBlob string `json:"encrypted_blob"`

	Processed bool `json:"processed"`

	DateCreated time.Time  `json:"date_created"`
	DateUpdated *time.Time `json:"date_updated"`
}

//GetCTests gets CTests for the current user.
// @Description Gets not processed ctests for a user.
// @Tags Covid19
// @ID getCTests
// @Accept  json
// @Param processed query bool false "select false value"
// @Success 200 {array} model.CTest
// @Security AppUserAuth
// @Router /covid19/ctests [get]
func (h ApisHandler) GetCTests(current model.User, w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["processed"]
	if !ok || len(keys[0]) < 1 {
		log.Println("url param 'processed' is missing")
		http.Error(w, "url param 'processed' is missing", http.StatusBadRequest)
		return
	}
	processed, _ := strconv.ParseBool(keys[0])

	ctests, providers, err := h.app.Services.GetCTests(current, processed)
	if err != nil {
		log.Println("Error on getting the ctests items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	ctestsLen := len(ctests)
	resultList := make([]getCTestsResponse, ctestsLen)
	if ctestsLen > 0 {
		for i, ctest := range ctests {
			provider := h.findProvider(ctest.ProviderID, providers)

			r := getCTestsResponse{ID: ctest.ID, ProviderID: provider.ID, ProviderName: provider.Name,
				UserID: ctest.UserID, EncryptedKey: ctest.EncryptedKey, EncryptedBlob: ctest.EncryptedBlob,
				Processed: ctest.Processed, DateCreated: ctest.DateCreated, DateUpdated: ctest.DateUpdated}
			resultList[i] = r
		}
	}
	data, err := json.Marshal(resultList)
	if err != nil {
		log.Println("Error on marshal the ctests items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateCTestRequest struct {
	Processed bool `json:"processed" validate:"required"`
} // @name updateCTestRequest

//UpdateCTest updates a CTests
// @Description  Mark ctest as processed.
// @Tags Covid19
// @ID updateCTest
// @Accept json
// @Produce json
// @Param data body updateCTestRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} model.CTest
// @Security AppUserAuth
// @Router /covid19/ctests/{id} [put]
func (h ApisHandler) UpdateCTest(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("id is required")
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update ctest item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateCTestRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update ctest item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update ctest data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctest, err := h.app.Services.UpdateCTest(current, ID, requestData.Processed)
	if err != nil {
		log.Printf("Error on updating the ctest item - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err = json.Marshal(ctest)
	if err != nil {
		log.Println("Error on marshal the ctest item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteCTests deletes CTests for the current user.
// @Description Deletes all ctests for a user
// @Tags Covid19
// @ID deleteCTests
// @Accept plain
// @Success 200 {object} string "Successfuly deleted [n] items"
// @Security AppUserAuth
// @Router /covid19/ctests [delete]
func (h ApisHandler) DeleteCTests(current model.User, w http.ResponseWriter, r *http.Request) {
	deletedCount, err := h.app.Services.DeleteCTests(current.ID)
	if err != nil {
		log.Printf("Error on deleting the ctests items - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	result := fmt.Sprintf("Successfuly deleted %d items", deletedCount)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

//GetResources gives the covid19 resources
// @Description Gives the covid19 resources
// @Tags Covid19
// @ID getResources
// @Accept  json
// @Success 200 {array} model.Resource
// @Security RokwireAuth
// @Router /covid19/resources [get]
func (h ApisHandler) GetResources(appVersion *string, w http.ResponseWriter, r *http.Request) {
	resources, err := h.app.Services.GetResources()
	if err != nil {
		log.Printf("Error on getting resources %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//sort
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].DisplayOrder < resources[j].DisplayOrder
	})

	data, err := json.Marshal(resources)
	if err != nil {
		log.Println("Error on marshal the resources")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetFAQ gives the covid19 FAQs
// @Description Gives the covid19 FAQs. The sections are sorted by the display order field. The questions within a section are sorted by the display order field.
// @Tags Covid19
// @ID getFAQ
// @Accept json
// @Success 200 {array} model.FAQ
// @Security RokwireAuth
// @Router /covid19/faq [get]
func (h ApisHandler) GetFAQ(appVersion *string, w http.ResponseWriter, r *http.Request) {
	faq, err := h.app.Services.GetFAQ()
	if err != nil {
		log.Printf("Error on getting faq %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//sort
	faq.Sort()

	data, err := json.Marshal(faq)
	if err != nil {
		log.Println("Error on marshal the faq")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetNews gives the covid19 news
// @Description Gives the covid19 news.
// @Tags Covid19
// @ID GetNews
// @Accept json
// @Success 200 {array} model.News
// @Security RokwireAuth
// @Router /covid19/news [get]
func (h ApisHandler) GetNews(appVersion *string, w http.ResponseWriter, r *http.Request) {
	var limit int64
	var err error

	limParam := r.URL.Query().Get("limit")
	if len(limParam) > 0 {
		limit, err = strconv.ParseInt(limParam, 10, 64)
		if err != nil {
			http.Error(w, "limit must be a number", http.StatusBadRequest)
		}
	}

	news, err := h.app.Services.GetNews(limit)
	if err != nil {
		log.Printf("Error on getting news %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(news)
	if err != nil {
		log.Println("Error on marshal the news")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type getStatusByCountyResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	HealthStatus string     `json:"health_status"`
	Date         *time.Time `json:"date"`
	NextStep     *string    `json:"next_step"`
	NextStepDate *time.Time `json:"next_step_date"`
	URL          *string    `json:"url"`
}

//GetStatusV2Deprecated gets the status for the current user
// @Deprecated
// @Description Gets the status for the current user
// @Tags Covid19
// @ID GetStatusV2Deprecated
// @Accept  json
// @Success 200 {object} model.EStatus
// @Success 404 {string} Not Found
// @Security AppUserAuth
// @Router /covid19/v2/statuses [get]
func (h ApisHandler) GetStatusV2Deprecated(current model.User, w http.ResponseWriter, r *http.Request) {
	h.processGetStatusV2(current, nil, w, r)
}

//GetStatusV2 gets the status for the current user for a specific app version
// @Description Gets the status for the current user for a specific app version
// @Tags Covid19
// @ID GetStatusV2
// @Accept json
// @Param app-version path string false "App version"
// @Success 200 {object} model.EStatus
// @Success 404 {string} Not Found
// @Security AppUserAuth
// @Router /covid19/v2/app-version/{app-version}/statuses [get]
func (h ApisHandler) GetStatusV2(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	appVersion := params["app-version"]
	h.processGetStatusV2(current, &appVersion, w, r)
}

func (h ApisHandler) processGetStatusV2(current model.User, appVersion *string, w http.ResponseWriter, r *http.Request) {
	status, err := h.app.Services.GetEStatusByUserID(current.ID, appVersion)
	if err != nil {
		log.Println("Error on getting a status")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//not found
	if status == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(status)
	if err != nil {
		log.Println("Error on marshal a status")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createOrUpdateStatusRequestV2 struct {
	Date          *time.Time `json:"date"`
	EncryptedKey  string     `json:"encrypted_key" validate:"required"`
	EncryptedBlob string     `json:"encrypted_blob" validate:"required"`
} // @name createOrUpdateStatusRequest

//CreateOrUpdateStatusV2Deprecated creates or updates the status for the current user
// @Deprecated
// @Description Updates the status for the user. it creates it if not already created.
// @Tags Covid19
// @ID CreateOrUpdateStatusV2Deprecated
// @Accept json
// @Produce json
// @Param data body createOrUpdateStatusRequestV2 true "body data"
// @Success 200 {object} model.EStatus
// @Security AppUserAuth
// @Router /covid19/v2/statuses [put]
func (h ApisHandler) CreateOrUpdateStatusV2Deprecated(current model.User, w http.ResponseWriter, r *http.Request) {
	h.processCreateOrUpdateStatusV2(current, nil, w, r)
}

//CreateOrUpdateStatusV2 creates or updates the status for the current user for a specific app version
// @Description Updates the status for the user for a specific app version. it creates it if not already created.
// @Tags Covid19
// @ID CreateOrUpdateStatusV2
// @Accept json
// @Produce json
// @Param app-version path string false "App version"
// @Param data body createOrUpdateStatusRequestV2 true "body data"
// @Success 200 {object} model.EStatus
// @Security AppUserAuth
// @Router /covid19/v2/app-version/{app-version}/statuses [put]
func (h ApisHandler) CreateOrUpdateStatusV2(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	appVersion := params["app-version"]
	h.processCreateOrUpdateStatusV2(current, &appVersion, w, r)
}

func (h ApisHandler) processCreateOrUpdateStatusV2(current model.User, appVersion *string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create or update a status - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createOrUpdateStatusRequestV2
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create or update status request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create or update status data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	date := requestData.Date
	encryptedKey := requestData.EncryptedKey
	encryptedBlob := requestData.EncryptedBlob

	status, err := h.app.Services.CreateOrUpdateEStatus(current.ID, appVersion, date, encryptedKey, encryptedBlob)
	if err != nil {
		log.Printf("Error on marshal a status - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

//DeleteStatusV2Deprecated deletes a status
// @Deprecated
// @Description Deletes the status for the user.
// @Tags Covid19
// @ID DeleteStatusV2Deprecated
// @Accept plain
// @Success 200 {object} string "Successfully deleted"
// @Security AppUserAuth
// @Router /covid19/v2/statuses [delete]
func (h ApisHandler) DeleteStatusV2Deprecated(current model.User, w http.ResponseWriter, r *http.Request) {
	h.processDeleteStatusV2(current, nil, w, r)
}

//DeleteStatusV2 deletes the status for a specific app version.
// @Description Deletes the status for the user for a specific app version.
// @Tags Covid19
// @ID DeleteStatusV2
// @Accept plain
// @Param app-version path string false "App version"
// @Success 200 {object} string "Successfully deleted"
// @Security AppUserAuth
// @Router /covid19/v2/app-version/{app-version}/statuses [delete]
func (h ApisHandler) DeleteStatusV2(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	appVersion := params["app-version"]
	h.processDeleteStatusV2(current, &appVersion, w, r)
}

func (h ApisHandler) processDeleteStatusV2(current model.User, appVersion *string, w http.ResponseWriter, r *http.Request) {
	err := h.app.Services.DeleteEStatus(current.ID, appVersion)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

type createЕHistory struct {
	Date          time.Time `json:"date" validate:"required"`
	Type          string    `json:"type" validate:"required"`
	EncryptedKey  string    `json:"encrypted_key" validate:"required"`
	EncryptedBlob string    `json:"encrypted_blob" validate:"required"`

	EncryptedImageKey  *string `json:"encrypted_image_key"`
	EncryptedImageBlob *string `json:"encrypted_image_blob"`
	LocationID         *string `json:"location_id"`
	CountyID           *string `json:"county_id"`
} // @name createHistoryRequest

//CreateHistoryV2 creates a new history
// @Description "date", "type", "encrypted_key" and "encrypted_blob" are mandatory fields. When the type is "unverified_manual_test" then the client must pass also "encrypted_image_key", "encrypted_image_blob" and ("location_id" or "county_id").
// @Tags Covid19
// @ID createHistoryV2
// @Accept json
// @Produce json
// @Param data body createЕHistory true "body data"
// @Success 200 {object} model.EHistory
// @Security AppUserAuth
// @Router /covid19/v2/histories [post]
func (h ApisHandler) CreateHistoryV2(current model.User, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a history - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createЕHistory
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create history request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create history data - %s - type - %s - date - %s\n", err.Error(), requestData.Type, requestData.Date)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var history *model.EHistory

	date := requestData.Date
	eType := requestData.Type
	encryptedKey := requestData.EncryptedKey
	encryptedBlob := requestData.EncryptedBlob

	if eType == "unverified_manual_test" {
		encryptedImageKey := requestData.EncryptedImageKey
		encryptedImageBlob := requestData.EncryptedImageBlob
		locationID := requestData.LocationID
		countyID := requestData.CountyID

		err := h.validateManualTestParamsV2(encryptedImageKey, encryptedImageBlob, locationID, countyID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		history, err = h.app.Services.CreateManualЕHistory(current.ID, date, encryptedKey, encryptedBlob, encryptedImageKey, encryptedImageBlob, countyID, locationID)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		history, err = h.app.Services.CreateЕHistory(current.ID, date, eType, encryptedKey, encryptedBlob)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	data, err = json.Marshal(history)
	if err != nil {
		log.Println("Error on marshal a history")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h ApisHandler) validateManualTestParamsV2(encryptedImageKey *string, encryptedImageBlob *string, locationID *string, countyID *string) error {
	//we always need an image
	if encryptedImageKey == nil || encryptedImageBlob == nil {
		return errors.New("manual test requires an image")
	}

	//we need county id or location id
	hasCounty := false
	if countyID != nil {
		hasCounty = true
	}
	hasLocation := false
	if locationID != nil {
		hasLocation = true
	}
	if !hasCounty && !hasLocation {
		return errors.New("required fields - county id or location id")
	}
	if hasCounty && hasLocation {
		return errors.New("required fields - county id or location id, not both")
	}

	return nil
}

func (h ApisHandler) getMapData(key string, mapData map[string]interface{}) *string {
	if mapData[key] != nil {
		value := mapData[key].(string)
		return &value
	}
	return nil
}

type updateЕHistory struct {
	Date          *time.Time `json:"date"`
	EncryptedKey  *string    `json:"encrypted_key"`
	EncryptedBlob *string    `json:"encrypted_blob"`
} // @name updateHistoryRequest

//UpdateHistoryV2 updates the history
// @Description "date", "encrypted_key" and "encrypted_blob" are optional. If a field is omitted then it will not be updated.
// @Tags Covid19
// @ID updateHistoryV2
// @Accept json
// @Produce json
// @Param data body updateЕHistory true "body data"
// @Param id path string true "ID"
// @Success 200 {object} model.EHistory
// @Security AppUserAuth
// @Router /covid19/v2/histories/{id} [put]
func (h ApisHandler) UpdateHistoryV2(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("History id is required")
		http.Error(w, "History id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the ehistorye item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateЕHistory
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the ehistory item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating ehistory data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date := requestData.Date
	encryptedKey := requestData.EncryptedKey
	encryptedBlob := requestData.EncryptedBlob
	history, err := h.app.Services.UpdateEHistory(current.ID, ID, date, encryptedKey, encryptedBlob)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err = json.Marshal(history)
	if err != nil {
		log.Println("Error on marshal the ehistory item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetHistoriesV2 gets all histories for the current user user
// @Description Gets all histories for the current user user
// @Tags Covid19
// @ID getHistoriesV2
// @Accept  json
// @Success 200 {array} model.EHistory
// @Security AppUserAuth
// @Router /covid19/v2/histories [get]
func (h ApisHandler) GetHistoriesV2(current model.User, w http.ResponseWriter, r *http.Request) {
	historiesItems, err := h.app.Services.GetEHistoriesByUserID(current.ID)
	if err != nil {
		log.Println("Error on getting the histories items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(historiesItems) == 0 {
		historiesItems = make([]*model.EHistory, 0)
	}

	data, err := json.Marshal(historiesItems)
	if err != nil {
		log.Println("Error on marshal the histories items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteHistoriesV2 deletes the user histories - for debug purposes only
// @Description Deletes the history items for an user.
// @Tags Covid19
// @ID deleteHistoriesV2
// @Accept plain
// @Success 200 {object} string "Successfully deleted [n] items"
// @Security AppUserAuth
// @Router /covid19/v2/histories [delete]
func (h ApisHandler) DeleteHistoriesV2(current model.User, w http.ResponseWriter, r *http.Request) {
	deletedCount, err := h.app.Services.DeleteEHitories(current.ID)
	if err != nil {
		log.Printf("Error on deleting the ehistories items - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	result := fmt.Sprintf("Successfuly deleted %d items", deletedCount)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

//GetProviders gets the providers
// @Description Gives all the providers
// @Tags Covid19
// @ID getProviders
// @Accept json
// @Success 200 {array} rest.providerResponse
// @Security RokwireAuth
// @Router /covid19/providers [get]
func (h ApisHandler) GetProviders(appVersion *string, w http.ResponseWriter, r *http.Request) {
	providers, err := h.app.Services.GetProviders()
	if err != nil {
		log.Println("Error on getting the providers items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var response []providerResponse
	if providers != nil {
		for _, provider := range providers {
			r := providerResponse{ID: provider.ID, ProviderName: provider.Name, ManualTest: provider.ManualTest, AvailableMechanisms: provider.AvailableMechanisms}
			response = append(response, r)
		}
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the providers items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetProvidersByCounty gets the providers which have locations in a specific county
// @Description Gives the providers which have locations in a specific county
// @Tags Covid19
// @ID getProvidersByCounty
// @Accept json
// @Param id path string true "County ID"
// @Success 200 {array} rest.providerResponse
// @Security RokwireAuth
// @Router /covid19/providers/county/{id} [get]
func (h ApisHandler) GetProvidersByCounty(appVersion *string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	countyID := params["county-id"]
	if len(countyID) <= 0 {
		log.Println("county id is required")
		http.Error(w, "county id is required", http.StatusBadRequest)
		return
	}

	locations, err := h.app.Services.GetLocationsByCountyID(countyID)
	if err != nil {
		log.Println("Error on getting the providers for a county")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	providersList := h.findCountyProviders(countyID, locations)

	data, err := json.Marshal(providersList)
	if err != nil {
		log.Println("Error on marshal the providers by county items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetProvidersByCounties gets the providers for a list of counties
// //TODO Description Gets the providers for a list of counties. The counties ids have to be comma separated.
// Tags Covid19
// ID getProvidersByCounties
// Accept json
// Param county-ids query string true "County IDs"
// TODO Success 200 {map} rest.providerResponse
// Security RokwireAuth
// Router /covid19/providers [get]
func (h ApisHandler) GetProvidersByCounties(appVersion *string, w http.ResponseWriter, r *http.Request) {
	countyIDsKeys, ok := r.URL.Query()["county-ids"]
	if !ok || len(countyIDsKeys[0]) < 1 {
		log.Println("url param 'county-ids' is missing")
		return
	}
	countyIDsKey := countyIDsKeys[0]
	countyIDs := strings.Split(countyIDsKey, ",")
	if len(countyIDs) == 0 {
		http.Error(w, "county-ids is required", http.StatusBadRequest)
		return
	}

	locations, err := h.app.Services.GetLocationsByCounties(countyIDs)
	if err != nil {
		log.Println("Error on getting the providers for a counties list")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := make(map[string][]providerResponse, len(countyIDs))
	for _, countyID := range countyIDs {
		providersList := h.findCountyProviders(countyID, locations)
		response[countyID] = providersList
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the providers by counties items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h ApisHandler) findCountyProviders(countyID string, locations []*model.Location) []providerResponse {
	if locations == nil || len(locations) == 0 {
		return make([]providerResponse, 0)
	}

	var resultList []providerResponse
	for _, loc := range locations {
		if loc.County.ID == countyID {
			provider := loc.Provider
			contains := h.containsProvider(provider.ID, resultList)
			if !contains {
				resEntity := providerResponse{ID: provider.ID, ProviderName: provider.Name, ManualTest: provider.ManualTest, AvailableMechanisms: provider.AvailableMechanisms}
				resultList = append(resultList, resEntity)
			}
		}
	}
	if len(resultList) == 0 {
		resultList = make([]providerResponse, 0)
	}
	return resultList
}

func (h ApisHandler) containsProvider(providerID string, list []providerResponse) bool {
	if list == nil || len(list) == 0 {
		return false
	}
	for _, p := range list {
		if p.ID == providerID {
			return true
		}
	}
	return false
}

type getRulesByCountyResponse struct {
	TestTypeID   string `json:"test_type_id"`
	TestTypeName string `json:"test_type"`
	Priority     *int   `json:"priority"`

	Results []getRulesByCountyResultResponse `json:"results"`
} // @name Rule

type getRulesByCountyResultResponse struct {
	ResultID                   string `json:"result_id"`
	ResultName                 string `json:"result"`
	ResultNextStep             string `json:"result_next_step"`
	ResultNextStepTimeInterval *int   `json:"result_next_step_time_interval"`

	CountyStatusID   string `json:"health_status_id"`
	CountyStatusName string `json:"health_status"`
} // @name TestTypeResultCountyStatus

//GetRulesByCounty gets the rules for a county
// @Description Gets the rules for a county
// @Tags Covid19
// @ID GetRulesByCounty
// @Accept json
// @Param id path string true "County ID"
// @Success 200 {array} getRulesByCountyResponse
// @Security RokwireAuth
// @Router /covid19/rules/county/{id} [get]
func (h ApisHandler) GetRulesByCounty(appVersion *string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	countyID := params["county-id"]
	if len(countyID) <= 0 {
		log.Println("county id is required")
		http.Error(w, "county id is required", http.StatusBadRequest)
		return
	}
	rules, countyStatuses, testTypes, err := h.app.Services.GetRulesByCounty(countyID)
	if err != nil {
		log.Printf("Error on getting the rules items - %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rulesLen := len(rules)
	responseList := make([]getRulesByCountyResponse, rulesLen)
	for i, rule := range rules {
		//test type
		testType := h.findTestType(rule.TestType.ID, testTypes)

		var results []getRulesByCountyResultResponse
		if rule.ResultsStates != nil {
			for _, rs := range rule.ResultsStates {
				//test type result
				testTypeResult := h.findTestTypeResult(rs.TestTypeResultID, testType.Results)
				//county status
				countyStatus := h.findCountyStatus(rs.CountyStatusID, countyStatuses)

				rsItem := getRulesByCountyResultResponse{ResultID: rs.TestTypeResultID, ResultName: testTypeResult.Name,
					ResultNextStep: testTypeResult.NextStep, ResultNextStepTimeInterval: testTypeResult.NextStepOffset,
					CountyStatusID: rs.CountyStatusID, CountyStatusName: countyStatus.Name}
				results = append(results, rsItem)
			}
		}

		ruleResponse := getRulesByCountyResponse{TestTypeID: rule.TestType.ID,
			TestTypeName: testType.Name, Priority: rule.Priority, Results: results}
		responseList[i] = ruleResponse
	}

	data, err := json.Marshal(responseList)
	if err != nil {
		log.Println("Error on marshal the rules items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h ApisHandler) findTestType(ID string, list []*model.TestType) *model.TestType {
	for _, testType := range list {
		if ID == testType.ID {
			return testType
		}
	}
	return nil
}

func (h ApisHandler) findTestTypeResult(ID string, list []model.TestTypeResult) *model.TestTypeResult {
	for _, testTypeResult := range list {
		if ID == testTypeResult.ID {
			return &testTypeResult
		}
	}
	return nil
}

func (h ApisHandler) findCountyStatus(ID string, list []*model.CountyStatus) *model.CountyStatus {
	for _, countyStatus := range list {
		if ID == countyStatus.ID {
			return countyStatus
		}
	}
	return nil
}

func (h ApisHandler) findProvider(ID string, list []*model.Provider) *model.Provider {
	for _, provider := range list {
		if ID == provider.ID {
			return provider
		}
	}
	return nil
}

func (h ApisHandler) findLocation(ID string, list []*model.Location) *model.Location {
	for _, location := range list {
		if ID == location.ID {
			return location
		}
	}
	return nil
}

//GetLocationsByCountyIDProviderID gets the locations for a specific county and provider
// @Description Gets locations for county and provider - pass county-id and provider-id params. Get locations for county - pass county-id param.
// @Tags Covid19
// @ID GetLocations
// @Accept json
// @Param county-id query string false "County ID"
// @Param provider-id query string false "Provider ID"
// @Success 200 {array} locationResponse
// @Security RokwireAuth
// @Router /covid19/locations [get]
func (h ApisHandler) GetLocationsByCountyIDProviderID(appVersion *string, w http.ResponseWriter, r *http.Request) {
	countyKeys, ok := r.URL.Query()["county-id"]
	if !ok || len(countyKeys[0]) < 1 {
		log.Println("url param 'county-id' is missing")
		return
	}
	providerKeys, ok := r.URL.Query()["provider-id"]
	if !ok || len(providerKeys[0]) < 1 {
		log.Println("url param 'provider-id' is missing")
		return
	}
	countyID := countyKeys[0]
	providerID := providerKeys[0]

	locations, err := h.app.Services.GetLocationsByProviderIDCountyID(providerID, countyID)
	if err != nil {
		log.Println("Error on getting the locations for provider and county")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var response []locationResponse
	if locations != nil {
		for _, location := range locations {
			var availableTestsRes []string
			if location.AvailableTests != nil {
				for _, testType := range location.AvailableTests {
					availableTestsRes = append(availableTestsRes, testType.ID)
				}
			}
			locItem := locationResponse{ID: location.ID, Name: location.Name, Address1: location.Address1, Address2: location.Address2,
				City: location.City, State: location.State, ZIP: location.ZIP, Country: location.Country, Latitude: location.Latitude, Longitude: location.Longitude,
				Contact: location.Contact, DaysOfOperation: convertFromDaysOfOperations(location.DaysOfOperation),
				URL: location.URL, Notes: location.Notes, ProviderID: location.Provider.ID, CountyID: location.County.ID, AvailableTests: availableTestsRes}

			response = append(response, locItem)
		}
	}
	if len(response) == 0 {
		response = make([]locationResponse, 0)
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the locations items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetLocationsByCountyID gets the locations for a specific county
func (h ApisHandler) GetLocationsByCountyID(appVersion *string, w http.ResponseWriter, r *http.Request) {
	countyKeys, ok := r.URL.Query()["county-id"]
	if !ok || len(countyKeys[0]) < 1 {
		log.Println("url param 'county-id' is missing")
		return
	}
	countyID := countyKeys[0]

	locations, err := h.app.Services.GetLocationsByCountyID(countyID)
	if err != nil {
		log.Println("Error on getting the locations for county")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var response []locationResponse
	if locations != nil {
		for _, location := range locations {
			var availableTestsRes []string
			if location.AvailableTests != nil {
				for _, testType := range location.AvailableTests {
					availableTestsRes = append(availableTestsRes, testType.ID)
				}
			}
			locItem := locationResponse{ID: location.ID, Name: location.Name, Address1: location.Address1, Address2: location.Address2,
				City: location.City, State: location.State, ZIP: location.ZIP, Country: location.Country, Latitude: location.Latitude, Longitude: location.Longitude,
				Contact: location.Contact, DaysOfOperation: convertFromDaysOfOperations(location.DaysOfOperation),
				URL: location.URL, Notes: location.Notes, ProviderID: location.Provider.ID, CountyID: location.County.ID, AvailableTests: availableTestsRes}

			response = append(response, locItem)
		}
	}
	if len(response) == 0 {
		response = make([]locationResponse, 0)
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the locations items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetLocation gets a location
// @Description Gets a location
// @Tags Covid19
// @ID getLocation
// @Accept json
// @Param id path string true "ID"
// @Success 200 {object} locationResponse
// @Security RokwireAuth
// @Router /covid19/locations/{id} [get]
func (h ApisHandler) GetLocation(appVersion *string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("location id is required")
		http.Error(w, "location id is required", http.StatusBadRequest)
		return
	}
	location, err := h.app.Services.GetLocation(ID)
	if err != nil {
		log.Printf("Error on getting the location- %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if location == nil {
		//not found
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var availableTestsRes []string
	if location.AvailableTests != nil {
		for _, testType := range location.AvailableTests {
			availableTestsRes = append(availableTestsRes, testType.ID)
		}
	}
	locItem := locationResponse{ID: location.ID, Name: location.Name, Address1: location.Address1, Address2: location.Address2,
		City: location.City, State: location.State, ZIP: location.ZIP, Country: location.Country, Latitude: location.Latitude, Longitude: location.Longitude,
		Contact: location.Contact, DaysOfOperation: convertFromDaysOfOperations(location.DaysOfOperation),
		URL: location.URL, Notes: location.Notes, ProviderID: location.Provider.ID, CountyID: location.County.ID, AvailableTests: availableTestsRes}
	data, err := json.Marshal(locItem)
	if err != nil {
		log.Println("Error on marshal a location")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type getMTestTypesResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Priority *int   `json:"priority"`

	Results []getMTestTypesResultResponse `json:"results"`
} // @name TestType

type getMTestTypesResultResponse struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	NextStep            string `json:"next_step"`
	NextStepOffset      *int   `json:"next_step_offset"`
	ResultExpiresOffset *int   `json:"result_expires_offset"`
} // @name TestTypeResult

//GetTestTypesByIDs gets the test types for the provided IDs
func (h ApisHandler) GetTestTypesByIDs(appVersion *string, w http.ResponseWriter, r *http.Request) {
	idsKeys, ok := r.URL.Query()["ids"]
	if !ok || len(idsKeys[0]) < 1 {
		log.Println("url param 'ids' is missing")
		return
	}
	idsKey := idsKeys[0]

	ids := strings.Split(idsKey, ",")
	testTypes, err := h.app.Services.GetTestTypesByIDs(ids)
	if err != nil {
		log.Println("Error on getting the test types for ids")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response []getMTestTypesResponse
	if testTypes != nil {
		for _, testType := range testTypes {
			var results []getMTestTypesResultResponse
			if len(testType.Results) > 0 {
				for _, result := range testType.Results {
					resItem := getMTestTypesResultResponse{ID: result.ID, Name: result.Name, NextStep: result.NextStep,
						NextStepOffset: result.NextStepOffset, ResultExpiresOffset: result.ResultExpiresOffset}
					results = append(results, resItem)
				}
			}
			r := getMTestTypesResponse{ID: testType.ID, Name: testType.Name, Priority: testType.Priority, Results: results}
			response = append(response, r)
		}
	}
	if len(response) == 0 {
		response = make([]getMTestTypesResponse, 0)
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the test types items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetTestTypes gets all test types
// @Description Gets test types. You can filter by id. The ids have to be comma separated.
// @Tags Covid19
// @ID GetTestTypes
// @Accept json
// @Param ids query string false "Test Type IDs"
// @Success 200 {array} getMTestTypesResponse
// @Security RokwireAuth
// @Router /covid19/test-types [get]
func (h ApisHandler) GetTestTypes(appVersion *string, w http.ResponseWriter, r *http.Request) {
	testTypes, err := h.app.Services.GetAllTestTypes()
	if err != nil {
		log.Println("Error on getting the test types")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response []getMTestTypesResponse
	if testTypes != nil {
		for _, testType := range testTypes {
			var results []getMTestTypesResultResponse
			if len(testType.Results) > 0 {
				for _, result := range testType.Results {
					resItem := getMTestTypesResultResponse{ID: result.ID, Name: result.Name, NextStep: result.NextStep,
						NextStepOffset: result.NextStepOffset, ResultExpiresOffset: result.ResultExpiresOffset}
					results = append(results, resItem)
				}
			}
			r := getMTestTypesResponse{ID: testType.ID, Name: testType.Name, Priority: testType.Priority, Results: results}
			response = append(response, r)
		}
	}
	if len(response) == 0 {
		response = make([]getMTestTypesResponse, 0)
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the test types items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type getMSymptomGroupsResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Symptoms []mSymptomResponse `json:"symptoms"`
} // @name SymptomGroup

type mSymptomResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
} // @name Symptom

//GetSymptomGroups gets the symptom groups
// @Deprecated
// @Description Gives the symptom groups
// @Tags Covid19
// @ID getSymptomGroups
// @Accept json
// @Success 200 {array} getMSymptomGroupsResponse
// @Security RokwireAuth
// @Router /covid19/symptom-groups [get]
func (h ApisHandler) GetSymptomGroups(appVersion *string, w http.ResponseWriter, r *http.Request) {
	symptomGroups, err := h.app.Services.GetSymptomGroups()
	if err != nil {
		log.Println("Error on getting the symptom groups items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var response []getMSymptomGroupsResponse
	if symptomGroups != nil {
		for _, sg := range symptomGroups {
			var symptoms []mSymptomResponse
			if sg.Symptoms != nil {
				for _, s := range sg.Symptoms {
					item := mSymptomResponse{ID: s.ID, Name: s.Name}
					symptoms = append(symptoms, item)
				}
			}
			r := getMSymptomGroupsResponse{ID: sg.ID, Name: sg.Name, Symptoms: symptoms}
			response = append(response, r)
		}
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the symptom groups items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetSymptoms gets the symptoms
// @Description Gives the symptoms
// @Tags Covid19
// @ID GetSymptoms
// @Accept json
// @Success 200 {object} string
// @Security RokwireAuth
// @Router /covid19/symptoms [get]
func (h ApisHandler) GetSymptoms(appVersion *string, w http.ResponseWriter, r *http.Request) {
	symptoms, err := h.app.Services.GetSymptoms(appVersion)
	if err != nil {
		log.Printf("Error on getting the symptoms - %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data := []byte(symptoms.Items)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type mSymptomRuleResponse struct {
	ID string `json:"id"`

	Gr1Count int `json:"gr1_count"`
	Gr2Count int `json:"gr2_count"`

	Items []mSymptomRuleItemResponse `json:"items"`
} // @name SymptomRule

type mSymptomRuleItemResponse struct {
	Gr1              bool   `json:"gr1"`
	Gr2              bool   `json:"gr2"`
	CountyStatusID   string `json:"county_status_id"`
	CountyStatusName string `json:"health_status"`
	NextStep         string `json:"next_step"`
} // @name SymptomRuleItem

//GetSymptomRuleByCounty give the symptom rule for a county
// @Deprecated
// @Description Gives the symptom rule for a county.
// @Tags Covid19
// @ID getSymptomRuleByCounty
// @Accept json
// @Param id path string true "County ID"
// @Success 200 {array} mSymptomRuleResponse
// @Security RokwireAuth
// @Router /covid19/symptom-rules/county/{id} [get]
func (h ApisHandler) GetSymptomRuleByCounty(appVersion *string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	countyID := params["county-id"]
	if len(countyID) <= 0 {
		log.Println("county id is required")
		http.Error(w, "county id is required", http.StatusBadRequest)
		return
	}
	symptomRule, countyStatuses, err := h.app.Services.GetSymptomRuleByCounty(countyID)
	if err != nil {
		log.Printf("Error on getting the symptom rule - %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if symptomRule == nil {
		//not found
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var rsResponseItems []mSymptomRuleItemResponse
	if symptomRule.Items != nil {
		for _, item := range symptomRule.Items {
			countyStatus := h.findCountyStatus(item.CountyStatus.ID, countyStatuses)
			r := mSymptomRuleItemResponse{Gr1: item.Gr1, Gr2: item.Gr2, CountyStatusID: item.CountyStatus.ID,
				CountyStatusName: countyStatus.Name, NextStep: item.NextStep}
			rsResponseItems = append(rsResponseItems, r)
		}
	}

	resultItem := mSymptomRuleResponse{ID: symptomRule.ID, Gr1Count: symptomRule.Gr1Count, Gr2Count: symptomRule.Gr2Count, Items: rsResponseItems}
	data, err := json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal a symptom rule")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetAccessRuleByCounty gets the access rule for a county
// @Description Gives the access rule for a county.
// @Tags Covid19
// @ID GetAccessRuleByCounty
// @Accept json
// @Param id path string true "County ID"
// @Success 200 {object} string "TODO"
// @Security RokwireAuth
// @Router /covid19/access-rules/county/{id} [get]
func (h ApisHandler) GetAccessRuleByCounty(appVersion *string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	countyID := params["county-id"]
	if len(countyID) <= 0 {
		log.Println("county id is required")
		http.Error(w, "county id is required", http.StatusBadRequest)
		return
	}

	accessRule, countyStatuses, err := h.app.Services.GetAccessRuleByCounty(countyID)
	if err != nil {
		log.Printf("Error on getting the access rule - %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if accessRule == nil {
		//not found
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	//the client format
	resp := make(map[string]string)
	if accessRule.Rules != nil {
		for _, item := range accessRule.Rules {
			countyStatus := h.findCountyStatus(item.CountyStatusID, countyStatuses)
			resp[countyStatus.Name] = item.Value
		}
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error on marshal an access rule")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type addTraceReportRequest []struct {
	Timestamp   int64  `json:"timestamp" validate:"required"`
	TEK         string `json:"tek" validate:"required"`
	Expirestamp *int64 `json:"expirestamp"`
} // @name addTraceReportRequest

//AddTraceReport adds a trace report
// @Description Adds contact tracing report. "timestamp" - Unix time, the number of milliseconds elapsed since January 1, 1970 UTC
// @Tags Covid19
// @ID AddTraceReport
// @Produce plain
// @Accept json
// @Param data body addTraceReportRequest true "body data"
// @Success 200 {object} string "Successfully added [n] items"
// @Security RokwireAuth
// @Router /covid19/trace/report [post]
func (h ApisHandler) AddTraceReport(appVersion *string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a trace report - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData addTraceReportRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create trace report request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Var(requestData, "required,min=1,dive")
	if err != nil {
		log.Printf("Error on validating create trace report data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//construct the trace exposure items
	var traceExposures []model.TraceExposure
	now := time.Now().UnixNano() / 1000000 //we need milliseconds
	for _, item := range requestData {
		traceExposure := model.TraceExposure{DateAdded: now, Timestamp: item.Timestamp, TEK: item.TEK, Expirestamp: item.Expirestamp}
		traceExposures = append(traceExposures, traceExposure)
	}

	//add it
	insertedCount, err := h.app.Services.AddTraceReport(traceExposures)
	if err != nil {
		log.Printf("Error on adding a trace report - %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Successfully added %d items", insertedCount)))
}

//GetExposures gets the exposures items
// @Description Gives the exposures records. "timestamp" and "date-added" params are optional. It is the time in milliseconds.
// @Tags Covid19
// @ID GetExposures
// @Accept json
// @Param timestamp query int false "timestamp"
// @Param date-added query string false "date-added"
// @Success 200 {array} model.TraceExposure
// @Security RokwireAuth
// @Router /covid19/trace/exposures [get]
func (h ApisHandler) GetExposures(appVersion *string, w http.ResponseWriter, r *http.Request) {
	var timestamp *int64
	timestampKeys, ок := r.URL.Query()["timestamp"]
	if ок && len(timestampKeys[0]) > 0 {
		//there is a param
		ts, err := strconv.ParseInt(timestampKeys[0], 10, 64)
		if err != nil {
			log.Println("bad timestamp value")
			http.Error(w, "bad timestamp value", http.StatusBadRequest)
			return
		}
		timestamp = &ts
	}

	var dateAdded *int64
	dateAddedKeys, ок := r.URL.Query()["date-added"]
	if ок && len(dateAddedKeys[0]) > 0 {
		//there is a param
		da, err := strconv.ParseInt(dateAddedKeys[0], 10, 64)
		if err != nil {
			log.Println("bad date-added value")
			http.Error(w, "bad date-added value", http.StatusBadRequest)
			return
		}
		dateAdded = &da
	}

	items, err := h.app.Services.GetExposures(timestamp, dateAdded)
	if err != nil {
		log.Printf("Error on getting the trace exposures items - %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if items == nil {
		items = make([]model.TraceExposure, 0)
	}

	data, err := json.Marshal(items)
	if err != nil {
		log.Println("Error on marshal the trace exposures items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//NewApisHandler creates new rest Handler instance
func NewApisHandler(app *core.Application) ApisHandler {
	return ApisHandler{app: app}
}
