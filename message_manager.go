/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"strings"
)

// MessageManager is a group of methods that does the magic to allow developers
// worry only by getting updated its provider and subscriber files
type MessageManager struct {
}

// Will call the publisher for a specified message and return the string with the
// message to be published
func (mm *MessageManager) preparePublishMessage(subject string, s *map[string]interface{}) (string, error) {
	var p Publisher

	return p.Process(s, subject)
}

// It gets a message subject and the body received and calls the necessary
// subscriber methods to read them into a service object
func (mm *MessageManager) getServiceFromMessage(subject string, body []byte) (map[string]interface{}, string, error) {
	var sub Subscriber

	if err := mm.validateSubject(subject); err != nil {
		return nil, "", errors.New("Message not supported")
	}

	s, err := mm.getService(body)
	if err != nil {
		return nil, "", errors.New("Message not supported")
	}

	supported, status := sub.Process(&s, subject, body)

	if status != "" {
		em := ErrorManager{}
		em.markAsFailed(&s, subject, body)
		return s, status, nil
	}

	if supported == false {
		return nil, "", errors.New("Message not supported")
	}

	return s, subject, nil
}

func (mm *MessageManager) validateSubject(subject string) error {
	parts := strings.Split(subject, ".")
	if len(parts) == 2 && parts[0] != "service" {
		return errors.New("Message not supported")
	}
	if len(parts) < 2 {
		return errors.New("Message not supported")
	}
	if parts[1] != "create" && parts[1] != "update" && parts[1] != "delete" && parts[1] != "patch" {
		return errors.New("Message not supported")
	}
	if subject == "service.create.done" || subject == "service.create.error" {
		return errors.New("Message not supported")
	}

	return nil
}

// Creates or gets a persisted service based on the service field of the
// message body
func (mm *MessageManager) getService(body []byte) (map[string]interface{}, error) {
	type InputMessage struct {
		ID      string `json:"id"`
		Service string `json:"service"`
	}

	m := InputMessage{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}

	serviceID := m.Service
	if serviceID == "" {
		serviceID = m.ID
	}
	if serviceID == "" {
		return nil, errors.New("Unsupported message")
	}

	s := p.getService(serviceID)

	return s, nil
}
