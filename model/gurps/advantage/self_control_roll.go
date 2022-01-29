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

package advantage

import (
	"strconv"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Possible SelfControlRoll values.
const (
	None SelfControlRoll = iota
	CR6
	CR9
	CR12
	CR15
)

type selfControlRollData struct {
	String     string
	Multiplier fixed.F64d4
	Value      int
}

// SelfControlRoll holds the information about a self-control roll, from B121.
type SelfControlRoll uint8

var selfControlRollValues = []*selfControlRollData{
	{
		String:     i18n.Text("None Required"),
		Multiplier: f64d4.One,
	},
	{
		String:     i18n.Text("CR: 6 (Resist rarely)"),
		Multiplier: fixed.F64d4FromInt64(2),
		Value:      6,
	},
	{
		String:     i18n.Text("CR: 9 (Resist fairly often)"),
		Multiplier: fixed.F64d4FromStringForced("1.5"),
		Value:      9,
	},
	{
		String:     i18n.Text("CR: 12 (Resist quite often)"),
		Multiplier: f64d4.One,
		Value:      12,
	},
	{
		String:     i18n.Text("CR: 15 (Resist almost all the time)"),
		Multiplier: fixed.F64d4FromStringForced("0.5"),
		Value:      15,
	},
}

// SelfControlRollFromJSON extracts a SelfControlRoll from a key and JSON data object.
func SelfControlRollFromJSON(key string, data map[string]interface{}) SelfControlRoll {
	if v, exists := data[key]; exists {
		if n, err := strconv.ParseInt(encoding.String(v), 10, 64); err == nil {
			num := int(n)
			for i, one := range selfControlRollValues[1:] {
				if num == one.Value {
					return SelfControlRoll(i)
				}
			}
		}
	}
	return None
}

// ToKeyedJSON writes the SelfControlRoll to JSON.
func (s SelfControlRoll) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	if resolved := s.EnsureValid(); resolved != None {
		encoder.KeyedNumber(key, fixed.F64d4FromInt64(int64(selfControlRollValues[resolved].Value)), false)
	}
}

// EnsureValid returns the first SelfControlRoll if this SelfControlRoll is not a known value.
func (s SelfControlRoll) EnsureValid() SelfControlRoll {
	if int(s) < len(selfControlRollValues) {
		return s
	}
	return None
}

// String implements fmt.Stringer.
func (s SelfControlRoll) String() string {
	return selfControlRollValues[s.EnsureValid()].String
}

// DescriptionWithCost returns a formatted description that includes the cost multiplier.
func (s SelfControlRoll) DescriptionWithCost() string {
	resolved := s.EnsureValid()
	if resolved == None {
		return ""
	}
	cr := selfControlRollValues[resolved]
	return cr.String + ", x" + cr.Multiplier.String()
}

// Multiplier returns the cost multiplier.
func (s SelfControlRoll) Multiplier() fixed.F64d4 {
	return selfControlRollValues[s.EnsureValid()].Multiplier
}

// MinimumRoll returns the minimum roll to retain control.
func (s SelfControlRoll) MinimumRoll() int {
	return selfControlRollValues[s.EnsureValid()].Value
}
