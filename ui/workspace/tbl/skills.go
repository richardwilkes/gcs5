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
	skillListColMap = map[int]int{
		0: gurps.SkillDescriptionColumn,
		1: gurps.SkillDifficultyColumn,
		2: gurps.SkillTagsColumn,
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

type skillsProvider struct {
	colMap   map[int]int
	provider gurps.SkillListProvider
	forPage  bool
}

// NewSkillsProvider creates a new table provider for skills.
func NewSkillsProvider(provider gurps.SkillListProvider, forPage bool) TableProvider {
	p := &skillsProvider{
		provider: provider,
		forPage:  forPage,
	}
	if forPage {
		if _, ok := provider.(*gurps.Entity); ok {
			p.colMap = entitySkillPageColMap
		} else {
			p.colMap = skillPageColMap
		}
	} else {
		p.colMap = skillListColMap
	}
	return p
}

func (p *skillsProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *skillsProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.SkillDescriptionColumn:
			headers = append(headers, NewHeader(i18n.Text("Skill / Technique"), "", p.forPage))
		case gurps.SkillDifficultyColumn:
			headers = append(headers, NewHeader(i18n.Text("Diff"), i18n.Text("Difficulty"), p.forPage))
		case gurps.SkillTagsColumn:
			headers = append(headers, NewHeader(i18n.Text("Tags"), "", p.forPage))
		case gurps.SkillReferenceColumn:
			headers = append(headers, NewPageRefHeader(p.forPage))
		case gurps.SkillLevelColumn:
			headers = append(headers, NewHeader(i18n.Text("SL"), i18n.Text("Skill Level"), p.forPage))
		case gurps.SkillRelativeLevelColumn:
			headers = append(headers, NewHeader(i18n.Text("RSL"), i18n.Text("Relative Skill Level"), p.forPage))
		case gurps.SkillPointsColumn:
			headers = append(headers, NewHeader(i18n.Text("Pts"), i18n.Text("Points"), p.forPage))
		default:
			jot.Fatalf(1, "invalid skill column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *skillsProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.SkillList()
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, p.colMap, one, p.forPage))
	}
	return rows
}

func (p *skillsProvider) SyncHeader(_ []unison.TableColumnHeader) {
}

func (p *skillsProvider) HierarchyColumnIndex() int {
	for k, v := range p.colMap {
		if v == gurps.SkillDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *skillsProvider) ExcessWidthColumnIndex() int {
	return p.HierarchyColumnIndex()
}

func (p *skillsProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table) {
	for _, row := range table.SelectedRows(false) {
		if node, ok := row.(*Node); ok {
			var s *gurps.Skill
			if s, ok = node.Data().(*gurps.Skill); ok {
				editors.EditSkill(owner, s)
			}
		}
	}
}

func (p *skillsProvider) CreateItem(owner widget.Rebuildable, table *unison.Table, container bool) {
	// TODO: Implement
}
