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
	Type                   FeatureType
	Limitation             AttributeBonusLimitation
	SkillSelectionType     SkillSelectionType
	SpellMatchType         SpellMatchType
	WeaponSelectionType    WeaponSelectionType
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
func NewFeature(featureType FeatureType, entity *Entity) *Feature {
	f := &Feature{
		Type:   featureType,
		Amount: LeveledAmount{Amount: 1},
	}
	switch featureType {
	case AttributeBonus:
		f.Attribute = DefaultAttributeIDFor(entity)
	case ConditionalModifierBonus:
		f.Situation = i18n.Text("triggering condition")
	case ContainedWeightReduction:
		f.Reduction = "0%"
	case CostReduction:
		f.Attribute = DefaultAttributeIDFor(entity)
		f.Amount.Amount = fixed.F64d4FromInt64(40)
	case DRBonus:
		f.Location = "torso"
		f.Specialization = All
	case ReactionBonus:
		f.Situation = i18n.Text("from others")
	case SkillBonus:
		f.SkillSelectionType = SkillsWithNameSkillSelect
		fallthrough
	case SkillPointBonus:
		f.NameCriteria.Type = criteria.Is
		f.SpecializationCriteria.Type = criteria.Any
		f.CategoryCriteria.Type = criteria.Any
	case SpellBonus, SpellPointBonus:
		f.SpellMatchType = AllColleges
		f.NameCriteria.Type = criteria.Is
		f.CategoryCriteria.Type = criteria.Any
	case WeaponDamageBonus:
		f.WeaponSelectionType = WeaponsWithRequiredSkillWeaponSelect
		f.NameCriteria.Type = criteria.Is
		f.SpecializationCriteria.Type = criteria.Any
		f.RelativeLevelCriteria.Type = criteria.AnyNumber
		f.CategoryCriteria.Type = criteria.Any
	default:
		jot.Fatal(1, "invalid feature type: ", f.Type)
	}
	return f
}

// NewFeatureFromJSON creates a new Feature from JSON.
func NewFeatureFromJSON(data map[string]interface{}) *Feature {
	f := &Feature{Type: FeatureTypeFromString(encoding.String(data[featureTypeKey]))}
	f.Amount.FromJSON(data)
	switch f.Type {
	case AttributeBonus:
		f.Attribute = encoding.String(data[featureAttributeKey])
		f.Limitation = AttributeBonusLimitationFromKey(encoding.String(data[featureLimitationKey]))
	case ConditionalModifierBonus, ReactionBonus:
		f.Situation = encoding.String(data[featureSituationKey])
	case ContainedWeightReduction:
		f.Reduction = strings.TrimSpace(encoding.String(data[featureReductionKey]))
	case CostReduction:
		f.Attribute = encoding.String(data[featureAttributeKey])
		f.Amount.Amount = encoding.Number(data[featurePercentageKey])
	case DRBonus:
		f.Location = encoding.String(data[featureLocationKey])
		f.Specialization = encoding.String(data[featureSpecializationKey])
	case SkillBonus:
		f.SkillSelectionType = SkillSelectionTypeFromString(encoding.String(data[featureSelectionTypeKey]))
		f.SpecializationCriteria.FromJSON(encoding.Object(data[featureSpecializationKey]))
		if f.SkillSelectionType != ThisWeaponSkillSelect {
			f.NameCriteria.FromJSON(encoding.Object(data[featureNameKey]))
			f.CategoryCriteria.FromJSON(encoding.Object(data[featureCategoryKey]))
		}
	case SkillPointBonus:
		f.NameCriteria.FromJSON(encoding.Object(data[featureNameKey]))
		f.SpecializationCriteria.FromJSON(encoding.Object(data[featureSpecializationKey]))
		f.CategoryCriteria.FromJSON(encoding.Object(data[featureCategoryKey]))
	case SpellBonus, SpellPointBonus:
		f.SpellMatchType = SpellMatchTypeFromString(encoding.String(data[featureMatchKey]))
		f.NameCriteria.FromJSON(encoding.Object(data[featureNameKey]))
		f.CategoryCriteria.FromJSON(encoding.Object(data[featureCategoryKey]))
	case WeaponDamageBonus:
		f.WeaponSelectionType = WeaponSelectionTypeFromString(encoding.String(data[featureSelectionTypeKey]))
		f.IsPercent = encoding.Bool(data[featureIsPercentKey])
		f.SpecializationCriteria.FromJSON(encoding.Object(data[featureSpecializationKey]))
		switch f.WeaponSelectionType {
		case WeaponsWithNameWeaponSelect:
			f.NameCriteria.FromJSON(encoding.Object(data[featureNameKey]))
			f.CategoryCriteria.FromJSON(encoding.Object(data[featureCategoryKey]))
		case WeaponsWithRequiredSkillWeaponSelect:
			f.NameCriteria.FromJSON(encoding.Object(data[featureNameKey]))
			f.RelativeLevelCriteria.FromJSON(encoding.Object(data[featureLevelKey]))
			f.CategoryCriteria.FromJSON(encoding.Object(data[featureCategoryKey]))
		}
	default:
		jot.Fatal(1, "invalid feature type: ", f.Type)
	}
	f.Normalize()
	return f
}

