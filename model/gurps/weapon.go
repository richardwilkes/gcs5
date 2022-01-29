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
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	weaponTypeKey            = "type"
	weaponDamageKey          = "damage"
	weaponMinimumStrengthKey = "strength"
	weaponReachKey           = "reach"
	weaponParryKey           = "parry"
	weaponBlockKey           = "block"
	weaponAccuracyKey        = "accuracy"
	weaponRangeKey           = "range"
	weaponRateOfFireKey      = "rate_of_fire"
	weaponShotsKey           = "shots"
	weaponBulkKey            = "bulk"
	weaponRecoilKey          = "recoil"
	weaponUsageKey           = "usage"
	weaponUsageNotesKey      = "usage_notes"
	weaponCalcLevelKey       = "level"
	weaponCalcParryKey       = "parry"
	weaponCalcBlockKey       = "block"
	weaponCalcDamageKey      = "damage"
	weaponCalcRangeKey       = "range"
)

// Weapon holds the stats for a weapon.
type Weapon struct {
	Type            weapon.Type
	Damage          *WeaponDamage
	MinimumStrength string
	Usage           string
	UsageNotes      string
	Reach           string
	Parry           string
	Block           string
	Accuracy        string
	Range           string
	RateOfFire      string
	Shots           string
	Bulk            string
	Recoil          string
	Defaults        []*SkillDefault
}

// NewWeaponFromJSON creates a new Weapon from a JSON object.
func NewWeaponFromJSON(data map[string]interface{}) *Weapon {
	w := &Weapon{
		Type:            weapon.TypeFromKey(encoding.String(data[weaponTypeKey])),
		MinimumStrength: encoding.String(data[weaponMinimumStrengthKey]),
		Usage:           encoding.String(data[weaponUsageKey]),
		UsageNotes:      encoding.String(data[weaponUsageNotesKey]),
	}
	w.Damage = NewWeaponDamageFromJSON(w, encoding.Object(data[weaponDamageKey]))
	switch w.Type {
	case weapon.Melee:
		w.Reach = encoding.String(data[weaponReachKey])
		w.Parry = encoding.String(data[weaponParryKey])
		w.Block = encoding.String(data[weaponBlockKey])
	case weapon.Ranged:
		w.Accuracy = encoding.String(data[weaponAccuracyKey])
		w.Range = encoding.String(data[weaponRangeKey])
		w.RateOfFire = encoding.String(data[weaponRateOfFireKey])
		w.Shots = encoding.String(data[weaponShotsKey])
		w.Bulk = encoding.String(data[weaponBulkKey])
		w.Recoil = encoding.String(data[weaponRecoilKey])
	}
	w.Defaults = SkillDefaultsListFromJSON(data)
	return w
}

// ToJSON emits this object as JSON.
func (w *Weapon) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(weaponTypeKey, w.Type.Key(), false, false)
	encoding.ToKeyedJSON(w.Damage, weaponDamageKey, encoder)
	encoder.KeyedString(weaponMinimumStrengthKey, w.MinimumStrength, true, true)
	encoder.KeyedString(weaponUsageKey, w.Usage, true, true)
	encoder.KeyedString(weaponUsageNotesKey, w.UsageNotes, true, true)
	switch w.Type {
	case weapon.Melee:
		encoder.KeyedString(weaponReachKey, w.Reach, true, true)
		encoder.KeyedString(weaponParryKey, w.Parry, true, true)
		encoder.KeyedString(weaponBlockKey, w.Block, true, true)
	case weapon.Ranged:
		encoder.KeyedString(weaponAccuracyKey, w.Accuracy, true, true)
		encoder.KeyedString(weaponRangeKey, w.Range, true, true)
		encoder.KeyedString(weaponRateOfFireKey, w.RateOfFire, true, true)
		encoder.KeyedString(weaponShotsKey, w.Shots, true, true)
		encoder.KeyedString(weaponBulkKey, w.Bulk, true, true)
		encoder.KeyedString(weaponRecoilKey, w.Recoil, true, true)
	}
	SkillDefaultsListToJSON(w.Defaults, encoder)
	// Emit the calculated values for third parties
	encoder.Key(calcKey)
	encoder.StartObject()
	encoder.KeyedNumber(weaponCalcLevelKey, fixed.F64d4FromInt64(int64(xmath.MaxInt(w.SkillLevel(), 0))), true)
	switch w.Type {
	case weapon.Melee:
		encoder.KeyedString(weaponCalcParryKey, w.ResolvedParry(nil), true, true)
		encoder.KeyedString(weaponCalcBlockKey, w.ResolvedBlock(nil), true, true)
	case weapon.Ranged:
		encoder.KeyedString(weaponCalcRangeKey, w.ResolvedRange(), true, true)
	}
	encoder.KeyedString(weaponCalcDamageKey, w.Damage.ResolvedDamage(nil), true, true)
	encoder.EndObject()
	encoder.EndObject()
}

// SkillLevel returns the resolved skill level.
func (w *Weapon) SkillLevel() int {
	// TODO: Implement
	return 0
}

// ResolvedParry returns the resolved parry level.
func (w *Weapon) ResolvedParry(tooltip *xio.ByteBuffer) string {
	// TODO: Implement
	return ""
}

// ResolvedBlock returns the resolved block level.
func (w *Weapon) ResolvedBlock(tooltip *xio.ByteBuffer) string {
	// TODO: Implement
	return ""
}

// ResolvedRange returns the range, fully resolved for the user's ST, if possible.
func (w *Weapon) ResolvedRange() string {
	// TODO: Implement
	return ""
}

// ResolvedMinimumStrength returns the resolved minimum strength required to use this weapon, or 0 if there is none.
func (w *Weapon) ResolvedMinimumStrength() int {
	started := false
	value := 0
	for _, ch := range w.MinimumStrength {
		if ch >= '0' && ch <= '9' {
			value *= 10
			value += int(ch - '0')
			started = true
		} else if started {
			break
		}
	}
	return value
}

// FillWithNameableKeys adds any nameable keys found in this Weapon to the provided map.
func (w *Weapon) FillWithNameableKeys(nameables map[string]string) {
	for _, one := range w.Defaults {
		one.FillWithNameableKeys(nameables)
	}
}

// ApplyNameableKeys replaces any nameable keys found in this Weapon with the corresponding values in the provided map.
func (w *Weapon) ApplyNameableKeys(nameables map[string]string) {
	for _, one := range w.Defaults {
		one.ApplyNameableKeys(nameables)
	}
}
