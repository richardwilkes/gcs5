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
	"strings"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// BodyPanel holds the contents of the body block on the sheet.
type BodyPanel struct {
	unison.Panel
	entity        *gurps.Entity
	row           []unison.Paneler
	sepLayoutData []*unison.FlexLayoutData
}

// NewBodyPanel creates a new body panel.
func NewBodyPanel(entity *gurps.Entity) *BodyPanel {
	p := &BodyPanel{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  6,
		HSpacing: 4,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		VSpan:  2,
	})
	locations := gurps.SheetSettingsFor(entity).HitLocations
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: locations.Name}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))
	p.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
		r := p.Children()[0].FrameRect()
		r.X = rect.X
		r.Width = rect.Width
		gc.DrawRect(r, theme.HeaderColor.Paint(gc, r, unison.Fill))
		for i, row := range p.row {
			var ink unison.Ink
			if i&1 == 1 {
				ink = unison.BandingColor
			} else {
				ink = unison.ContentColor
			}
			r = row.AsPanel().FrameRect()
			r.X = rect.X
			r.Width = rect.Width
			gc.DrawRect(r, ink.Paint(gc, r, unison.Fill))
		}
	}

	p.AddChild(widget.NewPageHeader(i18n.Text("Roll"), 1))
	p.AddChild(unison.NewPanel())
	p.AddChild(widget.NewPageHeader(i18n.Text("Location"), 2))
	p.AddChild(unison.NewPanel())
	p.AddChild(widget.NewPageHeader(i18n.Text("DR"), 1))

	p.row = nil
	p.sepLayoutData = nil
	p.addTable(locations, 0)

	for _, one := range p.sepLayoutData {
		one.VSpan = len(p.row)
	}

	return p
}

func (p *BodyPanel) addTable(bodyType *gurps.BodyType, depth int) {
	for i, location := range bodyType.Locations {
		prefix := strings.Repeat("   ", depth)
		p.AddChild(widget.NewPageLabelCenter(prefix + location.RollRange))

		if i == 0 {
			p.addSeparator()
		}

		name := widget.NewPageLabel(prefix + location.TableName)
		p.row = append(p.row, name)
		p.AddChild(name)
		p.AddChild(widget.NewNonEditablePageFieldEnd(fmt.Sprintf("%+d", location.HitPenalty), ""))

		if i == 0 {
			p.addSeparator()
		}

		var tooltip xio.ByteBuffer
		dr := location.DisplayDR(p.entity, &tooltip)
		p.AddChild(widget.NewNonEditablePageFieldCenter(dr,
			fmt.Sprintf(i18n.Text("The DR covering the %s hit location%s"), location.TableName, tooltip.String())))

		if location.SubTable != nil {
			p.addTable(location.SubTable, depth+1)
		}
	}
}

func (p *BodyPanel) addSeparator() {
	sep := unison.NewSeparator()
	sep.Vertical = true
	sep.LineInk = theme.HeaderColor
	layoutData := &unison.FlexLayoutData{
		HAlign: unison.MiddleAlignment,
		VAlign: unison.FillAlignment,
		VGrab:  true,
	}
	sep.SetLayoutData(layoutData)
	p.sepLayoutData = append(p.sepLayoutData, layoutData)
	p.AddChild(sep)
}
