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

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/eval/f64d4eval"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// ReservedAttributeDefIDs holds a list of IDs that aren't permitted for an AttributeDef.
var ReservedAttributeDefIDs = []string{"skill", "parry", "block", "dodge", "sm"}

// AttributeDef holds the definition of an attribute.
type AttributeDef struct {
	ID                  string           `json:"id"`
	Type                AttributeType    `json:"type"`
	Name                string           `json:"name"`
	FullName            string           `json:"full_name,omitempty"`
	AttributeBase       string           `json:"attribute_base,omitempty"`
	CostPerPoint        int              `json:"cost_per_point,omitempty"`
	CostAdjPercentPerSM int              `json:"cost_adj_percent_per_sm,omitempty"`
	Order               int              `json:"-"`
	Thresholds          []*PoolThreshold `json:"thresholds,omitempty"`
}

// CombinedName returns the combined FullName and Name, as appropriate.
func (a *AttributeDef) CombinedName() string {
	full := strings.TrimSpace(a.FullName)
	name := strings.TrimSpace(a.Name)
	if full == "" {
		return name
	}
	if name == "" || name == full {
		return full
	}
	return full + " (" + name + ")"
}

// Primary returns true if the base value is a non-derived value.
func (a *AttributeDef) Primary() bool {
	_, err := strconv.ParseInt(strings.TrimSpace(a.AttributeBase), 10, 64)
	return err == nil
}

// BaseValue returns the resolved base value.
func (a *AttributeDef) BaseValue(resolver eval.VariableResolver) fixed.F64d4 {
	result, err := f64d4eval.NewEvaluator(resolver, true).Evaluate(a.AttributeBase)
	if err != nil {
		jot.Warn(errs.NewWithCausef(err, "unable to resolve '%s'", a.AttributeBase))
		return 0
	}
	if value, ok := result.(fixed.F64d4); ok {
		return value
	}
	jot.Warn(errs.Newf("unable to resolve '%s' to a number", a.AttributeBase))
	return 0
}

// ComputeCost returns the value adjusted for a cost reduction.
func (a *AttributeDef) ComputeCost(entity *Entity, value fixed.F64d4, sizeModifier, costReduction int) int {
	cost := int(value.Mul(fixed.F64d4FromInt64(int64(a.CostPerPoint))).AsInt64())
	if sizeModifier > 0 && a.CostAdjPercentPerSM > 0 && !(a.ID == "hp" && entity.SheetSettings.DamageProgression == KnowingYourOwnStrength) {
		costReduction += sizeModifier * a.CostAdjPercentPerSM
		if costReduction < 0 {
			costReduction = 0
		} else if costReduction > 80 {
			costReduction = 80
		}
	}
	if costReduction != 0 {
		cost *= 100 - costReduction
		rem := cost % 100
		cost /= 100
		if rem > 49 {
			cost++
		} else if rem < -50 {
			cost--
		}
	}
	return cost
}
