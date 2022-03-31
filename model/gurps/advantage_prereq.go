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
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

var _ Prereq = &AdvantagePrereq{}

// AdvantagePrereq holds a prereq against an Advantage.
type AdvantagePrereq struct {
	Parent        *PrereqList      `json:"-"`
	Type          prereq.Type      `json:"type"`
	Has           bool             `json:"has"`
	NameCriteria  criteria.String  `json:"name,omitempty"`
	LevelCriteria criteria.Numeric `json:"level,omitempty"`
	NotesCriteria criteria.String  `json:"notes,omitempty"`
}

// NewAdvantagePrereq creates a new AdvantagePrereq.
func NewAdvantagePrereq() *AdvantagePrereq {
	return &AdvantagePrereq{
		Type: prereq.Advantage,
		NameCriteria: criteria.String{
			StringData: criteria.StringData{
				Compare: criteria.Is,
			},
		},
		LevelCriteria: criteria.Numeric{
			NumericData: criteria.NumericData{
				Compare: criteria.AtLeast,
			},
		},
		NotesCriteria: criteria.String{
			StringData: criteria.StringData{
				Compare: criteria.Any,
			},
		},
		Has: true,
	}
}

// Clone implements Prereq.
func (a *AdvantagePrereq) Clone(parent *PrereqList) Prereq {
	clone := *a
	clone.Parent = parent
	return &clone
}

// FillWithNameableKeys implements Prereq.
func (a *AdvantagePrereq) FillWithNameableKeys(m map[string]string) {
	nameables.Extract(a.NameCriteria.Qualifier, m)
	nameables.Extract(a.NotesCriteria.Qualifier, m)
}

// ApplyNameableKeys implements Prereq.
func (a *AdvantagePrereq) ApplyNameableKeys(m map[string]string) {
	a.NameCriteria.Qualifier = nameables.Apply(a.NameCriteria.Qualifier, m)
	a.NotesCriteria.Qualifier = nameables.Apply(a.NotesCriteria.Qualifier, m)
}

// Satisfied implements Prereq.
func (a *AdvantagePrereq) Satisfied(entity *Entity, exclude interface{}, tooltip *xio.ByteBuffer, prefix string) bool {
	satisfied := false
	TraverseAdvantages(func(adq *Advantage) bool {
		if exclude == adq || !a.NameCriteria.Matches(adq.Name) {
			return false
		}
		notes := adq.Notes()
		if modNotes := adq.ModifierNotes(); modNotes != "" {
			notes += "\n" + modNotes
		}
		if !a.NotesCriteria.Matches(notes) {
			return false
		}
		var levels f64d4.Int
		if adq.Levels != nil {
			levels = adq.Levels.Max(0)
		}
		satisfied = a.LevelCriteria.Matches(levels)
		return satisfied
	}, true, entity.Advantages...)
	if !a.Has {
		satisfied = !satisfied
	}
	if !satisfied && tooltip != nil {
		tooltip.WriteByte('\n')
		tooltip.WriteString(prefix)
		tooltip.WriteString(HasText(a.Has))
		tooltip.WriteString(i18n.Text(" an advantage whose name "))
		tooltip.WriteString(a.NameCriteria.String())
		if a.NotesCriteria.Compare != criteria.Any {
			tooltip.WriteString(i18n.Text(", notes "))
			tooltip.WriteString(a.NotesCriteria.String())
			tooltip.WriteByte(',')
		}
		tooltip.WriteString(i18n.Text(" and level "))
		tooltip.WriteString(a.LevelCriteria.String())
	}
	return satisfied
}
