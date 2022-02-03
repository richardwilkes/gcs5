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

	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/gurps/datafile"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

var _ eval.VariableResolver = &Entity{}

// Entity holds the base information for various types of entities: PC, NPC, Creature, etc.
type Entity struct {
	Type                  datafile.EntityType
	Profile               *Profile
	SheetSettings         *SheetSettings
	LiftingStrengthBonus  fixed.F64d4
	StrikingStrengthBonus fixed.F64d4
	ThrowingStrengthBonus fixed.F64d4
	ParryBonus            fixed.F64d4
	BlockBonus            fixed.F64d4
	Attributes            map[string]*Attribute
	Skills                []*Skill
	CarriedEquipment      []*Equipment
	featureMap            map[string][]feature.Feature
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

// AddDRBonusesFor locates any active DR bonuses and adds them to the map. If 'drMap' isn't nil, it will be returned.
func (e *Entity) AddDRBonusesFor(id string, tooltip *xio.ByteBuffer, drMap map[string]int) map[string]int {
	if drMap == nil {
		drMap = make(map[string]int)
	}
	if list, exists := e.featureMap[strings.ToLower(id)]; exists {
		for _, one := range list {
			if drBonus, ok := one.(*feature.DRBonus); ok {
				drMap[strings.ToLower(drBonus.Specialization)] += int(drBonus.AdjustedAmount().AsInt64())
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
			if !requirePoints || sk.Type == techniqueTypeKey || sk.AdjustedPoints() > 0 {
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

// NamedWeaponSkillBonusesFor returns the bonuses for matching weapons.
func (e *Entity) NamedWeaponSkillBonusesFor(id, nameQualifier, usageQualifier string, categoryQualifiers []string, tooltip *xio.ByteBuffer) []*feature.SkillBonus {
	list := e.featureMap[strings.ToLower(id)]
	if len(list) == 0 {
		return nil
	}
	var bonuses []*feature.SkillBonus
	for _, bonus := range list {
		if skillBonus, ok := bonus.(*feature.SkillBonus); ok &&
			skillBonus.SelectionType == skill.WeaponsWithName &&
			skillBonus.NameCriteria.Matches(nameQualifier) &&
			skillBonus.SpecializationCriteria.Matches(usageQualifier) &&
			skillBonus.CategoryCriteria.Matches(categoryQualifiers...) {
			bonuses = append(bonuses, skillBonus)
			skillBonus.AddToTooltip(tooltip)
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
	stAttr, ok := e.Attributes["st"]
	if !ok {
		return 0
	}
	st := (stAttr.Current() + e.LiftingStrengthBonus).Trunc()
	if IsThresholdOpMet(attribute.HalveST, e.Attributes) {
		st = st.Div(f64d4.Two)
		if st != st.Trunc() {
			st = st.Trunc() + f64d4.One
		}
	}
	if st < f64d4.One {
		return 0
	}
	var v fixed.F64d4
	if e.SheetSettings.DamageProgression == attribute.KnowingYourOwnStrength {
		var diff fixed.F64d4
		if st > f64d4.Nineteen {
			diff = st.Div(f64d4.Ten).Trunc() - f64d4.One
			st -= diff.Mul(f64d4.Ten)
		}
		v = fixed.F64d4FromFloat64(math.Pow(10, st.AsFloat64()/10)).Mul(f64d4.Two)
		if st <= f64d4.Six {
			v = f64d4.Round(v.Mul(f64d4.Ten)).Div(f64d4.Ten)
		} else {
			v = f64d4.Round(v)
		}
		v = v.Mul(fixed.F64d4FromFloat64(math.Pow(10, diff.AsFloat64())))
	} else {
		v = st.Mul(st).Div(f64d4.Five)
	}
	if v >= f64d4.Ten {
		v = f64d4.Round(v)
	}
	return measure.Weight(v.Mul(f64d4.Ten).Trunc().Div(f64d4.Ten))
}

// ResolveVariable implements eval.VariableResolver.
func (e *Entity) ResolveVariable(variableName string) string {
	// TODO implement me
	return variableName
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
	// TODO: Implement... should only return true for sheets
	return true
}
