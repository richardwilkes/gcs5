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
	"context"
	"io/fs"

	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
)

var _ Node[*Note] = &Note{}

// Columns that can be used with the note method .CellData()
const (
	NoteTextColumn = iota
	NoteReferenceColumn
)

const (
	noteListTypeKey = "note_list"
	noteTypeKey     = "note"
)

// Note holds a note.
type Note struct {
	NoteData
	Entity *Entity
}

type noteListData struct {
	Type    string  `json:"type"`
	Version int     `json:"version"`
	Rows    []*Note `json:"rows"`
}

// NewNotesFromFile loads an Note list from a file.
func NewNotesFromFile(fileSystem fs.FS, filePath string) ([]*Note, error) {
	var data noteListData
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause(gid.InvalidFileDataMsg, err)
	}
	if data.Type != noteListTypeKey {
		return nil, errs.New(gid.UnexpectedFileDataMsg)
	}
	if err := gid.CheckVersion(data.Version); err != nil {
		return nil, err
	}
	return data.Rows, nil
}

// SaveNotes writes the Note list to the file as JSON.
func SaveNotes(notes []*Note, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &noteListData{
		Type:    noteListTypeKey,
		Version: gid.CurrentDataVersion,
		Rows:    notes,
	})
}

// NewNote creates a new Note.
func NewNote(entity *Entity, parent *Note, container bool) *Note {
	n := &Note{
		NoteData: NoteData{
			ContainerBase: newContainerBase[*Note](noteTypeKey, container),
		},
		Entity: entity,
	}
	n.Text = n.Kind()
	n.parent = parent
	return n
}

// Clone this Note, assigning a new id, entity and parent.
func (n *Note) Clone(entity *Entity, parent *Note) *Note {
	other := NewNote(entity, parent, n.Container())
	other.IsOpen = n.IsOpen
	other.NoteEditData.CopyFrom(n)
	if n.HasChildren() {
		other.Children = make([]*Note, 0, len(n.Children))
		for _, child := range n.Children {
			other.Children = append(other.Children, child.Clone(entity, other))
		}
	}
	return other
}

// MarshalJSON implements json.Marshaler.
func (n *Note) MarshalJSON() ([]byte, error) {
	n.ClearUnusedFieldsForType()
	return json.Marshal(&n.NoteData)
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *Note) UnmarshalJSON(data []byte) error {
	n.NoteData = NoteData{}
	if err := json.Unmarshal(data, &n.NoteData); err != nil {
		return err
	}
	n.ClearUnusedFieldsForType()
	if n.Container() {
		for _, one := range n.Children {
			one.parent = n
		}
	}
	return nil
}

// CellData returns the cell data information for the given column.
func (n *Note) CellData(column int, data *CellData) {
	switch column {
	case NoteTextColumn:
		data.Type = Text
		data.Primary = n.Text
	case NoteReferenceColumn, PageRefCellAlias:
		data.Type = PageRef
		data.Primary = n.PageRef
		data.Secondary = n.Text
	}
}

// Depth returns the number of parents this node has.
func (n *Note) Depth() int {
	count := 0
	p := n.parent
	for p != nil {
		count++
		p = p.parent
	}
	return count
}

// OwningEntity returns the owning Entity.
func (n *Note) OwningEntity() *Entity {
	return n.Entity
}

// SetOwningEntity sets the owning entity and configures any sub-components as needed.
func (n *Note) SetOwningEntity(entity *Entity) {
	n.Entity = entity
	if n.Container() {
		for _, child := range n.Children {
			child.SetOwningEntity(entity)
		}
	}
}

// Enabled returns true if this node is enabled.
func (n *Note) Enabled() bool {
	return true
}
