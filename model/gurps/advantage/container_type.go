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
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible ContainerType values.
const (
	// Group is the standard grouping container type.
	Group = ContainerType("group")
	// MetaTrait acts as one normal trait, listed as an advantage if its point total is positive, or a disadvantage if it
	// is negative.
	MetaTrait = ContainerType("meta_trait")
	// Race tracks its point cost separately from normal advantages and disadvantages.
	Race = ContainerType("race")
	// AlternativeAbilities behaves similar to a MetaTrait , but applies the rules for alternative abilities (see B61 and
	// P11) to its immediate children
	AlternativeAbilities = ContainerType("alternative_abilities")
)

// AllContainerTypes is the complete set of ContainerType values.
var AllContainerTypes = []ContainerType{
	Group,
	MetaTrait,
	Race,
	AlternativeAbilities,
}

// ContainerType holds the type of an advantage container.
type ContainerType string

// EnsureValid ensures this is of a known value.
func (c ContainerType) EnsureValid() ContainerType {
	for _, one := range AllContainerTypes {
		if one == c {
			return c
		}
	}
	return AllContainerTypes[0]
}

// String implements fmt.Stringer.
func (c ContainerType) String() string {
	switch c {
	case Group:
		return i18n.Text("Group")
	case MetaTrait:
		return i18n.Text("Meta-Trait")
	case Race:
		return i18n.Text("Race")
	case AlternativeAbilities:
		return i18n.Text("Alternative Abilities")
	default:
		return Group.String()
	}
}
