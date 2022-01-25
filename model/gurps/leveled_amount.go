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

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	leveledAmountAmountKey   = "amount"
	leveledAmountPerLevelKey = "per_level"
)

// LeveledAmount holds an amount that can be either a fixed amount, or an amount per level.
type LeveledAmount struct {
	Amount   fixed.F64d4
	Level    fixed.F64d4
	PerLevel bool
}

// NewLeveledAmountFromJSON creates a new LeveledAmount from a JSON object.
func NewLeveledAmountFromJSON(data map[string]interface{}) *LeveledAmount {
	a := &LeveledAmount{
		Amount:   encoding.Number(data[leveledAmountAmountKey]),
		PerLevel: encoding.Bool(data[leveledAmountPerLevelKey]),
	}
	return a
}

// FromJSON replaces the current data with data from a JSON object.
func (a *LeveledAmount) FromJSON(data map[string]interface{}) {
	a.Amount = encoding.Number(data[leveledAmountAmountKey])
	a.Level = 0
	a.PerLevel = encoding.Bool(data[leveledAmountPerLevelKey])
}

// ToJSON implements Feature.
func (a *LeveledAmount) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	a.ToInlineJSON(encoder)
	encoder.EndObject()
}

// ToInlineJSON emits the JSON key values that comprise this object without the object wrapper.
func (a *LeveledAmount) ToInlineJSON(encoder *encoding.JSONEncoder) {
	encoder.KeyedNumber(leveledAmountAmountKey, a.Amount, false)
	encoder.KeyedBool(leveledAmountPerLevelKey, a.PerLevel, true)
}

// AdjustedAmount returns the amount, adjusted for level, if requested.
func (a *LeveledAmount) AdjustedAmount() fixed.F64d4 {
	if a.PerLevel {
		if a.Level < 0 {
			return 0
		}
		return a.Amount.Mul(a.Level)
	}
	return a.Amount
}

// Format the value.
func (a *LeveledAmount) Format(what string) string {
	str := a.Amount.String()
	if a.Amount >= 0 {
		return "+" + str
	}
	if a.PerLevel {
		full := a.AdjustedAmount()
		fullStr := full.String()
		if full >= 0 {
			fullStr = "+" + fullStr
		}
		return fmt.Sprintf(i18n.Text("%s (%s per %s)"), fullStr, str, what)
	}
	return str
}
