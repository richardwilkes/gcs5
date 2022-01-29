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

package advantage

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Affects values.
const (
	Total Affects = iota
	BaseOnly
	LevelsOnly
)

type affectsData struct {
	Key        string
	String     string
	ShortTitle string
}

// Affects describes how an AdvantageModifier affects the point cost.
type Affects uint8

var affectsValues = []*affectsData{
	{
		Key:        "total",
		String:     i18n.Text("to cost"),
		ShortTitle: "",
	},
	{
		Key:        "base_only",
		String:     i18n.Text("to base cost only"),
		ShortTitle: i18n.Text("(base only)"),
	},
	{
		Key:        "levels_only",
		String:     i18n.Text("to leveled cost only"),
		ShortTitle: i18n.Text("(levels only)"),
	},
}

// AffectsFromKey extracts a Affects from a key.
func AffectsFromKey(key string) Affects {
	for i, one := range affectsValues {
		if strings.EqualFold(key, one.Key) {
			return Affects(i)
		}
	}
	return 0
}

// EnsureValid returns the first Affects if this Affects is not a known value.
func (a Affects) EnsureValid() Affects {
	if int(a) < len(affectsValues) {
		return a
	}
	return 0
}

// Key returns the key used to represent this Affects.
func (a Affects) Key() string {
	return affectsValues[a.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (a Affects) String() string {
	return affectsValues[a.EnsureValid()].String
}

// ShortTitle returns the short version of the title.
func (a Affects) ShortTitle() string {
	return affectsValues[a.EnsureValid()].ShortTitle
}
