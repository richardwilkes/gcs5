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

import "github.com/richardwilkes/gcs/model/gurps/weapon"

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
	AdvantageList() []*Advantage
}

// EquipmentListProvider defines the method needed to access the equipment list data.
type EquipmentListProvider interface {
	CarriedEquipmentList() []*Equipment
	OtherEquipmentList() []*Equipment
}

// NoteListProvider defines the method needed to access the note list data.
type NoteListProvider interface {
	NoteList() []*Note
}

// SkillListProvider defines the method needed to access the skill list data.
type SkillListProvider interface {
	SkillList() []*Skill
}

// SpellListProvider defines the method needed to access the spell list data.
type SpellListProvider interface {
	SpellList() []*Spell
}

// AdvantageModifierListProvider defines the method needed to access the advantage modifier list data.
type AdvantageModifierListProvider interface {
	AdvantageModifierList() []*AdvantageModifier
}

// EquipmentModifierListProvider defines the method needed to access the equipment modifier list data.
type EquipmentModifierListProvider interface {
	EquipmentModifierList() []*EquipmentModifier
}

// ConditionalModifierListProvider defines the method needed to access the conditional modifier list data.
type ConditionalModifierListProvider interface {
	ConditionalModifiers() []*ConditionalModifier
}

// ReactionModifierListProvider defines the method needed to access the reaction modifier list data.
type ReactionModifierListProvider interface {
	Reactions() []*ConditionalModifier
}

// WeaponListProvider defines the method needed to access the weapon list data.
type WeaponListProvider interface {
	EquippedWeapons(weapon.Type) []*Weapon
}
