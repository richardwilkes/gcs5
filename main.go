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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/richardwilkes/gcs/internal/ui"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/log/jotrotate"
)

func main() {
	cmdline.AppName = "GCS"
	cmdline.AppCmdName = "gcs"
	cmdline.License = "Mozilla Public License, version 2.0"
	cmdline.CopyrightYears = fmt.Sprintf("1998-%d", time.Now().Year())
	cmdline.CopyrightHolder = "Richard A. Wilkes"
	cmdline.AppIdentifier = "com.trollworks.gcs"
	if cmdline.AppVersion == "" {
		cmdline.AppVersion = "0.0"
	}
	cl := cmdline.New(true)
	fileList := jotrotate.ParseAndSetup(cl)

	settings.Global() // Here to force early initialization

	name := "basic.adq"
	info, err := gurps.NewAdvantagesFromFile(os.DirFS("samples"), name)
	jot.FatalIfErr(err)
	jot.FatalIfErr(gurps.SaveAdvantages(info, "samples_converted/"+name))

	var m map[string]interface{}
	jot.FatalIfErr(jio.LoadFromFile(context.Background(), "samples/"+name, &m))
	jot.FatalIfErr(jio.SaveToFile(context.Background(), "samples/sorted-"+name, m))
	m = make(map[string]interface{})
	jot.FatalIfErr(jio.LoadFromFile(context.Background(), "samples_converted/"+name, &m))
	jot.FatalIfErr(jio.SaveToFile(context.Background(), "samples_converted/sorted-"+name, m))

	atexit.Exit(0)

	ui.Start(fileList) // Never returns
}
