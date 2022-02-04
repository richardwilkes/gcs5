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
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	spellTypeKey            = "spell"
	ritualMagicSpellTypeKey = "ritual_magic_spell"
)

// SpellItem holds the Spell data that only exists in non-containers.
type SpellItem struct {
	TechLevel       string              `json:"tech_level,omitempty"`
	Difficulty      AttributeDifficulty `json:"difficulty"`
	College         []string            `json:"college,omitempty"`
	PowerSource     string              `json:"power_source,omitempty"`
	Class           string              `json:"spell_class,omitempty"`
	Resist          string              `json:"resist,omitempty"`
	CastingCost     string              `json:"casting_cost,omitempty"`
	MaintenanceCost string              `json:"maintenance_cost,omitempty"`
	CastingTime     string              `json:"casting_time,omitempty"`
	Duration        string              `json:"duration,omitempty"`
	Points          fixed.F64d4         `json:"points,omitempty"`
	Prereq          *PrereqList         `json:"prereqs,omitempty"`
	Weapons         []*Weapon           `json:"weapons,omitempty"`
}

// SpellContainer holds the Spell data that only exists in containers.
type SpellContainer struct {
	Children []*Spell `json:"children,omitempty"`
	Open     bool     `json:"open,omitempty"`
}

// SpellData holds the Spell data that is written to disk.
type SpellData struct {
	Type            string    `json:"type"`
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name,omitempty"`
	PageRef         string    `json:"reference,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	VTTNotes        string    `json:"vtt_notes,omitempty"`
	Categories      []string  `json:"categories,omitempty"`
	*SpellItem      `json:",omitempty"`
	*SpellContainer `json:",omitempty"`
}

// Spell holds the data for a spell.
type Spell struct {
	SpellData
	Entity            *Entity
	Parent            *Spell
	Level             skill.Level
	UnsatisfiedReason string
	Satisfied         bool
}

type spellListData struct {
	Current []*Spell `json:"spells"`
}

// NewSpellsFromFile loads an Spell list from a file.
func NewSpellsFromFile(fileSystem fs.FS, filePath string) ([]*Spell, error) {
	var data struct {
		spellListData
		OldKey []*Spell `json:"rows"`
	}
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause("invalid spells file: "+filePath, err)
	}
	if len(data.Current) != 0 {
		return data.Current, nil
	}
	return data.OldKey, nil
}

// SaveSpells writes the Spell list to the file as JSON.
func SaveSpells(spells []*Spell, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &spellListData{Current: spells})
}

// NewSpell creates a new Spell.
func NewSpell(entity *Entity, parent *Spell, container bool) *Spell {
	s := Spell{
		SpellData: SpellData{
			Type: spellTypeKey,
			ID:   id.NewUUID(),
			Name: i18n.Text("Spell"),
		},
		Entity: entity,
		Parent: parent,
	}
	if container {
		s.Type += commonContainerKeyPostfix
		s.SpellContainer = &SpellContainer{Open: true}
	} else {
		s.SpellItem = &SpellItem{
			Difficulty: AttributeDifficulty{
				Attribute:  AttributeIDFor(entity, gid.Intelligence),
				Difficulty: skill.Hard,
			},
			PowerSource: i18n.Text("Arcane"),
			Class:       i18n.Text("Regular"),
			CastingCost: "1",
			CastingTime: "1 sec",
			Duration:    "Instant",
			Points:      fxp.One,
			Prereq:      NewPrereqList(),
		}
	}
	return &s
}

// MarshalJSON implements json.Marshaler.
func (s *Spell) MarshalJSON() ([]byte, error) {
	if s.Container() {
		s.SpellItem = nil
	} else {
		s.SpellContainer = nil
		if s.Level.Level > 0 {
			type calc struct {
				Level              fixed.F64d4 `json:"level"`
				RelativeSkillLevel string      `json:"rsl"`
			}
			data := struct {
				SpellData
				Calc calc `json:"calc"`
			}{
				SpellData: s.SpellData,
				Calc: calc{
					Level: s.Level.Level,
				},
			}
			rsl := s.AdjustedRelativeLevel()
			switch {
			case rsl == math.MinInt:
				data.Calc.RelativeSkillLevel = "-"
			case s.Type != ritualMagicSpellTypeKey:
				s.Type = ResolveAttributeName(s.Entity, s.Difficulty.Attribute) + rsl.StringWithSign()
			default:
				s.Type = rsl.StringWithSign()
			}
			return json.Marshal(&data)
		}
	}
	return json.Marshal(&s.SpellData)
}

// Container returns true if this is a container.
func (s *Spell) Container() bool {
	return strings.HasSuffix(s.Type, commonContainerKeyPostfix)
}

// AdjustedRelativeLevel returns the relative skill level.
func (s *Spell) AdjustedRelativeLevel() fixed.F64d4 {
	if s.Container() {
		return fixed.F64d4Min
	}
	if s.Entity != nil && s.Level.Level > 0 {
		return s.Level.RelativeLevel
	}
	// TODO: Old code had a case for templates... but can't see that being exercised in the actual display anywhere
	return fixed.F64d4Min
}
