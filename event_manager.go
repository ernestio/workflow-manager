/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"
	"log"
)

type eventManager struct {
}

// Manage a trigger based on a given definition
func (em *eventManager) manage(subject string, s *map[string]interface{}) (string, error) {
	err := em.move(s, subject)
	if err != nil {
		log.Println(err)
		return "", err
	}
	event := em.next(s)

	return event, err
}

// Prepares a proper message and sends the next event
func (em *eventManager) next(s *map[string]interface{}) string {
	w, _ := ParseWorkflow(s)
	status, _ := (*s)["status"].(string)
	event, err := w.nextEvent(status)
	if err != nil {
		log.Println(err)
		return ""
	}

	return event
}

// Moves a service to its next status and return a
// string with it
func (em *eventManager) move(s *map[string]interface{}, event string) error {
	// Is a valid transition?
	status, _ := (*s)["status"].(string)
	if status == "" {
		status = "created"
	}
	(*s)["status"] = status

	w, _ := ParseWorkflow(s)
	a, err := w.nextArc(status, event)
	if err != nil {
		return errors.New("Invalid status(" + status + ") event (" + event + ") pair")
	}

	// Update status
	(*s)["status"] = a.To

	// Return new status
	return nil
}
