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

// NewNotesPageList creates the notes page list.
func NewNotesPageList(entity *gurps.Entity) *PageList {
	return NewPageList(entity, []unison.TableColumnHeader{
		tbl.NewHeader(i18n.Text("Note"), "", true),
		tbl.NewPageRefHeader(true),
	}, func(table *unison.Table) []unison.TableRowData {
		//rows := make([]unison.TableRowData, 0, len(entity.Spells))
		//for _, one := range entity.Spells {
		//	rows = append(rows, NewAdvantagePageNode(table, nil, one))
		//}
		//return rows
		return nil
	})
}
