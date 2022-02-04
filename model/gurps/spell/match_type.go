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

package spell

// Possible MatchType values.
const (
	AllColleges = MatchType("all_colleges")
	CollegeName = MatchType("college_name")
	PowerSource = MatchType("power_source_name")
	Spell       = MatchType("spell_name")
)

// AllMatchTypes is the complete set of MatchType values.
var AllMatchTypes = []MatchType{
	AllColleges,
	CollegeName,
	PowerSource,
	Spell,
}

// MatchType holds the type of an attribute definition.
type MatchType string

// EnsureValid ensures this is of a known value.
func (c MatchType) EnsureValid() MatchType {
	for _, one := range AllMatchTypes {
		if one == c {
			return c
		}
	}
	return AllMatchTypes[0]
}
