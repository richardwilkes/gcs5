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
)

// Possible ModifierWeightValueType values.
const (
	WeightAddition ModifierWeightValueType = iota
	WeightPercentageAdder
	WeightPercentageMultiplier
	WeightMultiplier
)

type modifierWeightValueTypeData struct {
	Format func(fraction f64d4.Fraction) string
}

// ModifierWeightValueType describes how an EquipmentModifier's point cost is applied.
type ModifierWeightValueType uint8

var modifierWeightValueTypeValues = []*modifierWeightValueTypeData{
	{
		Format: func(fraction f64d4.Fraction) string { return fraction.StringWithSign() },
	},
	{
		Format: func(fraction f64d4.Fraction) string { return fraction.StringWithSign() + "%" },
	},
	{
		Format: func(fraction f64d4.Fraction) string {
			if fraction.Numerator <= 0 {
				fraction.Numerator = f64d4.Hundred
				fraction.Denominator = f64d4.One
			}
			return "x" + fraction.String() + "%"
		},
	},
	{
		Format: func(fraction f64d4.Fraction) string {
			if fraction.Numerator <= 0 {
				fraction.Numerator = f64d4.One
				fraction.Denominator = f64d4.One
			}
			return "x" + fraction.String()
		},
	},
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

// EnsureValid returns the first ModifierWeightValueType if this ModifierWeightValueType is not a known value.
func (m ModifierWeightValueType) EnsureValid() ModifierWeightValueType {
	if int(m) < len(modifierWeightValueTypeValues) {
		return m
	}
	return 0
}

// Format returns a formatted version of the value.
func (m ModifierWeightValueType) Format(fraction f64d4.Fraction) string {
	return modifierWeightValueTypeValues[m.EnsureValid()].Format(fraction)
}

// ExtractFraction from the string.
func (m ModifierWeightValueType) ExtractFraction(s string) f64d4.Fraction {
	fraction := f64d4.NewFractionFromString(s)
	revised := m.EnsureValid()
	switch revised {
	case WeightPercentageMultiplier:
		if fraction.Numerator <= 0 {
			fraction.Numerator = f64d4.Hundred
			fraction.Denominator = f64d4.One
		}
	case WeightMultiplier:
		if fraction.Numerator <= 0 {
			fraction.Numerator = f64d4.One
			fraction.Denominator = f64d4.One
		}
	default:
	}
	return fraction
}
