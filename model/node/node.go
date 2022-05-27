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

package node

import (
	"github.com/google/uuid"
)

// Node defines the methods required of nodes in our tables.
type Node interface {
	UUID() uuid.UUID
	Kind() string
	Container() bool
	NodeChildren() []Node
	Open() bool
	SetOpen(open bool)
	CellData(column int, data *CellData)
}

// EditorData defines the methods required of editor data.
type EditorData[T Node] interface {
	// CopyFrom copies the corresponding data from the node into this editor data.
	CopyFrom(T)
	// ApplyTo copes he editor data into the provided node.
	ApplyTo(T)
}
