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

// Format the length for this LengthUnits.
func (enum LengthUnits) Format(length Length) string {
	switch enum {
	case FeetAndInches:
		oneFoot := fixed.F64d4FromInt(12)
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
		return fixed.F64d4(length).String() + enum.Key()
	case Feet:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt(12)).String() + enum.Key()
	case Yard:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt(36)).String() + enum.Key()
	case Mile:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt(5280)).String() + enum.Key()
	case Centimeter:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt(36)).Mul(fixed.F64d4FromInt(100)).String() + enum.Key()
	case Kilometer:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt(36000)).String() + enum.Key()
	case Meter:
		return fixed.F64d4(length).Div(fixed.F64d4FromInt(36)).String() + enum.Key()
	default:
		return FeetAndInches.Format(length)
	}
}

// ToInches converts the length in this LengthUnits to inches.
func (enum LengthUnits) ToInches(length fixed.F64d4) fixed.F64d4 {
	switch enum {
	case FeetAndInches, Inch:
		return length
	case Feet:
		return length.Mul(fixed.F64d4FromInt(12))
	case Yard:
		return length.Mul(fixed.F64d4FromInt(36))
	case Mile:
		return length.Mul(fixed.F64d4FromInt(63360))
	case Centimeter:
		return length.Mul(fixed.F64d4FromInt(36)).Div(fixed.F64d4FromInt(100))
	case Kilometer:
		return length.Mul(fixed.F64d4FromInt(36000))
	case Meter:
		return length.Mul(fixed.F64d4FromInt(36))
	default:
		return FeetAndInches.ToInches(length)
	}
}
