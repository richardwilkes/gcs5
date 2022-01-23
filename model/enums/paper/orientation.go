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
	"strings"

	"github.com/richardwilkes/gcs/model/unit/length"
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Orientation values.
const (
	Portrait Orientation = iota
	Landscape
)

// Orientation holds the orientation of the page.
type Orientation uint8

// OrientationFromString extracts a Orientation from a string.
func OrientationFromString(str string) Orientation {
	for one := Portrait; one <= Landscape; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return Portrait
}

// Key returns the key used to represent this ThresholdOp.
func (a Orientation) Key() string {
	switch a {
	case Landscape:
		return "landscape"
	default: // Portrait
		return "portrait"
	}
}

// String implements fmt.Stringer.
func (a Orientation) String() string {
	switch a {
	case Landscape:
		return i18n.Text("Landscape")
	default: // Portrait
		return i18n.Text("Portrait")
	}
}

// Dimensions returns the paper dimensions after orienting the paper.
func (a Orientation) Dimensions(size Size) (width, height length.GURPS) {
	switch a {
	case Landscape:
		return height, width
	default: // Portrait
		return width, height
	}
}
