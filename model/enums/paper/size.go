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
	"github.com/richardwilkes/gcs/model/unit/length"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
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
func (s Size) Dimensions() (width, height length.GURPS) {
	switch s {
	case Legal:
		return length.GURPSFromFloat64(8.5, units.Inch), length.GURPSFromInt64(14, units.Inch)
	case Tabloid:
		return length.GURPSFromInt64(11, units.Inch), length.GURPSFromInt64(17, units.Inch)
	case A0:
		return mmToInches(841), mmToInches(1189)
	case A1:
		return mmToInches(594), mmToInches(841)
	case A2:
		return mmToInches(420), mmToInches(594)
	case A3:
		return mmToInches(297), mmToInches(420)
	case A4:
		return mmToInches(210), mmToInches(297)
	case A5:
		return mmToInches(148), mmToInches(210)
	case A6:
		return mmToInches(105), mmToInches(148)
	default: // Letter
		return length.GURPSFromFloat64(8.5, units.Inch), length.GURPSFromInt64(11, units.Inch)
	}
}

func mmToInches(mm int) length.GURPS {
	// Do our own conversion to inches, since the metric values are using GURPS' simplified values
	return length.GURPS(fixed.F64d4FromInt64(int64(mm)*10) / fixed.F64d4FromInt64(254))
}
