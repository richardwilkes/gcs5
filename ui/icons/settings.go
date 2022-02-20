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

var settingsSVG *unison.SVG

// SettingsSVG returns an SVG that holds an icon for settings.
func SettingsSVG() *unison.SVG {
	if settingsSVG == nil {
		var err error
		settingsSVG, err = unison.NewSVG(geom32.NewSize(512, 512), "M0 416c0-17.7 14.33-32 32-32h54.66C99 355.7 127.2 336 160 336c32.8 0 60.1 19.7 73.3 48H480c17.7 0 32 14.3 32 32s-14.3 32-32 32H233.3c-13.2 28.3-40.5 48-73.3 48s-61-19.7-73.34-48H32c-17.67 0-32-14.3-32-32zm192 0c0-17.7-14.3-32-32-32s-32 14.3-32 32 14.3 32 32 32 32-14.3 32-32zm160-240c32.8 0 60.1 19.7 73.3 48H480c17.7 0 32 14.3 32 32s-14.3 32-32 32h-54.7c-13.2 28.3-40.5 48-73.3 48s-61-19.7-73.3-48H32c-17.67 0-32-14.3-32-32s14.33-32 32-32h246.7c12.3-28.3 40.5-48 73.3-48zm32 80c0-17.7-14.3-32-32-32s-32 14.3-32 32 14.3 32 32 32 32-14.3 32-32zm96-192c17.7 0 32 14.33 32 32 0 17.7-14.3 32-32 32H265.3c-13.2 28.3-40.5 48-73.3 48s-61-19.7-73.3-48H32c-17.67 0-32-14.3-32-32 0-17.67 14.33-32 32-32h86.7C131 35.75 159.2 16 192 16s60.1 19.75 73.3 48H480zM160 96c0 17.7 14.3 32 32 32s32-14.3 32-32c0-17.67-14.3-32-32-32s-32 14.33-32 32z")
		jot.FatalIfErr(err)
	}
	return settingsSVG
}
