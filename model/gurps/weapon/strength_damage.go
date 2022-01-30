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

package weapon

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible StrengthDamage values.
const (
	None StrengthDamage = iota
	Thrust
	LeveledThrust
	Swing
	LeveledSwing
)

type strengthDamageData struct {
	Key    string
	String string
}

// StrengthDamage holds the type of strength dice to add to damage.
type StrengthDamage uint8

var strengthDamageValues = []*strengthDamageData{
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

// StrengthDamageFromKey extracts a StrengthDamage from a key.
func StrengthDamageFromKey(key string) StrengthDamage {
	for i, one := range strengthDamageValues {
		if strings.EqualFold(key, one.Key) {
			return StrengthDamage(i)
		}
	}
	return 0
}

// EnsureValid returns the first StrengthDamage if this StrengthDamage is not a known value.
func (s StrengthDamage) EnsureValid() StrengthDamage {
	if int(s) < len(strengthDamageValues) {
		return s
	}
	return 0
}

// Key returns the key used to represent this StrengthDamage.
func (s StrengthDamage) Key() string {
	return strengthDamageValues[s.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (s StrengthDamage) String() string {
	return strengthDamageValues[s.EnsureValid()].String
}
