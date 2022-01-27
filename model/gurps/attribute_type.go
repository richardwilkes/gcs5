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

type attributeTypeData struct {
	Key    string
	String string
}

// AttributeType holds the type of an attribute definition.
type AttributeType uint8

var attributeTypeValues = []*attributeTypeData{
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

// AttributeTypeFromKey extracts a AttributeType from a key.
func AttributeTypeFromKey(key string) AttributeType {
	for i, one := range attributeTypeValues {
		if strings.EqualFold(key, one.Key) {
			return AttributeType(i)
		}
	}
	return 0
}

// EnsureValid returns the first AttributeType if this AttributeType is not a known value.
func (a AttributeType) EnsureValid() AttributeType {
	if int(a) < len(attributeTypeValues) {
		return a
	}
	return 0
}

// Key returns the key used to represent this AttributeType.
func (a AttributeType) Key() string {
	return attributeTypeValues[a.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (a AttributeType) String() string {
	return attributeTypeValues[a.EnsureValid()].String
}
