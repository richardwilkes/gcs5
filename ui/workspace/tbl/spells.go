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
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/editors"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
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
		9:  gurps.SpellTagsColumn,
		10: gurps.SpellReferenceColumn,
	}
	entitySpellPageColMap = map[int]int{
		0: gurps.SpellDescriptionForPageColumn,
		1: gurps.SpellCollegeColumn,
		2: gurps.SpellLevelColumn,
		3: gurps.SpellRelativeLevelColumn,
		4: gurps.SpellPointsColumn,
		5: gurps.SpellReferenceColumn,
	}
	spellPageColMap = map[int]int{
		0: gurps.SpellDescriptionForPageColumn,
		1: gurps.SpellCollegeColumn,
		2: gurps.SpellDifficultyColumn,
		3: gurps.SpellPointsColumn,
		4: gurps.SpellReferenceColumn,
	}
)

type spellsProvider struct {
	colMap   map[int]int
	provider gurps.SpellListProvider
	forPage  bool
}

// NewSpellsProvider creates a new table provider for spells.
func NewSpellsProvider(provider gurps.SpellListProvider, forPage bool) TableProvider {
	p := &spellsProvider{
		provider: provider,
		forPage:  forPage,
	}
	if forPage {
		if _, ok := provider.(*gurps.Entity); ok {
			p.colMap = entitySpellPageColMap
		} else {
			p.colMap = spellPageColMap
		}
	} else {
		p.colMap = spellListColMap
	}
	return p
}

func (p *spellsProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *spellsProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.SpellDescriptionColumn, gurps.SpellDescriptionForPageColumn:
			headers = append(headers, NewHeader(i18n.Text("Spell"), "", p.forPage))
		case gurps.SpellResistColumn:
			headers = append(headers, NewHeader(i18n.Text("Resist"), i18n.Text("Resistance"), p.forPage))
		case gurps.SpellClassColumn:
			headers = append(headers, NewHeader(i18n.Text("Class"), "", p.forPage))
		case gurps.SpellCollegeColumn:
			headers = append(headers, NewHeader(i18n.Text("College"), "", p.forPage))
		case gurps.SpellCastCostColumn:
			headers = append(headers, NewHeader(i18n.Text("Cost"), i18n.Text("The mana cost to cast the spell"),
				p.forPage))
		case gurps.SpellMaintainCostColumn:
			headers = append(headers, NewHeader(i18n.Text("Maintain"), i18n.Text("The mana cost to maintain the spell"),
				p.forPage))
		case gurps.SpellCastTimeColumn:
			headers = append(headers, NewHeader(i18n.Text("Time"), i18n.Text("The time required to cast the spell"),
				p.forPage))
		case gurps.SpellDurationColumn:
			headers = append(headers, NewHeader(i18n.Text("Duration"), "", p.forPage))
		case gurps.SpellDifficultyColumn:
			headers = append(headers, NewHeader(i18n.Text("Diff"), i18n.Text("Difficulty"), p.forPage))
		case gurps.SpellTagsColumn:
			headers = append(headers, NewHeader(i18n.Text("Tags"), "", p.forPage))
		case gurps.SpellReferenceColumn:
			headers = append(headers, NewPageRefHeader(p.forPage))
		case gurps.SpellLevelColumn:
			headers = append(headers, NewHeader(i18n.Text("SL"), i18n.Text("Skill Level"), p.forPage))
		case gurps.SpellRelativeLevelColumn:
			headers = append(headers, NewHeader(i18n.Text("RSL"), i18n.Text("Relative Skill Level"), p.forPage))
		case gurps.SpellPointsColumn:
			headers = append(headers, NewHeader(i18n.Text("Pts"), i18n.Text("Points"), p.forPage))
		default:
			jot.Fatalf(1, "invalid spell column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *spellsProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.SpellList()
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, p.colMap, one, p.forPage))
	}
	return rows
}

func (p *spellsProvider) SyncHeader(_ []unison.TableColumnHeader) {
}

func (p *spellsProvider) HierarchyColumnIndex() int {
	for k, v := range p.colMap {
		if v == gurps.SpellDescriptionColumn || v == gurps.SpellDescriptionForPageColumn {
			return k
		}
	}
	return 0
}

func (p *spellsProvider) ExcessWidthColumnIndex() int {
	return p.HierarchyColumnIndex()
}

func (p *spellsProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table) {
	for _, row := range table.SelectedRows(false) {
		if node, ok := row.(*Node); ok {
			var s *gurps.Spell
			if s, ok = node.Data().(*gurps.Spell); ok {
				editors.EditSpell(owner, s)
			}
		}
	}
}

func (p *spellsProvider) CreateItem(owner widget.Rebuildable, table *unison.Table, container bool) {
	// TODO: Implement
}
