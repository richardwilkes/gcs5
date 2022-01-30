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
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	equipmentTypeKey                        = "equipment"
	equipmentEquippedKey                    = "equipped"
	equipmentQuantityKey                    = "quantity"
	equipmentDescriptionKey                 = "description"
	equipmentTechLevelKey                   = "tech_level"
	equipmentLegalityClassKey               = "legality_class"
	equipmentValueKey                       = "value"
	equipmentIgnoreWeightForSkillsKey       = "ignore_weight_for_skills"
	equipmentWeightKey                      = "weight"
	equipmentUsesKey                        = "uses"
	equipmentMaxUsesKey                     = "max_uses"
	equipmentPrereqsKey                     = "prereqs"
	equipmentCalcExtendedValueKey           = "extended_value"
	equipmentCalcExtendedWeightKey          = "extended_weight"
	equipmentCalcExtendedWeightForSkillsKey = "extended_weight_for_skills"
)

// Equipment holds a piece of equipment.
type Equipment struct {
	Common
	Parent                 *Equipment
	Quantity               fixed.F64d4
	Value                  fixed.F64d4
	Weight                 measure.Weight
	Uses                   int
	MaxUses                int
	TechLevel              string
	LegalityClass          string
	UnsatisfiedReason      string
	Weapons                []*Weapon
	Modifiers              []*EquipmentModifier
	Features               []*Feature
	Prereq                 *Prereq
	Categories             []string
	Children               []*Equipment
	Equipped               bool
	WeightIgnoredForSkills bool
	Satisfied              bool
}

// NewEquipment creates a new Equipment.
func NewEquipment(parent *Equipment, container bool) *Equipment {
	return &Equipment{
		Common: Common{
			ID:        id.NewUUID(),
			Name:      i18n.Text("Equipment"),
			Container: container,
			Open:      true,
		},
		Parent:        parent,
		Quantity:      f64d4.One,
		LegalityClass: "4",
		Prereq:        NewPrereq(prereq.List, nil),
		Equipped:      true,
		Satisfied:     true,
	}
}

// NewEquipmentFromJSON creates a new Equipment from a JSON object. 'entity' may be nil.
func NewEquipmentFromJSON(parent *Equipment, data map[string]interface{}, entity *Entity) *Equipment {
	e := &Equipment{Parent: parent}
	e.Common.FromJSON(equipmentTypeKey, data)
	if e.Name == "" { // If no name, then try to load from the old key
		e.Name = encoding.String(data[equipmentDescriptionKey])
	}
	e.Equipped = encoding.Bool(data[equipmentEquippedKey])
	e.TechLevel = encoding.String(data[equipmentTechLevelKey])
	e.LegalityClass = encoding.String(data[equipmentLegalityClassKey])
	e.Value = encoding.Number(data[equipmentValueKey])
	e.WeightIgnoredForSkills = encoding.Bool(data[equipmentIgnoreWeightForSkillsKey])
	defWeightUnits := SheetSettingsFor(entity).DefaultWeightUnits
	e.Weight = measure.WeightFromStringForced(encoding.String(data[equipmentWeightKey]), defWeightUnits)
	e.MaxUses = int(encoding.Number(data[equipmentMaxUsesKey]).Max(0).AsInt64())
	e.Uses = xmath.MinInt(int(encoding.Number(data[equipmentUsesKey]).Max(0).AsInt64()), e.MaxUses)
	e.Weapons = WeaponsListFromJSON(data)
	e.Modifiers = EquipmentModifiersListFromJSON(commonModifiersKey, data)
	e.Features = FeaturesListFromJSON(data)
	e.Prereq = NewPrereqFromJSON(encoding.Object(data[equipmentPrereqsKey]), entity)
	e.Categories = StringListFromJSON(commonCategoriesKey, true, data)
	if e.Container {
		array := encoding.Array(data[commonChildrenKey])
		if len(array) != 0 {
			e.Children = make([]*Equipment, len(array))
			for i, one := range array {
				e.Children[i] = NewEquipmentFromJSON(e, encoding.Object(one), entity)
			}
		}
		e.Quantity = f64d4.One
	} else {
		e.Quantity = encoding.Number(data[equipmentQuantityKey])
	}
	return e
}

