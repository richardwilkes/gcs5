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

package library

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/toolbox/xio"
)

const releaseFile = "release.txt"

var (
	masterLibrary = New(&Config{
		Title:  i18n.Text("Master Library"),
		GitHub: "richardwilkes",
		Repo:   "gcs_master_library",
		Path:   DefaultMasterLibraryPath(),
	})
	userLibrary = New(&Config{
		Title:  i18n.Text("User Library"),
		GitHub: "*",
		Repo:   "gcs_user_library",
		Path:   DefaultUserLibraryPath(),
	})
)

// Library holds information about a library of data files.
type Library struct {
	config  Config
	fs      fs.FS
	lock    sync.RWMutex
	upgrade *Release
}

// New creates a new library. A copy of the configuration is made, so it may be changed after this call without
// affecting the returned library.
func New(cfg *Config) *Library {
	return &Library{
		config: *cfg,
		fs:     os.DirFS(cfg.Path),
	}
}

// Master holds information about the master library.
func Master() *Library {
	return masterLibrary
}

// ReplaceMaster replaces the existing Master Library with the provided library if the GitHub and Repo internal fields
// are a match.
func ReplaceMaster(lib *Library) {
	if lib.IsMaster() {
		lib.config.Title = masterLibrary.config.Title
		masterLibrary = lib
	}
}

// User holds information about the user library.
func User() *Library {
	return userLibrary
}

// ReplaceUser replaces the existing User Library with the provided library if the GitHub and Repo internal fields
// are a match.
func ReplaceUser(lib *Library) {
	if lib.IsUser() {
		lib.config.Title = userLibrary.config.Title
		userLibrary = lib
	}
}

// PerformUpdateChecks checks each of the libraries for updates.
func PerformUpdateChecks(libraries []*Library) {
	client := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(len(libraries))
	for _, lib := range libraries {
		go func(l *Library) {
			defer wg.Done()
			l.CheckForAvailableUpgrade(ctx, client)
		}(lib)
	}
	wg.Wait()
}

// Title returns the title of this library.
func (l *Library) Title() string {
	return l.config.Title
}

// FS returns the file system for the library.
func (l *Library) FS() fs.FS {
	return l.fs
}

// Config returns a copy of the current configuration.
func (l *Library) Config() Config {
	return l.config
}

// MarshalJSON implements the json.Marshaler interface.
func (l *Library) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(&l.config)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return data, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *Library) UnmarshalJSON(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return errs.Wrap(err)
	}
	p, err := filepath.Abs(cfg.Path)
	if err != nil {
		return errs.NewWithCause("unable to resolve library path: "+cfg.Path, err)
	}
	cfg.Path = p
	l.config = cfg
	l.fs = os.DirFS(p)
	return nil
}

// UpdatePath updates the path to the Library as well as the version.
func (l *Library) UpdatePath(newPath string) error {
	p, err := filepath.Abs(newPath)
	if err != nil {
		return errs.NewWithCause("unable to update library path to "+newPath, err)
	}
	if l.config.Path != p {
		l.config.Path = p
		l.fs = os.DirFS(p)
		l.config.LastSeen = l.VersionOnDisk()
	}
	return nil
}

// IsMaster returns true if this is the Master Library.
func (l *Library) IsMaster() bool {
	return l == masterLibrary ||
		(l.config.GitHub == masterLibrary.config.GitHub && l.config.Repo == masterLibrary.config.Repo)
}

// IsUser returns true if this is the User Library.
func (l *Library) IsUser() bool {
	return l == userLibrary ||
		(l.config.GitHub == userLibrary.config.GitHub && l.config.Repo == userLibrary.config.Repo)
}

// CheckForAvailableUpgrade returns releases that can be upgraded to.
func (l *Library) CheckForAvailableUpgrade(ctx context.Context, client *http.Client) {
	l.lock.Lock()
	l.upgrade = nil
	l.lock.Unlock()
	available, err := LoadReleases(ctx, client, l.config.GitHub, l.config.Repo, l.VersionOnDisk(),
		func(version Version, notes string) bool {
			incompatible := IncompatibleFutureLibraryVersion()
			return version.Less(MinimumLibraryVersion()) || incompatible.Less(version) || incompatible == version
		})
	var upgrade *Release
	if err != nil {
		jot.Error(err)
		upgrade = &Release{CheckFailed: true}
	} else {
		switch len(available) {
		case 0:
			upgrade = &Release{}
		case 1:
			upgrade = &available[0]
		default:
			for _, one := range available[1:] {
				available[0].Notes += "\n\n## Version " + one.Version.String() + "\n" + one.Notes
			}
			upgrade = &available[0]
		}
	}
	l.lock.Lock()
	l.upgrade = upgrade
	l.lock.Unlock()
}

