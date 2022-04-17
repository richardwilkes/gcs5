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
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/datafile"
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/ui/workspace"
	gurpsui "github.com/richardwilkes/gcs/ui/workspace/lists"
	"github.com/richardwilkes/gcs/ui/workspace/sheet"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

func registerFileMenuActions() {
	settings.RegisterKeyBinding("new.char.sheet", NewCharacterSheet)
	settings.RegisterKeyBinding("new.char.template", NewCharacterTemplate)
	settings.RegisterKeyBinding("new.adq.lib", NewAdvantagesLibrary)
	settings.RegisterKeyBinding("new.adm.lib", NewAdvantageModifiersLibrary)
	settings.RegisterKeyBinding("new.eqp.lib", NewEquipmentLibrary)
	settings.RegisterKeyBinding("new.eqp.lib", NewEquipmentModifiersLibrary)
	settings.RegisterKeyBinding("new.not.lib", NewNotesLibrary)
	settings.RegisterKeyBinding("new.skl.lib", NewSkillsLibrary)
	settings.RegisterKeyBinding("new.spl.lib", NewSpellsLibrary)
	settings.RegisterKeyBinding("open", Open)
	settings.RegisterKeyBinding("close", CloseTab)
	settings.RegisterKeyBinding("save", Save)
	settings.RegisterKeyBinding("save_as", SaveAs)
	settings.RegisterKeyBinding("print", Print)
}

func setupFileMenu(bar unison.Menu) {
	f := bar.Factory()
	m := bar.Menu(unison.FileMenuID)
	newFileMenu := f.NewMenu(constants.NewFileMenuID, i18n.Text("New File…"), nil)
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
	m.InsertMenu(2, f.NewMenu(constants.RecentFilesMenuID, i18n.Text("Recent Files"), recentFilesUpdater))
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
	m.InsertMenu(i, f.NewMenu(constants.ExportToMenuID, i18n.Text("Export To…"), exportToUpdater))
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
	ID:         constants.NewSheetItemID,
	Title:      i18n.Text("New Character Sheet"),
	KeyBinding: unison.KeyBinding{KeyCode: unison.KeyN, Modifiers: unison.OSMenuCmdModifier()},
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		entity := gurps.NewEntity(datafile.PC)
		workspace.DisplayNewDockable(nil, sheet.NewSheet(entity.Profile.Name+".gcs", entity))
	},
}

// NewCharacterTemplate creates a new character template.
var NewCharacterTemplate = &unison.Action{
	ID:              constants.NewTemplateItemID,
	Title:           i18n.Text("New Character Template"),
	ExecuteCallback: unimplemented,
}

// NewAdvantagesLibrary creates a new advantages library.
var NewAdvantagesLibrary = &unison.Action{
	ID:    constants.NewAdvantagesLibraryItemID,
	Title: i18n.Text("New Advantages Library"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		workspace.DisplayNewDockable(nil, gurpsui.NewAdvantageTableDockable("Advantages.adq", nil))
	},
}

// NewAdvantageModifiersLibrary creates a new advantage modifiers library.
var NewAdvantageModifiersLibrary = &unison.Action{
	ID:    constants.NewAdvantageModifiersLibraryItemID,
	Title: i18n.Text("New Advantage Modifiers Library"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		workspace.DisplayNewDockable(nil, gurpsui.NewAdvantageModifierTableDockable("Advantage Modifiers.adm", nil))
	},
}

// NewEquipmentLibrary creates a new equipment library.
var NewEquipmentLibrary = &unison.Action{
	ID:    constants.NewEquipmentLibraryItemID,
	Title: i18n.Text("New Equipment Library"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		workspace.DisplayNewDockable(nil, gurpsui.NewEquipmentTableDockable("Equipment.eqp", nil))
	},
}

// NewEquipmentModifiersLibrary creates a new equipment modifiers library.
var NewEquipmentModifiersLibrary = &unison.Action{
	ID:    constants.NewEquipmentModifiersLibraryItemID,
	Title: i18n.Text("New Equipment Modifiers Library"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		workspace.DisplayNewDockable(nil, gurpsui.NewEquipmentModifierTableDockable("Equipment Modifiers.eqm", nil))
	},
}

// NewNotesLibrary creates a new notes library.
var NewNotesLibrary = &unison.Action{
	ID:    constants.NewNotesLibraryItemID,
	Title: i18n.Text("New Notes Library"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		workspace.DisplayNewDockable(nil, gurpsui.NewNoteTableDockable("Notes.not", nil))
	},
}

// NewSkillsLibrary creates a new skills library.
var NewSkillsLibrary = &unison.Action{
	ID:    constants.NewSkillsLibraryItemID,
	Title: i18n.Text("New Skills Library"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		workspace.DisplayNewDockable(nil, gurpsui.NewSkillTableDockable("Skills.skl", nil))
	},
}

// NewSpellsLibrary creates a new spells library.
var NewSpellsLibrary = &unison.Action{
	ID:    constants.NewSpellsLibraryItemID,
	Title: i18n.Text("New Spells Library"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		workspace.DisplayNewDockable(nil, gurpsui.NewSpellTableDockable("Spells.spl", nil))
	},
}

// Open a file.
var Open = &unison.Action{
	ID:         constants.OpenItemID,
	Title:      i18n.Text("Open…"),
	KeyBinding: unison.KeyBinding{KeyCode: unison.KeyO, Modifiers: unison.OSMenuCmdModifier()},
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

// CloseTab closes a workspace tab if the workspace is foremost, or the current window if not.
var CloseTab = &unison.Action{
	ID:         constants.CloseTabID,
	Title:      i18n.Text("Close"),
	KeyBinding: unison.KeyBinding{KeyCode: unison.KeyW, Modifiers: unison.OSMenuCmdModifier()},
	EnabledCallback: func(_ *unison.Action, _ interface{}) bool {
		if wnd := unison.ActiveWindow(); wnd != nil {
			if workspace.FromWindow(wnd) == nil {
				return true // not the workspace, so allow regular window close
			}
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
			if workspace.FromWindow(wnd) == nil {
				// not the workspace, so allow regular window close
				wnd.AttemptClose()
			} else if dc := unison.FocusedDockContainerFor(wnd); dc != nil {
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
	ID:              constants.SaveItemID,
	Title:           i18n.Text("Save"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyS, Modifiers: unison.OSMenuCmdModifier()},
	ExecuteCallback: unimplemented,
}

// SaveAs saves to a new file.
var SaveAs = &unison.Action{
	ID:              constants.SaveAsItemID,
	Title:           i18n.Text("Save As…"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyS, Modifiers: unison.ShiftModifier | unison.OSMenuCmdModifier()},
	ExecuteCallback: unimplemented,
}

// Print the content.
var Print = &unison.Action{
	ID:              constants.PrintItemID,
	Title:           i18n.Text("Print…"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyP, Modifiers: unison.OSMenuCmdModifier()},
	ExecuteCallback: unimplemented,
}
