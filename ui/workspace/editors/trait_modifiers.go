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
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/unison"
)

type traitModifiersPanel struct {
	unison.Panel
	entity    *gurps.Entity
	modifiers *[]*gurps.TraitModifier
	provider  widget.TableProvider[*Node[*gurps.TraitModifier]]
	table     *unison.Table[*Node[*gurps.TraitModifier]]
}

func newTraitModifiersPanel(entity *gurps.Entity, modifiers *[]*gurps.TraitModifier) *traitModifiersPanel {
	p := &traitModifiersPanel{
		entity:    entity,
		modifiers: modifiers,
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
	p.provider = NewTraitModifiersProvider(p, true)
	p.table = newTable(p.AsPanel(), p.provider)
	return p
}

func (p *traitModifiersPanel) Entity() *gurps.Entity {
	return p.entity
}

func (p *traitModifiersPanel) TraitModifierList() []*gurps.TraitModifier {
	return *p.modifiers
}

func (p *traitModifiersPanel) SetTraitModifierList(list []*gurps.TraitModifier) {
	*p.modifiers = list
	sel := RecordTableSelection(p.table)
	p.table.SyncToModel()
	ApplyTableSelection(p.table, sel)
}
