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

// LibraryExplorer holds settings for the library explorer view.
type LibraryExplorer struct {
	DividerPosition float32  `json:"divider_position"`
	OpenRowKeys     []string `json:"open_row_keys,omitempty"`
}

// FileRef holds a path to a file and an offset for all page references within that file.
type FileRef struct {
	Path   string `json:"path"`
	Offset int    `json:"offset"`
}

// ExportInfo holds information about a recent export so that it can be redone quickly.
type ExportInfo struct {
	TemplatePath string `json:"template_path"`
	ExportPath   string `json:"export_path"`
	LastUsed     int64  `json:"last_used"`
}
