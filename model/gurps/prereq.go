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
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/enum"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Already applied the HasPrereq and NameLevelPrereq, AdvantagePrereq, AttributePrereq, ContainedQuantityPrereq,
// ContainedWeightPrereq, SkillPrereq portions

const (
	prereqTypeKey           = "type"
	prereqHasKey            = "has"
	prereqNameKey           = "name"
	prereqLevelKey          = "level"
	prereqNotesKey          = "notes"
	prereqWhichKey          = "which"
	prereqCombinedWithKey   = "combined_with"
	prereqQualifierKey      = "qualifier"
	prereqSpecializationKey = "specialization"
)

// Prereq holds data necessary to track a prerequisite.
type Prereq struct {
	Type                   enum.PrereqType
	Has                    bool
	NameCriteria           criteria.String
	LevelCriteria          criteria.Numeric
	NotesCriteria          criteria.String
	Which                  string
	CombinedWith           string
	ValueCriteria          criteria.Numeric
	WeightCriteria         criteria.Weight
	SpecializationCriteria criteria.String
	Owner                  *Prereq // Only those of type PrereqList
}

// NewPrereq creates a new Prereq for the given entity, which may be nil.
func NewPrereq(prereqType enum.PrereqType, entity *Entity) *Prereq {
	p := &Prereq{
		Type: prereqType,
	}
	switch prereqType {
	case enum.AdvantagePrereq:
		p.Has = true
		p.NameCriteria.Type = enum.Is
		p.LevelCriteria.Type = enum.AtLeast
		p.NotesCriteria.Type = enum.Any
	case enum.AttributePrereq:
		p.Has = true
		p.ValueCriteria.Type = enum.AtLeast
		p.ValueCriteria.Qualifier = fixed.F64d4FromInt64(10)
		p.Which = DefaultAttributeIDFor(entity)
	case enum.ContainedQuantityPrereq:
		p.Has = true
		p.ValueCriteria.Type = enum.AtMost
		p.ValueCriteria.Qualifier = f64d4.One
	case enum.ContainedWeightPrereq:
		p.WeightCriteria.Type = enum.AtMost
		p.WeightCriteria.Qualifier = measure.WeightFromInt64(5, SheetSettingsFor(entity).DefaultWeightUnits)
		p.Has = true
	case enum.PrereqList:
	// TODO: Implement
	case enum.SkillPrereq:
		p.Has = true
		p.NameCriteria.Type = enum.Is
		p.LevelCriteria.Type = enum.AtLeast
		p.SpecializationCriteria.Type = enum.Any
	case enum.SpellPrereq:
		// TODO: Implement
		p.Has = true
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
	return p
}

// NewPrereqFromJSON creates a new Prereq from JSON.
func NewPrereqFromJSON(data map[string]interface{}, entity *Entity) *Prereq {
	p := &Prereq{Type: enum.PrereqTypeFromString(encoding.String(data[prereqTypeKey]))}
	switch p.Type {
	case enum.AdvantagePrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.NameCriteria.FromJSON(encoding.Object(data[prereqNameKey]))
		p.LevelCriteria.FromJSON(encoding.Object(data[prereqLevelKey]))
		p.NotesCriteria.FromJSON(encoding.Object(data[prereqNotesKey]))
	case enum.AttributePrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.Which = encoding.String(data[prereqWhichKey])
		p.CombinedWith = encoding.String(data[prereqCombinedWithKey])
		p.ValueCriteria.FromJSON(encoding.Object(data[prereqQualifierKey]))
	case enum.ContainedQuantityPrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.ValueCriteria.FromJSON(encoding.Object(data[prereqQualifierKey]))
	case enum.ContainedWeightPrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.WeightCriteria.FromJSON(encoding.Object(data[prereqQualifierKey]), SheetSettingsFor(entity).DefaultWeightUnits)
	case enum.PrereqList:
	// TODO: Implement
	case enum.SkillPrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.NameCriteria.FromJSON(encoding.Object(data[prereqNameKey]))
		p.LevelCriteria.FromJSON(encoding.Object(data[prereqLevelKey]))
		p.SpecializationCriteria.FromJSON(encoding.Object(data[prereqSpecializationKey]))
	case enum.SpellPrereq:
		// TODO: Implement
		p.Has = encoding.Bool(data[prereqHasKey])
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
	return p
}

// ToJSON emits this Feature as JSON.
func (p *Prereq) ToJSON(encoder *encoding.JSONEncoder) {
	encoder.StartObject()
	encoder.KeyedString(prereqTypeKey, p.Type.Key(), false, false)
	switch p.Type {
	case enum.AdvantagePrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoding.ToKeyedJSON(&p.NameCriteria, prereqNameKey, encoder)
		if p.LevelCriteria.Type != enum.AtLeast || p.LevelCriteria.Qualifier != 0 {
			encoding.ToKeyedJSON(&p.LevelCriteria, prereqLevelKey, encoder)
		}
		encoding.ToKeyedJSON(&p.NotesCriteria, prereqNotesKey, encoder)
	case enum.AttributePrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoder.KeyedString(prereqWhichKey, p.Which, true, true)
		encoder.KeyedString(prereqCombinedWithKey, p.CombinedWith, true, true)
		encoding.ToKeyedJSON(&p.ValueCriteria, prereqQualifierKey, encoder)
	case enum.ContainedQuantityPrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoding.ToKeyedJSON(&p.ValueCriteria, prereqQualifierKey, encoder)
	case enum.ContainedWeightPrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoding.ToKeyedJSON(&p.WeightCriteria, prereqQualifierKey, encoder)
	case enum.PrereqList:
	// TODO: Implement
	case enum.SkillPrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoding.ToKeyedJSON(&p.NameCriteria, prereqNameKey, encoder)
		if p.LevelCriteria.Type != enum.AtLeast || p.LevelCriteria.Qualifier != 0 {
			encoding.ToKeyedJSON(&p.LevelCriteria, prereqLevelKey, encoder)
		}
		encoding.ToKeyedJSON(&p.SpecializationCriteria, prereqSpecializationKey, encoder)
	case enum.SpellPrereq:
		// TODO: Implement
		encoder.KeyedBool(prereqHasKey, p.Has, false)
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
	encoder.EndObject()
}

// Satisfied returns true if this Prereq is satisfied by the specified Entity. 'buffer' will be used, if not nil, to
// write a description of what was unsatisfied. 'prefix' will be appended to each line of the description.
func (p *Prereq) Satisfied(entity *Entity, exclude interface{}, buffer *xio.ByteBuffer, prefix string) bool {
	switch p.Type {
	case enum.AdvantagePrereq:
	// TODO: Implement
	/*
	   boolean         satisfied     = false;
	   StringCriteria  nameCriteria  = getNameCriteria();
	   IntegerCriteria levelCriteria = getLevelCriteria();

	   for (Advantage advantage : character.getAdvantagesIterator(false)) {
	       if (exclude != advantage && nameCriteria.matches(advantage.getName())) {
	           String notes         = advantage.getNotes();
	           String modifierNotes = advantage.getModifierNotes();

	           if (!modifierNotes.isEmpty()) {
	               notes = modifierNotes + '\n' + notes;
	           }
	           if (mNotesCriteria.matches(notes)) {
	               int levels = advantage.getLevels();
	               if (levels < 0) {
	                   levels = 0;
	               }
	               satisfied = levelCriteria.matches(levels);
	               break;
	           }
	       }
	   }
	   if (!has()) {
	       satisfied = !satisfied;
	   }
	   if (!satisfied && builder != null) {
	       builder.append(MessageFormat.format(I18n.text("\n{0}{1} an advantage whose name {2}"), prefix, getHasText(), nameCriteria.toString()));
	       if (!mNotesCriteria.isTypeAnything()) {
	           builder.append(MessageFormat.format(I18n.text(", notes {0},"), mNotesCriteria.toString()));
	       }
	       builder.append(MessageFormat.format(I18n.text(" and level {0}"), levelCriteria.toString()));
	   }
	   return satisfied;
	*/
	case enum.AttributePrereq:
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
	   return satisfied;
	*/
	case enum.ContainedQuantityPrereq:
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
	case enum.ContainedWeightPrereq:
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
	case enum.PrereqList:
	// TODO: Implement
	case enum.SkillPrereq:
	// TODO: Implement
	/*
	   boolean         satisfied     = false;
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
	   return satisfied;
	*/
	case enum.SpellPrereq:
	// TODO: Implement
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
	return false
}

// FillWithNameableKeys adds any nameable keys found in this Prereq to the provided map.
func (p *Prereq) FillWithNameableKeys(nameables map[string]string) {
	switch p.Type {
	case enum.AdvantagePrereq:
		ExtractNameables(p.NameCriteria.Qualifier, nameables)
		ExtractNameables(p.NotesCriteria.Qualifier, nameables)
	case enum.AttributePrereq:
	case enum.ContainedQuantityPrereq:
	case enum.ContainedWeightPrereq:
	case enum.PrereqList:
	// TODO: Implement
	case enum.SkillPrereq:
		ExtractNameables(p.NameCriteria.Qualifier, nameables)
		ExtractNameables(p.SpecializationCriteria.Qualifier, nameables)
	case enum.SpellPrereq:
	// TODO: Implement
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
}

// ApplyNameableKeys replaces any nameable keys found in this Prereq with the corresponding values in the provided map.
func (p *Prereq) ApplyNameableKeys(nameables map[string]string) {
	switch p.Type {
	case enum.AdvantagePrereq:
		p.NameCriteria.Qualifier = ApplyNameables(p.NameCriteria.Qualifier, nameables)
		p.NotesCriteria.Qualifier = ApplyNameables(p.NotesCriteria.Qualifier, nameables)
	case enum.AttributePrereq:
	case enum.ContainedQuantityPrereq:
	case enum.ContainedWeightPrereq:
	case enum.PrereqList:
	// TODO: Implement
	case enum.SkillPrereq:
		p.NameCriteria.Qualifier = ApplyNameables(p.NameCriteria.Qualifier, nameables)
		p.SpecializationCriteria.Qualifier = ApplyNameables(p.SpecializationCriteria.Qualifier, nameables)
	case enum.SpellPrereq:
	// TODO: Implement
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
}
