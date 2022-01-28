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

// Possible WeaponType values.
const (
	Melee WeaponType = iota
	Ranged
)

type weaponTypeData struct {
	Key    string
	String string
}

// WeaponType holds the type of an weapon definition.
type WeaponType uint8

var weaponTypeValues = []*weaponTypeData{
	{
		Key:    "melee_weapon",
		String: i18n.Text("Melee Weapon"),
	},
	{
		Key:    "ranged_weapon",
		String: i18n.Text("Ranged Weapon"),
	},
}

// WeaponTypeFromKey extracts a WeaponType from a key.
func WeaponTypeFromKey(key string) WeaponType {
	for i, one := range weaponTypeValues {
		if strings.EqualFold(key, one.Key) {
			return WeaponType(i)
		}
	}
	return 0
}

// EnsureValid returns the first WeaponType if this WeaponType is not a known value.
func (a WeaponType) EnsureValid() WeaponType {
	if int(a) < len(weaponTypeValues) {
		return a
	}
	return 0
}

// Key returns the key used to represent this WeaponType.
func (a WeaponType) Key() string {
	return weaponTypeValues[a.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (a WeaponType) String() string {
	return weaponTypeValues[a.EnsureValid()].String
}
