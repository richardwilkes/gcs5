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

// HitLocations holds a set of hit locations.
type HitLocations struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Roll      string      `json:"roll"`
	Locations []*Location `json:"locations"`
}

// Location holds a single hit location.
type Location struct {
	ID          string       `json:"id"`
	ChoiceName  string       `json:"choice_name"`
	TableName   string       `json:"table_name"`
	Slots       int          `json:"slots"`
	HitPenalty  int          `json:"hit_penalty"`
	DRBonus     int          `json:"dr_bonus"`
	Description string       `json:"description"`
	Calc        LocationCalc `json:"calc"`
}

// LocationCalc holds values GCS calculates for a Location, but that we want to be present in any json output so that
// other uses of the data don't have to replicate the code to calculate it.
type LocationCalc struct {
	RollRange string `json:"roll_range"`
}

// FactoryHitLocations returns the hit location factory settings.
func FactoryHitLocations() *HitLocations {
	// TODO: Fill
	return &HitLocations{}
}
