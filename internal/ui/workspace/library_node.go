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
	"github.com/richardwilkes/gcs/internal/library"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var _ unison.TableRowData = &LibraryNode{}

// LibraryNode holds a library in the navigator.
type LibraryNode struct {
	nav      *Navigator
	library  *library.Library
	children []unison.TableRowData
	open     bool
}

// NewLibraryNode creates a new LibraryNode.
func NewLibraryNode(nav *Navigator, lib *library.Library) *LibraryNode {
	n := &LibraryNode{
		nav:     nav,
		library: lib,
	}
	n.Refresh()
	return n
}

// Refresh the contents of this node.
func (n *LibraryNode) Refresh() {
	n.children = refreshChildren(n.nav, n.library, ".")
}

// CanHaveChildRows always returns true.
func (n *LibraryNode) CanHaveChildRows() bool {
	return true
}

// ChildRows returns the children of this node.
func (n *LibraryNode) ChildRows() []unison.TableRowData {
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *LibraryNode) CellDataForSort(index int) string {
	switch index {
	case 0:
		return n.library.Title()
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *LibraryNode) ColumnCell(index int, _ bool) unison.Paneler {
	switch index {
	case 0:
		if n.open {
			return createNodeCell(library.OpenFolder, n.library.Title())
		}
		return createNodeCell(library.ClosedFolder, n.library.Title())
	default:
		jot.Errorf("column index out of range (0-0): %d", index)
		return unison.NewLabel()
	}
}

// IsOpen returns true if this node should display its children.
func (n *LibraryNode) IsOpen() bool {
	return n.open
}

// SetOpen sets the current open state for this node.
func (n *LibraryNode) SetOpen(open bool) {
	if open != n.open {
		n.open = open
		n.nav.adjustTableSize()
	}
}
