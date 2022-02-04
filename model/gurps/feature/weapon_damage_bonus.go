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
	"fmt"

	"github.com/richardwilkes/gcs/model/criteria"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
)

const (
	// ThisWeaponID holds the ID for "this weapon".
	ThisWeaponID = "\u0001"
	// WeaponNamedIDPrefix the prefix for "weapon named" IDs.
	WeaponNamedIDPrefix = "weapon_named."
)

var _ Bonus = &WeaponDamageBonus{}

// WeaponDamageBonus holds the data for an adjustment to weapon damage.
type WeaponDamageBonus struct {
	Type                   Type                 `json:"type"`
	Parent                 fmt.Stringer         `json:"-"`
	SelectionType          weapon.SelectionType `json:"selection_type"`
	NameCriteria           criteria.String      `json:"name,omitempty"`
	SpecializationCriteria criteria.String      `json:"specialization"`
	RelativeLevelCriteria  criteria.Numeric     `json:"level,omitempty"`
	CategoryCriteria       criteria.String      `json:"category,omitempty"`
	Percent                bool                 `json:"percent,omitempty"`
	LeveledAmount
}

// NewWeaponDamageBonus creates a new WeaponDamageBonus.
func NewWeaponDamageBonus() *WeaponDamageBonus {
	return &WeaponDamageBonus{
		Type:          WeaponDamageBonusType,
		SelectionType: weapon.WithRequiredSkill,
		NameCriteria: criteria.String{
			StringData: criteria.StringData{
				Compare: criteria.Is,
			},
		},
		SpecializationCriteria: criteria.String{
			StringData: criteria.StringData{
				Compare: criteria.Any,
			},
		},
		RelativeLevelCriteria: criteria.Numeric{
			NumericData: criteria.NumericData{
				Compare: criteria.AtLeast,
			},
		},
		CategoryCriteria: criteria.String{
			StringData: criteria.StringData{
				Compare: criteria.Any,
			},
		},
		LeveledAmount: LeveledAmount{Amount: fxp.One},
	}
}

// FeatureMapKey implements Feature.
func (w *WeaponDamageBonus) FeatureMapKey() string {
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

// FillWithNameableKeys implements Feature.
func (w *WeaponDamageBonus) FillWithNameableKeys(m map[string]string) {
	nameables.Extract(w.SpecializationCriteria.Qualifier, m)
	if w.SelectionType != weapon.ThisWeapon {
		nameables.Extract(w.NameCriteria.Qualifier, m)
		nameables.Extract(w.SpecializationCriteria.Qualifier, m)
		nameables.Extract(w.CategoryCriteria.Qualifier, m)
	}
}

// ApplyNameableKeys implements Feature.
func (w *WeaponDamageBonus) ApplyNameableKeys(m map[string]string) {
	w.SpecializationCriteria.Qualifier = nameables.Apply(w.SpecializationCriteria.Qualifier, m)
	if w.SelectionType != weapon.ThisWeapon {
		w.NameCriteria.Qualifier = nameables.Apply(w.NameCriteria.Qualifier, m)
		w.SpecializationCriteria.Qualifier = nameables.Apply(w.SpecializationCriteria.Qualifier, m)
		w.CategoryCriteria.Qualifier = nameables.Apply(w.CategoryCriteria.Qualifier, m)
	}
}

// AddToTooltip implements Bonus.
func (w *WeaponDamageBonus) AddToTooltip(buffer *xio.ByteBuffer) {
	if buffer != nil {
		buffer.WriteByte('\n')
		buffer.WriteString(parentName(w.Parent))
		buffer.WriteString(" [")
		buffer.WriteString(w.LeveledAmount.Format(i18n.Text("die")))
		if w.Percent {
			buffer.WriteByte('%')
		}
		buffer.WriteByte(']')
	}
}
