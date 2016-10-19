/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"strings"
)

// ErrorManager : manages error messages
type ErrorManager struct{}

// isAnErrorMessage : checks if the received message is an error or not
func (em *ErrorManager) isAnErrorMessage(subject string) bool {
	// Checking the last part of the messages subject to determine if there has been an error
	switch subject[len(subject)-14:] {
	case "s.create.error",
		"s.delete.error",
		"s.update.error":
		return true
	}
	return false
}

// markAsFailed : marks as message as failed
func (em *ErrorManager) markAsFailed(s *map[string]interface{}, subject string, body []byte) {
	parts := strings.Split(subject, ".")
	input := NewGenericComponentMsg(body)

	// Checking the last part of the messages subject to determine if there has been an error
	switch subject[len(subject)-14:] {
	case "s.create.error":
		TransferCreated(s, parts[0], input)
	case "s.delete.error":
		TransferUpdated(s, parts[0], input)
	case "s.update.error":
		TransferDeleted(s, parts[0], input)
	}

	(*s)["status"] = "pre-failed"
}
