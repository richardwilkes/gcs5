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
	"context"
	"io/fs"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/toolbox/i18n"
)

const (
	masterGitHubAccountName = "richardwilkes"
	masterRepoName          = "gcs_master_library"
	userGitHubAccountName   = "*"
	userRepoName            = "gcs_user_library"
)

// Libraries holds a Library set.
type Libraries map[string]*Library

// NewLibraries creates a new, empty, Libraries object.
func NewLibraries() Libraries {
	libs := Libraries(make(map[string]*Library))
	libs.Master()
	libs.User()
	return libs
}

// NewLibrariesFromFS creates a new Ancestry from a file.
func NewLibrariesFromFS(fileSystem fs.FS, filePath string) (Libraries, error) {
	var libs Libraries
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &libs); err != nil {
		return nil, err
	}
	for k, v := range libs {
		if !v.Valid() {
			delete(libs, k)
		}
		v.ConfigureForKey(k)
	}
	libs.Master()
	libs.User()
	return libs, nil
}

// Master holds information about the master library.
func (l Libraries) Master() *Library {
	lib, ok := l[masterGitHubAccountName+"/"+masterRepoName]
	if !ok {
		lib = &Library{
			Title:             i18n.Text("Master Library"),
			GitHubAccountName: masterGitHubAccountName,
			RepoName:          masterRepoName,
			PathOnDisk:        DefaultMasterLibraryPath(),
		}
		l[lib.Key()] = lib
	}
	return lib
}

// User holds information about the user library.
func (l Libraries) User() *Library {
	lib, ok := l[userGitHubAccountName+"/"+userRepoName]
	if !ok {
		lib = &Library{
			Title:             i18n.Text("User Library"),
			GitHubAccountName: userGitHubAccountName,
			RepoName:          userRepoName,
			PathOnDisk:        DefaultUserLibraryPath(),
		}
		l[lib.Key()] = lib
	}
	return lib
}

// List returns an ordered list of Library objects.
func (l Libraries) List() []*Library {
	libs := make([]*Library, 0, len(l))
	for _, lib := range l {
		libs = append(libs, lib)
	}
	sort.Slice(libs, func(i, j int) bool { return libs[i].Less(libs[j]) })
	return libs
}

// PerformUpdateChecks checks each of the libraries for updates.
func (l Libraries) PerformUpdateChecks() {
	client := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(len(l))
	for _, lib := range l {
		go func(l *Library) {
			defer wg.Done()
			l.CheckForAvailableUpgrade(ctx, client)
		}(lib)
	}
	wg.Wait()
}
