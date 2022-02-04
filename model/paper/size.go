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
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Size values.
const (
	Letter  = Size("letter")
	Legal   = Size("legal")
	Tabloid = Size("tabloid")
	A0      = Size("a0")
	A1      = Size("a1")
	A2      = Size("a2")
	A3      = Size("a3")
	A4      = Size("a4")
	A5      = Size("a5")
	A6      = Size("a6")
)

// AllSizes is the complete set of Size values.
var AllSizes = []Size{
	Letter,
	Legal,
	Tabloid,
	A0,
	A1,
	A2,
	A3,
	A4,
	A5,
	A6,
}

// Size holds a standard paper dimension.
type Size string

// EnsureValid ensures this is of a known value.
func (s Size) EnsureValid() Size {
	for _, one := range AllSizes {
		if one == s {
			return s
		}
	}
	return AllSizes[0]
}

// String implements fmt.Stringer.
func (s Size) String() string {
	switch s {
	case Letter:
		return i18n.Text("Letter")
	case Legal:
		return i18n.Text("Legal")
	case Tabloid:
		return i18n.Text("Tabloid")
	case A0:
		return i18n.Text("A0")
	case A1:
		return i18n.Text("A1")
	case A2:
		return i18n.Text("A2")
	case A3:
		return i18n.Text("A3")
	case A4:
		return i18n.Text("A4")
	case A5:
		return i18n.Text("A5")
	case A6:
		return i18n.Text("A6")
	default:
		return Letter.String()
	}
}

// Dimensions returns the paper dimensions.
func (s Size) Dimensions() (width, height Length) {
	switch s {
	case Letter:
		return Length{Length: 8.5, Units: Inch}, Length{Length: 11, Units: Inch}
	case Legal:
		return Length{Length: 8.5, Units: Inch}, Length{Length: 14, Units: Inch}
	case Tabloid:
		return Length{Length: 11, Units: Inch}, Length{Length: 17, Units: Inch}
	case A0:
		return Length{Length: 841, Units: Millimeter}, Length{Length: 1189, Units: Millimeter}
	case A1:
		return Length{Length: 594, Units: Millimeter}, Length{Length: 841, Units: Millimeter}
	case A2:
		return Length{Length: 420, Units: Millimeter}, Length{Length: 594, Units: Millimeter}
	case A3:
		return Length{Length: 297, Units: Millimeter}, Length{Length: 420, Units: Millimeter}
	case A4:
		return Length{Length: 210, Units: Millimeter}, Length{Length: 297, Units: Millimeter}
	case A5:
		return Length{Length: 148, Units: Millimeter}, Length{Length: 210, Units: Millimeter}
	case A6:
		return Length{Length: 105, Units: Millimeter}, Length{Length: 148, Units: Millimeter}
	default:
		return Letter.Dimensions()
	}
}
