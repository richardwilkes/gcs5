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
)

// SpellItem holds the Spell data that only exists in non-containers.
type SpellItem struct {
	Prereq  Prereq    `json:"prereqs,omitempty"`
	Weapons []*Weapon `json:"weapons,omitempty"`
}

// SpellContainer holds the Spell data that only exists in containers.
type SpellContainer struct {
	Children []*Spell `json:"children,omitempty"`
	Open     bool     `json:"open,omitempty"`
}

// SpellData holds the Spell data that is written to disk.
type SpellData struct {
	Type            string    `json:"type"`
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name,omitempty"`
	PageRef         string    `json:"reference,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	VTTNotes        string    `json:"vtt_notes,omitempty"`
	Categories      []string  `json:"categories,omitempty"`
	*SpellItem      `json:",omitempty"`
	*SpellContainer `json:",omitempty"`
}

// Spell holds the data for a spell.
type Spell struct {
	SpellData
	Entity            *Entity
	Parent            *Spell
	UnsatisfiedReason string
	Satisfied         bool
}

// Container returns true if this is a container.
func (s *Spell) Container() bool {
	return strings.HasSuffix(s.Type, commonContainerKeyPostfix)
}
