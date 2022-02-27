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
	"path"
	"strings"

	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/unison"
)

// Some special "extension" values.
const (
	GenericFile  = "file"
	ClosedFolder = "folder-closed"
	OpenFolder   = "folder-open"
)

// FileInfo contains some static information about a given file type.
type FileInfo struct {
	Extension    string
	SVG          *unison.SVG
	Load         func(filePath string) (unison.Dockable, error)
	IsSpecial    bool
	IsGCSData    bool
	IsImage      bool
	IsPDF        bool
	IsExportable bool
}

// FileTypes holds a map of keys to FileInfo data.
var FileTypes = make(map[string]*FileInfo)

// Register with the central registry.
func (f *FileInfo) Register() {
	FileTypes[f.Extension] = f
}

// FileInfoFor returns the FileInfo for the given file path's extension.
func FileInfoFor(filePath string) *FileInfo {
	info, ok := FileTypes[strings.ToLower(path.Ext(filePath))]
	if !ok {
		info = FileTypes[GenericFile]
	}
	return info
}

// AcceptableExtensions returns the file extensions that we should be able to open.
func AcceptableExtensions() []string {
	list := make([]string, 0, len(FileTypes))
	for k, v := range FileTypes {
		if !v.IsSpecial {
			list = append(list, k)
		}
	}
	txt.SortStringsNaturalAscending(list)
	return list
}
