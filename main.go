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

package main

import (
	"fmt"

	"github.com/richardwilkes/gcs/model/export"
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/setup"
	"github.com/richardwilkes/gcs/ui"
	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jotrotate"
	"github.com/richardwilkes/unison"
)

var dev = "1"

func main() {
	cmdline.AppName = "GCS"
	cmdline.AppCmdName = "gcs"
	cmdline.License = "Mozilla Public License, version 2.0"
	cmdline.CopyrightStartYear = "1998"
	cmdline.CopyrightHolder = "Richard A. Wilkes"
	cmdline.AppIdentifier = "com.trollworks.gcs"
	if dev == "1" {
		cmdline.AppVersion = "0.0"
	} else {
		cmdline.AppVersion = "5.0.0"
	}
	unison.AttachConsole()
	cl := cmdline.New(true)
	var textTmplPath string
	var showCopyrightDateAndExit bool
	cl.NewGeneralOption(&textTmplPath).SetName("text").SetSingle('x').SetArg("file").
		SetUsage(i18n.Text("Export sheets using the specified template file"))
	cl.NewGeneralOption(&showCopyrightDateAndExit).SetName("copyright-date")
	fileList := jotrotate.ParseAndSetup(cl)
	if showCopyrightDateAndExit {
		fmt.Print(cmdline.ResolveCopyrightYears())
		atexit.Exit(0)
	}
	setup.Setup()
	settings.Global() // Here to force early initialization
	if textTmplPath != "" {
		if len(fileList) == 0 {
			cl.FatalMsg(i18n.Text("No files to process."))
		}
		for _, one := range fileList {
			if !library.FileInfoFor(one).IsExportable {
				cl.FatalMsg(one + i18n.Text(" is not exportable."))
			}
		}
		if err := export.ToText(textTmplPath, fileList); err != nil {
			cl.FatalMsg(err.Error())
		}
	} else {
		ui.Start(fileList) // Never returns
	}
	atexit.Exit(0)
}
