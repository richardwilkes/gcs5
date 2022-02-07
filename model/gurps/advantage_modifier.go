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
	"github.com/richardwilkes/gcs/model/gurps/advantage"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const advantageModifierTypeKey = "modifier"

// AdvantageModifierItem holds the AdvantageModifier data that only exists in non-containers.
type AdvantageModifierItem struct {
	CostType advantage.ModifierCostType `json:"cost_type,omitempty"`
	Disabled bool                       `json:"disabled,omitempty"`
	Cost     fixed.F64d4                `json:"cost,omitempty"`
	Levels   fixed.F64d4                `json:"levels,omitempty"`
	Affects  *advantage.Affects         `json:"affects,omitempty"`
	Features feature.Features           `json:"features,omitempty"`
}

// AdvantageModifierContainer holds the AdvantageModifier data that only exists in containers.
type AdvantageModifierContainer struct {
	Children []*AdvantageModifier `json:"children,omitempty"`
	Open     bool                 `json:"open,omitempty"`
}

// AdvantageModifierData holds the AdvantageModifier data that is written to disk.
type AdvantageModifierData struct {
	Type                        string    `json:"type"`
	ID                          uuid.UUID `json:"id"`
	Name                        string    `json:"name,omitempty"`
	PageRef                     string    `json:"reference,omitempty"`
	Notes                       string    `json:"notes,omitempty"`
	VTTNotes                    string    `json:"vtt_notes,omitempty"`
	*AdvantageModifierItem      `json:",omitempty"`
	*AdvantageModifierContainer `json:",omitempty"`
}

// AdvantageModifier holds a modifier to an Advantage.
type AdvantageModifier struct {
	AdvantageModifierData
}

type advantageModifierListData struct {
	Current []*AdvantageModifier `json:"advantage_modifiers"`
}

// NewAdvantageModifiersFromFile loads an AdvantageModifier list from a file.
func NewAdvantageModifiersFromFile(fileSystem fs.FS, filePath string) ([]*AdvantageModifier, error) {
	var data struct {
		advantageModifierListData
		OldKey []*AdvantageModifier `json:"rows"`
	}
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause("invalid advantage modifiers file: "+filePath, err)
	}
	if len(data.Current) != 0 {
		return data.Current, nil
	}
	return data.OldKey, nil
}

// SaveAdvantageModifiers writes the AdvantageModifier list to the file as JSON.
func SaveAdvantageModifiers(modifiers []*AdvantageModifier, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &advantageModifierListData{Current: modifiers})
}

// NewAdvantageModifier creates an AdvantageModifier.
func NewAdvantageModifier(container bool) *AdvantageModifier {
	a := AdvantageModifier{
		AdvantageModifierData: AdvantageModifierData{
			Type: advantageModifierTypeKey,
			ID:   id.NewUUID(),
			Name: i18n.Text("Advantage Modifier"),
		},
	}
	if container {
		a.Type += commonContainerKeyPostfix
		a.AdvantageModifierContainer = &AdvantageModifierContainer{Open: true}
	} else {
		affects := advantage.Total
		a.AdvantageModifierItem = &AdvantageModifierItem{
			Affects: &affects,
		}
	}
	return &a
}

// MarshalJSON implements json.Marshaler.
func (a *AdvantageModifier) MarshalJSON() ([]byte, error) {
	if a.Container() {
		a.AdvantageModifierItem = nil
	} else {
		a.AdvantageModifierContainer = nil
	}
	return json.Marshal(&a.AdvantageModifierData)
}

// Container returns true if this is a container.
func (a *AdvantageModifier) Container() bool {
	return strings.HasSuffix(a.Type, commonContainerKeyPostfix)
}

// CostModifier returns the total cost modifier.
func (a *AdvantageModifier) CostModifier() fixed.F64d4 {
	if a.Levels > 0 {
		return a.Cost.Mul(a.Levels)
	}
	return a.Cost
}

// HasLevels returns true if this AdvantageModifier has levels.
func (a *AdvantageModifier) HasLevels() bool {
	return a.CostType == advantage.Percentage && a.Levels > 0
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

// FullDescription returns a full description. 'entity' may be nil.
func (a *AdvantageModifier) FullDescription(entity *Entity) string {
	var buffer strings.Builder
	buffer.WriteString(a.String())
	if a.Notes != "" {
		buffer.WriteString(" (")
		buffer.WriteString(a.Notes)
		buffer.WriteByte(')')
	}
	if entity != nil && SheetSettingsFor(entity).ShowAdvantageModifierAdj {
		buffer.WriteString(" [")
		buffer.WriteString(a.CostDescription())
		buffer.WriteByte(']')
	}
	return buffer.String()
}

// CostDescription returns the formatted cost.
func (a *AdvantageModifier) CostDescription() string {
	var buffer strings.Builder
	if a.CostType == advantage.Multiplier {
		buffer.WriteByte('x')
		buffer.WriteString(a.Cost.String())
	} else {
		if a.Cost >= 0 {
			buffer.WriteByte('+')
		}
		buffer.WriteString(a.Cost.String())
		if a.CostType == advantage.Percentage {
			buffer.WriteByte('%')
		}
		if desc := a.Affects.AltString(); desc != "" {
			buffer.WriteByte(' ')
			buffer.WriteString(desc)
		}
	}
	return buffer.String()
}

// FillWithNameableKeys adds any nameable keys found in this AdvantageModifier to the provided map.
func (a *AdvantageModifier) FillWithNameableKeys(m map[string]string) {
	if !a.Disabled {
		nameables.Extract(a.Name, m)
		nameables.Extract(a.Notes, m)
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
		a.Notes = nameables.Apply(a.Notes, m)
		a.VTTNotes = nameables.Apply(a.VTTNotes, m)
		for _, one := range a.Features {
			one.ApplyNameableKeys(m)
		}
	}
}
