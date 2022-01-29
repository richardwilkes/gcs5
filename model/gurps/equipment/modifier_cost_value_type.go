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

	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Possible ModifierCostValueType values.
const (
	Addition ModifierCostValueType = iota
	Percentage
	Multiplier
	CostFactor
)

type modifierCostValueTypeData struct {
	Format func(value fixed.F64d4) string
}

// ModifierCostValueType describes how an EquipmentModifier's point cost is applied.
type ModifierCostValueType uint8

var modifierCostValueTypeValues = []*modifierCostValueTypeData{
	{
		Format: func(value fixed.F64d4) string { return value.StringWithSign() },
	},
	{
		Format: func(value fixed.F64d4) string { return value.StringWithSign() + "%" },
	},
	{
		Format: func(value fixed.F64d4) string {
			if value <= 0 {
				value = f64d4.One
			}
			return "x" + value.String()
		},
	},
	{
		Format: func(value fixed.F64d4) string { return value.StringWithSign() + " CF" },
	},
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

// EnsureValid returns the first ModifierCostValueType if this ModifierCostValueType is not a known value.
func (a ModifierCostValueType) EnsureValid() ModifierCostValueType {
	if int(a) < len(modifierCostValueTypeValues) {
		return a
	}
	return 0
}

// Format returns a formatted version of the value.
func (a ModifierCostValueType) Format(value fixed.F64d4) string {
	return modifierCostValueTypeValues[a.EnsureValid()].Format(value)
}

// ExtractValue from the string.
func (a ModifierCostValueType) ExtractValue(s string) fixed.F64d4 {
	v := fixed.F64d4FromStringForced(s)
	if a.EnsureValid() == Multiplier && v <= 0 {
		v = f64d4.One
	}
	return v
}
