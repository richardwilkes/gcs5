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

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible SkillSelectionType values.
const (
	SkillsWithNameSkillSelect SkillSelectionType = iota
	ThisWeaponSkillSelect
	WeaponsWithNameSkillSelect
)

// SkillSelectionType holds the type of an attribute definition.
type SkillSelectionType uint8

// SkillSelectionTypeFromString extracts a SkillSelectionType from a string.
func SkillSelectionTypeFromString(str string) SkillSelectionType {
	for one := SkillsWithNameSkillSelect; one <= WeaponsWithNameSkillSelect; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return SkillsWithNameSkillSelect
}

// Key returns the key used to represent this SkillSelectionType.
func (a SkillSelectionType) Key() string {
	switch a {
	case ThisWeaponSkillSelect:
		return "this_weapon"
	case WeaponsWithNameSkillSelect:
		return "weapons_with_name"
	default: // SkillsWithName
		return "skills_with_name"
	}
}

// String implements fmt.Stringer.
func (a SkillSelectionType) String() string {
	switch a {
	case ThisWeaponSkillSelect:
		return i18n.Text("to this weapon")
	case WeaponsWithNameSkillSelect:
		return i18n.Text("to weapons whose name")
	default: // SkillsWithName
		return i18n.Text("to skills whose name")
	}
}
