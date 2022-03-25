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
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/unison"
)

// PointsPanel holds the contents of the points block on the sheet.
type PointsPanel struct {
	unison.Panel
	entity       *gurps.Entity
	pointsBorder *TitledBorder
	unspent      *widget.NumericField
}

// NewPointsPanel creates a new points panel.
func NewPointsPanel(entity *gurps.Entity) *PointsPanel {
	p := &PointsPanel{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.EndAlignment,
		VAlign: unison.FillAlignment,
		VSpan:  2,
	})
	p.pointsBorder = &TitledBorder{Title: fmt.Sprintf(i18n.Text("%s Points"), p.entity.TotalPoints.String())}
	p.SetBorder(unison.NewCompoundBorder(p.pointsBorder, unison.NewEmptyBorder(geom.Insets[float32]{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))
	p.DrawCallback = func(gc *unison.Canvas, rect geom.Rect[float32]) { drawBandedBackground(p, gc, rect, 0, 2) }

	p.unspent = widget.NewNumericPageField(func() fixed.F64d4 { return p.entity.UnspentPoints() },
		func(v fixed.F64d4) { p.entity.SetUnspentPoints(v) }, fixed.F64d4Min, fixed.F64d4Max, true)
	p.unspent.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	p.unspent.Tooltip = unison.NewTooltipWithText(i18n.Text("Points earned but not yet spent"))
	p.AddChild(p.unspent)
	p.AddChild(widget.NewPageLabel(i18n.Text("Unspent")))
	p.addPointsField(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		_, _, race, _ := p.entity.AdvantagePoints()
		if text := race.String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Race"), i18n.Text("Total points spent on a racial package"))
	p.addPointsField(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := p.entity.AttributePoints().String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Attributes"), i18n.Text("Total points spent on attributes"))
	p.addPointsField(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		ad, _, _, _ := p.entity.AdvantagePoints()
		if text := ad.String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Advantages"), i18n.Text("Total points spent on advantages"))
	p.addPointsField(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		_, disad, _, _ := p.entity.AdvantagePoints()
		if text := disad.String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Disadvantages"), i18n.Text("Total points spent on disadvantages"))
	p.addPointsField(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		_, _, _, quirk := p.entity.AdvantagePoints()
		if text := quirk.String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Quirks"), i18n.Text("Total points spent on quirks"))
	p.addPointsField(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := p.entity.SkillPoints().String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Skills"), i18n.Text("Total points spent on skills"))
	p.addPointsField(widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := p.entity.SpellPoints().String(); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	}), i18n.Text("Spells"), i18n.Text("Total points spent on spells"))

	return p
}

func (p *PointsPanel) addPointsField(field *widget.NonEditablePageField, title, tooltip string) {
	field.Tooltip = unison.NewTooltipWithText(tooltip)
	p.AddChild(field)
	label := widget.NewPageLabel(title)
	label.Tooltip = unison.NewTooltipWithText(tooltip)
	p.AddChild(label)
}

// Sync the panel to the current data.
func (p *PointsPanel) Sync() {
	p.unspent.Sync()
	p.pointsBorder.Title = fmt.Sprintf(i18n.Text("%s Points"), p.entity.TotalPoints.String())
}
