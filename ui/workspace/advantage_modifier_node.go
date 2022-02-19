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
	advantageModifierDescriptionColumn = iota
	advantageModifierCostColumn
	advantageModifierReferenceColumn
	advantageModifierColumnCount
)

var _ unison.TableRowData = &AdvantageModifierNode{}

// AdvantageModifierNode holds an advantage modifier in the advantage modifier list.
type AdvantageModifierNode struct {
	table     *unison.Table
	modifier  *gurps.AdvantageModifier
	children  []unison.TableRowData
	cellCache []*cellCache
}

// NewAdvantageModifierListDockable creates a new ListFileDockable for advantage modifier list files.
func NewAdvantageModifierListDockable(filePath string) (*ListFileDockable, error) {
	modifiers, err := gurps.NewAdvantageModifiersFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewListFileDockable(filePath, []unison.TableColumnHeader{
		unison.NewTableColumnHeader(i18n.Text("Modifier")),
		unison.NewTableColumnHeader(i18n.Text("Cost Modifier")),
		newPageReferenceHeader(),
	}, func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(modifiers))
		for _, one := range modifiers {
			rows = append(rows, NewAdvantageModifierNode(table, one))
		}
		return rows
	}), nil
}

// NewAdvantageModifierNode creates a new AdvantageModifierNode.
func NewAdvantageModifierNode(table *unison.Table, modifier *gurps.AdvantageModifier) *AdvantageModifierNode {
	n := &AdvantageModifierNode{
		table:     table,
		modifier:  modifier,
		cellCache: make([]*cellCache, advantageModifierColumnCount),
	}
	return n
}

// CanHaveChildRows always returns true.
func (n *AdvantageModifierNode) CanHaveChildRows() bool {
	return n.modifier.Container()
}

// ChildRows returns the children of this node.
func (n *AdvantageModifierNode) ChildRows() []unison.TableRowData {
	if n.modifier.Container() && n.children == nil {
		n.children = make([]unison.TableRowData, len(n.modifier.Children))
		for i, one := range n.modifier.Children {
			n.children[i] = NewAdvantageModifierNode(n.table, one)
		}
	}
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *AdvantageModifierNode) CellDataForSort(index int) string {
	switch index {
	case advantageModifierDescriptionColumn:
		text := n.modifier.Name
		secondary := n.modifier.SecondaryText()
		if secondary != "" {
			text += "\n" + secondary
		}
		return text
	case advantageModifierCostColumn:
		if n.modifier.Container() {
			return ""
		}
		return n.modifier.CostDescription()
	case advantageModifierReferenceColumn:
		return n.modifier.PageRef
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *AdvantageModifierNode) ColumnCell(row, col int, selected bool) unison.Paneler {
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
	case advantageModifierDescriptionColumn:
		createAndAddCellLabel(p, width, n.modifier.Name, unison.DefaultLabelTheme.Font, selected)
		if text := n.modifier.SecondaryText(); strings.TrimSpace(text) != "" {
			desc := unison.DefaultLabelTheme.Font.Descriptor()
			desc.Size--
			createAndAddCellLabel(p, width, text, desc.Font(), selected)
		}
	case advantageModifierReferenceColumn:
		createAndAddPageRefCellLabel(p, n.CellDataForSort(col), n.modifier.Name, unison.DefaultLabelTheme.Font, selected)
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
func (n *AdvantageModifierNode) IsOpen() bool {
	return n.modifier.Container() && n.modifier.Open
}

// SetOpen sets the current open state for this node.
func (n *AdvantageModifierNode) SetOpen(open bool) {
	if n.modifier.Container() && open != n.modifier.Open {
		n.modifier.Open = open
		n.table.SyncToModel()
		n.table.SizeColumnsToFit(true)
	}
}
