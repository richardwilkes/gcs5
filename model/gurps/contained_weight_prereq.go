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
	"github.com/richardwilkes/toolbox/i18n"
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

// PrereqType implements Prereq.
func (c *ContainedWeightPrereq) PrereqType() prereq.Type {
	return c.Type
}

// ParentList implements Prereq.
func (c *ContainedWeightPrereq) ParentList() *PrereqList {
	return c.Parent
}

// Clone implements Prereq.
func (c *ContainedWeightPrereq) Clone(parent *PrereqList) Prereq {
	clone := *c
	clone.Parent = parent
	return &clone
}

// FillWithNameableKeys implements Prereq.
func (c *ContainedWeightPrereq) FillWithNameableKeys(_ map[string]string) {
}

// ApplyNameableKeys implements Prereq.
func (c *ContainedWeightPrereq) ApplyNameableKeys(_ map[string]string) {
}

// Satisfied implements Prereq.
func (c *ContainedWeightPrereq) Satisfied(entity *Entity, exclude any, tooltip *xio.ByteBuffer, prefix string) bool {
	satisfied := false
	if eqp, ok := exclude.(*Equipment); ok {
		if satisfied = !eqp.Container(); !satisfied {
			units := SheetSettingsFor(entity).DefaultWeightUnits
			weight := eqp.ExtendedWeight(false, units) - eqp.AdjustedWeight(false, units)
			satisfied = c.WeightCriteria.Matches(weight)
		}
	}
	if !c.Has {
		satisfied = !satisfied
	}
	if !satisfied && tooltip != nil {
		tooltip.WriteString(prefix)
		tooltip.WriteString(HasText(c.Has))
		tooltip.WriteString(i18n.Text(" a contained weight which "))
		tooltip.WriteString(c.WeightCriteria.String())
	}
	return satisfied
}
