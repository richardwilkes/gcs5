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

var previousSVG *unison.SVG

// PreviousSVG returns an SVG that holds an icon for going to the previous element.
func PreviousSVG() *unison.SVG {
	if previousSVG == nil {
		var err error
		previousSVG, err = unison.NewSVG(geom32.NewSize(320, 512), "M31.1 64.03c-17.67 0-31.1 14.33-31.1 32v319.9c0 17.67 14.33 32 32 32 17.67-.83 32-14.33 32-32.83V96.03c0-17.67-14.33-32-32.9-32zm236.4 7.38-192 159.1C67.82 237.8 64 246.9 64 256c0 9.094 3.82 18.18 11.44 24.62l192 159.1c20.63 17.12 52.51 2.75 52.51-24.62V95.2c-.85-26.54-31.85-40.92-52.45-23.79z")
		jot.FatalIfErr(err)
	}
	return previousSVG
}
