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
	"fmt"
	"strings"

	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/enum"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
)

const (
	bonusAttributeKey      = "attribute"
	bonusCategoryKey       = "category"
	bonusIsPercentKey      = "percent"
	bonusLevelKey          = "level"
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
	SpellMatchType         enum.SpellMatchType           // Used by SpellBonusType, SpellPointBonusType
	WeaponSelectionType    enum.WeaponSelectionType      // Used by WeaponDamageBonusType
	IsPercent              bool                          // Used by WeaponDamageBonusType
	Amount                 LeveledAmount                 // Used by all
	Attribute              string                        // Used by AttributeBonusType
	Situation              string                        // Used by ConditionalModifierBonusType, ReactionBonusType
	Location               string                        // Used by DRBonusType
	Specialization         string                        // Used by DRBonusType
	NameCriteria           StringCriteria                // Used by SkillBonusType, SkillPointBonusType, SpellBonusType, SpellPointBonusType, WeaponDamageBonusType
	SpecializationCriteria StringCriteria                // Used by SkillBonusType, SkillPointBonusType, WeaponDamageBonusType
	CategoryCriteria       StringCriteria                // Used by SkillBonusType, SkillPointBonusType, SpellBonusType, SpellPointBonusType, WeaponDamageBonusType
	RelativeLevelCriteria  NumberCriteria                // Used by WeaponDamageBonusType
	Owner                  fmt.Stringer
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
		b.SkillSelectionType = enum.SkillsWithNameSkillSelect
		fallthrough
	case enum.SkillPointBonusType:
		b.NameCriteria.Type = enum.Is
		b.SpecializationCriteria.Type = enum.Any
		b.CategoryCriteria.Type = enum.Any
	case enum.SpellBonusType, enum.SpellPointBonusType:
		b.SpellMatchType = enum.AllColleges
		b.NameCriteria.Type = enum.Is
		b.CategoryCriteria.Type = enum.Any
	case enum.WeaponDamageBonusType:
		b.WeaponSelectionType = enum.WeaponsWithRequiredSkillWeaponSelect
		b.NameCriteria.Type = enum.Is
		b.SpecializationCriteria.Type = enum.Any
		b.RelativeLevelCriteria.Type = enum.AnyNumber
		b.CategoryCriteria.Type = enum.Any
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
		if b.SkillSelectionType != enum.ThisWeaponSkillSelect {
			b.NameCriteria.FromJSON(encoding.Object(data[bonusNameKey]))
			b.CategoryCriteria.FromJSON(encoding.Object(data[bonusCategoryKey]))
		}
	case enum.SkillPointBonusType:
		b.NameCriteria.FromJSON(encoding.Object(data[bonusNameKey]))
		b.SpecializationCriteria.FromJSON(encoding.Object(data[bonusSpecializationKey]))
		b.CategoryCriteria.FromJSON(encoding.Object(data[bonusCategoryKey]))
	case enum.SpellBonusType, enum.SpellPointBonusType:
		b.SpellMatchType = enum.SpellMatchTypeFromString(encoding.String(data[bonusMatchKey]))
		b.NameCriteria.FromJSON(encoding.Object(data[bonusNameKey]))
		b.CategoryCriteria.FromJSON(encoding.Object(data[bonusCategoryKey]))
	case enum.WeaponDamageBonusType:
		b.WeaponSelectionType = enum.WeaponSelectionTypeFromString(encoding.String(data[bonusSelectionTypeKey]))
		b.IsPercent = encoding.Bool(data[bonusIsPercentKey])
		b.SpecializationCriteria.FromJSON(encoding.Object(data[bonusSpecializationKey]))
		switch b.WeaponSelectionType {
		case enum.WeaponsWithNameWeaponSelect:
			b.NameCriteria.FromJSON(encoding.Object(data[bonusNameKey]))
			b.CategoryCriteria.FromJSON(encoding.Object(data[bonusCategoryKey]))
		case enum.WeaponsWithRequiredSkillWeaponSelect:
			b.NameCriteria.FromJSON(encoding.Object(data[bonusNameKey]))
			b.RelativeLevelCriteria.FromJSON(encoding.Object(data[bonusLevelKey]))
			b.CategoryCriteria.FromJSON(encoding.Object(data[bonusCategoryKey]))
		}
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
		if b.SkillSelectionType != enum.ThisWeaponSkillSelect {
			encoding.ToKeyedJSON(&b.NameCriteria, bonusNameKey, encoder)
			encoding.ToKeyedJSON(&b.CategoryCriteria, bonusCategoryKey, encoder)
		}
	case enum.SkillPointBonusType:
		encoding.ToKeyedJSON(&b.NameCriteria, bonusNameKey, encoder)
		encoding.ToKeyedJSON(&b.SpecializationCriteria, bonusSpecializationKey, encoder)
		encoding.ToKeyedJSON(&b.CategoryCriteria, bonusCategoryKey, encoder)
	case enum.SpellBonusType, enum.SpellPointBonusType:
		encoder.KeyedString(bonusMatchKey, b.SpellMatchType.Key(), false, false)
		encoding.ToKeyedJSON(&b.NameCriteria, bonusNameKey, encoder)
		encoding.ToKeyedJSON(&b.CategoryCriteria, bonusCategoryKey, encoder)
	case enum.WeaponDamageBonusType:
		encoder.KeyedString(bonusSelectionTypeKey, b.WeaponSelectionType.Key(), false, false)
		encoder.KeyedBool(bonusIsPercentKey, b.IsPercent, true)
		encoding.ToKeyedJSON(&b.SpecializationCriteria, bonusSpecializationKey, encoder)
		switch b.WeaponSelectionType {
		case enum.WeaponsWithNameWeaponSelect:
			encoding.ToKeyedJSON(&b.NameCriteria, bonusNameKey, encoder)
			encoding.ToKeyedJSON(&b.CategoryCriteria, bonusCategoryKey, encoder)
		case enum.WeaponsWithRequiredSkillWeaponSelect:
			encoding.ToKeyedJSON(&b.NameCriteria, bonusNameKey, encoder)
			encoding.ToKeyedJSON(&b.RelativeLevelCriteria, bonusLevelKey, encoder)
			encoding.ToKeyedJSON(&b.CategoryCriteria, bonusCategoryKey, encoder)
		}
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
		case enum.SkillsWithNameSkillSelect:
			return b.buildKey(SkillNameID, false)
		case enum.ThisWeaponSkillSelect:
			return ThisWeaponID
		case enum.WeaponsWithNameSkillSelect:
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
		if b.CategoryCriteria.Type != enum.Any {
			return SpellPointsID + "*"
		}
		switch b.SpellMatchType {
		case enum.AllColleges:
			return SpellCollegePointsID
		case enum.CollegeName:
			return b.buildKey(SpellCollegePointsID, true)
		case enum.PowerSourceName:
			return b.buildKey(SpellPowerSourcePointsID, true)
		case enum.SpellName:
			return b.buildKey(SpellPointsID, true)
		default:
			jot.Fatal(1, "invalid match type: ", b.SpellMatchType)
		}
	case enum.WeaponDamageBonusType:
		switch b.WeaponSelectionType {
		case enum.WeaponsWithRequiredSkillWeaponSelect:
			return b.buildKey(WeaponNamedIDPrefix, false)
		case enum.ThisWeaponWeaponSelect:
			return ThisWeaponID
		case enum.WeaponsWithNameWeaponSelect:
			return b.buildKey(SkillNameID, false)
		default:
			jot.Fatal(1, "invalid selection type: ", b.WeaponSelectionType)
		}
	}
	jot.Fatal(1, "invalid bonus type: ", b.Type)
	return "" // Never reached
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
		if b.SkillSelectionType != enum.ThisWeaponSkillSelect {
			ExtractNameables(b.NameCriteria.Qualifier, nameables)
			ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	case enum.SkillPointBonusType:
		ExtractNameables(b.NameCriteria.Qualifier, nameables)
		ExtractNameables(b.SpecializationCriteria.Qualifier, nameables)
		ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.SpellBonusType, enum.SpellPointBonusType:
		if b.SpellMatchType != enum.AllColleges {
			ExtractNameables(b.NameCriteria.Qualifier, nameables)
		}
		ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.WeaponDamageBonusType:
		ExtractNameables(b.SpecializationCriteria.Qualifier, nameables)
		if b.WeaponSelectionType != enum.ThisWeaponWeaponSelect {
			ExtractNameables(b.NameCriteria.Qualifier, nameables)
			ExtractNameables(b.SpecializationCriteria.Qualifier, nameables)
			ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	}
}

