package navigator

import (
	"github.com/richardwilkes/gcs/internal/library"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

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

// ColumnCell returns the cell for the given column index.
func (n *LibraryNode) ColumnCell(index int) unison.Paneler {
	switch index {
	case 0:
		if n.open {
			return createNodeCell(OpenFolder, n.library.Title())
		}
		return createNodeCell(ClosedFolder, n.library.Title())
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
