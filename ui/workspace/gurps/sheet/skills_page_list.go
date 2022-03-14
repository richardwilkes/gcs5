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

// NewSkillsPageList creates the skills page list.
func NewSkillsPageList(entity *gurps.Entity) *PageList {
	return NewPageList(entity, []unison.TableColumnHeader{
		NewPageListHeader(i18n.Text("Skill"), ""),
		NewPageListHeader(i18n.Text("SL"), i18n.Text("Skill Level")),
		NewPageListHeader(i18n.Text("RSL"), i18n.Text("Relative Skill Level")),
		NewPageListHeader(i18n.Text("Pts"), i18n.Text("Points")),
		NewPageReferenceHeader(),
	})
}
