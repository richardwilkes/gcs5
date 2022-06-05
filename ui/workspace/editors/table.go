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
	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/constants"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/unison"
	"golang.org/x/exp/slices"
)

type comparableNode interface {
	comparable
	gurps.Node
}

func newTable(parent *unison.Panel, provider TableProvider) *unison.Table {
	table := unison.NewTable()
	table.DividerInk = theme.HeaderColor
	table.Padding.Top = 0
	table.Padding.Bottom = 0
	table.HierarchyColumnIndex = provider.HierarchyColumnIndex()
	table.HierarchyIndent = unison.FieldFont.LineHeight()
	table.MinimumRowHeight = unison.FieldFont.LineHeight()
	headers := provider.Headers()
	widget.TableSetupColumnSizes(table, headers)
	table.SetLayoutData(&unison.FlexLayoutData{
		MinSize: geom.Size[float32]{Height: unison.FieldFont.LineHeight()},
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
		HGrab:   true,
		VGrab:   true,
	})
	widget.TableInstallStdCallbacks(table)
	table.FrameChangeCallback = func() {
		table.SizeColumnsToFitWithExcessIn(provider.ExcessWidthColumnIndex())
	}
	tableHeader := widget.TableCreateHeader(table, headers)
	tableHeader.BackgroundInk = theme.HeaderColor
	tableHeader.DividerInk = theme.HeaderColor
	tableHeader.HeaderBorder = unison.NewLineBorder(theme.HeaderColor, 0, unison.Insets{Bottom: 1}, false)
	tableHeader.SetBorder(tableHeader.HeaderBorder)
	tableHeader.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	table.SetTopLevelRows(provider.RowData(table))
	table.InstallCmdHandlers(constants.OpenEditorItemID, func(_ any) bool { return table.HasSelection() },
		func(_ any) { provider.OpenEditor(widget.FindRebuildable(table), table) })
	table.InstallCmdHandlers(constants.OpenOnePageReferenceItemID,
		func(_ any) bool { return CanOpenPageRef(table) },
		func(_ any) { OpenPageRef(table) })
	table.InstallCmdHandlers(constants.OpenEachPageReferenceItemID,
		func(_ any) bool { return CanOpenPageRef(table) },
		func(_ any) { OpenEachPageRef(table) })
	table.InstallCmdHandlers(unison.DeleteItemID,
		func(_ any) bool { return table.HasSelection() },
		func(_ any) { provider.DeleteSelection(table) })
	parent.AddChild(tableHeader)
	parent.AddChild(table)
	return table
}

func collectUUIDs(node gurps.Node, m map[uuid.UUID]bool) {
	m[node.UUID()] = true
	for _, child := range node.NodeChildren() {
		collectUUIDs(child, m)
	}
}

func deleteTableSelection[T comparableNode](table *unison.Table, topLevelRows []T, setTopLevelRows func(nodes []T), parentPtrFunc func(node T) *T, childrenPtrFunc func(node T) *[]T) {
	if sel := table.SelectedRows(true); len(sel) > 0 {
		ids := make(map[uuid.UUID]bool, len(sel))
		list := make([]T, 0, len(sel))
		for _, row := range sel {
			if target := ExtractFromRowData[T](row); !toolbox.IsNil(target) {
				list = append(list, target)
				collectUUIDs(target, ids)
			}
		}
		if !workspace.CloseUUID(ids) {
			return
		}
		needSet := false
		for _, target := range list {
			parentPtr := parentPtrFunc(target)
			if toolbox.IsNil(*parentPtr) {
				for i, one := range topLevelRows {
					if one == target {
						topLevelRows = slices.Delete(topLevelRows, i, i+1)
						needSet = true
						break
					}
				}
			} else {
				childrenPtr := childrenPtrFunc(*parentPtr)
				for i, one := range *childrenPtr {
					if one == target {
						*childrenPtr = slices.Delete(*childrenPtr, i, i+1)
						break
					}
				}
			}
		}
		if needSet {
			setTopLevelRows(topLevelRows)
		}
		if rebuilder := widget.FindRebuildable(table); rebuilder != nil {
			rebuilder.Rebuild(true)
		}
	}
}

// RecordTableSelection collects the currently selected row UUIDs.
func RecordTableSelection(table *unison.Table) map[uuid.UUID]bool {
	rows := table.SelectedRows(false)
	selection := make(map[uuid.UUID]bool, len(rows))
	for _, row := range rows {
		if node := ExtractFromRowData[gurps.Node](row); node != nil {
			selection[node.UUID()] = true
		}
	}
	return selection
}

// ApplyTableSelection locates the rows with the given UUIDs and selects them, replacing any existing selection.
func ApplyTableSelection(table *unison.Table, selection map[uuid.UUID]bool) {
	table.ClearSelection()
	if len(selection) != 0 {
		_, indexes := collectRowMappings(0, make([]int, 0, len(selection)), selection, table.TopLevelRows())
		if len(indexes) != 0 {
			table.SelectByIndex(indexes...)
		}
	}
}

func collectRowMappings(index int, indexes []int, selection map[uuid.UUID]bool, rows []unison.TableRowData) (updatedIndex int, updatedIndexes []int) {
	for _, row := range rows {
		if node := ExtractFromRowData[gurps.Node](row); node != nil {
			if selection[node.UUID()] {
				indexes = append(indexes, index)
			}
		}
		index++
		if row.IsOpen() {
			index, indexes = collectRowMappings(index, indexes, selection, row.ChildRows())
		}
	}
	return index, indexes
}
