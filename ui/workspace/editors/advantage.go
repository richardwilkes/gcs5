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
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/advantage"
	"github.com/richardwilkes/gcs/model/gurps/ancestry"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// EditAdvantage displays the editor for an advantage.
func EditAdvantage(owner widget.Rebuildable, advantage *gurps.Advantage) {
	displayEditor[*gurps.Advantage, *advantageEditorData](owner, advantage, initAdvantageEditor)
}

func initAdvantageEditor(e *editor[*gurps.Advantage, *advantageEditorData], content *unison.Panel) func() {
	content.AddChild(unison.NewPanel())
	addInvertedCheckBox(content, i18n.Text("Enabled"), &e.editorData.disabled)
	addNameLabelAndField(content, &e.editorData.name)
	addNotesLabelAndField(content, &e.editorData.notes)
	addVTTNotesLabelAndField(content, &e.editorData.vttNotes)
	addUserDescLabelAndField(content, &e.editorData.userDesc)
	addTagsLabelAndField(content, &e.editorData.tags)
	var levelField *widget.NumericField
	if !e.target.Container() {
		wrapper := addFlowWrapper(content, i18n.Text("Type"), 5)
		addCheckBox(wrapper, i18n.Text("Mental"), &e.editorData.mental)
		addCheckBox(wrapper, i18n.Text("Physical"), &e.editorData.physical)
		addCheckBox(wrapper, i18n.Text("Social"), &e.editorData.social)
		addCheckBox(wrapper, i18n.Text("Exotic"), &e.editorData.exotic)
		addCheckBox(wrapper, i18n.Text("Supernatural"), &e.editorData.supernatural)
		wrapper = addFlowWrapper(content, i18n.Text("Point Cost"), 8)
		pointCost := widget.NewNonEditableField(func(field *widget.NonEditableField) {
			field.Text = fxp.ApplyRounding(e.editorData.basePoints+e.editorData.levels.Mul(e.editorData.pointsPerLevel),
				e.editorData.roundCostDown).String()
			field.MarkForLayoutAndRedraw()
		})
		insets := pointCost.Border().Insets()
		pointCost.SetLayoutData(&unison.FlexLayoutData{
			MinSize: unison.NewSize(pointCost.Font.SimpleWidth((-fxp.MaxBasePoints*2).String())+insets.Left+insets.Right, 0),
		})
		wrapper.AddChild(pointCost)
		addCheckBox(wrapper, i18n.Text("Round Down"), &e.editorData.roundCostDown)
		baseCost := i18n.Text("Base Cost")
		wrapper = addFlowWrapper(content, baseCost, 8)
		addNumericField(wrapper, baseCost, "", &e.editorData.basePoints, -fxp.MaxBasePoints,
			fxp.MaxBasePoints)
		addLabelAndNumericField(wrapper, i18n.Text("Per Level"), "", &e.editorData.pointsPerLevel, -fxp.MaxBasePoints,
			fxp.MaxBasePoints)
		levelField = addLabelAndNumericField(wrapper, i18n.Text("Level"), "", &e.editorData.levels, 0, fxp.MaxBasePoints)
		if e.editorData.pointsPerLevel == 0 {
			disableAndBlankField(levelField)
		}
	}
	addLabelAndPopup(content, i18n.Text("Self-Control Roll"), "", advantage.AllSelfControlRolls, &e.editorData.cr)
	crAdjPopup := addLabelAndPopup(content, i18n.Text("CR Adjustment"), "", gurps.AllSelfControlRollAdj, &e.editorData.crAdj)
	if e.editorData.cr == advantage.None {
		crAdjPopup.SetEnabled(false)
	}
	var ancestryPopup *unison.PopupMenu[string]
	if e.target.Container() {
		addLabelAndPopup(content, i18n.Text("Container Type"), "", advantage.AllContainerType, &e.editorData.containerType)
		var choices []string
		for _, lib := range ancestry.AvailableAncestries(gurps.SettingsProvider.Libraries()) {
			for _, one := range lib.List {
				choices = append(choices, one.Name)
			}
		}
		ancestryPopup = addLabelAndPopup(content, i18n.Text("Ancestry"), "", choices, &e.editorData.ancestry)
		if e.editorData.containerType != advantage.Race {
			disableAndBlankPopup(ancestryPopup)
		}
	}
	addPageRefLabelAndField(content, &e.editorData.pageRef)
	return func() {
		if levelField != nil {
			if e.editorData.pointsPerLevel == 0 {
				disableAndBlankField(levelField)
			} else {
				enableAndUnblankField(levelField)
			}
		}
		if e.editorData.cr == advantage.None {
			crAdjPopup.SetEnabled(false)
			crAdjPopup.Select(gurps.NoCRAdj)
		} else {
			crAdjPopup.SetEnabled(true)
		}
		if ancestryPopup != nil {
			if e.editorData.containerType == advantage.Race {
				if !ancestryPopup.Enabled() {
					enableAndUnblankPopup(ancestryPopup)
					if ancestryPopup.IndexOfItem(e.editorData.ancestry) == -1 {
						e.editorData.ancestry = ancestry.Default
					}
					ancestryPopup.Select(e.editorData.ancestry)
				}
			} else {
				disableAndBlankPopup(ancestryPopup)
			}
		}
	}
}
