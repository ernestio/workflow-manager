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

// This method is mapping received messages to internal methods
func (sub *subscriber) MethodName(subject string) (string, error) {
	m := make(map[string]string)

	m["test.message"] = "DummyTest"
	m["service.create"] = "ServiceCreate"
	m["service.delete"] = "ServiceDelete"
	m["service.patch"] = "ServicePatch"
	m["routers.create.done"] = "RoutersCreateDone"
	m["routers.delete.done"] = "RoutersDeleteDone"
	m["networks.create.done"] = "NetworksCreateDone"
	m["networks.delete.done"] = "NetworksDeleteDone"
	m["instances.create.done"] = "InstancesCreateDone"
	m["instances.delete.done"] = "InstancesDeleteDone"
	m["instances.update.done"] = "InstancesUpdateDone"
	m["firewalls.create.done"] = "FirewallsCreateDone"
	m["firewalls.delete.done"] = "FirewallsDeleteDone"
	m["firewalls.update.done"] = "FirewallsUpdateDone"
	m["nats.create.done"] = "NatsCreateDone"
	m["nats.delete.done"] = "NatsDeleteDone"
	m["nats.update.done"] = "NatsUpdateDone"
	m["executions.create.done"] = "ExecutionsCreateDone"

	if val, ok := m[subject]; ok {
		return val, nil
	}
	return "", errors.New("Message not supported")
}

