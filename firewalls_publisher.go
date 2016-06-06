/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// Creates a CreateFirewalls struct based on a given service
func buildCreateFirewalls(s *service) FirewallsCreate {
	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Setting up firewalls:", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return buildFirewallsList(s, s.Firewalls.Items)
}

func buildUpdateFirewalls(s *service) FirewallsCreate {
	return buildCreateFirewalls(s)
}

func buildDeleteFirewalls(s *service) FirewallsCreate {
	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Deleting firewalls:", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return buildFirewallsList(s, s.FirewallsToDelete.Items)
}

func buildFirewallsList(s *service, inputList []firewall) FirewallsCreate {
	list := make([]firewall, len(inputList))
	copy(list, inputList)

	d := s.datacenter()

	m := FirewallsCreate{
		Service:   s.ID,
		Firewalls: list,
		Networks:  s.Networks.Items,
	}

	r := &router{}
	for i, f := range list {
		r = s.routerByName(f.RouterName)
		rules := make([]firewallRules, len(f.Rules))

		endpoint := r.IP
		if s.ServiceIP != "" {
			endpoint = s.ServiceIP
		}

		for j, rule := range f.Rules {
			destination := rule.Destination
			source := rule.Source
			if network := s.networkByName(rule.Destination); network != nil {
				destination = network.Range
			}
			if network := s.networkByName(rule.Source); network != nil {
				source = network.Range
			}
			if destination == "" {
				destination = endpoint
			}

			rules[j] = firewallRules{
				Destination:     destination,
				DestinationPort: rule.DestinationPort,
				Protocol:        rule.Protocol,
				Source:          source,
				SourcePort:      rule.SourcePort,
			}
		}

		m.Firewalls[i] = firewall{
			Name:               f.Name,
			Rules:              rules,
			RouterName:         r.Name,
			RouterType:         r.Type,
			RouterIP:           r.IP,
			ClientName:         s.ClientName,
			DatacenterName:     d.Name,
			DatacenterPassword: d.Password,
			DatacenterRegion:   d.Region,
			DatacenterType:     d.Type,
			DatacenterUsername: d.Username,
			ExternalNetwork:    d.ExternalNetwork,
			VCloudURL:          d.VCloudURL,
		}
		m.Firewalls[i].Status = f.Status
	}

	return m
}
