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
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// NewIconButton creates a standard icon button.
func NewIconButton(svg *unison.SVG, size float32) *unison.Button {
	b := unison.NewButton()
	b.Drawable = &unison.DrawableSVG{
		SVG:  svg,
		Size: geom32.NewSize(size, size),
	}
	return b
}