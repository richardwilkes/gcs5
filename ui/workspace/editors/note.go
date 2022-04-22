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
	"fmt"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/gcs/ui/workspace/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type noteEditor struct {
	Editor
	target     *gurps.Note
	beforeData noteEditorData
	editorData noteEditorData
}

// EditNote displays the editor for a note.
func EditNote(owner widget.Rebuildable, note *gurps.Note) {
	ws, dc, found := workspace.Activate(func(d unison.Dockable) bool {
		if editor, ok := d.(*noteEditor); ok {
			return editor.owner == owner && editor.target == note
		}
		return false
	})
	if !found && ws != nil {
		e := &noteEditor{
			Editor: Editor{
				owner: owner,
			},
			target: note,
		}
		e.Self = e
		e.beforeData.From(note)
		e.editorData = e.beforeData
		e.TabTitle = fmt.Sprintf(i18n.Text("%s Editor for %s"), note.Kind(), owner.String())
		e.Setup(ws, dc, e.initContent)
		e.IsModifiedCallback = func() bool { return e.beforeData != e.editorData }
		e.ApplyCallback = func() {
			if mgr := unison.UndoManagerFor(e.owner); mgr != nil {
				before := e.beforeData
				after := e.editorData
				mgr.Add(&unison.UndoEdit[*noteEditorData]{
					ID:       unison.NextUndoID(),
					EditName: fmt.Sprintf(i18n.Text("%s Changes"), note.Kind()),
					EditCost: 1,
					UndoFunc: func(edit *unison.UndoEdit[*noteEditorData]) {
						edit.BeforeData.Apply(note)
						owner.MarkForRebuild(true)
					},
					RedoFunc: func(edit *unison.UndoEdit[*noteEditorData]) {
						edit.AfterData.Apply(note)
						owner.MarkForRebuild(true)
					},
					BeforeData: &before,
					AfterData:  &after,
				})
			}
			e.editorData.Apply(e.target)
			e.owner.MarkForRebuild(true)
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
	noteField := widget.NewMultiLineStringField(noteLabel, func() string { return e.editorData.note },
		func(value string) {
			e.editorData.note = value
			content.MarkForLayoutAndRedraw()
			widget.MarkModified(content)
		})
	noteField.AutoScroll = false
	content.AddChild(noteField)

	pageLabel := i18n.Text("Page")
	label := widget.NewFieldLeadingLabel(pageLabel)
	label.Tooltip = unison.NewTooltipWithText(tbl.PageRefTooltipText)
	content.AddChild(label)
	field := widget.NewStringField(pageLabel, func() string { return e.editorData.pageRef },
		func(value string) {
			e.editorData.pageRef = value
			widget.MarkModified(content)
		})
	field.Tooltip = unison.NewTooltipWithText(tbl.PageRefTooltipText)
	content.AddChild(field)
}
