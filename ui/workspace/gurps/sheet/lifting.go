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

// LiftingPanel holds the contents of the lifting block on the sheet.
type LiftingPanel struct {
	unison.Panel
	entity *gurps.Entity
}

// NewLiftingPanel creates a new lifting panel.
func NewLiftingPanel(entity *gurps.Entity) *LiftingPanel {
	p := &LiftingPanel{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
		HAlign:   unison.MiddleAlignment,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Lifting & Moving Things")}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))
	p.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) { drawBandedBackground(p, gc, rect, 0, 2) }
	p.addFieldAndLabel(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := p.entity.BasicLift().String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Basic Lift"), i18n.Text("The weight that can be lifted overhead with one hand in one second"))
	p.addFieldAndLabel(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := p.entity.OneHandedLift().String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("One-Handed Lift"), i18n.Text("The weight that can be lifted overhead with one hand in two seconds"))
	p.addFieldAndLabel(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := p.entity.TwoHandedLift().String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Two-Handed Lift"),
		i18n.Text("The weight that can be lifted overhead with both hands in four seconds"))
	p.addFieldAndLabel(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := p.entity.ShoveAndKnockOver().String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Shove & Knock Over"), i18n.Text("The weight of an object that can be shoved and knocked over"))
	p.addFieldAndLabel(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := p.entity.RunningShoveAndKnockOver().String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Running Shove & Knock Over"),
		i18n.Text("The weight of an object that can be shoved and knocked over with a running start"))
	p.addFieldAndLabel(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := p.entity.CarryOnBack().String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Carry On Back"), i18n.Text("The weight that can be carried slung across the back"))
	p.addFieldAndLabel(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := p.entity.ShiftSlightly().String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Shift Slightly"), i18n.Text("The weight that can be shifted slightly on a floor"))
	return p
}

func (p *LiftingPanel) addFieldAndLabel(field *widget.NonEditablePageField, title, tooltip string) {
	field.Tooltip = unison.NewTooltipWithText(tooltip)
	p.AddChild(field)
	label := widget.NewPageLabel(title)
	label.Tooltip = unison.NewTooltipWithText(tooltip)
	p.AddChild(label)
}
