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
	"strings"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/gurps/enum"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
)

const (
	bonusAttributeKey      = "attribute"
	bonusCategoryKey       = "category"
	bonusLimitationKey     = "limitation"
	bonusLocationKey       = "location"
	bonusMatchKey          = "match"
	bonusNameKey           = "name"
	bonusSelectionTypeKey  = "selection_type"
	bonusSituationKey      = "situation"
	bonusSpecializationKey = "specialization"
)

const (
	// All is the DR specialization key for DR that affects everything.
	All = "all"
	// ThisWeaponID holds the ID for "this weapon".
	ThisWeaponID = "\u0001"
	// WeaponNamedIDPrefix the prefix for "weapon named" IDs.
	WeaponNamedIDPrefix = "weapon_named."
)

var _ Feature = &Bonus{}

// Bonus holds a numerical bonus to another object.
type Bonus struct {
	Type                   enum.BonusType                // Used by all
	Limitation             enum.AttributeBonusLimitation // Used by AttributeBonusType
	SkillSelectionType     enum.SkillSelectionType       // Used by SkillBonusType
	SpellMatchType         enum.SpellMatchType           // Used by SpellBonusType
	Amount                 LeveledAmount                 // Used by all
	Attribute              string                        // Used by AttributeBonusType
	Situation              string                        // Used by ConditionalModifierBonusType, ReactionBonusType
	Location               string                        // Used by DRBonusType
	Specialization         string                        // Used by DRBonusType
	NameCriteria           StringCriteria                // Used by SkillBonusType, SkillPointBonusType, SpellBonusType
	SpecializationCriteria StringCriteria                // Used by SkillBonusType, SkillPointBonusType
	CategoryCriteria       StringCriteria                // Used by SkillBonusType, SkillPointBonusType, SpellBonusType
}

// NewBonus creates a new Bonus for the given entity, which may be nil.
func NewBonus(bonusType enum.BonusType, entity *Entity) *Bonus {
	b := &Bonus{
		Type:   bonusType,
		Amount: LeveledAmount{Amount: 1},
	}
	switch bonusType {
	case enum.AttributeBonusType:
		b.Attribute = DefaultAttributeIDFor(entity)
	case enum.ConditionalModifierBonusType:
		b.Situation = i18n.Text("triggering condition")
	case enum.DRBonusType:
		b.Location = "torso"
		b.Specialization = All
	case enum.ReactionBonusType:
		b.Situation = i18n.Text("from others")
	case enum.SkillBonusType:
		b.SkillSelectionType = enum.SkillsWithName
		fallthrough
	case enum.SkillPointBonusType:
		b.NameCriteria.Type = enum.Is
		b.SpecializationCriteria.Type = enum.Any
		b.CategoryCriteria.Type = enum.Any
	case enum.SpellBonusType:
		b.SpellMatchType = enum.AllColleges
		b.NameCriteria.Type = enum.Is
		b.CategoryCriteria.Type = enum.Any
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
	// TODO: Implement
	default:
		jot.Fatal(1, "invalid bonus type: ", b.Type)
	}
	return b
}

