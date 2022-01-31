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
	"encoding/json"

	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// AttributeIDPrefix is the prefix all references to attribute IDs should use.
const AttributeIDPrefix = "attr."

// AttributeCalc holds the Attribute data that is only emitted for third parties.
type AttributeCalc struct {
	Value   fixed.F64d4  `json:"value"`
	Current *fixed.F64d4 `json:"current,omitempty"`
	Points  fixed.F64d4  `json:"points"`
}

// AttributeData holds the Attribute data that is written to disk.
type AttributeData struct {
	AttrID     string         `json:"attr_id"`
	Adjustment fixed.F64d4    `json:"adj"`
	Damage     fixed.F64d4    `json:"damage,omitempty"`
	Calc       *AttributeCalc `json:"calc,omitempty"`
}

// Attribute holds the current state of an AttributeDef.
type Attribute struct {
	AttributeData
	Entity        *Entity     `json:"-"`
	Bonus         fixed.F64d4 `json:"-"`
	CostReduction fixed.F64d4 `json:"-"`
}

// MarshalJSON implements json.Marshaler.
func (a *Attribute) MarshalJSON() ([]byte, error) {
	a.Calc = nil
	if a.Entity != nil {
		if def := a.AttributeDef(); def != nil {
			a.Calc = &AttributeCalc{
				Value:  a.Maximum(),
				Points: a.PointCost(),
			}
			if def.Type == attribute.Pool {
				current := a.Current()
				a.Calc.Current = &current
			}
		}
	}
	data, err := json.Marshal(a.AttributeData)
	a.Calc = nil
	return data, err
}

// ID returns the ID.
func (a *Attribute) ID() string {
	return a.AttrID
}

// SetID sets the ID, sanitizing it in the process (i.e. it may be changed from what you set -- read it back if you want
// to be sure of what it gets set to.
func (a *Attribute) SetID(value string) {
	a.AttrID = id.Sanitize(value, false, ReservedIDs...)
}

// AttributeDef looks up the AttributeDef this Attribute references from the Entity. May return nil.
func (a *Attribute) AttributeDef() *AttributeDef {
	if a.Entity == nil {
		return nil
	}
	return a.Entity.SheetSettings.Attributes.Set[a.AttrID]
}

// Maximum returns the maximum value of a pool or the adjusted attribute value for other types.
func (a *Attribute) Maximum() fixed.F64d4 {
	def := a.AttributeDef()
	if def == nil {
		return 0
	}
	max := def.BaseValue(a.Entity) + a.Adjustment + a.Bonus
	if def.Type != attribute.Decimal {
		max = max.Trunc()
	}
	return max
}

// Current returns the current value. Only valid for pools.
func (a *Attribute) Current() fixed.F64d4 {
	return a.Maximum() - a.Damage
}

// SetCurrent sets the current value.
func (a *Attribute) SetCurrent(value fixed.F64d4) {
	if a.Current() == value {
		return
	}
	if def := a.AttributeDef(); def != nil {
		a.Adjustment = value - (def.BaseValue(a.Entity) + a.Bonus)
	}
}

// CurrentThreshold return the current PoolThreshold, if any.
func (a *Attribute) CurrentThreshold() *PoolThreshold {
	def := a.AttributeDef()
	if def == nil {
		return nil
	}
	max := a.Maximum()
	cur := a.Current()
	for _, threshold := range def.Thresholds {
		if cur <= threshold.Threshold(max) {
			return threshold
		}
	}
	return nil
}

// PointCost returns the number of points spent on this Attribute.
func (a *Attribute) PointCost() fixed.F64d4 {
	def := a.AttributeDef()
	if def == nil {
		return 0
	}
	var sm fixed.F64d4
	if a.Entity != nil {
		sm = a.Entity.Profile.AdjustedSizeModifier()
	}
	return def.ComputeCost(a.Entity, a.Adjustment, sm, a.CostReduction)
}
