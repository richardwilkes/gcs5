package navigator

import (
	"path"
	"strings"

	"github.com/richardwilkes/gcs/internal/library"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

// FileNode holds a file in the navigator.
type FileNode struct {
	library *library.Library
	path    string
}

// NewFileNode creates a new FileNode.
func NewFileNode(lib *library.Library, filePath string) *FileNode {
	return &FileNode{
		library: lib,
		path:    filePath,
	}
}

// CanHaveChildRows always returns false.
func (n *FileNode) CanHaveChildRows() bool {
	return false
}

// ChildRows always returns nil.
func (n *FileNode) ChildRows() []unison.TableRowData {
	return nil
}

// ColumnCell returns the cell for the given column index.
func (n *FileNode) ColumnCell(index int) unison.Paneler {
	switch index {
	case 0:
		name := path.Base(n.path)
		return createNodeCell(strings.ToLower(path.Ext(name)), name)
	default:
		jot.Errorf("column index out of range (0-0): %d", index)
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
