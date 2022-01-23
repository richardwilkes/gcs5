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

package length

import (
	"fmt"
	"strings"

	"github.com/richardwilkes/gcs/model/enums/units"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// GURPS contains a fixed-point value in inches. Conversions to/from metric are done using the simplified GURPS metric
// conversion of 1 yd = 1 meter. For consistency, all metric lengths are converted to meters, then to yards, rather than
// the variations at different lengths that the GURPS rules suggest.
type GURPS fixed.F64d4

// GURPSFromInt64 creates a new GURPS.
func GURPSFromInt64(value int64, unit units.GURPSLength) GURPS {
	return convertGURPSToInches(fixed.F64d4FromInt64(value), unit)
}

// GURPSFromFloat64 creates a new GURPS.
func GURPSFromFloat64(value float64, unit units.GURPSLength) GURPS {
	return convertGURPSToInches(fixed.F64d4FromFloat64(value), unit)
}

// GURPSFromStringForced creates a new GURPS. May have any of the known Units suffixes, a feet and inches format (e.g.
// 6'2"), or no notation at all, in which case defaultUnits is used.
func GURPSFromStringForced(text string, defaultUnits units.GURPSLength) GURPS {
	length, err := GURPSFromString(text, defaultUnits)
	if err != nil {
		return 0
	}
	return length
}

// GURPSFromString creates a new GURPS. May have any of the known Units suffixes, a feet and inches format (e.g. 6'2"),
// or no notation at all, in which case defaultUnits is used.
func GURPSFromString(text string, defaultUnits units.GURPSLength) (GURPS, error) {
	text = strings.TrimLeft(strings.TrimSpace(text), "+")
	for unit := units.Centimeter; unit <= units.Mile; unit++ {
		if strings.HasSuffix(text, unit.Key()) {
			value, err := fixed.F64d4FromString(strings.TrimSpace(strings.TrimSuffix(text, unit.Key())))
			if err != nil {
				return 0, err
			}
			return convertGURPSToInches(value, unit), nil
		}
	}
	// Didn't match any of the Units types, let's try feet & inches
	feetIndex := strings.Index(text, "'")
	inchIndex := strings.Index(text, `"`)
	if feetIndex == -1 && inchIndex == -1 {
		// Nope, so let's use our passed-in default units
		value, err := fixed.F64d4FromString(strings.TrimSpace(text))
		if err != nil {
			return 0, err
		}
		return convertGURPSToInches(value, defaultUnits), nil
	}
	var feet, inches fixed.F64d4
	var err error
	if feetIndex != -1 {
		s := strings.TrimSpace(text[:feetIndex])
		feet, err = fixed.F64d4FromString(s)
		if err != nil {
			return 0, err
		}
	}
	if inchIndex != -1 {
		if feetIndex > inchIndex {
			return 0, errs.New(fmt.Sprintf("invalid format: %s", text))
		}
		s := strings.TrimSpace(text[feetIndex+1 : inchIndex])
		inches, err = fixed.F64d4FromString(s)
		if err != nil {
			return 0, err
		}
	}
	return GURPS(feet.Mul(fixed.F64d4FromInt64(12)) + inches), nil
}

func convertGURPSToInches(value fixed.F64d4, unit units.GURPSLength) GURPS {
	switch unit {
	case units.Centimeter:
		value = value.Mul(fixed.F64d4FromFloat64(36)).Div(fixed.F64d4FromInt64(100))
	case units.Feet:
		value = value.Mul(fixed.F64d4FromInt64(12))
	case units.Yard, units.Meter:
		value = value.Mul(fixed.F64d4FromInt64(36))
	case units.Kilometer:
		value = value.Mul(fixed.F64d4FromInt64(36000))
	case units.Mile:
		value = value.Mul(fixed.F64d4FromInt64(63360))
	default: // Same as Inch
	}
	return GURPS(value)
}

func (l GURPS) String() string {
	return l.Format(units.FeetAndInches)
}

// Format the length as the given units.
func (l GURPS) Format(unit units.GURPSLength) string {
	inches := fixed.F64d4(l)
	switch unit {
	case units.Centimeter:
		return inches.Div(fixed.F64d4FromInt64(36)).Mul(fixed.F64d4FromInt64(100)).String() + unit.Key()
	case units.Feet:
		return inches.Div(fixed.F64d4FromInt64(12)).String() + unit.Key()
	case units.Yard, units.Meter:
		return inches.Div(fixed.F64d4FromInt64(36)).String() + unit.Key()
	case units.Kilometer:
		return inches.Div(fixed.F64d4FromInt64(36000)).String() + unit.Key()
	case units.Mile:
		return inches.Div(fixed.F64d4FromInt64(5280)).String() + unit.Key()
	case units.FeetAndInches:
		oneFoot := fixed.F64d4FromInt64(12)
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
	default: // Same as Inch
		return inches.String() + units.Inch.Key()
	}
}
