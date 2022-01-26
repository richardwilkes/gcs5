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
	prereqAllKey            = "all"
	prereqChildrenKey       = "prereqs"
	prereqCombinedWithKey   = "combined_with"
	prereqHasKey            = "has"
	prereqLevelKey          = "level"
	prereqNameKey           = "name"
	prereqNotesKey          = "notes"
	prereqQualifierKey      = "qualifier"
	prereqQuantityKey       = "quantity"
	prereqSpecializationKey = "specialization"
	prereqSubTypeKey        = "sub_type"
	prereqTypeKey           = "type"
	prereqWhenTLKey         = "when_tl"
	prereqWhichKey          = "which"
)

// Prereq holds data necessary to track a prerequisite.
type Prereq struct {
	Type                   enum.PrereqType
	SubType                enum.SpellComparisonType
	Has                    bool
	WhenEnabled            bool
	All                    bool
	NameCriteria           criteria.String
	SpecializationCriteria criteria.String
	NotesCriteria          criteria.String
	NumericCriteria        criteria.Numeric
	WeightCriteria         criteria.Weight
	Which                  string
	CombinedWith           string
	Children               []*Prereq
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
		p.NumericCriteria.Type = enum.AtLeast
		p.NotesCriteria.Type = enum.Any
	case enum.AttributePrereq:
		p.Has = true
		p.NumericCriteria.Type = enum.AtLeast
		p.NumericCriteria.Qualifier = fixed.F64d4FromInt64(10)
		p.Which = DefaultAttributeIDFor(entity)
	case enum.ContainedQuantityPrereq:
		p.Has = true
		p.NumericCriteria.Type = enum.AtMost
		p.NumericCriteria.Qualifier = f64d4.One
	case enum.ContainedWeightPrereq:
		p.Has = true
		p.WeightCriteria.Type = enum.AtMost
		p.WeightCriteria.Qualifier = measure.WeightFromInt64(5, SheetSettingsFor(entity).DefaultWeightUnits)
	case enum.PrereqList:
		p.NumericCriteria.Type = enum.AtLeast
	case enum.SkillPrereq:
		p.Has = true
		p.NameCriteria.Type = enum.Is
		p.NumericCriteria.Type = enum.AtLeast
		p.SpecializationCriteria.Type = enum.Any
	case enum.SpellPrereq:
		p.Has = true
		p.SubType = enum.Name
		p.NameCriteria.Type = enum.Is
		p.NumericCriteria.Type = enum.AtLeast
		p.NumericCriteria.Qualifier = f64d4.One
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
		p.NumericCriteria.FromJSON(encoding.Object(data[prereqLevelKey]))
		p.NotesCriteria.FromJSON(encoding.Object(data[prereqNotesKey]))
	case enum.AttributePrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.Which = encoding.String(data[prereqWhichKey])
		p.CombinedWith = encoding.String(data[prereqCombinedWithKey])
		p.NumericCriteria.FromJSON(encoding.Object(data[prereqQualifierKey]))
	case enum.ContainedQuantityPrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.NumericCriteria.FromJSON(encoding.Object(data[prereqQualifierKey]))
	case enum.ContainedWeightPrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.WeightCriteria.FromJSON(encoding.Object(data[prereqQualifierKey]), SheetSettingsFor(entity).DefaultWeightUnits)
	case enum.PrereqList:
		p.All = encoding.Bool(data[prereqAllKey])
		if _, p.WhenEnabled = data[prereqWhenTLKey]; p.WhenEnabled {
			p.NumericCriteria.FromJSON(encoding.Object(data[prereqWhenTLKey]))
		}
		if array := encoding.Array(data[prereqChildrenKey]); len(array) != 0 {
			p.Children = make([]*Prereq, 0, len(array))
			for _, one := range array {
				p.Children = append(p.Children, NewPrereqFromJSON(encoding.Object(one), entity))
			}
		}
	case enum.SkillPrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.NameCriteria.FromJSON(encoding.Object(data[prereqNameKey]))
		p.NumericCriteria.FromJSON(encoding.Object(data[prereqLevelKey]))
		p.SpecializationCriteria.FromJSON(encoding.Object(data[prereqSpecializationKey]))
	case enum.SpellPrereq:
		p.Has = encoding.Bool(data[prereqHasKey])
		p.SubType = enum.SpellComparisonTypeFromString(encoding.String(data[prereqSubTypeKey]))
		p.NameCriteria.FromJSON(encoding.Object(data[prereqQualifierKey]))
		p.NumericCriteria.FromJSON(encoding.Object(data[prereqQuantityKey]))
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
		if p.NumericCriteria.Type != enum.AtLeast || p.NumericCriteria.Qualifier != 0 {
			encoding.ToKeyedJSON(&p.NumericCriteria, prereqLevelKey, encoder)
		}
		encoding.ToKeyedJSON(&p.NotesCriteria, prereqNotesKey, encoder)
	case enum.AttributePrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoder.KeyedString(prereqWhichKey, p.Which, true, true)
		encoder.KeyedString(prereqCombinedWithKey, p.CombinedWith, true, true)
		encoding.ToKeyedJSON(&p.NumericCriteria, prereqQualifierKey, encoder)
	case enum.ContainedQuantityPrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoding.ToKeyedJSON(&p.NumericCriteria, prereqQualifierKey, encoder)
	case enum.ContainedWeightPrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoding.ToKeyedJSON(&p.WeightCriteria, prereqQualifierKey, encoder)
	case enum.PrereqList:
		encoder.KeyedBool(prereqAllKey, p.All, false)
		if p.WhenEnabled {
			encoding.ToKeyedJSON(&p.NumericCriteria, prereqWhenTLKey, encoder)
		}
		if len(p.Children) != 0 {
			encoder.Key(prereqChildrenKey)
			encoder.StartArray()
			for _, one := range p.Children {
				one.ToJSON(encoder)
			}
			encoder.EndArray()
		}
	case enum.SkillPrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoding.ToKeyedJSON(&p.NameCriteria, prereqNameKey, encoder)
		if p.NumericCriteria.Type != enum.AtLeast || p.NumericCriteria.Qualifier != 0 {
			encoding.ToKeyedJSON(&p.NumericCriteria, prereqLevelKey, encoder)
		}
		encoding.ToKeyedJSON(&p.SpecializationCriteria, prereqSpecializationKey, encoder)
	case enum.SpellPrereq:
		encoder.KeyedBool(prereqHasKey, p.Has, false)
		encoder.KeyedString(prereqSubTypeKey, p.SubType.Key(), false, false)
		if p.SubType.UsesStringCriteria() {
			encoding.ToKeyedJSON(&p.NameCriteria, prereqQualifierKey, encoder)
		}
		encoding.ToKeyedJSON(&p.NumericCriteria, prereqQuantityKey, encoder)
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
	encoder.EndObject()
}

