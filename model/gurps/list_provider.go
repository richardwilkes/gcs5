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

// ListProvider defines the methods needed to access list data.
type ListProvider interface {
	AdvantageList() []*Advantage
	CarriedEquipmentList() []*Equipment
	OtherEquipmentList() []*Equipment
	NoteList() []*Note
	SkillListProvider
	SpellListProvider
}

// SkillListProvider defines the method needed to access the skill list data.
type SkillListProvider interface {
	SkillList() []*Skill
}

// SpellListProvider defines the method needed to access the spell list data.
type SpellListProvider interface {
	SpellList() []*Spell
}
