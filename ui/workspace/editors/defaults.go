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
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
	"golang.org/x/exp/slices"
)

var lastDefaultTypeUsed = gid.Dexterity

type defaultsPanel struct {
	unison.Panel
	entity   *gurps.Entity
	defaults *[]*gurps.SkillDefault
}

func newDefaultsPanel(entity *gurps.Entity, defaults *[]*gurps.SkillDefault) *defaultsPanel {
	p := &defaultsPanel{
		entity:   entity,
		defaults: defaults,
	}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  2,
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.SetBorder(unison.NewCompoundBorder(
		&widget.TitledBorder{
			Title: i18n.Text("Defaults"),
			Font:  unison.LabelFont,
		},
		unison.NewEmptyBorder(unison.NewUniformInsets(2))))
	p.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	}
	addButton := unison.NewSVGButton(res.CircledAddSVG)
	addButton.ClickCallback = func() {
		def := &gurps.SkillDefault{DefaultType: lastDefaultTypeUsed}
		*defaults = slices.Insert(*defaults, 0, def)
		p.insertDefaultsPanel(1, def)
		unison.DockContainerFor(p).MarkForLayoutRecursively()
		widget.MarkModified(p)
	}
	p.AddChild(addButton)
	for i, one := range *defaults {
		p.insertDefaultsPanel(i+1, one)
	}
	return p
}

func (p *defaultsPanel) insertDefaultsPanel(index int, def *gurps.SkillDefault) {
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  5,
		HAlign:   unison.FillAlignment,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})

	deleteButton := unison.NewSVGButton(res.TrashSVG)
	deleteButton.ClickCallback = func() {
		if i := slices.IndexFunc(*p.defaults, func(elem *gurps.SkillDefault) bool { return elem == def }); i != -1 {
			*p.defaults = slices.Delete(*p.defaults, i, i+1)
		}
		panel.RemoveFromParent()
		unison.DockContainerFor(p).MarkForLayoutRecursively()
		widget.MarkModified(p)
	}
	panel.AddChild(deleteButton)

	name := i18n.Text("Name")
	nameField := widget.NewStringField(name, func() string { return def.Name },
		func(s string) { def.Name = s })
	nameField.Watermark = name
	specialization := i18n.Text("Specialization")
	specializationField := widget.NewStringField(specialization,
		func() string { return def.Specialization }, func(s string) { def.Specialization = s })
	specializationField.Watermark = specialization
	modifierField := widget.NewDecimalField(i18n.Text("Modifier"),
		func() fxp.Int { return def.Modifier },
		func(v fxp.Int) { def.Modifier = v },
		-fxp.Thousand, fxp.Thousand, true, false)
	attrChoicePopup := addAttributeChoicePopup(panel, p.entity, "", &def.DefaultType,
		gurps.TenFlag|gurps.ParryFlag|gurps.BlockFlag|gurps.SkillFlag)
	callback := attrChoicePopup.SelectionCallback
	attrChoicePopup.SelectionCallback = func(index int, item *gurps.AttributeChoice) {
		lastDefaultTypeUsed = item.Key
		callback(index, item)
		adjustFieldBlank(nameField, item.Key != gid.Skill)
		adjustFieldBlank(specializationField, item.Key != gid.Skill)
	}
	panel.AddChild(nameField)
	panel.AddChild(specializationField)
	panel.AddChild(modifierField)
	adjustFieldBlank(nameField, def.DefaultType != gid.Skill)
	adjustFieldBlank(specializationField, def.DefaultType != gid.Skill)

	panel.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.AddChildAtIndex(panel, index)
}
