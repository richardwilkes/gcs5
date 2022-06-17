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
	"github.com/richardwilkes/gcs/ui/widget/ntable"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type toggleEquippedUndoEdit = *unison.UndoEdit[*toggleEquippedList]

type toggleEquippedList struct {
	Owner widget.Rebuildable
	List  []*equippedAdjuster
}

func (a *toggleEquippedList) Apply() {
	for _, one := range a.List {
		one.Apply()
	}
	a.Finish()
}

func (a *toggleEquippedList) Finish() {
	entity := a.List[0].Target.OwningEntity()
	if entity != nil {
		entity.Recalculate()
	}
	widget.MarkModified(a.Owner)
}

type equippedAdjuster struct {
	Target   *gurps.Equipment
	Equipped bool
}

func newEquippedAdjuster(target *gurps.Equipment) *equippedAdjuster {
	return &equippedAdjuster{
		Target:   target,
		Equipped: target.Equipped,
	}
}

func (a *equippedAdjuster) Apply() {
	a.Target.Equipped = a.Equipped
}

func canToggleEquipped(table *unison.Table[*ntable.Node[*gurps.Equipment]]) bool {
	for _, row := range table.SelectedRows(false) {
		if eqp := ntable.ExtractFromRowData[*gurps.Equipment](row); eqp != nil {
			return true
		}
	}
	return false
}

func toggleEquipped(owner widget.Rebuildable, table *unison.Table[*ntable.Node[*gurps.Equipment]]) {
	before := &toggleEquippedList{Owner: owner}
	after := &toggleEquippedList{Owner: owner}
	for _, row := range table.SelectedRows(false) {
		if eqp := ntable.ExtractFromRowData[*gurps.Equipment](row); eqp != nil {
			before.List = append(before.List, newEquippedAdjuster(eqp))
			eqp.Equipped = !eqp.Equipped
			after.List = append(after.List, newEquippedAdjuster(eqp))
		}
	}
	if len(before.List) > 0 {
		if mgr := unison.UndoManagerFor(table); mgr != nil {
			mgr.Add(&unison.UndoEdit[*toggleEquippedList]{
				ID:         unison.NextUndoID(),
				EditName:   i18n.Text("Toggle Equipped"),
				UndoFunc:   func(edit toggleEquippedUndoEdit) { edit.BeforeData.Apply() },
				RedoFunc:   func(edit toggleEquippedUndoEdit) { edit.AfterData.Apply() },
				BeforeData: before,
				AfterData:  after,
			})
		}
		before.Finish()
	}
}
