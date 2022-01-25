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
)

const (
	bonusAttributeKey      = "attribute"
	bonusLimitationKey     = "limitation"
	bonusSituationKey      = "situation"
	bonusLocationKey       = "location"
	bonusSpecializationKey = "specialization"
	bonusSelectionTypeKey  = "selection_type"
	bonusNameKey           = "name"
	bonusCategoryKey       = "category"
)

// All is the DR specialization key for DR that affects everything.
const All = "all"

const (
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
	SelectionType          enum.SkillSelectionType       // Used by SkillBonusType
	Amount                 LeveledAmount                 // Used by all
	Attribute              string                        // Used by AttributeBonusType
	Situation              string                        // Used by ConditionalModifierBonusType, ReactionBonusType
	Location               string                        // Used by DRBonusType
	Specialization         string                        // Used by DRBonusType
	NameCriteria           StringCriteria                // Used by SkillBonusType
	SpecializationCriteria StringCriteria                // Used by SkillBonusType
	CategoryCriteria       StringCriteria                // Used by SkillBonusType
}

// NewBonus creates a new Bonus for the given entity, which may be nil.
func NewBonus(bonusType enum.BonusType, entity *Entity) *Bonus {
	b := &Bonus{
		Type:   bonusType,
		Amount: LeveledAmount{Amount: 1},
	}
	switch bonusType {
	case enum.ConditionalModifierBonusType:
		b.Situation = i18n.Text("triggering condition")
	case enum.DRBonusType:
		b.Location = "torso"
		b.Specialization = All
	case enum.ReactionBonusType:
		b.Situation = i18n.Text("from others")
	case enum.SkillBonusType:
		b.SelectionType = enum.SkillsWithName
		b.NameCriteria.Type = enum.Is
		b.SpecializationCriteria.Type = enum.Any
		b.CategoryCriteria.Type = enum.Any
	case enum.SkillPointBonusType:
	// TODO: Implement
	case enum.SpellBonusType:
	// TODO: Implement
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
	// TODO: Implement
	default: // AttributeBonusType
		b.Attribute = DefaultAttributeIDFor(entity)
	}
	return b
}

// NewBonusFromJSON creates a new Bonus from JSON.
func NewBonusFromJSON(key string, data map[string]interface{}) *Bonus {
	b := &Bonus{Type: enum.BonusTypeFromString(key)}
	b.Amount.FromJSON(data)
	switch b.Type {
	case enum.ConditionalModifierBonusType, enum.ReactionBonusType:
		b.Situation = encoding.String(data[bonusSituationKey])
	case enum.DRBonusType:
		b.Location = encoding.String(data[bonusLocationKey])
		b.Specialization = encoding.String(data[bonusSpecializationKey])
	case enum.SkillBonusType:
		b.SelectionType = enum.SkillSelectionTypeFromString(encoding.String(data[bonusSelectionTypeKey]))
		b.SpecializationCriteria.FromJSON(encoding.Object(data[bonusSpecializationKey]))
		if b.SelectionType != enum.ThisWeapon {
			b.NameCriteria.FromJSON(encoding.Object(data[bonusNameKey]))
			b.CategoryCriteria.FromJSON(encoding.Object(data[bonusCategoryKey]))
		}
	case enum.SkillPointBonusType:
	// TODO: Implement
	case enum.SpellBonusType:
	// TODO: Implement
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
	// TODO: Implement
	default: // AttributeBonusType
		b.Attribute = encoding.String(data[bonusAttributeKey])
		b.Limitation = enum.AttributeBonusLimitationFromString(encoding.String(data[bonusLimitationKey]))
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
	case enum.ConditionalModifierBonusType, enum.ReactionBonusType:
		encoder.KeyedString(bonusSituationKey, b.Situation, true, true)
	case enum.DRBonusType:
		encoder.KeyedString(bonusLocationKey, b.Location, true, true)
		encoder.KeyedString(bonusSpecializationKey, b.Specialization, true, true)
	case enum.SkillBonusType:
		encoder.KeyedString(bonusSelectionTypeKey, b.SelectionType.Key(), false, false)
		encoding.ToKeyedJSON(&b.SpecializationCriteria, bonusSpecializationKey, encoder)
		if b.SelectionType != enum.ThisWeapon {
			encoding.ToKeyedJSON(&b.NameCriteria, bonusNameKey, encoder)
			encoding.ToKeyedJSON(&b.CategoryCriteria, bonusCategoryKey, encoder)
		}
	case enum.SkillPointBonusType:
	// TODO: Implement
	case enum.SpellBonusType:
	// TODO: Implement
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
	// TODO: Implement
	default: // AttributeBonusType
		encoder.KeyedString(bonusAttributeKey, b.Attribute, false, false)
		if b.Limitation != enum.None {
			encoder.KeyedString(bonusLimitationKey, b.Limitation.Key(), false, false)
		}
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
	case enum.ConditionalModifierBonusType:
		return enum.ConditionalModifierBonusType.Key()
	case enum.DRBonusType:
		return HitLocationPrefix + b.Location
	case enum.ReactionBonusType:
		return "reaction"
	case enum.SkillBonusType:
		switch b.SelectionType {
		case enum.ThisWeapon:
			return ThisWeaponID
		case enum.WeaponsWithName:
			return b.buildKey(WeaponNamedIDPrefix)
		default: // SkillsWithName
			return b.buildKey(SkillIDName)
		}
	case enum.SkillPointBonusType:
	// TODO: Implement
	case enum.SpellBonusType:
	// TODO: Implement
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
	// TODO: Implement
	default: // AttributeBonusType
		key := AttributeIDPrefix + b.Attribute
		if b.Limitation != enum.None {
			key += "." + b.Limitation.Key()
		}
		return key
	}
	// TODO: Eliminate
	return ""
}

