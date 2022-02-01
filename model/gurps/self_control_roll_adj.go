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
	"fmt"
	"strings"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/advantage"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Possible SelfControlRollAdj values.
const (
	NoCRAdj SelfControlRollAdj = iota
	ActionPenalty
	ReactionPenalty
	FrightCheckPenalty
	FrightCheckBonus
	MinorCostOfLivingIncrease
	MajorCostOfLivingIncrease
)

type selfControlRollAdjData struct {
	Key         string
	String      string
	Description func(cr advantage.SelfControlRoll) string
	Features    func(cr advantage.SelfControlRoll) feature.Features
}

// SelfControlRollAdj holds an Adjustment for a self-control roll.
type SelfControlRollAdj uint8

var selfControlRollAdjValues = []*selfControlRollAdjData{
	{
		Key:         "none",
		String:      i18n.Text("None"),
		Description: func(_ advantage.SelfControlRoll) string { return "" },
		Features:    func(_ advantage.SelfControlRoll) feature.Features { return nil },
	},
	{
		Key:    "action_penalty",
		String: i18n.Text("Includes an Action Penalty for Failure"),
		Description: func(cr advantage.SelfControlRoll) string {
			return fmt.Sprintf(i18n.Text("%d Action Penalty"), int(cr)-int(advantage.None))
		},
		Features: func(_ advantage.SelfControlRoll) feature.Features { return nil },
	},
	{
		Key:    "reaction_penalty",
		String: i18n.Text("Includes a Reaction Penalty for Failure"),
		Description: func(cr advantage.SelfControlRoll) string {
			return fmt.Sprintf(i18n.Text("%d Reaction Penalty"), int(cr)-int(advantage.None))
		},
		Features: func(_ advantage.SelfControlRoll) feature.Features { return nil },
	},
	{
		Key:    "fright_check_penalty",
		String: i18n.Text("Includes Fright Check Penalty"),
		Description: func(cr advantage.SelfControlRoll) string {
			return fmt.Sprintf(i18n.Text("%d Fright Check Penalty"), int(cr)-int(advantage.None))
		},
		Features: func(_ advantage.SelfControlRoll) feature.Features { return nil },
	},
	{
		Key:    "fright_check_bonus",
		String: i18n.Text("Includes Fright Check Bonus"),
		Description: func(cr advantage.SelfControlRoll) string {
			return fmt.Sprintf(i18n.Text("+%d Fright Check Bonus"), int(advantage.None)-int(cr))
		},
		Features: func(_ advantage.SelfControlRoll) feature.Features { return nil },
	},
	{
		Key:    "minor_cost_of_living_increase",
		String: i18n.Text("Includes a Minor Cost of Living Increase"),
		Description: func(cr advantage.SelfControlRoll) string {
			return fmt.Sprintf(i18n.Text("+%d%% Cost of Living Increase"), 5*(int(advantage.None)-int(cr)))
		},
		Features: func(_ advantage.SelfControlRoll) feature.Features { return nil },
	},
	{
		Key:    "major_cost_of_living_increase",
		String: i18n.Text("Includes a Major Cost of Living Increase and Merchant Skill Penalty"),
		Description: func(cr advantage.SelfControlRoll) string {
			return fmt.Sprintf(i18n.Text("+%d%% Cost of Living Increase"), 10*(1<<((int(advantage.None)-int(cr))-1)))
		},
		Features: func(cr advantage.SelfControlRoll) feature.Features {
			f := feature.NewSkillBonus()
			f.NameCriteria.Qualifier = "Merchant"
			f.Amount = fixed.F64d4FromInt64(int64(cr) - int64(advantage.None))
			return feature.Features{f}
		},
	},
}

// EnsureValid returns the first SelfControlRollAdj if this SelfControlRollAdj is not a known value.
func (s SelfControlRollAdj) EnsureValid() SelfControlRollAdj {
	if int(s) < len(selfControlRollAdjValues) {
		return s
	}
	return NoCRAdj
}

// SelfControlRollAdjFromKey extracts a SelfControlRollAdj from a key.
func SelfControlRollAdjFromKey(key string) SelfControlRollAdj {
	for i, one := range selfControlRollAdjValues {
		if strings.EqualFold(key, one.Key) {
			return SelfControlRollAdj(i)
		}
	}
	return 0
}

// ToKeyedJSON writes the SelfControlRollAdj to JSON.
func (s SelfControlRollAdj) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	if resolved := s.EnsureValid(); resolved != NoCRAdj {
		encoder.KeyedString(key, selfControlRollAdjValues[resolved].Key, false, false)
	}
}

// String implements fmt.Stringer.
func (s SelfControlRollAdj) String() string {
	return selfControlRollAdjValues[s.EnsureValid()].String
}

// Description returns a formatted description.
func (s SelfControlRollAdj) Description(cr advantage.SelfControlRoll) string {
	if cr == advantage.None {
		return ""
	}
	return selfControlRollAdjValues[s.EnsureValid()].Description(cr)
}

// Features returns the set of features to apply.
func (s SelfControlRollAdj) Features(cr advantage.SelfControlRoll) feature.Features {
	return selfControlRollAdjValues[s.EnsureValid()].Features(cr)
}
