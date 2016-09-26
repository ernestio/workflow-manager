/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// NetworksCreate : Message to create networks
type NetworksCreate struct {
	Service              string    `json:"service"`
	Networks             []network `json:"components"`
	Status               string    `json:"status"`
	ErrorCode            string    `json:"error_code"`
	ErrorMessage         string    `json:"error_message"`
	SequentialProcessing bool      `json:"sequential_processing"`
}

// GenericComponentCreate : Message to create instances
type GenericComponentMsg struct {
	Service              string        `json:"service"`
	Components           []interface{} `json:"components"`
	Status               string        `json:"status"`
	ErrorCode            string        `json:"error_code"`
	ErrorMessage         string        `json:"error_message"`
	SequentialProcessing bool          `json:"sequential_processing"`
}

// InstancesCreate : Message to create instances
type InstancesCreate struct {
	Service              string     `json:"service"`
	Instances            []instance `json:"components"`
	Status               string     `json:"status"`
	ErrorCode            string     `json:"error_code"`
	ErrorMessage         string     `json:"error_message"`
	SequentialProcessing bool       `json:"sequential_processing"`
}

// FirewallsCreate : Message to create firewalls
type FirewallsCreate struct {
	Service              string     `json:"service"`
	Firewalls            []firewall `json:"components"`
	Networks             []network  `json:"networks"`
	Status               string     `json:"status"`
	ErrorCode            string     `json:"error_code"`
	ErrorMessage         string     `json:"error_message"`
	SequentialProcessing bool       `json:"sequential_processing"`
}

// NatsCreate : Message to create nats
type NatsCreate struct {
	Service              string `json:"service"`
	Nats                 []nat  `json:"components"`
	Status               string `json:"status"`
	ErrorCode            string `json:"error_code"`
	ErrorMessage         string `json:"error_message"`
	SequentialProcessing bool   `json:"sequential_processing"`
}

// ExecutionsCreate : Message to create Executions
type ExecutionsCreate struct {
	Service     string      `json:"service"`
	ServiceName string      `json:"service_name"`
	ServiceType string      `json:"service_type"`
	Executions  []execution `json:"components"`
}

// ServiceOptions : Service options aka salt user password
type ServiceOptions struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// Monitor : Messages to be sent to monitoring service
type Monitor struct {
	Service  string           `json:"service"`
	Messages []MonitorMessage `json:"messages"`
}

// MonitorMessages : THe message to be sent
type MonitorMessage struct {
	Body  string `json:"body"`
	Level string `json:"level"`
}
