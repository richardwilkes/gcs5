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

import "strings"

// Possible PrereqType values.
const (
	PrereqList PrereqType = iota
	AdvantagePrereq
	AttributePrereq
	ContainedQuantityPrereq
	ContainedWeightPrereq
	SkillPrereq
	SpellPrereq
)

type prereqTypeData struct {
	Key string
}

// PrereqType holds the type of a Feature.
type PrereqType uint8

var prereqTypeValues = []*prereqTypeData{
	{
		Key: "prereq_list",
	},
	{
		Key: "advantage_prereq",
	},
	{
		Key: "attribute_prereq",
	},
	{
		Key: "contained_quantity_prereq",
	},
	{
		Key: "contained_weight_prereq",
	},
	{
		Key: "skill_prereq",
	},
	{
		Key: "spell_prereq",
	},
}

// PrereqTypeFromString extracts a PrereqType from a key.
func PrereqTypeFromString(key string) PrereqType {
	for i, one := range prereqTypeValues {
		if strings.EqualFold(key, one.Key) {
			return PrereqType(i)
		}
	}
	return 0
}

// EnsureValid returns the first PrereqType if this PrereqType is not a known value.
func (p PrereqType) EnsureValid() PrereqType {
	if int(p) < len(prereqTypeValues) {
		return p
	}
	return 0
}

// Key returns the key used to represent this PrereqType.
func (p PrereqType) Key() string {
	return prereqTypeValues[p.EnsureValid()].Key
}
