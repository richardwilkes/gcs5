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

// ModificationUpdater defines the methods required to update the modified status.
type ModificationUpdater interface {
	UpdateModified()
}

// Misc holds the Miscellaneous panel on a sheet.
type Misc struct {
	unison.Panel
	entity        *gurps.Entity
	ModifiedField *unison.Label
	Modified      bool
}

// NewMisc creates a new Miscellaneous panel for a sheet.
func NewMisc(entity *gurps.Entity) *Misc {
	m := &Misc{entity: entity}
	m.Self = m
	m.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
	})
	m.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Miscellaneous")}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))

	m.AddChild(widget.NewPageLabel(i18n.Text("Created")))
	disabledField := widget.NewNonEditablePageField(entity.CreatedOn.String(), unison.StartAlignment)
	disabledField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	m.AddChild(disabledField)

	m.AddChild(widget.NewPageLabel(i18n.Text("Modified")))
	m.ModifiedField = widget.NewNonEditablePageField(entity.ModifiedOn.String(), unison.StartAlignment)
	m.ModifiedField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	m.AddChild(m.ModifiedField)

	m.AddChild(widget.NewPageLabel(i18n.Text("Player")))
	field := widget.NewStringPageField(entity.Profile.PlayerName, func(v string) {
		entity.Profile.PlayerName = v
		MarkModified(m)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	m.AddChild(field)

	return m
}

// UpdateModified implements ModificationUpdater
func (m *Misc) UpdateModified() {
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
		if um, ok := panel.Self.(ModificationUpdater); ok {
			um.UpdateModified()
			break
		}
		panel = panel.Parent()
	}
}
