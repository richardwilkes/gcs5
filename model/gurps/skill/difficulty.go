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

package skill

import (
	"strings"
)

// Possible Difficulty values.
const (
	E  = Difficulty("e")
	A  = Difficulty("a")
	H  = Difficulty("h")
	VH = Difficulty("vh")
	W  = Difficulty("w")
)

// AllDifficulties is the complete set of Difficulty values.
var AllDifficulties = []Difficulty{E, A, H, VH, W}

// Difficulty holds the difficulty level of a skill.
type Difficulty string

// EnsureValid ensures this is of a known value.
func (d Difficulty) EnsureValid() Difficulty {
	for _, one := range AllDifficulties {
		if one == d {
			return d
		}
	}
	return AllDifficulties[0]
}

// String implements fmt.Stringer.
func (d Difficulty) String() string {
	return strings.ToUpper(string(d.EnsureValid()))
}

// BaseRelativeLevel returns the base relative skill level at 0 points.
func (d Difficulty) BaseRelativeLevel() int {
	switch d {
	case E:
		return 0
	case A:
		return -1
	case H:
		return -2
	case VH, W:
		return -3
	default:
		return E.BaseRelativeLevel()
	}
}
