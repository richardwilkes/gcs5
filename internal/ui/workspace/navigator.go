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

	"github.com/richardwilkes/gcs/internal/library"
	"github.com/richardwilkes/gcs/internal/settings"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var _ unison.Dockable = &Navigator{}

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
	rows := make([]unison.TableRowData, 0, len(settings.Global.Libraries))
	for _, one := range settings.Global.Libraries {
		rows = append(rows, NewLibraryNode(n, one))
	}
	n.table.SetTopLevelRows(rows)
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

func (n *Navigator) TitleIcon(suggestedSize geom32.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  unison.DocumentSVG(),
		Size: suggestedSize,
	}
}

func (n *Navigator) Title() string {
	return i18n.Text("Library Explorer")
}

func (n *Navigator) Tooltip() string {
	return ""
}

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
		wnd := n.Window()
		if workspace := FromWindow(wnd); workspace != nil {
			var defaultDockContainer *unison.DockContainer
			if focus := wnd.Focus(); focus != nil {
				if dc := unison.DockContainerFor(focus); dc != nil && dc.Dock == workspace.DocumentDock.Dock {
					defaultDockContainer = dc
				}
			}
			found := false
			filePath := path.Clean(path.Join(t.library.Config().Path, t.path))
			workspace.DocumentDock.RootDockLayout().ForEachDockContainer(func(dc *unison.DockContainer) bool {
				for _, d := range dc.Dockables() {
					switch actual := d.(type) {
					case *placeholder:
						if filePath == actual.path {
							found = true
							dc.SetCurrentDockable(d)
							dc.AcquireFocus()
							return true
						}
					case *ImageDockable:
						if filePath == actual.path {
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
				var d unison.Dockable
				if unison.EncodedImageFormatForPath(filePath).CanRead() {
					var err error
					if d, err = NewImageDockable(filePath); err != nil {
						unison.ErrorDialogWithMessage(i18n.Text("Unable to open image file"), err.Error())
						return
					}
				} else {
					d = NewPlaceholder(filePath)
				}
				if defaultDockContainer != nil {
					defaultDockContainer.Stack(d, -1)
				} else {
					workspace.DocumentDock.DockTo(d, nil, unison.LeftSide)
					d.AsPanel().RequestFocus()
				}
			}
		}
	}
}

func createNodeCell(ext, title string) *unison.Panel {
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
	return label.AsPanel()
}
