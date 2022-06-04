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

package editors

import (
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/unison"
)

type weaponsPanel struct {
	unison.Panel
	entity     *gurps.Entity
	weaponType weapon.Type
	allWeapons *[]*gurps.Weapon
	weapons    []*gurps.Weapon
}

func newWeaponsPanel(entity *gurps.Entity, weaponType weapon.Type, weapons *[]*gurps.Weapon) *weaponsPanel {
	p := &weaponsPanel{
		entity:     entity,
		weaponType: weaponType,
		allWeapons: weapons,
		weapons:    gurps.ExtractWeaponsOfType(weaponType, *weapons),
	}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{Columns: 1})
	p.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  2,
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.SetBorder(unison.NewLineBorder(theme.HeaderColor, 0, unison.NewUniformInsets(1), false))
	p.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	}
	newTable(p.AsPanel(), NewWeaponsProvider(p, p.weaponType, false))
	return p
}

func (p *weaponsPanel) Entity() *gurps.Entity {
	return p.entity
}

func (p *weaponsPanel) EquippedWeapons(_ weapon.Type) []*gurps.Weapon {
	return p.weapons
}
