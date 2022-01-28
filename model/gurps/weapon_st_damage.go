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

// Possible WeaponSTDamage values.
const (
	NoSTBasedDamage WeaponSTDamage = iota
	Thrust
	LeveledThrust
	Swing
	LeveledSwing
)

type weaponSTDamageData struct {
	Key    string
	String string
}

// WeaponSTDamage holds the type of strength dice to add to damage.
type WeaponSTDamage uint8

var weaponSTDamageValues = []*weaponSTDamageData{
	{
		Key:    "none",
		String: "",
	},
	{
		Key:    "thr",
		String: "thr",
	},
	{
		Key:    "thr_leveled",
		String: "thr " + i18n.Text("(leveled)"),
	},
	{
		Key:    "sw",
		String: "sw",
	},
	{
		Key:    "sw_leveled",
		String: "sw " + i18n.Text("(leveled)"),
	},
}

// WeaponSTDamageFromKey extracts a WeaponSTDamage from a key.
func WeaponSTDamageFromKey(key string) WeaponSTDamage {
	for i, one := range weaponSTDamageValues {
		if strings.EqualFold(key, one.Key) {
			return WeaponSTDamage(i)
		}
	}
	return 0
}

// EnsureValid returns the first WeaponSTDamage if this WeaponSTDamage is not a known value.
func (w WeaponSTDamage) EnsureValid() WeaponSTDamage {
	if int(w) < len(weaponSTDamageValues) {
		return w
	}
	return 0
}

// Key returns the key used to represent this WeaponSTDamage.
func (w WeaponSTDamage) Key() string {
	return weaponSTDamageValues[w.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (w WeaponSTDamage) String() string {
	return weaponSTDamageValues[w.EnsureValid()].String
}
