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

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/equipment"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const equipmentModifierTypeKey = "modifier"

// EquipmentModifierItem holds the EquipmentModifier data that only exists in non-containers.
type EquipmentModifierItem struct {
	TechLevel    string                       `json:"tech_level,omitempty"`
	CostType     equipment.ModifierCostType   `json:"cost_type,omitempty"`
	CostAmount   string                       `json:"cost,omitempty"`
	WeightType   equipment.ModifierWeightType `json:"weight_type,omitempty"`
	WeightAmount string                       `json:"weight,omitempty"`
	Features     feature.Features             `json:"features,omitempty"`
	Disabled     bool                         `json:"disabled,omitempty"`
}

// EquipmentModifierContainer holds the EquipmentModifier data that only exists in containers.
type EquipmentModifierContainer struct {
	Children []*EquipmentModifier `json:"children,omitempty"`
	Open     bool                 `json:"open,omitempty"`
}

// EquipmentModifierData holds the EquipmentModifier data that is written to disk.
type EquipmentModifierData struct {
	Type                        string    `json:"type"`
	ID                          uuid.UUID `json:"id"`
	Name                        string    `json:"name,omitempty"`
	PageRef                     string    `json:"reference,omitempty"`
	Notes                       string    `json:"notes,omitempty"`
	VTTNotes                    string    `json:"vtt_notes,omitempty"`
	*EquipmentModifierItem      `json:",omitempty"`
	*EquipmentModifierContainer `json:",omitempty"`
}

// EquipmentModifier holds a modifier to a piece of Equipment.
type EquipmentModifier struct {
	EquipmentModifierData
}

// NewEquipmentModifier creates an EquipmentModifier.
func NewEquipmentModifier(container bool) *EquipmentModifier {
	a := EquipmentModifier{
		EquipmentModifierData: EquipmentModifierData{
			Type: equipmentModifierTypeKey,
			ID:   id.NewUUID(),
			Name: i18n.Text("Advantage Modifier"),
		},
	}
	if container {
		a.Type += commonContainerKeyPostfix
		a.EquipmentModifierContainer = &EquipmentModifierContainer{Open: true}
	} else {
		a.EquipmentModifierItem = &EquipmentModifierItem{
			CostType:   equipment.OriginalCost,
			WeightType: equipment.OriginalWeight,
		}
	}
	return &a
}

// MarshalJSON implements json.Marshaler.
func (e *EquipmentModifier) MarshalJSON() ([]byte, error) {
	if e.Container() {
		e.EquipmentModifierItem = nil
	} else {
		e.EquipmentModifierContainer = nil
	}
	data, err := json.Marshal(&e.EquipmentModifierData)
	return data, err
}

// Container returns true if this is a container.
func (e *EquipmentModifier) Container() bool {
	return strings.HasSuffix(e.Type, commonContainerKeyPostfix)
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
	if e.Container() || (e.CostType == equipment.OriginalCost && e.CostAmount == "+0") {
		return ""
	}
	return e.CostType.Format(e.CostAmount) + " " + e.CostType.ShortString()
}

// WeightDescription returns the formatted weight.
func (e *EquipmentModifier) WeightDescription(entity *Entity) string {
	if e.Container() || (e.WeightType == equipment.OriginalWeight && (e.WeightAmount == "+0" || strings.HasPrefix(e.WeightAmount, "+0 "))) {
		return ""
	}
	return e.WeightType.Format(e.WeightAmount, SheetSettingsFor(entity).DefaultWeightUnits) + " " + e.WeightType.ShortString()
}

// FillWithNameableKeys adds any nameable keys found in this EquipmentModifier to the provided map.
func (e *EquipmentModifier) FillWithNameableKeys(m map[string]string) {
	if !e.Disabled {
		nameables.Extract(e.Name, m)
		nameables.Extract(e.Notes, m)
		nameables.Extract(e.VTTNotes, m)
		for _, one := range e.Features {
			one.FillWithNameableKeys(m)
		}
	}
}

// ApplyNameableKeys replaces any nameable keys found in this EquipmentModifier with the corresponding values in the provided map.
func (e *EquipmentModifier) ApplyNameableKeys(m map[string]string) {
	if !e.Disabled {
		e.Name = nameables.Apply(e.Name, m)
		e.Notes = nameables.Apply(e.Notes, m)
		e.VTTNotes = nameables.Apply(e.VTTNotes, m)
		for _, one := range e.Features {
			one.ApplyNameableKeys(m)
		}
	}
}

