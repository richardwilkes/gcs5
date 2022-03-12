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
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/mathf32"
	"github.com/richardwilkes/unison"
)

// NumericField holds the data for a numeric field.
type NumericField struct {
	*unison.Field
	applier func(v fixed.F64d4)
	minimum fixed.F64d4
	maximum fixed.F64d4
}

// NewNumericField creates a new field that holds a fixed-point number.
func NewNumericField(value, min, max fixed.F64d4, noMinWidth bool, applier func(fixed.F64d4)) *NumericField {
	f := &NumericField{
		Field:   unison.NewField(),
		applier: applier,
		minimum: min,
		maximum: max,
	}
	f.Self = f
	f.SetText(value.String())
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	if !noMinWidth && min != fixed.F64d4Min && max != fixed.F64d4Max {
		f.MinimumTextWidth = mathf32.Max(f.Font.SimpleWidth((min.Trunc() + fixed.F64d4One - 1).String()),
			f.Font.SimpleWidth((max.Trunc() + fixed.F64d4One - 1).String()))
	}
	return f
}

func (f *NumericField) trimmed() string {
	return strings.TrimSpace(f.Text())
}

func (f *NumericField) validate() bool {
	v, err := fixed.F64d4FromString(f.trimmed())
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
	if v, err := fixed.F64d4FromString(f.trimmed()); err == nil &&
		(f.minimum == fixed.F64d4Min || v >= f.minimum) &&
		(f.maximum == fixed.F64d4Max || v <= f.maximum) {
		f.applier(v)
	}
}
