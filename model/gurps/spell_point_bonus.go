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

// SpellPointBonus holds an adjustment to a spell's points.
type SpellPointBonus struct {
	Bonus
	SpellMatchType   spell.MatchType `json:"match,omitempty"`
	NameCriteria     criteria.String `json:"name,omitempty"`
	CategoryCriteria criteria.String `json:"category,omitempty"`
}

// NewSpellPointBonus creates a new SpellPointBonus.
func NewSpellPointBonus() *SpellPointBonus {
	s := &SpellPointBonus{
		Bonus: Bonus{
			Feature: Feature{
				Type: feature.SpellPointBonus,
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

func (s *SpellPointBonus) featureMapKey() string {
	if s.CategoryCriteria.Compare != criteria.Any {
		return SpellPointsID + "*"
	}
	switch s.SpellMatchType {
	case spell.AllColleges:
		return SpellCollegePointsID
	case spell.CollegeName:
		return s.buildKey(SpellCollegePointsID)
	case spell.PowerSource:
		return s.buildKey(SpellPowerSourcePointsID)
	case spell.Spell:
		return s.buildKey(SpellPointsID)
	default:
		jot.Fatal(1, "invalid match type: ", s.SpellMatchType)
		return ""
	}
}

func (s *SpellPointBonus) buildKey(prefix string) string {
	if s.NameCriteria.Compare == criteria.Is {
		return prefix + "/" + s.NameCriteria.Qualifier
	}
	return prefix + "*"
}

func (s *SpellPointBonus) fillWithNameableKeys(nameables map[string]string) {
	if s.SpellMatchType != spell.AllColleges {
		ExtractNameables(s.NameCriteria.Qualifier, nameables)
	}
	ExtractNameables(s.CategoryCriteria.Qualifier, nameables)
}

func (s *SpellPointBonus) applyNameableKeys(nameables map[string]string) {
	if s.SpellMatchType != spell.AllColleges {
		s.NameCriteria.Qualifier = ApplyNameables(s.NameCriteria.Qualifier, nameables)
	}
	s.CategoryCriteria.Qualifier = ApplyNameables(s.CategoryCriteria.Qualifier, nameables)
}
