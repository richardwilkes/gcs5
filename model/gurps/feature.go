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

	"github.com/richardwilkes/gcs/model/criteria"
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/enum"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

const (
	featureTypeKey           = "type"
	featureAttributeKey      = "attribute"
	featureCategoryKey       = "category"
	featureIsPercentKey      = "percent"
	featureLevelKey          = "level"
	featureLimitationKey     = "limitation"
	featureLocationKey       = "location"
	featureMatchKey          = "match"
	featureNameKey           = "name"
	featurePercentageKey     = "percentage"
	featureReductionKey      = "reduction"
	featureSelectionTypeKey  = "selection_type"
	featureSituationKey      = "situation"
	featureSpecializationKey = "specialization"
)

const (
	// All is the DR specialization key for DR that affects everything.
	All = "all"
	// ThisWeaponID holds the ID for "this weapon".
	ThisWeaponID = "\u0001"
	// WeaponNamedIDPrefix the prefix for "weapon named" IDs.
	WeaponNamedIDPrefix = "weapon_named."
	// ContainedWeightFeatureKey is the key used in the Feature map for things this Feature applies to.
	ContainedWeightFeatureKey = "equipment.weight.sum"
)

// Feature holds data that affects another object.
type Feature struct {
	Type                   enum.FeatureType
	Limitation             enum.AttributeBonusLimitation
	SkillSelectionType     enum.SkillSelectionType
	SpellMatchType         enum.SpellMatchType
	WeaponSelectionType    enum.WeaponSelectionType
	IsPercent              bool
	Amount                 LeveledAmount
	Attribute              string
	Situation              string
	Location               string
	Specialization         string
	Reduction              string
	NameCriteria           criteria.String
	SpecializationCriteria criteria.String
	CategoryCriteria       criteria.String
	RelativeLevelCriteria  criteria.Numeric
	Owner                  fmt.Stringer
}

