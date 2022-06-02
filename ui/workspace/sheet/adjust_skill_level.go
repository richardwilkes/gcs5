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

type adjustSkillLevelListUndoEdit = *unison.UndoEdit[*adjustRawPointsList]

type adjustRawPointsList struct {
	Owner widget.Rebuildable
	List  []*rawPointsAdjuster
}

func (a *adjustRawPointsList) Apply() {
	for _, one := range a.List {
		one.Apply()
	}
	entity := a.List[0].Target.OwningEntity()
	if entity != nil {
		entity.Recalculate()
	}
	widget.MarkModified(a.Owner)
}

type rawPointsAdjuster struct {
	Target gurps.RawPointsAdjuster
	Points fxp.Int
}

func newSkillLevelAdjuster(target gurps.RawPointsAdjuster) *rawPointsAdjuster {
	return &rawPointsAdjuster{
		Target: target,
		Points: target.RawPoints(),
	}
}

func (a *rawPointsAdjuster) Apply() {
	a.Target.SetRawPoints(a.Points)
}

func canAdjustSkillLevel(table *unison.Table, increment bool) bool {
	for _, row := range table.SelectedRows(false) {
		if provider := tbl.ExtractFromRowData[gurps.SkillAdjustmentProvider](row); provider != nil && !provider.Container() {
			if increment || provider.RawPoints() > 0 {
				return true
			}
		}
	}
	return false
}

func adjustSkillLevel(owner widget.Rebuildable, table *unison.Table, increment bool) {
	before := &adjustRawPointsList{Owner: owner}
	after := &adjustRawPointsList{Owner: owner}
	for _, row := range table.SelectedRows(false) {
		if provider := tbl.ExtractFromRowData[gurps.SkillAdjustmentProvider](row); provider != nil {
			if increment || provider.RawPoints() > 0 {
				before.List = append(before.List, newSkillLevelAdjuster(provider))
				if increment {
					provider.IncrementSkillLevel()
				} else {
					provider.DecrementSkillLevel()
				}
				after.List = append(after.List, newSkillLevelAdjuster(provider))
			}
		}
	}
	if len(before.List) > 0 {
		if mgr := unison.UndoManagerFor(table); mgr != nil {
			var name string
			if increment {
				name = i18n.Text("Increase Skill Level")
			} else {
				name = i18n.Text("Decrease Skill Level")
			}
			mgr.Add(&unison.UndoEdit[*adjustRawPointsList]{
				ID:         unison.NextUndoID(),
				EditName:   name,
				UndoFunc:   func(edit adjustSkillLevelListUndoEdit) { edit.BeforeData.Apply() },
				RedoFunc:   func(edit adjustSkillLevelListUndoEdit) { edit.AfterData.Apply() },
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
