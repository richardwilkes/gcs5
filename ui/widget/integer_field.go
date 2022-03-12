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
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/mathf32"
	"github.com/richardwilkes/unison"
)

// IntegerField holds the value for an integer field.
type IntegerField struct {
	*unison.Field
	value   *int
	minimum int
	maximum int
}

// NewIntegerField creates a new field that holds an integer.
func NewIntegerField(value *int, min, max int) *IntegerField {
	f := &IntegerField{
		Field:   unison.NewField(),
		value:   value,
		minimum: min,
		maximum: max,
	}
	f.Self = f
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	if min != math.MinInt && max != math.MaxInt {
		f.MinimumTextWidth = mathf32.Max(f.Font.SimpleWidth(strconv.Itoa(min)), f.Font.SimpleWidth(strconv.Itoa(max)))
	}
	f.Sync()
	return f
}

// Value returns the current value of the field.
func (f *IntegerField) Value() int {
	return *f.value
}

// SetValue sets the value of this field, marking the field and all of its parents as needing to be laid out again if the
// value is not what is currently in the field.
func (f *IntegerField) SetValue(value int) {
	if f.minimum != math.MinInt && value < f.minimum {
		value = f.minimum
	} else if f.maximum != math.MaxInt && value > f.maximum {
		value = f.maximum
	}
	SetFieldValue(f.Field, strconv.Itoa(value))
}

func (f *IntegerField) trimmed() string {
	return strings.TrimSpace(f.Text())
}

func (f *IntegerField) validate() bool {
	v, err := strconv.Atoi(f.trimmed())
	if err != nil {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Invalid integer"))
		return false
	}
	if f.minimum != math.MinInt && v < f.minimum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Integer must be at least %d"), f.minimum))
		return false
	}
	if f.maximum != math.MaxInt && v > f.maximum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Integer must be no more than %d"), f.maximum))
		return false
	}
	f.Tooltip = nil
	return true
}

func (f *IntegerField) modified() {
	if v, err := strconv.Atoi(f.trimmed()); err == nil &&
		(f.minimum == math.MinInt || v >= f.minimum) &&
		(f.maximum == math.MaxInt || v <= f.maximum) {
		*f.value = v
		MarkForLayoutWithinDockable(f)
		MarkModified(f)
	}
}

// Sync the field to the current value.
func (f *IntegerField) Sync() {
	f.SetValue(*f.value)
}
