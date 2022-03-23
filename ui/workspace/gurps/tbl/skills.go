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
	skillListColMap = map[int]int{
		0: gurps.SkillDescriptionColumn,
		1: gurps.SkillDifficultyColumn,
		2: gurps.SkillCategoryColumn,
		3: gurps.SkillReferenceColumn,
	}
	entitySkillPageColMap = map[int]int{
		0: gurps.SkillDescriptionColumn,
		1: gurps.SkillLevelColumn,
		2: gurps.SkillRelativeLevelColumn,
		3: gurps.SkillPointsColumn,
		4: gurps.SkillReferenceColumn,
	}
	skillPageColMap = map[int]int{
		0: gurps.SkillDescriptionColumn,
		1: gurps.SkillPointsColumn,
		2: gurps.SkillReferenceColumn,
	}
)

// NewSkillTableHeaders creates a new set of table column headers for skills.
func NewSkillTableHeaders(provider gurps.SkillListProvider, forPage bool) []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	headers = append(headers,
		NewHeader(i18n.Text("Skill / Technique"), "", forPage),
	)
	if forPage {
		if _, ok := provider.(*gurps.Entity); ok {
			headers = append(headers,
				NewHeader(i18n.Text("SL"), i18n.Text("Skill Level"), true),
				NewHeader(i18n.Text("RSL"), i18n.Text("Relative Skill Level"), true),
			)
		}
		headers = append(headers, NewHeader(i18n.Text("Pts"), i18n.Text("Points"), true))
	} else {
		headers = append(headers,
			NewHeader(i18n.Text("Diff"), i18n.Text("Difficulty"), false),
			NewHeader(i18n.Text("Category"), "", false),
		)
	}
	return append(headers,
		NewPageRefHeader(forPage),
	)
}

// NewSkillRowData creates a new table data provider function for skills.
func NewSkillRowData(provider gurps.SkillListProvider, forPage bool) func(table *unison.Table) []unison.TableRowData {
	return func(table *unison.Table) []unison.TableRowData {
		var colMap map[int]int
		if forPage {
			if _, ok := provider.(*gurps.Entity); ok {
				colMap = entitySkillPageColMap
			} else {
				colMap = skillPageColMap
			}
		} else {
			colMap = skillListColMap
		}
		data := provider.SkillList()
		rows := make([]unison.TableRowData, 0, len(data))
		for _, one := range data {
			rows = append(rows, NewNode(table, nil, colMap, one, forPage))
		}
		return rows
	}
}
