/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnMappedMessage(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service.json")

		Convey("When I try to get body for the unmapped message", func() {
			mm := messageManager{}
			message, err := mm.preparePublishMessage("test.message.invalid", &s)

			Convey("Then I'll receive an empty body and an error", func() {
				So(message, ShouldEqual, "")
				So(err, ShouldNotEqual, nil)

			})
		})
	})
}

func TestCreateRouters(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_create_routers.json")

		Convey("When I get the message for a routers.create event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("routers.create", &s)
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				r := m.Components[0].(map[string]interface{})
				So(len(m.Components), ShouldEqual, 1)
				So(r["name"].(string), ShouldEqual, s.RoutersToCreate.Items[0].Name)
				So(r["type"].(string), ShouldEqual, s.RoutersToCreate.Items[0].Type)

				d := s.Datacenters.Items[0]
				So(r["datacenter_name"].(string), ShouldEqual, d.Name)
				So(r["datacenter_password"].(string), ShouldEqual, d.Password)
				So(r["datacenter_region"].(string), ShouldEqual, d.Region)
				So(r["datacenter_type"].(string), ShouldEqual, d.Type)
				So(r["datacenter_username"].(string), ShouldEqual, d.Username)
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestPublisherCreateError(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_real_workflow.json")

		Convey("When I get the message for a services.create.error event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("service.create.error", &s)
			m := &service{}
			err = json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.ID, ShouldEqual, s.ID)
				So(m.Status, ShouldEqual, "errored")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestCreateNetworks(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_create_networks.json")

		Convey("When I get the message for a networks.create event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("networks.create", &s)
			m := &NetworksCreate{}
			d := s.Datacenters.Items[0]
			r := s.Routers.Items[0]
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Networks), ShouldEqual, 1)
				n := m.Networks[0]
				So(n.Name, ShouldEqual, s.NetworksToCreate.Items[0].Name)
				So(n.Range, ShouldEqual, s.NetworksToCreate.Items[0].Range)
				So(n.RouterName, ShouldEqual, r.Name)
				So(n.RouterType, ShouldEqual, r.Type)
				So(n.RouterIP, ShouldEqual, r.IP)
				So(n.ClientName, ShouldEqual, s.ClientName)
				So(n.DatacenterName, ShouldEqual, d.Name)
				So(n.DatacenterPassword, ShouldEqual, d.Password)
				So(n.DatacenterRegion, ShouldEqual, d.Region)
				So(n.DatacenterType, ShouldEqual, d.Type)
				So(n.DatacenterUsername, ShouldEqual, d.Username)

				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestDeleteNetworks(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_real_workflow.json")

		Convey("When I get the message for a networks.delete event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("networks.delete", &s)
			m := &NetworksCreate{}
			d := s.Datacenters.Items[0]
			r := s.Routers.Items[0]
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Networks), ShouldEqual, 1)
				n := m.Networks[0]
				So(n.Name, ShouldEqual, s.Networks.Items[0].Name)
				So(n.Range, ShouldEqual, s.Networks.Items[0].Range)
				So(n.RouterName, ShouldEqual, r.Name)
				So(n.RouterType, ShouldEqual, r.Type)
				So(n.RouterIP, ShouldEqual, r.IP)
				So(n.ClientName, ShouldEqual, s.ClientName)
				So(n.DatacenterName, ShouldEqual, d.Name)
				So(n.DatacenterPassword, ShouldEqual, d.Password)
				So(n.DatacenterRegion, ShouldEqual, d.Region)
				So(n.DatacenterType, ShouldEqual, d.Type)
				So(n.DatacenterUsername, ShouldEqual, d.Username)

				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestCreateInstances(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_create_instances.json")

		Convey("When I get the message for a instances.create event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("instances.create", &s)
			m := &InstancesCreate{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Instances), ShouldEqual, 2)
				i := m.Instances[0]
				So(i.Name, ShouldEqual, s.InstancesToCreate.Items[0].Name)

				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestCreateNats(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_create_nats.json")

		Convey("When I get the message for a nats.create event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("nats.create", &s)
			m := &NatsCreate{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				d := s.Datacenters.Items[0]
				ro := s.Routers.Items[0]

				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Nats), ShouldEqual, 1)
				n := m.Nats[0]
				So(n.Name, ShouldEqual, s.NatsToCreate.Items[0].Name)
				So(n.Status, ShouldEqual, s.NatsToCreate.Items[0].Status)
				So(n.RouterName, ShouldEqual, ro.Name)
				So(n.RouterType, ShouldEqual, ro.Type)
				So(n.RouterIP, ShouldEqual, ro.IP)
				So(n.ClientName, ShouldEqual, s.ClientName)
				So(n.DatacenterName, ShouldEqual, d.Name)
				So(n.DatacenterType, ShouldEqual, d.Type)
				So(n.DatacenterRegion, ShouldEqual, d.Region)
				So(n.DatacenterUsername, ShouldEqual, d.Username)
				So(n.DatacenterPassword, ShouldEqual, d.Password)
				r := n.Rules[0]
				So(r.Protocol, ShouldEqual, "protocol")
				So(r.Type, ShouldEqual, "type")
				So(r.Network, ShouldEqual, "network")
				So(r.OriginIP, ShouldEqual, "11.11.11.11/11")
				So(r.OriginPort, ShouldEqual, "1")
				So(r.TranslationIP, ShouldEqual, "10.10.10.10/10")
				So(r.TranslationPort, ShouldEqual, "1")

				So(err, ShouldEqual, nil)
			})
		})
	})
}