// Clone creates a new copy of this Prereq.
func (p *Prereq) Clone(owner *Prereq) *Prereq {
	clone := *p
	clone.Owner = owner
	if p.Type == enum.PrereqList {
		clone.Children = make([]*Prereq, 0, len(p.Children))
		for _, one := range p.Children {
			clone.Children = append(clone.Children, one.Clone(&clone))
		}
	}
	return &clone
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
	/*
	   if (mWhenEnabled) {
	       if (!mWhenTLCriteria.matches(Numbers.extractInteger(character.getProfile().getTechLevel(), 0, false))) {
	           return true;
	       }
	   }

	   int           satisfiedCount = 0;
	   int           total          = mPrereqs.size();
	   boolean       requiresAll    = requiresAll();
	   StringBuilder localBuilder   = builder != null ? new StringBuilder() : null;
	   for (Prereq prereq : mPrereqs) {
	       if (prereq.satisfied(character, exclude, localBuilder, prefix)) {
	           satisfiedCount++;
	       }
	   }
	   if (localBuilder != null && !localBuilder.isEmpty()) {
	       String indented = LINE_FEED_MATCHER.matcher(localBuilder.toString()).replaceAll("\n\u00a0\u00a0");
	       localBuilder.setLength(0);
	       localBuilder.append(indented);
	   }

	   boolean satisfied = satisfiedCount == total || !requiresAll && satisfiedCount > 0;
	   if (!satisfied && localBuilder != null) {
	       builder.append("\n");
	       builder.append(prefix);
	       builder.append(requiresAll ? I18n.text("Requires all of:") : I18n.text("Requires at least one of:"));
	       builder.append(localBuilder);
	   }
	   return satisfied;
	*/
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
	   return satisfied;
	*/
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
		for _, one := range p.Children {
			one.FillWithNameableKeys(nameables)
		}
	case enum.SkillPrereq:
		ExtractNameables(p.NameCriteria.Qualifier, nameables)
		ExtractNameables(p.SpecializationCriteria.Qualifier, nameables)
	case enum.SpellPrereq:
		if p.SubType.UsesStringCriteria() {
			ExtractNameables(p.NameCriteria.Qualifier, nameables)
		}
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
		for _, one := range p.Children {
			one.ApplyNameableKeys(nameables)
		}
	case enum.SkillPrereq:
		p.NameCriteria.Qualifier = ApplyNameables(p.NameCriteria.Qualifier, nameables)
		p.SpecializationCriteria.Qualifier = ApplyNameables(p.SpecializationCriteria.Qualifier, nameables)
	case enum.SpellPrereq:
		if p.SubType.UsesStringCriteria() {
			p.NameCriteria.Qualifier = ApplyNameables(p.NameCriteria.Qualifier, nameables)
		}
	default:
		jot.Fatal(1, "invalid prereq type: ", p.Type)
	}
}
