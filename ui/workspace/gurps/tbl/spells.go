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
	spellListColMap = map[int]int{
		0:  gurps.SpellDescriptionColumn,
		1:  gurps.SpellCollegeColumn,
		2:  gurps.SpellResistColumn,
		3:  gurps.SpellClassColumn,
		4:  gurps.SpellCastCostColumn,
		5:  gurps.SpellMaintainCostColumn,
		6:  gurps.SpellCastTimeColumn,
		7:  gurps.SpellDurationColumn,
		8:  gurps.SpellDifficultyColumn,
		9:  gurps.SpellCategoryColumn,
		10: gurps.SpellReferenceColumn,
	}
	spellPageColMap = map[int]int{
		0: gurps.SpellDescriptionForPageColumn,
		1: gurps.SpellCollegeColumn,
		2: gurps.SpellDifficultyColumn,
		3: gurps.SpellLevelColumn,
		4: gurps.SpellRelativeLevelColumn,
		5: gurps.SpellPointsColumn,
		6: gurps.SpellReferenceColumn,
	}
)

// NewSpellTableHeaders creates a new set of table column headers for spells.
func NewSpellTableHeaders(forPage bool) []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	headers = append(headers,
		NewHeader(i18n.Text("Spell"), "", forPage),
		NewHeader(i18n.Text("College"), "", forPage),
	)
	if !forPage {
		headers = append(headers,
			NewHeader(i18n.Text("Resist"), i18n.Text("Resistance"), forPage),
			NewHeader(i18n.Text("Class"), "", forPage),
			NewHeader(i18n.Text("Cost"), i18n.Text("The mana cost to cast the spell"), forPage),
			NewHeader(i18n.Text("Maintain"), i18n.Text("The mana cost to maintain the spell"), forPage),
			NewHeader(i18n.Text("Time"), i18n.Text("The time required to cast the spell"), forPage),
			NewHeader(i18n.Text("Duration"), "", forPage),
		)
	}
	headers = append(headers,
		NewHeader(i18n.Text("Diff"), i18n.Text("Difficulty"), forPage),
	)
	if forPage {
		headers = append(headers,
			NewHeader(i18n.Text("SL"), i18n.Text("Skill Level"), true),
			NewHeader(i18n.Text("RSL"), i18n.Text("Relative Skill Level"), true),
			NewHeader(i18n.Text("Pts"), i18n.Text("Points"), true),
		)
	} else {
		headers = append(headers,
			NewHeader(i18n.Text("Category"), "", false),
		)
	}
	return append(headers, NewPageRefHeader(forPage))
}

// NewSpellRowData creates a new table data provider function for spells.
func NewSpellRowData(topLevelRowsProvider func() []*gurps.Spell, forPage bool) func(table *unison.Table) []unison.TableRowData {
	return func(table *unison.Table) []unison.TableRowData {
		var colMap map[int]int
		if forPage {
			colMap = spellPageColMap
		} else {
			colMap = spellListColMap
		}
		data := topLevelRowsProvider()
		rows := make([]unison.TableRowData, 0, len(data))
		for _, one := range data {
			rows = append(rows, NewNode(table, nil, colMap, one, forPage))
		}
		return rows
	}
}
