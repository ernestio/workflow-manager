/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
)

type errorManager struct{}

func (em *errorManager) isAnErrorMessage(subject string) bool {
	switch subject {
	case
		"routers.create.error",
		"routers.delete.error",
		"networks.create.error",
		"networks.delete.error",
		"instances.create.error",
		"instances.delete.error",
		"instances.update.error",
		"firewalls.create.error",
		"firewalls.delete.error",
		"firewalls.update.error",
		"nats.create.error",
		"nats.delete.error",
		"nats.update.error",
		"executions.create.error":
		return true
	}

	return false
}

func (em *errorManager) markAsFailed(s *service, subject string, body []byte) *service {

	switch subject {
	case "routers.create.error":
		s = em.markRoutersCreationAsFailed(s, body)
	case "routers.delete.error":
		s = em.markRoutersDeletionAsFailed(s, body)
	case "networks.create.error":
		s = em.markNetworksCreationAsErrored(s, body)
	case "networks.delete.error":
		s = em.markNetworksDeletionAsErrored(s, body)
	case "instances.create.error":
		s = em.markInstancesCreationAsErrored(s, body)
	case "instances.delete.error":
		s = em.markInstancesDeletionAsErrored(s, body)
	case "instances.update.error":
		s = em.markInstancesUpdateAsErrored(s, body)
	case "firewalls.create.error":
		s = em.markFirewallsCreationAsErrored(s, body)
	case "firewalls.delete.error":
		s = em.markFirewallsDeletionAsErrored(s, body)
	case "firewalls.update.error":
		s = em.markFirewallsUpdateAsErrored(s, body)
	case "nats.create.error":
		s = em.markNatsCreationAsErrored(s, body)
	case "nats.delete.error":
		s = em.markNatsDeleteAsErrored(s, body)
	case "nats.update.error":
		s = em.markNatsUpdateAsErrored(s, body)
	case "executions.create.error":
		s = em.markExecutionsCreationAsErrored(s, body)
	}

	s.markAsFailed()

	return s
}

func (em *errorManager) markRoutersCreationAsFailed(s *service, body []byte) *service {
	m := RoutersCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	for i, sr := range s.Routers.Items {
		for _, mr := range m.Routers {
			if sr.Name == mr.Name {
				s.Routers.Items[i].Status = mr.Status
				if mr.Status == "errored" {
					msg := "Router " + mr.Name + " creation failed with: \n" + mr.ErrorMessage
					em.sendErrors(s.Channel(), msg)
				}
			}
		}
	}

	return s
}

func (em *errorManager) markRoutersDeletionAsFailed(s *service, body []byte) *service {
	m := RoutersCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	var errored []router
	for i, sr := range s.RoutersToDelete.Items {
		for _, mr := range m.Routers {
			if sr.Name == mr.Name {
				s.RoutersToDelete.Items[i].Status = mr.Status
				if mr.Status == "errored" {
					errored = append(errored, s.Routers.Items[i])
					msg := "Router " + mr.Name + " deletion failed with: \n" + mr.ErrorMessage
					em.sendErrors(s.Channel(), msg)
				}
			}
		}
	}

	s.Routers.Items = errored
	s.RoutersToDelete.Items = errored

	return s
}

func (em *errorManager) markNetworksCreationAsErrored(s *service, body []byte) *service {
	m := NetworksCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	if len(s.NetworksToCreate.Items) > 0 {
		for i, sr := range s.NetworksToCreate.Items {
			for _, mr := range m.Networks {
				if sr.Name == mr.Name {
					s.NetworksToCreate.Items[i].Status = mr.Status
					if mr.Status == "errored" {
						msg := "Network " + mr.Name + " creation failed with: \n" + mr.ErrorMessage
						em.sendErrors(s.Channel(), msg)
					}
				}
			}
		}
	} else {
		for i, sr := range s.Networks.Items {
			for _, mr := range m.Networks {
				if sr.Name == mr.Name {
					s.Networks.Items[i].Status = mr.Status
					if mr.Status == "errored" {
						msg := "Network " + mr.Name + " creation failed with: \n" + mr.ErrorMessage
						em.sendErrors(s.Channel(), msg)
					}
				}
			}
		}
	}

	return s
}

