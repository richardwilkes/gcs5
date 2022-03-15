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

// NewRangedWeaponsPageList creates the ranged weapons page list.
func NewRangedWeaponsPageList(entity *gurps.Entity) *PageList {
	return NewPageList(entity, []unison.TableColumnHeader{
		tbl.NewHeader(i18n.Text("Ranged Weapon"), "", true),
		tbl.NewHeader(i18n.Text("Usage"), "", true),
		tbl.NewHeader(i18n.Text("Lvl"), "", true),
		tbl.NewHeader(i18n.Text("Acc"), "", true),
		tbl.NewHeader(i18n.Text("Damage"), "", true),
		tbl.NewHeader(i18n.Text("Range"), "", true),
		tbl.NewHeader(i18n.Text("RoF"), "", true),
		tbl.NewHeader(i18n.Text("Shots"), "", true),
		tbl.NewHeader(i18n.Text("Bulk"), "", true),
		tbl.NewHeader(i18n.Text("Rcl"), "", true),
		tbl.NewHeader(i18n.Text("ST"), "", true),
	}, func(table *unison.Table) []unison.TableRowData {
		//rows := make([]unison.TableRowData, 0, len(entity.Spells))
		//for _, one := range entity.Spells {
		//	rows = append(rows, NewAdvantagePageNode(table, nil, one))
		//}
		//return rows
		return nil
	})
}
