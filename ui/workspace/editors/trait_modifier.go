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
	"github.com/richardwilkes/gcs/model/gurps/trait"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

// EditTraitModifier displays the editor for a trait modifier.
func EditTraitModifier(owner widget.Rebuildable, modifier *gurps.TraitModifier) {
	displayEditor[*gurps.TraitModifier, *gurps.TraitModifierEditData](owner, modifier,
		initTraitModifierEditor)
}

func initTraitModifierEditor(e *editor[*gurps.TraitModifier, *gurps.TraitModifierEditData], content *unison.Panel) func() {
	if !e.target.Container() {
		content.AddChild(unison.NewPanel())
		addInvertedCheckBox(content, i18n.Text("Enabled"), &e.editorData.Disabled)
	}
	addNameLabelAndField(content, &e.editorData.Name)
	addNotesLabelAndField(content, &e.editorData.LocalNotes)
	addVTTNotesLabelAndField(content, &e.editorData.VTTNotes)
	if !e.target.Container() {
		costLabel := i18n.Text("Cost")
		wrapper := addFlowWrapper(content, costLabel, 3)
		addDecimalField(wrapper, costLabel, "", &e.editorData.Cost, -fxp.MaxBasePoints, fxp.MaxBasePoints)
		costTypePopup := addCostTypePopup(wrapper, e)
		affectsPopup := addPopup(wrapper, trait.AllAffects, &e.editorData.Affects)
		levels := addLabelAndDecimalField(content, i18n.Text("Level"), "", &e.editorData.Levels, 0, fxp.Thousand)
		adjustFieldBlank(levels, !e.target.HasLevels())
		total := widget.NewNonEditableField(func(field *widget.NonEditableField) {
			enabled := true
			switch costTypePopup.SelectedIndex() - 1 {
			case -1:
				field.Text = e.editorData.Cost.Mul(e.editorData.Levels).StringWithSign() + trait.Percentage.String()
			case int(trait.Percentage):
				field.Text = e.editorData.Cost.StringWithSign() + trait.Percentage.String()
			case int(trait.Points):
				field.Text = e.editorData.Cost.StringWithSign()
			case int(trait.Multiplier):
				field.Text = trait.Multiplier.String() + e.editorData.Cost.String()
				affectsPopup.Select(trait.Total)
				enabled = false
			default:
				jot.Errorf("unhandled cost type popup index: %d", costTypePopup.SelectedIndex())
				field.Text = e.editorData.Cost.StringWithSign() + trait.Percentage.String()
			}
			affectsPopup.SetEnabled(enabled)
			field.MarkForLayoutAndRedraw()
		})
		insets := total.Border().Insets()
		total.SetLayoutData(&unison.FlexLayoutData{
			MinSize: unison.NewSize(total.Font.SimpleWidth((-fxp.MaxBasePoints*2).String())+insets.Left+insets.Right, 0),
		})
		content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Total")))
		content.AddChild(total)
		costTypePopup.SelectionCallback = func(index int, _ string) {
			if index == 0 {
				e.editorData.CostType = trait.Percentage
				if e.editorData.Levels < fxp.One {
					levels.SetText("1")
				}
			} else {
				e.editorData.CostType = trait.AllModifierCostType[index-1]
				e.editorData.Levels = 0
			}
			adjustFieldBlank(levels, index != 0)
			widget.MarkModified(wrapper)
		}
	}
	addTagsLabelAndField(content, &e.editorData.Tags)
	addPageRefLabelAndField(content, &e.editorData.PageRef)
	if !e.target.Container() {
		content.AddChild(newFeaturesPanel(e.target.Entity, e.target, &e.editorData.Features))
	}
	return nil
}

func addCostTypePopup(parent *unison.Panel, e *editor[*gurps.TraitModifier, *gurps.TraitModifierEditData]) *unison.PopupMenu[string] {
	popup := unison.NewPopupMenu[string]()
	popup.AddItem(i18n.Text("% per level"))
	for _, one := range trait.AllModifierCostType {
		popup.AddItem(one.String())
	}
	if e.target.HasLevels() {
		popup.SelectIndex(0)
	} else {
		popup.Select(e.editorData.CostType.String())
	}
	parent.AddChild(popup)
	return popup
}
