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
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	costReductionAttributeKey  = "attribute"
	costReductionPercentageKey = "percentage"
)

// CostReductionType is the data type key for a CostReduction.
const CostReductionType = "cost_reduction"

var _ Feature = &CostReduction{}

// CostReduction holds a cost reduction. */
type CostReduction struct {
	Attribute  string
	Percentage fixed.F64d4
}

// NewCostReduction creates a new CostReduction. 'entity' may be nil.
func NewCostReduction(entity *Entity) *CostReduction {
	c := &CostReduction{Percentage: fixed.F64d4FromInt64(40)}
	list := AttributeDefsFor(entity).List()
	if len(list) != 0 {
		c.Attribute = list[0].ID()
	} else {
		c.Attribute = "st"
	}
	return c
}

// NewCostReductionFromJSON creates a new CostReduction from a JSON object.
func NewCostReductionFromJSON(data map[string]interface{}) *CostReduction {
	return &CostReduction{
		Attribute:  encoding.String(data[costReductionAttributeKey]),
		Percentage: encoding.Number(data[costReductionPercentageKey]),
	}
}

// ToJSON implements Feature.
func (c *CostReduction) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(featureTypeKey, c.DataType(), false, false)
	encoder.KeyedString(costReductionAttributeKey, c.Attribute, false, false)
	encoder.KeyedNumber(costReductionPercentageKey, c.Percentage, false)
	encoder.EndObject()
}

// CloneFeature implements Feature.
func (c *CostReduction) CloneFeature() Feature {
	clone := *c
	return &clone
}

// DataType implements Feature.
func (c *CostReduction) DataType() string {
	return CostReductionType
}

// FeatureKey implements Feature.
func (c *CostReduction) FeatureKey() string {
	return AttributeIDPrefix + c.Attribute
}

// FillWithNameableKeys implements Feature.
func (c *CostReduction) FillWithNameableKeys(set map[string]bool) {
	// Does nothing
}

// ApplyNameableKeys implements Feature.
func (c *CostReduction) ApplyNameableKeys(nameables map[string]string) {
	// Does nothing
}
