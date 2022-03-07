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

// PrimaryAttrPanel holds the contents of the primary attributes block on the sheet.
type PrimaryAttrPanel struct {
	unison.Panel
	entity *gurps.Entity
}

// NewPrimaryAttrPanel creates a new primary attributes panel.
func NewPrimaryAttrPanel(entity *gurps.Entity) *PrimaryAttrPanel {
	p := &PrimaryAttrPanel{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  3,
		HSpacing: 4,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
	})
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Primary Attributes")}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))
	p.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	}

	for _, def := range gurps.SheetSettingsFor(entity).Attributes.List() {
		if def.Type == attribute.Pool || !def.Primary() {
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

		field := widget.NewNumericPageField(attr.Current(), 0, attr.Maximum(), func(v fixed.F64d4) {
			// TODO: Implement
		})
		p.AddChild(field)

		p.AddChild(widget.NewPageLabel(def.CombinedName()))
	}

	return p
}
