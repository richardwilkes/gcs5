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
	"github.com/richardwilkes/gcs/ui/workspace/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// EditNote displays the editor for a note.
func EditNote(owner widget.Rebuildable, note *gurps.Note) {
	displayEditor[*gurps.Note, *noteEditorData](owner, note, initNoteEditor)
}

func initNoteEditor(e *editor[*gurps.Note, *noteEditorData], content *unison.Panel) {
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
