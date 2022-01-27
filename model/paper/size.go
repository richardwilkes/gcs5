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

type sizeData struct {
	Key    string
	String string
	Width  Length
	Height Length
}

// Size holds a standard paper dimension.
type Size uint8

var sizeValues = []*sizeData{
	{
		Key:    "letter",
		String: i18n.Text("Letter"),
		Width:  Length{Length: 8.5, Units: Inch},
		Height: Length{Length: 11, Units: Inch},
	},
	{
		Key:    "legal",
		String: i18n.Text("Legal"),
		Width:  Length{Length: 8.5, Units: Inch},
		Height: Length{Length: 14, Units: Inch},
	},
	{
		Key:    "tabloid",
		String: i18n.Text("Tabloid"),
		Width:  Length{Length: 11, Units: Inch},
		Height: Length{Length: 17, Units: Inch},
	},
	{
		Key:    "a0",
		String: "A0",
		Width:  Length{Length: 841, Units: Millimeter},
		Height: Length{Length: 1189, Units: Millimeter},
	},
	{
		Key:    "a1",
		String: "A1",
		Width:  Length{Length: 594, Units: Millimeter},
		Height: Length{Length: 841, Units: Millimeter},
	},
	{
		Key:    "a2",
		String: "A2",
		Width:  Length{Length: 420, Units: Millimeter},
		Height: Length{Length: 594, Units: Millimeter},
	},
	{
		Key:    "a3",
		String: "A3",
		Width:  Length{Length: 297, Units: Millimeter},
		Height: Length{Length: 420, Units: Millimeter},
	},
	{
		Key:    "a4",
		String: "A4",
		Width:  Length{Length: 210, Units: Millimeter},
		Height: Length{Length: 297, Units: Millimeter},
	},
	{
		Key:    "a5",
		String: "A5",
		Width:  Length{Length: 148, Units: Millimeter},
		Height: Length{Length: 210, Units: Millimeter},
	},
	{
		Key:    "a6",
		String: "A6",
		Width:  Length{Length: 105, Units: Millimeter},
		Height: Length{Length: 148, Units: Millimeter},
	},
}

// SizeFromString extracts a Size from a key.
func SizeFromString(key string) Size {
	for i, one := range sizeValues {
		if strings.EqualFold(key, one.Key) {
			return Size(i)
		}
	}
	return 0
}

// EnsureValid returns the first Size if this Size is not a known value.
func (s Size) EnsureValid() Size {
	if int(s) < len(sizeValues) {
		return s
	}
	return 0
}

// Key returns the key used to represent this ThresholdOp.
func (s Size) Key() string {
	return sizeValues[s.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (s Size) String() string {
	return sizeValues[s.EnsureValid()].String
}

// Dimensions returns the paper dimensions.
func (s Size) Dimensions() (width, height Length) {
	one := sizeValues[s.EnsureValid()]
	return one.Width, one.Height
}
