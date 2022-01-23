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
	"github.com/richardwilkes/gcs/model/enums/paper"
	"github.com/richardwilkes/gcs/model/unit/length"
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
	Size         paper.Size
	Orientation  paper.Orientation
	TopMargin    length.Length
	LeftMargin   length.Length
	BottomMargin length.Length
	RightMargin  length.Length
}

// FactoryPageSettings returns a new PageSettings with factory defaults.
func FactoryPageSettings() *PageSettings {
	return &PageSettings{
		Size:         paper.Letter,
		Orientation:  paper.Portrait,
		TopMargin:    length.FromFloat64(0.25, length.Inch),
		LeftMargin:   length.FromFloat64(0.25, length.Inch),
		BottomMargin: length.FromFloat64(0.25, length.Inch),
		RightMargin:  length.FromFloat64(0.25, length.Inch),
	}
}

// NewPageSettingsFromJSON creates a new PageSettings from a JSON object.
func NewPageSettingsFromJSON(data map[string]interface{}) *PageSettings {
	s := FactoryPageSettings()
	s.Size = paper.SizeFromString(encoding.String(data[pageSettingsPaperSizeKey]))
	s.Orientation = paper.OrientationFromString(encoding.String(data[pageSettingsOrientationKey]))
	s.TopMargin = length.FromStringForced(encoding.String(data[pageSettingsTopMarginKey]), length.Inch)
	s.LeftMargin = length.FromStringForced(encoding.String(data[pageSettingsLeftMarginKey]), length.Inch)
	s.BottomMargin = length.FromStringForced(encoding.String(data[pageSettingsBottomMarginKey]), length.Inch)
	s.RightMargin = length.FromStringForced(encoding.String(data[pageSettingsRightMarginKey]), length.Inch)
	return s
}

// ToKeyedJSON emits this object as JSON with the specified key.
func (s *PageSettings) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	encoder.Key(key)
	s.ToJSON(encoder)
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