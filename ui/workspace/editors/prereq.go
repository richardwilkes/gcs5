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

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/gcs/model/gurps/spell"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
	"golang.org/x/exp/slices"
)

const noAndOr = ""

var lastPrereqTypeUsed = prereq.Advantage

type prereqPanel struct {
	unison.Panel
	entity   *gurps.Entity
	root     **gurps.PrereqList
	andOrMap map[gurps.Prereq]*unison.Label
}

func newPrereqPanel(entity *gurps.Entity, root **gurps.PrereqList) *prereqPanel {
	p := &prereqPanel{
		entity:   entity,
		root:     root,
		andOrMap: make(map[gurps.Prereq]*unison.Label),
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
	inFront := andOrText(list) != noAndOr
	if inFront {
		p.addAndOr(panel, list)
	}
	addNumericCriteriaPanel(panel, i18n.Text("When the Tech Level"), i18n.Text("When Tech Level"), &list.WhenTL, 0,
		fxp.Twelve, true, true, 1)
	popup := addBoolPopup(panel, i18n.Text("requires all of:"), i18n.Text("requires at least one of:"), &list.All)
	callback := popup.SelectionCallback
	popup.SelectionCallback = func(index int, item string) {
		callback(index, item)
		p.adjustAndOrForList(list)
	}
	if !inFront {
		p.addAndOr(panel, list)
	}
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HAlign:   unison.FillAlignment,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	for _, child := range list.Prereqs {
		p.addToList(panel, depth+1, child, false)
	}
	return panel
}

func (p *prereqPanel) addToList(parent *unison.Panel, depth int, child gurps.Prereq, first bool) {
	var panel *unison.Panel
	switch one := child.(type) {
	case *gurps.PrereqList:
		panel = p.createPrereqListPanel(depth, one)
	case *gurps.AdvantagePrereq:
		panel = p.createAdvantagePrereqPanel(depth, one)
	case *gurps.AttributePrereq:
		panel = p.createAttributePrereqPanel(depth, one)
	case *gurps.ContainedQuantityPrereq:
		panel = p.createContainedQuantityPrereqPanel(depth, one)
	case *gurps.ContainedWeightPrereq:
		panel = p.createContainedWeightPrereqPanel(depth, one)
	case *gurps.SkillPrereq:
		panel = p.createSkillPrereqPanel(depth, one)
	case *gurps.SpellPrereq:
		panel = p.createSpellPrereqPanel(depth, one)
	default:
		jot.Warn(errs.Newf("unknown prerequisite type: %s", reflect.TypeOf(child).String()))
	}
	if panel != nil {
		columns := parent.Layout().(*unison.FlexLayout).Columns
		panel.SetLayoutData(&unison.FlexLayoutData{
			HSpan:  columns,
			HAlign: unison.FillAlignment,
			HGrab:  true,
		})
		if first {
			parent.AddChildAtIndex(panel, columns)
		} else {
			parent.AddChild(panel)
		}
	}
}

func (p *prereqPanel) createButtonsPanel(parent *unison.Panel, depth int, data gurps.Prereq) {
	buttons := unison.NewPanel()
	buttons.SetBorder(unison.NewEmptyBorder(unison.Insets{Left: float32(depth * 20)}))
	parent.AddChild(buttons)
	if prereqList, ok := data.(*gurps.PrereqList); ok {
		addPrereqButton := unison.NewSVGButton(res.CircledAddSVG)
		addPrereqButton.ClickCallback = func() {
			var created gurps.Prereq
			switch lastPrereqTypeUsed {
			case prereq.Attribute:
				one := gurps.NewAttributePrereq(p.entity)
				one.Parent = prereqList
				created = one
			case prereq.ContainedQuantity:
				one := gurps.NewContainedQuantityPrereq()
				one.Parent = prereqList
				created = one
			case prereq.ContainedWeight:
				one := gurps.NewContainedWeightPrereq(p.entity)
				one.Parent = prereqList
				created = one
			case prereq.Skill:
				one := gurps.NewSkillPrereq()
				one.Parent = prereqList
				created = one
			case prereq.Spell:
				one := gurps.NewSpellPrereq()
				one.Parent = prereqList
				created = one
			default: // prereq.Advantage
				one := gurps.NewAdvantagePrereq()
				one.Parent = prereqList
				created = one
			}
			prereqList.Prereqs = slices.Insert(prereqList.Prereqs, 0, created)
			p.addToList(parent, depth+1, created, true)
			p.adjustAndOrForList(prereqList)
			unison.DockContainerFor(p).MarkForLayoutRecursively()
			widget.MarkModified(p)
		}
		buttons.AddChild(addPrereqButton)

		addPrereqListButton := unison.NewSVGButton(res.CircledVerticalElipsisSVG)
		addPrereqListButton.ClickCallback = func() {
			newList := gurps.NewPrereqList()
			newList.Parent = prereqList
			prereqList.Prereqs = slices.Insert(prereqList.Prereqs, 0, gurps.Prereq(newList))
			p.addToList(parent, depth+1, newList, true)
			p.adjustAndOrForList(prereqList)
			unison.DockContainerFor(p).MarkForLayoutRecursively()
			widget.MarkModified(p)
		}
		buttons.AddChild(addPrereqListButton)
	}
	parentList := data.ParentList()
	if parentList != nil {
		deleteButton := unison.NewSVGButton(res.TrashSVG)
		deleteButton.ClickCallback = func() {
			delete(p.andOrMap, data)
			if i := slices.IndexFunc(parentList.Prereqs, func(elem gurps.Prereq) bool { return elem == data }); i != -1 {
				parentList.Prereqs = slices.Delete(parentList.Prereqs, i, i+1)
			}
			parent.RemoveFromParent()
			p.adjustAndOrForList(parentList)
			unison.DockContainerFor(p).MarkForLayoutRecursively()
			widget.MarkModified(p)
		}
		buttons.AddChild(deleteButton)
	}
	buttons.SetLayout(&unison.FlexLayout{
		Columns: len(buttons.Children()),
	})
}

func (p *prereqPanel) addAndOr(parent *unison.Panel, data gurps.Prereq) {
	label := widget.NewFieldLeadingLabel(andOrText(data))
	parent.AddChild(label)
	p.andOrMap[data] = label
}

func (p *prereqPanel) adjustAndOrForList(list *gurps.PrereqList) {
	for _, one := range list.Prereqs {
		p.adjustAndOr(one)
	}
	p.MarkForLayoutRecursively()
}

func (p *prereqPanel) adjustAndOr(data gurps.Prereq) {
	if label, ok := p.andOrMap[data]; ok {
		if text := andOrText(data); text != label.Text {
			parent := label.Parent()
			label.RemoveFromParent()
			label.Text = text
			i := 1
			if text == noAndOr {
				i = parent.Layout().(*unison.FlexLayout).Columns - 1
			}
			parent.AddChildAtIndex(label, i)
		}
	}
}

func andOrText(pr gurps.Prereq) string {
	list := pr.ParentList()
	if list == nil || len(list.Prereqs) < 2 || list.Prereqs[0] == pr {
		return noAndOr
	}
	if list.All {
		return i18n.Text("and")
	}
	return i18n.Text("or")
}

func (p *prereqPanel) createAdvantagePrereqPanel(depth int, pr *gurps.AdvantagePrereq) *unison.Panel {
	panel := unison.NewPanel()
	p.createButtonsPanel(panel, depth, pr)
	inFront := andOrText(pr) != noAndOr
	if inFront {
		p.addAndOr(panel, pr)
	}
	addHasPopup(panel, &pr.Has)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	if !inFront {
		p.addAndOr(panel, pr)
	}
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	addNameCriteriaPanel(panel, &pr.NameCriteria, columns-1)
	addNotesCriteriaPanel(panel, &pr.NotesCriteria, columns-1)
	addLevelCriteriaPanel(panel, &pr.LevelCriteria, columns-1, true)
	return panel
}

func (p *prereqPanel) createAttributePrereqPanel(depth int, pr *gurps.AttributePrereq) *unison.Panel {
	panel := unison.NewPanel()
	p.createButtonsPanel(panel, depth, pr)
	inFront := andOrText(pr) != noAndOr
	if inFront {
		p.addAndOr(panel, pr)
	}
	addHasPopup(panel, &pr.Has)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	if !inFront {
		p.addAndOr(panel, pr)
	}
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	second := unison.NewPanel()
	second.SetLayoutData(&unison.FlexLayoutData{HSpan: columns - 1})
	addAttributeChoicePopup(second, p.entity, noAndOr, &pr.Which, false)
	addAttributeChoicePopup(second, p.entity, i18n.Text("combined with"), &pr.CombinedWith, true)
	addNumericCriteriaPanel(second, i18n.Text("which"), i18n.Text("Attribute Qualifier"), &pr.QualifierCriteria,
		fxp.Min, fxp.Max, false, false, 1)
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
	inFront := andOrText(pr) != noAndOr
	if inFront {
		p.addAndOr(panel, pr)
	}
	addHasPopup(panel, &pr.Has)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	addQuantityCriteriaPanel(panel, &pr.QualifierCriteria)
	if !inFront {
		p.addAndOr(panel, pr)
	}
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
	inFront := andOrText(pr) != noAndOr
	if inFront {
		p.addAndOr(panel, pr)
	}
	addHasPopup(panel, &pr.Has)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	if !inFront {
		p.addAndOr(panel, pr)
	}
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	second := unison.NewPanel()
	second.SetLayoutData(&unison.FlexLayoutData{HSpan: columns - 1})
	addWeightCriteriaPanel(second, p.entity, &pr.WeightCriteria)
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
	inFront := andOrText(pr) != noAndOr
	if inFront {
		p.addAndOr(panel, pr)
	}
	addHasPopup(panel, &pr.Has)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	if !inFront {
		p.addAndOr(panel, pr)
	}
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	addNameCriteriaPanel(panel, &pr.NameCriteria, columns-1)
	addSpecializationCriteriaPanel(panel, &pr.SpecializationCriteria, columns-1)
	addLevelCriteriaPanel(panel, &pr.LevelCriteria, columns-1, true)
	return panel
}

func (p *prereqPanel) createSpellPrereqPanel(depth int, pr *gurps.SpellPrereq) *unison.Panel {
	panel := unison.NewPanel()
	p.createButtonsPanel(panel, depth, pr)
	inFront := andOrText(pr) != noAndOr
	if inFront {
		p.addAndOr(panel, pr)
	}
	addHasPopup(panel, &pr.Has)
	addQuantityCriteriaPanel(panel, &pr.QuantityCriteria)
	addPopup[prereq.Type](panel, prereq.AllType[1:], &pr.Type)
	if !inFront {
		p.addAndOr(panel, pr)
	}
	columns := len(panel.Children())
	panel.SetLayout(&unison.FlexLayout{
		Columns:  columns,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	second := unison.NewPanel()
	second.SetLayoutData(&unison.FlexLayoutData{HSpan: columns - 1})
	subTypePopup := addPopup[spell.ComparisonType](second, spell.AllComparisonType, &pr.SubType)
	popup, field := addStringCriteriaPanel(second, "", i18n.Text("Spell Qualifier"), &pr.QualifierCriteria, 1)
	savedCallback := subTypePopup.SelectionCallback
	subTypePopup.SelectionCallback = func(index int, item spell.ComparisonType) {
		savedCallback(index, item)
		if pr.SubType == spell.Any || pr.SubType == spell.CollegeCount {
			disableAndBlankPopup(popup)
			disableAndBlankField(field)
		} else {
			enableAndUnblankPopup(popup)
			enableAndUnblankField(field)
		}
	}
	if pr.SubType == spell.Any || pr.SubType == spell.CollegeCount {
		disableAndBlankPopup(popup)
		disableAndBlankField(field)
	}
	second.SetLayout(&unison.FlexLayout{
		Columns:  len(second.Children()),
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	panel.AddChild(unison.NewPanel())
	panel.AddChild(second)
	return panel
}
