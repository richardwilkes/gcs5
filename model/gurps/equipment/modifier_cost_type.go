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

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Possible ModifierCostType values.
const (
	// OriginalCost modifies the original value stored in the equipment. Can be a ±value, ±% value, or a multiplier.
	OriginalCost = ModifierCostType("to_original_cost")
	// BaseCost modifies the base cost. Can be an additive multiplier or a CF value.
	BaseCost = ModifierCostType("to_base_cost")
	// FinalBaseCost modifies the final base cost. Can be a ±value, ±% value, or a multiplier.
	FinalBaseCost = ModifierCostType("to_final_base_cost")
	// FinalCost modifies the final cost. Can be a ±value, ±% value, or a multiplier.
	FinalCost = ModifierCostType("to_final_cost")
)

// AllModifierCostTypes is the complete set of ModifierCostType values.
var AllModifierCostTypes = []ModifierCostType{
	OriginalCost,
	BaseCost,
	FinalBaseCost,
	FinalCost,
}

// ModifierCostType describes how an EquipmentModifier's cost is applied.
type ModifierCostType string

// EnsureValid ensures this is of a known value.
func (m ModifierCostType) EnsureValid() ModifierCostType {
	for _, one := range AllModifierCostTypes {
		if one == m {
			return m
		}
	}
	return AllModifierCostTypes[0]
}

// ShortString returns the same thing as .String(), but without the example.
func (m ModifierCostType) ShortString() string {
	switch m {
	case OriginalCost:
		return i18n.Text("to original cost")
	case BaseCost:
		return i18n.Text("to base cost")
	case FinalBaseCost:
		return i18n.Text("to final base cost")
	case FinalCost:
		return i18n.Text("to final cost")
	default:
		return OriginalCost.String()
	}
}

// String implements fmt.Stringer.
func (m ModifierCostType) String() string {
	return fmt.Sprintf("%s (e.g. %s)", m.String(), m.Example())
}

// Example returns example values.
func (m ModifierCostType) Example() string {
	if m.EnsureValid() == BaseCost {
		return `"x2", "+2 CF", "-0.2 CF"`
	}
	return `"+5", "-5", "+10%", "-10%", "x3.2"`
}

// Permitted returns the permitted ModifierCostValueType values.
func (m ModifierCostType) Permitted() []ModifierCostValueType {
	if m.EnsureValid() == BaseCost {
		return []ModifierCostValueType{CostFactor, Multiplier}
	}
	return []ModifierCostValueType{Addition, Percentage, Multiplier}
}

// DetermineModifierCostValueTypeFromString examines a string to determine what type it is, but restricts the result to
// those allowed for this ModifierCostType.
func (m ModifierCostType) DetermineModifierCostValueTypeFromString(s string) ModifierCostValueType {
	cvt := DetermineModifierCostValueTypeFromString(s)
	permitted := m.Permitted()
	for _, one := range permitted {
		if one == cvt {
			return cvt
		}
	}
	return permitted[0]
}

// ExtractValue from the string.
func (m ModifierCostType) ExtractValue(s string) fixed.F64d4 {
	return m.DetermineModifierCostValueTypeFromString(s).ExtractValue(s)
}

// Format returns a formatted version of the value.
func (m ModifierCostType) Format(s string) string {
	cvt := m.DetermineModifierCostValueTypeFromString(s)
	return cvt.Format(cvt.ExtractValue(s))
}
