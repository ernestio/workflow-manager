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

// Publisher : When an event has happened and a a new message is about to be sent,
// this is where the message is prepared to be sent.
//
// In order to create your own "translators" just add your mapping for the
// received message to the internal method on the MethodName function.
//
// Then create your own method so, getting a current service, it will
// produce a specific message body to be sent to the dark side.
type Publisher struct {
}

// Process : starts message publication process
func (p *Publisher) Process(s *map[string]interface{}, subject string) (result string, err error) {
	if p.isSupportedMessage(s, subject) == false {
		return result, errors.New("Message not supported")
	}
	switch subject {
	case "service.create.error", "service.import.error":
		result = p.FinishProcessing(s, "errored")
	case "service.create.done", "service.import.done", "service.import.aws.done":
		result = p.FinishProcessing(s, "done")
	case "service.delete.error":
		result = p.FinishProcessing(s, "errored")
	case "service.delete.done":
		result = p.FinishProcessing(s, "done")
	default:
		return p.GenericHandler(s, subject)
	}

	return result, nil
}

// GenericHandler : Generates a GenericComponentMsg depending on the event thrown
func (p *Publisher) GenericHandler(s *map[string]interface{}, subject string) (string, error) {
	id, _ := (*s)["id"].(string)
	output := GenericComponentMsg{
		Service: id,
		Status:  "processing",
	}

	parts := strings.Split(subject, ".")
	if len(parts) == 2 {
		if parts[1] == "find" {
			body, err := json.Marshal(s)
			if err != nil {
				log.Println("Can't marshal current service")
			}
			data := string(body)
			t, _ := (*s)["type"]
			output.Type = t.(string)
			output.AWSAccessKeyID = MapString(data, "$(datacenters.items.0.aws_access_key_id)")
			output.AWSSecretAccessKey = MapString(data, "$(datacenters.items.0.aws_secret_access_key)")
			output.DatacenterRegion = MapString(data, "$(datacenters.items.0.region)")

			tags := make(map[string]string)
			tags["ernest.service"], _ = (*s)["name"].(string)
			output.Tags = tags

			marshalled, err := json.Marshal(output)
			if err != nil {
				log.Println(err)
				return "", errors.New(err.Error())
			}

			return string(marshalled), nil
		}
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

// MapString : fills a templated string field on its mapped value
func MapString(data string, value string) string {
	if len(value) > 3 && value[0:2] == "$(" && value[len(value)-1:len(value)] == ")" {
		q := gjson.Get(data, value[2:len(value)-1]).String()
		if len(q) > 3 && q[0:2] == "$(" && q[len(q)-1:len(q)] == ")" {
			return MapString(data, q)
		} else if q != "" && q != "null" {
			return q
		}
		return value
	}
	return value
}

// MapHash : finds and replaces templated values on a hash
func MapHash(data string, value map[string]interface{}) map[string]interface{} {
	for field, selector := range value {
		switch v := selector.(type) {
		case string:
			value[field] = MapString(data, v)
		case []interface{}:
			value[field] = MapSlice(data, v)
		case map[string]interface{}:
			value[field] = MapHash(data, v)
		}
	}
	return value
}

// MapSlice : finds and replace templated strings on a slice
func MapSlice(data string, values []interface{}) []interface{} {
	for i := 0; i < len(values); i++ {
		switch v := values[i].(type) {
		case string:
			values[i] = MapString(data, v)
		case []interface{}:
			values[i] = MapSlice(data, v)
		case map[string]interface{}:
			values[i] = MapHash(data, v)
		}
	}
	return values
}

// UpdateTemplateVariables : replaces any qjson queries in fields with information from the current service build
func (p *Publisher) UpdateTemplateVariables(items []interface{}, s *map[string]interface{}) []interface{} {
	body, err := json.Marshal(s)
	if err != nil {
		log.Println("Can't marshal current service")
		return items
	}
	data := string(body)

	for i, v := range items {
		item := v.(map[string]interface{})
		items[i] = MapHash(data, item)
	}

	return items
}

// isSupportedMessage : checks if a message is supported or not
func (p *Publisher) isSupportedMessage(s *map[string]interface{}, subject string) bool {
	w, _ := NewWorkflow(s)
	valid := w.transitions()
	for _, v := range valid {
		if v == subject {
			return true
		}
	}

	return false
}

// FinishProcessing : finishes a service processation setting the final status
func (p *Publisher) FinishProcessing(s *map[string]interface{}, status string) string {
	(*s)["status"] = status
	marshalled, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
		return ""
	}

	id, _ := (*s)["id"].(string)
	natsClient.Request("service.set", []byte(`{"id":"`+id+`","status":"`+status+`"}`), time.Second)

	return string(marshalled)
}
