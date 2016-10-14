/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// When an event has happened and a a new message is about to be sent,
// this is where the message is prepared to be sent.
//
// In order to create your own "translators" just add your mapping for the
// received message to the internal method on the MethodName function.
//
// Then create your own method so, getting a current service, it will
// produce a specific message body to be sent to the dark side.
type publisher struct {
}

func (p *publisher) Process(s *service, subject string) (result string, err error) {
	if p.isSupportedMessage(s, subject) == false {
		return result, errors.New("Message not supported")
	}
	switch subject {
	case "service.create.error":
		result = p.ServiceCreateError(s)
	case "service.create.done":
		result = p.ServiceCreateDone(s)
	case "service.delete.error":
		result = p.ServicesDeleteError(s)
	case "service.delete.done":
		result = p.ServiceDeleteDone(s)
	default:
		return p.GenericHandler(s, subject)
	}

	return result, nil
}

func (p *publisher) GenericHandler(s *service, subject string) (string, error) {
	output := GenericComponentMsg{
		Service: s.ID,
		Status:  "processing",
	}

	key := strings.Replace(subject, ".", "_to_", 1)

	mapped := s.asMap()
	m, ok := mapped[key]
	if ok == false {
		return "", errors.New("Component " + key + " not present")
	}
	list := m.(map[string]interface{})
	if list["items"] == nil {
		return "", errors.New("Could not handle components")
	}
	items := list["items"].([]interface{})
	items = p.UpdateTemplateVariables(items, s)
	output.Components = items

	processing, ok := list["sequential_processing"].(bool)
	if ok {
		output.SequentialProcessing = processing
	}

	marshalled, err := json.Marshal(output)
	if err != nil {
		log.Println(err)
		return "", errors.New(err.Error())
	}

	return string(marshalled), nil
}

func MapString(data string, value string) string {
	if len(value) > 3 && value[0:2] == "$(" && value[len(value)-1:len(value)] == ")" {
		q := gjson.Get(data, value[2:len(value)-1]).String()
		if q != "" && q != "null" {
			return q
		}
		return value
	}
	return value
}

func MapSlice(data string, values []interface{}) []interface{} {
	for i := 0; i < len(values); i++ {
		switch v := values[i].(type) {
		case string:
			values[i] = MapString(data, v)
		case map[string]interface{}:
			for field, selector := range v {
				vv, ok := selector.(string)
				if ok {
					v[field] = MapString(data, vv)
				}
			}
		}
	}
	return values
}

// UpdateTemplateVariables : replaces any qjson queries in fields with information from the current service build
func (p *publisher) UpdateTemplateVariables(items []interface{}, s *service) []interface{} {
	body, err := json.Marshal(s)
	if err != nil {
		log.Println("Can't marshal current service")
		return items
	}
	data := string(body)

	for _, v := range items {
		item := v.(map[string]interface{})

		for field, selector := range item {
			switch value := selector.(type) {
			case string:
				item[field] = MapString(data, value)
			case []interface{}:
				item[field] = MapSlice(data, value)
			}
		}
	}

	return items
}

func (p *publisher) isSupportedMessage(s *service, subject string) bool {
	valid := s.Workflow.transitions()
	for _, v := range valid {
		if v == subject {
			return true
		}
	}

	return false
}

func (p *publisher) ServiceCreateError(s *service) string {
	s.Status = "errored"
	marshalled, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
		return ""
	}

	natsClient.Request("service.set", []byte(`{"id":"`+s.ID+`","status":"errored"}`), time.Second)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "\nOops! Something went wrong. Please manually fix any errors shown above and re-apply your definition.", Level: "INFO"})
	messages = append(messages, MonitorMessage{Body: "error", Level: "ERROR"})
	UserOutput(s.Channel(), messages)

	return string(marshalled)
}

func (p *publisher) ServicesDeleteError(s *service) string {
	return p.ServiceCreateError(s)
}

// ServiceCreateDone
func (p *publisher) ServiceCreateDone(s *service) string {
	marshalled, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}

	natsClient.Request("service.set", []byte(`{"id":"`+s.ID+`","status":"done"}`), time.Second)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "SUCCESS: rules successfully applied", Level: "SUCCESS"})
	messages = append(messages, MonitorMessage{Body: "error", Level: "ERROR"})
	UserOutput(s.Channel(), messages)

	return string(marshalled)
}

// ServiceDeleteDone ...
func (p *publisher) ServiceDeleteDone(s *service) string {
	marshalled, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}

	natsClient.Request("service.set", []byte(`{"id":"`+s.ID+`","status":"done"}`), time.Second)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "SUCCESS: your environment has been successfully deleted", Level: "SUCCESS"})
	messages = append(messages, MonitorMessage{Body: "success", Level: "SUCCESS"})
	UserOutput(s.Channel(), messages)

	return string(marshalled)
}
