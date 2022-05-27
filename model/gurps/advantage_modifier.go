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
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/txt"
	"golang.org/x/exp/slices"
)

var _ node.Node = &AdvantageModifier{}

// Columns that can be used with the advantage modifier method .CellData()
const (
	AdvantageModifierDescriptionColumn = iota
	AdvantageModifierCostColumn
	AdvantageModifierTagsColumn
	AdvantageModifierReferenceColumn
)

const (
	advantageModifierListTypeKey = "modifier_list"
	advantageModifierTypeKey     = "modifier"
)

// AdvantageModifier holds a modifier to an Advantage.
type AdvantageModifier struct {
	AdvantageModifierData
	Entity *Entity
}

type advantageModifierListData struct {
	Type    string               `json:"type"`
	Version int                  `json:"version"`
	Rows    []*AdvantageModifier `json:"rows"`
}

// NewAdvantageModifiersFromFile loads an AdvantageModifier list from a file.
func NewAdvantageModifiersFromFile(fileSystem fs.FS, filePath string) ([]*AdvantageModifier, error) {
	var data advantageModifierListData
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause(gid.InvalidFileDataMsg, err)
	}
	if data.Type != advantageModifierListTypeKey {
		return nil, errs.New(gid.UnexpectedFileDataMsg)
	}
	if err := gid.CheckVersion(data.Version); err != nil {
		return nil, err
	}
	return data.Rows, nil
}

// SaveAdvantageModifiers writes the AdvantageModifier list to the file as JSON.
func SaveAdvantageModifiers(modifiers []*AdvantageModifier, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &advantageModifierListData{
		Type:    advantageModifierListTypeKey,
		Version: gid.CurrentDataVersion,
		Rows:    modifiers,
	})
}

// NewAdvantageModifier creates an AdvantageModifier.
func NewAdvantageModifier(entity *Entity, _ *AdvantageModifier, container bool) *AdvantageModifier {
	a := &AdvantageModifier{
		AdvantageModifierData: AdvantageModifierData{
			ContainerBase: newContainerBase[*AdvantageModifier](advantageModifierTypeKey, container),
		},
		Entity: entity,
	}
	a.Name = a.Kind()
	return a
}

// Clone creates a copy of this data.
func (a *AdvantageModifier) Clone() *AdvantageModifier {
	other := *a
	other.Tags = txt.CloneStringSlice(a.Tags)
	other.Features = a.Features.Clone()
	other.Children = nil
	if len(a.Children) != 0 {
		other.Children = make([]*AdvantageModifier, 0, len(a.Children))
		for _, one := range a.Children {
			other.Children = append(other.Children, one.Clone())
		}
	}
	return &other
}

// MarshalJSON implements json.Marshaler.
func (a *AdvantageModifier) MarshalJSON() ([]byte, error) {
	a.ClearUnusedFieldsForType()
	return json.Marshal(&a.AdvantageModifierData)
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AdvantageModifier) UnmarshalJSON(data []byte) error {
	var localData struct {
		AdvantageModifierData
		// Old data fields
		Categories []string `json:"categories"`
	}
	if err := json.Unmarshal(data, &localData); err != nil {
		return err
	}
	localData.ClearUnusedFieldsForType()
	a.AdvantageModifierData = localData.AdvantageModifierData
	a.Tags = convertOldCategoriesToTags(a.Tags, localData.Categories)
	slices.Sort(a.Tags)
	return nil
}

// CellData returns the cell data information for the given column.
func (a *AdvantageModifier) CellData(column int, data *node.CellData) {
	switch column {
	case AdvantageModifierDescriptionColumn:
		data.Type = node.Text
		data.Primary = a.Name
		data.Secondary = a.SecondaryText()
	case AdvantageModifierCostColumn:
		if !a.Container() {
			data.Type = node.Text
			data.Primary = a.CostDescription()
		}
	case AdvantageModifierTagsColumn:
		data.Type = node.Text
		data.Primary = CombineTags(a.Tags)
	case AdvantageModifierReferenceColumn, node.PageRefCellAlias:
		data.Type = node.PageRef
		data.Primary = a.PageRef
		data.Secondary = a.Name
	}
}

