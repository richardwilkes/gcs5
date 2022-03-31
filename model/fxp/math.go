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

package fxp

import (
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

// ApplyRounding truncates if 'roundDown' is true and performs a ceil() if false.
func ApplyRounding(value f64d4.Int, roundDown bool) f64d4.Int {
	if roundDown {
		return value.Trunc()
	}
	if value.Trunc() != value {
		return value.Trunc() + f64d4.One
	}
	return value
}

// ResetIfOutOfRange checks the value and if it is lower than min or greater than max, returns def, otherwise returns
// value.
func ResetIfOutOfRange(value, min, max, def f64d4.Int) f64d4.Int {
	if value < min || value > max {
		return def
	}
	return value
}

// ResetIfOutOfRangeInt checks the value and if it is lower than min or greater than max, returns def, otherwise returns
// value.
func ResetIfOutOfRangeInt(value, min, max, def int) int {
	if value < min || value > max {
		return def
	}
	return value
}
