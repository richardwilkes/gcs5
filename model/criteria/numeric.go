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
	"github.com/richardwilkes/gcs/model/gurps/enum"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	typeKey      = "compare"
	qualifierKey = "qualifier"
)

// Numeric holds the criteria for matching a number.
type Numeric struct {
	Type      enum.NumericCompareType
	Qualifier fixed.F64d4
}

// NewNumericFromJSON creates a new Numeric from a JSON object.
func NewNumericFromJSON(data map[string]interface{}) *Numeric {
	n := &Numeric{}
	n.FromJSON(data)
	return n
}

// FromJSON replaces the current data with data from a JSON object.
func (n *Numeric) FromJSON(data map[string]interface{}) {
	n.Type = enum.NumericCompareTypeFromString(encoding.String(data[typeKey]))
	n.Qualifier = encoding.Number(data[qualifierKey])
}

// ToJSON emits the JSON for this object.
func (n *Numeric) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	n.ToInlineJSON(encoder)
	encoder.EndObject()
}

// ToInlineJSON emits the JSON key values that comprise this object without the object wrapper.
func (n *Numeric) ToInlineJSON(encoder *encoding.JSONEncoder) {
	if n.Type != enum.AnyNumber {
		encoder.KeyedString(typeKey, n.Type.Key(), false, false)
		encoder.KeyedNumber(qualifierKey, n.Qualifier, false)
	}
}
