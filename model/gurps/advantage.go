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
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/advantage"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const advantageTypeKey = "advantage"

// AdvantageItem holds the Advantage data that only exists in non-containers.
type AdvantageItem struct {
	Levels         *fixed.F64d4 `json:"levels,omitempty"`
	BasePoints     fixed.F64d4  `json:"base_points"`
	PointsPerLevel fixed.F64d4  `json:"points_per_level,omitempty"`
	Prereq         *Prereq      `json:"prereqs,omitempty"`
	Weapons        []*Weapon    `json:"weapons,omitempty"`
	Features       []*Feature   `json:"features,omitempty"`
	Mental         bool         `json:"mental,omitempty"`
	Physical       bool         `json:"physical,omitempty"`
	Social         bool         `json:"social,omitempty"`
	Exotic         bool         `json:"exotic,omitempty"`
	Supernatural   bool         `json:"supernatural,omitempty"`
	Disabled       bool         `json:"disabled,omitempty"`
	RoundCostDown  bool         `json:"round_down,omitempty"`
}

// AdvantageContainer holds the Advantage data that only exists in containers.
type AdvantageContainer struct {
	ContainerType advantage.ContainerType `json:"container_type,omitempty"`
	Ancestry      string                  `json:"ancestry,omitempty"`
	Children      []*Advantage            `json:"children,omitempty"`
	Open          bool                    `json:"open,omitempty"`
}

// AdvantageCalc holds the Advantage data that is only emitted for third parties.
type AdvantageCalc struct {
	Points fixed.F64d4 `json:"points"`
}

// AdvantageData holds the Advantage data that is written to disk.
type AdvantageData struct {
	Type                string                    `json:"type"`
	ID                  uuid.UUID                 `json:"id"`
	Name                string                    `json:"name,omitempty"`
	PageRef             string                    `json:"reference,omitempty"`
	Notes               string                    `json:"notes,omitempty"`
	VTTNotes            string                    `json:"vtt_notes,omitempty"`
	CR                  advantage.SelfControlRoll `json:"cr,omitempty"`
	CRAdj               SelfControlRollAdj        `json:"cr_adj,omitempty"`
	Modifiers           []*AdvantageModifier      `json:"modifiers,omitempty"`
	UserDesc            string                    `json:"userdesc,omitempty"`
	Categories          []string                  `json:"categories,omitempty"`
	*AdvantageItem      `json:",omitempty"`
	*AdvantageContainer `json:",omitempty"`
	Calc                *AdvantageCalc `json:"calc,omitempty"`
}

// Advantage holds an advantage, disadvantage, quirk, or perk.
type Advantage struct {
	AdvantageData
	Entity            *Entity
	Parent            *Advantage
	UnsatisfiedReason string
	Satisfied         bool
}

// NewAdvantage creates a new Advantage.
func NewAdvantage(entity *Entity, parent *Advantage, container bool) *Advantage {
	a := Advantage{
		AdvantageData: AdvantageData{
			Type: advantageTypeKey,
			ID:   id.NewUUID(),
			Name: i18n.Text("Advantage"),
		},
		Entity: entity,
		Parent: parent,
	}
	if container {
		a.Type += commonContainerKeyPostfix
		a.AdvantageContainer = &AdvantageContainer{Open: true}
	} else {
		a.AdvantageItem = &AdvantageItem{
			Prereq:   NewPrereq(prereq.List, nil),
			Physical: true,
		}
	}
	return &a
}

// MarshalJSON implements json.Marshaler.
func (a *Advantage) MarshalJSON() ([]byte, error) {
	if a.Container() {
		a.AdvantageItem = nil
	} else {
		a.AdvantageContainer = nil
	}
	a.Calc = &AdvantageCalc{
		Points: a.AdjustedPoints(),
	}
	data, err := json.Marshal(&a.AdvantageData)
	a.Calc = nil
	return data, err
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Advantage) UnmarshalJSON(data []byte) error {
	a.AdvantageData = AdvantageData{}
	if err := json.Unmarshal(data, &a.AdvantageData); err != nil {
		return err
	}
	if a.Container() {
		for _, one := range a.Children {
			one.Parent = a
		}
	}
	return nil
}

// Container returns true if this is a container.
func (a *Advantage) Container() bool {
	return strings.HasSuffix(a.Type, commonContainerKeyPostfix)
}

// AdjustedPoints returns the total points, taking levels and modifiers into account. 'entity' may be nil.
func (a *Advantage) AdjustedPoints() fixed.F64d4 {
	if a.Disabled {
		return 0
	}
	if !a.Container() {
		var levels fixed.F64d4
		if a.Levels != nil {
			levels = *a.Levels
		}
		return AdjustedPoints(a.Entity, a.BasePoints, levels, a.PointsPerLevel, a.CR, a.AllModifiers(), a.RoundCostDown)
	}
	var points fixed.F64d4
	if a.ContainerType == advantage.AlternativeAbilities {
		values := make([]fixed.F64d4, len(a.Children))
		for i, one := range a.Children {
			values[i] = one.AdjustedPoints()
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
			points += one.AdjustedPoints()
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
	if a.Disabled {
		return false
	}
	p := a.Parent
	for p != nil {
		if p.Disabled {
			return false
		}
		p = p.Parent
	}
	return true
}

// TypeAsText returns the set of type bits that are set if this isn't a container, or an empty string if it is.
func (a *Advantage) TypeAsText() string {
	if a.Container() {
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
	if !a.Container() && a.Levels != nil && *a.Levels > 0 {
		buffer.WriteByte(' ')
		buffer.WriteString(a.Levels.String())
	}
	return buffer.String()
}

// FillWithNameableKeys adds any nameable keys found in this Advantage to the provided map.
func (a *Advantage) FillWithNameableKeys(nameables map[string]string) {
	ExtractNameables(a.Name, nameables)
	ExtractNameables(a.Notes, nameables)
	ExtractNameables(a.VTTNotes, nameables)
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
	a.Name = ApplyNameables(a.Name, nameables)
	a.Notes = ApplyNameables(a.Notes, nameables)
	a.VTTNotes = ApplyNameables(a.VTTNotes, nameables)
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
		if !one.Disabled && strings.EqualFold(one.Name, name) {
			return one
		}
	}
	return nil
}

// ModifierNotes returns the notes due to modifiers. 'entity' may be nil.
func (a *Advantage) ModifierNotes() string {
	var buffer strings.Builder
	if a.CR != advantage.None {
		buffer.WriteString(a.CR.String())
		if a.CRAdj != NoCRAdj {
			buffer.WriteString(", ")
			buffer.WriteString(a.CRAdj.Description(a.CR))
		}
	}
	for _, one := range a.Modifiers {
		if !one.Disabled {
			if buffer.Len() != 0 {
				buffer.WriteString("; ")
			}
			buffer.WriteString(one.FullDescription(a.Entity))
		}
	}
	return buffer.String()
}

// SecondaryText returns the "secondary" text: the text display below an Advantage.
func (a *Advantage) SecondaryText() string {
	var buffer strings.Builder
	settings := SheetSettingsFor(a.Entity)
	if a.UserDesc != "" && settings.UserDescriptionDisplay.Inline() {
		buffer.WriteString(a.UserDesc)
	}
	if settings.ModifiersDisplay.Inline() {
		if notes := a.ModifierNotes(); notes != "" {
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
		if !one.Disabled {
			modifier := one.CostModifier()
			switch one.CostType {
			case advantage.Percentage:
				switch *one.Affects {
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
				if *one.Affects == advantage.LevelsOnly {
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
