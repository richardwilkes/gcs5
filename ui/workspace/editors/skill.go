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
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// EditSkill displays the editor for an skill.
func EditSkill(owner widget.Rebuildable, skill *gurps.Skill) {
	displayEditor[*gurps.Skill, *gurps.SkillEditData](owner, skill, initSkillEditor)
}

func initSkillEditor(e *editor[*gurps.Skill, *gurps.SkillEditData], content *unison.Panel) func() {
	var dockableKind string
	if one, ok := e.owner.(widget.DockableKind); ok {
		dockableKind = one.DockableKind()
	}
	addNameLabelAndField(content, &e.editorData.Name)
	if !e.target.Container() {
		addSpecializationLabelAndField(content, &e.editorData.Specialization)
		addTechLevelRequired(content, &e.editorData.TechLevel, dockableKind == widget.SheetDockableKind)
	}
	addNotesLabelAndField(content, &e.editorData.LocalNotes)
	addVTTNotesLabelAndField(content, &e.editorData.VTTNotes)
	addTagsLabelAndField(content, &e.editorData.Tags)
	if !e.target.Container() {
		wrapper := addFlowWrapper(content, i18n.Text("Difficulty"), 3)
		choices := gurps.AttributeChoices(e.target.Entity)
		current := -1
		for i, one := range choices {
			if one.Key == e.editorData.Difficulty.Attribute {
				current = i
				break
			}
		}
		if current == -1 {
			current = len(choices)
			choices = append(choices, &gurps.AttributeChoice{
				Key:   e.editorData.Difficulty.Attribute,
				Title: e.editorData.Difficulty.Attribute,
			})
		}
		attrChoice := choices[current]
		attrChoicePopup := addPopup(wrapper, choices, &attrChoice)
		attrChoicePopup.SelectionCallback = func(_ int, item *gurps.AttributeChoice) {
			e.editorData.Difficulty.Attribute = item.Key
			widget.MarkModified(content)
		}

		wrapper.AddChild(widget.NewFieldTrailingLabel("/"))
		addPopup(wrapper, skill.AllDifficulty, &e.editorData.Difficulty.Difficulty)
		encLabel := i18n.Text("Encumbrance Penalty")

		wrapper = addFlowWrapper(content, encLabel, 2)
		addNumericField(wrapper, encLabel, "", &e.editorData.EncumbrancePenaltyMultiplier, 0, fxp.Nine)
		wrapper.AddChild(widget.NewFieldTrailingLabel(i18n.Text("times the current encumbrance level")))

		if dockableKind == widget.SheetDockableKind || dockableKind == widget.TemplateDockableKind {
			addLabelAndNumericField(content, i18n.Text("Points"), "", &e.editorData.Points, 0, fxp.MaxBasePoints)
		}
	}
	addPageRefLabelAndField(content, &e.editorData.PageRef)
	return nil
}