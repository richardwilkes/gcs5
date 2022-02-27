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
	"github.com/richardwilkes/gcs/ui/workspace/node"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var _ unison.Dockable = &Navigator{}

// Pather defines the method for returning a path from an object.
type Pather interface {
	Path() string
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
func (n *Navigator) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
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
	for _, row := range n.table.SelectedRows() {
		n.openRow(row)
	}
}

func (n *Navigator) openRow(row unison.TableRowData) {
	switch t := row.(type) {
	case *LibraryNode, *DirectoryNode:
		for _, child := range t.ChildRows() {
			n.openRow(child)
		}
	case *FileNode:
		OpenFile(n.Window(), path.Join(t.library.Path(), t.path))
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

// OpenFile attempts to open the given file path in the given window.
func OpenFile(wnd *unison.Window, filePath string) (dockable unison.Dockable, wasOpen bool) {
	workspace := FromWindow(wnd)
	if workspace == nil {
		return nil, false
	}
	var defaultDockContainer *unison.DockContainer
	if focus := wnd.Focus(); focus != nil {
		if dc := unison.DockContainerFor(focus); dc != nil && dc.Dock == workspace.DocumentDock.Dock {
			defaultDockContainer = dc
		}
	}
	var d unison.Dockable
	filePath = path.Clean(filePath)
	workspace.DocumentDock.RootDockLayout().ForEachDockContainer(func(dc *unison.DockContainer) bool {
		for _, one := range dc.Dockables() {
			if f, ok := one.(node.FileBackedDockable); ok {
				if filePath == f.BackingFilePath() {
					d = one
					dc.SetCurrentDockable(one)
					dc.AcquireFocus()
					return true
				}
			}
			if defaultDockContainer == nil {
				defaultDockContainer = dc
			}
		}
		return false
	})
	if d != nil {
		return d, true
	}
	fi := library.FileInfoFor(filePath)
	if fi.IsSpecial {
		return nil, false
	}
	var err error
	if d, err = fi.Load(filePath); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to open file"), err)
		return nil, false
	}
	if defaultDockContainer != nil {
		defaultDockContainer.Stack(d, -1)
	} else {
		workspace.DocumentDock.DockTo(d, nil, unison.LeftSide)
		d.AsPanel().RequestFocus()
	}
	return d, false
}

func createNodeCell(ext, title string, selected bool) *unison.Panel {
	size := unison.LabelFont.Size() + 5
	fi := library.FileInfoFor(ext)
	label := unison.NewLabel()
	label.Text = title
	label.Drawable = &unison.DrawableSVG{
		SVG:  fi.SVG,
		Size: geom32.NewSize(size, size),
	}
	if selected {
		label.LabelTheme.OnBackgroundInk = unison.OnSelectionColor
	}
	return label.AsPanel()
}
