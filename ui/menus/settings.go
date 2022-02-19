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
	"github.com/richardwilkes/gcs/ui/settings"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

func createSettingsMenu(f unison.MenuFactory) unison.Menu {
	m := f.NewMenu(SettingsMenuID, i18n.Text("Settings"), nil)
	m.InsertItem(-1, PerSheetSettings.NewMenuItem(f))
	m.InsertItem(-1, PerSheetAttributeSettings.NewMenuItem(f))
	m.InsertItem(-1, PerSheetBodyTypeSettings.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, DefaultSheetSettings.NewMenuItem(f))
	m.InsertItem(-1, DefaultAttributeSettings.NewMenuItem(f))
	m.InsertItem(-1, DefaultBodyTypeSettings.NewMenuItem(f))
	m.InsertSeparator(-1, false)
	m.InsertItem(-1, GeneralSettings.NewMenuItem(f))
	m.InsertItem(-1, PageRefMappings.NewMenuItem(f))
	m.InsertItem(-1, ColorSettings.NewMenuItem(f))
	m.InsertItem(-1, FontSettings.NewMenuItem(f))
	m.InsertItem(-1, MenuKeySettings.NewMenuItem(f))
	return m
}

// PerSheetSettings opens the settings for the foremost character sheet.
var PerSheetSettings = &unison.Action{
	ID:              PerSheetSettingsItemID,
	Title:           i18n.Text("Sheet Settings…"),
	HotKey:          unison.KeyComma,
	HotKeyMods:      unison.OSMenuCmdModifier() | unison.ShiftModifier,
	ExecuteCallback: unimplemented,
}

// DefaultSheetSettings opens the default settings for the character sheet.
var DefaultSheetSettings = &unison.Action{
	ID:              DefaultSheetSettingsItemID,
	Title:           i18n.Text("Default Sheet Settings…"),
	HotKey:          unison.KeyComma,
	HotKeyMods:      unison.OSMenuCmdModifier(),
	ExecuteCallback: unimplemented,
}

// PerSheetAttributeSettings opens the attributes settings for the foremost character sheet.
var PerSheetAttributeSettings = &unison.Action{
	ID:              PerSheetAttributeSettingsItemID,
	Title:           i18n.Text("Attributes…"),
	ExecuteCallback: unimplemented,
}

// DefaultAttributeSettings opens the default attributes settings.
var DefaultAttributeSettings = &unison.Action{
	ID:              DefaultAttributeSettingsItemID,
	Title:           i18n.Text("Default Attributes…"),
	ExecuteCallback: unimplemented,
}

// PerSheetBodyTypeSettings opens the body type settings for the foremost character sheet.
var PerSheetBodyTypeSettings = &unison.Action{
	ID:              PerSheetBodyTypeSettingsItemID,
	Title:           i18n.Text("Body Type…"),
	ExecuteCallback: unimplemented,
}

// DefaultBodyTypeSettings opens the default body type settings.
var DefaultBodyTypeSettings = &unison.Action{
	ID:              DefaultBodyTypeSettingsItemID,
	Title:           i18n.Text("Default Body Type…"),
	ExecuteCallback: unimplemented,
}

// GeneralSettings opens the general settings.
var GeneralSettings = &unison.Action{
	ID:              GeneralSettingsItemID,
	Title:           i18n.Text("General Settings…"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) { settings.ShowGeneralSettings() },
}

// PageRefMappings opens the page reference mappings.
var PageRefMappings = &unison.Action{
	ID:              PageRefMappingsItemID,
	Title:           i18n.Text("Page Reference Mappings…"),
	ExecuteCallback: unimplemented,
}

// ColorSettings opens the color settings.
var ColorSettings = &unison.Action{
	ID:              ColorSettingsItemID,
	Title:           i18n.Text("Colors…"),
	ExecuteCallback: unimplemented,
}

// FontSettings opens the font settings.
var FontSettings = &unison.Action{
	ID:              FontSettingsItemID,
	Title:           i18n.Text("Fonts…"),
	ExecuteCallback: unimplemented,
}

// MenuKeySettings opens the menu key settings.
var MenuKeySettings = &unison.Action{
	ID:              MenuKeySettingsItemID,
	Title:           i18n.Text("Menu Keys…"),
	ExecuteCallback: unimplemented,
}
