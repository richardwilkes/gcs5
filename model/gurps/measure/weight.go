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

	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Weight contains a fixed-point value in pounds.
type Weight fixed.F64d4

// WeightFromInt64 creates a new Weight.
func WeightFromInt64(value int64, unit WeightUnits) Weight {
	return convertToPounds(fixed.F64d4FromInt64(value), unit)
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
	for unit := Pound; unit <= Gram; unit++ {
		if strings.HasSuffix(text, unit.Key()) {
			value, err := fixed.F64d4FromString(strings.TrimSpace(strings.TrimSuffix(text, unit.Key())))
			if err != nil {
				return 0, err
			}
			return convertToPounds(value, unit), nil
		}
	}
	// No matches, so let's use our passed-in default units
	value, err := fixed.F64d4FromString(strings.TrimSpace(text))
	if err != nil {
		return 0, err
	}
	return convertToPounds(value, defaultUnits), nil
}

func convertToPounds(value fixed.F64d4, unit WeightUnits) Weight {
	switch unit {
	case Ounce:
		value = value.Div(fixed.F64d4FromInt64(16))
	case Ton:
		value = value.Mul(fixed.F64d4FromInt64(2000))
	case Kilogram:
		value = value.Mul(fixed.F64d4FromInt64(2))
	case Gram:
		value = value.Div(fixed.F64d4FromInt64(500))
	default: // Same as Pound
	}
	return Weight(value)
}

func (w Weight) String() string {
	return w.Format(Pound)
}

// Format the weight as the given units.
func (w Weight) Format(unit WeightUnits) string {
	pounds := fixed.F64d4(w)
	switch unit {
	case Ounce:
		return pounds.Mul(fixed.F64d4FromInt64(16)).String() + unit.Key()
	case Ton:
		return pounds.Div(fixed.F64d4FromInt64(2000)).String() + unit.Key()
	case Kilogram:
		return pounds.Div(fixed.F64d4FromInt64(2)).String() + unit.Key()
	case Gram:
		return pounds.Mul(fixed.F64d4FromInt64(500)).String() + unit.Key()
	default: // Same as Pound
		return pounds.String() + Pound.Key()
	}
}
