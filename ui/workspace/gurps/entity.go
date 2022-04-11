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
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/ui/workspace/gurps/sheet"
	"github.com/richardwilkes/gcs/ui/workspace/settings"
)

// ActiveEntity returns the currently active entity.
func ActiveEntity() *gurps.Entity {
	d := settings.ActiveDockable()
	if d == nil {
		return nil
	}
	if s, ok := d.(*sheet.Sheet); ok {
		return s.Entity()
	}
	return nil
}
