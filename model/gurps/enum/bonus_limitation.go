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

package enum

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible AttributeBonusLimitation values.
const (
	None AttributeBonusLimitation = iota
	StrikingOnly
	LiftingOnly
	ThrowingOnly
)

// AttributeBonusLimitation holds a limitation for an AttributeBonus.
type AttributeBonusLimitation uint8

// AttributeBonusLimitationFromString returns the AttributeBonusLimitation for the given key, or a default of None if
// nothing matches.
func AttributeBonusLimitationFromString(key string) AttributeBonusLimitation {
	for one := None; one <= ThrowingOnly; one++ {
		if strings.EqualFold(one.Key(), key) {
			return one
		}
	}
	return None
}

// Key returns the key used to represent this AttributeBonusLimitation.
func (a AttributeBonusLimitation) Key() string {
	switch a {
	case StrikingOnly:
		return "striking_only"
	case LiftingOnly:
		return "lifting_only"
	case ThrowingOnly:
		return "throwing_only"
	default: // None
		return "none"
	}
}

// String implements fmt.Stringer.
func (a AttributeBonusLimitation) String() string {
	switch a {
	case StrikingOnly:
		return i18n.Text("for striking only")
	case LiftingOnly:
		return i18n.Text("for lifting only")
	case ThrowingOnly:
		return i18n.Text("for throwing only")
	default: // None
		return " "
	}
}
