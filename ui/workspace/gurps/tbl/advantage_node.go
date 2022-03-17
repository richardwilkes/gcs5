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
	advantageDescriptionColumn = iota
	advantagePointsColumn
	advantageTypeColumn
	advantageCategoryColumn
	advantageReferenceColumn
)

var (
	_                   unison.TableRowData = &AdvantageNode{}
	_                   Matcher             = &AdvantageNode{}
	advantageListColMap                     = map[int]int{
		0: advantageDescriptionColumn,
		1: advantagePointsColumn,
		2: advantageTypeColumn,
		3: advantageCategoryColumn,
		4: advantageReferenceColumn,
	}
	advantagePageColMap = map[int]int{
		0: advantageDescriptionColumn,
		1: advantagePointsColumn,
		2: advantageReferenceColumn,
	}
)

// AdvantageNode holds an advantage in the advantage list.
type AdvantageNode struct {
	table     *unison.Table
	parent    *AdvantageNode
	advantage *gurps.Advantage
	children  []unison.TableRowData
	cellCache []*CellCache
	forPage   bool
}

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
		rows := make([]unison.TableRowData, 0, len(topLevelData))
		for _, one := range topLevelData {
			rows = append(rows, NewAdvantageNode(table, nil, one, forPage))
		}
		return rows
	}
}

// NewAdvantageNode creates a new AdvantageNode.
func NewAdvantageNode(table *unison.Table, parent *AdvantageNode, advantage *gurps.Advantage, forPage bool) *AdvantageNode {
	n := &AdvantageNode{
		table:     table,
		parent:    parent,
		advantage: advantage,
		forPage:   forPage,
	}
	n.cellCache = make([]*CellCache, len(n.colMap()))
	return n
}

func (n *AdvantageNode) colMap() map[int]int {
	if n.forPage {
		return advantagePageColMap
	}
	return advantageListColMap
}

// ParentRow returns the parent row, or nil if this is a root node.
func (n *AdvantageNode) ParentRow() unison.TableRowData {
	return n.parent
}

// CanHaveChildRows returns true if this is a container.
func (n *AdvantageNode) CanHaveChildRows() bool {
	return n.advantage.Container()
}

// ChildRows returns the children of this node.
func (n *AdvantageNode) ChildRows() []unison.TableRowData {
	if n.advantage.Container() && n.children == nil {
		n.children = make([]unison.TableRowData, len(n.advantage.Children))
		for i, one := range n.advantage.Children {
			n.children[i] = NewAdvantageNode(n.table, n, one, n.forPage)
		}
	}
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *AdvantageNode) CellDataForSort(index int) string {
	switch n.colMap()[index] {
	case advantageDescriptionColumn:
		text := n.advantage.Description()
		secondary := n.advantage.SecondaryText()
		if secondary != "" {
			text += "\n" + secondary
		}
		return text
	case advantagePointsColumn:
		return n.advantage.AdjustedPoints().String()
	case advantageTypeColumn:
		return n.advantage.TypeAsText()
	case advantageCategoryColumn:
		return strings.Join(n.advantage.Categories, ", ")
	case advantageReferenceColumn:
		return n.advantage.PageRef
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *AdvantageNode) ColumnCell(row, col int, selected bool) unison.Paneler {
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
	case advantageDescriptionColumn:
		CreateAndAddCellLabel(p, width, n.advantage.Description(), primaryFont, selected)
		if text := n.advantage.SecondaryText(); strings.TrimSpace(text) != "" {
			CreateAndAddCellLabel(p, width, text, secondaryFont, selected)
		}
	case advantagePointsColumn:
		CreateAndAddCellLabel(p, width, n.CellDataForSort(col), primaryFont, selected)
		layout.HAlign = unison.EndAlignment
	case advantageReferenceColumn:
		CreateAndAddPageRefCellLabel(p, n.CellDataForSort(col), n.advantage.Name, primaryFont, selected)
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
func (n *AdvantageNode) IsOpen() bool {
	return n.advantage.Container() && n.advantage.Open
}

// SetOpen sets the current open state for this node.
func (n *AdvantageNode) SetOpen(open bool) {
	if n.advantage.Container() && open != n.advantage.Open {
		n.advantage.Open = open
		n.table.SyncToModel()
	}
}

// Match implements Matcher.
func (n *AdvantageNode) Match(text string) bool {
	count := len(n.colMap())
	for i := 0; i < count; i++ {
		if strings.Contains(strings.ToLower(n.CellDataForSort(i)), text) {
			return true
		}
	}
	return false
}
