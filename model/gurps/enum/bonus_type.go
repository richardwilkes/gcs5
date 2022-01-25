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

// Possible BonusType values.
const (
	AttributeBonusType BonusType = iota
	ConditionalModifierBonusType
	DRBonusType
	ReactionBonusType
	SkillBonusType
	SkillPointBonusType
	SpellBonusType
	SpellPointBonusType
	WeaponDamageBonusType
)

// BonusType holds the type of a Bonus.
type BonusType uint8

// BonusTypeFromString extracts a BonusType from a string.
func BonusTypeFromString(str string) BonusType {
	for one := AttributeBonusType; one <= WeaponDamageBonusType; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return AttributeBonusType
}

// Key returns the key used to represent this BonusType.
func (b BonusType) Key() string {
	switch b {
	case ConditionalModifierBonusType:
		return "conditional_modifier"
	case DRBonusType:
		return "dr_bonus"
	case ReactionBonusType:
		return "reaction_bonus"
	case SkillBonusType:
		return "skill_bonus"
	case SkillPointBonusType:
		return "skill_point_bonus"
	case SpellBonusType:
		return "spell_bonus"
	case SpellPointBonusType:
		return "spell_point_bonus"
	case WeaponDamageBonusType:
		return "weapon_bonus"
	default: // AttributeBonusType
		return "attribute_bonus"
	}
}
