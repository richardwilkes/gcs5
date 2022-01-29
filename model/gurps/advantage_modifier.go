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

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/advantage"
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
	CostType advantage.ModifierCostType
	Affects  advantage.Affects
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
		a.CostType = advantage.ModifierCostTypeFromKey(encoding.String(data[advantageModifierCostTypeKey]))
		a.Cost = encoding.Number(data[advantageModifierCostKey])
		if a.CostType != advantage.Multiplier {
			a.Affects = advantage.AffectsFromKey(encoding.String(data[advantageModifierAffectsKey]))
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
		if a.CostType != advantage.Multiplier {
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

// HasLevels returns true if this AdvantageModifier has levels.
func (a *AdvantageModifier) HasLevels() bool {
	return a.CostType == advantage.Percentage && a.Levels > 0
}

func (a *AdvantageModifier) String() string {
	var buffer strings.Builder
	buffer.WriteString(a.Name)
	if a.HasLevels() {
		buffer.WriteByte(' ')
		buffer.WriteString(a.Levels.String())
	}
	return buffer.String()
}

// FullDescription returns a full description. 'entity' may be nil.
func (a *AdvantageModifier) FullDescription(entity *Entity) string {
	var buffer strings.Builder
	buffer.WriteString(a.String())
	if a.Notes != "" {
		buffer.WriteString(" (")
		buffer.WriteString(a.Notes)
		buffer.WriteByte(')')
	}
	if entity != nil && SheetSettingsFor(entity).ShowAdvantageModifierAdj {
		buffer.WriteString(" [")
		buffer.WriteString(a.CostDescription())
		buffer.WriteByte(']')
	}
	return buffer.String()
}

// CostDescription returns the formatted cost.
func (a *AdvantageModifier) CostDescription() string {
	var buffer strings.Builder
	if a.CostType == advantage.Multiplier {
		buffer.WriteByte('x')
		buffer.WriteString(a.Cost.String())
	} else {
		if a.Cost >= 0 {
			buffer.WriteByte('+')
		}
		buffer.WriteString(a.Cost.String())
		if a.CostType == advantage.Percentage {
			buffer.WriteByte('%')
		}
		if desc := a.Affects.ShortTitle(); desc != "" {
			buffer.WriteByte(' ')
			buffer.WriteString(desc)
		}
	}
	return buffer.String()
}

// FillWithNameableKeys adds any nameable keys found in this AdvantageModifier to the provided map.
func (a *AdvantageModifier) FillWithNameableKeys(nameables map[string]string) {
	if a.Enabled {
		a.Common.FillWithNameableKeys(nameables)
		for _, one := range a.Features {
			one.FillWithNameableKeys(nameables)
		}
	}
}

// ApplyNameableKeys replaces any nameable keys found in this AdvantageModifier with the corresponding values in the provided map.
func (a *AdvantageModifier) ApplyNameableKeys(nameables map[string]string) {
	if a.Enabled {
		a.Common.ApplyNameableKeys(nameables)
		for _, one := range a.Features {
			one.ApplyNameableKeys(nameables)
		}
	}
}