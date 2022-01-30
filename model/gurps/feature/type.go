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

package feature

import "strings"

// Possible Type values.
const (
	AttributeBonus Type = iota
	ConditionalModifierBonus
	ContainedWeightReduction
	CostReduction
	DRBonus
	ReactionBonus
	SkillBonus
	SkillPointBonus
	SpellBonus
	SpellPointBonus
	WeaponDamageBonus
)

type typeData struct {
	Key string
}

// Type holds the type of a Feature.
type Type uint8

var typeValues = []*typeData{
	{
		Key: "attribute_bonus",
	},
	{
		Key: "conditional_modifier",
	},
	{
		Key: "contained_weight_reduction",
	},
	{
		Key: "cost_reduction",
	},
	{
		Key: "dr_bonus",
	},
	{
		Key: "reaction_bonus",
	},
	{
		Key: "skill_bonus",
	},
	{
		Key: "skill_point_bonus",
	},
	{
		Key: "spell_bonus",
	},
	{
		Key: "spell_point_bonus",
	},
	{
		Key: "weapon_bonus",
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
