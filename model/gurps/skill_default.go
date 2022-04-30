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

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

var skillBasedDefaultTypes = map[string]bool{
	gid.Skill: true,
	gid.Parry: true,
	gid.Block: true,
}

// SkillDefault holds data for a Skill default.
type SkillDefault struct {
	DefaultType    string    `json:"type"`
	Name           string    `json:"name,omitempty"`
	Specialization string    `json:"specialization,omitempty"`
	Modifier       f64d4.Int `json:"modifier,omitempty"`
	Level          f64d4.Int `json:"level,omitempty"`
	AdjLevel       f64d4.Int `json:"adjusted_level,omitempty"`
	Points         f64d4.Int `json:"points,omitempty"`
}

// CloneWithoutLevelOrPoints creates a copy, but without the level or points set.
func (s *SkillDefault) CloneWithoutLevelOrPoints() *SkillDefault {
	clone := *s
	clone.Level = 0
	clone.AdjLevel = 0
	clone.Points = 0
	return &clone
}

// Equivalent returns true if this can be considered equivalent to other.
func (s *SkillDefault) Equivalent(other *SkillDefault) bool {
	return other != nil && s.DefaultType == other.DefaultType && s.Modifier == other.Modifier && s.Name == other.Name &&
		s.Specialization == other.Specialization
}

// Type returns the type of the SkillDefault.
func (s *SkillDefault) Type() string {
	return s.DefaultType
}

// SetType sets the type of the SkillDefault.
func (s *SkillDefault) SetType(t string) {
	s.DefaultType = id.Sanitize(t, true)
}

// FullName returns the full name of the skill to default from.
func (s *SkillDefault) FullName(entity *Entity) string {
	if s.SkillBased() {
		var buffer strings.Builder
		buffer.WriteString(s.Name)
		if s.Specialization != "" {
			buffer.WriteString(" (")
			buffer.WriteString(s.Specialization)
			buffer.WriteByte(')')
		}
		switch {
		case strings.EqualFold(gid.Dodge, s.DefaultType):
			buffer.WriteString(i18n.Text(" Dodge"))
		case strings.EqualFold(gid.Parry, s.DefaultType):
			buffer.WriteString(i18n.Text(" Parry"))
		case strings.EqualFold(gid.Block, s.DefaultType):
			buffer.WriteString(i18n.Text(" Block"))
		}
		return buffer.String()
	}
	return ResolveAttributeName(entity, s.DefaultType)
}

// FillWithNameableKeys adds any nameable keys found in this SkillDefault to the provided map.
func (s *SkillDefault) FillWithNameableKeys(m map[string]string) {
	nameables.Extract(s.Name, m)
	nameables.Extract(s.Specialization, m)
}

// ApplyNameableKeys replaces any nameable keys found in this SkillDefault with the corresponding values in the provided
// map.
func (s *SkillDefault) ApplyNameableKeys(m map[string]string) {
	s.Name = nameables.Apply(s.Name, m)
	s.Specialization = nameables.Apply(s.Specialization, m)
}

// ModifierAsString returns the modifier as a string suitable for appending.
func (s *SkillDefault) ModifierAsString() string {
	if s.Modifier != 0 {
		return s.Modifier.StringWithSign()
	}
	return ""
}

// SkillBased returns true if the Type() is Skill-based.
func (s *SkillDefault) SkillBased() bool {
	return skillBasedDefaultTypes[strings.ToLower(strings.TrimSpace(s.DefaultType))]
}

// SkillLevel returns the base skill level for this SkillDefault.
func (s *SkillDefault) SkillLevel(entity *Entity, requirePoints bool, excludes map[string]bool, ruleOf20 bool) f64d4.Int {
	switch s.Type() {
	case gid.Parry:
		best := s.best(entity, requirePoints, excludes)
		if best != f64d4.Min {
			best = best.Div(fxp.Two).Trunc() + fxp.Three + entity.ParryBonus
		}
		return s.finalLevel(best)
	case gid.Block:
		best := s.best(entity, requirePoints, excludes)
		if best != f64d4.Min {
			best = best.Div(fxp.Two).Trunc() + fxp.Three + entity.BlockBonus
		}
		return s.finalLevel(best)
	case gid.Skill:
		return s.finalLevel(s.best(entity, requirePoints, excludes))
	default:
		return s.SkillLevelFast(entity, requirePoints, excludes, ruleOf20)
	}
}

func (s *SkillDefault) best(entity *Entity, requirePoints bool, excludes map[string]bool) f64d4.Int {
	best := f64d4.Min
	for _, sk := range entity.SkillNamed(s.Name, s.Specialization, requirePoints, excludes) {
		if best < sk.LevelData.Level {
			level := sk.CalculateLevel().Level
			if best < level {
				best = level
			}
		}
	}
	return best
}

// SkillLevelFast returns the base skill level for this SkillDefault.
func (s *SkillDefault) SkillLevelFast(entity *Entity, requirePoints bool, excludes map[string]bool, ruleOf20 bool) f64d4.Int {
	switch s.Type() {
	case gid.Dodge:
		level := entity.Dodge(entity.EncumbranceLevel(true))
		if ruleOf20 && level > 20 {
			level = 20
		}
		return s.finalLevel(f64d4.FromInt(level))
	case gid.Parry:
		best := s.bestFast(entity, requirePoints, excludes)
		if best != f64d4.Min {
			best = best.Div(fxp.Two).Trunc() + fxp.Three + entity.ParryBonus
		}
		return s.finalLevel(best)
	case gid.Block:
		best := s.bestFast(entity, requirePoints, excludes)
		if best != f64d4.Min {
			best = best.Div(fxp.Two).Trunc() + fxp.Three + entity.BlockBonus
		}
		return s.finalLevel(best)
	case gid.Skill:
		return s.finalLevel(s.bestFast(entity, requirePoints, excludes))
	default:
		level := entity.ResolveAttributeCurrent(s.Type())
		if ruleOf20 {
			level = level.Min(fxp.Twenty)
		}
		return s.finalLevel(level)
	}
}

func (s *SkillDefault) bestFast(entity *Entity, requirePoints bool, excludes map[string]bool) f64d4.Int {
	best := f64d4.Min
	for _, sk := range entity.SkillNamed(s.Name, s.Specialization, requirePoints, excludes) {
		if best < sk.LevelData.Level {
			best = sk.LevelData.Level
		}
	}
	return best
}

func (s *SkillDefault) finalLevel(level f64d4.Int) f64d4.Int {
	if level != f64d4.Min {
		level += s.Modifier
	}
	return level
}
