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
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

const (
	equipmentModifierDescriptionColumn = iota
	equipmentModifierTechLevelColumn
	equipmentModifierCostColumn
	equipmentModifierWeightColumn
	equipmentModifierCategoryColumn
	equipmentModifierReferenceColumn
	equipmentModifierColumnCount
)

var (
	_ unison.TableRowData = &EquipmentModifierNode{}
	_ tbl.Matcher         = &EquipmentModifierNode{}
)

// EquipmentModifierNode holds an equipment modifier in the equipment modifier list.
type EquipmentModifierNode struct {
	table     *unison.Table
	parent    *EquipmentModifierNode
	modifier  *gurps.EquipmentModifier
	children  []unison.TableRowData
	cellCache []*tbl.CellCache
}

// NewEquipmentModifierListDockable creates a new unison.Dockable for equipment modifier list files.
func NewEquipmentModifierListDockable(filePath string) (unison.Dockable, error) {
	modifiers, err := gurps.NewEquipmentModifiersFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewTableDockable(filePath, []unison.TableColumnHeader{
		tbl.NewHeader(i18n.Text("Modifier"), "", false),
		tbl.NewHeader(i18n.Text("TL"), i18n.Text("Tech Level"), false),
		tbl.NewHeader(i18n.Text("Cost Adjustment"), "", false),
		tbl.NewHeader(i18n.Text("Weight Adjustment"), "", false),
		tbl.NewHeader(i18n.Text("Category"), "", false),
		tbl.NewPageRefHeader(false),
	}, func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(modifiers))
		for _, one := range modifiers {
			rows = append(rows, NewEquipmentModifierNode(table, nil, one))
		}
		return rows
	}), nil
}

// NewEquipmentModifierNode creates a new EquipmentModifierNode.
func NewEquipmentModifierNode(table *unison.Table, parent *EquipmentModifierNode, modifier *gurps.EquipmentModifier) *EquipmentModifierNode {
	n := &EquipmentModifierNode{
		table:     table,
		parent:    parent,
		modifier:  modifier,
		cellCache: make([]*tbl.CellCache, equipmentModifierColumnCount),
	}
	return n
}

// ParentRow returns the parent row, or nil if this is a root node.
func (n *EquipmentModifierNode) ParentRow() unison.TableRowData {
	return n.parent
}

// CanHaveChildRows always returns true.
func (n *EquipmentModifierNode) CanHaveChildRows() bool {
	return n.modifier.Container()
}

// ChildRows returns the children of this node.
func (n *EquipmentModifierNode) ChildRows() []unison.TableRowData {
	if n.modifier.Container() && n.children == nil {
		n.children = make([]unison.TableRowData, len(n.modifier.Children))
		for i, one := range n.modifier.Children {
			n.children[i] = NewEquipmentModifierNode(n.table, n, one)
		}
	}
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *EquipmentModifierNode) CellDataForSort(index int) string {
	switch index {
	case equipmentModifierDescriptionColumn:
		text := n.modifier.Name
		secondary := n.modifier.SecondaryText()
		if secondary != "" {
			text += "\n" + secondary
		}
		return text
	case equipmentModifierTechLevelColumn:
		if n.modifier.Container() {
			return ""
		}
		return n.modifier.TechLevel
	case equipmentModifierCostColumn:
		if n.modifier.Container() {
			return ""
		}
		return n.modifier.CostDescription()
	case equipmentModifierWeightColumn:
		if n.modifier.Container() {
			return ""
		}
		return n.modifier.WeightDescription()
	case equipmentModifierCategoryColumn:
		return strings.Join(n.modifier.Categories, ", ")
	case equipmentModifierReferenceColumn:
		return n.modifier.PageRef
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *EquipmentModifierNode) ColumnCell(row, col int, selected bool) unison.Paneler {
	width := n.table.CellWidth(row, col)
	data := n.CellDataForSort(col)
	if n.cellCache[col].MatchesOld(width, data) {
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
	case equipmentModifierDescriptionColumn:
		tbl.CreateAndAddCellLabel(p, width, n.modifier.Name, unison.DefaultLabelTheme.Font, selected)
		if text := n.modifier.SecondaryText(); strings.TrimSpace(text) != "" {
			desc := unison.DefaultLabelTheme.Font.Descriptor()
			desc.Size--
			tbl.CreateAndAddCellLabel(p, width, text, desc.Font(), selected)
		}
	case equipmentModifierReferenceColumn:
		tbl.CreateAndAddPageRefCellLabel(p, n.CellDataForSort(col), n.modifier.Name, unison.DefaultLabelTheme.Font, selected)
	default:
		tbl.CreateAndAddCellLabel(p, width, n.CellDataForSort(col), unison.DefaultLabelTheme.Font, selected)
	}
	n.cellCache[col] = &tbl.CellCache{
		Width: width,
		Data:  node.CellData{Primary: data},
		Panel: p,
	}
	return p
}

// IsOpen returns true if this node should display its children.
func (n *EquipmentModifierNode) IsOpen() bool {
	return n.modifier.Container() && n.modifier.Open
}

// SetOpen sets the current open state for this node.
func (n *EquipmentModifierNode) SetOpen(open bool) {
	if n.modifier.Container() && open != n.modifier.Open {
		n.modifier.Open = open
		n.table.SyncToModel()
	}
}

// Match implements Matcher.
func (n *EquipmentModifierNode) Match(text string) bool {
	return strings.Contains(strings.ToLower(n.modifier.Name), text) ||
		strings.Contains(strings.ToLower(n.modifier.SecondaryText()), text) ||
		(!n.modifier.Container() &&
			(strings.Contains(strings.ToLower(n.modifier.TechLevel), text) ||
				strings.Contains(strings.ToLower(n.modifier.CostDescription()), text) ||
				strings.Contains(strings.ToLower(n.modifier.WeightDescription()), text))) ||
		strings.Contains(strings.ToLower(n.modifier.PageRef), text) ||
		stringSliceContains(n.modifier.Categories, text)
}
