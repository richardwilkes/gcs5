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
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible ModifierCostType values.
const (
	// Percentage adds to the percentage multiplier.
	Percentage = ModifierCostType("percentage")
	// Points adds a constant to the base value prior to any multiplier or percentage adjustment.
	Points = ModifierCostType("points")
	// Multiplier multiplies the final cost by a constant.
	Multiplier = ModifierCostType("multiplier")
)

// AllModifierCostTypes is the complete set of ModifierCostType values.
var AllModifierCostTypes = []ModifierCostType{
	Percentage,
	Points,
	Multiplier,
}

// ModifierCostType describes how an AdvantageModifier's point cost is applied.
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

// String implements fmt.Stringer.
func (m ModifierCostType) String() string {
	switch m {
	case Percentage:
		return "%"
	case Points:
		return i18n.Text("points")
	case Multiplier:
		return "×"
	default:
		return Percentage.String()
	}
}
