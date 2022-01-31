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
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/gurps/feature"
)

// AttributeBonus holds the data for a bonus to an attribute.
type AttributeBonus struct {
	Bonus
	Attribute  string                    `json:"attribute"`
	Limitation attribute.BonusLimitation `json:"limitation,omitempty"`
}

// NewAttributeBonus creates a new AttributeBonus.
func NewAttributeBonus(entity *Entity) *AttributeBonus {
	a := &AttributeBonus{
		Bonus: Bonus{
			Feature: Feature{
				Type: feature.AttributeBonus,
			},
			LeveledAmount: LeveledAmount{Amount: f64d4.One},
		},
		Attribute:  DefaultAttributeIDFor(entity),
		Limitation: attribute.None,
	}
	a.Self = a
	return a
}

func (a *AttributeBonus) featureMapKey() string {
	key := AttributeIDPrefix + a.Attribute
	if a.Limitation != attribute.None {
		key += "." + string(a.Limitation)
	}
	return key
}
