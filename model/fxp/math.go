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
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// ApplyRounding truncates if 'roundDown' is true and performs a ceil() if false.
func ApplyRounding(value fixed.F64d4, roundDown bool) fixed.F64d4 {
	if roundDown {
		return value.Trunc()
	}
	if value.Trunc() != value {
		return value.Trunc() + fixed.F64d4One
	}
	return value
}

// ResetIfOutOfRange checks the value and if it is lower than min or greater than max, returns def, otherwise returns
// value.
func ResetIfOutOfRange(value, min, max, def fixed.F64d4) fixed.F64d4 {
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
