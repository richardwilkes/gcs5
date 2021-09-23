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

package ui

import (
	"fmt"

	"github.com/richardwilkes/gcs/internal/settings"
	"github.com/richardwilkes/gcs/internal/ui/menus"
	"github.com/richardwilkes/gcs/internal/ui/trampolines"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

// Start the UI.
func Start(s *settings.Settings, files []string) {
	unison.Start(
		unison.StartupFinishedCallback(func() {
			trampolines.MenuSetup = menus.Setup
			wnd, err := unison.NewWindow("GCS")
			jot.FatalIfErr(err)
			menus.Setup(wnd)
			wnd.SetFrameRect(unison.PrimaryDisplay().Usable)
			wnd.ToFront()
			openURLs(files)
		}),
		unison.OpenURLsCallback(openURLs),
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

func openURLs(urls []string) {
	for _, one := range urls {
		fmt.Println(one)
	}
}
