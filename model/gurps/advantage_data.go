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

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/gurps/advantage"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

var _ node.EditorData[*Advantage] = &AdvantageData{}

// AdvantageEditData holds the Advantage data that can be edited by the UI detail editor.
type AdvantageEditData struct {
	Name           string                    `json:"name,omitempty"`
	PageRef        string                    `json:"reference,omitempty"`
	LocalNotes     string                    `json:"notes,omitempty"`
	VTTNotes       string                    `json:"vtt_notes,omitempty"`
	Ancestry       string                    `json:"ancestry,omitempty"` // Container only
	UserDesc       string                    `json:"userdesc,omitempty"`
	Categories     []string                  `json:"categories,omitempty"`
	Modifiers      []*AdvantageModifier      `json:"modifiers,omitempty"`
	BasePoints     f64d4.Int                 `json:"base_points,omitempty"`      // Non-container only
	Levels         f64d4.Int                 `json:"levels,omitempty"`           // Non-container only
	PointsPerLevel f64d4.Int                 `json:"points_per_level,omitempty"` // Non-container only
	Prereq         *PrereqList               `json:"prereqs,omitempty"`          // Non-container only
	Weapons        []*Weapon                 `json:"weapons,omitempty"`          // Non-container only
	Features       feature.Features          `json:"features,omitempty"`         // Non-container only
	CR             advantage.SelfControlRoll `json:"cr,omitempty"`
	CRAdj          SelfControlRollAdj        `json:"cr_adj,omitempty"`
	ContainerType  advantage.ContainerType   `json:"container_type,omitempty"` // Container only
	Disabled       bool                      `json:"disabled,omitempty"`
	Mental         bool                      `json:"mental,omitempty"`       // Non-container only
	Physical       bool                      `json:"physical,omitempty"`     // Non-container only
	Social         bool                      `json:"social,omitempty"`       // Non-container only
	Exotic         bool                      `json:"exotic,omitempty"`       // Non-container only
	Supernatural   bool                      `json:"supernatural,omitempty"` // Non-container only
	RoundCostDown  bool                      `json:"round_down,omitempty"`   // Non-container only
}

// CopyFrom implements node.EditorData.
func (d *AdvantageEditData) CopyFrom(adv *Advantage) {
	d.copyFrom(&adv.AdvantageEditData)
}

// ApplyTo implements node.EditorData.
func (d *AdvantageEditData) ApplyTo(adv *Advantage) {
	adv.AdvantageEditData.copyFrom(d)
}

func (d *AdvantageEditData) copyFrom(other *AdvantageEditData) {
	*d = *other
	d.Categories = txt.CloneStringSlice(d.Categories)
	d.Modifiers = nil
	if len(other.Modifiers) != 0 {
		d.Modifiers = make([]*AdvantageModifier, 0, len(other.Modifiers))
		for _, one := range other.Modifiers {
			d.Modifiers = append(d.Modifiers, one.Clone())
		}
	}
	if d.Prereq != nil {
		d.Prereq = d.Prereq.CloneAsPrereqList(nil)
	}
	d.Weapons = nil
	if len(other.Weapons) != 0 {
		d.Weapons = make([]*Weapon, 0, len(other.Weapons))
		for _, one := range other.Weapons {
			d.Weapons = append(d.Weapons, one.Clone())
		}
	}
	d.Features = other.Features.Clone()
}

// AdvantageData holds the Advantage data that is written to disk.
type AdvantageData struct {
	ID   uuid.UUID `json:"id"`
	Type string    `json:"type"`
	AdvantageEditData
	Children []*Advantage `json:"children,omitempty"` // Container only
	IsOpen   bool         `json:"open,omitempty"`     // Container only
}

// UUID returns the UUID of this data.
func (a *AdvantageData) UUID() uuid.UUID {
	return a.ID
}

// Kind returns the kind of data.
func (a *AdvantageData) Kind() string {
	if a.Container() {
		return i18n.Text("Advantage Container")
	}
	return i18n.Text("Advantage")
}

// Container returns true if this is a container.
func (a *AdvantageData) Container() bool {
	return strings.HasSuffix(a.Type, commonContainerKeyPostfix)
}

// Open returns true if this node is currently open.
func (a *AdvantageData) Open() bool {
	return a.IsOpen && a.Container()
}

// SetOpen sets the current open state for this node.
func (a *AdvantageData) SetOpen(open bool) {
	a.IsOpen = open && a.Container()
}

// ClearUnusedFieldsForType zeroes out the fields that are not applicable to this type (container vs not-container).
func (a *AdvantageData) ClearUnusedFieldsForType() {
	if a.Container() {
		a.BasePoints = 0
		a.Levels = 0
		a.PointsPerLevel = 0
		a.Prereq = nil
		a.Weapons = nil
		a.Features = nil
		a.Mental = false
		a.Physical = false
		a.Social = false
		a.Exotic = false
		a.Supernatural = false
		a.RoundCostDown = false
	} else {
		a.ContainerType = 0
		a.IsOpen = false
		a.Ancestry = ""
		a.Children = nil
	}
}

// NodeChildren returns the children of this node, if any.
func (a *AdvantageData) NodeChildren() []node.Node {
	if a.Container() {
		children := make([]node.Node, len(a.Children))
		for i, child := range a.Children {
			children[i] = child
		}
		return children
	}
	return nil
}
