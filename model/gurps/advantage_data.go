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

package gurps

import (
	"github.com/richardwilkes/toolbox/i18n"
)

// AdvantageData holds the Advantage data that is written to disk.
type AdvantageData struct {
	ContainerBase[*Advantage]
	AdvantageEditData
}

// Kind returns the kind of data.
func (d *AdvantageData) Kind() string {
	return d.kind(i18n.Text("Advantage"))
}

// ClearUnusedFieldsForType zeroes out the fields that are not applicable to this type (container vs not-container).
func (d *AdvantageData) ClearUnusedFieldsForType() {
	d.clearUnusedFields()
	if d.Container() {
		d.BasePoints = 0
		d.Levels = 0
		d.PointsPerLevel = 0
		d.Prereq = nil
		d.Weapons = nil
		d.Features = nil
		d.RoundCostDown = false
	} else {
		d.ContainerType = 0
		d.Ancestry = ""
	}
}
