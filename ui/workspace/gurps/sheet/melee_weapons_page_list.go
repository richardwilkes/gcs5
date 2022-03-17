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

// NewMeleeWeaponsPageList creates the melee weapons page list.
func NewMeleeWeaponsPageList(entity *gurps.Entity) *PageList {
	return NewPageList([]unison.TableColumnHeader{
		tbl.NewHeader(i18n.Text("Melee Weapon"), "", true),
		tbl.NewHeader(i18n.Text("Usage"), "", true),
		tbl.NewHeader(i18n.Text("Lvl"), "", true),
		tbl.NewHeader(i18n.Text("Parry"), "", true),
		tbl.NewHeader(i18n.Text("Block"), "", true),
		tbl.NewHeader(i18n.Text("Damage"), "", true),
		tbl.NewHeader(i18n.Text("Reach"), "", true),
		tbl.NewHeader(i18n.Text("ST"), "", true),
	}, 0, func(table *unison.Table) []unison.TableRowData {
		//rows := make([]unison.TableRowData, 0, len(entity.Spells))
		//for _, one := range entity.Spells {
		//	rows = append(rows, NewAdvantagePageNode(table, nil, one))
		//}
		//return rows
		return nil
	})
}
