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
)

// Possible ModifierWeightValueType values.
const (
	WeightAddition             = ModifierWeightValueType("+")
	WeightPercentageAdder      = ModifierWeightValueType("%")
	WeightPercentageMultiplier = ModifierWeightValueType("x%")
	WeightMultiplier           = ModifierWeightValueType("x")
)

// AllModifierWeightValueTypes is the complete set of ModifierWeightValueType values.
var AllModifierWeightValueTypes = []ModifierWeightValueType{
	WeightAddition,
	WeightPercentageAdder,
	WeightPercentageMultiplier,
	WeightMultiplier,
}

// ModifierWeightValueType describes how an EquipmentModifier's point cost is applied.
type ModifierWeightValueType string

// EnsureValid ensures this is of a known value.
func (m ModifierWeightValueType) EnsureValid() ModifierWeightValueType {
	for _, one := range AllModifierWeightValueTypes {
		if one == m {
			return m
		}
	}
	return AllModifierWeightValueTypes[0]
}

// Format returns a formatted version of the value.
func (m ModifierWeightValueType) Format(fraction fxp.Fraction) string {
	switch m {
	case WeightAddition:
		return fraction.StringWithSign()
	case WeightPercentageAdder:
		return fraction.StringWithSign() + "%"
	case WeightPercentageMultiplier:
		if fraction.Numerator <= 0 {
			fraction.Numerator = fxp.Hundred
			fraction.Denominator = fxp.One
		}
		return "x" + fraction.String() + "%"
	case WeightMultiplier:
		if fraction.Numerator <= 0 {
			fraction.Numerator = fxp.One
			fraction.Denominator = fxp.One
		}
		return "x" + fraction.String()
	default:
		return WeightAddition.Format(fraction)
	}
}

// ExtractFraction from the string.
func (m ModifierWeightValueType) ExtractFraction(s string) fxp.Fraction {
	fraction := fxp.NewFractionFromString(s)
	revised := m.EnsureValid()
	switch revised {
	case WeightPercentageMultiplier:
		if fraction.Numerator <= 0 {
			fraction.Numerator = fxp.Hundred
			fraction.Denominator = fxp.One
		}
	case WeightMultiplier:
		if fraction.Numerator <= 0 {
			fraction.Numerator = fxp.One
			fraction.Denominator = fxp.One
		}
	default:
	}
	return fraction
}

// DetermineModifierWeightValueTypeFromString examines a string to determine what type it is.
func DetermineModifierWeightValueTypeFromString(s string) ModifierWeightValueType {
	s = strings.ToLower(strings.TrimSpace(s))
	switch {
	case strings.HasSuffix(s, "%"):
		if strings.HasPrefix(s, "x") {
			return WeightPercentageMultiplier
		}
		return WeightPercentageAdder
	case strings.HasPrefix(s, "x") || strings.HasSuffix(s, "x"):
		return WeightMultiplier
	default:
		return WeightAddition
	}
}
