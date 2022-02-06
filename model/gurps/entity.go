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
	"math"
	"strings"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/gurps/datafile"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

var _ eval.VariableResolver = &Entity{}

// EntityData holds the Entity data that is written to disk.
type EntityData struct {
	Type             datafile.EntityType
	Profile          *Profile
	SheetSettings    *SheetSettings
	Attributes       map[string]*Attribute
	Skills           []*Skill
	CarriedEquipment []*Equipment
	OtherEquipment   []*Equipment
	CreatedOn        jio.Time
	ModifiedOn       jio.Time
}

// Entity holds the base information for various types of entities: PC, NPC, Creature, etc.
type Entity struct {
	EntityData
	LiftingStrengthBonus       fixed.F64d4
	StrikingStrengthBonus      fixed.F64d4
	ThrowingStrengthBonus      fixed.F64d4
	ParryBonus                 fixed.F64d4
	BlockBonus                 fixed.F64d4
	featureMap                 map[string][]feature.Feature
	variableResolverExclusions map[string]bool
}

// MarshalJSON implements json.Marshaler.
func (e *Entity) MarshalJSON() ([]byte, error) {
	return json.Marshal(&e.EntityData)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *Entity) UnmarshalJSON(data []byte) error {
	e.EntityData = EntityData{}
	if err := json.Unmarshal(data, &e.EntityData); err != nil {
		return err
	}
	return nil
}

// StrengthOrZero returns the current ST value, or zero if no such attribute exists.
func (e *Entity) StrengthOrZero() fixed.F64d4 {
	if stAttr, exists := e.Attributes[gid.Strength]; exists {
		return stAttr.Current()
	}
	return 0
}

// Thrust returns the thrust value for the current strength.
func (e *Entity) Thrust() *dice.Dice {
	return e.ThrustFor((e.StrengthOrZero() + e.StrikingStrengthBonus).AsInt())
}

// ThrustFor returns the thrust value for the provided strength.
func (e *Entity) ThrustFor(st int) *dice.Dice {
	return e.SheetSettings.DamageProgression.Thrust(st)
}

// Swing returns the swing value for the current strength.
func (e *Entity) Swing() *dice.Dice {
	return e.SwingFor((e.StrengthOrZero() + e.StrikingStrengthBonus).AsInt())
}

// SwingFor returns the swing value for the provided strength.
func (e *Entity) SwingFor(st int) *dice.Dice {
	return e.SheetSettings.DamageProgression.Swing(st)
}

// SkillPointComparedBonusFor returns the total bonus for the matching skill point bonuses.
func (e *Entity) SkillPointComparedBonusFor(id, name, specialization string, categories []string, tooltip *xio.ByteBuffer) fixed.F64d4 {
	var total fixed.F64d4
	for _, f := range e.featureMap[strings.ToLower(id)] {
		if bonus, ok := f.(*feature.SkillPointBonus); ok &&
			bonus.NameCriteria.Matches(name) &&
			bonus.SpecializationCriteria.Matches(specialization) &&
			bonus.CategoryCriteria.Matches(categories...) {
			total += bonus.AdjustedAmount()
			bonus.AddToTooltip(tooltip)
		}
	}
	return total
}

// AddWeaponComparedDamageBonusesFor adds the bonuses for matching weapons that match to the map. If 'm' is nil, it will
// be created. The provided map (or the newly created one) will be returned.
func (e *Entity) AddWeaponComparedDamageBonusesFor(id, nameQualifier, specializationQualifier string, categoryQualifier []string, dieCount int, tooltip *xio.ByteBuffer, m map[*feature.WeaponDamageBonus]bool) map[*feature.WeaponDamageBonus]bool {
	if m == nil {
		m = make(map[*feature.WeaponDamageBonus]bool)
	}
	for _, one := range e.WeaponComparedDamageBonusesFor(id, nameQualifier, specializationQualifier, categoryQualifier, dieCount, tooltip) {
		m[one] = true
	}
	return m
}

// WeaponComparedDamageBonusesFor returns the bonuses for matching weapons that match.
func (e *Entity) WeaponComparedDamageBonusesFor(id, nameQualifier, specializationQualifier string, categoryQualifier []string, dieCount int, tooltip *xio.ByteBuffer) []*feature.WeaponDamageBonus {
	rsl := fixed.F64d4Min
	for _, sk := range e.SkillNamed(nameQualifier, specializationQualifier, true, nil) {
		if rsl < sk.Level.RelativeLevel {
			rsl = sk.Level.RelativeLevel
		}
	}
	if rsl == fixed.F64d4Min {
		return nil
	}
	var bonuses []*feature.WeaponDamageBonus
	for _, f := range e.featureMap[strings.ToLower(id)] {
		//nolint:gocritic // Don't want to invert the logic here
		if bonus, ok := f.(*feature.WeaponDamageBonus); ok &&
			bonus.NameCriteria.Matches(nameQualifier) &&
			bonus.SpecializationCriteria.Matches(specializationQualifier) &&
			bonus.RelativeLevelCriteria.Matches(rsl) &&
			bonus.CategoryCriteria.Matches(categoryQualifier...) {
			bonuses = append(bonuses, bonus)
			level := bonus.LeveledAmount.Level
			bonus.LeveledAmount.Level = fixed.F64d4FromInt(dieCount)
			bonus.AddToTooltip(tooltip)
			bonus.LeveledAmount.Level = level
		}
	}
	return bonuses
}

