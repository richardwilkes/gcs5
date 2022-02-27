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

var lockSVG *unison.SVG

// LockSVG returns an SVG that holds an icon of a lock.
func LockSVG() *unison.SVG {
	if lockSVG == nil {
		var err error
		lockSVG, err = unison.NewSVG(geom32.NewSize(448, 512), "M80 192v-48C80 64.47 144.5 0 224 0s144 64.47 144 144v48h16c35.3 0 64 28.7 64 64v192c0 35.3-28.7 64-64 64H64c-35.35 0-64-28.7-64-64V256c0-35.3 28.65-64 64-64h16zm64 0h160v-48c0-44.18-35.8-80-80-80s-80 35.82-80 80v48z")
		jot.FatalIfErr(err)
	}
	return lockSVG
}
