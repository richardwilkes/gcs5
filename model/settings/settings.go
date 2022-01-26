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

package settings

import (
	"path/filepath"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/xio/fs/paths"
)

const (
	settingsLastSeenGCSVersionKey = "last_seen_gcs_version"
	settingsGeneralKey            = "general"
	settingsLibrariesKey          = "libraries"
	settingsLibraryExplorerKey    = "library_explorer"
	settingsRecentFilesKey        = "recent_files"
	settingsLastDirsKey           = "last_dirs"
	settingsPageRefsKey           = "page_refs"
	settingsKeyBindingsKey        = "key_bindings"
	settingsWindowPositionsKey    = "window_positions"
	settingsThemeKey              = "theme"
	settingsQuickExportsKey       = "quick_exports"
	settingsSheetKey              = "sheet_settings"
)

var global *Settings

// Settings holds the application settings.
type Settings struct {
	LastSeenGCSVersion string
	General            *gurps.GeneralSettings
	Libraries          *library.Libraries
	LibraryExplorer    *NavigatorSettings
	RecentFiles        *RecentFiles
	LastDirs           *LastDirs
	PageRefs           *PageRefs
	KeyBindings        *KeyBindings
	WindowPositions    *WindowPositions
	Theme              *Theme
	QuickExports       *gurps.QuickExports
	Sheet              *gurps.SheetSettings
}

// Default returns new default settings.
func Default() *Settings {
	return &Settings{
		LastSeenGCSVersion: cmdline.AppVersion,
		General:            gurps.NewGeneralSettings(),
		Libraries:          library.NewLibraries(),
		LibraryExplorer:    NewNavigatorSettings(),
		RecentFiles:        NewRecentFiles(),
		LastDirs:           NewLastDirs(),
		PageRefs:           NewPageRefs(),
		KeyBindings:        NewKeyBindings(),
		WindowPositions:    NewWindowPositions(),
		Theme:              NewTheme(),
		QuickExports:       gurps.NewQuickExports(),
		Sheet:              gurps.FactorySheetSettings(),
	}
}

// Global returns the global settings.
func Global() *Settings {
	if global == nil {
		dice.GURPSFormat = true
		p := Path()
		if data, err := encoding.LoadJSON(p); err == nil {
			obj := encoding.Object(data)
			global = &Settings{
				LastSeenGCSVersion: encoding.String(obj[settingsLastSeenGCSVersionKey]),
				General:            gurps.NewGeneralSettingsFromJSON(encoding.Object(obj[settingsGeneralKey])),
				Libraries:          library.NewLibrariesFromJSON(encoding.Object(obj[settingsLibrariesKey])),
				LibraryExplorer:    NewNavigatorSettingsFromJSON(encoding.Object(obj[settingsLibraryExplorerKey])),
				RecentFiles:        NewRecentFilesFromJSON(encoding.Object(obj[settingsRecentFilesKey])),
				LastDirs:           NewLastDirsFromJSON(encoding.Object(obj[settingsLastDirsKey])),
				PageRefs:           NewPageRefsFromJSON(encoding.Object(obj[settingsPageRefsKey])),
				KeyBindings:        NewKeyBindingsFromJSON(encoding.Object(obj[settingsKeyBindingsKey])),
				WindowPositions:    NewWindowPositionsFromJSON(encoding.Object(obj[settingsWindowPositionsKey])),
				Theme:              NewThemeFromJSON(encoding.Object(obj[settingsThemeKey])),
				QuickExports:       gurps.NewQuickExportsFromJSON(encoding.Object(obj[settingsQuickExportsKey])),
				Sheet:              gurps.NewSheetSettingsFromJSON(encoding.Object(obj[settingsSheetKey]), nil),
			}
		} else {
			global = Default()
		}
		gurps.GlobalSheetSettingsProvider = func() *gurps.SheetSettings { return global.Sheet }
	}
	return global
}

// Save to the standard path.
func (s *Settings) Save() error {
	return encoding.SaveJSON(Path(), true, s.toJSON)
}

func (s *Settings) toJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(settingsLastSeenGCSVersionKey, s.LastSeenGCSVersion, false, false)
	encoding.ToKeyedJSON(s.General, settingsGeneralKey, encoder)
	encoding.ToKeyedJSON(s.Libraries, settingsLibrariesKey, encoder)
	encoding.ToKeyedJSON(s.LibraryExplorer, settingsLibraryExplorerKey, encoder)
	encoding.ToKeyedJSON(s.RecentFiles, settingsRecentFilesKey, encoder)
	encoding.ToKeyedJSON(s.LastDirs, settingsLastDirsKey, encoder)
	encoding.ToKeyedJSON(s.PageRefs, settingsPageRefsKey, encoder)
	encoding.ToKeyedJSON(s.KeyBindings, settingsKeyBindingsKey, encoder)
	encoding.ToKeyedJSON(s.WindowPositions, settingsWindowPositionsKey, encoder)
	encoding.ToKeyedJSON(s.Theme, settingsThemeKey, encoder)
	encoding.ToKeyedJSON(s.QuickExports, settingsQuickExportsKey, encoder)
	gurps.ToKeyedJSON(s.Sheet, settingsSheetKey, encoder, nil)
	encoder.EndObject()
}

// Path returns the path to our settings file.
func Path() string {
	return filepath.Join(paths.AppDataDir(), cmdline.AppCmdName+"_prefs.json")
}
