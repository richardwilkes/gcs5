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

package feature

import (
	"fmt"

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
)

// Feature holds data that affects another object.
type Feature interface {
	// FeatureMapKey returns the key used for matching within the feature map.
	FeatureMapKey() string
	// FillWithNameableKeys fills the map with nameable keys.
	FillWithNameableKeys(m map[string]string)
	// ApplyNameableKeys applies the nameable keys to this object.
	ApplyNameableKeys(m map[string]string)
}

// Bonus is an extension of a Feature, which provides a numerical bonus or penalty.
type Bonus interface {
	Feature
	// AddToTooltip adds this Bonus's details to the tooltip. 'buffer' may be nil.
	AddToTooltip(buffer *xio.ByteBuffer)
}

func basicAddToTooltip(parent fmt.Stringer, amt *LeveledAmount, buffer *xio.ByteBuffer) {
	if buffer != nil {
		buffer.WriteByte('\n')
		buffer.WriteString(parentName(parent))
		buffer.WriteString(" [")
		buffer.WriteString(amt.FormatWithLevel())
		buffer.WriteByte(']')
	}
}

func parentName(parent fmt.Stringer) string {
	if parent == nil {
		return i18n.Text("Unknown")
	}
	return parent.String()
}
