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
func (t Type) ToInlineJSON(encoder *encoding.JSONEncoder) {
	encoder.KeyedBool(mentalKey, t.Mental(), true)
	encoder.KeyedBool(physicalKey, t.Physical(), true)
	encoder.KeyedBool(socialKey, t.Social(), true)
	encoder.KeyedBool(exoticKey, t.Exotic(), true)
	encoder.KeyedBool(supernaturalKey, t.Supernatural(), true)
}

func (t Type) String() string {
	list := make([]string, 0, 5)
	if t.Mental() {
		list = append(list, i18n.Text("Mental"))
	}
	if t.Physical() {
		list = append(list, i18n.Text("Physical"))
	}
	if t.Social() {
		list = append(list, i18n.Text("Social"))
	}
	if t.Exotic() {
		list = append(list, i18n.Text("Exotic"))
	}
	if t.Supernatural() {
		list = append(list, i18n.Text("Supernatural"))
	}
	return strings.Join(list, ", ")
}

// Mental returns true if this Type has the Mental flag set.
func (t Type) Mental() bool {
	return t&Mental != 0
}

// Physical returns true if this Type has the Physical flag set.
func (t Type) Physical() bool {
	return t&Physical != 0
}

// Social returns true if this Type has the Social flag set.
func (t Type) Social() bool {
	return t&Social != 0
}

// Exotic returns true if this Type has the Exotic flag set.
func (t Type) Exotic() bool {
	return t&Exotic != 0
}

// Supernatural returns true if this Type has the Supernatural flag set.
func (t Type) Supernatural() bool {
	return t&Supernatural != 0
}
