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
	"github.com/richardwilkes/toolbox/i18n"
)

// ConditionalModifierType is the data type key for a ConditionalModifier.
const ConditionalModifierType = "conditional_modifier"

const (
	attributeBonusSituationKey = "situation"
)

var _ Bonus = &ConditionalModifier{}

// ConditionalModifier holds a modifier that only applies under specific situations.
type ConditionalModifier struct {
	Situation string
	Amount    *LeveledAmount
}

// NewConditionalModifier creates a new ConditionalModifier.
func NewConditionalModifier() *ConditionalModifier {
	return &ConditionalModifier{
		Situation: i18n.Text("triggering condition"),
		Amount:    &LeveledAmount{Amount: 1},
	}
}

// NewConditionalModifierFromJSON creates a new ConditionalModifier from JSON.
func NewConditionalModifierFromJSON(data map[string]interface{}) *ConditionalModifier {
	return &ConditionalModifier{
		Situation: encoding.String(data[attributeBonusSituationKey]),
		Amount:    NewLeveledAmountFromJSON(data),
	}
}

// ToJSON implements Feature.
func (m *ConditionalModifier) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(attributeBonusSituationKey, m.Situation, true, true)
	m.Amount.ToInlineJSON(encoder)
	encoder.EndObject()
}

// CloneFeature implements Feature.
func (m *ConditionalModifier) CloneFeature() Feature {
	amt := *m.Amount
	clone := *m
	clone.Amount = &amt
	return &clone
}

// DataType implements Feature.
func (m *ConditionalModifier) DataType() string {
	return ConditionalModifierType
}

// FeatureKey implements Feature.
func (m *ConditionalModifier) FeatureKey() string {
	return ConditionalModifierType
}

// FillWithNameableKeys implements Feature.
func (m *ConditionalModifier) FillWithNameableKeys(nameables map[string]string) {
	ExtractNameables(m.Situation, nameables)
}

// ApplyNameableKeys implements Feature.
func (m *ConditionalModifier) ApplyNameableKeys(nameables map[string]string) {
	m.Situation = ApplyNameables(m.Situation, nameables)
}
