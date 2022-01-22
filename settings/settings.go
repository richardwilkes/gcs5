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
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio/fs/paths"
	"github.com/richardwilkes/toolbox/xio/fs/safe"
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
				Sheet:              gurps.NewSheetSettingsFromJSON(encoding.Object(obj[settingsSheetKey])),
			}
		} else {
			global = Default()
		}
	}
	return global
}

func (s *Settings) toJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(settingsLastSeenGCSVersionKey, s.LastSeenGCSVersion, false, false)
	s.General.ToKeyedJSON(settingsGeneralKey, encoder)
	s.Libraries.ToKeyedJSON(settingsLibrariesKey, encoder)
	s.LibraryExplorer.ToKeyedJSON(settingsLibraryExplorerKey, encoder)
	s.RecentFiles.ToKeyedJSON(settingsRecentFilesKey, encoder)
	s.LastDirs.ToKeyedJSON(settingsLastDirsKey, encoder)
	s.PageRefs.ToKeyedJSON(settingsPageRefsKey, encoder)
	s.KeyBindings.ToKeyedJSON(settingsKeyBindingsKey, encoder)
	s.WindowPositions.ToKeyedJSON(settingsWindowPositionsKey, encoder)
	s.Theme.ToKeyedJSON(settingsThemeKey, encoder)
	s.QuickExports.ToKeyedJSON(settingsQuickExportsKey, encoder)
	s.Sheet.ToKeyedJSON(settingsSheetKey, encoder)
	encoder.EndObject()
}

// Save to the standard path.
func (s *Settings) Save() error {
	filePath := Path()
	if err := os.MkdirAll(filepath.Dir(filePath), 0o750); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	if err := safe.WriteFileWithMode(filePath, func(w io.Writer) error {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(s)
	}, 0o640); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	return nil
}

// Path returns the path to our settings file.
func Path() string {
	return filepath.Join(paths.AppDataDir(), cmdline.AppCmdName+"_prefs.json")
}
