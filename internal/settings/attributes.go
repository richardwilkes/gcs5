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

package settings

// Attribute holds the definition of an attribute.
type Attribute struct {
	ID                  string       `json:"id"`
	Type                string       `json:"type"`
	Name                string       `json:"name"`
	FullName            string       `json:"full_name"`
	AttributeBase       string       `json:"attribute_base"`
	CostPerPoint        int          `json:"cost_per_point"`
	CostAdjPercentPerSm int          `json:"cost_adj_percent_per_sm"`
	Thresholds          []*Threshold `json:"thresholds,omitempty"`
}

// Threshold holds a point within an attribute pool where changes in state occur.
type Threshold struct {
	State       string   `json:"state"`
	Explanation string   `json:"explanation"`
	Multiplier  int      `json:"multiplier"`
	Divisor     int      `json:"divisor"`
	Addition    int      `json:"addition"`
	Ops         []string `json:"ops"`
}

// FactoryAttributes returns the attribute factory settings.
func FactoryAttributes() []*Attribute {
	// TODO: Fill
	return nil
}
