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
	i++
	m.InsertSeparator(i, false)
	i++
	m.InsertItem(i, JumpToSearch.NewMenuItem(f))
}

// Undo the last action.
var Undo = &unison.Action{
	ID:              UndoItemID,
	Title:           i18n.Text("Undo"),
	HotKey:          unison.KeyZ,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// Redo the last action.
var Redo = &unison.Action{
	ID:              RedoItemID,
	Title:           i18n.Text("Redo"),
	HotKey:          unison.KeyZ,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// Duplicate the currently selected content.
var Duplicate = &unison.Action{
	ID:              DuplicateItemID,
	Title:           i18n.Text("Duplicate"),
	HotKey:          unison.KeyU,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// ConvertToContainer converts the currently selected item into a container.
var ConvertToContainer = &unison.Action{
	ID:              ConvertToContainerItemID,
	Title:           i18n.Text("Convert to Container"),
	ExecuteCallback: unimplemented,
}

// JumpToSearch switches the focus to the search widget..
var JumpToSearch = &unison.Action{
	ID:              JumpToSearchItemID,
	Title:           i18n.Text("Jump to Search"),
	HotKey:          unison.KeyJ,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}
