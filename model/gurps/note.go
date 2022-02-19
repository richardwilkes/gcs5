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
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/errs"
)

const noteTypeKey = "note"

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
	When           string    `json:"when,omitempty"`
	PageRef        string    `json:"reference,omitempty"`
	*NoteContainer `json:",omitempty"`
}

// Note holds a note.
type Note struct {
	NoteData
	Parent *Note
}

type noteListData struct {
	Current []*Note `json:"notes"`
}

// NewNotesFromFile loads an Note list from a file.
func NewNotesFromFile(fileSystem fs.FS, filePath string) ([]*Note, error) {
	var data struct {
		noteListData
		OldKey []*Note `json:"rows"`
	}
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, errs.NewWithCause("invalid notes file: "+filePath, err)
	}
	if len(data.Current) != 0 {
		return data.Current, nil
	}
	return data.OldKey, nil
}

// SaveNotes writes the Note list to the file as JSON.
func SaveNotes(notes []*Note, filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &noteListData{Current: notes})
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
		for _, one := range n.Children {
			one.Parent = n
		}
	}
	return nil
}

// Container returns true if this is a container.
func (n *Note) Container() bool {
	return strings.HasSuffix(n.Type, commonContainerKeyPostfix)
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
