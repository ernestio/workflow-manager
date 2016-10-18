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

func (p *publisher) Process(s *map[string]interface{}, subject string) (result string, err error) {
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

func (p *publisher) GenericHandler(s *map[string]interface{}, subject string) (string, error) {
	id, _ := (*s)["id"].(string)
	output := GenericComponentMsg{
		Service: id,
		Status:  "processing",
	}

	key := strings.Replace(subject, ".", "_to_", 1)

	m, ok := (*s)[key]
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
func (p *publisher) UpdateTemplateVariables(items []interface{}, s *map[string]interface{}) []interface{} {
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

func (p *publisher) isSupportedMessage(s *map[string]interface{}, subject string) bool {
	w, _ := ParseWorkflow(s)
	valid := w.transitions()
	for _, v := range valid {
		if v == subject {
			return true
		}
	}

	return false
}

func (p *publisher) ServiceCreateError(s *map[string]interface{}) string {
	(*s)["status"] = "errored"
	marshalled, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
		return ""
	}

	id, _ := (*s)["id"].(string)
	natsClient.Request("service.set", []byte(`{"id":"`+id+`","status":"errored"}`), time.Second)

	return string(marshalled)
}

func (p *publisher) ServicesDeleteError(s *map[string]interface{}) string {
	return p.ServiceCreateError(s)
}

// ServiceCreateDone
func (p *publisher) ServiceCreateDone(s *map[string]interface{}) string {
	marshalled, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}

	id, _ := (*s)["id"].(string)
	natsClient.Request("service.set", []byte(`{"id":"`+id+`","status":"done"}`), time.Second)

	return string(marshalled)
}

// ServiceDeleteDone ...
func (p *publisher) ServiceDeleteDone(s *map[string]interface{}) string {
	marshalled, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}

	id, _ := (*s)["id"].(string)
	natsClient.Request("service.set", []byte(`{"id":"`+id+`","status":"done"}`), time.Second)

	return string(marshalled)
}
