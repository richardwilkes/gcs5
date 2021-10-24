package navigator

import (
	"io/fs"
	"path"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/fa"
)

// DirectoryNode holds a directory in the navigator.
type DirectoryNode struct {
	nav      *Navigator
	fs       fs.FS
	path     string
	children []unison.TableRowData
	open     bool
}

// NewDirectoryNode creates a new DirectoryNode.
func NewDirectoryNode(nav *Navigator, owningFS fs.FS, dirPath string) *DirectoryNode {
	n := &DirectoryNode{
		nav:  nav,
		fs:   owningFS,
		path: dirPath,
	}
	n.Refresh()
	return n
}

// Refresh the contents of this node.
func (n *DirectoryNode) Refresh() {
	n.children = nil
	entries, err := fs.ReadDir(n.fs, n.path)
	if err != nil {
		jot.Error(errs.NewWithCausef(err, "unable to read the directory: %s", n.path))
		return
	}
	n.children = make([]unison.TableRowData, 0, len(entries))
	for _, entry := range entries {
		p := path.Join(n.path, entry.Name())
		if entry.IsDir() {
			n.children = append(n.children, NewDirectoryNode(n.nav, n.fs, p))
		} else {
			n.children = append(n.children, NewFileNode(n.fs, p))
		}
	}
}

// CanHaveChildRows always returns true.
func (n *DirectoryNode) CanHaveChildRows() bool {
	return true
}

// ChildRows returns the children of this node.
func (n *DirectoryNode) ChildRows() []unison.TableRowData {
	return n.children
}

// ColumnCell returns the cell for the given column index.
func (n *DirectoryNode) ColumnCell(index int) unison.Paneler {
	switch index {
	case 0:
		title := path.Base(n.path)
		if n.open {
			return createNodeLabel(fa.FolderOpen, title)
		}
		return createNodeLabel(fa.Folder, title)
	default:
		jot.Errorf("column index out of range (0-0): %d", index)
		return unison.NewLabel()
	}
}

// IsOpen returns true if this node should display its children.
func (n *DirectoryNode) IsOpen() bool {
	return n.open
}

// SetOpen sets the current open state for this node.
func (n *DirectoryNode) SetOpen(open bool) {
	if open != n.open {
		n.open = open
		n.nav.adjustTableSize()
	}
}
