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

package feature

import (
	"fmt"

	"github.com/richardwilkes/gcs/model/criteria"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	// SkillNameID holds the ID for skill name lookups.
	SkillNameID = "skill.name"
)

var _ Bonus = &SkillBonus{}

// SkillBonus holds an adjustment to a skill.
type SkillBonus struct {
	Parent                 fmt.Stringer        `json:"-"`
	Type                   Type                `json:"type"`
	SelectionType          skill.SelectionType `json:"selection_type"`
	NameCriteria           criteria.String     `json:"name,omitempty"`
	SpecializationCriteria criteria.String     `json:"specialization,omitempty"`
	CategoryCriteria       criteria.String     `json:"category,omitempty"`
	LeveledAmount
}

// NewSkillBonus creates a new SkillBonus.
func NewSkillBonus() *SkillBonus {
	return &SkillBonus{
		Type:          SkillBonusType,
		SelectionType: skill.SkillsWithName,
		NameCriteria: criteria.String{
			StringData: criteria.StringData{
				Compare: criteria.Is,
			},
		},
		SpecializationCriteria: criteria.String{
			StringData: criteria.StringData{
				Compare: criteria.Any,
			},
		},
		CategoryCriteria: criteria.String{
			StringData: criteria.StringData{
				Compare: criteria.Any,
			},
		},
		LeveledAmount: LeveledAmount{Amount: fxp.One},
	}
}

// FeatureMapKey implements Feature.
func (s *SkillBonus) FeatureMapKey() string {
	switch s.SelectionType {
	case skill.SkillsWithName:
		return s.buildKey(SkillNameID)
	case skill.ThisWeapon:
		return ThisWeaponID
	case skill.WeaponsWithName:
		return s.buildKey(WeaponNamedIDPrefix)
	default:
		jot.Fatal(1, "invalid selection type: ", s.SelectionType)
		return ""
	}
}

func (s *SkillBonus) buildKey(prefix string) string {
	if s.NameCriteria.Compare == criteria.Is &&
		(s.SpecializationCriteria.Compare == criteria.Any && s.CategoryCriteria.Compare == criteria.Any) {
		return prefix + "/" + s.NameCriteria.Qualifier
	}
	return prefix + "*"
}

// FillWithNameableKeys implements Feature.
func (s *SkillBonus) FillWithNameableKeys(m map[string]string) {
	nameables.Extract(s.SpecializationCriteria.Qualifier, m)
	if s.SelectionType != skill.ThisWeapon {
		nameables.Extract(s.NameCriteria.Qualifier, m)
		nameables.Extract(s.CategoryCriteria.Qualifier, m)
	}
}

// ApplyNameableKeys implements Feature.
func (s *SkillBonus) ApplyNameableKeys(m map[string]string) {
	s.SpecializationCriteria.Qualifier = nameables.Apply(s.SpecializationCriteria.Qualifier, m)
	if s.SelectionType != skill.ThisWeapon {
		s.NameCriteria.Qualifier = nameables.Apply(s.NameCriteria.Qualifier, m)
		s.CategoryCriteria.Qualifier = nameables.Apply(s.CategoryCriteria.Qualifier, m)
	}
}

// SetParent implements Bonus.
func (s *SkillBonus) SetParent(parent fmt.Stringer) {
	s.Parent = parent
}

// SetLevel implements Bonus.
func (s *SkillBonus) SetLevel(level fixed.F64d4) {
	s.Level = level
}

// AddToTooltip implements Bonus.
func (s *SkillBonus) AddToTooltip(buffer *xio.ByteBuffer) {
	basicAddToTooltip(s.Parent, &s.LeveledAmount, buffer)
}
