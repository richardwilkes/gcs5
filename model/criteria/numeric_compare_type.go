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
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Possible NumericCompareType values.
const (
	AnyNumber NumericCompareType = iota
	Equals
	NotEquals
	AtLeast
	AtMost
)

type numericCompareTypeData struct {
	Key     string
	String  string
	Matches func(qualifier, data fixed.F64d4) bool
}

// NumericCompareType holds the type for a numeric comparison.
type NumericCompareType uint8

var numericCompareTypeValues = []*numericCompareTypeData{
	{
		Key:     "any",
		String:  i18n.Text("is anything"),
		Matches: func(qualifier, data fixed.F64d4) bool { return true },
	},
	{
		Key:     "is",
		String:  i18n.Text("is"),
		Matches: func(qualifier, data fixed.F64d4) bool { return data == qualifier },
	},
	{
		Key:     "is_not",
		String:  i18n.Text("is not"),
		Matches: func(qualifier, data fixed.F64d4) bool { return data != qualifier },
	},
	{
		Key:     "at_least",
		String:  i18n.Text("at least"),
		Matches: func(qualifier, data fixed.F64d4) bool { return data >= qualifier },
	},
	{
		Key:     "at_most",
		String:  i18n.Text("at most"),
		Matches: func(qualifier, data fixed.F64d4) bool { return data <= qualifier },
	},
}

// NumericCompareTypeFromString extracts a NumericCompareType from a key.
func NumericCompareTypeFromString(key string) NumericCompareType {
	for i, one := range numericCompareTypeValues {
		if strings.EqualFold(key, one.Key) {
			return NumericCompareType(i)
		}
	}
	return 0
}

// EnsureValid returns the first NumericCompareType if this NumericCompareType is not a known value.
func (n NumericCompareType) EnsureValid() NumericCompareType {
	if int(n) < len(numericCompareTypeValues) {
		return n
	}
	return 0
}

// Key returns the key used to represent this NumericCompareType.
func (n NumericCompareType) Key() string {
	return numericCompareTypeValues[n.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (n NumericCompareType) String() string {
	return numericCompareTypeValues[n.EnsureValid()].String
}

// Describe returns a description of this NumericCompareType using a qualifier.
func (n NumericCompareType) Describe(qualifier fixed.F64d4) string {
	if n == AnyNumber {
		return n.String()
	}
	return n.String() + " " + qualifier.String()
}

// Matches performs a comparison and returns true if the data matches.
func (n NumericCompareType) Matches(qualifier, data fixed.F64d4) bool {
	return numericCompareTypeValues[n.EnsureValid()].Matches(qualifier, data)
}
