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

package measure

import (
	"strings"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64"
)

// Weight contains a fixed-point value in pounds.
type Weight fxp.Int

// WeightFromInt64 creates a new Weight.
func WeightFromInt64(value int64, unit WeightUnits) Weight {
	return Weight(unit.ToPounds(f64.From[fxp.DP](value)))
}

// WeightFromInt creates a new Weight.
func WeightFromInt(value int, unit WeightUnits) Weight {
	return Weight(unit.ToPounds(f64.From[fxp.DP](value)))
}

// WeightFromStringForced creates a new Weight. May have any of the known Weight suffixes or no notation at all, in which
// case defaultUnits is used.
func WeightFromStringForced(text string, defaultUnits WeightUnits) Weight {
	weight, err := WeightFromString(text, defaultUnits)
	if err != nil {
		return 0
	}
	return weight
}

// WeightFromString creates a new Weight. May have any of the known Weight suffixes or no notation at all, in which case
// defaultUnits is used.
func WeightFromString(text string, defaultUnits WeightUnits) (Weight, error) {
	text = strings.TrimLeft(strings.TrimSpace(text), "+")
	for _, unit := range AllWeightUnits {
		if strings.HasSuffix(text, unit.Key()) {
			value, err := f64.FromString[fxp.DP](strings.TrimSpace(strings.TrimSuffix(text, unit.Key())))
			if err != nil {
				return 0, err
			}
			return Weight(unit.ToPounds(value)), nil
		}
	}
	// No matches, so let's use our passed-in default units
	value, err := f64.FromString[fxp.DP](strings.TrimSpace(text))
	if err != nil {
		return 0, err
	}
	return Weight(defaultUnits.ToPounds(value)), nil
}

func (w Weight) String() string {
	return Pound.Format(w)
}

// MarshalJSON implements json.Marshaler.
func (w Weight) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (w *Weight) UnmarshalJSON(in []byte) error {
	var s string
	if err := json.Unmarshal(in, &s); err != nil {
		return err
	}
	var err error
	*w, err = WeightFromString(s, Pound)
	return err
}
