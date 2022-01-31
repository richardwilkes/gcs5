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

package f64d4

import "github.com/richardwilkes/toolbox/xmath/fixed"

// Common values that can be reused.
var (
	One           = fixed.F64d4FromInt64(1)
	NegOne        = fixed.F64d4FromInt64(-1)
	Two           = fixed.F64d4FromInt64(2)
	OneAndAHalf   = fixed.F64d4FromStringForced("1.5")
	Half          = fixed.F64d4FromStringForced("0.5")
	NegHalf       = fixed.F64d4FromStringForced("-0.5")
	NegPointEight = fixed.F64d4FromStringForced("-0.8")
	Twenty        = fixed.F64d4FromInt64(20)
	Eighty        = fixed.F64d4FromInt64(80)
	NegEighty     = fixed.F64d4FromInt64(-80)
	Hundred       = fixed.F64d4FromInt64(100)
)
