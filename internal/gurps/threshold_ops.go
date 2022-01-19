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

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible ThresholdOp values.
const (
	Unknown ThresholdOp = iota
	HalveMove
	HalveDodge
	HalveST
)

// ThresholdOp holds an operation to apply when a pool threshold is hit.
type ThresholdOp uint8

// ThresholdOpFromString extracts a ThresholdOp from a string.
func ThresholdOpFromString(str string) ThresholdOp {
	for op := Unknown; op <= HalveST; op++ {
		if strings.EqualFold(op.Key(), str) {
			return op
		}
	}
	return Unknown
}

// Key returns the key used to represent this ThresholdOp.
func (op ThresholdOp) Key() string {
	switch op {
	case HalveMove:
		return "halve_move"
	case HalveDodge:
		return "halve_dodge"
	case HalveST:
		return "halve_st"
	default: // Unknown
		return "unknown"
	}
}

// MarshalText implements encoding.TextMarshaler.
func (op ThresholdOp) MarshalText() (text []byte, err error) {
	return []byte(op.Key()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (op *ThresholdOp) UnmarshalText(text []byte) error {
	*op = ThresholdOpFromString(string(text))
	return nil
}

// String implements fmt.Stringer.
func (op ThresholdOp) String() string {
	switch op {
	case HalveMove:
		return i18n.Text("Halve Move")
	case HalveDodge:
		return i18n.Text("Halve Dodge")
	case HalveST:
		return i18n.Text("Halve ST")
	default: // Unknown
		return i18n.Text("Unknown")
	}
}

// Description of this ThresholdOp's function.
func (op ThresholdOp) Description() string {
	switch op {
	case HalveMove:
		return i18n.Text("Halve Move (round up)")
	case HalveDodge:
		return i18n.Text("Halve Dodge (round up)")
	case HalveST:
		return i18n.Text("Halve ST (round up; does not affect HP and damage)")
	default: // Unknown
		return i18n.Text("Unknown")
	}
}
