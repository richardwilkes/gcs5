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
	"time"

	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
)

var result = i18n.Text("Checking for GCS updates…")

func CheckForAppUpdates() {
	go func() {
		fmt.Println("current version: ", cmdline.AppVersion)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()
		releases, err := LoadReleases(ctx, &http.Client{}, "richardwilkes", "gcs", "4.0.0", nil)
		if err != nil {
			jot.Error(err)
			return
		}
		fmt.Println("Available GCS releases:")
		for _, one := range releases {
			fmt.Println(one.Version)
		}
	}()
}

func AppUpdateResult() string {
	return result
}

func AppUpdateAvailable() bool {
	return false
}

func AppUpdate() {

}
