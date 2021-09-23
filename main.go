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

package main

import (
	"fmt"
	"time"

	"github.com/richardwilkes/gcs/internal/settings"
	"github.com/richardwilkes/gcs/internal/ui"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/log/jotrotate"
)

func main() {
	cmdline.AppName = "GCS"
	cmdline.AppCmdName = "gcs"
	cmdline.License = "Mozilla Public License, version 2.0"
	cmdline.CopyrightYears = fmt.Sprintf("1998-%d", time.Now().Year())
	cmdline.CopyrightHolder = "Richard A. Wilkes"
	cmdline.AppIdentifier = "com.trollworks.gcs"
	cl := cmdline.New(true)
	fileList := jotrotate.ParseAndSetup(cl)
	s := settings.New()
	ui.Start(s, fileList) // Never returns
}
