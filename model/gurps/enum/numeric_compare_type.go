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

// NumericCompareType holds the type for a numeric comparison.
type NumericCompareType uint8

// NumericCompareTypeFromString extracts a NumericCompareType from a string.
func NumericCompareTypeFromString(str string) NumericCompareType {
	for one := AnyNumber; one <= AtMost; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return Equals
}

// Key returns the key used to represent this NumericCompareType.
func (n NumericCompareType) Key() string {
	switch n {
	case Equals:
		return "is"
	case NotEquals:
		return "is_not"
	case AtLeast:
		return "at_least"
	case AtMost:
		return "at_most"
	default: // AnyNumber
		return "any"
	}
}

// String implements fmt.Stringer.
func (n NumericCompareType) String() string {
	switch n {
	case Equals:
		return i18n.Text("is")
	case NotEquals:
		return i18n.Text("is not")
	case AtLeast:
		return i18n.Text("at least")
	case AtMost:
		return i18n.Text("at most")
	default: // AnyNumber
		return i18n.Text("is anything")
	}
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
	switch n {
	case Equals:
		return data == qualifier
	case NotEquals:
		return data != qualifier
	case AtLeast:
		return data >= qualifier
	case AtMost:
		return data <= qualifier
	default: // AnyNumber
		return true
	}
}
