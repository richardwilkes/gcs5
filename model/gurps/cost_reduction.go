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
	"github.com/richardwilkes/gcs/model/gurps/feature"
)

// CostReduction holds the data for a cost reduction.
type CostReduction struct {
	Feature
	Attribute  string `json:"attribute,omitempty"`
	Percentage int    `json:"percentage,omitempty"`
}

// NewCostReduction creates a new CostReduction.
func NewCostReduction(entity *Entity) *CostReduction {
	c := &CostReduction{
		Feature: Feature{
			Type: feature.CostReduction,
		},
		Attribute:  DefaultAttributeIDFor(entity),
		Percentage: 40,
	}
	c.Self = c
	return c
}

func (c *CostReduction) featureMapKey() string {
	return AttributeIDPrefix + c.Attribute
}
