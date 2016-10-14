/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tidwall/gjson"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVitamineTemplating(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		var p publisher
		var comp map[string]interface{}
		var incomp map[string]interface{}

		dataComplete := h.getFixture("./fixtures/publisher.json")
		dataIncomplete := h.getFixture("./fixtures/publisher_incomplete_firewalls.json")
		s, _ := h.getService("./fixtures/publisher.json")
		si, _ := h.getService("./fixtures/publisher_incomplete_firewalls.json")

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
		s, _ := h.getService("./fixtures/service.json")

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
		s, sBody := h.getService("./fixtures/service_create_routers.json")

		Convey("When I get the message for a routers.create event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("routers.create", &s)
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				r := m.Components[0].(map[string]interface{})
				So(len(m.Components), ShouldEqual, 1)
				So(r["name"].(string), ShouldEqual, gjson.Get(sBody, "routers_to_create.items.0.name").String())
				So(r["type"].(string), ShouldEqual, gjson.Get(sBody, "routers_to_create.items.0.type").String())

				So(r["datacenter_name"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.name").String())
				So(r["datacenter_password"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.password").String())
				So(r["datacenter_region"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.region").String())
				So(r["datacenter_type"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.type").String())
				So(r["datacenter_username"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.username").String())
				So(err, ShouldEqual, nil)
			})
		})
	})
}

func TestPublisherCreateError(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s, _ := h.getService("./fixtures/service_real_workflow.json")

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
		s, sBody := h.getService("./fixtures/service_create_networks.json")

		Convey("When I get the message for a networks.create event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("networks.create", &s)
			m := &GenericComponentMsg{}

			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(m.SequentialProcessing, ShouldBeTrue)
				So(len(m.Components), ShouldEqual, 1)
				n := m.Components[0].(map[string]interface{})
				So(n["name"].(string), ShouldEqual, gjson.Get(sBody, "networks_to_create.items.0.name").String())
				So(n["range"].(string), ShouldEqual, gjson.Get(sBody, "networks_to_create.items.0.range").String())
				So(n["router_name"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.name").String())
				So(n["router_type"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.type").String())
				So(n["router_ip"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.ip").String())
				So(n["client_name"], ShouldBeNil)
				So(n["datacenter_name"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.name").String())
				So(n["datacenter_password"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.password").String())
				So(n["datacenter_region"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.region").String())
				So(n["datacenter_type"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.type").String())
				So(n["datacenter_username"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.username").String())

				So(err, ShouldEqual, nil)
			})
		})
	})
}

func TestDeleteNetworks(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s, sBody := h.getService("./fixtures/service_real_workflow.json")

		Convey("When I get the message for a networks.delete event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("networks.delete", &s)
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				n := m.Components[0].(map[string]interface{})
				So(n["name"].(string), ShouldEqual, gjson.Get(sBody, "networks.items.0.name").String())
				So(n["range"].(string), ShouldEqual, gjson.Get(sBody, "networks.items.0.range").String())
				So(n["router_name"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.name").String())
				So(n["router_type"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.type").String())
				So(n["router_ip"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.ip").String())
				So(n["client_name"], ShouldBeNil)
				So(n["datacenter_name"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.name").String())
				So(n["datacenter_password"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.password").String())
				So(n["datacenter_region"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.region").String())
				So(n["datacenter_type"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.type").String())
				So(n["datacenter_username"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.username").String())

				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestCreateInstances(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s, sBody := h.getService("./fixtures/service_create_instances.json")

		Convey("When I get the message for a instances.create event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("instances.create", &s)
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 2)
				i := m.Components[0].(map[string]interface{})
				So(i["name"].(string), ShouldEqual, gjson.Get(sBody, "instances_to_create.items.0.name").String())
				So(i["type"].(string), ShouldEqual, gjson.Get(sBody, "instances_to_create.items.0.type").String())
				So(i["ip"].(string), ShouldEqual, gjson.Get(sBody, "instances_to_create.items.0.ip").String())
				So(i["datacenter_name"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.name").String())
				So(i["datacenter_password"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.password").String())
				So(i["datacenter_region"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.region").String())
				So(i["datacenter_type"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.type").String())
				So(i["datacenter_username"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.username").String())
				So(i["network_aws_id"].(string), ShouldEqual, gjson.Get(sBody, "networks.items.0.network_aws_id").String())
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestCreateNats(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s, sBody := h.getService("./fixtures/service_create_nats.json")

		Convey("When I get the message for a nats.create event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("nats.create", &s)
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)
			fmt.Println(body)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				n := m.Components[0].(map[string]interface{})
				So(n["name"].(string), ShouldEqual, gjson.Get(sBody, "nats_to_create.items.0.name").String())
				So(n["status"].(string), ShouldEqual, gjson.Get(sBody, "nats_to_create.items.0.status").String())
				So(n["router_name"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.name").String())
				So(n["router_type"], ShouldBeNil)
				So(n["router_ip"], ShouldBeNil)
				So(n["client_name"], ShouldBeNil)
				So(n["datacenter_name"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.name").String())
				So(n["datacenter_password"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.password").String())
				So(n["datacenter_region"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.region").String())
				So(n["datacenter_type"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.type").String())
				So(n["datacenter_username"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.username").String())
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
		s, sBody := h.getService("./fixtures/service_update_nats.json")

		Convey("When I get the message for a nats.update event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("nats.update", &s)
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)
			fmt.Println(body)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				n := m.Components[0].(map[string]interface{})

				So(n["name"].(string), ShouldEqual, gjson.Get(sBody, "nats_to_update.items.0.name").String())
				So(n["status"].(string), ShouldEqual, gjson.Get(sBody, "nats_to_update.items.0.status").String())
				So(n["router_name"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.name").String())
				So(n["router_type"], ShouldBeNil)
				So(n["router_ip"], ShouldBeNil)
				So(n["client_name"], ShouldBeNil)
				So(n["datacenter_name"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.name").String())
				So(n["datacenter_password"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.password").String())
				So(n["datacenter_region"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.region").String())
				So(n["datacenter_type"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.type").String())
				So(n["datacenter_username"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.username").String())
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
		s, sBody := h.getService("./fixtures/service_create_firewalls.json")

		Convey("When I get the message for a firewalls.create event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("firewalls.create", &s)
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				f := m.Components[0].(map[string]interface{})
				So(f["name"].(string), ShouldEqual, gjson.Get(sBody, "firewalls_to_create.items.0.name").String())
				So(f["status"].(string), ShouldEqual, gjson.Get(sBody, "firewalls_to_create.items.0.status").String())
				So(f["router_name"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.name").String())
				So(f["router_type"], ShouldBeNil)
				So(f["datacenter_name"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.name").String())
				So(f["datacenter_password"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.password").String())
				So(f["datacenter_region"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.region").String())
				So(f["datacenter_type"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.type").String())
				So(f["datacenter_username"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.username").String())
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
		s, sBody := h.getService("./fixtures/service_update_firewalls.json")

		Convey("When I get the message for a firewalls.update event", func() {
			mm := messageManager{}
			body, err := mm.preparePublishMessage("firewalls.update", &s)
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				So(m.Service, ShouldEqual, s.ID)
				So(len(m.Components), ShouldEqual, 1)
				f := m.Components[0].(map[string]interface{})
				So(f["name"].(string), ShouldEqual, gjson.Get(sBody, "firewalls_to_update.items.0.name").String())
				So(f["status"].(string), ShouldEqual, gjson.Get(sBody, "firewalls_to_update.items.0.status").String())
				So(f["router_name"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.name").String())
				So(f["router_type"], ShouldBeNil)
				So(f["datacenter_name"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.name").String())
				So(f["datacenter_password"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.password").String())
				So(f["datacenter_region"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.region").String())
				So(f["datacenter_type"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.type").String())
				So(f["datacenter_username"].(string), ShouldEqual, gjson.Get(sBody, "datacenters.items.0.username").String())
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
		s, sBody := h.getService("./fixtures/service_real_workflow.json")
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

					So(b["name"].(string), ShouldEqual, gjson.Get(sBody, "bootstraps_to_create.items.0.name").String())
					So(b["type"], ShouldEqual, "salt")
					So(b["payload"].(string), ShouldEqual, gjson.Get(sBody, "bootstraps_to_create.items.0.payload").String())
					So(b["target"].(string), ShouldEqual, gjson.Get(sBody, "bootstraps_to_create.items.0.target").String())
					So(b["status"].(string), ShouldEqual, gjson.Get(sBody, "bootstraps_to_create.items.0.status").String())
					So(b["user"], ShouldBeNil)
					So(b["password"], ShouldBeNil)
					So(b["service_endpoint"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.ip").String())

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
		s, sBody := h.getService("./fixtures/service_real_workflow.json")
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

				So(e["name"].(string), ShouldEqual, gjson.Get(sBody, "executions_to_create.items.0.name").String())
				So(e["type"], ShouldEqual, "salt")
				So(e["payload"].(string), ShouldEqual, gjson.Get(sBody, "executions_to_create.items.0.payload").String())
				So(e["target"].(string), ShouldEqual, gjson.Get(sBody, "executions_to_create.items.0.target").String())
				So(e["user"], ShouldBeNil)
				So(e["password"], ShouldBeNil)
				So(e["service_endpoint"].(string), ShouldEqual, gjson.Get(sBody, "routers.items.0.ip").String())
			})
		})
	})
}

func TestServiceDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s, _ := h.getService("./fixtures/service_real_workflow.json")
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
