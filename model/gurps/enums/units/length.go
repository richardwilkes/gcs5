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

package units

import "strings"

// Length holds the length unit type. Note that conversions to/from metric are done using the simplified GURPS metric
// conversion of 1 yd = 1 meter. For consistency, all metric lengths are converted to meters, then to yards, rather than
// the variations at different lengths that the GURPS rules suggest.
type Length uint8

// Possible Length values.
const (
	FeetAndInches Length = iota // This one is special and not a suffix
	Centimeter
	Inch
	Feet
	Yard
	Kilometer
	Meter // must come after Centimeter & Kilometer, as it's abbreviation is a subset
	Mile
)

// LengthFromString extracts a Length from a string.
func LengthFromString(str string) Length {
	for p := FeetAndInches; p <= Mile; p++ {
		if strings.EqualFold(p.Key(), str) {
			return p
		}
	}
	return FeetAndInches
}

// Key returns the key used to represent this Length.
func (l Length) Key() string {
	switch l {
	case Centimeter:
		return "cm"
	case Inch:
		return "in"
	case Feet:
		return "ft"
	case Yard:
		return "yd"
	case Kilometer:
		return "km"
	case Meter:
		return "m"
	case Mile:
		return "mi"
	default: // FeetAndInches
		return "ftin"
	}
}
