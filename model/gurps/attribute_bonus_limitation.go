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

// Possible AttributeBonusLimitation values.
const (
	None AttributeBonusLimitation = iota
	StrikingOnly
	LiftingOnly
	ThrowingOnly
)

type attributeBonusLimitationData struct {
	Key    string
	String string
}

// AttributeBonusLimitation holds a limitation for an AttributeBonus.
type AttributeBonusLimitation uint8

var attributeBonusLimitationValues = []*attributeBonusLimitationData{
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

// AttributeBonusLimitationFromKey returns the AttributeBonusLimitation from a key.
func AttributeBonusLimitationFromKey(key string) AttributeBonusLimitation {
	for i, one := range attributeBonusLimitationValues {
		if strings.EqualFold(key, one.Key) {
			return AttributeBonusLimitation(i)
		}
	}
	return 0
}

// EnsureValid returns the first AttributeBonusLimitation if this AttributeBonusLimitation is not a known value.
func (a AttributeBonusLimitation) EnsureValid() AttributeBonusLimitation {
	if int(a) < len(attributeBonusLimitationValues) {
		return a
	}
	return 0
}

// Key returns the key used to represent this AttributeBonusLimitation.
func (a AttributeBonusLimitation) Key() string {
	return attributeBonusLimitationValues[a.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (a AttributeBonusLimitation) String() string {
	return attributeBonusLimitationValues[a.EnsureValid()].String
}
