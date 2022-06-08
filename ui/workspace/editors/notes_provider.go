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
	provider gurps.NoteListProvider
	forPage  bool
}

// NewNotesProvider creates a new table provider for skills.
func NewNotesProvider(provider gurps.NoteListProvider, forPage bool) TableProvider {
	return &notesProvider{
		provider: provider,
		forPage:  forPage,
	}
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

func (p *notesProvider) ItemNames() (singular, plural string) {
	return i18n.Text("Note"), i18n.Text("Notes")
}

func (p *notesProvider) Headers() []unison.TableColumnHeader {
	var headers []unison.TableColumnHeader
	for i := 0; i < len(noteColMap); i++ {
		switch noteColMap[i] {
		case gurps.NoteTextColumn:
			headers = append(headers, NewHeader(i18n.Text("Note"), "", p.forPage))
		case gurps.NoteReferenceColumn:
			headers = append(headers, NewPageRefHeader(p.forPage))
		default:
			jot.Fatalf(1, "invalid note column: %d", noteColMap[i])
		}
	}
	return headers
}

func (p *notesProvider) RowData(table *unison.Table) []unison.TableRowData {
	data := p.provider.NoteList()
	rows := make([]unison.TableRowData, 0, len(data))
	for _, one := range data {
		rows = append(rows, NewNode(table, nil, noteColMap, one, p.forPage))
	}
	return rows
}

func (p *notesProvider) SyncHeader(_ []unison.TableColumnHeader) {
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

func (p *notesProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table) {
	OpenEditor[*gurps.Note](table, func(item *gurps.Note) { EditNote(owner, item) })
}

func (p *notesProvider) CreateItem(owner widget.Rebuildable, table *unison.Table, variant ItemVariant) {
	item := gurps.NewNote(p.Entity(), nil, variant == ContainerItemVariant)
	InsertItem[*gurps.Note](owner, table, item,
		func(target, parent *gurps.Note) { target.Parent = parent },
		func(target *gurps.Note) []*gurps.Note { return target.Children },
		func(target *gurps.Note, children []*gurps.Note) { target.Children = children },
		p.provider.NoteList, p.provider.SetNoteList, p.RowData,
		func(target *gurps.Note) uuid.UUID { return target.ID })
	EditNote(owner, item)
}

func (p *notesProvider) DeleteSelection(table *unison.Table) {
	deleteTableSelection(table, p.provider.NoteList(),
		func(nodes []*gurps.Note) { p.provider.SetNoteList(nodes) },
		func(node *gurps.Note) **gurps.Note { return &node.Parent },
		func(node *gurps.Note) *[]*gurps.Note { return &node.Children })
}
