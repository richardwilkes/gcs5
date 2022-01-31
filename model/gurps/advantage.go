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
	"github.com/richardwilkes/gcs/model/gurps/advantage"
	"github.com/richardwilkes/gcs/model/gurps/ancestry"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	advantageAncestryKey       = "ancestry"
	advantageBasePointsKey     = "base_points"
	advantageContainerTypeKey  = "container_type"
	advantageCRAdjKey          = "cr_adj"
	advantageCRKey             = "cr"
	advantageLevelsKey         = "levels"
	advantagePointsPerLevelKey = "points_per_level"
	advantagePrereqsKey        = "prereqs"
	advantageRoundCostDownKey  = "round_down"
	advantageTypeKey           = "advantage"
	advantageUserDescKey       = "userdesc"
	advantageCalcPointsKey     = "points"
)

// Advantage holds an advantage, disadvantage, quirk, or perk.
type Advantage struct {
	Common
	Parent            *Advantage
	Levels            fixed.F64d4
	BasePoints        fixed.F64d4
	PointsPerLevel    fixed.F64d4
	Features          []*Feature
	Modifiers         []*AdvantageModifier
	Prereq            *Prereq
	Weapons           []*Weapon
	Ancestry          *ancestry.Ancestry
	Children          []*Advantage
	UnsatisfiedReason string
	UserDesc          string
	Categories        []string
	ContainerType     advantage.ContainerType // TODO: Consider merging Container & ContainerType
	CR                advantage.SelfControlRoll
	CRAdj             SelfControlRollAdj
	Mental            bool
	Physical          bool
	Social            bool
	Exotic            bool
	Supernatural      bool
	Satisfied         bool
	RoundCostDown     bool
	SelfEnabled       bool
}

// NewAdvantage creates a new Advantage.
func NewAdvantage(parent *Advantage, container bool) *Advantage {
	return &Advantage{
		Common: Common{
			ID:        id.NewUUID(),
			Name:      i18n.Text("Advantage"),
			Container: container,
			Open:      true,
		},
		Parent:      parent,
		Levels:      f64d4.NegOne,
		Prereq:      NewPrereq(prereq.List, nil),
		Physical:    true,
		Satisfied:   true,
		SelfEnabled: true,
	}
}

// NewAdvantageFromJSON creates a new Advantage from a JSON object.
func NewAdvantageFromJSON(parent *Advantage, data map[string]interface{}, entity *Entity) *Advantage {
	a := &Advantage{Parent: parent}
	a.Common.FromJSON(advantageTypeKey, data)
	if a.Container {
		a.ContainerType = advantage.ContainerTypeFromKey(encoding.String(data[advantageContainerTypeKey]))
		array := encoding.Array(data[commonChildrenKey])
		if len(array) != 0 {
			a.Children = make([]*Advantage, len(array))
			for i, one := range array {
				a.Children[i] = NewAdvantageFromJSON(a, encoding.Object(one), entity)
			}
		}
	} else {
		a.SelfEnabled = !encoding.Bool(data[commonDisabledKey])
		a.RoundCostDown = encoding.Bool(data[advantageRoundCostDownKey])
		a.TypeBits = advantage.TypeFromJSON(data)
		if v, exists := data[advantageLevelsKey]; exists {
			a.Levels = encoding.Number(v)
		} else {
			a.Levels = f64d4.NegOne
		}
		a.BasePoints = encoding.Number(data[advantageBasePointsKey])
		a.PointsPerLevel = encoding.Number(data[advantagePointsPerLevelKey])
		a.Prereq = NewPrereqFromJSON(encoding.Object(data[advantagePrereqsKey]), entity)
		a.Weapons = WeaponsListFromJSON(data)
		a.Features = FeaturesListFromJSON(data)
	}
	a.CR = advantage.SelfControlRollFromJSON(advantageCRKey, data)
	a.CRAdj = SelfControlRollAdjFromKey(encoding.String(data[advantageCRAdjKey]))
	a.Modifiers = AdvantageModifiersListFromJSON(commonModifiersKey, data)
	a.UserDesc = encoding.String(data[advantageUserDescKey])
	a.Categories = StringListFromJSON(commonCategoriesKey, true, data)
	return a
}

