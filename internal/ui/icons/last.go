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

package icons

import (
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var lastSVG *unison.SVG

// LastSVG returns an SVG that holds an icon for going to the last element.
func LastSVG() *unison.SVG {
	if lastSVG == nil {
		var err error
		lastSVG, err = unison.NewSVG(geom32.NewSize(512, 512), "M512 96.03v319.9c0 17.67-14.33 31.1-31.1 31.1-18.6.07-32.9-13.43-32.9-31.93v-131L276.5 440.6c-20.6 17.1-52.5 2.7-52.5-25.5v-131L52.5 440.6C31.88 457.7 0 443.3 0 415.1V96.03c0-27.37 31.88-41.74 52.5-24.62L224 226.8V96.03c0-27.37 31.88-41.74 52.5-24.62L448 226.8V96.03c0-17.67 14.33-31.1 31.1-31.1 18.6-.9 32.9 13.43 32.9 31.1z")
		jot.FatalIfErr(err)
	}
	return lastSVG
}
