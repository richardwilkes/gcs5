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
	"github.com/richardwilkes/gcs/ui/widget/ntable"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type adjustTraitLevelListUndoEdit = *unison.UndoEdit[*adjustTraitLevelList]

type adjustTraitLevelList struct {
	Owner widget.Rebuildable
	List  []*traitLevelAdjuster
}

func (a *adjustTraitLevelList) Apply() {
	for _, one := range a.List {
		one.Apply()
	}
	a.Finish()
}

func (a *adjustTraitLevelList) Finish() {
	entity := a.List[0].Target.OwningEntity()
	if entity != nil {
		entity.Recalculate()
	}
	widget.MarkModified(a.Owner)
}

type traitLevelAdjuster struct {
	Target *gurps.Trait
	Levels fxp.Int
}

func newTraitLevelAdjuster(target *gurps.Trait) *traitLevelAdjuster {
	return &traitLevelAdjuster{
		Target: target,
		Levels: target.Levels,
	}
}

func (a *traitLevelAdjuster) Apply() {
	a.Target.Levels = a.Levels
}

func canAdjustTraitLevel(table *unison.Table[*ntable.Node[*gurps.Trait]], increment bool) bool {
	for _, row := range table.SelectedRows(false) {
		if t := ntable.ExtractFromRowData[*gurps.Trait](row); t != nil && t.IsLeveled() {
			if increment || t.Levels > 0 {
				return true
			}
		}
	}
	return false
}

func adjustTraitLevel(owner widget.Rebuildable, table *unison.Table[*ntable.Node[*gurps.Trait]], increment bool) {
	before := &adjustTraitLevelList{Owner: owner}
	after := &adjustTraitLevelList{Owner: owner}
	for _, row := range table.SelectedRows(false) {
		if t := ntable.ExtractFromRowData[*gurps.Trait](row); t != nil && t.IsLeveled() {
			if increment || t.Levels > 0 {
				before.List = append(before.List, newTraitLevelAdjuster(t))
				original := t.Levels
				levels := original.Trunc()
				if increment {
					levels += fxp.One
				} else if original == levels {
					levels -= fxp.One
				}
				t.Levels = levels.Max(0)
				after.List = append(after.List, newTraitLevelAdjuster(t))
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
			mgr.Add(&unison.UndoEdit[*adjustTraitLevelList]{
				ID:         unison.NextUndoID(),
				EditName:   name,
				UndoFunc:   func(edit adjustTraitLevelListUndoEdit) { edit.BeforeData.Apply() },
				RedoFunc:   func(edit adjustTraitLevelListUndoEdit) { edit.AfterData.Apply() },
				BeforeData: before,
				AfterData:  after,
			})
		}
		before.Finish()
	}
}
