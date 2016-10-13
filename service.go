/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
)

type status struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type datacenter struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	Region          string `json:"region"`
	Type            string `json:"type"`
	Username        string `json:"username"`
	ExternalNetwork string `json:"external_network"`
	Token           string `json:"token"`
	Secret          string `json:"secret"`
	VCloudURL       string `json:"vcloud_url"`
	VseURL          string `json:"vse_url"`
	status
}

type vpc struct {
	DatacenterName   string `json:"datacenter_name"`
	DatacenterRegion string `json:"datacenter_region"`
	DatacenterToken  string `json:"datacenter_access_token"`
	DatacenterSecret string `json:"datacenter_access_key"`
	VpcID            string `json:"vpc_id"`
	VpcSubnet        string `json:"vpc_subnet"`
	Type             string `json:"_type"`
	status
}

type executionReport struct {
	Instance   string `json:"instance"`
	ReturnCode int    `json:"return_code"`
	Stderr     string `json:"stderr"`
	Stdout     string `json:"stdout"`
}

type execution struct {
	Name               string            `json:"name"`
	Type               string            `json:"type"`
	Service            string            `json:"service"`
	Payload            string            `json:"payload"`
	ClientName         string            `json:"client_name,omitempty"`
	DatacenterType     string            `json:"datacenter_type,omitempty"`
	DatacenterName     string            `json:"datacenter_name,omitempty"`
	DatacenterUsername string            `json:"datacenter_username,omitempty"`
	DatacenterPassword string            `json:"datacenter_password,omitempty"`
	DatacenterRegion   string            `json:"datacenter_region,omitempty"`
	Target             string            `json:"target"`
	MatchedInstances   []string          `json:"matched_instances"`
	Reports            []executionReport `json:"reports"`
	ExecutionStatus    string            `json:"execution_status"`
	Created            bool              `json:"created"`
	User               string            `json:"user"`
	Password           string            `json:"password"`
	EndPoint           string            `json:"service_endpoint"`
	status
}

type firewallRules struct {
	Type            string `json:"type"`
	Destination     string `json:"destination_ip"`
	DestinationPort string `json:"destination_port"`
	Protocol        string `json:"protocol"`
	Source          string `json:"source_ip"`
	SourcePort      string `json:"source_port"`
}

type firewall struct {
	Type               string          `json:"type"`
	Name               string          `json:"name"`
	Rules              []firewallRules `json:"rules"`
	FirewallType       string          `json:"firewall_type"`
	Service            string          `json:"service"`
	ClientName         string          `json:"client_name"`
	RouterName         string          `json:"router_name"`
	RouterType         string          `json:"router_type"`
	RouterIP           string          `json:"router_ip"`
	DatacenterName     string          `json:"datacenter_name"`
	DatacenterPassword string          `json:"datacenter_password"`
	DatacenterRegion   string          `json:"datacenter_region"`
	DatacenterType     string          `json:"datacenter_type"`
	DatacenterUsername string          `json:"datacenter_username"`
	DatacenterToken    string          `json:"datacenter_token"`
	DatacenterSecret   string          `json:"datacenter_secret"`
	ExternalNetwork    string          `json:"external_network"`
	SecurityGroupAWSID string          `json:"security_group_aws_id"`
	VCloudURL          string          `json:"vcloud_url"`
	VpcID              string          `json:"vpc_id"`
	status
}

type instanceDisk struct {
	ID   int `json:"id"`
	Size int `json:"size"`
}

