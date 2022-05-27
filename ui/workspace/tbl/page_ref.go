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
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/gcs/ui/workspace/settings"
	"github.com/richardwilkes/unison"
)

// NewCanOpenPageRefFunc creates a new function for handling the action for opening a page reference.
func NewCanOpenPageRefFunc(table *unison.Table) func() bool {
	return func() bool {
		for _, row := range table.SelectedRows(false) {
			if n, ok := row.(*Node); ok {
				var data node.CellData
				n.Data().CellData(node.PageRefCellAlias, &data)
				if len(settings.ExtractPageReferences(data.Primary)) != 0 {
					return true
				}
			}
		}
		return false
	}
}

// NewOpenPageRefFunc creates a new function for handling the action for opening a page reference.
func NewOpenPageRefFunc(table *unison.Table) func() {
	return func() {
		promptCtx := make(map[string]bool)
		for _, row := range table.SelectedRows(false) {
			if n, ok := row.(*Node); ok {
				var data node.CellData
				n.Data().CellData(node.PageRefCellAlias, &data)
				for _, one := range settings.ExtractPageReferences(data.Primary) {
					if settings.OpenPageReference(table.Window(), one, data.Secondary, promptCtx) {
						return
					}
				}
			}
		}
	}
}

// NewOpenEachPageRefFunc creates a new function for handling the action for opening each page reference.
func NewOpenEachPageRefFunc(table *unison.Table) func() {
	return func() {
		promptCtx := make(map[string]bool)
		for _, row := range table.SelectedRows(false) {
			if n, ok := row.(*Node); ok {
				var data node.CellData
				n.Data().CellData(node.PageRefCellAlias, &data)
				for _, one := range settings.ExtractPageReferences(data.Primary) {
					if settings.OpenPageReference(table.Window(), one, data.Secondary, promptCtx) {
						return
					}
				}
			}
		}
	}
}
