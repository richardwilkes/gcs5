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
	"path/filepath"

	"github.com/richardwilkes/gcs/internal/library"
	"github.com/richardwilkes/gcs/internal/settings"
	"github.com/richardwilkes/gcs/internal/ui/menus"
	"github.com/richardwilkes/gcs/internal/ui/trampolines"
	"github.com/richardwilkes/gcs/internal/ui/workspace"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

// Start the UI.
func Start(files []string) {
	go library.PerformUpdateChecks(settings.Global.Libraries)
	unison.Start(
		unison.StartupFinishedCallback(func() {
			trampolines.MenuSetup = menus.Setup
			wnd, err := unison.NewWindow("GCS")
			jot.FatalIfErr(err)
			menus.Setup(wnd)
			workspace.NewWorkspace(wnd)
			wnd.SetFrameRect(unison.PrimaryDisplay().Usable)
			wnd.ToFront()
			openFiles(files)
		}),
		unison.OpenFilesCallback(openFiles),
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

func openFiles(urls []string) {
	for _, wnd := range unison.Windows() {
		if ws := workspace.FromWindow(wnd); ws != nil {
			for _, one := range urls {
				if p, err := filepath.Abs(one); err != nil {
					unison.ErrorDialogWithError(i18n.Text("Unable to open ")+one, err)
				} else {
					workspace.OpenFile(wnd, p)
				}
			}
		}
	}
}
