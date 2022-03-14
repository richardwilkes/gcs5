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

// PercentageField holds the value for a percentage field.
type PercentageField struct {
	*unison.Field
	get           func() int
	set           func(int)
	minimum       int
	maximum       int
	marksModified bool
}

// NewPercentageField creates a new field that holds a percentage (where 100 == 100%).
func NewPercentageField(get func() int, set func(int), min, max int) *PercentageField {
	f := &PercentageField{
		Field:         unison.NewField(),
		get:           get,
		set:           set,
		minimum:       min,
		maximum:       max,
		marksModified: true,
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

// SetMarksModified sets whether this field will attempt to mark its ModifiableRoot as modified. Default is true.
func (f *PercentageField) SetMarksModified(marksModified bool) {
	f.marksModified = marksModified
}

// Set the field value.
func (f *PercentageField) Set(value int) {
	SetFieldValue(f.Field, f.formatted(value))
}

func (f *PercentageField) formatted(value int) string {
	return strconv.Itoa(value) + "%"
}

func (f *PercentageField) trimmed(text string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(text), "%"))
}

func (f *PercentageField) validate() bool {
	v, err := strconv.Atoi(f.trimmed(f.Text()))
	if err != nil {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Invalid percentage"))
		return false
	}
	if v < f.minimum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Percentage must be at least %s"), f.formatted(f.minimum)))
		return false
	}
	if v > f.maximum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Percentage must be no more than %s"), f.formatted(f.maximum)))
		return false
	}
	f.Tooltip = nil
	return true
}

func (f *PercentageField) modified() {
	if v, err := strconv.Atoi(f.trimmed(f.Text())); err == nil &&
		(f.minimum == math.MinInt || v >= f.minimum) &&
		(f.maximum == math.MaxInt || v <= f.maximum) {
		if f.get() != v {
			f.set(v)
			MarkForLayoutWithinDockable(f)
			if f.marksModified {
				MarkModified(f)
			}
		}
	}
}

// Sync the field to the current value.
func (f *PercentageField) Sync() {
	value := f.get()
	if f.minimum != math.MinInt && value < f.minimum {
		value = f.minimum
	} else if f.maximum != math.MaxInt && value > f.maximum {
		value = f.maximum
	}
	f.Set(value)
}

func (f *PercentageField) runeTyped(ch rune) bool {
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
