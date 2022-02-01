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
	"encoding/json"

	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// WeaponData holds the Weapon data that is written to disk.
type WeaponData struct {
	Type            weapon.Type     `json:"type"`
	Damage          *WeaponDamage   `json:"damage,omitempty"`
	MinimumStrength string          `json:"strength,omitempty"`
	Usage           string          `json:"usage,omitempty"`
	UsageNotes      string          `json:"usage_notes,omitempty"`
	Reach           string          `json:"reach,omitempty"`
	Parry           string          `json:"parry,omitempty"`
	Block           string          `json:"block,omitempty"`
	Accuracy        string          `json:"accuracy,omitempty"`
	Range           string          `json:"range,omitempty"`
	RateOfFire      string          `json:"rate_of_fire,omitempty"`
	Shots           string          `json:"shots,omitempty"`
	Bulk            string          `json:"bulk,omitempty"`
	Recoil          string          `json:"recoil,omitempty"`
	Defaults        []*SkillDefault `json:"defaults,omitempty"`
}

// Weapon holds the stats for a weapon.
type Weapon struct {
	WeaponData
}

// MarshalJSON implements json.Marshaler.
func (w *Weapon) MarshalJSON() ([]byte, error) {
	type calc struct {
		Level  fixed.F64d4 `json:"level,omitempty"`
		Parry  string      `json:"parry,omitempty"`
		Block  string      `json:"block,omitempty"`
		Range  string      `json:"range,omitempty"`
		Damage string      `json:"damage,omitempty"`
	}
	data := struct {
		WeaponData
		Calc calc `json:"calc"`
	}{
		WeaponData: w.WeaponData,
		Calc: calc{
			Level:  fixed.F64d4FromInt64(int64(xmath.MaxInt(w.SkillLevel(), 0))),
			Damage: w.Damage.ResolvedDamage(nil),
		},
	}
	if w.Type == weapon.Melee {
		data.Calc.Parry = w.ResolvedParry(nil)
		data.Calc.Block = w.ResolvedBlock(nil)
	} else {
		data.Calc.Range = w.ResolvedRange()
	}
	return json.Marshal(&data)
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
func (w *Weapon) FillWithNameableKeys(m map[string]string) {
	for _, one := range w.Defaults {
		one.FillWithNameableKeys(m)
	}
}

// ApplyNameableKeys replaces any nameable keys found in this Weapon with the corresponding values in the provided map.
func (w *Weapon) ApplyNameableKeys(m map[string]string) {
	for _, one := range w.Defaults {
		one.ApplyNameableKeys(m)
	}
}
