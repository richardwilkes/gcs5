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
	"github.com/richardwilkes/gcs/model/undo"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/mathf32"
	"github.com/richardwilkes/unison"
)

// WeightField holds the value for a weight field.
type WeightField struct {
	*unison.Field
	undoID    int
	undoTitle string
	entity    *gurps.Entity
	get       func() measure.Weight
	set       func(measure.Weight)
	minimum   measure.Weight
	maximum   measure.Weight
	inUndo    bool
}

// NewWeightField creates a new field that holds a weight.
func NewWeightField(undoID int, undoTitle string, entity *gurps.Entity, get func() measure.Weight, set func(measure.Weight), min, max measure.Weight) *WeightField {
	f := &WeightField{
		Field:     unison.NewField(),
		undoID:    undoID,
		undoTitle: undoTitle,
		entity:    entity,
		get:       get,
		set:       set,
		minimum:   min,
		maximum:   max,
	}
	f.Self = f
	f.ModifiedCallback = f.modified
	f.ValidateCallback = f.validate
	units := gurps.SheetSettingsFor(f.entity).DefaultWeightUnits
	if min >= 0 && max > 0 {
		f.MinimumTextWidth = mathf32.Max(f.Font.SimpleWidth(units.Format(min)), f.Font.SimpleWidth(units.Format(max)))
	}
	f.Sync()
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
	if f.minimum >= 0 && v < f.minimum {
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Weight must be at least %s"), units.Format(f.minimum)))
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
	text := f.Text()
	if !f.inUndo && f.undoID != undo.NoneID {
		if mgr := undo.Manager(f); mgr != nil {
			mgr.Add(&undo.Edit[string]{
				ID:       f.undoID,
				EditName: f.undoTitle,
				EditCost: 1,
				UndoFunc: func(e *undo.Edit[string]) { f.setWithoutUndo(e.BeforeData, true) },
				RedoFunc: func(e *undo.Edit[string]) { f.setWithoutUndo(e.AfterData, true) },
				AbsorbFunc: func(e *undo.Edit[string], other unison.UndoEdit) bool {
					if e2, ok := other.(*undo.Edit[string]); ok && e2.ID == f.undoID {
						e.AfterData = e2.AfterData
						return true
					}
					return false
				},
				BeforeData: f.get().String(),
				AfterData:  text,
			})
		}
	}
	units := gurps.SheetSettingsFor(f.entity).DefaultWeightUnits
	if v, err := measure.WeightFromString(text, units); err == nil &&
		(f.minimum < 0 || v >= f.minimum) &&
		(f.maximum <= 0 || v <= f.maximum) && f.get() != v {
		f.set(v)
		MarkForLayoutWithinDockable(f)
		MarkModified(f)
	}
}

func (f *WeightField) setWithoutUndo(text string, focus bool) {
	f.inUndo = true
	f.SetText(text)
	f.inUndo = false
	if focus {
		f.RequestFocus()
		f.SelectAll()
	}
}

// Sync the field to the current value.
func (f *WeightField) Sync() {
	value := f.get()
	if f.minimum >= 0 && value < f.minimum {
		value = f.minimum
	} else if f.maximum > 0 && value > f.maximum {
		value = f.maximum
	}
	f.setWithoutUndo(gurps.SheetSettingsFor(f.entity).DefaultWeightUnits.Format(value), false)
}
