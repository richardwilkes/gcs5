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

import "github.com/richardwilkes/unison"

// TableProvider defines the methods a table provider must contain.
type TableProvider interface {
	Headers() []unison.TableColumnHeader
	RowData(table *unison.Table) []unison.TableRowData
	SyncHeader(headers []unison.TableColumnHeader)
	HierarchyColumnIndex() int
	ExcessWidthColumnIndex() int
}
