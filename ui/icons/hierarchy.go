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

var hierarchySVG *unison.SVG

// HierarchySVG returns an SVG that holds an icon toggling hierarchy.
func HierarchySVG() *unison.SVG {
	if hierarchySVG == nil {
		var err error
		hierarchySVG, err = unison.NewSVG(geom32.NewSize(576, 512), "M208 80c0-26.51 21.5-48 48-48h64c26.5 0 48 21.49 48 48v64c0 26.5-21.5 48-48 48h-8v40h152c30.9 0 56 25.1 56 56v32h8c26.5 0 48 21.5 48 48v64c0 26.5-21.5 48-48 48h-64c-26.5 0-48-21.5-48-48v-64c0-26.5 21.5-48 48-48h8v-32c0-4.4-3.6-8-8-8H312v40h8c26.5 0 48 21.5 48 48v64c0 26.5-21.5 48-48 48h-64c-26.5 0-48-21.5-48-48v-64c0-26.5 21.5-48 48-48h8v-40H112c-4.4 0-8 3.6-8 8v32h8c26.5 0 48 21.5 48 48v64c0 26.5-21.5 48-48 48H48c-26.51 0-48-21.5-48-48v-64c0-26.5 21.49-48 48-48h8v-32c0-30.9 25.07-56 56-56h152v-40h-8c-26.5 0-48-21.5-48-48V80z")
		jot.FatalIfErr(err)
	}
	return hierarchySVG
}
