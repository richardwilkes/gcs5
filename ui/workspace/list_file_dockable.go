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

package workspace

import (
	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// ListFileDockable holds the view for a file that contains a (potentially hierarchical) list of data.
type ListFileDockable struct {
	unison.Panel
	path   string
	scroll *unison.ScrollPanel
	table  *unison.Table
}

// NewListFileDockable creates a new ListFileDockable for list data files.
func NewListFileDockable(filePath string, columnHeaders []unison.TableColumnHeader, topLevelRows func(table *unison.Table) []unison.TableRowData) *ListFileDockable {
	d := &ListFileDockable{
		path:   filePath,
		scroll: unison.NewScrollPanel(),
		table:  unison.NewTable(),
	}
	d.Self = d

	d.table.ColumnSizes = make([]unison.ColumnSize, len(columnHeaders))
	for i := range d.table.ColumnSizes {
		d.table.ColumnSizes[i].AutoMaximum = 800
	}
	d.table.SetTopLevelRows(topLevelRows(d.table))
	d.table.SizeColumnsToFit(true)

	d.scroll.MouseWheelMultiplier = 4
	header := unison.NewTableHeader(d.table, columnHeaders...)
	header.Less = func(s1, s2 string) bool {
		if n1, err := fixed.F64d4FromString(s1); err == nil {
			var n2 fixed.F64d4
			if n2, err = fixed.F64d4FromString(s2); err == nil {
				return n1 < n2
			}
		}
		return txt.NaturalLess(s1, s2, true)
	}
	d.scroll.SetColumnHeader(header)
	d.scroll.SetContent(d.table, unison.FillBehavior)
	d.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})

	d.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
	})
	d.AddChild(d.scroll)
	return d
}

// TitleIcon implements FileBackedDockable
func (d *ListFileDockable) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements FileBackedDockable
func (d *ListFileDockable) Title() string {
	return fs.BaseName(d.path)
}

// Tooltip implements FileBackedDockable
func (d *ListFileDockable) Tooltip() string {
	return d.path
}

// BackingFilePath implements FileBackedDockable
func (d *ListFileDockable) BackingFilePath() string {
	return d.path
}

// Modified implements FileBackedDockable
func (d *ListFileDockable) Modified() bool {
	return false
}

// MayAttemptClose implements unison.TabCloser
func (d *ListFileDockable) MayAttemptClose() bool {
	return true
}

// AttemptClose implements unison.TabCloser
func (d *ListFileDockable) AttemptClose() {
	if dc := unison.DockContainerFor(d); dc != nil {
		dc.Close(d)
	}
}
