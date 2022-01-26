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

// Possible SpellComparisonType values.
const (
	AnySpell SpellComparisonType = iota
	Name
	Category
	College
	CollegeCount
)

// SpellComparisonType holds the type of an attribute definition.
type SpellComparisonType uint8

// SpellComparisonTypeFromString extracts a SpellComparisonType from a string.
func SpellComparisonTypeFromString(str string) SpellComparisonType {
	for one := AnySpell; one <= CollegeCount; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return Name
}

// Key returns the key used to represent this SpellComparisonType.
func (a SpellComparisonType) Key() string {
	switch a {
	case AnySpell:
		return "any"
	case Category:
		return "category"
	case College:
		return "college"
	case CollegeCount:
		return "college_count"
	default: // Name
		return "name"
	}
}

// UsesStringCriteria returns true if the comparison uses a string value.
func (a SpellComparisonType) UsesStringCriteria() bool {
	return a == Name || a == Category || a == College
}
