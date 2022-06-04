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

package editors

import (
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/settings"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/unison"
	"golang.org/x/exp/slices"
)

const excludeMarker = "exclude"

var (
	_ unison.TableRowData = &Node{}
	_ Matcher             = &Node{}
)

// Node represents a row in a table.
type Node struct {
	table     *unison.Table
	parent    unison.TableRowData
	data      gurps.Node
	children  []unison.TableRowData
	cellCache []*CellCache
	colMap    map[int]int
	forPage   bool
}

// NewNode creates a new node for a table.
func NewNode(table *unison.Table, parent unison.TableRowData, colMap map[int]int, data gurps.Node, forPage bool) *Node {
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

// Data returns the underlying data object.
func (n *Node) Data() gurps.Node {
	return n.data
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
func (n *Node) ColumnCell(row, col int, foreground, _ unison.Ink, _, _, _ bool) unison.Paneler {
	var cellData gurps.CellData
	if column, exists := n.colMap[col]; exists {
		n.data.CellData(column, &cellData)
	}
	width := n.table.CellWidth(row, col)
	if n.cellCache[col].Matches(width, &cellData) {
		applyForegroundInkRecursively(n.cellCache[col].Panel.AsPanel(), foreground)
		return n.cellCache[col].Panel
	}
	cell := n.CellFromCellData(&cellData, width, foreground)
	n.cellCache[col] = &CellCache{
		Panel: cell,
		Data:  cellData,
		Width: width,
	}
	return cell
}

func applyForegroundInkRecursively(panel *unison.Panel, foreground unison.Ink) {
	if label, ok := panel.Self.(*unison.Label); ok {
		if _, exists := label.ClientData()[excludeMarker]; !exists {
			label.OnBackgroundInk = foreground
		}
	}
	for _, child := range panel.Children() {
		applyForegroundInkRecursively(child, foreground)
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
		var data gurps.CellData
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
func (n *Node) CellFromCellData(c *gurps.CellData, width float32, foreground unison.Ink) unison.Paneler {
	switch c.Type {
	case gurps.Text:
		return n.createLabelCell(c, width, foreground)
	case gurps.Toggle:
		return n.createToggleCell(c, foreground)
	case gurps.PageRef:
		return n.createPageRefCell(c, foreground)
	default:
		return unison.NewPanel()
	}
}

func (n *Node) createLabelCell(c *gurps.CellData, width float32, foreground unison.Ink) unison.Paneler {
	p := unison.NewPanel()
	p.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  c.Alignment,
	})
	n.addLabelCell(c, p, width, c.Primary, n.primaryFieldFont(), foreground, true)
	if c.Secondary != "" {
		n.addLabelCell(c, p, width, c.Secondary, n.secondaryFieldFont(), foreground, false)
	}
	tooltip := c.Tooltip
	if c.UnsatisfiedReason != "" {
		label := unison.NewLabel()
		label.Font = n.secondaryFieldFont()
		baseline := label.Font.Baseline()
		label.Drawable = &unison.DrawableSVG{
			SVG:  unison.TriangleExclamationSVG(),
			Size: geom.NewSize(baseline, baseline),
		}
		label.Text = i18n.Text("Unsatisfied prerequisite(s)")
		label.HAlign = c.Alignment
		label.ClientData()[excludeMarker] = true
		label.OnBackgroundInk = unison.OnErrorColor
		label.SetBorder(unison.NewEmptyBorder(unison.Insets{
			Top:    1,
			Left:   4,
			Bottom: 0,
			Right:  4,
		}))
		label.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
			gc.DrawRect(rect, unison.ErrorColor.Paint(gc, rect, unison.Fill))
			label.DefaultDraw(gc, rect)
		}
		p.AddChild(label)
		tooltip = c.UnsatisfiedReason
	}
	if tooltip != "" {
		p.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	return p
}

func (n *Node) addLabelCell(c *gurps.CellData, parent *unison.Panel, width float32, text string, f unison.Font, foreground unison.Ink, primary bool) {
	decoration := &unison.TextDecoration{
		Font:          f,
		StrikeThrough: primary && c.Disabled,
	}
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
		label.StrikeThrough = primary && c.Disabled
		label.HAlign = c.Alignment
		label.OnBackgroundInk = foreground
		label.SetEnabled(!c.Dim)
		parent.AddChild(label)
	}
}

func (n *Node) createToggleCell(c *gurps.CellData, foreground unison.Ink) unison.Paneler {
	check := unison.NewLabel()
	check.Font = n.primaryFieldFont()
	check.SetBorder(unison.NewEmptyBorder(unison.Insets{Top: 1}))
	baseline := check.Font.Baseline()
	if c.Checked {
		check.Drawable = &unison.DrawableSVG{
			SVG:  res.CheckmarkSVG,
			Size: unison.Size{Width: baseline, Height: baseline},
		}
	}
	check.HAlign = c.Alignment
	check.VAlign = unison.StartAlignment
	check.OnBackgroundInk = foreground
	if c.Tooltip != "" {
		check.Tooltip = unison.NewTooltipWithText(c.Tooltip)
	}
	check.MouseDownCallback = func(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
		c.Checked = !c.Checked
		switch item := n.data.(type) {
		case *gurps.Equipment:
			item.Equipped = c.Checked
			if mgr := unison.UndoManagerFor(check); mgr != nil {
				owner := widget.FindRebuildable(check)
				mgr.Add(&unison.UndoEdit[*equipmentAdjuster]{
					ID:       unison.NextUndoID(),
					EditName: i18n.Text("Toggle Equipped"),
					UndoFunc: func(edit *unison.UndoEdit[*equipmentAdjuster]) { edit.BeforeData.Apply() },
					RedoFunc: func(edit *unison.UndoEdit[*equipmentAdjuster]) { edit.AfterData.Apply() },
					BeforeData: &equipmentAdjuster{
						Owner:    owner,
						Target:   item,
						Equipped: !item.Equipped,
					},
					AfterData: &equipmentAdjuster{
						Owner:    owner,
						Target:   item,
						Equipped: item.Equipped,
					},
				})
			}
			item.Entity.Recalculate()
		case *gurps.TraitModifier:
			item.Disabled = !c.Checked
			if mgr := unison.UndoManagerFor(check); mgr != nil {
				owner := widget.FindRebuildable(check)
				mgr.Add(&unison.UndoEdit[*traitModifierAdjuster]{
					ID:       unison.NextUndoID(),
					EditName: i18n.Text("Toggle Trait Modifier"),
					UndoFunc: func(edit *unison.UndoEdit[*traitModifierAdjuster]) { edit.BeforeData.Apply() },
					RedoFunc: func(edit *unison.UndoEdit[*traitModifierAdjuster]) { edit.AfterData.Apply() },
					BeforeData: &traitModifierAdjuster{
						Owner:    owner,
						Target:   item,
						Disabled: !item.Disabled,
					},
					AfterData: &traitModifierAdjuster{
						Owner:    owner,
						Target:   item,
						Disabled: item.Disabled,
					},
				})
			}
			item.Entity.Recalculate()
		case *gurps.EquipmentModifier:
			item.Disabled = !c.Checked
			if mgr := unison.UndoManagerFor(check); mgr != nil {
				owner := widget.FindRebuildable(check)
				mgr.Add(&unison.UndoEdit[*equipmentModifierAdjuster]{
					ID:       unison.NextUndoID(),
					EditName: i18n.Text("Toggle Equipment Modifier"),
					UndoFunc: func(edit *unison.UndoEdit[*equipmentModifierAdjuster]) { edit.BeforeData.Apply() },
					RedoFunc: func(edit *unison.UndoEdit[*equipmentModifierAdjuster]) { edit.AfterData.Apply() },
					BeforeData: &equipmentModifierAdjuster{
						Owner:    owner,
						Target:   item,
						Disabled: !item.Disabled,
					},
					AfterData: &equipmentModifierAdjuster{
						Owner:    owner,
						Target:   item,
						Disabled: item.Disabled,
					},
				})
			}
			item.Entity.Recalculate()
		}
		if c.Checked {
			check.Drawable = &unison.DrawableSVG{
				SVG:  res.CheckmarkSVG,
				Size: unison.Size{Width: baseline, Height: baseline},
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

type equipmentAdjuster struct {
	Owner    widget.Rebuildable
	Target   *gurps.Equipment
	Equipped bool
}

func (a *equipmentAdjuster) Apply() {
	a.Target.Equipped = a.Equipped
	if a.Target.Entity != nil {
		a.Target.Entity.Recalculate()
	}
	widget.MarkModified(a.Owner)
}

type equipmentModifierAdjuster struct {
	Owner    widget.Rebuildable
	Target   *gurps.EquipmentModifier
	Disabled bool
}

func (a *equipmentModifierAdjuster) Apply() {
	a.Target.Disabled = a.Disabled
	if a.Target.Entity != nil {
		a.Target.Entity.Recalculate()
	}
	widget.MarkModified(a.Owner)
}

type traitModifierAdjuster struct {
	Owner    widget.Rebuildable
	Target   *gurps.TraitModifier
	Disabled bool
}

func (a *traitModifierAdjuster) Apply() {
	a.Target.Disabled = a.Disabled
	if a.Target.Entity != nil {
		a.Target.Entity.Recalculate()
	}
	widget.MarkModified(a.Owner)
}

func (n *Node) createPageRefCell(c *gurps.CellData, foreground unison.Ink) unison.Paneler {
	label := unison.NewLabel()
	label.Font = n.primaryFieldFont()
	label.VAlign = unison.StartAlignment
	label.OnBackgroundInk = foreground
	label.SetEnabled(!c.Dim)
	parts := strings.FieldsFunc(c.Primary, func(ch rune) bool { return ch == ',' || ch == ';' || ch == ' ' })
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
		label.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
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
		label.MouseEnterCallback = func(where unison.Point, mod unison.Modifiers) bool {
			label.ClientData()[isLinkKey] = true
			label.MarkForRedraw()
			return true
		}
		label.MouseExitCallback = func() bool {
			delete(label.ClientData(), isLinkKey)
			label.MarkForRedraw()
			return true
		}
		label.MouseDownCallback = func(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
			list := settings.ExtractPageReferences(c.Primary)
			if len(list) != 0 {
				settings.OpenPageReference(label.Window(), list[0], c.Secondary, nil)
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

// FindRowIndexByID returns the row index of the row with the given ID in the given table.
func FindRowIndexByID(table *unison.Table, id uuid.UUID) int {
	_, i := rowIndex(id, 0, table.TopLevelRows())
	return i
}

func rowIndex(id uuid.UUID, startIndex int, rows []unison.TableRowData) (updatedStartIndex, result int) {
	for _, row := range rows {
		if n, ok := row.(*Node); ok {
			if id == n.Data().UUID() {
				return 0, startIndex
			}
		}
		startIndex++
		if row.IsOpen() {
			if startIndex, result = rowIndex(id, startIndex, row.ChildRows()); result != -1 {
				return 0, result
			}
		}
	}
	return startIndex, -1
}

// InsertItem inserts an item into a table.
func InsertItem[T comparable](owner widget.Rebuildable, table *unison.Table, item T, setParent func(target, parent T), childrenOf func(target T) []T, setChildren func(target T, children []T), topList func() []T, setTopList func([]T), rowData func(table *unison.Table) []unison.TableRowData, id func(T) uuid.UUID) {
	var target, zero T
	i := table.FirstSelectedRowIndex()
	if i != -1 {
		row := table.RowFromIndex(i)
		if target = ExtractFromRowData[T](row); target != zero {
			if row.CanHaveChildRows() {
				// Target is container, append to end of that container
				setParent(item, target)
				setChildren(target, append(childrenOf(target), item))
			} else {
				// Target isn't a container. If it has a parent, insert after the target within that parent.
				if parent := ExtractFromRowData[T](row.ParentRow()); parent != zero {
					setParent(item, parent)
					children := childrenOf(parent)
					setChildren(parent, slices.Insert(children, slices.Index(children, target)+1, item))
				} else {
					// Otherwise, insert after the target within the top-level list.
					setParent(item, zero)
					list := topList()
					setTopList(slices.Insert(list, slices.Index(list, target)+1, item))
				}
			}
		}
	}
	if target == zero {
		// There was no selection, so append to the end of the top-level list.
		setParent(item, zero)
		setTopList(append(topList(), item))
	}
	widget.MarkModified(table)
	table.SetTopLevelRows(rowData(table))
	index := FindRowIndexByID(table, id(item))
	table.SelectByIndex(index)
	table.ScrollRowCellIntoView(index, 0)
	table.RequestFocus()
	owner.Rebuild(true)
}

// ExtractFromRowData extracts a specific type of data from the row data.
func ExtractFromRowData[T any](row unison.TableRowData) T {
	if n, ok := row.(*Node); ok {
		var target T
		if target, ok = n.Data().(T); ok {
			return target
		}
	}
	var zero T
	return zero
}

// OpenEditor opens an editor for each selected row in the table.
func OpenEditor[T comparable](table *unison.Table, edit func(item T)) {
	var zero T
	for _, row := range table.SelectedRows(false) {
		if a := ExtractFromRowData[T](row); a != zero {
			edit(a)
		}
	}
}
