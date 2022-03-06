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

type Identity struct {
	unison.Panel
	entity *gurps.Entity
}

func NewIdentity(entity *gurps.Entity) *Identity {
	p := &Identity{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
	})
	p.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Identity")}, unison.NewEmptyBorder(geom32.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))

	p.AddChild(widget.NewPageLabel(i18n.Text("Name")))
	field := widget.NewStringPageField(entity.Profile.Name, func(v string) {
		entity.Profile.Name = v
		MarkModified(p)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	p.AddChild(field)

	p.AddChild(widget.NewPageLabel(i18n.Text("Title")))
	field = widget.NewStringPageField(entity.Profile.Title, func(v string) {
		entity.Profile.Title = v
		MarkModified(p)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	p.AddChild(field)

	p.AddChild(widget.NewPageLabel(i18n.Text("Organization")))
	field = widget.NewStringPageField(entity.Profile.Organization, func(v string) {
		entity.Profile.Organization = v
		MarkModified(p)
	})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})
	p.AddChild(field)
	return p
}
