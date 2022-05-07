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
	"strconv"
	"strings"

	"github.com/richardwilkes/gcs/model/crc"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64"
)

// ReservedIDs holds a list of IDs that are reserved for internal use.
var ReservedIDs = []string{gid.Skill, gid.Parry, gid.Block, "dodge", "sm"}

// AttributeDef holds the definition of an attribute.
type AttributeDef struct {
	DefID               string           `json:"id"`
	Type                attribute.Type   `json:"type"`
	Name                string           `json:"name"`
	FullName            string           `json:"full_name,omitempty"`
	AttributeBase       string           `json:"attribute_base"`
	CostPerPoint        fxp.Int          `json:"cost_per_point"`
	CostAdjPercentPerSM fxp.Int          `json:"cost_adj_percent_per_sm,omitempty"`
	Thresholds          []*PoolThreshold `json:"thresholds,omitempty"`
	Order               int              `json:"-"`
}

// Clone a copy of this.
func (a *AttributeDef) Clone() *AttributeDef {
	clone := *a
	if a.Thresholds != nil {
		clone.Thresholds = make([]*PoolThreshold, len(a.Thresholds))
		for i, one := range a.Thresholds {
			clone.Thresholds[i] = one.Clone()
		}
	}
	return &clone
}

// ID returns the ID.
func (a *AttributeDef) ID() string {
	return a.DefID
}

// SetID sets the ID, sanitizing it in the process (i.e. it may be changed from what you set -- read it back if you want
// to be sure of what it gets set to.
func (a *AttributeDef) SetID(value string) {
	a.DefID = id.Sanitize(value, false, ReservedIDs...)
}

// ResolveFullName returns the full name, using the short name if full name is empty.
func (a *AttributeDef) ResolveFullName() string {
	if a.FullName == "" {
		return a.Name
	}
	return a.FullName
}

// CombinedName returns the combined FullName and Name, as appropriate.
func (a *AttributeDef) CombinedName() string {
	if a.FullName == "" {
		return a.Name
	}
	if a.Name == "" || a.Name == a.FullName {
		return a.FullName
	}
	return a.FullName + " (" + a.Name + ")"
}

// Primary returns true if the base value is a non-derived value.
func (a *AttributeDef) Primary() bool {
	_, err := strconv.ParseInt(strings.TrimSpace(a.AttributeBase), 10, 64)
	return err == nil
}

// BaseValue returns the resolved base value.
func (a *AttributeDef) BaseValue(resolver eval.VariableResolver) fxp.Int {
	return fxp.EvaluateToNumber(a.AttributeBase, resolver)
}

// ComputeCost returns the value adjusted for a cost reduction.
func (a *AttributeDef) ComputeCost(entity *Entity, value, costReduction fxp.Int, sizeModifier int) fxp.Int {
	cost := value.Mul(a.CostPerPoint)
	if sizeModifier > 0 && a.CostAdjPercentPerSM > 0 && !(a.DefID == "hp" && entity.SheetSettings.DamageProgression == attribute.KnowingYourOwnStrength) {
		costReduction += f64.From[fxp.DP](sizeModifier).Mul(a.CostAdjPercentPerSM)
	}
	if costReduction > 0 {
		if costReduction > fxp.Eighty {
			costReduction = fxp.Eighty
		}
		cost = cost.Mul(fxp.Hundred - costReduction).Div(fxp.Hundred)
	}
	return cost.Round()
}

func (a *AttributeDef) crc64(c uint64) uint64 {
	c = crc.String(c, a.DefID)
	c = crc.Byte(c, byte(a.Type))
	c = crc.String(c, a.Name)
	c = crc.String(c, a.FullName)
	c = crc.String(c, a.AttributeBase)
	c = crc.Number(c, a.CostPerPoint)
	c = crc.Number(c, a.CostAdjPercentPerSM)
	c = crc.Number(c, len(a.Thresholds))
	for _, one := range a.Thresholds {
		c = one.crc64(c)
	}
	return c
}