// ValueAdjustedForModifiers returns the value after adjusting it for a set of modifiers.
func ValueAdjustedForModifiers(value fixed.F64d4, modifiers []*EquipmentModifier) fixed.F64d4 {
	// Apply all equipment.OriginalCost
	cost := processNonCFStep(equipment.OriginalCost, value, modifiers)

	// Apply all equipment.BaseCost
	var cf fixed.F64d4
	for _, one := range modifiers {
		if !one.Disabled && one.CostType == equipment.BaseCost {
			t := equipment.BaseCost.DetermineModifierCostValueTypeFromString(one.CostAmount)
			cf += t.ExtractValue(one.CostAmount)
			if t == equipment.Multiplier {
				cf -= f64d4.One
			}
		}
	}
	if cf != 0 {
		cf = cf.Max(f64d4.NegPointEight)
		cost = cost.Mul(cf.Max(f64d4.NegPointEight) + f64d4.One)
	}

	// Apply all equipment.FinalBaseCost
	cost = processNonCFStep(equipment.FinalBaseCost, cost, modifiers)

	// Apply all equipment.FinalCost
	cost = processNonCFStep(equipment.FinalCost, cost, modifiers)

	return cost.Max(0)
}

func processNonCFStep(costType equipment.ModifierCostType, value fixed.F64d4, modifiers []*EquipmentModifier) fixed.F64d4 {
	var percentages, additions fixed.F64d4
	cost := value
	for _, one := range modifiers {
		if !one.Disabled && one.CostType == costType {
			t := costType.DetermineModifierCostValueTypeFromString(one.CostAmount)
			amt := t.ExtractValue(one.CostAmount)
			switch t {
			case equipment.Addition:
				additions += amt
			case equipment.Percentage:
				percentages += amt
			case equipment.Multiplier:
				cost = cost.Mul(amt)
			}
		}
	}
	cost += additions
	if percentages != 0 {
		cost += value.Mul(percentages.Div(f64d4.Hundred))
	}
	return cost
}

// WeightAdjustedForModifiers returns the weight after adjusting it for a set of modifiers.
func WeightAdjustedForModifiers(weight measure.Weight, modifiers []*EquipmentModifier, defUnits measure.WeightUnits) measure.Weight {
	var percentages fixed.F64d4
	w := fixed.F64d4(weight)

	// Apply all equipment.OriginalWeight
	for _, one := range modifiers {
		if !one.Disabled && one.WeightType == equipment.OriginalWeight {
			t := equipment.OriginalWeight.DetermineModifierWeightValueTypeFromString(one.WeightAmount)
			amt := t.ExtractFraction(one.WeightAmount).Value()
			if t == equipment.WeightAddition {
				w += measure.TrailingWeightUnitsFromString(one.WeightAmount, defUnits).ToPounds(amt)
			} else {
				percentages += amt
			}
		}
	}
	if percentages != 0 {
		w += fixed.F64d4(weight).Mul(percentages.Div(f64d4.Hundred))
	}

	// Apply all equipment.BaseWeight
	w = processMultiplyAddWeightStep(equipment.BaseWeight, w, defUnits, modifiers)

	// Apply all equipment.FinalBaseWeight
	w = processMultiplyAddWeightStep(equipment.FinalBaseWeight, w, defUnits, modifiers)

	// Apply all equipment.FinalWeight
	w = processMultiplyAddWeightStep(equipment.FinalWeight, w, defUnits, modifiers)

	return measure.Weight(w.Max(0))
}

func processMultiplyAddWeightStep(weightType equipment.ModifierWeightType, weight fixed.F64d4, defUnits measure.WeightUnits, modifiers []*EquipmentModifier) fixed.F64d4 {
	var sum fixed.F64d4
	for _, one := range modifiers {
		if !one.Disabled && one.WeightType == weightType {
			t := weightType.DetermineModifierWeightValueTypeFromString(one.WeightAmount)
			f := t.ExtractFraction(one.WeightAmount)
			switch t {
			case equipment.WeightAddition:
				sum += measure.TrailingWeightUnitsFromString(one.WeightAmount, defUnits).ToPounds(f.Value())
			case equipment.WeightPercentageMultiplier:
				weight = weight.Mul(f.Numerator).Div(f.Denominator.Mul(f64d4.Hundred))
			case equipment.WeightMultiplier:
				weight = weight.Mul(f.Numerator).Div(f.Denominator)
			}
		}
	}
	return weight + sum
}
