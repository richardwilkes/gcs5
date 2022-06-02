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
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type adjustQuantityListUndoEdit = *unison.UndoEdit[*adjustQuantityList]

type adjustQuantityList struct {
	Owner widget.Rebuildable
	List  []*quantityAdjuster
}

func (a *adjustQuantityList) Apply() {
	for _, one := range a.List {
		one.Apply()
	}
	entity := a.List[0].Target.OwningEntity()
	if entity != nil {
		entity.Recalculate()
	}
	widget.MarkModified(a.Owner)
}

type quantityAdjuster struct {
	Target   *gurps.Equipment
	Quantity fxp.Int
}

func newQuantityAdjuster(target *gurps.Equipment) *quantityAdjuster {
	return &quantityAdjuster{
		Target:   target,
		Quantity: target.Quantity,
	}
}

func (a *quantityAdjuster) Apply() {
	a.Target.Quantity = a.Quantity
}

func canAdjustQuantity(table *unison.Table, increment bool) bool {
	for _, row := range table.SelectedRows(false) {
		if eqp := tbl.ExtractFromRowData[*gurps.Equipment](row); eqp != nil {
			if increment || eqp.Quantity > 0 {
				return true
			}
		}
	}
	return false
}

func adjustQuantity(owner widget.Rebuildable, table *unison.Table, increment bool) {
	before := &adjustQuantityList{Owner: owner}
	after := &adjustQuantityList{Owner: owner}
	for _, row := range table.SelectedRows(false) {
		if eqp := tbl.ExtractFromRowData[*gurps.Equipment](row); eqp != nil {
			if increment || eqp.Quantity > 0 {
				before.List = append(before.List, newQuantityAdjuster(eqp))
				original := eqp.Quantity
				qty := original.Trunc()
				if increment {
					qty += fxp.One
				} else if original == qty {
					qty -= fxp.One
				}
				eqp.Quantity = qty.Max(0)
				after.List = append(after.List, newQuantityAdjuster(eqp))
			}
		}
	}
	if len(before.List) > 0 {
		if mgr := unison.UndoManagerFor(table); mgr != nil {
			var name string
			if increment {
				name = i18n.Text("Increment Quantity")
			} else {
				name = i18n.Text("Decrement Quantity")
			}
			mgr.Add(&unison.UndoEdit[*adjustQuantityList]{
				ID:         unison.NextUndoID(),
				EditName:   name,
				UndoFunc:   func(edit adjustQuantityListUndoEdit) { edit.BeforeData.Apply() },
				RedoFunc:   func(edit adjustQuantityListUndoEdit) { edit.AfterData.Apply() },
				BeforeData: before,
				AfterData:  after,
			})
		}
		entity := before.List[0].Target.OwningEntity()
		if entity != nil {
			entity.Recalculate()
		}
		widget.MarkModified(before.Owner)
	}
}
