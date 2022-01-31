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
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/id"
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
	PageRef        string    `json:"reference,omitempty"`
	*NoteContainer `json:",omitempty"`
}

// Note holds a note.
type Note struct {
	NoteData
	Parent *Note
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
