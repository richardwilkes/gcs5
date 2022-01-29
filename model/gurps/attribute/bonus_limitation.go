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

type bonusLimitationData struct {
	Key    string
	String string
}

// BonusLimitation holds a limitation for an AttributeBonus.
type BonusLimitation uint8

var bonusLimitationValues = []*bonusLimitationData{
	{
		Key:    "none",
		String: " ",
	},
	{
		Key:    "striking_only",
		String: i18n.Text("for striking only"),
	},
	{
		Key:    "lifting_only",
		String: i18n.Text("for lifting only"),
	},
	{
		Key:    "throwing_only",
		String: i18n.Text("for throwing only"),
	},
}

// BonusLimitationFromKey returns the BonusLimitation from a key.
func BonusLimitationFromKey(key string) BonusLimitation {
	for i, one := range bonusLimitationValues {
		if strings.EqualFold(key, one.Key) {
			return BonusLimitation(i)
		}
	}
	return 0
}

// EnsureValid returns the first BonusLimitation if this BonusLimitation is not a known value.
func (a BonusLimitation) EnsureValid() BonusLimitation {
	if int(a) < len(bonusLimitationValues) {
		return a
	}
	return 0
}

// Key returns the key used to represent this BonusLimitation.
func (a BonusLimitation) Key() string {
	return bonusLimitationValues[a.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (a BonusLimitation) String() string {
	return bonusLimitationValues[a.EnsureValid()].String
}
