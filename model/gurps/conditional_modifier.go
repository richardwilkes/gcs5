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
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/toolbox/i18n"
)

// ConditionalModifier holds the data for a conditional modifier.
type ConditionalModifier struct {
	Bonus
	Situation string `json:"situation,omitempty"`
}

// NewConditionalModifierBonus creates a new ConditionalModifier.
func NewConditionalModifierBonus() *ConditionalModifier {
	c := &ConditionalModifier{
		Bonus: Bonus{
			Feature: Feature{
				Type: feature.ConditionalModifier,
			},
			LeveledAmount: LeveledAmount{Amount: f64d4.One},
		},
		Situation: i18n.Text("triggering condition"),
	}
	c.Self = c
	return c
}

func (c *ConditionalModifier) fillWithNameableKeys(nameables map[string]string) {
	ExtractNameables(c.Situation, nameables)
}

func (c *ConditionalModifier) applyNameableKeys(nameables map[string]string) {
	c.Situation = ApplyNameables(c.Situation, nameables)
}
