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
	uisettings "github.com/richardwilkes/gcs/ui/workspace/settings"
	"github.com/richardwilkes/gcs/ui/workspace/sheet"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

func registerSettingsMenuActions() {
	settings.RegisterKeyBinding("settings.sheet.per_sheet", PerSheetSettings)
	settings.RegisterKeyBinding("settings.attributes.per_sheet", PerSheetAttributeSettings)
	settings.RegisterKeyBinding("settings.body_type.per_sheet", PerSheetBodyTypeSettings)
	settings.RegisterKeyBinding("settings.sheet.default", DefaultSheetSettings)
	settings.RegisterKeyBinding("settings.attributes.default", DefaultAttributeSettings)
	settings.RegisterKeyBinding("settings.body_type.default", DefaultBodyTypeSettings)
	settings.RegisterKeyBinding("settings.general", GeneralSettings)
	settings.RegisterKeyBinding("settings.pagerefs", PageRefMappings)
	settings.RegisterKeyBinding("settings.colors", ColorSettings)
	settings.RegisterKeyBinding("settings.fonts", FontSettings)
	settings.RegisterKeyBinding("settings.keys", MenuKeySettings)
}

func createSettingsMenu(f unison.MenuFactory) unison.Menu {
	m := f.NewMenu(constants.SettingsMenuID, i18n.Text("Settings"), nil)
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

// PerSheetSettings opens the settings for the front character sheet.
var PerSheetSettings = &unison.Action{
	ID:              constants.PerSheetSettingsItemID,
	Title:           i18n.Text("Sheet Settings…"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyComma, Modifiers: unison.ShiftModifier | unison.OSMenuCmdModifier()},
	EnabledCallback: func(_ *unison.Action, _ interface{}) bool { return sheet.ActiveEntity() != nil },
	ExecuteCallback: func(_ *unison.Action, _ interface{}) {
		if entity := sheet.ActiveEntity(); entity != nil {
			uisettings.ShowSheetSettings(entity)
		}
	},
}

// DefaultSheetSettings opens the default settings for the character sheet.
var DefaultSheetSettings = &unison.Action{
	ID:              constants.DefaultSheetSettingsItemID,
	Title:           i18n.Text("Default Sheet Settings…"),
	KeyBinding:      unison.KeyBinding{KeyCode: unison.KeyComma, Modifiers: unison.OSMenuCmdModifier()},
	ExecuteCallback: func(_ *unison.Action, _ interface{}) { uisettings.ShowSheetSettings(nil) },
}

// PerSheetAttributeSettings opens the attributes settings for the foremost character sheet.
var PerSheetAttributeSettings = &unison.Action{
	ID:              constants.PerSheetAttributeSettingsItemID,
	Title:           i18n.Text("Attributes…"),
	ExecuteCallback: unimplemented,
}

// DefaultAttributeSettings opens the default attributes settings.
var DefaultAttributeSettings = &unison.Action{
	ID:              constants.DefaultAttributeSettingsItemID,
	Title:           i18n.Text("Default Attributes…"),
	ExecuteCallback: unimplemented,
}

// PerSheetBodyTypeSettings opens the body type settings for the foremost character sheet.
var PerSheetBodyTypeSettings = &unison.Action{
	ID:              constants.PerSheetBodyTypeSettingsItemID,
	Title:           i18n.Text("Body Type…"),
	ExecuteCallback: unimplemented,
}

// DefaultBodyTypeSettings opens the default body type settings.
var DefaultBodyTypeSettings = &unison.Action{
	ID:              constants.DefaultBodyTypeSettingsItemID,
	Title:           i18n.Text("Default Body Type…"),
	ExecuteCallback: unimplemented,
}

// GeneralSettings opens the general settings.
var GeneralSettings = &unison.Action{
	ID:              constants.GeneralSettingsItemID,
	Title:           i18n.Text("General Settings…"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) { uisettings.ShowGeneralSettings() },
}

// PageRefMappings opens the page reference mappings.
var PageRefMappings = &unison.Action{
	ID:              constants.PageRefMappingsItemID,
	Title:           i18n.Text("Page Reference Mappings…"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) { uisettings.ShowPageRefMappings() },
}

// ColorSettings opens the color settings.
var ColorSettings = &unison.Action{
	ID:              constants.ColorSettingsItemID,
	Title:           i18n.Text("Colors…"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) { uisettings.ShowColorSettings() },
}

// FontSettings opens the font settings.
var FontSettings = &unison.Action{
	ID:              constants.FontSettingsItemID,
	Title:           i18n.Text("Fonts…"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) { uisettings.ShowFontSettings() },
}

// MenuKeySettings opens the menu key settings.
var MenuKeySettings = &unison.Action{
	ID:              constants.MenuKeySettingsItemID,
	Title:           i18n.Text("Menu Keys…"),
	ExecuteCallback: func(_ *unison.Action, _ interface{}) { uisettings.ShowMenuKeySettings() },
}
