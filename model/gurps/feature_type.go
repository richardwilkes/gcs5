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

type featureTypeData struct {
	Key string
}

// FeatureType holds the type of a Feature.
type FeatureType uint8

var featureTypeValues = []*featureTypeData{
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

// FeatureTypeFromString extracts a FeatureType from a key.
func FeatureTypeFromString(key string) FeatureType {
	for i, one := range featureTypeValues {
		if strings.EqualFold(key, one.Key) {
			return FeatureType(i)
		}
	}
	return 0
}

// EnsureValid returns the first FeatureType if this FeatureType is not a known value.
func (f FeatureType) EnsureValid() FeatureType {
	if int(f) < len(featureTypeValues) {
		return f
	}
	return 0
}

// Key returns the key used to represent this FeatureType.
func (f FeatureType) Key() string {
	return featureTypeValues[f.EnsureValid()].Key
}
