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
	Name SpellComparisonType = iota
	Category
	College
	CollegeCount
	AnySpell
)

type spellComparisonTypeData struct {
	Key                string
	UsesStringCriteria bool
}

// SpellComparisonType holds the type of an attribute definition.
type SpellComparisonType uint8

var spellComparisonTypeValues = []*spellComparisonTypeData{
	{
		Key:                "name",
		UsesStringCriteria: true,
	},
	{
		Key:                "category",
		UsesStringCriteria: true,
	},
	{
		Key:                "college",
		UsesStringCriteria: true,
	},
	{
		Key:                "college_count",
		UsesStringCriteria: false,
	},
	{
		Key:                "any",
		UsesStringCriteria: false,
	},
}

// SpellComparisonTypeFromString extracts a SpellComparisonType from a key.
func SpellComparisonTypeFromString(key string) SpellComparisonType {
	for i, one := range spellComparisonTypeValues {
		if strings.EqualFold(key, one.Key) {
			return SpellComparisonType(i)
		}
	}
	return 0
}

// EnsureValid returns the first SpellComparisonType if this SpellComparisonType is not a known value.
func (s SpellComparisonType) EnsureValid() SpellComparisonType {
	if int(s) < len(spellComparisonTypeValues) {
		return s
	}
	return 0
}

// Key returns the key used to represent this SpellComparisonType.
func (s SpellComparisonType) Key() string {
	return spellComparisonTypeValues[s.EnsureValid()].Key
}

// UsesStringCriteria returns true if the comparison uses a string value.
func (s SpellComparisonType) UsesStringCriteria() bool {
	return spellComparisonTypeValues[s.EnsureValid()].UsesStringCriteria
}
