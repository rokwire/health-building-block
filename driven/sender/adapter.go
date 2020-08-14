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

package sender

import (
	"bytes"
	"fmt"
	"health/core/model"
	"log"
	"net/smtp"
)

//Adapter implements the Sender interface
type Adapter struct {
	smptHost string
	smtpPort string
	user     string
	password string

	from string
	to   []string
}

//SendForNews sends emails to the recepients for new news
func (a *Adapter) SendForNews(newsList []*model.News) {
	if len(a.to) <= 0 {
		log.Println("SendForNews() -> there is no recepients")
		return
	}
	log.Printf("SendForNews() -> sending data for %d items\n", len(newsList))

	body := a.constructNewsBody(newsList)
	for _, recipient := range a.to {
		go a.send("New COVID-19 News Item", recipient, body)
	}
}

//SendForResources sends emails to the recepients for new resources
func (a *Adapter) SendForResources(resourcesList []*model.Resource) {
	if len(a.to) <= 0 {
		log.Println("SendForResources() -> there is no recepients")
		return
	}
	log.Printf("SendForResources() -> sending data for %d items\n", len(resourcesList))

	body := a.constructResourcesBody(resourcesList)
	for _, recipient := range a.to {
		go a.send("New COVID-19 Resource Item", recipient, body)
	}
}

func (a *Adapter) constructResourcesBody(resourcesList []*model.Resource) string {
	buf := bytes.Buffer{}

	buf.WriteString("Added " + fmt.Sprintf("%d", len(resourcesList)) + " items\n\n\n")

	for _, resource := range resourcesList {
		buf.WriteString("\n\nTitle\n")
		buf.WriteString(resource.Title)
		buf.WriteString("\n\nLink\n")
		buf.WriteString(resource.Link)
		buf.WriteString("\n\n")
	}

	return buf.String()
}

func (a *Adapter) constructNewsBody(newsList []*model.News) string {
	buf := bytes.Buffer{}

	buf.WriteString("Added " + fmt.Sprintf("%d", len(newsList)) + " items\n\n\n")

	for _, news := range newsList {
		buf.WriteString("\n\nTitle\n")
		buf.WriteString(news.Title)
		buf.WriteString("\n\nDescription\n")
		buf.WriteString(news.Description)
		buf.WriteString("\n\nHTML Content\n")
		buf.WriteString(news.HTMLContent)
		buf.WriteString("\n\n")
	}

	return buf.String()
}

func (a *Adapter) send(subject string, to string, body string) {
	log.Printf("sending to... %s\n", to)

	auth := smtp.PlainAuth(
		"",
		a.user,
		a.password,
		a.smptHost,
	)

	from := a.from
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(
		a.smptHost+":"+a.smtpPort,
		auth,
		from,
		[]string{to},
		[]byte(msg),
	)
	if err != nil {
		log.Printf("error sending email - %s", err)
	}
}

//NewSenderAdapter creates a new sender adapter instance
func NewSenderAdapter(smptHost string, smtpPort string, user string, password string,
	from string, to []string) *Adapter {
	return &Adapter{smptHost: smptHost, smtpPort: smtpPort, user: user, password: password,
		from: from, to: to}
}
