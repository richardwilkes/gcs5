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
