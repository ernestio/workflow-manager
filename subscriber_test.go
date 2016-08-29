/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubscriberMappedMessage(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := []byte("{\"service\":\"1\"}")

		Convey("When I try to get body for the mapped message", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("test.message", body)

			Convey("Then I'll receive the valid body", func() {
				So(s.Name, ShouldEqual, "hello world from subscriber!")
				So(subject, ShouldEqual, "test.message")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestSubscriberUnMappedMessage(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := []byte("")

		Convey("When I try to get body for the unmapped message", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("test.message.invalid", body)

			Convey("Then I'll receive an empty body and an error", func() {
				So(s, ShouldEqual, nil)
				So(subject, ShouldEqual, "")
				So(err, ShouldNotEqual, nil)

			})
		})
	})
}

func TestRoutersCreateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/routers_create_done.json")
		s := h.getService("./fixtures/service.json")
		s.save()

		Convey("When I try to get body for the mapped message routers.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("routers.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Routers.Items), ShouldEqual, 1)
				So(s.Routers.Items[0].IP, ShouldEqual, "31.210.241.2")
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
		s := h.getService("./fixtures/service.json")
		s.save()

		Convey("When I try to get body for the mapped message networks.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("networks.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Routers.Items), ShouldEqual, 1)
				So(s.Networks.Items[0].Name, ShouldEqual, "network_test")
				So(s.Networks.Items[0].Range, ShouldEqual, "10.64.4.0/24")
				So(s.Networks.Items[0].Netmask, ShouldEqual, "netmask")
				So(s.Networks.Items[0].StartAddress, ShouldEqual, "start")
				So(s.Networks.Items[0].EndAddress, ShouldEqual, "end")
				So(s.Networks.Items[0].Gateway, ShouldEqual, "gateway")
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
		s := h.getService("./fixtures/service.json")
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
		s := h.getService("./fixtures/service.json")
		s.save()

		Convey("When I try to get body for the mapped message instances.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("instances.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Instances.Items), ShouldEqual, 2)
				So(subject, ShouldEqual, "instances.create.done")
				So(err, ShouldEqual, nil)
				So(s.Instances.Items[0].Status, ShouldEqual, "completed")
			})
		})
	})
}

func TestInstancesUpdateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/instances_create_done.json")
		s := h.getService("./fixtures/service.json")
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
		s := h.getService("./fixtures/service.json")
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
		s := h.getService("./fixtures/service.json")
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
		s := h.getService("./fixtures/service.json")
		s.Status = "bootstrapping"
		s.save()

		Convey("When I try to get body for the mapped message executions.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("executions.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Bootstraps.Items), ShouldEqual, 1)
				So(s.Bootstraps.Items[0].Reports[0].ReturnCode, ShouldEqual, 0)
				So(s.Bootstraps.Items[0].Reports[0].Stdout, ShouldEqual, "")
				So(s.Bootstraps.Items[0].MatchedInstances[0], ShouldEqual, "test")
				So(subject, ShouldEqual, "executions.create.done")
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
		s := h.getService("./fixtures/service.json")
		s.Status = "running_executions"
		s.save()

		Convey("When I try to get body for the mapped message executions.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("executions.create.done", body)

			Convey("Then I'll receive the valid body", func() {
				So(len(s.Executions.Items), ShouldEqual, 1)
				So(s.Executions.Items[0].Reports[0].ReturnCode, ShouldEqual, 0)
				So(s.Executions.Items[0].Reports[0].Stdout, ShouldEqual, "")
				So(s.Executions.Items[0].MatchedInstances[0], ShouldEqual, "test")
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

			Convey("Then I'll receive the valid body", func() {
				So(s.Name, ShouldEqual, "test")
				So(s.Status, ShouldEqual, "pre-failed")
				So(subject, ShouldEqual, "to_error")
				So(err, ShouldEqual, nil)
			})
		})
	})
}