// ToJSON emits this object as JSON.
func (a *Advantage) ToJSON(encoder *encoding.JSONEncoder, entity *Entity) {
	encoder.StartObject()
	a.Common.ToInlineJSON(advantageTypeKey, encoder)
	if a.Container {
		if a.ContainerType != advantage.Group {
			encoder.KeyedString(advantageContainerTypeKey, a.ContainerType.Key(), false, false)
			if a.ContainerType == advantage.Race {
				encoder.KeyedString(advantageAncestryKey, a.Ancestry.Name, true, true)
			}
		}
		if len(a.Children) != 0 {
			encoder.Key(commonChildrenKey)
			encoder.StartArray()
			for _, one := range a.Children {
				one.ToJSON(encoder, entity)
			}
			encoder.EndArray()
		}
	} else {
		encoder.KeyedBool(commonDisabledKey, !a.SelfEnabled, true)
		encoder.KeyedBool(advantageRoundCostDownKey, a.RoundCostDown, true)
		a.TypeBits.ToInlineJSON(encoder)
		if a.Levels >= 0 {
			encoder.KeyedNumber(advantageLevelsKey, a.Levels, false)
		}
		encoder.KeyedNumber(advantageBasePointsKey, a.BasePoints, true)
		encoder.KeyedNumber(advantagePointsPerLevelKey, a.PointsPerLevel, true)
		encoding.ToKeyedJSON(a.Prereq, advantagePrereqsKey, encoder)
		WeaponsListToJSON(a.Weapons, encoder)
		FeaturesListToJSON(a.Features, encoder)
	}
	a.CR.ToKeyedJSON(advantageCRKey, encoder)
	if a.CR != advantage.None {
		a.CRAdj.ToKeyedJSON(advantageCRAdjKey, encoder)
	}
	AdvantageModifiersListToJSON(commonModifiersKey, a.Modifiers, encoder)
	if entity != nil && entity.PreservesUserDesc() {
		encoder.KeyedString(advantageUserDescKey, a.UserDesc, true, true)
	}
	StringListToJSON(commonCategoriesKey, a.Categories, encoder)
	// Emit the calculated values for third parties
	encoder.Key(commonCalcKey)
	encoder.StartObject()
	encoder.KeyedNumber(advantageCalcPointsKey, a.AdjustedPoints(entity), false)
	encoder.EndObject()
	encoder.EndObject()
}

// AdjustedPoints returns the total points, taking levels and modifiers into account. 'entity' may be nil.
func (a *Advantage) AdjustedPoints(entity *Entity) fixed.F64d4 {
	if !a.SelfEnabled {
		return 0
	}
	if !a.Container {
		return AdjustedPoints(entity, a.BasePoints, a.Levels, a.PointsPerLevel, a.CR, a.AllModifiers(), a.RoundCostDown)
	}
	var points fixed.F64d4
	if a.ContainerType == advantage.AlternativeAbilities {
		values := make([]fixed.F64d4, len(a.Children))
		for i, one := range a.Children {
			values[i] = one.AdjustedPoints(entity)
			if values[i] > points {
				points = values[i]
			}
		}
		max := points
		found := false
		for _, v := range values {
			if !found && max == v {
				found = true
			} else {
				points += f64d4.ApplyRounding(calculateModifierPoints(v, f64d4.Twenty), a.RoundCostDown)
			}
		}
	} else {
		for _, one := range a.Children {
			points += one.AdjustedPoints(entity)
		}
	}
	return points
}

// AllModifiers returns the modifiers plus any inherited from parents.
func (a *Advantage) AllModifiers() []*AdvantageModifier {
	all := make([]*AdvantageModifier, len(a.Modifiers))
	copy(all, a.Modifiers)
	p := a.Parent
	for p != nil {
		all = append(all, p.Modifiers...)
		p = p.Parent
	}
	return all
}

