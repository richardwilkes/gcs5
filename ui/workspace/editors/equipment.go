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

package editors

import (
	"strconv"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
	"github.com/richardwilkes/unison"
)

// EditEquipment displays the editor for equipment.
func EditEquipment(owner widget.Rebuildable, equipment *gurps.Equipment, carried bool) {
	displayEditor[*gurps.Equipment, *gurps.EquipmentEditData](owner, equipment, func(e *editor[*gurps.Equipment, *gurps.EquipmentEditData], content *unison.Panel) func() {
		addNameLabelAndField(content, &e.editorData.Name)
		addNotesLabelAndField(content, &e.editorData.LocalNotes)
		addVTTNotesLabelAndField(content, &e.editorData.VTTNotes)
		if !e.target.Container() {
			qtyLabel := i18n.Text("Quantity")
			if carried {
				wrapper := addFlowWrapper(content, qtyLabel, 2)
				addNumericField(wrapper, qtyLabel, "", &e.editorData.Quantity, 0, f64d4.Max-1)
				addCheckBox(wrapper, i18n.Text("Equipped"), &e.editorData.Equipped)
			} else {
				addLabelAndNumericField(content, qtyLabel, "", &e.editorData.Quantity, 0, f64d4.Max-1)
			}
		}
		addLabelAndStringField(content, i18n.Text("Tech Level"), gurps.TechLevelInfo, &e.editorData.TechLevel)
		addLabelAndStringField(content, i18n.Text("Legality Class"), gurps.LegalityClassInfo, &e.editorData.LegalityClass)
		valueLabel := i18n.Text("Value")
		wrapper := addFlowWrapper(content, valueLabel, 3)
		addNumericField(wrapper, valueLabel, "", &e.editorData.Value, 0, f64d4.Max-1)
		wrapper.AddChild(widget.NewFieldInteriorLeadingLabel(i18n.Text("Extended")))
		wrapper.AddChild(widget.NewNonEditableField(func(field *widget.NonEditableField) {
			var value f64d4.Int
			if e.editorData.Quantity > 0 {
				value = gurps.ValueAdjustedForModifiers(e.editorData.Value, e.editorData.Modifiers).Mul(e.editorData.Quantity)
				if e.target.Container() {
					for _, one := range e.target.Children {
						value += one.ExtendedValue()
					}
				}
			}
			field.Text = value.Comma()
			field.MarkForLayoutAndRedraw()
		}))
		weightLabel := i18n.Text("Weight")
		wrapper = addFlowWrapper(content, weightLabel, 3)
		addWeightField(wrapper, weightLabel, "", e.target.Entity, &e.editorData.Weight)
		wrapper.AddChild(widget.NewFieldInteriorLeadingLabel(i18n.Text("Extended")))
		wrapper.AddChild(widget.NewNonEditableField(func(field *widget.NonEditableField) {
			var weight measure.Weight
			defUnits := gurps.SheetSettingsFor(e.target.Entity).DefaultWeightUnits
			if e.editorData.Quantity > 0 {
				weight = gurps.ExtendedWeightAdjustedForModifiers(defUnits, e.editorData.Quantity, e.editorData.Weight,
					e.editorData.Modifiers, e.editorData.Features, e.target.Children, false, false)
			}
			field.Text = defUnits.Format(weight)
			field.MarkForLayoutAndRedraw()
		}))
		content.AddChild(unison.NewPanel())
		addCheckBox(content, i18n.Text("Ignore weight for skills"), &e.editorData.WeightIgnoredForSkills)
		usesLabel := i18n.Text("Uses")
		wrapper = addFlowWrapper(content, usesLabel, 3)
		usesField := addIntegerField(wrapper, usesLabel, "", &e.editorData.Uses, 0, 9999999)
		maxUsesLabel := i18n.Text("Maximum Uses")
		wrapper.AddChild(widget.NewFieldInteriorLeadingLabel(maxUsesLabel))
		addIntegerField(wrapper, maxUsesLabel, "", &e.editorData.MaxUses, 0, 9999999)
		addTagsLabelAndField(content, &e.editorData.Tags)
		addPageRefLabelAndField(content, &e.editorData.PageRef)
		if e.editorData.MaxUses <= 0 {
			disableAndBlankField(usesField)
		}
		return func() {
			if e.editorData.Uses > e.editorData.MaxUses {
				usesField.SetText(strconv.Itoa(e.editorData.MaxUses))
			}
			if e.editorData.MaxUses > 0 {
				enableAndUnblankField(usesField)
			} else {
				disableAndBlankField(usesField)
			}
		}
	})
}