func TestUpdateNats(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_update_nats.json")

		Convey("When I get the message for a nats.update event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("nats.update", &s)
			m := &NatsCreate{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				d := s.Datacenters.Items[0]
				ro := s.Routers.Items[0]

				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Nats), ShouldEqual, 1)
				n := m.Nats[0]
				So(n.Name, ShouldEqual, s.NatsToUpdate.Items[0].Name)
				So(n.RouterName, ShouldEqual, ro.Name)
				So(n.RouterType, ShouldEqual, ro.Type)
				So(n.RouterIP, ShouldEqual, ro.IP)
				So(n.ClientName, ShouldEqual, s.ClientName)
				So(n.DatacenterName, ShouldEqual, d.Name)
				So(n.DatacenterType, ShouldEqual, d.Type)
				So(n.DatacenterRegion, ShouldEqual, d.Region)
				So(n.DatacenterUsername, ShouldEqual, d.Username)
				So(n.DatacenterPassword, ShouldEqual, d.Password)
				r := n.Rules[0]
				So(r.Protocol, ShouldEqual, "protocol")
				So(r.Type, ShouldEqual, "type")
				So(r.Network, ShouldEqual, "network")
				So(r.OriginIP, ShouldEqual, "11.11.11.11/11")
				So(r.OriginPort, ShouldEqual, "1")
				So(r.TranslationIP, ShouldEqual, "10.10.10.10/10")
				So(r.TranslationPort, ShouldEqual, "1")

				So(err, ShouldEqual, nil)
			})
		})
	})
}

func TestCreateFirewalls(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_create_firewalls.json")

		Convey("When I get the message for a firewalls.create event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("firewalls.create", &s)
			m := &FirewallsCreate{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Firewalls), ShouldEqual, 1)
				So(len(m.Networks), ShouldEqual, 1)
				f := m.Firewalls[0]
				n := m.Networks[0]
				So(f.Name, ShouldEqual, s.FirewallsToCreate.Items[0].Name)
				So(f.Status, ShouldEqual, s.FirewallsToCreate.Items[0].Status)
				So(n.Name, ShouldEqual, s.Networks.Items[0].Name)
				So(n.Range, ShouldEqual, s.Networks.Items[0].Range)
				r := f.Rules[0]
				So(r.Source, ShouldEqual, "11.11.11.11/11")
				So(r.SourcePort, ShouldEqual, "source_port")
				So(r.Protocol, ShouldEqual, "protocol")
				So(r.Destination, ShouldEqual, "any")
				So(r.DestinationPort, ShouldEqual, "destination_port")

				So(err, ShouldEqual, nil)
			})
		})
	})
}

func TestUpdateFirewalls(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_update_firewalls.json")

		Convey("When I get the message for a firewalls.update event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("firewalls.update", &s)
			m := &FirewallsCreate{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Firewalls), ShouldEqual, 1)
				So(len(m.Networks), ShouldEqual, 1)
				f := m.Firewalls[0]
				n := m.Networks[0]
				So(f.Name, ShouldEqual, s.Firewalls.Items[0].Name)
				So(f.Status, ShouldEqual, s.Firewalls.Items[0].Status)
				So(n.Name, ShouldEqual, s.Networks.Items[0].Name)
				So(n.Range, ShouldEqual, s.Networks.Items[0].Range)
				r := f.Rules[0]
				So(r.Source, ShouldEqual, "11.11.11.11/11")
				So(r.SourcePort, ShouldEqual, "source_port")
				So(r.Protocol, ShouldEqual, "protocol")
				So(r.Destination, ShouldEqual, "any")
				So(r.DestinationPort, ShouldEqual, "destination_port")

				So(err, ShouldEqual, nil)
			})
		})
	})
}

