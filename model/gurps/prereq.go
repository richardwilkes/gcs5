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
	"github.com/richardwilkes/gcs/model/criteria"
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/enum"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
)

// Already applied the HasPrereq and NameLevelPrereq portions

const (
	prereqTypeKey  = "type"
	prereqHasKey   = "has"
	prereqNameKey  = "name"
	prereqLevelKey = "level"
)

// Prereq holds data necessary to track a prerequisite.
type Prereq struct {
	Type          enum.PrereqType
	Has           bool
	NameCriteria  criteria.String
	LevelCriteria criteria.Numeric
	Owner         *Prereq // Only those of type PrereqList
}

// NewPrereq creates a new Prereq.
func NewPrereq(prereqType enum.PrereqType) *Prereq {
	p := &Prereq{
		Type: prereqType,
	}
	// TODO: Implement
	switch prereqType {
	case enum.AdvantagePrereq:
		p.Has = true
		p.NameCriteria.Type = enum.Is
		p.LevelCriteria.Type = enum.AtLeast
	case enum.AttributePrereq:
		p.Has = true
	case enum.ContainedQuantityPrereq:
		p.Has = true
	case enum.ContainedWeightPrereq:
		p.Has = true
	case enum.PrereqList:
	case enum.SkillPrereq:
		p.Has = true
		p.NameCriteria.Type = enum.Is
		p.LevelCriteria.Type = enum.AtLeast
	case enum.SpellPrereq:
		p.Has = true
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
	return p
}

// NewPrereqFromJSON creates a new Prereq from JSON.
func NewPrereqFromJSON(data map[string]interface{}) *Prereq {
	p := &Prereq{Type: enum.PrereqTypeFromString(encoding.String(data[prereqTypeKey]))}
	// TODO: Implement
	switch p.Type {
	case enum.AdvantagePrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.NameCriteria.FromJSON(encoding.Object(data[prereqNameKey]))
		p.LevelCriteria.FromJSON(encoding.Object(data[prereqLevelKey]))
	case enum.AttributePrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
	case enum.ContainedQuantityPrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
	case enum.ContainedWeightPrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
	case enum.PrereqList:
	case enum.SkillPrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.NameCriteria.FromJSON(encoding.Object(data[prereqNameKey]))
		p.LevelCriteria.FromJSON(encoding.Object(data[prereqLevelKey]))
	case enum.SpellPrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
	return p
}

// ToJSON emits this Feature as JSON.
func (p *Prereq) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(prereqTypeKey, p.Type.Key(), false, false)
	// TODO: Implement
	switch p.Type {
	case enum.AdvantagePrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoding.ToKeyedJSON(&p.NameCriteria, prereqNameKey, encoder)
		if p.LevelCriteria.Type != enum.AtLeast || p.LevelCriteria.Qualifier != 0 {
			encoding.ToKeyedJSON(&p.LevelCriteria, prereqLevelKey, encoder)
		}
	case enum.AttributePrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
	case enum.ContainedQuantityPrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
	case enum.ContainedWeightPrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
	case enum.PrereqList:
	case enum.SkillPrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoding.ToKeyedJSON(&p.NameCriteria, prereqNameKey, encoder)
		if p.LevelCriteria.Type != enum.AtLeast || p.LevelCriteria.Qualifier != 0 {
			encoding.ToKeyedJSON(&p.LevelCriteria, prereqLevelKey, encoder)
		}
	case enum.SpellPrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
	encoder.EndObject()
}

// Satisfied returns true if this Prereq is satisfied by the specified Entity. 'buffer' will be used, if not nil, to
// write a description of what was unsatisfied. 'prefix' will be appended to each line of the description.
func (p *Prereq) Satisfied(entity *Entity, exclude interface{}, buffer *xio.ByteBuffer, prefix string) bool {
	// TODO: Implement
	switch p.Type {
	case enum.AdvantagePrereq:
	case enum.AttributePrereq:
	case enum.ContainedQuantityPrereq:
	case enum.ContainedWeightPrereq:
	case enum.PrereqList:
	case enum.SkillPrereq:
	case enum.SpellPrereq:
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
	return false
}

// FillWithNameableKeys adds any nameable keys found in this Prereq to the provided map.
func (p *Prereq) FillWithNameableKeys(nameables map[string]string) {
	// TODO: Implement
	switch p.Type {
	case enum.AdvantagePrereq:
		ExtractNameables(p.NameCriteria.Qualifier, nameables)
	case enum.AttributePrereq:
	case enum.ContainedQuantityPrereq:
	case enum.ContainedWeightPrereq:
	case enum.PrereqList:
	case enum.SkillPrereq:
		ExtractNameables(p.NameCriteria.Qualifier, nameables)
	case enum.SpellPrereq:
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
}

// ApplyNameableKeys replaces any nameable keys found in this Prereq with the corresponding values in the provided map.
func (p *Prereq) ApplyNameableKeys(nameables map[string]string) {
	// TODO: Implement
	switch p.Type {
	case enum.AdvantagePrereq:
		p.NameCriteria.Qualifier = ApplyNameables(p.NameCriteria.Qualifier, nameables)
	case enum.AttributePrereq:
	case enum.ContainedQuantityPrereq:
	case enum.ContainedWeightPrereq:
	case enum.PrereqList:
	case enum.SkillPrereq:
		p.NameCriteria.Qualifier = ApplyNameables(p.NameCriteria.Qualifier, nameables)
	case enum.SpellPrereq:
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
}
