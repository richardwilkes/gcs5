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

// IdentityPanel holds the contents of the identity block on the sheet.
type IdentityPanel struct {
	unison.Panel
	entity *gurps.Entity
}

// NewIdentityPanel creates a new identity panel.
func NewIdentityPanel(entity *gurps.Entity) *IdentityPanel {
	p := &IdentityPanel{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("IdentityPanel")}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))

	p.AddChild(widget.NewPageLabelEnd(i18n.Text("Name")))
	p.AddChild(widget.NewStringPageField(entity.Profile.Name, func(v string) {
		entity.Profile.Name = v
		MarkModified(p)
	}))

	p.AddChild(widget.NewPageLabelEnd(i18n.Text("Title")))
	p.AddChild(widget.NewStringPageField(entity.Profile.Title, func(v string) {
		entity.Profile.Title = v
		MarkModified(p)
	}))

	p.AddChild(widget.NewPageLabelEnd(i18n.Text("Organization")))
	p.AddChild(widget.NewStringPageField(entity.Profile.Organization, func(v string) {
		entity.Profile.Organization = v
		MarkModified(p)
	}))
	return p
}
