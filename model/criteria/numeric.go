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

package criteria

import (
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Numeric holds the criteria for matching a number.
type Numeric struct {
	Compare   NumericCompareType `json:"compare,omitempty"`
	Qualifier fixed.F64d4        `json:"qualifier,omitempty"`
}

// Normalize the data.
func (n *Numeric) Normalize() {
	if n.Compare = n.Compare.EnsureValid(); n.Compare == AnyNumber {
		n.Qualifier = 0
	}
}
