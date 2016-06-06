/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"os"

	"github.com/nats-io/nats"
)

// Config : struct representation of service configuration
type Config struct {
	Nats struct {
		URL string `json:"url"`
	} `json:"nats"`
	SaltAuthentication struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}
}

// Load : will load configuration from the given path
func (c *Config) Load() {
	c.Nats.URL = os.Getenv("NATS_URI")
	c.SaltAuthentication.User = ""
	c.SaltAuthentication.Password = ""
}

// NatsClient : creates a new nats client
func (c *Config) NatsClient() *nats.Conn {
	n, err := nats.Connect(c.Nats.URL)
	if err != nil {
		log.Println("Could not connect to NATS server")
		panic(err)
	}

	return n
}
