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
	"strconv"
	"strings"

	"github.com/richardwilkes/gcs/model/enums/units"
)

// Length contains a real-world length value with an attached units.
type Length struct {
	Length float64
	Units  units.Length
}

// LengthFromString creates a new Length. May have any of the known units.Length suffixes or no notation at all, in which
// case units.Inch is used.
func LengthFromString(text string) Length {
	text = strings.TrimLeft(strings.TrimSpace(text), "+")
	for unit := units.Millimeter; unit <= units.Inch; unit++ {
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
		return Length{Units: units.Inch}
	}
	return Length{Length: value, Units: units.Inch}
}

func (l Length) String() string {
	return strconv.FormatFloat(l.Length, 'f', -1, 64) + l.Units.Key()
}

// Pixels returns the number of 72-pixels-per-inch pixels this represents.
func (l Length) Pixels() float32 {
	length := l.Length * 72
	switch l.Units {
	case units.Millimeter:
		return float32(length / 25.4)
	case units.Centimeter:
		return float32(length / 2.54)
	default:
		return float32(length)
	}
}