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
	"reflect"

	"github.com/richardwilkes/gcs/model/criteria"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

type prereqPanel struct {
	unison.Panel
	entity *gurps.Entity
	root   **gurps.PrereqList
}

func newPrereqPanel(entity *gurps.Entity, root **gurps.PrereqList) *prereqPanel {
	p := &prereqPanel{
		entity: entity,
		root:   root,
	}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{Columns: 1})
	p.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  2,
		HAlign: unison.FillAlignment,
	})
	p.SetBorder(unison.NewCompoundBorder(
		&widget.TitledBorder{
			Title: i18n.Text("Prerequisites"),
			Font:  unison.LabelFont,
		},
		unison.NewEmptyBorder(unison.NewUniformInsets(2))))
	p.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	}
	p.AddChild(p.createPrereqListPanel(0, *root))
	return p
}

func (p *prereqPanel) createPrereqListPanel(depth int, list *gurps.PrereqList) *unison.Panel {
	panel := unison.NewPanel()
	p.createButtonsPanel(panel, depth, list)
	addNumericCriteriaPanel(panel, i18n.Text("When the Tech Level"), i18n.Text("When Tech Level"), &list.WhenTL, 0,
		fxp.Twelve, true, 1)
	addBoolPopup(panel, i18n.Text("requires all of:"), i18n.Text("requires at least one of:"), &list.All)
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HAlign:   unison.FillAlignment,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	for _, child := range list.Prereqs {
		var childPanel *unison.Panel
		switch one := child.(type) {
		case *gurps.PrereqList:
			childPanel = p.createPrereqListPanel(depth+1, one)
		case *gurps.AdvantagePrereq:
			childPanel = p.createAdvantagePrereqPanel(depth+1, one)
		case *gurps.AttributePrereq:
			childPanel = p.createAttributePrereqPanel(depth+1, one)
		case *gurps.ContainedQuantityPrereq:
			childPanel = p.createContainedQuantityPrereqPanel(depth+1, one)
		case *gurps.ContainedWeightPrereq:
			childPanel = p.createContainedWeightPrereqPanel(depth+1, one)
		case *gurps.SkillPrereq:
			childPanel = p.createSkillPrereqPanel(depth+1, one)
		case *gurps.SpellPrereq:
			childPanel = p.createSpellPrereqPanel(depth+1, one)
		default:
			jot.Warn(errs.Newf("unknown prerequisite type: %s", reflect.TypeOf(child).String()))
		}
		if childPanel != nil {
			childPanel.SetLayoutData(&unison.FlexLayoutData{
				HSpan:  columns,
				HAlign: unison.FillAlignment,
				HGrab:  true,
			})
			panel.AddChild(childPanel)
		}
	}
	return panel
}

func (p *prereqPanel) createButtonsPanel(parent *unison.Panel, depth int, prereq gurps.Prereq) {
	buttons := unison.NewPanel()
	buttons.SetBorder(unison.NewEmptyBorder(unison.Insets{Left: float32(depth * 20)}))
	parent.AddChild(buttons)
	if _, ok := prereq.(*gurps.PrereqList); ok {
		addPrereqButton := unison.NewSVGButton(res.CircledAddSVG)
		// TODO: Add button action
		buttons.AddChild(addPrereqButton)
		addPrereqListButton := unison.NewSVGButton(res.CircledVerticalElipsisSVG)
		// TODO: Add button action
		buttons.AddChild(addPrereqListButton)
	}
	parentList := prereq.ParentList()
	if parentList != nil {
		deleteButton := unison.NewSVGButton(res.TrashSVG)
		// TODO: Add button action
		buttons.AddChild(deleteButton)
		if parentList.Prereqs[0] != prereq {
			var text string
			if parentList.All {
				text = i18n.Text("and")
			} else {
				text = i18n.Text("or")
			}
			label := widget.NewFieldLeadingLabel(text)
			parent.AddChild(label)
		}
	}
	buttons.SetLayout(&unison.FlexLayout{
		Columns: len(buttons.Children()),
	})
}

func (p *prereqPanel) createAdvantagePrereqPanel(depth int, pr *gurps.AdvantagePrereq) *unison.Panel {
	panel := unison.NewPanel()
	p.createButtonsPanel(panel, depth, pr)
	addHasPopup(panel, &pr.Has)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	addNameCriteriaPanel(panel, &pr.NameCriteria, columns-1)
	addNotesCriteriaPanel(panel, &pr.NotesCriteria, columns-1)
	addLevelCriteriaPanel(panel, &pr.LevelCriteria, columns-1)
	return panel
}