func (b *Bonus) buildKey(prefix string) string {
	var buffer strings.Builder
	buffer.WriteString(prefix)
	if b.NameCriteria.Type == enum.Is && b.SpecializationCriteria.Type == enum.Any && b.CategoryCriteria.Type == enum.Any {
		buffer.WriteByte('/')
		buffer.WriteString(b.NameCriteria.Qualifier)
	} else {
		buffer.WriteByte('*')
	}
	return buffer.String()
}

// FillWithNameableKeys implements Feature.
func (b *Bonus) FillWithNameableKeys(nameables map[string]string) {
	switch b.Type {
	case enum.ConditionalModifierBonusType, enum.ReactionBonusType:
		ExtractNameables(b.Situation, nameables)
	case enum.SkillBonusType:
		ExtractNameables(b.SpecializationCriteria.Qualifier, nameables)
		if b.SelectionType != enum.ThisWeapon {
			ExtractNameables(b.NameCriteria.Qualifier, nameables)
			ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	case enum.SkillPointBonusType:
	// TODO: Implement
	case enum.SpellBonusType:
	// TODO: Implement
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
	// TODO: Implement
	default: // AttributeBonusType, DRBonusType
		// Does nothing
	}
}

// ApplyNameableKeys implements Feature.
func (b *Bonus) ApplyNameableKeys(nameables map[string]string) {
	switch b.Type {
	case enum.ConditionalModifierBonusType, enum.ReactionBonusType:
		b.Situation = ApplyNameables(b.Situation, nameables)
	case enum.SkillBonusType:
		b.SpecializationCriteria.Qualifier = ApplyNameables(b.SpecializationCriteria.Qualifier, nameables)
		if b.SelectionType != enum.ThisWeapon {
			b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
			b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	case enum.SkillPointBonusType:
	// TODO: Implement
	case enum.SpellBonusType:
	// TODO: Implement
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
	// TODO: Implement
	default: // AttributeBonusType, DRBonusType
		// Does nothing
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
	case enum.SkillBonusType:
	// TODO: Implement
	case enum.SkillPointBonusType:
	// TODO: Implement
	case enum.SpellBonusType:
	// TODO: Implement
	case enum.SpellPointBonusType:
	// TODO: Implement
	case enum.WeaponDamageBonusType:
	// TODO: Implement
	default: // AttributeBonusType, ConditionalModifierBonusType, ReactionBonusType
		// Does nothing
	}
}
