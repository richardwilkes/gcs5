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

// PercentageField holds the data for a percentage field.
type PercentageField struct {
	*unison.Field
	applier func(v fixed.F64d4)
	minimum fixed.F64d4
	maximum fixed.F64d4
}

// NewPercentageField creates a new field that holds a percentage (where 100 == 100%).
func NewPercentageField(value, min, max fixed.F64d4, applier func(fixed.F64d4)) *PercentageField {
	f := &PercentageField{
		Field:   unison.NewField(),
		applier: applier,
		minimum: min,
		maximum: max,
	}
	f.Self = f
	f.SetText(value.String() + "%")
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	f.MinimumTextWidth = mathf32.Max(f.Font.Extents((min.Trunc()+fixed.F64d4One-1).String()+"%").Width,
		f.Font.Extents((max.Trunc()+fixed.F64d4One-1).String()+"%").Width)
	return f
}

func (f *PercentageField) trimmed() string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(f.Text()), "%"))
}

func (f *PercentageField) validate() bool {
	v, err := fixed.F64d4FromString(f.trimmed())
	if err != nil {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Invalid percentage"))
		return false
	}
	if v < f.minimum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Percentage must be at least %s%%"), f.minimum.String()))
		return false
	}
	if v > f.maximum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Percentage must be no more than %s%%"), f.maximum.String()))
		return false
	}
	f.Tooltip = nil
	return true
}

func (f *PercentageField) modified() {
	if v, err := fixed.F64d4FromString(f.trimmed()); err == nil && v >= f.minimum && v <= f.maximum {
		f.applier(v)
	}
}
