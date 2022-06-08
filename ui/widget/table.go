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
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/unison"
)

// TableDragData holds the data from a table row drag.
type TableDragData struct {
	Table *unison.Table
	Rows  []unison.TableRowData
}

// InstallTableDragSupport installs drag support into a table.
func InstallTableDragSupport(table *unison.Table, svg *unison.SVG, dragKey, singularName, pluralName string) {
	orig := table.MouseDragCallback
	table.MouseDragCallback = func(where unison.Point, button int, mod unison.Modifiers) bool {
		if orig != nil && orig(where, button, mod) {
			return true
		}
		if table.HasSelection() && table.IsDragGesture(where) {
			data := &TableDragData{
				Table: table,
				Rows:  table.SelectedRows(true),
			}
			drawable := NewTableDragDrawable(data, svg, singularName, pluralName)
			size := drawable.LogicalSize()
			table.StartDataDrag(&unison.DragData{
				Data:     map[string]any{dragKey: data},
				Drawable: drawable,
				Ink:      table.OnBackgroundInk,
				Offset:   unison.Point{X: 0, Y: -size.Height / 2},
			})
		}
		return false
	}
}

// TableSetupColumnSizes sets the standard column sizing.
func TableSetupColumnSizes(table *unison.Table, headers []unison.TableColumnHeader) {
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
func TableInstallStdCallbacks(table *unison.Table) {
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
func TableCreateHeader(table *unison.Table, headers []unison.TableColumnHeader) *unison.TableHeader {
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
