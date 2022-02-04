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
	AttributeBonusType           = Type("attribute_bonus")
	ConditionalModifierType      = Type("conditional_modifier")
	ContainedWeightReductionType = Type("contained_weight_reduction")
	CostReductionType            = Type("cost_reduction")
	DRBonusType                  = Type("dr_bonus")
	ReactionBonusType            = Type("reaction_bonus")
	SkillBonusType               = Type("skill_bonus")
	SkillPointBonusType          = Type("skill_point_bonus")
	SpellBonusType               = Type("spell_bonus")
	SpellPointBonusType          = Type("spell_point_bonus")
	WeaponDamageBonusType        = Type("weapon_bonus")
)

// AllTypes is the complete set of Type values.
var AllTypes = []Type{
	AttributeBonusType,
	ConditionalModifierType,
	ContainedWeightReductionType,
	CostReductionType,
	DRBonusType,
	ReactionBonusType,
	SkillBonusType,
	SkillPointBonusType,
	SpellBonusType,
	SpellPointBonusType,
	WeaponDamageBonusType,
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
