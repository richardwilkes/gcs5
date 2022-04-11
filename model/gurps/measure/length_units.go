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

	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

// Format the length for this LengthUnits.
func (enum LengthUnits) Format(length Length) string {
	inches := f64d4.Int(length)
	switch enum {
	case FeetAndInches:
		oneFoot := f64d4.FromInt(12)
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
		return inches.String() + " " + enum.Key()
	case Feet:
		return inches.Div(f64d4.FromInt(12)).String() + " " + enum.Key()
	case Yard, Meter:
		return inches.Div(f64d4.FromInt(36)).String() + " " + enum.Key()
	case Mile:
		return inches.Div(f64d4.FromInt(63360)).String() + " " + enum.Key()
	case Centimeter:
		return inches.Div(f64d4.FromInt(36)).Mul(f64d4.FromInt(100)).String() + " " + enum.Key()
	case Kilometer:
		return inches.Div(f64d4.FromInt(36000)).String() + " " + enum.Key()
	default:
		return FeetAndInches.Format(length)
	}
}

// ToInches converts the length in this LengthUnits to inches.
func (enum LengthUnits) ToInches(length f64d4.Int) f64d4.Int {
	switch enum {
	case FeetAndInches, Inch:
		return length
	case Feet:
		return length.Mul(f64d4.FromInt(12))
	case Yard:
		return length.Mul(f64d4.FromInt(36))
	case Mile:
		return length.Mul(f64d4.FromInt(63360))
	case Centimeter:
		return length.Mul(f64d4.FromInt(36)).Div(f64d4.FromInt(100))
	case Kilometer:
		return length.Mul(f64d4.FromInt(36000))
	case Meter:
		return length.Mul(f64d4.FromInt(36))
	default:
		return FeetAndInches.ToInches(length)
	}
}
