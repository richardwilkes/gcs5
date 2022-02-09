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
	"fmt"
	"io/fs"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/datafile"
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

// SkillItem holds the Skill data that only exists in non-containers.
type SkillItem struct {
	Specialization               string              `json:"specialization,omitempty"`
	TechLevel                    *string             `json:"tech_level,omitempty"`
	Difficulty                   AttributeDifficulty `json:"difficulty"`
	Points                       fixed.F64d4         `json:"points,omitempty"`
	EncumbrancePenaltyMultiplier fixed.F64d4         `json:"encumbrance_penalty_multiplier,omitempty"`
	DefaultedFrom                *SkillDefault       `json:"defaulted_from,omitempty"`
	Defaults                     []*SkillDefault     `json:"defaults,omitempty"`
	TechniqueDefault             *SkillDefault       `json:"default,omitempty"`
	TechniqueLimitModifier       *fixed.F64d4        `json:"limit,omitempty"`
	Prereq                       *PrereqList         `json:"prereqs,omitempty"`
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
	LocalNotes      string    `json:"notes,omitempty"`
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
	LevelData         skill.Level
	UnsatisfiedReason string
	Satisfied         bool
}

type skillListData struct {
	Current []*Skill `json:"skills"`
}

// NewSkillsFromFile loads an Skill list from a file.
func NewSkillsFromFile(fileSystem fs.FS, filePath string) ([]*Skill, error) {
	var data struct {
		skillListData
		OldKey []*Skill `json:"rows"`
	}
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause("invalid skills file: "+filePath, err)
	}
	if len(data.Current) != 0 {
		return data.Current, nil
	}
	return data.OldKey, nil
}

// SaveSkills writes the Skill list to the file as JSON.
func SaveSkills(skills []*Skill, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &skillListData{Current: skills})
}

// NewSkill creates a new Skill.
func NewSkill(entity *Entity, parent *Skill, container bool) *Skill {
	s := Skill{
		SkillData: SkillData{
			Type: gid.Skill,
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
				Attribute:  AttributeIDFor(entity, gid.Dexterity),
				Difficulty: skill.Average,
			},
			Points: fxp.One,
			Prereq: NewPrereqList(),
		}
	}
	return &s
}

// NewTechnique creates a new technique (i.e. a specialized use of a Skill). All parameters may be nil or empty.
func NewTechnique(entity *Entity, parent *Skill, skillName string) *Skill {
	t := NewSkill(entity, parent, false)
	t.Type = gid.Technique
	t.Name = i18n.Text("Technique")
	if skillName == "" {
		skillName = i18n.Text("Skill")
	}
	t.TechniqueDefault = &SkillDefault{
		DefaultType: gid.Skill,
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
		if s.LevelData.Level > 0 {
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
					Level: s.LevelData.Level,
				},
			}
			rsl := s.AdjustedRelativeLevel()
			switch {
			case rsl == math.MinInt:
				data.Calc.RelativeSkillLevel = "-"
			case s.Type != gid.Technique:
				data.Calc.RelativeSkillLevel = ResolveAttributeName(s.Entity, s.Difficulty.Attribute) + rsl.StringWithSign()
			default:
				data.Calc.RelativeSkillLevel = rsl.StringWithSign()
			}
			return json.Marshal(&data)
		}
	}
	return json.Marshal(&s.SkillData)
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Skill) UnmarshalJSON(data []byte) error {
	s.SkillData = SkillData{}
	if err := json.Unmarshal(data, &s.SkillData); err != nil {
		return err
	}
	if s.Container() {
		for _, one := range s.Children {
			one.Parent = s
		}
	} else if s.Prereq == nil {
		s.Prereq = NewPrereqList()
	}
	return nil
}

// Container returns true if this is a container.
func (s *Skill) Container() bool {
	return strings.HasSuffix(s.Type, commonContainerKeyPostfix)
}

// OwningEntity returns the owning Entity.
func (s *Skill) OwningEntity() *Entity {
	return s.Entity
}

// SetOwningEntity sets the owning entity and configures any sub-components as needed.
func (s *Skill) SetOwningEntity(entity *Entity) {
	s.Entity = entity
	if s.Container() {
		for _, child := range s.Children {
			child.SetOwningEntity(entity)
		}
	} else {
		for _, w := range s.Weapons {
			w.SetOwner(s)
		}
	}
}

