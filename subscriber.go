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

// Subscriber : When a message is received from the interspace the FSM can't
// process it directly this is where the subscriber comes and translates the
// received message into an fsm understandable service
//
// In order to create your own "translators" just add your mapping for the
// received message to the internal method on the MethodName function.
//
// Then create your own method in order to attach the received information
// to the stored service and return this service, it will be persisted at
// the other side
type Subscriber struct{}

// Process : starts message subscription processing
func (sub *Subscriber) Process(s *map[string]interface{}, subject string, body []byte) (bool, string) {
	e := ErrorManager{}
	if e.isAnErrorMessage(subject) {
		return true, "to_error"
	}

	if sub.isSupportedMessage(s, subject) == false {
		return false, ""
	}

	switch subject {
	case "service.create", "service.import":
		sub.ServiceCreate(s, subject, body)
	case "service.delete":
		sub.ServiceDelete(s, subject, body)
	case "service.patch":
		sub.ServicePatch(s, subject, body)
	default:
		parts := strings.Split(subject, ".")
		if len(parts) != 3 || parts[0] == "service" {
			log.Println("Message not supported : " + subject)
			return false, ""
		}
		input := NewGenericComponentMsg(body)
		switch parts[1] {
		case "create":
			TransferCreated(s, parts[0], input)
		case "update":
			TransferUpdated(s, parts[0], input)
		case "delete":
			TransferDeleted(s, parts[0], input)
		case "find":
			TransferFound(s, parts[0], input)
		default:
			log.Println("Message not supported")
			return false, ""
		}
	}

	return true, ""
}

// isSupportedMessage : checks if a message is supported or not based on the service workflow
func (sub *Subscriber) isSupportedMessage(s *map[string]interface{}, subject string) bool {
	w, _ := NewWorkflow(s)
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

// ServiceCreate : Entry point to the flow environment creation, it will create
// the service and attach a default workflow to it
func (sub *Subscriber) ServiceCreate(s *map[string]interface{}, subject string, body []byte) *map[string]interface{} {

	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		return nil
	}

	id, _ := (*s)["id"].(string)
	natsClient.Request("service.set", []byte(`{"id":"`+id+`","status":"in_progress"}`), time.Second)

	return s
}

// ServiceDelete : Entry point to the flow environment deletion, it will trigger
// a cleanup of the entire service
func (sub *Subscriber) ServiceDelete(s *map[string]interface{}, subject string, body []byte) *map[string]interface{} {
	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		return nil
	}
	id, _ := (*s)["id"].(string)
	(*s)["status"] = "created"
	natsClient.Request("service.set", []byte(`{"id":"`+id+`","status":"in_progress"}`), time.Second)

	return s
}

// ServicePatch Entry point to the flow environment patching, it will create the service and attach
// a default workflow to it
func (sub *Subscriber) ServicePatch(s *map[string]interface{}, subject string, body []byte) *map[string]interface{} {
	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		return nil
	}
	(*s)["status"] = ""

	return s
}
