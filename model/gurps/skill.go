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
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	skillTypeKey     = "skill"
	techniqueTypeKey = "technique"
)

// SkillItem holds the Skill data that only exists in non-containers.
type SkillItem struct {
	Specialization               string              `json:"specialization,omitempty"`
	TechLevel                    string              `json:"tech_level,omitempty"`
	Difficulty                   AttributeDifficulty `json:"difficulty"`
	Points                       fixed.F64d4         `json:"points,omitempty"`
	EncumbrancePenaltyMultiplier fixed.F64d4         `json:"encumbrance_penalty_multiplier,omitempty"`
	DefaultedFrom                *SkillDefault       `json:"defaulted_from,omitempty"`
	Defaults                     []*SkillDefault     `json:"defaults,omitempty"`
	TechniqueDefault             *SkillDefault       `json:"default,omitempty"`
	TechniqueLimitModifier       *fixed.F64d4        `json:"limit,omitempty"`
	Prereq                       Prereq              `json:"prereqs,omitempty"`
	Weapons                      []*Weapon           `json:"weapons,omitempty"`
	Features                     feature.Features    `json:"features,omitempty"`
}

// SkillContainer holds the Skill data that only exists in containers.
type SkillContainer struct {
	Children []*Skill `json:"children,omitempty"`
	Open     bool     `json:"open,omitempty"`
}

// SkillData holds the Skill data that is written to disk.
type SkillData struct {
	Type            string    `json:"type"`
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name,omitempty"`
	PageRef         string    `json:"reference,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	VTTNotes        string    `json:"vtt_notes,omitempty"`
	Categories      []string  `json:"categories,omitempty"`
	*SkillItem      `json:",omitempty"`
	*SkillContainer `json:",omitempty"`
}

// Skill holds the data for a skill.
type Skill struct {
	SkillData
	Entity            *Entity
	Parent            *Skill
	Level             skill.Level
	UnsatisfiedReason string
	Satisfied         bool
}

// NewSkill creates a new Skill.
func NewSkill(entity *Entity, parent *Skill, container bool) *Skill {
	s := Skill{
		SkillData: SkillData{
			Type: skillTypeKey,
			ID:   id.NewUUID(),
			Name: i18n.Text("Skill"),
		},
		Entity: entity,
		Parent: parent,
	}
	if container {
		s.Type += commonContainerKeyPostfix
		s.SkillContainer = &SkillContainer{Open: true}
	} else {
		s.SkillItem = &SkillItem{
			Difficulty: AttributeDifficulty{
				Attribute:  AttributeIDFor(entity, "dx"),
				Difficulty: skill.Average,
			},
			Points: f64d4.One,
			Prereq: NewPrereqList(),
		}
	}
	return &s
}

// NewTechnique creates a new technique (i.e. a specialized use of a Skill). All parameters may be nil or empty.
func NewTechnique(entity *Entity, parent *Skill, skillName string) *Skill {
	t := NewSkill(entity, parent, false)
	t.Type = techniqueTypeKey
	t.Name = i18n.Text("Technique")
	if skillName == "" {
		skillName = i18n.Text("Skill")
	}
	t.TechniqueDefault = &SkillDefault{
		DefaultType: skillTypeKey,
		Name:        skillName,
	}
	return t
}

// MarshalJSON implements json.Marshaler.
func (s *Skill) MarshalJSON() ([]byte, error) {
	if s.Container() {
		s.SkillItem = nil
	} else {
		s.SkillContainer = nil
		if s.Level.Level > 0 {
			type calc struct {
				Level              fixed.F64d4 `json:"level"`
				RelativeSkillLevel string      `json:"rsl"`
			}
			data := struct {
				SkillData
				Calc calc `json:"calc"`
			}{
				SkillData: s.SkillData,
				Calc: calc{
					Level: s.Level.Level,
				},
			}
			rsl := s.AdjustedRelativeLevel()
			switch {
			case rsl == math.MinInt:
				data.Calc.RelativeSkillLevel = "-"
			case s.Type != techniqueTypeKey:
				s.Type = ResolveAttributeName(s.Entity, s.Difficulty.Attribute) + rsl.StringWithSign()
			default:
				s.Type = rsl.StringWithSign()
			}
			return json.Marshal(&data)
		}
	}
	return json.Marshal(&s.SkillData)
}

// Container returns true if this is a container.
func (s *Skill) Container() bool {
	return strings.HasSuffix(s.Type, commonContainerKeyPostfix)
}

// AdjustedRelativeLevel returns the relative skill level.
func (s *Skill) AdjustedRelativeLevel() fixed.F64d4 {
	if s.Container() {
		return fixed.F64d4Min
	}
	if s.Entity != nil && s.Level.Level > 0 {
		if s.Type == techniqueTypeKey {
			return s.Level.RelativeLevel + s.TechniqueDefault.Modifier
		}
		return s.Level.RelativeLevel
	}
	// TODO: Old code had a case for templates... but can't see that being exercised in the actual display anywhere
	return fixed.F64d4Min
}
