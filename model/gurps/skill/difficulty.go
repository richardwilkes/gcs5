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

package skill

import (
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

// AllTechniqueDifficulty holds all possible values when used with Techniques.
var AllTechniqueDifficulty = []Difficulty{
	Average,
	Hard,
}

// BaseRelativeLevel returns the base relative skill level at 0 points.
func (enum Difficulty) BaseRelativeLevel() f64d4.Int {
	switch enum {
	case Easy:
		return 0
	case Average:
		return fxp.NegOne
	case Hard:
		return fxp.NegTwo
	case VeryHard, Wildcard:
		return fxp.NegThree
	default:
		return Easy.BaseRelativeLevel()
	}
}
