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

package res

import (
	_ "embed"

	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

// Various images.
var (
	//go:embed "default_portrait.png"
	DefaultPortraitData []byte
	DefaultPortrait     = mustImg(DefaultPortraitData)
)

func mustImg(data []byte) *unison.Image {
	img, err := unison.NewImageFromBytes(DefaultPortraitData, 0.5)
	jot.FatalIfErr(err)
	return img
}
