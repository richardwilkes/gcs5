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
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// PointPoolsPanel holds the contents of the point pools block on the sheet.
type PointPoolsPanel struct {
	unison.Panel
	entity *gurps.Entity
}

// NewPointPoolsPanel creates a new point pools panel.
func NewPointPoolsPanel(entity *gurps.Entity) *PointPoolsPanel {
	p := &PointPoolsPanel{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  6,
		HSpacing: 4,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HSpan:  2,
	})
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Point Pools")}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))
	p.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	}

	// TODO: Need to CRC64 this so that we can swap out full data when attribute list changes
	for _, def := range gurps.SheetSettingsFor(entity).Attributes.List() {
		if def.Type != attribute.Pool {
			continue
		}
		attr, ok := entity.Attributes.Set[def.ID()]
		if !ok {
			jot.Warnf("unable to locate attribute data for '%s'", def.ID())
			continue
		}
		p.AddChild(p.createPointsField(attr))

		current := attr.Current()
		currentField := widget.NewNumericPageField(&current, fixed.F64d4Min, attr.Maximum(), true,
			func() { attr.Damage = attr.Maximum() - current })
		p.AddChild(currentField)

		p.AddChild(widget.NewPageLabel(i18n.Text("of")))

		maximum := attr.Maximum()
		maximumField := widget.NewNumericPageField(&maximum, fixed.F64d4Min, fixed.F64d4Max, true, func() {
			attr.SetMaximum(maximum)
			currentField.SetMaximum(maximum)
			currentField.SetValue(currentField.Value())
		})
		p.AddChild(maximumField)

		name := widget.NewPageLabel(def.Name)
		if def.FullName != "" {
			name.Tooltip = unison.NewTooltipWithText(def.FullName)
		}
		p.AddChild(name)

		if threshold := attr.CurrentThreshold(); threshold != nil {
			state := widget.NewPageLabel("[" + threshold.State + "]")
			if threshold.Explanation != "" {
				state.Tooltip = unison.NewTooltipWithText(threshold.Explanation)
			}
			p.AddChild(state)
		} else {
			p.AddChild(unison.NewPanel())
		}
	}

	return p
}

func (p *PointPoolsPanel) createPointsField(attr *gurps.Attribute) *widget.NonEditablePageField {
	field := widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := "[" + attr.PointCost().String() + "]"; text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
		if def := attr.AttributeDef(); def != nil {
			f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("Points spent on %s"), def.CombinedName()))
		}
	})
	field.Font = theme.PageFieldSecondaryFont
	return field
}
