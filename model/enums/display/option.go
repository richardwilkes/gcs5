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

// Option holds a display option.
type Option uint8

// OptionFromString extracts an Option from a string.
func OptionFromString(str string, def Option) Option {
	for one := NotShown; one <= InlineAndTooltip; one++ {
		if strings.EqualFold(one.Key(), str) {
			return one
		}
	}
	return def
}

// Key returns the key used to represent this Option.
func (o Option) Key() string {
	switch o {
	case Inline:
		return "inline"
	case Tooltip:
		return "tooltip"
	case InlineAndTooltip:
		return "inline_and_tooltip"
	default: // NotShown
		return "not_shown"
	}
}

// String implements fmt.Stringer.
func (o Option) String() string {
	switch o {
	case Inline:
		return i18n.Text("Inline")
	case Tooltip:
		return i18n.Text("Tooltip")
	case InlineAndTooltip:
		return i18n.Text("Inline & Tooltip")
	default: // NotShown
		return i18n.Text("Not Shown")
	}
}

// Inline returns true if inline notes should be shown.
func (o Option) Inline() bool {
	return o == Inline || o == InlineAndTooltip
}

// Tooltip returns true if tooltips should be shown.
func (o Option) Tooltip() bool {
	return o == Tooltip || o == InlineAndTooltip
}
