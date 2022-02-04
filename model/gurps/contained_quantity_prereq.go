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
	"github.com/richardwilkes/toolbox/xio"
)

var _ Prereq = &ContainedQuantityPrereq{}

// ContainedQuantityPrereq holds a prerequisite for an equipment contained quantity.
type ContainedQuantityPrereq struct {
	Parent            *PrereqList      `json:"-"`
	Type              prereq.Type      `json:"type"`
	QualifierCriteria criteria.Numeric `json:"qualifier,omitempty"`
	Has               bool             `json:"has"`
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
func (c *ContainedQuantityPrereq) Satisfied(entity *Entity, exclude interface{}, buffer *xio.ByteBuffer, prefix string) bool {
	satisfied := false
	// TODO: Implement
	/*
	   boolean satisfied = false;
	   if (exclude instanceof Equipment equipment) {
	       satisfied = !equipment.canHaveChildren();
	       if (!satisfied) {
	           int qty = 0;
	           for (Row child : equipment.getChildren()) {
	               if (child instanceof Equipment) {
	                   qty += ((Equipment) child).getQuantity();
	               }
	           }
	           satisfied = mQuantityCompare.matches(qty);
	       }
	   }
	   if (!has()) {
	       satisfied = !satisfied;
	   }
	   if (!satisfied && builder != null) {
	       builder.append(MessageFormat.format(I18n.text("\n{0}{1} a contained quantity which {2}"), prefix, getHasText(), mQuantityCompare));
	   }
	   return satisfied;
	*/
	return satisfied
}
