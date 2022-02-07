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
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Description returns a formatted description.
func (enum SelfControlRollAdj) Description(cr advantage.SelfControlRoll) string {
	if cr == advantage.None {
		return ""
	}
	switch enum {
	case NoCRAdj:
		return enum.AltString()
	case ActionPenalty:
		return fmt.Sprintf(enum.AltString(), cr.Index()-len(advantage.AllSelfControlRolls))
	case ReactionPenalty:
		return fmt.Sprintf(enum.AltString(), cr.Index()-len(advantage.AllSelfControlRolls))
	case FrightCheckPenalty:
		return fmt.Sprintf(enum.AltString(), cr.Index()-len(advantage.AllSelfControlRolls))
	case FrightCheckBonus:
		return fmt.Sprintf(enum.AltString(), len(advantage.AllSelfControlRolls)-cr.Index())
	case MinorCostOfLivingIncrease:
		return fmt.Sprintf(enum.AltString(), 5*(len(advantage.AllSelfControlRolls)-cr.Index()))
	case MajorCostOfLivingIncrease:
		return fmt.Sprintf(enum.AltString(), 10*(1<<((cr.Index()-len(advantage.AllSelfControlRolls))-1)))
	default:
		return NoCRAdj.Description(cr)
	}
}

// Features returns the set of features to apply.
func (enum SelfControlRollAdj) Features(cr advantage.SelfControlRoll) feature.Features {
	if enum.EnsureValid() != MajorCostOfLivingIncrease {
		return nil
	}
	f := feature.NewSkillBonus()
	f.NameCriteria.Qualifier = "Merchant"
	f.Amount = fixed.F64d4FromInt(cr.Index() - len(advantage.AllSelfControlRolls))
	return feature.Features{f}
}
