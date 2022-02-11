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
	"path/filepath"
	"strings"
	"time"

	"github.com/richardwilkes/gcs/internal/export"
	"github.com/richardwilkes/gcs/internal/ui"
	gsettings "github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/gcs/model/paper"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/i18n"
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
	var textTmplPath, paperSize, topMargin, bottomMargin, leftMargin, rightMargin, orientation string
	var png, webp bool
	marginUsage := i18n.Text("Sets the %s margin to use, rather than the one embedded in the file when exporting from the command line. May be specified in inches (in), centimeters (cm), or millimeters (mm). No suffix will imply inches.")
	cl.NewStringOption(&topMargin).SetName("top").SetUsage(fmt.Sprintf(marginUsage, i18n.Text("top")))
	cl.NewStringOption(&bottomMargin).SetName("bottom").SetUsage(fmt.Sprintf(marginUsage, i18n.Text("bottom")))
	cl.NewStringOption(&leftMargin).SetName("left").SetUsage(fmt.Sprintf(marginUsage, i18n.Text("left")))
	cl.NewStringOption(&rightMargin).SetName("right").SetUsage(fmt.Sprintf(marginUsage, i18n.Text("right")))
	orientations := make([]string, len(paper.AllOrientation))
	for i, one := range paper.AllOrientation {
		orientations[i] = one.String()
	}
	cl.NewStringOption(&orientation).SetName("orientation").SetUsage(i18n.Text("Sets the page orientation. Valid choices are: ") + strings.Join(orientations, ", "))
	sizes := make([]string, len(paper.AllSize))
	for i, one := range paper.AllSize {
		sizes[i] = one.String()
	}
	cl.NewStringOption(&paperSize).SetName("paper").SetSingle('p').SetArg("size").SetUsage(i18n.Text("Sets the paper size to use, rather than the one embedded in the file when exporting from the command line. Valid choices are: ") + strings.Join(sizes, ", "))
	cl.NewStringOption(&textTmplPath).SetName("text").SetSingle('x').SetArg("file").SetUsage(i18n.Text("Export sheets using the specified template file"))
	cl.NewBoolOption(&png).SetName("png").SetUsage(i18n.Text("Export sheets to PNG"))
	cl.NewBoolOption(&webp).SetName("webp").SetUsage(i18n.Text("Export sheets to WebP"))
	fileList := jotrotate.ParseAndSetup(cl)

	exportRequests := 0
	if png {
		exportRequests++
	}
	if webp {
		exportRequests++
	}
	if textTmplPath != "" {
		exportRequests++
	}
	settings.Global() // Here to force early initialization
	switch {
	case exportRequests == 0:
		ui.Start(fileList) // Never returns
	case exportRequests > 1:
		cl.FatalMsg(i18n.Text("Only one of --text, --png, or --webp may be specified."))
	default:
		if len(fileList) == 0 {
			cl.FatalMsg(i18n.Text("No files to process."))
		}
		for _, one := range fileList {
			if strings.ToLower(filepath.Ext(one)) != ".gcs" {
				cl.FatalMsg(i18n.Text("Only .gcs files may be exported."))
			}
		}
		pageOverrides := &gsettings.PageOverrides{}
		pageOverrides.ParseSize(paperSize)
		pageOverrides.ParseOrientation(orientation)
		pageOverrides.ParseTopMargin(topMargin)
		pageOverrides.ParseLeftMargin(leftMargin)
		pageOverrides.ParseBottomMargin(bottomMargin)
		pageOverrides.ParseRightMargin(rightMargin)
		switch {
		case textTmplPath != "":
			export.ToText(textTmplPath, pageOverrides, fileList)
		case webp:
			export.ToWebP(pageOverrides, fileList)
		case png:
			export.ToPNG(pageOverrides, fileList)
		default:
			jot.Fatal(1, "unexpected case")
		}
	}
	atexit.Exit(0)
}
