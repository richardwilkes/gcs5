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
)

// All is the DR specialization key for DR that affects everything.
const All = "all"

// DRBonusType is the data type key for a DRBonus.
const DRBonusType = "dr_bonus"

const (
	drBonusLocationKey       = "location"
	drBonusSpecializationKey = "specialization"
)

var _ Bonus = &DRBonus{}

// DRBonus holds a bonus to DR.
type DRBonus struct {
	Location       string
	Specialization string
	Amount         *LeveledAmount
}

// NewDRBonus creates a new DRBonus.
func NewDRBonus() *DRBonus {
	return &DRBonus{
		Location:       "torso",
		Specialization: All,
		Amount:         &LeveledAmount{Amount: 1},
	}
}

// NewDRBonusFromJSON creates a new DRBonus from JSON.
func NewDRBonusFromJSON(data map[string]interface{}) *DRBonus {
	b := &DRBonus{
		Location:       encoding.String(data[drBonusLocationKey]),
		Specialization: encoding.String(data[drBonusSpecializationKey]),
		Amount:         NewLeveledAmountFromJSON(data),
	}
	b.Normalize()
	return b
}

// ToJSON implements Feature.
func (b *DRBonus) ToJSON(encoder *encoding.JSONEncoder) {
	b.Normalize()
	encoder.StartObject()
	encoder.KeyedString(drBonusLocationKey, b.Location, true, true)
	encoder.KeyedString(drBonusSpecializationKey, b.Specialization, true, true)
	b.Amount.ToInlineJSON(encoder)
	encoder.EndObject()
}

// CloneFeature implements Feature.
func (b *DRBonus) CloneFeature() Feature {
	amt := *b.Amount
	clone := *b
	clone.Amount = &amt
	return &clone
}

// DataType implements Feature.
func (b *DRBonus) DataType() string {
	return DRBonusType
}

// FeatureKey implements Feature.
func (b *DRBonus) FeatureKey() string {
	return HitLocationPrefix + b.Location
}

// FillWithNameableKeys implements Feature.
func (b *DRBonus) FillWithNameableKeys(_ map[string]string) {
	// Unused
}

// ApplyNameableKeys implements Feature.
func (b *DRBonus) ApplyNameableKeys(_ map[string]string) {
	// Unused
}

// Normalize implements Feature.
func (b *DRBonus) Normalize() {
	s := strings.TrimSpace(b.Specialization)
	if s == "" || strings.EqualFold(s, All) {
		s = All
	}
	b.Specialization = s
}
