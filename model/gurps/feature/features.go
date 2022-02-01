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

package feature

import (
	"encoding/json"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
)

// Features holds a list of features.
type Features []Feature

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *Features) UnmarshalJSON(data []byte) error {
	var s []*json.RawMessage
	if err := json.Unmarshal(data, &s); err != nil {
		return errs.Wrap(err)
	}
	*f = make([]Feature, len(s))
	for i, one := range s {
		var typeData struct {
			Type Type `json:"type"`
		}
		if err := json.Unmarshal(*one, &typeData); err != nil {
			return errs.Wrap(err)
		}
		var feature Feature
		switch typeData.Type {
		case AttributeBonusType:
			feature = &AttributeBonus{}
		case ConditionalModifierType:
			feature = &ConditionalModifier{}
		case ContainedWeightReductionType:
			feature = &ContainedWeightReduction{}
		case CostReductionType:
			feature = &CostReduction{}
		case DRBonusType:
			feature = &DRBonus{}
		case ReactionBonusType:
			feature = &ReactionBonus{}
		case SkillBonusType:
			feature = &SkillBonus{}
		case SkillPointBonusType:
			feature = &SkillPointBonus{}
		case SpellBonusType:
			feature = &SpellBonus{}
		case SpellPointBonusType:
			feature = &SpellPointBonus{}
		case WeaponDamageBonusType:
			feature = &WeaponDamageBonus{}
		default:
			return errs.Newf(i18n.Text("Unknown feature type: %s"), typeData.Type)
		}
		if err := json.Unmarshal(*one, &feature); err != nil {
			return errs.Wrap(err)
		}
		(*f)[i] = feature
	}
	return nil
}
