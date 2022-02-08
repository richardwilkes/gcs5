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

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

var _ WeaponOwner = &Equipment{}

const equipmentTypeKey = "equipment"

// EquipmentContainer holds the Equipment data that only exists in containers.
type EquipmentContainer struct {
	Children []*Equipment `json:"children,omitempty"`
	Open     bool         `json:"open,omitempty"`
}

// EquipmentData holds the Equipment data that is written to disk.
type EquipmentData struct {
	Type                   string               `json:"type"`
	ID                     uuid.UUID            `json:"id"`
	Name                   string               `json:"description,omitempty"`
	PageRef                string               `json:"reference,omitempty"`
	LocalNotes             string               `json:"notes,omitempty"`
	VTTNotes               string               `json:"vtt_notes,omitempty"`
	TechLevel              string               `json:"tech_level,omitempty"`
	LegalityClass          string               `json:"legality_class,omitempty"`
	Quantity               fixed.F64d4          `json:"quantity,omitempty"`
	Value                  fixed.F64d4          `json:"value,omitempty"`
	Weight                 measure.Weight       `json:"weight,omitempty"`
	MaxUses                int                  `json:"max_uses,omitempty"`
	Uses                   int                  `json:"uses,omitempty"`
	Weapons                []*Weapon            `json:"weapons,omitempty"`
	Modifiers              []*EquipmentModifier `json:"modifiers,omitempty"`
	Features               feature.Features     `json:"features,omitempty"`
	Prereq                 *PrereqList          `json:"prereqs,omitempty"`
	Categories             []string             `json:"categories,omitempty"`
	Equipped               bool                 `json:"equipped,omitempty"`
	WeightIgnoredForSkills bool                 `json:"ignore_weight_for_skills,omitempty"`
	*EquipmentContainer    `json:",omitempty"`
}

// Equipment holds a piece of equipment.
type Equipment struct {
	EquipmentData
	Entity            *Entity
	Parent            *Equipment
	UnsatisfiedReason string
	Satisfied         bool
}

type equipmentListData struct {
	Current []*Equipment `json:"equipment"`
}

// NewEquipmentFromFile loads an Equipment list from a file.
func NewEquipmentFromFile(fileSystem fs.FS, filePath string) ([]*Equipment, error) {
	var data struct {
		equipmentListData
		OldKey []*Equipment `json:"rows"`
	}
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause("invalid equipment file: "+filePath, err)
	}
	if len(data.Current) != 0 {
		return data.Current, nil
	}
	return data.OldKey, nil
}

// SaveEquipment writes the Equipment list to the file as JSON.
func SaveEquipment(equipment []*Equipment, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &equipmentListData{Current: equipment})
}

// NewEquipment creates a new Equipment.
func NewEquipment(entity *Entity, parent *Equipment, container bool) *Equipment {
	e := Equipment{
		EquipmentData: EquipmentData{
			Type:          equipmentTypeKey,
			ID:            id.NewUUID(),
			Name:          i18n.Text("Equipment"),
			LegalityClass: "4",
			Prereq:        NewPrereqList(),
			Quantity:      fxp.One,
			Equipped:      true,
		},
		Entity: entity,
		Parent: parent,
	}
	if container {
		e.Type += commonContainerKeyPostfix
		e.EquipmentContainer = &EquipmentContainer{Open: true}
	}
	return &e
}

// MarshalJSON implements json.Marshaler.
func (e *Equipment) MarshalJSON() ([]byte, error) {
	type calc struct {
		ExtendedValue           fixed.F64d4     `json:"extended_value"`
		ExtendedWeight          measure.Weight  `json:"extended_weight"`
		ExtendedWeightForSkills *measure.Weight `json:"extended_weight_for_skills,omitempty"`
	}
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
	if !e.Container() {
		data.EquipmentContainer = nil
	}
	return json.Marshal(&data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *Equipment) UnmarshalJSON(data []byte) error {
	e.EquipmentData = EquipmentData{}
	if err := json.Unmarshal(data, &e.EquipmentData); err != nil {
		return err
	}
	if e.Prereq == nil {
		e.Prereq = NewPrereqList()
	}
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

// Container returns true if this is a container.
func (e *Equipment) Container() bool {
	return strings.HasSuffix(e.Type, commonContainerKeyPostfix)
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

// CategoryList returns the list of categories.
func (e *Equipment) CategoryList() []string {
	return e.Categories
}

// AdjustedValue returns the value after adjustments for any modifiers. Does not include the value of children.
func (e *Equipment) AdjustedValue() fixed.F64d4 {
	return ValueAdjustedForModifiers(e.Value, e.Modifiers)
}

// ExtendedValue returns the extended value.
func (e *Equipment) ExtendedValue() fixed.F64d4 {
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
	if e.Quantity <= 0 {
		return 0
	}
	base := fixed.F64d4(e.AdjustedWeight(forSkills, defUnits)).Mul(e.Quantity)
	if e.Container() {
		var contained fixed.F64d4
		for _, one := range e.Children {
			contained += fixed.F64d4(one.AdjustedWeight(forSkills, defUnits))
		}
		var percentage, reduction fixed.F64d4
		for _, one := range e.Features {
			if cwr, ok := one.(*feature.ContainedWeightReduction); ok {
				if cwr.IsPercentageReduction() {
					percentage += cwr.PercentageReduction()
				} else {
					reduction += fixed.F64d4(cwr.FixedReduction(defUnits))
				}
			}
		}
		for _, one := range e.Modifiers {
			if !one.Disabled {
				for _, f := range e.Features {
					if cwr, ok := f.(*feature.ContainedWeightReduction); ok {
						if cwr.IsPercentageReduction() {
							percentage += cwr.PercentageReduction()
						} else {
							reduction += fixed.F64d4(cwr.FixedReduction(defUnits))
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
	e.Prereq.FillWithNameableKeys(m)
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
	e.Prereq.ApplyNameableKeys(m)
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
			buffer.WriteString(one.FullDescription(e.Entity))
		}
	}
	return buffer.String()
}
