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
)

// Possible StringCompareType values.
const (
	Any StringCompareType = iota
	Is
	IsNot
	Contains
	DoesNotContain
	StartsWith
	DoesNotStartWith
	EndsWith
	DoesNotEndWith
)

type stringCompareTypeData struct {
	Key     string
	String  string
	Matches func(qualifier, data string) bool
}

// StringCompareType holds the type for a string comparison.
type StringCompareType uint8

var stringCompareTypeValues = []*stringCompareTypeData{
	{
		Key:     "any",
		String:  i18n.Text("is anything"),
		Matches: func(qualifier, data string) bool { return true },
	},
	{
		Key:     "is",
		String:  i18n.Text("is"),
		Matches: func(qualifier, data string) bool { return strings.EqualFold(data, qualifier) },
	},
	{
		Key:     "is_not",
		String:  i18n.Text("is not"),
		Matches: func(qualifier, data string) bool { return !strings.EqualFold(data, qualifier) },
	},
	{
		Key:    "contains",
		String: i18n.Text("contains"),
		Matches: func(qualifier, data string) bool {
			return strings.Contains(strings.ToLower(data), strings.ToLower(qualifier))
		},
	},
	{
		Key:    "does_not_contain",
		String: i18n.Text("does not contain"),
		Matches: func(qualifier, data string) bool {
			return !strings.Contains(strings.ToLower(data), strings.ToLower(qualifier))
		},
	},
	{
		Key:    "starts_with",
		String: i18n.Text("starts with"),
		Matches: func(qualifier, data string) bool {
			return strings.HasPrefix(strings.ToLower(data), strings.ToLower(qualifier))
		},
	},
	{
		Key:    "does_not_start_with",
		String: i18n.Text("does not start with"),
		Matches: func(qualifier, data string) bool {
			return !strings.HasPrefix(strings.ToLower(data), strings.ToLower(qualifier))
		},
	},
	{
		Key:    "ends_with",
		String: i18n.Text("ends with"),
		Matches: func(qualifier, data string) bool {
			return strings.HasSuffix(strings.ToLower(data), strings.ToLower(qualifier))
		},
	},
	{
		Key:    "does_not_end_with",
		String: i18n.Text("does not end with"),
		Matches: func(qualifier, data string) bool {
			return !strings.HasSuffix(strings.ToLower(data), strings.ToLower(qualifier))
		},
	},
}

// StringCompareTypeFromString extracts a StringCompareType from a key.
func StringCompareTypeFromString(key string) StringCompareType {
	for i, one := range stringCompareTypeValues {
		if strings.EqualFold(key, one.Key) {
			return StringCompareType(i)
		}
	}
	return 0
}

// EnsureValid returns the first StringCompareType if this StringCompareType is not a known value.
func (s StringCompareType) EnsureValid() StringCompareType {
	if int(s) < len(stringCompareTypeValues) {
		return s
	}
	return 0
}

// Key returns the key used to represent this StringCompareType.
func (s StringCompareType) Key() string {
	return stringCompareTypeValues[s.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (s StringCompareType) String() string {
	return stringCompareTypeValues[s.EnsureValid()].String
}

// Describe returns a description of this StringCompareType using a qualifier.
func (s StringCompareType) Describe(qualifier string) string {
	if s == Any {
		return s.String()
	}
	return s.String() + ` "` + qualifier + `"`
}

// Matches performs a comparison and returns true if the data matches.
func (s StringCompareType) Matches(qualifier, data string) bool {
	return stringCompareTypeValues[s.EnsureValid()].Matches(qualifier, data)
}
