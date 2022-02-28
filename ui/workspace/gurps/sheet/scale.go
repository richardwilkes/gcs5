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
	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/unison"
)

// ScaleProvider should be implemented by root components that can provide a scaling factor.
type ScaleProvider interface {
	CurrentScale() float32
}

// DetermineScale walks the hierarchy looking for a ScaleProvider. If found, returns its current scale, otherwise returns
// 1.
func DetermineScale(p unison.Paneler) float32 {
	for !toolbox.IsNil(p) {
		if s, ok := p.(ScaleProvider); ok {
			return s.CurrentScale()
		}
		p = p.AsPanel().Parent()
	}
	return 1
}
