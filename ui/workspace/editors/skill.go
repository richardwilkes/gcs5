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
	"strings"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
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
	isTechnique := strings.HasPrefix(e.target.Type, gid.Technique)
	if !e.target.Container() && !isTechnique {
		addSpecializationLabelAndField(content, &e.editorData.Specialization)
		addTechLevelRequired(content, &e.editorData.TechLevel, dockableKind == widget.SheetDockableKind)
	}
	addNotesLabelAndField(content, &e.editorData.LocalNotes)
	addVTTNotesLabelAndField(content, &e.editorData.VTTNotes)
	addTagsLabelAndField(content, &e.editorData.Tags)
	if !e.target.Container() {
		difficultyLabel := i18n.Text("Difficulty")
		choices := gurps.AttributeChoices(e.target.Entity, isTechnique)
		current := -1
		if isTechnique {
			wrapper := addFlowWrapper(content, i18n.Text("Defaults To"), 4)
			wrapper.SetLayoutData(&unison.FlexLayoutData{
				HAlign: unison.FillAlignment,
				HGrab:  true,
			})
			for i, one := range choices {
				if one.Key == e.editorData.TechniqueDefault.DefaultType {
					current = i
					break
				}
			}
			if current == -1 {
				current = len(choices)
				choices = append(choices, &gurps.AttributeChoice{
					Key:   e.editorData.TechniqueDefault.DefaultType,
					Title: e.editorData.TechniqueDefault.DefaultType,
				})
			}
			attrChoice := choices[current]
			attrChoicePopup := addPopup(wrapper, choices, &attrChoice)
			skillDefNameField := addStringField(wrapper, i18n.Text("Technique Default Skill Name"),
				i18n.Text("Skill Name"), &e.editorData.TechniqueDefault.Name)
			skillDefNameField.Watermark = i18n.Text("Skill")
			skillDefNameField.SetLayoutData(&unison.FlexLayoutData{
				HAlign: unison.FillAlignment,
				HGrab:  true,
			})
			skillDefSpecialtyField := addStringField(wrapper, i18n.Text("Technique Default Skill Specialization"),
				i18n.Text("Skill Specialization"), &e.editorData.TechniqueDefault.Specialization)
			skillDefSpecialtyField.Watermark = i18n.Text("Specialization")
			skillDefSpecialtyField.SetLayoutData(&unison.FlexLayoutData{
				HAlign: unison.FillAlignment,
				HGrab:  true,
			})
			lastWasSkillBased := skill.DefaultTypeIsSkillBased(e.editorData.TechniqueDefault.DefaultType)
			if !lastWasSkillBased {
				skillDefNameField.RemoveFromParent()
				skillDefSpecialtyField.RemoveFromParent()
			}
			addNumericField(wrapper, i18n.Text("Technique Default Adjustment"), i18n.Text("Default Adjustment"),
				&e.editorData.TechniqueDefault.Modifier, -fxp.NinetyNine, fxp.NinetyNine)
			attrChoicePopup.SelectionCallback = func(_ int, item *gurps.AttributeChoice) {
				e.editorData.TechniqueDefault.DefaultType = item.Key
				if skillBased := skill.DefaultTypeIsSkillBased(e.editorData.TechniqueDefault.DefaultType); skillBased != lastWasSkillBased {
					lastWasSkillBased = skillBased
					if skillBased {
						wrapper.AddChildAtIndex(skillDefNameField, len(wrapper.Children())-1)
						wrapper.AddChildAtIndex(skillDefSpecialtyField, len(wrapper.Children())-1)
					} else {
						skillDefNameField.RemoveFromParent()
						skillDefSpecialtyField.RemoveFromParent()
					}
				}
				widget.MarkModified(content)
			}
			wrapper2 := addFlowWrapper(content, "", 2)
			limitField := widget.NewNumericField(i18n.Text("Limit"), func() f64d4.Int {
				if e.editorData.TechniqueLimitModifier != nil {
					return *e.editorData.TechniqueLimitModifier
				}
				return 0
			}, func(value f64d4.Int) {
				if e.editorData.TechniqueLimitModifier != nil {
					*e.editorData.TechniqueLimitModifier = value
				}
				widget.MarkModified(wrapper2)
			}, -fxp.NinetyNine, fxp.NinetyNine, false)
			wrapper2.AddChild(widget.NewCheckBox(i18n.Text("Cannot exceed default skill level by more than"),
				e.editorData.TechniqueLimitModifier != nil, func(b bool) {
					if b {
						if e.editorData.TechniqueLimitModifier == nil {
							var limit f64d4.Int
							e.editorData.TechniqueLimitModifier = &limit
						}
						enableAndUnblankField(limitField)
					} else {
						e.editorData.TechniqueLimitModifier = nil
						disableAndBlankField(limitField)
					}
					widget.MarkModified(wrapper2)
				}))
			if e.editorData.TechniqueLimitModifier == nil {
				disableAndBlankField(limitField)
			}
			wrapper2.AddChild(limitField)
			addLabelAndPopup(content, difficultyLabel, "", skill.AllTechniqueDifficulty, &e.editorData.Difficulty.Difficulty)
		} else {
			wrapper := addFlowWrapper(content, difficultyLabel, 3)
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
		}

		if dockableKind == widget.SheetDockableKind || dockableKind == widget.TemplateDockableKind {
			pointsLabel := i18n.Text("Points")
			wrapper := addFlowWrapper(content, pointsLabel, 3)
			addNumericField(wrapper, pointsLabel, "", &e.editorData.Points, 0, fxp.MaxBasePoints)
			wrapper.AddChild(widget.NewFieldInteriorLeadingLabel(i18n.Text("Level")))
			levelField := widget.NewNonEditableField(func(field *widget.NonEditableField) {
				points := gurps.AdjustedPointsForNonContainerSkillOrTechnique(e.target.Entity, e.editorData.Points,
					e.editorData.Name, e.editorData.Specialization, e.editorData.Tags)
				var level skill.Level
				if e.target.Type == gid.Skill {
					level = gurps.CalculateSkillLevel(e.target.Entity, e.editorData.Name, e.editorData.Specialization,
						e.editorData.Tags, e.editorData.DefaultedFrom, e.editorData.Difficulty, points,
						e.editorData.EncumbrancePenaltyMultiplier)
				} else {
					level = gurps.CalculateTechniqueLevel(e.target.Entity, e.editorData.Name,
						e.editorData.Specialization, e.editorData.Tags, e.editorData.TechniqueDefault,
						e.editorData.Difficulty.Difficulty, points, true, e.editorData.TechniqueLimitModifier)
				}
				lvl := level.Level.Trunc()
				if lvl <= 0 {
					field.Text = "-"
				} else {
					rsl := level.RelativeLevel
					if isTechnique {
						rsl += e.editorData.TechniqueDefault.Modifier
					}
					field.Text = lvl.String() + "/" + gurps.FormatRelativeSkill(e.target.Entity, e.target.Type,
						e.editorData.Difficulty, rsl)
				}
				field.MarkForLayoutAndRedraw()
			})
			insets := levelField.Border().Insets()
			levelField.SetLayoutData(&unison.FlexLayoutData{
				MinSize: unison.NewSize(levelField.Font.SimpleWidth((-fxp.MaxBasePoints*2).String())+insets.Left+insets.Right, 0),
			})
			wrapper.AddChild(levelField)
		}
	}
	addPageRefLabelAndField(content, &e.editorData.PageRef)
	return nil
}