// ApplyNameableKeys implements Feature.
func (b *Bonus) ApplyNameableKeys(nameables map[string]string) {
	switch b.Type {
	case enum.ConditionalModifierBonusType, enum.ReactionBonusType:
		b.Situation = ApplyNameables(b.Situation, nameables)
	case enum.SkillBonusType:
		b.SpecializationCriteria.Qualifier = ApplyNameables(b.SpecializationCriteria.Qualifier, nameables)
		if b.SkillSelectionType != enum.ThisWeaponSkillSelect {
			b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
			b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	case enum.SkillPointBonusType:
		b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
		b.SpecializationCriteria.Qualifier = ApplyNameables(b.SpecializationCriteria.Qualifier, nameables)
		b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.SpellBonusType, enum.SpellPointBonusType:
		if b.SpellMatchType != enum.AllColleges {
			b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
		}
		b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.WeaponDamageBonusType:
		b.SpecializationCriteria.Qualifier = ApplyNameables(b.SpecializationCriteria.Qualifier, nameables)
		if b.WeaponSelectionType != enum.ThisWeaponWeaponSelect {
			b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
			b.SpecializationCriteria.Qualifier = ApplyNameables(b.SpecializationCriteria.Qualifier, nameables)
			b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	}
}

// Normalize implements Feature.
func (b *Bonus) Normalize() {
	if b.Type == enum.DRBonusType {
		s := strings.TrimSpace(b.Specialization)
		if s == "" || strings.EqualFold(s, All) {
			s = All
		}
		b.Specialization = s
	}
}

// AddToTooltip adds this bonus' details to the tooltip.
func (b *Bonus) AddToTooltip(tooltip *xio.ByteBuffer) {
	if tooltip == nil || b.Owner == nil {
		return
	}
	tooltip.WriteByte('\n')
	tooltip.WriteString(b.Owner.String())
	tooltip.WriteString(" [")
	if b.Type == enum.WeaponDamageBonusType {
		tooltip.WriteString(b.Amount.Format(i18n.Text("die")))
		if b.IsPercent {
			tooltip.WriteByte('%')
		}
	} else {
		tooltip.WriteString(b.Amount.Format(i18n.Text("level")))
	}
	switch b.Type {
	case enum.DRBonusType:
		tooltip.WriteString(i18n.Text(" against "))
		tooltip.WriteString(b.Specialization)
		tooltip.WriteString(i18n.Text(" attacks"))
	case enum.SkillPointBonusType:
		if b.Amount.Amount == f64d4.One {
			tooltip.WriteString(i18n.Text(" pt"))
		} else {
			tooltip.WriteString(i18n.Text(" pts"))
		}
	}
	tooltip.WriteByte(']')
}
