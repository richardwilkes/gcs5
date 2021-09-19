/*
 * Copyright ©1998-2021 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package internal

import (
	"github.com/richardwilkes/gcs/internal/menus"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

// NewWorkspace creates the workspace window.
func NewWorkspace() *unison.Window {
	wnd, err := unison.NewWindow("GCS")
	jot.FatalIfErr(err)
	menus.Setup(wnd)
	wnd.SetFrameRect(unison.PrimaryDisplay().Usable)
	wnd.ToFront()
	return wnd
}
