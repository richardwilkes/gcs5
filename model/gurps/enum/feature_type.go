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

package enum

import "strings"

// Possible FeatureType values.
const (
	AttributeBonus FeatureType = iota
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

// FeatureType holds the type of a Feature.
type FeatureType uint8

// FeatureTypeFromString extracts a FeatureType from a string.
func FeatureTypeFromString(str string) FeatureType {
	for one := AttributeBonus; one <= WeaponDamageBonus; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return AttributeBonus
}

// Key returns the key used to represent this FeatureType.
func (b FeatureType) Key() string {
	switch b {
	case ConditionalModifierBonus:
		return "conditional_modifier"
	case ContainedWeightReduction:
		return "contained_weight_reduction"
	case CostReduction:
		return "cost_reduction"
	case DRBonus:
		return "dr_bonus"
	case ReactionBonus:
		return "reaction_bonus"
	case SkillBonus:
		return "skill_bonus"
	case SkillPointBonus:
		return "skill_point_bonus"
	case SpellBonus:
		return "spell_bonus"
	case SpellPointBonus:
		return "spell_point_bonus"
	case WeaponDamageBonus:
		return "weapon_bonus"
	default: // AttributeBonusType
		return "attribute_bonus"
	}
}
