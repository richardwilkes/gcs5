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
	"github.com/richardwilkes/gcs/ui/workspace/gurps/tbl"
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

var (
	_ unison.TableRowData = &AdvantageNode{}
	_ tbl.Matcher         = &AdvantageNode{}
)

// AdvantageNode holds an advantage in the advantage list.
type AdvantageNode struct {
	table     *unison.Table
	parent    *AdvantageNode
	advantage *gurps.Advantage
	children  []unison.TableRowData
	cellCache []*tbl.CellCache
}

// NewAdvantageListDockable creates a new unison.Dockable for advantage list files.
func NewAdvantageListDockable(filePath string) (unison.Dockable, error) {
	advantages, err := gurps.NewAdvantagesFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewListFileDockable(filePath, []unison.TableColumnHeader{
		tbl.NewHeader(i18n.Text("Advantage / Disadvantage"), i18n.Text("Points"), false),
		tbl.NewHeader(i18n.Text("Pts"), "", false),
		tbl.NewHeader(i18n.Text("Type"), "", false),
		tbl.NewHeader(i18n.Text("Category"), "", false),
		tbl.NewPageRefHeader(false),
	}, func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(advantages))
		for _, one := range advantages {
			rows = append(rows, NewAdvantageNode(table, nil, one))
		}
		return rows
	}), nil
}

// NewAdvantageNode creates a new AdvantageNode.
func NewAdvantageNode(table *unison.Table, parent *AdvantageNode, advantage *gurps.Advantage) *AdvantageNode {
	n := &AdvantageNode{
		table:     table,
		parent:    parent,
		advantage: advantage,
		cellCache: make([]*tbl.CellCache, advantageColumnCount),
	}
	return n
}

// ParentRow returns the parent row, or nil if this is a root node.
func (n *AdvantageNode) ParentRow() unison.TableRowData {
	return n.parent
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
			n.children[i] = NewAdvantageNode(n.table, n, one)
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
		tbl.CreateAndAddCellLabel(p, width, n.advantage.Description(), unison.DefaultLabelTheme.Font, selected)
		if text := n.advantage.SecondaryText(); strings.TrimSpace(text) != "" {
			desc := unison.DefaultLabelTheme.Font.Descriptor()
			desc.Size--
			tbl.CreateAndAddCellLabel(p, width, text, desc.Font(), selected)
		}
	case advantagePointsColumn:
		tbl.CreateAndAddCellLabel(p, width, n.CellDataForSort(col), unison.DefaultLabelTheme.Font, selected)
		layout.HAlign = unison.EndAlignment
	case advantageReferenceColumn:
		tbl.CreateAndAddPageRefCellLabel(p, n.CellDataForSort(col), n.advantage.Name, unison.DefaultLabelTheme.Font, selected)
	default:
		tbl.CreateAndAddCellLabel(p, width, n.CellDataForSort(col), unison.DefaultLabelTheme.Font, selected)
	}
	n.cellCache[col] = &tbl.CellCache{
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
	return strings.Contains(strings.ToLower(n.advantage.Description()), text) ||
		strings.Contains(strings.ToLower(n.advantage.SecondaryText()), text) ||
		strings.Contains(strings.ToLower(n.advantage.AdjustedPoints().String()), text) ||
		strings.Contains(strings.ToLower(n.advantage.TypeAsText()), text) ||
		strings.Contains(strings.ToLower(n.advantage.PageRef), text) ||
		stringSliceContains(n.advantage.Categories, text)
}

func stringSliceContains(strs []string, text string) bool {
	for _, s := range strs {
		if strings.Contains(strings.ToLower(s), text) {
			return true
		}
	}
	return false
}
