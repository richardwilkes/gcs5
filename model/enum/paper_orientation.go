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

package enum

import (
	"strings"

	"github.com/richardwilkes/gcs/model/measure"
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible PaperOrientation values.
const (
	Portrait PaperOrientation = iota
	Landscape
)

// PaperOrientation holds the orientation of the page.
type PaperOrientation uint8

// PaperOrientationFromString extracts a PaperOrientation from a string.
func PaperOrientationFromString(str string) PaperOrientation {
	for one := Portrait; one <= Landscape; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return Portrait
}

// Key returns the key used to represent this ThresholdOp.
func (a PaperOrientation) Key() string {
	switch a {
	case Landscape:
		return "landscape"
	default: // Portrait
		return "portrait"
	}
}

// String implements fmt.Stringer.
func (a PaperOrientation) String() string {
	switch a {
	case Landscape:
		return i18n.Text("Landscape")
	default: // Portrait
		return i18n.Text("Portrait")
	}
}

// Dimensions returns the paper dimensions after orienting the paper.
func (a PaperOrientation) Dimensions(size PaperSize) (width, height measure.Length) {
	width, height = size.Dimensions()
	switch a {
	case Landscape:
		return height, width
	default: // Portrait
		return width, height
	}
}