type instance struct {
	Service             string         `json:"service"`
	Name                string         `json:"name"`
	Type                string         `json:"type"`
	IP                  string         `json:"ip"`
	CPU                 int            `json:"cpus"`
	RAM                 int            `json:"ram"`
	Catalog             string         `json:"reference_catalog"`
	Image               string         `json:"reference_image"`
	Disks               []instanceDisk `json:"disks"`
	PublicIP            string         `json:"public_ip"`
	InstanceAWSID       string         `json:"instance_aws_id"`
	RouterName          string         `json:"router_name"`
	RouterType          string         `json:"router_type"`
	RouterIP            string         `json:"router_ip"`
	ClientName          string         `json:"client_name"`
	DatacenterName      string         `json:"datacenter_name"`
	DatacenterPassword  string         `json:"datacenter_password"`
	DatacenterRegion    string         `json:"datacenter_region"`
	DatacenterType      string         `json:"datacenter_type"`
	DatacenterUsername  string         `json:"datacenter_username"`
	DatacenterToken     string         `json:"datacenter_token"`
	DatacenterSecret    string         `json:"datacenter_secret"`
	NetworkName         string         `json:"network_name"`
	NetworkIsPublic     bool           `json:"network_is_public"`
	NetworkAWSID        string         `json:"network_aws_id"`
	KeyPair             string         `json:"key_pair"`
	AssignElasticIP     bool           `json:"assign_elastic_ip"`
	SecurityGroups      []string       `json:"security_groups"`
	SecurityGroupAWSIDs []string       `json:"security_group_aws_ids"`
	VCloudURL           string         `json:"vcloud_url"`
	VpcID               string         `json:"vpc_id"`
	status
}

type loadbalancer struct {
	Instance string `json:"instance"`
	Name     string `json:"name"`
	Ports    []struct {
		From int `json:"from"`
		To   int `json:"to"`
	} `json:"ports"`
	Router  string `json:"router"`
	Service string `json:"service"`
	status
}

type natRule struct {
	Network         string `json:"network"`
	OriginIP        string `json:"origin_ip"`
	OriginPort      string `json:"origin_port"`
	Protocol        string `json:"protocol"`
	TranslationIP   string `json:"translation_ip"`
	TranslationPort string `json:"translation_port"`
	Type            string `json:"type"`
}

type nat struct {
	Service                string    `json:"service"`
	Name                   string    `json:"name"`
	Rules                  []natRule `json:"rules"`
	NatType                string    `json:"nat_type"`
	NetworkName            string    `json:"network_name"`
	PublicNetwork          string    `json:"public_network"`
	RoutedNetworks         []string  `json:"routed_networks"`
	RoutedNetworkAWSIDs    []string  `json:"routed_networks_aws_ids"`
	PublicNetworkAWSID     string    `json:"public_network_aws_id"`
	NatGatewayAWSID        string    `json:"nat_gateway_aws_id"`
	NatGatewayAllocationID string    `json:"nat_gateway_allocation_id"`
	NatGatewayAllocationIP string    `json:"nat_gateway_allocation_ip"`
	RouterName             string    `json:"router_name"`
	RouterType             string    `json:"router_type"`
	RouterIP               string    `json:"router_ip"`
	ClientID               string    `json:"client_id"`
	ClientName             string    `json:"client_name"`
	DatacenterType         string    `json:"datacenter_type"`
	DatacenterName         string    `json:"datacenter_name"`
	DatacenterUsername     string    `json:"datacenter_username"`
	DatacenterPassword     string    `json:"datacenter_password"`
	DatacenterRegion       string    `json:"datacenter_region,omitempty"`
	DatacenterToken        string    `json:"datacenter_token"`
	DatacenterSecret       string    `json:"datacenter_secret"`
	ExternalNetwork        string    `json:"external_network"`
	VCloudURL              string    `json:"vcloud_url"`
	VpcID                  string    `json:"vpc_id"`
	status
}