// Notes implements WeaponOwner.
func (s *Skill) Notes() string {
	return s.LocalNotes
}

// FeatureList returns the list of Features.
func (s *Skill) FeatureList() feature.Features {
	return s.Features
}

// CategoryList returns the list of categories.
func (s *Skill) CategoryList() []string {
	return s.Categories
}

// Description implements WeaponOwner.
func (s *Skill) Description() string {
	return s.String()
}

func (s *Skill) String() string {
	var buffer strings.Builder
	buffer.WriteString(s.Name)
	if !s.Container() {
		if s.TechLevel != nil {
			buffer.WriteString("/TL")
			buffer.WriteString(*s.TechLevel)
		}
		if s.Specialization != "" {
			buffer.WriteString(" (")
			buffer.WriteString(s.Specialization)
			buffer.WriteByte(')')
		}
	}
	return buffer.String()
}

// AdjustedRelativeLevel returns the relative skill level.
func (s *Skill) AdjustedRelativeLevel() fixed.F64d4 {
	if s.Container() {
		return fixed.F64d4Min
	}
	if s.Entity != nil && s.LevelData.Level > 0 {
		if s.Type == gid.Technique {
			return s.LevelData.RelativeLevel + s.TechniqueDefault.Modifier
		}
		return s.LevelData.RelativeLevel
	}
	// TODO: Old code had a case for templates... but can't see that being exercised in the actual display anywhere
	return fixed.F64d4Min
}

// AdjustedPoints returns the points, adjusted for any bonuses.
func (s *Skill) AdjustedPoints() fixed.F64d4 {
	if s.Container() {
		var total fixed.F64d4
		for _, one := range s.Children {
			total += one.AdjustedPoints()
		}
		return total
	}
	points := s.Points
	if s.Entity != nil && s.Entity.Type == datafile.PC {
		points += s.Entity.SkillPointComparedBonusFor(feature.SkillPointsID+"*", s.Name, s.Specialization, s.Categories, nil)
		points += s.Entity.BonusFor(feature.SkillPointsID+"/"+strings.ToLower(s.Name), nil)
		if points < 0 {
			points = 0
		}
	}
	return points
}

// Level returns the computed level.
func (s *Skill) Level(excludes map[string]bool) fixed.F64d4 {
	return s.CalculateLevel(s.Name, s.Specialization, s.Categories, s.Difficulty, s.Points,
		s.EncumbrancePenaltyMultiplier).Level
}

// CalculateLevel computes the level.
func (s *Skill) CalculateLevel(name, specialization string, categories []string, attrDiff AttributeDifficulty, points, encPenaltyMult fixed.F64d4) skill.Level {
	var tooltip xio.ByteBuffer
	relativeLevel := attrDiff.Difficulty.BaseRelativeLevel()
	level := s.Entity.ResolveAttributeCurrent(attrDiff.Attribute)
	if level != fixed.F64d4Min {
		if attrDiff.Difficulty == skill.Wildcard {
			points = points.Div(fxp.Three).Trunc()
		} else if s.DefaultedFrom != nil && s.DefaultedFrom.Points > 0 {
			points += s.DefaultedFrom.Points
		}
		switch {
		case points == fxp.One:
		// relativeLevel is preset to this point value
		case points > 0 && points < fxp.Four:
			relativeLevel += fxp.One
		case points > 0:
			relativeLevel += fxp.One + points.Div(fxp.Four).Trunc()
		case s.DefaultedFrom != nil && s.DefaultedFrom.Points < 0:
			relativeLevel = s.DefaultedFrom.AdjLevel - level
		default:
			level = fixed.F64d4Min
			relativeLevel = 0
		}
		if level != fixed.F64d4Min {
			level += relativeLevel
			if s.DefaultedFrom != nil && level < s.DefaultedFrom.AdjLevel {
				level = s.DefaultedFrom.AdjLevel
			}
			if s.Entity != nil {
				bonus := s.Entity.SkillComparedBonusFor(feature.SkillNameID+"*", name, specialization, categories, &tooltip)
				level += bonus
				relativeLevel += bonus
				bonus = s.Entity.BonusFor(feature.SkillNameID+"/"+strings.ToLower(name), &tooltip)
				level += bonus
				relativeLevel += bonus
				bonus = s.Entity.EncumbranceLevel(true).Penalty().Mul(encPenaltyMult)
				level += bonus
				if bonus != 0 {
					fmt.Fprintf(&tooltip, i18n.Text("\nEncumbrance [%s]"), bonus.StringWithSign())
				}
			}
		}
	}
	return skill.Level{
		Level:         level,
		RelativeLevel: relativeLevel,
		Tooltip:       tooltip.String(),
	}
}

