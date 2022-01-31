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
	"strings"

	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// ContainedWeightReduction holds the data for a weight reduction that can be applied to a container's contents.
type ContainedWeightReduction struct {
	Feature
	Reduction string `json:"reduction"`
}

// NewContainedWeightReduction creates a new ContainedWeightReduction.
func NewContainedWeightReduction() *ContainedWeightReduction {
	c := &ContainedWeightReduction{
		Feature: Feature{
			Type: feature.ContainedWeightReduction,
		},
		Reduction: "0%",
	}
	c.Self = c
	return c
}

func (c *ContainedWeightReduction) featureMapKey() string {
	return "equipment.weight.sum"
}

// IsPercentageReduction returns true if this is a percentage reduction and not a fixed amount.
func (c *ContainedWeightReduction) IsPercentageReduction() bool {
	return strings.HasSuffix(c.Reduction, "%")
}

// PercentageReduction returns the percentage (where 1% is 1, not 0.01) the weight should be reduced by. Will return 0 if
// this is not a percentage.
func (c *ContainedWeightReduction) PercentageReduction() fixed.F64d4 {
	if !c.IsPercentageReduction() {
		return 0
	}
	return fixed.F64d4FromStringForced(c.Reduction[:len(c.Reduction)-1])
}

// FixedReduction returns the fixed amount the weight should be reduced by. Will return 0 if this is a percentage.
func (c *ContainedWeightReduction) FixedReduction(defUnits measure.WeightUnits) measure.Weight {
	if c.IsPercentageReduction() {
		return 0
	}
	return measure.WeightFromStringForced(c.Reduction, defUnits)
}
