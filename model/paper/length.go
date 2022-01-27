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

import (
	"strconv"
	"strings"
)

// Length contains a real-world length value with an attached units.
type Length struct {
	Length float64
	Units  Units
}

// LengthFromString creates a new Length. May have any of the known units.Units suffixes or no notation at all, in which
// case units.Inch is used.
func LengthFromString(text string) Length {
	text = strings.TrimLeft(strings.TrimSpace(text), "+")
	for unit := Millimeter; unit <= Inch; unit++ {
		if strings.HasSuffix(text, unit.Key()) {
			value, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(text, unit.Key())), 64)
			if err != nil {
				return Length{Units: unit}
			}
			return Length{Length: value, Units: unit}
		}
	}
	// Didn't match any of the Units types, assume the default
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return Length{Units: Inch}
	}
	return Length{Length: value, Units: Inch}
}

func (l Length) String() string {
	return strconv.FormatFloat(l.Length, 'f', -1, 64) + l.Units.Key()
}

// Pixels returns the number of 72-pixels-per-inch pixels this represents.
func (l Length) Pixels() float32 {
	return l.Units.ToPixels(l.Length)
}
