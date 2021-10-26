package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/richardwilkes/gcs/internal/library"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xio/fs/paths"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

// Versions of the settings that are supported.
const (
	CurrentSettingsVersion = 3
	MinimumSettingsVersion = 3
)

// Global settings.
var Global = newGlobal()

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

// Path returns the path to our settings file.
func Path() string {
	return filepath.Join(paths.AppDataDir(), cmdline.AppCmdName+"_prefs.json")
}
