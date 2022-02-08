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
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// SpellItem holds the Spell data that only exists in non-containers.
type SpellItem struct {
	TechLevel       *string             `json:"tech_level,omitempty"`
	Difficulty      AttributeDifficulty `json:"difficulty"`
	College         []string            `json:"college,omitempty"`
	PowerSource     string              `json:"power_source,omitempty"`
	Class           string              `json:"spell_class,omitempty"`
	Resist          string              `json:"resist,omitempty"`
	CastingCost     string              `json:"casting_cost,omitempty"`
	MaintenanceCost string              `json:"maintenance_cost,omitempty"`
	CastingTime     string              `json:"casting_time,omitempty"`
	Duration        string              `json:"duration,omitempty"`
	RitualSkillName string              `json:"base_skill,omitempty"`
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
	LevelData         skill.Level
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
			Type: gid.Spell,
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
		if s.LevelData.Level > 0 {
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
					Level: s.LevelData.Level,
				},
			}
			rsl := s.AdjustedRelativeLevel()
			switch {
			case rsl == math.MinInt:
				data.Calc.RelativeSkillLevel = "-"
			case s.Type != gid.RitualMagicSpell:
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
	if s.Entity != nil && s.LevelData.Level > 0 {
		return s.LevelData.RelativeLevel
	}
	// TODO: Old code had a case for templates... but can't see that being exercised in the actual display anywhere
	return fixed.F64d4Min
}

// UpdateLevel updates the level of the spell, returning true if it has changed.
func (s *Spell) UpdateLevel() bool {
	saved := s.LevelData
	s.LevelData = s.CalculateLevel(s.Points, s.Difficulty, s.College, s.Categories, s.PowerSource, s.Name)
	return saved != s.LevelData
}

// Level returns the computed level.
func (s *Spell) Level() fixed.F64d4 {
	return s.CalculateLevel(s.Points, s.Difficulty, s.College, s.Categories, s.PowerSource, s.Name).Level
}

// CalculateLevel computes the level.
func (s *Spell) CalculateLevel(points fixed.F64d4, attrDiff AttributeDifficulty, colleges, categories []string, powerSource, name string) skill.Level {
	var tooltip xio.ByteBuffer
	relativeLevel := attrDiff.Difficulty.BaseRelativeLevel()
	level := fxp.NegOne
	if s.Entity != nil {
		level = s.Entity.ResolveAttributeCurrent(attrDiff.Attribute)
		if attrDiff.Difficulty == skill.Wildcard {
			points = points.Div(fxp.Three).Trunc()
		}
		switch {
		case points < fxp.One:
			level = fxp.NegOne
			relativeLevel = 0
		case points == fxp.One:
		// relativeLevel is preset to this point value
		case points < fxp.Four:
			relativeLevel += fxp.One
		default:
			relativeLevel += fxp.One + points.Div(fxp.Four).Trunc()
		}
		if level != fxp.One {
			relativeLevel += s.BestCollegeSpellBonus(categories, colleges, &tooltip)
			relativeLevel += s.SpellBonusesFor(feature.SpellPowerSourceID, powerSource, categories, &tooltip)
			relativeLevel += s.SpellBonusesFor(feature.SpellNameID, name, categories, &tooltip)
			level += relativeLevel
		}
	}
	return skill.Level{
		Level:         level,
		RelativeLevel: relativeLevel,
		Tooltip:       tooltip.String(),
	}
}

// BestCollegeSpellBonus returns the best college spell bonus for this spell.
func (s *Spell) BestCollegeSpellBonus(categories, colleges []string, tooltip *xio.ByteBuffer) fixed.F64d4 {
	best := fixed.F64d4Min
	var bestTooltip string
	for _, college := range colleges {
		var buffer *xio.ByteBuffer
		if tooltip != nil {
			buffer = &xio.ByteBuffer{}
		}
		if pts := s.SpellBonusesFor(feature.SpellCollegeID, college, categories, buffer); best < pts {
			best = pts
			if buffer != nil {
				bestTooltip = buffer.String()
			}
		}
	}
	if tooltip != nil {
		tooltip.WriteString(bestTooltip)
	}
	if best == fixed.F64d4Min {
		best = 0
	}
	return best
}

// SpellBonusesFor returns the bonus for this spell.
func (s *Spell) SpellBonusesFor(featureID, qualifier string, categories []string, tooltip *xio.ByteBuffer) fixed.F64d4 {
	level := s.Entity.BonusFor(featureID, tooltip)
	level += s.Entity.BonusFor(featureID+"/"+strings.ToLower(qualifier), tooltip)
	level += s.Entity.SpellComparedBonusFor(featureID+"*", qualifier, categories, tooltip)
	return level
}

// RitualMagicSatisfied returns true if the Ritual Magic Spell is satisfied.
func (s *Spell) RitualMagicSatisfied(tooltip *xio.ByteBuffer, prefix string) bool {
	if s.Type != gid.RitualMagicSpell {
		return true
	}
	if len(s.College) == 0 {
		if tooltip != nil {
			tooltip.WriteString(prefix)
			tooltip.WriteString(i18n.Text("Must be assigned to a college"))
		}
		return false
	}
	for _, college := range s.College {
		if s.Entity.BestSkillNamed(s.RitualSkillName, college, false, nil) != nil {
			return true
		}
	}
	if s.Entity.BestSkillNamed(s.RitualSkillName, "", false, nil) != nil {
		return true
	}
	if tooltip != nil {
		tooltip.WriteString(prefix)
		tooltip.WriteString(i18n.Text("Requires a skill named "))
		tooltip.WriteString(s.RitualSkillName)
		tooltip.WriteString(" (")
		tooltip.WriteString(s.College[0])
		tooltip.WriteByte(')')
		for _, college := range s.College[1:] {
			tooltip.WriteString(i18n.Text(" or "))
			tooltip.WriteString(s.RitualSkillName)
			tooltip.WriteString(" (")
			tooltip.WriteString(college)
			tooltip.WriteByte(')')
		}
	}
	return false
}
