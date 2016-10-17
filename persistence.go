/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats"
)

// Wrapper for redis in order to easily store / recover persisted
// services
type storage struct {
	Nats *nats.Conn
}

type serviceMessage struct {
	ID      string `json:"id"`
	Mapping string `json:"mapping"`
}

// Prepares the connection based on a given config file
func (s *storage) load(n *nats.Conn) {
	s.Nats = n
}

// Get the value for a given key
func (s *storage) get(key string) string {
	if key == "" {
		return ""
	}
	msg, err := natsClient.Request("service.get.mapping", []byte(`{"id":"`+key+`"}`), 1*time.Second)
	if err != nil {
		log.Println(err)
		return ""
	}
	if string(msg.Data) == `{"error":"not found"}` {
		return ""
	}

	return string(msg.Data)
}

// Gets a service object for a given key
func (s *storage) getService(key string) *map[string]interface{} {
	var srv map[string]interface{}

	body := s.get(key)

	if err := json.Unmarshal([]byte(body), &srv); err != nil {
		return &srv
	}
	if srv == nil {
		return &srv
	}

	return &srv
}

// Set a value for a given key
func (s *storage) set(key string, value string) error {
	sm := serviceMessage{}
	sm.ID = key
	sm.Mapping = value
	body, err := json.Marshal(sm)
	_, err = natsClient.Request("service.set.mapping", body, 1*time.Second)
	if err != nil {
		log.Println(err)
		log.Panic("Data can't be stored")
	}
	return err
}

func (s *storage) del(key string) error {
	_, err := natsClient.Request("service.del", []byte(`{"id":"`+key+`"}`), 1*time.Second)
	if err != nil {
		log.Println(err)
	}

	return nil
}
