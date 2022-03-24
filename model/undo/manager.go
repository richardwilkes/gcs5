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

package undo

import "github.com/richardwilkes/unison"

type Provider interface {
	UndoManager() *unison.UndoManager
}

func Manager(paneler unison.Paneler) *unison.UndoManager {
	p := paneler.AsPanel()
	for p != nil {
		if mgr, ok := p.Self.(Provider); ok {
			return mgr.UndoManager()
		}
		p = p.Parent()
	}
	return nil
}
