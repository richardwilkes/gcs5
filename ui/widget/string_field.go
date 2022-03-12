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
	"github.com/richardwilkes/unison"
)

// NewStringField creates a new string field.
func NewStringField(value string, applier func(string)) *unison.Field {
	f := unison.NewField()
	f.SetText(value)
	f.ModifiedCallback = func() { applier(f.Text()) }
	f.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	return f
}

// SetFieldValue sets the value of this field, marking the field and all of its parents as needing to be laid out again
// if the value is not what is currently in the field.
func SetFieldValue(field *unison.Field, value string) {
	if value != field.Text() {
		field.SetText(value)
		MarkForLayoutWithinDockable(field)
	}
}

// MarkForLayoutWithinDockable sets the NeedsLayout flag on the provided panel and all of its parents up to the first
// Dockable.
func MarkForLayoutWithinDockable(panel unison.Paneler) {
	p := panel.AsPanel()
	for p != nil {
		p.NeedsLayout = true
		if _, ok := p.Self.(unison.Dockable); ok {
			break
		}
		p = p.Parent()
	}
}
