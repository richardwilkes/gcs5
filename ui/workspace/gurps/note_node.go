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

package gurps

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

const (
	noteDescriptionColumn = iota
	noteWhenColumn
	noteCategoryColumn
	noteReferenceColumn
	noteColumnCount
)

var (
	_ unison.TableRowData = &NoteNode{}
	_ tbl.Matcher         = &NoteNode{}
)

// NoteNode holds a note in the note list.
type NoteNode struct {
	table     *unison.Table
	parent    *NoteNode
	note      *gurps.Note
	children  []unison.TableRowData
	cellCache []*tbl.CellCache
}

// NewNoteListDockable creates a new unison.Dockable for note list files.
func NewNoteListDockable(filePath string) (unison.Dockable, error) {
	notes, err := gurps.NewNotesFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	return NewTableDockable(filePath, []unison.TableColumnHeader{
		tbl.NewHeader(i18n.Text("Note"), "", false),
		tbl.NewHeader(i18n.Text("When"), "", false),
		tbl.NewHeader(i18n.Text("Category"), "", false),
		tbl.NewPageRefHeader(false),
	}, func(table *unison.Table) []unison.TableRowData {
		rows := make([]unison.TableRowData, 0, len(notes))
		for _, one := range notes {
			rows = append(rows, NewNoteNode(table, nil, one))
		}
		return rows
	}), nil
}

// NewNoteNode creates a new NoteNode.
func NewNoteNode(table *unison.Table, parent *NoteNode, note *gurps.Note) *NoteNode {
	n := &NoteNode{
		table:     table,
		parent:    parent,
		note:      note,
		cellCache: make([]*tbl.CellCache, noteColumnCount),
	}
	return n
}

// ParentRow returns the parent row, or nil if this is a root node.
func (n *NoteNode) ParentRow() unison.TableRowData {
	return n.parent
}

// CanHaveChildRows always returns true.
func (n *NoteNode) CanHaveChildRows() bool {
	return n.note.Container()
}

// ChildRows returns the children of this node.
func (n *NoteNode) ChildRows() []unison.TableRowData {
	if n.note.Container() && n.children == nil {
		n.children = make([]unison.TableRowData, len(n.note.Children))
		for i, one := range n.note.Children {
			n.children[i] = NewNoteNode(n.table, n, one)
		}
	}
	return n.children
}

// CellDataForSort returns the string that represents the data in the specified cell.
func (n *NoteNode) CellDataForSort(index int) string {
	switch index {
	case noteDescriptionColumn:
		return n.note.Text
	case noteWhenColumn:
		return n.note.When
	case noteCategoryColumn:
		return strings.Join(n.note.Categories, ", ")
	case noteReferenceColumn:
		return n.note.PageRef
	default:
		return ""
	}
}

// ColumnCell returns the cell for the given column index.
func (n *NoteNode) ColumnCell(row, col int, selected bool) unison.Paneler {
	width := n.table.CellWidth(row, col)
	data := n.CellDataForSort(col)
	if n.cellCache[col].Matches(width, data) {
		color := unison.DefaultLabelTheme.OnBackgroundInk
		if selected {
			color = unison.OnSelectionColor
		}
		for _, child := range n.cellCache[col].Panel.Children() {
			child.Self.(*unison.Label).LabelTheme.OnBackgroundInk = color
		}
		return n.cellCache[col].Panel
	}
	p := &unison.Panel{}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{Columns: 1})
	switch col {
	case noteReferenceColumn:
		tbl.CreateAndAddPageRefCellLabel(p, n.CellDataForSort(col), "", unison.DefaultLabelTheme.Font, selected)
	default:
		tbl.CreateAndAddCellLabel(p, width, n.CellDataForSort(col), unison.DefaultLabelTheme.Font, selected)
	}
	n.cellCache[col] = &tbl.CellCache{
		Width: width,
		Data:  data,
		Panel: p,
	}
	return p
}

// IsOpen returns true if this node should display its children.
func (n *NoteNode) IsOpen() bool {
	return n.note.Container() && n.note.Open
}

// SetOpen sets the current open state for this node.
func (n *NoteNode) SetOpen(open bool) {
	if n.note.Container() && open != n.note.Open {
		n.note.Open = open
		n.table.SyncToModel()
	}
}

// Match implements Matcher.
func (n *NoteNode) Match(text string) bool {
	return strings.Contains(strings.ToLower(n.note.Text), text) ||
		strings.Contains(strings.ToLower(n.note.When), text) ||
		strings.Contains(strings.ToLower(n.note.PageRef), text) ||
		stringSliceContains(n.note.Categories, text)
}
