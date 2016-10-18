/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"os"
	"runtime"

	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"
)

var natsClient *nats.Conn
var em = eventManager{}
var p = storage{}
var cfg *ecc.Config

// Receives a message, updates the related service on the FSM
// and emits the relative message
func manageInputMessage(m *nats.Msg) {
	var service map[string]interface{}
	mm := MessageManager{}

	service, subject, err := mm.getServiceFromMessage(m.Subject, m.Data)
	if err == nil {
		subject, err := em.manage(subject, &service)
		if err := SaveService(&service); err != nil {
			log.Println("[ERROR] : " + err.Error())
		}
		message, err := mm.preparePublishMessage(subject, &service)

		if err != nil {
			log.Println(err)
		} else {
			em.move(&service, subject)
			log.Println("[PROCESSED]", m.Subject)
			if err := SaveService(&service); err != nil {
				log.Println("[ERROR] : " + err.Error())
			}
			natsClient.Publish(subject, []byte(message))
			log.Println("[EMITTED]", subject)
		}
	}
	service = nil
}

// Setup the listeners for all messages on the platform
func main() {
	cfg = ecc.NewConfig(os.Getenv("NATS_URI"))
	natsClient = cfg.Nats()
	p.load(natsClient)

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
		mm := MessageManager{}
		s, err := mm.getService(m.Data)
		if err != nil {
			log.Println("Service not found")
		} else {
			ServiceDel(&s)
		}
	})

	runtime.Goexit()
}
