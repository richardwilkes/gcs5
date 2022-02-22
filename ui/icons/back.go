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

var backSVG *unison.SVG

// BackSVG returns an SVG that holds an icon for going back in history one step.
func BackSVG() *unison.SVG {
	if backSVG == nil {
		var err error
		backSVG, err = unison.NewSVG(geom32.NewSize(256, 512), "M137.4 406.6 9.4 279.5C3.125 272.4 0 264.2 0 255.1s3.125-16.38 9.375-22.63l128-127.1c9.156-9.156 22.91-11.9 34.88-6.943S192 115.1 192 128v255.1c0 12.94-7.781 24.62-19.75 29.58s-25.75 3.12-34.85-6.08z")
		jot.FatalIfErr(err)
	}
	return backSVG
}
