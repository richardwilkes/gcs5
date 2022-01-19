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
	"context"

	"github.com/goccy/go-json"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Attribute holds the current state of an AttributeDef.
type Attribute struct {
	AttributeStorage
	Bonus         fixed.F64d4
	CostReduction int
}

// AttributeStorage defines the current Attribute data format.
type AttributeStorage struct {
	ID         string      `json:"attr_id"`
	Adjustment fixed.F64d4 `json:"adj,omitempty"`
	Damage     int         `json:"damage,omitempty"`
}

// MarshalJSON implements json.MarshalerContext.
func (a *Attribute) MarshalJSON(ctx context.Context) ([]byte, error) {
	entity, ok := ctx.Value(EntityCtxKey).(*Entity)
	if !ok {
		return nil, errs.New("missing context data")
	}
	def := entity.SheetSettings.Attributes.Lookup(a.ID)
	if def == nil {
		return nil, errs.New("reference to undefined attribute: " + a.ID)
	}
	points := def.ComputeCost(entity, a.Adjustment, entity.Profile.SM(), a.CostReduction)
	if def.Type == PoolAttributeType {
		var output struct {
			AttributeStorage `json:",inline"`
			Calc             struct {
				Value   fixed.F64d4 `json:"value"`
				Current int         `json:"current"`
				Points  int         `json:"points,omitempty"`
			} `json:"calc"`
		}
		output.AttributeStorage = a.AttributeStorage
		output.Calc.Points = points
		output.Calc.Value = (def.BaseValue(entity) + a.Adjustment + a.Bonus).Trunc()
		output.Calc.Current = int((output.Calc.Value - fixed.F64d4FromInt64(int64(a.Damage))).AsInt64())
		return json.MarshalContext(ctx, &output)
	}
	var output struct {
		AttributeStorage `json:",inline"`
		Calc             struct {
			Value  fixed.F64d4 `json:"value"`
			Points int         `json:"points,omitempty"`
		} `json:"calc"`
	}
	output.AttributeStorage = a.AttributeStorage
	output.Damage = 0
	output.Calc.Points = points
	output.Calc.Value = def.BaseValue(entity) + a.Adjustment + a.Bonus
	if def.Type == IntegerAttributeType {
		output.Calc.Value = output.Calc.Value.Trunc()
	}
	return json.MarshalContext(ctx, &output)
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Attribute) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &a.AttributeStorage)
}
