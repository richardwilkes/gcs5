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
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// ConditionalModifier holds data for a reaction or conditional modifier.
type ConditionalModifier struct {
	From    string
	Amounts []fixed.F64d4
	Sources []string
}

// NewReaction creates a new ConditionalModifier.
func NewReaction(source, from string, amt fixed.F64d4) *ConditionalModifier {
	return &ConditionalModifier{
		From:    from,
		Amounts: []fixed.F64d4{amt},
		Sources: []string{source},
	}
}

// Add another source.
func (r *ConditionalModifier) Add(source string, amt fixed.F64d4) {
	r.Amounts = append(r.Amounts, amt)
	r.Sources = append(r.Sources, source)
}

// Total returns the total of all amounts.
func (r *ConditionalModifier) Total() fixed.F64d4 {
	var total fixed.F64d4
	for _, amt := range r.Amounts {
		total += amt
	}
	return total
}

// Less returns true if this should be sorted above the other.
func (r *ConditionalModifier) Less(other *ConditionalModifier) bool {
	if txt.NaturalLess(r.From, other.From, true) {
		return true
	}
	if r.From != other.From {
		return false
	}
	if r.Total() < other.Total() {
		return true
	}
	return false
}
