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

const (
	weaponTypeKey       = "type"
	weaponReachKey      = "reach"
	weaponParryKey      = "parry"
	weaponBlockKey      = "block"
	weaponAccuracyKey   = "accuracy"
	weaponRangeKey      = "range"
	weaponRateOfFireKey = "rate_of_fire"
	weaponShotsKey      = "shots"
	weaponBulkKey       = "bulk"
	weaponRecoilKey     = "recoil"
	weaponDefaultsKey   = "defaults"
	weaponStrengthKey   = "strength"
	weaponUsageKey      = "usage"
	weaponUsageNotesKey = "usage_notes"
)

// Weapon holds the stats for a weapon.
type Weapon struct {
	Type            WeaponType
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
	Damage          *WeaponDamage
	Defaults        []*SkillDefault
}
