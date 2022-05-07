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
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64"
)

var _ Prereq = &AttributePrereq{}

// AttributePrereq holds a prerequisite for an attribute.
type AttributePrereq struct {
	Parent            *PrereqList      `json:"-"`
	Type              prereq.Type      `json:"type"`
	Has               bool             `json:"has"`
	CombinedWith      string           `json:"combined_with,omitempty"`
	QualifierCriteria criteria.Numeric `json:"qualifier,omitempty"`
	Which             string           `json:"which"`
}

// NewAttributePrereq creates a new AttributePrereq. 'entity' may be nil.
func NewAttributePrereq(entity *Entity) *AttributePrereq {
	return &AttributePrereq{
		Type: prereq.Attribute,
		QualifierCriteria: criteria.Numeric{
			NumericData: criteria.NumericData{
				Compare:   criteria.AtLeast,
				Qualifier: f64.From[fxp.DP](10),
			},
		},
		Which: AttributeIDFor(entity, gid.Strength),
		Has:   true,
	}
}

// Clone implements Prereq.
func (a *AttributePrereq) Clone(parent *PrereqList) Prereq {
	clone := *a
	clone.Parent = parent
	return &clone
}

// FillWithNameableKeys implements Prereq.
func (a *AttributePrereq) FillWithNameableKeys(_ map[string]string) {
}

// ApplyNameableKeys implements Prereq.
func (a *AttributePrereq) ApplyNameableKeys(_ map[string]string) {
}

// Satisfied implements Prereq.
func (a *AttributePrereq) Satisfied(entity *Entity, _ interface{}, tooltip *xio.ByteBuffer, prefix string) bool {
	value := entity.ResolveAttributeCurrent(a.Which)
	if a.CombinedWith != "" {
		value += entity.ResolveAttributeCurrent(a.CombinedWith)
	}
	satisfied := a.QualifierCriteria.Matches(value)
	if !a.Has {
		satisfied = !satisfied
	}
	if !satisfied && tooltip != nil {
		tooltip.WriteByte('\n')
		tooltip.WriteString(prefix)
		tooltip.WriteString(HasText(a.Has))
		tooltip.WriteByte(' ')
		tooltip.WriteString(entity.ResolveAttributeName(a.Which))
		if a.CombinedWith != "" {
			tooltip.WriteByte('+')
			tooltip.WriteString(entity.ResolveAttributeName(a.CombinedWith))
		}
		tooltip.WriteString(i18n.Text(" which "))
		tooltip.WriteString(a.QualifierCriteria.String())
	}
	return satisfied
}
