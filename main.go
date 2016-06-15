/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"runtime"
	"time"

	"github.com/nats-io/nats"
)

var natsClient *nats.Conn
var em = eventManager{}
var p = storage{}
var c = Config{}

// Receives a message, updates the related service on the FSM
// and emits the relative message
func manageInputMessage(m *nats.Msg) {
	mm := messageManager{}

	log.Println("[RECEIVED]", m.Subject)
	s, subject, err := mm.getServiceFromMessage(m.Subject, m.Data)
	if err == nil {
		subject, s, err := em.manage(subject, s)
		s.save()
		message, err := mm.preparePublishMessage(subject, s)

		if err != nil {
			log.Println(err)
		} else {
			em.move(s, subject)
			s.save()
			natsClient.Publish(subject, []byte(message))
			log.Println("[EMITTED]", subject)
		}
	}
}

// Setup the listeners for all messages on the platform
func main() {
	c.Load()
	natsClient = c.NatsClient()
	p.load(natsClient)

	saltCfg, err := natsClient.Request("config.get.salt", []byte(""), 1*time.Second)
	if err == nil {
		json.Unmarshal(saltCfg.Data, &c.SaltAuthentication)
	}

	// Messages matching *.* are always actions
	natsClient.Subscribe("*.*", func(m *nats.Msg) {
		manageInputMessage(m)
	})

	// Messages with *.*.* are results
	natsClient.Subscribe("*.*.*", func(m *nats.Msg) {
		manageInputMessage(m)
	})

	// Service delete
	natsClient.Subscribe("service.delete.done", func(m *nats.Msg) {
		mm := messageManager{}
		s, err := mm.getService(m.Data)
		if err != nil {
			log.Println("Service not found")
		} else {
			s.del()
		}
	})

	runtime.Goexit()
}
