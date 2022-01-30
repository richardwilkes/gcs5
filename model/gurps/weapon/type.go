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

package weapon

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Type values.
const (
	Melee Type = iota
	Ranged
)

type typeData struct {
	Key    string
	String string
}

// Type holds the type of an weapon definition.
type Type uint8

var typeValues = []*typeData{
	{
		Key:    "melee_weapon",
		String: i18n.Text("Melee Weapon"),
	},
	{
		Key:    "ranged_weapon",
		String: i18n.Text("Ranged Weapon"),
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
