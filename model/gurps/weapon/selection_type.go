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
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible SelectionType values.
const (
	WithRequiredSkill = SelectionType("weapons_with_required_skill")
	ThisWeapon        = SelectionType("this_weapon")
	WithName          = SelectionType("weapons_with_name")
)

// AllSelectionTypes is the complete set of SelectionType values.
var AllSelectionTypes = []SelectionType{
	WithRequiredSkill,
	ThisWeapon,
	WithName,
}

// SelectionType holds the type of an attribute definition.
type SelectionType string

// EnsureValid ensures this is of a known value.
func (s SelectionType) EnsureValid() SelectionType {
	for _, one := range AllSelectionTypes {
		if one == s {
			return s
		}
	}
	return AllSelectionTypes[0]
}

// String implements fmt.Stringer.
func (s SelectionType) String() string {
	switch s {
	case WithRequiredSkill:
		return i18n.Text("to weapons whose required skill name")
	case ThisWeapon:
		return i18n.Text("to this weapon")
	case WithName:
		return i18n.Text("to weapons whose name")
	default:
		return WithRequiredSkill.String()
	}
}
