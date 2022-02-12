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
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

func createItemMenu(f unison.MenuFactory) unison.Menu {
	m := f.NewMenu(ItemMenuID, i18n.Text("Item"), nil)
	m.InsertMenu(-1, createAdvantagesMenu(f))
	m.InsertMenu(-1, createSkillsMenu(f))
	m.InsertMenu(-1, createSpellsMenu(f))
	m.InsertMenu(-1, createEquipmentMenu(f))
	m.InsertMenu(-1, createNotesMenu(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, OpenEditor.NewMenuItem(f))
	m.InsertItem(-1, CopyToSheet.NewMenuItem(f))
	m.InsertItem(-1, CopyToTemplate.NewMenuItem(f))
	m.InsertItem(-1, ApplyTemplate.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, OpenOnePageReference.NewMenuItem(f))
	m.InsertItem(-1, OpenEachPageReference.NewMenuItem(f))
	return m
}

func createAdvantagesMenu(f unison.MenuFactory) unison.Menu {
	m := f.NewMenu(AdvantagesMenuID, i18n.Text("Advantages"), nil)
	m.InsertItem(-1, NewAdvantage.NewMenuItem(f))
	m.InsertItem(-1, NewAdvantageContainer.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, NewAdvantageModifier.NewMenuItem(f))
	m.InsertItem(-1, NewAdvantageContainerModifier.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, AddNaturalAttacksAdvantage.NewMenuItem(f))
	return m
}

func createSkillsMenu(f unison.MenuFactory) unison.Menu {
	m := f.NewMenu(SkillsMenuID, i18n.Text("Skills"), nil)
	m.InsertItem(-1, NewSkill.NewMenuItem(f))
	m.InsertItem(-1, NewSkillContainer.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, NewTechnique.NewMenuItem(f))
	return m
}

func createSpellsMenu(f unison.MenuFactory) unison.Menu {
	m := f.NewMenu(SkillsMenuID, i18n.Text("Spells"), nil)
	m.InsertItem(-1, NewSpell.NewMenuItem(f))
	m.InsertItem(-1, NewSpellContainer.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, NewRitualMagicSpell.NewMenuItem(f))
	return m
}

func createEquipmentMenu(f unison.MenuFactory) unison.Menu {
	m := f.NewMenu(EquipmentMenuID, i18n.Text("Equipment"), nil)
	m.InsertItem(-1, NewCarriedEquipment.NewMenuItem(f))
	m.InsertItem(-1, NewCarriedEquipmentContainer.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, NewOtherEquipment.NewMenuItem(f))
	m.InsertItem(-1, NewOtherEquipmentContainer.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, NewEquipmentModifier.NewMenuItem(f))
	m.InsertItem(-1, NewEquipmentContainerModifier.NewMenuItem(f))
	return m
}

func createNotesMenu(f unison.MenuFactory) unison.Menu {
	m := f.NewMenu(NotesMenuID, i18n.Text("Notes"), nil)
	m.InsertItem(-1, NewNote.NewMenuItem(f))
	m.InsertItem(-1, NewNoteContainer.NewMenuItem(f))
	return m
}

// NewAdvantage creates a new advantage.
var NewAdvantage = &unison.Action{
	ID:              NewAdvantageItemID,
	Title:           i18n.Text("New Advantage"),
	HotKey:          unison.KeyD,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// NewAdvantageContainer creates a new advantage container.
var NewAdvantageContainer = &unison.Action{
	ID:              NewAdvantageContainerItemID,
	Title:           i18n.Text("New Advantage Container"),
	HotKey:          unison.KeyD,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// NewAdvantageModifier creates a new advantage modifier.
var NewAdvantageModifier = &unison.Action{
	ID:              NewAdvantageModifierItemID,
	Title:           i18n.Text("New Advantage Modifier"),
	HotKey:          unison.KeyM,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.OptionModifier,
	ExecuteCallback: unimplemented,
}

// NewAdvantageContainerModifier creates a new advantage container modifier.
var NewAdvantageContainerModifier = &unison.Action{
	ID:              NewAdvantageContainerModifierItemID,
	Title:           i18n.Text("New Advantage Container Modifier"),
	HotKey:          unison.KeyM,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.OptionModifier | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// AddNaturalAttacksAdvantage creates the natural attacks advantage.
var AddNaturalAttacksAdvantage = &unison.Action{
	ID:              AddNaturalAttacksAdvantageItemID,
	Title:           i18n.Text("Add Natural Attacks Advantage"),
	ExecuteCallback: unimplemented,
}

// NewSkill creates a new skill.
var NewSkill = &unison.Action{
	ID:              NewSkillItemID,
	Title:           i18n.Text("New Skill"),
	HotKey:          unison.KeyK,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// NewSkillContainer creates a new skill container.
var NewSkillContainer = &unison.Action{
	ID:              NewSkillContainerItemID,
	Title:           i18n.Text("New Skill Container"),
	HotKey:          unison.KeyK,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// NewTechnique creates a new technique.
var NewTechnique = &unison.Action{
	ID:              NewTechniqueItemID,
	Title:           i18n.Text("New Technique"),
	HotKey:          unison.KeyT,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// NewSpell creates a new spell.
var NewSpell = &unison.Action{
	ID:              NewSpellItemID,
	Title:           i18n.Text("New Spell"),
	HotKey:          unison.KeyB,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// NewSpellContainer creates a new spell container.
var NewSpellContainer = &unison.Action{
	ID:              NewSpellContainerItemID,
	Title:           i18n.Text("New Spell Container"),
	HotKey:          unison.KeyB,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// NewRitualMagicSpell creates a new ritual magic spell.
var NewRitualMagicSpell = &unison.Action{
	ID:              NewRitualMagicSpellItemID,
	Title:           i18n.Text("New Ritual Magic Spell"),
	HotKey:          unison.KeyB,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier | unison.OptionModifier,
	ExecuteCallback: unimplemented,
}

// NewCarriedEquipment creates a new equipment item.
var NewCarriedEquipment = &unison.Action{
	ID:              NewCarriedEquipmentItemID,
	Title:           i18n.Text("New Carried Equipment"),
	HotKey:          unison.KeyE,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// NewCarriedEquipmentContainer creates a new equipment container.
var NewCarriedEquipmentContainer = &unison.Action{
	ID:              NewCarriedEquipmentContainerItemID,
	Title:           i18n.Text("New Carried Equipment Container"),
	HotKey:          unison.KeyE,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// NewOtherEquipment creates a new equipment item.
var NewOtherEquipment = &unison.Action{
	ID:              NewOtherEquipmentItemID,
	Title:           i18n.Text("New Other Equipment"),
	HotKey:          unison.KeyE,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.OptionModifier,
	ExecuteCallback: unimplemented,
}

// NewOtherEquipmentContainer creates a new equipment container.
var NewOtherEquipmentContainer = &unison.Action{
	ID:              NewOtherEquipmentContainerItemID,
	Title:           i18n.Text("New Other Equipment Container"),
	HotKey:          unison.KeyE,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier | unison.OptionModifier,
	ExecuteCallback: unimplemented,
}

// NewEquipmentModifier creates a new equipment modifier.
var NewEquipmentModifier = &unison.Action{
	ID:              NewEquipmentModifierItemID,
	Title:           i18n.Text("New Equipment Modifier"),
	HotKey:          unison.KeyM,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// NewEquipmentContainerModifier creates a new equipment container modifier.
var NewEquipmentContainerModifier = &unison.Action{
	ID:              NewEquipmentContainerModifierItemID,
	Title:           i18n.Text("New Equipment Container Modifier"),
	HotKey:          unison.KeyM,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// NewNote creates a new note.
var NewNote = &unison.Action{
	ID:              NewNoteItemID,
	Title:           i18n.Text("New Note"),
	HotKey:          unison.KeyN,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// NewNoteContainer creates a new note container.
var NewNoteContainer = &unison.Action{
	ID:              NewNoteContainerItemID,
	Title:           i18n.Text("New Note Container"),
	HotKey:          unison.KeyN,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier | unison.OptionModifier,
	ExecuteCallback: unimplemented,
}

// OpenEditor opens an editor for the selected item(s).
var OpenEditor = &unison.Action{
	ID:              OpenEditorItemID,
	Title:           i18n.Text("Open Detail Editor"),
	HotKey:          unison.KeyI,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// CopyToSheet copies the selected items to the foremost character sheet.
var CopyToSheet = &unison.Action{
	ID:              CopyToSheetItemID,
	Title:           i18n.Text("Copy to Character Sheet"),
	HotKey:          unison.KeyC,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// CopyToTemplate copies the selected items to the foremost template.
var CopyToTemplate = &unison.Action{
	ID:              CopyToTemplateItemID,
	Title:           i18n.Text("Copy to Template"),
	HotKey:          unison.KeyT,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// ApplyTemplate applies the foremost template to the foremost character sheet.
var ApplyTemplate = &unison.Action{
	ID:              ApplyTemplateItemID,
	Title:           i18n.Text("Apply Template To Character Sheet"),
	HotKey:          unison.KeyA,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// OpenOnePageReference opens the first page reference for each selected item.
var OpenOnePageReference = &unison.Action{
	ID:              OpenOnePageReferenceItemID,
	Title:           i18n.Text("Open Page Reference"),
	HotKey:          unison.KeyG,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// OpenEachPageReference opens each page reference associated with the selected items.
var OpenEachPageReference = &unison.Action{
	ID:              OpenEachPageReferenceItemID,
	Title:           i18n.Text("Open Each Page Reference"),
	HotKey:          unison.KeyG,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}
