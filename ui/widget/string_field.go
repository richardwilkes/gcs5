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

package widget

import (
	"github.com/richardwilkes/gcs/model/undo"
	"github.com/richardwilkes/unison"
)

// StringField holds the value for a string field.
type StringField struct {
	*unison.Field
	undoID    int
	undoTitle string
	get       func() string
	set       func(string)
	inUndo    bool
}

// NewStringField creates a new field for editing a string.
func NewStringField(undoID int, undoTitle string, get func() string, set func(string)) *StringField {
	f := &StringField{
		Field:     unison.NewField(),
		undoID:    undoID,
		undoTitle: undoTitle,
		get:       get,
		set:       set,
	}
	f.Self = f
	f.ModifiedCallback = f.modified
	f.Sync()
	f.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	return f
}

func (f *StringField) modified() {
	text := f.Text()
	if !f.inUndo && f.undoID != undo.NoneID {
		if mgr := undo.Manager(f); mgr != nil {
			mgr.Add(&undo.Edit[string]{
				ID:       f.undoID,
				EditName: f.undoTitle,
				EditCost: 1,
				UndoFunc: func(e *undo.Edit[string]) { f.setWithoutUndo(e.BeforeData, true) },
				RedoFunc: func(e *undo.Edit[string]) { f.setWithoutUndo(e.AfterData, true) },
				AbsorbFunc: func(e *undo.Edit[string], other unison.UndoEdit) bool {
					if e2, ok := other.(*undo.Edit[string]); ok && e2.ID == f.undoID {
						e.AfterData = e2.AfterData
						return true
					}
					return false
				},
				BeforeData: f.get(),
				AfterData:  text,
			})
		}
	}
	f.set(text)
	MarkForLayoutWithinDockable(f)
	MarkModified(f)
}

func (f *StringField) setWithoutUndo(text string, focus bool) {
	f.inUndo = true
	f.SetText(text)
	f.inUndo = false
	if focus {
		f.RequestFocus()
		f.SelectAll()
	}
}

// Sync the field to the current value.
func (f *StringField) Sync() {
	f.setWithoutUndo(f.get(), false)
}
