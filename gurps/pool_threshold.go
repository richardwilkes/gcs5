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

// PoolThreshold holds a point within an attribute pool where changes in state occur.
type PoolThreshold struct {
	State       string        `json:"state"`
	Explanation string        `json:"explanation,omitempty"`
	Multiplier  int           `json:"multiplier,omitempty"`
	Divisor     int           `json:"divisor,omitempty"`
	Addition    int           `json:"addition,omitempty"`
	Ops         []ThresholdOp `json:"ops,omitempty"`
}

// Threshold returns the threshold value for the given maximum.
func (p *PoolThreshold) Threshold(max int) int {
	threshold := max * p.Multiplier
	if p.Divisor > 1 {
		threshold /= p.Divisor
		if max%p.Divisor != 0 {
			threshold++
		}
		threshold--
		if threshold < 0 {
			threshold = 0
		}
	}
	return threshold + p.Addition
}
