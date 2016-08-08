/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

func buildCreateNats(s *service) NatsCreate {
	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Configuring nats", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return buildNatsList(s, s.NatsToCreate.Items)
}

func buildUpdateNats(s *service) NatsCreate {
	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Configuring nats", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return buildNatsList(s, s.NatsToUpdate.Items)
}

func buildDeleteNats(s *service) NatsCreate {
	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Deleting nats", Level: "INFO"})
	UserOutput(s.Channel(), messages)

	return buildNatsList(s, s.NatsToDelete.Items)
}

// Creates a CreateNetworks struct based on a given service
func buildNatsList(s *service, inputList []nat) NatsCreate {
	list := make([]nat, len(inputList))
	copy(list, inputList)

	d := s.datacenter()

	m := NatsCreate{
		Service: s.ID,
		Nats:    list,
	}

	r := &router{}

	for i, n := range list {
		var endpoint string

		r = s.routerByName(n.RouterName)
		net := s.networkByName(n.NetworkName)

		if s.ServiceIP != "" {
			endpoint = s.ServiceIP
		} else if r != nil {
			endpoint = r.IP
		}

		m.Nats[i] = nat{
			Service:            s.ID,
			Name:               n.Name,
			ClientName:         s.ClientName,
			DatacenterName:     d.Name,
			DatacenterPassword: d.Password,
			DatacenterRegion:   d.Region,
			DatacenterType:     d.Type,
			DatacenterUsername: d.Username,
			DatacenterToken:    d.Token,
			DatacenterSecret:   d.Secret,
			ExternalNetwork:    d.ExternalNetwork,
			VCloudURL:          d.VCloudURL,
		}

		m.Nats[i].RouterName = r.Name
		m.Nats[i].RouterType = r.Type
		m.Nats[i].RouterIP = r.IP

		if net != nil {
			m.Nats[i].NetworkAWSID = net.NetworkAWSID
		}

		m.Nats[i].Status = n.Status
		m.Nats[i].Rules = n.Rules
		for x := 0; x < len(n.Rules); x++ {
			m.Nats[i].Rules[x].Type = n.Rules[x].Type
			m.Nats[i].Rules[x].Protocol = n.Rules[x].Protocol
			m.Nats[i].Rules[x].Network = n.Rules[x].Network
			origin := n.Rules[x].OriginIP
			if origin == "" {
				m.Nats[i].Rules[x].OriginIP = endpoint
			} else {
				m.Nats[i].Rules[x].OriginIP = origin
			}
			m.Nats[i].Rules[x].OriginPort = n.Rules[x].OriginPort
			translation := n.Rules[x].TranslationIP
			if translation == "" {
				m.Nats[i].Rules[x].TranslationIP = endpoint
			} else {
				m.Nats[i].Rules[x].TranslationIP = translation
			}
			m.Nats[i].Rules[x].TranslationPort = n.Rules[x].TranslationPort
		}
	}

	return m
}
