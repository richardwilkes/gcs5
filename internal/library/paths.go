package library

import (
	"os"
	"os/user"
	"path/filepath"
)

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
