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
	"encoding/binary"
	"hash"

	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

// PoolThreshold holds a point within an attribute pool where changes in state occur.
type PoolThreshold struct {
	State       string                  `json:"state"`
	Explanation string                  `json:"explanation,omitempty"`
	Multiplier  f64d4.Int               `json:"multiplier"`
	Divisor     f64d4.Int               `json:"divisor"`
	Addition    f64d4.Int               `json:"addition,omitempty"`
	Ops         []attribute.ThresholdOp `json:"ops,omitempty"`
	// TODO: Turn the Multiplier, Divisor & Addition fields into an expression widget instead
}

// Clone a copy of this.
func (p *PoolThreshold) Clone() *PoolThreshold {
	clone := *p
	if p.Ops != nil {
		clone.Ops = make([]attribute.ThresholdOp, len(p.Ops))
		copy(clone.Ops, p.Ops)
	}
	return &clone
}

// Threshold returns the threshold value for the given maximum.
func (p *PoolThreshold) Threshold(max f64d4.Int) f64d4.Int {
	divisor := p.Divisor //nolint:ifshort // bad recommendation
	if divisor == 0 {
		divisor = f64d4.One
	}
	// TODO: Check that rounding here is correct for our purposes
	return (max.Mul(p.Multiplier).Div(divisor) + p.Addition).Round()
}

// ContainsOp returns true if this PoolThreshold contains the specified ThresholdOp.
func (p *PoolThreshold) ContainsOp(op attribute.ThresholdOp) bool {
	for _, one := range p.Ops {
		if one == op {
			return true
		}
	}
	return false
}

func (p *PoolThreshold) crc64(h hash.Hash64) {
	h.Write([]byte(p.State))
	h.Write([]byte(p.Explanation))
	var buffer [8]byte
	binary.LittleEndian.PutUint64(buffer[:], uint64(p.Multiplier))
	h.Write(buffer[:])
	binary.LittleEndian.PutUint64(buffer[:], uint64(p.Divisor))
	h.Write(buffer[:])
	binary.LittleEndian.PutUint64(buffer[:], uint64(p.Addition))
	h.Write(buffer[:])
	binary.LittleEndian.PutUint64(buffer[:], uint64(len(p.Ops)))
	h.Write(buffer[:])
	for _, one := range p.Ops {
		h.Write([]byte{byte(one)})
	}
}
