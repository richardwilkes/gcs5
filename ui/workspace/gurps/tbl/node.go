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
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/settings"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var (
	_ unison.TableRowData = &Node{}
	_ Matcher             = &Node{}
)

// Node represents a row in a table.
type Node struct {
	table     *unison.Table
	parent    unison.TableRowData
	data      node.Node
	children  []unison.TableRowData
	cellCache []*CellCache
	colMap    map[int]int
	forPage   bool
}

// NewNode creates a new node for a table.
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
func (n *Node) ParentRow() unison.TableRowData {
	return n.parent
}

// CanHaveChildRows returns true if this is a container.
func (n *Node) CanHaveChildRows() bool {
	return n.data.Container()
}

// ChildRows returns the children of this node.
func (n *Node) ChildRows() []unison.TableRowData {
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
func (n *Node) ColumnCell(row, col int, selected bool) unison.Paneler {
	var cellData node.CellData
	if column, exists := n.colMap[col]; exists {
		n.data.CellData(column, &cellData)
	}
	width := n.table.CellWidth(row, col)
	if n.cellCache[col].Matches(width, &cellData) {
		applyBackgroundInkRecursively(n.cellCache[col].Panel.AsPanel(), selected)
		return n.cellCache[col].Panel
	}
	cell := n.CellFromCellData(&cellData, width, selected)
	n.cellCache[col] = &CellCache{
		Panel: cell,
		Data:  cellData,
		Width: width,
	}
	return cell
}

func applyBackgroundInkRecursively(panel *unison.Panel, selected bool) {
	if label, ok := panel.Self.(*unison.Label); ok {
		if selected {
			label.OnBackgroundInk = unison.OnSelectionColor
		} else {
			label.OnBackgroundInk = unison.DefaultLabelTheme.OnBackgroundInk
		}
	}
	for _, child := range panel.Children() {
		applyBackgroundInkRecursively(child, selected)
	}
}

// IsOpen returns true if this node should display its children.
func (n *Node) IsOpen() bool {
	return n.data.Container() && n.data.Open()
}

// SetOpen sets the current open state for this node.
func (n *Node) SetOpen(open bool) {
	if n.data.Container() && open != n.data.Open() {
		n.data.SetOpen(open)
		n.table.SyncToModel()
	}
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *Node) CellDataForSort(index int) string {
	if column, exists := n.colMap[index]; exists {
		var data node.CellData
		n.data.CellData(column, &data)
		return data.ForSort()
	}
	return ""
}

// Match implements Matcher.
func (n *Node) Match(text string) bool {
	count := len(n.colMap)
	for i := 0; i < count; i++ {
		if strings.Contains(strings.ToLower(n.CellDataForSort(i)), text) {
			return true
		}
	}
	return false
}

// CellFromCellData creates a new panel for the given cell data.
func (n *Node) CellFromCellData(c *node.CellData, width float32, selected bool) unison.Paneler {
	switch c.Type {
	case node.Text:
		return n.createLabelCell(c, width, selected)
	case node.Toggle:
		return n.createToggleCell(c, selected)
	case node.PageRef:
		return n.createPageRefCell(c.Primary, c.Secondary, selected)
	default:
		return unison.NewPanel()
	}
}

func (n *Node) createLabelCell(c *node.CellData, width float32, selected bool) unison.Paneler {
	p := unison.NewPanel()
	p.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  c.Alignment,
	})
	n.addLabelCell(c, p, width, c.Primary, n.primaryFieldFont(), selected)
	if c.Secondary != "" {
		n.addLabelCell(c, p, width, c.Secondary, n.secondaryFieldFont(), selected)
	}
	if c.Tooltip != "" {
		p.Tooltip = unison.NewTooltipWithText(c.Tooltip)
	}
	return p
}

func (n *Node) addLabelCell(c *node.CellData, parent *unison.Panel, width float32, text string, f unison.Font, selected bool) {
	decoration := &unison.TextDecoration{Font: f}
	var lines []*unison.Text
	if width > 0 {
		lines = unison.NewTextWrappedLines(text, decoration, width)
	} else {
		lines = unison.NewTextLines(text, decoration)
	}
	for _, line := range lines {
		label := unison.NewLabel()
		label.Text = line.String()
		label.Font = f
		label.HAlign = c.Alignment
		if selected {
			label.OnBackgroundInk = unison.OnSelectionColor
		}
		parent.AddChild(label)
	}
}

