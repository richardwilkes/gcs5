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
