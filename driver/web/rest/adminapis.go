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
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

//AdminApisHandler handles the admin rest APIs implementation
type AdminApisHandler struct {
	app *core.Application
}

//GetCovid19Config gets the covid19 config
func (h AdminApisHandler) GetCovid19Config(current model.User, w http.ResponseWriter, r *http.Request) {
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
func (h AdminApisHandler) UpdateCovid19Config(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Router /admin/news [get]
func (h AdminApisHandler) GetNews(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Router /admin/news [post]
func (h AdminApisHandler) CreateNews(current model.User, w http.ResponseWriter, r *http.Request) {
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

	news, err := h.app.Administration.CreateNews(date, title, description, htmlContent, nil)
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
// @Security AppUserAuth
// @Router /admin/news/{id} [put]
func (h AdminApisHandler) UpdateNews(current model.User, w http.ResponseWriter, r *http.Request) {
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

	news, err := h.app.Administration.UpdateNews(ID, requestData.Date, requestData.Title,
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
// @Router /admin/news/{id} [delete]
func (h AdminApisHandler) DeleteNews(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("News item id is required")
		http.Error(w, "News item id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteNews(ID)
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
// @Router /admin/resources [get]
func (h AdminApisHandler) GetResources(current model.User, w http.ResponseWriter, r *http.Request) {
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
	Title        string `json:"title"`
	Link         string `json:"link"`
	DisplayOrder int    `json:"display_order"`
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
// @Router /admin/resources [post]
func (h AdminApisHandler) CreateResources(current model.User, w http.ResponseWriter, r *http.Request) {
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

	resource, err := h.app.Administration.CreateResource(title, link, displayOrder)
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
	Title        string `json:"title"`
	Link         string `json:"link"`
	DisplayOrder int    `json:"display_order"`
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
// @Security AppUserAuth
// @Router /admin/resources/{id} [put]
func (h AdminApisHandler) UpdateResource(current model.User, w http.ResponseWriter, r *http.Request) {
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

	resource, err := h.app.Administration.UpdateResource(ID, requestData.Title, requestData.Link, displayOrder)
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
// @Router /admin/resources/{id} [delete]]
func (h AdminApisHandler) DeleteResource(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Resource item id is required")
		http.Error(w, "Resource item id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteResource(ID)
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
// @Router /admin/resources/display-order [post]
func (h AdminApisHandler) UpdateDisplaOrderResources(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Router /admin/faq [get]
func (h AdminApisHandler) GetFAQs(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Router /admin/faq [post]
func (h AdminApisHandler) CreateFAQItem(current model.User, w http.ResponseWriter, r *http.Request) {
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

	err = h.app.Administration.CreateFAQ(section, sdo, title, description, qdo)
	if err != nil {
		log.Println("Error on creating a faq item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

type updateFAQ struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
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
// @Security AppUserAuth
// @Router /admin/faq/{id} [put]
func (h AdminApisHandler) UpdateFAQItem(current model.User, w http.ResponseWriter, r *http.Request) {
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

	err = h.app.Administration.UpdateFAQ(ID, requestData.Title, requestData.Description, requestData.DisplayOrder)
	if err != nil {
		log.Println("Error on updating the FAQ item")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

type updateFAQSection struct {
	Title        string `json:"title"`
	DisplayOrder int    `json:"display_order"`
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
// @Security AppUserAuth
// @Router /admin/faq/section/{id} [put]
func (h AdminApisHandler) UpdateFAQSection(current model.User, w http.ResponseWriter, r *http.Request) {
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

	err = h.app.Administration.UpdateFAQSection(ID, requestData.Title, requestData.DisplayOrder)
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
// @Router /admin/faq/{id} [delete]
func (h AdminApisHandler) DeleteFAQItem(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("FAQ item id is required")
		http.Error(w, "FAQ item id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteFAQ(ID)
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
// @Router /admin/faq/section/{id} [delete]
func (h AdminApisHandler) DeleteFAQSection(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("FAQ section id is required")
		http.Error(w, "FAQ section id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteFAQSection(ID)
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
// @Router /admin/providers [get]
func (h AdminApisHandler) GetProviders(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Router /admin/providers [post]
func (h AdminApisHandler) CreateProvider(current model.User, w http.ResponseWriter, r *http.Request) {
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
	providerName := requestData.ProviderName
	manualTest := requestData.ManualTest
	mechanisms := requestData.AvailableMechanisms

	provider, err := h.app.Administration.CreateProvider(providerName, *manualTest, mechanisms)
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
// @Security AppUserAuth
// @Router /admin/providers/{id} [put]
func (h AdminApisHandler) UpdateProvider(current model.User, w http.ResponseWriter, r *http.Request) {
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

	provider, err := h.app.Administration.UpdateProvider(ID, requestData.ProviderName, *requestData.ManualTest, requestData.AvailableMechanisms)
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
// @Router /admin/providers/{id} [delete]
func (h AdminApisHandler) DeleteProvider(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Provider id is required")
		http.Error(w, "Provider id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteProvider(ID)
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
	Name          string `json:"name" validate:"required"`
	StateProvince string `json:"state_province" validate:"required"`
	Country       string `json:"country" validate:"required"`
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
// @Router /admin/counties [post]
func (h AdminApisHandler) CreateCounty(current model.User, w http.ResponseWriter, r *http.Request) {
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
	name := requestData.Name
	stateProvince := requestData.StateProvince
	country := requestData.Country

	county, err := h.app.Administration.CreateCounty(current, name, stateProvince, country)
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
	Name          string `json:"name" validate:"required"`
	StateProvince string `json:"state_province" validate:"required"`
	Country       string `json:"country" validate:"required"`
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
// @Security AppUserAuth
// @Router /admin/counties/{id} [put]
func (h AdminApisHandler) UpdateCounty(current model.User, w http.ResponseWriter, r *http.Request) {
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

	county, err := h.app.Administration.UpdateCounty(current, ID, requestData.Name,
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
// @Router /admin/counties/{id} [delete]
func (h AdminApisHandler) DeleteCounty(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("County id is required")
		http.Error(w, "County id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteCounty(current, ID)
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
// @Router /admin/counties [get]
func (h AdminApisHandler) GetCounties(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Router /admin/guidelines [post]
func (h AdminApisHandler) CreateGuideline(current model.User, w http.ResponseWriter, r *http.Request) {
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

	guideline, err := h.app.Administration.CreateGuideline(countyID, name, description, items)
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
// @Security AppUserAuth
// @Router /admin/guidelines/{id} [put]
func (h AdminApisHandler) UpdateGuideline(current model.User, w http.ResponseWriter, r *http.Request) {
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

	guideline, err := h.app.Administration.UpdateGuideline(current, ID, name, description, items)
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
// @Router /admin/guidelines/{id} [delete]
func (h AdminApisHandler) DeleteGuideline(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Guideline id is required")
		http.Error(w, "Guideline id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteGuideline(ID)
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
// @Router /admin/guidelines [get]
func (h AdminApisHandler) GetGuidelinesByCountyID(current model.User, w http.ResponseWriter, r *http.Request) {
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
	CountyID    string `json:"county_id" validate:"uuid"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
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
// @Router /admin/county-statuses [post]
func (h AdminApisHandler) CreateCountyStatus(current model.User, w http.ResponseWriter, r *http.Request) {
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

	countyID := requestData.CountyID
	name := requestData.Name
	description := requestData.Description

	countyStatus, err := h.app.Administration.CreateCountyStatus(countyID, name, description)
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
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
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
// @Security AppUserAuth
// @Router /admin/county-statuses/{id} [put]
func (h AdminApisHandler) UpdateCountyStatus(current model.User, w http.ResponseWriter, r *http.Request) {
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

	name := requestData.Name
	description := requestData.Description

	countyStatus, err := h.app.Administration.UpdateCountyStatus(ID, name, description)
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
// @Router /admin/county-statuses/{id} [delete]
func (h AdminApisHandler) DeleteCountyStatus(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("County status id is required")
		http.Error(w, "County status id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteCountyStatus(ID)
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
// @Router /admin/county-statuses [get]
func (h AdminApisHandler) GetCountyStatusesByCountyID(current model.User, w http.ResponseWriter, r *http.Request) {
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
	Name     string `json:"name" validate:"required"`
	Priority *int   `json:"priority"`
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
// @Router /admin/test-types [post]
func (h AdminApisHandler) CreateTestType(current model.User, w http.ResponseWriter, r *http.Request) {
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
	name := requestData.Name
	priority := requestData.Priority

	testType, err := h.app.Administration.CreateTestType(name, priority)
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
	Name     string `json:"name" validate:"required"`
	Priority *int   `json:"priority"`
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
// @Security AppUserAuth
// @Router /admin/test-types/{id} [put]
func (h AdminApisHandler) UpdateTestType(current model.User, w http.ResponseWriter, r *http.Request) {
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

	testType, err := h.app.Administration.UpdateTestType(ID, requestData.Name, requestData.Priority)
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
// @Router /admin/test-types/{id} [delete]
func (h AdminApisHandler) DeleteTestType(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Test type id is required")
		http.Error(w, "Test type id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteTestType(ID)
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
// @Router /admin/test-types [get]
func (h AdminApisHandler) GetTestTypes(current model.User, w http.ResponseWriter, r *http.Request) {
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
	TestTypeID          string `json:"test_type_id" validate:"uuid"`
	Name                string `json:"name" validate:"required"`
	NextStep            string `json:"next_step" validate:"required"`
	NextStepOffset      *int   `json:"next_step_offset"`
	ResultExpiresOffset *int   `json:"result_expires_offset"`
} //@name createTestTypeResultRequest

type createTestTypeResultResponse struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	NextStep            string `json:"next_step"`
	NextStepOffset      *int   `json:"next_step_offset"`
	ResultExpiresOffset *int   `json:"result_expires_offset"`
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
// @Router /admin/test-type-results [post]
func (h AdminApisHandler) CreateTestTypeResult(current model.User, w http.ResponseWriter, r *http.Request) {
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

	testTypeID := requestData.TestTypeID
	name := requestData.Name
	nextStep := requestData.NextStep
	nextStepOffset := requestData.NextStepOffset
	resultExpiresOffset := requestData.ResultExpiresOffset

	testTypeResult, err := h.app.Administration.CreateTestTypeResult(testTypeID, name, nextStep, nextStepOffset, resultExpiresOffset)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultItem := createTestTypeResultResponse{ID: testTypeResult.ID, Name: testTypeResult.Name, NextStep: testTypeResult.NextStep,
		NextStepOffset: testTypeResult.NextStepOffset, ResultExpiresOffset: testTypeResult.ResultExpiresOffset}
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
	Name                string `json:"name" validate:"required"`
	NextStep            string `json:"next_step" validate:"required"`
	NextStepOffset      *int   `json:"next_step_offset"`
	ResultExpiresOffset *int   `json:"result_expires_offset"`
} // @name updateTestTypeResultRequest

type updateTestTypeResultResponse struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	NextStep            string `json:"next_step"`
	NextStepOffset      *int   `json:"next_step_offset"`
	ResultExpiresOffset *int   `json:"result_expires_offset"`
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
// @Security AppUserAuth
// @Router /admin/test-type-results/{id} [put]
func (h AdminApisHandler) UpdateTestTypeResult(current model.User, w http.ResponseWriter, r *http.Request) {
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

	name := requestData.Name
	nextStep := requestData.NextStep
	nextStepOffset := requestData.NextStepOffset
	resultExpiresOffset := requestData.ResultExpiresOffset

	testTypeResult, err := h.app.Administration.UpdateTestTypeResult(ID, name, nextStep, nextStepOffset, resultExpiresOffset)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultItem := updateTestTypeResultResponse{ID: testTypeResult.ID, Name: testTypeResult.Name,
		NextStep: testTypeResult.NextStep, NextStepOffset: testTypeResult.NextStepOffset, ResultExpiresOffset: testTypeResult.ResultExpiresOffset}
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
// @Router /admin/test-type-results/{id} [delete]
func (h AdminApisHandler) DeleteTestTypeResult(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Test type result id is required")
		http.Error(w, "Test type result id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteTestTypeResult(ID)
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
	ID                  string `json:"id"`
	Name                string `json:"name"`
	NextStep            string `json:"next_step"`
	NextStepOffset      *int   `json:"next_step_offset"`
	ResultExpiresOffset *int   `json:"result_expires_offset"`
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
// @Router /admin/test-type-results [get]
func (h AdminApisHandler) GetTestTypeResultsByTestTypeID(current model.User, w http.ResponseWriter, r *http.Request) {
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
			r := getTestTypeResultsResponse{ID: item.ID, Name: item.Name, NextStep: item.NextStep,
				NextStepOffset: item.NextStepOffset, ResultExpiresOffset: item.ResultExpiresOffset}
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
	CountyID   string `json:"county_id" validate:"required,uuid"`
	TestTypeID string `json:"test_type_id" validate:"required,uuid"`
	Priority   *int   `json:"priority"`

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
// @Router /admin/rules [post]
func (h AdminApisHandler) CreateRule(current model.User, w http.ResponseWriter, r *http.Request) {
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

	rule, err := h.app.Administration.CreateRule(countyID, testTypeID, priority, rsItems)
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
	Priority *int `json:"priority"`

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
// @Security AppUserAuth
// @Router /admin/rules/{id} [put]
func (h AdminApisHandler) UpdateRule(current model.User, w http.ResponseWriter, r *http.Request) {
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

	priority := requestData.Priority
	resultsStatuses := requestData.ResultsStatuses

	var rsItems []model.TestTypeResultCountyStatus
	if resultsStatuses != nil {
		for _, rs := range resultsStatuses {
			r := model.TestTypeResultCountyStatus{TestTypeResultID: rs.TestTypeResultID, CountyStatusID: rs.CountyStatusID}
			rsItems = append(rsItems, r)
		}
	}

	rule, err := h.app.Administration.UpdateRule(ID, priority, rsItems)
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
// @Router /admin/rules/{id} [delete]
func (h AdminApisHandler) DeleteRule(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Rule item id is required")
		http.Error(w, "Rule item id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteRule(ID)
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
// @Router /admin/rules [get]
func (h AdminApisHandler) GetRules(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Router /admin/locations [post]
func (h AdminApisHandler) CreateLocation(current model.User, w http.ResponseWriter, r *http.Request) {
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

	providerID := requestData.ProviderID
	countyID := requestData.CountyID

	availableTests := requestData.AvailableTests

	location, err := h.app.Administration.CreateLocation(providerID, countyID, name, address1, address2, city,
		state, zip, country, latitude, longitude, contact, daysOfOperation, url, notes, availableTests)
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
		Country: location.Country, Contact: location.Contact, DaysOfOperation: convertFromDaysOfOperations(location.DaysOfOperation),
		URL: location.URL, Notes: location.Notes, ProviderID: location.Provider.ID, CountyID: location.County.ID, AvailableTests: availableTestsRes}
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
// @Security AppUserAuth
// @Router /admin/locations/{id} [put]
func (h AdminApisHandler) UpdateLocation(current model.User, w http.ResponseWriter, r *http.Request) {
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

	availableTests := requestData.AvailableTests

	location, err := h.app.Administration.UpdateLocation(ID, name, address1, address2, city,
		state, zip, country, latitude, longitude, contact, daysOfOperation, url, notes, availableTests)
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
		Country: location.Country, Contact: location.Contact, DaysOfOperation: convertFromDaysOfOperations(location.DaysOfOperation),
		URL: location.URL, Notes: location.Notes, ProviderID: location.Provider.ID, CountyID: location.County.ID, AvailableTests: availableTestsRes}
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
// @Router /admin/locations [get]
func (h AdminApisHandler) GetLocations(current model.User, w http.ResponseWriter, r *http.Request) {
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
				Contact: location.Contact, DaysOfOperation: convertFromDaysOfOperations(location.DaysOfOperation),
				URL: location.URL, Notes: location.Notes, ProviderID: location.Provider.ID, CountyID: location.County.ID, AvailableTests: availableTestsRes}
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
// @Router /admin/locations/{id} [delete]
func (h AdminApisHandler) DeleteLocation(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Location id is required")
		http.Error(w, "Location id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteLocation(ID)
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
// @Description Creates a symptom
// @Tags Admin
// @ID CreateSymptom
// @Accept json
// @Produce json
// @Param data body createSymptomRequest true "body data"
// @Success 200 {object} symptomResponse
// @Security AdminUserAuth
// @Router /admin/symptoms [post]
func (h AdminApisHandler) CreateSymptom(current model.User, w http.ResponseWriter, r *http.Request) {
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

	symptom, err := h.app.Administration.CreateSymptom(name, symptomGroup)
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
// @Description Updates a symptom.
// @Tags Admin
// @ID UpdateSymptom
// @Accept json
// @Produce json
// @Param data body updateSymptomRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} symptomResponse
// @Security AppUserAuth
// @Router /admin/symptoms/{id} [put]
func (h AdminApisHandler) UpdateSymptom(current model.User, w http.ResponseWriter, r *http.Request) {
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
	symptom, err := h.app.Administration.UpdateSymptom(ID, name)
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
// @Description Deletes a symptom
// @Tags Admin
// @ID deleteSymptom
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Router /admin/symptoms/{id} [delete]
func (h AdminApisHandler) DeleteSymptom(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Symptom id is required")
		http.Error(w, "Symptom id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteSymptom(ID)
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
// @Description Gives the symptom groups
// @Tags Admin
// @ID getSymptomGroups
// @Accept  json
// @Success 200 {array} getSymptomGroupsResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Router /admin/symptom-groups [get]
func (h AdminApisHandler) GetSymptomGroups(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Description Creates a symptom rule.
// @Tags Admin
// @ID CreateSymptomRule
// @Accept json
// @Produce json
// @Param data body createSymptomRuleRequest true "body data"
// @Success 200 {object} symptomRuleResponse
// @Security AdminUserAuth
// @Router /admin/symptom-rules [post]
func (h AdminApisHandler) CreateSymptomRule(current model.User, w http.ResponseWriter, r *http.Request) {
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

	symptomRule, err := h.app.Administration.CreateSymptomRule(countyID, gr1Count, gr2Count, rsItems)
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
// @Description Updates a symptom rule.
// @Tags Admin
// @ID UpdateSymptomRule
// @Accept json
// @Produce json
// @Param data body updateSymptomRuleRequest true "body data"
// @Param id path string true "ID"
// @Success 200 {object} symptomRuleResponse
// @Security AppUserAuth
// @Router /admin/symptom-rules/{id} [put]
func (h AdminApisHandler) UpdateSymptomRule(current model.User, w http.ResponseWriter, r *http.Request) {
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

	symptomRule, err := h.app.Administration.UpdateSymptomRule(ID, countyID, gr1Count, gr2Count, rsItems)
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
// @Description Gives the symptom rules
// @Tags Admin
// @ID getSymptomRules
// @Accept  json
// @Success 200 {array} symptomRuleResponse
// @Failure 400 {object} string "Authentication error"
// @Failure 404 {object} string "Not Found"
// @Failure 500 {object} string "Internal Server error"
// @Security AdminUserAuth
// @Router /admin/symptom-rules [get]
func (h AdminApisHandler) GetSymptomRules(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Description Deletes a symptom rule.
// @Tags Admin
// @ID deleteSymptomRule
// @Accept plain
// @Param id path string true "ID"
// @Success 200 {object} string "Successfuly deleted"
// @Security AdminUserAuth
// @Router /admin/symptom-rules/{id} [delete]
func (h AdminApisHandler) DeleteSymptomRule(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Symptom rule id is required")
		http.Error(w, "Symptom rule id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteSymptomRule(ID)
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
// @Router /admin/manual-tests [get]
func (h AdminApisHandler) GetManualTestsByCountyID(current model.User, w http.ResponseWriter, r *http.Request) {
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
			userResponse := AppUserResponse{UUID: user.UUID, PublicKey: user.PublicKey,
				Consent: user.Consent, ExposureNotification: user.ExposureNotification, RePost: user.RePost,
				EncryptedKey: user.EncryptedKey, EncryptedBlob: user.EncryptedBlob}

			r := eManualTestResponse{ID: item.ID, HistoryID: item.HistoryID, LocationID: item.LocationID,
				CountyID: item.CountyID, EncryptedKey: item.EncryptedKey, EncryptedBlob: item.EncryptedBlob,
				Status: item.Status, Date: item.Date, User: userResponse}
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
// @Security AppUserAuth
// @Router /admin/manual-tests/{id}/process [put]
func (h AdminApisHandler) ProcessManualTest(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Router /admin/manual-tests/{id}/image [get]
func (h AdminApisHandler) GetManualTestImage(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Router /admin/access-rules [post]
func (h AdminApisHandler) CreateAccessRule(current model.User, w http.ResponseWriter, r *http.Request) {
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

	countyID := requestData.CountyID
	rules := requestData.Rules

	var arRules []model.AccessRuleCountyStatus
	if rules != nil {
		for _, ar := range rules {
			r := model.AccessRuleCountyStatus{CountyStatusID: ar.CountyStatusID, Value: ar.Value}
			arRules = append(arRules, r)
		}
	}

	accessRule, err := h.app.Administration.CreateAccessRule(countyID, arRules)
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
// @Router /admin/access-rules [get]
func (h AdminApisHandler) GetAccessRules(current model.User, w http.ResponseWriter, r *http.Request) {
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
// @Security AppUserAuth
// @Router /admin/access-rules/{id} [put]
func (h AdminApisHandler) UpdateAccessRule(current model.User, w http.ResponseWriter, r *http.Request) {
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

	countyID := requestData.CountyID
	rules := requestData.Rules

	var arRules []model.AccessRuleCountyStatus
	if rules != nil {
		for _, ar := range rules {
			r := model.AccessRuleCountyStatus{CountyStatusID: ar.CountyStatusID, Value: ar.Value}
			arRules = append(arRules, r)
		}
	}

	accessRule, err := h.app.Administration.UpdateAccessRule(ID, countyID, arRules)
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
// @Router /admin/access-rules/{id} [delete]
func (h AdminApisHandler) DeleteAccessRule(current model.User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["id"]
	if len(ID) <= 0 {
		log.Println("Access rule id is required")
		http.Error(w, "Access rule id is required", http.StatusBadRequest)
		return
	}
	err := h.app.Administration.DeleteAccessRule(ID)
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
// @Router /admin/user [get]
func (h AdminApisHandler) GetUserByExternalID(current model.User, w http.ResponseWriter, r *http.Request) {
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
	responseUser := AppUserResponse{ID: user.ID, UUID: user.UUID, PublicKey: user.PublicKey,
		Consent: user.Consent, ExposureNotification: user.ExposureNotification, RePost: user.RePost,
		EncryptedKey: user.EncryptedKey, EncryptedBlob: user.EncryptedBlob}
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

type createActionRequest struct {
	ProviderID    string `json:"provider_id" validate:"required"`
	UserID        string `json:"user_id" validate:"required"`
	EncryptedKey  string `json:"encrypted_key" validate:"required"`
	EncryptedBlob string `json:"encrypted_blob" validate:"required"`
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
// @Router /admin/actions [post]
func (h ApisHandler) CreateAction(current model.User, w http.ResponseWriter, r *http.Request) {
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

	providerID := requestData.ProviderID
	userID := requestData.UserID
	encryptedKey := requestData.EncryptedKey
	encryptedBlob := requestData.EncryptedBlob

	item, err := h.app.Administration.CreateAction(providerID, userID, encryptedKey, encryptedBlob)
	if err != nil {
		log.Printf("Error on creating an action - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err = json.Marshal(item)
	if err != nil {
		log.Println("Error on marshal an action")
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
