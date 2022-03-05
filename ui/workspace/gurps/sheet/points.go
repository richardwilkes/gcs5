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
	"fmt"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// Points displays the character point totals and allows editing of the unspent point total.
type Points struct {
	unison.Panel
	entity *gurps.Entity
}

// NewPoints creates a new Points display.
func NewPoints(entity *gurps.Entity) *Points {
	p := &Points{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
	})
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: fmt.Sprintf(i18n.Text("%s Points"), p.entity.TotalPoints.String())}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))

	field := widget.NewNumericPageField(entity.UnspentPoints(), fixed.F64d4FromInt(-999999), fixed.F64d4FromInt(999999),
		entity.SetUnspentPoints)
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	field.Tooltip = unison.NewTooltipWithText(i18n.Text("Points earned but not yet spent"))
	p.AddChild(field)
	p.AddChild(widget.NewPageLabel(i18n.Text("Unspent")))

	ad, disad, race, quirk := entity.AdvantagePoints()
	disabledField := widget.NewNonEditablePageField(race.String(), unison.EndAlignment)
	disabledField.Tooltip = unison.NewTooltipWithText(i18n.Text("Total points spent on a racial package"))
	disabledField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	p.AddChild(disabledField)
	p.AddChild(widget.NewPageLabel(i18n.Text("Race")))

	disabledField = widget.NewNonEditablePageField(entity.AttributePoints().String(), unison.EndAlignment)
	disabledField.Tooltip = unison.NewTooltipWithText(i18n.Text("Total points spent on attributes"))
	disabledField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	p.AddChild(disabledField)
	p.AddChild(widget.NewPageLabel(i18n.Text("Attributes")))

	disabledField = widget.NewNonEditablePageField(ad.String(), unison.EndAlignment)
	disabledField.Tooltip = unison.NewTooltipWithText(i18n.Text("Total points spent on advantages"))
	disabledField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	p.AddChild(disabledField)
	p.AddChild(widget.NewPageLabel(i18n.Text("Advantages")))

	disabledField = widget.NewNonEditablePageField(disad.String(), unison.EndAlignment)
	disabledField.Tooltip = unison.NewTooltipWithText(i18n.Text("Total points spent on disadvantages"))
	disabledField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	p.AddChild(disabledField)
	p.AddChild(widget.NewPageLabel(i18n.Text("Disadvantages")))

	disabledField = widget.NewNonEditablePageField(quirk.String(), unison.EndAlignment)
	disabledField.Tooltip = unison.NewTooltipWithText(i18n.Text("Total points spent on quirks"))
	disabledField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	p.AddChild(disabledField)
	p.AddChild(widget.NewPageLabel(i18n.Text("Quirks")))

	disabledField = widget.NewNonEditablePageField(entity.SkillPoints().String(), unison.EndAlignment)
	disabledField.Tooltip = unison.NewTooltipWithText(i18n.Text("Total points spent on skills"))
	disabledField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	p.AddChild(disabledField)
	p.AddChild(widget.NewPageLabel(i18n.Text("Skills")))

	disabledField = widget.NewNonEditablePageField(entity.SpellPoints().String(), unison.EndAlignment)
	disabledField.Tooltip = unison.NewTooltipWithText(i18n.Text("Total points spent on spells"))
	disabledField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	p.AddChild(disabledField)
	p.AddChild(widget.NewPageLabel(i18n.Text("Spells")))

	return p
}
