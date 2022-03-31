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
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

// Weight holds the criteria for matching a number.
type Weight struct {
	WeightData
}

// WeightData holds the criteria for matching a number that should be written to disk.
type WeightData struct {
	Compare   NumericCompareType `json:"compare,omitempty"`
	Qualifier measure.Weight     `json:"qualifier,omitempty"`
}

// ShouldOmit implements json.Omitter.
func (w Weight) ShouldOmit() bool {
	return w.Compare.EnsureValid() == AnyNumber
}

// UnmarshalJSON implements json.Unmarshaler.
func (w *Weight) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &w.WeightData)
	w.Compare = w.Compare.EnsureValid()
	return err
}

// Matches performs a comparison and returns true if the data matches.
func (w Weight) Matches(value measure.Weight) bool {
	return w.Compare.Matches(f64d4.Int(w.Qualifier), f64d4.Int(value))
}

func (w Weight) String() string {
	return w.Compare.Describe(f64d4.Int(w.Qualifier))
}
