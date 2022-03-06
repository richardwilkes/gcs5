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
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/toolbox/xmath/mathf32"
	"github.com/richardwilkes/unison"
)

var nonEditableFieldColor = unison.NewDynamicColor(func() unison.Color {
	return unison.OnContentColor.GetColor().SetAlphaIntensity(0.375)
})

// NewPageLabel creates a new start-aligned field label for a sheet page.
func NewPageLabel(title string) *unison.Label {
	label := unison.NewLabel()
	label.OnBackgroundInk = unison.OnContentColor
	label.Text = title
	label.Font = theme.PageLabelPrimaryFont
	label.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	return label
}

// NewPageLabelEnd creates a new end-aligned field label for a sheet page.
func NewPageLabelEnd(title string) *unison.Label {
	label := unison.NewLabel()
	label.OnBackgroundInk = unison.OnContentColor
	label.Text = title
	label.Font = theme.PageLabelPrimaryFont
	label.HAlign = unison.EndAlignment
	label.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	return label
}

// NewStringPageField creates a new text entry field for a sheet page.
func NewStringPageField(value string, applier func(string)) *unison.Field {
	field := NewStringField(value, applier)
	field.Font = theme.PageFieldPrimaryFont
	field.FocusedBorder = unison.NewLineBorder(theme.AccentColor, 0, geom32.Insets{Bottom: 1}, false)
	field.UnfocusedBorder = unison.NewLineBorder(unison.ControlEdgeColor, 0, geom32.Insets{Bottom: 1}, false)
	field.SetBorder(field.UnfocusedBorder)
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	return field
}

// NewStringPageFieldNoGrab creates a new text entry field for a sheet page, but with HGrab set to false.
func NewStringPageFieldNoGrab(value string, applier func(string)) *unison.Field {
	field := NewStringField(value, applier)
	field.Font = theme.PageFieldPrimaryFont
	field.FocusedBorder = unison.NewLineBorder(theme.AccentColor, 0, geom32.Insets{Bottom: 1}, false)
	field.UnfocusedBorder = unison.NewLineBorder(unison.ControlEdgeColor, 0, geom32.Insets{Bottom: 1}, false)
	field.SetBorder(field.UnfocusedBorder)
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	return field
}

// NewNumericPageField creates a new numeric text entry field for a sheet page.
func NewNumericPageField(value, min, max fixed.F64d4, applier func(fixed.F64d4)) *NumericField {
	field := NewNumericField(value, min, max, applier)
	field.HAlign = unison.EndAlignment
	field.Font = theme.PageFieldPrimaryFont
	field.FocusedBorder = unison.NewLineBorder(theme.AccentColor, 0, geom32.Insets{Bottom: 1}, false)
	field.UnfocusedBorder = unison.NewLineBorder(unison.ControlEdgeColor, 0, geom32.Insets{Bottom: 1}, false)
	field.SetBorder(field.UnfocusedBorder)
	// Override to ignore fractional values
	field.MinimumTextWidth = mathf32.Max(field.Font.SimpleWidth(min.Trunc().String()),
		field.Font.SimpleWidth(max.Trunc().String()))
	return field
}

// NewNonEditablePageField creates a new start-aligned non-editable field that uses the same font and size as the page
// field.
func NewNonEditablePageField(title, tooltip string) *unison.Label {
	return newNonEditablePageField(title, tooltip, unison.StartAlignment)
}

// NewNonEditablePageFieldEnd creates a new end-aligned non-editable field that uses the same font and size as the page
// field.
func NewNonEditablePageFieldEnd(title, tooltip string) *unison.Label {
	return newNonEditablePageField(title, tooltip, unison.EndAlignment)
}

func newNonEditablePageField(title, tooltip string, hAlign unison.Alignment) *unison.Label {
	label := unison.NewLabel()
	label.OnBackgroundInk = nonEditableFieldColor
	label.Text = title
	label.Font = theme.PageFieldPrimaryFont
	label.HAlign = hAlign
	label.SetBorder(unison.NewEmptyBorder(geom32.NewHorizontalInsets(1))) // Account for cursor in normal fields
	label.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	return label
}