// ToJSON emits this Feature as JSON.
func (f *Feature) ToJSON(encoder *encoding.JSONEncoder) {
	f.Normalize()
	encoder.StartObject()
	encoder.KeyedString(featureTypeKey, f.Type.Key(), false, false)
	f.Amount.ToInlineJSON(encoder)
	switch f.Type {
	case AttributeBonus:
		encoder.KeyedString(featureAttributeKey, f.Attribute, false, false)
		if f.Limitation != None {
			encoder.KeyedString(featureLimitationKey, f.Limitation.Key(), false, false)
		}
	case ConditionalModifierBonus, ReactionBonus:
		encoder.KeyedString(featureSituationKey, f.Situation, true, true)
	case ContainedWeightReduction:
		encoder.KeyedString(featureReductionKey, f.Reduction, true, true)
	case CostReduction:
		encoder.KeyedString(featureAttributeKey, f.Attribute, false, false)
		encoder.KeyedNumber(featurePercentageKey, f.Amount.Amount, false)
	case DRBonus:
		encoder.KeyedString(featureLocationKey, f.Location, true, true)
		encoder.KeyedString(featureSpecializationKey, f.Specialization, true, true)
	case SkillBonus:
		encoder.KeyedString(featureSelectionTypeKey, f.SkillSelectionType.Key(), false, false)
		encoding.ToKeyedJSON(&f.SpecializationCriteria, featureSpecializationKey, encoder)
		if f.SkillSelectionType != ThisWeaponSkillSelect {
			encoding.ToKeyedJSON(&f.NameCriteria, featureNameKey, encoder)
			encoding.ToKeyedJSON(&f.CategoryCriteria, featureCategoryKey, encoder)
		}
	case SkillPointBonus:
		encoding.ToKeyedJSON(&f.NameCriteria, featureNameKey, encoder)
		encoding.ToKeyedJSON(&f.SpecializationCriteria, featureSpecializationKey, encoder)
		encoding.ToKeyedJSON(&f.CategoryCriteria, featureCategoryKey, encoder)
	case SpellBonus, SpellPointBonus:
		encoder.KeyedString(featureMatchKey, f.SpellMatchType.Key(), false, false)
		encoding.ToKeyedJSON(&f.NameCriteria, featureNameKey, encoder)
		encoding.ToKeyedJSON(&f.CategoryCriteria, featureCategoryKey, encoder)
	case WeaponDamageBonus:
		encoder.KeyedString(featureSelectionTypeKey, f.WeaponSelectionType.Key(), false, false)
		encoder.KeyedBool(featureIsPercentKey, f.IsPercent, true)
		encoding.ToKeyedJSON(&f.SpecializationCriteria, featureSpecializationKey, encoder)
		switch f.WeaponSelectionType {
		case WeaponsWithNameWeaponSelect:
			encoding.ToKeyedJSON(&f.NameCriteria, featureNameKey, encoder)
			encoding.ToKeyedJSON(&f.CategoryCriteria, featureCategoryKey, encoder)
		case WeaponsWithRequiredSkillWeaponSelect:
			encoding.ToKeyedJSON(&f.NameCriteria, featureNameKey, encoder)
			encoding.ToKeyedJSON(&f.RelativeLevelCriteria, featureLevelKey, encoder)
			encoding.ToKeyedJSON(&f.CategoryCriteria, featureCategoryKey, encoder)
		}
	default:
		jot.Fatal(1, "invalid feature type: ", f.Type)
	}
	encoder.EndObject()
}

