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
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Possible NumericCompareType values.
const (
	AnyNumber = NumericCompareType("any")
	Equals    = NumericCompareType("is")
	NotEquals = NumericCompareType("is_not")
	AtLeast   = NumericCompareType("at_least")
	AtMost    = NumericCompareType("at_most")
)

// AllNumericCompareTypes is the complete set of NumericCompareType values.
var AllNumericCompareTypes = []NumericCompareType{
	AnyNumber,
	Equals,
	NotEquals,
	AtLeast,
	AtMost,
}

// NumericCompareType holds the type for a numeric comparison.
type NumericCompareType string

// EnsureValid ensures this is of a known value.
func (n NumericCompareType) EnsureValid() NumericCompareType {
	for _, one := range AllNumericCompareTypes {
		if one == n {
			return n
		}
	}
	return AllNumericCompareTypes[0]
}

// String implements fmt.Stringer.
func (n NumericCompareType) String() string {
	switch n {
	case AnyNumber:
		return i18n.Text("is anything")
	case Equals:
		return i18n.Text("is")
	case NotEquals:
		return i18n.Text("is not")
	case AtLeast:
		return i18n.Text("at least")
	case AtMost:
		return i18n.Text("at most")
	default:
		return AnyNumber.String()
	}
}

// Describe returns a description of this NumericCompareType using a qualifier.
func (n NumericCompareType) Describe(qualifier fixed.F64d4) string {
	v := n.EnsureValid()
	if v == AnyNumber {
		return v.String()
	}
	return v.String() + " " + qualifier.String()
}

// Matches performs a comparison and returns true if the data matches.
func (n NumericCompareType) Matches(qualifier, data fixed.F64d4) bool {
	switch n {
	case AnyNumber:
		return true
	case Equals:
		return data == qualifier
	case NotEquals:
		return data != qualifier
	case AtLeast:
		return data >= qualifier
	case AtMost:
		return data <= qualifier
	default:
		return AnyNumber.Matches(qualifier, data)
	}
}
