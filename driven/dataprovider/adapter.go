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

package dataprovider

import (
	"encoding/xml"
	"health/core"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//Adapter implements the data provider interface
type Adapter struct {
	newsRSSURL   string
	resourcesURL string
}

//LoadNews loads the provider news
func (a *Adapter) LoadNews() ([]core.ProviderNews, error) {
	log.Println("LoadNews() -> start loading data provider news...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", a.newsRSSURL, nil)
	if err != nil {
		log.Printf("Error creating news rss request:%s\n", err.Error())
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error loading news %s\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error getting the body for the get news response %s\n", err.Error())
		return nil, err
	}

	var getNewsResponse rss
	err = xml.Unmarshal(body, &getNewsResponse)
	if err != nil {
		log.Printf("Error unmarshal the news %s\n", err.Error())
		return nil, err
	}

	var result []core.ProviderNews

	items := getNewsResponse.Channel.Item
	if items != nil {
		for _, item := range items {
			//Mon, 23 Mar 2020 16:38:23 +0000
			layout := time.RFC1123Z
			parsedDate, err := time.Parse(layout, item.PubDate)
			if err != nil {
				//do not add the item if the time parsing fails
				log.Printf("time parse error %s", err)
			} else {

				/* we need to workaround .000Z
				It appears when it is marshaled in Golang!!
				this is good
				2020-03-20T10:00:00.001Z

				this is bad
				2020-03-20T10:00:00.000Z
				*/
				modDate := parsedDate.Add(time.Millisecond * 1)

				pubDate := modDate
				title := item.Title
				description := item.Description
				contentEncoded := item.Encoded

				cItem := core.ProviderNews{PubDate: pubDate, Title: title,
					Description: description, ContentEncoded: contentEncoded}
				result = append(result, cItem)
			}
		}
	}

	log.Println("LoadNews() -> end loading data provider news...")
	return result, nil
}

//LoadResources loads the provider resources
func (a *Adapter) LoadResources() ([]core.ProviderResource, error) {
	log.Println("LoadResources() -> start loading data provider resources...")

	doc, err := goquery.NewDocument(a.resourcesURL)
	if err != nil {
		log.Printf("LoadResources() -> error creating reader from url - %s\n", err)
		//there is no what to do so return the input
		return nil, err
	}

	var resources []core.ProviderResource

	doc.Find(".menu-resources-container a").Each(func(_ int, link *goquery.Selection) {
		text := strings.TrimSpace(link.Text())

		href, ok := link.Attr("href")
		if ok && len(href) > 0 {

			item := core.ProviderResource{Title: text, Link: href}
			resources = append(resources, item)
		}
	})

	log.Println("LoadResources() -> end loading data provider resources...")
	return resources, nil
}

//NewDataProviderAdapter creates a new provider adapter instance
func NewDataProviderAdapter(newsRSSURL string, resourcesURL string) *Adapter {
	return &Adapter{newsRSSURL: newsRSSURL, resourcesURL: resourcesURL}
}

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Item []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			PubDate     string `xml:"pubDate"`
			Description string `xml:"description"`
			Encoded     string `xml:"encoded"`
		} `xml:"item"`
	} `xml:"channel"`
}
