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

package setup

import (
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/gcs/ui/workspace/external"
	"github.com/richardwilkes/gcs/ui/workspace/gurps"
)

// Setup the application. This code is here to break circular dependencies.
func Setup() {
	workspace.RegisterFileTypes()
	external.RegisterFileTypes()
	gurps.RegisterFileTypes()
}
