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

package workspace

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

const (
	advantageDescriptionColumn = iota
	advantagePointsColumn
	advantageTypeColumn
	advantageCategoryColumn
	advantageReferenceColumn
	advantageColumnCount
)

var _ unison.TableRowData = &AdvantageNode{}

// AdvantageNode holds a advantage in the advantage list.
type AdvantageNode struct {
	table     *unison.Table
	advantage *gurps.Advantage
	children  []unison.TableRowData
	cellCache []*cellCache
}

// NewAdvantageListDockable creates a new ListFileDockable for advantage list files.
func NewAdvantageListDockable(filePath string) (*ListFileDockable, error) {
	advantages, err := gurps.NewAdvantagesFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewListFileDockable(filePath, []unison.TableColumnHeader{
		unison.NewTableColumnHeader(i18n.Text("Advantage / Disadvantage"), ""),
		unison.NewTableColumnHeader(i18n.Text("Pts"), i18n.Text("Points")),
		unison.NewTableColumnHeader(i18n.Text("Type"), ""),
		unison.NewTableColumnHeader(i18n.Text("Category"), ""),
		newPageReferenceHeader(),
	}, func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(advantages))
		for _, one := range advantages {
			rows = append(rows, NewAdvantageNode(table, one))
		}
		return rows
	}), nil
}

// NewAdvantageNode creates a new AdvantageNode.
func NewAdvantageNode(table *unison.Table, advantage *gurps.Advantage) *AdvantageNode {
	n := &AdvantageNode{
		table:     table,
		advantage: advantage,
		cellCache: make([]*cellCache, advantageColumnCount),
	}
	return n
}

// CanHaveChildRows always returns true.
func (n *AdvantageNode) CanHaveChildRows() bool {
	return n.advantage.Container()
}

// ChildRows returns the children of this node.
func (n *AdvantageNode) ChildRows() []unison.TableRowData {
	if n.advantage.Container() && n.children == nil {
		n.children = make([]unison.TableRowData, len(n.advantage.Children))
		for i, one := range n.advantage.Children {
			n.children[i] = NewAdvantageNode(n.table, one)
		}
	}
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *AdvantageNode) CellDataForSort(index int) string {
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
	layout := &unison.FlexLayout{Columns: 1}
	p.SetLayout(layout)
	switch col {
	case advantageDescriptionColumn:
		createAndAddCellLabel(p, width, n.advantage.Description(), unison.DefaultLabelTheme.Font, selected)
		if text := n.advantage.SecondaryText(); strings.TrimSpace(text) != "" {
			desc := unison.DefaultLabelTheme.Font.Descriptor()
			desc.Size--
			createAndAddCellLabel(p, width, text, desc.Font(), selected)
		}
	case advantagePointsColumn:
		createAndAddCellLabel(p, width, n.CellDataForSort(col), unison.DefaultLabelTheme.Font, selected)
		layout.HAlign = unison.EndAlignment
	case advantageReferenceColumn:
		createAndAddPageRefCellLabel(p, n.CellDataForSort(col), n.advantage.Name, unison.DefaultLabelTheme.Font, selected)
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
func (n *AdvantageNode) IsOpen() bool {
	return n.advantage.Container() && n.advantage.Open
}

// SetOpen sets the current open state for this node.
func (n *AdvantageNode) SetOpen(open bool) {
	if n.advantage.Container() && open != n.advantage.Open {
		n.advantage.Open = open
		n.table.SyncToModel()
		n.table.SizeColumnsToFit(true)
	}
}
