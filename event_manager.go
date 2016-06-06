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
func (em *eventManager) manage(subject string, s *service) (string, *service, error) {
	err := em.move(s, subject)
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	event := em.next(s)

	return event, s, err
}

// Prepares a proper message and sends the next event
func (em *eventManager) next(s *service) string {
	event, err := s.Workflow.nextEvent(s.Status)
	if err != nil {
		log.Println(err)
		return ""
	}

	return event
}

// Moves a service to its next status and return a
// string with it
func (em *eventManager) move(s *service, event string) error {
	// Is a valid transition?
	if s.Status == "" {
		s.Status = "created"
	}

	a, err := s.Workflow.nextArc(s.Status, event)
	if err != nil {
		return errors.New("Invalid status(" + s.Status + ") event (" + event + ") pair")
	}

	// Update status
	s.Status = a.To

	// Return new status
	return nil
}
