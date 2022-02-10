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

package images

import (
	_ "embed"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

var (
	// DefaultPortraitData holds the default portrait image data.
	//go:embed "default_portrait.png"
	DefaultPortraitData []byte
	defaultPortrait     *unison.Image
)

// DefaultPortrait returns the default portrait image.
func DefaultPortrait() *unison.Image {
	if defaultPortrait == nil {
		var err error
		if defaultPortrait, err = unison.NewImageFromBytes(DefaultPortraitData, 0.5); err != nil {
			jot.Fatal(1, errs.NewWithCause("unable to load default portrait data", err))
		}
	}
	return defaultPortrait
}
