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
	"github.com/richardwilkes/gcs/model/gurps/equipment"
)

const (
	equipmentModifierCostKey       = "cost"
	equipmentModifierCostTypeKey   = "cost_type"
	equipmentModifierWeightKey     = "weight"
	equipmentModifierWeightTypeKey = "weight_type"
	equipmentModifierTechLevelKey  = "tech_level"
	equipmentModifierTypeKey       = "modifier"
)

// EquipmentModifier holds a modifier to a piece of Equipment.
type EquipmentModifier struct {
	Common
	TechLevel    string
	CostAmount   string
	WeightAmount string
	Features     []*Feature
	Children     []*EquipmentModifier
	CostType     equipment.ModifierCostType
	WeightType   equipment.ModifierWeightType
	Enabled      bool
}

// NewEquipmentModifierFromJSON creates a new EquipmentModifier from a JSON object.
func NewEquipmentModifierFromJSON(data map[string]interface{}) *EquipmentModifier {
	e := &EquipmentModifier{}
	e.Common.FromJSON(equipmentModifierTypeKey, data)
	if e.Container {
		e.Children = EquipmentModifiersListFromJSON(commonChildrenKey, data)
	} else {
		e.Enabled = !encoding.Bool(data[commonDisabledKey])
		if v, ok := data[equipmentModifierCostTypeKey]; ok {
			e.CostType = equipment.ModifierCostTypeFromKey(encoding.String(v))
			e.CostAmount = encoding.String(data[equipmentModifierCostKey])
		}
		if v, ok := data[equipmentModifierWeightTypeKey]; ok {
			e.WeightType = equipment.ModifierWeightTypeFromKey(encoding.String(v))
			e.WeightAmount = encoding.String(data[equipmentModifierWeightKey])
		}
		e.TechLevel = encoding.String(data[equipmentModifierTechLevelKey])
		e.Features = FeaturesListFromJSON(data)
	}
	return e
}

// ToJSON emits this object as JSON.
func (e *EquipmentModifier) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	e.Common.ToInlineJSON(equipmentModifierTypeKey, encoder)
	if e.Container {
		EquipmentModifiersListToJSON(commonChildrenKey, e.Children, encoder)
	} else {
		encoder.KeyedBool(commonDisabledKey, !e.Enabled, true)
		if e.CostType != equipment.OriginalCost || e.CostAmount != "+0" {
			encoder.KeyedString(equipmentModifierCostTypeKey, e.CostType.Key(), false, false)
			encoder.KeyedString(equipmentModifierCostKey, e.CostAmount, true, true)
		}
		if e.WeightType != equipment.OriginalWeight || (e.WeightAmount != "+0" && !strings.HasPrefix(e.WeightAmount, "+0 ")) {
			encoder.KeyedString(equipmentModifierWeightTypeKey, e.WeightType.Key(), false, false)
			encoder.KeyedString(equipmentModifierWeightKey, e.WeightAmount, true, true)
		}
		encoder.KeyedString(equipmentModifierTechLevelKey, e.TechLevel, true, true)
		FeaturesListToJSON(e.Features, encoder)
	}
	encoder.EndObject()
}

func (e *EquipmentModifier) String() string {
	return e.Name
}

// FullDescription returns a full description. 'entity' may be nil.
func (e *EquipmentModifier) FullDescription(entity *Entity) string {
	var buffer strings.Builder
	buffer.WriteString(e.String())
	if e.Notes != "" {
		buffer.WriteString(" (")
		buffer.WriteString(e.Notes)
		buffer.WriteByte(')')
	}
	if entity != nil && SheetSettingsFor(entity).ShowEquipmentModifierAdj {
		costDesc := e.CostDescription()
		weightDesc := e.WeightDescription(entity)
		if costDesc != "" || weightDesc != "" {
			buffer.WriteString(" [")
			buffer.WriteString(costDesc)
			if weightDesc != "" {
				if costDesc != "" {
					buffer.WriteString("; ")
				}
				buffer.WriteString(weightDesc)
			}
			buffer.WriteByte(']')
		}
	}
	return buffer.String()
}

// CostDescription returns the formatted cost.
func (e *EquipmentModifier) CostDescription() string {
	if e.Container || (e.CostType == equipment.OriginalCost && e.CostAmount == "+0") {
		return ""
	}
	return e.CostType.Format(e.CostAmount) + " " + e.CostType.ShortString()
}

// WeightDescription returns the formatted weight.
func (e *EquipmentModifier) WeightDescription(entity *Entity) string {
	if e.Container || (e.WeightType == equipment.OriginalWeight && (e.WeightAmount == "+0" || strings.HasPrefix(e.WeightAmount, "+0 "))) {
		return ""
	}
	return e.WeightType.Format(e.WeightAmount, SheetSettingsFor(entity).DefaultWeightUnits) + " " + e.WeightType.ShortString()
}

// FillWithNameableKeys adds any nameable keys found in this EquipmentModifier to the provided map.
func (e *EquipmentModifier) FillWithNameableKeys(nameables map[string]string) {
	if e.Enabled {
		e.Common.FillWithNameableKeys(nameables)
		for _, one := range e.Features {
			one.FillWithNameableKeys(nameables)
		}
	}
}

// ApplyNameableKeys replaces any nameable keys found in this EquipmentModifier with the corresponding values in the provided map.
func (e *EquipmentModifier) ApplyNameableKeys(nameables map[string]string) {
	if e.Enabled {
		e.Common.ApplyNameableKeys(nameables)
		for _, one := range e.Features {
			one.ApplyNameableKeys(nameables)
		}
	}
}
