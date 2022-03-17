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

var (
	advantageListColMap = map[int]int{
		0: gurps.AdvantageDescriptionColumn,
		1: gurps.AdvantagePointsColumn,
		2: gurps.AdvantageTypeColumn,
		3: gurps.AdvantageCategoryColumn,
		4: gurps.AdvantageReferenceColumn,
	}
	advantagePageColMap = map[int]int{
		0: gurps.AdvantageDescriptionColumn,
		1: gurps.AdvantagePointsColumn,
		2: gurps.AdvantageReferenceColumn,
	}
)

func NewAdvantageTableHeaders(forPage bool) []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	headers = append(headers,
		NewHeader(i18n.Text("Advantage / Disadvantage"), "", forPage),
		NewHeader(i18n.Text("Pts"), i18n.Text("Points"), forPage),
	)
	if !forPage {
		headers = append(headers,
			NewHeader(i18n.Text("Type"), "", false),
			NewHeader(i18n.Text("Category"), "", false),
		)
	}
	return append(headers, NewPageRefHeader(forPage))
}

func NewAdvantageRowData(topLevelData []*gurps.Advantage, forPage bool) func(table *unison.Table) []unison.TableRowData {
	return func(table *unison.Table) []unison.TableRowData {
		var colMap map[int]int
		if forPage {
			colMap = advantagePageColMap
		} else {
			colMap = advantageListColMap
		}
		rows := make([]unison.TableRowData, 0, len(topLevelData))
		for _, one := range topLevelData {
			rows = append(rows, NewNode(table, nil, colMap, one, forPage))
		}
		return rows
	}
}
