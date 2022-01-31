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
	FeetAndInches = LengthUnits("ftin") // This one is special and not a suffix
	Inch          = LengthUnits("in")
	Feet          = LengthUnits("ft")
	Yard          = LengthUnits("yd")
	Mile          = LengthUnits("mi")
	Centimeter    = LengthUnits("cm")
	Kilometer     = LengthUnits("km")
	Meter         = LengthUnits("m") // must come after Centimeter & Kilometer, as it's abbreviation is a subset
)

// AllLengthUnits is the complete set of LengthUnits values.
var AllLengthUnits = []LengthUnits{
	FeetAndInches,
	Inch,
	Feet,
	Yard,
	Mile,
	Centimeter,
	Kilometer,
	Meter,
}

// LengthUnits holds the length unit type. Note that conversions to/from metric are done using the simplified GURPS
// metric conversion of 1 yd = 1 meter. For consistency, all metric lengths are converted to meters, then to yards,
// rather than the variations at different lengths that the GURPS rules suggest.
type LengthUnits string

// EnsureValid ensures this is of a known value.
func (l LengthUnits) EnsureValid() LengthUnits {
	for _, one := range AllLengthUnits {
		if one == l {
			return l
		}
	}
	return AllLengthUnits[0]
}

// Format the length for this LengthUnits.
func (l LengthUnits) Format(length Length) string {
	switch l {
	case FeetAndInches:
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
	case Inch:
		return fixed.F64d4(length).String() + string(l)
	case Feet:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt64(12)).String() + string(l)
	case Yard:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt64(36)).String() + string(l)
	case Mile:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt64(5280)).String() + string(l)
	case Centimeter:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt64(36)).Mul(fixed.F64d4FromInt64(100)).String() + string(l)
	case Kilometer:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt64(36000)).String() + string(l)
	case Meter:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt64(36)).String() + string(l)
	default:
		return FeetAndInches.Format(length)
	}
}

// ToInches converts the length in this LengthUnits to inches.
func (l LengthUnits) ToInches(length fixed.F64d4) fixed.F64d4 {
	switch l {
	case FeetAndInches, Inch:
		return length
	case Feet:
		return length.Mul(fixed.F64d4FromInt64(12))
	case Yard:
		return length.Mul(fixed.F64d4FromInt64(36))
	case Mile:
		return length.Mul(fixed.F64d4FromInt64(63360))
	case Centimeter:
		return length.Mul(fixed.F64d4FromFloat64(36)).Div(fixed.F64d4FromInt64(100))
	case Kilometer:
		return length.Mul(fixed.F64d4FromInt64(36000))
	case Meter:
		return length.Mul(fixed.F64d4FromInt64(36))
	default:
		return FeetAndInches.ToInches(length)
	}
}
