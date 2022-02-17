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
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// WeightMultiplier returns the weight multiplier associated with the Encumbrance level.
func (enum Encumbrance) WeightMultiplier() fixed.F64d4 {
	switch enum {
	case None:
		return fixed.F64d4One
	case Light:
		return fxp.Two
	case Medium:
		return fxp.Three
	case Heavy:
		return fxp.Six
	case ExtraHeavy:
		return fxp.Ten
	default:
		return None.WeightMultiplier()
	}
}

// Penalty returns the penalty associated with the Encumbrance level.
func (enum Encumbrance) Penalty() fixed.F64d4 {
	switch enum {
	case None:
		return 0
	case Light:
		return fxp.NegOne
	case Medium:
		return fxp.NegTwo
	case Heavy:
		return fxp.NegThree
	case ExtraHeavy:
		return fxp.NegFour
	default:
		return None.Penalty()
	}
}
