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

func (sub *subscriber) Process(s *service, subject string, body []byte) (*service, bool, string) {
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
	case "executions.create.done":
		sub.ExecutionsCreateDone(s, subject, body)
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

func (sub *subscriber) isSupportedMessage(s *service, subject string) bool {
	valid := s.Workflow.transitions()
	for _, v := range valid {
		if v == subject {
			return true
		}
	}
	if subject == "service.delete" {
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
func (sub *subscriber) GenericCreation(s *service, subject string, body []byte) *service {
	parts := strings.Split(subject, ".")
	input := sub.getInputList(body)
	s.transferCreated(parts[0], input)

	return s
}

// GenericModification : Will process generic messages
func (sub *subscriber) GenericModification(s *service, subject string, body []byte) *service {
	parts := strings.Split(subject, ".")
	input := sub.getInputList(body)
	s.transferUpdated(parts[0], input)

	return s
}

// GenericDeletion : Will process generic messages
func (sub *subscriber) GenericDeletion(s *service, subject string, body []byte) *service {
	parts := strings.Split(subject, ".")
	input := sub.getInputList(body)
	s.transferDeleted(parts[0], input)

	return s
}

// Entry point to the flow environment creation, it will create the service and attach
// a default workflow to it
func (sub *subscriber) ServiceCreate(s *service, subject string, body []byte) *service {

	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		return nil
	}

	w := &s.Workflow
	if len(w.Arcs) == 0 {
		w = &workflow{}
		w.loadDefault()
		s.Workflow = *w
	}
	natsClient.Request("service.set", []byte(`{"id":"`+s.ID+`","status":"in_progress"}`), time.Second)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Starting environment creation", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// Entry point to the flow environment deletion, it will trigger a cleanup of the
// entire service
func (sub *subscriber) ServiceDelete(s *service, subject string, body []byte) *service {
	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		return nil
	}
	s.Status = "created"
	natsClient.Request("service.set", []byte(`{"id":"`+s.ID+`","status":"in_progress"}`), time.Second)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Starting environment deletion", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// ServicePatch Entry point to the flow environment patching, it will create the service and attach
// a default workflow to it
func (sub *subscriber) ServicePatch(s *service, subject string, body []byte) *service {
	if err := json.Unmarshal(body, &s); err != nil {
		log.Println(err)
		return nil
	}
	s.Status = ""

	return s
}

// A executions.create.done event is emmited when all bootstraps/executions have
// been created, so in this method we will be processing this
// message and storing the executions data
// When all executions have been completed, service.create.done will be emitted
func (sub *subscriber) ExecutionsCreateDone(s *service, subject string, body []byte) *service {
	m := ExecutionsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Println(err)
		return nil
	}

	if s.Status == "bootstrapping" {
		for i, sr := range s.Bootstraps.Items {
			s.Bootstraps.Finished = "yes"
			for _, mr := range m.Executions {
				if sr.Name == mr.Name {
					s.Bootstraps.Items[i].MatchedInstances = mr.MatchedInstances
					s.Bootstraps.Items[i].ExecutionStatus = mr.ExecutionStatus
					s.Bootstraps.Items[i].Status = mr.Status
					// s.Bootstraps.Items[i].Reports = mr.Reports
				}
			}
		}
		if len(s.Bootstraps.Items) > 0 {
			messages := []MonitorMessage{}
			messages = append(messages, MonitorMessage{Body: "Instances bootstrapped", Level: "INFO"})
			UserOutput(s.Channel(), messages)
		}
	} else if s.Status == "running_executions" {
		for _, mr := range m.Executions {
			if sr := s.executionByName(mr.Name); sr != nil {
				sr.Payload = mr.Payload
				sr.Target = mr.Target
				sr.MatchedInstances = mr.MatchedInstances
				sr.ExecutionStatus = mr.ExecutionStatus
				sr.Status = mr.Status
				// s.Executions.Items[i].Reports = mr.Reports
			} else {
				ex := execution{
					Type:             mr.Type,
					Name:             mr.Name,
					Payload:          mr.Payload,
					Target:           mr.Target,
					MatchedInstances: mr.MatchedInstances,
					ExecutionStatus:  mr.ExecutionStatus,
					Created:          mr.Created,
				}
				ex.Status = mr.Status
				s.Executions.Items = append(s.Executions.Items, ex)
			}
		}
		if len(s.ExecutionsToCreate.Items) > 0 {
			messages := []MonitorMessage{}
			messages = append(messages, MonitorMessage{Body: "Executions ran", Level: "INFO"})
			UserOutput(s.Channel(), messages)
		}

		// Clear executions
		s.ExecutionsToCreate.Items = []execution{}
	}

	return s
}