// NewFeature creates a new Feature for the given entity, which may be nil.
func NewFeature(bonusType enum.FeatureType, entity *Entity) *Feature {
	b := &Feature{
		Type:   bonusType,
		Amount: LeveledAmount{Amount: 1},
	}
	switch bonusType {
	case enum.AttributeBonus:
		b.Attribute = DefaultAttributeIDFor(entity)
	case enum.ConditionalModifierBonus:
		b.Situation = i18n.Text("triggering condition")
	case enum.ContainedWeightReduction:
		b.Reduction = "0%"
	case enum.CostReduction:
		b.Attribute = DefaultAttributeIDFor(entity)
		b.Amount.Amount = fixed.F64d4FromInt64(40)
	case enum.DRBonus:
		b.Location = "torso"
		b.Specialization = All
	case enum.ReactionBonus:
		b.Situation = i18n.Text("from others")
	case enum.SkillBonus:
		b.SkillSelectionType = enum.SkillsWithNameSkillSelect
		fallthrough
	case enum.SkillPointBonus:
		b.NameCriteria.Type = enum.Is
		b.SpecializationCriteria.Type = enum.Any
		b.CategoryCriteria.Type = enum.Any
	case enum.SpellBonus, enum.SpellPointBonus:
		b.SpellMatchType = enum.AllColleges
		b.NameCriteria.Type = enum.Is
		b.CategoryCriteria.Type = enum.Any
	case enum.WeaponDamageBonus:
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

// NewFeatureFromJSON creates a new Feature from JSON.
func NewFeatureFromJSON(data map[string]interface{}) *Feature {
	b := &Feature{Type: enum.FeatureTypeFromString(encoding.String(data[featureTypeKey]))}
	b.Amount.FromJSON(data)
	switch b.Type {
	case enum.AttributeBonus:
		b.Attribute = encoding.String(data[featureAttributeKey])
		b.Limitation = enum.AttributeBonusLimitationFromString(encoding.String(data[featureLimitationKey]))
	case enum.ConditionalModifierBonus, enum.ReactionBonus:
		b.Situation = encoding.String(data[featureSituationKey])
	case enum.ContainedWeightReduction:
		b.Reduction = strings.TrimSpace(encoding.String(data[featureReductionKey]))
	case enum.CostReduction:
		b.Attribute = encoding.String(data[featureAttributeKey])
		b.Amount.Amount = encoding.Number(data[featurePercentageKey])
	case enum.DRBonus:
		b.Location = encoding.String(data[featureLocationKey])
		b.Specialization = encoding.String(data[featureSpecializationKey])
	case enum.SkillBonus:
		b.SkillSelectionType = enum.SkillSelectionTypeFromString(encoding.String(data[featureSelectionTypeKey]))
		b.SpecializationCriteria.FromJSON(encoding.Object(data[featureSpecializationKey]))
		if b.SkillSelectionType != enum.ThisWeaponSkillSelect {
			b.NameCriteria.FromJSON(encoding.Object(data[featureNameKey]))
			b.CategoryCriteria.FromJSON(encoding.Object(data[featureCategoryKey]))
		}
	case enum.SkillPointBonus:
		b.NameCriteria.FromJSON(encoding.Object(data[featureNameKey]))
		b.SpecializationCriteria.FromJSON(encoding.Object(data[featureSpecializationKey]))
		b.CategoryCriteria.FromJSON(encoding.Object(data[featureCategoryKey]))
	case enum.SpellBonus, enum.SpellPointBonus:
		b.SpellMatchType = enum.SpellMatchTypeFromString(encoding.String(data[featureMatchKey]))
		b.NameCriteria.FromJSON(encoding.Object(data[featureNameKey]))
		b.CategoryCriteria.FromJSON(encoding.Object(data[featureCategoryKey]))
	case enum.WeaponDamageBonus:
		b.WeaponSelectionType = enum.WeaponSelectionTypeFromString(encoding.String(data[featureSelectionTypeKey]))
		b.IsPercent = encoding.Bool(data[featureIsPercentKey])
		b.SpecializationCriteria.FromJSON(encoding.Object(data[featureSpecializationKey]))
		switch b.WeaponSelectionType {
		case enum.WeaponsWithNameWeaponSelect:
			b.NameCriteria.FromJSON(encoding.Object(data[featureNameKey]))
			b.CategoryCriteria.FromJSON(encoding.Object(data[featureCategoryKey]))
		case enum.WeaponsWithRequiredSkillWeaponSelect:
			b.NameCriteria.FromJSON(encoding.Object(data[featureNameKey]))
			b.RelativeLevelCriteria.FromJSON(encoding.Object(data[featureLevelKey]))
			b.CategoryCriteria.FromJSON(encoding.Object(data[featureCategoryKey]))
		}
	default:
		jot.Fatal(1, "invalid bonus type: ", b.Type)
	}
	b.Normalize()
	return b
}

// ToJSON emits this Feature as JSON.
func (b *Feature) ToJSON(encoder *encoding.JSONEncoder) {
	b.Normalize()
	encoder.StartObject()
	encoder.KeyedString(featureTypeKey, b.Type.Key(), false, false)
	b.Amount.ToInlineJSON(encoder)
	switch b.Type {
	case enum.AttributeBonus:
		encoder.KeyedString(featureAttributeKey, b.Attribute, false, false)
		if b.Limitation != enum.None {
			encoder.KeyedString(featureLimitationKey, b.Limitation.Key(), false, false)
		}
	case enum.ConditionalModifierBonus, enum.ReactionBonus:
		encoder.KeyedString(featureSituationKey, b.Situation, true, true)
	case enum.ContainedWeightReduction:
		encoder.KeyedString(featureReductionKey, b.Reduction, true, true)
	case enum.CostReduction:
		encoder.KeyedString(featureAttributeKey, b.Attribute, false, false)
		encoder.KeyedNumber(featurePercentageKey, b.Amount.Amount, false)
	case enum.DRBonus:
		encoder.KeyedString(featureLocationKey, b.Location, true, true)
		encoder.KeyedString(featureSpecializationKey, b.Specialization, true, true)
	case enum.SkillBonus:
		encoder.KeyedString(featureSelectionTypeKey, b.SkillSelectionType.Key(), false, false)
		encoding.ToKeyedJSON(&b.SpecializationCriteria, featureSpecializationKey, encoder)
		if b.SkillSelectionType != enum.ThisWeaponSkillSelect {
			encoding.ToKeyedJSON(&b.NameCriteria, featureNameKey, encoder)
			encoding.ToKeyedJSON(&b.CategoryCriteria, featureCategoryKey, encoder)
		}
	case enum.SkillPointBonus:
		encoding.ToKeyedJSON(&b.NameCriteria, featureNameKey, encoder)
		encoding.ToKeyedJSON(&b.SpecializationCriteria, featureSpecializationKey, encoder)
		encoding.ToKeyedJSON(&b.CategoryCriteria, featureCategoryKey, encoder)
	case enum.SpellBonus, enum.SpellPointBonus:
		encoder.KeyedString(featureMatchKey, b.SpellMatchType.Key(), false, false)
		encoding.ToKeyedJSON(&b.NameCriteria, featureNameKey, encoder)
		encoding.ToKeyedJSON(&b.CategoryCriteria, featureCategoryKey, encoder)
	case enum.WeaponDamageBonus:
		encoder.KeyedString(featureSelectionTypeKey, b.WeaponSelectionType.Key(), false, false)
		encoder.KeyedBool(featureIsPercentKey, b.IsPercent, true)
		encoding.ToKeyedJSON(&b.SpecializationCriteria, featureSpecializationKey, encoder)
		switch b.WeaponSelectionType {
		case enum.WeaponsWithNameWeaponSelect:
			encoding.ToKeyedJSON(&b.NameCriteria, featureNameKey, encoder)
			encoding.ToKeyedJSON(&b.CategoryCriteria, featureCategoryKey, encoder)
		case enum.WeaponsWithRequiredSkillWeaponSelect:
			encoding.ToKeyedJSON(&b.NameCriteria, featureNameKey, encoder)
			encoding.ToKeyedJSON(&b.RelativeLevelCriteria, featureLevelKey, encoder)
			encoding.ToKeyedJSON(&b.CategoryCriteria, featureCategoryKey, encoder)
		}
	default:
		jot.Fatal(1, "invalid bonus type: ", b.Type)
	}
	encoder.EndObject()
}

// Key returns the key used in the Feature map for things this Feature applies to.
func (b *Feature) Key() string {
	switch b.Type {
	case enum.AttributeBonus:
		key := AttributeIDPrefix + b.Attribute
		if b.Limitation != enum.None {
			key += "." + b.Limitation.Key()
		}
		return key
	case enum.ConditionalModifierBonus:
		return enum.ConditionalModifierBonus.Key()
	case enum.ContainedWeightReduction:
		return ContainedWeightFeatureKey
	case enum.CostReduction:
		return AttributeIDPrefix + b.Attribute
	case enum.DRBonus:
		return HitLocationPrefix + b.Location
	case enum.ReactionBonus:
		return "reaction"
	case enum.SkillBonus:
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
	case enum.SkillPointBonus:
		return b.buildKey(SkillPointsID, false)
	case enum.SpellBonus:
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
	case enum.SpellPointBonus:
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
	case enum.WeaponDamageBonus:
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

func (b *Feature) buildKey(prefix string, considerNameCriteriaOnly bool) string {
	if b.NameCriteria.Type == enum.Is && (considerNameCriteriaOnly ||
		(b.SpecializationCriteria.Type == enum.Any && b.CategoryCriteria.Type == enum.Any)) {
		return prefix + "/" + b.NameCriteria.Qualifier
	}
	return prefix + "*"
}

// FillWithNameableKeys adds any nameable keys found in this Feature to the provided map.
func (b *Feature) FillWithNameableKeys(nameables map[string]string) {
	switch b.Type {
	case enum.ConditionalModifierBonus, enum.ReactionBonus:
		ExtractNameables(b.Situation, nameables)
	case enum.SkillBonus:
		ExtractNameables(b.SpecializationCriteria.Qualifier, nameables)
		if b.SkillSelectionType != enum.ThisWeaponSkillSelect {
			ExtractNameables(b.NameCriteria.Qualifier, nameables)
			ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	case enum.SkillPointBonus:
		ExtractNameables(b.NameCriteria.Qualifier, nameables)
		ExtractNameables(b.SpecializationCriteria.Qualifier, nameables)
		ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.SpellBonus, enum.SpellPointBonus:
		if b.SpellMatchType != enum.AllColleges {
			ExtractNameables(b.NameCriteria.Qualifier, nameables)
		}
		ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.WeaponDamageBonus:
		ExtractNameables(b.SpecializationCriteria.Qualifier, nameables)
		if b.WeaponSelectionType != enum.ThisWeaponWeaponSelect {
			ExtractNameables(b.NameCriteria.Qualifier, nameables)
			ExtractNameables(b.SpecializationCriteria.Qualifier, nameables)
			ExtractNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	}
}

// ApplyNameableKeys replaces any nameable keys found in this Feature with the corresponding values in the provided map.
func (b *Feature) ApplyNameableKeys(nameables map[string]string) {
	switch b.Type {
	case enum.ConditionalModifierBonus, enum.ReactionBonus:
		b.Situation = ApplyNameables(b.Situation, nameables)
	case enum.SkillBonus:
		b.SpecializationCriteria.Qualifier = ApplyNameables(b.SpecializationCriteria.Qualifier, nameables)
		if b.SkillSelectionType != enum.ThisWeaponSkillSelect {
			b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
			b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	case enum.SkillPointBonus:
		b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
		b.SpecializationCriteria.Qualifier = ApplyNameables(b.SpecializationCriteria.Qualifier, nameables)
		b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.SpellBonus, enum.SpellPointBonus:
		if b.SpellMatchType != enum.AllColleges {
			b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
		}
		b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
	case enum.WeaponDamageBonus:
		b.SpecializationCriteria.Qualifier = ApplyNameables(b.SpecializationCriteria.Qualifier, nameables)
		if b.WeaponSelectionType != enum.ThisWeaponWeaponSelect {
			b.NameCriteria.Qualifier = ApplyNameables(b.NameCriteria.Qualifier, nameables)
			b.SpecializationCriteria.Qualifier = ApplyNameables(b.SpecializationCriteria.Qualifier, nameables)
			b.CategoryCriteria.Qualifier = ApplyNameables(b.CategoryCriteria.Qualifier, nameables)
		}
	}
}

// Normalize adjusts the data to it preferred representation.
func (b *Feature) Normalize() {
	if b.Type == enum.DRBonus {
		s := strings.TrimSpace(b.Specialization)
		if s == "" || strings.EqualFold(s, All) {
			s = All
		}
		b.Specialization = s
	}
}

// AddToTooltip adds this feature's bonus details to the tooltip.
func (b *Feature) AddToTooltip(tooltip *xio.ByteBuffer) {
	if tooltip == nil || b.Owner == nil || b.Type == enum.CostReduction || b.Type == enum.ContainedWeightReduction {
		return
	}
	tooltip.WriteByte('\n')
	tooltip.WriteString(b.Owner.String())
	tooltip.WriteString(" [")
	if b.Type == enum.WeaponDamageBonus {
		tooltip.WriteString(b.Amount.Format(i18n.Text("die")))
		if b.IsPercent {
			tooltip.WriteByte('%')
		}
	} else {
		tooltip.WriteString(b.Amount.Format(i18n.Text("level")))
	}
	switch b.Type {
	case enum.DRBonus:
		tooltip.WriteString(i18n.Text(" against "))
		tooltip.WriteString(b.Specialization)
		tooltip.WriteString(i18n.Text(" attacks"))
	case enum.SkillPointBonus:
		if b.Amount.Amount == f64d4.One {
			tooltip.WriteString(i18n.Text(" pt"))
		} else {
			tooltip.WriteString(i18n.Text(" pts"))
		}
	}
	tooltip.WriteByte(']')
}

// IsPercentageReduction returns true if this is a percentage reduction and not a fixed amount. Only applicable to
// enum.ContainedWeightReduction.
func (b *Feature) IsPercentageReduction() bool {
	return strings.HasSuffix(b.Reduction, "%")
}

// PercentageReduction returns the percentage (where 1% is 1, not 0.01) the weight should be reduced by. Will return 0 if
// this is not a percentage. Only applicable to enum.ContainedWeightReduction.
func (b *Feature) PercentageReduction() fixed.F64d4 {
	if !b.IsPercentageReduction() {
		return 0
	}
	return fixed.F64d4FromStringForced(b.Reduction[:len(b.Reduction)-1])
}

// FixedReduction returns the fixed amount the weight should be reduced by. Will return 0 if this is a percentage. Only
// applicable to enum.ContainedWeightReduction.
func (b *Feature) FixedReduction(defUnits measure.WeightUnits) measure.Weight {
	if b.IsPercentageReduction() {
		return 0
	}
	return measure.WeightFromStringForced(b.Reduction, defUnits)
}
