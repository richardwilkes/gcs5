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
	"github.com/richardwilkes/gcs/model/gurps/advantage"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
	"github.com/richardwilkes/unison"
)

var (
	_ WeaponOwner = &Advantage{}
	_ node.Node   = &Advantage{}
)

// Columns that can be used with the advantage method .CellData()
const (
	AdvantageDescriptionColumn = iota
	AdvantagePointsColumn
	AdvantageTypeColumn
	AdvantageCategoryColumn
	AdvantageReferenceColumn
)

const (
	advantageListTypeKey = "advantage_list"
	advantageTypeKey     = "advantage"
)

// Advantage holds an advantage, disadvantage, quirk, or perk.
type Advantage struct {
	AdvantageData
	Entity            *Entity
	Parent            *Advantage
	UnsatisfiedReason string
	Satisfied         bool
}

type advantageListData struct {
	Type    string       `json:"type"`
	Version int          `json:"version"`
	Rows    []*Advantage `json:"rows"`
}

// NewAdvantagesFromFile loads an Advantage list from a file.
func NewAdvantagesFromFile(fileSystem fs.FS, filePath string) ([]*Advantage, error) {
	var data advantageListData
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause(gid.InvalidFileDataMsg, err)
	}
	if data.Type != advantageListTypeKey {
		return nil, errs.New(gid.UnexpectedFileDataMsg)
	}
	if err := gid.CheckVersion(data.Version); err != nil {
		return nil, err
	}
	return data.Rows, nil
}

// SaveAdvantages writes the Advantage list to the file as JSON.
func SaveAdvantages(advantages []*Advantage, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &advantageListData{
		Type:    advantageListTypeKey,
		Version: gid.CurrentDataVersion,
		Rows:    advantages,
	})
}

// NewAdvantage creates a new Advantage.
func NewAdvantage(entity *Entity, parent *Advantage, container bool) *Advantage {
	a := Advantage{
		AdvantageData: AdvantageData{
			Type: advantageTypeKey,
			ID:   id.NewUUID(),
			AdvantageEditData: AdvantageEditData{
				Name: i18n.Text("Advantage"),
			},
		},
		Entity: entity,
		Parent: parent,
	}
	if container {
		a.Type += commonContainerKeyPostfix
		a.IsOpen = true
	}
	return &a
}

