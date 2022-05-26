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
	"github.com/richardwilkes/gcs/model/gurps/weapon"
)

// ListProvider defines the methods needed to access list data.
type ListProvider interface {
	AdvantageListProvider
	EquipmentListProvider
	NoteListProvider
	SkillListProvider
	SpellListProvider
}

// AdvantageListProvider defines the method needed to access the advantage list data.
type AdvantageListProvider interface {
	EntityProvider
	AdvantageList() []*Advantage
	SetAdvantageList(list []*Advantage)
}

// EquipmentListProvider defines the method needed to access the equipment list data.
type EquipmentListProvider interface {
	EntityProvider
	CarriedEquipmentList() []*Equipment
	SetCarriedEquipmentList(list []*Equipment)
	OtherEquipmentList() []*Equipment
	SetOtherEquipmentList(list []*Equipment)
}

// NoteListProvider defines the method needed to access the note list data.
type NoteListProvider interface {
	EntityProvider
	NoteList() []*Note
	SetNoteList(list []*Note)
}

// SkillListProvider defines the method needed to access the skill list data.
type SkillListProvider interface {
	EntityProvider
	SkillList() []*Skill
	SetSkillList(list []*Skill)
}

// SpellListProvider defines the method needed to access the spell list data.
type SpellListProvider interface {
	EntityProvider
	SpellList() []*Spell
	SetSpellList(list []*Spell)
}

// AdvantageModifierListProvider defines the method needed to access the advantage modifier list data.
type AdvantageModifierListProvider interface {
	EntityProvider
	AdvantageModifierList() []*AdvantageModifier
	SetAdvantageModifierList(list []*AdvantageModifier)
}

// EquipmentModifierListProvider defines the method needed to access the equipment modifier list data.
type EquipmentModifierListProvider interface {
	EntityProvider
	EquipmentModifierList() []*EquipmentModifier
	SetEquipmentModifierList(list []*EquipmentModifier)
}

// ConditionalModifierListProvider defines the method needed to access the conditional modifier list data.
type ConditionalModifierListProvider interface {
	EntityProvider
	ConditionalModifiers() []*ConditionalModifier
}

// ReactionModifierListProvider defines the method needed to access the reaction modifier list data.
type ReactionModifierListProvider interface {
	EntityProvider
	Reactions() []*ConditionalModifier
}

// WeaponListProvider defines the method needed to access the weapon list data.
type WeaponListProvider interface {
	EntityProvider
	EquippedWeapons(weapon.Type) []*Weapon
}