// Enabled returns true if this Advantage and all of its parents are enabled.
func (a *Advantage) Enabled() bool {
	if !a.SelfEnabled {
		return false
	}
	p := a.Parent
	for p != nil {
		if !p.SelfEnabled {
			return false
		}
		p = p.Parent
	}
	return true
}

// TypeAsText returns the set of type bits that are set if this isn't a container, or an empty string if it is.
func (a *Advantage) TypeAsText() string {
	if a.Container {
		return ""
	}
	list := make([]string, 0, 5)
	if a.Mental {
		list = append(list, i18n.Text("Mental"))
	}
	if a.Physical {
		list = append(list, i18n.Text("Physical"))
	}
	if a.Social {
		list = append(list, i18n.Text("Social"))
	}
	if a.Exotic {
		list = append(list, i18n.Text("Exotic"))
	}
	if a.Supernatural {
		list = append(list, i18n.Text("Supernatural"))
	}
	return strings.Join(list, ", ")
}

func (a *Advantage) String() string {
	var buffer strings.Builder
	buffer.WriteString(a.Name)
	if !a.Container && a.Levels > 0 {
		buffer.WriteByte(' ')
		buffer.WriteString(a.Levels.String())
	}
	return buffer.String()
}

// FillWithNameableKeys adds any nameable keys found in this Advantage to the provided map.
func (a *Advantage) FillWithNameableKeys(nameables map[string]string) {
	a.Common.FillWithNameableKeys(nameables)
	a.Prereq.FillWithNameableKeys(nameables)
	for _, one := range a.Features {
		one.FillWithNameableKeys(nameables)
	}
	for _, one := range a.Weapons {
		one.FillWithNameableKeys(nameables)
	}
	for _, one := range a.Modifiers {
		one.FillWithNameableKeys(nameables)
	}
}

// ApplyNameableKeys replaces any nameable keys found in this Advantage with the corresponding values in the provided map.
func (a *Advantage) ApplyNameableKeys(nameables map[string]string) {
	a.Common.ApplyNameableKeys(nameables)
	a.Prereq.ApplyNameableKeys(nameables)
	for _, one := range a.Features {
		one.ApplyNameableKeys(nameables)
	}
	for _, one := range a.Weapons {
		one.ApplyNameableKeys(nameables)
	}
	for _, one := range a.Modifiers {
		one.ApplyNameableKeys(nameables)
	}
}

// ActiveModifierFor returns the first modifier that matches the name (case-insensitive).
func (a *Advantage) ActiveModifierFor(name string) *AdvantageModifier {
	for _, one := range a.Modifiers {
		if one.Enabled && strings.EqualFold(one.Name, name) {
			return one
		}
	}
	return nil
}

// ModifierNotes returns the notes due to modifiers. 'entity' may be nil.
func (a *Advantage) ModifierNotes(entity *Entity) string {
	var buffer strings.Builder
	if a.CR != advantage.None {
		buffer.WriteString(a.CR.String())
		if a.CRAdj != NoCRAdj {
			buffer.WriteString(", ")
			buffer.WriteString(a.CRAdj.Description(a.CR))
		}
	}
	for _, one := range a.Modifiers {
		if one.Enabled {
			if buffer.Len() != 0 {
				buffer.WriteString("; ")
			}
			buffer.WriteString(one.FullDescription(entity))
		}
	}
	return buffer.String()
}

