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

package gurps

import (
	"os"
	"os/user"
)

// GeneralSettings holds general settings for a sheet.
type GeneralSettings struct {
	DefaultPlayerName           string  `json:"default_player_name"`
	DefaultTechLevel            string  `json:"default_tech_level"`
	PDFViewer                   string  `json:"pdf_viewer"`
	InitialPoints               int     `json:"initial_points"`
	TooltipTimeout              int     `json:"tooltip_timeout"`
	ImageResolution             int     `json:"image_resolution"`
	InitialUIScale              float32 `json:"initial_ui_scale"`
	AutoFillProfile             bool    `json:"auto_fill_profile"`
	IncludeUnspentPointsInTotal bool    `json:"include_unspent_points_in_total"`
}

// NewGeneralSettings return new GeneralSettings.
func NewGeneralSettings() *GeneralSettings {
	var name string
	if u, err := user.Current(); err != nil {
		name = os.Getenv("USER")
	} else {
		name = u.Name
	}
	return &GeneralSettings{
		DefaultPlayerName:           name,
		DefaultTechLevel:            "3",
		PDFViewer:                   "", // TODO: get default for platform
		InitialPoints:               250,
		TooltipTimeout:              60,
		ImageResolution:             200,
		InitialUIScale:              1.25,
		AutoFillProfile:             true,
		IncludeUnspentPointsInTotal: true,
	}
}