func (em *errorManager) markNetworksDeletionAsErrored(s *service, body []byte) *service {
	m := NetworksCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	var errored []network
	for i, sr := range s.NetworksToDelete.Items {
		for _, mr := range m.Networks {
			if sr.Name == mr.Name {
				s.NetworksToDelete.Items[i].Status = mr.Status
				if mr.Status == "errored" {
					errored = append(errored, s.Networks.Items[i])
					msg := "Network " + mr.Name + " deletion failed with: \n" + mr.ErrorMessage
					em.sendErrors(s.Channel(), msg)
				}
			}
		}
	}
	s.NetworksToDelete.Items = errored
	s.Networks.Items = errored

	return s
}

func (em *errorManager) markInstancesCreationAsErrored(s *service, body []byte) *service {
	m := InstancesCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	if len(s.InstancesToCreate.Items) > 0 {
		for i, sr := range s.InstancesToCreate.Items {
			for _, mr := range m.Instances {
				if sr.Name == mr.Name {
					s.InstancesToCreate.Items[i].Status = mr.Status
					if mr.Status == "errored" {
						msg := "Instance " + mr.Name + " creation failed with: \n" + mr.ErrorMessage
						em.sendErrors(s.Channel(), msg)
					}
				}
			}
		}
	} else {
		for i, sr := range s.Instances.Items {
			for _, mr := range m.Instances {
				if sr.Name == mr.Name {
					s.Instances.Items[i].Status = mr.Status
					if mr.Status == "errored" {
						msg := "Instance " + mr.Name + " creation failed with: \n" + mr.ErrorMessage
						em.sendErrors(s.Channel(), msg)
					}
				}
			}
		}
	}

	return s
}

func (em *errorManager) markInstancesUpdateAsErrored(s *service, body []byte) *service {
	m := InstancesCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	for i, sr := range s.InstancesToUpdate.Items {
		for _, mr := range m.Instances {
			if sr.Name == mr.Name {
				s.InstancesToUpdate.Items[i].Status = mr.Status
				if mr.Status == "errored" {
					msg := "Instance " + mr.Name + " modification failed with: \n" + mr.ErrorMessage
					em.sendErrors(s.Channel(), msg)
				}
			}
		}
	}

	return s
}

func (em *errorManager) markInstancesDeletionAsErrored(s *service, body []byte) *service {
	m := InstancesCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	errored := make([]instance, 0)
	for i, sr := range s.InstancesToDelete.Items {
		for _, mr := range m.Instances {
			if sr.Name == mr.Name {
				s.InstancesToDelete.Items[i].Status = mr.Status
				if mr.Status == "errored" {
					errored = append(errored, s.Instances.Items[i])
					msg := "Instance " + mr.Name + " removal failed with: \n" + mr.ErrorMessage
					em.sendErrors(s.Channel(), msg)
				}
			}
		}
	}

	s.InstancesToDelete.Items = errored
	s.Instances.Items = errored

	return s
}

func (em *errorManager) markFirewallsCreationAsErrored(s *service, body []byte) *service {
	m := FirewallsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	for i, sr := range s.Firewalls.Items {
		for _, mr := range m.Firewalls {
			if sr.Name == mr.Name {
				s.Firewalls.Items[i].Status = mr.Status
				if mr.Status == "errored" {
					msg := "Firewall " + mr.Name + " creation failed with: \n" + mr.ErrorMessage
					em.sendErrors(s.Channel(), msg)
				}
			}
		}
	}

	return s
}