// UpdateLevel updates the level of the skill, returning true if it has changed.
func (s *Skill) UpdateLevel() bool {
	saved := s.LevelData
	s.DefaultedFrom = s.bestDefaultWithPoints(nil)
	s.LevelData = s.CalculateLevel(s.Name, s.Specialization, s.Categories, s.Difficulty, s.Points,
		s.EncumbrancePenaltyMultiplier)
	return saved != s.LevelData
}

func (s *Skill) bestDefaultWithPoints(excluded *SkillDefault) *SkillDefault {
	best := s.bestDefault(excluded)
	if best != nil {
		baseLine := s.Entity.ResolveAttributeCurrent(s.Difficulty.Attribute) + s.Difficulty.Difficulty.BaseRelativeLevel()
		level := best.Level
		best.AdjLevel = level
		switch {
		case level == baseLine:
			best.Points = fxp.One
		case level == baseLine+fxp.One:
			best.Points = fxp.Two
		case level > baseLine+fxp.One:
			best.Points = fxp.Four.Mul(level - (baseLine + fxp.One))
		default:
			best.Points = -level.Max(0)
		}
	}
	return best
}

func (s *Skill) bestDefault(excluded *SkillDefault) *SkillDefault {
	if s.Entity == nil || len(s.Defaults) == 0 {
		return nil
	}
	excludes := make(map[string]bool)
	excludes[s.String()] = true
	var bestDef *SkillDefault
	best := fixed.F64d4Min
	for _, def := range s.Defaults {
		// For skill-based defaults, prune out any that already use a default that we are involved with
		if def.Equivalent(excluded) || s.inDefaultChain(def, make(map[*Skill]bool)) {
			continue
		}
		level := def.SkillLevel(s.Entity, true, excludes, s.Type != gid.Technique)
		if def.SkillBased() {
			if other := s.Entity.BestSkillNamed(def.Name, def.Specialization, true, excludes); other != nil {
				level -= s.Entity.SkillComparedBonusFor(feature.SkillNameID+"*", def.Name, def.Specialization, s.Categories, nil)
				level -= s.Entity.BonusFor(feature.SkillNameID+"/"+strings.ToLower(def.Name), nil)
			}
		}
		if best < level {
			best = s.LevelData.Level
			bestDef = def.CloneWithoutLevelOrPoints()
			bestDef.Level = level
		}
	}
	return bestDef
}

func (s *Skill) inDefaultChain(def *SkillDefault, lookedAt map[*Skill]bool) bool {
	if s.Entity == nil || def == nil || !def.SkillBased() {
		return false
	}
	hadOne := false
	for _, one := range s.Entity.SkillNamed(def.Name, def.Specialization, true, nil) {
		if one == s {
			return true
		}
		if _, has := lookedAt[one]; !has {
			lookedAt[one] = true
			if s.inDefaultChain(one.DefaultedFrom, lookedAt) {
				return true
			}
		}
		hadOne = true
	}
	return !hadOne
}

// TechniqueSatisfied returns true if the Technique is satisfied.
func (s *Skill) TechniqueSatisfied(tooltip *xio.ByteBuffer, prefix string) bool {
	if s.Type != gid.Technique || !s.TechniqueDefault.SkillBased() {
		return true
	}
	sk := s.Entity.BestSkillNamed(s.TechniqueDefault.Name, s.TechniqueDefault.Specialization, false, nil)
	satisfied := sk != nil && (sk.Type == gid.Technique || sk.Points > 0)
	if !satisfied && tooltip != nil {
		tooltip.WriteString(prefix)
		if sk == nil {
			tooltip.WriteString(i18n.Text("Requires a skill named "))
		} else {
			tooltip.WriteString(i18n.Text("Requires at least 1 point in the skill named "))
		}
		tooltip.WriteString(s.TechniqueDefault.FullName(s.Entity))
	}
	return satisfied
}
