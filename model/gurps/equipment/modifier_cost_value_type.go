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
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Possible ModifierCostValueType values.
const (
	Addition   = ModifierCostValueType("+")
	Percentage = ModifierCostValueType("%")
	Multiplier = ModifierCostValueType("x")
	CostFactor = ModifierCostValueType("cf")
)

// AllModifierCostValueTypes is the complete set of ModifierCostValueType values.
var AllModifierCostValueTypes = []ModifierCostValueType{
	Addition,
	Percentage,
	Multiplier,
	CostFactor,
}

// ModifierCostValueType describes how an EquipmentModifier's point cost is applied.
type ModifierCostValueType string

// EnsureValid ensures this is of a known value.
func (m ModifierCostValueType) EnsureValid() ModifierCostValueType {
	for _, one := range AllModifierCostValueTypes {
		if one == m {
			return m
		}
	}
	return AllModifierCostValueTypes[0]
}

// Format returns a formatted version of the value.
func (m ModifierCostValueType) Format(value fixed.F64d4) string {
	switch m {
	case Addition:
		return value.StringWithSign()
	case Percentage:
		return value.StringWithSign() + "%"
	case Multiplier:
		if value <= 0 {
			value = fxp.One
		}
		return "x" + value.String()
	case CostFactor:
		return value.StringWithSign() + " CF"
	default:
		return Addition.Format(value)
	}
}

// ExtractValue from the string.
func (m ModifierCostValueType) ExtractValue(s string) fixed.F64d4 {
	v := fixed.F64d4FromStringForced(s)
	if m.EnsureValid() == Multiplier && v <= 0 {
		v = fxp.One
	}
	return v
}

// DetermineModifierCostValueTypeFromString examines a string to determine what type it is.
func DetermineModifierCostValueTypeFromString(s string) ModifierCostValueType {
	s = strings.ToLower(strings.TrimSpace(s))
	switch {
	case strings.HasSuffix(s, "cf"):
		return CostFactor
	case strings.HasSuffix(s, "%"):
		return Percentage
	case strings.HasPrefix(s, "x") || strings.HasSuffix(s, "x"):
		return Multiplier
	default:
		return Addition
	}
}
