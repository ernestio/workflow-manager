/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"log"
	"time"
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

// This method is mapping messages to be sent to internal methods
func (p *publisher) MethodName(subject string) (string, error) {
	m := make(map[string]string)

	m["test.message"] = "DummyTest"
	m["routers.create"] = "CreateRouters"
	m["routers.delete"] = "DeleteRouters"
	m["service.create.error"] = "ServiceCreateError"
	m["service.create.done"] = "ServiceCreateDone"
	m["service.delete.error"] = "ServicesDeleteError"
	m["service.delete.done"] = "ServiceDeleteDone"
	m["networks.create"] = "NetworksCreate"
	m["networks.delete"] = "NetworksDelete"
	m["instances.create"] = "InstancesCreate"
	m["instances.delete"] = "InstancesDelete"
	m["instances.update"] = "InstancesUpdate"
	m["nats.create"] = "NatsCreate"
	m["nats.delete"] = "NatsDelete"
	m["nats.update"] = "NatsUpdate"
	m["firewalls.create"] = "FirewallsCreate"
	m["firewalls.delete"] = "FirewallsDelete"
	m["firewalls.update"] = "FirewallsUpdate"
	m["executions.create"] = "ExecutionsCreate"

	if val, ok := m[subject]; ok {
		return val, nil
	}

	return "", errors.New("Message not supported")
}

// This method is here just for testing / educational purposes
func (p *publisher) DummyTest(s *service) string {
	return "hello world from publisher!"
}

// Prepares a message to create routers
func (p *publisher) CreateRouters(s *service) string {
	m := buildCreateRouters(s)
	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

// Prepares a message to delete routers
func (p *publisher) DeleteRouters(s *service) string {
	m := buildDeleteRouters(s)
	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
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

func (p *publisher) NetworksCreate(s *service) string {
	m := buildCreateNetworks(s)
	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
}

func (p *publisher) NetworksDelete(s *service) string {
	m := buildDeleteNetworks(s)
	marshalled, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(marshalled)
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
	if s.Bootstraps.Finished == "yes" || len(s.Bootstraps.Items) == 0 {
		m = buildCreateExecutions(s)
	} else {
		m = buildCreateBootstraps(s)
	}

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
