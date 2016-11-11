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
		var p Publisher
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
			items := p.UpdateTemplateVariables(x, s)

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
			items := p.UpdateTemplateVariables(x, si)

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

		Convey("When i try and template fields nested inside of another structure multiple levels deep", func() {
			x := comp["route53s"].(map[string]interface{})["items"].([]interface{})
			items := p.UpdateTemplateVariables(x, si)

			Convey("It should not have mapped fields where there was a result", func() {
				collection, ok := items[0].(map[string]interface{})
				So(ok, ShouldBeTrue)
				records, ok := collection["records"].([]interface{})
				So(ok, ShouldBeTrue)
				record, ok := records[0].(map[string]interface{})
				So(ok, ShouldBeTrue)
				fmt.Println(record)
				values, ok := record["values"].([]interface{})
				So(ok, ShouldBeTrue)
				So(len(values), ShouldEqual, 1)
				So(values[0].(string), ShouldEqual, "8.8.8.8")
			})
		})
	})
}

func TestPublisherCreateError(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s, _ := h.getService("./fixtures/service_components.json")

		Convey("When I get the message for a services.create.error event", func() {
			mm := MessageManager{}
			body, err := mm.preparePublishMessage("service.create.error", s)
			m := &service{}
			err = json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				id, _ := (*s)["id"]
				So(m.ID, ShouldEqual, id)
				So(m.Status, ShouldEqual, "errored")
				So(err, ShouldEqual, nil)

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
			mm := MessageManager{}
			message, err := mm.preparePublishMessage("test.message.invalid", s)

			Convey("Then I'll receive an empty body and an error", func() {
				So(message, ShouldEqual, "")
				So(err, ShouldNotEqual, nil)

			})
		})
	})
}

func TestCreateComponents(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s, sBody := h.getService("./fixtures/service_components.json")

		Convey("When I get the message for a components.create event", func() {
			mm := MessageManager{}
			body, err := mm.preparePublishMessage("components.create", s)
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				r := m.Components[0].(map[string]interface{})
				So(len(m.Components), ShouldEqual, 2)
				So(r["name"].(string), ShouldEqual, gjson.Get(sBody, "components_to_create.items.0.name").String())
				So(r["type"].(string), ShouldEqual, gjson.Get(sBody, "components_to_create.items.0.type").String())
				So(err, ShouldEqual, nil)
			})
		})
	})
}

func TestUpdateComponents(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s, sBody := h.getService("./fixtures/service_components.json")

		Convey("When I get the message for a components.create event", func() {
			mm := MessageManager{}
			body, err := mm.preparePublishMessage("components.update", s)
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				r := m.Components[0].(map[string]interface{})
				So(len(m.Components), ShouldEqual, 1)
				So(r["name"].(string), ShouldEqual, gjson.Get(sBody, "components_to_update.items.0.name").String())
				So(r["type"].(string), ShouldEqual, gjson.Get(sBody, "components_to_update.items.0.type").String())
				So(err, ShouldEqual, nil)
			})
		})
	})
}

func TestDeleteComponents(t *testing.T) {
	Convey("Given I have a valid service", t, func() {
		setup()

		p.load(natsClient)
		s, sBody := h.getService("./fixtures/service_components.json")

		Convey("When I get the message for a components.create event", func() {
			mm := MessageManager{}
			body, err := mm.preparePublishMessage("components.delete", s)
			m := &GenericComponentMsg{}
			json.Unmarshal([]byte(body), &m)

			Convey("Then I'll receive a valid json string", func() {
				r := m.Components[0].(map[string]interface{})
				So(len(m.Components), ShouldEqual, 1)
				So(r["name"].(string), ShouldEqual, gjson.Get(sBody, "components_to_delete.items.0.name").String())
				So(r["type"].(string), ShouldEqual, gjson.Get(sBody, "components_to_delete.items.0.type").String())
				So(err, ShouldEqual, nil)
			})
		})
	})
}
