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

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	skillDefaultTypeKey           = "type"
	skillDefaultNameKey           = "name"
	skillDefaultSpecializationKey = "specialization"
	skillDefaultModifierKey       = "modifier"
	skillDefaultLevelKey          = "level"
	skillDefaultAdjustedLevelKey  = "adjusted_level"
	skillDefaultPointsKey         = "points"
)

// SkillDefault holds data for a Skill default.
type SkillDefault struct {
	defaultType    string
	Name           string
	Specialization string
	Modifier       fixed.F64d4
	Level          fixed.F64d4
	AdjLevel       fixed.F64d4
	Points         fixed.F64d4
}

// NewSkillDefaultFromJSON creates a new SkillDefault from a JSON object.
func NewSkillDefaultFromJSON(full bool, data map[string]interface{}) *SkillDefault {
	s := &SkillDefault{
		Name:           encoding.String(data[skillDefaultNameKey]),
		Specialization: encoding.String(data[skillDefaultSpecializationKey]),
		Modifier:       encoding.Number(data[skillDefaultModifierKey]),
	}
	s.SetType(encoding.String(data[skillDefaultTypeKey]))
	if full {
		s.Level = encoding.Number(data[skillDefaultLevelKey])
		s.AdjLevel = encoding.Number(data[skillDefaultAdjustedLevelKey])
		s.Points = encoding.Number(data[skillDefaultPointsKey])
	}
	return s
}

// ToJSON emits this object as JSON.
func (s *SkillDefault) ToJSON(full bool, encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(skillDefaultTypeKey, s.defaultType, true, true)
	if skill.DefaultTypeIsSkillBased(s.defaultType) {
		encoder.KeyedString(skillDefaultNameKey, s.Name, true, true)
		encoder.KeyedString(skillDefaultSpecializationKey, s.Specialization, true, true)
	}
	encoder.KeyedNumber(skillDefaultModifierKey, s.Modifier, true)
	if full {
		encoder.KeyedNumber(skillDefaultLevelKey, s.Level, true)
		encoder.KeyedNumber(skillDefaultAdjustedLevelKey, s.AdjLevel, true)
		encoder.KeyedNumber(skillDefaultPointsKey, s.Points, true)
	}
	encoder.EndObject()
}

// Type returns the type of the SkillDefault.
func (s *SkillDefault) Type() string {
	return s.defaultType
}

// SetType sets the type of the SkillDefault.
func (s *SkillDefault) SetType(t string) {
	s.defaultType = id.Sanitize(t, true)
}

// FullName returns the full name of the skill to default from.
func (s *SkillDefault) FullName(entity *Entity) string {
	if skill.DefaultTypeIsSkillBased(s.defaultType) {
		var buffer strings.Builder
		buffer.WriteString(s.Name)
		if s.Specialization != "" {
			buffer.WriteString(" (")
			buffer.WriteString(s.Specialization)
			buffer.WriteByte(')')
		}
		if strings.EqualFold("parry", s.defaultType) {
			buffer.WriteString(i18n.Text(" Parry"))
		} else if strings.EqualFold("block", s.defaultType) {
			buffer.WriteString(i18n.Text(" Block"))
		}
		return buffer.String()
	}
	return ResolveAttributeName(entity, s.defaultType)
}

// FillWithNameableKeys adds any nameable keys found in this SkillDefault to the provided map.
func (s *SkillDefault) FillWithNameableKeys(nameables map[string]string) {
	ExtractNameables(s.Name, nameables)
	ExtractNameables(s.Specialization, nameables)
}

// ApplyNameableKeys replaces any nameable keys found in this SkillDefault with the corresponding values in the provided
// map.
func (s *SkillDefault) ApplyNameableKeys(nameables map[string]string) {
	s.Name = ApplyNameables(s.Name, nameables)
	s.Specialization = ApplyNameables(s.Specialization, nameables)
}

// ModifierAsString returns the modifier as a string suitable for appending.
func (s *SkillDefault) ModifierAsString() string {
	switch {
	case s.Modifier > 0:
		return " + " + s.Modifier.String()
	case s.Modifier < 0:
		return " - " + s.Modifier.String()
	default:
		return ""
	}
}
