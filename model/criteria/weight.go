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
	"github.com/richardwilkes/gcs/model/gurps/measure"
)

// Weight holds the criteria for matching a number.
type Weight struct {
	Compare   NumericCompareType `json:"compare,omitempty"`
	Qualifier measure.Weight     `json:"qualifier,omitempty"`
}

// Normalize the data.
func (w *Weight) Normalize() {
	if w.Compare = w.Compare.EnsureValid(); w.Compare == AnyNumber {
		w.Qualifier = 0
	}
}
