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

type ReactionBonus struct {
	Bonus
	Situation string `json:"situation,omitempty"`
}

// NewReactionBonus creates a new ReactionBonus.
func NewReactionBonus() *ReactionBonus {
	r := &ReactionBonus{
		Bonus: Bonus{
			Feature: Feature{
				Type: feature.ReactionBonus,
			},
			LeveledAmount: LeveledAmount{Amount: f64d4.One},
		},
		Situation: i18n.Text("from others"),
	}
	r.Self = r
	return r
}

func (r *ReactionBonus) featureMapKey() string {
	return "reaction"
}

func (r *ReactionBonus) fillWithNameableKeys(nameables map[string]string) {
	ExtractNameables(r.Situation, nameables)
}

func (r *ReactionBonus) applyNameableKeys(nameables map[string]string) {
	r.Situation = ApplyNameables(r.Situation, nameables)
}
