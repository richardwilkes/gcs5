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

package weight

import (
	"strings"

	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Units holds the weight unit type. Note that conversions to/from metric are done using the simplified GURPS metric
// conversion of 1# = 0.5kg. For consistency, all metric weights are converted to kilograms, then to pounds, rather than
// the variations at different weights that the GURPS rules suggest.
type Units string

// Possible Units values.
const (
	Gram     = Units("g")
	Ounce    = Units("oz")
	Pound    = Units("#")
	PoundAlt = Units("lb")
	Kilogram = Units("kg")
	Ton      = Units("ton")
)

// Weight contains a fixed-point value in pounds.
type Weight fixed.F64d4

// FromInt64 creates a new Weight.
func FromInt64(value int64, unit Units) Weight {
	return convertToPounds(fixed.F64d4FromInt64(value), unit)
}

// FromFloat64 creates a new Weight.
func FromFloat64(value float64, unit Units) Weight {
	return convertToPounds(fixed.F64d4FromFloat64(value), unit)
}

// FromString creates a new Weight. May have any of the known Units suffixes or no notation at all, in which case
// defaultUnits is used.
func FromString(text string, defaultUnits Units) (Weight, error) {
	text = strings.TrimLeft(strings.TrimSpace(text), "+")
	for _, unit := range []Units{
		Ounce,
		Pound,
		PoundAlt,
		Kilogram,
		Gram, // must come after Kilogram, as it's abbreviation is a subset
		Ton,
	} {
		if strings.HasSuffix(text, string(unit)) {
			value, err := fixed.F64d4FromString(strings.TrimSpace(strings.TrimSuffix(text, string(unit))))
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

func convertToPounds(value fixed.F64d4, unit Units) Weight {
	switch unit {
	case Gram:
		value = value.Div(fixed.F64d4FromInt64(500))
	case Ounce:
		value = value.Div(fixed.F64d4FromInt64(16))
	case Kilogram:
		value = value.Mul(fixed.F64d4FromInt64(2))
	case Ton:
		value = value.Mul(fixed.F64d4FromInt64(2000))
	default: // Same as Pound
	}
	return Weight(value)
}

func (w Weight) String() string {
	return w.Format(Pound)
}

// Format the weight as the given units.
func (w Weight) Format(unit Units) string {
	pounds := fixed.F64d4(w)
	switch unit {
	case Gram:
		return pounds.Mul(fixed.F64d4FromInt64(500)).String() + string(unit)
	case Ounce:
		return pounds.Mul(fixed.F64d4FromInt64(16)).String() + string(unit)
	case Kilogram:
		return pounds.Div(fixed.F64d4FromInt64(2)).String() + string(unit)
	case Ton:
		return pounds.Div(fixed.F64d4FromInt64(2000)).String() + string(unit)
	default: // Same as Pound
		return pounds.String() + string(Pound)
	}
}

// MarshalText implements the encoding.TextMarshaler interface.
func (w Weight) MarshalText() (text []byte, err error) {
	return []byte(w.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (w *Weight) UnmarshalText(text []byte) error {
	var err error
	if *w, err = FromString(string(text), Pound); err != nil {
		return err
	}
	return nil
}
