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

package advantage

import (
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

// Possible SelfControlRoll values.
const (
	None = SelfControlRoll(0)
	CR6  = SelfControlRoll(6)
	CR9  = SelfControlRoll(9)
	CR12 = SelfControlRoll(12)
	CR15 = SelfControlRoll(15)
)

// AllSelfControlRolls is the complete set of SelfControlRoll values.
var AllSelfControlRolls = []SelfControlRoll{
	None,
	CR6,
	CR9,
	CR12,
	CR15,
}

// SelfControlRoll holds the information about a self-control roll, from B121.
type SelfControlRoll int

// EnsureValid ensures this is of a known value.
func (s SelfControlRoll) EnsureValid() SelfControlRoll {
	for _, one := range AllSelfControlRolls {
		if one == s {
			return s
		}
	}
	return AllSelfControlRolls[0]
}

// Index returns of the SelfControlRoll within AllSelfControlRolls.
func (s SelfControlRoll) Index() int {
	for i, one := range AllSelfControlRolls {
		if one == s {
			return i
		}
	}
	return 0
}

// String implements fmt.Stringer.
func (s SelfControlRoll) String() string {
	switch s {
	case None:
		return i18n.Text("None Required")
	case CR6:
		return i18n.Text("CR: 6 (Resist rarely)")
	case CR9:
		return i18n.Text("CR: 9 (Resist fairly often)")
	case CR12:
		return i18n.Text("CR: 12 (Resist quite often)")
	case CR15:
		return i18n.Text("CR: 15 (Resist almost all the time)")
	default:
		return None.String()
	}
}

// DescriptionWithCost returns a formatted description that includes the cost multiplier.
func (s SelfControlRoll) DescriptionWithCost() string {
	v := s.EnsureValid()
	if v == None {
		return ""
	}
	return v.String() + ", x" + v.Multiplier().String()
}

// Multiplier returns the cost multiplier.
func (s SelfControlRoll) Multiplier() f64d4.Int {
	switch s {
	case None:
		return f64d4.One
	case CR6:
		return fxp.Two
	case CR9:
		return fxp.OneAndAHalf
	case CR12:
		return f64d4.One
	case CR15:
		return f64d4.Half
	default:
		return None.Multiplier()
	}
}

// MinimumRoll returns the minimum roll to retain control.
func (s SelfControlRoll) MinimumRoll() int {
	return int(s.EnsureValid())
}
