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

package editors

import "github.com/richardwilkes/gcs/model/gurps"

var _ editorData[*gurps.Note] = &noteEditorData{}

type noteEditorData struct {
	note    string
	pageRef string
}

func (d *noteEditorData) From(note *gurps.Note) {
	d.note = note.Text
	d.pageRef = note.PageRef
}

func (d *noteEditorData) Apply(note *gurps.Note) {
	note.Text = d.note
	note.PageRef = d.pageRef
}
