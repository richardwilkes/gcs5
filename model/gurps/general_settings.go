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
	generalSettingsDefaultPlayerNameKey           = "default_player_name"
	generalSettingsDefaultTechLevelKey            = "default_tech_level"
	generalSettingsCalendarRefKey                 = "calendar_ref"
	generalSettingsGCalcKey                       = "gurps_calculator_key"
	generalSettingsInitialPointsKey               = "initial_points"
	generalSettingsTooltipDelayKey                = "tooltip_initial_delay_milliseconds"
	generalSettingsTooltipDismissalKey            = "tooltip_dismiss_delay_seconds"
	generalSettingsImageResolutionKey             = "image_resolution"
	generalSettingsInitialUIScaleKey              = "initial_ui_scale"
	generalSettingsAutoFillProfileKey             = "auto_fill_profile"
	generalSettingsIncludeUnspentPointsInTotalKey = "include_unspent_points_in_total"
)

// GeneralSettings holds general settings for a sheet.
type GeneralSettings struct {
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

// NewGeneralSettingsFromFile loads a new GeneralSettings from a file.
func NewGeneralSettingsFromFile(fsys fs.FS, filePath string) (*GeneralSettings, error) {
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
	return NewGeneralSettingsFromJSON(obj), nil
}

// NewGeneralSettingsFromJSON creates a new GeneralSettings from a JSON object.
func NewGeneralSettingsFromJSON(data map[string]interface{}) *GeneralSettings {
	return &GeneralSettings{
		DefaultPlayerName:           encoding.String(data[generalSettingsDefaultPlayerNameKey]),
		DefaultTechLevel:            encoding.String(data[generalSettingsDefaultTechLevelKey]),
		CalendarRef:                 encoding.String(data[generalSettingsCalendarRefKey]),
		GCalcKey:                    encoding.String(data[generalSettingsGCalcKey]),
		InitialPoints:               int(encoding.Number(data[generalSettingsInitialPointsKey]).AsInt64()),
		ToolTipDelayMillis:          int(encoding.Number(data[generalSettingsTooltipDelayKey]).AsInt64()),
		ToolTipDismissalSeconds:     int(encoding.Number(data[generalSettingsTooltipDismissalKey]).AsInt64()),
		ImageResolution:             int(encoding.Number(data[generalSettingsImageResolutionKey]).AsInt64()),
		InitialUIScale:              int(encoding.Number(data[generalSettingsInitialUIScaleKey]).AsInt64()),
		AutoFillProfile:             encoding.Bool(data[generalSettingsAutoFillProfileKey]),
		IncludeUnspentPointsInTotal: encoding.Bool(data[generalSettingsIncludeUnspentPointsInTotalKey]),
	}
}

// Save writes the GeneralSettings to the file as JSON.
func (s *GeneralSettings) Save(filePath string) error {
	return encoding.SaveJSON(filePath, true, func(encoder *encoding.JSONEncoder) {
		s.ToJSON(encoder)
	})
}

// ToKeyedJSON emits this object as JSON with the specified key.
func (s *GeneralSettings) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	encoder.Key(key)
	s.ToJSON(encoder)
}

// ToJSON emits this object as JSON.
func (s *GeneralSettings) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(generalSettingsDefaultPlayerNameKey, s.DefaultPlayerName, true, true)
	encoder.KeyedString(generalSettingsDefaultTechLevelKey, s.DefaultTechLevel, true, true)
	encoder.KeyedString(generalSettingsCalendarRefKey, s.CalendarRef, true, true)
	encoder.KeyedString(generalSettingsGCalcKey, s.GCalcKey, true, true)
	encoder.KeyedNumber(generalSettingsInitialPointsKey, fixed.F64d4FromInt64(int64(s.InitialPoints)), true)
	encoder.KeyedNumber(generalSettingsTooltipDelayKey, fixed.F64d4FromInt64(int64(s.ToolTipDelayMillis)), true)
	encoder.KeyedNumber(generalSettingsTooltipDismissalKey, fixed.F64d4FromInt64(int64(s.ToolTipDismissalSeconds)), true)
	encoder.KeyedNumber(generalSettingsImageResolutionKey, fixed.F64d4FromInt64(int64(s.ImageResolution)), true)
	encoder.KeyedNumber(generalSettingsInitialUIScaleKey, fixed.F64d4FromInt64(int64(s.InitialUIScale)), true)
	encoder.KeyedBool(generalSettingsAutoFillProfileKey, s.AutoFillProfile, true)
	encoder.KeyedBool(generalSettingsIncludeUnspentPointsInTotalKey, s.IncludeUnspentPointsInTotal, true)
	encoder.EndObject()
}

// UpdateToolTipTiming updates the default tooltip theme to use the timing values from this object.
func (s *GeneralSettings) UpdateToolTipTiming() {
	unison.DefaultTooltipTheme.Delay = time.Duration(s.ToolTipDelayMillis) * time.Millisecond
	unison.DefaultTooltipTheme.Dismissal = time.Duration(s.ToolTipDismissalSeconds) * time.Second
}
