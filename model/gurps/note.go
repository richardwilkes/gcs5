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
	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/id"
)

const (
	noteTextKey = "text"
	noteTypeKey = "note"
)

// Note holds a note.
type Note struct {
	Parent    *Note
	ID        uuid.UUID
	Text      string
	PageRef   string
	Children  []*Note
	Container bool
	Open      bool
}

// NewNoteFromJSON creates a new Note from a JSON object.
func NewNoteFromJSON(parent *Note, data map[string]interface{}) *Note {
	n := &Note{Parent: parent}
	n.Container = encoding.String(data[commonTypeKey]) == noteTypeKey+commonContainerKeyPostfix
	n.ID = id.ParseOrNewUUID(encoding.String(data[commonIDKey]))
	n.Text = encoding.String(data[noteTextKey])
	n.PageRef = encoding.String(data[commonPageRefKey])
	if n.Container {
		n.Open = encoding.Bool(data[commonOpenKey])
		array := encoding.Array(data[commonChildrenKey])
		if len(array) != 0 {
			n.Children = make([]*Note, len(array))
			for i, one := range array {
				n.Children[i] = NewNoteFromJSON(n, encoding.Object(one))
			}
		}
	}
	return n
}

// ToJSON emits this object as JSON.
func (n *Note) ToJSON(encoder *encoding.JSONEncoder, entity *Entity) {
	encoder.StartObject()
	typeString := noteTypeKey
	if n.Container {
		typeString += commonContainerKeyPostfix
	}
	encoder.KeyedString(commonTypeKey, typeString, false, false)
	encoder.KeyedString(commonIDKey, n.ID.String(), false, false)
	encoder.KeyedString(noteTextKey, n.Text, true, true)
	encoder.KeyedString(commonPageRefKey, n.PageRef, true, true)
	if n.Container {
		encoder.KeyedBool(commonOpenKey, n.Open, true)
		if len(n.Children) != 0 {
			encoder.Key(commonChildrenKey)
			encoder.StartArray()
			for _, one := range n.Children {
				one.ToJSON(encoder, entity)
			}
			encoder.EndArray()
		}
	}
	encoder.EndObject()
}
