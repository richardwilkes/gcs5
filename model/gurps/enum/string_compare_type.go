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

package enum

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

// StringCompareType holds the type for a string comparison.
type StringCompareType uint8

// StringCompareTypeFromString extracts a StringCompareType from a string.
func StringCompareTypeFromString(str string) StringCompareType {
	for one := Any; one <= DoesNotEndWith; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return Any
}

// Key returns the key used to represent this StringCompareType.
func (s StringCompareType) Key() string {
	switch s {
	case Is:
		return "is"
	case IsNot:
		return "is_not"
	case Contains:
		return "contains"
	case DoesNotContain:
		return "does_not_contain"
	case StartsWith:
		return "starts_with"
	case DoesNotStartWith:
		return "does_not_start_with"
	case EndsWith:
		return "ends_with"
	case DoesNotEndWith:
		return "does_not_end_with"
	default: // Any
		return "any"
	}
}

// String implements fmt.Stringer.
func (s StringCompareType) String() string {
	switch s {
	case Is:
		return i18n.Text("is")
	case IsNot:
		return i18n.Text("is not")
	case Contains:
		return i18n.Text("contains")
	case DoesNotContain:
		return i18n.Text("does not contain")
	case StartsWith:
		return i18n.Text("starts with")
	case DoesNotStartWith:
		return i18n.Text("does not start with")
	case EndsWith:
		return i18n.Text("ends with")
	case DoesNotEndWith:
		return i18n.Text("does not end with")
	default: // Any
		return i18n.Text("is anything")
	}
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
	switch s {
	case Is:
		return strings.EqualFold(data, qualifier)
	case IsNot:
		return !strings.EqualFold(data, qualifier)
	case Contains:
		return strings.Contains(strings.ToLower(data), strings.ToLower(qualifier))
	case DoesNotContain:
		return !strings.Contains(strings.ToLower(data), strings.ToLower(qualifier))
	case StartsWith:
		return strings.HasPrefix(strings.ToLower(data), strings.ToLower(qualifier))
	case DoesNotStartWith:
		return !strings.HasPrefix(strings.ToLower(data), strings.ToLower(qualifier))
	case EndsWith:
		return strings.HasSuffix(strings.ToLower(data), strings.ToLower(qualifier))
	case DoesNotEndWith:
		return !strings.HasSuffix(strings.ToLower(data), strings.ToLower(qualifier))
	default: // Any
		return true
	}
}
