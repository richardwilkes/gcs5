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
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Type values.
const (
	Integer Type = iota
	Decimal
	Pool
)

type typeData struct {
	Key    string
	String string
}

// Type holds the type of an attribute definition.
type Type uint8

var typeValues = []*typeData{
	{
		Key:    "integer",
		String: i18n.Text("Integer"),
	},
	{
		Key:    "decimal",
		String: i18n.Text("Decimal"),
	},
	{
		Key:    "pool",
		String: i18n.Text("Pool"),
	},
}

// TypeFromKey extracts a Type from a key.
func TypeFromKey(key string) Type {
	for i, one := range typeValues {
		if strings.EqualFold(key, one.Key) {
			return Type(i)
		}
	}
	return 0
}

// EnsureValid returns the first Type if this Type is not a known value.
func (t Type) EnsureValid() Type {
	if int(t) < len(typeValues) {
		return t
	}
	return 0
}

// Key returns the key used to represent this Type.
func (t Type) Key() string {
	return typeValues[t.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (t Type) String() string {
	return typeValues[t.EnsureValid()].String
}