func (n *Node) createToggleCell(c *node.CellData, selected bool) unison.Paneler {
	check := unison.NewLabel()
	check.Font = n.primaryFieldFont()
	check.SetBorder(unison.NewEmptyBorder(geom32.Insets{Top: 1}))
	baseline := check.Font.Baseline()
	if c.Checked {
		check.Drawable = &unison.DrawableSVG{
			SVG:  res.CheckmarkSVG,
			Size: geom32.Size{Width: baseline, Height: baseline},
		}
	}
	check.HAlign = c.Alignment
	check.VAlign = unison.StartAlignment
	if selected {
		check.OnBackgroundInk = unison.OnSelectionColor
	}
	if c.Tooltip != "" {
		check.Tooltip = unison.NewTooltipWithText(c.Tooltip)
	}
	check.MouseDownCallback = func(where geom32.Point, button, clickCount int, mod unison.Modifiers) bool {
		c.Checked = !c.Checked
		// Currently, there is only one thing, carried equipment, that generates a toggle field and that only generates
		// a single cell, so I'm going to hard-code that into this logic. If we get more, we'll have to find a better
		// approach.
		if equipment, ok := n.data.(*gurps.Equipment); ok {
			equipment.Equipped = c.Checked
			equipment.Entity.Recalculate()
		}
		if c.Checked {
			check.Drawable = &unison.DrawableSVG{
				SVG:  res.CheckmarkSVG,
				Size: geom32.Size{Width: baseline, Height: baseline},
			}
		} else {
			check.Drawable = nil
		}
		check.MarkForLayoutAndRedraw()
		widget.MarkModified(check)
		return true
	}
	return check
}

func (n *Node) createPageRefCell(text, highlight string, selected bool) unison.Paneler {
	label := unison.NewLabel()
	label.Font = n.primaryFieldFont()
	label.VAlign = unison.StartAlignment
	if selected {
		label.OnBackgroundInk = unison.OnSelectionColor
	}
	parts := strings.FieldsFunc(text, func(ch rune) bool { return ch == ',' || ch == ';' || ch == ' ' })
	switch len(parts) {
	case 0:
	case 1:
		label.Text = parts[0]
	default:
		label.Text = parts[0] + "+"
		label.Tooltip = unison.NewTooltipWithText(strings.Join(parts, "\n"))
	}
	if label.Text != "" {
		const isLinkKey = "is_link"
		label.DrawCallback = func(gc *unison.Canvas, rect geom32.Rect) {
			if _, exists := label.ClientData()[isLinkKey]; exists {
				gc.DrawRect(rect, theme.LinkColor.Paint(gc, rect, unison.Fill))
				save := label.OnBackgroundInk
				label.OnBackgroundInk = theme.OnLinkColor
				label.DefaultDraw(gc, rect)
				label.OnBackgroundInk = save
			} else {
				label.DefaultDraw(gc, rect)
			}
		}
		label.MouseEnterCallback = func(where geom32.Point, mod unison.Modifiers) bool {
			label.ClientData()[isLinkKey] = true
			label.MarkForRedraw()
			return true
		}
		label.MouseExitCallback = func() bool {
			delete(label.ClientData(), isLinkKey)
			label.MarkForRedraw()
			return true
		}
		label.MouseDownCallback = func(where geom32.Point, button, clickCount int, mod unison.Modifiers) bool {
			var list []string
			for _, one := range strings.FieldsFunc(text, func(ch rune) bool { return ch == ',' || ch == ';' || ch == ' ' }) {
				if one = strings.TrimSpace(one); one != "" {
					list = append(list, one)
				}
			}
			if len(list) != 0 {
				settings.OpenPageReference(label.Window(), list[0], highlight)
			}
			return true
		}
	}
	return label
}

func (n *Node) primaryFieldFont() unison.Font {
	if n.forPage {
		return theme.PageFieldPrimaryFont
	}
	return unison.FieldFont
}

func (n *Node) secondaryFieldFont() unison.Font {
	if n.forPage {
		return theme.PageFieldSecondaryFont
	}
	return theme.FieldSecondaryFont
}
