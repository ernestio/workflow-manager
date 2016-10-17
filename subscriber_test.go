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

func TestComponentsCreateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/components_create_done.json")
		s, _ := h.getService("./fixtures/service_components.json")
		SaveService(s)

		Convey("When I try to get body for the mapped message components.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("components.create.done", body)
			b, _ := json.Marshal(s)
			sBody := string(b)

			Convey("Then I'll receive the valid body", func() {
				So(len(gjson.Get(sBody, "components.items").Array()), ShouldEqual, 2)
				So(len(gjson.Get(sBody, "components_to_create.items").Array()), ShouldEqual, 0)
				So(gjson.Get(sBody, "components.items.0.name").String(), ShouldEqual, "existing")
				So(gjson.Get(sBody, "components.items.0.field").String(), ShouldEqual, "existing")
				So(gjson.Get(sBody, "components.items.1.name").String(), ShouldEqual, "created")
				So(gjson.Get(sBody, "components.items.1.field").String(), ShouldEqual, "created")
				So(subject, ShouldEqual, "components.create.done")
				So(err, ShouldEqual, nil)

			})
		})
	})
}

func TestComponentsUpdateDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/components_update_done.json")
		s, _ := h.getService("./fixtures/service_components.json")
		(*s)["status"] = "updating_components"
		(*s)["components"] = (*s)["components_to_create"]
		SaveService(s)

		Convey("When I try to get body for the mapped message components.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("components.update.done", body)
			b, _ := json.Marshal(s)
			sBody := string(b)

			Convey("Then I'll receive the valid body", func() {
				So(len(gjson.Get(sBody, "components.items").Array()), ShouldEqual, 2)
				So(len(gjson.Get(sBody, "components_to_create.items").Array()), ShouldEqual, 2)
				So(len(gjson.Get(sBody, "components_to_update.items").Array()), ShouldEqual, 0)
				So(gjson.Get(sBody, "components.items.0.name").String(), ShouldEqual, "added")
				So(gjson.Get(sBody, "components.items.0.field").String(), ShouldEqual, "created")
				So(gjson.Get(sBody, "components.items.1.name").String(), ShouldEqual, "updated")
				So(gjson.Get(sBody, "components.items.1.field").String(), ShouldEqual, "updated")
				So(subject, ShouldEqual, "components.update.done")
				So(err, ShouldEqual, nil)
			})
		})
	})
}

func TestComponentsDeleteDone(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		body := h.getFixture("./fixtures/components_delete_done.json")
		s, _ := h.getService("./fixtures/service_components.json")
		(*s)["status"] = "deleting_components"
		SaveService(s)

		Convey("When I try to get body for the mapped message components.create.done", func() {
			mm := messageManager{}
			s, subject, err := mm.getServiceFromMessage("components.delete.done", body)
			b, _ := json.Marshal(s)
			sBody := string(b)

			Convey("Then I'll receive the valid body", func() {
				So(len(gjson.Get(sBody, "components.items").Array()), ShouldEqual, 0)
				So(len(gjson.Get(sBody, "components_to_create.items").Array()), ShouldEqual, 2)
				So(len(gjson.Get(sBody, "components_to_update.items").Array()), ShouldEqual, 1)
				So(len(gjson.Get(sBody, "components_to_delete.items").Array()), ShouldEqual, 0)
				So(subject, ShouldEqual, "components.delete.done")
				So(err, ShouldEqual, nil)
			})
		})
	})
}
