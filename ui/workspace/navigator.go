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
	"path"
	"path/filepath"

	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

var _ unison.Dockable = &Navigator{}

// Pather defines the method for returning a path from an object.
type Pather interface {
	Path() string
}

// FileBackedDockable defines methods a Dockable that is based on a file should implement.
type FileBackedDockable interface {
	unison.Dockable
	BackingFilePath() string
}

// Navigator holds the workspace navigation panel.
type Navigator struct {
	unison.Panel
	scroll *unison.ScrollPanel
	table  *unison.Table
}

// RegisterFileTypes registers special navigator file types.
func RegisterFileTypes() {
	registerSpecialFileInfo(library.ClosedFolder, res.ClosedFolderSVG)
	registerSpecialFileInfo(library.OpenFolder, res.OpenFolderSVG)
	registerSpecialFileInfo(library.GenericFile, res.GenericFileSVG)
}

func registerSpecialFileInfo(key string, svg *unison.SVG) {
	library.FileInfo{
		Extension: key,
		SVG:       svg,
		IsSpecial: true,
	}.Register()
}

func newNavigator() *Navigator {
	n := &Navigator{
		scroll: unison.NewScrollPanel(),
		table:  unison.NewTable(),
	}
	n.Self = n

	n.table.ColumnSizes = make([]unison.ColumnSize, 1)
	globalSettings := settings.Global()
	libs := globalSettings.LibrarySet.List()
	rows := make([]unison.TableRowData, 0, len(libs))
	for _, one := range libs {
		rows = append(rows, NewLibraryNode(n, one))
	}
	n.table.SetTopLevelRows(rows)
	n.ApplyDisclosedPaths(globalSettings.LibraryExplorer.OpenRowKeys)
	n.table.SizeColumnsToFit(true)

	n.scroll.SetContent(n.table, unison.FillBehavior)
	n.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})

	n.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
	})
	n.AddChild(n.scroll)

	n.table.SelectionDoubleClickCallback = n.handleSelectionDoubleClick
	return n
}

func (n *Navigator) adjustTableSize() {
	n.table.SyncToModel()
	n.table.SizeColumnsToFit(true)
}

// TitleIcon implements unison.Dockable
func (n *Navigator) TitleIcon(suggestedSize unison.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  unison.DocumentSVG(),
		Size: suggestedSize,
	}
}

// Title implements unison.Dockable
func (n *Navigator) Title() string {
	return i18n.Text("Library Explorer")
}

// Tooltip implements unison.Dockable
func (n *Navigator) Tooltip() string {
	return ""
}

// Modified implements unison.Dockable
func (n *Navigator) Modified() bool {
	return false
}

func (n *Navigator) handleSelectionDoubleClick() {
	for _, row := range n.table.SelectedRows(false) {
		n.openRow(row)
	}
}

func (n *Navigator) openRow(row unison.TableRowData) {
	if node, ok := row.(*FileNode); ok {
		OpenFile(n.Window(), path.Join(node.library.Path(), node.path))
	}
}

// DisclosedPaths returns a list of paths that are currently disclosed.
func (n *Navigator) DisclosedPaths() []string {
	return n.accumulateDisclosedPaths(n.table.TopLevelRows(), nil)
}

func (n *Navigator) accumulateDisclosedPaths(rows []unison.TableRowData, disclosedPaths []string) []string {
	for _, row := range rows {
		if row.IsOpen() {
			if p, ok := row.(Pather); ok {
				disclosedPaths = append(disclosedPaths, p.Path())
			}
		}
		disclosedPaths = n.accumulateDisclosedPaths(row.ChildRows(), disclosedPaths)
	}
	return disclosedPaths
}

// ApplyDisclosedPaths closes all nodes except the ones provided, which are explicitly opened.
func (n *Navigator) ApplyDisclosedPaths(paths []string) {
	m := make(map[string]bool, len(paths))
	for _, one := range paths {
		m[one] = true
	}
	n.applyDisclosedPaths(n.table.TopLevelRows(), m)
}

func (n *Navigator) applyDisclosedPaths(rows []unison.TableRowData, paths map[string]bool) {
	for _, row := range rows {
		if p, ok := row.(Pather); ok {
			open := paths[p.Path()]
			if row.IsOpen() != open {
				row.SetOpen(open)
			}
		}
		n.applyDisclosedPaths(row.ChildRows(), paths)
	}
}

// OpenFiles attempts to open the given file paths.
func OpenFiles(filePaths []string) {
	for _, wnd := range unison.Windows() {
		if ws := FromWindow(wnd); ws != nil {
			for _, one := range filePaths {
				if p, err := filepath.Abs(one); err != nil {
					unison.ErrorDialogWithError(i18n.Text("Unable to open ")+one, err)
				} else {
					OpenFile(wnd, p)
				}
			}
		}
	}
}

// DisplayNewDockable adds the Dockable to the dock and gives it the focus.
func DisplayNewDockable(wnd *unison.Window, dockable unison.Dockable) {
	ws := FromWindowOrAny(wnd)
	if ws == nil {
		ShowUnableToLocateWorkspaceError()
		return
	}
	defer func() { dockable.AsPanel().RequestFocus() }()
	if fbd, ok := dockable.(FileBackedDockable); ok {
		fi := library.FileInfoFor(fbd.BackingFilePath())
		if dc := ws.CurrentlyFocusedDockContainer(); dc != nil && DockContainerHoldsExtension(dc, fi.ExtensionsToGroupWith...) {
			dc.Stack(dockable, -1)
			return
		} else if dc = ws.LocateDockContainerForExtension(fi.ExtensionsToGroupWith...); dc != nil {
			dc.Stack(dockable, -1)
			return
		}
	}
	ws.DocumentDock.DockTo(dockable, nil, unison.RightSide)
}

// OpenFile attempts to open the given file path in the given window, which should contain a workspace. May pass nil for
// wnd to let it pick the first such window it discovers.
func OpenFile(wnd *unison.Window, filePath string) (dockable unison.Dockable, wasOpen bool) {
	ws := FromWindowOrAny(wnd)
	if ws == nil {
		ShowUnableToLocateWorkspaceError()
		return nil, false
	}
	var err error
	if filePath, err = filepath.Abs(filePath); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to resolve path"), err)
		return nil, false
	}
	if d := ws.LocateFileBackedDockable(filePath); d != nil {
		dc := unison.DockContainerFor(d)
		dc.SetCurrentDockable(d)
		dc.AcquireFocus()
		return d, true
	}
	fi := library.FileInfoFor(filePath)
	if fi.IsSpecial {
		return nil, false
	}
	var d unison.Dockable
	if d, err = fi.Load(filePath); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to open file"), err)
		return nil, false
	}
	DisplayNewDockable(wnd, d)
	return d, false
}

func createNodeCell(ext, title string, selected bool) *unison.Panel {
	size := unison.LabelFont.Size() + 5
	fi := library.FileInfoFor(ext)
	label := unison.NewLabel()
	label.Text = title
	label.Drawable = &unison.DrawableSVG{
		SVG:  fi.SVG,
		Size: unison.NewSize(size, size),
	}
	if selected {
		label.LabelTheme.OnBackgroundInk = unison.OnSelectionColor
	}
	return label.AsPanel()
}