// This method is here just for testing / educational purposes
func (sub *subscriber) DummyTest(s *service, subject string, body []byte) *service {
	s.Name = "hello world from subscriber!"

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

// A routers.create.done event is emmited when all routers have
// been created, so in this method we will be processing this
// message and storing the routers data
func (sub *subscriber) RoutersCreateDone(s *service, subject string, body []byte) *service {
	m := RoutersCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Println(err)
		return nil
	}

	messages := []MonitorMessage{}
	for i, sr := range s.Routers.Items {
		for _, mr := range m.Routers {
			if sr.Name == mr.Name {
				s.Routers.Items[i] = mr
				s.Endpoint = mr.IP
				messages = append(messages, MonitorMessage{Body: "\t" + mr.IP, Level: ""})
			}
		}
	}
	messages = append(messages, MonitorMessage{Body: "Routers successfully created", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// A networks.create.done event is emmited when all networks have
// been created, so in this method we will be processing this
// message and storing the networks data
func (sub *subscriber) NetworksCreateDone(s *service, subject string, body []byte) *service {
	m := NetworksCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	for i, sr := range s.Networks.Items {
		for _, mr := range m.Networks {
			if sr.Name == mr.Name {
				s.Networks.Items[i].Name = mr.Name
				s.Networks.Items[i].Range = mr.Range
				s.Networks.Items[i].Netmask = mr.Netmask
				s.Networks.Items[i].StartAddress = mr.StartAddress
				s.Networks.Items[i].EndAddress = mr.EndAddress
				s.Networks.Items[i].Gateway = mr.Gateway
				s.Networks.Items[i].Status = mr.Status
			}
		}
	}

	s.Networks.Status = m.Status
	s.Networks.ErrorCode = m.ErrorCode
	s.Networks.ErrorMessage = m.ErrorMessage

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Networks successfully created", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// A routers.delete.done event is emmited when all networks have
// been deleted, so in this method we will be processing this
// message and deleting the networks data
func (sub *subscriber) RoutersDeleteDone(s *service, subject string, body []byte) *service {
	s.RoutersToDelete.Items = make([]router, 0)
	s.Routers.Items = make([]router, 0)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Routers deleted", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// A networks.delete.done event is emmited when all networks have
// been deleted, so in this method we will be processing this
// message and deleting the networks data
func (sub *subscriber) NetworksDeleteDone(s *service, subject string, body []byte) *service {
	s.NetworksToDelete.Items = make([]network, 0)
	s.Networks.Items = make([]network, 0)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Networks deleted", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// A instances.create.done event is emmited when all instances have
// been created, so in this method we will be processing this
// message and storing the instances data
func (sub *subscriber) InstancesCreateDone(s *service, subject string, body []byte) *service {
	m := InstancesCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}
	for i := range s.Instances.Items {
		for j := range m.Instances {
			if s.Instances.Items[i].Name == m.Instances[j].Name {
				s.Instances.Items[i].Status = m.Instances[j].Status
			}
		}
	}

	s.Instances.Status = m.Status
	s.Instances.ErrorCode = m.ErrorCode
	s.Instances.ErrorMessage = m.ErrorMessage

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Instances successfully created", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

func (sub *subscriber) InstancesUpdateDone(s *service, subject string, body []byte) *service {
	s.InstancesToUpdate.Items = make([]instance, 0)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Instances successfully updated", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

func (sub *subscriber) InstancesDeleteDone(s *service, subject string, body []byte) *service {
	s.InstancesToDelete.Items = make([]instance, 0)

	var instances []instance
	for i, instance := range s.Instances.Items {
		deleted := false
		for _, d := range s.InstancesToDelete.Items {
			if instance.Name == d.Name {
				deleted = true
				s.Instances.Items[i].Status = d.Status
			}
		}
		if deleted == false {
			instances = append(instances, instance)
		}
	}
	s.Instances.Items = instances

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Instances deleted", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// A firewalls.create.done event is emmited when all firewalls have
// been created, so in this method we will be processing this
// message and storing the firewalls data
func (sub *subscriber) FirewallsCreateDone(s *service, subject string, body []byte) *service {
	m := FirewallsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Println(err)
		return nil
	}

	for i, sr := range s.Firewalls.Items {
		for _, mr := range m.Firewalls {
			if sr.Name == mr.Name {
				s.Firewalls.Items[i].Status = mr.Status
			}
		}
	}

	s.Firewalls.Status = m.Status
	s.Firewalls.ErrorCode = m.ErrorCode
	s.Firewalls.ErrorMessage = m.ErrorMessage

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Firewalls Created", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// A firewalls.update.done event is emmited when all firewalls have
// been created, so in this method we will be processing this
// message and storing the firewalls data
func (sub *subscriber) FirewallsUpdateDone(s *service, subject string, body []byte) *service {
	m := FirewallsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Println(err)
		return nil
	}

	for i, sr := range s.Firewalls.Items {
		for _, mr := range m.Firewalls {
			if sr.Name == mr.Name {
				s.Firewalls.Items[i].Status = mr.Status
			}
		}
	}
	s.Firewalls.Status = m.Status
	s.Firewalls.ErrorCode = m.ErrorCode
	s.Firewalls.ErrorMessage = m.ErrorMessage

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Firewalls Updated", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// A firewalls.delete.done event is emmited when all firewalls have
// been deleted, so in this method we will be processing this
// message and storing the firewalls data
func (sub *subscriber) FirewallsDeleteDone(s *service, subject string, body []byte) *service {
	s.FirewallsToDelete.Items = make([]firewall, 0)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Firewalls Deleted", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// A nats.create.done event is emmited when all nats have
// been created, so in this method we will be processing this
// message and storing the nats data
func (sub *subscriber) NatsCreateDone(s *service, subject string, body []byte) *service {
	m := NatsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Println(err)
		return nil
	}

	for i, sr := range s.Nats.Items {
		for _, mr := range m.Nats {
			if sr.Name == mr.Name {
				s.Nats.Items[i].Status = mr.Status
			}
		}
	}

	s.Nats.Status = m.Status
	s.Nats.ErrorCode = m.ErrorCode
	s.Nats.ErrorMessage = m.ErrorMessage

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Nats Created", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// A nats.update.done event is emmited when all nats have
// been created, so in this method we will be processing this
// message and storing the nats data
func (sub *subscriber) NatsUpdateDone(s *service, subject string, body []byte) *service {
	m := NatsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		log.Println(err)
		return nil
	}

	for i, sr := range s.Nats.Items {
		for _, mr := range m.Nats {
			if sr.Name == mr.Name {
				s.Nats.Items[i].Status = mr.Status
			}
		}
	}

	s.Nats.Status = m.Status
	s.Nats.ErrorCode = m.ErrorCode
	s.Nats.ErrorMessage = m.ErrorMessage

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Nats Updated", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return s
}

// A nats.delete.done event is emmited when all nats have
// been deleted, so in this method we will be processing this
// message and storing the nats data
func (sub *subscriber) NatsDeleteDone(s *service, subject string, body []byte) *service {
	s.NatsToDelete.Items = make([]nat, 0)

	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Nats Deleted", Level: "INFO"})
	UserOutput(s.Channel(), messages)

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
