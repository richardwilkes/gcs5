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

var forwardSVG *unison.SVG

// ForwardSVG returns an SVG that holds an icon for going forward in history one step.
func ForwardSVG() *unison.SVG {
	if forwardSVG == nil {
		var err error
		forwardSVG, err = unison.NewSVG(geom32.NewSize(256, 512), "m118.6 105.4 128 127.1c6.3 7.1 9.4 15.3 9.4 22.6s-3.125 16.38-9.375 22.63l-128 127.1c-9.156 9.156-22.91 11.9-34.88 6.943S64 396.9 64 383.1V128c0-12.94 7.781-24.62 19.75-29.58s25.75-2.19 34.85 6.98z")
		jot.FatalIfErr(err)
	}
	return forwardSVG
}
