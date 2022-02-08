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

// Ceil returns the value rounded up to the nearest whole number.
func Ceil(value fixed.F64d4) fixed.F64d4 {
	v := value.Trunc()
	if value != v {
		v += One
	}
	return v
}

// Round the given value.
func Round(value fixed.F64d4) fixed.F64d4 {
	rem := value
	value = value.Trunc()
	rem -= value
	if rem >= Half {
		value += One
	} else if rem < NegHalf {
		value -= One
	}
	return value
}

// ApplyRounding truncates if 'roundDown' is true and performs a ceil() if false.
func ApplyRounding(value fixed.F64d4, roundDown bool) fixed.F64d4 {
	if roundDown {
		return value.Trunc()
	}
	if value.Trunc() != value {
		return value.Trunc() + One
	}
	return value
}

// Mod returns the remainder of x/y.
func Mod(x, y fixed.F64d4) fixed.F64d4 {
	return x - (y.Mul(x.Div(y).Trunc()))
}
