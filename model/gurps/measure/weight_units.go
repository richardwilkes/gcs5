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
	"strings"

	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Possible WeightUnits values.
const (
	Pound WeightUnits = iota
	PoundAlt
	Ounce
	Ton
	Kilogram
	Gram // must come after Kilogram, as it's abbreviation is a subset
)

type weightUnitsData struct {
	Key      string
	Format   func(weight Weight) string
	ToPounds func(value fixed.F64d4) fixed.F64d4
}

// WeightUnits holds the weight unit type. Note that conversions to/from metric are done using the simplified GURPS
// metric conversion of 1# = 0.5kg. For consistency, all metric weights are converted to kilograms, then to pounds,
// rather than the variations at different weights that the GURPS rules suggest.
type WeightUnits uint8

var weightUnitsValues = []*weightUnitsData{
	{
		Key:      "#",
		Format:   func(weight Weight) string { return fixed.F64d4(weight).String() + "#" },
		ToPounds: func(value fixed.F64d4) fixed.F64d4 { return value },
	},
	{
		Key:      "lb",
		Format:   func(weight Weight) string { return fixed.F64d4(weight).String() + "lb" },
		ToPounds: func(value fixed.F64d4) fixed.F64d4 { return value },
	},
	{
		Key:      "oz",
		Format:   func(weight Weight) string { return fixed.F64d4(weight).Mul(fixed.F64d4FromInt64(16)).String() + "oz" },
		ToPounds: func(value fixed.F64d4) fixed.F64d4 { return value.Div(fixed.F64d4FromInt64(16)) },
	},
	{
		Key:      "tn",
		Format:   func(weight Weight) string { return fixed.F64d4(weight).Div(fixed.F64d4FromInt64(2000)).String() + "tn" },
		ToPounds: func(value fixed.F64d4) fixed.F64d4 { return value.Mul(fixed.F64d4FromInt64(2000)) },
	},
	{
		Key:      "kg",
		Format:   func(weight Weight) string { return fixed.F64d4(weight).Div(fixed.F64d4FromInt64(2)).String() + "kg" },
		ToPounds: func(value fixed.F64d4) fixed.F64d4 { return value.Mul(fixed.F64d4FromInt64(2)) },
	},
	{
		Key:      "g",
		Format:   func(weight Weight) string { return fixed.F64d4(weight).Mul(fixed.F64d4FromInt64(500)).String() + "g" },
		ToPounds: func(value fixed.F64d4) fixed.F64d4 { return value.Div(fixed.F64d4FromInt64(500)) },
	},
}

// WeightUnitsFromString extracts a WeightUnits from a key.
func WeightUnitsFromString(key string) WeightUnits {
	for i, one := range weightUnitsValues {
		if strings.EqualFold(key, one.Key) {
			return WeightUnits(i)
		}
	}
	return 0
}

// EnsureValid returns the first WeightUnits if this WeightUnits is not a known value.
func (w WeightUnits) EnsureValid() WeightUnits {
	if int(w) < len(weightUnitsValues) {
		return w
	}
	return 0
}

// Key returns the key used to represent this WeightUnits.
func (w WeightUnits) Key() string {
	return weightUnitsValues[w.EnsureValid()].Key
}

// Format the weight for this WeightUnits.
func (w WeightUnits) Format(weight Weight) string {
	return weightUnitsValues[w.EnsureValid()].Format(weight)
}

// ToPounds the weight for this WeightUnits.
func (w WeightUnits) ToPounds(weight fixed.F64d4) fixed.F64d4 {
	return weightUnitsValues[w.EnsureValid()].ToPounds(weight)
}
