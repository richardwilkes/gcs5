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

package display

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Option values.
const (
	NotShown Option = iota
	Inline
	Tooltip
	InlineAndTooltip
)

type optionData struct {
	Key     string
	String  string
	Inline  bool
	Tooltip bool
}

// Option holds a display option.
type Option uint8

var optionValues = []*optionData{
	{
		Key:    "not_shown",
		String: i18n.Text("Not Shown"),
	},
	{
		Key:    "inline",
		String: i18n.Text("Inline"),
		Inline: true,
	},
	{
		Key:     "tooltip",
		String:  i18n.Text("Tooltip"),
		Tooltip: true,
	},
	{
		Key:     "inline_and_tooltip",
		String:  i18n.Text("Inline & Tooltip"),
		Inline:  true,
		Tooltip: true,
	},
}

// OptionFromString extracts an Option from a key.
func OptionFromString(key string, def Option) Option {
	for i, one := range optionValues {
		if strings.EqualFold(key, one.Key) {
			return Option(i)
		}
	}
	return def.EnsureValid()
}

// EnsureValid returns the first Option if this Option is not a known value.
func (o Option) EnsureValid() Option {
	if int(o) < len(optionValues) {
		return o
	}
	return 0
}

// Key returns the key used to represent this Option.
func (o Option) Key() string {
	return optionValues[o.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (o Option) String() string {
	return optionValues[o.EnsureValid()].String
}

// Inline returns true if inline notes should be shown.
func (o Option) Inline() bool {
	return optionValues[o.EnsureValid()].Inline
}

// Tooltip returns true if tooltips should be shown.
func (o Option) Tooltip() bool {
	return optionValues[o.EnsureValid()].Tooltip
}
