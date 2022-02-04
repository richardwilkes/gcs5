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

package attribute

import (
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible ThresholdOp values.
const (
	Unknown    = ThresholdOp("unknown")
	HalveMove  = ThresholdOp("halve_move")
	HalveDodge = ThresholdOp("halve_dodge")
	HalveST    = ThresholdOp("halve_st")
)

// AllThresholdOps is the complete set of ThresholdOp values.
var AllThresholdOps = []ThresholdOp{
	Unknown,
	HalveMove,
	HalveDodge,
	HalveST,
}

// ThresholdOp holds an operation to apply when a pool threshold is hit.
type ThresholdOp string

// EnsureValid ensures this is of a known value.
func (t ThresholdOp) EnsureValid() ThresholdOp {
	for _, one := range AllThresholdOps {
		if one == t {
			return t
		}
	}
	return AllThresholdOps[0]
}

// String implements fmt.Stringer.
func (t ThresholdOp) String() string {
	switch t {
	case Unknown:
		return i18n.Text("Unknown")
	case HalveMove:
		return i18n.Text("Halve Move")
	case HalveDodge:
		return i18n.Text("Halve Dodge")
	case HalveST:
		return i18n.Text("Halve Strength")
	default:
		return Unknown.String()
	}
}

// Description of this ThresholdOp's function.
func (t ThresholdOp) Description() string {
	switch t {
	case Unknown:
		return i18n.Text("Unknown")
	case HalveMove:
		return i18n.Text("Halve Move (round up)")
	case HalveDodge:
		return i18n.Text("Halve Dodge (round up)")
	case HalveST:
		return i18n.Text("Halve Strength (round up; does not affect HP and damage)")
	default:
		return Unknown.String()
	}
}