// NewBonusFromJSON creates a new Bonus from JSON.
func NewBonusFromJSON(key string, data map[string]interface{}) *Bonus {
	b := &Bonus{Type: enum.BonusTypeFromString(key)}
	b.Amount.FromJSON(data)
	switch b.Type {
	case enum.AttributeBonusType:
		b.Attribute = encoding.String(data[bonusAttributeKey])
		b.Limitation = enum.AttributeBonusLimitationFromString(encoding.String(data[bonusLimitationKey]))
	case enum.ConditionalModifierBonusType, enum.ReactionBonusType:
		b.Situation = encoding.String(data[bonusSituationKey])
	case enum.DRBonusType:
		b.Location = encoding.String(data[bonusLocationKey])
		b.Specialization = encoding.String(data[bonusSpecializationKey])
	case enum.SkillBonusType:
		b.SkillSelectionType = enum.SkillSelectionTypeFromString(encoding.String(data[bonusSelectionTypeKey]))
		b.SpecializationCriteria.FromJSON(encoding.Object(data[bonusSpecializationKey]))
		if b.SkillSelectionType != enum.ThisWeapon {
			b.NameCriteria.FromJSON(encoding.Object(data[bonusNameKey]))
			b.CategoryCriteria.FromJSON(encoding.Object(data[bonusCategoryKey]))
		}
	case enum.SkillPointBonusType:
		b.NameCriteria.FromJSON(encoding.Object(data[bonusNameKey]))
		b.SpecializationCriteria.FromJSON(encoding.Object(data[bonusSpecializationKey]))
		b.CategoryCriteria.FromJSON(encoding.Object(data[bonusCategoryKey]))
	case enum.SpellBonusType:
		b.SpellMatchType = enum.SpellMatchTypeFromString(encoding.String(data[bonusMatchKey]))
		b.NameCriteria.FromJSON(encoding.Object(data[bonusNameKey]))
		b.CategoryCriteria.FromJSON(encoding.Object(data[bonusCategoryKey]))
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
	// TODO: Implement
	default:
		jot.Fatal(1, "invalid bonus type: ", b.Type)
	}
	b.Normalize()
	return b
}

// ToJSON implements Feature.
func (b *Bonus) ToJSON(encoder *encoding.JSONEncoder) {
	b.Normalize()
	encoder.StartObject()
	b.Amount.ToInlineJSON(encoder)
	switch b.Type {
	case enum.AttributeBonusType:
		encoder.KeyedString(bonusAttributeKey, b.Attribute, false, false)
		if b.Limitation != enum.None {
			encoder.KeyedString(bonusLimitationKey, b.Limitation.Key(), false, false)
		}
	case enum.ConditionalModifierBonusType, enum.ReactionBonusType:
		encoder.KeyedString(bonusSituationKey, b.Situation, true, true)
	case enum.DRBonusType:
		encoder.KeyedString(bonusLocationKey, b.Location, true, true)
		encoder.KeyedString(bonusSpecializationKey, b.Specialization, true, true)
	case enum.SkillBonusType:
		encoder.KeyedString(bonusSelectionTypeKey, b.SkillSelectionType.Key(), false, false)
		encoding.ToKeyedJSON(&b.SpecializationCriteria, bonusSpecializationKey, encoder)
		if b.SkillSelectionType != enum.ThisWeapon {
			encoding.ToKeyedJSON(&b.NameCriteria, bonusNameKey, encoder)
			encoding.ToKeyedJSON(&b.CategoryCriteria, bonusCategoryKey, encoder)
		}
	case enum.SkillPointBonusType:
		encoding.ToKeyedJSON(&b.NameCriteria, bonusNameKey, encoder)
		encoding.ToKeyedJSON(&b.SpecializationCriteria, bonusSpecializationKey, encoder)
		encoding.ToKeyedJSON(&b.CategoryCriteria, bonusCategoryKey, encoder)
	case enum.SpellBonusType:
		encoder.KeyedString(bonusMatchKey, b.SpellMatchType.Key(), false, false)
		encoding.ToKeyedJSON(&b.NameCriteria, bonusNameKey, encoder)
		encoding.ToKeyedJSON(&b.CategoryCriteria, bonusCategoryKey, encoder)
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
	// TODO: Implement
	default:
		jot.Fatal(1, "invalid bonus type: ", b.Type)
	}
	encoder.EndObject()
}

// CloneFeature implements Feature.
func (b *Bonus) CloneFeature() Feature {
	clone := *b
	return &clone
}

// DataType implements Feature.
func (b *Bonus) DataType() string {
	return b.Type.Key()
}

