/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
)

// GenericComponentMsg : Message to create instances
type GenericComponentMsg struct {
	Service              string            `json:"service"`
	Components           []interface{}     `json:"components"`
	AWSAccessKeyID       string            `json:"aws_access_key_id"`
	AWSSecretAccessKey   string            `json:"aws_secret_access_key"`
	DatacenterRegion     string            `json:"datacenter_region"`
	Tags                 map[string]string `json:"tags"`
	Status               string            `json:"status"`
	ErrorCode            string            `json:"error_code"`
	ErrorMessage         string            `json:"error"`
	SequentialProcessing bool              `json:"sequential_processing"`
	Type                 string            `json:"_type"`
}

// NewGenericComponentMsg : GenericComponentMsg constructor
func NewGenericComponentMsg(body []byte) GenericComponentMsg {
	input := GenericComponentMsg{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Panic(err.Error())
	}

	return input
}
