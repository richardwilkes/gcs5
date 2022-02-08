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
	"context"
	"fmt"
	"io/fs"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/advantage"
	"github.com/richardwilkes/gcs/model/gurps/ancestry"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/gurps/datafile"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/gcs/model/id"
	"github.com/richardwilkes/gcs/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/eval"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

var _ eval.VariableResolver = &Entity{}

// EntityData holds the Entity data that is written to disk.
type EntityData struct {
	Type             datafile.Type          `json:"type"`
	ID               uuid.UUID              `json:"id"`
	TotalPoints      fixed.F64d4            `json:"total_points"`
	Profile          *Profile               `json:"profile,omitempty"`
	SheetSettings    *SheetSettings         `json:"settings,omitempty"`
	Attributes       *Attributes            `json:"attributes,omitempty"`
	Advantages       []*Advantage           `json:"advantages,omitempty"`
	Skills           []*Skill               `json:"skills,omitempty"`
	Spells           []*Spell               `json:"spells,omitempty"`
	CarriedEquipment []*Equipment           `json:"equipment,omitempty"`
	OtherEquipment   []*Equipment           `json:"other_equipment,omitempty"`
	Notes            []*Note                `json:"notes,omitempty"`
	CreatedOn        jio.Time               `json:"created_date"`
	ModifiedOn       jio.Time               `json:"modified_date"`
	ThirdParty       map[string]interface{} `json:"third_party,omitempty"`
}

// Entity holds the base information for various types of entities: PC, NPC, Creature, etc.
type Entity struct {
	EntityData
	LiftingStrengthBonus       fixed.F64d4
	StrikingStrengthBonus      fixed.F64d4
	ThrowingStrengthBonus      fixed.F64d4
	DodgeBonus                 fixed.F64d4
	ParryBonus                 fixed.F64d4
	BlockBonus                 fixed.F64d4
	featureMap                 map[string][]feature.Feature
	variableResolverExclusions map[string]bool
}

// NewEntityFromFile loads an Entity from a file.
func NewEntityFromFile(fileSystem fs.FS, filePath string) (*Entity, error) {
	var entity Entity
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &entity); err != nil {
		return nil, errs.NewWithCause("invalid entity file: "+filePath, err)
	}
	return &entity, nil
}

// NewEntity creates a new Entity.
func NewEntity(entityType datafile.Type) *Entity {
	entity := &Entity{
		EntityData: EntityData{
			Type:        entityType,
			ID:          id.NewUUID(),
			TotalPoints: fixed.F64d4FromInt(SettingsProvider.GeneralSettings().InitialPoints),
			Profile:     &Profile{},
			Advantages:  nil,
			CreatedOn:   jio.Now(),
		},
	}
	entity.SheetSettings = SettingsProvider.SheetSettings().Clone(entity)
	entity.Attributes = NewAttributes(entity)
	if SettingsProvider.GeneralSettings().AutoFillProfile {
		entity.Profile.AutoFill(entity)
	}
	entity.ModifiedOn = entity.CreatedOn
	entity.Recalculate()
	return entity
}

// Save the Entity to a file as JSON.
func (e *Entity) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, e)
}

