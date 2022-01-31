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
	"github.com/richardwilkes/gcs/model/criteria"
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
)

// WeaponDamageBonus holds the data for an adjustment to weapon damage.
type WeaponDamageBonus struct {
	Bonus
	SelectionType          weapon.SelectionType `json:"selection_type"`
	NameCriteria           criteria.String      `json:"name"`
	SpecializationCriteria criteria.String      `json:"specialization"`
	RelativeLevelCriteria  criteria.Numeric     `json:"level"`
	CategoryCriteria       criteria.String      `json:"category"`
	Percent                bool                 `json:"percent,omitempty"`
}

// NewWeaponDamageBonus creates a new WeaponDamageBonus.
func NewWeaponDamageBonus() *WeaponDamageBonus {
	s := &WeaponDamageBonus{
		Bonus: Bonus{
			Feature: Feature{
				Type: feature.WeaponDamageBonus,
			},
			LeveledAmount: LeveledAmount{Amount: f64d4.One},
		},
		SelectionType: weapon.WithRequiredSkill,
		NameCriteria: criteria.String{
			Compare: criteria.Is,
		},
		SpecializationCriteria: criteria.String{
			Compare: criteria.Any,
		},
		RelativeLevelCriteria: criteria.Numeric{
			Compare: criteria.AtLeast,
		},
		CategoryCriteria: criteria.String{
			Compare: criteria.Any,
		},
	}
	s.Self = s
	return s
}

func (w *WeaponDamageBonus) featureMapKey() string {
	switch w.SelectionType {
	case weapon.WithRequiredSkill:
		return w.buildKey(WeaponNamedIDPrefix)
	case weapon.ThisWeapon:
		return ThisWeaponID
	case weapon.WithName:
		return w.buildKey(SkillNameID)
	default:
		jot.Fatal(1, "invalid selection type: ", w.SelectionType)
		return ""
	}
}

func (w *WeaponDamageBonus) buildKey(prefix string) string {
	if w.NameCriteria.Compare == criteria.Is &&
		(w.SpecializationCriteria.Compare == criteria.Any && w.CategoryCriteria.Compare == criteria.Any) {
		return prefix + "/" + w.NameCriteria.Qualifier
	}
	return prefix + "*"
}

func (w *WeaponDamageBonus) fillWithNameableKeys(nameables map[string]string) {
	ExtractNameables(w.SpecializationCriteria.Qualifier, nameables)
	if w.SelectionType != weapon.ThisWeapon {
		ExtractNameables(w.NameCriteria.Qualifier, nameables)
		ExtractNameables(w.SpecializationCriteria.Qualifier, nameables)
		ExtractNameables(w.CategoryCriteria.Qualifier, nameables)
	}
}

func (w *WeaponDamageBonus) applyNameableKeys(nameables map[string]string) {
	w.SpecializationCriteria.Qualifier = ApplyNameables(w.SpecializationCriteria.Qualifier, nameables)
	if w.SelectionType != weapon.ThisWeapon {
		w.NameCriteria.Qualifier = ApplyNameables(w.NameCriteria.Qualifier, nameables)
		w.SpecializationCriteria.Qualifier = ApplyNameables(w.SpecializationCriteria.Qualifier, nameables)
		w.CategoryCriteria.Qualifier = ApplyNameables(w.CategoryCriteria.Qualifier, nameables)
	}
}

func (w *WeaponDamageBonus) addToTooltip(buffer *xio.ByteBuffer) {
	if buffer != nil {
		buffer.WriteByte('\n')
		buffer.WriteString(w.ParentName())
		buffer.WriteString(" [")
		buffer.WriteString(w.LeveledAmount.Format(i18n.Text("die")))
		if w.Percent {
			buffer.WriteByte('%')
		}
		buffer.WriteByte(']')
	}
}
