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
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
)

// SkillPointBonus holds an adjustment to a skill's points.
type SkillPointBonus struct {
	Bonus
	NameCriteria           criteria.String `json:"name"`
	SpecializationCriteria criteria.String `json:"specialization"`
	CategoryCriteria       criteria.String `json:"category"`
}

// NewSkillPointBonus creates a new SkillPointBonus.
func NewSkillPointBonus() *SkillPointBonus {
	s := &SkillPointBonus{
		Bonus: Bonus{
			Feature: Feature{
				Type: feature.SkillPointBonus,
			},
			LeveledAmount: LeveledAmount{Amount: f64d4.One},
		},
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

func (s *SkillPointBonus) featureMapKey() string {
	if s.NameCriteria.Compare == criteria.Is &&
		(s.SpecializationCriteria.Compare == criteria.Any && s.CategoryCriteria.Compare == criteria.Any) {
		return SkillPointsID + "/" + s.NameCriteria.Qualifier
	}
	return SkillPointsID + "*"
}

func (s *SkillPointBonus) fillWithNameableKeys(nameables map[string]string) {
	ExtractNameables(s.NameCriteria.Qualifier, nameables)
	ExtractNameables(s.SpecializationCriteria.Qualifier, nameables)
	ExtractNameables(s.CategoryCriteria.Qualifier, nameables)
}

func (s *SkillPointBonus) applyNameableKeys(nameables map[string]string) {
	s.NameCriteria.Qualifier = ApplyNameables(s.NameCriteria.Qualifier, nameables)
	s.SpecializationCriteria.Qualifier = ApplyNameables(s.SpecializationCriteria.Qualifier, nameables)
	s.CategoryCriteria.Qualifier = ApplyNameables(s.CategoryCriteria.Qualifier, nameables)
}

func (s *SkillPointBonus) addToTooltip(buffer *xio.ByteBuffer) {
	if buffer != nil {
		buffer.WriteByte('\n')
		buffer.WriteString(s.ParentName())
		buffer.WriteString(" [")
		buffer.WriteString(s.LeveledAmount.Format(i18n.Text("level")))
		if s.AdjustedAmount() == f64d4.One {
			buffer.WriteString(i18n.Text(" pt]"))
		} else {
			buffer.WriteString(i18n.Text(" pts]"))
		}
	}
}
