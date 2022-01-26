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

// SpellMatchType holds the type of an attribute definition.
type SpellMatchType uint8

// SpellMatchTypeFromString extracts a SpellMatchType from a string.
func SpellMatchTypeFromString(str string) SpellMatchType {
	for one := AllColleges; one <= SpellName; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return AllColleges
}

// Key returns the key used to represent this SpellMatchType.
func (s SpellMatchType) Key() string {
	switch s {
	case CollegeName:
		return "college_name"
	case PowerSourceName:
		return "power_source_name"
	case SpellName:
		return "spell_name"
	default: // AllColleges
		return "all_colleges"
	}
}
