package navigator

import (
	"io/fs"
	"os"

	"github.com/richardwilkes/gcs/internal/library"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/fa"
)

// LibraryNode holds a library in the navigator.
type LibraryNode struct {
	nav      *Navigator
	library  *library.Library
	fs       fs.FS
	children []unison.TableRowData
	open     bool
}

// NewLibraryNode creates a new LibraryNode.
func NewLibraryNode(nav *Navigator, lib *library.Library) *LibraryNode {
	n := &LibraryNode{
		nav:     nav,
		library: lib,
		fs:      os.DirFS(lib.Path), // Should change this to obtain it from the library on demand
	}
	n.Refresh()
	return n
}

// Refresh the contents of this node.
func (n *LibraryNode) Refresh() {
	n.children = nil
	entries, err := fs.ReadDir(n.fs, ".")
	if err != nil {
		jot.Error(errs.NewWithCausef(err, "unable to read the directory: %s", n.library.Path))
		return
	}
	n.children = make([]unison.TableRowData, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			n.children = append(n.children, NewDirectoryNode(n.nav, n.fs, entry.Name()))
		} else {
			n.children = append(n.children, NewFileNode(n.fs, entry.Name()))
		}
	}
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
			return createNodeLabel(fa.FolderOpen, n.library.Title)
		}
		return createNodeLabel(fa.Folder, n.library.Title)
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