func (em *errorManager) markFirewallsDeletionAsErrored(s *service, body []byte) *service {
	m := FirewallsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	for i, sr := range s.FirewallsToDelete.Items {
		for _, mr := range m.Firewalls {
			if sr.Name == mr.Name {
				s.FirewallsToDelete.Items[i].Status = mr.Status
				if mr.Status == "errored" {
					msg := "Firewall " + mr.Name + " creation failed with: \n" + mr.ErrorMessage
					em.sendErrors(s.Channel(), msg)
				}
			}
		}
	}

	return s
}

func (em *errorManager) markFirewallsUpdateAsErrored(s *service, body []byte) *service {
	m := FirewallsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	for i, sr := range s.Firewalls.Items {
		for _, mr := range m.Firewalls {
			if sr.Name == mr.Name {
				s.Firewalls.Items[i].Status = mr.Status
				if mr.Status == "errored" {
					msg := "Firewall " + mr.Name + " modification failed with: \n" + mr.ErrorMessage
					em.sendErrors(s.Channel(), msg)
				}
			}
		}
	}

	return s
}

func (em *errorManager) markNatsCreationAsErrored(s *service, body []byte) *service {
	m := NatsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	for i, sr := range s.Nats.Items {
		for _, mr := range m.Nats {
			if sr.Name == mr.Name {
				s.Nats.Items[i].Status = mr.Status
				if mr.Status == "errored" {
					msg := "Nats " + mr.Name + " creation failed with: \n" + mr.ErrorMessage
					em.sendErrors(s.Channel(), msg)
				}
			}
		}
	}

	return s
}

func (em *errorManager) markNatsUpdateAsErrored(s *service, body []byte) *service {
	m := NatsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	for i, sr := range s.Nats.Items {
		for _, mr := range m.Nats {
			if sr.Name == mr.Name {
				s.Nats.Items[i].Status = mr.Status
				if mr.Status == "errored" {
					msg := "Nats " + mr.Name + " modification failed with: \n" + mr.ErrorMessage
					em.sendErrors(s.Channel(), msg)
				}
			}
		}
	}

	return s
}

func (em *errorManager) markNatsDeleteAsErrored(s *service, body []byte) *service {
	m := NatsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	for i, sr := range s.NatsToDelete.Items {
		for _, mr := range m.Nats {
			if sr.Name == mr.Name {
				s.NatsToDelete.Items[i].Status = mr.Status
				if mr.Status == "errored" {
					msg := "Nats " + mr.Name + " modification failed with: \n" + mr.ErrorMessage
					em.sendErrors(s.Channel(), msg)
				}
			}
		}
	}

	return s
}

func (em *errorManager) markExecutionsCreationAsErrored(s *service, body []byte) *service {
	m := ExecutionsCreate{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}

	if s.Status == "bootstrapping" {
		for i, sr := range s.Bootstraps.Items {
			for _, mr := range m.Executions {
				if sr.Name == mr.Name {
					s.Bootstraps.Items[i].Status = mr.Status
					if mr.Status == "errored" {
						msg := "Bootstrapping " + mr.Payload + " failed with: \n" + mr.ErrorMessage
						em.sendErrors(s.Channel(), msg)
					}
				}
			}
		}
	} else if s.Status == "running_executions" {
		for i, sr := range s.ExecutionsToCreate.Items {
			for _, mr := range m.Executions {
				if sr.Name == mr.Name {
					s.ExecutionsToCreate.Items[i].Status = mr.Status
					if mr.Status == "errored" {
						msg := "Executing " + mr.Payload + " failed with: \n" + mr.ErrorMessage
						em.sendErrors(s.Channel(), msg)
					}
				}
			}
		}
	}

	return s
}

func (em *errorManager) sendErrors(channel string, message string) {
	messages := []MonitorMessage{}
	messages = append(messages, MonitorMessage{Body: message, Level: "ERROR"})
	UserOutput(channel, messages)
}
