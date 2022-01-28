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

package gurps

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible AdvantageModifierCostType values.
const (
	Percentage AdvantageModifierCostType = iota // Adds to the percentage multiplier
	Points                                      // Adds a constant to the base value prior to any multiplier or percentage adjustment
	Multiplier                                  // Multiplies the final cost by a constant
)

type advantageModifierCostTypeData struct {
	Key    string
	String string
}

// AdvantageModifierCostType eescribes how an AdvantageModifier's point cost is applied.
type AdvantageModifierCostType uint8

var advantageModifierCostTypeValues = []*advantageModifierCostTypeData{
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

// AdvantageModifierCostTypeFromKey extracts a AdvantageModifierCostType from a key.
func AdvantageModifierCostTypeFromKey(key string) AdvantageModifierCostType {
	for i, one := range advantageModifierCostTypeValues {
		if strings.EqualFold(key, one.Key) {
			return AdvantageModifierCostType(i)
		}
	}
	return 0
}

// EnsureValid returns the first AdvantageModifierCostType if this AdvantageModifierCostType is not a known value.
func (a AdvantageModifierCostType) EnsureValid() AdvantageModifierCostType {
	if int(a) < len(advantageModifierCostTypeValues) {
		return a
	}
	return 0
}

// Key returns the key used to represent this AdvantageModifierCostType.
func (a AdvantageModifierCostType) Key() string {
	return advantageModifierCostTypeValues[a.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (a AdvantageModifierCostType) String() string {
	return advantageModifierCostTypeValues[a.EnsureValid()].String
}
