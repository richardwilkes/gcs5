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
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	weaponDamageTypeKey                      = "type"
	weaponDamageStrengthTypeKey              = "st"
	weaponDamageBaseKey                      = "base"
	weaponDamageArmorDivisorKey              = "armor_divisor"
	weaponDamageFragmentationKey             = "fragmentation"
	weaponDamageFragmentationArmorDivisorKey = "fragmentation_armor_divisor"
	weaponDamageFragmentationTypeKey         = "fragmentation_type"
	weaponDamageModifierPerDieKey            = "modifier_per_die"
)

// WeaponDamage holds the damage information for a weapon.
type WeaponDamage struct {
	Owner                     *Weapon
	Type                      string
	StrengthType              WeaponSTDamage
	Base                      *dice.Dice
	ArmorDivisor              fixed.F64d4
	Fragmentation             *dice.Dice
	FragmentationArmorDivisor fixed.F64d4
	FragmentationType         string
	ModifierPerDie            fixed.F64d4
}

// NewWeaponDamageFromJSON creates a new WeaponDamage from a JSON object.
func NewWeaponDamageFromJSON(owner *Weapon, data map[string]interface{}) *WeaponDamage {
	w := &WeaponDamage{
		Owner:          owner,
		Type:           encoding.String(data[weaponDamageTypeKey]),
		StrengthType:   WeaponSTDamageFromKey(encoding.String(data[weaponDamageStrengthTypeKey])),
		ArmorDivisor:   encoding.Number(data[weaponDamageArmorDivisorKey]),
		ModifierPerDie: encoding.Number(data[weaponDamageModifierPerDieKey]),
	}
	if w.ArmorDivisor == 0 {
		w.ArmorDivisor = f64d4.One
	}
	if v, exists := data[weaponDamageBaseKey]; exists {
		w.Base = dice.New(encoding.String(v))
	}
	if v, exists := data[weaponDamageFragmentationKey]; exists {
		w.Fragmentation = dice.New(encoding.String(v))
		w.FragmentationType = encoding.String(data[weaponDamageFragmentationTypeKey])
		w.FragmentationArmorDivisor = encoding.Number(data[weaponDamageFragmentationArmorDivisorKey])
		if w.FragmentationArmorDivisor == 0 {
			w.FragmentationArmorDivisor = f64d4.One
		}
	}
	return w
}

// ToJSON emits this object as JSON.
func (w *WeaponDamage) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(weaponDamageTypeKey, w.Type, true, true)
	if w.StrengthType != NoSTBasedDamage {
		encoder.KeyedString(weaponDamageStrengthTypeKey, w.StrengthType.Key(), false, false)
	}
	if w.Base != nil {
		if s := w.Base.String(); s != "0" {
			encoder.KeyedString(weaponDamageBaseKey, s, false, false)
		}
	}
	if w.ArmorDivisor != f64d4.One {
		encoder.KeyedNumber(weaponDamageArmorDivisorKey, w.ArmorDivisor, true)
	}
	if w.Fragmentation != nil {
		if s := w.Fragmentation.String(); s != "0" {
			encoder.KeyedString(weaponDamageFragmentationKey, s, false, false)
			if w.FragmentationArmorDivisor != f64d4.One {
				encoder.KeyedNumber(weaponDamageFragmentationArmorDivisorKey, w.FragmentationArmorDivisor, true)
			}
			encoder.KeyedString(weaponDamageFragmentationTypeKey, w.FragmentationType, true, true)
		}
	}
	encoder.KeyedNumber(weaponDamageModifierPerDieKey, w.ModifierPerDie, true)
	encoder.EndObject()
}

// DamageTooltip returns a formatted tooltip for the damage.
func (w *WeaponDamage) DamageTooltip() string {
	var tooltip xio.ByteBuffer
	w.ResolvedDamage(&tooltip)
	if tooltip.Len() == 0 {
		return i18n.Text("No additional modifiers")
	}
	return i18n.Text("Includes modifiers from") + tooltip.String()
}

// ResolvedDamage returns the damage, fully resolved for the user's sw or thr, if possible.
func (w *WeaponDamage) ResolvedDamage(tooltip *xio.ByteBuffer) string {
	// TODO: Implement
	return ""
}
