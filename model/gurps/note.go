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
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
)

var _ node.Node = &Note{}

// Columns that can be used with the note method .CellData()
const (
	NoteTextColumn = iota
	NoteReferenceColumn
)

const (
	noteListTypeKey = "note_list"
	noteTypeKey     = "note"
)

// NoteContainer holds the Note data that only exists in containers.
type NoteContainer struct {
	Children []*Note `json:"children,omitempty"`
	Open     bool    `json:"open,omitempty"`
}

// NoteData holds the Note data that is written to disk.
type NoteData struct {
	Type           string    `json:"type"`
	ID             uuid.UUID `json:"id"`
	Text           string    `json:"text,omitempty"`
	PageRef        string    `json:"reference,omitempty"`
	*NoteContainer `json:",omitempty"`
}

// Note holds a note.
type Note struct {
	NoteData
	Parent *Note
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
func NewNote(parent *Note, container bool) *Note {
	n := Note{
		NoteData: NoteData{
			Type: noteTypeKey,
			ID:   id.NewUUID(),
		},
		Parent: parent,
	}
	if container {
		n.Type += commonContainerKeyPostfix
		n.NoteContainer = &NoteContainer{Open: true}
	}
	return &n
}

// MarshalJSON implements json.Marshaler.
func (n *Note) MarshalJSON() ([]byte, error) {
	if !n.Container() {
		n.NoteContainer = nil
	}
	return json.Marshal(&n.NoteData)
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *Note) UnmarshalJSON(data []byte) error {
	n.NoteData = NoteData{}
	if err := json.Unmarshal(data, &n.NoteData); err != nil {
		return err
	}
	if n.Container() {
		if n.NoteContainer == nil {
			n.NoteContainer = &NoteContainer{}
		}
		for _, one := range n.Children {
			one.Parent = n
		}
	}
	return nil
}

// UUID returns the UUID of this data.
func (n *Note) UUID() uuid.UUID {
	return n.ID
}

// Kind returns the kind of data.
func (n *Note) Kind() string {
	if n.Container() {
		return i18n.Text("Note Container")
	}
	return i18n.Text("Note")
}

// Container returns true if this is a container.
func (n *Note) Container() bool {
	return strings.HasSuffix(n.Type, commonContainerKeyPostfix)
}

// Open returns true if this node is currently open.
func (n *Note) Open() bool {
	if n.Container() {
		return n.NoteContainer.Open
	}
	return false
}

// SetOpen sets the current open state for this node.
func (n *Note) SetOpen(open bool) {
	if n.Container() {
		n.NoteContainer.Open = open
	}
}

// NodeChildren returns the children of this node, if any.
func (n *Note) NodeChildren() []node.Node {
	if n.Container() {
		children := make([]node.Node, len(n.Children))
		for i, child := range n.Children {
			children[i] = child
		}
		return children
	}
	return nil
}

// CellData returns the cell data information for the given column.
func (n *Note) CellData(column int, data *node.CellData) {
	switch column {
	case NoteTextColumn:
		data.Type = node.Text
		data.Primary = n.Text
	case NoteReferenceColumn:
		data.Type = node.PageRef
		data.Primary = n.PageRef
		data.Secondary = n.Text
	}
}

// Depth returns the number of parents this node has.
func (n *Note) Depth() int {
	count := 0
	p := n.Parent
	for p != nil {
		count++
		p = p.Parent
	}
	return count
}
