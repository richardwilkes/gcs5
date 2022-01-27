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

import "strings"

// Possible Units values.
const (
	Inch Units = iota
	Centimeter
	Millimeter
)

type unitsData struct {
	Key      string
	ToPixels func(value float64) float32
}

// Units holds the real-world length unit type.
type Units uint8

var unitsValues = []*unitsData{
	{
		Key:      "in",
		ToPixels: func(length float64) float32 { return float32(length * 72) },
	},
	{
		Key:      "cm",
		ToPixels: func(length float64) float32 { return float32((length * 72) / 2.54) },
	},
	{
		Key:      "mm",
		ToPixels: func(length float64) float32 { return float32((length * 72) / 25.4) },
	},
}

// UnitsFromString extracts a Units from a key.
func UnitsFromString(key string) Units {
	for i, one := range unitsValues {
		if strings.EqualFold(key, one.Key) {
			return Units(i)
		}
	}
	return 0
}

// EnsureValid returns the first Units if this Units is not a known value.
func (u Units) EnsureValid() Units {
	if int(u) < len(unitsValues) {
		return u
	}
	return 0
}

// Key returns the key used to represent this.
func (u Units) Key() string {
	return unitsValues[u.EnsureValid()].Key
}

// ToPixels converts the given length in this Units to the number of 72-pixels-per-inch pixels it represents.
func (u Units) ToPixels(length float64) float32 {
	return unitsValues[u.EnsureValid()].ToPixels(length)
}
