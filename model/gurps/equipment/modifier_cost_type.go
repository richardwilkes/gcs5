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

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Possible ModifierCostType values.
const (
	// OriginalCost modifies the original value stored in the equipment. Can be a ±value, ±% value, or a multiplier.
	OriginalCost ModifierCostType = iota
	// BaseCost modifies the base cost. Can be an additive multiplier or a CF value.
	BaseCost
	// FinalBaseCost modifies the final base cost. Can be a ±value, ±% value, or a multiplier.
	FinalBaseCost
	// FinalCost modifies the final cost. Can be a ±value, ±% value, or a multiplier.
	FinalCost
)

type modifierCostTypeData struct {
	Key         string
	Description string
	Example     string
	Permitted   []ModifierCostValueType
}

// ModifierCostType describes how an EquipmentModifier's cost is applied.
type ModifierCostType uint8

var modifierCostTypeValues = []*modifierCostTypeData{
	{
		Key:         "to_original_cost",
		Description: i18n.Text("to original cost"),
		Example:     `"+5", "-5", "+10%", "-10%", "x3.2"`,
		Permitted:   []ModifierCostValueType{Addition, Percentage, Multiplier},
	},
	{
		Key:         "to_base_cost",
		Description: i18n.Text("to base cost"),
		Example:     `"x2", "+2 CF", "-0.2 CF"`,
		Permitted:   []ModifierCostValueType{CostFactor, Multiplier},
	},
	{
		Key:         "to_final_base_cost",
		Description: i18n.Text("to final base cost"),
		Example:     `"+5", "-5", "+10%", "-10%", "x3.2"`,
		Permitted:   []ModifierCostValueType{Addition, Percentage, Multiplier},
	},
	{
		Key:         "to_final_cost",
		Description: i18n.Text("to final cost"),
		Example:     `"+5", "-5", "+10%", "-10%", "x3.2"`,
		Permitted:   []ModifierCostValueType{Addition, Percentage, Multiplier},
	},
}

// ModifierCostTypeFromKey extracts a ModifierCostType from a key.
func ModifierCostTypeFromKey(key string) ModifierCostType {
	for i, one := range modifierCostTypeValues {
		if strings.EqualFold(key, one.Key) {
			return ModifierCostType(i)
		}
	}
	return 0
}

// DetermineModifierCostValueTypeFromString examines a string to determine what type it is, but restricts the result to
// those allowed for this ModifierCostType.
func (m ModifierCostType) DetermineModifierCostValueTypeFromString(s string) ModifierCostValueType {
	t := DetermineModifierCostValueTypeFromString(s)
	permitted := modifierCostTypeValues[m.EnsureValid()].Permitted
	for _, one := range permitted {
		if one == t {
			return t
		}
	}
	return permitted[0]
}

// EnsureValid returns the first ModifierCostType if this ModifierCostType is not a known value.
func (m ModifierCostType) EnsureValid() ModifierCostType {
	if int(m) < len(modifierCostTypeValues) {
		return m
	}
	return 0
}

// Key returns the key used to represent this ModifierCostType.
func (m ModifierCostType) Key() string {
	return modifierCostTypeValues[m.EnsureValid()].Key
}

// ShortString returns the same thing as .String(), but without the example.
func (m ModifierCostType) ShortString() string {
	return modifierCostTypeValues[m.EnsureValid()].Description
}

// String implements fmt.Stringer.
func (m ModifierCostType) String() string {
	data := modifierCostTypeValues[m.EnsureValid()]
	return fmt.Sprintf("%s (e.g. %s)", data.Description, data.Example)
}

// ExtractValue from the string.
func (m ModifierCostType) ExtractValue(s string) fixed.F64d4 {
	return m.DetermineModifierCostValueTypeFromString(s).ExtractValue(s)
}

// Format returns a formatted version of the value.
func (m ModifierCostType) Format(s string) string {
	t := m.DetermineModifierCostValueTypeFromString(s)
	return t.Format(t.ExtractValue(s))
}
