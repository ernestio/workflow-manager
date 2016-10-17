/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

// Workflow object is a representation for the json that represents the
// service creation workflow graph
type workflow struct {
	Arcs []arc `json:"arcs"`
}

// Arc or transition is the definition of an event that happens when the
// service on a status "from" receives an "event" and becomes on status
// "to"
type arc struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Event string `json:"event"`
}

// Get next arc for the current workflow definition for a given status and
// event
func (w *workflow) nextArc(status string, event string) (*arc, error) {
	var a = &arc{}
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

// Get next event for the current workflow definition for a given status
func (w *workflow) nextEvent(status string) (string, error) {
	for i := 0; i < len(w.Arcs); i++ {
		if w.Arcs[i].From == status {
			return w.Arcs[i].Event, nil
		}
	}
	return "", errors.New("No new event defined")
}

// Check if the current workflow is valid
func (w *workflow) valid() bool {
	// TODO - Check all steps can be reachable
	return true
}

// Check if a status name is configured for the current workflow
func (w *workflow) validStatus(status string) bool {
	for i := 0; i < len(w.Arcs); i++ {
		if status == w.Arcs[i].To {
			return true
		}
	}
	return false
}

// Loads the default workflow (workflow.json) on the current object
func (w *workflow) loadDefault() {
	w.loadWorkflow("workflow.json")
}

func (w *workflow) loadWorkflow(source string) {
	absPath, _ := filepath.Abs(source)
	file, err := os.Open(absPath)
	log.Printf("Reading config from: %s", source)
	if err != nil {
		log.Panic("error:", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&w)
	if err != nil {
		log.Println("Workflow file is invalid")
		log.Panic("error:", err)
	}
}

func (w *workflow) transitions() (transitions []string) {
	for _, a := range w.Arcs {
		transitions = append(transitions, a.Event)
	}
	return transitions
}

func ParseWorkflow(s *map[string]interface{}) (w workflow, err error) {
	var ser service

	raw, _ := json.Marshal(s)
	json.Unmarshal(raw, &ser)

	return ser.Workflow, err
}