func (p *prereqPanel) createAttributePrereqPanel(depth int, pr *gurps.AttributePrereq) *unison.Panel {
	panel := unison.NewPanel()
	p.createButtonsPanel(panel, depth, pr)
	addHasPopup(panel, &pr.Has)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	second := unison.NewPanel()
	second.SetLayoutData(&unison.FlexLayoutData{HSpan: columns - 1})
	choices := gurps.AttributeChoices(p.entity, "", true)
	var current *gurps.AttributeChoice
	for _, choice := range choices {
		if choice.Key == pr.Which {
			current = choice
		}
	}
	popup := addPopup[*gurps.AttributeChoice](second, choices, &current)
	popup.SelectionCallback = func(index int, _ *gurps.AttributeChoice) {
		pr.Which = choices[index].Key
		widget.MarkModified(second)
	}
	combinedWithChoices := append([]*gurps.AttributeChoice{{}},
		gurps.AttributeChoices(p.entity, i18n.Text("combined with"), true)...)
	var currentCombinedWith *gurps.AttributeChoice
	for _, choice := range choices {
		if choice.Key == pr.CombinedWith {
			currentCombinedWith = choice
		}
	}
	popup = addPopup[*gurps.AttributeChoice](second, combinedWithChoices, &currentCombinedWith)
	popup.SelectionCallback = func(index int, _ *gurps.AttributeChoice) {
		pr.CombinedWith = combinedWithChoices[index].Key
		widget.MarkModified(second)
	}
	var field unison.Paneler
	popupCriteria := unison.NewPopupMenu[string]()
	for _, one := range criteria.PrefixedNumericCompareTypeChoices(i18n.Text("which")) {
		popupCriteria.AddItem(one)
	}
	popupCriteria.SelectIndex(criteria.ExtractNumericCompareTypeIndex(string(pr.QualifierCriteria.Compare)))
	popupCriteria.SelectionCallback = func(index int, _ string) {
		pr.QualifierCriteria.Compare = criteria.AllNumericCompareTypes[index]
		if pr.QualifierCriteria.Compare == criteria.AnyNumber {
			disableAndBlankField(field)
		} else {
			enableAndUnblankField(field)
		}
		widget.MarkModified(second)
	}
	second.AddChild(popupCriteria)
	field = addNumericField(second, i18n.Text("Attribute Criteria"), "", &pr.QualifierCriteria.Qualifier, fxp.Min, fxp.Max)
	second.SetLayout(&unison.FlexLayout{
		Columns:  len(second.Children()),
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	panel.AddChild(unison.NewPanel())
	panel.AddChild(second)
	return panel
}

func (p *prereqPanel) createContainedQuantityPrereqPanel(depth int, pr *gurps.ContainedQuantityPrereq) *unison.Panel {
	panel := unison.NewPanel()
	p.createButtonsPanel(panel, depth, pr)
	addHasPopup(panel, &pr.Has)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	addQuantityCriteriaPanel(panel, &pr.QualifierCriteria)
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	return panel
}

func (p *prereqPanel) createContainedWeightPrereqPanel(depth int, pr *gurps.ContainedWeightPrereq) *unison.Panel {
	panel := unison.NewPanel()
	p.createButtonsPanel(panel, depth, pr)
	addHasPopup(panel, &pr.Has)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	second := unison.NewPanel()
	second.SetLayoutData(&unison.FlexLayoutData{HSpan: columns - 1})
	var field unison.Paneler
	popup := unison.NewPopupMenu[string]()
	for _, one := range criteria.PrefixedNumericCompareTypeChoices(i18n.Text("which")) {
		popup.AddItem(one)
	}
	popup.SelectIndex(criteria.ExtractNumericCompareTypeIndex(string(pr.WeightCriteria.Compare)))
	popup.SelectionCallback = func(index int, _ string) {
		pr.WeightCriteria.Compare = criteria.AllNumericCompareTypes[index]
		if pr.WeightCriteria.Compare == criteria.AnyNumber {
			disableAndBlankField(field)
		} else {
			enableAndUnblankField(field)
		}
		widget.MarkModified(second)
	}
	second.AddChild(popup)
	field = addWeightField(second, i18n.Text("Weight Criteria"), "", p.entity, &pr.WeightCriteria.Qualifier)
	second.SetLayout(&unison.FlexLayout{
		Columns:  len(second.Children()),
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	panel.AddChild(unison.NewPanel())
	panel.AddChild(second)
	return panel
}

func (p *prereqPanel) createSkillPrereqPanel(depth int, pr *gurps.SkillPrereq) *unison.Panel {
	panel := unison.NewPanel()
	p.createButtonsPanel(panel, depth, pr)
	addHasPopup(panel, &pr.Has)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	addNameCriteriaPanel(panel, &pr.NameCriteria, columns-1)
	addSpecializationCriteriaPanel(panel, &pr.SpecializationCriteria, columns-1)
	addLevelCriteriaPanel(panel, &pr.LevelCriteria, columns-1)
	return panel
}

func (p *prereqPanel) createSpellPrereqPanel(depth int, pr *gurps.SpellPrereq) *unison.Panel {
	panel := unison.NewPanel()
	p.createButtonsPanel(panel, depth, pr)
	addHasPopup(panel, &pr.Has)
	addQuantityCriteriaPanel(panel, &pr.QuantityCriteria)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	// TODO: Add other bits here
	return panel
}
