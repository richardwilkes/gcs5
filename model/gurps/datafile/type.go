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

package datafile

import "github.com/richardwilkes/toolbox/i18n"

// Possible EntityType values.
const (
	PC       = EntityType("character")
	Template = EntityType("template")
)

// AllEntityTypes is the complete set of EntityType values.
var AllEntityTypes = []EntityType{
	PC,
	Template,
}

// EntityType holds the type of an Entity.
type EntityType string

// EnsureValid ensures this is of a known value.
func (t EntityType) EnsureValid() EntityType {
	for _, one := range AllEntityTypes {
		if one == t {
			return t
		}
	}
	return AllEntityTypes[0]
}

// String implements fmt.Stringer.
func (t EntityType) String() string {
	switch t {
	case PC:
		return i18n.Text("PC")
	case Template:
		return i18n.Text("Template")
	default:
		return PC.String()
	}
}
