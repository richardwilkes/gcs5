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

package gurps

import (
	"fmt"
	"strings"

	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// Matcher defines the methods that rows that can participate in searching must implement.
type Matcher interface {
	// Match should look for the text in the object and return true if it is present. Note that calls to this method
	// should always pass in text that has already been run through strings.ToLower().
	Match(text string) bool
}

// ListFileDockable holds the view for a file that contains a (potentially hierarchical) list of data.
type ListFileDockable struct {
	unison.Panel
	path            string
	lockButton      *unison.Button
	hierarchyButton *unison.Button
	sizeToFitButton *unison.Button
	backButton      *unison.Button
	forwardButton   *unison.Button
	searchField     *unison.Field
	matchesLabel    *unison.Label
	scroll          *unison.ScrollPanel
	table           *unison.Table
	searchResult    []unison.TableRowData
	searchIndex     int
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
	d.lockButton.ClickCallback = d.toggleLock

	d.hierarchyButton = unison.NewSVGButton(res.HierarchySVG)
	d.hierarchyButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Opens/closes all hierarchical rows"))
	d.hierarchyButton.ClickCallback = d.toggleHierarchy

	d.sizeToFitButton = unison.NewSVGButton(res.SizeToFitSVG)
	d.sizeToFitButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Sets the width of each column to fit its contents"))
	d.sizeToFitButton.ClickCallback = d.sizeToFit

	d.backButton = unison.NewSVGButton(res.BackSVG)
	d.backButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Previous Match"))
	d.backButton.ClickCallback = d.previousMatch
	d.backButton.SetEnabled(false)

	d.forwardButton = unison.NewSVGButton(res.ForwardSVG)
	d.forwardButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Next Match"))
	d.forwardButton.ClickCallback = d.nextMatch
	d.forwardButton.SetEnabled(false)

	d.searchField = unison.NewField()
	search := i18n.Text("Search")
	d.searchField.Watermark = search
	d.searchField.Tooltip = unison.NewTooltipWithText(search)
	d.searchField.ModifiedCallback = d.searchModified
	d.searchField.KeyDownCallback = func(keyCode unison.KeyCode, mod unison.Modifiers, repeat bool) bool {
		if keyCode == unison.KeyReturn || keyCode == unison.KeyNumPadEnter {
			if mod.ShiftDown() {
				d.previousMatch()
			} else {
				d.nextMatch()
			}
			return true
		}
		return d.searchField.DefaultKeyDown(keyCode, mod, repeat)
	}
	d.searchField.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
		HGrab:  true,
	})

	d.matchesLabel = unison.NewLabel()
	d.matchesLabel.Text = "-"
	d.matchesLabel.Tooltip = unison.NewTooltipWithText(i18n.Text("Number of matches found"))

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
	toolbar.AddChild(d.backButton)
	toolbar.AddChild(d.forwardButton)
	toolbar.AddChild(d.searchField)
	toolbar.AddChild(d.matchesLabel)
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})

	d.AddChild(toolbar)
	d.AddChild(d.scroll)
	return d
}

// TitleIcon implements node.FileBackedDockable
func (d *ListFileDockable) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  library.FileInfoFor(d.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements node.FileBackedDockable
func (d *ListFileDockable) Title() string {
	return fs.BaseName(d.path)
}

// Tooltip implements node.FileBackedDockable
func (d *ListFileDockable) Tooltip() string {
	return d.path
}

// BackingFilePath implements node.FileBackedDockable
func (d *ListFileDockable) BackingFilePath() string {
	return d.path
}

// Modified implements node.FileBackedDockable
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

func (d *ListFileDockable) searchModified() {
	d.searchIndex = 0
	d.searchResult = nil
	text := strings.ToLower(d.searchField.Text())
	for _, row := range d.table.TopLevelRows() {
		d.search(text, row)
	}
	d.adjustForMatch()
}

func (d *ListFileDockable) search(text string, row unison.TableRowData) {
	if matcher, ok := row.(Matcher); ok {
		if matcher.Match(text) {
			d.searchResult = append(d.searchResult, row)
		}
	}
	if row.CanHaveChildRows() {
		for _, child := range row.ChildRows() {
			d.search(text, child)
		}
	}
}

func (d *ListFileDockable) previousMatch() {
	if d.searchIndex > 0 {
		d.searchIndex--
		d.adjustForMatch()
	}
}

func (d *ListFileDockable) nextMatch() {
	if d.searchIndex < len(d.searchResult)-1 {
		d.searchIndex++
		d.adjustForMatch()
	}
}

func (d *ListFileDockable) adjustForMatch() {
	d.backButton.SetEnabled(d.searchIndex != 0)
	d.forwardButton.SetEnabled(len(d.searchResult) != 0 && d.searchIndex != len(d.searchResult)-1)
	if len(d.searchResult) != 0 {
		d.matchesLabel.Text = fmt.Sprintf(i18n.Text("%d of %d"), d.searchIndex+1, len(d.searchResult))
		row := d.searchResult[d.searchIndex]
		d.table.DiscloseRow(row, false)
		d.table.ClearSelection()
		rowIndex := d.table.RowToIndex(row)
		d.table.SelectByIndex(rowIndex)
		d.table.ScrollRowIntoView(rowIndex)
	} else {
		d.matchesLabel.Text = "-"
	}
	d.matchesLabel.Parent().MarkForLayoutAndRedraw()
}
