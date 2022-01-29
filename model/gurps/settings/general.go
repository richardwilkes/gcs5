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
	"io/fs"
	"os"
	"os/user"
	"time"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/unison"
)

const (
	generalDefaultPlayerNameKey           = "default_player_name"
	generalDefaultTechLevelKey            = "default_tech_level"
	generalCalendarRefKey                 = "calendar_ref"
	generalGCalcKey                       = "gurps_calculator_key"
	generalInitialPointsKey               = "initial_points"
	generalTooltipDelayKey                = "tooltip_initial_delay_milliseconds"
	generalTooltipDismissalKey            = "tooltip_dismiss_delay_seconds"
	generalImageResolutionKey             = "image_resolution"
	generalInitialUIScaleKey              = "initial_ui_scale"
	generalAutoFillProfileKey             = "auto_fill_profile"
	generalIncludeUnspentPointsInTotalKey = "include_unspent_points_in_total"
)

// General holds settings for a sheet.
type General struct {
	DefaultPlayerName           string
	DefaultTechLevel            string
	CalendarRef                 string
	GCalcKey                    string
	InitialPoints               int
	ToolTipDelayMillis          int
	ToolTipDismissalSeconds     int
	ImageResolution             int
	InitialUIScale              int
	AutoFillProfile             bool
	IncludeUnspentPointsInTotal bool
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
		CalendarRef:                 "",
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
	data, err := encoding.LoadJSONFromFS(fsys, filePath)
	if err != nil {
		return nil, err
	}
	obj := encoding.Object(data)
	// Check for older formats
	var exists bool
	if data, exists = obj["general"]; exists {
		obj = encoding.Object(data)
	}
	if obj == nil {
		return nil, errs.New("invalid general settings file: " + filePath)
	}
	return NewGeneralFromJSON(obj), nil
}

// NewGeneralFromJSON creates new settings from a JSON object.
func NewGeneralFromJSON(data map[string]interface{}) *General {
	return &General{
		DefaultPlayerName:           encoding.String(data[generalDefaultPlayerNameKey]),
		DefaultTechLevel:            encoding.String(data[generalDefaultTechLevelKey]),
		CalendarRef:                 encoding.String(data[generalCalendarRefKey]),
		GCalcKey:                    encoding.String(data[generalGCalcKey]),
		InitialPoints:               int(encoding.Number(data[generalInitialPointsKey]).AsInt64()),
		ToolTipDelayMillis:          int(encoding.Number(data[generalTooltipDelayKey]).AsInt64()),
		ToolTipDismissalSeconds:     int(encoding.Number(data[generalTooltipDismissalKey]).AsInt64()),
		ImageResolution:             int(encoding.Number(data[generalImageResolutionKey]).AsInt64()),
		InitialUIScale:              int(encoding.Number(data[generalInitialUIScaleKey]).AsInt64()),
		AutoFillProfile:             encoding.Bool(data[generalAutoFillProfileKey]),
		IncludeUnspentPointsInTotal: encoding.Bool(data[generalIncludeUnspentPointsInTotalKey]),
	}
}

// Save writes the settings to the file as JSON.
func (s *General) Save(filePath string) error {
	return encoding.SaveJSON(filePath, true, func(encoder *encoding.JSONEncoder) {
		s.ToJSON(encoder)
	})
}

// ToJSON emits this object as JSON.
func (s *General) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(generalDefaultPlayerNameKey, s.DefaultPlayerName, true, true)
	encoder.KeyedString(generalDefaultTechLevelKey, s.DefaultTechLevel, true, true)
	encoder.KeyedString(generalCalendarRefKey, s.CalendarRef, true, true)
	encoder.KeyedString(generalGCalcKey, s.GCalcKey, true, true)
	encoder.KeyedNumber(generalInitialPointsKey, fixed.F64d4FromInt64(int64(s.InitialPoints)), true)
	encoder.KeyedNumber(generalTooltipDelayKey, fixed.F64d4FromInt64(int64(s.ToolTipDelayMillis)), true)
	encoder.KeyedNumber(generalTooltipDismissalKey, fixed.F64d4FromInt64(int64(s.ToolTipDismissalSeconds)), true)
	encoder.KeyedNumber(generalImageResolutionKey, fixed.F64d4FromInt64(int64(s.ImageResolution)), true)
	encoder.KeyedNumber(generalInitialUIScaleKey, fixed.F64d4FromInt64(int64(s.InitialUIScale)), true)
	encoder.KeyedBool(generalAutoFillProfileKey, s.AutoFillProfile, true)
	encoder.KeyedBool(generalIncludeUnspentPointsInTotalKey, s.IncludeUnspentPointsInTotal, true)
	encoder.EndObject()
}

// UpdateToolTipTiming updates the default tooltip theme to use the timing values from this object.
func (s *General) UpdateToolTipTiming() {
	unison.DefaultTooltipTheme.Delay = time.Duration(s.ToolTipDelayMillis) * time.Millisecond
	unison.DefaultTooltipTheme.Dismissal = time.Duration(s.ToolTipDismissalSeconds) * time.Second
}
