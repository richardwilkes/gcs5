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

package gurps

import (
	"strings"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/i18n"
)

// Masks for the various AdvantageTypeBits.
const (
	MentalTypeMask AdvantageTypeBits = 1 << iota
	PhysicalTypeMask
	SocialTypeMask
	ExoticTypeMask
	SupernaturalTypeMask
)

const (
	advantageMentalKey       = "mental"
	advantagePhysicalKey     = "physical"
	advantageSocialKey       = "social"
	advantageExoticKey       = "exotic"
	advantageSupernaturalKey = "supernatural"
)

// AdvantageTypeBits holds the various type flags for an Advantage.
type AdvantageTypeBits uint8

// AdvantageTypeBitsFromJSON loads an AdvantageTypeBits from JSON.
func AdvantageTypeBitsFromJSON(data map[string]interface{}) AdvantageTypeBits {
	var bits AdvantageTypeBits
	if encoding.Bool(data[advantageMentalKey]) {
		bits |= MentalTypeMask
	}
	if encoding.Bool(data[advantagePhysicalKey]) {
		bits |= PhysicalTypeMask
	}
	if encoding.Bool(data[advantageSocialKey]) {
		bits |= SocialTypeMask
	}
	if encoding.Bool(data[advantageExoticKey]) {
		bits |= ExoticTypeMask
	}
	if encoding.Bool(data[advantageSupernaturalKey]) {
		bits |= SupernaturalTypeMask
	}
	return bits
}

// ToInlineJSON emits the AdvantageTypeBits into JSON.
func (a AdvantageTypeBits) ToInlineJSON(encoder *encoding.JSONEncoder) {
	encoder.KeyedBool(advantageMentalKey, a&MentalTypeMask != 0, true)
	encoder.KeyedBool(advantagePhysicalKey, a&PhysicalTypeMask != 0, true)
	encoder.KeyedBool(advantageSocialKey, a&SocialTypeMask != 0, true)
	encoder.KeyedBool(advantageExoticKey, a&ExoticTypeMask != 0, true)
	encoder.KeyedBool(advantageSupernaturalKey, a&SupernaturalTypeMask != 0, true)
}

func (a AdvantageTypeBits) String() string {
	list := make([]string, 0, 5)
	if a&MentalTypeMask != 0 {
		list = append(list, i18n.Text("Mental"))
	}
	if a&PhysicalTypeMask != 0 {
		list = append(list, i18n.Text("Physical"))
	}
	if a&SocialTypeMask != 0 {
		list = append(list, i18n.Text("Social"))
	}
	if a&ExoticTypeMask != 0 {
		list = append(list, i18n.Text("Exotic"))
	}
	if a&SupernaturalTypeMask != 0 {
		list = append(list, i18n.Text("Supernatural"))
	}
	return strings.Join(list, ", ")
}
