/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
)

type components struct {
	Finished             string        `json:"finished"`
	Items                []interface{} `json:"items"`
	Started              string        `json:"started"`
	Status               string        `json:"status"`
	SequentialProcessing bool          `json:"sequential_processing"`
}

// This is the object representation for a service inside the
// FSM, it has appended the workflow the service needs to
// follow to be built
type service struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Body          string   `json:"body"`
	Workflow      workflow `json:"workflow"`
	Started       string   `json:"started"`
	Finished      string   `json:"finished"`
	Status        string   `json:"status"`
	Type          string   `json:"type"`
	ClientName    string   `json:"client_name"`
	Parent        string   `json:"-"`
	Bootstrapping string   `json:"bootstrapping"`
	ErnestIP      []string `json:"ernest_ip"`
	ServiceIP     string   `json:"service_ip"`
	Options       struct {
		Password string `json:"password"`
		User     string `json:"user"`
	} `json:"options"`
	Datacenters        components `json:"datacenters"`
	VPCs               components `json:"vpcs"`
	VPCsToCreate       components `json:"vpcs_to_create"`
	VPCsToDelete       components `json:"vpcs_to_delete"`
	Bootstraps         components `json:"bootstraps"`
	BootstrapsToCreate components `json:"bootstraps_to_create"`
	Executions         components `json:"executions"`
	ExecutionsToCreate components `json:"executions_to_create"`
	Firewalls          components `json:"firewalls"`
	FirewallsToCreate  components `json:"firewalls_to_create"`
	FirewallsToUpdate  components `json:"firewalls_to_update"`
	FirewallsToDelete  components `json:"firewalls_to_delete"`
	Instances          components `json:"instances"`
	InstancesToCreate  components `json:"instances_to_create"`
	InstancesToUpdate  components `json:"instances_to_update"`
	InstancesToDelete  components `json:"instances_to_delete"`
	Nats               components `json:"nats"`
	NatsToCreate       components `json:"nats_to_create"`
	NatsToUpdate       components `json:"nats_to_update"`
	NatsToDelete       components `json:"nats_to_delete"`
	Networks           components `json:"networks"`
	NetworksToCreate   components `json:"networks_to_create"`
	NetworksToUpdate   components `json:"networks_to_update"`
	NetworksToDelete   components `json:"networks_to_delete"`
	Routers            components `json:"routers"`
	RoutersToCreate    components `json:"routers_to_create"`
	RoutersToDelete    components `json:"routers_to_delete"`
	ELBs               components `json:"elbs"`
	ELBsToCreate       components `json:"elbs_to_create"`
	ELBsToUpdate       components `json:"elbs_to_update"`
	ELBsToDelete       components `json:"elbs_to_delete"`
}

type message struct {
	Service string `json:"service"`
}

func (s *service) markAsFailed() {
	s.Status = "pre-failed"
}

// Persist current service
func (s *service) save() error {
	json, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
		return err
	}
	err = p.set(s.ID, string(json))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *service) del() {
	p.del(s.ID)
}

func (s *service) Channel() string {
	return s.ID
}

func (s *service) asMap() (mapped map[string]interface{}) {
	body, err := json.Marshal(s)
	if err != nil {
		log.Panic(err.Error())
	}
	err = json.Unmarshal(body, &mapped)
	if err != nil {
		log.Panic(err.Error())
	}
	return mapped
}

func (s *service) loadFromMap(mapped map[string]interface{}) {
	body, err := json.Marshal(mapped)
	if err != nil {
		log.Panic(err.Error())
	}
	err = json.Unmarshal(body, s)
}

func (s *service) getComponentList(cType string) []interface{} {
	tmp := s.asMap()
	cList := tmp[cType].(map[string]interface{})
	list := cList["items"].([]interface{})

	return list
}

