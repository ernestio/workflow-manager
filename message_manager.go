/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
)

// Message manager is a group of methods that does the magic to allow developers
// worry only by getting updated its provider and subscriber files
type messageManager struct {
}

// Will call the publisher for a specified message and return the string with the
// message to be published
func (mm *messageManager) preparePublishMessage(subject string, s *service) (string, error) {
	var p publisher
	methodName, err := p.MethodName(subject)
	if err != nil {
		log.Printf("Message not supported: %s", subject)
		return "", errors.New("Message not supported")
	}

	inputs := make([]reflect.Value, 1)
	inputs[0] = reflect.ValueOf(s)
	outputs := reflect.ValueOf(&p).MethodByName(methodName).Call(inputs)
	return outputs[0].String(), nil
}

// It gets a message subject and the body received and calls the necessary
// subscriber methods to read them into a service object
func (mm *messageManager) getServiceFromMessage(subject string, body []byte) (*service, string, error) {

	var sub subscriber
	methodName, err := sub.MethodName(subject)
	if err != nil {
		e := errorManager{}
		if e.isAnErrorMessage(subject) {
			s, err := mm.getService(body)
			if err == nil {
				s = e.markAsFailed(s, subject, body)
			}
			return s, "to_error", nil
		}
		log.Printf("Message not supported: %s", subject)
		return nil, "", errors.New("Message not supported")
	}
	s, err := mm.getService(body)

	inputs := make([]reflect.Value, 3)
	inputs[0] = reflect.ValueOf(s)
	inputs[1] = reflect.ValueOf(subject)
	inputs[2] = reflect.ValueOf(body)

	outputs := reflect.ValueOf(&sub).MethodByName(methodName).Call(inputs)
	s = outputs[0].Interface().(*service)

	return s, subject, nil
}

// Creates or gets a persisted service based on the service field of the
// message body
func (mm *messageManager) getService(body []byte) (*service, error) {
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