type network struct {
	Name               string   `json:"name"`
	Type               string   `json:"type"`
	Service            string   `json:"service"`
	Range              string   `json:"range"`
	Subnet             string   `json:"subnet"`
	Netmask            string   `json:"netmask"`
	StartAddress       string   `json:"start_address"`
	EndAddress         string   `json:"end_address"`
	Gateway            string   `json:"gateway"`
	IsPublic           bool     `json:"is_public"`
	RouterName         string   `json:"router_name"`
	RouterType         string   `json:"router_type"`
	RouterIP           string   `json:"router_ip"`
	ClientName         string   `json:"client_name"`
	DatacenterType     string   `json:"datacenter_type"`
	DatacenterName     string   `json:"datacenter_name"`
	DatacenterUsername string   `json:"datacenter_username"`
	DatacenterPassword string   `json:"datacenter_password"`
	DatacenterRegion   string   `json:"datacenter_region"`
	DatacenterToken    string   `json:"datacenter_token"`
	DatacenterSecret   string   `json:"datacenter_secret"`
	NetworkType        string   `json:"network_type"`
	NetworkSubnet      string   `json:"network_subnet"`
	NetworkAWSID       string   `json:"network_aws_id"`
	VpcID              string   `json:"vpc_id"`
	DNS                []string `json:"DNS"`
	VCloudURL          string   `json:"vcloud_url"`
	status
}

type router struct {
	Service            string `json:"service"`
	Type               string `json:"type"`
	IP                 string `json:"ip"`
	Name               string `json:"name"`
	ClientName         string `json:"client_name"`
	DatacenterName     string `json:"datacenter_name"`
	DatacenterPassword string `json:"datacenter_password"`
	DatacenterRegion   string `json:"datacenter_region"`
	DatacenterType     string `json:"datacenter_type"`
	DatacenterUsername string `json:"datacenter_username"`
	ExternalNetwork    string `json:"external_network"`
	VCloudURL          string `json:"vcloud_url"`
	VseURL             string `json:"vse_url"`
	Created            bool   `json:"created"`
	status
}

// elbListener ...
type elbListener struct {
	FromPort int    `json:"from_port"`
	ToPort   int    `json:"to_port"`
	Protocol string `json:"protocol"`
	SSLCert  string `json:"ssl_cert"`
}

// ELB : Mapping for a elb component
type elb struct {
	Type                string        `json:"_type"`
	Name                string        `json:"name"`
	IsPrivate           bool          `json:"is_private"`
	DNSName             string        `json:"dns_name"`
	Listeners           []elbListener `json:"listeners"`
	NetworkAWSIDs       []string      `json:"network_aws_ids"`
	Instances           []string      `json:"instances"`
	InstanceNames       []string      `json:"instance_names"`
	InstanceAWSIDs      []string      `json:"instance_aws_ids"`
	SecurityGroups      []string      `json:"security_groups"`
	SecurityGroupAWSIDs []string      `json:"security_group_aws_ids"`
	DatacenterType      string        `json:"datacenter_type,omitempty"`
	DatacenterName      string        `json:"datacenter_name,omitempty"`
	DatacenterRegion    string        `json:"datacenter_region"`
	DatacenterToken     string        `json:"datacenter_token"`
	DatacenterSecret    string        `json:"datacenter_secret"`
	VpcID               string        `json:"vpc_id"`
	Service             string        `json:"service"`
	Status              string        `json:"status"`
	Exists              bool
}