// Key returns the key used in the Feature map for things this Feature applies to.
func (f *Feature) Key() string {
	switch f.Type {
	case AttributeBonus:
		key := AttributeIDPrefix + f.Attribute
		if f.Limitation != None {
			key += "." + f.Limitation.Key()
		}
		return key
	case ConditionalModifierBonus:
		return ConditionalModifierBonus.Key()
	case ContainedWeightReduction:
		return ContainedWeightFeatureKey
	case CostReduction:
		return AttributeIDPrefix + f.Attribute
	case DRBonus:
		return HitLocationPrefix + f.Location
	case ReactionBonus:
		return "reaction"
	case SkillBonus:
		switch f.SkillSelectionType {
		case SkillsWithNameSkillSelect:
			return f.buildKey(SkillNameID, false)
		case ThisWeaponSkillSelect:
			return ThisWeaponID
		case WeaponsWithNameSkillSelect:
			return f.buildKey(WeaponNamedIDPrefix, false)
		default:
			jot.Fatal(1, "invalid selection type: ", f.SkillSelectionType)
		}
	case SkillPointBonus:
		return f.buildKey(SkillPointsID, false)
	case SpellBonus:
		if f.CategoryCriteria.Type != criteria.Any {
			return SpellNameID + "*"
		}
		switch f.SpellMatchType {
		case AllColleges:
			return SpellCollegeID
		case CollegeName:
			return f.buildKey(SpellCollegeID, true)
		case PowerSourceName:
			return f.buildKey(SpellPowerSourceID, true)
		case SpellName:
			return f.buildKey(SpellNameID, true)
		default:
			jot.Fatal(1, "invalid match type: ", f.SpellMatchType)
		}
	case SpellPointBonus:
		if f.CategoryCriteria.Type != criteria.Any {
			return SpellPointsID + "*"
		}
		switch f.SpellMatchType {
		case AllColleges:
			return SpellCollegePointsID
		case CollegeName:
			return f.buildKey(SpellCollegePointsID, true)
		case PowerSourceName:
			return f.buildKey(SpellPowerSourcePointsID, true)
		case SpellName:
			return f.buildKey(SpellPointsID, true)
		default:
			jot.Fatal(1, "invalid match type: ", f.SpellMatchType)
		}
	case WeaponDamageBonus:
		switch f.WeaponSelectionType {
		case WeaponsWithRequiredSkillWeaponSelect:
			return f.buildKey(WeaponNamedIDPrefix, false)
		case ThisWeaponWeaponSelect:
			return ThisWeaponID
		case WeaponsWithNameWeaponSelect:
			return f.buildKey(SkillNameID, false)
		default:
			jot.Fatal(1, "invalid selection type: ", f.WeaponSelectionType)
		}
	}
	jot.Fatal(1, "invalid feature type: ", f.Type)
	return "" // Never reached
}

func (f *Feature) buildKey(prefix string, considerNameCriteriaOnly bool) string {
	if f.NameCriteria.Type == criteria.Is && (considerNameCriteriaOnly ||
		(f.SpecializationCriteria.Type == criteria.Any && f.CategoryCriteria.Type == criteria.Any)) {
		return prefix + "/" + f.NameCriteria.Qualifier
	}
	return prefix + "*"
}

// FillWithNameableKeys adds any nameable keys found in this Feature to the provided map.
func (f *Feature) FillWithNameableKeys(nameables map[string]string) {
	switch f.Type {
	case ConditionalModifierBonus, ReactionBonus:
		ExtractNameables(f.Situation, nameables)
	case SkillBonus:
		ExtractNameables(f.SpecializationCriteria.Qualifier, nameables)
		if f.SkillSelectionType != ThisWeaponSkillSelect {
			ExtractNameables(f.NameCriteria.Qualifier, nameables)
			ExtractNameables(f.CategoryCriteria.Qualifier, nameables)
		}
	case SkillPointBonus:
		ExtractNameables(f.NameCriteria.Qualifier, nameables)
		ExtractNameables(f.SpecializationCriteria.Qualifier, nameables)
		ExtractNameables(f.CategoryCriteria.Qualifier, nameables)
	case SpellBonus, SpellPointBonus:
		if f.SpellMatchType != AllColleges {
			ExtractNameables(f.NameCriteria.Qualifier, nameables)
		}
		ExtractNameables(f.CategoryCriteria.Qualifier, nameables)
	case WeaponDamageBonus:
		ExtractNameables(f.SpecializationCriteria.Qualifier, nameables)
		if f.WeaponSelectionType != ThisWeaponWeaponSelect {
			ExtractNameables(f.NameCriteria.Qualifier, nameables)
			ExtractNameables(f.SpecializationCriteria.Qualifier, nameables)
			ExtractNameables(f.CategoryCriteria.Qualifier, nameables)
		}
	}
}

