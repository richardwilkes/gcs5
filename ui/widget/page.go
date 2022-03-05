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

// NewPageLabel creates a new field label for a sheet page.
func NewPageLabel(title string) *unison.Label {
	label := unison.NewLabel()
	label.OnBackgroundInk = unison.OnContentColor
	label.Text = title
	label.Font = theme.PageLabelPrimaryFont
	return label
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

// NewNonEditablePageField creates a new non-editable field that uses the same font and size as the page field.
func NewNonEditablePageField(title string, alignment unison.Alignment) *unison.Label {
	label := unison.NewLabel()
	label.OnBackgroundInk = unison.OnContentColor
	label.Text = title
	label.Font = theme.PageFieldPrimaryFont
	label.HAlign = alignment
	label.SetBorder(unison.NewEmptyBorder(geom32.NewHorizontalInsets(1))) // Account for cursor in normal fields
	return label
}
