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

func newTable[T gurps.NodeConstraint[T]](parent *unison.Panel, provider widget.TableProvider[*Node[T]]) *unison.Table[*Node[T]] {
	table := unison.NewTable[*Node[T]](provider)
	provider.SetTable(table)
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
	table.SyncToModel()
	table.InstallCmdHandlers(constants.OpenEditorItemID, func(_ any) bool { return table.HasSelection() },
		func(_ any) { provider.OpenEditor(unison.AncestorOrSelf[widget.Rebuildable](table), table) })
	table.InstallCmdHandlers(constants.OpenOnePageReferenceItemID,
		func(_ any) bool { return CanOpenPageRef(table) },
		func(_ any) { OpenPageRef(table) })
	table.InstallCmdHandlers(constants.OpenEachPageReferenceItemID,
		func(_ any) bool { return CanOpenPageRef(table) },
		func(_ any) { OpenEachPageRef(table) })
	table.InstallCmdHandlers(unison.DeleteItemID,
		func(_ any) bool { return table.HasSelection() },
		func(_ any) { provider.DeleteSelection(table) })
	table.InstallCmdHandlers(constants.DuplicateItemID,
		func(_ any) bool { return table.HasSelection() },
		func(_ any) { provider.DuplicateSelection(table) })
	parent.AddChild(tableHeader)
	parent.AddChild(table)
	singular, plural := provider.ItemNames()
	table.InstallDragSupport(provider.DragSVG(), provider.DragKey(), singular, plural)
	widget.InstallTableDropSupport(table, provider)
	return table
}

func collectUUIDs[T gurps.NodeConstraint[T]](node T, m map[uuid.UUID]bool) {
	m[node.UUID()] = true
	for _, child := range node.NodeChildren() {
		collectUUIDs(child, m)
	}
}

func duplicateTableSelection[T gurps.NodeConstraint[T]](table *unison.Table[*Node[T]], topLevelRows []T, setTopLevelRows func(nodes []T), childrenPtrFunc func(node T) *[]T) {
	if table.HasSelection() {
		var zero T
		needSet := false
		sel := table.SelectedRows(true)
		selMap := make(map[uuid.UUID]bool, len(sel))
		for _, row := range sel {
			if target := ExtractFromRowData[T](row); !toolbox.IsNil(target) {
				parent := target.Parent()
				clone := target.Clone(target.OwningEntity(), parent, false)
				selMap[clone.UUID()] = true
				if parent == zero {
					for i, child := range topLevelRows {
						if child == target {
							topLevelRows = slices.Insert(topLevelRows, i+1, clone)
							needSet = true
							break
						}
					}
				} else {
					childrenPtr := childrenPtrFunc(parent)
					for i, child := range *childrenPtr {
						if child == target {
							*childrenPtr = slices.Insert(*childrenPtr, i+1, clone)
							break
						}
					}
				}
			}
		}
		if needSet {
			setTopLevelRows(topLevelRows)
		}
		table.SyncToModel()
		table.SetSelectionMap(selMap)
		if rebuilder := unison.AncestorOrSelf[widget.Rebuildable](table); rebuilder != nil {
			rebuilder.Rebuild(true)
		}
	}
}

func deleteTableSelection[T gurps.NodeConstraint[T]](table *unison.Table[*Node[T]], topLevelRows []T, setTopLevelRows func(nodes []T), childrenPtrFunc func(node T) *[]T) {
	if table.HasSelection() {
		sel := table.SelectedRows(true)
		ids := make(map[uuid.UUID]bool, len(sel))
		list := make([]T, 0, len(sel))
		for _, row := range sel {
			if target := ExtractFromRowData[T](row); !toolbox.IsNil(target) {
				list = append(list, target)
				collectUUIDs[T](target, ids)
			}
		}
		if !workspace.CloseUUID(ids) {
			return
		}
		var zero T
		needSet := false
		for _, target := range list {
			parent := target.Parent()
			if parent == zero {
				for i, one := range topLevelRows {
					if one == target {
						topLevelRows = slices.Delete(topLevelRows, i, i+1)
						needSet = true
						break
					}
				}
			} else {
				childrenPtr := childrenPtrFunc(parent)
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
		if rebuilder := unison.AncestorOrSelf[widget.Rebuildable](table); rebuilder != nil {
			rebuilder.Rebuild(true)
		}
	}
}

// RecordTableSelection collects the currently selected row UUIDs.
func RecordTableSelection[T gurps.NodeConstraint[T]](table *unison.Table[*Node[T]]) map[uuid.UUID]bool {
	var zero T
	rows := table.SelectedRows(false)
	selection := make(map[uuid.UUID]bool, len(rows))
	for _, row := range rows {
		if node := ExtractFromRowData[T](row); node != zero {
			selection[node.UUID()] = true
		}
	}
	return selection
}

// ApplyTableSelection locates the rows with the given UUIDs and selects them, replacing any existing selection.
func ApplyTableSelection[T gurps.NodeConstraint[T]](table *unison.Table[*Node[T]], selection map[uuid.UUID]bool) {
	table.ClearSelection()
	if len(selection) != 0 {
		_, indexes := collectRowMappings(0, make([]int, 0, len(selection)), selection, table.RootRows())
		if len(indexes) != 0 {
			table.SelectByIndex(indexes...)
		}
	}
}

func collectRowMappings[T gurps.NodeConstraint[T]](index int, indexes []int, selection map[uuid.UUID]bool, rows []*Node[T]) (updatedIndex int, updatedIndexes []int) {
	var zero T
	for _, row := range rows {
		if node := ExtractFromRowData[T](row); node != zero {
			if selection[node.UUID()] {
				indexes = append(indexes, index)
			}
		}
		index++
		if row.IsOpen() {
			index, indexes = collectRowMappings(index, indexes, selection, row.Children())
		}
	}
	return index, indexes
}
