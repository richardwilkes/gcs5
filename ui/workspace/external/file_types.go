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

package external

import (
	"github.com/richardwilkes/gcs/model/library"
	"github.com/richardwilkes/gcs/res"
)

// RegisterFileTypes registers external file types.
func RegisterFileTypes() {
	registerPDFFileInfo()
	registerImageFileInfo(".gif")
	registerImageFileInfo(".jpeg")
	registerImageFileInfo(".jpg")
	registerImageFileInfo(".png")
	registerImageFileInfo(".webp")
}

func registerImageFileInfo(ext string) {
	library.FileInfo{
		Extension: ext,
		SVG:       res.ImageFileSVG,
		Load:      NewImageDockable,
		IsImage:   true,
	}.Register()
}

func registerPDFFileInfo() {
	library.FileInfo{
		Extension: ".pdf",
		SVG:       res.PDFFileSVG,
		Load:      NewPDFDockable,
		IsPDF:     true,
	}.Register()
}