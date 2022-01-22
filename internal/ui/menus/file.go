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
	"github.com/richardwilkes/gcs/internal/ui/workspace"
	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

func setupFileMenu(bar unison.Menu) {
	f := bar.Factory()
	m := bar.Menu(unison.FileMenuID)
	newFileMenu := f.NewMenu(NewFileMenuID, i18n.Text("New File…"), nil)
	newFileMenu.InsertItem(-1, NewCharacterSheet.NewMenuItem(f))
	newFileMenu.InsertItem(-1, NewCharacterTemplate.NewMenuItem(f))
	newFileMenu.InsertItem(-1, NewAdvantagesLibrary.NewMenuItem(f))
	newFileMenu.InsertItem(-1, NewAdvantageModifiersLibrary.NewMenuItem(f))
	newFileMenu.InsertItem(-1, NewEquipmentLibrary.NewMenuItem(f))
	newFileMenu.InsertItem(-1, NewEquipmentModifiersLibrary.NewMenuItem(f))
	newFileMenu.InsertItem(-1, NewNotesLibrary.NewMenuItem(f))
	newFileMenu.InsertItem(-1, NewSkillsLibrary.NewMenuItem(f))
	newFileMenu.InsertItem(-1, NewSpellsLibrary.NewMenuItem(f))
	m.InsertMenu(0, newFileMenu)
	m.InsertItem(1, Open.NewMenuItem(f))
	m.InsertMenu(2, f.NewMenu(RecentFilesMenuID, i18n.Text("Recent Files"), recentFilesUpdater))
	i := m.Item(unison.CloseItemID).Index()
	m.RemoveItem(i)
	m.InsertItem(i, CloseTab.NewMenuItem(f))
	i++
	m.InsertSeparator(i, false)
	i++
	m.InsertItem(i, Save.NewMenuItem(f))
	i++
	m.InsertItem(i, SaveAs.NewMenuItem(f))
	i++
	m.InsertMenu(i, f.NewMenu(ExportToMenuID, i18n.Text("Export To…"), exportToUpdater))
	i++
	m.InsertSeparator(i, false)
	i++
	m.InsertItem(i, Print.NewMenuItem(f))
}

func recentFilesUpdater(_ unison.Menu) {
	// TODO: Implement
}

func exportToUpdater(_ unison.Menu) {
	// TODO: Implement
}

// NewCharacterSheet creates a new character sheet.
var NewCharacterSheet = &unison.Action{
	ID:              NewSheetItemID,
	Title:           i18n.Text("New Character Sheet"),
	HotKey:          unison.KeyN,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// NewCharacterTemplate creates a new character template.
var NewCharacterTemplate = &unison.Action{
	ID:              NewTemplateItemID,
	Title:           i18n.Text("New Character Template"),
	ExecuteCallback: unimplemented,
}

// NewAdvantagesLibrary creates a new advantages library.
var NewAdvantagesLibrary = &unison.Action{
	ID:              NewAdvantagesLibraryItemID,
	Title:           i18n.Text("New Advantages Library"),
	ExecuteCallback: unimplemented,
}

// NewAdvantageModifiersLibrary creates a new advantage modifiers library.
var NewAdvantageModifiersLibrary = &unison.Action{
	ID:              NewAdvantageModifiersLibraryItemID,
	Title:           i18n.Text("New Advantage Modifiers Library"),
	ExecuteCallback: unimplemented,
}

// NewEquipmentLibrary creates a new equipment library.
var NewEquipmentLibrary = &unison.Action{
	ID:              NewEquipmentLibraryItemID,
	Title:           i18n.Text("New Equipment Library"),
	ExecuteCallback: unimplemented,
}

// NewEquipmentModifiersLibrary creates a new equipment modifiers library.
var NewEquipmentModifiersLibrary = &unison.Action{
	ID:              NewEquipmentModifiersLibraryItemID,
	Title:           i18n.Text("New Equipment Modifiers Library"),
	ExecuteCallback: unimplemented,
}

// NewNotesLibrary creates a new notes library.
var NewNotesLibrary = &unison.Action{
	ID:              NewNotesLibraryItemID,
	Title:           i18n.Text("New Notes Library"),
	ExecuteCallback: unimplemented,
}

// NewSkillsLibrary creates a new skills library.
var NewSkillsLibrary = &unison.Action{
	ID:              NewSkillsLibraryItemID,
	Title:           i18n.Text("New Skills Library"),
	ExecuteCallback: unimplemented,
}

// NewSpellsLibrary creates a new spells library.
var NewSpellsLibrary = &unison.Action{
	ID:              NewSpellsLibraryItemID,
	Title:           i18n.Text("New Spells Library"),
	ExecuteCallback: unimplemented,
}

// Open a file.
var Open = &unison.Action{
	ID:         OpenItemID,
	Title:      i18n.Text("Open…"),
	HotKey:     unison.KeyO,
	HotKeyMods: unison.OSMenuCmdModifier(),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		dialog := unison.NewOpenDialog()
		dialog.SetAllowsMultipleSelection(true)
		dialog.SetResolvesAliases(true)
		dialog.SetAllowedExtensions(library.AcceptableExtensions()...)
		if dialog.RunModal() {
			workspace.OpenFiles(dialog.Paths())
		}
	},
}

// CloseTab closes a workspace tab.
var CloseTab = &unison.Action{
	ID:         CloseTabID,
	Title:      i18n.Text("Close Tab"),
	HotKey:     unison.KeyW,
	HotKeyMods: unison.OSMenuCmdModifier(),
	EnabledCallback: func(_ *unison.Action, _ interface{}) bool {
		if wnd := unison.ActiveWindow(); wnd != nil {
			if dc := unison.FocusedDockContainerFor(wnd); dc != nil {
				if current := dc.CurrentDockable(); current != nil {
					if _, ok := current.(unison.TabCloser); ok {
						return true
					}
				}
			}
		}
		return false
	},
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		if wnd := unison.ActiveWindow(); wnd != nil {
			if dc := unison.FocusedDockContainerFor(wnd); dc != nil {
				if current := dc.CurrentDockable(); current != nil {
					if closer, ok := current.(unison.TabCloser); ok {
						closer.AttemptClose()
					}
				}
			}
		}
	},
}

// Save a file.
var Save = &unison.Action{
	ID:              SaveItemID,
	Title:           i18n.Text("Save"),
	HotKey:          unison.KeyS,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// SaveAs saves to a new file.
var SaveAs = &unison.Action{
	ID:              SaveAsItemID,
	Title:           i18n.Text("Save As…"),
	HotKey:          unison.KeyS,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// Print the content.
var Print = &unison.Action{
	ID:              PrintItemID,
	Title:           i18n.Text("Print…"),
	HotKey:          unison.KeyP,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}
