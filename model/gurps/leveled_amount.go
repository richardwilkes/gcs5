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
	"fmt"

	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// LeveledAmount holds an amount that can be either a fixed amount, or an amount per level.
type LeveledAmount struct {
	Level    fixed.F64d4 `json:"-"`
	Amount   fixed.F64d4 `json:"amount"`
	PerLevel bool        `json:"per_level,omitempty"`
}

// AdjustedAmount returns the amount, adjusted for level, if requested.
func (l *LeveledAmount) AdjustedAmount() fixed.F64d4 {
	if l.PerLevel {
		if l.Level < 0 {
			return 0
		}
		return l.Amount.Mul(l.Level)
	}
	return l.Amount
}

// Format the value.
func (l *LeveledAmount) Format(what string) string {
	str := l.Amount.StringWithSign()
	if l.PerLevel {
		return fmt.Sprintf(i18n.Text("%s (%s per %s)"), l.AdjustedAmount().StringWithSign(), str, what)
	}
	return str
}
