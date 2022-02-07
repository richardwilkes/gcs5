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
	"github.com/richardwilkes/gcs/model/criteria"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/toolbox/xio"
)

var _ Prereq = &ContainedWeightPrereq{}

// ContainedWeightPrereq holds a prerequisite for an equipment contained weight.
type ContainedWeightPrereq struct {
	Parent         *PrereqList     `json:"-"`
	Type           prereq.Type     `json:"type"`
	Has            bool            `json:"has"`
	WeightCriteria criteria.Weight `json:"qualifier,omitempty"`
}

// NewContainedWeightPrereq creates a new ContainedWeightPrereq.
func NewContainedWeightPrereq(entity *Entity) *ContainedWeightPrereq {
	return &ContainedWeightPrereq{
		Type: prereq.ContainedWeight,
		WeightCriteria: criteria.Weight{
			WeightData: criteria.WeightData{
				Compare:   criteria.AtMost,
				Qualifier: measure.WeightFromInt(5, SheetSettingsFor(entity).DefaultWeightUnits),
			},
		},
		Has: true,
	}
}

// Clone implements Prereq.
func (c *ContainedWeightPrereq) Clone(parent *PrereqList) Prereq {
	clone := *c
	clone.Parent = parent
	return &clone
}

// FillWithNameableKeys implements Prereq.
func (c *ContainedWeightPrereq) FillWithNameableKeys(m map[string]string) {
}

// ApplyNameableKeys implements Prereq.
func (c *ContainedWeightPrereq) ApplyNameableKeys(m map[string]string) {
}

// Satisfied implements Prereq.
func (c *ContainedWeightPrereq) Satisfied(entity *Entity, exclude interface{}, buffer *xio.ByteBuffer, prefix string) bool {
	satisfied := false
	// TODO: Implement
	/*
	   boolean satisfied = false;
	   if (exclude instanceof Equipment equipment) {
	       satisfied = !equipment.canHaveChildren();
	       if (!satisfied) {
	           WeightValue weight = new WeightValue(equipment.getExtendedWeight(false));
	           weight.subtract(equipment.getAdjustedWeight(false));
	           satisfied = mWeightCompare.matches(weight);
	       }
	   }
	   if (!has()) {
	       satisfied = !satisfied;
	   }
	   if (!satisfied && builder != null) {
	       builder.append(MessageFormat.format(I18n.text("\n{0}{1} a contained weight which {2}"), prefix, getHasText(), mWeightCompare));
	   }
	   return satisfied;
	*/
	return satisfied
}
