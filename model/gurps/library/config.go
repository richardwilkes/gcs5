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

// Config holds the configuration information for a library of data files.
type Config struct {
	Title    string `json:"title"`
	GitHub   string `json:"github"`
	Repo     string `json:"repo"`
	Path     string `json:"path"`
	LastSeen string `json:"last_seen"`
}
