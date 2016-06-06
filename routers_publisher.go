/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// Creates a CreateRouters struct based on a given service
func buildCreateRouters(s *service) RoutersCreate {
	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Creating routers:", Level: "INFO"})

	m := buildRoutersList(s, s.Routers.Items)

	UserOutput(s.Channel(), messages)

	return m
}

// Creates a DeleteRouters struct based on a given service
func buildDeleteRouters(s *service) RoutersCreate {
	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: "Deleting router:", Level: "INFO"})

	m := buildRoutersList(s, s.RoutersToDelete.Items)

	UserOutput(s.Channel(), messages)

	return m
}

func buildRoutersList(s *service, list []router) RoutersCreate {
	d := s.datacenter()

	m := RoutersCreate{
		Service: s.ID,
		Routers: list,
		Status:  s.Routers.Status,
	}

	for i := range m.Routers {

		m.Routers[i].ClientName = s.ClientName
		m.Routers[i].DatacenterName = d.Name
		m.Routers[i].DatacenterPassword = d.Password
		m.Routers[i].DatacenterRegion = d.Region
		m.Routers[i].DatacenterType = d.Type
		m.Routers[i].DatacenterUsername = d.Username
		m.Routers[i].ExternalNetwork = d.ExternalNetwork
		m.Routers[i].VCloudURL = d.VCloudURL
		m.Routers[i].VseURL = d.VseURL
	}

	return m
}
