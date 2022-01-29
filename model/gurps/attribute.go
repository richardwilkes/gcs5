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
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	calcKey                 = "calc"
	attributeIDKey          = "attr_id"
	attributeAdjKey         = "adj"
	attributeDamageKey      = "damage"
	attributeCalcValueKey   = "value"
	attributeCalcPointsKey  = "points"
	attributeCalcCurrentKey = "current"
)

// AttributeIDPrefix is the prefix all references to attribute IDs should use.
const AttributeIDPrefix = "attr."

// Attribute holds the current state of an AttributeDef.
type Attribute struct {
	id            string
	adjustment    fixed.F64d4
	Bonus         fixed.F64d4
	CostReduction fixed.F64d4
	Damage        fixed.F64d4
}

// NewAttributeFromJSON creates a new Attribute from a JSON object.
func NewAttributeFromJSON(data map[string]interface{}) *Attribute {
	a := &Attribute{
		adjustment: encoding.Number(data[attributeAdjKey]),
		Damage:     encoding.Number(data[attributeDamageKey]),
	}
	a.SetID(encoding.String(data[attributeIDKey]))
	return a
}

// ToJSON emits this object as JSON.
func (a *Attribute) ToJSON(encoder *encoding.JSONEncoder, entity *Entity) {
	if def := a.AttributeDef(entity); def != nil {
		encoder.StartObject()
		encoder.KeyedString(attributeIDKey, a.id, false, false)
		encoder.KeyedNumber(attributeAdjKey, a.adjustment, false)
		if def.Type == attribute.Pool {
			encoder.KeyedNumber(attributeDamageKey, a.Damage, false)
		}

		// Emit calculated values for third parties
		encoder.Key(calcKey)
		encoder.StartObject()
		encoder.KeyedNumber(attributeCalcValueKey, a.Maximum(entity), false)
		if def.Type == attribute.Pool {
			encoder.KeyedNumber(attributeCalcCurrentKey, a.Current(entity), false)
		}
		encoder.KeyedNumber(attributeCalcPointsKey, a.PointCost(entity), false)
		encoder.EndObject()

		encoder.EndObject()
	}
}

// ID returns the ID.
func (a *Attribute) ID() string {
	return a.id
}

// SetID sets the ID, sanitizing it in the process (i.e. it may be changed from what you set -- read it back if you want
// to be sure of what it gets set to.
func (a *Attribute) SetID(value string) {
	a.id = id.Sanitize(value, false, ReservedIDs...)
}

// AttributeDef looks up the AttributeDef this Attribute references from the Entity. May return nil.
func (a *Attribute) AttributeDef(entity *Entity) *AttributeDef {
	return entity.SheetSettings.Attributes.Set[a.id]
}

// Maximum returns the maximum value of a pool or the adjusted attribute value for other types.
func (a *Attribute) Maximum(entity *Entity) fixed.F64d4 {
	def := a.AttributeDef(entity)
	if def == nil {
		return 0
	}
	max := def.BaseValue(entity) + a.adjustment + a.Bonus
	if def.Type != attribute.Decimal {
		max = max.Trunc()
	}
	return max
}

// Current returns the current value. Only valid for pools.
func (a *Attribute) Current(entity *Entity) fixed.F64d4 {
	return a.Maximum(entity) - a.Damage
}

// SetCurrent sets the current value.
func (a *Attribute) SetCurrent(entity *Entity, value fixed.F64d4) {
	if a.Current(entity) == value {
		return
	}
	if def := a.AttributeDef(entity); def != nil {
		a.adjustment = value - (def.BaseValue(entity) + a.Bonus)
	}
}

// CurrentThreshold return the current PoolThreshold, if any.
func (a *Attribute) CurrentThreshold(entity *Entity) *PoolThreshold {
	def := a.AttributeDef(entity)
	if def == nil {
		return nil
	}
	max := a.Maximum(entity)
	cur := a.Current(entity)
	for _, threshold := range def.Thresholds {
		if cur <= threshold.Threshold(max) {
			return threshold
		}
	}
	return nil
}

// PointCost returns the number of points spent on this Attribute.
func (a *Attribute) PointCost(entity *Entity) fixed.F64d4 {
	def := a.AttributeDef(entity)
	if def == nil {
		return 0
	}
	return def.ComputeCost(entity, a.adjustment, entity.Profile.AdjustedSizeModifier(), a.CostReduction)
}
