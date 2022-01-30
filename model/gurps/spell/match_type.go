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

import (
	"strings"
)

// Possible MatchType values.
const (
	AllColleges MatchType = iota
	CollegeName
	PowerSource
	Spell
)

type matchTypeData struct {
	Key string
}

// MatchType holds the type of an attribute definition.
type MatchType uint8

var matchTypeValues = []*matchTypeData{
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

// MatchTypeFromString extracts a MatchType from a key.
func MatchTypeFromString(key string) MatchType {
	for i, one := range matchTypeValues {
		if strings.EqualFold(key, one.Key) {
			return MatchType(i)
		}
	}
	return 0
}

// EnsureValid returns the first MatchType if this MatchType is not a known value.
func (m MatchType) EnsureValid() MatchType {
	if int(m) < len(matchTypeValues) {
		return m
	}
	return 0
}

// Key returns the key used to represent this MatchType.
func (m MatchType) Key() string {
	return matchTypeValues[m.EnsureValid()].Key
}
