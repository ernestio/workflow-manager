/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"testing"

	"github.com/tidwall/gjson"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRoutersCreateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/routers_create_done.json")
		s, _ := h.getService("./fixtures/service_real_workflow.json")
		s.save()

		Convey("When I try to get body for the mapped message routers.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("routers.create.done", body)
			b, _ := json.Marshal(s)
			sBody := string(b)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Routers.Items), ShouldEqual, 2)
				So(gjson.Get(sBody, "routers.items.1.ip").String(), ShouldEqual, "31.210.241.2")
				So(subject, ShouldEqual, "routers.create.done")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestCreateErrors(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/routers_create_done.json")
		s, _ := h.getService("./fixtures/service.json")
		s.save()

		Convey("When I try to get body for the non mapped message firewalls.create.error", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("firewalls.create.error", body)

			Convey("Then I'll receive the valid body", func() {
				So(s.Name, ShouldEqual, "test")
				So(subject, ShouldEqual, "to_error")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestNetworksCreateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/networks_create_done.json")
		s, _ := h.getService("./fixtures/service_real_workflow.json")
		s.save()

		Convey("When I try to get body for the mapped message networks.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("networks.create.done", body)
			b, _ := json.Marshal(s)
			sBody := string(b)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Networks.Items), ShouldEqual, 2)
				So(gjson.Get(sBody, "networks.items.1.name").String(), ShouldEqual, "network_test_2")
				So(gjson.Get(sBody, "networks.items.1.range").String(), ShouldEqual, "10.64.4.0/24")
				So(gjson.Get(sBody, "networks.items.1.netmask").String(), ShouldEqual, "netmask")
				So(gjson.Get(sBody, "networks.items.1.start_address").String(), ShouldEqual, "start")
				So(gjson.Get(sBody, "networks.items.1.end_address").String(), ShouldEqual, "end")
				So(gjson.Get(sBody, "networks.items.1.gateway").String(), ShouldEqual, "gateway")
				So(subject, ShouldEqual, "networks.create.done")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestNetworksDeleteDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/networks_delete_done.json")
		s, _ := h.getService("./fixtures/service_real_workflow.json")
		s.save()

		Convey("When I try to get body for the mapped message networks.delete.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("networks.delete.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.NetworksToDelete.Items), ShouldEqual, 0)
				So(subject, ShouldEqual, "networks.delete.done")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestInstancesCreateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/instances_create_done.json")
		s, _ := h.getService("./fixtures/service_real_workflow.json")
		s.save()

		Convey("When I try to get body for the mapped message instances.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("instances.create.done", body)
			b, _ := json.Marshal(s)
			sBody := string(b)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Instances.Items), ShouldEqual, 3)
				So(subject, ShouldEqual, "instances.create.done")
				So(err, ShouldEqual, nil)
				So(gjson.Get(sBody, "instances.items.2.status").String(), ShouldEqual, "completed")
				So(gjson.Get(sBody, "instances.items.2.name").String(), ShouldEqual, "test_instance_2")
				So(gjson.Get(sBody, "instances.items.2.ram").String(), ShouldEqual, "1024")
				So(gjson.Get(sBody, "instances.items.2.ip").String(), ShouldEqual, "10.64.4.101")
			})
		})
	})
}

func TestInstancesUpdateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/instances_create_done.json")
		s, _ := h.getService("./fixtures/service_real_workflow.json")
		s.save()

		Convey("When I try to get body for the mapped message instances.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("instances.update.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Instances.Items), ShouldEqual, 2)
				So(subject, ShouldEqual, "instances.update.done")
				So(err, ShouldEqual, nil)
				So(len(s.InstancesToUpdate.Items), ShouldEqual, 0)
			})
		})
	})
}

func TestFirewallsCreateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/firewalls_create_done.json")
		s, _ := h.getService("./fixtures/service_create_firewalls.json")
		s.save()

		Convey("When I try to get body for the mapped message firewalls.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("firewalls.create.done", body)
			b, _ := json.Marshal(s)
			sBody := string(b)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Firewalls.Items), ShouldEqual, 1)
				So(subject, ShouldEqual, "firewalls.create.done")
				So(err, ShouldEqual, nil)
				So(gjson.Get(sBody, "firewalls.items.0.status").String(), ShouldEqual, "completed")

			})
		})
	})
}

func TestFirewallsUpdateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/firewalls_create_done.json")
		s, _ := h.getService("./fixtures/service_real_workflow.json")
		s.save()

		Convey("When I try to get body for the mapped message firewalls.update.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("firewalls.update.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Firewalls.Items), ShouldEqual, 1)
				So(subject, ShouldEqual, "firewalls.update.done")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestNatsCreateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/nats_create_done.json")
		s, _ := h.getService("./fixtures/service_create_nats.json")
		s.save()

		Convey("When I try to get body for the mapped message nats.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("nats.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Nats.Items), ShouldEqual, 1)
				So(subject, ShouldEqual, "nats.create.done")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestNatUpdateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/nats_create_done.json")
		s, _ := h.getService("./fixtures/service_real_workflow.json")
		s.save()

		Convey("When I try to get body for the mapped message nats.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("nats.update.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Nats.Items), ShouldEqual, 1)
				So(subject, ShouldEqual, "nats.update.done")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestBootstrapsCreateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/executions_create_done.json")
		s, sBody := h.getService("./fixtures/service_real_workflow.json")
		s.Status = "bootstrapping"
		s.save()

		Convey("When I try to get body for the mapped message executions.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("bootstraps.create.done", body)
			Convey("Then I'll receive the valid body", func() {
				So(len(s.Bootstraps.Items), ShouldEqual, 2)
				So(gjson.Get(sBody, "bootstraps.items.0.reports.0.return_code").String(), ShouldEqual, "0")
				So(gjson.Get(sBody, "bootstraps.items.0.reports.0.stdout").String(), ShouldEqual, "")
				So(gjson.Get(sBody, "bootstraps.items.0.matched_instances.0").String(), ShouldEqual, "")
				So(subject, ShouldEqual, "bootstraps.create.done")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestExecutionsCreateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/executions_create_done.json")
		s, sBody := h.getService("./fixtures/service_real_workflow.json")
		s.Status = "running_executions"
		s.save()

		Convey("When I try to get body for the mapped message executions.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("executions.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Executions.Items), ShouldEqual, 2)
				So(gjson.Get(sBody, "executions.items.0.reports.0.return_code").String(), ShouldEqual, "0")
				So(gjson.Get(sBody, "executions.items.0.reports.0.stdout").String(), ShouldEqual, "")
				So(gjson.Get(sBody, "executions.items.0.matched_instances.0").String(), ShouldEqual, "")
				So(subject, ShouldEqual, "executions.create.done")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestExecutionsCreateError(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/executions_create_done.json")
		s, _ := h.getService("./fixtures/service.json")
		s.Status = "running_executions"
		s.save()

		Convey("When I try to get body for the mapped message executions.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("executions.create.error", body)
			So(err, ShouldBeNil)

			Convey("Then I'll receive the valid body", func() {
				So(s.Name, ShouldEqual, "test")
				So(s.Status, ShouldEqual, "pre-failed")
				So(subject, ShouldEqual, "to_error")
				So(err, ShouldEqual, nil)
			})
		})
	})
}
