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

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/mathf32"
	"github.com/richardwilkes/unison"
)

// PercentageField holds the data for a percentage field.
type PercentageField struct {
	*unison.Field
	applier func(v int)
	value   int
	minimum int
	maximum int
}

// NewPercentageField creates a new field that holds a percentage (where 100 == 100%).
func NewPercentageField(value, min, max int, applier func(int)) *PercentageField {
	f := &PercentageField{
		Field:   unison.NewField(),
		applier: applier,
		minimum: min,
		maximum: max,
	}
	f.Self = f
	f.SetValue(value)
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	f.RuneTypedCallback = f.runeTyped
	f.MinimumTextWidth = mathf32.Max(f.Font.SimpleWidth(strconv.Itoa(min)+"%"), f.Font.SimpleWidth(strconv.Itoa(max)+"%"))
	return f
}

// Value returns the current value of the field.
func (f *PercentageField) Value() int {
	return f.value
}

// SetValue sets the current value of the field. Will be limited to the minimum and maximum values.
func (f *PercentageField) SetValue(value int) {
	f.value = xmath.MinInt(xmath.MaxInt(value, f.minimum), f.maximum)
	f.SetText(strconv.Itoa(f.value) + "%")
}

func (f *PercentageField) validate() bool {
	v, err := strconv.Atoi(f.trimmed())
	if err != nil {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Invalid percentage"))
		return false
	}
	if v < f.minimum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Percentage must be at least %d%%"), f.minimum))
		return false
	}
	if v > f.maximum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Percentage must be no more than %d%%"), f.maximum))
		return false
	}
	f.Tooltip = nil
	return true
}

func (f *PercentageField) modified() {
	if v, err := strconv.Atoi(f.trimmed()); err == nil && v >= f.minimum && v <= f.maximum {
		f.value = v
		f.applier(v)
	}
}

func (f *PercentageField) trimmed() string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(f.Text()), "%"))
}

func (f *PercentageField) runeTyped(ch rune) bool {
	switch {
	case ch >= '0' && ch <= '9':
		return f.DefaultRuneTyped(ch)
	case ch == '-' || ch == '_':
		f.SetValue(f.Value() - 10)
		return true
	case ch == '=' || ch == '+':
		f.SetValue(f.Value() + 10)
		return true
	case ch == '%':
		if strings.Contains(f.SelectedText(), "%") || !strings.Contains(f.Text(), "%") {
			return f.DefaultRuneTyped(ch)
		}
		fallthrough
	default:
		unison.Beep()
		return false
	}
}
