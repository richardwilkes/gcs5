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

// IntegerField holds the data for an integer field.
type IntegerField struct {
	*unison.Field
	applier func(v int)
	minimum int
	maximum int
}

// NewIntegerField creates a new field that holds an integer.
func NewIntegerField(value, min, max int, applier func(int)) *IntegerField {
	f := &IntegerField{
		Field:   unison.NewField(),
		applier: applier,
		minimum: min,
		maximum: max,
	}
	f.Self = f
	f.SetText(strconv.Itoa(value))
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	if min != math.MinInt && max != math.MaxInt {
		f.MinimumTextWidth = mathf32.Max(f.Font.SimpleWidth(strconv.Itoa(min)), f.Font.SimpleWidth(strconv.Itoa(max)))
	}
	return f
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
		f.applier(v)
	}
}
