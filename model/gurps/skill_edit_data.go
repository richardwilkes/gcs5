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
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

var _ node.EditorData[*Skill] = &SkillEditData{}

// SkillEditData holds the Skill data that can be edited by the UI detail editor.
type SkillEditData struct {
	Name                         string              `json:"name,omitempty"`
	PageRef                      string              `json:"reference,omitempty"`
	LocalNotes                   string              `json:"notes,omitempty"`
	VTTNotes                     string              `json:"vtt_notes,omitempty"`
	Tags                         []string            `json:"tags,omitempty"`
	Specialization               string              `json:"specialization,omitempty"`                 // Non-container only
	TechLevel                    *string             `json:"tech_level,omitempty"`                     // Non-container only
	Difficulty                   AttributeDifficulty `json:"difficulty"`                               // Non-container only
	Points                       f64d4.Int           `json:"points,omitempty"`                         // Non-container only
	EncumbrancePenaltyMultiplier f64d4.Int           `json:"encumbrance_penalty_multiplier,omitempty"` // Non-container only
	DefaultedFrom                *SkillDefault       `json:"defaulted_from,omitempty"`                 // Non-container only
	Defaults                     []*SkillDefault     `json:"defaults,omitempty"`                       // Non-container only
	TechniqueDefault             *SkillDefault       `json:"default,omitempty"`                        // Non-container only
	TechniqueLimitModifier       *f64d4.Int          `json:"limit,omitempty"`                          // Non-container only
	Prereq                       *PrereqList         `json:"prereqs,omitempty"`                        // Non-container only
	Weapons                      []*Weapon           `json:"weapons,omitempty"`                        // Non-container only
	Features                     feature.Features    `json:"features,omitempty"`                       // Non-container only
}

// CopyFrom implements node.EditorData.
func (d *SkillEditData) CopyFrom(s *Skill) {
	d.copyFrom(&s.SkillEditData)
}

// ApplyTo implements node.EditorData.
func (d *SkillEditData) ApplyTo(s *Skill) {
	s.SkillEditData.copyFrom(d)
}

func (d *SkillEditData) copyFrom(other *SkillEditData) {
	*d = *other
	d.Tags = txt.CloneStringSlice(d.Tags)
	if other.TechLevel != nil {
		tl := *other.TechLevel
		d.TechLevel = &tl
	}
	if other.DefaultedFrom != nil {
		def := *other.DefaultedFrom
		d.DefaultedFrom = &def
	}
	d.Defaults = nil
	if len(other.Defaults) != 0 {
		d.Defaults = make([]*SkillDefault, len(other.Defaults))
		for i, def := range other.Defaults {
			def2 := *def
			d.Defaults[i] = &def2
		}
	}
	if other.TechniqueDefault != nil {
		def := *other.TechniqueDefault
		d.TechniqueDefault = &def
	}
	if other.TechniqueLimitModifier != nil {
		mod := *other.TechniqueLimitModifier
		d.TechniqueLimitModifier = &mod
	}
	if d.Prereq != nil {
		d.Prereq = d.Prereq.CloneAsPrereqList(nil)
	}
	d.Weapons = nil
	if len(other.Weapons) != 0 {
		d.Weapons = make([]*Weapon, 0, len(other.Weapons))
		for _, one := range other.Weapons {
			d.Weapons = append(d.Weapons, one.Clone())
		}
	}
	d.Features = other.Features.Clone()
}
