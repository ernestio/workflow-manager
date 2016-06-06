/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
)

// UserOutput ...
func UserOutput(service string, messages []MonitorMessage) {
	m := Monitor{
		Service:  service,
		Messages: messages,
	}

	body, err := json.Marshal(m)
	if err != nil {
		log.Println("Can't send logs to the final user")
	} else {
		natsClient.Publish("monitor.user", body)
	}
}