// FeatureKey implements Feature.
func (b *Bonus) FeatureKey() string {
	switch b.Type {
	case enum.AttributeBonusType:
		key := AttributeIDPrefix + b.Attribute
		if b.Limitation != enum.None {
			key += "." + b.Limitation.Key()
		}
		return key
	case enum.ConditionalModifierBonusType:
		return enum.ConditionalModifierBonusType.Key()
	case enum.DRBonusType:
		return HitLocationPrefix + b.Location
	case enum.ReactionBonusType:
		return "reaction"
	case enum.SkillBonusType:
		switch b.SkillSelectionType {
		case enum.SkillsWithName:
			return b.buildKey(SkillNameID, false)
		case enum.ThisWeapon:
			return ThisWeaponID
		case enum.WeaponsWithName:
			return b.buildKey(WeaponNamedIDPrefix, false)
		default:
			jot.Fatal(1, "invalid selection type: ", b.SkillSelectionType)
		}
	case enum.SkillPointBonusType:
		return b.buildKey(SkillPointsID, false)
	case enum.SpellBonusType:
		if b.CategoryCriteria.Type != enum.Any {
			return SpellNameID + "*"
		}
		switch b.SpellMatchType {
		case enum.AllColleges:
			return SpellCollegeID
		case enum.CollegeName:
			return b.buildKey(SpellCollegeID, true)
		case enum.PowerSourceName:
			return b.buildKey(SpellPowerSourceID, true)
		case enum.SpellName:
			return b.buildKey(SpellNameID, true)
		default:
			jot.Fatal(1, "invalid match type: ", b.SpellMatchType)
		}
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
	// TODO: Implement
	default:
		jot.Fatal(1, "invalid bonus type: ", b.Type)
	}
	// TODO: Eliminate
	return ""
}

func (b *Bonus) buildKey(prefix string, considerNameCriteriaOnly bool) string {
	if b.NameCriteria.Type == enum.Is && (considerNameCriteriaOnly ||
		(b.SpecializationCriteria.Type == enum.Any && b.CategoryCriteria.Type == enum.Any)) {
		return prefix + "/" + b.NameCriteria.Qualifier
	}
	return prefix + "*"
}

// FillWithNameableKeys implements Feature.
func (b *Bonus) FillWithNameableKeys(nameables map[string]string) {
	switch b.Type {
	case enum.ConditionalModifierBonusType, enum.ReactionBonusType:
		ExtractNameables(b.Situation, nameables)
	case enum.SkillBonusType:
		ExtractNameables(b.SpecializationCriteria.Qualifier, nameables)
		if b.SkillSelectionType != enum.ThisWeapon {
			ExtractNameables(b.NameCriteria.Qualifier, nameables)
			ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	case enum.SkillPointBonusType:
		ExtractNameables(b.NameCriteria.Qualifier, nameables)
		ExtractNameables(b.SpecializationCriteria.Qualifier, nameables)
		ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.SpellBonusType:
		if b.SpellMatchType != enum.AllColleges {
			ExtractNameables(b.NameCriteria.Qualifier, nameables)
		}
		ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
		// TODO: Implement
	}
}

// ApplyNameableKeys implements Feature.
func (b *Bonus) ApplyNameableKeys(nameables map[string]string) {
	switch b.Type {
	case enum.ConditionalModifierBonusType, enum.ReactionBonusType:
		b.Situation = ApplyNameables(b.Situation, nameables)
	case enum.SkillBonusType:
		b.SpecializationCriteria.Qualifier = ApplyNameables(b.SpecializationCriteria.Qualifier, nameables)
		if b.SkillSelectionType != enum.ThisWeapon {
			b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
			b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	case enum.SkillPointBonusType:
		b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
		b.SpecializationCriteria.Qualifier = ApplyNameables(b.SpecializationCriteria.Qualifier, nameables)
		b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.SpellBonusType:
		if b.SpellMatchType != enum.AllColleges {
			b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
		}
		b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
		// TODO: Implement
	}
}

// Normalize implements Feature.
func (b *Bonus) Normalize() {
	switch b.Type {
	case enum.DRBonusType:
		s := strings.TrimSpace(b.Specialization)
		if s == "" || strings.EqualFold(s, All) {
			s = All
		}
		b.Specialization = s
	case enum.SpellBonusType:
	// TODO: Implement
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
		// TODO: Implement
	}
}