// MarshalJSON implements json.Marshaler.
func (a *Advantage) MarshalJSON() ([]byte, error) {
	type calc struct {
		Points f64d4.Int `json:"points"`
	}
	a.ClearUnusedFieldsForType()
	data := struct {
		AdvantageData
		Calc calc `json:"calc"`
	}{
		AdvantageData: a.AdvantageData,
		Calc: calc{
			Points: a.AdjustedPoints(),
		},
	}
	return json.Marshal(&data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Advantage) UnmarshalJSON(data []byte) error {
	a.AdvantageData = AdvantageData{}
	if err := json.Unmarshal(data, &a.AdvantageData); err != nil {
		return err
	}
	a.ClearUnusedFieldsForType()
	if a.Container() {
		for _, one := range a.Children {
			one.Parent = a
		}
	}
	return nil
}

// CellData returns the cell data information for the given column.
func (a *Advantage) CellData(column int, data *node.CellData) {
	switch column {
	case AdvantageDescriptionColumn:
		data.Type = node.Text
		data.Primary = a.String()
		data.Secondary = a.SecondaryText()
		data.Disabled = a.Disabled
	case AdvantagePointsColumn:
		data.Type = node.Text
		data.Primary = a.AdjustedPoints().String()
		data.Alignment = unison.EndAlignment
	case AdvantageTypeColumn:
		data.Type = node.Text
		data.Primary = a.TypeAsText()
	case AdvantageCategoryColumn:
		data.Type = node.Text
		data.Primary = CombineTags(a.Categories)
	case AdvantageReferenceColumn:
		data.Type = node.PageRef
		data.Primary = a.PageRef
		data.Secondary = a.Name
	}
}

// Depth returns the number of parents this node has.
func (a *Advantage) Depth() int {
	count := 0
	p := a.Parent
	for p != nil {
		count++
		p = p.Parent
	}
	return count
}

// OwningEntity returns the owning Entity.
func (a *Advantage) OwningEntity() *Entity {
	return a.Entity
}

// SetOwningEntity sets the owning entity and configures any sub-components as needed.
func (a *Advantage) SetOwningEntity(entity *Entity) {
	a.Entity = entity
	if a.Container() {
		for _, child := range a.Children {
			child.SetOwningEntity(entity)
		}
	} else {
		for _, w := range a.Weapons {
			w.SetOwner(a)
		}
	}
	for _, w := range a.Modifiers {
		w.SetOwningEntity(entity)
	}
}

// Notes returns the local notes.
func (a *Advantage) Notes() string {
	return a.LocalNotes
}

// IsLeveled returns true if the Advantage is capable of having levels.
func (a *Advantage) IsLeveled() bool {
	return !a.Container() && a.PointsPerLevel != 0
}

// AdjustedPoints returns the total points, taking levels and modifiers into account.
func (a *Advantage) AdjustedPoints() f64d4.Int {
	if a.Disabled {
		return 0
	}
	if !a.Container() {
		return AdjustedPoints(a.Entity, a.BasePoints, a.Levels, a.PointsPerLevel, a.CR, a.AllModifiers(), a.RoundCostDown)
	}
	var points f64d4.Int
	if a.ContainerType == advantage.AlternativeAbilities {
		values := make([]f64d4.Int, len(a.Children))
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
				points += fxp.ApplyRounding(calculateModifierPoints(v, fxp.Twenty), a.RoundCostDown)
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
	return CombineTags(list)
}

// Description returns a description, which doesn't include any levels.
func (a *Advantage) Description() string {
	return a.Name
}

// String implements fmt.Stringer.
func (a *Advantage) String() string {
	var buffer strings.Builder
	buffer.WriteString(a.Name)
	if a.IsLeveled() && a.Levels > 0 {
		buffer.WriteByte(' ')
		buffer.WriteString(a.Levels.String())
	}
	return buffer.String()
}

// FeatureList returns the list of Features.
func (a *Advantage) FeatureList() feature.Features {
	return a.Features
}

// CategoryList returns the list of categories.
func (a *Advantage) CategoryList() []string {
	return a.Categories
}

// FillWithNameableKeys adds any nameable keys found in this Advantage to the provided map.
func (a *Advantage) FillWithNameableKeys(m map[string]string) {
	nameables.Extract(a.Name, m)
	nameables.Extract(a.LocalNotes, m)
	nameables.Extract(a.VTTNotes, m)
	if a.Prereq != nil {
		a.Prereq.FillWithNameableKeys(m)
	}
	for _, one := range a.Features {
		one.FillWithNameableKeys(m)
	}
	for _, one := range a.Weapons {
		one.FillWithNameableKeys(m)
	}
	for _, one := range a.Modifiers {
		one.FillWithNameableKeys(m)
	}
}

// ApplyNameableKeys replaces any nameable keys found in this Advantage with the corresponding values in the provided map.
func (a *Advantage) ApplyNameableKeys(m map[string]string) {
	a.Name = nameables.Apply(a.Name, m)
	a.LocalNotes = nameables.Apply(a.LocalNotes, m)
	a.VTTNotes = nameables.Apply(a.VTTNotes, m)
	if a.Prereq != nil {
		a.Prereq.ApplyNameableKeys(m)
	}
	for _, one := range a.Features {
		one.ApplyNameableKeys(m)
	}
	for _, one := range a.Weapons {
		one.ApplyNameableKeys(m)
	}
	for _, one := range a.Modifiers {
		one.ApplyNameableKeys(m)
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

// ModifierNotes returns the notes due to modifiers.
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
			buffer.WriteString(one.FullDescription())
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
	if a.LocalNotes != "" && settings.NotesDisplay.Inline() {
		if buffer.Len() != 0 {
			buffer.WriteByte('\n')
		}
		buffer.WriteString(a.LocalNotes)
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

// CombineTags combines multiple tags into a single string.
func CombineTags(tags []string) string {
	return strings.Join(tags, ", ")
}

// ExtractTags from a combined tags string.
func ExtractTags(tags string) []string {
	var list []string
	for _, one := range strings.Split(tags, ",") {
		if one = strings.TrimSpace(one); one != "" {
			list = append(list, one)
		}
	}
	return list
}

// AdjustedPoints returns the total points, taking levels and modifiers into account. 'entity' may be nil.
func AdjustedPoints(entity *Entity, basePoints, levels, pointsPerLevel f64d4.Int, cr advantage.SelfControlRoll, modifiers []*AdvantageModifier, roundCostDown bool) f64d4.Int {
	var baseEnh, levelEnh, baseLim, levelLim f64d4.Int
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
				modifiedBasePoints = modifyPoints(modifyPoints(modifiedBasePoints+leveledPoints, baseEnh), fxp.NegEighty.Max(baseLim))
			} else {
				modifiedBasePoints = modifyPoints(modifyPoints(modifiedBasePoints, baseEnh), fxp.NegEighty.Max(baseLim)) +
					modifyPoints(modifyPoints(leveledPoints, levelEnh), fxp.NegEighty.Max(levelLim))
			}
		} else {
			baseMod := fxp.NegEighty.Max(baseEnh + baseLim)
			levelMod := fxp.NegEighty.Max(levelEnh + levelLim)
			if baseMod == levelMod {
				modifiedBasePoints = modifyPoints(modifiedBasePoints+leveledPoints, baseMod)
			} else {
				modifiedBasePoints = modifyPoints(modifiedBasePoints, baseMod) + modifyPoints(leveledPoints, levelMod)
			}
		}
	} else {
		modifiedBasePoints += leveledPoints
	}
	return fxp.ApplyRounding(modifiedBasePoints.Mul(multiplier), roundCostDown)
}

func modifyPoints(points, modifier f64d4.Int) f64d4.Int {
	return points + calculateModifierPoints(points, modifier)
}

func calculateModifierPoints(points, modifier f64d4.Int) f64d4.Int {
	return points.Mul(modifier).Div(fxp.Hundred)
}
