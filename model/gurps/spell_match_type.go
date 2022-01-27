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
)

// Possible SpellMatchType values.
const (
	AllColleges SpellMatchType = iota
	CollegeName
	PowerSourceName
	SpellName
)

type spellMatchTypeData struct {
	Key string
}

// SpellMatchType holds the type of an attribute definition.
type SpellMatchType uint8

var spellMatchTypeValues = []*spellMatchTypeData{
	{
		Key: "all_colleges",
	},
	{
		Key: "college_name",
	},
	{
		Key: "power_source_name",
	},
	{
		Key: "spell_name",
	},
}

// SpellMatchTypeFromString extracts a SpellMatchType from a key.
func SpellMatchTypeFromString(key string) SpellMatchType {
	for i, one := range spellMatchTypeValues {
		if strings.EqualFold(key, one.Key) {
			return SpellMatchType(i)
		}
	}
	return 0
}

// EnsureValid returns the first SpellMatchType if this SpellMatchType is not a known value.
func (s SpellMatchType) EnsureValid() SpellMatchType {
	if int(s) < len(spellMatchTypeValues) {
		return s
	}
	return 0
}

// Key returns the key used to represent this SpellMatchType.
func (s SpellMatchType) Key() string {
	return attributeTypeValues[s.EnsureValid()].Key
}
