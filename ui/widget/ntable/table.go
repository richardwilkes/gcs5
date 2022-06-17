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

package ntable

import (
	"github.com/richardwilkes/gcs/constants"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/unison"
)

// ItemVariant holds the type of item variant to create.
type ItemVariant int

// Possible values for ItemVariant.
const (
	NoItemVariant ItemVariant = iota
	ContainerItemVariant
	AlternateItemVariant
)

// TableProvider defines the methods a table provider must contain.
type TableProvider[T gurps.NodeConstraint[T]] interface {
	unison.TableModel[*Node[T]]
	gurps.EntityProvider
	SetTable(table *unison.Table[*Node[T]])
	DragKey() string
	DragSVG() *unison.SVG
	DropShouldMoveData(from, to *unison.Table[*Node[T]]) bool
	ItemNames() (singular, plural string)
	Headers() []unison.TableColumnHeader[*Node[T]]
	SyncHeader(headers []unison.TableColumnHeader[*Node[T]])
	HierarchyColumnIndex() int
	ExcessWidthColumnIndex() int
	OpenEditor(owner widget.Rebuildable, table *unison.Table[*Node[T]])
	CreateItem(owner widget.Rebuildable, table *unison.Table[*Node[T]], variant ItemVariant)
	DuplicateSelection(table *unison.Table[*Node[T]])
	DeleteSelection(table *unison.Table[*Node[T]])
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}

// NewNodeTable creates a new node table of the specified type, returning the header and table. Pass nil for 'font' if
// this should be a standalone top-level table for a dockable. Otherwise, pass in the typical font used for a cell.
func NewNodeTable[T gurps.NodeConstraint[T]](provider TableProvider[T], font unison.Font) (header *unison.TableHeader[*Node[T]], table *unison.Table[*Node[T]]) {
	table = unison.NewTable[*Node[T]](provider)
	provider.SetTable(table)
	table.HierarchyColumnIndex = provider.HierarchyColumnIndex()
	table.DividerInk = theme.HeaderColor
	layoutData := &unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	}
	if font != nil {
		table.Padding.Top = 0
		table.Padding.Bottom = 0
		table.HierarchyIndent = font.LineHeight()
		table.MinimumRowHeight = font.LineHeight()
		layoutData.MinSize = unison.Size{Height: 4 + theme.PageFieldPrimaryFont.LineHeight()}
	}
	table.SetLayoutData(layoutData)

	headers := provider.Headers()
	table.ColumnSizes = make([]unison.ColumnSize, len(headers))
	for i := range table.ColumnSizes {
		_, pref, _ := headers[i].AsPanel().Sizes(unison.Size{})
		pref.Width += table.Padding.Left + table.Padding.Right
		table.ColumnSizes[i].AutoMinimum = pref.Width
		table.ColumnSizes[i].AutoMaximum = 800
		table.ColumnSizes[i].Minimum = pref.Width
		table.ColumnSizes[i].Maximum = 10000
	}
	header = unison.NewTableHeader(table, headers...)
	header.Less = flexibleLess
	if font != nil {
		header.BackgroundInk = theme.HeaderColor
		header.DividerInk = theme.HeaderColor
		header.HeaderBorder = unison.NewLineBorder(theme.HeaderColor, 0, unison.Insets{Bottom: 1}, false)
		header.SetBorder(header.HeaderBorder)
	}
	header.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})

	mouseDownCallback := table.MouseDownCallback
	table.MouseDownCallback = func(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
		table.RequestFocus()
		return mouseDownCallback(where, button, clickCount, mod)
	}
	table.DoubleClickCallback = func() { table.PerformCmd(nil, constants.OpenEditorItemID) }
	keydownCallback := table.KeyDownCallback
	table.KeyDownCallback = func(keyCode unison.KeyCode, mod unison.Modifiers, repeat bool) bool {
		if mod == 0 && (keyCode == unison.KeyBackspace || keyCode == unison.KeyDelete) {
			table.PerformCmd(table, unison.DeleteItemID)
			return true
		}
		return keydownCallback(keyCode, mod, repeat)
	}
	singular, plural := provider.ItemNames()
	table.InstallDragSupport(provider.DragSVG(), provider.DragKey(), singular, plural)
	if font != nil {
		table.FrameChangeCallback = func() {
			table.SizeColumnsToFitWithExcessIn(provider.ExcessWidthColumnIndex())
		}
	}
	return header, table
}

func flexibleLess(s1, s2 string) bool {
	if n1, err := fxp.FromString(s1); err == nil {
		var n2 fxp.Int
		if n2, err = fxp.FromString(s2); err == nil {
			return n1 < n2
		}
	}
	return txt.NaturalLess(s1, s2, true)
}
