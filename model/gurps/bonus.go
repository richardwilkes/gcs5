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
	"github.com/richardwilkes/toolbox/xio"
)

type addToTooltiper interface {
	addToTooltip(buffer *xio.ByteBuffer)
}

// Bonus is an extension of a Feature, which provides a numerical bonus or penalty.
type Bonus struct {
	Feature
	LeveledAmount
}

// ParentName returns the name of the parent.
func (b *Bonus) ParentName() string {
	if b.Parent == nil {
		return i18n.Text("Unknown")
	}
	return b.Parent.String()
}

// AddToTooltip adds this Bonus's details to the tooltip. 'buffer' may be nil.
func (b *Bonus) AddToTooltip(buffer *xio.ByteBuffer) {
	if buffer != nil {
		if t, ok := b.Self.(addToTooltiper); ok {
			t.addToTooltip(buffer)
		} else {
			buffer.WriteByte('\n')
			buffer.WriteString(b.ParentName())
			buffer.WriteString(" [")
			buffer.WriteString(b.LeveledAmount.Format(i18n.Text("level")))
			buffer.WriteByte(']')
		}
	}
}
