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

	"github.com/richardwilkes/gcs/model/gurps/advantage"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Possible SelfControlRollAdj values.
const (
	NoCRAdj                   = SelfControlRollAdj("")
	ActionPenalty             = SelfControlRollAdj("action_penalty")
	ReactionPenalty           = SelfControlRollAdj("reaction_penalty")
	FrightCheckPenalty        = SelfControlRollAdj("fright_check_penalty")
	FrightCheckBonus          = SelfControlRollAdj("fright_check_bonus")
	MinorCostOfLivingIncrease = SelfControlRollAdj("minor_cost_of_living_increase")
	MajorCostOfLivingIncrease = SelfControlRollAdj("major_cost_of_living_increase")
)

// AllSelfControlRollAdjs is the complete set of SelfControlRollAdj values.
var AllSelfControlRollAdjs = []SelfControlRollAdj{
	NoCRAdj,
	ActionPenalty,
	ReactionPenalty,
	FrightCheckPenalty,
	FrightCheckBonus,
	MinorCostOfLivingIncrease,
	MajorCostOfLivingIncrease,
}

// SelfControlRollAdj holds an Adjustment for a self-control roll.
type SelfControlRollAdj string

// EnsureValid ensures this is of a known value.
func (s SelfControlRollAdj) EnsureValid() SelfControlRollAdj {
	for _, one := range AllSelfControlRollAdjs {
		if one == s {
			return s
		}
	}
	return AllSelfControlRollAdjs[0]
}

// String implements fmt.Stringer.
func (s SelfControlRollAdj) String() string {
	switch s {
	case NoCRAdj:
		return i18n.Text("None")
	case ActionPenalty:
		return i18n.Text("Includes an Action Penalty for Failure")
	case ReactionPenalty:
		return i18n.Text("Includes a Reaction Penalty for Failure")
	case FrightCheckPenalty:
		return i18n.Text("Includes Fright Check Penalty")
	case FrightCheckBonus:
		return i18n.Text("Includes Fright Check Bonus")
	case MinorCostOfLivingIncrease:
		return i18n.Text("Includes a Minor Cost of Living Increase")
	case MajorCostOfLivingIncrease:
		return i18n.Text("Includes a Major Cost of Living Increase and Merchant Skill Penalty")
	default:
		return NoCRAdj.String()
	}
}

// Description returns a formatted description.
func (s SelfControlRollAdj) Description(cr advantage.SelfControlRoll) string {
	if cr == advantage.None {
		return ""
	}
	switch s {
	case NoCRAdj:
		return i18n.Text("None")
	case ActionPenalty:
		return fmt.Sprintf(i18n.Text("%d Action Penalty"), cr.Index()-len(advantage.AllSelfControlRolls))
	case ReactionPenalty:
		return fmt.Sprintf(i18n.Text("%d Reaction Penalty"), cr.Index()-len(advantage.AllSelfControlRolls))
	case FrightCheckPenalty:
		return fmt.Sprintf(i18n.Text("%d Fright Check Penalty"), cr.Index()-len(advantage.AllSelfControlRolls))
	case FrightCheckBonus:
		return fmt.Sprintf(i18n.Text("+%d Fright Check Bonus"), len(advantage.AllSelfControlRolls)-cr.Index())
	case MinorCostOfLivingIncrease:
		return fmt.Sprintf(i18n.Text("+%d%% Cost of Living Increase"), 5*(len(advantage.AllSelfControlRolls)-cr.Index()))
	case MajorCostOfLivingIncrease:
		return fmt.Sprintf(i18n.Text("+%d%% Cost of Living Increase"), 10*(1<<((cr.Index()-len(advantage.AllSelfControlRolls))-1)))
	default:
		return NoCRAdj.String()
	}
}

// Features returns the set of features to apply.
func (s SelfControlRollAdj) Features(cr advantage.SelfControlRoll) feature.Features {
	if s.EnsureValid() != MajorCostOfLivingIncrease {
		return nil
	}
	f := feature.NewSkillBonus()
	f.NameCriteria.Qualifier = "Merchant"
	f.Amount = fixed.F64d4FromInt64(int64(cr.Index() - len(advantage.AllSelfControlRolls)))
	return feature.Features{f}
}
