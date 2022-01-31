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
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/library"
	"github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/toolbox/xio/fs/paths"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

const maxRecentFiles = 20

var global *Settings

// PageRef holds a path to a file and an offset for all page references within that file.
type PageRef struct {
	Path   string `json:"path,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

// WindowPosition holds a window's last known frame and when the frame's size or position was last altered.
type WindowPosition struct {
	Frame       geom32.Rect `json:"frame"`
	LastUpdated time.Time   `json:"last_updated"`
}

// NavigatorSettings holds settings for the navigator view.
type NavigatorSettings struct {
	DividerPosition int      `json:"divider_position"`
	OpenRowKeys     []string `json:"open_row_keys,omitempty"`
}

// Settings holds the application settings.
type Settings struct {
	LastSeenGCSVersion string                           `json:"last_seen_gcs_version,omitempty"`
	General            *settings.General                `json:"general,omitempty"`
	Libraries          *library.Libraries               `json:"libraries,omitempty"`
	LibraryExplorer    NavigatorSettings                `json:"library_explorer"`
	RecentFiles        []string                         `json:"recent_files,omitempty"`
	LastDirs           map[string]string                `json:"last_dirs,omitempty"`
	PageRefs           map[string]*PageRef              `json:"page_refs,omitempty"`
	KeyBindings        map[string]string                `json:"key_bindings,omitempty"`
	WindowPositions    map[string]*WindowPosition       `json:"window_positions,omitempty"`
	Colors             map[string]unison.Color          `json:"colors,omitempty"`
	Fonts              map[string]unison.FontDescriptor `json:"fonts,omitempty"`
	QuickExports       *gurps.QuickExports              `json:"quick_exports,omitempty"`
	Sheet              *gurps.SheetSettings             `json:"sheet_settings,omitempty"`
}

// Default returns new default settings.
func Default() *Settings {
	return &Settings{
		LastSeenGCSVersion: cmdline.AppVersion,
		General:            settings.NewGeneral(),
		Libraries:          library.NewLibraries(),
		LibraryExplorer:    NavigatorSettings{DividerPosition: 300},
		LastDirs:           make(map[string]string),
		PageRefs:           make(map[string]*PageRef),
		KeyBindings:        make(map[string]string),
		WindowPositions:    make(map[string]*WindowPosition),
		Colors:             make(map[string]unison.Color),
		Fonts:              make(map[string]unison.FontDescriptor),
		QuickExports:       gurps.NewQuickExports(),
		Sheet:              gurps.FactorySheetSettings(),
	}
}

// Global returns the global settings.
func Global() *Settings {
	if global == nil {
		dice.GURPSFormat = true
		if err := fs.LoadJSON(Path(), &global); err != nil {
			global = Default()
		}
		gurps.GlobalSheetSettingsProvider = func() *gurps.SheetSettings { return global.Sheet }
	}
	return global
}

// Save to the standard path.
func (s *Settings) Save() error {
	return fs.SaveJSON(Path(), s, true)
}

// LookupPageRef the PageRef for the given ID. If not found or if the path it points to isn't a readable file, returns
// nil.
func (s *Settings) LookupPageRef(id string) *PageRef {
	if ref, ok := s.PageRefs[id]; ok && fs.FileIsReadable(ref.Path) {
		return ref
	}
	return nil
}

// ListRecentFiles returns the current list of recently opened files. Files that are no longer readable for any reason
// are omitted.
func (s *Settings) ListRecentFiles() []string {
	list := make([]string, 0, len(s.RecentFiles))
	for _, one := range s.RecentFiles {
		if fs.FileIsReadable(one) {
			list = append(list, one)
		}
	}
	if len(list) != len(s.RecentFiles) {
		s.RecentFiles = make([]string, len(list))
		copy(s.RecentFiles, list)
	}
	return list
}

// AddRecentFile adds a file path to the list of recently opened files.
func (s *Settings) AddRecentFile(filePath string) {
	ext := path.Ext(filePath)
	if runtime.GOOS == toolbox.MacOS || runtime.GOOS == toolbox.WindowsOS {
		ext = strings.ToLower(ext)
	}
	for _, one := range library.AcceptableExtensions() {
		if one == ext {
			full, err := filepath.Abs(filePath)
			if err != nil {
				return
			}
			if fs.FileIsReadable(full) {
				for i, f := range s.RecentFiles {
					if f == full {
						copy(s.RecentFiles[i:], s.RecentFiles[i+1:])
						s.RecentFiles[len(s.RecentFiles)-1] = ""
						s.RecentFiles = s.RecentFiles[:len(s.RecentFiles)-1]
						break
					}
				}
				s.RecentFiles = append(s.RecentFiles, "")
				copy(s.RecentFiles[1:], s.RecentFiles)
				s.RecentFiles[0] = full
				if len(s.RecentFiles) > maxRecentFiles {
					s.RecentFiles = s.RecentFiles[:maxRecentFiles]
				}
			}
			return
		}
	}
}

// Path returns the path to our settings file.
func Path() string {
	return filepath.Join(paths.AppDataDir(), cmdline.AppCmdName+"_prefs.json")
}
