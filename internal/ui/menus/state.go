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

func createStateMenu(f unison.MenuFactory) unison.Menu {
	m := f.NewMenu(StateMenuID, i18n.Text("State…"), nil)
	m.InsertItem(-1, ToggleState.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, Increment.NewMenuItem(f))
	m.InsertItem(-1, Decrement.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, IncreaseUses.NewMenuItem(f))
	m.InsertItem(-1, DecreaseUses.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, IncreaseSkillLevel.NewMenuItem(f))
	m.InsertItem(-1, DecreaseSkillLevel.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, IncreaseTechLevel.NewMenuItem(f))
	m.InsertItem(-1, DecreaseTechLevel.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, SwapDefaults.NewMenuItem(f))
	return m
}

// ToggleState switches the focus to the search field.
var ToggleState = &unison.Action{
	ID:              ToggleStateItemID,
	Title:           i18n.Text("Toggle State"),
	HotKey:          unison.KeyApostrophe,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// Increment the points of the selection.
var Increment = &unison.Action{
	ID:              IncrementItemID,
	Title:           i18n.Text("Increment"),
	HotKey:          unison.KeyEqual,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// Decrement the points of the selection.
var Decrement = &unison.Action{
	ID:              DecrementItemID,
	Title:           i18n.Text("Decrement"),
	HotKey:          unison.KeyMinus,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// IncreaseUses increments the uses of the selection.
var IncreaseUses = &unison.Action{
	ID:              IncrementUsesItemID,
	Title:           i18n.Text("Increase Uses"),
	HotKey:          unison.KeyUp,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// DecreaseUses decrements the uses of the selection.
var DecreaseUses = &unison.Action{
	ID:              DecrementUsesItemID,
	Title:           i18n.Text("Decrease Uses"),
	HotKey:          unison.KeyDown,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// IncreaseSkillLevel increments the uses of the skill level.
var IncreaseSkillLevel = &unison.Action{
	ID:              IncrementSkillLevelItemID,
	Title:           i18n.Text("Increase Skill Level"),
	HotKey:          unison.KeySlash,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// DecreaseSkillLevel decrements the uses of the skill level.
var DecreaseSkillLevel = &unison.Action{
	ID:              DecrementSkillLevelItemID,
	Title:           i18n.Text("Decrease Skill Level"),
	HotKey:          unison.KeyPeriod,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// IncreaseTechLevel increments the uses of the tech level.
var IncreaseTechLevel = &unison.Action{
	ID:              IncrementTechLevelItemID,
	Title:           i18n.Text("Increase Tech Level"),
	HotKey:          unison.KeyCloseBracket,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// DecreaseTechLevel decrements the uses of the tech level.
var DecreaseTechLevel = &unison.Action{
	ID:              DecrementTechLevelItemID,
	Title:           i18n.Text("Decrease Tech Level"),
	HotKey:          unison.KeyOpenBracket,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// SwapDefaults swaps the defaults of the selected skill.
var SwapDefaults = &unison.Action{
	ID:              SwapDefaultsItemID,
	Title:           i18n.Text("Swap Defaults"),
	HotKey:          unison.KeyX,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}
