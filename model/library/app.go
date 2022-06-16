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
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/unison"
)

type appUpdater struct {
	lock    sync.RWMutex
	result  string
	release *Release
}

var appUpdate appUpdater

func (u *appUpdater) Reset() {
	u.lock.Lock()
	u.result = fmt.Sprintf(i18n.Text("Checking for %s updates…"), cmdline.AppName)
	u.release = nil
	u.lock.Unlock()
}

func (u *appUpdater) Result() (string, *Release) {
	u.lock.RLock()
	defer u.lock.RUnlock()
	var release *Release
	if u.release != nil {
		other := *u.release
		release = &other
	}
	return u.result, release
}

func (u *appUpdater) SetResult(str string) {
	u.lock.Lock()
	u.result = str
	u.lock.Unlock()
}

func (u *appUpdater) SetRelease(release *Release) {
	u.lock.Lock()
	u.result = fmt.Sprintf(i18n.Text("%s v%s is available!"), cmdline.AppName, release.Version)
	other := *release
	u.release = &other
	u.lock.Unlock()
}

func (u *appUpdater) NotifyOfUpdate() {
	if title, release := u.Result(); release != nil {
		// TODO: Show release notes of all releases between the version running and the version available.
		unison.QuestionDialog(title, "")
	}
}

// CheckForAppUpdates initiates a fresh check for application updates.
func CheckForAppUpdates() {
	if cmdline.AppVersion == "0.0" {
		appUpdate.SetResult(fmt.Sprintf(i18n.Text("Development versions don't look for %s updates"), cmdline.AppName))
		return
	}
	appUpdate.Reset()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()
		releases, err := LoadReleases(ctx, &http.Client{}, "richardwilkes", "gcs", cmdline.AppVersion,
			func(version, notes string) bool {
				return txt.NaturalLess(version, "5.0.0", true)
			})
		if err != nil {
			appUpdate.SetResult(fmt.Sprintf(i18n.Text("Unable to access the %s update site"), cmdline.AppName))
			jot.Error(err)
			return
		}
		if len(releases) == 0 || releases[0].Version == cmdline.AppVersion {
			appUpdate.SetResult(fmt.Sprintf(i18n.Text("%s has no update available"), cmdline.AppName))
			return
		}
		appUpdate.SetRelease(&releases[0])
		unison.InvokeTask(appUpdate.NotifyOfUpdate)
	}()
}

// AppUpdateResult returns the current results of any outstanding app update check.
func AppUpdateResult() (string, *Release) {
	return appUpdate.Result()
}

// AppUpdate will perform the application update or take the user to the website so they can do it themselves.
func AppUpdate() {
	if _, release := appUpdate.Result(); release != nil {
		// TODO: Do the update
	}
}
