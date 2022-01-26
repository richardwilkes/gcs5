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

package ui

import (
	"github.com/richardwilkes/gcs/internal/ui/menus"
	"github.com/richardwilkes/gcs/internal/ui/trampolines"
	"github.com/richardwilkes/gcs/internal/ui/workspace"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

// Start the UI.
func Start(files []string) {
	libs := settings.Global().Libraries
	go libs.PerformUpdateChecks()
	unison.Start(
		unison.StartupFinishedCallback(func() {
			trampolines.MenuSetup = menus.Setup
			wnd, err := unison.NewWindow("GCS")
			jot.FatalIfErr(err)
			menus.Setup(wnd)
			workspace.NewWorkspace(wnd)
			wnd.SetFrameRect(unison.PrimaryDisplay().Usable)
			wnd.ToFront()
			workspace.OpenFiles(files)
		}),
		unison.OpenFilesCallback(workspace.OpenFiles),
		unison.AllowQuitCallback(func() bool {
			for _, wnd := range unison.Windows() {
				wnd.AttemptClose()
				if wnd.IsValid() {
					return false
				}
			}
			return true
		}),
	) // Never returns
}
