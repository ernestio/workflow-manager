/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"
)

type testHelper struct{}

func (t *testHelper) getService(source string) (*map[string]interface{}, string) {
	var s map[string]interface{}

	absPath, _ := filepath.Abs(source)
	file, err := os.Open(absPath)
	log.Printf("Reading config from: %s", source)
	if err != nil {
		log.Panic("error:", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s)
	if err != nil {
		log.Println("Definition file is invalid")
		log.Panic("error:", err)
	}
	content, _ := ioutil.ReadFile(source)

	return &s, string(content)
}

func (t *testHelper) getServiceBody(source string) string {
	absPath, _ := filepath.Abs(source)
	content, _ := ioutil.ReadFile(absPath)

	return string(content)
}

func (t *testHelper) manage(subject string, s *map[string]interface{}) (string, error) {
	em := eventManager{}
	return em.manage(subject, s)
}

func (t *testHelper) getFixture(source string) []byte {
	absPath, _ := filepath.Abs(source)
	content, _ := ioutil.ReadFile(absPath)

	return []byte(content)
}

var store = make(map[string]string)
var listeners = false

func runListenerMocks() {
	natsClient.Subscribe("service.get.mapping", func(m *nats.Msg) {
		sm := serviceMessage{}
		json.Unmarshal(m.Data, &sm)
		natsClient.Publish(m.Reply, []byte(store[sm.ID]))
	})

	natsClient.Subscribe("service.set.mapping", func(m *nats.Msg) {
		sm := serviceMessage{}
		err := json.Unmarshal(m.Data, &sm)
		if err != nil {
			println(err.Error())
		}
		store[sm.ID] = sm.Mapping
		manageInputMessage(m)
		natsClient.Publish(m.Reply, []byte(store[sm.ID]))
	})

	natsClient.Subscribe("service.del.mapping", func(m *nats.Msg) {
		sm := serviceMessage{}
		json.Unmarshal(m.Data, &sm)
	})
}

func setup() {
	if listeners == false {
		natsClient = ecc.NewConfig(os.Getenv("NATS_URI")).Nats()
		runListenerMocks()
		listeners = true
	}
}
