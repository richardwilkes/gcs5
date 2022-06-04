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
	"github.com/richardwilkes/gcs/ui/workspace/editors"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type adjustAdvantageLevelListUndoEdit = *unison.UndoEdit[*adjustAdvantageLevelList]

type adjustAdvantageLevelList struct {
	Owner widget.Rebuildable
	List  []*advantageLevelAdjuster
}

func (a *adjustAdvantageLevelList) Apply() {
	for _, one := range a.List {
		one.Apply()
	}
	a.Finish()
}

func (a *adjustAdvantageLevelList) Finish() {
	entity := a.List[0].Target.OwningEntity()
	if entity != nil {
		entity.Recalculate()
	}
	widget.MarkModified(a.Owner)
}

type advantageLevelAdjuster struct {
	Target *gurps.Advantage
	Levels fxp.Int
}

func newAdvantageLevelAdjuster(target *gurps.Advantage) *advantageLevelAdjuster {
	return &advantageLevelAdjuster{
		Target: target,
		Levels: target.Levels,
	}
}

func (a *advantageLevelAdjuster) Apply() {
	a.Target.Levels = a.Levels
}

func canAdjustAdvantageLevel(table *unison.Table, increment bool) bool {
	for _, row := range table.SelectedRows(false) {
		if adv := editors.ExtractFromRowData[*gurps.Advantage](row); adv != nil && adv.IsLeveled() {
			if increment || adv.Levels > 0 {
				return true
			}
		}
	}
	return false
}

func adjustAdvantageLevel(owner widget.Rebuildable, table *unison.Table, increment bool) {
	before := &adjustAdvantageLevelList{Owner: owner}
	after := &adjustAdvantageLevelList{Owner: owner}
	for _, row := range table.SelectedRows(false) {
		if adv := editors.ExtractFromRowData[*gurps.Advantage](row); adv != nil && adv.IsLeveled() {
			if increment || adv.Levels > 0 {
				before.List = append(before.List, newAdvantageLevelAdjuster(adv))
				original := adv.Levels
				levels := original.Trunc()
				if increment {
					levels += fxp.One
				} else if original == levels {
					levels -= fxp.One
				}
				adv.Levels = levels.Max(0)
				after.List = append(after.List, newAdvantageLevelAdjuster(adv))
			}
		}
	}
	if len(before.List) > 0 {
		if mgr := unison.UndoManagerFor(table); mgr != nil {
			var name string
			if increment {
				name = i18n.Text("Increment Level")
			} else {
				name = i18n.Text("Decrement Level")
			}
			mgr.Add(&unison.UndoEdit[*adjustAdvantageLevelList]{
				ID:         unison.NextUndoID(),
				EditName:   name,
				UndoFunc:   func(edit adjustAdvantageLevelListUndoEdit) { edit.BeforeData.Apply() },
				RedoFunc:   func(edit adjustAdvantageLevelListUndoEdit) { edit.AfterData.Apply() },
				BeforeData: before,
				AfterData:  after,
			})
		}
		before.Finish()
	}
}
