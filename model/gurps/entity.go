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
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/xio"
)

var _ eval.VariableResolver = &Entity{}

// Entity holds the base information for various types of entities: PC, NPC, Creature, etc.
type Entity struct {
	Profile       *Profile
	SheetSettings *SheetSettings
	featureMap    map[string][]*Feature
}

// AddDRBonusesFor locates any active DR bonuses and adds them to the map. If 'drMap' isn't nil, it will be returned.
func (e *Entity) AddDRBonusesFor(id string, tooltip *xio.ByteBuffer, drMap map[string]int) map[string]int {
	if drMap == nil {
		drMap = make(map[string]int)
	}
	if list, exists := e.featureMap[strings.ToLower(id)]; exists {
		for _, one := range list {
			if one.Type == feature.DRBonus {
				drBonus := one.Self.(*DRBonus)
				drMap[strings.ToLower(drBonus.Specialization)] += int(drBonus.LeveledAmount.AdjustedAmount().AsInt64())
				drBonus.AddToTooltip(tooltip)
			}
		}
	}
	return drMap
}

// ResolveVariable implements eval.VariableResolver.
func (e *Entity) ResolveVariable(variableName string) string {
	// TODO implement me
	return variableName
}

// PreservesUserDesc returns true if the user description field should be preserved when written to disk. Normally, only
// character sheets should return true for this.
func (e *Entity) PreservesUserDesc() bool {
	// TODO: Implement... should only return true for sheets
	return true
}
