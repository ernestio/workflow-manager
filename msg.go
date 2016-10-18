/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

// GenericComponentCreate : Message to create instances
type GenericComponentMsg struct {
	Service              string        `json:"service"`
	Components           []interface{} `json:"components"`
	Status               string        `json:"status"`
	ErrorCode            string        `json:"error_code"`
	ErrorMessage         string        `json:"error_message"`
	SequentialProcessing bool          `json:"sequential_processing"`
}
