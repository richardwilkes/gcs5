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
	"github.com/richardwilkes/gcs/model/gurps/enum"
)

const (
	stringCriteriaTypeKey      = "compare"
	stringCriteriaQualifierKey = "qualifier"
)

// StringCriteria holds the criteria for matching a string.
type StringCriteria struct {
	Type      enum.StringCompareType
	Qualifier string
}

// NewStringCriteriaFromJSON creates a new StringCriteria from a JSON object.
func NewStringCriteriaFromJSON(data map[string]interface{}) *StringCriteria {
	s := &StringCriteria{}
	s.FromJSON(data)
	return s
}

// FromJSON replaces the current data with data from a JSON object.
func (s *StringCriteria) FromJSON(data map[string]interface{}) {
	s.Type = enum.StringCompareTypeFromString(encoding.String(data[stringCriteriaTypeKey]))
	s.Qualifier = encoding.String(data[stringCriteriaQualifierKey])
}

// ToJSON emits the JSON for this object.
func (s *StringCriteria) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	s.ToInlineJSON(encoder)
	encoder.EndObject()
}

// ToInlineJSON emits the JSON key values that comprise this object without the object wrapper.
func (s *StringCriteria) ToInlineJSON(encoder *encoding.JSONEncoder) {
	if s.Type != enum.Any {
		encoder.KeyedString(stringCriteriaTypeKey, s.Type.Key(), false, false)
		encoder.KeyedString(stringCriteriaQualifierKey, s.Qualifier, true, true)
	}
}
