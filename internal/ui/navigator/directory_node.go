package navigator

import (
	"io/fs"
	"path"
	"strings"

	"github.com/richardwilkes/gcs/internal/library"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var _ unison.TableRowData = &DirectoryNode{}

// DirectoryNode holds a directory in the navigator.
type DirectoryNode struct {
	nav      *Navigator
	library  *library.Library
	path     string
	children []unison.TableRowData
	open     bool
}

// NewDirectoryNode creates a new DirectoryNode.
func NewDirectoryNode(nav *Navigator, lib *library.Library, dirPath string) *DirectoryNode {
	n := &DirectoryNode{
		nav:     nav,
		library: lib,
		path:    dirPath,
	}
	n.Refresh()
	return n
}

// Refresh the contents of this node.
func (n *DirectoryNode) Refresh() {
	n.children = refreshChildren(n.nav, n.library, n.path)
}

func refreshChildren(nav *Navigator, lib *library.Library, dirPath string) []unison.TableRowData {
	libFS := lib.FS()
	entries, err := fs.ReadDir(libFS, dirPath)
	if err != nil {
		jot.Error(errs.NewWithCausef(err, "unable to read the directory: %s", dirPath))
		return nil
	}
	children := make([]unison.TableRowData, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, ".") {
			p := path.Join(dirPath, name)
			isDir := entry.IsDir()
			if entry.Type() == fs.ModeSymlink {
				var sub []fs.DirEntry
				if sub, err = fs.ReadDir(libFS, p); err == nil && len(sub) > 0 {
					isDir = true
				}
			}
			if isDir {
				dirNode := NewDirectoryNode(nav, lib, p)
				if dirNode.recursiveFileCount() > 0 {
					children = append(children, dirNode)
				}
			} else if _, exists := fileTypes[strings.ToLower(path.Ext(name))]; exists {
				children = append(children, NewFileNode(lib, p))
			}
		}
	}
	return children
}

// CanHaveChildRows always returns true.
func (n *DirectoryNode) CanHaveChildRows() bool {
	return true
}

// ChildRows returns the children of this node.
func (n *DirectoryNode) ChildRows() []unison.TableRowData {
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *DirectoryNode) CellDataForSort(index int) string {
	switch index {
	case 0:
		return path.Base(n.path)
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *DirectoryNode) ColumnCell(index int, _ bool) unison.Paneler {
	switch index {
	case 0:
		title := path.Base(n.path)
		if n.open {
			return createNodeCell(OpenFolder, title)
		}
		return createNodeCell(ClosedFolder, title)
	default:
		jot.Fatalf(1, "column index out of range (0-0): %d", index)
		return nil
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

func (n *DirectoryNode) recursiveFileCount() int {
	count := 0
	for _, child := range n.children {
		switch node := child.(type) {
		case *FileNode:
			count++
		case *DirectoryNode:
			count += node.recursiveFileCount()
		}
	}
	return count
}
