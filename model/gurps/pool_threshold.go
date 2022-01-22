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
	"sort"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	poolThresholdStateKey       = "state"
	poolThresholdExplanationKey = "explanation"
	poolThresholdMultiplierKey  = "multiplier"
	poolThresholdDivisorKey     = "divisor"
	poolThresholdAdditionKey    = "addition"
	poolThresholdOpsKey         = "ops"
)

// PoolThreshold holds a point within an attribute pool where changes in state occur.
type PoolThreshold struct {
	State       string
	Explanation string
	Multiplier  fixed.F64d4
	Divisor     fixed.F64d4
	Addition    fixed.F64d4
	Ops         []ThresholdOp
	// TODO: Turn the Multiplier, Divisor & Addition fields into an expression field instead
}

// NewPoolThresholdFromJSON creates a new PoolThreshold from a JSON object.
func NewPoolThresholdFromJSON(data map[string]interface{}) *PoolThreshold {
	p := &PoolThreshold{
		State:       encoding.String(data[poolThresholdStateKey]),
		Explanation: encoding.String(data[poolThresholdExplanationKey]),
		Multiplier:  encoding.Number(data[poolThresholdMultiplierKey]),
		Divisor:     encoding.Number(data[poolThresholdDivisorKey]),
		Addition:    encoding.Number(data[poolThresholdAdditionKey]),
	}
	ops := encoding.Array(data[poolThresholdOpsKey])
	if len(ops) != 0 {
		p.Ops = make([]ThresholdOp, 0, len(ops))
		for _, one := range ops {
			p.Ops = append(p.Ops, ThresholdOpFromString(encoding.String(one)))
		}
	}
	return p
}

// ToJSON emits this object as JSON.
func (p *PoolThreshold) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(poolThresholdStateKey, p.State, false, false)
	encoder.KeyedString(poolThresholdExplanationKey, p.Explanation, true, true)
	encoder.KeyedNumber(poolThresholdMultiplierKey, p.Multiplier, false)
	encoder.KeyedNumber(poolThresholdDivisorKey, p.Divisor, false)
	encoder.KeyedNumber(poolThresholdAdditionKey, p.Addition, true)
	if len(p.Ops) != 0 {
		encoder.Key(poolThresholdOpsKey)
		encoder.StartArray()
		sort.Slice(p.Ops, func(i, j int) bool { return p.Ops[i] < p.Ops[j] })
		for _, op := range p.Ops {
			encoder.String(op.Key())
		}
		encoder.EndArray()
	}
	encoder.EndObject()
}

// Threshold returns the threshold value for the given maximum.
func (p *PoolThreshold) Threshold(max fixed.F64d4) fixed.F64d4 {
	divisor := p.Divisor
	if divisor == 0 {
		divisor = f64d4.One
	}
	// TODO: Check that rounding here is correct for our purposes
	return f64d4.Round(max.Mul(p.Multiplier).Div(divisor) + p.Addition)
}
