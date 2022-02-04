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

package attribute

import (
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible BonusLimitation values.
const (
	None         = BonusLimitation("")
	StrikingOnly = BonusLimitation("striking_only")
	LiftingOnly  = BonusLimitation("lifting_only")
	ThrowingOnly = BonusLimitation("throwing_only")
)

// AllBonusLimitations is the complete set of BonusLimitation values.
var AllBonusLimitations = []BonusLimitation{
	None,
	StrikingOnly,
	LiftingOnly,
	ThrowingOnly,
}

// BonusLimitation holds a limitation for an AttributeBonus.
type BonusLimitation string

// EnsureValid ensures this is of a known value.
func (b BonusLimitation) EnsureValid() BonusLimitation {
	for _, one := range AllBonusLimitations {
		if one == b {
			return b
		}
	}
	return AllBonusLimitations[0]
}

// String implements fmt.Stringer.
func (b BonusLimitation) String() string {
	switch b {
	case None:
		return ""
	case StrikingOnly:
		return i18n.Text("for striking only")
	case LiftingOnly:
		return i18n.Text("for lifting only")
	case ThrowingOnly:
		return i18n.Text("for throwing only")
	default:
		return None.String()
	}
}