// MarshalJSON implements json.Marshaler.
func (e *Entity) MarshalJSON() ([]byte, error) {
	type calc struct {
		Swing                 *dice.Dice     `json:"swing"`
		Thrust                *dice.Dice     `json:"thrust"`
		BasicLift             measure.Weight `json:"basic_lift"`
		LiftingStrengthBonus  fixed.F64d4    `json:"lifting_st_bonus,omitempty"`
		StrikingStrengthBonus fixed.F64d4    `json:"striking_st_bonus,omitempty"`
		ThrowingStrengthBonus fixed.F64d4    `json:"throwing_st_bonus,omitempty"`
		DodgeBonus            fixed.F64d4    `json:"dodge_bonus,omitempty"`
		ParryBonus            fixed.F64d4    `json:"parry_bonus,omitempty"`
		BlockBonus            fixed.F64d4    `json:"block_bonus,omitempty"`
		Move                  []int          `json:"move"`
		Dodge                 []int          `json:"dodge"`
	}
	data := struct {
		EntityData
		Calc calc `json:"calc"`
	}{
		EntityData: e.EntityData,
		Calc: calc{
			Swing:                 e.Swing(),
			Thrust:                e.Thrust(),
			BasicLift:             e.BasicLift(),
			LiftingStrengthBonus:  e.LiftingStrengthBonus,
			StrikingStrengthBonus: e.StrikingStrengthBonus,
			ThrowingStrengthBonus: e.ThrowingStrengthBonus,
			DodgeBonus:            e.DodgeBonus,
			ParryBonus:            e.ParryBonus,
			BlockBonus:            e.BlockBonus,
			Move:                  make([]int, len(datafile.AllEncumbrance)),
			Dodge:                 make([]int, len(datafile.AllEncumbrance)),
		},
	}
	for i, one := range datafile.AllEncumbrance {
		data.Calc.Move[i] = e.Move(one)
		data.Calc.Dodge[i] = e.Dodge(one)
	}
	return json.Marshal(&data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *Entity) UnmarshalJSON(data []byte) error {
	e.EntityData = EntityData{}
	if err := json.Unmarshal(data, &e.EntityData); err != nil {
		return err
	}
	return nil
}

// Recalculate the statistics.
func (e *Entity) Recalculate() {
	e.UpdateSkills()
	e.UpdateSpells()
	for i := 0; i < 5; i++ {
		// Unfortunately, there are what amount to circular references in the GURPS logic, so we need to potentially run
		// though this process a few times until things stabilize. To avoid a potential endless loop, though, we cap the
		// iterations.
		e.processFeatures()
		e.processPrereqs()
		skillsChanged := e.UpdateSkills()
		spellsChanged := e.UpdateSpells()
		if !skillsChanged && !spellsChanged {
			break
		}
	}
}

func (e *Entity) processFeatures() {
	m := make(map[string][]feature.Feature)
	TraverseAdvantages(func(a *Advantage) bool {
		for _, f := range a.Features {
			processFeature(a, m, f, a.Levels.Max(0))
		}
		for _, f := range a.CRAdj.Features(a.CR) {
			processFeature(a, m, f, a.Levels.Max(0))
		}
		for _, mod := range a.Modifiers {
			if !mod.Disabled {
				for _, f := range mod.Features {
					processFeature(a, m, f, mod.Levels)
				}
			}
		}
		return false
	}, true, e.Advantages...)
	TraverseSkills(func(s *Skill) bool {
		for _, f := range s.Features {
			processFeature(s, m, f, 0)
		}
		return false
	}, e.Skills...)
	TraverseEquipment(func(eqp *Equipment) bool {
		if !eqp.Equipped || eqp.Quantity <= 0 {
			return false
		}
		for _, f := range eqp.Features {
			processFeature(eqp, m, f, 0)
		}
		for _, mod := range eqp.Modifiers {
			if !mod.Disabled {
				for _, f := range mod.Features {
					processFeature(eqp, m, f, 0)
				}
			}
		}
		return false
	}, e.CarriedEquipment...)
}

func processFeature(parent fmt.Stringer, m map[string][]feature.Feature, f feature.Feature, levels fixed.F64d4) {
	key := strings.ToLower(f.FeatureMapKey())
	list := m[key]
	if bonus, ok := f.(feature.Bonus); ok {
		bonus.SetParent(parent)
		bonus.SetLevel(levels)
	}
	list = append(list, f)
	m[key] = list
}

func (e *Entity) processPrereqs() {
	const prefix = "\n- "
	notMetPrefix := i18n.Text("Prerequisites have not been met:")
	TraverseAdvantages(func(a *Advantage) bool {
		var tooltip xio.ByteBuffer
		if a.Satisfied = a.Prereq.Satisfied(e, a, &tooltip, prefix); a.Satisfied {
			a.UnsatisfiedReason = ""
		} else {
			a.UnsatisfiedReason = notMetPrefix + tooltip.String()
		}
		return false
	}, true, e.Advantages...)
	TraverseSkills(func(s *Skill) bool {
		var tooltip xio.ByteBuffer
		s.Satisfied = s.Prereq.Satisfied(e, s, &tooltip, prefix)
		if s.Satisfied && s.Type == gid.Technique {
			s.Satisfied = s.TechniqueSatisfied(&tooltip, prefix)
		}
		if s.Satisfied {
			s.UnsatisfiedReason = ""
		} else {
			s.UnsatisfiedReason = notMetPrefix + tooltip.String()
		}
		return false
	}, e.Skills...)
	TraverseSpells(func(s *Spell) bool {
		var tooltip xio.ByteBuffer
		s.Satisfied = s.Prereq.Satisfied(e, s, &tooltip, prefix)
		if s.Satisfied && s.Type == gid.RitualMagicSpell {
			s.Satisfied = s.RitualMagicSatisfied(&tooltip, prefix)
		}
		if s.Satisfied {
			s.UnsatisfiedReason = ""
		} else {
			s.UnsatisfiedReason = notMetPrefix + tooltip.String()
		}
		return false
	}, e.Spells...)
	equipmentFunc := func(eqp *Equipment) bool {
		var tooltip xio.ByteBuffer
		if eqp.Satisfied = eqp.Prereq.Satisfied(e, eqp, &tooltip, prefix); eqp.Satisfied {
			eqp.UnsatisfiedReason = ""
		} else {
			eqp.UnsatisfiedReason = notMetPrefix + tooltip.String()
		}
		return false
	}
	TraverseEquipment(equipmentFunc, e.CarriedEquipment...)
	TraverseEquipment(equipmentFunc, e.OtherEquipment...)
}

// UpdateSkills updates the levels of all skills.
func (e *Entity) UpdateSkills() bool {
	changed := false
	TraverseSkills(func(s *Skill) bool {
		if s.UpdateLevel() {
			changed = true
		}
		return false
	}, e.Skills...)
	return changed
}

// UpdateSpells updates the levels of all spells.
func (e *Entity) UpdateSpells() bool {
	changed := false
	TraverseSpells(func(s *Spell) bool {
		if s.UpdateLevel() {
			changed = true
		}
		return false
	}, e.Spells...)
	return changed
}

// AttributePoints returns the number of points spent on attributes.
func (e *Entity) AttributePoints() fixed.F64d4 {
	var total fixed.F64d4
	for _, attr := range e.Attributes.Set {
		total += attr.PointCost()
	}
	return total
}

// AdvantagePoints returns the number of points spent on advantages.
func (e *Entity) AdvantagePoints() (ad, disad, race, perk, quirk fixed.F64d4) {
	for _, one := range e.Advantages {
		a, d, r, p, q := calculateSingleAdvantagePoints(one)
		ad += a
		disad += d
		race += r
		perk += p
		quirk += q
	}
	return
}

func calculateSingleAdvantagePoints(adq *Advantage) (ad, disad, race, perk, quirk fixed.F64d4) {
	if adq.Container() {
		switch adq.ContainerType {
		case advantage.Group:
			for _, child := range adq.Children {
				a, d, r, p, q := calculateSingleAdvantagePoints(child)
				ad += a
				disad += d
				race += r
				perk += p
				quirk += q
			}
			return
		case advantage.Race:
			return 0, 0, adq.AdjustedPoints(), 0, 0
		}
	}
	pts := adq.AdjustedPoints()
	switch {
	case pts == fxp.One:
		perk += pts
	case pts == fxp.NegOne:
		quirk += pts
	case pts > 0:
		ad += pts
	case pts < 0:
		disad += pts
	}
	return
}

// SkillPoints returns the number of points spent on skills.
func (e *Entity) SkillPoints() fixed.F64d4 {
	var total fixed.F64d4
	TraverseSkills(func(s *Skill) bool {
		if !s.Container() {
			total += s.Points
		}
		return false
	}, e.Skills...)
	return total
}

// SpellPoints returns the number of points spent on spells.
func (e *Entity) SpellPoints() fixed.F64d4 {
	var total fixed.F64d4
	TraverseSpells(func(s *Spell) bool {
		if !s.Container() {
			total += s.Points
		}
		return false
	}, e.Spells...)
	return total
}

// WealthCarried returns the current wealth being carried.
func (e *Entity) WealthCarried() fixed.F64d4 {
	var value fixed.F64d4
	for _, one := range e.CarriedEquipment {
		value += one.ExtendedValue()
	}
	return value
}

// WealthNotCarried returns the current wealth not being carried.
func (e *Entity) WealthNotCarried() fixed.F64d4 {
	var value fixed.F64d4
	for _, one := range e.OtherEquipment {
		value += one.ExtendedValue()
	}
	return value
}

// StrengthOrZero returns the current ST value, or zero if no such attribute exists.
func (e *Entity) StrengthOrZero() fixed.F64d4 {
	return e.ResolveAttributeCurrent(gid.Strength).Max(0)
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

// AddWeaponComparedDamageBonusesFor adds the bonuses for matching weapons that match to the map. If 'm' is nil, it will
// be created. The provided map (or the newly created one) will be returned.
func (e *Entity) AddWeaponComparedDamageBonusesFor(featureID, nameQualifier, specializationQualifier string, categoryQualifier []string, dieCount int, tooltip *xio.ByteBuffer, m map[*feature.WeaponDamageBonus]bool) map[*feature.WeaponDamageBonus]bool {
	if m == nil {
		m = make(map[*feature.WeaponDamageBonus]bool)
	}
	for _, one := range e.WeaponComparedDamageBonusesFor(featureID, nameQualifier, specializationQualifier, categoryQualifier, dieCount, tooltip) {
		m[one] = true
	}
	return m
}

// WeaponComparedDamageBonusesFor returns the bonuses for matching weapons that match.
func (e *Entity) WeaponComparedDamageBonusesFor(featureID, nameQualifier, specializationQualifier string, categoryQualifier []string, dieCount int, tooltip *xio.ByteBuffer) []*feature.WeaponDamageBonus {
	rsl := fixed.F64d4Min
	for _, sk := range e.SkillNamed(nameQualifier, specializationQualifier, true, nil) {
		if rsl < sk.LevelData.RelativeLevel {
			rsl = sk.LevelData.RelativeLevel
		}
	}
	if rsl == fixed.F64d4Min {
		return nil
	}
	var bonuses []*feature.WeaponDamageBonus
	for _, f := range e.featureMap[strings.ToLower(featureID)] {
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
func (e *Entity) BonusFor(featureID string, tooltip *xio.ByteBuffer) fixed.F64d4 {
	var total fixed.F64d4
	for _, f := range e.featureMap[strings.ToLower(featureID)] {
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
func (e *Entity) AddDRBonusesFor(featureID string, tooltip *xio.ByteBuffer, drMap map[string]int) map[string]int {
	if drMap == nil {
		drMap = make(map[string]int)
	}
	if list, exists := e.featureMap[strings.ToLower(featureID)]; exists {
		for _, one := range list {
			if drBonus, ok := one.(*feature.DRBonus); ok {
				drMap[strings.ToLower(drBonus.Specialization)] += drBonus.AdjustedAmount().AsInt()
				drBonus.AddToTooltip(tooltip)
			}
		}
	}
	return drMap
}

// BestSkillNamed returns the best skill that matches.
func (e *Entity) BestSkillNamed(name, specialization string, requirePoints bool, excludes map[string]bool) *Skill {
	var best *Skill
	level := fixed.F64d4Min
	for _, sk := range e.SkillNamed(name, specialization, requirePoints, excludes) {
		skillLevel := sk.Level(excludes)
		if best == nil || level < skillLevel {
			best = sk
			level = skillLevel
		}
	}
	return best
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

// SkillComparedBonusFor returns the total bonus for the matching skill bonuses.
func (e *Entity) SkillComparedBonusFor(featureID, name, specialization string, categories []string, tooltip *xio.ByteBuffer) fixed.F64d4 {
	var total fixed.F64d4
	for _, f := range e.featureMap[strings.ToLower(featureID)] {
		if bonus, ok := f.(*feature.SkillBonus); ok &&
			bonus.NameCriteria.Matches(name) &&
			bonus.SpecializationCriteria.Matches(specialization) &&
			bonus.CategoryCriteria.Matches(categories...) {
			total += bonus.AdjustedAmount()
			bonus.AddToTooltip(tooltip)
		}
	}
	return total
}

// SkillPointComparedBonusFor returns the total bonus for the matching skill point bonuses.
func (e *Entity) SkillPointComparedBonusFor(featureID, name, specialization string, categories []string, tooltip *xio.ByteBuffer) fixed.F64d4 {
	var total fixed.F64d4
	for _, f := range e.featureMap[strings.ToLower(featureID)] {
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

// SpellComparedBonusFor returns the total bonus for the matching spell bonuses.
func (e *Entity) SpellComparedBonusFor(featureID, name string, categories []string, tooltip *xio.ByteBuffer) fixed.F64d4 {
	var total fixed.F64d4
	for _, f := range e.featureMap[strings.ToLower(featureID)] {
		if bonus, ok := f.(*feature.SpellBonus); ok &&
			bonus.NameCriteria.Matches(name) &&
			bonus.CategoryCriteria.Matches(categories...) {
			total += bonus.AdjustedAmount()
			bonus.AddToTooltip(tooltip)
		}
	}
	return total
}

// AddNamedWeaponDamageBonusesFor adds the bonuses for matching weapons that match to the map. If 'm' is nil, it will
// be created. The provided map (or the newly created one) will be returned.
func (e *Entity) AddNamedWeaponDamageBonusesFor(featureID, nameQualifier, usageQualifier string, categoryQualifier []string, dieCount int, tooltip *xio.ByteBuffer, m map[*feature.WeaponDamageBonus]bool) map[*feature.WeaponDamageBonus]bool {
	if m == nil {
		m = make(map[*feature.WeaponDamageBonus]bool)
	}
	for _, one := range e.NamedWeaponDamageBonusesFor(featureID, nameQualifier, usageQualifier, categoryQualifier, dieCount, tooltip) {
		m[one] = true
	}
	return m
}

// NamedWeaponDamageBonusesFor returns the bonuses for matching weapons.
func (e *Entity) NamedWeaponDamageBonusesFor(featureID, nameQualifier, usageQualifier string, categoryQualifiers []string, dieCount int, tooltip *xio.ByteBuffer) []*feature.WeaponDamageBonus {
	list := e.featureMap[strings.ToLower(featureID)]
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
func (e *Entity) NamedWeaponSkillBonusesFor(featureID, nameQualifier, usageQualifier string, categoryQualifiers []string, tooltip *xio.ByteBuffer) []*feature.SkillBonus {
	list := e.featureMap[strings.ToLower(featureID)]
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

// Move returns the current Move value for the given Encumbrance.
func (e *Entity) Move(enc datafile.Encumbrance) int {
	initialMove := e.ResolveAttributeCurrent(gid.BasicMove).Max(0)
	divisor := 2 * xmath.MinInt(CountThresholdOpMet(attribute.HalveMove, e.Attributes), 2)
	if divisor > 0 {
		initialMove = fxp.Ceil(initialMove.Div(fixed.F64d4FromInt(divisor)))
	}
	move := initialMove.Mul(fxp.Ten + fxp.Two.Mul(enc.Penalty())).Div(fxp.Ten).Trunc()
	if move < fxp.One {
		if initialMove > 0 {
			return 1
		}
		return 0
	}
	return move.AsInt()
}

// Dodge returns the current Dodge value for the given Encumbrance.
func (e *Entity) Dodge(enc datafile.Encumbrance) int {
	dodge := e.ResolveAttributeCurrent(gid.BasicSpeed).Max(0)
	divisor := 2 * xmath.MinInt(CountThresholdOpMet(attribute.HalveDodge, e.Attributes), 2)
	if divisor > 0 {
		dodge = fxp.Ceil(dodge.Div(fixed.F64d4FromInt(divisor)))
	}
	return (dodge + enc.Penalty()).Max(fxp.One).AsInt()
}

// EncumbranceLevel returns the current Encumbrance level.
func (e *Entity) EncumbranceLevel(forSkills bool) datafile.Encumbrance {
	carried := e.WeightCarried(forSkills)
	for _, one := range datafile.AllEncumbrance {
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
	attr := e.Attributes.Set[parts[0]]
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

// ResolveAttributeDef resolves the given attribute ID to its AttributeDef, or nil.
func (e *Entity) ResolveAttributeDef(attrID string) *AttributeDef {
	if e != nil && e.Type == datafile.PC {
		if a, ok := e.Attributes.Set[attrID]; ok {
			return a.AttributeDef()
		}
	}
	return nil
}

// ResolveAttributeName resolves the given attribute ID to its name, or <unknown>.
func (e *Entity) ResolveAttributeName(attrID string) string {
	if def := e.ResolveAttributeDef(attrID); def != nil {
		return def.Name
	}
	return i18n.Text("<unknown>")
}

// ResolveAttribute resolves the given attribute ID to its Attribute, or nil.
func (e *Entity) ResolveAttribute(attrID string) *Attribute {
	if e != nil && e.Type == datafile.PC {
		if a, ok := e.Attributes.Set[attrID]; ok {
			return a
		}
	}
	return nil
}

// ResolveAttributeCurrent resolves the given attribute ID to its current value, or fixed.F64d4Min.
func (e *Entity) ResolveAttributeCurrent(attrID string) fixed.F64d4 {
	if a := e.ResolveAttribute(attrID); a != nil {
		return a.Current()
	}
	if v, err := fixed.F64d4FromString(attrID); err == nil {
		return v
	}
	return fixed.F64d4Min
}

// PreservesUserDesc returns true if the user description field should be preserved when written to disk. Normally, only
// character sheets should return true for this.
func (e *Entity) PreservesUserDesc() bool {
	return e.Type == datafile.PC
}

// Ancestry returns the current Ancestry.
func (e *Entity) Ancestry() *ancestry.Ancestry {
	var anc *ancestry.Ancestry
	TraverseAdvantages(func(adq *Advantage) bool {
		if adq.Container() && adq.ContainerType == advantage.Race {
			if anc = ancestry.Lookup(adq.Ancestry, SettingsProvider.Libraries()); anc != nil {
				return true
			}
		}
		return false
	}, true, e.Advantages...)
	if anc == nil {
		if anc = ancestry.Lookup("Human", SettingsProvider.Libraries()); anc == nil {
			jot.Fatal(1, "unable to load default ancestry (Human)")
		}
	}
	return anc
}
