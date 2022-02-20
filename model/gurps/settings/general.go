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
	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/unison"
)

// Default, min & max values for the general numeric settings
var (
	InitialPointsDef    = fixed.F64d4FromInt(150)
	InitialPointsMin    fixed.F64d4
	InitialPointsMax    = fixed.F64d4FromInt(9999999)
	TooltipDelayDef     = fixed.F64d4FromStringForced("0.75")
	TooltipDelayMin     fixed.F64d4
	TooltipDelayMax     = fxp.Thirty
	TooltipDismissalDef = fixed.F64d4FromInt(60)
	TooltipDismissalMin = fixed.F64d4One
	TooltipDismissalMax = fixed.F64d4FromInt(3600)
	ImageResolutionDef  = 200
	ImageResolutionMin  = 50
	ImageResolutionMax  = 400
	InitialUIScaleDef   = fixed.F64d4FromInt(125)
	InitialUIScaleMin   = fxp.Ten
	InitialUIScaleMax   = fixed.F64d4FromInt(999)
)

// General holds settings for a sheet.
type General struct {
	DefaultPlayerName           string      `json:"default_player_name,omitempty"`
	DefaultTechLevel            string      `json:"default_tech_level,omitempty"`
	CalendarName                string      `json:"calendar_ref,omitempty"`
	GCalcKey                    string      `json:"gurps_calculator_key,omitempty"`
	InitialPoints               fixed.F64d4 `json:"initial_points,omitempty"`
	TooltipDelay                fixed.F64d4 `json:"tooltip_delay,omitempty"`
	TooltipDismissal            fixed.F64d4 `json:"tooltip_dismissal,omitempty"`
	InitialUIScale              fixed.F64d4 `json:"initial_scale,omitempty"`
	ImageResolution             int         `json:"image_resolution,omitempty"`
	AutoFillProfile             bool        `json:"auto_fill_profile,omitempty"`
	IncludeUnspentPointsInTotal bool        `json:"include_unspent_points_in_total,omitempty"`
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
		InitialUIScale:              InitialUIScaleDef,
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
	s.InitialPoints = fxp.ResetIfOutOfRange(s.InitialPoints, InitialPointsMin, InitialPointsMax, InitialPointsDef)
	s.TooltipDelay = fxp.ResetIfOutOfRange(s.TooltipDelay, TooltipDelayMin, TooltipDelayMax, TooltipDelayDef)
	s.TooltipDismissal = fxp.ResetIfOutOfRange(s.TooltipDismissal, TooltipDismissalMin, TooltipDismissalMax, TooltipDismissalDef)
	s.ImageResolution = fxp.ResetIfOutOfRangeInt(s.ImageResolution, ImageResolutionMin, ImageResolutionMax, ImageResolutionDef)
	s.InitialUIScale = fxp.ResetIfOutOfRange(s.InitialUIScale, InitialUIScaleMin, InitialUIScaleMax, InitialUIScaleDef)
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
