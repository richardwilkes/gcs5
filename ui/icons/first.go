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

var firstSVG *unison.SVG

// FirstSVG returns an SVG that holds an icon for going to the first element.
func FirstSVG() *unison.SVG {
	if firstSVG == nil {
		var err error
		firstSVG, err = unison.NewSVG(geom32.NewSize(512, 512), "M0 415.1V96.03c0-17.67 14.33-31.1 31.1-31.1 18.57-.9 32.9 13.43 32.9 31.1v131.8l171.5-156.5c20.6-17.05 52.5-2.67 52.5 24.7v131.9l171.5-156.5c20.6-17.15 52.5-2.77 52.5 24.6v319.9c0 27.37-31.88 41.74-52.5 24.62L288 285.2v130.7c0 27.37-31.88 41.74-52.5 24.62L64 285.2v130.7c0 17.67-14.33 31.1-31.1 31.1-18.57.1-32.9-13.4-32.9-31.9z")
		jot.FatalIfErr(err)
	}
	return firstSVG
}
