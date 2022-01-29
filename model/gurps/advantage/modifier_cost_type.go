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

package advantage

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible ModifierCostType values.
const (
	Percentage ModifierCostType = iota // Adds to the percentage multiplier
	Points                             // Adds a constant to the base value prior to any multiplier or percentage adjustment
	Multiplier                         // Multiplies the final cost by a constant
)

type modifierCostTypeData struct {
	Key    string
	String string
}

// ModifierCostType eescribes how an AdvantageModifier's point cost is applied.
type ModifierCostType uint8

var modifierCostTypeValues = []*modifierCostTypeData{
	{
		Key:    "percentage",
		String: "%",
	},
	{
		Key:    "points",
		String: i18n.Text("points"),
	},
	{
		Key:    "multiplier",
		String: "×",
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

// EnsureValid returns the first ModifierCostType if this ModifierCostType is not a known value.
func (a ModifierCostType) EnsureValid() ModifierCostType {
	if int(a) < len(modifierCostTypeValues) {
		return a
	}
	return 0
}

// Key returns the key used to represent this ModifierCostType.
func (a ModifierCostType) Key() string {
	return modifierCostTypeValues[a.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (a ModifierCostType) String() string {
	return modifierCostTypeValues[a.EnsureValid()].String
}
