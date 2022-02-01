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
	"github.com/richardwilkes/toolbox/xio"
)

var _ Prereq = &AdvantagePrereq{}

// AdvantagePrereq holds a prereq against an Advantage.
type AdvantagePrereq struct {
	Parent        *PrereqList      `json:"-"`
	Type          prereq.Type      `json:"type"`
	NameCriteria  criteria.String  `json:"name"`
	LevelCriteria criteria.Numeric `json:"level"`
	NotesCriteria criteria.String  `json:"notes"`
	Has           bool             `json:"has,omitempty"`
}

// NewAdvantagePrereq creates a new AdvantagePrereq.
func NewAdvantagePrereq() *AdvantagePrereq {
	return &AdvantagePrereq{
		Type: prereq.Advantage,
		NameCriteria: criteria.String{
			Compare: criteria.Is,
		},
		LevelCriteria: criteria.Numeric{
			Compare: criteria.AtLeast,
		},
		NotesCriteria: criteria.String{
			Compare: criteria.Any,
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
func (a *AdvantagePrereq) Satisfied(entity *Entity, exclude interface{}, buffer *xio.ByteBuffer, prefix string) bool {
	satisfied := false
	// TODO: Implement
	/*
	   StringCriteria  nameCriteria  = getNameCriteria();
	   IntegerCriteria levelCriteria = getLevelCriteria();

	   for (Advantage advantage : character.getAdvantagesIterator(false)) {
	       if (exclude != advantage && nameCriteria.matches(advantage.getName())) {
	           String notes         = advantage.getNotes();
	           String modifierNotes = advantage.getModifierNotes();

	           if (!modifierNotes.isEmpty()) {
	               notes = modifierNotes + '\n' + notes;
	           }
	           if (mNotesCriteria.matches(notes)) {
	               int levels = advantage.getLevels();
	               if (levels < 0) {
	                   levels = 0;
	               }
	               satisfied = levelCriteria.matches(levels);
	               break;
	           }
	       }
	   }
	   if (!has()) {
	       satisfied = !satisfied;
	   }
	   if (!satisfied && builder != null) {
	       builder.append(MessageFormat.format(I18n.text("\n{0}{1} an advantage whose name {2}"), prefix, getHasText(), nameCriteria.toString()));
	       if (!mNotesCriteria.isTypeAnything()) {
	           builder.append(MessageFormat.format(I18n.text(", notes {0},"), mNotesCriteria.toString()));
	       }
	       builder.append(MessageFormat.format(I18n.text(" and level {0}"), levelCriteria.toString()));
	   }
	*/
	return satisfied
}
