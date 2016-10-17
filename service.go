/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
)

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
}

// SaveService : persists the service
func SaveService(s *map[string]interface{}) error {
	json, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
		return err
	}
	id, _ := (*s)["id"].(string)
	err = p.set(id, string(json))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// ServiceDel : removes the current service
func ServiceDel(s *map[string]interface{}) {
	id, _ := (*s)["id"].(string)
	p.del(id)
}

// TransferCreated : transferst the components_to_created to components array
func TransferCreated(s *map[string]interface{}, cType string, input GenericComponentMsg) {
	var components []interface{}
	var erroredComponents []interface{}

	inputComponents := input.Components
	currentComponents := (*s)[cType].(map[string]interface{})
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
	if componentsToBeProcessed, ok := (*s)[cType+"_to_create"].(map[string]interface{}); ok {
		componentsToBeProcessed["items"] = erroredComponents
		componentsToBeProcessed["status"] = input.Status
		componentsToBeProcessed["error_code"] = input.ErrorCode
		componentsToBeProcessed["error_message"] = input.ErrorMessage
	}
}

// TransferUpdated : updates components with components_to_update data
func TransferUpdated(s *map[string]interface{}, cType string, input GenericComponentMsg) {
	var components []interface{}
	var erroredComponents []interface{}

	inputComponents := input.Components

	currentComponents := (*s)[cType].(map[string]interface{})
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
	if componentsToBeProcessed, ok := (*s)[cType+"_to_update"].(map[string]interface{}); ok {
		componentsToBeProcessed["items"] = erroredComponents
		componentsToBeProcessed["status"] = input.Status
		componentsToBeProcessed["error_code"] = input.ErrorCode
		componentsToBeProcessed["error_message"] = input.ErrorMessage
	}
}

// TrasnferDeleted : removes from components received components_to_delete componets
func TransferDeleted(s *map[string]interface{}, cType string, input GenericComponentMsg) {
	var components []interface{}
	var remanentComponents []interface{}
	var erroredComponents []interface{}

	inputComponents := input.Components

	currentComponents := (*s)[cType].(map[string]interface{})
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
	if componentsToBeProcessed, ok := (*s)[cType+"_to_delete"].(map[string]interface{}); ok {
		componentsToBeProcessed["items"] = erroredComponents
		if len(erroredComponents) > 0 {
			componentsToBeProcessed["status"] = "errored"
		} else {
			componentsToBeProcessed["status"] = "completed"
		}
		componentsToBeProcessed["error_code"] = input.ErrorCode
		componentsToBeProcessed["error_message"] = input.ErrorMessage
	}
}
