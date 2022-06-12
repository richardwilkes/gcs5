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

package widget

import (
	"github.com/richardwilkes/gcs/constants"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
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
type TableProvider[T unison.TableRowConstraint[T]] interface {
	unison.TableModel[T]
	gurps.EntityProvider
	SetTable(table *unison.Table[T])
	DragKey() string
	DragSVG() *unison.SVG
	DropShouldMoveData(from, to *unison.Table[T]) bool
	ItemNames() (singular, plural string)
	Headers() []unison.TableColumnHeader[T]
	SyncHeader(headers []unison.TableColumnHeader[T])
	HierarchyColumnIndex() int
	ExcessWidthColumnIndex() int
	OpenEditor(owner Rebuildable, table *unison.Table[T])
	CreateItem(owner Rebuildable, table *unison.Table[T], variant ItemVariant)
	DeleteSelection(table *unison.Table[T])
}

// TableSetupColumnSizes sets the standard column sizing.
func TableSetupColumnSizes[T unison.TableRowConstraint[T]](table *unison.Table[T], headers []unison.TableColumnHeader[T]) {
	table.ColumnSizes = make([]unison.ColumnSize, len(headers))
	for i := range table.ColumnSizes {
		_, pref, _ := headers[i].AsPanel().Sizes(unison.Size{})
		pref.Width += table.Padding.Left + table.Padding.Right
		table.ColumnSizes[i].AutoMinimum = pref.Width
		table.ColumnSizes[i].AutoMaximum = 800
		table.ColumnSizes[i].Minimum = pref.Width
		table.ColumnSizes[i].Maximum = 10000
	}
}

// TableInstallStdCallbacks installs the standard callbacks.
func TableInstallStdCallbacks[T unison.TableRowConstraint[T]](table *unison.Table[T]) {
	mouseDownCallback := table.MouseDownCallback
	table.MouseDownCallback = func(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
		table.RequestFocus()
		return mouseDownCallback(where, button, clickCount, mod)
	}
	table.SelectionDoubleClickCallback = func() { table.PerformCmd(nil, constants.OpenEditorItemID) }
	table.KeyDownCallback = func(keyCode unison.KeyCode, mod unison.Modifiers, _ bool) bool {
		if mod == 0 && (keyCode == unison.KeyBackspace || keyCode == unison.KeyDelete) {
			table.PerformCmd(table, unison.DeleteItemID)
			return true
		}
		return false
	}
}

// TableCreateHeader creates the standard table header with a flexible sorting mechanism.
func TableCreateHeader[T unison.TableRowConstraint[T]](table *unison.Table[T], headers []unison.TableColumnHeader[T]) *unison.TableHeader[T] {
	tableHeader := unison.NewTableHeader(table, headers...)
	tableHeader.Less = func(s1, s2 string) bool {
		if n1, err := fxp.FromString(s1); err == nil {
			var n2 fxp.Int
			if n2, err = fxp.FromString(s2); err == nil {
				return n1 < n2
			}
		}
		return txt.NaturalLess(s1, s2, true)
	}
	return tableHeader
}

// InstallTableDropSupport installs our standard drop support on a table.
func InstallTableDropSupport[T unison.TableRowConstraint[T]](table *unison.Table[T], provider TableProvider[T]) {
	// TODO: Revisit this for undo
	unison.InstallDropSupport[T, any](table, provider.DragKey(), provider.DropShouldMoveData, nil, nil)
	table.DragRemovedRowsCallback = func() { MarkModified(table) }
	table.DropOccurredCallback = func() { MarkModified(table) }
}
