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

package editors

import (
	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
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
	table    *unison.Table[*Node[*gurps.Skill]]
	colMap   map[int]int
	provider gurps.SkillListProvider
	forPage  bool
}

// NewSkillsProvider creates a new table provider for skills.
func NewSkillsProvider(provider gurps.SkillListProvider, forPage bool) widget.TableProvider[*Node[*gurps.Skill]] {
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

func (p *skillsProvider) SetTable(table *unison.Table[*Node[*gurps.Skill]]) {
	p.table = table
}

func (p *skillsProvider) RootRowCount() int {
	return len(p.provider.SkillList())
}

func (p *skillsProvider) RootRows() []*Node[*gurps.Skill] {
	data := p.provider.SkillList()
	rows := make([]*Node[*gurps.Skill], 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode[*gurps.Skill](p.table, nil, p.colMap, one, p.forPage))
	}
	return rows
}

func (p *skillsProvider) SetRootRows(rows []*Node[*gurps.Skill]) {
	p.provider.SetSkillList(ExtractNodeDataFromList(rows))
}

func (p *skillsProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *skillsProvider) DragKey() string {
	return gid.Skill
}

func (p *skillsProvider) DragSVG() *unison.SVG {
	return res.GCSSkillsSVG
}

func (p *skillsProvider) DropShouldMoveData(from, to *unison.Table[*Node[*gurps.Skill]]) bool {
	return from == to
}

func (p *skillsProvider) ItemNames() (singular, plural string) {
	return i18n.Text("Skill"), i18n.Text("Skills")
}

func (p *skillsProvider) Headers() []unison.TableColumnHeader[*Node[*gurps.Skill]] {
	var headers []unison.TableColumnHeader[*Node[*gurps.Skill]]
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.SkillDescriptionColumn:
			headers = append(headers, NewHeader[*gurps.Skill](i18n.Text("Skill / Technique"), "", p.forPage))
		case gurps.SkillDifficultyColumn:
			headers = append(headers, NewHeader[*gurps.Skill](i18n.Text("Diff"), i18n.Text("Difficulty"), p.forPage))
		case gurps.SkillTagsColumn:
			headers = append(headers, NewHeader[*gurps.Skill](i18n.Text("Tags"), "", p.forPage))
		case gurps.SkillReferenceColumn:
			headers = append(headers, NewPageRefHeader[*gurps.Skill](p.forPage))
		case gurps.SkillLevelColumn:
			headers = append(headers, NewHeader[*gurps.Skill](i18n.Text("SL"), i18n.Text("Skill Level"), p.forPage))
		case gurps.SkillRelativeLevelColumn:
			headers = append(headers, NewHeader[*gurps.Skill](i18n.Text("RSL"), i18n.Text("Relative Skill Level"), p.forPage))
		case gurps.SkillPointsColumn:
			headers = append(headers, NewHeader[*gurps.Skill](i18n.Text("Pts"), i18n.Text("Points"), p.forPage))
		default:
			jot.Fatalf(1, "invalid skill column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *skillsProvider) SyncHeader(_ []unison.TableColumnHeader[*Node[*gurps.Skill]]) {
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

func (p *skillsProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table[*Node[*gurps.Skill]]) {
	OpenEditor[*gurps.Skill](table, func(item *gurps.Skill) { EditSkill(owner, item) })
}

func (p *skillsProvider) CreateItem(owner widget.Rebuildable, table *unison.Table[*Node[*gurps.Skill]], variant widget.ItemVariant) {
	var item *gurps.Skill
	switch variant {
	case widget.NoItemVariant:
		item = gurps.NewSkill(p.Entity(), nil, false)
	case widget.ContainerItemVariant:
		item = gurps.NewSkill(p.Entity(), nil, true)
	case widget.AlternateItemVariant:
		item = gurps.NewTechnique(p.Entity(), nil, "")
	default:
		jot.Fatal(1, "unhandled variant")
	}
	InsertItem[*gurps.Skill](owner, table, item,
		func(target *gurps.Skill) []*gurps.Skill { return target.Children },
		func(target *gurps.Skill, children []*gurps.Skill) { target.Children = children },
		p.provider.SkillList, p.provider.SetSkillList,
		func(_ *unison.Table[*Node[*gurps.Skill]]) []*Node[*gurps.Skill] { return p.RootRows() },
		func(target *gurps.Skill) uuid.UUID { return target.ID })
	EditSkill(owner, item)
}

func (p *skillsProvider) DeleteSelection(table *unison.Table[*Node[*gurps.Skill]]) {
	deleteTableSelection(table, p.provider.SkillList(),
		func(nodes []*gurps.Skill) { p.provider.SetSkillList(nodes) },
		func(node *gurps.Skill) *[]*gurps.Skill { return &node.Children })
}
