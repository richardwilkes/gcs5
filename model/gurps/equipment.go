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
	"strconv"
	"strings"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
	"golang.org/x/exp/slices"
)

var (
	_ WeaponOwner = &Equipment{}
	_ node.Node   = &Equipment{}
)

// Columns that can be used with the equipment method .CellData()
const (
	EquipmentEquippedColumn = iota
	EquipmentQuantityColumn
	EquipmentDescriptionColumn
	EquipmentUsesColumn
	EquipmentMaxUsesColumn
	EquipmentTLColumn
	EquipmentLCColumn
	EquipmentCostColumn
	EquipmentExtendedCostColumn
	EquipmentWeightColumn
	EquipmentExtendedWeightColumn
	EquipmentTagsColumn
	EquipmentReferenceColumn
)

var (
	// TechLevelInfo holds the general TL age list
	TechLevelInfo = i18n.Text(`TL0: Stone Age (Prehistory)
TL1: Bronze Age (3500 B.C.+)
TL2: Iron Age (1200 B.C.+)
TL3: Medieval (600 A.D.+)
TL4: Age of Sail (1450+)
TL5: Industrial Revolution (1730+)
TL6: Mechanized Age (1880+)
TL7: Nuclear Age (1940+)
TL8: Digital Age (1980+)
TL9: Microtech Age (2025+?)
TL10: Robotic Age (2070+?)
TL11: Age of Exotic Matter
TL12: Anything Goes`)

	// LegalityClassInfo holds the LC list
	LegalityClassInfo = i18n.Text(`LC0: Banned
LC1: Military
LC2: Restricted
LC3: Licensed
LC4: Open`)
)

const (
	equipmentListTypeKey = "equipment_list"
	equipmentTypeKey     = "equipment"
)

// Equipment holds a piece of equipment.
type Equipment struct {
	EquipmentData
	Entity            *Entity
	Parent            *Equipment
	UnsatisfiedReason string
	Satisfied         bool
}

type equipmentListData struct {
	Type    string       `json:"type"`
	Version int          `json:"version"`
	Rows    []*Equipment `json:"rows"`
}

// NewEquipmentFromFile loads an Equipment list from a file.
func NewEquipmentFromFile(fileSystem fs.FS, filePath string) ([]*Equipment, error) {
	var data equipmentListData
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause(gid.InvalidFileDataMsg, err)
	}
	if data.Type != equipmentListTypeKey {
		return nil, errs.New(gid.UnexpectedFileDataMsg)
	}
	if err := gid.CheckVersion(data.Version); err != nil {
		return nil, err
	}
	return data.Rows, nil
}

// SaveEquipment writes the Equipment list to the file as JSON.
func SaveEquipment(equipment []*Equipment, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &equipmentListData{
		Type:    equipmentListTypeKey,
		Version: gid.CurrentDataVersion,
		Rows:    equipment,
	})
}

// NewEquipment creates a new Equipment.
func NewEquipment(entity *Entity, parent *Equipment, container bool) *Equipment {
	e := Equipment{
		EquipmentData: EquipmentData{
			ContainerBase: newContainerBase[*Equipment](equipmentTypeKey, container),
			EquipmentEditData: EquipmentEditData{
				LegalityClass: "4",
				Quantity:      fxp.One,
				Equipped:      true,
			},
		},
		Entity: entity,
		Parent: parent,
	}
	e.Name = e.Kind()
	return &e
}

