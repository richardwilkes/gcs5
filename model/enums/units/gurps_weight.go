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

// GURPSWeight holds the weight unit type. Note that conversions to/from metric are done using the simplified GURPS
// metric conversion of 1# = 0.5kg. For consistency, all metric weights are converted to kilograms, then to pounds,
// rather than the variations at different weights that the GURPS rules suggest.
type GURPSWeight uint8

// Possible GURPSWeight values.
const (
	Pound GURPSWeight = iota
	PoundAlt
	Ounce
	Ton
	Kilogram
	Gram // must come after Kilogram, as it's abbreviation is a subset
)

// GURPSWeightFromString extracts a GURPSLength from a string.
func GURPSWeightFromString(str string) GURPSWeight {
	for p := Pound; p <= Gram; p++ {
		if strings.EqualFold(p.Key(), str) {
			return p
		}
	}
	return Pound
}

// Key returns the key used to represent this GURPSWeight.
func (l GURPSWeight) Key() string {
	switch l {
	case PoundAlt:
		return "lb"
	case Ounce:
		return "oz"
	case Ton:
		return "tn"
	case Kilogram:
		return "kg"
	case Gram:
		return "g"
	default: // Pound
		return "#"
	}
}
