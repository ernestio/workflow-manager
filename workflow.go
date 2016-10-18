/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
)

// Workflow : object is a representation for the json that represents the
// service creation workflow graph
type Workflow struct {
	Arcs []Arc `json:"arcs"`
}

// Arc : or transition is the definition of an event that happens when the
// service on a status "from" receives an "event" and becomes on status
// "to"
type Arc struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Event string `json:"event"`
}

// NewWorkflow : generates a workflow based on the input service
func NewWorkflow(s *map[string]interface{}) (Workflow, error) {
	var ser service

	raw, _ := json.Marshal(s)
	err := json.Unmarshal(raw, &ser)

	return ser.Workflow, err
}

// nextArc : Get next arc for the current workflow definition for a given status and
// event
func (w *Workflow) nextArc(status string, event string) (*Arc, error) {
	var a = &Arc{}
	var err = errors.New("No arcs matching your request")
	for i := 0; i < len(w.Arcs); i++ {
		if event == w.Arcs[i].Event && w.Arcs[i].From == status {
			a = &w.Arcs[i]
			err = nil
			break
		}
	}

	return a, err
}

// nextEvent : Get next event for the current workflow definition for a given status
func (w *Workflow) nextEvent(status string) (string, error) {
	for i := 0; i < len(w.Arcs); i++ {
		if w.Arcs[i].From == status {
			return w.Arcs[i].Event, nil
		}
	}
	return "", errors.New("No new event defined")
}

// transitions : gets all the events on current workflow
func (w *Workflow) transitions() (transitions []string) {
	for _, a := range w.Arcs {
		transitions = append(transitions, a.Event)
	}
	return transitions
}
