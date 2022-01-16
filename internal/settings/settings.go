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
	"sort"

	"github.com/richardwilkes/gcs/internal/library"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xio/fs/paths"
	"github.com/richardwilkes/toolbox/xio/fs/safe"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// Versions of the settings that are supported.
const (
	CurrentSettingsVersion = 3
	MinimumSettingsVersion = 3
)

var global *Settings

// Settings holds the application settings.
type Settings struct {
	Version            int                              `json:"version"`
	LastSeenGCSVersion library.Version                  `json:"last_seen_gcs_version"`
	General            *General                         `json:"general"`
	Libraries          []*library.Library               `json:"libraries,omitempty"`
	LibraryExplorer    LibraryExplorer                  `json:"library_explorer"`
	RecentFiles        []string                         `json:"recent_files,omitempty"`
	PageRefs           map[string]FileRef               `json:"page_refs,omitempty"`
	KeyBindings        map[string]string                `json:"key_bindings,omitempty"`
	WindowPositions    map[string]geom32.Rect           `json:"window_positions,omitempty"`
	Colors             map[string]unison.Color          `json:"colors,omitempty"`
	Fonts              map[string]unison.FontDescriptor `json:"fonts,omitempty"`
	QuickExports       map[string]*ExportInfo           `json:"quick_exports,omitempty"`
	Sheet              *Sheet                           `json:"sheet_settings"`
}

// Default returns new default settings.
func Default() *Settings {
	return &Settings{
		Version:            CurrentSettingsVersion,
		LastSeenGCSVersion: library.VersionFromString(cmdline.AppVersion),
		General:            NewGeneral(),
		Libraries:          []*library.Library{library.Master(), library.User()},
		LibraryExplorer:    LibraryExplorer{DividerPosition: 300},
		PageRefs:           make(map[string]FileRef),
		KeyBindings:        make(map[string]string),
		WindowPositions:    make(map[string]geom32.Rect),
		Colors:             make(map[string]unison.Color),
		Fonts:              make(map[string]unison.FontDescriptor),
		QuickExports:       make(map[string]*ExportInfo),
		Sheet:              NewSheet(),
	}
}

// Global returns the global settings.
func Global() *Settings {
	if global == nil {
		global = newGlobal()
	}
	return global
}

func newGlobal() *Settings {
	s := Default()
	s.Libraries = nil // reset so that we don't overwrite our master and user libraries
	p := Path()
	if f, err := os.Open(p); err == nil {
		defer xio.CloseIgnoringErrors(f)
		if err = json.NewDecoder(f).Decode(s); err != nil {
			jot.Error(errs.NewWithCause("unable to read settings from "+p, err))
			return Default()
		}
		if s.Version < MinimumSettingsVersion {
			jot.Error(errs.New("unable to read settings: too old"))
			return Default()
		}
		if s.Version > CurrentSettingsVersion {
			jot.Error(errs.New("unable to read settings: too new"))
			return Default()
		}
	}
	hadMaster := false
	hadUser := false
	for _, lib := range s.Libraries {
		if lib.IsMaster() {
			library.ReplaceMaster(lib)
			hadMaster = true
			if hadUser {
				break
			}
		} else if lib.IsUser() {
			library.ReplaceUser(lib)
			hadUser = true
			if hadMaster {
				break
			}
		}
	}
	if !hadMaster {
		s.Libraries = append(s.Libraries, library.Master())
	}
	if !hadUser {
		s.Libraries = append(s.Libraries, library.User())
	}
	sort.Slice(s.Libraries, func(i, j int) bool { return s.Libraries[i].Less(s.Libraries[j]) })
	return s
}

// Save to the standard path.
func (s *Settings) Save() {
	s.SaveTo(Path())
}

// SaveTo the provided path.
func (s *Settings) SaveTo(filePath string) {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o750); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to create settings directory"), err)
		return
	}
	if err := safe.WriteFileWithMode(filePath, func(w io.Writer) error {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(s)
	}, 0o640); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to save settings"), err)
	}
}

// Path returns the path to our settings file.
func Path() string {
	return filepath.Join(paths.AppDataDir(), cmdline.AppCmdName+"_prefs.json")
}
