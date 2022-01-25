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
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/enum"
	"github.com/richardwilkes/gcs/model/measure"
)

const (
	pageSettingsPaperSizeKey    = "paper_size"
	pageSettingsOrientationKey  = "orientation"
	pageSettingsTopMarginKey    = "top_margin"
	pageSettingsLeftMarginKey   = "left_margin"
	pageSettingsBottomMarginKey = "bottom_margin"
	pageSettingsRightMarginKey  = "right_margin"
)

// PageSettings holds page settings.
type PageSettings struct {
	Size         enum.PaperSize
	Orientation  enum.PaperOrientation
	TopMargin    measure.Length
	LeftMargin   measure.Length
	BottomMargin measure.Length
	RightMargin  measure.Length
}

// FactoryPageSettings returns a new PageSettings with factory defaults.
func FactoryPageSettings() *PageSettings {
	return &PageSettings{
		Size:         enum.Letter,
		Orientation:  enum.Portrait,
		TopMargin:    measure.Length{Length: 0.25, Units: measure.Inch},
		LeftMargin:   measure.Length{Length: 0.25, Units: measure.Inch},
		BottomMargin: measure.Length{Length: 0.25, Units: measure.Inch},
		RightMargin:  measure.Length{Length: 0.25, Units: measure.Inch},
	}
}

// NewPageSettingsFromJSON creates a new PageSettings from a JSON object.
func NewPageSettingsFromJSON(data map[string]interface{}) *PageSettings {
	s := FactoryPageSettings()
	s.Size = enum.PaperSizeFromString(encoding.String(data[pageSettingsPaperSizeKey]))
	s.Orientation = enum.PaperOrientationFromString(encoding.String(data[pageSettingsOrientationKey]))
	s.TopMargin = measure.LengthFromString(encoding.String(data[pageSettingsTopMarginKey]))
	s.LeftMargin = measure.LengthFromString(encoding.String(data[pageSettingsLeftMarginKey]))
	s.BottomMargin = measure.LengthFromString(encoding.String(data[pageSettingsBottomMarginKey]))
	s.RightMargin = measure.LengthFromString(encoding.String(data[pageSettingsRightMarginKey]))
	return s
}

// ToJSON emits this object as JSON.
func (s *PageSettings) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(pageSettingsPaperSizeKey, s.Size.Key(), false, false)
	encoder.KeyedString(pageSettingsOrientationKey, s.Orientation.Key(), false, false)
	encoder.KeyedString(pageSettingsTopMarginKey, s.TopMargin.String(), false, false)
	encoder.KeyedString(pageSettingsLeftMarginKey, s.LeftMargin.String(), false, false)
	encoder.KeyedString(pageSettingsBottomMarginKey, s.BottomMargin.String(), false, false)
	encoder.KeyedString(pageSettingsRightMarginKey, s.RightMargin.String(), false, false)
	encoder.EndObject()
}
