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
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/gurps/spell"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
)

const (
	// SpellPointsID holds the ID for spell point lookups.
	SpellPointsID = "spell.points"
	// SpellCollegePointsID holds the ID for spell college point lookups.
	SpellCollegePointsID = "spell.college.points"
	// SpellPowerSourcePointsID holds the ID for spell power source point lookups.
	SpellPowerSourcePointsID = "spell.power_source.points"
)

var _ Bonus = &SpellPointBonus{}

// SpellPointBonus holds an adjustment to a spell's points.
type SpellPointBonus struct {
	Type             Type            `json:"type"`
	Parent           fmt.Stringer    `json:"-"`
	SpellMatchType   spell.MatchType `json:"match,omitempty"`
	NameCriteria     criteria.String `json:"name,omitempty"`
	CategoryCriteria criteria.String `json:"category,omitempty"`
	LeveledAmount
}

// NewSpellPointBonus creates a new SpellPointBonus.
func NewSpellPointBonus() *SpellPointBonus {
	return &SpellPointBonus{
		Type:           SpellPointBonusType,
		SpellMatchType: spell.AllColleges,
		NameCriteria: criteria.String{
			Compare: criteria.Is,
		},
		CategoryCriteria: criteria.String{
			Compare: criteria.Any,
		},
		LeveledAmount: LeveledAmount{Amount: f64d4.One},
	}
}

// FeatureMapKey implements Feature.
func (s *SpellPointBonus) FeatureMapKey() string {
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

// FillWithNameableKeys implements Feature.
func (s *SpellPointBonus) FillWithNameableKeys(m map[string]string) {
	if s.SpellMatchType != spell.AllColleges {
		nameables.Extract(s.NameCriteria.Qualifier, m)
	}
	nameables.Extract(s.CategoryCriteria.Qualifier, m)
}

// ApplyNameableKeys implements Feature.
func (s *SpellPointBonus) ApplyNameableKeys(m map[string]string) {
	if s.SpellMatchType != spell.AllColleges {
		s.NameCriteria.Qualifier = nameables.Apply(s.NameCriteria.Qualifier, m)
	}
	s.CategoryCriteria.Qualifier = nameables.Apply(s.CategoryCriteria.Qualifier, m)
}

// AddToTooltip implements Bonus.
func (s *SpellPointBonus) AddToTooltip(buffer *xio.ByteBuffer) {
	basicAddToTooltip(s.Parent, &s.LeveledAmount, buffer)
}