// This is the object representation for a service inside the
// FSM, it has appended the workflow the service needs to
// follow to be built
type service struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Body          string   `json:"body"`
	Workflow      workflow `json:"workflow"`
	Started       string   `json:"started"`
	Finished      string   `json:"finished"`
	Status        string   `json:"status"`
	Type          string   `json:"type"`
	Endpoint      string   `json:"endpoint"`
	ClientName    string   `json:"client_name"`
	Parent        string   `json:"-"`
	Bootstrapping string   `json:"bootstrapping"`
	ErnestIP      []string `json:"ernest_ip"`
	ServiceIP     string   `json:"service_ip"`
	Options       struct {
		Password string `json:"password"`
		User     string `json:"user"`
	} `json:"options"`
	Datacenters struct {
		Finished             string       `json:"finished"`
		Items                []datacenter `json:"items"`
		Started              string       `json:"started"`
		Status               string       `json:"status"`
		SequentialProcessing bool         `json:"sequential_processing"`
	} `json:"datacenters"`
	VPCs struct {
		Finished             string `json:"finished"`
		Items                []vpc  `json:"items"`
		Started              string `json:"started"`
		Status               string `json:"status"`
		SequentialProcessing bool   `json:"sequential_processing"`
	} `json:"vpcs"`
	VPCsToCreate struct {
		Finished             string `json:"finished"`
		Items                []vpc  `json:"items"`
		Started              string `json:"started"`
		Status               string `json:"status"`
		SequentialProcessing bool   `json:"sequential_processing"`
	} `json:"vpcs_to_create"`
	VPCsToDelete struct {
		Finished             string `json:"finished"`
		Items                []vpc  `json:"items"`
		Started              string `json:"started"`
		Status               string `json:"status"`
		SequentialProcessing bool   `json:"sequential_processing"`
	} `json:"vpcs_to_delete"`
	Bootstraps struct {
		Finished             string      `json:"finished"`
		Items                []execution `json:"items"`
		Started              string      `json:"started"`
		Status               string      `json:"status"`
		SequentialProcessing bool        `json:"sequential_processing"`
	} `json:"bootstraps"`
	BootstrapsToCreate struct {
		Finished             string      `json:"finished"`
		Items                []execution `json:"items"`
		Started              string      `json:"started"`
		Status               string      `json:"status"`
		SequentialProcessing bool        `json:"sequential_processing"`
	} `json:"bootstraps_to_create"`
	Executions struct {
		Finished             string      `json:"finished"`
		Items                []execution `json:"items"`
		Started              string      `json:"started"`
		Status               string      `json:"status"`
		SequentialProcessing bool        `json:"sequential_processing"`
	} `json:"executions"`
	ExecutionsToCreate struct {
		Finished             string      `json:"finished"`
		Items                []execution `json:"items"`
		Started              string      `json:"started"`
		Status               string      `json:"status"`
		SequentialProcessing bool        `json:"sequential_processing"`
	} `json:"executions_to_create"`
	Firewalls struct {
		Finished string     `json:"finished"`
		Items    []firewall `json:"items"`
		Started  string     `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"firewalls"`
	FirewallsToCreate struct {
		Finished string     `json:"finished"`
		Items    []firewall `json:"items"`
		Started  string     `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"firewalls_to_create"`
	FirewallsToUpdate struct {
		Finished string     `json:"finished"`
		Items    []firewall `json:"items"`
		Started  string     `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"firewalls_to_update"`
	FirewallsToDelete struct {
		Finished string     `json:"finished"`
		Items    []firewall `json:"items"`
		Started  string     `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"firewalls_to_delete"`
	Instances struct {
		Finished string     `json:"finished"`
		Items    []instance `json:"items"`
		Started  string     `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"instances"`
	InstancesToCreate struct {
		Finished string     `json:"finished"`
		Items    []instance `json:"items"`
		Started  string     `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"instances_to_create"`
	InstancesToDelete struct {
		Finished string     `json:"finished"`
		Items    []instance `json:"items"`
		Started  string     `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"instances_to_delete"`
	InstancesToUpdate struct {
		Finished             string     `json:"finished"`
		Items                []instance `json:"items"`
		Started              string     `json:"started"`
		Status               string     `json:"status"`
		SequentialProcessing bool       `json:"sequential_processing"`
	} `json:"instances_to_update"`
	Loadbalancer struct {
		Finished string         `json:"finished"`
		Items    []loadbalancer `json:"items"`
		Started  string         `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"loadbalancer"`
	Nats struct {
		Finished string `json:"finished"`
		Items    []nat  `json:"items"`
		Started  string `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"nats"`
	NatsToCreate struct {
		Finished string `json:"finished"`
		Items    []nat  `json:"items"`
		Started  string `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"nats_to_create"`
	NatsToUpdate struct {
		Finished string `json:"finished"`
		Items    []nat  `json:"items"`
		Started  string `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"nats_to_update"`
	NatsToDelete struct {
		Finished string `json:"finished"`
		Items    []nat  `json:"items"`
		Started  string `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"nats_to_delete"`
	Networks struct {
		Finished string    `json:"finished"`
		Items    []network `json:"items"`
		Started  string    `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"networks"`
	NetworksToCreate struct {
		Finished string    `json:"finished"`
		Items    []network `json:"items"`
		Started  string    `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"networks_to_create"`
	NetworksToDelete struct {
		Finished string    `json:"finished"`
		Items    []network `json:"items"`
		Started  string    `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"networks_to_delete"`
	Routers struct {
		Finished string   `json:"finished"`
		Items    []router `json:"items"`
		Started  string   `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"routers"`
	RoutersToCreate struct {
		Finished string   `json:"finished"`
		Items    []router `json:"items"`
		Started  string   `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"routers_to_create"`
	RoutersToDelete struct {
		Finished string   `json:"finished"`
		Items    []router `json:"items"`
		Started  string   `json:"started"`
		status
		SequentialProcessing bool `json:"sequential_processing"`
	} `json:"routers_to_delete"`
	ELBs struct {
		Started              string `json:"started"`
		Finished             string `json:"finished"`
		Status               string `json:"status"`
		Items                []elb  `json:"items"`
		SequentialProcessing bool   `json:"sequential_processing"`
	} `json:"elbs"`
	ELBsToCreate struct {
		Started              string `json:"started"`
		Finished             string `json:"finished"`
		Status               string `json:"status"`
		Items                []elb  `json:"items"`
		SequentialProcessing bool   `json:"sequential_processing"`
	} `json:"elbs_to_create"`
	ELBsToUpdate struct {
		Started              string `json:"started"`
		Finished             string `json:"finished"`
		Status               string `json:"status"`
		Items                []elb  `json:"items"`
		SequentialProcessing bool   `json:"sequential_processing"`
	} `json:"elbs_to_update"`
	ELBsToDelete struct {
		Started              string `json:"started"`
		Finished             string `json:"finished"`
		Status               string `json:"status"`
		Items                []elb  `json:"items"`
		SequentialProcessing bool   `json:"sequential_processing"`
	} `json:"elbs_to_delete"`
}

type message struct {
	Service string `json:"service"`
}

// Validates if a service is valid or not
func (s *service) valid() bool {
	// TODO : Check if a service is valid or nt
	return true
}

func (s *service) markAsFailed() {
	s.Status = "pre-failed"
}

func (s *service) toJSON() string {
	json, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}
	return string(json)
}

// Persist current service
func (s *service) save() error {
	json, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
		return err
	}
	err = p.set(s.ID, string(json))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *service) del() {
	p.del(s.ID)
}

func (s *service) datacenter() datacenter {
	return s.Datacenters.Items[0]
}

func (s *service) routerByName(name string) *router {
	for _, r := range s.Routers.Items {
		if r.Name == name {
			return &r
		}
	}

	return nil
}

func (s *service) networkByName(name string) *network {
	for _, n := range s.Networks.Items {
		if n.Name == name {
			return &n
		}
	}

	return nil
}

func (s *service) firewallByName(name string) *firewall {
	for _, f := range s.Firewalls.Items {
		if f.Name == name {
			return &f
		}
	}

	return nil
}

func (s *service) executionByName(name string) *execution {
	for i, e := range s.Executions.Items {
		if e.Name == name {
			return &s.Executions.Items[i]
		}
	}

	return nil
}

func (s *service) saltMaster() *instance {
	for _, instance := range s.Instances.Items {
		if instance.Type != "salt" {
			return &instance
		}
	}
	return nil
}

func (s *service) Channel() string {
	return s.ID
}

func (s *service) asMap() (mapped map[string]interface{}) {
	body, err := json.Marshal(s)
	if err != nil {
		log.Panic(err.Error())
	}
	err = json.Unmarshal(body, &mapped)
	if err != nil {
		log.Panic(err.Error())
	}
	return mapped
}

func (s *service) loadFromMap(mapped map[string]interface{}) {
	body, err := json.Marshal(mapped)
	if err != nil {
		log.Panic(err.Error())
	}
	err = json.Unmarshal(body, s)
}

func (s *service) getComponentList(cType string) []interface{} {
	tmp := s.asMap()
	cList := tmp[cType].(map[string]interface{})
	list := cList["items"].([]interface{})

	return list
}

func (s *service) transferCreated(cType string, input GenericComponentMsg) {
	var components []interface{}
	var erroredComponents []interface{}

	tmp := s.asMap()
	inputComponents := input.Components
	currentComponents := tmp[cType].(map[string]interface{})
	if currentComponents["items"] != nil {
		components = currentComponents["items"].([]interface{})
	}

	// Append new components
	for _, c := range inputComponents {
		inHash := c.(map[string]interface{})
		status := inHash["status"].(string)
		if status == "errored" {
			erroredComponents = append(erroredComponents, c)
		} else {
			components = append(components, c)
		}
	}
	currentComponents["status"] = "completed"
	currentComponents["items"] = components

	// Remove to be created components
	componentsToBeProcessed := tmp[cType+"_to_create"].(map[string]interface{})
	componentsToBeProcessed["items"] = erroredComponents
	componentsToBeProcessed["status"] = input.Status
	componentsToBeProcessed["error_code"] = input.ErrorCode
	componentsToBeProcessed["error_message"] = input.ErrorMessage

	// Save result
	s.loadFromMap(tmp)
}

func (s *service) transferUpdated(cType string, input GenericComponentMsg) {
	var components []interface{}
	var erroredComponents []interface{}

	tmp := s.asMap()
	inputComponents := input.Components

	currentComponents := tmp[cType].(map[string]interface{})
	if currentComponents["items"] != nil {
		components = currentComponents["items"].([]interface{})
	}

	// Append new components
	for _, c := range inputComponents {
		for i, v := range components {
			inHash := c.(map[string]interface{})
			exHash := v.(map[string]interface{})
			iName := inHash["name"].(string)
			name := exHash["name"].(string)
			if iName == name {
				status := inHash["status"].(string)
				if status == "completed" {
					components[i] = c
				} else {
					erroredComponents = append(erroredComponents, c)
				}
			}
		}
	}
	currentComponents["status"] = "completed"
	currentComponents["items"] = components

	// Remove to be created components
	componentsToBeProcessed := tmp[cType+"_to_update"].(map[string]interface{})

	componentsToBeProcessed["items"] = erroredComponents
	componentsToBeProcessed["status"] = input.Status
	componentsToBeProcessed["error_code"] = input.ErrorCode
	componentsToBeProcessed["error_message"] = input.ErrorMessage

	// Save result
	s.loadFromMap(tmp)
}

func (s *service) transferDeleted(cType string, input GenericComponentMsg) {
	var components []interface{}
	var remanentComponents []interface{}
	var erroredComponents []interface{}

	tmp := s.asMap()
	inputComponents := input.Components

	currentComponents := tmp[cType].(map[string]interface{})
	if currentComponents["items"] != nil {
		components = currentComponents["items"].([]interface{})
	}

	for _, v := range components {
		sw := false
		exHash := v.(map[string]interface{})
		name := exHash["name"].(string)
		for _, c := range inputComponents {
			inHash := c.(map[string]interface{})
			iName := inHash["name"].(string)
			if iName == name {
				status := inHash["status"].(string)
				if status == "errored" {
					erroredComponents = append(erroredComponents, c)
				} else {
					sw = true
				}
			}
		}
		if sw == false {
			remanentComponents = append(remanentComponents, v)
		}
	}
	currentComponents["status"] = "completed"
	currentComponents["items"] = remanentComponents

	// Remove to be created components
	componentsToBeProcessed := tmp[cType+"_to_delete"].(map[string]interface{})

	componentsToBeProcessed["items"] = erroredComponents
	if len(erroredComponents) > 0 {
		componentsToBeProcessed["status"] = "errored"
	} else {
		componentsToBeProcessed["status"] = "completed"
	}
	componentsToBeProcessed["error_code"] = input.ErrorCode
	componentsToBeProcessed["error_message"] = input.ErrorMessage

	// Save result
	s.loadFromMap(tmp)
}
