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
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// EditAdvantage displays the editor for an advantage.
func EditAdvantage(owner widget.Rebuildable, advantage *gurps.Advantage) {
	displayEditor[*gurps.Advantage, *gurps.AdvantageEditData](owner, advantage, initAdvantageEditor)
}

func initAdvantageEditor(e *editor[*gurps.Advantage, *gurps.AdvantageEditData], content *unison.Panel) func() {
	content.AddChild(unison.NewPanel())
	addInvertedCheckBox(content, i18n.Text("Enabled"), &e.editorData.Disabled)
	addNameLabelAndField(content, &e.editorData.Name)
	addNotesLabelAndField(content, &e.editorData.LocalNotes)
	addVTTNotesLabelAndField(content, &e.editorData.VTTNotes)
	addUserDescLabelAndField(content, &e.editorData.UserDesc)
	addTagsLabelAndField(content, &e.editorData.Tags)
	var levelField *widget.DecimalField
	if !e.target.Container() {
		wrapper := addFlowWrapper(content, i18n.Text("Point Cost"), 8)
		cost := widget.NewNonEditableField(func(field *widget.NonEditableField) {
			field.Text = gurps.AdjustedPoints(e.target.Entity, e.editorData.BasePoints, e.editorData.Levels,
				e.editorData.PointsPerLevel, e.editorData.CR, e.editorData.Modifiers,
				e.editorData.RoundCostDown).String()
			field.MarkForLayoutAndRedraw()
		})
		insets := cost.Border().Insets()
		cost.SetLayoutData(&unison.FlexLayoutData{
			MinSize: unison.NewSize(cost.Font.SimpleWidth((-fxp.MaxBasePoints*2).String())+insets.Left+insets.Right, 0),
		})
		wrapper.AddChild(cost)
		addCheckBox(wrapper, i18n.Text("Round Down"), &e.editorData.RoundCostDown)
		baseCost := i18n.Text("Base Cost")
		wrapper = addFlowWrapper(content, baseCost, 8)
		addDecimalField(wrapper, baseCost, "", &e.editorData.BasePoints, -fxp.MaxBasePoints,
			fxp.MaxBasePoints)
		addLabelAndDecimalField(wrapper, i18n.Text("Per Level"), "", &e.editorData.PointsPerLevel, -fxp.MaxBasePoints,
			fxp.MaxBasePoints)
		levelField = addLabelAndDecimalField(wrapper, i18n.Text("Level"), "", &e.editorData.Levels, 0,
			fxp.MaxBasePoints)
		adjustFieldBlank(levelField, e.editorData.PointsPerLevel == 0)
	}
	addLabelAndPopup(content, i18n.Text("Self-Control Roll"), "", advantage.AllSelfControlRolls, &e.editorData.CR)
	crAdjPopup := addLabelAndPopup(content, i18n.Text("CR Adjustment"), i18n.Text("Self-Control Roll Adjustment"),
		gurps.AllSelfControlRollAdj, &e.editorData.CRAdj)
	if e.editorData.CR == advantage.None {
		crAdjPopup.SetEnabled(false)
	}
	var ancestryPopup *unison.PopupMenu[string]
	if e.target.Container() {
		addLabelAndPopup(content, i18n.Text("Container Type"), "", advantage.AllContainerType,
			&e.editorData.ContainerType)
		var choices []string
		for _, lib := range ancestry.AvailableAncestries(gurps.SettingsProvider.Libraries()) {
			for _, one := range lib.List {
				choices = append(choices, one.Name)
			}
		}
		ancestryPopup = addLabelAndPopup(content, i18n.Text("Ancestry"), "", choices, &e.editorData.Ancestry)
		adjustPopupBlank(ancestryPopup, e.editorData.ContainerType != advantage.Race)
	}
	addPageRefLabelAndField(content, &e.editorData.PageRef)
	if e.target.Container() {
		content.AddChild(newAdvantageModifiersPanel(e.target.Entity, e.target, e.target.Modifiers))
	} else {
		content.AddChild(newPrereqPanel(e.target.Entity, &e.editorData.Prereq))
		content.AddChild(newFeaturesPanel(e.target.Entity, e.target, &e.editorData.Features))
		content.AddChild(newAdvantageModifiersPanel(e.target.Entity, e.target, e.target.Modifiers))
		content.AddChild(newMeleeWeaponsPanel(e.target.Entity, e.target,
			gurps.ExtractWeaponsOfType(weapon.Melee, e.target.Weapons)))
		content.AddChild(newRangedWeaponsPanel(e.target.Entity, e.target,
			gurps.ExtractWeaponsOfType(weapon.Ranged, e.target.Weapons)))
	}
	return func() {
		if levelField != nil {
			adjustFieldBlank(levelField, e.editorData.PointsPerLevel == 0)
		}
		if e.editorData.CR == advantage.None {
			crAdjPopup.SetEnabled(false)
			crAdjPopup.Select(gurps.NoCRAdj)
		} else {
			crAdjPopup.SetEnabled(true)
		}
		if ancestryPopup != nil {
			if e.editorData.ContainerType == advantage.Race {
				if !ancestryPopup.Enabled() {
					adjustPopupBlank(ancestryPopup, false)
					if ancestryPopup.IndexOfItem(e.editorData.Ancestry) == -1 {
						e.editorData.Ancestry = ancestry.Default
					}
					ancestryPopup.Select(e.editorData.Ancestry)
				}
			} else {
				adjustPopupBlank(ancestryPopup, true)
			}
		}
	}
}
