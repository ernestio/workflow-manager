/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// Creates a CreateNetworks struct based on a given service
func buildCreateNetworks(s *service) NetworksCreate {
	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Creating networks:", Level: "INFO"})

	res := NetworksCreate{}

	if len(s.NetworksToCreate.Items) > 0 {
		res, messages = buildNetworksList(s, s.NetworksToCreate.Items, messages)
	} else {
		res, messages = buildNetworksList(s, s.Networks.Items, messages)
	}

	UserOutput(s.Channel(), messages)

	return res
}

// Creates a DeleteNetworks struct based on a given service
func buildDeleteNetworks(s *service) NetworksCreate {
	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Deleting networks:", Level: "INFO"})

	res, messages := buildNetworksList(s, s.NetworksToDelete.Items, messages)

	UserOutput(s.Channel(), messages)

	return res
}

func buildNetworksList(s *service, list []network, messages []MonitorMessage) (NetworksCreate, []MonitorMessage) {
	d := s.datacenter()

	m := NetworksCreate{
		Service:              s.ID,
		Networks:             list,
		SequentialProcessing: true,
	}

	r := &router{}
	for i, n := range list {
		messages = append(messages, MonitorMessage{Body: "\t- " + n.Range, Level: ""})

		r = s.routerByName(n.RouterName)

		m.Networks[i].RouterName = r.Name
		m.Networks[i].RouterType = r.Type
		m.Networks[i].RouterIP = r.IP
		m.Networks[i].ClientName = s.ClientName
		m.Networks[i].DatacenterName = d.Name
		m.Networks[i].DatacenterPassword = d.Password
		m.Networks[i].DatacenterRegion = d.Region
		m.Networks[i].DatacenterType = d.Type
		m.Networks[i].DatacenterUsername = d.Username
		m.Networks[i].VCloudURL = d.VCloudURL
	}

	return m, messages
}
