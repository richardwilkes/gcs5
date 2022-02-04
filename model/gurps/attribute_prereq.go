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
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

var _ Prereq = &AttributePrereq{}

// AttributePrereq holds a prerequisite for an attribute.
type AttributePrereq struct {
	Parent               *PrereqList      `json:"-"`
	Type                 prereq.Type      `json:"type"`
	CombinedWithCriteria criteria.String  `json:"combined_with,omitempty"`
	QualifierCriteria    criteria.Numeric `json:"qualifier,omitempty"`
	Which                string           `json:"which"`
	Has                  bool             `json:"has"`
}

// NewAttributePrereq creates a new AttributePrereq. 'entity' may be nil.
func NewAttributePrereq(entity *Entity) *AttributePrereq {
	return &AttributePrereq{
		Type: prereq.Attribute,
		CombinedWithCriteria: criteria.String{
			StringData: criteria.StringData{
				Compare: criteria.Is,
			},
		},
		QualifierCriteria: criteria.Numeric{
			NumericData: criteria.NumericData{
				Compare:   criteria.AtLeast,
				Qualifier: fixed.F64d4FromInt64(10),
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
func (a *AttributePrereq) FillWithNameableKeys(m map[string]string) {
}

// ApplyNameableKeys implements Prereq.
func (a *AttributePrereq) ApplyNameableKeys(m map[string]string) {
}

// Satisfied implements Prereq.
func (a *AttributePrereq) Satisfied(entity *Entity, exclude interface{}, buffer *xio.ByteBuffer, prefix string) bool {
	satisfied := false
	// TODO: Implement
	/*
	   boolean satisfied = mValueCompare.matches(character.getAttributeIntValue(mWhich) + (mCombinedWith != null ? character.getAttributeIntValue(mCombinedWith) : 0));
	   if (!has()) {
	       satisfied = !satisfied;
	   }
	   if (!satisfied && builder != null) {
	       Map<String, AttributeDef> attributes = character.getSheetSettings().getAttributes();
	       AttributeDef              def        = attributes.get(mWhich);
	       String                    text       = def != null ? def.getName() : "<unknown>";
	       if (mCombinedWith != null) {
	           def = attributes.get(mCombinedWith);
	           text += "+" + (def != null ? def.getName() : "<unknown>");
	       }
	       builder.append(MessageFormat.format(I18n.text("{0}{1} {2} which {3}\n"), prefix, getHasText(), text, mValueCompare.toString()));
	   }
	*/
	return satisfied
}
