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

package workspace

import (
	"strconv"

	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/unison"
)

// OpenReference opens the given page reference.
func OpenReference(wnd *unison.Window, ref, highlight string) {
	i := len(ref) - 1
	for i >= 0 {
		ch := ref[i]
		if ch >= '0' && ch <= '9' {
			i--
		} else {
			i++
			break
		}
	}
	if i > 0 {
		page, err := strconv.Atoi(ref[i:])
		if err != nil {
			return
		}
		key := ref[:i]
		s := settings.Global()
		pageRef := s.PageRefs.Lookup(key)
		if pageRef == nil {
			// TODO: Need to let the user know *what* the dialog is for!
			dialog := unison.NewOpenDialog()
			dialog.SetAllowsMultipleSelection(false)
			dialog.SetResolvesAliases(true)
			dialog.SetAllowedExtensions("pdf")
			if dialog.RunModal() {
				pageRef = &settings.PageRef{
					ID:   key,
					Path: dialog.Paths()[0],
				}
				s.PageRefs.Set(pageRef)
				// TODO: once the settings window for page references is working, need to rebuild it when a ref changes
				// PageRefSettingsWindow.rebuild()
			}
		}
		if pageRef != nil {
			if d, wasOpen := OpenFile(wnd, pageRef.Path); d != nil {
				if pdfDockable, ok := d.(*PDFDockable); ok {
					pdfDockable.SetSearchText(highlight)
					pdfDockable.LoadPage(page + pageRef.Offset - 1) // The pdf package uses 0 for the first page, not 1
					if !wasOpen {
						pdfDockable.ClearHistory()
					}
				}
			}
		}
	}
}
