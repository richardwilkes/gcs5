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
	"strconv"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/settings"
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
		Columns:  3,
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
	d.AddChild(d.createColumn1())
	d.AddChild(d.createColumn2())
	d.AddChild(d.createColumn3())
	return d
}

func createColumn() *unison.Panel {
	p := unison.NewPanel()
	p.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	return p
}

func (d *DescriptionPanel) createColumn1() *unison.Panel {
	column := createColumn()

	genderField := widget.NewStringPageField(d.entity.Profile.Gender, func(v string) {
		d.entity.Profile.Gender = v
		MarkModified(d)
	})
	column.AddChild(widget.NewPageLabelWithRandomizer(i18n.Text("Gender"),
		i18n.Text("Randomize the gender using the current ancestry"), func() {
			d.entity.Profile.Gender = d.entity.Ancestry().RandomGender(d.entity.Profile.Gender)
			SetTextAndMarkModified(genderField, d.entity.Profile.Gender)
		}))
	column.AddChild(genderField)

	ageField := widget.NewStringPageField(d.entity.Profile.Age, func(v string) {
		d.entity.Profile.Age = v
		MarkModified(d)
	})
	column.AddChild(widget.NewPageLabelWithRandomizer(i18n.Text("Age"),
		i18n.Text("Randomize the age using the current ancestry"), func() {
			age, _ := strconv.Atoi(d.entity.Profile.Age) //nolint:errcheck // A default of 0 is ok here on error
			d.entity.Profile.Age = strconv.Itoa(d.entity.Ancestry().RandomAge(d.entity, d.entity.Profile.Gender, age))
			SetTextAndMarkModified(ageField, d.entity.Profile.Age)
		}))
	column.AddChild(ageField)

	birthdayField := widget.NewStringPageField(d.entity.Profile.Birthday, func(v string) {
		d.entity.Profile.Birthday = v
		MarkModified(d)
	})
	column.AddChild(widget.NewPageLabelWithRandomizer(i18n.Text("Birthday"),
		i18n.Text("Randomize the birthday using the current calendar"), func() {
			global := settings.Global()
			d.entity.Profile.Birthday = global.General.CalendarRef(global.LibrarySet).RandomBirthday(d.entity.Profile.Birthday)
			SetTextAndMarkModified(birthdayField, d.entity.Profile.Birthday)
		}))
	column.AddChild(birthdayField)

	column.AddChild(widget.NewPageLabelEnd(i18n.Text("Religion")))
	column.AddChild(widget.NewStringPageField(d.entity.Profile.Religion, func(v string) {
		d.entity.Profile.Religion = v
		MarkModified(d)
	}))

	return column
}

func (d *DescriptionPanel) createColumn2() *unison.Panel {
	column := createColumn()

	heightField := widget.NewHeightPageField(d.entity, d.entity.Profile.Height, 0, func(v measure.Length) {
		d.entity.Profile.Height = v
		MarkModified(d)
	})
	column.AddChild(widget.NewPageLabelWithRandomizer(i18n.Text("Height"),
		i18n.Text("Randomize the height using the current ancestry"), func() {
			d.entity.Profile.Height = d.entity.Ancestry().RandomHeight(d.entity, d.entity.Profile.Gender, d.entity.Profile.Height)
			SetTextAndMarkModified(heightField.Field, d.entity.Profile.Height.String())
		}))
	column.AddChild(heightField)

	weightField := widget.NewWeightPageField(d.entity, d.entity.Profile.Weight, 0, func(v measure.Weight) {
		d.entity.Profile.Weight = v
		MarkModified(d)
	})
	column.AddChild(widget.NewPageLabelWithRandomizer(i18n.Text("Weight"),
		i18n.Text("Randomize the weight using the current ancestry"), func() {
			d.entity.Profile.Weight = d.entity.Ancestry().RandomWeight(d.entity, d.entity.Profile.Gender, d.entity.Profile.Weight)
			SetTextAndMarkModified(weightField.Field, d.entity.Profile.Weight.String())
		}))
	column.AddChild(weightField)

	column.AddChild(widget.NewPageLabelEnd(i18n.Text("Size")))
	field := widget.NewSignedIntegerPageField(d.entity.Profile.AdjustedSizeModifier(), -99, 99, func(v int) {
		d.entity.Profile.SetAdjustedSizeModifier(v)
		MarkModified(d)
	})
	field.HAlign = unison.StartAlignment
	column.AddChild(field)

	column.AddChild(widget.NewPageLabelEnd(i18n.Text("TL")))
	column.AddChild(widget.NewStringPageField(d.entity.Profile.TechLevel, func(v string) {
		d.entity.Profile.TechLevel = v
		MarkModified(d)
	}))

	return column
}

func (d *DescriptionPanel) createColumn3() *unison.Panel {
	column := createColumn()

	hairField := widget.NewStringPageField(d.entity.Profile.Hair, func(v string) {
		d.entity.Profile.Hair = v
		MarkModified(d)
	})
	column.AddChild(widget.NewPageLabelWithRandomizer(i18n.Text("Hair"),
		i18n.Text("Randomize the hair using the current ancestry"), func() {
			d.entity.Profile.Hair = d.entity.Ancestry().RandomHair(d.entity.Profile.Gender, d.entity.Profile.Hair)
			SetTextAndMarkModified(hairField, d.entity.Profile.Hair)
		}))
	column.AddChild(hairField)

	eyesField := widget.NewStringPageField(d.entity.Profile.Eyes, func(v string) {
		d.entity.Profile.Eyes = v
		MarkModified(d)
	})
	column.AddChild(widget.NewPageLabelWithRandomizer(i18n.Text("Eyes"),
		i18n.Text("Randomize the eyes using the current ancestry"), func() {
			d.entity.Profile.Eyes = d.entity.Ancestry().RandomEyes(d.entity.Profile.Gender, d.entity.Profile.Eyes)
			SetTextAndMarkModified(eyesField, d.entity.Profile.Eyes)
		}))
	column.AddChild(eyesField)

	skinField := widget.NewStringPageField(d.entity.Profile.Skin, func(v string) {
		d.entity.Profile.Skin = v
		MarkModified(d)
	})
	column.AddChild(widget.NewPageLabelWithRandomizer(i18n.Text("Skin"),
		i18n.Text("Randomize the skin using the current ancestry"), func() {
			d.entity.Profile.Skin = d.entity.Ancestry().RandomSkin(d.entity.Profile.Gender, d.entity.Profile.Skin)
			SetTextAndMarkModified(skinField, d.entity.Profile.Skin)
		}))
	column.AddChild(skinField)

	handField := widget.NewStringPageField(d.entity.Profile.Handedness, func(v string) {
		d.entity.Profile.Handedness = v
		MarkModified(d)
	})
	column.AddChild(widget.NewPageLabelWithRandomizer(i18n.Text("Hand"),
		i18n.Text("Randomize the handedness using the current ancestry"), func() {
			d.entity.Profile.Handedness = d.entity.Ancestry().RandomHandedness(d.entity.Profile.Gender, d.entity.Profile.Handedness)
			SetTextAndMarkModified(handField, d.entity.Profile.Handedness)
		}))
	column.AddChild(handField)

	return column
}
