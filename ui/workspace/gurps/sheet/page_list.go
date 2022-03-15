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

package sheet

import (
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// PageList holds a list for a sheet page.
type PageList struct {
	unison.Panel
	tableHeader *unison.TableHeader
	table       *unison.Table
}

// NewPageList creates a new list for a sheet page.
func NewPageList(entity *gurps.Entity, columnHeaders []unison.TableColumnHeader, topLevelRows func(table *unison.Table) []unison.TableRowData) *PageList {
	p := &PageList{
		table: unison.NewTable(),
	}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{Columns: 1})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.StartAlignment,
		HGrab:  true,
	})
	p.table.MinimumRowHeight = theme.PageFieldPrimaryFont.LineHeight()
	p.table.Padding.Top = 0
	p.table.Padding.Bottom = 0
	p.table.ColumnSizes = make([]unison.ColumnSize, len(columnHeaders))
	for i := range p.table.ColumnSizes {
		_, pref, _ := columnHeaders[i].AsPanel().Sizes(geom32.Size{})
		p.table.ColumnSizes[i].AutoMinimum = pref.Width
		p.table.ColumnSizes[i].AutoMaximum = 800
		p.table.ColumnSizes[i].Minimum = pref.Width
		p.table.ColumnSizes[i].Maximum = 10000
	}
	p.table.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.tableHeader = unison.NewTableHeader(p.table, columnHeaders...)
	p.tableHeader.BackgroundInk = theme.HeaderColor
	p.tableHeader.Less = func(s1, s2 string) bool {
		if n1, err := fixed.F64d4FromString(s1); err == nil {
			var n2 fixed.F64d4
			if n2, err = fixed.F64d4FromString(s2); err == nil {
				return n1 < n2
			}
		}
		return txt.NaturalLess(s1, s2, true)
	}
	p.tableHeader.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.table.SetTopLevelRows(topLevelRows(p.table))
	p.table.SizeColumnsToFit(true)
	p.AddChild(p.tableHeader)
	p.AddChild(p.table)
	return p
}