// SecondaryText returns the "secondary" text: the text display below an Advantage.
func (a *Advantage) SecondaryText(entity *Entity) string {
	var buffer strings.Builder
	settings := SheetSettingsFor(entity)
	if a.UserDesc != "" && settings.UserDescriptionDisplay.Inline() {
		buffer.WriteString(a.UserDesc)
	}
	if settings.ModifiersDisplay.Inline() {
		if notes := a.ModifierNotes(entity); notes != "" {
			if buffer.Len() != 0 {
				buffer.WriteByte('\n')
			}
			buffer.WriteString(notes)
		}
	}
	if a.Notes != "" && settings.NotesDisplay.Inline() {
		if buffer.Len() != 0 {
			buffer.WriteByte('\n')
		}
		buffer.WriteString(a.Notes)
	}
	return buffer.String()
}

// HasCategory returns true if 'category' is present in 'categories'. This check both ignores case and can check for
// subsets that are colon-separated.
func HasCategory(category string, categories []string) bool {
	category = strings.TrimSpace(category)
	for _, one := range categories {
		for _, part := range strings.Split(one, ":") {
			if strings.EqualFold(category, strings.TrimSpace(part)) {
				return true
			}
		}
	}
	return false
}

// AdjustedPoints returns the total points, taking levels and modifiers into account. 'entity' may be nil.
func AdjustedPoints(entity *Entity, basePoints, levels, pointsPerLevel fixed.F64d4, cr advantage.SelfControlRoll, modifiers []*AdvantageModifier, roundCostDown bool) fixed.F64d4 {
	var baseEnh, levelEnh, baseLim, levelLim fixed.F64d4
	multiplier := cr.Multiplier()
	for _, one := range modifiers {
		if one.Enabled {
			modifier := one.CostModifier()
			switch one.CostType {
			case advantage.Percentage:
				switch one.Affects {
				case advantage.Total:
					if modifier < 0 {
						baseLim += modifier
						levelLim += modifier
					} else {
						baseEnh += modifier
						levelEnh += modifier
					}
				case advantage.BaseOnly:
					if modifier < 0 {
						baseLim += modifier
					} else {
						baseEnh += modifier
					}
				case advantage.LevelsOnly:
					if modifier < 0 {
						levelLim += modifier
					} else {
						levelEnh += modifier
					}
				}
			case advantage.Points:
				if one.Affects == advantage.LevelsOnly {
					pointsPerLevel += modifier
				} else {
					basePoints += modifier
				}
			case advantage.Multiplier:
				multiplier = multiplier.Mul(modifier)
			}
		}
	}
	modifiedBasePoints := basePoints
	leveledPoints := pointsPerLevel.Mul(levels)
	if baseEnh != 0 || baseLim != 0 || levelEnh != 0 || levelLim != 0 {
		if SheetSettingsFor(entity).UseMultiplicativeModifiers {
			if baseEnh == levelEnh && baseLim == levelLim {
				modifiedBasePoints = modifyPoints(modifyPoints(modifiedBasePoints+leveledPoints, baseEnh), f64d4.NegEighty.Max(baseLim))
			} else {
				modifiedBasePoints = modifyPoints(modifyPoints(modifiedBasePoints, baseEnh), f64d4.NegEighty.Max(baseLim)) +
					modifyPoints(modifyPoints(leveledPoints, levelEnh), f64d4.NegEighty.Max(levelLim))
			}
		} else {
			baseMod := f64d4.NegEighty.Max(baseEnh + baseLim)
			levelMod := f64d4.NegEighty.Max(levelEnh + levelLim)
			if baseMod == levelMod {
				modifiedBasePoints = modifyPoints(modifiedBasePoints+leveledPoints, baseMod)
			} else {
				modifiedBasePoints = modifyPoints(modifiedBasePoints, baseMod) + modifyPoints(leveledPoints, levelMod)
			}
		}
	} else {
		modifiedBasePoints += leveledPoints
	}
	return f64d4.ApplyRounding(modifiedBasePoints.Mul(multiplier), roundCostDown)
}

func modifyPoints(points, modifier fixed.F64d4) fixed.F64d4 {
	return points + calculateModifierPoints(points, modifier)
}

func calculateModifierPoints(points, modifier fixed.F64d4) fixed.F64d4 {
	return points.Mul(modifier).Div(f64d4.Hundred)
}
