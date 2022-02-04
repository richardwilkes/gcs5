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

package attribute

import (
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Type values.
const (
	Integer = Type("integer")
	Decimal = Type("decimal")
	Pool    = Type("pool")
)

// AllTypes is the complete set of Type values.
var AllTypes = []Type{
	Integer,
	Decimal,
	Pool,
}

// Type holds the type of an attribute definition.
type Type string

// EnsureValid ensures this is of a known value.
func (t Type) EnsureValid() Type {
	for _, one := range AllTypes {
		if one == t {
			return t
		}
	}
	return AllTypes[0]
}

// String implements fmt.Stringer.
func (t Type) String() string {
	switch t {
	case Integer:
		return i18n.Text("Integer")
	case Decimal:
		return i18n.Text("Decimal")
	case Pool:
		return i18n.Text("Pool")
	default:
		return Integer.String()
	}
}
