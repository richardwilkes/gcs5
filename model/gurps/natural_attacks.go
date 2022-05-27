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
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/i18n"
)

// NewNaturalAttacks creates a new "Natural Attacks" advantage.
func NewNaturalAttacks(entity *Entity, parent *Advantage) *Advantage {
	a := NewAdvantage(entity, parent, false)
	a.Name = i18n.Text("Natural Attacks")
	a.PageRef = "B271"
	a.Weapons = []*Weapon{newBite(a), newPunch(a), newKick(a)}
	return a
}

func newBite(owner WeaponOwner) *Weapon {
	no := i18n.Text("No")
	bite := &Weapon{
		WeaponData: WeaponData{
			Type:            weapon.Melee,
			MinimumStrength: "",
			Usage:           i18n.Text("Bite"),
			Reach:           "C",
			Parry:           no,
			Block:           no,
			Defaults: []*SkillDefault{
				{
					DefaultType: gid.Dexterity,
				},
				{
					DefaultType: gid.Skill,
					Name:        "Brawling",
				},
			},
		},
		Owner: owner,
	}
	bite.Damage = &WeaponDamage{
		WeaponDamageData: WeaponDamageData{
			Type:         "cr",
			StrengthType: weapon.Thrust,
			Base: &dice.Dice{
				Sides:      6,
				Modifier:   -1,
				Multiplier: 1,
			},
			ArmorDivisor:              fxp.One,
			FragmentationArmorDivisor: fxp.One,
		},
		Owner: bite,
	}
	return bite
}

func newPunch(owner WeaponOwner) *Weapon {
	punch := &Weapon{
		WeaponData: WeaponData{
			Type:            weapon.Melee,
			MinimumStrength: "",
			Usage:           i18n.Text("Punch"),
			Reach:           "C",
			Parry:           "0",
			Defaults: []*SkillDefault{
				{
					DefaultType: gid.Dexterity,
				},
				{
					DefaultType: gid.Skill,
					Name:        "Boxing",
				},
				{
					DefaultType: gid.Skill,
					Name:        "Brawling",
				},
				{
					DefaultType: gid.Skill,
					Name:        "Karate",
				},
			},
		},
		Owner: owner,
	}
	punch.Damage = &WeaponDamage{
		WeaponDamageData: WeaponDamageData{
			Type:         "cr",
			StrengthType: weapon.Thrust,
			Base: &dice.Dice{
				Sides:      6,
				Modifier:   -1,
				Multiplier: 1,
			},
			ArmorDivisor:              fxp.One,
			FragmentationArmorDivisor: fxp.One,
		},
		Owner: punch,
	}
	return punch
}

func newKick(owner WeaponOwner) *Weapon {
	punch := &Weapon{
		WeaponData: WeaponData{
			Type:            weapon.Melee,
			MinimumStrength: "",
			Usage:           i18n.Text("Kick"),
			Reach:           "C,1",
			Parry:           i18n.Text("No"),
			Defaults: []*SkillDefault{
				{
					DefaultType: gid.Dexterity,
					Modifier:    -fxp.Two,
				},
				{
					DefaultType: gid.Skill,
					Name:        "Brawling",
					Modifier:    -fxp.Two,
				},
				{
					DefaultType: gid.Skill,
					Name:        "Kicking",
				},
				{
					DefaultType: gid.Skill,
					Name:        "Karate",
					Modifier:    -fxp.Two,
				},
			},
		},
		Owner: owner,
	}
	punch.Damage = &WeaponDamage{
		WeaponDamageData: WeaponDamageData{
			Type:                      "cr",
			StrengthType:              weapon.Thrust,
			ArmorDivisor:              fxp.One,
			FragmentationArmorDivisor: fxp.One,
		},
		Owner: punch,
	}
	return punch
}
