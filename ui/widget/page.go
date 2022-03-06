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
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/res"
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
	label.SetBorder(unison.NewEmptyBorder(geom32.Insets{Bottom: 1})) // To match field underline spacing
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
	label.SetBorder(unison.NewEmptyBorder(geom32.Insets{Bottom: 1})) // To match field underline spacing
	return label
}

// NewPageLabelWithRandomizer creates a new end-aligned field label for a sheet page that includes a randomization
// button.
func NewPageLabelWithRandomizer(title, tooltip string, clickCallback func()) *unison.Panel {
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
	})
	wrapper.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	b := unison.NewButton()
	b.ButtonTheme = unison.DefaultSVGButtonTheme
	b.DrawableOnlyVMargin = 1
	b.DrawableOnlyHMargin = 1
	b.HideBase = true
	baseline := theme.PageLabelPrimaryFont.Baseline()
	size := geom32.NewSize(baseline, baseline)
	b.Drawable = &unison.DrawableSVG{
		SVG:  res.RandomizeSVG,
		Size: *size.GrowToInteger(),
	}
	if tooltip != "" {
		b.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	b.ClickCallback = clickCallback
	b.SetLayoutData(&unison.FlexLayoutData{HGrab: true})
	wrapper.AddChild(b)
	wrapper.AddChild(NewPageLabelEnd(title))
	return wrapper
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

// NewHeightPageField creates a new height entry field for a sheet page.
func NewHeightPageField(entity *gurps.Entity, value, max measure.Length, applier func(measure.Length)) *HeightField {
	field := NewHeightField(entity, value, max, applier)
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

// NewWeightPageField creates a new weight entry field for a sheet page.
func NewWeightPageField(entity *gurps.Entity, value, max measure.Weight, applier func(measure.Weight)) *WeightField {
	field := NewWeightField(entity, value, max, applier)
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

// NewSignedIntegerPageField creates a new signed integer entry field for a sheet page.
func NewSignedIntegerPageField(value, min, max int, applier func(int)) *SignedIntegerField {
	field := NewSignedIntegerField(value, min, max, applier)
	field.HAlign = unison.EndAlignment
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
	label.SetBorder(unison.NewEmptyBorder(geom32.Insets{
		Left:   1,
		Bottom: 1,
		Right:  1,
	})) // Match normal fields
	label.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	return label
}
