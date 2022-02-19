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
	"strconv"
	"strings"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/icons"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

const (
	equipmentDescriptionColumn = iota
	equipmentUsesColumn
	equipmentTLColumn
	equipmentLCColumn
	equipmentCostColumn
	equipmentWeightColumn
	equipmentCategoryColumn
	equipmentReferenceColumn
	equipmentColumnCount
)

var _ unison.TableRowData = &EquipmentNode{}

// EquipmentNode holds equipment in the equipment list.
type EquipmentNode struct {
	table     *unison.Table
	equipment *gurps.Equipment
	children  []unison.TableRowData
	cellCache []*cellCache
}

// NewEquipmentListDockable creates a new ListFileDockable for equipment list files.
func NewEquipmentListDockable(filePath string) (*ListFileDockable, error) {
	eqp, err := gurps.NewEquipmentFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	tlHdr := unison.NewTableColumnHeader(i18n.Text("TL"))
	tlHdr.Tooltip = unison.NewTooltipWithText(i18n.Text("Tech Level"))
	lcHdr := unison.NewTableColumnHeader(i18n.Text("LC"))
	lcHdr.Tooltip = unison.NewTooltipWithText(i18n.Text("Legality Class"))
	costHdr := unison.NewTableColumnHeader(i18n.Text("$"))
	costHdr.Tooltip = unison.NewTooltipWithText(i18n.Text("Cost"))
	weightHdr := unison.NewTableColumnHeader("")
	baseline := weightHdr.Font.Baseline()
	weightHdr.Drawable = &unison.DrawableSVG{
		SVG:  icons.WeightSVG(),
		Size: geom32.NewSize(baseline, baseline),
	}
	weightHdr.Tooltip = unison.NewTooltipWithText(i18n.Text("Weight"))
	return NewListFileDockable(filePath, []unison.TableColumnHeader{
		unison.NewTableColumnHeader(i18n.Text("Equipment")),
		unison.NewTableColumnHeader(i18n.Text("Uses")),
		tlHdr,
		lcHdr,
		costHdr,
		weightHdr,
		unison.NewTableColumnHeader(i18n.Text("Category")),
		newPageReferenceHeader(),
	}, func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(eqp))
		for _, one := range eqp {
			rows = append(rows, NewEquipmentNode(table, one))
		}
		return rows
	}), nil
}

// NewEquipmentNode creates a new EquipmentNode.
func NewEquipmentNode(table *unison.Table, eqp *gurps.Equipment) *EquipmentNode {
	n := &EquipmentNode{
		table:     table,
		equipment: eqp,
		cellCache: make([]*cellCache, equipmentColumnCount),
	}
	return n
}

// CanHaveChildRows always returns true.
func (n *EquipmentNode) CanHaveChildRows() bool {
	return n.equipment.Container()
}

// ChildRows returns the children of this node.
func (n *EquipmentNode) ChildRows() []unison.TableRowData {
	if n.equipment.Container() && n.children == nil {
		n.children = make([]unison.TableRowData, len(n.equipment.Children))
		for i, one := range n.equipment.Children {
			n.children[i] = NewEquipmentNode(n.table, one)
		}
	}
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *EquipmentNode) CellDataForSort(index int) string {
	switch index {
	case equipmentDescriptionColumn:
		text := n.equipment.Description()
		secondary := n.equipment.SecondaryText()
		if secondary != "" {
			text += "\n" + secondary
		}
		return text
	case equipmentUsesColumn:
		if n.equipment.MaxUses > 0 {
			return strconv.Itoa(n.equipment.Uses)
		}
		return ""
	case equipmentTLColumn:
		return n.equipment.TechLevel
	case equipmentLCColumn:
		return n.equipment.LegalityClass
	case equipmentCostColumn:
		return n.equipment.AdjustedValue().String()
	case equipmentWeightColumn:
		return n.equipment.AdjustedWeight(false, gurps.SheetSettingsFor(n.equipment.Entity).DefaultWeightUnits).String()
	case equipmentCategoryColumn:
		return strings.Join(n.equipment.Categories, ", ")
	case equipmentReferenceColumn:
		return n.equipment.PageRef
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *EquipmentNode) ColumnCell(row, col int, selected bool) unison.Paneler {
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
	case equipmentDescriptionColumn:
		createAndAddCellLabel(p, width, n.equipment.Description(), unison.DefaultLabelTheme.Font, selected)
		if text := n.equipment.SecondaryText(); strings.TrimSpace(text) != "" {
			desc := unison.DefaultLabelTheme.Font.Descriptor()
			desc.Size--
			createAndAddCellLabel(p, width, text, desc.Font(), selected)
		}
	case equipmentUsesColumn, equipmentTLColumn, equipmentLCColumn, equipmentCostColumn, equipmentWeightColumn:
		createAndAddCellLabel(p, width, n.CellDataForSort(col), unison.DefaultLabelTheme.Font, selected)
		layout.HAlign = unison.EndAlignment
	case equipmentReferenceColumn:
		createAndAddPageRefCellLabel(p, n.CellDataForSort(col), n.equipment.Name, unison.DefaultLabelTheme.Font, selected)
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
func (n *EquipmentNode) IsOpen() bool {
	return n.equipment.Container() && n.equipment.Open
}

// SetOpen sets the current open state for this node.
func (n *EquipmentNode) SetOpen(open bool) {
	if n.equipment.Container() && open != n.equipment.Open {
		n.equipment.Open = open
		n.table.SyncToModel()
		n.table.SizeColumnsToFit(true)
	}
}
