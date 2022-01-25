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

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/enums/units"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const containedWeightReductionKey = "reduction"

const (
	// ContainedWeightType is the data type key for a ContainedWeight.
	ContainedWeightType = "contained_weight_reduction"
	// ContainedWeightFeatureKey is the key used in the Feature map for things this Feature applies to.
	ContainedWeightFeatureKey = "equipment.weight.sum"
)

var _ Feature = &ContainedWeight{}

// ContainedWeight holds a cost reduction. */
type ContainedWeight struct {
	Reduction string
}

// NewContainedWeight creates a new ContainedWeight.
func NewContainedWeight() *ContainedWeight {
	return &ContainedWeight{Reduction: "0%"}
}

// NewContainedWeightFromJSON creates a new ContainedWeight from a JSON object.
func NewContainedWeightFromJSON(data map[string]interface{}) *ContainedWeight {
	return &ContainedWeight{Reduction: strings.TrimSpace(encoding.String(data[containedWeightReductionKey]))}
}

// ToJSON implements Feature.
func (c *ContainedWeight) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(featureTypeKey, c.DataType(), false, false)
	encoder.KeyedString(containedWeightReductionKey, c.Reduction, true, true)
	encoder.EndObject()
}

// CloneFeature implements Feature.
func (c *ContainedWeight) CloneFeature() Feature {
	clone := *c
	return &clone
}

// DataType implements Feature.
func (c *ContainedWeight) DataType() string {
	return ContainedWeightType
}

// FeatureKey implements Feature.
func (c *ContainedWeight) FeatureKey() string {
	return ContainedWeightFeatureKey
}

// FillWithNameableKeys implements Feature.
func (c *ContainedWeight) FillWithNameableKeys(_ map[string]string) {
	// Does nothing
}

// ApplyNameableKeys implements Feature.
func (c *ContainedWeight) ApplyNameableKeys(_ map[string]string) {
	// Does nothing
}

// Normalize implements Feature.
func (c *ContainedWeight) Normalize() {
	// Unused
}

// IsPercentageReduction returns true if this is a percentage reduction and not a fixed amount.
func (c *ContainedWeight) IsPercentageReduction() bool {
	return strings.HasSuffix(c.Reduction, "%")
}

// PercentageReduction returns the percentage (where 1% is 1, not 0.01) the weight should be reduced by. Will return 0 if
// this is not a percentage.
func (c *ContainedWeight) PercentageReduction() fixed.F64d4 {
	if !c.IsPercentageReduction() {
		return 0
	}
	return fixed.F64d4FromStringForced(c.Reduction[:len(c.Reduction)-1])
}

// FixedReduction returns the fixed amount the weight should be reduced by. Will return 0 if this is a percentage.
func (c *ContainedWeight) FixedReduction(defUnits units.Weight) measure.Weight {
	if c.IsPercentageReduction() {
		return 0
	}
	return measure.WeightFromStringForced(c.Reduction, defUnits)
}
