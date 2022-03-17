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

package tbl

import (
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

var noteColMap = map[int]int{
	0: gurps.NoteTextColumn,
	1: gurps.NoteReferenceColumn,
}

// NewNoteTableHeaders creates a new set of table column headers for notes.
func NewNoteTableHeaders(forPage bool) []unison.TableColumnHeader {
	return []unison.TableColumnHeader{
		NewHeader(i18n.Text("Note"), "", forPage),
		NewPageRefHeader(forPage),
	}
}

// NewNoteRowData creates a new table data provider function for notes.
func NewNoteRowData(topLevelData []*gurps.Note, forPage bool) func(table *unison.Table) []unison.TableRowData {
	return func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(topLevelData))
		for _, one := range topLevelData {
			rows = append(rows, NewNode(table, nil, noteColMap, one, forPage))
		}
		return rows
	}
}
