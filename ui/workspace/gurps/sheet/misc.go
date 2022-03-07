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
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// MiscPanel holds the contents of the miscellaneous block on the sheet.
type MiscPanel struct {
	unison.Panel
	entity        *gurps.Entity
	ModifiedField *unison.Label
	Modified      bool
}

// NewMiscPanel creates a new miscellaneous panel.
func NewMiscPanel(entity *gurps.Entity) *MiscPanel {
	m := &MiscPanel{entity: entity}
	m.Self = m
	m.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
	})
	m.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
	})
	m.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Miscellaneous")}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))
	m.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	}

	m.AddChild(widget.NewPageLabelEnd(i18n.Text("Created")))
	m.AddChild(widget.NewNonEditablePageField(entity.CreatedOn.String(), ""))

	m.AddChild(widget.NewPageLabelEnd(i18n.Text("Modified")))
	m.ModifiedField = widget.NewNonEditablePageField(entity.ModifiedOn.String(), "")
	m.AddChild(m.ModifiedField)

	m.AddChild(widget.NewPageLabelEnd(i18n.Text("Player")))
	m.AddChild(widget.NewStringPageFieldNoGrab(entity.Profile.PlayerName, func(v string) {
		entity.Profile.PlayerName = v
		MarkModified(m)
	}))

	return m
}

func (m *MiscPanel) updateModified() {
	m.Modified = true
	m.entity.ModifiedOn = jio.Now()
	text := m.entity.ModifiedOn.String()
	if text != m.ModifiedField.Text {
		m.ModifiedField.Text = text
		m.ModifiedField.MarkForRedraw()
		p := m.AsPanel()
		for p != nil {
			if d, ok := p.Self.(unison.Dockable); ok {
				if dc := unison.DockContainerFor(m); dc != nil {
					dc.UpdateTitle(d)
				}
				break
			}
			p = p.Parent()
		}
	}
}

// MarkModified updates the modification timestamp and marks the entity as modified.
func MarkModified(p unison.Paneler) {
	panel := p.AsPanel()
	for panel != nil {
		if um, ok := panel.Self.(*Sheet); ok {
			um.MiscPanel.updateModified()
			break
		}
		panel = panel.Parent()
	}
}

// SetTextAndMarkModified sets the field to the given text, selects it, requests focus, then calls MarkModified().
func SetTextAndMarkModified(field *unison.Field, text string) {
	field.SetText(text)
	field.SelectAll()
	field.RequestFocus()
	field.Parent().MarkForLayoutAndRedraw()
	MarkModified(field)
}
