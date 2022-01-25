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
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/enum"
)

// AttributeBonusType is the data type key for an AttributeBonus.
const AttributeBonusType = "attribute_bonus"

const (
	attributeBonusAttributeKey  = "attribute"
	attributeBonusLimitationKey = "limitation"
)

var _ Bonus = &AttributeBonus{}

// AttributeBonus holds a bonus to an Attribute.
type AttributeBonus struct {
	Attribute  string
	Limitation enum.AttributeBonusLimitation
	Amount     *LeveledAmount
}

// NewAttributeBonus creates a new AttributeBonus for the given entity, which may be nil.
func NewAttributeBonus(entity *Entity) *AttributeBonus {
	return &AttributeBonus{
		Attribute: DefaultAttributeIDFor(entity),
		Amount:    &LeveledAmount{Amount: 1},
	}
}

// NewAttributeBonusFromJSON creates a new AttributeBonus from JSON.
func NewAttributeBonusFromJSON(data map[string]interface{}) *AttributeBonus {
	return &AttributeBonus{
		Attribute:  encoding.String(data[attributeBonusAttributeKey]),
		Limitation: enum.AttributeBonusLimitationFromString(encoding.String(data[attributeBonusLimitationKey])),
		Amount:     NewLeveledAmountFromJSON(data),
	}
}

// ToJSON implements Feature.
func (a *AttributeBonus) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(attributeBonusAttributeKey, a.Attribute, false, false)
	if a.Limitation != enum.None {
		encoder.KeyedString(attributeBonusLimitationKey, a.Limitation.Key(), false, false)
	}
	a.Amount.ToInlineJSON(encoder)
	encoder.EndObject()
}

// CloneFeature implements Feature.
func (a *AttributeBonus) CloneFeature() Feature {
	amt := *a.Amount
	clone := *a
	clone.Amount = &amt
	return &clone
}

// DataType implements Feature.
func (a *AttributeBonus) DataType() string {
	return AttributeBonusType
}

// FeatureKey implements Feature.
func (a *AttributeBonus) FeatureKey() string {
	key := AttributeIDPrefix + a.Attribute
	if a.Limitation != enum.None {
		key += "." + a.Limitation.Key()
	}
	return key
}

// FillWithNameableKeys implements Feature.
func (a *AttributeBonus) FillWithNameableKeys(_ map[string]string) {
	// Does nothing
}

// ApplyNameableKeys implements Feature.
func (a *AttributeBonus) ApplyNameableKeys(_ map[string]string) {
	// Does nothing
}

// Normalize implements Feature.
func (a *AttributeBonus) Normalize() {
	// Unused
}
