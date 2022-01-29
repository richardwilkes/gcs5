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

type thresholdOpData struct {
	Key         string
	String      string
	Description string
}

// ThresholdOp holds an operation to apply when a pool threshold is hit.
type ThresholdOp uint8

var thresholdOpValues = []*thresholdOpData{
	{
		Key:         "unknown",
		String:      i18n.Text("Unknown"),
		Description: i18n.Text("Unknown"),
	},
	{
		Key:         "halve_move",
		String:      i18n.Text("Halve Move"),
		Description: i18n.Text("Halve Move (round up)"),
	},
	{
		Key:         "halve_dodge",
		String:      i18n.Text("Halve Dodge"),
		Description: i18n.Text("Halve Dodge (round up)"),
	},
	{
		Key:         "halve_st",
		String:      i18n.Text("Halve ST"),
		Description: i18n.Text("Halve ST (round up; does not affect HP and damage)"),
	},
}

// ThresholdOpFromString extracts a ThresholdOp from a key.
func ThresholdOpFromString(key string) ThresholdOp {
	for i, one := range thresholdOpValues {
		if strings.EqualFold(key, one.Key) {
			return ThresholdOp(i)
		}
	}
	return 0
}

// EnsureValid returns the first ThresholdOp if this ThresholdOp is not a known value.
func (t ThresholdOp) EnsureValid() ThresholdOp {
	if int(t) < len(thresholdOpValues) {
		return t
	}
	return 0
}

// Key returns the key used to represent this ThresholdOp.
func (t ThresholdOp) Key() string {
	return thresholdOpValues[t.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (t ThresholdOp) String() string {
	return thresholdOpValues[t.EnsureValid()].String
}

// Description of this ThresholdOp's function.
func (t ThresholdOp) Description() string {
	return thresholdOpValues[t.EnsureValid()].Description
}
