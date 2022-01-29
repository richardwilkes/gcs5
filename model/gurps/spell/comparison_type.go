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

package spell

import (
	"strings"
)

// Possible ComparisonType values.
const (
	Name ComparisonType = iota
	Category
	College
	CollegeCount
	Any
)

type comparisonTypeData struct {
	Key                string
	UsesStringCriteria bool
}

// ComparisonType holds the type of an attribute definition.
type ComparisonType uint8

var comparisonTypeValues = []*comparisonTypeData{
	{
		Key:                "name",
		UsesStringCriteria: true,
	},
	{
		Key:                "category",
		UsesStringCriteria: true,
	},
	{
		Key:                "college",
		UsesStringCriteria: true,
	},
	{
		Key:                "college_count",
		UsesStringCriteria: false,
	},
	{
		Key:                "any",
		UsesStringCriteria: false,
	},
}

// ComparisonTypeFromString extracts a ComparisonType from a key.
func ComparisonTypeFromString(key string) ComparisonType {
	for i, one := range comparisonTypeValues {
		if strings.EqualFold(key, one.Key) {
			return ComparisonType(i)
		}
	}
	return 0
}

// EnsureValid returns the first ComparisonType if this ComparisonType is not a known value.
func (s ComparisonType) EnsureValid() ComparisonType {
	if int(s) < len(comparisonTypeValues) {
		return s
	}
	return 0
}

// Key returns the key used to represent this ComparisonType.
func (s ComparisonType) Key() string {
	return comparisonTypeValues[s.EnsureValid()].Key
}

// UsesStringCriteria returns true if the comparison uses a string value.
func (s ComparisonType) UsesStringCriteria() bool {
	return comparisonTypeValues[s.EnsureValid()].UsesStringCriteria
}
