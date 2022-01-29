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

package skill

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible SelectionType values.
const (
	SkillsWithName SelectionType = iota
	ThisWeapon
	WeaponsWithName
)

type selectionTypeData struct {
	Key    string
	String string
}

// SelectionType holds the type of an attribute definition.
type SelectionType uint8

var selectionTypeValues = []*selectionTypeData{
	{
		Key:    "skills_with_name",
		String: i18n.Text("to skills whose name"),
	},
	{
		Key:    "this_weapon",
		String: i18n.Text("to this weapon"),
	},
	{
		Key:    "weapons_with_name",
		String: i18n.Text("to weapons whose name"),
	},
}

// SelectionTypeFromString extracts a SelectionType from a key.
func SelectionTypeFromString(key string) SelectionType {
	for i, one := range selectionTypeValues {
		if strings.EqualFold(key, one.Key) {
			return SelectionType(i)
		}
	}
	return 0
}

// EnsureValid returns the first SelectionType if this SelectionType is not a known value.
func (s SelectionType) EnsureValid() SelectionType {
	if int(s) < len(selectionTypeValues) {
		return s
	}
	return 0
}

// Key returns the key used to represent this SelectionType.
func (s SelectionType) Key() string {
	return selectionTypeValues[s.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (s SelectionType) String() string {
	return selectionTypeValues[s.EnsureValid()].String
}
