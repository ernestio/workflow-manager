/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"strings"
)

// ErrorSubjects : Subjects to be exported
var ErrorSubjects = []string{"s.create.error", "s.delete.error", "s.update.error", "s.find.error"}

// ErrorManager : manages error messages
type ErrorManager struct{}

// isAnErrorMessage : checks if the received message is an error or not
func (em *ErrorManager) isAnErrorMessage(subject string) bool {
	// Checking the last part of the messages subject to determine if there has been an error
	if getErrorType(subject) != "" {
		return true
	}
	return false
}

// markAsFailed : marks as message as failed
func (em *ErrorManager) markAsFailed(s *map[string]interface{}, subject string, body []byte) {
	parts := strings.Split(subject, ".")
	input := NewGenericComponentMsg(body)

	// Checking the last part of the messages subject to determine if there has been an error
	switch getErrorType(subject) {
	case "s.create.error":
		TransferCreated(s, parts[0], input)
	case "s.delete.error":
		TransferUpdated(s, parts[0], input)
	case "s.update.error":
		TransferDeleted(s, parts[0], input)
	case "s.find.error":
		TransferFound(s, parts[0], input)
	}

	(*s)["last_known_error"] = em.getErrorMessage(input)
	(*s)["status"] = "pre-failed"
}

func (em *ErrorManager) getErrorMessage(input GenericComponentMsg) string {
	for _, c := range input.Components {
		inHash := c.(map[string]interface{})
		status := inHash["status"].(string)
		if status == "errored" {
			err, ok := inHash["error"].(string)
			if ok {
				return err
			}
			return "Internal error: 00001"
		}
	}
	return ""
}

func getErrorType(subject string) string {
	for _, v := range ErrorSubjects {
		if strings.Contains(subject, v) {
			return v
		}
	}
	return ""
}
