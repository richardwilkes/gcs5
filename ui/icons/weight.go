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

var weightSVG *unison.SVG

// WeightSVG returns an SVG that holds an icon for weight.
func WeightSVG() *unison.SVG {
	if weightSVG == nil {
		var err error
		weightSVG, err = unison.NewSVG(geom32.NewSize(512, 512), "m510.3 445.9-73-292.1c-3.8-15.3-16.5-25.8-30.9-25.8h-60.3c3.625-9.1 5.875-20.75 5.875-32 0-53-42.1-96-96-96S159.1 43 159.1 96c0 11.25 2.25 22 5.875 32H105.6c-14.38 0-27.13 10.5-30.88 25.75L1.71 445.85C-6.641 479.1 16.36 512 47.99 512h416c31.61 0 54.61-32.9 46.31-66.1zM256 128c-17.6 0-32.9-14.4-32.9-32s15.3-32 32.9-32c17.63 0 32 14.38 32 32s-14.4 32-32 32z")
		jot.FatalIfErr(err)
	}
	return weightSVG
}
