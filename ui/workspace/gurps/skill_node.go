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

package gurps

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

const (
	skillDescriptionColumn = iota
	skillDifficultyColumn
	skillCategoryColumn
	skillReferenceColumn
	skillColumnCount
)

var (
	_ unison.TableRowData = &SkillNode{}
	_ Matcher             = &SkillNode{}
)

// SkillNode holds a skill in the skill list.
type SkillNode struct {
	table     *unison.Table
	parent    *SkillNode
	skill     *gurps.Skill
	children  []unison.TableRowData
	cellCache []*cellCache
}

// NewSkillListDockable creates a new unison.Dockable for skill list files.
func NewSkillListDockable(filePath string) (unison.Dockable, error) {
	skills, err := gurps.NewSkillsFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	diffHdr := unison.NewTableColumnHeader(i18n.Text("Diff"))
	diffHdr.Tooltip = unison.NewTooltipWithText(i18n.Text("Difficulty"))
	return NewListFileDockable(filePath, []unison.TableColumnHeader{
		unison.NewTableColumnHeader(i18n.Text("Skill / Technique")),
		diffHdr,
		unison.NewTableColumnHeader(i18n.Text("Category")),
		newPageReferenceHeader(),
	}, func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(skills))
		for _, one := range skills {
			rows = append(rows, NewSkillNode(table, nil, one))
		}
		return rows
	}), nil
}

// NewSkillNode creates a new SkillNode.
func NewSkillNode(table *unison.Table, parent *SkillNode, skill *gurps.Skill) *SkillNode {
	n := &SkillNode{
		table:     table,
		parent:    parent,
		skill:     skill,
		cellCache: make([]*cellCache, skillColumnCount),
	}
	return n
}

// ParentRow returns the parent row, or nil if this is a root node.
func (n *SkillNode) ParentRow() unison.TableRowData {
	return n.parent
}

// CanHaveChildRows always returns true.
func (n *SkillNode) CanHaveChildRows() bool {
	return n.skill.Container()
}

// ChildRows returns the children of this node.
func (n *SkillNode) ChildRows() []unison.TableRowData {
	if n.skill.Container() && n.children == nil {
		n.children = make([]unison.TableRowData, len(n.skill.Children))
		for i, one := range n.skill.Children {
			n.children[i] = NewSkillNode(n.table, n, one)
		}
	}
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *SkillNode) CellDataForSort(index int) string {
	switch index {
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
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *SkillNode) ColumnCell(row, col int, selected bool) unison.Paneler {
	width := n.table.CellWidth(row, col)
	data := n.CellDataForSort(col)
	if n.cellCache[col].matches(width, data) {
		color := unison.DefaultLabelTheme.OnBackgroundInk
		if selected {
			color = unison.OnSelectionColor
		}
		for _, child := range n.cellCache[col].panel.Children() {
			child.Self.(*unison.Label).LabelTheme.OnBackgroundInk = color
		}
		return n.cellCache[col].panel
	}
	p := &unison.Panel{}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{Columns: 1})
	switch col {
	case skillDescriptionColumn:
		createAndAddCellLabel(p, width, n.skill.Description(), unison.DefaultLabelTheme.Font, selected)
		if text := n.skill.SecondaryText(); strings.TrimSpace(text) != "" {
			desc := unison.DefaultLabelTheme.Font.Descriptor()
			desc.Size--
			createAndAddCellLabel(p, width, text, desc.Font(), selected)
		}
	case skillReferenceColumn:
		createAndAddPageRefCellLabel(p, n.CellDataForSort(col), n.skill.Name, unison.DefaultLabelTheme.Font, selected)
	default:
		createAndAddCellLabel(p, width, n.CellDataForSort(col), unison.DefaultLabelTheme.Font, selected)
	}
	n.cellCache[col] = &cellCache{
		width: width,
		data:  data,
		panel: p,
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
	return strings.Contains(strings.ToLower(n.skill.Description()), text) ||
		strings.Contains(strings.ToLower(n.skill.SecondaryText()), text) ||
		(!n.skill.Container() && strings.Contains(strings.ToLower(n.skill.Difficulty.Description(n.skill.Entity)), text)) ||
		strings.Contains(strings.ToLower(n.skill.PageRef), text) ||
		stringSliceContains(n.skill.Categories, text)
}
