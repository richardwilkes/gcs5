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

// SecondaryAttrPanel holds the contents of the secondary attributes block on the sheet.
type SecondaryAttrPanel struct {
	unison.Panel
	entity *gurps.Entity
}

// NewSecondaryAttrPanel creates a new secondary attributes panel.
func NewSecondaryAttrPanel(entity *gurps.Entity) *SecondaryAttrPanel {
	p := &SecondaryAttrPanel{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  3,
		HSpacing: 4,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		VSpan:  2,
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
	})
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Secondary Attributes")}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))
	p.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	}

	for _, def := range gurps.SheetSettingsFor(entity).Attributes.List() {
		if def.Type == attribute.Pool || def.Primary() {
			continue
		}
		attr, ok := entity.Attributes.Set[def.ID()]
		if !ok {
			jot.Warnf("unable to locate attribute data for '%s'", def.ID())
			continue
		}
		pts := widget.NewNonEditablePageFieldEnd("["+attr.PointCost().String()+"]",
			fmt.Sprintf(i18n.Text("Points spent on %s"), def.CombinedName()))
		pts.Font = theme.PageFieldSecondaryFont
		p.AddChild(pts)

		field := widget.NewNumericPageField(attr.Current(), 0, attr.Maximum(), true, func(v fixed.F64d4) {
			// TODO: Implement
		})
		p.AddChild(field)

		p.AddChild(widget.NewPageLabel(def.CombinedName()))
	}

	return p
}

// Sync the panel to the current data.
func (p *SecondaryAttrPanel) Sync() {
	// TODO: Sync each attribute and points field
}
