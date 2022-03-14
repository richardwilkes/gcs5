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

package sheet

import (
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// NewRangedWeaponsPageList creates the ranged weapons page list.
func NewRangedWeaponsPageList(entity *gurps.Entity) *PageList {
	return NewPageList(entity, []unison.TableColumnHeader{
		NewPageListHeader(i18n.Text("Ranged Weapon"), ""),
		NewPageListHeader(i18n.Text("Usage"), ""),
		NewPageListHeader(i18n.Text("Lvl"), ""),
		NewPageListHeader(i18n.Text("Acc"), ""),
		NewPageListHeader(i18n.Text("Damage"), ""),
		NewPageListHeader(i18n.Text("Range"), ""),
		NewPageListHeader(i18n.Text("RoF"), ""),
		NewPageListHeader(i18n.Text("Shots"), ""),
		NewPageListHeader(i18n.Text("Bulk"), ""),
		NewPageListHeader(i18n.Text("Rcl"), ""),
		NewPageListHeader(i18n.Text("ST"), ""),
	})
}
