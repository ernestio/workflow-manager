/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRoutersCreateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/routers_create_done.json")
		s := h.getService("./fixtures/service_real_workflow.json")
		s.save()

		Convey("When I try to get body for the mapped message routers.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("routers.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Routers.Items), ShouldEqual, 2)
				router := s.Routers.Items[len(s.Routers.Items)-1]
				So(router.IP, ShouldEqual, "31.210.241.2")
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
		s := h.getService("./fixtures/service.json")
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
		s := h.getService("./fixtures/service_real_workflow.json")
		s.save()

		Convey("When I try to get body for the mapped message networks.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("networks.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Networks.Items), ShouldEqual, 2)
				network := s.Networks.Items[len(s.Networks.Items)-1]
				So(network.Name, ShouldEqual, "network_test_2")
				So(network.Range, ShouldEqual, "10.64.4.0/24")
				So(network.Netmask, ShouldEqual, "netmask")
				So(network.StartAddress, ShouldEqual, "start")
				So(network.EndAddress, ShouldEqual, "end")
				So(network.Gateway, ShouldEqual, "gateway")
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
		s := h.getService("./fixtures/service_real_workflow.json")
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
		s := h.getService("./fixtures/service_real_workflow.json")
		s.save()

		Convey("When I try to get body for the mapped message instances.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("instances.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Instances.Items), ShouldEqual, 3)
				So(subject, ShouldEqual, "instances.create.done")
				So(err, ShouldEqual, nil)
				lastInstance := s.Instances.Items[len(s.Instances.Items)-1]
				So(lastInstance.Status, ShouldEqual, "completed")
				So(lastInstance.Name, ShouldEqual, "test_instance_2")
				So(lastInstance.RAM, ShouldEqual, 1024)
				So(lastInstance.IP, ShouldEqual, "10.64.4.101")
			})
		})
	})
}

func TestInstancesUpdateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/instances_create_done.json")
		s := h.getService("./fixtures/service_real_workflow.json")
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
		s := h.getService("./fixtures/service_create_firewalls.json")
		s.save()

		Convey("When I try to get body for the mapped message firewalls.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("firewalls.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Firewalls.Items), ShouldEqual, 1)
				So(subject, ShouldEqual, "firewalls.create.done")
				So(err, ShouldEqual, nil)
				So(s.Firewalls.Items[0].Status, ShouldEqual, "completed")

			})
		})
	})
}

func TestFirewallsUpdateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/firewalls_create_done.json")
		s := h.getService("./fixtures/service_real_workflow.json")
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
		s := h.getService("./fixtures/service_create_nats.json")
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
		s := h.getService("./fixtures/service_real_workflow.json")
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
		s := h.getService("./fixtures/service_real_workflow.json")
		s.Status = "bootstrapping"
		s.save()

		Convey("When I try to get body for the mapped message executions.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("bootstraps.create.done", body)
			Convey("Then I'll receive the valid body", func() {
				So(len(s.Bootstraps.Items), ShouldEqual, 2)
				b := s.Bootstraps.Items[len(s.Bootstraps.Items)-1]
				So(b.Reports[0].ReturnCode, ShouldEqual, 0)
				So(b.Reports[0].Stdout, ShouldEqual, "test")
				So(b.MatchedInstances[0], ShouldEqual, "test")
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
		s := h.getService("./fixtures/service_real_workflow.json")
		s.Status = "running_executions"
		s.save()

		Convey("When I try to get body for the mapped message executions.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("executions.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Executions.Items), ShouldEqual, 2)
				e := s.Executions.Items[len(s.Executions.Items)-1]
				So(e.Reports[0].ReturnCode, ShouldEqual, 0)
				So(e.Reports[0].Stdout, ShouldEqual, "test")
				So(e.MatchedInstances[0], ShouldEqual, "test")
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
		s := h.getService("./fixtures/service.json")
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
