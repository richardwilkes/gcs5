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

package menus

import (
	"github.com/richardwilkes/gcs/constants"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

func registerItemMenuActions() {
	settings.RegisterKeyBinding("new.adq", NewTrait)
	settings.RegisterKeyBinding("new.adq.container", NewTraitContainer)
	settings.RegisterKeyBinding("new.adm", NewTraitModifier)
	settings.RegisterKeyBinding("new.adm.container", NewTraitContainerModifier)
	settings.RegisterKeyBinding("add.natural.attacks", AddNaturalAttacks)
	settings.RegisterKeyBinding("new.skl", NewSkill)
	settings.RegisterKeyBinding("new.skl.container", NewSkillContainer)
	settings.RegisterKeyBinding("new.skl.technique", NewTechnique)
	settings.RegisterKeyBinding("new.spl", NewSpell)
	settings.RegisterKeyBinding("new.spl.container", NewSpellContainer)
	settings.RegisterKeyBinding("new.spl.ritual", NewRitualMagicSpell)
	settings.RegisterKeyBinding("new.eqp", NewCarriedEquipment)
	settings.RegisterKeyBinding("new.eqp.container", NewCarriedEquipmentContainer)
	settings.RegisterKeyBinding("new.eqp.other", NewOtherEquipment)
	settings.RegisterKeyBinding("new.eqp.other.container", NewOtherEquipmentContainer)
	settings.RegisterKeyBinding("new.eqm", NewEquipmentModifier)
	settings.RegisterKeyBinding("new.eqm.container", NewEquipmentContainerModifier)
	settings.RegisterKeyBinding("new.not", NewNote)
	settings.RegisterKeyBinding("new.not.container", NewNoteContainer)
	settings.RegisterKeyBinding("new.melee", NewMeleeWeapon)
	settings.RegisterKeyBinding("new.ranged", NewRangedWeapon)
	settings.RegisterKeyBinding("pageref.open.first", OpenOnePageReference)
	settings.RegisterKeyBinding("pageref.open.all", OpenEachPageReference)
}

func createItemMenu(f unison.MenuFactory) unison.Menu {
	m := f.NewMenu(constants.ItemMenuID, i18n.Text("Item"), nil)

	m.InsertItem(-1, NewTrait.NewMenuItem(f))
	m.InsertItem(-1, NewTraitContainer.NewMenuItem(f))
	m.InsertItem(-1, NewTraitModifier.NewMenuItem(f))
	m.InsertItem(-1, NewTraitContainerModifier.NewMenuItem(f))
	m.InsertItem(-1, AddNaturalAttacks.NewMenuItem(f))

	m.InsertSeparator(-1, false)
	m.InsertItem(-1, NewSkill.NewMenuItem(f))
	m.InsertItem(-1, NewSkillContainer.NewMenuItem(f))
	m.InsertItem(-1, NewTechnique.NewMenuItem(f))

	m.InsertSeparator(-1, false)
	m.InsertItem(-1, NewSpell.NewMenuItem(f))
	m.InsertItem(-1, NewSpellContainer.NewMenuItem(f))
	m.InsertItem(-1, NewRitualMagicSpell.NewMenuItem(f))

	m.InsertSeparator(-1, false)
	m.InsertItem(-1, NewCarriedEquipment.NewMenuItem(f))
	m.InsertItem(-1, NewCarriedEquipmentContainer.NewMenuItem(f))
	m.InsertItem(-1, NewOtherEquipment.NewMenuItem(f))
	m.InsertItem(-1, NewOtherEquipmentContainer.NewMenuItem(f))
	m.InsertItem(-1, NewEquipmentModifier.NewMenuItem(f))
	m.InsertItem(-1, NewEquipmentContainerModifier.NewMenuItem(f))

	m.InsertSeparator(-1, false)
	m.InsertItem(-1, NewNote.NewMenuItem(f))
	m.InsertItem(-1, NewNoteContainer.NewMenuItem(f))

	m.InsertSeparator(-1, false)
	m.InsertItem(-1, NewMeleeWeapon.NewMenuItem(f))
	m.InsertItem(-1, NewRangedWeapon.NewMenuItem(f))

	m.InsertSeparator(-1, false)
	m.InsertItem(-1, OpenOnePageReference.NewMenuItem(f))
	m.InsertItem(-1, OpenEachPageReference.NewMenuItem(f))
	return m
}

// NewTrait creates a new trait.
var NewTrait = &unison.Action{
	ID:              constants.NewTraitItemID,
	Title:           i18n.Text("New Trait"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyD, Modifiers: unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewTraitContainer creates a new trait container.
var NewTraitContainer = &unison.Action{
	ID:              constants.NewTraitContainerItemID,
	Title:           i18n.Text("New Trait Container"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyD, Modifiers: unison.ShiftModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewTraitModifier creates a new trait modifier.
var NewTraitModifier = &unison.Action{
	ID:              constants.NewTraitModifierItemID,
	Title:           i18n.Text("New Trait Modifier"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyM, Modifiers: unison.OptionModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewTraitContainerModifier creates a new trait container modifier.
var NewTraitContainerModifier = &unison.Action{
	ID:              constants.NewTraitContainerModifierItemID,
	Title:           i18n.Text("New Trait Modifier Container"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyM, Modifiers: unison.ShiftModifier | unison.OptionModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// AddNaturalAttacks creates the natural attacks.
var AddNaturalAttacks = &unison.Action{
	ID:              constants.AddNaturalAttacksItemID,
	Title:           i18n.Text("Add Natural Attacks"),
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewSkill creates a new skill.
var NewSkill = &unison.Action{
	ID:              constants.NewSkillItemID,
	Title:           i18n.Text("New Skill"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyK, Modifiers: unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewSkillContainer creates a new skill container.
var NewSkillContainer = &unison.Action{
	ID:              constants.NewSkillContainerItemID,
	Title:           i18n.Text("New Skill Container"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyK, Modifiers: unison.ShiftModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewTechnique creates a new technique.
var NewTechnique = &unison.Action{
	ID:              constants.NewTechniqueItemID,
	Title:           i18n.Text("New Technique"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyT, Modifiers: unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewSpell creates a new spell.
var NewSpell = &unison.Action{
	ID:              constants.NewSpellItemID,
	Title:           i18n.Text("New Spell"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyB, Modifiers: unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewSpellContainer creates a new spell container.
var NewSpellContainer = &unison.Action{
	ID:              constants.NewSpellContainerItemID,
	Title:           i18n.Text("New Spell Container"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyB, Modifiers: unison.ShiftModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewRitualMagicSpell creates a new ritual magic spell.
var NewRitualMagicSpell = &unison.Action{
	ID:              constants.NewRitualMagicSpellItemID,
	Title:           i18n.Text("New Ritual Magic Spell"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyB, Modifiers: unison.ShiftModifier | unison.OptionModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewCarriedEquipment creates a new equipment item.
var NewCarriedEquipment = &unison.Action{
	ID:              constants.NewCarriedEquipmentItemID,
	Title:           i18n.Text("New Carried Equipment"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyE, Modifiers: unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewCarriedEquipmentContainer creates a new equipment container.
var NewCarriedEquipmentContainer = &unison.Action{
	ID:              constants.NewCarriedEquipmentContainerItemID,
	Title:           i18n.Text("New Carried Equipment Container"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyE, Modifiers: unison.ShiftModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewOtherEquipment creates a new equipment item.
var NewOtherEquipment = &unison.Action{
	ID:              constants.NewOtherEquipmentItemID,
	Title:           i18n.Text("New Other Equipment"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyE, Modifiers: unison.OptionModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewOtherEquipmentContainer creates a new equipment container.
var NewOtherEquipmentContainer = &unison.Action{
	ID:              constants.NewOtherEquipmentContainerItemID,
	Title:           i18n.Text("New Other Equipment Container"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyE, Modifiers: unison.ShiftModifier | unison.OptionModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewEquipmentModifier creates a new equipment modifier.
var NewEquipmentModifier = &unison.Action{
	ID:              constants.NewEquipmentModifierItemID,
	Title:           i18n.Text("New Equipment Modifier"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyF, Modifiers: unison.OptionModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewEquipmentContainerModifier creates a new equipment container modifier.
var NewEquipmentContainerModifier = &unison.Action{
	ID:              constants.NewEquipmentContainerModifierItemID,
	Title:           i18n.Text("New Equipment Modifier Container"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyF, Modifiers: unison.ShiftModifier | unison.OptionModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewNote creates a new note.
var NewNote = &unison.Action{
	ID:              constants.NewNoteItemID,
	Title:           i18n.Text("New Note"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyN, Modifiers: unison.ShiftModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewNoteContainer creates a new note container.
var NewNoteContainer = &unison.Action{
	ID:              constants.NewNoteContainerItemID,
	Title:           i18n.Text("New Note Container"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyN, Modifiers: unison.ShiftModifier | unison.OptionModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewMeleeWeapon creates a new melee weapon.
var NewMeleeWeapon = &unison.Action{
	ID:              constants.NewMeleeWeaponItemID,
	Title:           i18n.Text("New Melee Weapon"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyM, Modifiers: unison.ShiftModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// NewRangedWeapon creates a new ranged weapon.
var NewRangedWeapon = &unison.Action{
	ID:              constants.NewRangedWeaponItemID,
	Title:           i18n.Text("New Ranged Weapon"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyR, Modifiers: unison.ShiftModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// OpenOnePageReference opens the first page reference for each selected item.
var OpenOnePageReference = &unison.Action{
	ID:              constants.OpenOnePageReferenceItemID,
	Title:           i18n.Text("Open Page Reference"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyG, Modifiers: unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}

// OpenEachPageReference opens each page reference associated with the selected items.
var OpenEachPageReference = &unison.Action{
	ID:              constants.OpenEachPageReferenceItemID,
	Title:           i18n.Text("Open Each Page Reference"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyG, Modifiers: unison.ShiftModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: unison.RouteActionToFocusEnabledFunc,
	ExecuteCallback: unison.RouteActionToFocusExecuteFunc,
}
