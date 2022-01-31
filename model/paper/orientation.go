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

package paper

import (
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Orientation values.
const (
	Portrait  = Orientation("portrait")
	Landscape = Orientation("landscape")
)

// AllOrientations is the complete set of Orientation values.
var AllOrientations = []Orientation{
	Portrait,
	Landscape,
}

// Orientation holds the orientation of the page.
type Orientation string

// EnsureValid ensures this is of a known value.
func (o Orientation) EnsureValid() Orientation {
	for _, one := range AllOrientations {
		if one == o {
			return o
		}
	}
	return AllOrientations[0]
}

// String implements fmt.Stringer.
func (o Orientation) String() string {
	switch o {
	case Portrait:
		return i18n.Text("Portrait")
	case Landscape:
		return i18n.Text("Landscape")
	default:
		return Portrait.String()
	}
}

// Dimensions returns the paper dimensions after orienting the paper.
func (o Orientation) Dimensions(size Size) (width, height Length) {
	switch o {
	case Portrait:
		return size.Dimensions()
	case Landscape:
		width, height = size.Dimensions()
		return height, width
	default:
		return Portrait.Dimensions(size)
	}
}
