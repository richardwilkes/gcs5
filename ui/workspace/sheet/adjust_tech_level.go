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

type adjustTechLevelListUndoEdit = *unison.UndoEdit[*adjustTechLevelList]

type adjustTechLevelList struct {
	Owner widget.Rebuildable
	List  []*techLevelAdjuster
}

func (a *adjustTechLevelList) Apply() {
	for _, one := range a.List {
		one.Apply()
	}
	entity := a.List[0].Target.OwningEntity()
	if entity != nil {
		entity.Recalculate()
	}
	widget.MarkModified(a.Owner)
}

type techLevelAdjuster struct {
	Target    gurps.TechLevelProvider
	TechLevel string
}

func newTechLevelAdjuster(target gurps.TechLevelProvider) *techLevelAdjuster {
	return &techLevelAdjuster{
		Target:    target,
		TechLevel: target.TL(),
	}
}

func (a *techLevelAdjuster) Apply() {
	a.Target.SetTL(a.TechLevel)
}

func canAdjustTechLevel(table *unison.Table, amount fxp.Int) bool {
	for _, row := range table.SelectedRows(false) {
		if provider := tbl.ExtractFromRowData[gurps.TechLevelProvider](row); provider != nil {
			if _, changed := gurps.AdjustTechLevel(provider.TL(), amount); changed {
				return true
			}
		}
	}
	return false
}

func adjustTechLevel(owner widget.Rebuildable, table *unison.Table, amount fxp.Int) {
	before := &adjustTechLevelList{Owner: owner}
	after := &adjustTechLevelList{Owner: owner}
	for _, row := range table.SelectedRows(false) {
		if provider := tbl.ExtractFromRowData[gurps.TechLevelProvider](row); provider != nil {
			if tl, changed := gurps.AdjustTechLevel(provider.TL(), amount); changed {
				before.List = append(before.List, newTechLevelAdjuster(provider))
				provider.SetTL(tl)
				after.List = append(after.List, newTechLevelAdjuster(provider))
			}
		}
	}
	if len(before.List) > 0 {
		if mgr := unison.UndoManagerFor(table); mgr != nil {
			var name string
			if amount < 0 {
				name = i18n.Text("Decrease Tech Level")
			} else {
				name = i18n.Text("Increase Tech Level")
			}
			mgr.Add(&unison.UndoEdit[*adjustTechLevelList]{
				ID:         unison.NextUndoID(),
				EditName:   name,
				UndoFunc:   func(edit adjustTechLevelListUndoEdit) { edit.BeforeData.Apply() },
				RedoFunc:   func(edit adjustTechLevelListUndoEdit) { edit.AfterData.Apply() },
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
