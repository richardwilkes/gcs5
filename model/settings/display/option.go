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
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Option values.
const (
	NotShown         = Option("not_shown")
	Inline           = Option("inline")
	Tooltip          = Option("tooltip")
	InlineAndTooltip = Option("inline_and_tooltip")
)

// AllOptions is the complete set of Option values.
var AllOptions = []Option{
	NotShown,
	Inline,
	Tooltip,
	InlineAndTooltip,
}

// Option holds a display option.
type Option string

// EnsureValid ensures this is of a known value.
func (o Option) EnsureValid() Option {
	for _, one := range AllOptions {
		if one == o {
			return o
		}
	}
	return AllOptions[0]
}

// String implements fmt.Stringer.
func (o Option) String() string {
	switch o {
	case NotShown:
		return i18n.Text("Not Shown")
	case Inline:
		return i18n.Text("Inline")
	case Tooltip:
		return i18n.Text("Tooltip")
	case InlineAndTooltip:
		return i18n.Text("Inline & Tooltip")
	default:
		return NotShown.String()
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
