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

package advantage

import (
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Affects values.
const (
	Total      = Affects("total")
	BaseOnly   = Affects("base_only")
	LevelsOnly = Affects("levels_only")
)

// AllAffects is the complete set of Affects values.
var AllAffects = []Affects{
	Total,
	BaseOnly,
	LevelsOnly,
}

// Affects describes how an AdvantageModifier affects the point cost.
type Affects string

// EnsureValid ensures this is of a known value.
func (a Affects) EnsureValid() Affects {
	for _, one := range AllAffects {
		if one == a {
			return a
		}
	}
	return AllAffects[0]
}

// String implements fmt.Stringer.
func (a Affects) String() string {
	switch a {
	case Total:
		return i18n.Text("to cost")
	case BaseOnly:
		return i18n.Text("to base cost only")
	case LevelsOnly:
		return i18n.Text("to leveled cost only")
	default:
		return Total.String()
	}
}

// ShortTitle returns the short version of the title.
func (a Affects) ShortTitle() string {
	switch a {
	case Total:
		return ""
	case BaseOnly:
		return i18n.Text("(base only)")
	case LevelsOnly:
		return i18n.Text("(levels only)")
	default:
		return Total.String()
	}
}
