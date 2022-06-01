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
	"context"
	"io/fs"
	"strings"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/equipment"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/txt"
	"golang.org/x/exp/slices"
)

var _ Node = &EquipmentModifier{}

// Columns that can be used with the equipment modifier method .CellData()
const (
	EquipmentModifierDescriptionColumn = iota
	EquipmentModifierTechLevelColumn
	EquipmentModifierCostColumn
	EquipmentModifierWeightColumn
	EquipmentModifierTagsColumn
	EquipmentModifierReferenceColumn
)

const (
	equipmentModifierListTypeKey = "eqp_modifier_list"
	equipmentModifierTypeKey     = "eqp_modifier"
)

// EquipmentModifier holds a modifier to a piece of Equipment.
type EquipmentModifier struct {
	EquipmentModifierData
	Entity *Entity
}

type equipmentModifierListData struct {
	Type    string               `json:"type"`
	Version int                  `json:"version"`
	Rows    []*EquipmentModifier `json:"rows"`
}

// NewEquipmentModifiersFromFile loads an EquipmentModifier list from a file.
func NewEquipmentModifiersFromFile(fileSystem fs.FS, filePath string) ([]*EquipmentModifier, error) {
	var data equipmentModifierListData
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause(gid.InvalidFileDataMsg, err)
	}
	if data.Type != equipmentModifierListTypeKey {
		return nil, errs.New(gid.UnexpectedFileDataMsg)
	}
	if err := gid.CheckVersion(data.Version); err != nil {
		return nil, err
	}
	return data.Rows, nil
}

// SaveEquipmentModifiers writes the EquipmentModifier list to the file as JSON.
func SaveEquipmentModifiers(modifiers []*EquipmentModifier, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &equipmentModifierListData{
		Type:    equipmentModifierListTypeKey,
		Version: gid.CurrentDataVersion,
		Rows:    modifiers,
	})
}

// NewEquipmentModifier creates an EquipmentModifier.
func NewEquipmentModifier(entity *Entity, _ *EquipmentModifier, container bool) *EquipmentModifier {
	a := &EquipmentModifier{
		EquipmentModifierData: EquipmentModifierData{
			ContainerBase: newContainerBase[*EquipmentModifier](equipmentModifierTypeKey, container),
		},
		Entity: entity,
	}
	a.Name = a.Kind()
	return a
}

// Clone creates a copy of this data.
func (e *EquipmentModifier) Clone() *EquipmentModifier {
	other := *e
	other.Tags = txt.CloneStringSlice(e.Tags)
	other.Features = e.Features.Clone()
	other.Children = nil
	if len(e.Children) != 0 {
		other.Children = make([]*EquipmentModifier, 0, len(e.Children))
		for _, one := range e.Children {
			other.Children = append(other.Children, one.Clone())
		}
	}
	return &other
}

// MarshalJSON implements json.Marshaler.
func (e *EquipmentModifier) MarshalJSON() ([]byte, error) {
	e.ClearUnusedFieldsForType()
	return json.Marshal(&e.EquipmentModifierData)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *EquipmentModifier) UnmarshalJSON(data []byte) error {
	var localData struct {
		EquipmentModifierData
		// Old data fields
		Categories []string `json:"categories"`
	}
	if err := json.Unmarshal(data, &localData); err != nil {
		return err
	}
	localData.ClearUnusedFieldsForType()
	e.EquipmentModifierData = localData.EquipmentModifierData
	e.Tags = convertOldCategoriesToTags(e.Tags, localData.Categories)
	slices.Sort(e.Tags)
	return nil
}

// CellData returns the cell data information for the given column.
func (e *EquipmentModifier) CellData(column int, data *CellData) {
	switch column {
	case EquipmentModifierDescriptionColumn:
		data.Type = Text
		data.Primary = e.Name
		data.Secondary = e.SecondaryText()
	case EquipmentModifierTechLevelColumn:
		if !e.Container() {
			data.Type = Text
			data.Primary = e.TechLevel
		}
	case EquipmentModifierCostColumn:
		if !e.Container() {
			data.Type = Text
			data.Primary = e.CostDescription()
		}
	case EquipmentModifierWeightColumn:
		if !e.Container() {
			data.Type = Text
			data.Primary = e.WeightDescription()
		}
	case EquipmentModifierTagsColumn:
		data.Type = Text
		data.Primary = CombineTags(e.Tags)
	case EquipmentModifierReferenceColumn, PageRefCellAlias:
		data.Type = PageRef
		data.Primary = e.PageRef
		data.Secondary = e.Name
	}
}

// OwningEntity returns the owning Entity.
func (e *EquipmentModifier) OwningEntity() *Entity {
	return e.Entity
}

// SetOwningEntity sets the owning entity and configures any sub-components as needed.
func (e *EquipmentModifier) SetOwningEntity(entity *Entity) {
	e.Entity = entity
	if e.Container() {
		for _, child := range e.Children {
			child.SetOwningEntity(entity)
		}
	}
}

func (e *EquipmentModifier) String() string {
	return e.Name
}

// SecondaryText returns the "secondary" text: the text display below an Advantage.
func (e *EquipmentModifier) SecondaryText() string {
	var buffer strings.Builder
	settings := SheetSettingsFor(e.Entity)
	if e.LocalNotes != "" && settings.NotesDisplay.Inline() {
		if buffer.Len() != 0 {
			buffer.WriteByte('\n')
		}
		buffer.WriteString(e.LocalNotes)
	}
	return buffer.String()
}

