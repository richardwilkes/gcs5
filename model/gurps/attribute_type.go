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
	Integer AttributeType = iota
	Decimal
	Pool
)

// AttributeType holds the type of an attribute definition.
type AttributeType uint8

// AttributeTypeFromString extracts a AttributeType from a string.
func AttributeTypeFromString(str string) AttributeType {
	for one := Integer; one <= Pool; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return Integer
}

// Key returns the key used to represent this AttributeType.
func (a AttributeType) Key() string {
	switch a {
	case Decimal:
		return "decimal"
	case Pool:
		return "pool"
	default: // Integer
		return "integer"
	}
}

// String implements fmt.Stringer.
func (a AttributeType) String() string {
	switch a {
	case Decimal:
		return i18n.Text("Decimal")
	case Pool:
		return i18n.Text("Pool")
	default: // Integer
		return i18n.Text("Integer")
	}
}
