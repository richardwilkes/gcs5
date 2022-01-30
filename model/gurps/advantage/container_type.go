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

// Possible ContainerType values.
const (
	// Group is the standard grouping container type.
	Group ContainerType = iota
	// MetaTrait acts as one normal trait, listed as an advantage if its point total is positive, or a disadvantage if it
	// is negative.
	MetaTrait
	// Race tracks its point cost separately from normal advantages and disadvantages.
	Race
	// AlternativeAbilities behaves similar to a MetaTrait , but applies the rules for alternative abilities (see B61 and
	// P11) to its immediate children
	AlternativeAbilities
)

type containerTypeData struct {
	Key    string
	String string
}

// ContainerType holds the type of an advantage container.
type ContainerType uint8

var containerTypeValues = []*containerTypeData{
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

// ContainerTypeFromKey extracts a ContainerType from a key.
func ContainerTypeFromKey(key string) ContainerType {
	for i, one := range containerTypeValues {
		if strings.EqualFold(key, one.Key) {
			return ContainerType(i)
		}
	}
	return 0
}

// EnsureValid returns the first ContainerType if this ContainerType is not a known value.
func (c ContainerType) EnsureValid() ContainerType {
	if int(c) < len(containerTypeValues) {
		return c
	}
	return 0
}

// Key returns the key used to represent this ContainerType.
func (c ContainerType) Key() string {
	return containerTypeValues[c.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (c ContainerType) String() string {
	return containerTypeValues[c.EnsureValid()].String
}