// FullDescription returns a full description.
func (e *EquipmentModifier) FullDescription() string {
	var buffer strings.Builder
	buffer.WriteString(e.String())
	if e.LocalNotes != "" {
		buffer.WriteString(" (")
		buffer.WriteString(e.LocalNotes)
		buffer.WriteByte(')')
	}
	if SheetSettingsFor(e.Entity).ShowEquipmentModifierAdj {
		costDesc := e.CostDescription()
		weightDesc := e.WeightDescription()
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
	if e.Container() || (e.CostType == equipment.OriginalCost && (e.CostAmount == "" || e.CostAmount == "+0")) {
		return ""
	}
	return e.CostType.Format(e.CostAmount) + " " + e.CostType.String()
}

// WeightDescription returns the formatted weight.
func (e *EquipmentModifier) WeightDescription() string {
	if e.Container() || (e.WeightType == equipment.OriginalWeight && (e.WeightAmount == "" || strings.HasPrefix(e.WeightAmount, "+0 "))) {
		return ""
	}
	return e.WeightType.Format(e.WeightAmount, SheetSettingsFor(e.Entity).DefaultWeightUnits) + " " + e.WeightType.String()
}

// FillWithNameableKeys adds any nameable keys found in this EquipmentModifier to the provided map.
func (e *EquipmentModifier) FillWithNameableKeys(m map[string]string) {
	if !e.Disabled {
		nameables.Extract(e.Name, m)
		nameables.Extract(e.LocalNotes, m)
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
		e.LocalNotes = nameables.Apply(e.LocalNotes, m)
		e.VTTNotes = nameables.Apply(e.VTTNotes, m)
		for _, one := range e.Features {
			one.ApplyNameableKeys(m)
		}
	}
}

// ValueAdjustedForModifiers returns the value after adjusting it for a set of modifiers.
func ValueAdjustedForModifiers(value fxp.Int, modifiers []*EquipmentModifier) fxp.Int {
	// Apply all equipment.OriginalCost
	cost := processNonCFStep(equipment.OriginalCost, value, modifiers)

	// Apply all equipment.BaseCost
	var cf fxp.Int
	for _, one := range modifiers {
		if !one.Disabled && one.CostType == equipment.BaseCost {
			t := equipment.BaseCost.DetermineModifierCostValueTypeFromString(one.CostAmount)
			cf += t.ExtractValue(one.CostAmount)
			if t == equipment.Multiplier {
				cf -= fxp.One
			}
		}
	}
	if cf != 0 {
		cf = cf.Max(fxp.NegPointEight)
		cost = cost.Mul(cf.Max(fxp.NegPointEight) + fxp.One)
	}

	// Apply all equipment.FinalBaseCost
	cost = processNonCFStep(equipment.FinalBaseCost, cost, modifiers)

	// Apply all equipment.FinalCost
	cost = processNonCFStep(equipment.FinalCost, cost, modifiers)

	return cost.Max(0)
}

func processNonCFStep(costType equipment.ModifierCostType, value fxp.Int, modifiers []*EquipmentModifier) fxp.Int {
	var percentages, additions fxp.Int
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
		cost += value.Mul(percentages.Div(fxp.Hundred))
	}
	return cost
}

// WeightAdjustedForModifiers returns the weight after adjusting it for a set of modifiers.
func WeightAdjustedForModifiers(weight measure.Weight, modifiers []*EquipmentModifier, defUnits measure.WeightUnits) measure.Weight {
	var percentages fxp.Int
	w := fxp.Int(weight)

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
		w += fxp.Int(weight).Mul(percentages.Div(fxp.Hundred))
	}

	// Apply all equipment.BaseWeight
	w = processMultiplyAddWeightStep(equipment.BaseWeight, w, defUnits, modifiers)

	// Apply all equipment.FinalBaseWeight
	w = processMultiplyAddWeightStep(equipment.FinalBaseWeight, w, defUnits, modifiers)

	// Apply all equipment.FinalWeight
	w = processMultiplyAddWeightStep(equipment.FinalWeight, w, defUnits, modifiers)

	return measure.Weight(w.Max(0))
}

func processMultiplyAddWeightStep(weightType equipment.ModifierWeightType, weight fxp.Int, defUnits measure.WeightUnits, modifiers []*EquipmentModifier) fxp.Int {
	var sum fxp.Int
	for _, one := range modifiers {
		if !one.Disabled && one.WeightType == weightType {
			t := weightType.DetermineModifierWeightValueTypeFromString(one.WeightAmount)
			f := t.ExtractFraction(one.WeightAmount)
			switch t {
			case equipment.WeightAddition:
				sum += measure.TrailingWeightUnitsFromString(one.WeightAmount, defUnits).ToPounds(f.Value())
			case equipment.WeightPercentageMultiplier:
				weight = weight.Mul(f.Numerator).Div(f.Denominator.Mul(fxp.Hundred))
			case equipment.WeightMultiplier:
				weight = weight.Mul(f.Numerator).Div(f.Denominator)
			}
		}
	}
	return weight + sum
}
