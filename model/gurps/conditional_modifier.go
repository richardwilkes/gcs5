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
	"github.com/richardwilkes/gcs/model/node"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
	"github.com/richardwilkes/unison"
)

var _ node.Node = &ConditionalModifier{}

// Columns that can be used with the conditional modifier method .CellData()
const (
	ConditionalModifierValueColumn = iota
	ConditionalModifierDescriptionColumn
)

// ConditionalModifier holds data for a reaction or conditional modifier.
type ConditionalModifier struct {
	From    string
	Amounts []f64d4.Int
	Sources []string
}

// NewReaction creates a new ConditionalModifier.
func NewReaction(source, from string, amt f64d4.Int) *ConditionalModifier {
	return &ConditionalModifier{
		From:    from,
		Amounts: []f64d4.Int{amt},
		Sources: []string{source},
	}
}

// Add another source.
func (m *ConditionalModifier) Add(source string, amt f64d4.Int) {
	m.Amounts = append(m.Amounts, amt)
	m.Sources = append(m.Sources, source)
}

// Total returns the total of all amounts.
func (m *ConditionalModifier) Total() f64d4.Int {
	var total f64d4.Int
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

// Container returns true if this is a container.
func (m *ConditionalModifier) Container() bool {
	return false
}

// Open returns true if this node is currently open.
func (m *ConditionalModifier) Open() bool {
	return false
}

// SetOpen sets the current open state for this node.
func (m *ConditionalModifier) SetOpen(open bool) {
}

// NodeChildren returns the children of this node, if any.
func (m *ConditionalModifier) NodeChildren() []node.Node {
	return nil
}

// CellData returns the cell data information for the given column.
func (m *ConditionalModifier) CellData(column int, data *node.CellData) {
	switch column {
	case ConditionalModifierValueColumn:
		data.Type = node.Text
		data.Primary = m.Total().StringWithSign()
		data.Alignment = unison.EndAlignment
	case ConditionalModifierDescriptionColumn:
		data.Type = node.Text
		data.Primary = m.From
	}
}
