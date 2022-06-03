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

package widget

import (
	"fmt"

	"github.com/richardwilkes/unison"
)

// Rebuildable defines the methods a rebuildable panel should provide.
type Rebuildable interface {
	unison.Paneler
	fmt.Stringer
	Rebuild(full bool)
}

// FindRebuildable looks a Rebuildable starting at 'startAt' and moving up the parent chain. May return nil if one isn't
// found.
func FindRebuildable(startAt unison.Paneler) Rebuildable {
	p := startAt.AsPanel()
	for p != nil {
		if r, ok := p.Self.(Rebuildable); ok {
			return r
		}
		p = p.Parent()
	}
	return nil
}
