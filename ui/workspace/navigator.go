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
	"strings"

	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var _ unison.Dockable = &Navigator{}

// FileBackedDockable defines methods a Dockable that is based on a file should implement.
type FileBackedDockable interface {
	unison.Dockable
	BackingFilePath() string
}

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

	n.scroll.MouseWheelMultiplier = 2
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
func OpenFile(wnd *unison.Window, filePath string) {
	workspace := FromWindow(wnd)
	if workspace == nil {
		return
	}
	var defaultDockContainer *unison.DockContainer
	if focus := wnd.Focus(); focus != nil {
		if dc := unison.DockContainerFor(focus); dc != nil && dc.Dock == workspace.DocumentDock.Dock {
			defaultDockContainer = dc
		}
	}
	found := false
	filePath = path.Clean(filePath)
	workspace.DocumentDock.RootDockLayout().ForEachDockContainer(func(dc *unison.DockContainer) bool {
		for _, d := range dc.Dockables() {
			if f, ok := d.(FileBackedDockable); ok {
				if filePath == f.BackingFilePath() {
					found = true
					dc.SetCurrentDockable(d)
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
	if !found {
		var err error
		var d unison.Dockable
		if unison.EncodedImageFormatForPath(filePath).CanRead() {
			if d, err = NewImageDockable(filePath); err != nil {
				unison.ErrorDialogWithMessage(i18n.Text("Unable to open image file"), err.Error())
				return
			}
		} else {
			switch strings.ToLower(path.Ext(filePath)) {
			case ".pdf":
				if d, err = NewPDFDockable(filePath); err != nil {
					unison.ErrorDialogWithMessage(i18n.Text("Unable to open PDF"), err.Error())
					return
				}
			case ".skl":
				if d, err = NewSkillListDockable(filePath); err != nil {
					unison.ErrorDialogWithMessage(i18n.Text("Unable to open skills list"), err.Error())
					return
				}
			case ".spl":
				if d, err = NewSpellListDockable(filePath); err != nil {
					unison.ErrorDialogWithMessage(i18n.Text("Unable to open spells list"), err.Error())
					return
				}
			default:
				unison.ErrorDialogWithMessage(i18n.Text("Unable to open file"), filePath)
				return
			}
		}
		if defaultDockContainer != nil {
			defaultDockContainer.Stack(d, -1)
		} else {
			workspace.DocumentDock.DockTo(d, nil, unison.LeftSide)
			d.AsPanel().RequestFocus()
		}
	}
}

func createNodeCell(ext, title string, selected bool) *unison.Panel {
	size := unison.LabelFont.Size() + 5
	info, ok := library.FileTypes[ext]
	if !ok {
		info = library.FileTypes[library.GenericFile]
	}
	label := unison.NewLabel()
	label.Text = title
	label.Drawable = &unison.DrawableSVG{
		SVG:  info.SVG,
		Size: geom32.NewSize(size, size),
	}
	if selected {
		label.LabelTheme.OnBackgroundInk = unison.OnSelectionColor
	}
	return label.AsPanel()
}
