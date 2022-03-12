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

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/mathf32"
	"github.com/richardwilkes/unison"
)

// HeightField holds the value for a height field.
type HeightField struct {
	*unison.Field
	entity  *gurps.Entity
	applier func(v measure.Length)
	maximum measure.Length
}

// NewHeightField creates a new field that holds a height.
func NewHeightField(entity *gurps.Entity, value, max measure.Length, applier func(measure.Length)) *HeightField {
	f := &HeightField{
		Field:   unison.NewField(),
		entity:  entity,
		applier: applier,
		maximum: max,
	}
	f.Self = f
	f.SetText(value.String())
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	units := gurps.SheetSettingsFor(f.entity).DefaultLengthUnits
	if f.maximum > 0 {
		f.MinimumTextWidth = mathf32.Max(f.Font.SimpleWidth(units.Format(0)), f.Font.SimpleWidth(units.Format(f.maximum)))
	}
	return f
}

func (f *HeightField) validate() bool {
	units := gurps.SheetSettingsFor(f.entity).DefaultLengthUnits
	v, err := measure.LengthFromString(f.Text(), units)
	if err != nil {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Invalid length"))
		return false
	}
	if v < 0 {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Length may not be negative"))
		return false
	}
	if f.maximum > 0 && v > f.maximum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Length must be no more than %s"), units.Format(f.maximum)))
		return false
	}
	f.Tooltip = nil
	return true
}

func (f *HeightField) modified() {
	units := gurps.SheetSettingsFor(f.entity).DefaultLengthUnits
	if v, err := measure.LengthFromString(f.Text(), units); err == nil && v >= 0 && v <= f.maximum {
		f.applier(v)
	}
}
