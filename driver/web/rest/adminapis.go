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
	"health/core"
	"health/core/model"
	"health/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

//AdminApisHandler handles the admin rest APIs implementation
type AdminApisHandler struct {
	app *core.Application
}

//GetCovid19Config gets the covid19 config
func (h AdminApisHandler) GetCovid19Config(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	config, err := h.app.Administration.GetCovid19Config()
	if err != nil {
		log.Printf("Error on getting covid19 config - %s\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(config)
	if err != nil {
		log.Println("Error on marshal the covid19 config")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//UpdateCovid19Config updates the covid19 config
func (h AdminApisHandler) UpdateCovid19Config(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var config model.COVID19Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.app.Administration.UpdateCovid19Config(&config)
	if err != nil {
		log.Printf("Error on updating covid19 config - %s\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully updated"))
}

//GetAppVersions gives the supported app versions
// @Description Gives the supported app versions
// @Tags Admin
// @ID GetAppVersions
// @Accept  json
// @Success 200 {array} string
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/app-versions [get]
func (h AdminApisHandler) GetAppVersions(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	appVersions, err := h.app.Administration.GetAppVersions()
	if err != nil {
		log.Println("Error on getting the app versions")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(appVersions)
	if err != nil {
		log.Println("Error on marshal the app versions")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createAppVersionRequest struct {
	Audit   *string `json:"audit"`
	Version string  `json:"version" validate:"required"`
} //@name createAppVersionRequest

//CreateAppVersion creates an app version
// @Description Creates an app version. The supported version format is x.x.x or x.x which is the short for x.x.0
// @Tags Admin
// @ID CreateAppVersion
// @Accept json
// @Produce json
// @Param data body createAppVersionRequest true "body data"
// @Success 200 {object} string
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/app-versions [post]
func (h AdminApisHandler) CreateAppVersion(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create app version - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createAppVersionRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create app version data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create app version data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	version := requestData.Version

	err = h.app.Administration.CreateAppVersion(current, group, audit, version)
	if err != nil {
		log.Printf("Error on creating app version - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully created"))
}

//GetNews gets news
// @Description Gives news.
// @Tags Admin
// @ID GetNews
// @Accept  json
// @Success 200 {array} model.News
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/news [get]
func (h AdminApisHandler) GetNews(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	newsItems, err := h.app.Administration.GetNews()
	if err != nil {
		log.Println("Error on getting the news items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(newsItems)
	if err != nil {
		log.Println("Error on marshal the news items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createNews struct {
	Audit       *string   `json:"audit"`
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	HTMLContent string    `json:"htmlContent"`
} // @name createNewsRequest

//CreateNews creates a news
// @Description Creates news
// @Tags Admin
// @ID CreateNews
// @Accept json
// @Produce json
// @Param data body createNews true "body data"
// @Success 200 {object} model.News
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/news [post]
func (h AdminApisHandler) CreateNews(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a news - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createNews
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create news request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	audit := requestData.Audit
	date := requestData.Date
	title := requestData.Title
	description := requestData.Description
	htmlContent := requestData.HTMLContent
	if len(title) <= 0 {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}
	if len(description) <= 0 {
		http.Error(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}
	if len(htmlContent) <= 0 {
		http.Error(w, "HTML content cannot be empty", http.StatusBadRequest)
		return
	}

	news, err := h.app.Administration.CreateNews(current, group, audit, date, title, description, htmlContent, nil)
	if err != nil {
		log.Println("Error on creating a new")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err = json.Marshal(news)
	if err != nil {
		log.Println("Error on marshal a new")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateNews struct {
	Audit       *string   `json:"audit"`
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	HTMLContent string    `json:"htmlContent"`
} // @name updateNewsRequest

//UpdateNews updates news
// @Description Updates news.
// @Tags Admin
// @ID UpdateNews
// @Accept json
// @Produce json
// @Param data body updateNews true "body data"
// @Param id path string true "ID"
// @Success 200 {object} model.News
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/news/{id} [put]
func (h AdminApisHandler) UpdateNews(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("News item id is required")
		http.Error(w, "News item id is required", http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update news item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateNews
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update news item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	news, err := h.app.Administration.UpdateNews(current, group, audit, ID, requestData.Date, requestData.Title,
		requestData.Description, requestData.HTMLContent, nil)
	if err != nil {
		log.Println("Error on updating the news item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err = json.Marshal(news)
	if err != nil {
		log.Println("Error on marshal the news item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteNews deletes a news
// @Description Deletes news
// @Tags Admin
// @ID deleteNews
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted new items"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/news/{id} [delete]
func (h AdminApisHandler) DeleteNews(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("News item id is required")
		http.Error(w, "News item id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteNews(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted news item"))
}

//GetResources gets the resources
// @Description Gives the resources.
// @Tags Admin
// @ID getResources
// @Accept  json
// @Success 200 {array} model.Resource
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/resources [get]
func (h AdminApisHandler) GetResources(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	resources, err := h.app.Administration.GetResources()
	if err != nil {
		log.Println("Error on getting the resource items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(resources)
	if err != nil {
		log.Println("Error on marshal the resource items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createResource struct {
	Audit        *string `json:"audit"`
	Title        string  `json:"title"`
	Link         string  `json:"link"`
	DisplayOrder int     `json:"display_order"`
} // @name createResourceRequest

//CreateResources creates a new resource
// @Description Creates a resource
// @Tags Admin
// @ID createResources
// @Accept json
// @Produce json
// @Param data body createResource true "body data"
// @Success 200 {object} model.Resource
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/resources [post]
func (h AdminApisHandler) CreateResources(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a resource - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createResource
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create resource request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	audit := requestData.Audit
	title := requestData.Title
	link := requestData.Link
	displayOrder := requestData.DisplayOrder
	if len(title) <= 0 {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}
	if len(link) <= 0 {
		http.Error(w, "Link cannot be empty", http.StatusBadRequest)
		return
	}
	if displayOrder <= 0 {
		http.Error(w, "display order cannot be <= 0", http.StatusBadRequest)
		return
	}

	resource, err := h.app.Administration.CreateResource(current, group, audit, title, link, displayOrder)
	if err != nil {
		log.Println("Error on creating a resource")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err = json.Marshal(resource)
	if err != nil {
		log.Println("Error on marshal a resource")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateResource struct {
	Audit        *string `json:"audit"`
	Title        string  `json:"title"`
	Link         string  `json:"link"`
	DisplayOrder int     `json:"display_order"`
} // @name updateResourceRequest

//UpdateResource updates a resource
// @Description Updates a resource.
// @Tags Admin
// @ID UpdateResource
// @Accept json
// @Produce json
// @Param data body updateResource true "body data"
// @Param id path string true "ID"
// @Success 200 {object} model.Resource
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/resources/{id} [put]
func (h AdminApisHandler) UpdateResource(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Resource item id is required")
		http.Error(w, "Resource item id is required", http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update resource item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateResource
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update resource item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	displayOrder := requestData.DisplayOrder
	if displayOrder <= 0 {
		log.Println("display order cannot be <= 0")
		http.Error(w, "display order cannot be <= 0", http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	resource, err := h.app.Administration.UpdateResource(current, group, audit, ID, requestData.Title, requestData.Link, displayOrder)
	if err != nil {
		log.Println("Error on updating the resource item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err = json.Marshal(resource)
	if err != nil {
		log.Println("Error on marshal the resource item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteResource deletes a resource
// @Description Deletes a resource.
// @Tags Admin
// @ID deleteResource
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted resource item"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/resources/{id} [delete]]
func (h AdminApisHandler) DeleteResource(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Resource item id is required")
		http.Error(w, "Resource item id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteResource(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted resource item"))
}

type updateDisplayOrderResource struct {
	IDs []string `json:"ids"`
} // @name updateDisplayOrderResourceRequest

//UpdateDisplaOrderResources updates the display order for all resources
// @Description Updates the display order for all resources.
// @Tags Admin
// @ID updateDisplaOrderResources
// @Accept json
// @Produce json
// @Param data body updateDisplayOrderResource true "body data"
// @Success 200 {object} string "Successfully updated resource items"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/resources/display-order [post]
func (h AdminApisHandler) UpdateDisplaOrderResources(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal update resources display order - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateDisplayOrderResource
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update resources display order - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	IDs := requestData.IDs
	if len(IDs) <= 0 {
		http.Error(w, "IDs cannot be empty", http.StatusBadRequest)
		return
	}

	err = h.app.Administration.UpdateResourceDisplayOrder(IDs)
	if err != nil {
		log.Println("Error on updating resources display order")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully updated resource items"))
}

type createFAQ struct {
	Audit        *string           `json:"audit"`
	Section      string            `json:"section"`
	DisplayOrder int               `json:"display_order"`
	Question     createFAQQuestion `json:"question"`
} //@name createFAQRequest

type createFAQQuestion struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
} //@name createFAQQuestionRequest

//GetFAQs gives FAQs list
// @Description Gives FAQs list
// @Tags Admin
// @ID getFAQs
// @Accept  json
// @Success 200 {array} model.FAQ
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/faq [get]
func (h AdminApisHandler) GetFAQs(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	faq, err := h.app.Administration.GetFAQs()
	if err != nil {
		log.Println("Error on getting the faqs items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//sort
	faq.Sort()

	data, err := json.Marshal(faq.Sections)
	if err != nil {
		log.Println("Error on marshal the faqs items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//CreateFAQItem creates a faq item
// @Description Creates a faq item
// @Tags Admin
// @ID createFAQItem
// @Accept json
// @Produce json
// @Param data body createFAQ true "body data"
// @Success 200 {string} Successfully created
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/faq [post]
func (h AdminApisHandler) CreateFAQItem(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a faq item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createFAQ
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create faq item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	audit := requestData.Audit
	section := requestData.Section
	sdo := requestData.DisplayOrder
	title := requestData.Question.Title
	description := requestData.Question.Description
	qdo := requestData.Question.DisplayOrder
	if len(section) <= 0 {
		http.Error(w, "Section cannot be empty", http.StatusBadRequest)
		return
	}
	if len(title) <= 0 {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}
	if len(description) <= 0 {
		http.Error(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}

	err = h.app.Administration.CreateFAQ(current, group, audit, section, sdo, title, description, qdo)
	if err != nil {
		log.Println("Error on creating a faq item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

type updateFAQ struct {
	Audit        *string `json:"audit"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	DisplayOrder int     `json:"display_order"`
} // @name updateFAQRequest

//UpdateFAQItem updates a faq item
// @Description Updates a faq item.
// @Tags Admin
// @ID UpdateFAQItem
// @Accept json
// @Produce json
// @Param data body updateFAQ true "body data"
// @Param id path string true "ID"
// @Success 200 {string} Successfully updated
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/faq/{id} [put]
func (h AdminApisHandler) UpdateFAQItem(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("FAQ item id is required")
		http.Error(w, "FAQ item id is required", http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update FAQ item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var requestData updateFAQ
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update FAQ item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	title := requestData.Title
	description := requestData.Description
	if len(title) <= 0 {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}
	if len(description) <= 0 {
		http.Error(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	err = h.app.Administration.UpdateFAQ(current, group, audit, ID, requestData.Title, requestData.Description, requestData.DisplayOrder)
	if err != nil {
		log.Println("Error on updating the FAQ item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

type updateFAQSection struct {
	Audit        *string `json:"audit"`
	Title        string  `json:"title"`
	DisplayOrder int     `json:"display_order"`
} // @name updateFAQSection

//UpdateFAQSection updates a faq section
// @Description Updates a faq section.
// @Tags Admin
// @ID UpdateFAQSection
// @Accept json
// @Produce json
// @Param data body updateFAQSection true "body data"
// @Param id path string true "ID"
// @Success 200 {string} Successfully updated
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/faq/section/{id} [put]
func (h AdminApisHandler) UpdateFAQSection(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("FAQ section id is required")
		http.Error(w, "FAQ section id is required", http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update FAQ section - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var requestData updateFAQSection
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update FAQ section request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	audit := requestData.Audit
	title := requestData.Title
	displayOrder := requestData.DisplayOrder
	if len(title) <= 0 {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}
	if displayOrder < 0 {
		http.Error(w, "Display order cannot be < 0", http.StatusBadRequest)
		return
	}

	err = h.app.Administration.UpdateFAQSection(current, group, audit, ID, requestData.Title, requestData.DisplayOrder)
	if err != nil {
		log.Println("Error on updating the FAQ section")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

//DeleteFAQItem deletes a faq item
// @Description Deletes a faq item
// @Tags Admin
// @ID deleteFAQItem
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted FAQ item"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/faq/{id} [delete]
func (h AdminApisHandler) DeleteFAQItem(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("FAQ item id is required")
		http.Error(w, "FAQ item id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteFAQ(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted FAQ item"))
}

//DeleteFAQSection deletes a faq section
// @Description Deletes a faq section
// @Tags Admin
// @ID deleteFAQSection
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/faq/section/{id} [delete]
func (h AdminApisHandler) DeleteFAQSection(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("FAQ section id is required")
		http.Error(w, "FAQ section id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteFAQSection(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

//GetProviders gets the providers
// @Description Gives the providers list
// @Tags Admin
// @ID getProviders
// @Accept  json
// @Success 200 {array} providerResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/providers [get]
func (h AdminApisHandler) GetProviders(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	providers, err := h.app.Administration.GetProviders()
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

type createProviderRequest struct {
	Audit               *string  `json:"audit"`
	ProviderName        string   `json:"provider_name" validate:"required"`
	ManualTest          *bool    `json:"manual_test" validate:"required"`
	AvailableMechanisms []string `json:"available_mechanisms"`
} // @name createProviderRequest

//CreateProvider creates a provider
// @Description Creates a provider
// @Tags Admin
// @ID createProvider
// @Accept json
// @Produce json
// @Param data body createProviderRequest true "body data"
// @Success 200 {object} providerResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/providers [post]
func (h AdminApisHandler) CreateProvider(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a status - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createProviderRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create provider request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create provider data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validate.Var(requestData.AvailableMechanisms, "required,dive,eq=Epic|eq=McKinley|eq=None")
	if err != nil {
		log.Printf("Error on validating create provider mechanisms data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	audit := requestData.Audit
	providerName := requestData.ProviderName
	manualTest := requestData.ManualTest
	mechanisms := requestData.AvailableMechanisms

	provider, err := h.app.Administration.CreateProvider(current, group, audit, providerName, *manualTest, mechanisms)
	if err != nil {
		log.Println("Error on creating a provider")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := providerResponse{ID: provider.ID, ProviderName: provider.Name, ManualTest: provider.ManualTest, AvailableMechanisms: provider.AvailableMechanisms}
	data, err = json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal a provider")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateProviderRequest struct {
	Audit               *string  `json:"audit"`
	ProviderName        string   `json:"provider_name" validate:"required"`
	ManualTest          *bool    `json:"manual_test" validate:"required"`
	AvailableMechanisms []string `json:"available_mechanisms"`
} // @name updateProviderRequest

//UpdateProvider updates a provider
// @Description Updates a provider.
// @Tags Admin
// @ID UpdateProvider
// @Accept json
// @Produce json
// @Param data body updateProviderRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} providerResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/providers/{id} [put]
func (h AdminApisHandler) UpdateProvider(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Provider id is required")
		http.Error(w, "Provider id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update provider item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateProviderRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update provider item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update provider data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validate.Var(requestData.AvailableMechanisms, "required,dive,eq=Epic|eq=McKinley|eq=None")
	if err != nil {
		log.Printf("Error on validating update provider mechanisms data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	provider, err := h.app.Administration.UpdateProvider(current, group, audit, ID, requestData.ProviderName, *requestData.ManualTest, requestData.AvailableMechanisms)
	if err != nil {
		log.Println("Error on updating the provider item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := providerResponse{ID: provider.ID, ProviderName: provider.Name, ManualTest: provider.ManualTest, AvailableMechanisms: provider.AvailableMechanisms}
	data, err = json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the provider item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteProvider deletes a provider
// @Description Deletes a provider
// @Tags Admin
// @ID deleteProvider
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/providers/{id} [delete]
func (h AdminApisHandler) DeleteProvider(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Provider id is required")
		http.Error(w, "Provider id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteProvider(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

type createCountyRequest struct {
	Audit         *string `json:"audit"`
	Name          string  `json:"name" validate:"required"`
	StateProvince string  `json:"state_province" validate:"required"`
	Country       string  `json:"country" validate:"required"`
} //@name createCountyRequest

type createCountyResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	StateProvince string `json:"state_province"`
	Country       string `json:"country"`
} // @name ACounty

//CreateCounty creates a county
// @Description Creates a county
// @Tags Admin
// @ID CreateCounty
// @Accept json
// @Produce json
// @Param data body createCountyRequest true "body data"
// @Success 200 {object} createCountyResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/counties [post]
func (h AdminApisHandler) CreateCounty(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a county - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createCountyRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create county request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create county data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	audit := requestData.Audit
	name := requestData.Name
	stateProvince := requestData.StateProvince
	country := requestData.Country

	county, err := h.app.Administration.CreateCounty(current, group, audit, name, stateProvince, country)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := createCountyResponse{ID: county.ID, Name: county.Name,
		StateProvince: county.StateProvince, Country: county.Country}
	data, err = json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal a county")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateCountyRequest struct {
	Audit         *string `json:"audit"`
	Name          string  `json:"name" validate:"required"`
	StateProvince string  `json:"state_province" validate:"required"`
	Country       string  `json:"country" validate:"required"`
} // @name updateCountyRequest

type updateCountyResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	StateProvince string `json:"state_province"`
	Country       string `json:"country"`
} // @name ACounty

//UpdateCounty updates a county
// @Description Updates a county
// @Tags Admin
// @ID UpdateCounty
// @Accept json
// @Produce json
// @Param data body updateCountyRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} updateCountyResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/counties/{id} [put]
func (h AdminApisHandler) UpdateCounty(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("County id is required")
		http.Error(w, "County id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update county item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateCountyRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update county item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update county data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	county, err := h.app.Administration.UpdateCounty(current, group, audit, ID, requestData.Name,
		requestData.StateProvince, requestData.Country)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := updateCountyResponse{ID: county.ID, Name: county.Name,
		StateProvince: county.StateProvince, Country: county.Country}
	data, err = json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the county item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteCounty deletes a county
// @Description Deletes a county
// @Tags Admin
// @ID deleteCounty
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/counties/{id} [delete]
func (h AdminApisHandler) DeleteCounty(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("County id is required")
		http.Error(w, "County id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteCounty(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

type getCountiesResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	StateProvince string `json:"state_province"`
	Country       string `json:"country"`
} // @name ACounty

//GetCounties gets the counties
// @Description Gives the counties list
// @Tags Admin
// @ID getCounties
// @Accept  json
// @Success 200 {array} getCountiesResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/counties [get]
func (h AdminApisHandler) GetCounties(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	counties, err := h.app.Administration.FindCounties(nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var response []getCountiesResponse
	if counties != nil {
		for _, county := range counties {
			r := getCountiesResponse{ID: county.ID, Name: county.Name,
				StateProvince: county.StateProvince, Country: county.Country}
			response = append(response, r)
		}
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the counties items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createGuidelineRequest struct {
	Audit       *string                       `json:"audit"`
	CountyID    string                        `json:"county_id" validate:"uuid"`
	Name        string                        `json:"name" validate:"required"`
	Description string                        `json:"description"`
	Items       []createGuidelineItemsRequest `json:"items" validate:"required,dive"`
} //@name createGuidelineRequest

type createGuidelineItemsRequest struct {
	Icon        string `json:"icon" validate:"required"`
	Description string `json:"description" validate:"required"`
	Type        string `json:"type" validate:"required"`
} //@name createGuidelineItemsRequest

type createGuidelineResponse struct {
	ID          string                         `json:"id"`
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	Items       []createGuidelineItemsResponse `json:"items"`
} //@name Guideline

type createGuidelineItemsResponse struct {
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Type        string `json:"type"`
} //@name GuidelineItem

//CreateGuideline creates a guideline
// @Description Creates a guideline
// @Tags Admin
// @ID CreateGuideline
// @Accept json
// @Produce json
// @Param data body createGuidelineRequest true "body data"
// @Success 200 {object} createGuidelineResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/guidelines [post]
func (h AdminApisHandler) CreateGuideline(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a guideline - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createGuidelineRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create guideline request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create guideline data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	countyID := requestData.CountyID
	name := requestData.Name
	description := requestData.Description
	reqItems := requestData.Items
	var items []model.GuidelineItem
	for _, reqItem := range reqItems {
		itemType := model.GuidelineItemType{Value: reqItem.Type}
		r := model.GuidelineItem{Icon: reqItem.Icon,
			Description: reqItem.Description, Type: itemType}
		items = append(items, r)
	}

	guideline, err := h.app.Administration.CreateGuideline(current, group, audit, countyID, name, description, items)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var guidelineItems []createGuidelineItemsResponse
	for _, item := range guideline.Items {
		r := createGuidelineItemsResponse{Icon: item.Icon, Description: item.Description, Type: item.Type.Value}
		guidelineItems = append(guidelineItems, r)
	}
	resultItem := createGuidelineResponse{ID: guideline.ID, Name: guideline.Name,
		Description: guideline.Description, Items: guidelineItems}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal a guideline")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateGuidelineRequest struct {
	Audit       *string                       `json:"audit"`
	Name        string                        `json:"name" validate:"required"`
	Description string                        `json:"description"`
	Items       []updateGuidelineItemsRequest `json:"items" validate:"required,dive"`
} // @name updateGuidelineRequest

type updateGuidelineItemsRequest struct {
	Icon        string `json:"icon" validate:"required"`
	Description string `json:"description" validate:"required"`
	Type        string `json:"type" validate:"required"`
} // @name updateGuidelineItemsRequest

type updateGuidelineResponse struct {
	ID          string                         `json:"id"`
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	Items       []updateGuidelineItemsResponse `json:"items"`
} // @name Guideline

type updateGuidelineItemsResponse struct {
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Type        string `json:"type"`
} // @name GuidelineItem

//UpdateGuideline updates a guideline
// @Description Updates a guideline.
// @Tags Admin
// @ID UpdateGuideline
// @Accept json
// @Produce json
// @Param data body updateGuidelineRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} updateGuidelineResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/guidelines/{id} [put]
func (h AdminApisHandler) UpdateGuideline(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Guideline id is required")
		http.Error(w, "Guideline id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal update a guideline - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateGuidelineRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create guideline request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create guideline data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	name := requestData.Name
	description := requestData.Description
	reqItems := requestData.Items
	var items []model.GuidelineItem
	for _, reqItem := range reqItems {
		itemType := model.GuidelineItemType{Value: reqItem.Type}
		r := model.GuidelineItem{Icon: reqItem.Icon,
			Description: reqItem.Description, Type: itemType}
		items = append(items, r)
	}

	guideline, err := h.app.Administration.UpdateGuideline(current, group, audit, ID, name, description, items)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var guidelineItems []updateGuidelineItemsResponse
	for _, item := range guideline.Items {
		r := updateGuidelineItemsResponse{Icon: item.Icon, Description: item.Description, Type: item.Type.Value}
		guidelineItems = append(guidelineItems, r)
	}
	resultItem := updateGuidelineResponse{ID: guideline.ID, Name: guideline.Name,
		Description: guideline.Description, Items: guidelineItems}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal a guideline")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteGuideline deletes a guideline
// @Description Deletes a guideline.
// @Tags Admin
// @ID deleteGuideline
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/guidelines/{id} [delete]
func (h AdminApisHandler) DeleteGuideline(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Guideline id is required")
		http.Error(w, "Guideline id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteGuideline(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

//GetGuidelinesByCountyID gets all guidelines for a county
// @Description Gets all guidelines for a county
// @Tags Admin
// @ID getGuidelinesByCountyID
// @Accept  json
// @Param county-id query string true "County ID"
// @Success 200 {array} guidelinesResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/guidelines [get]
func (h AdminApisHandler) GetGuidelinesByCountyID(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["county-id"]
	if !ok || len(keys[0]) < 1 {
		log.Println("url param 'county-id' is missing")
		return
	}
	countyID := keys[0]

	guidelines, err := h.app.Administration.GetGuidelinesByCountyID(countyID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resultList []guidelinesResponse
	if guidelines != nil {
		for _, item := range guidelines {
			var items []guidelinesItemsResponse
			if item.Items != nil {
				for _, c := range item.Items {
					items = append(items, guidelinesItemsResponse{Icon: c.Icon, Description: c.Description, Type: c.Type.Value})
				}
			}
			r := guidelinesResponse{ID: item.ID, Name: item.Name, Description: item.Description, Items: items}
			resultList = append(resultList, r)
		}
	}
	data, err := json.Marshal(resultList)
	if err != nil {
		log.Println("Error on marshal the guidelines items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createCountyStatusRequest struct {
	Audit       *string `json:"audit"`
	CountyID    string  `json:"county_id" validate:"uuid"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
} //@name createCountyStatusRequest

type createCountyStatusResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
} //@name CountyStatus

//CreateCountyStatus creates a county status
// @Description Creates a county status.
// @Tags Admin
// @ID CreateCountyStatus
// @Accept json
// @Produce json
// @Param data body createCountyStatusRequest true "body data"
// @Success 200 {object} createCountyStatusResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/county-statuses [post]
func (h AdminApisHandler) CreateCountyStatus(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a county status - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createCountyStatusRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create county status request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create county status data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	countyID := requestData.CountyID
	name := requestData.Name
	description := requestData.Description

	countyStatus, err := h.app.Administration.CreateCountyStatus(current, group, audit, countyID, name, description)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resultItem := createCountyStatusResponse{ID: countyStatus.ID, Name: countyStatus.Name, Description: countyStatus.Description}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal a county status")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateCountyStatusRequest struct {
	Audit       *string `json:"audit"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
} // @name updateCountyStatusRequest

type updateCountyStatusResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
} // @name CountyStatus

//UpdateCountyStatus updates a county status
// @Description Updates a county status.
// @Tags Admin
// @ID UpdateCountyStatus
// @Accept json
// @Produce json
// @Param data body updateCountyStatusRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} updateCountyStatusResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/county-statuses/{id} [put]
func (h AdminApisHandler) UpdateCountyStatus(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("County status id is required")
		http.Error(w, "County status id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal update a county status - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateCountyStatusRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create county status request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update county status data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	name := requestData.Name
	description := requestData.Description

	countyStatus, err := h.app.Administration.UpdateCountyStatus(current, group, audit, ID, name, description)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resultItem := updateCountyStatusResponse{ID: countyStatus.ID,
		Name: countyStatus.Name, Description: countyStatus.Description}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal a county status")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteCountyStatus deletes a county status
// @Description Deletes a county status
// @Tags Admin
// @ID deleteCountyStatus
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/county-statuses/{id} [delete]
func (h AdminApisHandler) DeleteCountyStatus(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("County status id is required")
		http.Error(w, "County status id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteCountyStatus(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

type getCountyStatusesResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
} // @name CountyStatus

//GetCountyStatusesByCountyID gets all county statuses for a county
// @Description Gets all county statuses for a county.
// @Tags Admin
// @ID getCountyStatusesByCountyID
// @Accept json
// @Param county-id query string true "County ID"
// @Success 200 {array} getCountyStatusesResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/county-statuses [get]
func (h AdminApisHandler) GetCountyStatusesByCountyID(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["county-id"]
	if !ok || len(keys[0]) < 1 {
		log.Println("url param 'county-id' is missing")
		return
	}
	countyID := keys[0]

	countyStatuses, err := h.app.Administration.GetCountyStatusByCountyID(countyID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resultList []getCountyStatusesResponse
	if countyStatuses != nil {
		for _, item := range countyStatuses {
			r := getCountyStatusesResponse{ID: item.ID, Name: item.Name, Description: item.Description}
			resultList = append(resultList, r)
		}
	}
	data, err := json.Marshal(resultList)
	if err != nil {
		log.Println("Error on marshal the county statuses items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createTestTypeRequest struct {
	Audit    *string `json:"audit"`
	Name     string  `json:"name" validate:"required"`
	Priority *int    `json:"priority"`
} //@name createTestTypeRequest

type createTestTypeResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Priority *int   `json:"priority"`
} //@name ATestType

//CreateTestType creates a test type
// @Description Creates a test type.
// @Tags Admin
// @ID createTestType
// @Accept json
// @Produce json
// @Param data body createTestTypeRequest true "body data"
// @Success 200 {object} createTestTypeResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/test-types [post]
func (h AdminApisHandler) CreateTestType(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a test type - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createTestTypeRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create test type request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create test type data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	audit := requestData.Audit
	name := requestData.Name
	priority := requestData.Priority

	testType, err := h.app.Administration.CreateTestType(current, group, audit, name, priority)
	if err != nil {
		log.Println("Error on creating a test type")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := createTestTypeResponse{ID: testType.ID, Name: testType.Name, Priority: testType.Priority}
	data, err = json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal a test type")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateTestTypeRequest struct {
	Audit    *string `json:"audit"`
	Name     string  `json:"name" validate:"required"`
	Priority *int    `json:"priority"`
} // @name updateTestTypeRequest

type updateTestTypeResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Priority *int   `json:"priority"`
} // @name ATestType

//UpdateTestType updates a test type
// @Description Updates a test type.
// @Tags Admin
// @ID UpdateTestType
// @Accept json
// @Produce json
// @Param data body updateTestTypeRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} updateTestTypeResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/test-types/{id} [put]
func (h AdminApisHandler) UpdateTestType(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Test type id is required")
		http.Error(w, "Test type id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update test type item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateTestTypeRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update test type item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update test type data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	testType, err := h.app.Administration.UpdateTestType(current, group, audit, ID, requestData.Name, requestData.Priority)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := updateTestTypeResponse{ID: testType.ID, Name: testType.Name, Priority: testType.Priority}
	data, err = json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the test type item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteTestType deletes a test type
// @Description Deletes a test type
// @Tags Admin
// @ID deleteTestType
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/test-types/{id} [delete]
func (h AdminApisHandler) DeleteTestType(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Test type id is required")
		http.Error(w, "Test type id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteTestType(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

type getTestTypesResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Priority *int   `json:"priority"`
} // @name ATestType

//GetTestTypes gets the test types
// @Description Gives the test types
// @Tags Admin
// @ID getTestTypes
// @Accept  json
// @Success 200 {array} getTestTypesResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/test-types [get]
func (h AdminApisHandler) GetTestTypes(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	/*if !current.IsAdmin() {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	} */

	testTypes, err := h.app.Administration.GetTestTypes()
	if err != nil {
		log.Println("Error on getting the test types items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var response []getTestTypesResponse
	if testTypes != nil {
		for _, testType := range testTypes {
			r := getTestTypesResponse{ID: testType.ID, Name: testType.Name, Priority: testType.Priority}
			response = append(response, r)
		}
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

type createTestTypeResultRequest struct {
	Audit      *string `json:"audit"`
	TestTypeID string  `json:"test_type_id" validate:"uuid"`
	Name       string  `json:"name" validate:"required"`
} //@name createTestTypeResultRequest

type createTestTypeResultResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
} //@name TestTypeResult

//CreateTestTypeResult creates a test type result for a specific test type
// @Description Creates a test type result for a specific test type.
// @Tags Admin
// @ID CreateTestTypeResult
// @Accept json
// @Produce json
// @Param data body createTestTypeResultRequest true "body data"
// @Success 200 {object} createTestTypeResultResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/test-type-results [post]
func (h AdminApisHandler) CreateTestTypeResult(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a test type result - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createTestTypeResultRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create test type result request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create test type result data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	testTypeID := requestData.TestTypeID
	name := requestData.Name

	testTypeResult, err := h.app.Administration.CreateTestTypeResult(current, group, audit, testTypeID, name, "", nil, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultItem := createTestTypeResultResponse{ID: testTypeResult.ID, Name: testTypeResult.Name}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal a test type result")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateTestTypeResultRequest struct {
	Audit *string `json:"audit"`
	Name  string  `json:"name" validate:"required"`
} // @name updateTestTypeResultRequest

type updateTestTypeResultResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
} // @name TestTypeResult

//UpdateTestTypeResult updates test type result
// @Description Updates test type result.
// @Tags Admin
// @ID UpdateTestTypeResult
// @Accept json
// @Produce json
// @Param data body updateTestTypeResultRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} updateTestTypeResultResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/test-type-results/{id} [put]
func (h AdminApisHandler) UpdateTestTypeResult(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Test type result id is required")
		http.Error(w, "Test type result id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal update a test type result - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateTestTypeResultRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update test type result request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update test type result data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	name := requestData.Name

	testTypeResult, err := h.app.Administration.UpdateTestTypeResult(current, group, audit, ID, name, "", nil, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultItem := updateTestTypeResultResponse{ID: testTypeResult.ID, Name: testTypeResult.Name}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal a test type result")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteTestTypeResult deletes test type result
// @Description Deletes a test type result
// @Tags Admin
// @ID deleteTestTypeResult
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/test-type-results/{id} [delete]
func (h AdminApisHandler) DeleteTestTypeResult(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Test type result id is required")
		http.Error(w, "Test type result id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteTestTypeResult(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

type getTestTypeResultsResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
} // @name TestTypeResult

//GetTestTypeResultsByTestTypeID gets all test type results for a test type
// @Description Gets all test type results for a test type.
// @Tags Admin
// @ID getTestTypeResultsByTestTypeID
// @Accept  json
// @Param test-type-id query string true "Test Type ID"
// @Success 200 {array} getTestTypeResultsResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/test-type-results [get]
func (h AdminApisHandler) GetTestTypeResultsByTestTypeID(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["test-type-id"]
	if !ok || len(keys[0]) < 1 {
		log.Println("url param 'test-type-id' is missing")
		return
	}
	testTypeID := keys[0]

	testTypeResults, err := h.app.Administration.GetTestTypeResultsByTestTypeID(testTypeID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resultList []getTestTypeResultsResponse
	if testTypeResults != nil {
		for _, item := range testTypeResults {
			r := getTestTypeResultsResponse{ID: item.ID, Name: item.Name}
			resultList = append(resultList, r)
		}
	}
	data, err := json.Marshal(resultList)
	if err != nil {
		log.Println("Error on marshal the test type result items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createRuleRequest struct {
	Audit      *string `json:"audit"`
	CountyID   string  `json:"county_id" validate:"required,uuid"`
	TestTypeID string  `json:"test_type_id" validate:"required,uuid"`
	Priority   *int    `json:"priority"`

	ResultsStatuses []createRuleResultsStatusesTypeRequest `json:"results_statuses" validate:"required,dive"`
} //@name createRuleRequest

type createRuleResultsStatusesTypeRequest struct {
	TestTypeResultID string `json:"test_type_result_id" validate:"required,uuid"`
	CountyStatusID   string `json:"county_status_id" validate:"required,uuid"`
} // @name createRuleResultsStatusesTypeRequest

type createRuleResponse struct {
	ID         string `json:"id"`
	CountyID   string `json:"county_id"`
	TestTypeID string `json:"test_type_id"`
	Priority   *int   `json:"priority"`

	ResultsStatuses []createRuleResultsStatusesTypeResponse `json:"results_statuses"`
} //@name ARule

type createRuleResultsStatusesTypeResponse struct {
	TestTypeResultID string `json:"test_type_result_id"`
	CountyStatusID   string `json:"county_status_id"`
} // @name ATestTypeResultCountyStatus

//CreateRule creates a rule
// @Description Creates a rule
// @Tags Admin
// @ID CreateRule
// @Accept json
// @Produce json
// @Param data body createRuleRequest true "body data"
// @Success 200 {object} createRuleResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/rules [post]
func (h AdminApisHandler) CreateRule(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a rule - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createRuleRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create a rule request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create a rule data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	countyID := requestData.CountyID
	testTypeID := requestData.TestTypeID
	priority := requestData.Priority
	resultsStatuses := requestData.ResultsStatuses

	var rsItems []model.TestTypeResultCountyStatus
	if resultsStatuses != nil {
		for _, rs := range resultsStatuses {
			r := model.TestTypeResultCountyStatus{TestTypeResultID: rs.TestTypeResultID, CountyStatusID: rs.CountyStatusID}
			rsItems = append(rsItems, r)
		}
	}

	rule, err := h.app.Administration.CreateRule(current, group, audit, countyID, testTypeID, priority, rsItems)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rsResponseItems []createRuleResultsStatusesTypeResponse
	if rule.ResultsStates != nil {
		for _, item := range rule.ResultsStates {
			r := createRuleResultsStatusesTypeResponse{TestTypeResultID: item.TestTypeResultID, CountyStatusID: item.CountyStatusID}
			rsResponseItems = append(rsResponseItems, r)
		}
	}

	resultItem := createRuleResponse{ID: rule.ID, CountyID: rule.County.ID,
		TestTypeID: rule.TestType.ID, Priority: rule.Priority, ResultsStatuses: rsResponseItems}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal a rule")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateRuleRequest struct {
	Audit    *string `json:"audit"`
	Priority *int    `json:"priority"`

	ResultsStatuses []createRuleResultsStatusesTypeRequest `json:"results_statuses" validate:"required,dive"`
} //@name updateRuleRequest

type updateRuleResultsStatusesTypeRequest struct {
	TestTypeResultID string `json:"test_type_result_id" validate:"required,uuid"`
	CountyStatusID   string `json:"county_status_id" validate:"required,uuid"`
}

type updateRuleResponse struct {
	ID         string `json:"id"`
	CountyID   string `json:"county_id"`
	TestTypeID string `json:"test_type_id"`
	Priority   *int   `json:"priority"`

	ResultsStatuses []updateRuleResultsStatusesTypeResponse `json:"results_statuses"`
} // @name ARule

type updateRuleResultsStatusesTypeResponse struct {
	TestTypeResultID string `json:"test_type_result_id"`
	CountyStatusID   string `json:"county_status_id"`
} // @name ATestTypeResultCountyStatus

//UpdateRule updates a rule
// @Description Updates a rule.
// @Tags Admin
// @ID UpdateRule
// @Accept json
// @Produce json
// @Param data body updateRuleRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} updateRuleResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/rules/{id} [put]
func (h AdminApisHandler) UpdateRule(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Rule id is required")
		http.Error(w, "Rule id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update rule item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateRuleRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update rule item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update rule data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	priority := requestData.Priority
	resultsStatuses := requestData.ResultsStatuses

	var rsItems []model.TestTypeResultCountyStatus
	if resultsStatuses != nil {
		for _, rs := range resultsStatuses {
			r := model.TestTypeResultCountyStatus{TestTypeResultID: rs.TestTypeResultID, CountyStatusID: rs.CountyStatusID}
			rsItems = append(rsItems, r)
		}
	}

	rule, err := h.app.Administration.UpdateRule(current, group, audit, ID, priority, rsItems)
	if err != nil {
		log.Printf("Error on updating the rule item - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var rsResponseItems []updateRuleResultsStatusesTypeResponse
	if rule.ResultsStates != nil {
		for _, item := range rule.ResultsStates {
			r := updateRuleResultsStatusesTypeResponse{TestTypeResultID: item.TestTypeResultID, CountyStatusID: item.CountyStatusID}
			rsResponseItems = append(rsResponseItems, r)
		}
	}

	resultItem := updateRuleResponse{ID: rule.ID, CountyID: rule.County.ID,
		TestTypeID: rule.TestType.ID, Priority: rule.Priority, ResultsStatuses: rsResponseItems}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal a rule")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteRule deletes a rule
// @Description Deletes a rule
// @Tags Admin
// @ID deleteRule
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/rules/{id} [delete]
func (h AdminApisHandler) DeleteRule(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Rule item id is required")
		http.Error(w, "Rule item id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteRule(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

type getRulesResponse struct {
	ID         string `json:"id"`
	CountyID   string `json:"county_id"`
	TestTypeID string `json:"test_type_id"`
	Priority   *int   `json:"priority"`

	ResultsStatuses []getRulesResultsStatusesTypeResponse `json:"results_statuses"`
} // @name ARule

type getRulesResultsStatusesTypeResponse struct {
	TestTypeResultID string `json:"test_type_result_id"`
	CountyStatusID   string `json:"county_status_id"`
} // @name ATestTypeResultCountyStatus

//GetRules gets the rules
// @Description Gives the rules list
// @Tags Admin
// @ID getRules
// @Accept  json
// @Success 200 {array} getRulesResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/rules [get]
func (h AdminApisHandler) GetRules(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	rules, err := h.app.Administration.GetRules()
	if err != nil {
		log.Println("Error on getting the rules items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var response []getRulesResponse
	if rules != nil {
		for _, rule := range rules {
			var resSt []getRulesResultsStatusesTypeResponse
			if rule.ResultsStates != nil {
				for _, rs := range rule.ResultsStates {
					rsItem := getRulesResultsStatusesTypeResponse{TestTypeResultID: rs.TestTypeResultID, CountyStatusID: rs.CountyStatusID}
					resSt = append(resSt, rsItem)
				}
			}
			r := getRulesResponse{ID: rule.ID, CountyID: rule.County.ID, TestTypeID: rule.TestType.ID,
				Priority: rule.Priority, ResultsStatuses: resSt}
			response = append(response, r)
		}
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the rules items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createLocationRequest struct {
	Audit           *string                       `json:"audit"`
	Name            string                        `json:"name" validate:"required"`
	Address1        string                        `json:"address_1"`
	Address2        string                        `json:"address_2"`
	City            string                        `json:"city"`
	State           string                        `json:"state"`
	ZIP             string                        `json:"zip"`
	Contry          string                        `json:"country"`
	Latitude        float64                       `json:"latitude" validate:"required"`
	Longitude       float64                       `json:"longitude" validate:"required"`
	Contact         string                        `json:"contact"`
	DaysOfOperation []locationOperationDayRequest `json:"days_of_operation"`
	URL             string                        `json:"url"`
	Notes           string                        `json:"notes"`
	WaitTimeColor   *string                       `json:"wait_time_color"`

	ProviderID string `json:"provider_id" validate:"required"`
	CountyID   string `json:"county_id" validate:"required"`

	AvailableTests []string `json:"available_tests" validate:"required"`
} //@name createLocationRequest

type locationOperationDayRequest struct {
	Name      string `json:"name"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
} //@name locationOperationDayRequest

//CreateLocation creates a location
// @Description Creates a location.
// @Tags Admin
// @ID CreateLocation
// @Accept json
// @Produce json
// @Param data body createLocationRequest true "body data"
// @Success 200 {object} locationResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/locations [post]
func (h AdminApisHandler) CreateLocation(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a lcoation - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createLocationRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create location request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create location data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	audit := requestData.Audit
	name := requestData.Name
	address1 := requestData.Address1
	address2 := requestData.Address2
	city := requestData.City
	state := requestData.State
	zip := requestData.ZIP
	country := requestData.Contry
	latitude := requestData.Latitude
	longitude := requestData.Longitude
	contact := requestData.Contact
	daysOfOperation := convertToDaysOfOperations(requestData.DaysOfOperation)
	url := requestData.URL
	notes := requestData.Notes
	waitTimeColor := requestData.WaitTimeColor

	providerID := requestData.ProviderID
	countyID := requestData.CountyID

	availableTests := requestData.AvailableTests

	location, err := h.app.Administration.CreateLocation(current, group, audit, providerID, countyID, name, address1, address2, city,
		state, zip, country, latitude, longitude, contact, daysOfOperation, url, notes, waitTimeColor, availableTests)
	if err != nil {
		log.Printf("Error on creating a location - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var availableTestsRes []string
	if location.AvailableTests != nil {
		for _, testType := range location.AvailableTests {
			availableTestsRes = append(availableTestsRes, testType.ID)
		}
	}
	response := locationResponse{ID: location.ID, Name: location.Name, Address1: location.Address1, Address2: location.Address2,
		City: location.City, State: location.State, ZIP: location.ZIP, Latitude: location.Latitude, Longitude: location.Longitude,
		Timezone: location.Timezone, Country: location.Country, Contact: location.Contact, DaysOfOperation: convertFromDaysOfOperations(location.DaysOfOperation),
		URL: location.URL, Notes: location.Notes, WaitTimeColor: location.WaitTimeColor, ProviderID: location.Provider.ID,
		CountyID: location.County.ID, AvailableTests: availableTestsRes}
	data, err = json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal a location")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateLocationRequest struct {
	Audit           *string                       `json:"audit"`
	Name            string                        `json:"name" validate:"required"`
	Address1        string                        `json:"address_1"`
	Address2        string                        `json:"address_2"`
	City            string                        `json:"city"`
	State           string                        `json:"state"`
	ZIP             string                        `json:"zip"`
	Country         string                        `json:"country"`
	Latitude        float64                       `json:"latitude" validate:"required"`
	Longitude       float64                       `json:"longitude" validate:"required"`
	Contact         string                        `json:"contact"`
	DaysOfOperation []locationOperationDayRequest `json:"days_of_operation"`
	URL             string                        `json:"url"`
	Notes           string                        `json:"notes"`
	WaitTimeColor   *string                       `json:"wait_time_color"`

	AvailableTests []string `json:"available_tests" validate:"required"`
} //@name updateLocationRequest

//UpdateLocation updates a location
// @Description Updates a location.
// @Tags Admin
// @ID UpdateLocation
// @Accept json
// @Produce json
// @Param data body updateLocationRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} locationResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/locations/{id} [put]
func (h AdminApisHandler) UpdateLocation(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Location id is required")
		http.Error(w, "Location id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal update a location - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateLocationRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update location request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update location data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	audit := requestData.Audit
	name := requestData.Name
	address1 := requestData.Address1
	address2 := requestData.Address2
	city := requestData.City
	state := requestData.State
	zip := requestData.ZIP
	country := requestData.Country
	latitude := requestData.Latitude
	longitude := requestData.Longitude
	contact := requestData.Contact
	daysOfOperation := convertToDaysOfOperations(requestData.DaysOfOperation)
	url := requestData.URL
	notes := requestData.Notes
	waitTimeColor := requestData.WaitTimeColor

	availableTests := requestData.AvailableTests

	location, err := h.app.Administration.UpdateLocation(current, group, audit, ID, name, address1, address2, city,
		state, zip, country, latitude, longitude, contact, daysOfOperation, url, notes, waitTimeColor, availableTests)
	if err != nil {
		log.Printf("Error on creating a location - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var availableTestsRes []string
	if location.AvailableTests != nil {
		for _, testType := range location.AvailableTests {
			availableTestsRes = append(availableTestsRes, testType.ID)
		}
	}
	response := locationResponse{ID: location.ID, Name: location.Name, Address1: location.Address1, Address2: location.Address2,
		City: location.City, State: location.State, ZIP: location.ZIP, Latitude: location.Latitude, Longitude: location.Longitude,
		Timezone: location.Timezone, Country: location.Country, Contact: location.Contact, DaysOfOperation: convertFromDaysOfOperations(location.DaysOfOperation),
		URL: location.URL, Notes: location.Notes, WaitTimeColor: location.WaitTimeColor, ProviderID: location.Provider.ID,
		CountyID: location.County.ID, AvailableTests: availableTestsRes}
	data, err = json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal a location")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetLocations gets the locations
// @Description Gives the locations list
// @Tags Admin
// @ID getLocations
// @Accept  json
// @Success 200 {array} locationResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/locations [get]
func (h AdminApisHandler) GetLocations(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	locations, err := h.app.Administration.GetLocations()
	if err != nil {
		log.Println("Error on getting the lcoations items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var responseList []locationResponse
	if locations != nil {
		for _, location := range locations {
			var availableTestsRes []string
			if location.AvailableTests != nil {
				for _, testType := range location.AvailableTests {
					availableTestsRes = append(availableTestsRes, testType.ID)
				}
			}
			loc := locationResponse{ID: location.ID, Name: location.Name, Address1: location.Address1, Address2: location.Address2,
				City: location.City, State: location.State, ZIP: location.ZIP, Country: location.Country, Latitude: location.Latitude, Longitude: location.Longitude,
				Timezone: location.Timezone, Contact: location.Contact, DaysOfOperation: convertFromDaysOfOperations(location.DaysOfOperation),
				URL: location.URL, Notes: location.Notes, WaitTimeColor: location.WaitTimeColor, ProviderID: location.Provider.ID,
				CountyID: location.County.ID, AvailableTests: availableTestsRes}
			responseList = append(responseList, loc)
		}
	}
	data, err := json.Marshal(responseList)
	if err != nil {
		log.Println("Error on marshal the locations items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteLocation deletes a location
// @Description Deletes a location
// @Tags Admin
// @ID deleteLocation
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/locations/{id} [delete]
func (h AdminApisHandler) DeleteLocation(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Location id is required")
		http.Error(w, "Location id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteLocation(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

type createSymptomRequest struct {
	Name         string `json:"name" validate:"required"`
	SymptomGroup string `json:"symptom_group" validate:"required,oneof=gr1 gr2"`
} //@name createSymptomRequest

//CreateSymptom creates a symptom
// @Deprecated
// @Description Creates a symptom
// @Tags Admin
// @ID CreateSymptom
// @Accept json
// @Produce json
// @Param data body createSymptomRequest true "body data"
// @Success 200 {object} symptomResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/symptoms [post]
func (h AdminApisHandler) CreateSymptom(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a symptom - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createSymptomRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create symptom request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create symptom data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := requestData.Name
	symptomGroup := requestData.SymptomGroup

	symptom, err := h.app.Administration.CreateSymptom(current, group, name, symptomGroup)
	if err != nil {
		log.Printf("Error on creating a location - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := symptomResponse{ID: symptom.ID, Name: symptom.Name}
	data, err = json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal a symptom")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateSymptomRequest struct {
	Name string `json:"name" validate:"required"`
} //@name updateSymptomRequest

//UpdateSymptom updates a symptom
// @Deprecated
// @Description Updates a symptom.
// @Tags Admin
// @ID UpdateSymptom
// @Accept json
// @Produce json
// @Param data body updateSymptomRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} symptomResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/symptoms/{id} [put]
func (h AdminApisHandler) UpdateSymptom(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Symptom id is required")
		http.Error(w, "Symptom id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal update a symptom - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateSymptomRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update symptom request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update test type result data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := requestData.Name
	symptom, err := h.app.Administration.UpdateSymptom(current, group, ID, name)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := symptomResponse{ID: symptom.ID, Name: symptom.Name}
	data, err = json.Marshal(result)
	if err != nil {
		log.Println("Error on marshal a symptom")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteSymptom deletes a symptom
// @Deprecated
// @Description Deletes a symptom
// @Tags Admin
// @ID deleteSymptom
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/symptoms/{id} [delete]
func (h AdminApisHandler) DeleteSymptom(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Symptom id is required")
		http.Error(w, "Symptom id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteSymptom(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

type getSymptomGroupsResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Symptoms []symptomResponse `json:"symptoms"`
} // @name SymptomGroup

//GetSymptomGroups gets the symptom groups
// @Deprecated
// @Description Gives the symptom groups
// @Tags Admin
// @ID getSymptomGroups
// @Accept  json
// @Success 200 {array} getSymptomGroupsResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/symptom-groups [get]
func (h AdminApisHandler) GetSymptomGroups(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	symptomGroups, err := h.app.Administration.GetSymptomGroups()
	if err != nil {
		log.Println("Error on getting the symptom groups items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var response []getSymptomGroupsResponse
	if symptomGroups != nil {
		for _, sg := range symptomGroups {
			var symptoms []symptomResponse
			if sg.Symptoms != nil {
				for _, s := range sg.Symptoms {
					item := symptomResponse{ID: s.ID, Name: s.Name}
					symptoms = append(symptoms, item)
				}
			}
			r := getSymptomGroupsResponse{ID: sg.ID, Name: sg.Name, Symptoms: symptoms}
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

type createSymptomRuleRequest struct {
	CountyID string `json:"county_id" validate:"required,uuid"`

	Gr1Count int `json:"gr1_count" validate:"required"`
	Gr2Count int `json:"gr2_count" validate:"required"`

	Items []createSymptomRuleItemRequest `json:"items" validate:"required,max=4,min=4,dive"`
} //@name createSymptomRuleRequest

type createSymptomRuleItemRequest struct {
	Gr1            *bool  `json:"gr1" validate:"required"`
	Gr2            *bool  `json:"gr2" validate:"required"`
	CountyStatusID string `json:"county_status_id" validate:"required"`
	NextStep       string `json:"next_step" validate:"required"`
} // @name createSymptomRuleItemRequest

type symptomRuleResponse struct {
	ID       string `json:"id"`
	CountyID string `json:"county_id"`

	Gr1Count int `json:"gr1_count"`
	Gr2Count int `json:"gr2_count"`

	Items []symptomRuleItemResponse `json:"items"`
} // @name SymptomRule

type symptomRuleItemResponse struct {
	Gr1            bool   `json:"gr1"`
	Gr2            bool   `json:"gr2"`
	CountyStatusID string `json:"county_status_id"`
	NextStep       string `json:"next_step"`
} // @name SymptomRuleItem

//CreateSymptomRule creates a symptom rule
// @Deprecated
// @Description Creates a symptom rule.
// @Tags Admin
// @ID CreateSymptomRule
// @Accept json
// @Produce json
// @Param data body createSymptomRuleRequest true "body data"
// @Success 200 {object} symptomRuleResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/symptom-rules [post]
func (h AdminApisHandler) CreateSymptomRule(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a symptom rule - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createSymptomRuleRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create a symptom rule request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create a symptom rule data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	countyID := requestData.CountyID
	gr1Count := requestData.Gr1Count
	gr2Count := requestData.Gr2Count
	items := requestData.Items

	var rsItems []model.SymptomRuleItem
	if items != nil {
		for _, rs := range items {
			countyStatus := model.CountyStatus{ID: rs.CountyStatusID}
			r := model.SymptomRuleItem{Gr1: *rs.Gr1, Gr2: *rs.Gr2, CountyStatus: countyStatus, NextStep: rs.NextStep}
			rsItems = append(rsItems, r)
		}
	}

	symptomRule, err := h.app.Administration.CreateSymptomRule(current, group, countyID, gr1Count, gr2Count, rsItems)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rsResponseItems []symptomRuleItemResponse
	if symptomRule.Items != nil {
		for _, item := range symptomRule.Items {
			r := symptomRuleItemResponse{Gr1: item.Gr1, Gr2: item.Gr2, CountyStatusID: item.CountyStatus.ID, NextStep: item.NextStep}
			rsResponseItems = append(rsResponseItems, r)
		}
	}

	resultItem := symptomRuleResponse{ID: symptomRule.ID, CountyID: symptomRule.County.ID,
		Gr1Count: symptomRule.Gr1Count, Gr2Count: symptomRule.Gr2Count, Items: rsResponseItems}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal a symptom rule")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateSymptomRuleRequest struct {
	CountyID string `json:"county_id" validate:"required,uuid"`

	Gr1Count int `json:"gr1_count" validate:"required"`
	Gr2Count int `json:"gr2_count" validate:"required"`

	Items []updateSymptomRuleItemRequest `json:"items" validate:"required,max=4,min=4,dive"`
} //@name updateSymptomRuleRequest

type updateSymptomRuleItemRequest struct {
	Gr1            *bool  `json:"gr1" validate:"required"`
	Gr2            *bool  `json:"gr2" validate:"required"`
	CountyStatusID string `json:"county_status_id" validate:"required"`
	NextStep       string `json:"next_step" validate:"required"`
} //@name updateSymptomRuleItemRequest

//UpdateSymptomRule updates a symptom rule
// @Deprecated
// @Description Updates a symptom rule.
// @Tags Admin
// @ID UpdateSymptomRule
// @Accept json
// @Produce json
// @Param data body updateSymptomRuleRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} symptomRuleResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/symptom-rules/{id} [put]
func (h AdminApisHandler) UpdateSymptomRule(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Symptom rule id is required")
		http.Error(w, "Symptom rule id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update symptom rule item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateSymptomRuleRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update symptom rule item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update rule data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	countyID := requestData.CountyID
	gr1Count := requestData.Gr1Count
	gr2Count := requestData.Gr2Count
	items := requestData.Items

	var rsItems []model.SymptomRuleItem
	if items != nil {
		for _, rs := range items {
			countyStatus := model.CountyStatus{ID: rs.CountyStatusID}
			r := model.SymptomRuleItem{Gr1: *rs.Gr1, Gr2: *rs.Gr2, CountyStatus: countyStatus, NextStep: rs.NextStep}
			rsItems = append(rsItems, r)
		}
	}

	symptomRule, err := h.app.Administration.UpdateSymptomRule(current, group, ID, countyID, gr1Count, gr2Count, rsItems)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rsResponseItems []symptomRuleItemResponse
	if symptomRule.Items != nil {
		for _, item := range symptomRule.Items {
			r := symptomRuleItemResponse{Gr1: item.Gr1, Gr2: item.Gr2, CountyStatusID: item.CountyStatus.ID, NextStep: item.NextStep}
			rsResponseItems = append(rsResponseItems, r)
		}
	}

	resultItem := symptomRuleResponse{ID: symptomRule.ID, CountyID: symptomRule.County.ID,
		Gr1Count: symptomRule.Gr1Count, Gr2Count: symptomRule.Gr2Count, Items: rsResponseItems}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal a symptom rule")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetSymptomRules gets the symptom rules
// @Deprecated
// @Description Gives the symptom rules
// @Tags Admin
// @ID getSymptomRules
// @Accept  json
// @Success 200 {array} symptomRuleResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/symptom-rules [get]
func (h AdminApisHandler) GetSymptomRules(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	symptomRules, err := h.app.Administration.GetSymptomRules()
	if err != nil {
		log.Println("Error on getting the rules items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var response []symptomRuleResponse
	if symptomRules != nil {
		for _, symptomRule := range symptomRules {
			var rsResponseItems []symptomRuleItemResponse
			if symptomRule.Items != nil {
				for _, item := range symptomRule.Items {
					r := symptomRuleItemResponse{Gr1: item.Gr1, Gr2: item.Gr2, CountyStatusID: item.CountyStatus.ID, NextStep: item.NextStep}
					rsResponseItems = append(rsResponseItems, r)
				}
			}

			resultItem := symptomRuleResponse{ID: symptomRule.ID, CountyID: symptomRule.County.ID,
				Gr1Count: symptomRule.Gr1Count, Gr2Count: symptomRule.Gr2Count, Items: rsResponseItems}
			response = append(response, resultItem)
		}
	}
	if len(response) == 0 {
		response = make([]symptomRuleResponse, 0)
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the symptom rules items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteSymptomRule deletes a symptom rule
// @Deprecated
// @Description Deletes a symptom rule.
// @Tags Admin
// @ID deleteSymptomRule
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/symptom-rules/{id} [delete]
func (h AdminApisHandler) DeleteSymptomRule(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Symptom rule id is required")
		http.Error(w, "Symptom rule id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteSymptomRule(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

//GetManualTestsByCountyID gets the manual tests for a county
// @Description Gives the manual tests for a county
// @Tags Admin
// @ID getManualTestsByCountyID
// @Accept  json
// @Success 200 {array} eManualTestResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/manual-tests [get]
func (h AdminApisHandler) GetManualTestsByCountyID(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	//county id
	countyIDKeys, ok := r.URL.Query()["county-id"]
	if !ok || len(countyIDKeys[0]) < 1 {
		log.Println("url param 'county-id' is missing")
		return
	}
	countyID := countyIDKeys[0]

	//status
	var status *string
	statusKeys, ok := r.URL.Query()["status"]
	if ok && len(statusKeys[0]) > 0 {
		status = &statusKeys[0]
	}

	manualTests, err := h.app.Administration.GetManualTestByCountyID(countyID, status)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var resultList []eManualTestResponse

	if manualTests != nil {
		for _, item := range manualTests {
			user := item.User

			accounts := make([]AppUserAccountResponse, len(user.Accounts))
			if len(user.Accounts) > 0 {
				for i, c := range user.Accounts {
					accounts[i] = AppUserAccountResponse{ID: c.ID, ExternalID: c.ExternalID, Default: c.Default, Active: c.Active,
						FirstName: c.FirstName, MiddleName: c.MiddleName, LastName: c.LastName, BirthDate: c.BirthDate, Gender: c.Gender, Address1: c.Address1,
						Address2: c.Address2, Address3: c.Address3, City: c.City, State: c.State, ZipCode: c.ZipCode, Phone: c.Phone, Email: c.Email}
				}
			}

			userResponse := AppUserResponse{ID: user.ID, ExternalID: user.ExternalID, UUID: user.UUID, PublicKey: user.PublicKey,
				Consent: user.Consent, ExposureNotification: user.ExposureNotification, RePost: user.RePost,
				EncryptedKey: user.EncryptedKey, EncryptedBlob: user.EncryptedBlob, Accounts: accounts}

			r := eManualTestResponse{ID: item.ID, HistoryID: item.HistoryID, LocationID: item.LocationID,
				CountyID: item.CountyID, EncryptedKey: item.EncryptedKey, EncryptedBlob: item.EncryptedBlob,
				Status: item.Status, Date: item.Date, User: userResponse, AccountID: item.AccountID}
			resultList = append(resultList, r)
		}
	} else {
		resultList = make([]eManualTestResponse, 0)
	}
	data, err := json.Marshal(resultList)
	if err != nil {
		log.Println("Error on marshal the manual tests")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type processManualTestRequest struct {
	Status        string  `json:"status" validate:"required,oneof=unverified verified rejected"`
	EncryptedKey  *string `json:"encrypted_key"`
	EncryptedBlob *string `json:"encrypted_blob"`
} //@name processManualTestRequest

//ProcessManualTest processes manual test
// @Description Processes manual test.
// @Tags Admin
// @ID ProcessManualTest
// @Accept json
// @Produce json
// @Param data body processManualTestRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {string} Successfully processed
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/manual-tests/{id}/process [put]
func (h AdminApisHandler) ProcessManualTest(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Manual test id is required")
		http.Error(w, "Manual test id is required", http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal process manual test - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData processManualTestRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the process manual test request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating manual test data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	status := requestData.Status
	encryptedKey := requestData.EncryptedKey
	encryptedBlob := requestData.EncryptedBlob
	if status == "verified" && (encryptedKey == nil || encryptedBlob == nil) {
		http.Error(w, "encrypted key and encrypted blob are required when the status is verified", http.StatusBadRequest)
		return
	}

	err = h.app.Administration.ProcessManualTest(ID, status, encryptedKey, encryptedBlob)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully processed"))
}

//GetManualTestImage gets the image for a manual test
// @Description Gives the image for a manual test
// @Tags Admin
// @ID getManualTestImage
// @Accept  json
// @Param id path string true "Manual Test ID"
// @Success 200 {object} manualTestImageResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/manual-tests/{id}/image [get]
func (h AdminApisHandler) GetManualTestImage(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("id is required")
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	encryptedImageKey, encryptedImageBlob, err := h.app.Administration.GetManualTestImage(ID)
	if err != nil {
		log.Printf("Error on getting the manual test image - %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if encryptedImageKey == nil || encryptedImageBlob == nil {
		//not found
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	responseItem := manualTestImageResponse{EncryptedImageKey: *encryptedImageKey, EncryptedImageBlob: *encryptedImageBlob}
	data, err := json.Marshal(responseItem)
	if err != nil {
		log.Println("Error on marshal the manual test image")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createAccessRuleRequest struct {
	Audit    *string                       `json:"audit"`
	CountyID string                        `json:"county_id" validate:"required,uuid"`
	Rules    []createAccessRuleItemRequest `json:"rules" validate:"required,dive"`
} // @name createAccessRuleRequest

type createAccessRuleItemRequest struct {
	CountyStatusID string `json:"county_status_id" validate:"required"`
	Value          string `json:"value" validate:"required,oneof=granted denied"`
} // @name createAccessRuleItemRequest

type accessRuleResponse struct {
	ID       string `json:"id"`
	CountyID string `json:"county_id"`

	Rules []accessRuleItemResponse `json:"rules"`
} // @name AccessRule

type accessRuleItemResponse struct {
	CountyStatusID string `json:"county_status_id"`
	Value          string `json:"value"`
} // @name AccessRuleCountyStatus

//CreateAccessRule creates an access rule
// @Description Creates an access rule.
// @Tags Admin
// @ID CreateAccessRule
// @Accept json
// @Produce json
// @Param data body createAccessRuleRequest true "body data"
// @Success 200 {object} accessRuleResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/access-rules [post]
func (h AdminApisHandler) CreateAccessRule(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a access rule - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createAccessRuleRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create an access rule request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create an access rule data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	countyID := requestData.CountyID
	rules := requestData.Rules

	var arRules []model.AccessRuleCountyStatus
	if rules != nil {
		for _, ar := range rules {
			r := model.AccessRuleCountyStatus{CountyStatusID: ar.CountyStatusID, Value: ar.Value}
			arRules = append(arRules, r)
		}
	}

	accessRule, err := h.app.Administration.CreateAccessRule(current, group, audit, countyID, arRules)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rsResponseItems []accessRuleItemResponse
	if accessRule.Rules != nil {
		for _, item := range accessRule.Rules {
			r := accessRuleItemResponse{CountyStatusID: item.CountyStatusID, Value: item.Value}
			rsResponseItems = append(rsResponseItems, r)
		}
	}

	resultItem := accessRuleResponse{ID: accessRule.ID, CountyID: accessRule.County.ID, Rules: rsResponseItems}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal an access rule")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetAccessRules gives the access rules
// @Description Gives the access rules
// @Tags Admin
// @ID GetAccessRules
// @Accept  json
// @Success 200 {array} accessRuleResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/access-rules [get]
func (h AdminApisHandler) GetAccessRules(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	accessRules, err := h.app.Administration.GetAccessRules()
	if err != nil {
		log.Println("Error on getting the access rules items")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var response []accessRuleResponse
	if accessRules != nil {
		for _, accessRule := range accessRules {
			var rsResponseItems []accessRuleItemResponse
			if accessRule.Rules != nil {
				for _, item := range accessRule.Rules {
					r := accessRuleItemResponse{CountyStatusID: item.CountyStatusID, Value: item.Value}
					rsResponseItems = append(rsResponseItems, r)
				}
			}

			resultItem := accessRuleResponse{ID: accessRule.ID, CountyID: accessRule.County.ID, Rules: rsResponseItems}
			response = append(response, resultItem)
		}
	}
	if len(response) == 0 {
		response = make([]accessRuleResponse, 0)
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println("Error on marshal the access rules items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateAccessRuleRequest struct {
	Audit    *string                       `json:"audit"`
	CountyID string                        `json:"county_id" validate:"required,uuid"`
	Rules    []updateAccessRuleItemRequest `json:"rules" validate:"required,dive"`
} //@name updateAccessRuleRequest

type updateAccessRuleItemRequest struct {
	CountyStatusID string `json:"county_status_id" validate:"required"`
	Value          string `json:"value" validate:"required,oneof=granted denied"`
} // @name updateAccessRuleItemRequest

//UpdateAccessRule updates an access rule
// @Description Updates an access rule.
// @Tags Admin
// @ID UpdateAccessRule
// @Accept json
// @Produce json
// @Param data body updateAccessRuleRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} accessRuleResponse
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/access-rules/{id} [put]
func (h AdminApisHandler) UpdateAccessRule(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Access rule id is required")
		http.Error(w, "Access rule id is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update access rule item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateAccessRuleRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update access rule item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update an access rule data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	countyID := requestData.CountyID
	rules := requestData.Rules

	var arRules []model.AccessRuleCountyStatus
	if rules != nil {
		for _, ar := range rules {
			r := model.AccessRuleCountyStatus{CountyStatusID: ar.CountyStatusID, Value: ar.Value}
			arRules = append(arRules, r)
		}
	}

	accessRule, err := h.app.Administration.UpdateAccessRule(current, group, audit, ID, countyID, arRules)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rsResponseItems []accessRuleItemResponse
	if accessRule.Rules != nil {
		for _, item := range accessRule.Rules {
			r := accessRuleItemResponse{CountyStatusID: item.CountyStatusID, Value: item.Value}
			rsResponseItems = append(rsResponseItems, r)
		}
	}

	resultItem := accessRuleResponse{ID: accessRule.ID, CountyID: accessRule.County.ID, Rules: rsResponseItems}
	data, err = json.Marshal(resultItem)
	if err != nil {
		log.Println("Error on marshal an access rule")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteAccessRule deletes an access rule
// @Description Deletes an access rule
// @Tags Admin
// @ID deleteAccessRule
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/access-rules/{id} [delete]
func (h AdminApisHandler) DeleteAccessRule(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Access rule id is required")
		http.Error(w, "Access rule id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteAccessRule(current, group, ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

//GetCRules gets the rules
// @Description Gives the rules
// @Tags Admin
// @ID GetCRules
// @Accept json
// @Param county-id query string false "County ID"
// @Param app-version query string false "App version"
// @Success 200 {object} string
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/crules [get]
func (h AdminApisHandler) GetCRules(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	countyKeys, ok := r.URL.Query()["county-id"]
	if !ok || len(countyKeys[0]) < 1 {
		log.Println("url param 'county-id' is missing")
		return
	}
	appVersionKeys, ok := r.URL.Query()["app-version"]
	if !ok || len(appVersionKeys[0]) < 1 {
		log.Println("url param 'app-version' is missing")
		return
	}
	countyID := countyKeys[0]
	appVersion := appVersionKeys[0]

	cRules, err := h.app.Administration.GetCRules(countyID, appVersion)
	if err != nil {
		log.Printf("Error on getting crules - %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cRules == nil {
		//not found
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data := []byte(cRules.Data)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createOrUpdateCRulesRequest struct {
	Audit      *string `json:"audit"`
	AppVersion string  `json:"app_version" validate:"required"`
	CountyID   string  `json:"county_id" validate:"required"`
	Data       string  `json:"data" validate:"required"`
} //@name createOrUpdateCRulesRequest

//CreateOrUpdateCRules creates rules, updates them if already created
// @Description Creates rules, updates them if already created.
// @Tags Admin
// @ID CreateOrUpdateCRules
// @Accept json
// @Produce json
// @Param data body createOrUpdateCRulesRequest true "body data"
// @Success 200 {object} string
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/crules [put]
func (h AdminApisHandler) CreateOrUpdateCRules(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update crules items  - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createOrUpdateCRulesRequest
	err = json.Unmarshal(bodyData, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update crules items request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update crules data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	appVersion := requestData.AppVersion
	countyID := requestData.CountyID
	data := requestData.Data

	err = h.app.Administration.CreateOrUpdateCRules(current, group, audit, countyID, appVersion, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully processed"))
}

//GetSymptoms gets the symptoms
// @Description Gives the symptoms
// @Tags Admin
// @ID GetASymptoms
// @Accept json
// @Param app-version query string false "App version"
// @Success 200 {object} string
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/symptoms [get]
func (h AdminApisHandler) GetSymptoms(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	appVersionKeys, ok := r.URL.Query()["app-version"]
	if !ok || len(appVersionKeys[0]) < 1 {
		log.Println("url param 'app-version' is missing")
		return
	}
	appVersion := appVersionKeys[0]

	symptoms, err := h.app.Administration.GetSymptoms(appVersion)
	if err != nil {
		log.Printf("Error on getting symptoms - %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if symptoms == nil {
		//not found
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data := []byte(symptoms.Items)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createOrUpdateSymptomsRequest struct {
	Audit      *string `json:"audit"`
	AppVersion string  `json:"app_version" validate:"required"`
	Items      string  `json:"items" validate:"required"`
} //@name createOrUpdateSymptomsRequest

//CreateOrUpdateSymptoms creates symptoms or update them if already created
// @Description Creates symptoms or update them if already created.
// @Tags Admin
// @ID CreateorUpdateSymptoms
// @Accept json
// @Produce json
// @Param data body createOrUpdateSymptomsRequest true "body data"
// @Success 200 {string} Successfully processed
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/symptoms [put]
func (h AdminApisHandler) CreateOrUpdateSymptoms(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update symptoms items  - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createOrUpdateSymptomsRequest
	err = json.Unmarshal(bodyData, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update symptoms items request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update symptoms data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	appVersion := requestData.AppVersion
	items := requestData.Items

	err = h.app.Administration.CreateOrUpdateSymptoms(current, group, audit, appVersion, items)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully processed"))
}

//GetUINOverrides gives uin override items
// @Description Gives uin override items
// @Tags Admin
// @ID GetUINOverrides
// @Accept json
// @Param uin query string false "UIN"
// @Param sort query string false "Sort by uin or category"
// @Success 200 {array} model.UINOverride
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/uin-overrides [get]
func (h AdminApisHandler) GetUINOverrides(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	//uin
	var uin *string
	uinKeys, ok := r.URL.Query()["uin"]
	if ok && len(uinKeys[0]) > 0 {
		uin = &uinKeys[0]
	}

	//sort by
	var sort *string
	sortByKeys, ok := r.URL.Query()["sort"]
	if ok {
		sort = &sortByKeys[0]
	}

	uinOverrides, err := h.app.Administration.GetUINOverrides(uin, sort)
	if err != nil {
		log.Println("Error on getting the uin overrides items")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(uinOverrides) == 0 {
		uinOverrides = make([]*model.UINOverride, 0)
	}
	data, err := json.Marshal(uinOverrides)
	if err != nil {
		log.Println("Error on marshal the uin overrides items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createUINOverrideRequest struct {
	Audit      *string    `json:"audit"`
	UIN        string     `json:"uin" validate:"required"`
	Interval   int        `json:"interval" validate:"required"`
	Category   *string    `json:"category"`
	Expiration *time.Time `json:"expiration"`
} // @name createUINOverrideRequest

//CreateUINOverride creates an uin override
// @Description Creates an uin override.
// @Tags Admin
// @ID CreateUINOverride
// @Accept json
// @Produce json
// @Param data body createUINOverrideRequest true "body data"
// @Success 200 {object} model.UINOverride
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/uin-overrides [post]
func (h AdminApisHandler) CreateUINOverride(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create uin override - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createUINOverrideRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create an uin override request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create an uin override data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	uin := requestData.UIN
	interval := requestData.Interval
	category := requestData.Category
	expiration := requestData.Expiration

	uinOverride, err := h.app.Administration.CreateUINOverride(current, group, audit, uin, interval, category, expiration)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err = json.Marshal(uinOverride)
	if err != nil {
		log.Println("Error on marshal an uin override")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type updateUINOverrideRequest struct {
	Audit      *string    `json:"audit"`
	Interval   int        `json:"interval" validate:"required"`
	Category   *string    `json:"category"`
	Expiration *time.Time `json:"expiration"`
} // @name updateUINOverrideRequest

//UpdateUINOverride updates uin override
// @Description Updates uin override.
// @Tags Admin
// @ID UpdateUINOverride
// @Accept json
// @Produce json
// @Param data body updateUINOverrideRequest true "body data"
// @Param uin path string true "UIN"
// @Success 200 {object} string
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/uin-overrides/uin/{uin} [put]
func (h AdminApisHandler) UpdateUINOverride(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uin := params["uin"]
	if len(uin) <= 0 {
		log.Println("uin is required")
		http.Error(w, "uin is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update uin override item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateUINOverrideRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update uin override item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update an uin override data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	interval := requestData.Interval
	category := requestData.Category
	expiration := requestData.Expiration

	uinOverride, err := h.app.Administration.UpdateUINOverride(current, group, audit, uin, interval, category, expiration)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err = json.Marshal(uinOverride)
	if err != nil {
		log.Println("Error on marshal an uin override")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//DeleteUINOverride deletes an uin override
// @Description Deletes an uin override
// @Tags Admin
// @ID DeleteUINOverride
// @Accept plain
// @Param uin path string true "UIN"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/uin-overrides/uin/{uin} [delete]
func (h AdminApisHandler) DeleteUINOverride(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uin := params["uin"]
	if len(uin) <= 0 {
		log.Println("uin is required")
		http.Error(w, "uin is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteUINOverride(current, group, uin)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

//GetUserByExternalID gets the user by external id
// @Description Gets the user by external id.
// @Tags Admin
// @ID GetUserByExternalID
// @Accept json
// @Param external-id query string true "External ID"
// @Success 200 {object} AppUserResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/user [get]
func (h AdminApisHandler) GetUserByExternalID(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["external-id"]
	if !ok || len(keys[0]) < 1 {
		log.Println("url param 'external-id' is missing")
		return
	}
	externalID := keys[0]

	user, err := h.app.Administration.GetUserByExternalID(externalID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		//not found
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	accounts := make([]AppUserAccountResponse, len(user.Accounts))
	if len(user.Accounts) > 0 {
		for i, c := range user.Accounts {
			accounts[i] = AppUserAccountResponse{ID: c.ID, ExternalID: c.ExternalID, Default: c.Default, Active: c.Active,
				FirstName: c.FirstName, MiddleName: c.MiddleName, LastName: c.LastName, BirthDate: c.BirthDate, Gender: c.Gender, Address1: c.Address1,
				Address2: c.Address2, Address3: c.Address3, City: c.City, State: c.State, ZipCode: c.ZipCode, Phone: c.Phone, Email: c.Email}
		}
	}

	responseUser := AppUserResponse{ID: user.ID, ExternalID: user.ExternalID, UUID: user.UUID, PublicKey: user.PublicKey,
		Consent: user.Consent, ExposureNotification: user.ExposureNotification, RePost: user.RePost,
		EncryptedKey: user.EncryptedKey, EncryptedBlob: user.EncryptedBlob, Accounts: accounts}
	data, err := json.Marshal(responseUser)
	if err != nil {
		log.Println("Error on marshal the test type result items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type createRosterRequest struct {
	Audit      *string `json:"audit"`
	Phone      string  `json:"phone" validate:"required"`
	UIN        string  `json:"uin" validate:"required"`
	FirstName  string  `json:"first_name"`
	MiddleName string  `json:"middle_name"`
	LastName   string  `json:"last_name"`
	BirthDate  string  `json:"birth_date"`
	Gender     string  `json:"gender"`
	Address1   string  `json:"address1"`
	Address2   string  `json:"address2"`
	Address3   string  `json:"address3"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	ZipCode    string  `json:"zip_code"`
	Email      string  `json:"email"`
	BadgeType  string  `json:"badge_type"`
} // @name createRosterRequest

//CreateRoster creates a roster
// @Description Creates a roster
// @Tags Admin
// @ID CreateRoster
// @Accept json
// @Produce json
// @Param data body createRosterRequest true "body data"
// @Success 200 {string} Successfully created
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/rosters [post]
func (h AdminApisHandler) CreateRoster(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a roster - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createRosterRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create roster request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create roster data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	audit := requestData.Audit
	phone := requestData.Phone
	uin := requestData.UIN
	firstName := requestData.FirstName
	middleName := requestData.MiddleName
	lastName := requestData.LastName
	birthDate := requestData.BirthDate
	gender := requestData.Gender
	address1 := requestData.Address1
	address2 := requestData.Address2
	address3 := requestData.Address3
	city := requestData.City
	state := requestData.State
	zipCode := requestData.ZipCode
	email := requestData.Email
	badgeType := requestData.BadgeType

	err = h.app.Administration.CreateRoster(current, group, audit, phone, uin, firstName, middleName, lastName,
		birthDate, gender, address1, address2, address3, city, state, zipCode, email, badgeType)
	if err != nil {
		log.Printf("Error on creating a roster - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully created"))
}

type createRosterItemsRequest struct {
	Audit *string `json:"audit"`
	Items []struct {
		Phone      string `json:"phone" validate:"required"`
		UIN        string `json:"uin" validate:"required"`
		FirstName  string `json:"first_name"`
		MiddleName string `json:"middle_name"`
		LastName   string `json:"last_name"`
		BirthDate  string `json:"birth_date"`
		Gender     string `json:"gender"`
		Address1   string `json:"address1"`
		Address2   string `json:"address2"`
		Address3   string `json:"address3"`
		City       string `json:"city"`
		State      string `json:"state"`
		ZipCode    string `json:"zip_code"`
		Email      string `json:"email"`
		BadgeType  string `json:"badge_type"`
	} `json:"items" validate:"required,min=1"`
} // @name createRosterItemsRequest

//CreateRosterItems creates many roster items
// @Description Creates many roster items
// @Tags Admin
// @ID CreateRosterItems
// @Accept json
// @Produce json
// @Param data body createRosterItemsRequest true "body data"
// @Success 200 {string} Successfully created
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/roster-items [post]
func (h AdminApisHandler) CreateRosterItems(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create roster items - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createRosterItemsRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create roster items request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create roster items data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validate.Var(requestData.Items, "required,dive")
	if err != nil {
		log.Printf("Error on validating create roster items - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	items := requestData.Items

	//prepare the items
	itemsList := make([]map[string]string, len(items))
	for i, current := range items {
		item := map[string]string{"phone": current.Phone, "uin": current.UIN, "first_name": current.FirstName, "middle_name": current.MiddleName,
			"last_name": current.LastName, "birth_date": current.BirthDate, "gender": current.Gender, "address1": current.Address1, "address2": current.Address2,
			"address3": current.Address3, "city": current.City, "state": current.State, "zip_code": current.ZipCode, "email": current.Email, "badge_type": current.BadgeType}
		itemsList[i] = item
	}

	err = h.app.Administration.CreateRosterItems(current, group, audit, itemsList)
	if err != nil {
		log.Printf("Error on creating roster items - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully created"))
}

//GetRosters returns the roster members matching filters, sorted, and paginated
// @Description Gives the roster members matching filters, sorted, and paginated
// @Tags Admin
// @ID GetRosters
// @Accept json
// @Param phone query string false "Phone"
// @Param uin query string false "UIN"
// @Param sortBy query string false "Sort By"
// @Param orderBy query string false "Order By"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/rosters [get]
func (h AdminApisHandler) GetRosters(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	sortBy := "phone"
	sortOrder := 1
	limit := 20
	offset := 0
	var filter utils.Filter
	for key, value := range r.URL.Query() {
		if len(value) < 1 || len(value[0]) < 1 {
			continue
		}
		switch key {
		case "sortBy":
			sortBy = value[0]
		case "sortOrder":
			if value[0] == "1" {
				sortOrder = 1
			} else if value[0] == "-1" {
				sortOrder = -1
			} else {
				log.Println("Invalid 'sortOrder' value - " + value[0])
				http.Error(w, "Invalid 'sortOrder' value - Must be 1 or -1", http.StatusBadRequest)
				return
			}
		case "limit":
			limitValue, err := strconv.Atoi(value[0])
			if err == nil {
				if limitValue < 1 || limitValue > 50 {
					log.Println("Invalid 'limit' value - " + value[0])
					http.Error(w, "Invalid 'limit' value - Must be an integer between 1 and 50", http.StatusBadRequest)
					return
				}
				limit = limitValue
			} else {
				log.Printf("error parsing limit - %s\n", err)
			}
		case "offset":
			offsetValue, err := strconv.Atoi(value[0])
			if err == nil {
				offset = offsetValue
			} else {
				log.Printf("error parsing offset - %s\n", err)
			}
		default:
			filter.Items = append(filter.Items, utils.FilterItem{Field: key, Value: value})
		}
	}

	roster, err := h.app.Administration.GetRosters(&filter, sortBy, sortOrder, limit, offset)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(roster)
	if err != nil {
		log.Println("Error on marshal roster")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

//DeleteRosterByPhone deletes a roster by phone
// @Description Deletes a roster by phone
// @Tags Admin
// @ID DeleteRosterByPhone
// @Accept plain
// @Param phone path string true "Phone"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/rosters/phone/{phone} [delete]
func (h AdminApisHandler) DeleteRosterByPhone(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	phone := params["phone"]
	if len(phone) <= 0 {
		log.Println("phone is required")
		http.Error(w, "phone is required", http.StatusBadRequest)
		return
	}

	err := h.app.Administration.DeleteRosterByPhone(current, group, phone)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

//DeleteRosterByUIN deletes a roster by uin
// @Description Deletes a roster by uin
// @Tags Admin
// @ID DeleteRosterByUIN
// @Accept plain
// @Param uin path string true "UIN"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/rosters/uin/{uin} [delete]
func (h AdminApisHandler) DeleteRosterByUIN(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uin := params["uin"]
	if len(uin) <= 0 {
		log.Println("uin is required")
		http.Error(w, "uin is required", http.StatusBadRequest)
		return
	}

	err := h.app.Administration.DeleteRosterByUIN(current, group, uin)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

//DeleteAllRosters deletes all rosters
// @Description Deletes all rosters
// @Tags Admin
// @ID DeleteAllRosters
// @Accept plain
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/rosters [delete]
func (h AdminApisHandler) DeleteAllRosters(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	err := h.app.Administration.DeleteAllRosters(current, group)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

type updateRosterRequest struct {
	Audit      *string `json:"audit"`
	FirstName  string  `json:"first_name"`
	MiddleName string  `json:"middle_name"`
	LastName   string  `json:"last_name"`
	BirthDate  string  `json:"birth_date"`
	Gender     string  `json:"gender"`
	Address1   string  `json:"address1"`
	Address2   string  `json:"address2"`
	Address3   string  `json:"address3"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	ZipCode    string  `json:"zip_code"`
	Email      string  `json:"email"`
	BadgeType  string  `json:"badge_type"`
} // @name updateRosterRequest

//UpdateRoster updates a roster
// @Description Updates a roster.
// @Tags Admin
// @ID UpdateRoster
// @Accept json
// @Produce json
// @Param data body updateRosterRequest true "body data"
// @Param uin path string true "UIN"
// @Success 200 {string} Successfully updated
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/rosters/uin/{id} [put]
func (h AdminApisHandler) UpdateRoster(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uin := params["uin"]
	if len(uin) <= 0 {
		log.Println("uin is required")
		http.Error(w, "uin is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update roster item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateRosterRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update roster item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update roster data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	firstName := requestData.FirstName
	middleName := requestData.MiddleName
	lastName := requestData.LastName
	birthDate := requestData.BirthDate
	gender := requestData.Gender
	address1 := requestData.Address1
	address2 := requestData.Address2
	address3 := requestData.Address3
	city := requestData.City
	state := requestData.State
	zipCode := requestData.ZipCode
	email := requestData.Email
	badgeType := requestData.BadgeType
	err = h.app.Administration.UpdateRoster(current, group, audit, uin, firstName, middleName, lastName, birthDate, gender,
		address1, address2, address3, city, state, zipCode, email, badgeType)
	if err != nil {
		log.Printf("Error on updating roster - %s\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully updated"))
}

type createRawSubAccountItemsRequest struct {
	Audit *string `json:"audit"`
	Items []struct {
		UIN        string `json:"uin" validate:"required"`
		FirstName  string `json:"first_name"`
		MiddleName string `json:"middle_name"`
		LastName   string `json:"last_name"`
		BirthDate  string `json:"birth_date"`
		Gender     string `json:"gender"`
		Address1   string `json:"address1"`
		Address2   string `json:"address2"`
		Address3   string `json:"address3"`
		City       string `json:"city"`
		State      string `json:"state"`
		ZipCode    string `json:"zip_code"`
		Phone      string `json:"phone"  validate:"required"`
		NetID      string `json:"net_id"`
		Email      string `json:"email"`

		PrimaryAccount string `json:"primary_account" validate:"required"`
	} `json:"items" validate:"required,min=1"`
} // @name createRawSubAccountItemsRequest

//CreateSubAccountItems creates sub account items
// @Description Creates sub account items
// @Tags Admin
// @ID CreateSubAccountItems
// @Accept json
// @Produce json
// @Param data body createRawSubAccountItemsRequest true "body data"
// @Success 200 {string} Successfully created
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/raw-sub-account-items [post]
func (h AdminApisHandler) CreateSubAccountItems(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create raw sub accounts items - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createRawSubAccountItemsRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create raw sub account items request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create raw sub account items data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validate.Var(requestData.Items, "required,dive")
	if err != nil {
		log.Printf("Error on validating create raw sub account items - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	items := requestData.Items

	//prepare the items
	itemsList := make([]model.RawSubAccount, len(items))
	for i, current := range items {
		item := model.RawSubAccount{UIN: current.UIN, FirstName: current.FirstName, MiddleName: current.MiddleName,
			LastName: current.LastName, BirthDate: current.BirthDate, Gender: current.Gender, Address1: current.Address1,
			Address2: current.Address2, Address3: current.Address3, City: current.City, State: current.State, ZipCode: current.ZipCode,
			Phone: current.Phone, NetID: current.NetID, Email: current.Email, PrimaryAccount: current.PrimaryAccount}
		itemsList[i] = item
	}

	err = h.app.Administration.CreateRawSubAccountItems(current, group, audit, itemsList)
	if err != nil {
		log.Printf("Error on creating raw sub account items - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully created"))
}

//GetSubAccounts gives sub accounts
// @Description Gives the sub accounts matching filters, sorted, and paginated
// @Tags Admin
// @ID GetSubAccounts
// @Accept json
// @Param phone query string false "Phone"
// @Param uin query string false "UIN"
// @Param sortBy query string false "Sort By"
// @Param orderBy query string false "Order By"
// @Param limit query string false "Limit"
// @Param offset query string false "Offset"
// @Success 200 {array} model.RawSubAccount
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/raw-sub-accounts [get]
func (h AdminApisHandler) GetSubAccounts(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	//TODO
	sortBy := "uin"
	sortOrder := 1
	limit := 20
	offset := 0
	var filter utils.Filter
	for key, value := range r.URL.Query() {
		if len(value) < 1 || len(value[0]) < 1 {
			continue
		}
		switch key {
		case "sortBy":
			sortBy = value[0]
		case "sortOrder":
			if value[0] == "1" {
				sortOrder = 1
			} else if value[0] == "-1" {
				sortOrder = -1
			} else {
				log.Println("Invalid 'sortOrder' value - " + value[0])
				http.Error(w, "Invalid 'sortOrder' value - Must be 1 or -1", http.StatusBadRequest)
				return
			}
		case "limit":
			limitValue, err := strconv.Atoi(value[0])
			if err == nil {
				if limitValue < 1 || limitValue > 50 {
					log.Println("Invalid 'limit' value - " + value[0])
					http.Error(w, "Invalid 'limit' value - Must be an integer between 1 and 50", http.StatusBadRequest)
					return
				}
				limit = limitValue
			} else {
				log.Printf("error parsing limit - %s\n", err)
			}
		case "offset":
			offsetValue, err := strconv.Atoi(value[0])
			if err == nil {
				offset = offsetValue
			} else {
				log.Printf("error parsing offset - %s\n", err)
			}
		default:
			filter.Items = append(filter.Items, utils.FilterItem{Field: key, Value: value})
		}
	}

	subAccounts, err := h.app.Administration.GetRawSubAccounts(&filter, sortBy, sortOrder, limit, offset)
	if err != nil {
		log.Printf("error getting the raw sub accounts - %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(subAccounts)
	if err != nil {
		log.Println("Error on marshal sub accounts")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

type updateRawSubAccountRequest struct {
	Audit      *string `json:"audit"`
	FirstName  string  `json:"first_name"`
	MiddleName string  `json:"middle_name"`
	LastName   string  `json:"last_name"`
	BirthDate  string  `json:"birth_date"`
	Gender     string  `json:"gender"`
	Address1   string  `json:"address1"`
	Address2   string  `json:"address2"`
	Address3   string  `json:"address3"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	ZipCode    string  `json:"zip_code"`
	NetID      string  `json:"net_id"`
	Email      string  `json:"email"`
} // @name updateRawSubAccountRequest

//UpdateSubAccount updates sub account
// @Description Updates a sub account.
// @Tags Admin
// @ID UpdateSubAccount
// @Accept json
// @Produce json
// @Param data body updateRawSubAccountRequest true "body data"
// @Param uin path string true "UIN"
// @Success 200 {string} Successfully updated
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/raw-sub-accounts/uin/{id} [put]
func (h AdminApisHandler) UpdateSubAccount(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uin := params["uin"]
	if len(uin) <= 0 {
		log.Println("uin is required")
		http.Error(w, "uin is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal the update raw sub account item - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData updateRawSubAccountRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the update raw sub account item request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating update raw sub account data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	firstName := requestData.FirstName
	middleName := requestData.MiddleName
	lastName := requestData.LastName
	birthDate := requestData.BirthDate
	gender := requestData.Gender
	address1 := requestData.Address1
	address2 := requestData.Address2
	address3 := requestData.Address3
	city := requestData.City
	state := requestData.State
	zipCode := requestData.ZipCode
	netID := requestData.NetID
	email := requestData.Email

	err = h.app.Administration.UpdateRawSubAccount(current, group, audit, uin, firstName, middleName, lastName, birthDate, gender,
		address1, address2, address3, city, state, zipCode, netID, email)
	if err != nil {
		log.Printf("Error on updating raw sub account - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully updated"))
}

//DeleteSubAccountByUIN deletes a sub account by uin
// @Description Deletes a sub account by uin
// @Tags Admin
// @ID DeleteSubAccountByUIN
// @Accept plain
// @Param uin path string true "UIN"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/raw-sub-accounts/uin/{uin} [delete]
func (h AdminApisHandler) DeleteSubAccountByUIN(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uin := params["uin"]
	if len(uin) <= 0 {
		log.Println("uin is required")
		http.Error(w, "uin is required", http.StatusBadRequest)
		return
	}

	err := h.app.Administration.DeleteRawSubAccountByUIN(current, group, uin)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

//DeleteAllSubAccounts deletes all sub accounts
// @Description Deletes all sub accounts
// @Tags Admin
// @ID DeleteAllSubAccounts
// @Accept plain
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/raw-sub-accounts [delete]
func (h AdminApisHandler) DeleteAllSubAccounts(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	err := h.app.Administration.DeleteAllRawSubAccounts(current, group)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted"))
}

type createActionRequest struct {
	Audit         *string `json:"audit"`
	ProviderID    string  `json:"provider_id" validate:"required"`
	AccountID     string  `json:"account_id" validate:"required"`
	EncryptedKey  string  `json:"encrypted_key" validate:"required"`
	EncryptedBlob string  `json:"encrypted_blob" validate:"required"`
} // @name createActionRequest

//CreateAction creates an action
// @Description Creates an action
// @Tags Admin
// @ID CreateAction
// @Accept json
// @Produce json
// @Param data body createActionRequest true "body data"
// @Success 200 {object} model.CTest
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/actions [post]
func (h ApisHandler) CreateAction(current model.User, group string, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create an action - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var requestData createActionRequest
	err = json.Unmarshal(data, &requestData)
	if err != nil {
		log.Printf("Error on unmarshal the create action request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		log.Printf("Error on validating create action data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	audit := requestData.Audit
	providerID := requestData.ProviderID
	accountID := requestData.AccountID
	encryptedKey := requestData.EncryptedKey
	encryptedBlob := requestData.EncryptedBlob

	item, err := h.app.Administration.CreateAction(current, group, audit, providerID, accountID, encryptedKey, encryptedBlob)
	if err != nil {
		log.Printf("Error on creating an action - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseItem := cTestResponse{ID: item.ID, ProviderID: item.ProviderID, AccountID: item.UserID, EncryptedKey: item.EncryptedKey,
		EncryptedBlob: item.EncryptedBlob, OrderNumber: item.OrderNumber, Processed: item.Processed, DateCreated: item.DateCreated, DateUpdated: item.DateUpdated}

	data, err = json.Marshal(responseItem)
	if err != nil {
		log.Println("Error on marshal an action")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//GetAudit gets the audilt/log history
// @Description Gives the audilt/log history
// @Tags Admin
// @ID GetAudit
// @Accept json
// @Param user-identifier query string false "User identifier"
// @Param entity query string false "Entity"
// @Param entity-id query string false "Entity ID"
// @Param operation query string false "Operation"
// @Param client-data query string false "Client data"
// @Param created-at query string false "Created At"
// @Param sort query string false "Sort By"
// @Param asc query string false "Ascending"
// @Param limit query string false "Limit"
// @Success 200 {array} core.AuditEntity
// @Security AdminUserAuth
// @Security AdminGroupAuth
// @Router /admin/audit [get]
func (h ApisHandler) GetAudit(current model.User, group string, w http.ResponseWriter, r *http.Request) {

	//user identifier
	var userIdentifier *string
	userIdentifierKeys, ok := r.URL.Query()["user-identifier"]
	if ok {
		userIdentifier = &userIdentifierKeys[0]
	}

	//entity
	var entity *string
	entityKeys, ok := r.URL.Query()["entity"]
	if ok {
		entity = &entityKeys[0]
	}

	//entity id
	var entityID *string
	entityIDKeys, ok := r.URL.Query()["entity-id"]
	if ok {
		entityID = &entityIDKeys[0]
	}

	//operation
	var operation *string
	operationKeys, ok := r.URL.Query()["operation"]
	if ok {
		operation = &operationKeys[0]
	}

	//client data
	var clientData *string
	clientDataKeys, ok := r.URL.Query()["client-data"]
	if ok {
		clientData = &clientDataKeys[0]
	}

	//created at
	var createdAt *time.Time
	createdAtKeys, ok := r.URL.Query()["created-at"]
	if ok {
		createdAtValue := &createdAtKeys[0]
		layout := "2006-01-02T15:04:05.000Z"
		t, err := time.Parse(layout, *createdAtValue)
		if err == nil {
			createdAt = &t
		} else {
			log.Printf("error parsing date - %s\n", err)
		}
	}

	//sort by
	var sort *string
	sortByKeys, ok := r.URL.Query()["sort"]
	if ok {
		sort = &sortByKeys[0]
	}

	//asc
	var asc *bool
	ascKeys, ok := r.URL.Query()["asc"]
	if ok {
		ascValue, err := strconv.ParseBool(ascKeys[0])
		if err == nil {
			asc = &ascValue
		} else {
			log.Printf("error parsing asc - %s\n", err)
		}
	}

	//limit
	var limit *int64
	limitKeys, ok := r.URL.Query()["limit"]
	if ok {
		limitValue, err := strconv.ParseInt(limitKeys[0], 10, 64)
		if err == nil {
			limit = &limitValue
		} else {
			log.Printf("error parsing limit - %s\n", err)
		}
	}

	items, err := h.app.Administration.GetAudit(current, group, userIdentifier, entity, entityID, operation, clientData, createdAt, sort, asc, limit)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if len(items) == 0 {
		items = make([]*core.AuditEntity, 0)
	}
	data, err := json.Marshal(items)
	if err != nil {
		log.Println("Error on marshal the audit items")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//NewAdminApisHandler creates new admin rest Handler instance
func NewAdminApisHandler(app *core.Application) AdminApisHandler {
	return AdminApisHandler{app: app}
}
