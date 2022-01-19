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
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/richardwilkes/gcs/internal/id"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/eval/float64eval"
	"github.com/richardwilkes/toolbox/log/jot"
	xfs "github.com/richardwilkes/toolbox/xio/fs"
)

// ReservedAttributeDefIDs holds a list of IDs that aren't permitted for an AttributeDef.
var ReservedAttributeDefIDs = []string{"skill", "parry", "block", "dodge", "sm"}

// AttributeDefs holds a slice of AttributeDef.
type AttributeDefs []*AttributeDef

// AttributeDef holds the definition of an attribute.
type AttributeDef struct {
	ID                  string           `json:"id"`
	Type                string           `json:"type"`
	Name                string           `json:"name"`
	FullName            string           `json:"full_name,omitempty"`
	AttributeBase       string           `json:"attribute_base,omitempty"`
	CostPerPoint        int              `json:"cost_per_point,omitempty"`
	CostAdjPercentPerSM int              `json:"cost_adj_percent_per_sm,omitempty"`
	Order               int              `json:"-"`
	Thresholds          []*PoolThreshold `json:"thresholds,omitempty"`
}

// FactoryAttributeDefs returns the attribute factory settings.
func FactoryAttributeDefs() AttributeDefs {
	var defs AttributeDefs
	jot.FatalIfErr(xfs.LoadJSONFromFS(embeddedFS, "data/standard.attr", &defs))
	return defs
}

// UnmarshalJSON implements json.Unmarshaler. Loads the current format as well as older variants.
func (a *AttributeDefs) UnmarshalJSON(data []byte) error {
	var current []*AttributeDef
	if err := json.Unmarshal(data, &current); err != nil {
		var variants struct {
			JavaVersion []*AttributeDef `json:"attributes"`
		}
		if err2 := json.Unmarshal(data, &variants); err2 != nil {
			return err
		}
		*a = variants.JavaVersion
	} else {
		*a = current
	}
	set := make(map[string]bool)
	for _, one := range *a {
		one.ID = id.Sanitize(one.ID, false, ReservedAttributeDefIDs...)
		if set[one.ID] {
			return errs.New("duplicate ID in attributes: " + one.ID)
		}
		set[one.ID] = true
	}
	return nil
}

// SaveTo saves the AttributeDefs data to the specified file.
func (a AttributeDefs) SaveTo(filePath string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o750); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	return xfs.SaveJSONWithMode(filePath, a, true, 0o640)
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
func (a *AttributeDef) BaseValue(resolver eval.VariableResolver) float64 {
	result, err := float64eval.NewEvaluator(resolver, true).Evaluate(a.AttributeBase)
	if err != nil {
		jot.Warn(errs.NewWithCausef(err, "unable to resolve '%s'", a.AttributeBase))
		return 0
	}
	if value, ok := result.(float64); ok {
		return value
	}
	jot.Warn(errs.Newf("unable to resolve '%s' to a number", a.AttributeBase))
	return 0
}

// ComputeCost returns the value adjusted for a cost reduction.
func (a *AttributeDef) ComputeCost(entity *Entity, value float64, sizeModifier, costReduction int) int {
	cost := int(float64(a.CostPerPoint) * value)
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
