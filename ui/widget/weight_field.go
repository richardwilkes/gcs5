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

// WeightField holds the data for a weight field.
type WeightField struct {
	*unison.Field
	entity  *gurps.Entity
	applier func(v measure.Weight)
	maximum measure.Weight
}

// NewWeightField creates a new field that holds a weight.
func NewWeightField(entity *gurps.Entity, value, max measure.Weight, applier func(measure.Weight)) *WeightField {
	f := &WeightField{
		Field:   unison.NewField(),
		entity:  entity,
		applier: applier,
		maximum: max,
	}
	f.Self = f
	f.SetText(value.String())
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	units := gurps.SheetSettingsFor(f.entity).DefaultWeightUnits
	if f.maximum > 0 {
		f.MinimumTextWidth = mathf32.Max(f.Font.SimpleWidth(units.Format(0)), f.Font.SimpleWidth(units.Format(f.maximum)))
	}
	return f
}

func (f *WeightField) validate() bool {
	units := gurps.SheetSettingsFor(f.entity).DefaultWeightUnits
	v, err := measure.WeightFromString(f.Text(), units)
	if err != nil {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Invalid weight"))
		return false
	}
	if v < 0 {
		f.Tooltip = unison.NewTooltipWithText(i18n.Text("Weight may not be negative"))
		return false
	}
	if f.maximum > 0 && v > f.maximum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Weight must be no more than %s"), units.Format(f.maximum)))
		return false
	}
	f.Tooltip = nil
	return true
}

func (f *WeightField) modified() {
	units := gurps.SheetSettingsFor(f.entity).DefaultWeightUnits
	if v, err := measure.WeightFromString(f.Text(), units); err == nil && v >= 0 && v <= f.maximum {
		f.applier(v)
	}
}
