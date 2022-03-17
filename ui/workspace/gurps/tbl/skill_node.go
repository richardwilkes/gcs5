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
	"strings"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

const (
	skillDescriptionColumn = iota
	skillDifficultyColumn
	skillCategoryColumn
	skillReferenceColumn
	skillLevelColumn
	skillRelativeLevelColumn
	skillPointsColumn
)

var (
	_               unison.TableRowData = &SkillNode{}
	_               Matcher             = &SkillNode{}
	skillListColMap                     = map[int]int{
		0: skillDescriptionColumn,
		1: skillDifficultyColumn,
		2: skillCategoryColumn,
		3: skillReferenceColumn,
	}
	skillPageColMap = map[int]int{
		0: skillDescriptionColumn,
		1: skillLevelColumn,
		2: skillRelativeLevelColumn,
		3: skillPointsColumn,
		4: skillReferenceColumn,
	}
)

// SkillNode holds a skill in the skill list.
type SkillNode struct {
	table     *unison.Table
	parent    *SkillNode
	skill     *gurps.Skill
	children  []unison.TableRowData
	cellCache []*CellCache
	forPage   bool
}

func NewSkillTableHeaders(forPage bool) []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	headers = append(headers,
		NewHeader(i18n.Text("Skill / Technique"), "", forPage),
	)
	if forPage {
		headers = append(headers,
			NewHeader(i18n.Text("SL"), i18n.Text("Skill Level"), true),
			NewHeader(i18n.Text("RSL"), i18n.Text("Relative Skill Level"), true),
			NewHeader(i18n.Text("Pts"), i18n.Text("Points"), true),
		)
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

func NewSkillRowData(topLevelData []*gurps.Skill, forPage bool) func(table *unison.Table) []unison.TableRowData {
	return func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(topLevelData))
		for _, one := range topLevelData {
			rows = append(rows, NewSkillNode(table, nil, one, forPage))
		}
		return rows
	}
}

// NewSkillNode creates a new SkillNode.
func NewSkillNode(table *unison.Table, parent *SkillNode, skill *gurps.Skill, forPage bool) *SkillNode {
	n := &SkillNode{
		table:   table,
		parent:  parent,
		skill:   skill,
		forPage: forPage,
	}
	n.cellCache = make([]*CellCache, len(n.colMap()))
	return n
}

func (n *SkillNode) colMap() map[int]int {
	if n.forPage {
		return skillPageColMap
	}
	return skillListColMap
}

// ParentRow returns the parent row, or nil if this is a root node.
func (n *SkillNode) ParentRow() unison.TableRowData {
	return n.parent
}

// CanHaveChildRows returns true if this is a container.
func (n *SkillNode) CanHaveChildRows() bool {
	return n.skill.Container()
}

// ChildRows returns the children of this node.
func (n *SkillNode) ChildRows() []unison.TableRowData {
	if n.skill.Container() && n.children == nil {
		n.children = make([]unison.TableRowData, len(n.skill.Children))
		for i, one := range n.skill.Children {
			n.children[i] = NewSkillNode(n.table, n, one, n.forPage)
		}
	}
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *SkillNode) CellDataForSort(index int) string {
	switch n.colMap()[index] {
	case skillDescriptionColumn:
		text := n.skill.Description()
		secondary := n.skill.SecondaryText()
		if secondary != "" {
			text += "\n" + secondary
		}
		return text
	case skillDifficultyColumn:
		if n.skill.Container() {
			return ""
		}
		return n.skill.Difficulty.Description(n.skill.Entity)
	case skillCategoryColumn:
		return strings.Join(n.skill.Categories, ", ")
	case skillReferenceColumn:
		return n.skill.PageRef
	case skillLevelColumn:
		return n.skill.LevelAsString()
	case skillRelativeLevelColumn:
		return n.skill.AdjustedRelativeLevel().String()
	case skillPointsColumn:
		return n.skill.AdjustedPoints().String()
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *SkillNode) ColumnCell(row, col int, selected bool) unison.Paneler {
	width := n.table.CellWidth(row, col)
	data := n.CellDataForSort(col)
	if n.cellCache[col].Matches(width, data) {
		color := unison.DefaultLabelTheme.OnBackgroundInk
		if selected {
			color = unison.OnSelectionColor
		}
		for _, child := range n.cellCache[col].Panel.Children() {
			child.Self.(*unison.Label).LabelTheme.OnBackgroundInk = color
		}
		return n.cellCache[col].Panel
	}
	p := &unison.Panel{}
	p.Self = p
	layout := &unison.FlexLayout{Columns: 1}
	p.SetLayout(layout)
	var primaryFont, secondaryFont unison.Font
	if n.forPage {
		primaryFont = theme.PageFieldPrimaryFont
		secondaryFont = theme.PageFieldSecondaryFont
	} else {
		primaryFont = unison.FieldFont
		secondaryFont = theme.FieldSecondaryFont
	}
	switch n.colMap()[col] {
	case skillDescriptionColumn:
		CreateAndAddCellLabel(p, width, n.skill.Description(), primaryFont, selected)
		if text := n.skill.SecondaryText(); strings.TrimSpace(text) != "" {
			CreateAndAddCellLabel(p, width, text, secondaryFont, selected)
		}
	case skillReferenceColumn:
		CreateAndAddPageRefCellLabel(p, n.CellDataForSort(col), n.skill.Name, primaryFont, selected)
	case skillLevelColumn, skillRelativeLevelColumn, skillPointsColumn:
		CreateAndAddCellLabel(p, width, n.CellDataForSort(col), primaryFont, selected)
		layout.HAlign = unison.EndAlignment
	default:
		CreateAndAddCellLabel(p, width, n.CellDataForSort(col), primaryFont, selected)
	}
	n.cellCache[col] = &CellCache{
		Width: width,
		Data:  data,
		Panel: p,
	}
	return p
}

// IsOpen returns true if this node should display its children.
func (n *SkillNode) IsOpen() bool {
	return n.skill.Container() && n.skill.Open
}

// SetOpen sets the current open state for this node.
func (n *SkillNode) SetOpen(open bool) {
	if n.skill.Container() && open != n.skill.Open {
		n.skill.Open = open
		n.table.SyncToModel()
	}
}

// Match implements Matcher.
func (n *SkillNode) Match(text string) bool {
	count := len(n.colMap())
	for i := 0; i < count; i++ {
		if strings.Contains(strings.ToLower(n.CellDataForSort(i)), text) {
			return true
		}
	}
	return false
}
