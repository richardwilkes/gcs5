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

package equipment

import (
	"strings"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64"
)

// Format returns a formatted version of the value.
func (enum ModifierCostValueType) Format(value fxp.Int) string {
	switch enum {
	case Addition:
		return value.StringWithSign()
	case Percentage:
		return value.StringWithSign() + enum.String()
	case Multiplier:
		if value <= 0 {
			value = fxp.One
		}
		return enum.String() + value.String()
	case CostFactor:
		return value.StringWithSign() + " " + enum.String()
	default:
		return Addition.Format(value)
	}
}

// ExtractValue from the string.
func (enum ModifierCostValueType) ExtractValue(s string) fxp.Int {
	v := f64.FromStringForced[fxp.DP](s)
	if enum.EnsureValid() == Multiplier && v <= 0 {
		v = fxp.One
	}
	return v
}

// DetermineModifierCostValueTypeFromString examines a string to determine what type it is.
func DetermineModifierCostValueTypeFromString(s string) ModifierCostValueType {
	s = strings.ToLower(strings.TrimSpace(s))
	switch {
	case strings.HasSuffix(s, CostFactor.Key()):
		return CostFactor
	case strings.HasSuffix(s, Percentage.Key()):
		return Percentage
	case strings.HasPrefix(s, Multiplier.Key()) || strings.HasSuffix(s, Multiplier.Key()):
		return Multiplier
	default:
		return Addition
	}
}
