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
	"github.com/richardwilkes/gcs/ui/workspace/gurps/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// NewSpellsPageList creates the spells page list.
func NewSpellsPageList(entity *gurps.Entity) *PageList {
	return NewPageList(entity, []unison.TableColumnHeader{
		tbl.NewHeader(i18n.Text("Spell"), "", true),
		tbl.NewHeader(i18n.Text("Resist"), "", true),
		tbl.NewHeader(i18n.Text("Class"), "", true),
		tbl.NewHeader(i18n.Text("Cost"), "", true),
		tbl.NewHeader(i18n.Text("Maintain"), "", true),
		tbl.NewHeader(i18n.Text("Time"), "", true),
		tbl.NewHeader(i18n.Text("Duration"), "", true),
		tbl.NewHeader(i18n.Text("SL"), i18n.Text("Skill Level"), true),
		tbl.NewHeader(i18n.Text("RSL"), i18n.Text("Relative Skill Level"), true),
		tbl.NewHeader(i18n.Text("Pts"), i18n.Text("Points"), true),
		tbl.NewPageRefHeader(true),
	}, func(table *unison.Table) []unison.TableRowData {
		//rows := make([]unison.TableRowData, 0, len(entity.Spells))
		//for _, one := range entity.Spells {
		//	rows = append(rows, NewSpellPageNode(table, nil, one))
		//}
		//return rows
		return nil
	})
}
