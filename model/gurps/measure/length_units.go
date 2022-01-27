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

// Possible LengthUnits values.
const (
	FeetAndInches LengthUnits = iota // This one is special and not a suffix
	Centimeter
	Inch
	Feet
	Yard
	Kilometer
	Meter // must come after Centimeter & Kilometer, as it's abbreviation is a subset
	Mile
)

type lengthUnitsData struct {
	Key      string
	Format   func(length Length) string
	ToInches func(value fixed.F64d4) fixed.F64d4
}

// LengthUnits holds the length unit type. Note that conversions to/from metric are done using the simplified GURPS
// metric conversion of 1 yd = 1 meter. For consistency, all metric lengths are converted to meters, then to yards,
// rather than the variations at different lengths that the GURPS rules suggest.
type LengthUnits uint8

var lengthUnitsValues = []*lengthUnitsData{
	{
		Key: "ftin",
		Format: func(length Length) string {
			oneFoot := fixed.F64d4FromInt64(12)
			inches := fixed.F64d4(length)
			feet := inches.Div(oneFoot).Trunc()
			inches -= feet.Mul(oneFoot)
			if feet == 0 && inches == 0 {
				return "0'"
			}
			var buffer strings.Builder
			if feet > 0 {
				buffer.WriteString(feet.String())
				buffer.WriteByte('\'')
			}
			if inches > 0 {
				buffer.WriteString(inches.String())
				buffer.WriteByte('"')
			}
			return buffer.String()
		},
		ToInches: func(value fixed.F64d4) fixed.F64d4 { return value },
	},
	{
		Key: "cm",
		Format: func(length Length) string {
			return fixed.F64d4(length).Div(fixed.F64d4FromInt64(36)).Mul(fixed.F64d4FromInt64(100)).String() + "cm"
		},
		ToInches: func(value fixed.F64d4) fixed.F64d4 {
			return value.Mul(fixed.F64d4FromFloat64(36)).Div(fixed.F64d4FromInt64(100))
		},
	},
	{
		Key:      "in",
		Format:   func(length Length) string { return fixed.F64d4(length).String() + "in" },
		ToInches: func(value fixed.F64d4) fixed.F64d4 { return value },
	},
	{
		Key:      "ft",
		Format:   func(length Length) string { return fixed.F64d4(length).Div(fixed.F64d4FromInt64(12)).String() + "ft" },
		ToInches: func(value fixed.F64d4) fixed.F64d4 { return value.Mul(fixed.F64d4FromInt64(12)) },
	},
	{
		Key:      "yd",
		Format:   func(length Length) string { return fixed.F64d4(length).Div(fixed.F64d4FromInt64(36)).String() + "yd" },
		ToInches: func(value fixed.F64d4) fixed.F64d4 { return value.Mul(fixed.F64d4FromInt64(36)) },
	},
	{
		Key: "km",
		Format: func(length Length) string {
			return fixed.F64d4(length).Div(fixed.F64d4FromInt64(36000)).String() + "km"
		},
		ToInches: func(value fixed.F64d4) fixed.F64d4 { return value.Mul(fixed.F64d4FromInt64(36000)) },
	},
	{
		Key:      "m",
		Format:   func(length Length) string { return fixed.F64d4(length).Div(fixed.F64d4FromInt64(36)).String() + "m" },
		ToInches: func(value fixed.F64d4) fixed.F64d4 { return value.Mul(fixed.F64d4FromInt64(36)) },
	},
	{
		Key:      "mi",
		Format:   func(length Length) string { return fixed.F64d4(length).Div(fixed.F64d4FromInt64(5280)).String() + "mi" },
		ToInches: func(value fixed.F64d4) fixed.F64d4 { return value.Mul(fixed.F64d4FromInt64(63360)) },
	},
}

// LengthUnitsFromString extracts a LengthUnits from a key.
func LengthUnitsFromString(key string) LengthUnits {
	for i, one := range lengthUnitsValues {
		if strings.EqualFold(key, one.Key) {
			return LengthUnits(i)
		}
	}
	return 0
}

// EnsureValid returns the first LengthUnits if this LengthUnits is not a known value.
func (l LengthUnits) EnsureValid() LengthUnits {
	if int(l) < len(lengthUnitsValues) {
		return l
	}
	return 0
}

// Key returns the key used to represent this LengthUnits.
func (l LengthUnits) Key() string {
	return lengthUnitsValues[l.EnsureValid()].Key
}

// Format the length for this LengthUnits.
func (l LengthUnits) Format(length Length) string {
	return lengthUnitsValues[l.EnsureValid()].Format(length)
}

// ToInches converts the length in this LengthUnits to inches.
func (l LengthUnits) ToInches(length fixed.F64d4) fixed.F64d4 {
	return lengthUnitsValues[l.EnsureValid()].ToInches(length)
}
