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

package feature

import (
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
)

// AttributeIDPrefix is the prefix all references to attribute IDs should use.
const AttributeIDPrefix = "attr."

var _ Feature = &CostReduction{}

// CostReduction holds the data for a cost reduction.
type CostReduction struct {
	Type       Type      `json:"type"`
	Attribute  string    `json:"attribute,omitempty"`
	Percentage f64d4.Int `json:"percentage,omitempty"`
}

// NewCostReduction creates a new CostReduction.
func NewCostReduction(attrID string) *CostReduction {
	return &CostReduction{
		Type:       CostReductionType,
		Attribute:  attrID,
		Percentage: fxp.Forty,
	}
}

// Clone implements Feature.
func (c *CostReduction) Clone() Feature {
	other := *c
	return &other
}

// FeatureMapKey implements Feature.
func (c *CostReduction) FeatureMapKey() string {
	return AttributeIDPrefix + c.Attribute
}

// FillWithNameableKeys implements Feature.
func (c *CostReduction) FillWithNameableKeys(_ map[string]string) {
}

// ApplyNameableKeys implements Feature.
func (c *CostReduction) ApplyNameableKeys(_ map[string]string) {
}
