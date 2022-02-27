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
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// CategoryProvider defines the methods objects that can provide categories must implement.
type CategoryProvider interface {
	Categories() []string
}

// ListFileDockable holds the view for a file that contains a (potentially hierarchical) list of data.
type ListFileDockable struct {
	unison.Panel
	path            string
	lockButton      *unison.Button
	hierarchyButton *unison.Button
	sizeToFitButton *unison.Button
	categoryPopup   *unison.PopupMenu
	scroll          *unison.ScrollPanel
	table           *unison.Table
	locked          bool
}

// NewListFileDockable creates a new ListFileDockable for list data files.
func NewListFileDockable(filePath string, columnHeaders []unison.TableColumnHeader, topLevelRows func(table *unison.Table) []unison.TableRowData) *ListFileDockable {
	d := &ListFileDockable{
		path:   filePath,
		scroll: unison.NewScrollPanel(),
		table:  unison.NewTable(),
	}
	d.Self = d
	d.SetLayout(&unison.FlexLayout{Columns: 1})

	d.table.ColumnSizes = make([]unison.ColumnSize, len(columnHeaders))
	for i := range d.table.ColumnSizes {
		_, pref, _ := columnHeaders[i].AsPanel().Sizes(geom32.Size{})
		d.table.ColumnSizes[i].AutoMinimum = pref.Width
		d.table.ColumnSizes[i].AutoMaximum = 800
		d.table.ColumnSizes[i].Minimum = pref.Width
		d.table.ColumnSizes[i].Maximum = 10000
	}
	d.table.SetTopLevelRows(topLevelRows(d.table))
	d.table.SizeColumnsToFit(true)

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

	d.lockButton = unison.NewSVGButton(res.LockSVG)
	d.toggleLock()
	d.lockButton.ClickCallback = func() { d.toggleLock() }

	d.hierarchyButton = unison.NewSVGButton(res.HierarchySVG)
	d.hierarchyButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Opens/closes all hierarchical rows"))
	d.hierarchyButton.ClickCallback = func() { d.toggleHierarchy() }

	d.sizeToFitButton = unison.NewSVGButton(res.SizeToFitSVG)
	d.sizeToFitButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Sets the width of each column to fit its contents"))
	d.sizeToFitButton.ClickCallback = func() { d.sizeToFit() }

	d.categoryPopup = unison.NewPopupMenu()
	d.categoryPopup.AddItem(i18n.Text("Any Category"))
	d.categoryPopup.AddSeparator()
	for _, one := range d.categoryList() {
		d.categoryPopup.AddItem(one)
	}

	toolbar := unison.NewPanel()
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, geom32.Insets{Bottom: 1}, false),
		unison.NewEmptyBorder(geom32.Insets{
			Top:    unison.StdVSpacing,
			Left:   unison.StdHSpacing,
			Bottom: unison.StdVSpacing,
			Right:  unison.StdHSpacing,
		})))
	toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	toolbar.AddChild(d.lockButton)
	toolbar.AddChild(d.hierarchyButton)
	toolbar.AddChild(d.sizeToFitButton)
	toolbar.AddChild(d.categoryPopup)
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})

	d.AddChild(toolbar)
	d.AddChild(d.scroll)
	return d
}

func (d *ListFileDockable) categoryList() []string {
	m := make(map[string]bool)
	for _, row := range d.table.TopLevelRows() {
		extractCategories(row, m)
	}
	list := make([]string, 0, len(m))
	for one := range m {
		list = append(list, one)
	}
	txt.SortStringsNaturalAscending(list)
	return list
}

func extractCategories(row unison.TableRowData, categories map[string]bool) {
	if provider, ok := row.(CategoryProvider); ok {
		for _, one := range provider.Categories() {
			categories[one] = true
		}
	}
	if row.CanHaveChildRows() {
		for _, child := range row.ChildRows() {
			extractCategories(child, categories)
		}
	}
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

func (d *ListFileDockable) toggleLock() {
	d.locked = !d.locked
	if dsvg, ok := d.lockButton.Drawable.(*unison.DrawableSVG); ok {
		if d.locked {
			dsvg.SVG = res.LockSVG
			d.lockButton.Tooltip = unison.NewTooltipWithSecondaryText(i18n.Text("Locked"), i18n.Text("Click to enable editing"))
		} else {
			dsvg.SVG = res.UnlockedSVG
			d.lockButton.Tooltip = unison.NewTooltipWithSecondaryText(i18n.Text("Unlocked"), i18n.Text("Click to disable editing"))
		}
	}
	d.lockButton.MarkForRedraw()
}

func (d *ListFileDockable) toggleHierarchy() {
	first := true
	open := false
	for _, row := range d.table.TopLevelRows() {
		if row.CanHaveChildRows() {
			if first {
				first = false
				open = !row.IsOpen()
			}
			setRowOpen(row, open)
		}
	}
	d.table.SyncToModel()
	d.table.MarkForRedraw()
}

func setRowOpen(row unison.TableRowData, open bool) {
	row.SetOpen(open)
	for _, child := range row.ChildRows() {
		if child.CanHaveChildRows() {
			setRowOpen(child, open)
		}
	}
}

func (d *ListFileDockable) sizeToFit() {
	d.table.SizeColumnsToFit(true)
	d.table.MarkForRedraw()
}
