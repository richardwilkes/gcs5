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
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/equipment"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// EditEquipmentModifier displays the editor for an equipment modifier.
func EditEquipmentModifier(owner widget.Rebuildable, modifier *gurps.EquipmentModifier) {
	displayEditor[*gurps.EquipmentModifier, *gurps.EquipmentModifierEditData](owner, modifier, initEquipmentModifierEditor)
}

func initEquipmentModifierEditor(e *editor[*gurps.EquipmentModifier, *gurps.EquipmentModifierEditData], content *unison.Panel) func() {
	if !e.target.Container() {
		content.AddChild(unison.NewPanel())
		addInvertedCheckBox(content, i18n.Text("Enabled"), &e.editorData.Disabled)
	}
	addNameLabelAndField(content, &e.editorData.Name)
	if !e.target.Container() {
		addLabelAndStringField(content, i18n.Text("Tech Level"), gurps.TechLevelInfo, &e.editorData.TechLevel)
	}
	addNotesLabelAndField(content, &e.editorData.LocalNotes)
	addVTTNotesLabelAndField(content, &e.editorData.VTTNotes)
	if !e.target.Container() {
		addEquipmentCostFields(content, e)
		addEquipmentWeightFields(content, e)
	}
	addTagsLabelAndField(content, &e.editorData.Tags)
	addPageRefLabelAndField(content, &e.editorData.PageRef)
	return nil
}

func addEquipmentCostFields(parent *unison.Panel, e *editor[*gurps.EquipmentModifier, *gurps.EquipmentModifierEditData]) {
	label := i18n.Text("Cost Modifier")
	wrapper := addFlowWrapper(parent, label, 2)
	first := true
	var field *widget.StringField
	field = widget.NewStringField(label,
		func() string {
			if first {
				first = false
				return e.editorData.CostType.Format(e.editorData.CostAmount)
			}
			return field.Text()
		},
		func(value string) {
			e.editorData.CostAmount = e.editorData.CostType.Format(value)
			widget.MarkModified(parent)
		})
	field.MinimumTextWidth = 50
	wrapper.AddChild(field)
	popup := unison.NewPopupMenu[string]()
	for _, one := range equipment.AllModifierCostType {
		popup.AddItem(one.StringWithExample())
	}
	popup.SelectIndex(int(e.editorData.CostType))
	wrapper.AddChild(popup)
	popup.SelectionCallback = func(index int, _ string) {
		e.editorData.CostType = equipment.AllModifierCostType[index]
		field.SetText(e.editorData.CostType.Format(field.Text()))
		widget.MarkModified(wrapper)
	}
}

func addEquipmentWeightFields(parent *unison.Panel, e *editor[*gurps.EquipmentModifier, *gurps.EquipmentModifierEditData]) {
	units := gurps.SheetSettingsFor(e.target.Entity).DefaultWeightUnits
	label := i18n.Text("Weight Modifier")
	wrapper := addFlowWrapper(parent, label, 2)
	first := true
	var field *widget.StringField
	field = widget.NewStringField(label,
		func() string {
			if first {
				first = false
				return e.editorData.WeightType.Format(e.editorData.WeightAmount, units)
			}
			return field.Text()
		},
		func(value string) {
			e.editorData.WeightAmount = e.editorData.WeightType.Format(value, units)
			widget.MarkModified(parent)
		})
	field.MinimumTextWidth = 50
	wrapper.AddChild(field)
	popup := unison.NewPopupMenu[string]()
	for _, one := range equipment.AllModifierWeightType {
		popup.AddItem(one.StringWithExample())
	}
	popup.SelectIndex(int(e.editorData.WeightType))
	wrapper.AddChild(popup)
	popup.SelectionCallback = func(index int, _ string) {
		e.editorData.WeightType = equipment.AllModifierWeightType[index]
		field.SetText(e.editorData.WeightType.Format(field.Text(), units))
		widget.MarkModified(wrapper)
	}
}
