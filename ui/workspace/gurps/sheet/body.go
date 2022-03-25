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
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/unison"
)

// BodyPanel holds the contents of the body block on the sheet.
type BodyPanel struct {
	unison.Panel
	entity        *gurps.Entity
	row           []unison.Paneler
	sepLayoutData []*unison.FlexLayoutData
	crc           uint64
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
		VSpan:  3,
	})
	locations := gurps.SheetSettingsFor(entity).HitLocations
	p.crc = locations.CRC64()
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: locations.Name}, unison.NewEmptyBorder(geom.Insets[float32]{
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))
	p.DrawCallback = func(gc *unison.Canvas, rect geom.Rect[float32]) {
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
	p.addContent(locations)
	return p
}

func (p *BodyPanel) addContent(locations *gurps.BodyType) {
	p.RemoveAllChildren()
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
		p.AddChild(p.createHitPenaltyField(location))

		if i == 0 {
			p.addSeparator()
		}

		p.AddChild(p.createDRField(location))

		if location.SubTable != nil {
			p.addTable(location.SubTable, depth+1)
		}
	}
}

func (p *BodyPanel) createHitPenaltyField(location *gurps.HitLocation) unison.Paneler {
	return widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		f.Text = fmt.Sprintf("%+d", location.HitPenalty)
		widget.MarkForLayoutWithinDockable(f)
	})
}

func (p *BodyPanel) createDRField(location *gurps.HitLocation) unison.Paneler {
	return widget.NewNonEditablePageFieldCenter(func(f *widget.NonEditablePageField) {
		var tooltip xio.ByteBuffer
		f.Text = location.DisplayDR(p.entity, &tooltip)
		f.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("The DR covering the %s hit location%s"),
			location.TableName, tooltip.String()))
		widget.MarkForLayoutWithinDockable(f)
	})
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

// Sync the panel to the current data.
func (p *BodyPanel) Sync() {
	locations := gurps.SheetSettingsFor(p.entity).HitLocations
	if crc := locations.CRC64(); crc != p.crc {
		p.crc = crc
		p.addContent(locations)
		widget.MarkForLayoutWithinDockable(p)
	}
}