func TestCreateBootstraps(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_real_workflow.json")
		s.Status = "nats_created"
		s.save()

		Convey("When I get the message for a executions.create event and status nats_created", func() {
			Convey("If a service_ip is not provided", func() {
				mm := messageManager{}
				body, err := mm.preparePublishMessage("bootstraps.create", &s)
				m := &ExecutionsCreate{}
				json.Unmarshal([]byte(body), &m)

				Convey("Then I'll receive a valid json string", func() {
					So(m.Service, ShouldEqual, s.ID)
					So(m.ServiceName, ShouldEqual, s.Name)
					So(m.ServiceType, ShouldEqual, s.Type)
					So(len(m.Executions), ShouldEqual, 1)

					e := m.Executions[0]
					So(e.Name, ShouldEqual, s.BootstrapsToCreate.Items[0].Name)
					So(e.Type, ShouldEqual, "salt")
					So(e.Payload, ShouldEqual, s.BootstrapsToCreate.Items[0].Payload)
					So(e.Target, ShouldEqual, s.BootstrapsToCreate.Items[0].Target)
					So(e.Status, ShouldEqual, s.BootstrapsToCreate.Items[0].Status)
					So(e.User, ShouldEqual, "")
					So(e.Password, ShouldEqual, "")
					So(e.EndPoint, ShouldEqual, s.Routers.Items[0].IP)
					So(e.ClientName, ShouldEqual, s.ClientName)

					d := s.Datacenters.Items[0]
					So(e.DatacenterName, ShouldEqual, d.Name)
					So(e.DatacenterPassword, ShouldEqual, d.Password)
					So(e.DatacenterRegion, ShouldEqual, d.Region)
					So(e.DatacenterType, ShouldEqual, d.Type)
					So(e.DatacenterUsername, ShouldEqual, d.Username)

					So(err, ShouldEqual, nil)
				})
			})

			Convey("If a service_ip is provided", func() {
				mm := messageManager{}
				s.ServiceIP = "1.1.1.1"
				body, err := mm.preparePublishMessage("executions.create", &s)
				m := &ExecutionsCreate{}
				json.Unmarshal([]byte(body), &m)

				Convey("Then I'll receive a valid json string", func() {
					So(m.Service, ShouldEqual, s.ID)
					So(m.ServiceName, ShouldEqual, s.Name)
					So(m.ServiceType, ShouldEqual, s.Type)
					So(len(m.Executions), ShouldEqual, 1)

					e := m.Executions[0]
					So(e.Name, ShouldEqual, s.ExecutionsToCreate.Items[0].Name)
					So(e.Type, ShouldEqual, "salt")
					So(e.Payload, ShouldEqual, s.ExecutionsToCreate.Items[0].Payload)
					So(e.Target, ShouldEqual, s.ExecutionsToCreate.Items[0].Target)
					So(e.Status, ShouldEqual, s.ExecutionsToCreate.Items[0].Status)
					So(e.User, ShouldEqual, "")
					So(e.Password, ShouldEqual, "")
					So(e.EndPoint, ShouldEqual, "1.1.1.1")
					So(e.ClientName, ShouldEqual, s.ClientName)

					d := s.Datacenters.Items[0]
					So(e.DatacenterName, ShouldEqual, d.Name)
					So(e.DatacenterPassword, ShouldEqual, d.Password)
					So(e.DatacenterRegion, ShouldEqual, d.Region)
					So(e.DatacenterType, ShouldEqual, d.Type)
					So(e.DatacenterUsername, ShouldEqual, d.Username)

					So(err, ShouldEqual, nil)
				})
			})

		})
	})
}

func TestCreateExecutions(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_real_workflow.json")
		s.Status = "bootstrap_ran"
		s.Bootstraps.Finished = "yes"
		s.save()

		Convey("When I get the message for a executions.create event and status nats_created", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("executions.create", &s)
			m := &ExecutionsCreate{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(err, ShouldEqual, nil)

				So(m.Service, ShouldEqual, s.ID)
				So(m.ServiceName, ShouldEqual, s.Name)
				So(m.ServiceType, ShouldEqual, s.Type)
				So(len(m.Executions), ShouldEqual, 1)

				e := m.Executions[0]
				So(e.Payload, ShouldEqual, s.ExecutionsToCreate.Items[0].Payload)
				So(e.Name, ShouldEqual, s.ExecutionsToCreate.Items[0].Name)
				So(e.Type, ShouldEqual, "salt")
				So(e.User, ShouldEqual, "")
				So(e.Password, ShouldEqual, "")
				So(e.Target, ShouldEqual, s.ExecutionsToCreate.Items[0].Target)
				So(e.ClientName, ShouldEqual, s.ClientName)
				So(e.EndPoint, ShouldEqual, s.Routers.Items[0].IP)

				d := s.Datacenters.Items[0]
				So(e.DatacenterName, ShouldEqual, d.Name)
				So(e.DatacenterPassword, ShouldEqual, d.Password)
				So(e.DatacenterRegion, ShouldEqual, d.Region)
				So(e.DatacenterType, ShouldEqual, d.Type)
				So(e.DatacenterUsername, ShouldEqual, d.Username)

			})
		})
	})
}

func TestServiceDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s := h.getService("./fixtures/service_real_workflow.json")
		s.Status = "executions_ran"
		s.save()

		Convey("When I get the message for a executions.create event and status nats_created", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("service.create.done", &s)
			m := service{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(err, ShouldEqual, nil)

				So(m.ID, ShouldEqual, s.ID)
				So(m.Name, ShouldEqual, s.Name)
				So(m.Status, ShouldEqual, s.Status)
			})
		})
	})

}