func (s *service) transferCreated(cType string, input GenericComponentMsg) {
	var components []interface{}
	var erroredComponents []interface{}

	tmp := s.asMap()
	inputComponents := input.Components
	currentComponents := tmp[cType].(map[string]interface{})
	if currentComponents["items"] != nil {
		components = currentComponents["items"].([]interface{})
	}

	// Append new components
	for _, c := range inputComponents {
		inHash := c.(map[string]interface{})
		status := inHash["status"].(string)
		if status == "errored" {
			erroredComponents = append(erroredComponents, c)
		} else {
			components = append(components, c)
		}
	}
	currentComponents["status"] = "completed"
	currentComponents["items"] = components

	// Remove to be created components
	componentsToBeProcessed := tmp[cType+"_to_create"].(map[string]interface{})
	componentsToBeProcessed["items"] = erroredComponents
	componentsToBeProcessed["status"] = input.Status
	componentsToBeProcessed["error_code"] = input.ErrorCode
	componentsToBeProcessed["error_message"] = input.ErrorMessage

	// Save result
	s.loadFromMap(tmp)
}

func (s *service) transferUpdated(cType string, input GenericComponentMsg) {
	var components []interface{}
	var erroredComponents []interface{}

	tmp := s.asMap()
	inputComponents := input.Components

	currentComponents := tmp[cType].(map[string]interface{})
	if currentComponents["items"] != nil {
		components = currentComponents["items"].([]interface{})
	}

	// Append new components
	for _, c := range inputComponents {
		for i, v := range components {
			inHash := c.(map[string]interface{})
			exHash := v.(map[string]interface{})
			if inHash["name"] != nil && exHash["name"] != nil {
				iName := inHash["name"].(string)
				name := exHash["name"].(string)
				if iName == name {
					status := inHash["status"].(string)
					if status == "completed" {
						components[i] = c
					} else {
						erroredComponents = append(erroredComponents, c)
					}
				}
			}
		}
	}
	currentComponents["status"] = "completed"
	currentComponents["items"] = components

	// Remove to be created components
	componentsToBeProcessed := tmp[cType+"_to_update"].(map[string]interface{})

	componentsToBeProcessed["items"] = erroredComponents
	componentsToBeProcessed["status"] = input.Status
	componentsToBeProcessed["error_code"] = input.ErrorCode
	componentsToBeProcessed["error_message"] = input.ErrorMessage

	// Save result
	s.loadFromMap(tmp)
}

func (s *service) transferDeleted(cType string, input GenericComponentMsg) {
	var components []interface{}
	var remanentComponents []interface{}
	var erroredComponents []interface{}

	tmp := s.asMap()
	inputComponents := input.Components

	currentComponents := tmp[cType].(map[string]interface{})
	if currentComponents["items"] != nil {
		components = currentComponents["items"].([]interface{})
	}

	for _, v := range components {
		sw := false
		exHash := v.(map[string]interface{})
		name := exHash["name"].(string)
		for _, c := range inputComponents {
			inHash := c.(map[string]interface{})
			iName := inHash["name"].(string)
			if iName == name {
				status := inHash["status"].(string)
				if status == "errored" {
					erroredComponents = append(erroredComponents, c)
				} else {
					sw = true
				}
			}
		}
		if sw == false {
			remanentComponents = append(remanentComponents, v)
		}
	}
	currentComponents["status"] = "completed"
	currentComponents["items"] = remanentComponents

	// Remove to be created components
	componentsToBeProcessed := tmp[cType+"_to_delete"].(map[string]interface{})

	componentsToBeProcessed["items"] = erroredComponents
	if len(erroredComponents) > 0 {
		componentsToBeProcessed["status"] = "errored"
	} else {
		componentsToBeProcessed["status"] = "completed"
	}
	componentsToBeProcessed["error_code"] = input.ErrorCode
	componentsToBeProcessed["error_message"] = input.ErrorMessage

	// Save result
	s.loadFromMap(tmp)
}
