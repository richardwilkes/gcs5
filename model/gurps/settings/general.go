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

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
	"github.com/richardwilkes/unison"
)

// Default, min & max values for the general numeric settings
var (
	InitialPointsDef       = f64d4.FromInt(150)
	InitialPointsMin       f64d4.Int
	InitialPointsMax       = f64d4.FromInt(9999999)
	TooltipDelayDef        = f64d4.FromStringForced("0.75")
	TooltipDelayMin        f64d4.Int
	TooltipDelayMax        = fxp.Thirty
	TooltipDismissalDef    = f64d4.FromInt(60)
	TooltipDismissalMin    = f64d4.One
	TooltipDismissalMax    = f64d4.FromInt(3600)
	ImageResolutionDef     = 200
	ImageResolutionMin     = 50
	ImageResolutionMax     = 400
	InitialUIScaleMin      = 50
	InitialUIScaleMax      = 400
	InitialListUIScaleDef  = 100
	InitialSheetUIScaleDef = 133
)

// General holds settings for a sheet.
type General struct {
	DefaultPlayerName           string    `json:"default_player_name,omitempty"`
	DefaultTechLevel            string    `json:"default_tech_level,omitempty"`
	CalendarName                string    `json:"calendar_ref,omitempty"`
	GCalcKey                    string    `json:"gurps_calculator_key,omitempty"`
	InitialPoints               f64d4.Int `json:"initial_points"`
	TooltipDelay                f64d4.Int `json:"tooltip_delay"`
	TooltipDismissal            f64d4.Int `json:"tooltip_dismissal"`
	InitialListUIScale          int       `json:"initial_list_scale"`
	InitialSheetUIScale         int       `json:"initial_sheet_scale"`
	ImageResolution             int       `json:"image_resolution"`
	AutoFillProfile             bool      `json:"auto_fill_profile,omitempty"`
	IncludeUnspentPointsInTotal bool      `json:"include_unspent_points_in_total,omitempty"`
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
		InitialPoints:               InitialPointsDef,
		TooltipDelay:                TooltipDelayDef,
		TooltipDismissal:            TooltipDismissalDef,
		InitialListUIScale:          InitialListUIScaleDef,
		InitialSheetUIScale:         InitialSheetUIScaleDef,
		ImageResolution:             ImageResolutionDef,
		AutoFillProfile:             true,
		IncludeUnspentPointsInTotal: true,
	}
}

// NewGeneralFromFile loads new settings from a file.
func NewGeneralFromFile(fileSystem fs.FS, filePath string) (*General, error) {
	var data struct {
		General     `json:",inline"`
		OldLocation *General `json:"general"`
	}
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &data); err != nil {
		return nil, err
	}
	var s *General
	if data.OldLocation != nil {
		s = data.OldLocation
	} else {
		settings := data.General
		s = &settings
	}
	s.EnsureValidity()
	return s, nil
}

// Save writes the settings to the file as JSON.
func (s *General) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, s)
}

// UpdateToolTipTiming updates the default tooltip theme to use the timing values from this object.
func (s *General) UpdateToolTipTiming() {
	unison.DefaultTooltipTheme.Delay = time.Duration(s.TooltipDelay.Mul(fxp.Thousand).AsInt64()) * time.Millisecond
	unison.DefaultTooltipTheme.Dismissal = time.Duration(s.TooltipDismissal.Mul(fxp.Thousand).AsInt64()) * time.Millisecond
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

// EnsureValidity checks the current settings for validity and if they aren't valid, makes them so.
func (s *General) EnsureValidity() {
	s.InitialPoints = fxp.ResetIfOutOfRange(s.InitialPoints, InitialPointsMin, InitialPointsMax, InitialPointsDef)
	s.TooltipDelay = fxp.ResetIfOutOfRange(s.TooltipDelay, TooltipDelayMin, TooltipDelayMax, TooltipDelayDef)
	s.TooltipDismissal = fxp.ResetIfOutOfRange(s.TooltipDismissal, TooltipDismissalMin, TooltipDismissalMax, TooltipDismissalDef)
	s.ImageResolution = fxp.ResetIfOutOfRangeInt(s.ImageResolution, ImageResolutionMin, ImageResolutionMax, ImageResolutionDef)
	s.InitialListUIScale = fxp.ResetIfOutOfRangeInt(s.InitialListUIScale, InitialUIScaleMin, InitialUIScaleMax, InitialListUIScaleDef)
	s.InitialSheetUIScale = fxp.ResetIfOutOfRangeInt(s.InitialSheetUIScale, InitialUIScaleMin, InitialUIScaleMax, InitialSheetUIScaleDef)
}
