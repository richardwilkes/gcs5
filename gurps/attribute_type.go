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

package gurps

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible AttributeType values.
const (
	IntegerAttributeType AttributeType = iota
	DecimalAttributeType
	PoolAttributeType
)

// AttributeType holds the type of an AttributeDef.
type AttributeType uint8

// AttributeTypeFromString extracts a AttributeType from a string.
func AttributeTypeFromString(str string) AttributeType {
	for op := IntegerAttributeType; op <= PoolAttributeType; op++ {
		if strings.EqualFold(op.Key(), str) {
			return op
		}
	}
	return IntegerAttributeType
}

// Key returns the key used to represent this ThresholdOp.
func (a AttributeType) Key() string {
	switch a {
	case DecimalAttributeType:
		return "decimal"
	case PoolAttributeType:
		return "pool"
	default: // IntegerAttributeType
		return "integer"
	}
}

// MarshalText implements encoding.TextMarshaler.
func (a AttributeType) MarshalText() (text []byte, err error) {
	return []byte(a.Key()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (a *AttributeType) UnmarshalText(text []byte) error {
	*a = AttributeTypeFromString(string(text))
	return nil
}

// String implements fmt.Stringer.
func (a AttributeType) String() string {
	switch a {
	case DecimalAttributeType:
		return i18n.Text("Decimal")
	case PoolAttributeType:
		return i18n.Text("Pool")
	default: // IntegerAttributeType
		return i18n.Text("Integer")
	}
}
