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

package editors

import (
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type noteEditor struct {
	Editor
	owner    widget.Rebuildable
	note     *gurps.Note
	noteText string
	pageRef  string
}

// EditNote displays the editor for a note.
func EditNote(owner widget.Rebuildable, note *gurps.Note) {
	ws, dc, found := workspace.Activate(func(d unison.Dockable) bool {
		if editor, ok := d.(*noteEditor); ok {
			return editor.owner == owner && editor.note == note
		}
		return false
	})
	if !found && ws != nil {
		e := &noteEditor{
			owner:    owner,
			note:     note,
			noteText: note.Text,
			pageRef:  note.PageRef,
		}
		e.Self = e
		e.TabTitle = i18n.Text("Note Editor")
		e.TabTitle += i18n.Text(" for ") + owner.String()
		e.Setup(ws, dc, e.initContent)
		e.IsModifiedCallback = func() bool {
			return e.note.Text != e.noteText || e.note.PageRef != e.pageRef
		}
		e.ApplyCallback = func() {
			if mgr := unison.UndoManagerFor(e.owner); mgr != nil {
				mgr.Add(&unison.UndoEdit[[]string]{
					ID:       unison.NextUndoID(),
					EditName: i18n.Text("Note Changes"),
					EditCost: 1,
					UndoFunc: func(edit *unison.UndoEdit[[]string]) {
						note.Text = edit.BeforeData[0]
						note.PageRef = edit.BeforeData[1]
						owner.MarkForRebuild(false)
					},
					RedoFunc: func(edit *unison.UndoEdit[[]string]) {
						note.Text = edit.AfterData[0]
						note.PageRef = edit.AfterData[1]
						owner.MarkForRebuild(false)
					},
					BeforeData: []string{note.Text, note.PageRef},
					AfterData:  []string{e.noteText, e.pageRef},
				})
			}
			e.note.Text = e.noteText
			e.note.PageRef = e.pageRef
			e.owner.MarkForRebuild(false)
		}
	}
}

func (e *noteEditor) initContent(content *unison.Panel) {
	content.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})

	noteLabel := i18n.Text("Note")
	content.AddChild(widget.NewFieldLeadingLabel(noteLabel))
	noteField := widget.NewMultiLineStringField(noteLabel, func() string { return e.noteText },
		func(value string) {
			e.noteText = value
			content.MarkForLayoutAndRedraw()
			widget.MarkModified(content)
		})
	noteField.AutoScroll = false
	content.AddChild(noteField)

	pageLabel := i18n.Text("Page")
	content.AddChild(widget.NewFieldLeadingLabel(pageLabel))
	content.AddChild(widget.NewStringField(pageLabel, func() string { return e.pageRef },
		func(value string) {
			e.pageRef = value
			widget.MarkModified(content)
		}))
}
