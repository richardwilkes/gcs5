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

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible WeaponSelectionType values.
const (
	WeaponsWithRequiredSkillWeaponSelect WeaponSelectionType = iota
	ThisWeaponWeaponSelect
	WeaponsWithNameWeaponSelect
)

// WeaponSelectionType holds the type of an attribute definition.
type WeaponSelectionType uint8

// WeaponSelectionTypeFromString extracts a WeaponSelectionType from a string.
func WeaponSelectionTypeFromString(str string) WeaponSelectionType {
	for one := WeaponsWithRequiredSkillWeaponSelect; one <= WeaponsWithNameWeaponSelect; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return WeaponsWithRequiredSkillWeaponSelect
}

// Key returns the key used to represent this WeaponSelectionType.
func (a WeaponSelectionType) Key() string {
	switch a {
	case ThisWeaponWeaponSelect:
		return "this_weapon"
	case WeaponsWithNameWeaponSelect:
		return "weapons_with_name"
	default: // WeaponsWithRequiredSkillWeaponSelect
		return "weapons_with_required_skill"
	}
}

// String implements fmt.Stringer.
func (a WeaponSelectionType) String() string {
	switch a {
	case ThisWeaponWeaponSelect:
		return i18n.Text("to this weapon")
	case WeaponsWithNameWeaponSelect:
		return i18n.Text("to weapons whose name")
	default: // WeaponsWithRequiredSkillWeaponSelect
		return i18n.Text("to weapons whose required skill name")
	}
}
