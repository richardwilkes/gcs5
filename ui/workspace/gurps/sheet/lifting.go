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
	p.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
		children := p.Children()
		for i := 0; i < len(children); i += 2 {
			var ink unison.Ink
			if (i/2)&1 == 1 {
				ink = unison.BandingColor
			} else {
				ink = unison.ContentColor
			}
			r := children[i].FrameRect()
			r.X = rect.X
			r.Width = rect.Width
			if i == 0 {
				r.Y--
				r.Height++
			} else if i == len(children)-2 {
				r.Height++
			}
			gc.DrawRect(r, ink.Paint(gc, r, unison.Fill))
		}
	}
	p.createRow(entity.BasicLift(), i18n.Text("Basic Lift"), i18n.Text("The weight that can be lifted overhead with one hand in one second"))
	p.createRow(entity.OneHandedLift(), i18n.Text("One-Handed Lift"), i18n.Text("The weight that can be lifted overhead with one hand in two seconds"))
	p.createRow(entity.TwoHandedLift(), i18n.Text("Two-Handed Lift"), i18n.Text("The weight that can be lifted overhead with both hands in four seconds"))
	p.createRow(entity.ShoveAndKnockOver(), i18n.Text("Shove & Knock Over"), i18n.Text("The weight of an object that can be shoved and knocked over"))
	p.createRow(entity.RunningShoveAndKnockOver(), i18n.Text("Running Shove & Knock Over"), i18n.Text("The weight of an object that can be shoved and knocked over with a running start"))
	p.createRow(entity.CarryOnBack(), i18n.Text("Carry On Back"), i18n.Text("The weight that can be carried slung across the back"))
	p.createRow(entity.ShiftSlightly(), i18n.Text("Shift Slightly"), i18n.Text("The weight that can be shifted slightly on a floor"))
	return p
}

func (p *LiftingPanel) createRow(weight measure.Weight, title, tooltip string) {
	p.AddChild(widget.NewNonEditablePageFieldEnd(weight.String(), tooltip))
	label := widget.NewPageLabel(title)
	label.Tooltip = unison.NewTooltipWithText(tooltip)
	p.AddChild(label)
}
