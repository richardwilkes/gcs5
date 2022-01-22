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
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/i18n"
)

const (
	masterGitHubAccountName = "richardwilkes"
	masterRepoName          = "gcs_master_library"
	userGitHubAccountName   = "*"
	userRepoName            = "gcs_user_library"
)

type Libraries struct {
	libs map[string]*Library
}

// NewLibraries creates a new, empty, Libraries object.
func NewLibraries() *Libraries {
	l := &Libraries{libs: make(map[string]*Library)}
	l.Master()
	l.User()
	return l
}

// NewLibrariesFromJSON creates a new Libraries from a JSON object.
func NewLibrariesFromJSON(data map[string]interface{}) *Libraries {
	l := &Libraries{libs: make(map[string]*Library)}
	for k, v := range data {
		if lib := NewLibraryFromJSON(k, encoding.Object(v)); lib != nil {
			l.libs[lib.Key()] = lib
		}
	}
	l.Master()
	l.User()
	return l
}

// ToKeyedJSON emits this object as JSON with the specified key.
func (l *Libraries) ToKeyedJSON(key string, encoder *encoding.JSONEncoder) {
	encoder.Key(key)
	l.ToJSON(encoder)
}

// ToJSON emits this object as JSON.
func (l *Libraries) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	for _, lib := range l.List() {
		lib.ToKeyedJSON(lib.Key(), encoder)
	}
	encoder.EndObject()
}

// Master holds information about the master library.
func (l *Libraries) Master() *Library {
	lib, ok := l.libs[masterGitHubAccountName+"/"+masterRepoName]
	if !ok {
		lib = &Library{
			Title:             i18n.Text("Master Library"),
			GitHubAccountName: masterGitHubAccountName,
			RepoName:          masterRepoName,
			path:              DefaultMasterLibraryPath(),
		}
		l.libs[lib.Key()] = lib
	}
	return lib
}

// User holds information about the user library.
func (l *Libraries) User() *Library {
	lib, ok := l.libs[userGitHubAccountName+"/"+userRepoName]
	if !ok {
		lib = &Library{
			Title:             i18n.Text("User Library"),
			GitHubAccountName: userGitHubAccountName,
			RepoName:          userRepoName,
			path:              DefaultUserLibraryPath(),
		}
		l.libs[lib.Key()] = lib
	}
	return lib
}

func (l *Libraries) List() []*Library {
	libs := make([]*Library, 0, len(l.libs))
	for _, lib := range l.libs {
		libs = append(libs, lib)
	}
	sort.Slice(libs, func(i, j int) bool { return libs[i].Less(libs[j]) })
	return libs
}

// PerformUpdateChecks checks each of the libraries for updates.
func (l *Libraries) PerformUpdateChecks() {
	client := &http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(len(l.libs))
	for _, lib := range l.libs {
		go func(l *Library) {
			defer wg.Done()
			l.CheckForAvailableUpgrade(ctx, client)
		}(lib)
	}
	wg.Wait()
}
