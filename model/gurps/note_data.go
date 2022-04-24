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
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/toolbox/i18n"
)

var _ node.EditorData[*Note] = &NoteEditData{}

// NoteEditData holds the Note data that can be edited by the UI detail editor.
type NoteEditData struct {
	Text    string `json:"text,omitempty"`
	PageRef string `json:"reference,omitempty"`
}

// CopyFrom implements node.EditorData.
func (d *NoteEditData) CopyFrom(note *Note) {
	d.copyFrom(&note.NoteEditData)
}

// ApplyTo implements node.EditorData.
func (d *NoteEditData) ApplyTo(note *Note) {
	note.NoteEditData.copyFrom(d)
}

func (d *NoteEditData) copyFrom(other *NoteEditData) {
	*d = *other
}

// NoteData holds the Note data that is written to disk.
type NoteData struct {
	ID   uuid.UUID `json:"id"`
	Type string    `json:"type"`
	NoteEditData
	IsOpen   bool    `json:"open,omitempty"`     // Container only
	Children []*Note `json:"children,omitempty"` // Container only
}

// UUID returns the UUID of this data.
func (n *NoteData) UUID() uuid.UUID {
	return n.ID
}

// Kind returns the kind of data.
func (n *NoteData) Kind() string {
	if n.Container() {
		return i18n.Text("Note Container")
	}
	return i18n.Text("Note")
}

// Container returns true if this is a container.
func (n *NoteData) Container() bool {
	return strings.HasSuffix(n.Type, commonContainerKeyPostfix)
}

// Open returns true if this node is currently open.
func (n *NoteData) Open() bool {
	return n.IsOpen && n.Container()
}

// SetOpen sets the current open state for this node.
func (n *NoteData) SetOpen(open bool) {
	n.IsOpen = open && n.Container()
}

// ClearUnusedFieldsForType zeroes out the fields that are not applicable to this type (container vs not-container).
func (n *NoteData) ClearUnusedFieldsForType() {
	if !n.Container() {
		n.Children = nil
		n.IsOpen = false
	}
}

// NodeChildren returns the children of this node, if any.
func (n *NoteData) NodeChildren() []node.Node {
	if n.Container() {
		children := make([]node.Node, len(n.Children))
		for i, child := range n.Children {
			children[i] = child
		}
		return children
	}
	return nil
}
