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
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible StrengthDamage values.
const (
	None          = StrengthDamage("none")
	Thrust        = StrengthDamage("thr")
	LeveledThrust = StrengthDamage("thr_leveled")
	Swing         = StrengthDamage("sw")
	LeveledSwing  = StrengthDamage("sw_leveled")
)

// AllStrengthDamages is the complete set of StrengthDamage values.
var AllStrengthDamages = []StrengthDamage{
	None,
	Thrust,
	LeveledThrust,
	Swing,
	LeveledSwing,
}

// StrengthDamage holds the type of strength dice to add to damage.
type StrengthDamage string

// EnsureValid ensures this is of a known value.
func (s StrengthDamage) EnsureValid() StrengthDamage {
	for _, one := range AllStrengthDamages {
		if one == s {
			return s
		}
	}
	return AllStrengthDamages[0]
}

// String implements fmt.Stringer.
func (s StrengthDamage) String() string {
	switch s {
	case None:
		return ""
	case Thrust:
		return "thr"
	case LeveledThrust:
		return "thr " + i18n.Text("(leveled)")
	case Swing:
		return "sw"
	case LeveledSwing:
		return "sw " + i18n.Text("(leveled)")
	default:
		return None.String()
	}
}
