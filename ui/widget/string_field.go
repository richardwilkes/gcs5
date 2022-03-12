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

// StringField holds the value for a string field.
type StringField struct {
	*unison.Field
	value *string
}

// NewStringField creates a new field for editing a string.
func NewStringField(data *string) *StringField {
	f := &StringField{
		Field: unison.NewField(),
		value: data,
	}
	f.Self = f
	f.ModifiedCallback = func() {
		*f.value = f.Text()
		MarkForLayoutWithinDockable(f)
		MarkModified(f)
	}
	f.Sync()
	f.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	return f
}

// Sync the field to the current value.
func (f *StringField) Sync() {
	f.SetText(*f.value)
}
