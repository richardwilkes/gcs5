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

package paper

// Possible Units values.
const (
	Inch       = Units("in")
	Centimeter = Units("cm")
	Millimeter = Units("mm")
)

// AllUnits is the complete set of Unit values.
var AllUnits = []Units{
	Inch,
	Centimeter,
	Millimeter,
}

// Units holds the real-world length unit type.
type Units string

// EnsureValid ensures this is of a known value.
func (u Units) EnsureValid() Units {
	for _, one := range AllUnits {
		if one == u {
			return u
		}
	}
	return AllUnits[0]
}

// ToPixels converts the given length in this Units to the number of 72-pixels-per-inch pixels it represents.
func (u Units) ToPixels(length float64) float32 {
	switch u {
	case Inch:
		return float32(length * 72)
	case Centimeter:
		return float32((length * 72) / 2.54)
	case Millimeter:
		return float32((length * 72) / 25.4)
	default:
		return Inch.ToPixels(length)
	}
}
