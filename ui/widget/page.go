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
	return unison.OnEditableColor.GetColor().SetAlphaIntensity(0.6)
})

// NewPageHeader creates a new center-aligned header for a sheet page.
func NewPageHeader(title string, hSpan int) *unison.Label {
	label := unison.NewLabel()
	label.OnBackgroundInk = theme.OnHeaderColor
	label.Text = title
	label.Font = theme.PageLabelPrimaryFont
	label.HAlign = unison.MiddleAlignment
	label.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  hSpan,
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	label.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		gc.DrawRect(rect, theme.HeaderColor.Paint(gc, rect, unison.Fill))
		label.DefaultDraw(gc, rect)
	}
	return label
}

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

// NewPageLabelCenter creates a new center-aligned field label for a sheet page.
func NewPageLabelCenter(title string) *unison.Label {
	label := unison.NewLabel()
	label.OnBackgroundInk = unison.OnContentColor
	label.Text = title
	label.Font = theme.PageLabelPrimaryFont
	label.HAlign = unison.MiddleAlignment
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
	b.SetFocusable(false)
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
func NewStringPageField(get func() string, set func(string)) *StringField {
	field := NewStringField(get, set)
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
func NewStringPageFieldNoGrab(get func() string, set func(string)) *StringField {
	field := NewStringField(get, set)
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

// NewIntegerPageField creates a new integer entry field for a sheet page.
func NewIntegerPageField(get func() int, set func(int), min, max int, showSign bool) *IntegerField {
	field := NewIntegerField(get, set, min, max, showSign)
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
func NewNumericPageField(get func() fixed.F64d4, set func(fixed.F64d4), min, max fixed.F64d4, noMinWidth bool) *NumericField {
	field := NewNumericField(get, set, min, max, noMinWidth)
	field.HAlign = unison.EndAlignment
	field.Font = theme.PageFieldPrimaryFont
	field.FocusedBorder = unison.NewLineBorder(theme.AccentColor, 0, geom32.Insets{Bottom: 1}, false)
	field.UnfocusedBorder = unison.NewLineBorder(unison.ControlEdgeColor, 0, geom32.Insets{Bottom: 1}, false)
	field.SetBorder(field.UnfocusedBorder)
	if !noMinWidth && min != fixed.F64d4Min && max != fixed.F64d4Max {
		// Override to ignore fractional values
		field.MinimumTextWidth = mathf32.Max(field.Font.SimpleWidth(min.Trunc().String()),
			field.Font.SimpleWidth(max.Trunc().String()))
	}
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	return field
}