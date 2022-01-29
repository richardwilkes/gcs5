/*
 * Copyright ©1998-2022 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package advantage

import (
	"strings"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/i18n"
)

// Masks for the various Type bits.
const (
	Mental Type = 1 << iota
	Physical
	Social
	Exotic
	Supernatural
)

const (
	mentalKey       = "mental"
	physicalKey     = "physical"
	socialKey       = "social"
	exoticKey       = "exotic"
	supernaturalKey = "supernatural"
)

// Type holds the various type bits for an Advantage.
type Type uint8

// TypeFromJSON loads a Type from JSON.
func TypeFromJSON(data map[string]interface{}) Type {
	var bits Type
	if encoding.Bool(data[mentalKey]) {
		bits |= Mental
	}
	if encoding.Bool(data[physicalKey]) {
		bits |= Physical
	}
	if encoding.Bool(data[socialKey]) {
		bits |= Social
	}
	if encoding.Bool(data[exoticKey]) {
		bits |= Exotic
	}
	if encoding.Bool(data[supernaturalKey]) {
		bits |= Supernatural
	}
	return bits
}

// ToInlineJSON emits the Type into JSON.
func (a Type) ToInlineJSON(encoder *encoding.JSONEncoder) {
	encoder.KeyedBool(mentalKey, a.Mental(), true)
	encoder.KeyedBool(physicalKey, a.Physical(), true)
	encoder.KeyedBool(socialKey, a.Social(), true)
	encoder.KeyedBool(exoticKey, a.Exotic(), true)
	encoder.KeyedBool(supernaturalKey, a.Supernatural(), true)
}

func (a Type) String() string {
	list := make([]string, 0, 5)
	if a.Mental() {
		list = append(list, i18n.Text("Mental"))
	}
	if a.Physical() {
		list = append(list, i18n.Text("Physical"))
	}
	if a.Social() {
		list = append(list, i18n.Text("Social"))
	}
	if a.Exotic() {
		list = append(list, i18n.Text("Exotic"))
	}
	if a.Supernatural() {
		list = append(list, i18n.Text("Supernatural"))
	}
	return strings.Join(list, ", ")
}

// Mental returns true if this Type has the Mental flag set.
func (a Type) Mental() bool {
	return a&Mental != 0
}

// Physical returns true if this Type has the Physical flag set.
func (a Type) Physical() bool {
	return a&Physical != 0
}

// Social returns true if this Type has the Social flag set.
func (a Type) Social() bool {
	return a&Social != 0
}

// Exotic returns true if this Type has the Exotic flag set.
func (a Type) Exotic() bool {
	return a&Exotic != 0
}

// Supernatural returns true if this Type has the Supernatural flag set.
func (a Type) Supernatural() bool {
	return a&Supernatural != 0
}
