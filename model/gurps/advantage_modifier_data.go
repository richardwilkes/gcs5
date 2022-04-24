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

// AdvantageModifierData holds the AdvantageModifier data that is written to disk.
type AdvantageModifierData struct {
	ContainerBase[*AdvantageModifier]
	AdvantageModifierEditData
}

// Kind returns the kind of data.
func (d *AdvantageModifierData) Kind() string {
	return d.kind(i18n.Text("Advantage Modifier"))
}

// ClearUnusedFieldsForType zeroes out the fields that are not applicable to this type (container vs not-container).
func (d *AdvantageModifierData) ClearUnusedFieldsForType() {
	d.clearUnusedFields()
	if d.Container() {
		d.CostType = 0
		d.Disabled = false
		d.Cost = 0
		d.Levels = 0
		d.Affects = 0
		d.Features = nil
	}
}
