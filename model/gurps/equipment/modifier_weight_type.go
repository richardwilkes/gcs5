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
	"fmt"
	"strings"

	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible ModifierWeightType values.
const (
	// OriginalWeight modifies the original value stored in the equipment. Can be a ±value or a ±% value.
	OriginalWeight ModifierWeightType = iota
	// BaseWeight modifies the base weight. Can be a ±value, a ±% value, or a multiplier.
	BaseWeight
	// FinalBaseWeight modifies the final base weight. Can be a ±value, a ±% value, or a multiplier.
	FinalBaseWeight
	// FinalWeight modifies the final weight. Can be a ±value, a ±% value, or a multiplier.
	FinalWeight
)

type modifierWeightTypeData struct {
	Key         string
	Description string
	Example     string
	Permitted   []ModifierWeightValueType
}

// ModifierWeightType describes how an EquipmentModifier's cost is applied.
type ModifierWeightType uint8

var modifierWeightTypeValues = []*modifierWeightTypeData{
	{
		Key:         "to_original_weight",
		Description: i18n.Text("to original weight"),
		Example:     `"+5 lb", "-5 lb", "+10%", "-10%"`,
		Permitted:   []ModifierWeightValueType{WeightAddition, WeightPercentageAdder},
	},
	{
		Key:         "to_base_weight",
		Description: i18n.Text("to base weight"),
		Example:     `"+5 lb", "-5 lb", "x10%", "x3", "x2/3"`,
		Permitted:   []ModifierWeightValueType{WeightAddition, WeightPercentageMultiplier, WeightMultiplier},
	},
	{
		Key:         "to_final_base_weight",
		Description: i18n.Text("to final base weight"),
		Example:     `"+5 lb", "-5 lb", "x10%", "x3", "x2/3"`,
		Permitted:   []ModifierWeightValueType{WeightAddition, WeightPercentageMultiplier, WeightMultiplier},
	},
	{
		Key:         "to_final_weight",
		Description: i18n.Text("to final weight"),
		Example:     `"+5 lb", "-5 lb", "x10%", "x3", "x2/3"`,
		Permitted:   []ModifierWeightValueType{WeightAddition, WeightPercentageMultiplier, WeightMultiplier},
	},
}

// ModifierWeightTypeFromKey extracts a ModifierWeightType from a key.
func ModifierWeightTypeFromKey(key string) ModifierWeightType {
	for i, one := range modifierWeightTypeValues {
		if strings.EqualFold(key, one.Key) {
			return ModifierWeightType(i)
		}
	}
	return 0
}

// DetermineModifierWeightValueTypeFromString examines a string to determine what type it is, but restricts the result to
// those allowed for this ModifierWeightType.
func (a ModifierWeightType) DetermineModifierWeightValueTypeFromString(s string) ModifierWeightValueType {
	t := DetermineModifierWeightValueTypeFromString(s)
	permitted := modifierWeightTypeValues[a.EnsureValid()].Permitted
	for _, one := range permitted {
		if one == t {
			return t
		}
	}
	return permitted[0]
}

// EnsureValid returns the first ModifierWeightType if this ModifierWeightType is not a known value.
func (a ModifierWeightType) EnsureValid() ModifierWeightType {
	if int(a) < len(modifierWeightTypeValues) {
		return a
	}
	return 0
}

// Key returns the key used to represent this ModifierWeightType.
func (a ModifierWeightType) Key() string {
	return modifierWeightTypeValues[a.EnsureValid()].Key
}

// ShortString returns the same thing as .String(), but without the example.
func (a ModifierWeightType) ShortString() string {
	return modifierWeightTypeValues[a.EnsureValid()].Description
}

// String implements fmt.Stringer.
func (a ModifierWeightType) String() string {
	data := modifierWeightTypeValues[a.EnsureValid()]
	return fmt.Sprintf("%s (e.g. %s)", data.Description, data.Example)
}

// ExtractFraction from the string.
func (a ModifierWeightType) ExtractFraction(s string) f64d4.Fraction {
	return a.DetermineModifierWeightValueTypeFromString(s).ExtractFraction(s)
}

// Format returns a formatted version of the value.
func (a ModifierWeightType) Format(s string, defUnits measure.WeightUnits) string {
	t := a.DetermineModifierWeightValueTypeFromString(s)
	result := t.Format(t.ExtractFraction(s))
	if t == WeightAddition {
		result += " " + measure.TrailingWeightUnitsFromString(s, defUnits).Key()
	}
	return result
}
