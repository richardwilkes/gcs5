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

	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/unison"
)

var (
	_ unison.TableRowData = &Node{}
	_ Matcher             = &Node{}
)

type Node struct {
	table     *unison.Table
	parent    unison.TableRowData
	data      node.Node
	children  []unison.TableRowData
	cellCache []*CellCache
	colMap    map[int]int
	forPage   bool
}

func NewNode(table *unison.Table, parent unison.TableRowData, colMap map[int]int, data node.Node, forPage bool) *Node {
	return &Node{
		table:     table,
		parent:    parent,
		data:      data,
		cellCache: make([]*CellCache, len(colMap)),
		colMap:    colMap,
		forPage:   forPage,
	}
}

// ParentRow returns the parent row, or nil if this is a root node.
func (n Node) ParentRow() unison.TableRowData {
	return n.parent
}

// CanHaveChildRows returns true if this is a container.
func (n Node) CanHaveChildRows() bool {
	return n.data.Container()
}

// ChildRows returns the children of this node.
func (n Node) ChildRows() []unison.TableRowData {
	if n.data.Container() && n.children == nil {
		children := n.data.NodeChildren()
		n.children = make([]unison.TableRowData, len(children))
		for i, one := range children {
			n.children[i] = NewNode(n.table, n, n.colMap, one, n.forPage)
		}
	}
	return n.children
}

// ColumnCell returns the cell for the given column index.
func (n Node) ColumnCell(row, col int, selected bool) unison.Paneler {
	var cellData node.CellData
	if column, exists := n.colMap[col]; exists {
		n.data.CellData(column, &cellData)
	}
	width := n.table.CellWidth(row, col)
	if n.cellCache[col].Matches(width, &cellData) {
		color := unison.DefaultLabelTheme.OnBackgroundInk
		if selected {
			color = unison.OnSelectionColor
		}
		for _, child := range n.cellCache[col].Panel.Children() {
			if label, ok := child.Self.(*unison.Label); ok {
				label.OnBackgroundInk = color
			}
		}
		return n.cellCache[col].Panel
	}
	cell := CellFromCellData(&cellData, width, n.forPage, selected)
	n.cellCache[col] = &CellCache{
		Panel: cell,
		Data:  cellData,
		Width: width,
	}
	return cell
}

// IsOpen returns true if this node should display its children.
func (n Node) IsOpen() bool {
	return n.data.Container() && n.data.Open()
}

// SetOpen sets the current open state for this node.
func (n Node) SetOpen(open bool) {
	if n.data.Container() && open != n.data.Open() {
		n.data.SetOpen(open)
		n.table.SyncToModel()
	}
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n Node) CellDataForSort(index int) string {
	if column, exists := n.colMap[index]; exists {
		var data node.CellData
		n.data.CellData(column, &data)
		return data.ForSort()
	}
	return ""
}

// Match implements Matcher.
func (n Node) Match(text string) bool {
	count := len(n.colMap)
	for i := 0; i < count; i++ {
		if strings.Contains(strings.ToLower(n.CellDataForSort(i)), text) {
			return true
		}
	}
	return false
}
