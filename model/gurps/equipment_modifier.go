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

import "github.com/richardwilkes/gcs/model/gurps/equipment"

// EquipmentModifier holds a modifier to a piece of Equipment.
type EquipmentModifier struct {
	Common
	TechLevel    string
	CostAmount   string
	WeightAmount string
	Features     []*Feature
	Children     []*EquipmentModifier
	CostType     equipment.ModifierCostType
	WeightType   equipment.ModifierWeightType
	Enabled      bool
}

// FillWithNameableKeys adds any nameable keys found in this AdvantageModifier to the provided map.
func (e *EquipmentModifier) FillWithNameableKeys(nameables map[string]string) {
	if e.Enabled {
		e.Common.FillWithNameableKeys(nameables)
		for _, one := range e.Features {
			one.FillWithNameableKeys(nameables)
		}
	}
}

// ApplyNameableKeys replaces any nameable keys found in this AdvantageModifier with the corresponding values in the provided map.
func (e *EquipmentModifier) ApplyNameableKeys(nameables map[string]string) {
	if e.Enabled {
		e.Common.ApplyNameableKeys(nameables)
		for _, one := range e.Features {
			one.ApplyNameableKeys(nameables)
		}
	}
}
