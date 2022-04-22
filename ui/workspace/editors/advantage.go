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

// EditAdvantage displays the editor for an advantage.
func EditAdvantage(owner widget.Rebuildable, advantage *gurps.Advantage) {
	displayEditor[*gurps.Advantage, *advantageEditorData](owner, advantage, initAdvantageEditor)
}

func initAdvantageEditor(e *editor[*gurps.Advantage, *advantageEditorData], content *unison.Panel) {
	content.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})

	content.AddChild(unison.NewPanel())
	content.AddChild(widget.NewCheckBox(i18n.Text("Enabled"), !e.editorData.disabled, func(b bool) {
		e.editorData.disabled = !b
		widget.MarkModified(content)
	}))

	labelText := i18n.Text("Name")
	content.AddChild(widget.NewFieldLeadingLabel(labelText))
	field := widget.NewStringField(labelText, func() string { return e.editorData.name },
		func(value string) {
			e.editorData.name = value
			widget.MarkModified(content)
		})
	content.AddChild(field)

	labelText = i18n.Text("Notes")
	content.AddChild(widget.NewFieldLeadingLabel(labelText))
	field = widget.NewMultiLineStringField(labelText, func() string { return e.editorData.notes },
		func(value string) {
			e.editorData.notes = value
			content.MarkForLayoutAndRedraw()
			widget.MarkModified(content)
		})
	field.AutoScroll = false
	content.AddChild(field)

	labelText = i18n.Text("VTT Notes")
	tooltip := i18n.Text("Any notes for VTT use; see the instructions for your VVT to determine if/how these can be used")
	label := widget.NewFieldLeadingLabel(labelText)
	label.Tooltip = unison.NewTooltipWithText(tooltip)
	content.AddChild(label)
	field = widget.NewMultiLineStringField(labelText, func() string { return e.editorData.vttNotes },
		func(value string) {
			e.editorData.vttNotes = value
			content.MarkForLayoutAndRedraw()
			widget.MarkModified(content)
		})
	field.Tooltip = unison.NewTooltipWithText(tooltip)
	field.AutoScroll = false
	content.AddChild(field)

	labelText = i18n.Text("User Description")
	tooltip = i18n.Text("Additional notes for your own reference. These only exist in character sheets and will be removed if transferred to a data list or template")
	label = widget.NewFieldLeadingLabel(labelText)
	label.Tooltip = unison.NewTooltipWithText(tooltip)
	content.AddChild(label)
	field = widget.NewMultiLineStringField(labelText, func() string { return e.editorData.userDesc },
		func(value string) {
			e.editorData.userDesc = value
			content.MarkForLayoutAndRedraw()
			widget.MarkModified(content)
		})
	field.Tooltip = unison.NewTooltipWithText(tooltip)
	field.AutoScroll = false
	content.AddChild(field)

	labelText = i18n.Text("Tags")
	tooltip = i18n.Text("Separate multiple tags with commas")
	label = widget.NewFieldLeadingLabel(labelText)
	label.Tooltip = unison.NewTooltipWithText(tooltip)
	content.AddChild(label)
	field = widget.NewMultiLineStringField(labelText, func() string { return e.editorData.tags },
		func(value string) {
			e.editorData.tags = value
			content.MarkForLayoutAndRedraw()
			widget.MarkModified(content)
		})
	field.Tooltip = unison.NewTooltipWithText(tooltip)
	field.AutoScroll = false
	content.AddChild(field)

	if !e.target.Container() {
		content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Type")))
		wrapper := unison.NewPanel()
		wrapper.SetLayout(&unison.FlowLayout{
			HSpacing: unison.StdHSpacing,
			VSpacing: unison.StdVSpacing,
		})
		content.AddChild(wrapper)
		wrapper.AddChild(widget.NewCheckBox(i18n.Text("Mental"), e.editorData.mental, func(b bool) {
			e.editorData.mental = b
			widget.MarkModified(content)
		}))
		wrapper.AddChild(widget.NewCheckBox(i18n.Text("Physical"), e.editorData.physical, func(b bool) {
			e.editorData.physical = b
			widget.MarkModified(content)
		}))
		wrapper.AddChild(widget.NewCheckBox(i18n.Text("Social"), e.editorData.social, func(b bool) {
			e.editorData.social = b
			widget.MarkModified(content)
		}))
		wrapper.AddChild(widget.NewCheckBox(i18n.Text("Exotic"), e.editorData.exotic, func(b bool) {
			e.editorData.exotic = b
			widget.MarkModified(content)
		}))
		wrapper.AddChild(widget.NewCheckBox(i18n.Text("Supernatural"), e.editorData.supernatural, func(b bool) {
			e.editorData.supernatural = b
			widget.MarkModified(content)
		}))
	}

	labelText = i18n.Text("Page")
	label = widget.NewFieldLeadingLabel(labelText)
	label.Tooltip = unison.NewTooltipWithText(tbl.PageRefTooltipText)
	content.AddChild(label)
	field = widget.NewStringField(labelText, func() string { return e.editorData.pageRef },
		func(value string) {
			e.editorData.pageRef = value
			widget.MarkModified(content)
		})
	field.Tooltip = unison.NewTooltipWithText(tbl.PageRefTooltipText)
	content.AddChild(field)
}
