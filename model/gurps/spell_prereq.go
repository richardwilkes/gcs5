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
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/gcs/model/gurps/spell"
	"github.com/richardwilkes/toolbox/xio"
)

var _ Prereq = &SpellPrereq{}

// SpellPrereq holds a prerequisite for a spell.
type SpellPrereq struct {
	Parent            *PrereqList          `json:"-"`
	Type              prereq.Type          `json:"type"`
	SubType           spell.ComparisonType `json:"sub_type"`
	Has               bool                 `json:"has"`
	QualifierCriteria criteria.String      `json:"qualifier,omitempty"`
	QuantityCriteria  criteria.Numeric     `json:"quantity,omitempty"`
}

// NewSpellPrereq creates a new SpellPrereq.
func NewSpellPrereq() *SpellPrereq {
	return &SpellPrereq{
		Type:    prereq.Spell,
		SubType: spell.Name,
		QualifierCriteria: criteria.String{
			StringData: criteria.StringData{
				Compare: criteria.Is,
			},
		},
		QuantityCriteria: criteria.Numeric{
			NumericData: criteria.NumericData{
				Compare:   criteria.AtLeast,
				Qualifier: fxp.One,
			},
		},
		Has: true,
	}
}

// Clone implements Prereq.
func (s *SpellPrereq) Clone(parent *PrereqList) Prereq {
	clone := *s
	clone.Parent = parent
	return &clone
}

// FillWithNameableKeys implements Prereq.
func (s *SpellPrereq) FillWithNameableKeys(m map[string]string) {
	if s.SubType.UsesStringCriteria() {
		nameables.Extract(s.QualifierCriteria.Qualifier, m)
	}
}

// ApplyNameableKeys implements Prereq.
func (s *SpellPrereq) ApplyNameableKeys(m map[string]string) {
	if s.SubType.UsesStringCriteria() {
		s.QualifierCriteria.Qualifier = nameables.Apply(s.QualifierCriteria.Qualifier, m)
	}
}

// Satisfied implements Prereq.
func (s *SpellPrereq) Satisfied(entity *Entity, exclude interface{}, buffer *xio.ByteBuffer, prefix string) bool {
	satisfied := false
	// TODO: Implement
	/*
	   Set<String> colleges  = new HashSet<>();
	   String      techLevel = null;
	   int         count     = 0;
	   boolean     satisfied;
	   if (exclude instanceof Spell) {
	       techLevel = ((Spell) exclude).getTechLevel();
	   }
	   for (Spell spell : character.getSpellsIterator()) {
	       if (exclude != spell && spell.getPoints() > 0) {
	           boolean ok;
	           if (techLevel != null) {
	               String otherTL = spell.getTechLevel();

	               ok = otherTL == null || techLevel.equals(otherTL);
	           } else {
	               ok = true;
	           }
	           if (ok) {
	               if (KEY_NAME.equals(mType)) {
	                   if (mStringCriteria.matches(spell.getName())) {
	                       count++;
	                   }
	               } else if (KEY_ANY.equals(mType)) {
	                   count++;
	               } else if (KEY_CATEGORY.equals(mType)) {
	                   for (String category : spell.getCategories()) {
	                       if (mStringCriteria.matches(category)) {
	                           count++;
	                           break;
	                       }
	                   }
	               } else if (KEY_COLLEGE.equals(mType)) {
	                   for (String college : spell.getColleges()) {
	                       if (mStringCriteria.matches(college)) {
	                           count++;
	                           break;
	                       }
	                   }
	               } else if (Objects.equals(mType, KEY_COLLEGE_COUNT)) {
	                   colleges.addAll(spell.getColleges());
	               }
	           }
	       }
	   }

	   if (Objects.equals(mType, KEY_COLLEGE_COUNT)) {
	       count = colleges.size();
	   }

	   satisfied = mQuantityCriteria.matches(count);
	   if (!has()) {
	       satisfied = !satisfied;
	   }
	   if (!satisfied && builder != null) {
	       String oneSpell       = I18n.text("spell");
	       String multipleSpells = I18n.text("spells");
	       if (Objects.equals(mType, KEY_NAME)) {
	           builder.append(MessageFormat.format(I18n.text("\n{0}{1} {2} {3} whose name {4}"), prefix, getHasText(), mQuantityCriteria.toString(""), mQuantityCriteria.getQualifier() == 1 ? oneSpell : multipleSpells, mStringCriteria.toString()));
	       } else if (Objects.equals(mType, KEY_ANY)) {
	           builder.append(MessageFormat.format(I18n.text("\n{0}{1} {2} {3} of any kind"), prefix, getHasText(), mQuantityCriteria.toString(""), mQuantityCriteria.getQualifier() == 1 ? oneSpell : multipleSpells));
	       } else if (Objects.equals(mType, KEY_CATEGORY)) {
	           builder.append(MessageFormat.format(I18n.text("\n{0}{1} {2} {3} whose category {4}"), prefix, getHasText(), mQuantityCriteria.toString(""), mQuantityCriteria.getQualifier() == 1 ? oneSpell : multipleSpells, mStringCriteria.toString()));
	       } else if (Objects.equals(mType, KEY_COLLEGE)) {
	           builder.append(MessageFormat.format(I18n.text("\n{0}{1} {2} {3} whose college {4}"), prefix, getHasText(), mQuantityCriteria.toString(""), mQuantityCriteria.getQualifier() == 1 ? oneSpell : multipleSpells, mStringCriteria.toString()));
	       } else if (Objects.equals(mType, KEY_COLLEGE_COUNT)) {
	           builder.append(MessageFormat.format(I18n.text("\n{0}{1} college count which {2}"), prefix, getHasText(), mQuantityCriteria.toString()));
	       }
	   }
	*/
	return satisfied
}
