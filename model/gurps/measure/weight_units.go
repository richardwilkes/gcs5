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
	Pound    = WeightUnits("#")
	PoundAlt = WeightUnits("lb")
	Ounce    = WeightUnits("oz")
	Ton      = WeightUnits("tn")
	Kilogram = WeightUnits("kg")
	Gram     = WeightUnits("g") // must come after Kilogram, as it's abbreviation is a subset
)

// AllWeightUnits is the complete set of WeightUnits values.
var AllWeightUnits = []WeightUnits{
	Pound,
	PoundAlt,
	Ounce,
	Ton,
	Kilogram,
	Gram,
}

// WeightUnits holds the weight unit type. Note that conversions to/from metric are done using the simplified GURPS
// metric conversion of 1# = 0.5kg. For consistency, all metric weights are converted to kilograms, then to pounds,
// rather than the variations at different weights that the GURPS rules suggest.
type WeightUnits string

// TrailingWeightUnitsFromString extracts a trailing WeightUnits from a string.
func TrailingWeightUnitsFromString(s string, defUnits WeightUnits) WeightUnits {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, one := range AllWeightUnits {
		if strings.HasSuffix(s, string(one)) {
			return one
		}
	}
	return defUnits
}

// EnsureValid ensures this is of a known value.
func (w WeightUnits) EnsureValid() WeightUnits {
	for _, one := range AllWeightUnits {
		if one == w {
			return w
		}
	}
	return AllWeightUnits[0]
}

// Format the weight for this WeightUnits.
func (w WeightUnits) Format(weight Weight) string {
	switch w {
	case Pound, PoundAlt:
		return fixed.F64d4(weight).String() + string(w)
	case Ounce:
		return fixed.F64d4(weight).Mul(fixed.F64d4FromInt(16)).String() + string(w)
	case Ton:
		return fixed.F64d4(weight).Div(fixed.F64d4FromInt(2000)).String() + string(w)
	case Kilogram:
		return fixed.F64d4(weight).Div(fixed.F64d4FromInt(2)).String() + string(w)
	case Gram:
		return fixed.F64d4(weight).Mul(fixed.F64d4FromInt(500)).String() + string(w)
	default:
		return Pound.Format(weight)
	}
}

// ToPounds the weight for this WeightUnits.
func (w WeightUnits) ToPounds(weight fixed.F64d4) fixed.F64d4 {
	switch w {
	case Pound, PoundAlt:
		return weight
	case Ounce:
		return weight.Div(fixed.F64d4FromInt(16))
	case Ton:
		return weight.Mul(fixed.F64d4FromInt(2000))
	case Kilogram:
		return weight.Mul(fixed.F64d4FromInt(2))
	case Gram:
		return weight.Div(fixed.F64d4FromInt(500))
	default:
		return Pound.ToPounds(weight)
	}
}
