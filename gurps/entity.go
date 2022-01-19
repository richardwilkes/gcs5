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
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/xio"
)

var _ eval.VariableResolver = &Entity{}

// Entity holds the base information for various types of entities: PC, NPC, Creature, etc.
type Entity struct {
	Profile       PCProfile
	SheetSettings *SheetSettings
}

// AddDRBonusesFor locates any active DR bonuses and adds them to the map. If 'drMap' isn't nil, it will be returned.
func (e *Entity) AddDRBonusesFor(id string, tooltip *xio.ByteBuffer, drMap map[string]int) map[string]int {
	if drMap == nil {
		drMap = make(map[string]int)
	}
	// TODO: Implement
	/*
	   List<Feature> list = mFeatureMap.get(id.toLowerCase());
	   if (list != null) {
	       for (Feature feature : list) {
	           if (feature instanceof DRBonus bonus) {
	               String  specialization = bonus.getSpecialization();
	               int     amt            = bonus.getAmount().getIntegerAdjustedAmount();
	               Integer value          = dr.get(specialization);
	               if (value == null) {
	                   value = Integer.valueOf(amt);
	               } else {
	                   value = Integer.valueOf(value.intValue() + amt);
	               }
	               dr.put(specialization, value);
	               bonus.addToToolTip(tooltip);
	           }
	       }
	   }
	*/
	return drMap
}

// ResolveVariable implements eval.VariableResolver.
func (e *Entity) ResolveVariable(variableName string) string {
	// TODO implement me
	return variableName
}