// MarshalJSON implements json.Marshaler.
func (e *Equipment) MarshalJSON() ([]byte, error) {
	type calc struct {
		ExtendedValue           fxp.Int         `json:"extended_value"`
		ExtendedWeight          measure.Weight  `json:"extended_weight"`
		ExtendedWeightForSkills *measure.Weight `json:"extended_weight_for_skills,omitempty"`
	}
	e.ClearUnusedFieldsForType()
	defUnits := SheetSettingsFor(e.Entity).DefaultWeightUnits
	data := struct {
		EquipmentData
		Calc calc `json:"calc"`
	}{
		EquipmentData: e.EquipmentData,
		Calc: calc{
			ExtendedValue:           e.ExtendedValue(),
			ExtendedWeight:          e.ExtendedWeight(false, defUnits),
			ExtendedWeightForSkills: nil,
		},
	}
	if e.WeightIgnoredForSkills {
		w := e.ExtendedWeight(true, defUnits)
		data.Calc.ExtendedWeightForSkills = &w
	}
	return json.Marshal(&data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *Equipment) UnmarshalJSON(data []byte) error {
	var localData struct {
		EquipmentData
		// Old data fields
		Categories []string `json:"categories"`
	}
	if err := json.Unmarshal(data, &localData); err != nil {
		return err
	}
	localData.ClearUnusedFieldsForType()
	e.EquipmentData = localData.EquipmentData
	e.Tags = convertOldCategoriesToTags(e.Tags, localData.Categories)
	slices.Sort(e.Tags)
	if e.Container() {
		if e.Quantity == 0 {
			// Old formats omitted the quantity for containers. Try to see if it was omitted or if it was explicitly
			// set to zero.
			m := make(map[string]interface{})
			if err := json.Unmarshal(data, &m); err == nil {
				if _, exists := m["quantity"]; !exists {
					e.Quantity = fxp.One
				}
			}
		}
		for _, one := range e.Children {
			one.Parent = e
		}
	}
	return nil
}

// CellData returns the cell data information for the given column.
func (e *Equipment) CellData(column int, data *node.CellData) {
	switch column {
	case EquipmentEquippedColumn:
		data.Type = node.Toggle
		data.Checked = e.Equipped
		data.Alignment = unison.MiddleAlignment
	case EquipmentQuantityColumn:
		data.Type = node.Text
		data.Primary = e.Quantity.String()
		data.Alignment = unison.EndAlignment
	case EquipmentDescriptionColumn:
		data.Type = node.Text
		data.Primary = e.Description()
		data.Secondary = e.SecondaryText()
	case EquipmentUsesColumn:
		if e.MaxUses > 0 {
			data.Type = node.Text
			data.Primary = strconv.Itoa(e.Uses)
			data.Alignment = unison.EndAlignment
		}
	case EquipmentMaxUsesColumn:
		if e.MaxUses > 0 {
			data.Type = node.Text
			data.Primary = strconv.Itoa(e.MaxUses)
			data.Alignment = unison.EndAlignment
		}
	case EquipmentTLColumn:
		data.Type = node.Text
		data.Primary = e.TechLevel
		data.Alignment = unison.EndAlignment
	case EquipmentLCColumn:
		data.Type = node.Text
		data.Primary = e.LegalityClass
		data.Alignment = unison.EndAlignment
	case EquipmentCostColumn:
		data.Type = node.Text
		data.Primary = e.AdjustedValue().String()
		data.Alignment = unison.EndAlignment
	case EquipmentExtendedCostColumn:
		data.Type = node.Text
		data.Primary = e.ExtendedValue().String()
		data.Alignment = unison.EndAlignment
	case EquipmentWeightColumn:
		data.Type = node.Text
		units := SheetSettingsFor(e.Entity).DefaultWeightUnits
		data.Primary = units.Format(e.AdjustedWeight(false, units))
		data.Alignment = unison.EndAlignment
	case EquipmentExtendedWeightColumn:
		data.Type = node.Text
		units := SheetSettingsFor(e.Entity).DefaultWeightUnits
		data.Primary = units.Format(e.ExtendedWeight(false, units))
		data.Alignment = unison.EndAlignment
	case EquipmentTagsColumn:
		data.Type = node.Text
		data.Primary = CombineTags(e.Tags)
	case EquipmentReferenceColumn:
		data.Type = node.PageRef
		data.Primary = e.PageRef
		data.Secondary = e.Name
	}
}

// Depth returns the number of parents this node has.
func (e *Equipment) Depth() int {
	count := 0
	p := e.Parent
	for p != nil {
		count++
		p = p.Parent
	}
	return count
}

// OwningEntity returns the owning Entity.
func (e *Equipment) OwningEntity() *Entity {
	return e.Entity
}

// SetOwningEntity sets the owning entity and configures any sub-components as needed.
func (e *Equipment) SetOwningEntity(entity *Entity) {
	e.Entity = entity
	for _, w := range e.Weapons {
		w.SetOwner(e)
	}
	if e.Container() {
		for _, child := range e.Children {
			child.SetOwningEntity(entity)
		}
	}
}

// Description returns a description.
func (e *Equipment) Description() string {
	return e.Name
}

// SecondaryText returns the "secondary" text: the text display below the description.
func (e *Equipment) SecondaryText() string {
	var buffer strings.Builder
	settings := SheetSettingsFor(e.Entity)
	if settings.ModifiersDisplay.Inline() {
		if notes := e.ModifierNotes(); notes != "" {
			if buffer.Len() != 0 {
				buffer.WriteByte('\n')
			}
			buffer.WriteString(notes)
		}
	}
	if e.LocalNotes != "" && settings.NotesDisplay.Inline() {
		if buffer.Len() != 0 {
			buffer.WriteByte('\n')
		}
		buffer.WriteString(e.LocalNotes)
	}
	return buffer.String()
}

// String implements fmt.Stringer.
func (e *Equipment) String() string {
	return e.Name
}

// Notes returns the local notes.
func (e *Equipment) Notes() string {
	return e.LocalNotes
}

// FeatureList returns the list of Features.
func (e *Equipment) FeatureList() feature.Features {
	return e.Features
}

// TagList returns the list of tags.
func (e *Equipment) TagList() []string {
	return e.Tags
}

// AdjustedValue returns the value after adjustments for any modifiers. Does not include the value of children.
func (e *Equipment) AdjustedValue() fxp.Int {
	return ValueAdjustedForModifiers(e.Value, e.Modifiers)
}

// ExtendedValue returns the extended value.
func (e *Equipment) ExtendedValue() fxp.Int {
	if e.Quantity <= 0 {
		return 0
	}
	value := e.AdjustedValue().Mul(e.Quantity)
	if e.Container() {
		for _, one := range e.Children {
			value += one.ExtendedValue()
		}
	}
	return value
}

// AdjustedWeight returns the weight after adjustments for any modifiers. Does not include the weight of children.
func (e *Equipment) AdjustedWeight(forSkills bool, defUnits measure.WeightUnits) measure.Weight {
	if forSkills && e.WeightIgnoredForSkills {
		return 0
	}
	return WeightAdjustedForModifiers(e.Weight, e.Modifiers, defUnits)
}

// ExtendedWeight returns the extended weight.
func (e *Equipment) ExtendedWeight(forSkills bool, defUnits measure.WeightUnits) measure.Weight {
	return ExtendedWeightAdjustedForModifiers(defUnits, e.Quantity, e.Weight, e.Modifiers, e.Features, e.Children, forSkills, e.WeightIgnoredForSkills)
}

// ExtendedWeightAdjustedForModifiers calculates the extended weight.
func ExtendedWeightAdjustedForModifiers(defUnits measure.WeightUnits, qty fxp.Int, baseWeight measure.Weight, modifiers []*EquipmentModifier, features feature.Features, children []*Equipment, forSkills, weightIgnoredForSkills bool) measure.Weight {
	if qty <= 0 {
		return 0
	}
	var base fxp.Int
	if !forSkills || !weightIgnoredForSkills {
		base = fxp.Int(WeightAdjustedForModifiers(baseWeight, modifiers, defUnits)).Mul(qty)
	}
	if len(children) != 0 {
		var contained fxp.Int
		for _, one := range children {
			contained += fxp.Int(one.ExtendedWeight(forSkills, defUnits))
		}
		var percentage, reduction fxp.Int
		for _, one := range features {
			if cwr, ok := one.(*feature.ContainedWeightReduction); ok {
				if cwr.IsPercentageReduction() {
					percentage += cwr.PercentageReduction()
				} else {
					reduction += fxp.Int(cwr.FixedReduction(defUnits))
				}
			}
		}
		for _, one := range modifiers {
			if !one.Disabled {
				for _, f := range one.Features {
					if cwr, ok := f.(*feature.ContainedWeightReduction); ok {
						if cwr.IsPercentageReduction() {
							percentage += cwr.PercentageReduction()
						} else {
							reduction += fxp.Int(cwr.FixedReduction(defUnits))
						}
					}
				}
			}
		}
		if percentage >= fxp.Hundred {
			contained = 0
		} else if percentage > 0 {
			contained -= contained.Mul(percentage).Div(fxp.Hundred)
		}
		base += (contained - reduction).Max(0)
	}
	return measure.Weight(base)
}

// FillWithNameableKeys adds any nameable keys found in this Advantage to the provided map.
func (e *Equipment) FillWithNameableKeys(m map[string]string) {
	nameables.Extract(e.Name, m)
	nameables.Extract(e.LocalNotes, m)
	nameables.Extract(e.VTTNotes, m)
	if e.Prereq != nil {
		e.Prereq.FillWithNameableKeys(m)
	}
	for _, one := range e.Features {
		one.FillWithNameableKeys(m)
	}
	for _, one := range e.Weapons {
		one.FillWithNameableKeys(m)
	}
	for _, one := range e.Modifiers {
		one.FillWithNameableKeys(m)
	}
}

// ApplyNameableKeys replaces any nameable keys found in this Advantage with the corresponding values in the provided map.
func (e *Equipment) ApplyNameableKeys(m map[string]string) {
	e.Name = nameables.Apply(e.Name, m)
	e.LocalNotes = nameables.Apply(e.LocalNotes, m)
	e.VTTNotes = nameables.Apply(e.VTTNotes, m)
	if e.Prereq != nil {
		e.Prereq.ApplyNameableKeys(m)
	}
	for _, one := range e.Features {
		one.ApplyNameableKeys(m)
	}
	for _, one := range e.Weapons {
		one.ApplyNameableKeys(m)
	}
	for _, one := range e.Modifiers {
		one.ApplyNameableKeys(m)
	}
}

// DisplayLegalityClass returns a display version of the LegalityClass.
func (e *Equipment) DisplayLegalityClass() string {
	lc := strings.TrimSpace(e.LegalityClass)
	switch lc {
	case "0":
		return i18n.Text("LC0: Banned")
	case "1":
		return i18n.Text("LC1: Military")
	case "2":
		return i18n.Text("LC2: Restricted")
	case "3":
		return i18n.Text("LC3: Licensed")
	case "4":
		return i18n.Text("LC4: Open")
	default:
		return lc
	}
}

// ActiveModifierFor returns the first modifier that matches the name (case-insensitive).
func (e *Equipment) ActiveModifierFor(name string) *EquipmentModifier {
	for _, one := range e.Modifiers {
		if !one.Disabled && strings.EqualFold(one.Name, name) {
			return one
		}
	}
	return nil
}

// ModifierNotes returns the notes due to modifiers.
func (e *Equipment) ModifierNotes() string {
	var buffer strings.Builder
	for _, one := range e.Modifiers {
		if !one.Disabled {
			if buffer.Len() != 0 {
				buffer.WriteString("; ")
			}
			buffer.WriteString(one.FullDescription())
		}
	}
	return buffer.String()
}
