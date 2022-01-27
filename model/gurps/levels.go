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

// Possible Levels values.
const (
	NoLevels Levels = iota
	HasLevels
	HasHalfLevels
)

type levelsData struct {
	Key    string
	String string
}

// Levels holds the type of leveling that can be done.
type Levels uint8

var levelsValues = []*levelsData{
	{
		Key:    "no_levels",
		String: i18n.Text("Has No Levels"),
	},
	{
		Key:    "has_levels",
		String: i18n.Text("Has Levels"),
	},
	{
		Key:    "has_half_levels",
		String: i18n.Text("Has Half Levels"),
	},
}

// LevelsFromKey extracts a Levels from a key.
func LevelsFromKey(key string) Levels {
	for i, one := range levelsValues {
		if strings.EqualFold(key, one.Key) {
			return Levels(i)
		}
	}
	return 0
}

// EnsureValid returns the first Levels if this Levels is not a known value.
func (l Levels) EnsureValid() Levels {
	if int(l) < len(levelsValues) {
		return l
	}
	return 0
}

// Key returns the key used to represent this Levels.
func (l Levels) Key() string {
	return levelsValues[l.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (l Levels) String() string {
	return levelsValues[l.EnsureValid()].String
}
