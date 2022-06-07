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

// NewNaturalAttacks creates a new "Natural Attacks" trait.
func NewNaturalAttacks(entity *Entity, parent *Trait) *Trait {
	a := NewTrait(entity, parent, false)
	a.Name = i18n.Text("Natural Attacks")
	a.PageRef = "B271"
	a.Weapons = []*Weapon{newBite(a), newPunch(a), newKick(a)}
	return a
}

func newBite(owner WeaponOwner) *Weapon {
	no := i18n.Text("No")
	bite := NewWeapon(owner, weapon.Melee)
	bite.Usage = i18n.Text("Bite")
	bite.Reach = "C"
	bite.Parry = no
	bite.Block = no
	bite.Defaults = []*SkillDefault{
		{
			DefaultType: gid.Dexterity,
		},
		{
			DefaultType: gid.Skill,
			Name:        "Brawling",
		},
	}
	bite.Damage.Type = "cr"
	bite.Damage.StrengthType = weapon.Thrust
	bite.Damage.Base = &dice.Dice{
		Sides:      6,
		Modifier:   -1,
		Multiplier: 1,
	}
	bite.Damage.ArmorDivisor = fxp.One
	bite.Damage.FragmentationArmorDivisor = fxp.One
	bite.Damage.Owner = bite
	return bite
}

func newPunch(owner WeaponOwner) *Weapon {
	punch := NewWeapon(owner, weapon.Melee)
	punch.Usage = i18n.Text("Punch")
	punch.Reach = "C"
	punch.Parry = "0"
	punch.Defaults = []*SkillDefault{
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
	}
	punch.Damage.Type = "cr"
	punch.Damage.StrengthType = weapon.Thrust
	punch.Damage.Base = &dice.Dice{
		Sides:      6,
		Modifier:   -1,
		Multiplier: 1,
	}
	punch.Damage.ArmorDivisor = fxp.One
	punch.Damage.FragmentationArmorDivisor = fxp.One
	punch.Damage.Owner = punch
	return punch
}

func newKick(owner WeaponOwner) *Weapon {
	kick := NewWeapon(owner, weapon.Melee)
	kick.Usage = i18n.Text("Kick")
	kick.Reach = "C,1"
	kick.Parry = i18n.Text("No")
	kick.Defaults = []*SkillDefault{
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
	}
	kick.Damage.Type = "cr"
	kick.Damage.StrengthType = weapon.Thrust
	kick.Damage.ArmorDivisor = fxp.One
	kick.Damage.FragmentationArmorDivisor = fxp.One
	kick.Damage.Owner = kick
	return kick
}
