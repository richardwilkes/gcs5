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

// TrailingWeightUnitsFromString extracts a trailing WeightUnits from a string.
func TrailingWeightUnitsFromString(s string, defUnits WeightUnits) WeightUnits {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, one := range AllWeightUnits {
		if strings.HasSuffix(s, one.Key()) {
			return one
		}
	}
	return defUnits
}

// Format the weight for this WeightUnits.
func (enum WeightUnits) Format(weight Weight) string {
	switch enum {
	case Pound, PoundAlt:
		return fixed.F64d4(weight).String() + " " + enum.Key()
	case Ounce:
		return fixed.F64d4(weight).Mul(fixed.F64d4FromInt(16)).String() + " " + enum.Key()
	case Ton:
		return fixed.F64d4(weight).Div(fixed.F64d4FromInt(2000)).String() + " " + enum.Key()
	case Kilogram:
		return fixed.F64d4(weight).Div(fixed.F64d4FromInt(2)).String() + " " + enum.Key()
	case Gram:
		return fixed.F64d4(weight).Mul(fixed.F64d4FromInt(500)).String() + " " + enum.Key()
	default:
		return Pound.Format(weight)
	}
}

// ToPounds the weight for this WeightUnits.
func (enum WeightUnits) ToPounds(weight fixed.F64d4) fixed.F64d4 {
	switch enum {
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
