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
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	navigatorSettingsDividerPositionKey = "divider_position"
	navigatorSettingsOpenRowKeysKey     = "open_row_keys"
)

// NavigatorSettings holds settings for the navigator view.
type NavigatorSettings struct {
	DividerPosition int
	OpenRowKeys     []string
}

func NewNavigatorSettings() *NavigatorSettings {
	return &NavigatorSettings{DividerPosition: 300}
}

// NewNavigatorSettingsFromJSON creates a new NavigatorSettings from a JSON object.
func NewNavigatorSettingsFromJSON(data map[string]interface{}) *NavigatorSettings {
	s := &NavigatorSettings{
		DividerPosition: int(encoding.Number(data[navigatorSettingsDividerPositionKey]).AsFloat64()),
	}
	array := encoding.Array(data[navigatorSettingsOpenRowKeysKey])
	if len(array) != 0 {
		s.OpenRowKeys = make([]string, 0, len(array))
		for _, k := range array {
			s.OpenRowKeys = append(s.OpenRowKeys, encoding.String(k))
		}
	}
	return s
}

// ToKeyedJSON emits this object as JSON with the specified key.
func (s *NavigatorSettings) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	encoder.Key(key)
	s.ToJSON(encoder)
}

// ToJSON emits this object as JSON.
func (s *NavigatorSettings) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedNumber(navigatorSettingsDividerPositionKey, fixed.F64d4FromInt64(int64(s.DividerPosition)), false)
	if len(s.OpenRowKeys) != 0 {
		encoder.Key(navigatorSettingsOpenRowKeysKey)
		encoder.StartArray()
		for _, key := range s.OpenRowKeys {
			encoder.String(key)
		}
		encoder.EndArray()
	}
	encoder.EndObject()
}
