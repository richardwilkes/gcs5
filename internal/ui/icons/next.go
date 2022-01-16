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

var nextSVG *unison.SVG

// NextSVG returns an SVG that holds an icon for going to the next element.
func NextSVG() *unison.SVG {
	if nextSVG == nil {
		var err error
		nextSVG, err = unison.NewSVG(geom32.NewSize(320, 512), "M287.1 447.1c17.67 0 31.1-14.33 31.1-32V96.03c0-17.67-14.33-32-32-32s-31.1 14.33-31.1 31.1v319.9c0 18.57 15.2 32.07 32 32.07zm-234.59-6.5 192-159.1c7.625-6.436 11.43-15.53 11.43-24.62 0-9.094-3.809-18.18-11.43-24.62l-192-159.1C31.88 54.28 0 68.66 0 96.03v319.9c0 27.37 31.88 41.77 52.51 24.67z")
		jot.FatalIfErr(err)
	}
	return nextSVG
}
