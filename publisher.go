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
	case "instances.create":
		result = p.InstancesCreate(s)
	case "instances.delete":
		result = p.InstancesDelete(s)
	case "instances.update":
		result = p.InstancesUpdate(s)
	case "nats.create":
		result = p.NatsCreate(s)
	case "nats.delete":
		result = p.NatsDelete(s)
	case "nats.update":
		result = p.NatsUpdate(s)
	case "firewalls.create":
		result = p.FirewallsCreate(s)
	case "firewalls.delete":
		result = p.FirewallsDelete(s)
	case "firewalls.update":
		result = p.FirewallsUpdate(s)
	case "executions.create":
		result = p.ExecutionsCreate(s)
	case "bootstraps.create":
		result = p.BootstrapsCreate(s)
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
	items := list["items"].([]interface{})
	items = p.Vitamine(items, s)
	output.Components = items

	marshalled, err := json.Marshal(output)
	if err != nil {
		log.Println(err)
		return "", errors.New(err.Error())
	}

	return string(marshalled), nil
}

func (p *publisher) Vitamine(items []interface{}, s *service) []interface{} {
	body, err := json.Marshal(s)
	if err != nil {
		log.Println("Can't marshal current service")
		return items
	}
	json := string(body)

	for _, v := range items {
		item := v.(map[string]interface{})
		for field, selector := range item {
			value, err := selector.(string)
			if err == true && value != "" {
				if value[0:2] == "$(" && value[len(value)-1:len(value)] == ")" {
					item[field] = gjson.Get(json, value[2:len(value)-1]).String()
				}
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

func (p *publisher) InstancesUpdate(s *service) string {
	m := buildUpdateInstances(s)
	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

func (p *publisher) InstancesCreate(s *service) string {
	m := buildCreateInstances(s)
	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

func (p *publisher) InstancesDelete(s *service) string {
	m := buildDeleteInstances(s)
	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

func (p *publisher) NatsCreate(s *service) string {
	m := buildCreateNats(s)
	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

func (p *publisher) NatsDelete(s *service) string {
	m := buildDeleteNats(s)
	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

func (p *publisher) NatsUpdate(s *service) string {
	m := buildUpdateNats(s)
	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

func (p *publisher) FirewallsCreate(s *service) string {
	m := buildCreateFirewalls(s)

	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

func (p *publisher) FirewallsUpdate(s *service) string {
	m := buildUpdateFirewalls(s)

	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

func (p *publisher) FirewallsDelete(s *service) string {
	m := buildDeleteFirewalls(s)

	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

func (p *publisher) ExecutionsCreate(s *service) string {
	m := ExecutionsCreate{}
	m = buildCreateExecutions(s)

	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

func (p *publisher) BootstrapsCreate(s *service) string {
	m := ExecutionsCreate{}
	m = buildCreateBootstraps(s)

	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

// ServiceCreateDone
func (p *publisher) ServiceCreateDone(s *service) string {
	marshalled, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}
	if len(s.Routers.Items) > 0 {
		if s.Endpoint == "" {
			s.Endpoint = s.Routers.Items[0].IP
		}
		if s.ServiceIP == "" {
			s.ServiceIP = s.Routers.Items[0].IP
		}
	}

	natsClient.Request("service.set", []byte(`{"id":"`+s.ID+`","status":"done"}`), time.Second)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "SUCCESS: rules successfully applied", Level: "SUCCESS"})
	messages = append(messages, MonitorMessage{Body: "Your environment endpoint is: " + s.Endpoint, Level: "SUCCESS"})
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