// OwningEntity returns the owning Entity.
func (a *AdvantageModifier) OwningEntity() *Entity {
	return a.Entity
}

// SetOwningEntity sets the owning entity and configures any sub-components as needed.
func (a *AdvantageModifier) SetOwningEntity(entity *Entity) {
	a.Entity = entity
	if a.Container() {
		for _, child := range a.Children {
			child.SetOwningEntity(entity)
		}
	}
}

// CostModifier returns the total cost modifier.
func (a *AdvantageModifier) CostModifier() fxp.Int {
	if a.Levels > 0 {
		return a.Cost.Mul(a.Levels)
	}
	return a.Cost
}

// HasLevels returns true if this AdvantageModifier has levels.
func (a *AdvantageModifier) HasLevels() bool {
	return !a.Container() && a.CostType == advantage.Percentage && a.Levels > 0
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

// SecondaryText returns the "secondary" text: the text display below an Advantage.
func (a *AdvantageModifier) SecondaryText() string {
	var buffer strings.Builder
	settings := SheetSettingsFor(a.Entity)
	if a.LocalNotes != "" && settings.NotesDisplay.Inline() {
		if buffer.Len() != 0 {
			buffer.WriteByte('\n')
		}
		buffer.WriteString(a.LocalNotes)
	}
	return buffer.String()
}

// FullDescription returns a full description.
func (a *AdvantageModifier) FullDescription() string {
	var buffer strings.Builder
	buffer.WriteString(a.String())
	if a.LocalNotes != "" {
		buffer.WriteString(" (")
		buffer.WriteString(a.LocalNotes)
		buffer.WriteByte(')')
	}
	if SheetSettingsFor(a.Entity).ShowAdvantageModifierAdj {
		buffer.WriteString(" [")
		buffer.WriteString(a.CostDescription())
		buffer.WriteByte(']')
	}
	return buffer.String()
}

// CostDescription returns the formatted cost.
func (a *AdvantageModifier) CostDescription() string {
	if a.Container() {
		return ""
	}
	var base string
	switch a.CostType {
	case advantage.Percentage:
		if a.HasLevels() {
			base = a.Cost.Mul(a.Levels).StringWithSign()
		} else {
			base = a.Cost.StringWithSign()
		}
		base += advantage.Percentage.String()
	case advantage.Points:
		base = a.Cost.StringWithSign()
	case advantage.Multiplier:
		return a.CostType.String() + a.Cost.String()
	default:
		jot.Errorf("unhandled cost type: %d", a.CostType)
		base = a.Cost.StringWithSign() + advantage.Percentage.String()
	}
	if desc := a.Affects.AltString(); desc != "" {
		base += " " + desc
	}
	return base
}

// FillWithNameableKeys adds any nameable keys found in this AdvantageModifier to the provided map.
func (a *AdvantageModifier) FillWithNameableKeys(m map[string]string) {
	if !a.Disabled {
		nameables.Extract(a.Name, m)
		nameables.Extract(a.LocalNotes, m)
		nameables.Extract(a.VTTNotes, m)
		for _, one := range a.Features {
			one.FillWithNameableKeys(m)
		}
	}
}

// ApplyNameableKeys replaces any nameable keys found in this AdvantageModifier with the corresponding values in the provided map.
func (a *AdvantageModifier) ApplyNameableKeys(m map[string]string) {
	if !a.Disabled {
		a.Name = nameables.Apply(a.Name, m)
		a.LocalNotes = nameables.Apply(a.LocalNotes, m)
		a.VTTNotes = nameables.Apply(a.VTTNotes, m)
		for _, one := range a.Features {
			one.ApplyNameableKeys(m)
		}
	}
}
