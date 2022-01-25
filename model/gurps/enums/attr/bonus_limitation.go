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

package attr

import (
	"strings"

	"github.com/richardwilkes/toolbox/i18n"
)

// Possible BonusLimitation values.
const (
	None BonusLimitation = iota
	StrikingOnly
	LiftingOnly
	ThrowingOnly
)

// BonusLimitation holds a limitation for an AttributeBonus.
type BonusLimitation uint8

// BonusLimitationFromString returns the BonusLimitation for the given key, or a default of None if nothing matches.
func BonusLimitationFromString(key string) BonusLimitation {
	for one := None; one <= ThrowingOnly; one++ {
		if strings.EqualFold(one.Key(), key) {
			return one
		}
	}
	return None
}

// Key returns the key used to represent this BonusLimitation.
func (l BonusLimitation) Key() string {
	switch l {
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
func (l BonusLimitation) String() string {
	switch l {
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
