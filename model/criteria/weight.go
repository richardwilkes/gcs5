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

package criteria

import (
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/measure"
)

// Weight holds the criteria for matching a number.
type Weight struct {
	Type      NumericCompareType
	Qualifier measure.Weight
}

// NewWeightFromJSON creates a new Weight from a JSON object.
func NewWeightFromJSON(data map[string]interface{}, defUnits measure.WeightUnits) *Weight {
	n := &Weight{}
	n.FromJSON(data, defUnits)
	return n
}

// FromJSON replaces the current data with data from a JSON object.
func (n *Weight) FromJSON(data map[string]interface{}, defUnits measure.WeightUnits) {
	n.Type = NumericCompareTypeFromString(encoding.String(data[typeKey]))
	n.Qualifier = measure.WeightFromStringForced(encoding.String(data[qualifierKey]), defUnits)
}

// ToJSON emits the JSON for this object.
func (n *Weight) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	n.ToInlineJSON(encoder)
	encoder.EndObject()
}

// ToInlineJSON emits the JSON key values that comprise this object without the object wrapper.
func (n *Weight) ToInlineJSON(encoder *encoding.JSONEncoder) {
	if n.Type != AnyNumber {
		encoder.KeyedString(typeKey, n.Type.Key(), false, false)
		encoder.KeyedString(qualifierKey, n.Qualifier.String(), false, false)
	}
}
