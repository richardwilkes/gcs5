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

import "strings"

// Possible Type values.
const (
	List Type = iota
	Advantage
	Attribute
	ContainedQuantity
	ContainedWeight
	Skill
	Spell
)

type typeData struct {
	Key string
}

// Type holds the type of a Prereq.
type Type uint8

var typeValues = []*typeData{
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

// TypeFromString extracts a Type from a key.
func TypeFromString(key string) Type {
	for i, one := range typeValues {
		if strings.EqualFold(key, one.Key) {
			return Type(i)
		}
	}
	return 0
}

// EnsureValid returns the first Type if this Type is not a known value.
func (t Type) EnsureValid() Type {
	if int(t) < len(typeValues) {
		return t
	}
	return 0
}

// Key returns the key used to represent this Type.
func (t Type) Key() string {
	return typeValues[t.EnsureValid()].Key
}
