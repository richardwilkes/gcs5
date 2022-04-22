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

type advantageEditorData struct {
	name         string
	notes        string
	vttNotes     string
	userDesc     string
	tags         string
	pageRef      string
	disabled     bool
	mental       bool
	physical     bool
	social       bool
	exotic       bool
	supernatural bool
}

func (d *advantageEditorData) From(advantage *gurps.Advantage) {
	d.name = advantage.Name
	d.notes = advantage.LocalNotes
	d.vttNotes = advantage.VTTNotes
	d.userDesc = advantage.UserDesc
	d.tags = gurps.CombineTags(advantage.Categories)
	d.pageRef = advantage.PageRef
	d.disabled = advantage.Disabled
	if !advantage.Container() {
		d.mental = advantage.Mental
		d.physical = advantage.Physical
		d.social = advantage.Social
		d.exotic = advantage.Exotic
		d.supernatural = advantage.Supernatural
	}
}

func (d *advantageEditorData) Apply(advantage *gurps.Advantage) {
	advantage.Name = d.name
	advantage.LocalNotes = d.notes
	advantage.VTTNotes = d.vttNotes
	advantage.UserDesc = d.userDesc
	advantage.Categories = gurps.ExtractTags(d.tags)
	advantage.PageRef = d.pageRef
	advantage.Disabled = d.disabled
	if !advantage.Container() {
		advantage.Mental = d.mental
		advantage.Physical = d.physical
		advantage.Social = d.social
		advantage.Exotic = d.exotic
		advantage.Supernatural = d.supernatural
	}
}
