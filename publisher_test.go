/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVitamineTemplating(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		var p publisher
		var comp map[string]interface{}
		var incomp map[string]interface{}

		dataComplete := h.getFixture("./fixtures/publisher.json")
		dataIncomplete := h.getFixture("./fixtures/publisher_incomplete_firewalls.json")
		s := h.getService("./fixtures/publisher.json")
		si := h.getService("./fixtures/publisher_incomplete_firewalls.json")

		json.Unmarshal(dataComplete, &comp)
		json.Unmarshal(dataIncomplete, &incomp)

		Convey("When i try and template fields on an collection of instances where all fields are known", func() {
			x := comp["instances"].(map[string]interface{})["items"].([]interface{})
			items := p.UpdateTemplateVariables(x, &s)

			Convey("It should have mapped all string fields", func() {
				collection, ok := items[0].(map[string]interface{})
				So(ok, ShouldBeTrue)
				item, ok := collection["network_aws_id"].(string)
				So(ok, ShouldBeTrue)
				So(item, ShouldEqual, "network-1-id")
			})

			Convey("It should have mapped all slice fields", func() {
				collection, ok := items[0].(map[string]interface{})
				So(ok, ShouldBeTrue)
				itemsl, ok := collection["security_group_aws_ids"].([]interface{})
				So(ok, ShouldBeTrue)
				item, ok := itemsl[0].(string)
				So(ok, ShouldBeTrue)
				So(item, ShouldEqual, "firewall-1-id")
			})

		})

		Convey("When i try and template fields on an collection of instances where not all fields are known", func() {
			x := incomp["instances"].(map[string]interface{})["items"].([]interface{})
			items := p.UpdateTemplateVariables(x, &si)

			Convey("It should not have mapped fields where there was no result", func() {
				collection, ok := items[0].(map[string]interface{})
				So(ok, ShouldBeTrue)
				itemsl, ok := collection["security_group_aws_ids"].([]interface{})
				So(ok, ShouldBeTrue)
				item, ok := itemsl[0].(string)
				So(ok, ShouldBeTrue)
				So(item, ShouldNotEqual, "")
				So(item, ShouldNotEqual, "null")
				fmt.Println(item)
			})
		})

	})

}

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
			m := &GenericComponentMsg{}

			d := s.Datacenters.Items[0]
			r := s.Routers.Items[0]
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				n := m.Components[0].(map[string]interface{})
				So(n["name"].(string), ShouldEqual, s.NetworksToCreate.Items[0].Name)
				So(n["range"].(string), ShouldEqual, s.NetworksToCreate.Items[0].Range)
				So(n["router_name"].(string), ShouldEqual, r.Name)
				So(n["router_type"].(string), ShouldEqual, r.Type)
				So(n["router_ip"].(string), ShouldEqual, r.IP)
				So(n["client_name"].(string), ShouldEqual, s.ClientName)
				So(n["datacenter_name"].(string), ShouldEqual, d.Name)
				So(n["datacenter_password"].(string), ShouldEqual, d.Password)
				So(n["datacenter_region"].(string), ShouldEqual, d.Region)
				So(n["datacenter_type"].(string), ShouldEqual, d.Type)
				So(n["datacenter_username"].(string), ShouldEqual, d.Username)

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
			m := &GenericComponentMsg{}
			d := s.Datacenters.Items[0]
			r := s.Routers.Items[0]
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				n := m.Components[0].(map[string]interface{})
				So(n["name"].(string), ShouldEqual, s.Networks.Items[0].Name)
				So(n["range"].(string), ShouldEqual, s.Networks.Items[0].Range)
				So(n["router_name"].(string), ShouldEqual, r.Name)
				So(n["router_type"].(string), ShouldEqual, r.Type)
				So(n["router_ip"].(string), ShouldEqual, r.IP)
				So(n["client_name"].(string), ShouldEqual, s.ClientName)
				So(n["datacenter_name"].(string), ShouldEqual, d.Name)
				So(n["datacenter_password"].(string), ShouldEqual, d.Password)
				So(n["datacenter_region"].(string), ShouldEqual, d.Region)
				So(n["datacenter_type"].(string), ShouldEqual, d.Type)
				So(n["datacenter_username"].(string), ShouldEqual, d.Username)

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
			m := &GenericComponentMsg{}
			d := s.Datacenters.Items[0]
			n := s.Networks.Items[0]
			sg := s.Firewalls.Items[0]
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 2)
				i := m.Components[0].(map[string]interface{})
				So(i["name"], ShouldEqual, s.InstancesToCreate.Items[0].Name)
				So(i["type"], ShouldEqual, s.InstancesToCreate.Items[0].Type)
				So(i["ip"], ShouldEqual, s.InstancesToCreate.Items[0].IP)
				So(i["datacenter_name"].(string), ShouldEqual, d.Name)
				So(i["datacenter_password"].(string), ShouldEqual, d.Password)
				So(i["datacenter_region"].(string), ShouldEqual, d.Region)
				So(i["datacenter_type"].(string), ShouldEqual, d.Type)
				So(i["datacenter_username"].(string), ShouldEqual, d.Username)
				So(i["network_aws_id"].(string), ShouldEqual, n.NetworkAWSID)
				So(i["security_group_aws_ids"].([]interface{})[0].(string), ShouldEqual, sg.SecurityGroupAWSID)
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
			m := &GenericComponentMsg{}
			d := s.Datacenters.Items[0]
			ro := s.Routers.Items[0]
			json.Unmarshal([]byte(body), &m)
			fmt.Println(body)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				n := m.Components[0].(map[string]interface{})
				So(n["name"], ShouldEqual, s.NatsToCreate.Items[0].Name)
				So(n["status"], ShouldEqual, s.NatsToCreate.Items[0].Status)
				So(n["router_name"], ShouldEqual, ro.Name)
				So(n["router_type"], ShouldEqual, ro.Type)
				So(n["router_ip"], ShouldEqual, ro.IP)
				So(n["client_name"], ShouldEqual, s.ClientName)
				So(n["datacenter_name"], ShouldEqual, d.Name)
				So(n["datacenter_type"], ShouldEqual, d.Type)
				So(n["datacenter_region"], ShouldEqual, d.Region)
				So(n["datacenter_username"], ShouldEqual, d.Username)
				So(n["datacenter_password"], ShouldEqual, d.Password)
				r := n["rules"].([]interface{})[0].(map[string]interface{})
				So(r["protocol"], ShouldEqual, "protocol")
				So(r["type"], ShouldEqual, "type")
				So(r["network"], ShouldEqual, "network")
				So(r["origin_ip"], ShouldEqual, "11.11.11.11/11")
				So(r["origin_port"], ShouldEqual, "1")
				So(r["translation_ip"], ShouldEqual, "10.10.10.10/10")
				So(r["translation_port"], ShouldEqual, "1")

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
			m := &GenericComponentMsg{}
			d := s.Datacenters.Items[0]
			ro := s.Routers.Items[0]
			json.Unmarshal([]byte(body), &m)
			fmt.Println(body)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				n := m.Components[0].(map[string]interface{})
				So(n["name"], ShouldEqual, s.NatsToUpdate.Items[0].Name)
				So(n["status"], ShouldEqual, s.NatsToUpdate.Items[0].Status)
				So(n["router_name"], ShouldEqual, ro.Name)
				So(n["router_type"], ShouldEqual, ro.Type)
				So(n["router_ip"], ShouldEqual, ro.IP)
				So(n["client_name"], ShouldEqual, s.ClientName)
				So(n["datacenter_name"], ShouldEqual, d.Name)
				So(n["datacenter_type"], ShouldEqual, d.Type)
				So(n["datacenter_region"], ShouldEqual, d.Region)
				So(n["datacenter_username"], ShouldEqual, d.Username)
				So(n["datacenter_password"], ShouldEqual, d.Password)
				r := n["rules"].([]interface{})[0].(map[string]interface{})
				So(r["protocol"], ShouldEqual, "protocol")
				So(r["type"], ShouldEqual, "type")
				So(r["network"], ShouldEqual, "network")
				So(r["origin_ip"], ShouldEqual, "11.11.11.11/11")
				So(r["origin_port"], ShouldEqual, "1")
				So(r["translation_ip"], ShouldEqual, "10.10.10.10/10")
				So(r["translation_port"], ShouldEqual, "1")

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
			m := &GenericComponentMsg{}
			d := s.Datacenters.Items[0]
			ro := s.Routers.Items[0]
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				f := m.Components[0].(map[string]interface{})
				So(f["name"], ShouldEqual, s.FirewallsToCreate.Items[0].Name)
				So(f["status"], ShouldEqual, s.FirewallsToCreate.Items[0].Status)
				So(f["router_name"], ShouldEqual, ro.Name)
				So(f["router_type"], ShouldEqual, ro.Type)
				So(f["datacenter_name"], ShouldEqual, d.Name)
				So(f["datacenter_type"], ShouldEqual, d.Type)
				So(f["datacenter_region"], ShouldEqual, d.Region)
				So(f["datacenter_username"], ShouldEqual, d.Username)
				So(f["datacenter_password"], ShouldEqual, d.Password)
				r := f["rules"].([]interface{})[0].(map[string]interface{})
				So(r["source_ip"], ShouldEqual, "11.11.11.11/11")
				So(r["source_port"], ShouldEqual, "source_port")
				So(r["protocol"], ShouldEqual, "protocol")
				So(r["destination_ip"], ShouldEqual, "any")
				So(r["destination_port"], ShouldEqual, "destination_port")

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
			m := &GenericComponentMsg{}
			d := s.Datacenters.Items[0]
			ro := s.Routers.Items[0]
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				f := m.Components[0].(map[string]interface{})
				So(f["name"], ShouldEqual, s.FirewallsToUpdate.Items[0].Name)
				So(f["status"], ShouldEqual, s.FirewallsToUpdate.Items[0].Status)
				So(f["router_name"], ShouldEqual, ro.Name)
				So(f["router_type"], ShouldEqual, ro.Type)
				So(f["datacenter_name"], ShouldEqual, d.Name)
				So(f["datacenter_type"], ShouldEqual, d.Type)
				So(f["datacenter_region"], ShouldEqual, d.Region)
				So(f["datacenter_username"], ShouldEqual, d.Username)
				So(f["datacenter_password"], ShouldEqual, d.Password)
				r := f["rules"].([]interface{})[0].(map[string]interface{})
				So(r["source_ip"], ShouldEqual, "11.11.11.11/11")
				So(r["source_port"], ShouldEqual, "source_port")
				So(r["protocol"], ShouldEqual, "protocol")
				So(r["destination_ip"], ShouldEqual, "any")
				So(r["destination_port"], ShouldEqual, "destination_port")

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
				m := &GenericComponentMsg{}
				json.Unmarshal([]byte(body), &m)

				Convey("Then I'll receive a valid json string", func() {
					So(m.Service, ShouldEqual, s.ID)
					So(len(m.Components), ShouldEqual, 1)
					b := m.Components[0].(map[string]interface{})

					So(b["name"], ShouldEqual, s.BootstrapsToCreate.Items[0].Name)
					So(b["type"], ShouldEqual, "salt")
					So(b["payload"], ShouldEqual, s.BootstrapsToCreate.Items[0].Payload)
					So(b["target"], ShouldEqual, s.BootstrapsToCreate.Items[0].Target)
					So(b["status"], ShouldEqual, s.BootstrapsToCreate.Items[0].Status)
					So(b["user"], ShouldEqual, "")
					So(b["password"], ShouldEqual, "")
					So(b["service_endpoint"], ShouldEqual, s.Routers.Items[0].IP)

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
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(err, ShouldEqual, nil)

				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				e := m.Components[0].(map[string]interface{})

				So(e["payload"], ShouldEqual, s.ExecutionsToCreate.Items[0].Payload)
				So(e["name"], ShouldEqual, s.ExecutionsToCreate.Items[0].Name)
				So(e["type"], ShouldEqual, "salt")
				So(e["user"], ShouldEqual, "")
				So(e["password"], ShouldEqual, "")
				So(e["target"], ShouldEqual, s.ExecutionsToCreate.Items[0].Target)
				So(e["service_endpoint"], ShouldEqual, s.Routers.Items[0].IP)
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
