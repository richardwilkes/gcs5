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
	"unicode"

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/mathf32"
	"github.com/richardwilkes/unison"
)

// IntegerField holds the value for an integer field.
type IntegerField struct {
	*unison.Field
	applier  func()
	value    *int
	minimum  int
	maximum  int
	showSign bool
}

// NewIntegerField creates a new field that holds an integer.
func NewIntegerField(value *int, min, max int, showSign bool, applier func()) *IntegerField {
	f := &IntegerField{
		Field:    unison.NewField(),
		applier:  applier,
		value:    value,
		minimum:  min,
		maximum:  max,
		showSign: showSign,
	}
	f.Self = f
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	f.RuneTypedCallback = f.runeTyped
	if min != math.MinInt && max != math.MaxInt {
		f.MinimumTextWidth = mathf32.Max(f.Font.SimpleWidth(f.formatted(min)), f.Font.SimpleWidth(f.formatted(max)))
	}
	f.Sync()
	return f
}

func (f *IntegerField) formatted(value int) string {
	if f.showSign {
		return fmt.Sprintf("%+d", value)
	}
	return strconv.Itoa(value)
}

// Value returns the current value of the field.
func (f *IntegerField) Value() int {
	return *f.value
}

// SetValue sets the value of this field, applying any constraints.
func (f *IntegerField) SetValue(value int) {
	if f.minimum != math.MinInt && value < f.minimum {
		value = f.minimum
	} else if f.maximum != math.MaxInt && value > f.maximum {
		value = f.maximum
	}
	SetFieldValue(f.Field, f.formatted(value))
}

func (f *IntegerField) trimmed(text string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(text), "+"))
}

func (f *IntegerField) validate() bool {
	v, err := strconv.Atoi(f.trimmed(f.Text()))
	if err != nil {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Invalid integer"))
		return false
	}
	if f.minimum != math.MinInt && v < f.minimum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Integer must be at least %s"), f.formatted(f.minimum)))
		return false
	}
	if f.maximum != math.MaxInt && v > f.maximum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Integer must be no more than %s"), f.formatted(f.maximum)))
		return false
	}
	f.Tooltip = nil
	return true
}

func (f *IntegerField) modified() {
	if v, err := strconv.Atoi(f.trimmed(f.Text())); err == nil &&
		(f.minimum == math.MinInt || v >= f.minimum) &&
		(f.maximum == math.MaxInt || v <= f.maximum) {
		*f.value = v
		if f.applier != nil {
			f.applier()
		}
		MarkForLayoutWithinDockable(f)
		MarkModified(f)
	}
}

// Sync the field to the current value.
func (f *IntegerField) Sync() {
	f.SetValue(*f.value)
}

func (f *IntegerField) runeTyped(ch rune) bool {
	if !unicode.IsControl(ch) {
		if f.minimum >= 0 && ch == '-' {
			unison.Beep()
			return false
		}
		if text := f.trimmed(string(f.RunesIfPasted([]rune{ch}))); text != "-" {
			if _, err := strconv.Atoi(text); err != nil {
				unison.Beep()
				return false
			}
		}
	}
	return f.DefaultRuneTyped(ch)
}
