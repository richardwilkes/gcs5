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
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

var _ Prereq = &ContainedQuantityPrereq{}

// ContainedQuantityPrereq holds a prerequisite for an equipment contained quantity.
type ContainedQuantityPrereq struct {
	Parent            *PrereqList      `json:"-"`
	Type              prereq.Type      `json:"type"`
	Has               bool             `json:"has"`
	QualifierCriteria criteria.Numeric `json:"qualifier,omitempty"`
}

// NewContainedQuantityPrereq creates a new ContainedQuantityPrereq.
func NewContainedQuantityPrereq() *ContainedQuantityPrereq {
	return &ContainedQuantityPrereq{
		Type: prereq.ContainedQuantity,
		QualifierCriteria: criteria.Numeric{
			NumericData: criteria.NumericData{
				Compare:   criteria.AtMost,
				Qualifier: fxp.One,
			},
		},
		Has: true,
	}
}

// Clone implements Prereq.
func (c *ContainedQuantityPrereq) Clone(parent *PrereqList) Prereq {
	clone := *c
	clone.Parent = parent
	return &clone
}

// FillWithNameableKeys implements Prereq.
func (c *ContainedQuantityPrereq) FillWithNameableKeys(m map[string]string) {
}

// ApplyNameableKeys implements Prereq.
func (c *ContainedQuantityPrereq) ApplyNameableKeys(m map[string]string) {
}

// Satisfied implements Prereq.
func (c *ContainedQuantityPrereq) Satisfied(_ *Entity, exclude interface{}, tooltip *xio.ByteBuffer, prefix string) bool {
	satisfied := false
	if eqp, ok := exclude.(*Equipment); ok {
		if satisfied = !eqp.Container(); !satisfied {
			var qty fixed.F64d4
			for _, child := range eqp.Children {
				qty += child.Quantity
			}
			satisfied = c.QualifierCriteria.Matches(qty)
		}
	}
	if !c.Has {
		satisfied = !satisfied
	}
	if !satisfied && tooltip != nil {
		tooltip.WriteByte('\n')
		tooltip.WriteString(prefix)
		tooltip.WriteString(HasText(c.Has))
		tooltip.WriteString(i18n.Text(" a contained quantity which "))
		tooltip.WriteString(c.QualifierCriteria.String())
	}
	return satisfied
}
