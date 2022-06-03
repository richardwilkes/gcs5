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
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type advantageModifiersPanel struct {
	unison.Panel
	owner *gurps.Advantage
}

func newAdvantageModifiersPanel(owner *gurps.Advantage) *advantageModifiersPanel {
	p := &advantageModifiersPanel{
		owner: owner,
	}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  2,
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.SetBorder(unison.NewCompoundBorder(
		&widget.TitledBorder{
			Title: i18n.Text("Modifiers"),
			Font:  unison.LabelFont,
		},
		unison.NewEmptyBorder(unison.NewUniformInsets(2))))
	p.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	}

	// TODO: Implement
	addEditorNotYetImplementedBlock(p)

	return p
}

func (p *advantageModifiersPanel) createTable() {
	table := unison.NewTable()
	p.AddChild(table)
}

func (p *advantageModifiersPanel) Entity() *gurps.Entity {
	return p.owner.Entity
}

func (p *advantageModifiersPanel) AdvantageModifierList() []*gurps.AdvantageModifier {
	return p.owner.Modifiers
}

func (p *advantageModifiersPanel) SetAdvantageModifierList(list []*gurps.AdvantageModifier) {
	p.owner.Modifiers = list
}
