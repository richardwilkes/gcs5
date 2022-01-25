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
	"strings"

	"github.com/richardwilkes/gcs/model/enums/units"
	"github.com/richardwilkes/gcs/model/measure"
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Size values.
const (
	Letter Size = iota
	Legal
	Tabloid
	A0
	A1
	A2
	A3
	A4
	A5
	A6
)

// Size holds a standard paper dimension.
type Size uint8

// SizeFromString extracts a Size from a string.
func SizeFromString(str string) Size {
	for one := Letter; one <= A6; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return Letter
}

// Key returns the key used to represent this ThresholdOp.
func (s Size) Key() string {
	switch s {
	case Legal:
		return "legal"
	case Tabloid:
		return "tabloid"
	case A0:
		return "a0"
	case A1:
		return "a1"
	case A2:
		return "a2"
	case A3:
		return "a3"
	case A4:
		return "a4"
	case A5:
		return "a5"
	case A6:
		return "a6"
	default: // Letter
		return "letter"
	}
}

// String implements fmt.Stringer.
func (s Size) String() string {
	switch s {
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
	default: // Letter
		return i18n.Text("Letter")
	}
}

// Dimensions returns the paper dimensions.
func (s Size) Dimensions() (width, height measure.Length) {
	switch s {
	case Legal:
		return measure.LengthFromString("8.5in"), measure.Length{Length: 14, Units: units.Inch}
	case Tabloid:
		return measure.Length{Length: 11, Units: units.Inch}, measure.Length{Length: 17, Units: units.Inch}
	case A0:
		return measure.Length{Length: 841, Units: units.Millimeter}, measure.Length{Length: 1189, Units: units.Millimeter}
	case A1:
		return measure.Length{Length: 594, Units: units.Millimeter}, measure.Length{Length: 841, Units: units.Millimeter}
	case A2:
		return measure.Length{Length: 420, Units: units.Millimeter}, measure.Length{Length: 594, Units: units.Millimeter}
	case A3:
		return measure.Length{Length: 297, Units: units.Millimeter}, measure.Length{Length: 420, Units: units.Millimeter}
	case A4:
		return measure.Length{Length: 210, Units: units.Millimeter}, measure.Length{Length: 297, Units: units.Millimeter}
	case A5:
		return measure.Length{Length: 148, Units: units.Millimeter}, measure.Length{Length: 210, Units: units.Millimeter}
	case A6:
		return measure.Length{Length: 105, Units: units.Millimeter}, measure.Length{Length: 148, Units: units.Millimeter}
	default: // Letter
		return measure.LengthFromString("8.5in"), measure.Length{Length: 11, Units: units.Inch}
	}
}
