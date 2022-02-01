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
	"github.com/richardwilkes/gcs/model/gurps/nameables"
	"github.com/richardwilkes/gcs/model/gurps/prereq"
	"github.com/richardwilkes/toolbox/xio"
)

var _ Prereq = &SkillPrereq{}

// SkillPrereq holds a prerequisite for a skill.
type SkillPrereq struct {
	Parent                 *PrereqList      `json:"-"`
	Type                   prereq.Type      `json:"type"`
	NameCriteria           criteria.String  `json:"name"`
	LevelCriteria          criteria.Numeric `json:"level"`
	SpecializationCriteria criteria.String  `json:"specialization"`
	Has                    bool             `json:"has,omitempty"`
}

// NewSkillPrereq creates a new SkillPrereq.
func NewSkillPrereq() *SkillPrereq {
	return &SkillPrereq{
		Type: prereq.Skill,
		NameCriteria: criteria.String{
			Compare: criteria.Is,
		},
		LevelCriteria: criteria.Numeric{
			Compare: criteria.AtLeast,
		},
		SpecializationCriteria: criteria.String{
			Compare: criteria.Any,
		},
		Has: true,
	}
}

// Clone implements Prereq.
func (s *SkillPrereq) Clone(parent *PrereqList) Prereq {
	clone := *s
	clone.Parent = parent
	return &clone
}

// FillWithNameableKeys implements Prereq.
func (s *SkillPrereq) FillWithNameableKeys(m map[string]string) {
	nameables.Extract(s.NameCriteria.Qualifier, m)
	nameables.Extract(s.SpecializationCriteria.Qualifier, m)
}

// ApplyNameableKeys implements Prereq.
func (s *SkillPrereq) ApplyNameableKeys(m map[string]string) {
	s.NameCriteria.Qualifier = nameables.Apply(s.NameCriteria.Qualifier, m)
	s.SpecializationCriteria.Qualifier = nameables.Apply(s.SpecializationCriteria.Qualifier, m)
}

// Satisfied implements Prereq.
func (s *SkillPrereq) Satisfied(entity *Entity, exclude interface{}, buffer *xio.ByteBuffer, prefix string) bool {
	satisfied := false
	// TODO: Implement
	/*
	   String          techLevel     = null;
	   StringCriteria  nameCriteria  = getNameCriteria();
	   IntegerCriteria levelCriteria = getLevelCriteria();

	   if (exclude instanceof Skill) {
	       techLevel = ((Skill) exclude).getTechLevel();
	   }

	   for (Skill skill : character.getSkillsIterator()) {
	       if (exclude != skill && nameCriteria.matches(skill.getName()) && mSpecializationCriteria.matches(skill.getSpecialization())) {
	           satisfied = levelCriteria.matches(skill.getLevel());
	           if (satisfied && techLevel != null) {
	               String otherTL = skill.getTechLevel();
	               satisfied = otherTL == null || techLevel.equals(otherTL);
	           }
	           if (satisfied) {
	               break;
	           }
	       }
	   }
	   if (!has()) {
	       satisfied = !satisfied;
	   }
	   if (!satisfied && builder != null) {
	       builder.append(MessageFormat.format(I18n.text("\n{0}{1} a skill whose name {2}"), prefix, getHasText(), nameCriteria.toString()));
	       boolean notAnySpecialization = !mSpecializationCriteria.isTypeAnything();
	       if (notAnySpecialization) {
	           builder.append(MessageFormat.format(I18n.text(", specialization {0},"), mSpecializationCriteria.toString()));
	       }
	       if (techLevel == null) {
	           builder.append(MessageFormat.format(I18n.text(" and level {0}"), levelCriteria.toString()));
	       } else {
	           if (notAnySpecialization) {
	               builder.append(",");
	           }
	           builder.append(MessageFormat.format(I18n.text(" level {0} and tech level matches"), levelCriteria.toString()));
	       }
	   }
	*/
	return satisfied
}
