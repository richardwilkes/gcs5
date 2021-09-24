package library

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
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
	// MinimumLibraryVersion is the oldest version of the library data that can be loaded.
	MinimumLibraryVersion = Version{}
	// IncompatibleFutureLibraryVersion is the newest version at which the library data can no longer be loaded.
	IncompatibleFutureLibraryVersion = Version{Major: 4}
	// Master holds information about the master library.
	Master = &Library{
		Title:  i18n.Text("Main Library"),
		GitHub: "richardwilkes",
		Repo:   "gcs_master_library",
		Path:   DefaultMasterLibraryPath(),
	}
	// User holds information about the user library.
	User = &Library{
		Title:  i18n.Text("User Library"),
		GitHub: "*",
		Repo:   "gcs_user_library",
		Path:   DefaultUserLibraryPath(),
	}
)

// Library holds information about a library of data files.
type Library struct {
	Title    string  `json:"title"`
	GitHub   string  `json:"github"`
	Repo     string  `json:"repo"`
	Path     string  `json:"path"`
	LastSeen Version `json:"last_seen"`
	lock     sync.RWMutex
	upgrade  *Release
}

// DefaultRootLibraryPath returns the default root library path.
func DefaultRootLibraryPath() string {
	var home string
	if u, err := user.Current(); err != nil {
		home = os.Getenv("HOME")
	} else {
		home = u.HomeDir
	}
	return filepath.Join(home, "GCS")
}

// DefaultMasterLibraryPath returns the default master library path.
func DefaultMasterLibraryPath() string {
	return filepath.Join(DefaultRootLibraryPath(), "Master Library")
}

// DefaultUserLibraryPath returns the default user library path.
func DefaultUserLibraryPath() string {
	return filepath.Join(DefaultRootLibraryPath(), "User Library")
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

// UpdatePath updates the path to the Library as well as the version.
func (l *Library) UpdatePath(newPath string) {
	p, err := filepath.Abs(newPath)
	if err != nil {
		jot.Warn(errs.NewWithCause("unable to update path to "+newPath, err))
		return
	}
	if l.Path != p {
		l.Path = p
		l.LastSeen = l.VersionOnDisk()
	}
}

// CheckForAvailableUpgrade returns releases that can be upgraded to.
func (l *Library) CheckForAvailableUpgrade(ctx context.Context, client *http.Client) {
	l.lock.Lock()
	l.upgrade = nil
	l.lock.Unlock()
	available, err := LoadReleases(ctx, client, l.GitHub, l.Repo, l.VersionOnDisk(), func(version Version, notes string) bool {
		return version.Less(MinimumLibraryVersion) || IncompatibleFutureLibraryVersion.Less(version) || IncompatibleFutureLibraryVersion == version
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
	if Master == l {
		return true
	}
	if other == l {
		return false
	}
	if User == l {
		return true
	}
	if other == l {
		return false
	}
	if txt.NaturalLess(l.GitHub, other.GitHub, true) {
		return true
	}
	if l.GitHub != other.GitHub {
		return false
	}
	return txt.NaturalLess(l.Repo, other.Repo, true)
}

// VersionOnDisk returns the version of the data on disk, if it can be determined.
func (l *Library) VersionOnDisk() Version {
	releaseFilePath := filepath.Join(l.Path, releaseFile)
	data, err := ioutil.ReadFile(releaseFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			jot.Warn(errs.NewWithCause("unable to load "+releaseFilePath, err))
		}
		return Version{}
	}
	return VersionFromString(strings.TrimSpace(string(bytes.SplitN(data, []byte{'\n'}, 2)[0])))
}

// Download the release onto the local disk.
func (l *Library) Download(ctx context.Context, client *http.Client, release Release) error {
	if err := os.MkdirAll(l.Path, 0o755); err != nil {
		return errs.NewWithCause("unable to create "+l.Path, err)
	}
	data, err := l.downloadRelease(ctx, client, release)
	if err != nil {
		return err
	}
	var zr *zip.Reader
	if zr, err = zip.NewReader(bytes.NewReader(data), int64(len(data))); err != nil {
		return errs.NewWithCause("unable to open archive "+release.ZipFileURL, err)
	}
	root := filepath.Clean(l.Path)
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
			if err = os.MkdirAll(parent, 0o755); err != nil {
				return errs.NewWithCause("unable to create "+parent, err)
			}
			if err = extractFile(f, fullPath); err != nil {
				return errs.NewWithCause("unable to create "+fullPath, err)
			}
		}
	}
	f := filepath.Join(root, releaseFile)
	if err = ioutil.WriteFile(f, []byte(release.Version.String()+"\n"), 0o644); err != nil {
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
	if file, err = os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.FileInfo().Mode().Perm()&0o755); err != nil {
		return errs.Wrap(err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = errs.Wrap(closeErr)
		}
	}()
	if _, err = io.Copy(file, r); err != nil {
		err = errs.Wrap(err)
	}
	return
}

func (l *Library) downloadRelease(ctx context.Context, client *http.Client, release Release) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, release.ZipFileURL, nil)
	if err != nil {
		return nil, errs.NewWithCause("unable to create request for "+release.ZipFileURL, err)
	}
	var rsp *http.Response
	if rsp, err = client.Do(req); err != nil {
		return nil, errs.NewWithCause("unable to connect to "+release.ZipFileURL, err)
	}
	defer xio.DiscardAndCloseIgnoringErrors(rsp.Body)
	var data []byte
	if data, err = ioutil.ReadAll(rsp.Body); err != nil {
		return nil, errs.NewWithCause("unable to download "+release.ZipFileURL, err)
	}
	return data, nil
}
