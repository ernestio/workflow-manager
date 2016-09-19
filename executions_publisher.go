/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// Creates a CreateExecutions with Executions struct based on a given service
func buildCreateExecutions(s *service) ExecutionsCreate {
	if len(s.ExecutionsToCreate.Items) > 0 {
		messages := []MonitorMessage{}
		messages = append(messages, MonitorMessage{Body: "Running executions", Level: "INFO"})
		UserOutput(s.Channel(), messages)
	}
	res := buildBasicExecutionsCreate(s, s.ExecutionsToCreate.Items)

	return res
}

// Creates a CreateExecutions with Bootstraps struct based on a given service
func buildCreateBootstraps(s *service) ExecutionsCreate {
	if len(s.BootstrapsToCreate.Items) > 0 {
		messages := []MonitorMessage{}
		messages = append(messages, MonitorMessage{Body: "Bootstrapping", Level: "INFO"})
		UserOutput(s.Channel(), messages)
	}
	res := buildBasicExecutionsCreate(s, s.BootstrapsToCreate.Items)

	return res
}

func buildBasicExecutionsCreate(s *service, inputItems []execution) ExecutionsCreate {
	items := make([]execution, len(inputItems))
	copy(items, inputItems)

	d := s.datacenter()

	// TODO: This should be modified once we support multiple routers
	endpoint := ""
	if len(s.Routers.Items) > 0 {
		r := s.Routers.Items[0]
		endpoint = r.IP
	}

	if s.ServiceIP != "" {
		endpoint = s.ServiceIP
	}

	m := ExecutionsCreate{
		Service:     s.ID,
		ServiceName: s.Name,
		ServiceType: s.Type,
		Executions:  items,
	}

	for i, e := range items {
		m.Executions[i].Name = e.Name
		m.Executions[i].Type = "salt"
		if d.Type == "fake" {
			m.Executions[i].Type = d.Type
		}
		m.Executions[i].Payload = e.Payload
		m.Executions[i].Target = e.Target
		m.Executions[i].ClientName = s.ClientName
		m.Executions[i].DatacenterName = d.Name
		m.Executions[i].DatacenterPassword = d.Password
		m.Executions[i].DatacenterRegion = d.Region
		m.Executions[i].DatacenterType = d.Type
		m.Executions[i].DatacenterUsername = d.Username
		m.Executions[i].User = c.SaltAuthentication.User
		m.Executions[i].Password = c.SaltAuthentication.Password
		m.Executions[i].Status = e.Status
		m.Executions[i].EndPoint = endpoint
	}

	return m
}
