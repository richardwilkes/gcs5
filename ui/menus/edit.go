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

func registerEditMenuActions() {
	settings.RegisterKeyBinding("undo", Undo)
	settings.RegisterKeyBinding("redo", Redo)
	settings.RegisterKeyBinding("duplicate", Duplicate)
	settings.RegisterKeyBinding("cut", unison.CutAction)
	settings.RegisterKeyBinding("copy", unison.CopyAction)
	settings.RegisterKeyBinding("paste", unison.PasteAction)
	settings.RegisterKeyBinding("delete", unison.DeleteAction)
	settings.RegisterKeyBinding("select_all", unison.SelectAllAction)
	settings.RegisterKeyBinding("convert_to_container", ConvertToContainer)
}

func setupEditMenu(bar unison.Menu) {
	f := bar.Factory()
	m := bar.Menu(unison.EditMenuID)
	m.InsertItem(0, Undo.NewMenuItem(f))
	m.InsertItem(1, Redo.NewMenuItem(f))
	m.InsertSeparator(2, false)
	m.InsertItem(m.Item(unison.DeleteItemID).Index(), Duplicate.NewMenuItem(f))
	i := m.Item(unison.SelectAllItemID).Index() + 1
	m.InsertSeparator(i, false)
	i++
	m.InsertItem(i, ConvertToContainer.NewMenuItem(f))
	i++
	m.InsertSeparator(i, false)
	i++
	m.InsertMenu(i, createStateMenu(f))
}

// Undo the last action.
var Undo = &unison.Action{
	ID:         constants.UndoItemID,
	Title:      unison.CannotUndoTitle(),
	KeyBinding: unison.KeyBinding{KeyCode: unison.KeyZ, Modifiers: unison.OSMenuCmdModifier()},
	EnabledCallback: func(action *unison.Action, _ interface{}) bool {
		if wnd := unison.ActiveWindow(); wnd != nil {
			if mgr := wnd.UndoManager(); mgr != nil {
				action.Title = mgr.UndoTitle()
				return mgr.CanUndo()
			}
		}
		action.Title = unison.CannotUndoTitle()
		return false
	},
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		if wnd := unison.ActiveWindow(); wnd != nil {
			if mgr := wnd.UndoManager(); mgr != nil {
				mgr.Undo()
			}
		}
	},
}

// Redo the last action.
var Redo = &unison.Action{
	ID:         constants.RedoItemID,
	Title:      unison.CannotRedoTitle(),
	KeyBinding: unison.KeyBinding{KeyCode: unison.KeyY, Modifiers: unison.OSMenuCmdModifier()},
	EnabledCallback: func(action *unison.Action, _ interface{}) bool {
		if wnd := unison.ActiveWindow(); wnd != nil {
			if mgr := wnd.UndoManager(); mgr != nil {
				action.Title = mgr.RedoTitle()
				return mgr.CanRedo()
			}
		}
		action.Title = unison.CannotRedoTitle()
		return false
	},
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		if wnd := unison.ActiveWindow(); wnd != nil {
			if mgr := wnd.UndoManager(); mgr != nil {
				mgr.Redo()
			}
		}
	},
}

// Duplicate the currently selected content.
var Duplicate = &unison.Action{
	ID:              constants.DuplicateItemID,
	Title:           i18n.Text("Duplicate"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyU, Modifiers: unison.OSMenuCmdModifier()},
	EnabledCallback: notEnabled,
	ExecuteCallback: unimplemented,
}

// ConvertToContainer converts the currently selected item into a container.
var ConvertToContainer = &unison.Action{
	ID:              constants.ConvertToContainerItemID,
	Title:           i18n.Text("Convert to Container"),
	EnabledCallback: notEnabled,
	ExecuteCallback: unimplemented,
}