// BonusFor returns the total bonus for the given ID.
func (e *Entity) BonusFor(id string, tooltip *xio.ByteBuffer) fixed.F64d4 {
	var total fixed.F64d4
	for _, f := range e.featureMap[strings.ToLower(id)] {
		if bonus, ok := f.(feature.Bonus); ok {
			if _, ok = bonus.(*feature.WeaponDamageBonus); !ok {
				total += bonus.AdjustedAmount()
				bonus.AddToTooltip(tooltip)
			}
		}
	}
	return total
}

// AddDRBonusesFor locates any active DR bonuses and adds them to the map. If 'drMap' is nil, it will be created. The
// provided map (or the newly created one) will be returned.
func (e *Entity) AddDRBonusesFor(id string, tooltip *xio.ByteBuffer, drMap map[string]int) map[string]int {
	if drMap == nil {
		drMap = make(map[string]int)
	}
	if list, exists := e.featureMap[strings.ToLower(id)]; exists {
		for _, one := range list {
			if drBonus, ok := one.(*feature.DRBonus); ok {
				drMap[strings.ToLower(drBonus.Specialization)] += drBonus.AdjustedAmount().AsInt()
				drBonus.AddToTooltip(tooltip)
			}
		}
	}
	return drMap
}

// SkillNamed returns a list of skills that match.
func (e *Entity) SkillNamed(name, specialization string, requirePoints bool, excludes map[string]bool) []*Skill {
	var list []*Skill
	TraverseSkills(func(sk *Skill) bool {
		if !sk.Container() && !excludes[sk.String()] {
			if !requirePoints || sk.Type == gid.Technique || sk.AdjustedPoints() > 0 {
				if strings.EqualFold(sk.Name, name) {
					if specialization == "" || strings.EqualFold(sk.Specialization, specialization) {
						list = append(list, sk)
					}
				}
			}
		}
		return false
	}, e.Skills...)
	return list
}

// AddNamedWeaponDamageBonusesFor adds the bonuses for matching weapons that match to the map. If 'm' is nil, it will
// be created. The provided map (or the newly created one) will be returned.
func (e *Entity) AddNamedWeaponDamageBonusesFor(id, nameQualifier, usageQualifier string, categoryQualifier []string, dieCount int, tooltip *xio.ByteBuffer, m map[*feature.WeaponDamageBonus]bool) map[*feature.WeaponDamageBonus]bool {
	if m == nil {
		m = make(map[*feature.WeaponDamageBonus]bool)
	}
	for _, one := range e.NamedWeaponDamageBonusesFor(id, nameQualifier, usageQualifier, categoryQualifier, dieCount, tooltip) {
		m[one] = true
	}
	return m
}

// NamedWeaponDamageBonusesFor returns the bonuses for matching weapons.
func (e *Entity) NamedWeaponDamageBonusesFor(id, nameQualifier, usageQualifier string, categoryQualifiers []string, dieCount int, tooltip *xio.ByteBuffer) []*feature.WeaponDamageBonus {
	list := e.featureMap[strings.ToLower(id)]
	if len(list) == 0 {
		return nil
	}
	var bonuses []*feature.WeaponDamageBonus
	for _, one := range list {
		//nolint:gocritic // Don't want to invert the logic here
		if bonus, ok := one.(*feature.WeaponDamageBonus); ok &&
			bonus.SelectionType == weapon.WithName &&
			bonus.NameCriteria.Matches(nameQualifier) &&
			bonus.SpecializationCriteria.Matches(usageQualifier) &&
			bonus.CategoryCriteria.Matches(categoryQualifiers...) {
			bonuses = append(bonuses, bonus)
			level := bonus.LeveledAmount.Level
			bonus.LeveledAmount.Level = fixed.F64d4FromInt(dieCount)
			bonus.AddToTooltip(tooltip)
			bonus.LeveledAmount.Level = level
		}
	}
	return bonuses
}

