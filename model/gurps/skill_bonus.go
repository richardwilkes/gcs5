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
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/toolbox/log/jot"
)

// SkillBonus holds an adjustment to a skill.
type SkillBonus struct {
	Bonus
	SelectionType          skill.SelectionType `json:"selection_type"`
	NameCriteria           criteria.String     `json:"name,omitempty"`
	SpecializationCriteria criteria.String     `json:"specialization"`
	CategoryCriteria       criteria.String     `json:"category,omitempty"`
}

// NewSkillBonus creates a new SkillBonus.
func NewSkillBonus() *SkillBonus {
	s := &SkillBonus{
		Bonus: Bonus{
			Feature: Feature{
				Type: feature.SkillBonus,
			},
			LeveledAmount: LeveledAmount{Amount: f64d4.One},
		},
		SelectionType: skill.SkillsWithName,
		NameCriteria: criteria.String{
			Compare: criteria.Is,
		},
		SpecializationCriteria: criteria.String{
			Compare: criteria.Any,
		},
		CategoryCriteria: criteria.String{
			Compare: criteria.Any,
		},
	}
	s.Self = s
	return s
}

func (s *SkillBonus) featureMapKey() string {
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

func (s *SkillBonus) fillWithNameableKeys(nameables map[string]string) {
	ExtractNameables(s.SpecializationCriteria.Qualifier, nameables)
	if s.SelectionType != skill.ThisWeapon {
		ExtractNameables(s.NameCriteria.Qualifier, nameables)
		ExtractNameables(s.CategoryCriteria.Qualifier, nameables)
	}
}

func (s *SkillBonus) applyNameableKeys(nameables map[string]string) {
	s.SpecializationCriteria.Qualifier = ApplyNameables(s.SpecializationCriteria.Qualifier, nameables)
	if s.SelectionType != skill.ThisWeapon {
		s.NameCriteria.Qualifier = ApplyNameables(s.NameCriteria.Qualifier, nameables)
		s.CategoryCriteria.Qualifier = ApplyNameables(s.CategoryCriteria.Qualifier, nameables)
	}
}
