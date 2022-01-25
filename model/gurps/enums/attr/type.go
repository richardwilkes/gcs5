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

package attr

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Type values.
const (
	Integer Type = iota
	Decimal
	Pool
)

// Type holds the type of an attribute definition.
type Type uint8

// TypeFromString extracts a Type from a string.
func TypeFromString(str string) Type {
	for op := Integer; op <= Pool; op++ {
		if strings.EqualFold(op.Key(), str) {
			return op
		}
	}
	return Integer
}

// Key returns the key used to represent this Type.
func (a Type) Key() string {
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
func (a Type) String() string {
	switch a {
	case Decimal:
		return i18n.Text("Decimal")
	case Pool:
		return i18n.Text("Pool")
	default: // Integer
		return i18n.Text("Integer")
	}
}