// NamedWeaponSkillBonusesFor returns the bonuses for matching weapons.
func (e *Entity) NamedWeaponSkillBonusesFor(id, nameQualifier, usageQualifier string, categoryQualifiers []string, tooltip *xio.ByteBuffer) []*feature.SkillBonus {
	list := e.featureMap[strings.ToLower(id)]
	if len(list) == 0 {
		return nil
	}
	var bonuses []*feature.SkillBonus
	for _, one := range list {
		if bonus, ok := one.(*feature.SkillBonus); ok &&
			bonus.SelectionType == skill.WeaponsWithName &&
			bonus.NameCriteria.Matches(nameQualifier) &&
			bonus.SpecializationCriteria.Matches(usageQualifier) &&
			bonus.CategoryCriteria.Matches(categoryQualifiers...) {
			bonuses = append(bonuses, bonus)
			bonus.AddToTooltip(tooltip)
		}
	}
	return bonuses
}

// EncumbranceLevel returns the current Encumbrance level.
func (e *Entity) EncumbranceLevel(forSkills bool) datafile.Encumbrance {
	carried := e.WeightCarried(forSkills)
	for _, one := range datafile.AllEncumbrances {
		if carried <= e.MaximumCarry(one) {
			return one
		}
	}
	return datafile.ExtraHeavy
}

// WeightCarried returns the carried weight.
func (e *Entity) WeightCarried(forSkills bool) measure.Weight {
	var total measure.Weight
	for _, one := range e.CarriedEquipment {
		total += one.ExtendedWeight(forSkills, e.SheetSettings.DefaultWeightUnits)
	}
	return total
}

// MaximumCarry returns the maximum amount the Entity can carry for the specified encumbrance level.
func (e *Entity) MaximumCarry(encumbrance datafile.Encumbrance) measure.Weight {
	return measure.Weight(fixed.F64d4(e.BasicLift()).Mul(encumbrance.WeightMultiplier()))
}

// BasicLift returns the entity's Basic Lift.
func (e *Entity) BasicLift() measure.Weight {
	st := (e.StrengthOrZero() + e.LiftingStrengthBonus).Trunc()
	if IsThresholdOpMet(attribute.HalveST, e.Attributes) {
		st = st.Div(fxp.Two)
		if st != st.Trunc() {
			st = st.Trunc() + fxp.One
		}
	}
	if st < fxp.One {
		return 0
	}
	var v fixed.F64d4
	if e.SheetSettings.DamageProgression == attribute.KnowingYourOwnStrength {
		var diff fixed.F64d4
		if st > fxp.Nineteen {
			diff = st.Div(fxp.Ten).Trunc() - fxp.One
			st -= diff.Mul(fxp.Ten)
		}
		v = fixed.F64d4FromFloat64(math.Pow(10, st.AsFloat64()/10)).Mul(fxp.Two)
		if st <= fxp.Six {
			v = fxp.Round(v.Mul(fxp.Ten)).Div(fxp.Ten)
		} else {
			v = fxp.Round(v)
		}
		v = v.Mul(fixed.F64d4FromFloat64(math.Pow(10, diff.AsFloat64())))
	} else {
		v = st.Mul(st).Div(fxp.Five)
	}
	if v >= fxp.Ten {
		v = fxp.Round(v)
	}
	return measure.Weight(v.Mul(fxp.Ten).Trunc().Div(fxp.Ten))
}

// ResolveVariable implements eval.VariableResolver.
func (e *Entity) ResolveVariable(variableName string) string {
	if e.variableResolverExclusions[variableName] {
		jot.Warn("attempt to resolve variable via itself: $" + variableName)
		return ""
	}
	if e.variableResolverExclusions == nil {
		e.variableResolverExclusions = make(map[string]bool)
	}
	e.variableResolverExclusions[variableName] = true
	defer func() { delete(e.variableResolverExclusions, variableName) }()
	if gid.SizeModifier == variableName {
		return e.Profile.AdjustedSizeModifier().String()
	}
	parts := strings.SplitN(variableName, ".", 2)
	attr := e.Attributes[parts[0]]
	if attr == nil {
		jot.Warn("no such variable: $" + variableName)
		return ""
	}
	def := attr.AttributeDef()
	if def == nil {
		jot.Warn("no such variable definition: $" + variableName)
		return ""
	}
	if def.Type == attribute.Pool && len(parts) > 1 {
		switch parts[1] {
		case "current":
			return attr.Current().Trunc().String()
		case "maximum":
			return attr.Maximum().Trunc().String()
		default:
			jot.Warn("no such variable: $" + variableName)
			return ""
		}
	}
	return attr.Current().String()
}

// ResolveAttribute resolves the given attribute ID to its current value, or fixed.F64d4Min if it doesn't exist.
func (e *Entity) ResolveAttribute(attrID string) fixed.F64d4 {
	if e != nil && e.Type == datafile.PC {
		if a, ok := e.Attributes[attrID]; ok {
			return a.Current()
		}
		if v, err := fixed.F64d4FromString(attrID); err == nil {
			return v
		}
	}
	return fixed.F64d4Min
}

// PreservesUserDesc returns true if the user description field should be preserved when written to disk. Normally, only
// character sheets should return true for this.
func (e *Entity) PreservesUserDesc() bool {
	return e.Type == datafile.PC
}
