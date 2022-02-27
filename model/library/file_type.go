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

	"github.com/richardwilkes/gcs/res"
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
	SVG          *unison.SVG
	IsSpecial    bool
	IsGCSData    bool
	IsImage      bool
	IsPDF        bool
	IsExportable bool
}

// FileTypes holds a map of keys to FileInfo data.
var FileTypes = make(map[string]*FileInfo)

func init() {
	FileTypes[ClosedFolder] = &FileInfo{SVG: res.ClosedFolderSVG, IsSpecial: true}
	FileTypes[OpenFolder] = &FileInfo{SVG: res.OpenFolderSVG, IsSpecial: true}
	FileTypes[GenericFile] = &FileInfo{SVG: res.GenericFileSVG, IsSpecial: true}
	FileTypes[".gcs"] = &FileInfo{SVG: res.GCSSheet, IsGCSData: true, IsExportable: true}
	FileTypes[".gct"] = &FileInfo{SVG: res.GCSTemplate, IsGCSData: true}
	FileTypes[".adq"] = &FileInfo{SVG: res.GCSAdvantages, IsGCSData: true}
	FileTypes[".adm"] = &FileInfo{SVG: res.GCSAdvantageModifiers, IsGCSData: true}
	FileTypes[".eqp"] = &FileInfo{SVG: res.GCSEquipment, IsGCSData: true}
	FileTypes[".eqm"] = &FileInfo{SVG: res.GCSEquipmentModifiers, IsGCSData: true}
	FileTypes[".skl"] = &FileInfo{SVG: res.GCSSkills, IsGCSData: true}
	FileTypes[".spl"] = &FileInfo{SVG: res.GCSSpells, IsGCSData: true}
	FileTypes[".not"] = &FileInfo{SVG: res.GCSNotes, IsGCSData: true}
	FileTypes[".pdf"] = &FileInfo{SVG: res.PDFFileSVG, IsPDF: true}
	FileTypes[".png"] = &FileInfo{SVG: res.ImageFileSVG, IsImage: true}
	FileTypes[".jpg"] = &FileInfo{SVG: res.ImageFileSVG, IsImage: true}
	FileTypes[".jpeg"] = &FileInfo{SVG: res.ImageFileSVG, IsImage: true}
	FileTypes[".webp"] = &FileInfo{SVG: res.ImageFileSVG, IsImage: true}
	FileTypes[".gif"] = &FileInfo{SVG: res.ImageFileSVG, IsImage: true}
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
