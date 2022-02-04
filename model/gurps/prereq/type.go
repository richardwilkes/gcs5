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

package prereq

// Possible Type values.
const (
	List              = Type("prereq_list")
	Advantage         = Type("advantage_prereq")
	Attribute         = Type("attribute_prereq")
	ContainedQuantity = Type("contained_quantity_prereq")
	ContainedWeight   = Type("contained_weight_prereq")
	Skill             = Type("skill_prereq")
	Spell             = Type("spell_prereq")
)

// AllTypes is the complete set of Type values.
var AllTypes = []Type{
	List,
	Advantage,
	Attribute,
	ContainedQuantity,
	ContainedWeight,
	Skill,
	Spell,
}

// Type holds the type of a Prereq.
type Type string

// EnsureValid ensures this is of a known value.
func (t Type) EnsureValid() Type {
	for _, one := range AllTypes {
		if one == t {
			return t
		}
	}
	return AllTypes[0]
}
