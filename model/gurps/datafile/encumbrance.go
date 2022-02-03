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

package datafile

import (
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Possible Encumbrance values.
const (
	None       = Encumbrance("none")
	Light      = Encumbrance("light")
	Medium     = Encumbrance("medium")
	Heavy      = Encumbrance("heavy")
	ExtraHeavy = Encumbrance("extra_heavy")
)

// AllEncumbrances is the complete set of Encumbrance values.
var AllEncumbrances = []Encumbrance{
	None,
	Light,
	Medium,
	Heavy,
	ExtraHeavy,
}

// Encumbrance holds the encumbrance level.
type Encumbrance string

// EnsureValid ensures this is of a known value.
func (e Encumbrance) EnsureValid() Encumbrance {
	for _, one := range AllEncumbrances {
		if one == e {
			return e
		}
	}
	return AllEncumbrances[0]
}

// String implements fmt.Stringer.
func (e Encumbrance) String() string {
	switch e {
	case None:
		return i18n.Text("None")
	case Light:
		return i18n.Text("Light")
	case Medium:
		return i18n.Text("Medium")
	case Heavy:
		return i18n.Text("Heavy")
	case ExtraHeavy:
		return i18n.Text("X-Heavy")
	default:
		return None.String()
	}
}

// WeightMultiplier returns the weight multiplier associated with the Encumbrance level.
func (e Encumbrance) WeightMultiplier() fixed.F64d4 {
	switch e {
	case None:
		return f64d4.One
	case Light:
		return f64d4.Two
	case Medium:
		return f64d4.Three
	case Heavy:
		return f64d4.Six
	case ExtraHeavy:
		return f64d4.Ten
	default:
		return None.WeightMultiplier()
	}
}

// Penalty returns the penalty associated with the Encumbrance level.
func (e Encumbrance) Penalty() fixed.F64d4 {
	switch e {
	case None:
		return 0
	case Light:
		return f64d4.NegOne
	case Medium:
		return f64d4.NegTwo
	case Heavy:
		return f64d4.NegThree
	case ExtraHeavy:
		return f64d4.NegFour
	default:
		return None.Penalty()
	}
}
