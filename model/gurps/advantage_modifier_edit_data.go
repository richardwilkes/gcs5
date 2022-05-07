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

var _ node.EditorData[*AdvantageModifier] = &AdvantageModifierEditData{}

// AdvantageModifierEditData holds the AdvantageModifier data that can be edited by the UI detail editor.
type AdvantageModifierEditData struct {
	Name       string                     `json:"name,omitempty"`
	PageRef    string                     `json:"reference,omitempty"`
	LocalNotes string                     `json:"notes,omitempty"`
	VTTNotes   string                     `json:"vtt_notes,omitempty"`
	Tags       []string                   `json:"tags,omitempty"`
	Cost       fxp.Int                    `json:"cost,omitempty"`      // Non-container only
	Levels     fxp.Int                    `json:"levels,omitempty"`    // Non-container only
	Affects    advantage.Affects          `json:"affects,omitempty"`   // Non-container only
	CostType   advantage.ModifierCostType `json:"cost_type,omitempty"` // Non-container only
	Disabled   bool                       `json:"disabled,omitempty"`  // Non-container only
	Features   feature.Features           `json:"features,omitempty"`  // Non-container only
}

// CopyFrom implements node.EditorData.
func (d *AdvantageModifierEditData) CopyFrom(adv *AdvantageModifier) {
	d.copyFrom(&adv.AdvantageModifierEditData)
}

// ApplyTo implements node.EditorData.
func (d *AdvantageModifierEditData) ApplyTo(adv *AdvantageModifier) {
	adv.AdvantageModifierEditData.copyFrom(d)
}

func (d *AdvantageModifierEditData) copyFrom(other *AdvantageModifierEditData) {
	*d = *other
	d.Tags = txt.CloneStringSlice(d.Tags)
	d.Features = other.Features.Clone()
}
