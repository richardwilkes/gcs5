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

package sheet

import (
	"strings"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

const (
	advantageDescriptionColumn = iota
	advantagePointsColumn
	advantageReferenceColumn
	advantageColumnCount
)

var (
	_ unison.TableRowData = &AdvantagePageNode{}
	_ tbl.Matcher         = &AdvantagePageNode{}
)

// AdvantagePageNode holds an advantage in the advantage page list.
type AdvantagePageNode struct {
	table     *unison.Table
	parent    *AdvantagePageNode
	advantage *gurps.Advantage
	children  []unison.TableRowData
	cellCache []*tbl.CellCache
}

// NewAdvantagesPageList creates the advantages page list.
func NewAdvantagesPageList(entity *gurps.Entity) *PageList {
	return NewPageList(entity, []unison.TableColumnHeader{
		tbl.NewHeader(i18n.Text("Advantage / Disadvantage"), "", true),
		tbl.NewHeader(i18n.Text("Pts"), i18n.Text("Points"), true),
		tbl.NewPageRefHeader(true),
	}, func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(entity.Advantages))
		for _, one := range entity.Advantages {
			rows = append(rows, NewAdvantagePageNode(table, nil, one))
		}
		return rows
	})
}

// NewAdvantagePageNode creates a new AdvantagePageNode.
func NewAdvantagePageNode(table *unison.Table, parent *AdvantagePageNode, advantage *gurps.Advantage) *AdvantagePageNode {
	n := &AdvantagePageNode{
		table:     table,
		parent:    parent,
		advantage: advantage,
		cellCache: make([]*tbl.CellCache, advantageColumnCount),
	}
	return n
}

// ParentRow returns the parent row, or nil if this is a root node.
func (n *AdvantagePageNode) ParentRow() unison.TableRowData {
	return n.parent
}

// CanHaveChildRows always returns true.
func (n *AdvantagePageNode) CanHaveChildRows() bool {
	return n.advantage.Container()
}

// ChildRows returns the children of this node.
func (n *AdvantagePageNode) ChildRows() []unison.TableRowData {
	if n.advantage.Container() && n.children == nil {
		n.children = make([]unison.TableRowData, len(n.advantage.Children))
		for i, one := range n.advantage.Children {
			n.children[i] = NewAdvantagePageNode(n.table, n, one)
		}
	}
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *AdvantagePageNode) CellDataForSort(index int) string {
	switch index {
	case advantageDescriptionColumn:
		text := n.advantage.Description()
		secondary := n.advantage.SecondaryText()
		if secondary != "" {
			text += "\n" + secondary
		}
		return text
	case advantagePointsColumn:
		return n.advantage.AdjustedPoints().String()
	case advantageReferenceColumn:
		return n.advantage.PageRef
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *AdvantagePageNode) ColumnCell(row, col int, selected bool) unison.Paneler {
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
	switch col {
	case advantageDescriptionColumn:
		tbl.CreateAndAddCellLabel(p, width, n.advantage.Description(), theme.PageFieldPrimaryFont, selected)
		if text := n.advantage.SecondaryText(); strings.TrimSpace(text) != "" {
			tbl.CreateAndAddCellLabel(p, width, text, theme.PageFieldSecondaryFont, selected)
		}
	case advantagePointsColumn:
		tbl.CreateAndAddCellLabel(p, width, n.CellDataForSort(col), theme.PageFieldPrimaryFont, selected)
		layout.HAlign = unison.EndAlignment
	case advantageReferenceColumn:
		tbl.CreateAndAddPageRefCellLabel(p, n.CellDataForSort(col), n.advantage.Name, theme.PageFieldPrimaryFont, selected)
	default:
		tbl.CreateAndAddCellLabel(p, width, n.CellDataForSort(col), theme.PageFieldPrimaryFont, selected)
	}
	n.cellCache[col] = &tbl.CellCache{
		Width: width,
		Data:  data,
		Panel: p,
	}
	return p
}

// IsOpen returns true if this node should display its children.
func (n *AdvantagePageNode) IsOpen() bool {
	return n.advantage.Container() && n.advantage.Open
}

// SetOpen sets the current open state for this node.
func (n *AdvantagePageNode) SetOpen(open bool) {
	if n.advantage.Container() && open != n.advantage.Open {
		n.advantage.Open = open
		n.table.SyncToModel()
	}
}

// Match implements Matcher.
func (n *AdvantagePageNode) Match(text string) bool {
	return strings.Contains(strings.ToLower(n.advantage.Description()), text) ||
		strings.Contains(strings.ToLower(n.advantage.SecondaryText()), text) ||
		strings.Contains(strings.ToLower(n.advantage.AdjustedPoints().String()), text) ||
		strings.Contains(strings.ToLower(n.advantage.PageRef), text)
}
