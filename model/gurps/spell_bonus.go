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
	"github.com/richardwilkes/gcs/model/gurps/spell"
	"github.com/richardwilkes/toolbox/log/jot"
)

// SpellBonus holds the data for a bonus to a spell.
type SpellBonus struct {
	Bonus
	SpellMatchType   spell.MatchType `json:"match,omitempty"`
	NameCriteria     criteria.String `json:"name,omitempty"`
	CategoryCriteria criteria.String `json:"category,omitempty"`
}

// NewSpellBonus creates a new SpellBonus.
func NewSpellBonus() *SpellBonus {
	s := &SpellBonus{
		Bonus: Bonus{
			Feature: Feature{
				Type: feature.SpellBonus,
			},
			LeveledAmount: LeveledAmount{Amount: f64d4.One},
		},
		SpellMatchType: spell.AllColleges,
		NameCriteria: criteria.String{
			Compare: criteria.Is,
		},
		CategoryCriteria: criteria.String{
			Compare: criteria.Any,
		},
	}
	s.Self = s
	return s
}

func (s *SpellBonus) featureMapKey() string {
	if s.CategoryCriteria.Compare != criteria.Any {
		return SpellNameID + "*"
	}
	switch s.SpellMatchType {
	case spell.AllColleges:
		return SpellCollegeID
	case spell.CollegeName:
		return s.buildKey(SpellCollegeID)
	case spell.PowerSource:
		return s.buildKey(SpellPowerSourceID)
	case spell.Spell:
		return s.buildKey(SpellNameID)
	default:
		jot.Fatal(1, "invalid match type: ", s.SpellMatchType)
		return ""
	}
}

func (s *SpellBonus) buildKey(prefix string) string {
	if s.NameCriteria.Compare == criteria.Is {
		return prefix + "/" + s.NameCriteria.Qualifier
	}
	return prefix + "*"
}

func (s *SpellBonus) fillWithNameableKeys(nameables map[string]string) {
	if s.SpellMatchType != spell.AllColleges {
		ExtractNameables(s.NameCriteria.Qualifier, nameables)
	}
	ExtractNameables(s.CategoryCriteria.Qualifier, nameables)
}

func (s *SpellBonus) applyNameableKeys(nameables map[string]string) {
	if s.SpellMatchType != spell.AllColleges {
		s.NameCriteria.Qualifier = ApplyNameables(s.NameCriteria.Qualifier, nameables)
	}
	s.CategoryCriteria.Qualifier = ApplyNameables(s.CategoryCriteria.Qualifier, nameables)
}
