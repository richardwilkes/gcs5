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

// Possible ComparisonType values.
const (
	Name         = ComparisonType("name")
	Category     = ComparisonType("category")
	College      = ComparisonType("college")
	CollegeCount = ComparisonType("college_count")
	Any          = ComparisonType("any")
)

// AllComparisonTypes is the complete set of ComparisonType values.
var AllComparisonTypes = []ComparisonType{
	Name,
	Category,
	College,
	CollegeCount,
	Any,
}

// ComparisonType holds the type of an attribute definition.
type ComparisonType string

// EnsureValid ensures this is of a known value.
func (c ComparisonType) EnsureValid() ComparisonType {
	for _, one := range AllComparisonTypes {
		if one == c {
			return c
		}
	}
	return AllComparisonTypes[0]
}

// UsesStringCriteria returns true if the comparison uses a string value.
func (c ComparisonType) UsesStringCriteria() bool {
	v := c.EnsureValid()
	return v == Name || v == Category || v == College
}
