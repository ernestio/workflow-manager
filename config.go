/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// Config : struct representation of service configuration
type Config struct {
	SaltAuthentication struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}
}

// Load : will load configuration from the given path
func (c *Config) Load() {
	c.SaltAuthentication.User = ""
	c.SaltAuthentication.Password = ""
}
