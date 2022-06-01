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
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/unison"
)

var _ Node = &ConditionalModifier{}

// Columns that can be used with the conditional modifier method .CellData()
const (
	ConditionalModifierValueColumn = iota
	ConditionalModifierDescriptionColumn
)

// ConditionalModifier holds data for a reaction or conditional modifier.
type ConditionalModifier struct {
	From    string
	Amounts []fxp.Int
	Sources []string
}

// NewReaction creates a new ConditionalModifier.
func NewReaction(source, from string, amt fxp.Int) *ConditionalModifier {
	return &ConditionalModifier{
		From:    from,
		Amounts: []fxp.Int{amt},
		Sources: []string{source},
	}
}

// Add another source.
func (m *ConditionalModifier) Add(source string, amt fxp.Int) {
	m.Amounts = append(m.Amounts, amt)
	m.Sources = append(m.Sources, source)
}

// Total returns the total of all amounts.
func (m *ConditionalModifier) Total() fxp.Int {
	var total fxp.Int
	for _, amt := range m.Amounts {
		total += amt
	}
	return total
}

// Less returns true if this should be sorted above the other.
func (m *ConditionalModifier) Less(other *ConditionalModifier) bool {
	if txt.NaturalLess(m.From, other.From, true) {
		return true
	}
	if m.From != other.From {
		return false
	}
	if m.Total() < other.Total() {
		return true
	}
	return false
}

// UUID returns the UUID of this data.
func (m *ConditionalModifier) UUID() uuid.UUID {
	return uuid.UUID{}
}

// Kind returns the kind of data.
func (m *ConditionalModifier) Kind() string {
	return i18n.Text("Conditional Modifier")
}

// Container returns true if this is a container.
func (m *ConditionalModifier) Container() bool {
	return false
}

// Open returns true if this node is currently open.
func (m *ConditionalModifier) Open() bool {
	return false
}

// SetOpen sets the current open state for this node.
func (m *ConditionalModifier) SetOpen(_ bool) {
}

// NodeChildren returns the children of this node, if any.
func (m *ConditionalModifier) NodeChildren() []Node {
	return nil
}

// CellData returns the cell data information for the given column.
func (m *ConditionalModifier) CellData(column int, data *CellData) {
	switch column {
	case ConditionalModifierValueColumn:
		data.Type = Text
		data.Primary = m.Total().StringWithSign()
		data.Alignment = unison.EndAlignment
	case ConditionalModifierDescriptionColumn:
		data.Type = Text
		data.Primary = m.From
	case PageRefCellAlias:
		data.Type = PageRef
	}
}

// OwningEntity returns the owning Entity.
func (m *ConditionalModifier) OwningEntity() *Entity {
	return nil
}

// SetOwningEntity sets the owning entity and configures any sub-components as needed.
func (m *ConditionalModifier) SetOwningEntity(_ *Entity) {
}
