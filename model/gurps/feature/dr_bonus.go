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
	"strings"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
)

// HitLocationPrefix is the prefix used on all hit locations for DR bonuses.
const HitLocationPrefix = "hit_location."

var _ Bonus = &DRBonus{}

// DRBonusData is split out so that it can be adjusted before and after being serialized.
type DRBonusData struct {
	Type           Type   `json:"type"`
	Location       string `json:"location"`
	Specialization string `json:"specialization,omitempty"`
	LeveledAmount
}

// DRBonus holds the data for a DR adjustment.
type DRBonus struct {
	DRBonusData
	Parent fmt.Stringer
}

// NewDRBonus creates a new DRBonus.
func NewDRBonus() *DRBonus {
	return &DRBonus{
		DRBonusData: DRBonusData{
			Type:           DRBonusType,
			Location:       "torso",
			Specialization: gid.All,
			LeveledAmount:  LeveledAmount{Amount: fxp.One},
		},
	}
}

// FeatureType implements Feature.
func (d *DRBonus) FeatureType() Type {
	return d.Type
}

// Clone implements Feature.
func (d *DRBonus) Clone() Feature {
	other := *d
	return &other
}

// Normalize adjusts the data to it preferred representation.
func (d *DRBonus) Normalize() {
	s := strings.TrimSpace(d.Specialization)
	if s == "" || strings.EqualFold(s, gid.All) {
		s = gid.All
	}
	d.Specialization = s
}

// FeatureMapKey implements Feature.
func (d *DRBonus) FeatureMapKey() string {
	return HitLocationPrefix + d.Location
}

// FillWithNameableKeys implements Feature.
func (d *DRBonus) FillWithNameableKeys(_ map[string]string) {
}

// ApplyNameableKeys implements Feature.
func (d *DRBonus) ApplyNameableKeys(_ map[string]string) {
}

// SetParent implements Bonus.
func (d *DRBonus) SetParent(parent fmt.Stringer) {
	d.Parent = parent
}

// SetLevel implements Bonus.
func (d *DRBonus) SetLevel(level fxp.Int) {
	d.Level = level
}

// AddToTooltip implements Bonus.
func (d *DRBonus) AddToTooltip(buffer *xio.ByteBuffer) {
	if buffer != nil {
		d.Normalize()
		buffer.WriteByte('\n')
		buffer.WriteString(parentName(d.Parent))
		buffer.WriteString(" [")
		buffer.WriteString(d.LeveledAmount.FormatWithLevel())
		buffer.WriteString(i18n.Text(" against "))
		buffer.WriteString(d.Specialization)
		buffer.WriteString(i18n.Text(" attacks]"))
	}
}

// MarshalJSON implements json.Marshaler.
func (d *DRBonus) MarshalJSON() ([]byte, error) {
	d.Normalize()
	if d.Specialization == gid.All {
		d.Specialization = ""
	}
	data, err := json.Marshal(&d.DRBonusData)
	d.Normalize()
	return data, err
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DRBonus) UnmarshalJSON(data []byte) error {
	d.DRBonusData = DRBonusData{}
	if err := json.Unmarshal(data, &d.DRBonusData); err != nil {
		return err
	}
	d.Normalize()
	return nil
}
