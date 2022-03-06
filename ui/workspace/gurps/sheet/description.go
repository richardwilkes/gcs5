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
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// DescriptionPanel holds the contents of the description block on the sheet.
type DescriptionPanel struct {
	unison.Panel
	entity *gurps.Entity
}

// NewDescriptionPanel creates a new description panel.
func NewDescriptionPanel(entity *gurps.Entity) *DescriptionPanel {
	d := &DescriptionPanel{entity: entity}
	d.Self = d
	d.SetLayout(&unison.FlexLayout{
		Columns:  6,
		HSpacing: 4,
	})
	d.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  2,
		HAlign: unison.FillAlignment,
		VAlign: unison.StartAlignment,
		HGrab:  true,
	})
	d.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("DescriptionPanel")}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("Gender")))
	d.AddChild(widget.NewStringPageField(entity.Profile.Gender, func(v string) {
		entity.Profile.Gender = v
		MarkModified(d)
	}))

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("Height")))
	d.AddChild(widget.NewHeightPageField(entity, entity.Profile.Height, 0, func(v measure.Length) {
		entity.Profile.Height = v
		MarkModified(d)
	}))

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("Hair")))
	d.AddChild(widget.NewStringPageField(entity.Profile.Hair, func(v string) {
		entity.Profile.Hair = v
		MarkModified(d)
	}))

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("Age")))
	d.AddChild(widget.NewStringPageField(entity.Profile.Age, func(v string) {
		entity.Profile.Age = v
		MarkModified(d)
	}))

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("Weight")))
	d.AddChild(widget.NewWeightPageField(entity, entity.Profile.Weight, 0, func(v measure.Weight) {
		entity.Profile.Weight = v
		MarkModified(d)
	}))

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("Eyes")))
	d.AddChild(widget.NewStringPageField(entity.Profile.Eyes, func(v string) {
		entity.Profile.Eyes = v
		MarkModified(d)
	}))

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("Birthday")))
	d.AddChild(widget.NewStringPageField(entity.Profile.Birthday, func(v string) {
		entity.Profile.Birthday = v
		MarkModified(d)
	}))

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("Size")))
	field := widget.NewSignedIntegerPageField(entity.Profile.AdjustedSizeModifier(), -99, 99, func(v int) {
		entity.Profile.SetAdjustedSizeModifier(v)
		MarkModified(d)
	})
	field.HAlign = unison.StartAlignment
	d.AddChild(field)

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("Skin")))
	d.AddChild(widget.NewStringPageField(entity.Profile.Skin, func(v string) {
		entity.Profile.Skin = v
		MarkModified(d)
	}))

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("Religion")))
	d.AddChild(widget.NewStringPageField(entity.Profile.Religion, func(v string) {
		entity.Profile.Religion = v
		MarkModified(d)
	}))

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("TL")))
	d.AddChild(widget.NewStringPageField(entity.Profile.TechLevel, func(v string) {
		entity.Profile.TechLevel = v
		MarkModified(d)
	}))

	d.AddChild(widget.NewPageLabelEnd(i18n.Text("Hand")))
	d.AddChild(widget.NewStringPageField(entity.Profile.Handedness, func(v string) {
		entity.Profile.Handedness = v
		MarkModified(d)
	}))
	return d
}
