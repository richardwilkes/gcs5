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
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
	"github.com/richardwilkes/unison"
)

func addNameLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndStringField(parent, i18n.Text("Name"), "", fieldData)
}

func addPageRefLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndStringField(parent, i18n.Text("Page"), tbl.PageRefTooltipText, fieldData)
}

func addNotesLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndMultiLineStringField(parent, i18n.Text("Notes"), "", fieldData)
}

func addVTTNotesLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndMultiLineStringField(parent, i18n.Text("VTT Notes"),
		i18n.Text("Any notes for VTT use; see the instructions for your VVT to determine if/how these can be used"),
		fieldData)
}

func addUserDescLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndMultiLineStringField(parent, i18n.Text("User Description"),
		i18n.Text("Additional notes for your own reference. These only exist in character sheets and will be removed if transferred to a data list or template"),
		fieldData)
}

func addTagsLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndMultiLineStringField(parent, i18n.Text("Tags"), i18n.Text("Separate multiple tags with commas"),
		fieldData)
}

func addLabelAndStringField(parent *unison.Panel, labelText, tooltip string, fieldData *string) {
	label := widget.NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	field := widget.NewStringField(labelText, func() string { return *fieldData },
		func(value string) {
			*fieldData = value
			widget.MarkModified(parent)
		})
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
}

func addLabelAndMultiLineStringField(parent *unison.Panel, labelText, tooltip string, fieldData *string) {
	label := widget.NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	field := widget.NewMultiLineStringField(labelText, func() string { return *fieldData },
		func(value string) {
			*fieldData = value
			parent.MarkForLayoutAndRedraw()
			widget.MarkModified(parent)
		})
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	field.AutoScroll = false
	parent.AddChild(field)
}

func addLabelAndNumericField(parent *unison.Panel, labelText, tooltip string, fieldData *f64d4.Int, min, max f64d4.Int) *widget.NumericField {
	label := widget.NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	field := widget.NewNumericField(labelText, func() f64d4.Int { return *fieldData },
		func(value f64d4.Int) {
			*fieldData = value
			widget.MarkModified(parent)
		}, min, max, false)
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
	return field
}

func addNumericField(parent *unison.Panel, labelText, tooltip string, fieldData *f64d4.Int, min, max f64d4.Int) {
	field := widget.NewNumericField(labelText, func() f64d4.Int { return *fieldData },
		func(value f64d4.Int) {
			*fieldData = value
			widget.MarkModified(parent)
		}, min, max, false)
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
}

func addCheckBox(parent *unison.Panel, labelText string, fieldData *bool) {
	parent.AddChild(widget.NewCheckBox(labelText, *fieldData, func(b bool) {
		*fieldData = b
		widget.MarkModified(parent)
	}))
}

func addInvertedCheckBox(parent *unison.Panel, labelText string, fieldData *bool) {
	parent.AddChild(widget.NewCheckBox(labelText, !*fieldData, func(b bool) {
		*fieldData = !b
		widget.MarkModified(parent)
	}))
}

func addFlowWrapper(parent *unison.Panel, labelText string, count int) *unison.Panel {
	parent.AddChild(widget.NewFieldLeadingLabel(labelText))
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{
		Columns:  count,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
		VAlign:   unison.MiddleAlignment,
	})
	parent.AddChild(wrapper)
	return wrapper
}

func addLabelAndPopup[T comparable](parent *unison.Panel, labelText, tooltip string, choices []T, fieldData *T) *unison.PopupMenu[T] {
	label := widget.NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	popup := unison.NewPopupMenu[T]()
	for _, one := range choices {
		popup.AddItem(one)
	}
	popup.Select(*fieldData)
	popup.SelectionCallback = func(_ int, item T) {
		*fieldData = item
		widget.MarkModified(parent)
	}
	parent.AddChild(popup)
	return popup
}

func disableAndBlankField(field *widget.NumericField) {
	field.SetEnabled(false)
	field.DrawOverCallback = func(gc *unison.Canvas, rect unison.Rect) {
		rect = field.ContentRect(false)
		gc.DrawRect(rect, field.BackgroundInk.Paint(gc, rect, unison.Fill))
	}
}

func enableAndUnblankField(field *widget.NumericField) {
	field.SetEnabled(true)
	field.DrawOverCallback = nil
}

func disableAndBlankPopup[T comparable](popup *unison.PopupMenu[T]) {
	popup.SetEnabled(false)
	popup.DrawOverCallback = func(gc *unison.Canvas, rect unison.Rect) {
		rect = popup.ContentRect(false)
		unison.DrawRoundedRectBase(gc, rect, popup.CornerRadius, 1, popup.BackgroundInk, popup.EdgeInk)
	}
}

func enableAndUnblankPopup[T comparable](popup *unison.PopupMenu[T]) {
	popup.SetEnabled(true)
	popup.DrawOverCallback = nil
}
