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

package settings

import (
	"context"
	"io/fs"
	"os"
	"os/user"
	"time"

	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

// General holds settings for a sheet.
type General struct {
	DefaultPlayerName           string `json:"default_player_name,omitempty"`
	DefaultTechLevel            string `json:"default_tech_level,omitempty"`
	CalendarName                string `json:"calendar_ref,omitempty"`
	GCalcKey                    string `json:"gurps_calculator_key,omitempty"`
	InitialPoints               int    `json:"initial_points,omitempty"`
	ToolTipDelayMillis          int    `json:"tooltip_initial_delay_milliseconds,omitempty"`
	ToolTipDismissalSeconds     int    `json:"tooltip_dismiss_delay_seconds,omitempty"`
	ImageResolution             int    `json:"image_resolution,omitempty"`
	InitialUIScale              int    `json:"initial_ui_scale,omitempty"`
	AutoFillProfile             bool   `json:"auto_fill_profile,omitempty"`
	IncludeUnspentPointsInTotal bool   `json:"include_unspent_points_in_total,omitempty"`
}

// NewGeneral creates settings with factory defaults.
func NewGeneral() *General {
	var name string
	if u, err := user.Current(); err != nil {
		name = os.Getenv("USER")
	} else {
		name = u.Name
	}
	return &General{
		DefaultPlayerName:           name,
		DefaultTechLevel:            "3",
		InitialPoints:               250,
		ToolTipDelayMillis:          750,
		ToolTipDismissalSeconds:     60,
		ImageResolution:             200,
		InitialUIScale:              125,
		AutoFillProfile:             true,
		IncludeUnspentPointsInTotal: true,
	}
}

// NewGeneralFromFile loads new settings from a file.
func NewGeneralFromFile(fsys fs.FS, filePath string) (*General, error) {
	var data struct {
		General     `json:",inline"`
		OldLocation *General `json:"general"`
	}
	if err := jio.LoadFromFS(context.Background(), fsys, filePath, &data); err != nil {
		return nil, err
	}
	if data.OldLocation != nil {
		return data.OldLocation, nil
	}
	settings := data.General
	return &settings, nil
}

// Save writes the settings to the file as JSON.
func (s *General) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, s)
}

// UpdateToolTipTiming updates the default tooltip theme to use the timing values from this object.
func (s *General) UpdateToolTipTiming() {
	unison.DefaultTooltipTheme.Delay = time.Duration(s.ToolTipDelayMillis) * time.Millisecond
	unison.DefaultTooltipTheme.Dismissal = time.Duration(s.ToolTipDismissalSeconds) * time.Second
}

// CalendarRef returns the CalendarRef these settings refer to.
func (s *General) CalendarRef(libraries library.Libraries) *CalendarRef {
	ref := LookupCalendarRef(s.CalendarName, libraries)
	if ref == nil {
		if ref = LookupCalendarRef("Gregorian", libraries); ref == nil {
			jot.Fatal(1, "unable to load default calendar (Gregorian)")
		}
	}
	return ref
}
