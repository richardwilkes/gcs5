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

var bookmarkSVG *unison.SVG

// BookmarkSVG returns an SVG that holds an icon for a bookmark / page reference.
func BookmarkSVG() *unison.SVG {
	if bookmarkSVG == nil {
		var err error
		bookmarkSVG, err = unison.NewSVG(geom32.NewSize(384, 512), "M384 48v464L192 400 0 512V48C0 21.5 21.5 0 48 0h288c26.5 0 48 21.5 48 48z")
		jot.FatalIfErr(err)
	}
	return bookmarkSVG
}
