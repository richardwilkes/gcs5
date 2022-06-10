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
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/widget"
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
type TableProvider interface {
	gurps.EntityProvider
	DragKey() string
	DragSVG() *unison.SVG
	DropShouldMoveData(drop *unison.TableDrop) bool
	ItemNames() (singular, plural string)
	Headers() []unison.TableColumnHeader
	RowData(table *unison.Table) []unison.TableRowData
	SyncHeader(headers []unison.TableColumnHeader)
	HierarchyColumnIndex() int
	ExcessWidthColumnIndex() int
	OpenEditor(owner widget.Rebuildable, table *unison.Table)
	CreateItem(owner widget.Rebuildable, table *unison.Table, variant ItemVariant)
	DeleteSelection(table *unison.Table)
}

// InstallTableDropSupport installs our standard drop support on a table.
func InstallTableDropSupport(table *unison.Table, provider TableProvider) {
	table.InstallDropSupport(provider.DragKey(), provider.DropShouldMoveData,
		func(drop *unison.TableDrop) {
			// TODO: copyCallback
		}, func(drop *unison.TableDrop, row, newParent unison.TableRowData) {
			// TODO: setRowParentCallback
		}, func(drop *unison.TableDrop, row unison.TableRowData, children []unison.TableRowData) {
			// TODO: setChildRowsCallback
		})
}
