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
	"encoding/json"
	"strings"

	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
)

// All is the DR specialization key for DR that affects everything.
const All = "all"

// DRBonusData is split out so that it can be adjusted before and after being serialized.
type DRBonusData struct {
	Bonus
	Location       string `json:"location"`
	Specialization string `json:"specialization,omitempty"`
}

// DRBonus holds the data for a DR adjustment.
type DRBonus struct {
	DRBonusData
}

// NewDRBonus creates a new DRBonus.
func NewDRBonus() *DRBonus {
	d := &DRBonus{
		DRBonusData: DRBonusData{
			Bonus: Bonus{
				Feature: Feature{
					Type: feature.DRBonus,
				},
				LeveledAmount: LeveledAmount{Amount: f64d4.One},
			},
			Location:       "torso",
			Specialization: All,
		},
	}
	d.Self = d
	return d
}

// Normalize adjusts the data to it preferred representation.
func (d *DRBonus) Normalize() {
	s := strings.TrimSpace(d.Specialization)
	if s == "" || strings.EqualFold(s, All) {
		s = All
	}
	d.Specialization = s
}

func (d *DRBonus) featureMapKey() string {
	return HitLocationPrefix + d.Location
}

func (d *DRBonus) addToTooltip(buffer *xio.ByteBuffer) {
	d.Normalize()
	buffer.WriteByte('\n')
	buffer.WriteString(d.ParentName())
	buffer.WriteString(" [")
	buffer.WriteString(d.LeveledAmount.Format(i18n.Text("level")))
	buffer.WriteString(i18n.Text(" against "))
	buffer.WriteString(d.Specialization)
	buffer.WriteString(i18n.Text(" attacks]"))
}

// MarshalJSON implements json.Marshaler.
func (d *DRBonus) MarshalJSON() ([]byte, error) {
	d.Normalize()
	if d.Specialization == All {
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
