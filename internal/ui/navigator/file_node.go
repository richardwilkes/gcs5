package navigator

import (
	"io/fs"
	"path"

	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/fa"
)

// FileNode holds a file in the navigator.
type FileNode struct {
	fs   fs.FS
	path string
}

// NewFileNode creates a new FileNode.
func NewFileNode(owningFS fs.FS, filePath string) *FileNode {
	return &FileNode{
		fs:   owningFS,
		path: filePath,
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
		return createNodeLabel(fa.File, path.Base(n.path))
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

func createNodeLabel(icon, title string) *unison.Panel {
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 8,
	})
	label := unison.NewLabel()
	label.Text = icon
	faDesc := unison.FontDescriptor{
		Family:  unison.FontAwesomeFreeFamilyName,
		Size:    10,
		Weight:  unison.BlackFontWeight,
		Spacing: unison.StandardSpacing,
		Slant:   unison.NoSlant,
	}
	label.Font = faDesc.Font()
	label.SetLayoutData(&unison.FlexLayoutData{
		HSpan: 1,
		VSpan: 1,
	})
	panel.AddChild(label)
	label = unison.NewLabel()
	label.Text = title
	label.SetLayoutData(&unison.FlexLayoutData{
		HSpan: 1,
		VSpan: 1,
		HGrab: true,
	})
	panel.AddChild(label)
	panel.NeedsLayout = true
	return panel
}
