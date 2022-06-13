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
	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/constants"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
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
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
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

const tableProviderClientKey = "table-provider"

// InstallTableDropSupport installs our standard drop support on a table.
func InstallTableDropSupport[T unison.TableRowConstraint[T]](table *unison.Table[T], provider TableProvider[T]) {
	table.ClientData()[tableProviderClientKey] = provider
	unison.InstallDropSupport[T, *tableDragUndoEditData[T]](table, provider.DragKey(), provider.DropShouldMoveData,
		willDropCallback[T], didDropCallback[T])
	table.DragRemovedRowsCallback = func() { MarkModified(table) }
	table.DropOccurredCallback = func() { MarkModified(table) }
}

func willDropCallback[T unison.TableRowConstraint[T]](from, to *unison.Table[T], move bool) *unison.UndoEdit[*tableDragUndoEditData[T]] {
	mgr := unison.UndoManagerFor(from)
	if mgr == nil {
		return nil
	}
	data := newTableDragUndoEditData(from, to, move)
	if data == nil {
		return nil
	}
	return &unison.UndoEdit[*tableDragUndoEditData[T]]{
		ID:         unison.NextUndoID(),
		EditName:   i18n.Text("Drag"),
		UndoFunc:   func(e *unison.UndoEdit[*tableDragUndoEditData[T]]) { e.BeforeData.apply() },
		RedoFunc:   func(e *unison.UndoEdit[*tableDragUndoEditData[T]]) { e.AfterData.apply() },
		AbsorbFunc: func(e *unison.UndoEdit[*tableDragUndoEditData[T]], other unison.Undoable) bool { return false },
		BeforeData: data,
	}
}

func didDropCallback[T unison.TableRowConstraint[T]](undo *unison.UndoEdit[*tableDragUndoEditData[T]], from, to *unison.Table[T], move bool) {
	if undo == nil {
		return
	}
	mgr := unison.UndoManagerFor(from)
	if mgr == nil {
		return
	}
	undo.AfterData = newTableDragUndoEditData(from, to, move)
	if undo.AfterData != nil {
		mgr.Add(undo)
	}
}

type tableDragUndoEditData[T unison.TableRowConstraint[T]] struct {
	From       *unison.Table[T]
	To         *unison.Table[T]
	FromData   []byte
	FromSelMap map[uuid.UUID]bool
	ToData     []byte
	ToSelMap   map[uuid.UUID]bool
	Move       bool
}

func newTableDragUndoEditData[T unison.TableRowConstraint[T]](from, to *unison.Table[T], move bool) *tableDragUndoEditData[T] {
	data, err := collectTableData(to)
	if err != nil {
		jot.Error(err)
		return nil
	}
	undo := &tableDragUndoEditData[T]{
		From:     from,
		To:       to,
		ToData:   data,
		ToSelMap: to.CopySelectionMap(),
		Move:     move,
	}
	if move && from != to {
		if data, err = collectTableData(from); err != nil {
			jot.Error(err)
			return nil
		}
		undo.FromData = data
		undo.FromSelMap = from.CopySelectionMap()
	}
	return undo
}

func (t *tableDragUndoEditData[T]) apply() {
	applyTableData(t.To, t.ToData, t.ToSelMap)
	if t.Move && t.From != t.To {
		applyTableData(t.From, t.FromData, t.FromSelMap)
	}
}

func collectTableData[T unison.TableRowConstraint[T]](table *unison.Table[T]) ([]byte, error) {
	provider, ok := table.ClientData()[tableProviderClientKey].(TableProvider[T])
	if !ok {
		return nil, errs.New("unable to locate provider")
	}
	return provider.Serialize()
}

func applyTableData[T unison.TableRowConstraint[T]](table *unison.Table[T], data []byte, selMap map[uuid.UUID]bool) {
	provider, ok := table.ClientData()[tableProviderClientKey].(TableProvider[T])
	if !ok {
		jot.Error(errs.New("unable to locate provider"))
		return
	}
	if err := provider.Deserialize(data); err != nil {
		jot.Error(err)
		return
	}
	table.SyncToModel()
	MarkModified(table)
	table.SetSelectionMap(selMap)
}
