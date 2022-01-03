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

package navigator

import (
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

func createNodeCell(ext, title string) *unison.Panel {
	size := unison.LabelFont.ResolvedFont().Size() + 5
	info, ok := fileTypes[ext]
	if !ok {
		info, ok = fileTypes["file"]
	}
	label := unison.NewLabel()
	label.Text = title
	label.Drawable = &unison.DrawableSVG{
		SVG:  info.svg,
		Size: geom32.NewSize(size, size),
	}
	return label.AsPanel()
}
