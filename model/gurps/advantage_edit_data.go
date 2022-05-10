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
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/advantage"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/toolbox/txt"
)

var _ node.EditorData[*Advantage] = &AdvantageEditData{}

// AdvantageEditData holds the Advantage data that can be edited by the UI detail editor.
type AdvantageEditData struct {
	Name           string                    `json:"name,omitempty"`
	PageRef        string                    `json:"reference,omitempty"`
	LocalNotes     string                    `json:"notes,omitempty"`
	VTTNotes       string                    `json:"vtt_notes,omitempty"`
	Ancestry       string                    `json:"ancestry,omitempty"` // Container only
	UserDesc       string                    `json:"userdesc,omitempty"`
	Tags           []string                  `json:"tags,omitempty"`
	Modifiers      []*AdvantageModifier      `json:"modifiers,omitempty"`
	BasePoints     fxp.Int                   `json:"base_points,omitempty"`      // Non-container only
	Levels         fxp.Int                   `json:"levels,omitempty"`           // Non-container only
	PointsPerLevel fxp.Int                   `json:"points_per_level,omitempty"` // Non-container only
	Prereq         *PrereqList               `json:"prereqs,omitempty"`          // Non-container only
	Weapons        []*Weapon                 `json:"weapons,omitempty"`          // Non-container only
	Features       feature.Features          `json:"features,omitempty"`         // Non-container only
	CR             advantage.SelfControlRoll `json:"cr,omitempty"`
	CRAdj          SelfControlRollAdj        `json:"cr_adj,omitempty"`
	ContainerType  advantage.ContainerType   `json:"container_type,omitempty"` // Container only
	Disabled       bool                      `json:"disabled,omitempty"`
	RoundCostDown  bool                      `json:"round_down,omitempty"` // Non-container only
}

// CopyFrom implements node.EditorData.
func (d *AdvantageEditData) CopyFrom(adv *Advantage) {
	d.copyFrom(&adv.AdvantageEditData, adv.Container(), false)
}

// ApplyTo implements node.EditorData.
func (d *AdvantageEditData) ApplyTo(adv *Advantage) {
	adv.AdvantageEditData.copyFrom(d, adv.Container(), true)
}

func (d *AdvantageEditData) copyFrom(other *AdvantageEditData, isContainer, isApply bool) {
	*d = *other
	d.Tags = txt.CloneStringSlice(d.Tags)
	d.Modifiers = nil
	if len(other.Modifiers) != 0 {
		d.Modifiers = make([]*AdvantageModifier, 0, len(other.Modifiers))
		for _, one := range other.Modifiers {
			d.Modifiers = append(d.Modifiers, one.Clone())
		}
	}
	d.Prereq = d.Prereq.CloneResolvingEmpty(isContainer, isApply)
	d.Weapons = nil
	if len(other.Weapons) != 0 {
		d.Weapons = make([]*Weapon, 0, len(other.Weapons))
		for _, one := range other.Weapons {
			d.Weapons = append(d.Weapons, one.Clone())
		}
	}
	d.Features = other.Features.Clone()
}