// ApplyNameableKeys replaces any nameable keys found in this Feature with the corresponding values in the provided map.
func (f *Feature) ApplyNameableKeys(nameables map[string]string) {
	switch f.Type {
	case ConditionalModifierBonus, ReactionBonus:
		f.Situation = ApplyNameables(f.Situation, nameables)
	case SkillBonus:
		f.SpecializationCriteria.Qualifier = ApplyNameables(f.SpecializationCriteria.Qualifier, nameables)
		if f.SkillSelectionType != ThisWeaponSkillSelect {
			f.NameCriteria.Qualifier = ApplyNameables(f.NameCriteria.Qualifier, nameables)
			f.CategoryCriteria.Qualifier = ApplyNameables(f.CategoryCriteria.Qualifier, nameables)
		}
	case SkillPointBonus:
		f.NameCriteria.Qualifier = ApplyNameables(f.NameCriteria.Qualifier, nameables)
		f.SpecializationCriteria.Qualifier = ApplyNameables(f.SpecializationCriteria.Qualifier, nameables)
		f.CategoryCriteria.Qualifier = ApplyNameables(f.CategoryCriteria.Qualifier, nameables)
	case SpellBonus, SpellPointBonus:
		if f.SpellMatchType != AllColleges {
			f.NameCriteria.Qualifier = ApplyNameables(f.NameCriteria.Qualifier, nameables)
		}
		f.CategoryCriteria.Qualifier = ApplyNameables(f.CategoryCriteria.Qualifier, nameables)
	case WeaponDamageBonus:
		f.SpecializationCriteria.Qualifier = ApplyNameables(f.SpecializationCriteria.Qualifier, nameables)
		if f.WeaponSelectionType != ThisWeaponWeaponSelect {
			f.NameCriteria.Qualifier = ApplyNameables(f.NameCriteria.Qualifier, nameables)
			f.SpecializationCriteria.Qualifier = ApplyNameables(f.SpecializationCriteria.Qualifier, nameables)
			f.CategoryCriteria.Qualifier = ApplyNameables(f.CategoryCriteria.Qualifier, nameables)
		}
	}
}

// Normalize adjusts the data to it preferred representation.
func (f *Feature) Normalize() {
	if f.Type == DRBonus {
		s := strings.TrimSpace(f.Specialization)
		if s == "" || strings.EqualFold(s, All) {
			s = All
		}
		f.Specialization = s
	}
}

// AddToTooltip adds this feature's bonus details to the tooltip.
func (f *Feature) AddToTooltip(tooltip *xio.ByteBuffer) {
	if tooltip == nil || f.Owner == nil || f.Type == CostReduction || f.Type == ContainedWeightReduction {
		return
	}
	tooltip.WriteByte('\n')
	tooltip.WriteString(f.Owner.String())
	tooltip.WriteString(" [")
	if f.Type == WeaponDamageBonus {
		tooltip.WriteString(f.Amount.Format(i18n.Text("die")))
		if f.IsPercent {
			tooltip.WriteByte('%')
		}
	} else {
		tooltip.WriteString(f.Amount.Format(i18n.Text("level")))
	}
	switch f.Type {
	case DRBonus:
		tooltip.WriteString(i18n.Text(" against "))
		tooltip.WriteString(f.Specialization)
		tooltip.WriteString(i18n.Text(" attacks"))
	case SkillPointBonus:
		if f.Amount.Amount == f64d4.One {
			tooltip.WriteString(i18n.Text(" pt"))
		} else {
			tooltip.WriteString(i18n.Text(" pts"))
		}
	}
	tooltip.WriteByte(']')
}

// IsPercentageReduction returns true if this is a percentage reduction and not a fixed amount. Only applicable to
// enum.ContainedWeightReduction.
func (f *Feature) IsPercentageReduction() bool {
	return strings.HasSuffix(f.Reduction, "%")
}

// PercentageReduction returns the percentage (where 1% is 1, not 0.01) the weight should be reduced by. Will return 0 if
// this is not a percentage. Only applicable to enum.ContainedWeightReduction.
func (f *Feature) PercentageReduction() fixed.F64d4 {
	if !f.IsPercentageReduction() {
		return 0
	}
	return fixed.F64d4FromStringForced(f.Reduction[:len(f.Reduction)-1])
}

// FixedReduction returns the fixed amount the weight should be reduced by. Will return 0 if this is a percentage. Only
// applicable to enum.ContainedWeightReduction.
func (f *Feature) FixedReduction(defUnits measure.WeightUnits) measure.Weight {
	if f.IsPercentageReduction() {
		return 0
	}
	return measure.WeightFromStringForced(f.Reduction, defUnits)
}
