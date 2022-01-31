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

// Possible Type values.
const (
	AttributeBonus           = Type("attribute_bonus")
	ConditionalModifier      = Type("conditional_modifier")
	ContainedWeightReduction = Type("contained_weight_reduction")
	CostReduction            = Type("cost_reduction")
	DRBonus                  = Type("dr_bonus")
	ReactionBonus            = Type("reaction_bonus")
	SkillBonus               = Type("skill_bonus")
	SkillPointBonus          = Type("skill_point_bonus")
	SpellBonus               = Type("spell_bonus")
	SpellPointBonus          = Type("spell_point_bonus")
	WeaponDamageBonus        = Type("weapon_bonus")
)

// AllTypes is the complete set of Type values.
var AllTypes = []Type{
	AttributeBonus,
	ConditionalModifier,
	ContainedWeightReduction,
	CostReduction,
	DRBonus,
	ReactionBonus,
	SkillBonus,
	SkillPointBonus,
	SpellBonus,
	SpellPointBonus,
	WeaponDamageBonus,
}

// Type holds the type of a Feature.
type Type string

// EnsureValid ensures this is of a known value.
func (a Type) EnsureValid() Type {
	for _, one := range AllTypes {
		if one == a {
			return a
		}
	}
	return AllTypes[0]
}
