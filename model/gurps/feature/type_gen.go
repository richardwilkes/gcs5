// Code generated from "enum.go.tmpl" - DO NOT EDIT.

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

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible values.
const (
	AttributeBonusType Type = iota
	ConditionalModifierType
	ContainedWeightReductionType
	CostReductionType
	DRBonusType
	ReactionBonusType
	SkillBonusType
	SkillPointBonusType
	SpellBonusType
	SpellPointBonusType
	WeaponBonusType
	LastType = WeaponBonusType
)

var (
	// AllType holds all possible values.
	AllType = []Type{
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
		WeaponBonusType,
	}
	typeData = []struct {
		key    string
		string string
	}{
		{
			key:    "attribute_bonus",
			string: i18n.Text("Attribute Bonus"),
		},
		{
			key:    "conditional_modifier",
			string: i18n.Text("Conditional Modifier"),
		},
		{
			key:    "contained_weight_reduction",
			string: i18n.Text("Contained Weight Reduction"),
		},
		{
			key:    "cost_reduction",
			string: i18n.Text("Cost Reduction"),
		},
		{
			key:    "dr_bonus",
			string: i18n.Text("DR Bonus"),
		},
		{
			key:    "reaction_bonus",
			string: i18n.Text("Reaction Bonus"),
		},
		{
			key:    "skill_bonus",
			string: i18n.Text("Skill Bonus"),
		},
		{
			key:    "skill_point_bonus",
			string: i18n.Text("Skill Point Bonus"),
		},
		{
			key:    "spell_bonus",
			string: i18n.Text("Spell Bonus"),
		},
		{
			key:    "spell_point_bonus",
			string: i18n.Text("Spell Point Bonus"),
		},
		{
			key:    "weapon_bonus",
			string: i18n.Text("Weapon Bonus"),
		},
	}
)

// Type holds the type of a Feature.
type Type byte

// EnsureValid ensures this is of a known value.
func (enum Type) EnsureValid() Type {
	if enum <= LastType {
		return enum
	}
	return 0
}

// Key returns the key used in serialization.
func (enum Type) Key() string {
	return typeData[enum.EnsureValid()].key
}

// String implements fmt.Stringer.
func (enum Type) String() string {
	return typeData[enum.EnsureValid()].string
}

// ExtractType extracts the value from a string.
func ExtractType(str string) Type {
	for i, one := range typeData {
		if strings.EqualFold(one.key, str) {
			return Type(i)
		}
	}
	return 0
}

// MarshalText implements the encoding.TextMarshaler interface.
func (enum Type) MarshalText() (text []byte, err error) {
	return []byte(enum.Key()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (enum *Type) UnmarshalText(text []byte) error {
	*enum = ExtractType(string(text))
	return nil
}
