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
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/i18n"
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
	Points                       int                 `json:"points,omitempty"`
	EncumbrancePenaltyMultiplier int                 `json:"encumbrance_penalty_multiplier,omitempty"`
	DefaultedFrom                *SkillDefault       `json:"defaulted_from,omitempty"`
	Defaults                     []*SkillDefault     `json:"defaults,omitempty"`
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
				Attribute:  "dx",
				Difficulty: skill.A,
			},
			Points: 1,
			Prereq: NewPrereqList(),
		}
	}
	return &s
}

// MarshalJSON implements json.Marshaler.
func (s *Skill) MarshalJSON() ([]byte, error) {
	if s.Container() {
		s.SkillItem = nil
	} else {
		s.SkillContainer = nil
	}
	if s.Level.Level > 0 {
		type calc struct {
			Level              int    `json:"level"`
			RelativeSkillLevel string `json:"rsl"`
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
			s.Type = fmt.Sprintf("%s%+d", ResolveAttributeName(s.Entity, s.Difficulty.Attribute), rsl)
		default:
			s.Type = fmt.Sprintf("%+d", rsl)
		}
		return json.Marshal(&data)
	}
	return json.Marshal(&s.SkillData)
}

// Container returns true if this is a container.
func (s *Skill) Container() bool {
	return strings.HasSuffix(s.Type, commonContainerKeyPostfix)
}

// AdjustedRelativeLevel returns the relative skill level.
func (s *Skill) AdjustedRelativeLevel() int {
	if s.Container() {
		return math.MinInt
	}
	// TODO: Implement
	/*
	   if (getCharacter() != null) {
	       if (getLevel() < 0) {
	           return Integer.MIN_VALUE;
	       }
	       int level = getRelativeLevel();
	       if (this instanceof Technique) {
	           level += ((Technique) this).getDefault().getModifier();
	       }
	       return level;
	   } else if (getTemplate() != null) {
	       int points = getPoints();
	       if (points > 0) {
	           SkillDifficulty difficulty = getDifficulty();
	           int             level;
	           if (this instanceof Technique) {
	               if (difficulty != SkillDifficulty.A) {
	                   points--;
	               }
	               return points + ((Technique) this).getDefault().getModifier();
	           }
	           level = difficulty.getBaseRelativeLevel();
	           if (difficulty == SkillDifficulty.W) {
	               points /= 3;
	           }
	           if (points > 1) {
	               if (points < 4) {
	                   level++;
	               } else {
	                   level += 1 + points / 4;
	               }
	           }
	           return level;
	       }
	   }
	*/
	return math.MinInt
}
