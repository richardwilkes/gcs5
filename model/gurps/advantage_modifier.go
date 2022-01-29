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
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	advantageModifierAffectsKey  = "affects"
	advantageModifierCostKey     = "cost"
	advantageModifierCostTypeKey = "cost_type"
	advantageModifierLevelsKey   = "levels"
	advantageModifierTypeKey     = "modifier"
)

// AdvantageModifier holds a modifier to an Advantage.
type AdvantageModifier struct {
	Common
	Cost     fixed.F64d4
	Levels   fixed.F64d4
	Features []*Feature
	Children []*AdvantageModifier
	CostType AdvantageModifierCostType
	Affects  Affects
	Enabled  bool
}

// NewAdvantageModifierFromJSON creates a new AdvantageModifier from a JSON object.
func NewAdvantageModifierFromJSON(data map[string]interface{}) *AdvantageModifier {
	a := &AdvantageModifier{}
	a.Common.FromJSON(advantageModifierTypeKey, data)
	if a.Container {
		a.Children = AdvantageModifiersListFromJSON(commonChildrenKey, data)
	} else {
		a.Enabled = !encoding.Bool(data[commonDisabledKey])
		a.CostType = AdvantageModifierCostTypeFromKey(encoding.String(data[advantageModifierCostTypeKey]))
		a.Cost = encoding.Number(data[advantageModifierCostKey])
		if a.CostType != Multiplier {
			a.Affects = AffectsFromKey(encoding.String(data[advantageModifierAffectsKey]))
		}
		a.Levels = encoding.Number(data[advantageModifierLevelsKey])
		a.Features = FeaturesListFromJSON(data)
	}
	return a
}

// ToJSON emits this object as JSON.
func (a *AdvantageModifier) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	a.Common.ToInlineJSON(advantageModifierTypeKey, encoder)
	if a.Container {
		AdvantageModifiersListToJSON(commonChildrenKey, a.Children, encoder)
	} else {
		encoder.KeyedBool(commonDisabledKey, !a.Enabled, true)
		encoder.KeyedString(advantageModifierCostTypeKey, a.CostType.Key(), false, false)
		encoder.KeyedNumber(advantageModifierCostKey, a.Cost, false)
		if a.CostType != Multiplier {
			encoder.KeyedString(advantageModifierAffectsKey, a.Affects.Key(), false, false)
		}
		encoder.KeyedNumber(advantageModifierLevelsKey, a.Levels, true)
		FeaturesListToJSON(a.Features, encoder)
	}
	encoder.EndObject()
}

// CostModifier returns the total cost modifier.
func (a *AdvantageModifier) CostModifier() fixed.F64d4 {
	if a.Levels > 0 {
		return a.Cost.Mul(a.Levels)
	}
	return a.Cost
}
