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

package editors

import (
	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var noteColMap = map[int]int{
	0: gurps.NoteTextColumn,
	1: gurps.NoteReferenceColumn,
}

type notesProvider struct {
	table    *unison.Table[*Node[*gurps.Note]]
	provider gurps.NoteListProvider
	forPage  bool
}

// NewNotesProvider creates a new table provider for notes.
func NewNotesProvider(provider gurps.NoteListProvider, forPage bool) widget.TableProvider[*Node[*gurps.Note]] {
	return &notesProvider{
		provider: provider,
		forPage:  forPage,
	}
}

func (p *notesProvider) SetTable(table *unison.Table[*Node[*gurps.Note]]) {
	p.table = table
}

func (p *notesProvider) RootRowCount() int {
	return len(p.provider.NoteList())
}

func (p *notesProvider) RootRows() []*Node[*gurps.Note] {
	data := p.provider.NoteList()
	rows := make([]*Node[*gurps.Note], 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode[*gurps.Note](p.table, nil, noteColMap, one, p.forPage))
	}
	return rows
}

func (p *notesProvider) SetRootRows(rows []*Node[*gurps.Note]) {
	p.provider.SetNoteList(ExtractNodeDataFromList(rows))
}

func (p *notesProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *notesProvider) DragKey() string {
	return gid.Note
}

func (p *notesProvider) DragSVG() *unison.SVG {
	return res.GCSNotesSVG
}

func (p *notesProvider) DropShouldMoveData(drop *unison.TableDrop[*Node[*gurps.Note]]) bool {
	return drop.Table == drop.TableDragData.Table
}

func (p *notesProvider) DropCopyRow(drop *unison.TableDrop[*Node[*gurps.Note]], row *Node[*gurps.Note]) *Node[*gurps.Note] {
	note := ExtractFromRowData[*gurps.Note](row).Clone(p.provider.Entity(), nil)
	return NewNode[*gurps.Note](drop.Table, nil, noteColMap, note, p.forPage)
}

func (p *notesProvider) DropSetRowChildren(_ *unison.TableDrop[*Node[*gurps.Note]], row *Node[*gurps.Note], children []*Node[*gurps.Note]) {
	list := make([]*gurps.Note, 0, len(children))
	for _, child := range children {
		list = append(list, ExtractFromRowData[*gurps.Note](child))
	}
	if row == nil {
		p.provider.SetNoteList(list)
	} else {
		ExtractFromRowData[*gurps.Note](row).Children = list
		row.children = nil
	}
}

func (p *notesProvider) ItemNames() (singular, plural string) {
	return i18n.Text("Note"), i18n.Text("Notes")
}

func (p *notesProvider) Headers() []unison.TableColumnHeader[*Node[*gurps.Note]] {
	var headers []unison.TableColumnHeader[*Node[*gurps.Note]]
	for i := 0; i < len(noteColMap); i++ {
		switch noteColMap[i] {
		case gurps.NoteTextColumn:
			headers = append(headers, NewHeader[*gurps.Note](i18n.Text("Note"), "", p.forPage))
		case gurps.NoteReferenceColumn:
			headers = append(headers, NewPageRefHeader[*gurps.Note](p.forPage))
		default:
			jot.Fatalf(1, "invalid note column: %d", noteColMap[i])
		}
	}
	return headers
}

func (p *notesProvider) SyncHeader(_ []unison.TableColumnHeader[*Node[*gurps.Note]]) {
}

func (p *notesProvider) HierarchyColumnIndex() int {
	for k, v := range noteColMap {
		if v == gurps.NoteTextColumn {
			return k
		}
	}
	return 0
}

func (p *notesProvider) ExcessWidthColumnIndex() int {
	return p.HierarchyColumnIndex()
}

func (p *notesProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table[*Node[*gurps.Note]]) {
	OpenEditor[*gurps.Note](table, func(item *gurps.Note) { EditNote(owner, item) })
}

func (p *notesProvider) CreateItem(owner widget.Rebuildable, table *unison.Table[*Node[*gurps.Note]], variant widget.ItemVariant) {
	item := gurps.NewNote(p.Entity(), nil, variant == widget.ContainerItemVariant)
	InsertItem[*gurps.Note](owner, table, item,
		func(target *gurps.Note) []*gurps.Note { return target.Children },
		func(target *gurps.Note, children []*gurps.Note) { target.Children = children },
		p.provider.NoteList, p.provider.SetNoteList,
		func(_ *unison.Table[*Node[*gurps.Note]]) []*Node[*gurps.Note] { return p.RootRows() },
		func(target *gurps.Note) uuid.UUID { return target.ID })
	EditNote(owner, item)
}

func (p *notesProvider) DeleteSelection(table *unison.Table[*Node[*gurps.Note]]) {
	deleteTableSelection(table, p.provider.NoteList(),
		func(nodes []*gurps.Note) { p.provider.SetNoteList(nodes) },
		func(node *gurps.Note) *[]*gurps.Note { return &node.Children })
}
