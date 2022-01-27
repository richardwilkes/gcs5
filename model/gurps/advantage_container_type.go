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

// Possible AdvantageContainerType values.
const (
	// Group is the standard grouping container type.
	Group AdvantageContainerType = iota
	// MetaTrait acts as one normal trait, listed as an advantage if its point total is positive, or a disadvantage if it
	// is negative.
	MetaTrait
	// Race tracks its point cost separately from normal advantages and disadvantages.
	Race
	// AlternativeAbilities behaves similar to a MetaTrait , but applies the rules for alternative abilities (see B61 and
	// P11) to its immediate children
	AlternativeAbilities
)

type advantageContainerTypeData struct {
	Key    string
	String string
}

// AdvantageContainerType holds the type of an advantage container.
type AdvantageContainerType uint8

var advantageContainerTypeValues = []*advantageContainerTypeData{
	{
		Key:    "group",
		String: i18n.Text("Group"),
	},
	{
		Key:    "meta_trait",
		String: i18n.Text("Meta-Trait"),
	},
	{
		Key:    "race",
		String: i18n.Text("Race"),
	},
	{
		Key:    "alternative_abilities",
		String: i18n.Text("Alternative Abilities"),
	},
}

// AdvantageContainerTypeFromKey extracts a AdvantageContainerType from a key.
func AdvantageContainerTypeFromKey(key string) AdvantageContainerType {
	for i, one := range advantageContainerTypeValues {
		if strings.EqualFold(key, one.Key) {
			return AdvantageContainerType(i)
		}
	}
	return 0
}

// EnsureValid returns the first AdvantageContainerType if this AdvantageContainerType is not a known value.
func (a AdvantageContainerType) EnsureValid() AdvantageContainerType {
	if int(a) < len(advantageContainerTypeValues) {
		return a
	}
	return 0
}

// Key returns the key used to represent this AdvantageContainerType.
func (a AdvantageContainerType) Key() string {
	return advantageContainerTypeValues[a.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (a AdvantageContainerType) String() string {
	return advantageContainerTypeValues[a.EnsureValid()].String
}
