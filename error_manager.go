/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"strings"
)

type errorManager struct{}

func (em *errorManager) isAnErrorMessage(subject string) bool {
	switch subject[len(subject)-14:] {
	case "s.create.error",
		"s.delete.error",
		"s.update.error":
		return true
	}
	return false
}

func (em *errorManager) markAsFailed(s *service, subject string, body []byte) *service {
	switch subject[len(subject)-14:] {
	case "s.create.error":
		em.markComponentCreationAsFailed(s, subject, body)
	case "s.delete.error":
		em.markComponentDeletionAsFailed(s, subject, body)
	case "s.update.error":
		em.markComponentModificationAsFailed(s, subject, body)
	}
	s.markAsFailed()

	return s
}

func (em *errorManager) getInputList(body []byte) GenericComponentMsg {
	input := GenericComponentMsg{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Panic(err.Error())
	}

	return input
}

func (em *errorManager) markComponentCreationAsFailed(s *service, subject string, body []byte) *service {
	parts := strings.Split(subject, ".")
	input := em.getInputList(body)
	s.transferCreated(parts[0], input)

	return s
}

func (em *errorManager) markComponentModificationAsFailed(s *service, subject string, body []byte) *service {
	parts := strings.Split(subject, ".")
	input := em.getInputList(body)
	s.transferUpdated(parts[0], input)

	return s
}

func (em *errorManager) markComponentDeletionAsFailed(s *service, subject string, body []byte) *service {
	parts := strings.Split(subject, ".")
	input := em.getInputList(body)
	s.transferDeleted(parts[0], input)

	return s
}
