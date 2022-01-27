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

// Possible SkillSelectionType values.
const (
	SkillsWithNameSkillSelect SkillSelectionType = iota
	ThisWeaponSkillSelect
	WeaponsWithNameSkillSelect
)

type skillSelectionTypeData struct {
	Key    string
	String string
}

// SkillSelectionType holds the type of an attribute definition.
type SkillSelectionType uint8

var skillSelectionTypeValues = []*skillSelectionTypeData{
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

// SkillSelectionTypeFromString extracts a SkillSelectionType from a key.
func SkillSelectionTypeFromString(key string) SkillSelectionType {
	for i, one := range skillSelectionTypeValues {
		if strings.EqualFold(key, one.Key) {
			return SkillSelectionType(i)
		}
	}
	return 0
}

// EnsureValid returns the first SkillSelectionType if this SkillSelectionType is not a known value.
func (s SkillSelectionType) EnsureValid() SkillSelectionType {
	if int(s) < len(skillSelectionTypeValues) {
		return s
	}
	return 0
}

// Key returns the key used to represent this SkillSelectionType.
func (s SkillSelectionType) Key() string {
	return skillSelectionTypeValues[s.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (s SkillSelectionType) String() string {
	return skillSelectionTypeValues[s.EnsureValid()].String
}
