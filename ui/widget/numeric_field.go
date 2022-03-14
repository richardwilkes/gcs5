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
	"strconv"
	"strings"
	"unicode"

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/mathf32"
	"github.com/richardwilkes/unison"
)

// NumericField holds the value for a numeric field.
type NumericField struct {
	*unison.Field
	get        func() fixed.F64d4
	set        func(fixed.F64d4)
	minimum    fixed.F64d4
	maximum    fixed.F64d4
	noMinWidth bool
}

// NewNumericField creates a new field that holds a fixed-point number.
func NewNumericField(get func() fixed.F64d4, set func(fixed.F64d4), min, max fixed.F64d4, noMinWidth bool) *NumericField {
	f := &NumericField{
		Field:      unison.NewField(),
		get:        get,
		set:        set,
		minimum:    min,
		noMinWidth: noMinWidth,
	}
	f.Self = f
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	f.RuneTypedCallback = f.runeTyped
	f.SetMaximum(max)
	f.Sync()
	return f
}

// SetMaximum sets the maximum value allowed.
func (f *NumericField) SetMaximum(maximum fixed.F64d4) {
	f.maximum = maximum
	if !f.noMinWidth && f.minimum != fixed.F64d4Min && f.maximum != fixed.F64d4Max {
		f.MinimumTextWidth = mathf32.Max(f.Font.SimpleWidth((f.minimum.Trunc() + fixed.F64d4One - 1).String()),
			f.Font.SimpleWidth((f.maximum.Trunc() + fixed.F64d4One - 1).String()))
	}
}

func (f *NumericField) trimmed(text string) string {
	return strings.TrimSpace(text)
}

func (f *NumericField) validate() bool {
	v, err := fixed.F64d4FromString(f.trimmed(f.Text()))
	if err != nil {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Invalid number"))
		return false
	}
	if f.minimum != fixed.F64d4Min && v < f.minimum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Number must be at least %s"), f.minimum.String()))
		return false
	}
	if f.maximum != fixed.F64d4Max && v > f.maximum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Number must be no more than %s"), f.maximum.String()))
		return false
	}
	f.Tooltip = nil
	return true
}

func (f *NumericField) modified() {
	if v, err := fixed.F64d4FromString(f.trimmed(f.Text())); err == nil &&
		(f.minimum == fixed.F64d4Min || v >= f.minimum) &&
		(f.maximum == fixed.F64d4Max || v <= f.maximum) {
		if f.get() != v {
			f.set(v)
			MarkForLayoutWithinDockable(f)
			MarkModified(f)
		}
	}
}

// Sync the field to the current value.
func (f *NumericField) Sync() {
	value := f.get()
	if f.minimum != fixed.F64d4Min && value < f.minimum {
		value = f.minimum
	} else if f.maximum != fixed.F64d4Max && value > f.maximum {
		value = f.maximum
	}
	SetFieldValue(f.Field, value.String())
}

func (f *NumericField) runeTyped(ch rune) bool {
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
