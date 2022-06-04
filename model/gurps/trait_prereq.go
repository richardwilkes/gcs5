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
)

var _ Prereq = &TraitPrereq{}

// TraitPrereq holds a prereq against a Trait.
type TraitPrereq struct {
	Parent        *PrereqList      `json:"-"`
	Type          prereq.Type      `json:"type"`
	Has           bool             `json:"has"`
	NameCriteria  criteria.String  `json:"name,omitempty"`
	LevelCriteria criteria.Numeric `json:"level,omitempty"`
	NotesCriteria criteria.String  `json:"notes,omitempty"`
}

// NewTraitPrereq creates a new TraitPrereq.
func NewTraitPrereq() *TraitPrereq {
	return &TraitPrereq{
		Type: prereq.Trait,
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

// PrereqType implements Prereq.
func (a *TraitPrereq) PrereqType() prereq.Type {
	return a.Type
}

// ParentList implements Prereq.
func (a *TraitPrereq) ParentList() *PrereqList {
	return a.Parent
}

// Clone implements Prereq.
func (a *TraitPrereq) Clone(parent *PrereqList) Prereq {
	clone := *a
	clone.Parent = parent
	return &clone
}

// FillWithNameableKeys implements Prereq.
func (a *TraitPrereq) FillWithNameableKeys(m map[string]string) {
	nameables.Extract(a.NameCriteria.Qualifier, m)
	nameables.Extract(a.NotesCriteria.Qualifier, m)
}

// ApplyNameableKeys implements Prereq.
func (a *TraitPrereq) ApplyNameableKeys(m map[string]string) {
	a.NameCriteria.Qualifier = nameables.Apply(a.NameCriteria.Qualifier, m)
	a.NotesCriteria.Qualifier = nameables.Apply(a.NotesCriteria.Qualifier, m)
}

// Satisfied implements Prereq.
func (a *TraitPrereq) Satisfied(entity *Entity, exclude any, tooltip *xio.ByteBuffer, prefix string) bool {
	satisfied := false
	TraverseTraits(func(t *Trait) bool {
		if exclude == t || !a.NameCriteria.Matches(t.Name) {
			return false
		}
		notes := t.Notes()
		if modNotes := t.ModifierNotes(); modNotes != "" {
			notes += "\n" + modNotes
		}
		if !a.NotesCriteria.Matches(notes) {
			return false
		}
		satisfied = a.LevelCriteria.Matches(t.Levels.Max(0))
		return satisfied
	}, true, entity.Traits...)
	if !a.Has {
		satisfied = !satisfied
	}
	if !satisfied && tooltip != nil {
		tooltip.WriteString(prefix)
		tooltip.WriteString(HasText(a.Has))
		tooltip.WriteString(i18n.Text(" a trait whose name "))
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
