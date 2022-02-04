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

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible ModifierWeightType values.
const (
	// OriginalWeight modifies the original value stored in the equipment. Can be a ±value or a ±% value.
	OriginalWeight = ModifierWeightType("to_original_weight")
	// BaseWeight modifies the base weight. Can be a ±value, a ±% value, or a multiplier.
	BaseWeight = ModifierWeightType("to_base_weight")
	// FinalBaseWeight modifies the final base weight. Can be a ±value, a ±% value, or a multiplier.
	FinalBaseWeight = ModifierWeightType("to_final_base_weight")
	// FinalWeight modifies the final weight. Can be a ±value, a ±% value, or a multiplier.
	FinalWeight = ModifierWeightType("to_final_weight")
)

// AllModifierWeightTypes is the complete set of ModifierWeightType values.
var AllModifierWeightTypes = []ModifierWeightType{
	OriginalWeight,
	BaseWeight,
	FinalBaseWeight,
	FinalWeight,
}

// ModifierWeightType describes how an EquipmentModifier's cost is applied.
type ModifierWeightType string

// EnsureValid ensures this is of a known value.
func (m ModifierWeightType) EnsureValid() ModifierWeightType {
	for _, one := range AllModifierWeightTypes {
		if one == m {
			return m
		}
	}
	return AllModifierWeightTypes[0]
}

// ShortString returns the same thing as .String(), but without the example.
func (m ModifierWeightType) ShortString() string {
	switch m {
	case OriginalWeight:
		return i18n.Text("to original weight")
	case BaseWeight:
		return i18n.Text("to base weight")
	case FinalBaseWeight:
		return i18n.Text("to final base weight")
	case FinalWeight:
		return i18n.Text("to final weight")
	default:
		return OriginalWeight.ShortString()
	}
}

// String implements fmt.Stringer.
func (m ModifierWeightType) String() string {
	return fmt.Sprintf("%s (e.g. %s)", m.ShortString(), m.Example())
}

// Example returns example values.
func (m ModifierWeightType) Example() string {
	if m.EnsureValid() == OriginalWeight {
		return `"+5 lb", "-5 lb", "+10%", "-10%"`
	}
	return `"+5 lb", "-5 lb", "x10%", "x3", "x2/3"`
}

// Permitted returns the permitted ModifierCostValueType values.
func (m ModifierWeightType) Permitted() []ModifierWeightValueType {
	if m.EnsureValid() == OriginalWeight {
		return []ModifierWeightValueType{WeightAddition, WeightPercentageAdder}
	}
	return []ModifierWeightValueType{WeightAddition, WeightPercentageMultiplier, WeightMultiplier}
}

// DetermineModifierWeightValueTypeFromString examines a string to determine what type it is, but restricts the result to
// those allowed for this ModifierWeightType.
func (m ModifierWeightType) DetermineModifierWeightValueTypeFromString(s string) ModifierWeightValueType {
	mvt := DetermineModifierWeightValueTypeFromString(s)
	permitted := m.Permitted()
	for _, one := range permitted {
		if one == mvt {
			return mvt
		}
	}
	return permitted[0]
}

// ExtractFraction from the string.
func (m ModifierWeightType) ExtractFraction(s string) fxp.Fraction {
	return m.DetermineModifierWeightValueTypeFromString(s).ExtractFraction(s)
}

// Format returns a formatted version of the value.
func (m ModifierWeightType) Format(s string, defUnits measure.WeightUnits) string {
	t := m.DetermineModifierWeightValueTypeFromString(s)
	result := t.Format(t.ExtractFraction(s))
	if t == WeightAddition {
		result += " " + string(measure.TrailingWeightUnitsFromString(s, defUnits))
	}
	return result
}
