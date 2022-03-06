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

package sheet

import (
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

type Description struct {
	unison.Panel
	entity *gurps.Entity
}

func NewDescription(entity *gurps.Entity) *Description {
	d := &Description{entity: entity}
	d.Self = d
	d.SetLayout(&unison.FlexLayout{
		Columns:  6,
		HSpacing: 4,
	})
	d.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Description")}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))

	d.AddChild(widget.NewPageLabel(i18n.Text("Gender")))
	field := widget.NewStringPageField(entity.Profile.Gender, func(v string) {
		entity.Profile.Gender = v
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)

	d.AddChild(widget.NewPageLabel(i18n.Text("Height")))
	field = widget.NewStringPageField(entity.Profile.Height.String(), func(v string) {
		// TODO: Implement length editor
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)

	d.AddChild(widget.NewPageLabel(i18n.Text("Hair")))
	field = widget.NewStringPageField(entity.Profile.Hair, func(v string) {
		entity.Profile.Hair = v
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)

	d.AddChild(widget.NewPageLabel(i18n.Text("Age")))
	field = widget.NewStringPageField(entity.Profile.Age, func(v string) {
		entity.Profile.Age = v
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)

	d.AddChild(widget.NewPageLabel(i18n.Text("Weight")))
	field = widget.NewStringPageField(entity.Profile.Weight.String(), func(v string) {
		// TODO: Implement weight editor
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)

	d.AddChild(widget.NewPageLabel(i18n.Text("Eyes")))
	field = widget.NewStringPageField(entity.Profile.Eyes, func(v string) {
		entity.Profile.Eyes = v
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)

	d.AddChild(widget.NewPageLabel(i18n.Text("Birthday")))
	field = widget.NewStringPageField(entity.Profile.Birthday, func(v string) {
		entity.Profile.Birthday = v
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)

	d.AddChild(widget.NewPageLabel(i18n.Text("Size")))
	field = widget.NewStringPageField(entity.Profile.AdjustedSizeModifier().String(), func(v string) {
		// TODO: Use different editor for numbers
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)

	d.AddChild(widget.NewPageLabel(i18n.Text("Skin")))
	field = widget.NewStringPageField(entity.Profile.Skin, func(v string) {
		entity.Profile.Skin = v
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)

	d.AddChild(widget.NewPageLabel(i18n.Text("Religion")))
	field = widget.NewStringPageField(entity.Profile.Religion, func(v string) {
		entity.Profile.Religion = v
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)

	d.AddChild(widget.NewPageLabel(i18n.Text("TL")))
	field = widget.NewStringPageField(entity.Profile.TechLevel, func(v string) {
		entity.Profile.TechLevel = v
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)

	d.AddChild(widget.NewPageLabel(i18n.Text("Hand")))
	field = widget.NewStringPageField(entity.Profile.Handedness, func(v string) {
		entity.Profile.Handedness = v
		MarkModified(d)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	d.AddChild(field)
	return d
}
