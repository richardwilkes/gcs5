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

type weaponSelectionTypeData struct {
	Key    string
	String string
}

// WeaponSelectionType holds the type of an attribute definition.
type WeaponSelectionType uint8

var weaponSelectionTypeValues = []*weaponSelectionTypeData{
	{
		Key:    "weapons_with_required_skill",
		String: i18n.Text("to weapons whose required skill name"),
	},
	{
		Key:    "this_weapon",
		String: i18n.Text("to this weapon"),
	},
	{
		Key:    "weapons_with_name",
		String: i18n.Text("to weapons whose name"),
	},
}

// WeaponSelectionTypeFromString extracts a WeaponSelectionType from a key.
func WeaponSelectionTypeFromString(key string) WeaponSelectionType {
	for i, one := range weaponSelectionTypeValues {
		if strings.EqualFold(key, one.Key) {
			return WeaponSelectionType(i)
		}
	}
	return 0
}

// EnsureValid returns the first WeaponSelectionType if this WeaponSelectionType is not a known value.
func (w WeaponSelectionType) EnsureValid() WeaponSelectionType {
	if int(w) < len(weaponSelectionTypeValues) {
		return w
	}
	return 0
}

// Key returns the key used to represent this WeaponSelectionType.
func (w WeaponSelectionType) Key() string {
	return weaponSelectionTypeValues[w.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (w WeaponSelectionType) String() string {
	return weaponSelectionTypeValues[w.EnsureValid()].String
}
