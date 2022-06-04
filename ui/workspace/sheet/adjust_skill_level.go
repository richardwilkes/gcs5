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
	"github.com/richardwilkes/gcs/ui/workspace/editors"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type adjustSkillLevelListUndoEdit = *unison.UndoEdit[*adjustRawPointsList]

func canAdjustSkillLevel(table *unison.Table, increment bool) bool {
	for _, row := range table.SelectedRows(false) {
		if provider := editors.ExtractFromRowData[gurps.SkillAdjustmentProvider](row); provider != nil && !provider.Container() {
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
		if provider := editors.ExtractFromRowData[gurps.SkillAdjustmentProvider](row); provider != nil {
			if increment || provider.RawPoints() > 0 {
				before.List = append(before.List, newRawPointsAdjuster(provider))
				if increment {
					provider.IncrementSkillLevel()
				} else {
					provider.DecrementSkillLevel()
				}
				after.List = append(after.List, newRawPointsAdjuster(provider))
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
		before.Finish()
	}
}
