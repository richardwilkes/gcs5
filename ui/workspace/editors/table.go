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
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/widget/ntable"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/unison"
	"golang.org/x/exp/slices"
)

func newTable[T gurps.NodeConstraint[T]](parent *unison.Panel, provider ntable.TableProvider[T]) *unison.Table[*ntable.Node[T]] {
	header, table := ntable.NewNodeTable[T](provider, unison.FieldFont)
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
	ntable.InstallTableDropSupport(table, provider)
	table.SyncToModel()
	parent.AddChild(header)
	parent.AddChild(table)
	return table
}

func duplicateTableSelection[T gurps.NodeConstraint[T]](table *unison.Table[*ntable.Node[T]], topLevelRows []T, setTopLevelRows func(nodes []T), childrenPtrFunc func(node T) *[]T) {
	if table.HasSelection() {
		var zero T
		needSet := false
		sel := table.SelectedRows(true)
		selMap := make(map[uuid.UUID]bool, len(sel))
		for _, row := range sel {
			if target := row.Data(); target != zero {
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

func deleteTableSelection[T gurps.NodeConstraint[T]](table *unison.Table[*ntable.Node[T]], topLevelRows []T, setTopLevelRows func(nodes []T), childrenPtrFunc func(node T) *[]T) {
	if table.HasSelection() {
		sel := table.SelectedRows(true)
		ids := make(map[uuid.UUID]bool, len(sel))
		list := make([]T, 0, len(sel))
		var zero T
		for _, row := range sel {
			unison.CollectUUIDsFromRow(row, ids)
			if target := row.Data(); target != zero {
				list = append(list, target)
			}
		}
		if !workspace.CloseUUID(ids) {
			return
		}
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
