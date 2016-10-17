/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

// When a message is received from the interspace the FSM can't process it
// directly this is where the subscriber comes and translates the received
// message into an fsm understandable service
//
// In order to create your own "translators" just add your mapping for the
// received message to the internal method on the MethodName function.
//
// Then create your own method in order to attach the received information
// to the stored service and return this service, it will be persisted at
// the other side
type subscriber struct {
}

func (sub *subscriber) Process(s *map[string]interface{}, subject string, body []byte) (*map[string]interface{}, bool, string) {
	e := errorManager{}
	if e.isAnErrorMessage(subject) {
		return s, true, "to_error"
	}

	if sub.isSupportedMessage(s, subject) == false {
		return nil, false, ""
	}

	switch subject {
	case "service.create":
		sub.ServiceCreate(s, subject, body)
	case "service.delete":
		sub.ServiceDelete(s, subject, body)
	case "service.patch":
		sub.ServicePatch(s, subject, body)
	default:
		parts := strings.Split(subject, ".")
		if len(parts) != 3 || parts[0] == "service" {
			log.Println("Message not supported : " + subject)
			return s, false, ""
		}
		switch parts[1] {
		case "create":
			sub.GenericCreation(s, subject, body)
		case "update":
			sub.GenericModification(s, subject, body)
		case "delete":
			sub.GenericDeletion(s, subject, body)
		default:
			log.Println("Message not supported")
			return s, false, ""
		}
	}

	return s, true, ""
}

func (sub *subscriber) isSupportedMessage(s *map[string]interface{}, subject string) bool {
	w, _ := ParseWorkflow(s)
	valid := w.transitions()
	for _, v := range valid {
		if v == subject {
			return true
		}
	}
	if subject == "service.delete" {
		return true
	}
	if subject == "service.patch" {
		return true
	}

	return false
}

func (sub *subscriber) getInputList(body []byte) GenericComponentMsg {
	input := GenericComponentMsg{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Panic(err.Error())
	}

	return input
}

// GenericCreation : Will process generic messages
func (sub *subscriber) GenericCreation(s *map[string]interface{}, subject string, body []byte) *map[string]interface{} {
	parts := strings.Split(subject, ".")
	input := sub.getInputList(body)
	TransferCreated(s, parts[0], input)

	return s
}

// GenericModification : Will process generic messages
func (sub *subscriber) GenericModification(s *map[string]interface{}, subject string, body []byte) *map[string]interface{} {
	parts := strings.Split(subject, ".")
	input := sub.getInputList(body)
	TransferUpdated(s, parts[0], input)

	return s
}

// GenericDeletion : Will process generic messages
func (sub *subscriber) GenericDeletion(s *map[string]interface{}, subject string, body []byte) *map[string]interface{} {
	parts := strings.Split(subject, ".")
	input := sub.getInputList(body)
	TransferDeleted(s, parts[0], input)

	return s
}

// Entry point to the flow environment creation, it will create the service and attach
// a default workflow to it
func (sub *subscriber) ServiceCreate(s *map[string]interface{}, subject string, body []byte) *map[string]interface{} {

	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		return nil
	}

	id, _ := (*s)["id"].(string)
	natsClient.Request("service.set", []byte(`{"id":"`+id+`","status":"in_progress"}`), time.Second)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Starting environment creation", Level: "INFO"})
	UserOutput(id, messages)

	return s
}

// Entry point to the flow environment deletion, it will trigger a cleanup of the
// entire service
func (sub *subscriber) ServiceDelete(s *map[string]interface{}, subject string, body []byte) *map[string]interface{} {
	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		return nil
	}
	id, _ := (*s)["id"].(string)
	(*s)["status"] = "created"
	natsClient.Request("service.set", []byte(`{"id":"`+id+`","status":"in_progress"}`), time.Second)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Starting environment deletion", Level: "INFO"})
	UserOutput(id, messages)

	return s
}

// ServicePatch Entry point to the flow environment patching, it will create the service and attach
// a default workflow to it
func (sub *subscriber) ServicePatch(s *map[string]interface{}, subject string, body []byte) *map[string]interface{} {
	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		return nil
	}
	(*s)["status"] = ""

	return s
}
