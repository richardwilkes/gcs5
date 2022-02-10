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

package theme

import "github.com/richardwilkes/unison"

// Additional colors over and above what unison provides by default.
var (
	HeaderColor                = &unison.ThemeColor{Light: unison.RGB(43, 43, 43), Dark: unison.RGB(64, 64, 64)}
	OnHeaderColor              = &unison.ThemeColor{Light: unison.RGB(255, 255, 255), Dark: unison.RGB(192, 192, 192)}
	EditableBorderColor        = &unison.ThemeColor{Light: unison.RGB(192, 192, 192), Dark: unison.RGB(96, 96, 96)}
	EditableBorderFocusedColor = &unison.ThemeColor{Light: unison.RGB(0, 0, 192), Dark: unison.RGB(0, 102, 102)}
	AccentColor                = &unison.ThemeColor{Light: unison.RGB(0, 102, 102), Dark: unison.RGB(100, 153, 153)}
	SearchListColor            = &unison.ThemeColor{Light: unison.RGB(224, 255, 255), Dark: unison.RGB(0, 43, 43)}
	OnSearchListColor          = &unison.ThemeColor{Light: unison.RGB(0, 0, 0), Dark: unison.RGB(204, 204, 204)}
	PageColor                  = &unison.ThemeColor{Light: unison.RGB(255, 255, 255), Dark: unison.RGB(16, 16, 16)}
	OnPageColor                = &unison.ThemeColor{Light: unison.RGB(0, 0, 0), Dark: unison.RGB(160, 160, 160)}
	PageVoidColor              = &unison.ThemeColor{Light: unison.RGB(128, 128, 128), Dark: unison.RGB(0, 0, 0)}
	MarkerColor                = &unison.ThemeColor{Light: unison.RGB(252, 242, 196), Dark: unison.RGB(0, 51, 0)}
	OnMarkerColor              = &unison.ThemeColor{Light: unison.RGB(0, 0, 0), Dark: unison.RGB(221, 221, 221)}
	OverloadedColor            = &unison.ThemeColor{Light: unison.RGB(192, 64, 64), Dark: unison.RGB(115, 37, 37)}
	OnOverloadedColor          = &unison.ThemeColor{Light: unison.RGB(255, 255, 255), Dark: unison.RGB(221, 221, 221)}
	HintColor                  = &unison.ThemeColor{Light: unison.RGB(128, 128, 128), Dark: unison.RGB(64, 64, 64)}
)
