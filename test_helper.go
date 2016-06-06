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
)

type testHelper struct{}

var redisCfg = []byte(`{"addr":"localhost:6379","password":"","DB":0}`)

func (t *testHelper) getService(source string) service {
	s := service{}
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

	return s
}

func (t *testHelper) getServiceBody(source string) string {
	absPath, _ := filepath.Abs(source)
	content, _ := ioutil.ReadFile(absPath)

	return string(content)
}

func (t *testHelper) manage(subject string, s service) (string, *service, error) {
	em := eventManager{}
	return em.manage(subject, &s)
}

func (t *testHelper) getFixture(source string) []byte {
	absPath, _ := filepath.Abs(source)
	content, _ := ioutil.ReadFile(absPath)

	return []byte(content)
}
