/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var h = testHelper{}

func TestWithInvalidTransition(t *testing.T) {
	t.Parallel()
	Convey("Given a valid service input", t, func() {
		p.load(natsClient)
		s, _ := h.getService("./fixtures/service.json")

		Convey("When a message with an unexisting transition is received", func() {
			subject, service, err := h.manage("hello", s)

			Convey("Then should return an error", func() {
				So(subject, ShouldEqual, "")
				So(service, ShouldEqual, nil)
				So(err, ShouldNotEqual, nil)
			})
		})
	})
}

func TestWithValidTransitionButNotRelativeToCurrentStatus(t *testing.T) {
	t.Parallel()
	Convey("Given a valid service input", t, func() {
		p.load(natsClient)
		s, _ := h.getService("./fixtures/service.json")

		Convey("When a message with an existing transition is received", func() {
			subject, service, err := h.manage("to_done", s)

			Convey("Then should return an error", func() {
				So(subject, ShouldEqual, "")
				So(service, ShouldEqual, nil)
				So(err, ShouldNotEqual, nil)
			})
		})
	})
}

func TestWithValidTransitionAndStatus(t *testing.T) {
	t.Parallel()
	Convey("Given a valid service input", t, func() {
		p.load(natsClient)
		s, _ := h.getService("./fixtures/service.json")
		s.Status = "created"

		Convey("When a message with an existing transition is received", func() {
			subject, service, err := h.manage("start", s)

			Convey("Then should return an error", func() {
				So(subject, ShouldEqual, "to_in_progress")
				So(service.Status, ShouldEqual, "started")
				So(err, ShouldEqual, nil)
			})
		})
	})
}

func TestOnStartingStatus(t *testing.T) {
	t.Parallel()
	Convey("Given a valid service input", t, func() {
		p.load(natsClient)
		s, _ := h.getService("./fixtures/service.json")

		Convey("When a message with an existing transition is received and not set status", func() {
			subject, service, err := h.manage("start", s)

			Convey("Then should return an error", func() {
				So(subject, ShouldEqual, "to_in_progress")
				So(service.Status, ShouldEqual, "started")
				So(err, ShouldEqual, nil)
			})
		})
	})
}
func TestOnFinalStatus(t *testing.T) {
	t.Parallel()
	Convey("Given a valid service input", t, func() {
		p.load(natsClient)
		s, _ := h.getService("./fixtures/service.json")
		s.Status = "uat"

		Convey("When a message with an existing transition is received and not set status", func() {
			subject, service, err := h.manage("to_done", s)

			Convey("Then should return an error", func() {
				So(subject, ShouldEqual, "")
				So(service.Status, ShouldEqual, "done")
				So(err, ShouldEqual, nil)
			})
		})
	})
}

func TestOnEntryPoint(t *testing.T) {
	t.Parallel()
	Convey("Given a valid service input", t, func() {
		p.load(natsClient)
		s, _ := h.getService("./fixtures/service.json")

		Convey("When a message with an existing transition is received", func() {
			subject, service, err := h.manage("start", s)

			Convey("Then should return an error", func() {
				So(subject, ShouldEqual, "to_in_progress")
				So(service.Status, ShouldEqual, "started")
				So(err, ShouldEqual, nil)
			})
		})
	})
}
