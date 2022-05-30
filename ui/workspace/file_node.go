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
	"strings"

	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var _ unison.TableRowData = &FileNode{}

// FileNode holds a file in the navigator.
type FileNode struct {
	library *library.Library
	path    string
	parent  unison.TableRowData
}

// NewFileNode creates a new FileNode.
func NewFileNode(lib *library.Library, filePath string, parent unison.TableRowData) *FileNode {
	return &FileNode{
		library: lib,
		path:    filePath,
		parent:  parent,
	}
}

// ParentRow returns the parent row, or nil if this is a root node.
func (n *FileNode) ParentRow() unison.TableRowData {
	return n.parent
}

// CanHaveChildRows always returns false.
func (n *FileNode) CanHaveChildRows() bool {
	return false
}

// ChildRows always returns nil.
func (n *FileNode) ChildRows() []unison.TableRowData {
	return nil
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *FileNode) CellDataForSort(index int) string {
	switch index {
	case 0:
		return path.Base(n.path)
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *FileNode) ColumnCell(_, col int, foreground, _ unison.Ink, _, _, _ bool) unison.Paneler {
	switch col {
	case 0:
		name := path.Base(n.path)
		return createNodeCell(strings.ToLower(path.Ext(name)), name, foreground)
	default:
		jot.Errorf("column index out of range (0-0): %d", col)
		return unison.NewLabel()
	}
}

// IsOpen always returns false.
func (n *FileNode) IsOpen() bool {
	return false
}

// SetOpen does nothing.
func (n *FileNode) SetOpen(_ bool) {
}