// ToJSON emits this object as JSON.
func (e *Equipment) ToJSON(encoder *encoding.JSONEncoder, entity *Entity) {
	defUnits := SheetSettingsFor(entity).DefaultWeightUnits
	encoder.StartObject()
	e.Common.ToInlineJSON(equipmentTypeKey, encoder)
	if entity != nil {
		encoder.KeyedBool(equipmentEquippedKey, e.Equipped, true)
	}
	encoder.KeyedString(equipmentTechLevelKey, e.TechLevel, true, true)
	encoder.KeyedString(equipmentLegalityClassKey, e.LegalityClass, true, true)
	encoder.KeyedNumber(equipmentValueKey, e.Value, true)
	encoder.KeyedBool(equipmentIgnoreWeightForSkillsKey, e.WeightIgnoredForSkills, true)
	encoder.KeyedString(equipmentWeightKey, e.Weight.String(), false, false)
	encoder.KeyedNumber(equipmentMaxUsesKey, fixed.F64d4FromInt64(int64(e.MaxUses)), true)
	encoder.KeyedNumber(equipmentUsesKey, fixed.F64d4FromInt64(int64(e.Uses)), true)
	WeaponsListToJSON(e.Weapons, encoder)
	EquipmentModifiersListToJSON(commonModifiersKey, e.Modifiers, encoder)
	FeaturesListToJSON(e.Features, encoder)
	encoding.ToKeyedJSON(e.Prereq, equipmentPrereqsKey, encoder)
	StringListToJSON(commonCategoriesKey, e.Categories, encoder)
	if e.Container {
		encoder.Key(commonChildrenKey)
		encoder.StartArray()
		for _, one := range e.Children {
			one.ToJSON(encoder, entity)
		}
		encoder.EndArray()
	} else {
		encoder.KeyedNumber(equipmentQuantityKey, e.Quantity, true)
	}
	// Emit the calculated values for third parties
	encoder.Key(commonCalcKey)
	encoder.StartObject()
	encoder.KeyedNumber(equipmentCalcExtendedValueKey, e.ExtendedValue(), true)
	encoder.KeyedString(equipmentCalcExtendedWeightKey, e.ExtendedWeight(false, defUnits).String(), false, false)
	if e.WeightIgnoredForSkills {
		encoder.KeyedString(equipmentCalcExtendedWeightForSkillsKey, e.ExtendedWeight(true, defUnits).String(), false, false)
	}
	encoder.EndObject()
	encoder.EndObject()
}

// AdjustedValue returns the value after adjustments for any modifiers. Does not include the value of children.
func (e *Equipment) AdjustedValue() fixed.F64d4 {
	return ValueAdjustedForModifiers(e.Value, e.Modifiers)
}

// ExtendedValue returns the extended value.
func (e *Equipment) ExtendedValue() fixed.F64d4 {
	value := e.Quantity.Mul(e.AdjustedValue())
	for _, one := range e.Children {
		value += one.ExtendedValue()
	}
	return value
}

// AdjustedWeight returns the weight after adjustments for any modifiers. Does not include the weight of children.
// 'entity' may be nil.
func (e *Equipment) AdjustedWeight(forSkills bool, defUnits measure.WeightUnits) measure.Weight {
	if forSkills && e.WeightIgnoredForSkills {
		return 0
	}
	return WeightAdjustedForModifiers(e.Weight, e.Modifiers, defUnits)
}

// ExtendedWeight returns the extended weight.
func (e *Equipment) ExtendedWeight(forSkills bool, defUnits measure.WeightUnits) measure.Weight {
	var contained fixed.F64d4
	for _, one := range e.Children {
		contained += fixed.F64d4(one.AdjustedWeight(forSkills, defUnits))
	}
	var percentage, reduction fixed.F64d4
	for _, one := range e.Features {
		if one.Type == feature.ContainedWeightReduction {
			if one.IsPercentageReduction() {
				percentage += one.PercentageReduction()
			} else {
				reduction += fixed.F64d4(one.FixedReduction(defUnits))
			}
		}
	}
	for _, one := range e.Modifiers {
		if one.Enabled {
			for _, f := range e.Features {
				if f.Type == feature.ContainedWeightReduction {
					if f.IsPercentageReduction() {
						percentage += f.PercentageReduction()
					} else {
						reduction += fixed.F64d4(f.FixedReduction(defUnits))
					}
				}
			}
		}
	}
	if percentage >= f64d4.Hundred {
		contained = 0
	} else if percentage > 0 {
		contained -= contained.Mul(percentage).Div(f64d4.Hundred)
	}
	contained -= reduction
	return measure.Weight(fixed.F64d4(e.AdjustedWeight(forSkills, defUnits)).Mul(e.Quantity) + contained.Max(0))
}

// FillWithNameableKeys adds any nameable keys found in this Advantage to the provided map.
func (e *Equipment) FillWithNameableKeys(nameables map[string]string) {
	e.Common.FillWithNameableKeys(nameables)
	e.Prereq.FillWithNameableKeys(nameables)
	for _, one := range e.Features {
		one.FillWithNameableKeys(nameables)
	}
	for _, one := range e.Weapons {
		one.FillWithNameableKeys(nameables)
	}
	for _, one := range e.Modifiers {
		one.FillWithNameableKeys(nameables)
	}
}

// ApplyNameableKeys replaces any nameable keys found in this Advantage with the corresponding values in the provided map.
func (e *Equipment) ApplyNameableKeys(nameables map[string]string) {
	e.Common.ApplyNameableKeys(nameables)
	e.Prereq.ApplyNameableKeys(nameables)
	for _, one := range e.Features {
		one.ApplyNameableKeys(nameables)
	}
	for _, one := range e.Weapons {
		one.ApplyNameableKeys(nameables)
	}
	for _, one := range e.Modifiers {
		one.ApplyNameableKeys(nameables)
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
		if one.Enabled && strings.EqualFold(one.Name, name) {
			return one
		}
	}
	return nil
}

// ModifierNotes returns the notes due to modifiers. 'entity' may be nil.
func (e *Equipment) ModifierNotes(entity *Entity) string {
	var buffer strings.Builder
	for _, one := range e.Modifiers {
		if one.Enabled {
			if buffer.Len() != 0 {
				buffer.WriteString("; ")
			}
			buffer.WriteString(one.FullDescription(entity))
		}
	}
	return buffer.String()
}