// AvailableUpdate returns the available release that can be updated to.
func (l *Library) AvailableUpdate() *Release {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if l.upgrade == nil {
		return nil
	}
	r := *l.upgrade
	return &r
}

// Less returns true if this Library should be placed before the other Library.
func (l *Library) Less(other *Library) bool {
	if masterLibrary == l {
		return true
	}
	if other == l {
		return false
	}
	if userLibrary == l {
		return true
	}
	if other == l {
		return false
	}
	if txt.NaturalLess(l.config.GitHub, other.config.GitHub, true) {
		return true
	}
	if l.config.GitHub != other.config.GitHub {
		return false
	}
	return txt.NaturalLess(l.config.Repo, other.config.Repo, true)
}

// VersionOnDisk returns the version of the data on disk, if it can be determined.
func (l *Library) VersionOnDisk() Version {
	data, err := fs.ReadFile(l.fs, releaseFile)
	if err != nil {
		if !os.IsNotExist(err) {
			jot.Warn(errs.NewWithCause("unable to load "+releaseFile+" from library: "+l.config.Title, err))
		}
		return Version{}
	}
	return VersionFromString(strings.TrimSpace(string(bytes.SplitN(data, []byte{'\n'}, 2)[0])))
}

// Download the release onto the local disk.
func (l *Library) Download(ctx context.Context, client *http.Client, release Release) error {
	if err := os.MkdirAll(l.config.Path, 0o750); err != nil {
		return errs.NewWithCause("unable to create "+l.config.Path, err)
	}
	data, err := l.downloadRelease(ctx, client, release)
	if err != nil {
		return err
	}
	var zr *zip.Reader
	if zr, err = zip.NewReader(bytes.NewReader(data), int64(len(data))); err != nil {
		return errs.NewWithCause("unable to open archive "+release.ZipFileURL, err)
	}
	root := filepath.Clean(l.config.Path)
	rootWithTrailingSep := root
	if !strings.HasSuffix(rootWithTrailingSep, string(filepath.Separator)) {
		rootWithTrailingSep += string(filepath.Separator)
	}
	for _, f := range zr.File {
		fi := f.FileInfo()
		mode := fi.Mode()
		if mode&os.ModeType == 0 { // normal files only
			parts := strings.SplitN(filepath.ToSlash(f.Name), "/", 3)
			if len(parts) != 3 {
				continue
			}
			if !strings.EqualFold("Library", parts[1]) {
				continue
			}
			fullPath := filepath.Join(root, parts[2])
			if !strings.HasPrefix(fullPath, rootWithTrailingSep) {
				return errs.Newf("path outside of root is not permitted: %s", fullPath)
			}
			parent := filepath.Dir(fullPath)
			if err = os.MkdirAll(parent, 0o750); err != nil {
				return errs.NewWithCause("unable to create "+parent, err)
			}
			if err = extractFile(f, fullPath); err != nil {
				return errs.NewWithCause("unable to create "+fullPath, err)
			}
		}
	}
	f := filepath.Join(root, releaseFile)
	if err = os.WriteFile(f, []byte(release.Version.String()+"\n"), 0o640); err != nil {
		return errs.NewWithCause("unable to create "+f, err)
	}
	return nil
}

func extractFile(f *zip.File, dst string) (err error) {
	var r io.ReadCloser
	if r, err = f.Open(); err != nil {
		return errs.Wrap(err)
	}
	defer xio.CloseIgnoringErrors(r)
	var file *os.File
	if file, err = os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.FileInfo().Mode().Perm()&0o750); err != nil {
		return errs.Wrap(err)
	}
	if _, err = io.Copy(file, r); err != nil {
		err = errs.Wrap(err)
	}
	if closeErr := file.Close(); closeErr != nil && err == nil {
		err = errs.Wrap(closeErr)
	}
	return
}

func (l *Library) downloadRelease(ctx context.Context, client *http.Client, release Release) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, release.ZipFileURL, http.NoBody)
	if err != nil {
		return nil, errs.NewWithCause("unable to create request for "+release.ZipFileURL, err)
	}
	var rsp *http.Response
	if rsp, err = client.Do(req); err != nil {
		return nil, errs.NewWithCause("unable to connect to "+release.ZipFileURL, err)
	}
	defer xio.DiscardAndCloseIgnoringErrors(rsp.Body)
	var data []byte
	if data, err = io.ReadAll(rsp.Body); err != nil {
		return nil, errs.NewWithCause("unable to download "+release.ZipFileURL, err)
	}
	return data, nil
}
