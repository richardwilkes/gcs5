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

	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/eval/f64d4eval"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// ReservedIDs holds a list of IDs that are reserved for internal use.
var ReservedIDs = []string{"skill", "parry", "block", "dodge", "sm"}

// AttributeDefData holds the AttributeDef data that is written to disk.
type AttributeDefData struct {
	DefID               string           `json:"id"`
	Type                attribute.Type   `json:"type"`
	Name                string           `json:"name"`
	FullName            string           `json:"full_name,omitempty"`
	AttributeBase       string           `json:"attribute_base"`
	CostPerPoint        fixed.F64d4      `json:"cost_per_point"`
	CostAdjPercentPerSM fixed.F64d4      `json:"cost_adj_percent_per_sm,omitempty"`
	Thresholds          []*PoolThreshold `json:"thresholds,omitempty"`
}

// AttributeDef holds the definition of an attribute.
type AttributeDef struct {
	AttributeDefData
	Order int `json:"-"`
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
func (a *AttributeDef) ComputeCost(entity *Entity, value, sizeModifier, costReduction fixed.F64d4) fixed.F64d4 {
	cost := value.Mul(a.CostPerPoint)
	if sizeModifier > 0 && a.CostAdjPercentPerSM > 0 && !(a.DefID == "hp" && entity.SheetSettings.DamageProgression == attribute.KnowingYourOwnStrength) {
		costReduction += sizeModifier.Mul(a.CostAdjPercentPerSM)
	}
	if costReduction > 0 {
		if costReduction > f64d4.Eighty {
			costReduction = f64d4.Eighty
		}
		cost = cost.Mul(f64d4.Hundred - costReduction)
	}
	return f64d4.Round(cost)
}
