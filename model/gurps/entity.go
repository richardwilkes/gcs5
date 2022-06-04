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
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/richardwilkes/gcs/model/crc"
	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/ancestry"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/gurps/datafile"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/gurps/trait"
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
)

var (
	_ eval.VariableResolver = &Entity{}
	_ ListProvider          = &Entity{}
)

// EntityProvider provides a way to retrieve a (possibly nil) Entity.
type EntityProvider interface {
	Entity() *Entity
}

// EntityData holds the Entity data that is written to disk.
type EntityData struct {
	Type             datafile.Type  `json:"type"`
	Version          int            `json:"version"`
	ID               uuid.UUID      `json:"id"`
	TotalPoints      fxp.Int        `json:"total_points"`
	Profile          *Profile       `json:"profile,omitempty"`
	SheetSettings    *SheetSettings `json:"settings,omitempty"`
	Attributes       *Attributes    `json:"attributes,omitempty"`
	Traits           []*Trait       `json:"advantages,omitempty"`
	Skills           []*Skill       `json:"skills,omitempty"`
	Spells           []*Spell       `json:"spells,omitempty"`
	CarriedEquipment []*Equipment   `json:"equipment,omitempty"`
	OtherEquipment   []*Equipment   `json:"other_equipment,omitempty"`
	Notes            []*Note        `json:"notes,omitempty"`
	CreatedOn        jio.Time       `json:"created_date"`
	ModifiedOn       jio.Time       `json:"modified_date"`
	ThirdParty       map[string]any `json:"third_party,omitempty"`
}

// Entity holds the base information for various types of entities: PC, NPC, Creature, etc.
type Entity struct {
	EntityData
	LiftingStrengthBonus       fxp.Int
	StrikingStrengthBonus      fxp.Int
	ThrowingStrengthBonus      fxp.Int
	DodgeBonus                 fxp.Int
	ParryBonus                 fxp.Int
	BlockBonus                 fxp.Int
	featureMap                 map[string][]feature.Feature
	variableResolverExclusions map[string]bool
}

// NewEntityFromFile loads an Entity from a file.
func NewEntityFromFile(fileSystem fs.FS, filePath string) (*Entity, error) {
	var entity Entity
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &entity); err != nil {
		return nil, errs.NewWithCause(gid.InvalidFileDataMsg, err)
	}
	if err := gid.CheckVersion(entity.Version); err != nil {
		return nil, err
	}
	return &entity, nil
}

// NewEntity creates a new Entity.
func NewEntity(entityType datafile.Type) *Entity {
	entity := &Entity{
		EntityData: EntityData{
			Type:        entityType,
			ID:          id.NewUUID(),
			TotalPoints: SettingsProvider.GeneralSettings().InitialPoints,
			Profile:     &Profile{},
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

// Entity implements EntityProvider.
func (e *Entity) Entity() *Entity {
	return e
}

// Save the Entity to a file as JSON.
func (e *Entity) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, e)
}

// MarshalJSON implements json.Marshaler.
func (e *Entity) MarshalJSON() ([]byte, error) {
	e.Recalculate()
	type calc struct {
		Swing                 *dice.Dice     `json:"swing"`
		Thrust                *dice.Dice     `json:"thrust"`
		BasicLift             measure.Weight `json:"basic_lift"`
		LiftingStrengthBonus  fxp.Int        `json:"lifting_st_bonus,omitempty"`
		StrikingStrengthBonus fxp.Int        `json:"striking_st_bonus,omitempty"`
		ThrowingStrengthBonus fxp.Int        `json:"throwing_st_bonus,omitempty"`
		DodgeBonus            fxp.Int        `json:"dodge_bonus,omitempty"`
		ParryBonus            fxp.Int        `json:"parry_bonus,omitempty"`
		BlockBonus            fxp.Int        `json:"block_bonus,omitempty"`
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
	data.Version = gid.CurrentDataVersion
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
	if e.SheetSettings == nil {
		e.SheetSettings = SettingsProvider.SheetSettings().Clone(e)
	}
	if e.Profile == nil {
		e.Profile = &Profile{}
	}
	if e.Attributes == nil {
		e.Attributes = NewAttributes(e)
	}
	e.Recalculate()
	return nil
}

// Recalculate the statistics.
func (e *Entity) Recalculate() {
	e.ensureAttachments()
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

func (e *Entity) ensureAttachments() {
	e.SheetSettings.SetOwningEntity(e)
	for _, attr := range e.Attributes.Set {
		attr.Entity = e
	}
	for _, one := range e.Traits {
		one.SetOwningEntity(e)
	}
	for _, one := range e.Skills {
		one.SetOwningEntity(e)
	}
	for _, one := range e.Spells {
		one.SetOwningEntity(e)
	}
	for _, one := range e.CarriedEquipment {
		one.SetOwningEntity(e)
	}
	for _, one := range e.OtherEquipment {
		one.SetOwningEntity(e)
	}
	for _, one := range e.Notes {
		one.SetOwningEntity(e)
	}
}

func (e *Entity) processFeatures() {
	m := make(map[string][]feature.Feature)
	TraverseTraits(func(a *Trait) bool {
		if !a.Container() {
			for _, f := range a.Features {
				processFeature(a, m, f, a.Levels.Max(0))
			}
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
	}, true, e.Traits...)
	TraverseSkills(func(s *Skill) bool {
		if !s.Container() {
			for _, f := range s.Features {
				processFeature(s, m, f, 0)
			}
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
	e.featureMap = m
	e.LiftingStrengthBonus = e.BonusFor(feature.AttributeIDPrefix+gid.Strength+"."+attribute.LiftingOnly.Key(), nil).Trunc()
	e.StrikingStrengthBonus = e.BonusFor(feature.AttributeIDPrefix+gid.Strength+"."+attribute.StrikingOnly.Key(), nil).Trunc()
	e.ThrowingStrengthBonus = e.BonusFor(feature.AttributeIDPrefix+gid.Strength+"."+attribute.ThrowingOnly.Key(), nil).Trunc()
	for _, attr := range e.Attributes.Set {
		if def := attr.AttributeDef(); def != nil {
			attrID := feature.AttributeIDPrefix + attr.AttrID
			attr.Bonus = e.BonusFor(attrID, nil)
			if def.Type != attribute.Decimal {
				attr.Bonus = attr.Bonus.Trunc()
			}
			attr.CostReduction = e.CostReductionFor(attrID)
		} else {
			attr.Bonus = 0
			attr.CostReduction = 0
		}
	}
	e.Profile.Update(e)
	e.DodgeBonus = e.BonusFor(feature.AttributeIDPrefix+gid.Dodge, nil).Trunc()
	e.ParryBonus = e.BonusFor(feature.AttributeIDPrefix+gid.Parry, nil).Trunc()
	e.BlockBonus = e.BonusFor(feature.AttributeIDPrefix+gid.Block, nil).Trunc()
}

func processFeature(parent fmt.Stringer, m map[string][]feature.Feature, f feature.Feature, levels fxp.Int) {
	key := strings.ToLower(f.FeatureMapKey())
	list := m[key]
	if bonus, ok := f.(feature.Bonus); ok {
		bonus.SetParent(parent)
		bonus.SetLevel(levels)
	}
	m[key] = append(list, f)
}

func (e *Entity) processPrereqs() {
	const prefix = "\n● "
	notMetPrefix := i18n.Text("Prerequisites have not been met:")
	TraverseTraits(func(a *Trait) bool {
		a.UnsatisfiedReason = ""
		if !a.Container() && a.Prereq != nil {
			var tooltip xio.ByteBuffer
			if !a.Prereq.Satisfied(e, a, &tooltip, prefix) {
				a.UnsatisfiedReason = notMetPrefix + tooltip.String()
			}
		}
		return false
	}, true, e.Traits...)
	TraverseSkills(func(s *Skill) bool {
		s.UnsatisfiedReason = ""
		if !s.Container() {
			var tooltip xio.ByteBuffer
			satisfied := true
			if s.Prereq != nil {
				satisfied = s.Prereq.Satisfied(e, s, &tooltip, prefix)
			}
			if satisfied && s.Type == gid.Technique {
				satisfied = s.TechniqueSatisfied(&tooltip, prefix)
			}
			if !satisfied {
				s.UnsatisfiedReason = notMetPrefix + tooltip.String()
			}
		}
		return false
	}, e.Skills...)
	TraverseSpells(func(s *Spell) bool {
		s.UnsatisfiedReason = ""
		if !s.Container() {
			var tooltip xio.ByteBuffer
			satisfied := true
			if s.Prereq != nil {
				satisfied = s.Prereq.Satisfied(e, s, &tooltip, prefix)
			}
			if satisfied && s.Type == gid.RitualMagicSpell {
				satisfied = s.RitualMagicSatisfied(&tooltip, prefix)
			}
			if !satisfied {
				s.UnsatisfiedReason = notMetPrefix + tooltip.String()
			}
		}
		return false
	}, e.Spells...)
	equipmentFunc := func(eqp *Equipment) bool {
		eqp.UnsatisfiedReason = ""
		if eqp.Prereq != nil {
			var tooltip xio.ByteBuffer
			if !eqp.Prereq.Satisfied(e, eqp, &tooltip, prefix) {
				eqp.UnsatisfiedReason = notMetPrefix + tooltip.String()
			}
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

// SpentPoints returns the number of spent points.
func (e *Entity) SpentPoints() fxp.Int {
	total := e.AttributePoints()
	ad, disad, race, quirk := e.TraitPoints()
	total += ad + disad + race + quirk
	total += e.SkillPoints()
	total += e.SpellPoints()
	return total
}

// UnspentPoints returns the number of unspent points.
func (e *Entity) UnspentPoints() fxp.Int {
	return e.TotalPoints - e.SpentPoints()
}

// SetUnspentPoints sets the number of unspent points.
func (e *Entity) SetUnspentPoints(unspent fxp.Int) {
	if unspent != e.UnspentPoints() {
		// TODO: Need undo logic
		e.TotalPoints = unspent + e.SpentPoints()
	}
}

// AttributePoints returns the number of points spent on attributes.
func (e *Entity) AttributePoints() fxp.Int {
	var total fxp.Int
	for _, attr := range e.Attributes.Set {
		total += attr.PointCost()
	}
	return total
}

// TraitPoints returns the number of points spent on traits.
func (e *Entity) TraitPoints() (ad, disad, race, quirk fxp.Int) {
	for _, one := range e.Traits {
		a, d, r, q := calculateSingleTraitPoints(one)
		ad += a
		disad += d
		race += r
		quirk += q
	}
	return
}

func calculateSingleTraitPoints(t *Trait) (ad, disad, race, quirk fxp.Int) {
	if t.Container() {
		switch t.ContainerType {
		case trait.Group:
			for _, child := range t.Children {
				a, d, r, q := calculateSingleTraitPoints(child)
				ad += a
				disad += d
				race += r
				quirk += q
			}
			return
		case trait.Race:
			return 0, 0, t.AdjustedPoints(), 0
		}
	}
	pts := t.AdjustedPoints()
	switch {
	case pts == -fxp.One:
		quirk += pts
	case pts > 0:
		ad += pts
	case pts < 0:
		disad += pts
	}
	return
}

// SkillPoints returns the number of points spent on skills.
func (e *Entity) SkillPoints() fxp.Int {
	var total fxp.Int
	TraverseSkills(func(s *Skill) bool {
		if !s.Container() {
			total += s.Points
		}
		return false
	}, e.Skills...)
	return total
}

// SpellPoints returns the number of points spent on spells.
func (e *Entity) SpellPoints() fxp.Int {
	var total fxp.Int
	TraverseSpells(func(s *Spell) bool {
		if !s.Container() {
			total += s.Points
		}
		return false
	}, e.Spells...)
	return total
}

// WealthCarried returns the current wealth being carried.
func (e *Entity) WealthCarried() fxp.Int {
	var value fxp.Int
	for _, one := range e.CarriedEquipment {
		value += one.ExtendedValue()
	}
	return value
}

// WealthNotCarried returns the current wealth not being carried.
func (e *Entity) WealthNotCarried() fxp.Int {
	var value fxp.Int
	for _, one := range e.OtherEquipment {
		value += one.ExtendedValue()
	}
	return value
}

// StrengthOrZero returns the current ST value, or zero if no such attribute exists.
func (e *Entity) StrengthOrZero() fxp.Int {
	return e.ResolveAttributeCurrent(gid.Strength).Max(0)
}

// Thrust returns the thrust value for the current strength.
func (e *Entity) Thrust() *dice.Dice {
	return e.ThrustFor(fxp.As[int](e.StrengthOrZero() + e.StrikingStrengthBonus))
}

// ThrustFor returns the thrust value for the provided strength.
func (e *Entity) ThrustFor(st int) *dice.Dice {
	return e.SheetSettings.DamageProgression.Thrust(st)
}

// Swing returns the swing value for the current strength.
func (e *Entity) Swing() *dice.Dice {
	return e.SwingFor(fxp.As[int](e.StrengthOrZero() + e.StrikingStrengthBonus))
}

// SwingFor returns the swing value for the provided strength.
func (e *Entity) SwingFor(st int) *dice.Dice {
	return e.SheetSettings.DamageProgression.Swing(st)
}

// AddWeaponComparedDamageBonusesFor adds the bonuses for matching weapons that match to the map. If 'm' is nil, it will
// be created. The provided map (or the newly created one) will be returned.
func (e *Entity) AddWeaponComparedDamageBonusesFor(featureID, nameQualifier, specializationQualifier string, tagsQualifier []string, dieCount int, tooltip *xio.ByteBuffer, m map[*feature.WeaponDamageBonus]bool) map[*feature.WeaponDamageBonus]bool {
	if m == nil {
		m = make(map[*feature.WeaponDamageBonus]bool)
	}
	for _, one := range e.WeaponComparedDamageBonusesFor(featureID, nameQualifier, specializationQualifier, tagsQualifier, dieCount, tooltip) {
		m[one] = true
	}
	return m
}

// WeaponComparedDamageBonusesFor returns the bonuses for matching weapons that match.
func (e *Entity) WeaponComparedDamageBonusesFor(featureID, nameQualifier, specializationQualifier string, tagsQualifier []string, dieCount int, tooltip *xio.ByteBuffer) []*feature.WeaponDamageBonus {
	rsl := fxp.Min
	for _, sk := range e.SkillNamed(nameQualifier, specializationQualifier, true, nil) {
		if rsl < sk.LevelData.RelativeLevel {
			rsl = sk.LevelData.RelativeLevel
		}
	}
	if rsl == fxp.Min {
		return nil
	}
	var bonuses []*feature.WeaponDamageBonus
	for _, f := range e.featureMap[strings.ToLower(featureID)] {
		//nolint:gocritic // Don't want to invert the logic here
		if bonus, ok := f.(*feature.WeaponDamageBonus); ok &&
			bonus.NameCriteria.Matches(nameQualifier) &&
			bonus.SpecializationCriteria.Matches(specializationQualifier) &&
			bonus.RelativeLevelCriteria.Matches(rsl) &&
			bonus.TagsCriteria.Matches(tagsQualifier...) {
			bonuses = append(bonuses, bonus)
			level := bonus.LeveledAmount.Level
			bonus.LeveledAmount.Level = fxp.From(dieCount)
			bonus.AddToTooltip(tooltip)
			bonus.LeveledAmount.Level = level
		}
	}
	return bonuses
}

// SpellBonusesFor returns the total bonuses for a spell.
func (e *Entity) SpellBonusesFor(featureID, qualifier string, tags []string, tooltip *xio.ByteBuffer) fxp.Int {
	level := e.BonusFor(featureID, tooltip)
	level += e.BonusFor(featureID+"/"+strings.ToLower(qualifier), tooltip)
	level += e.SpellComparedBonusFor(featureID+"*", qualifier, tags, tooltip)
	return level
}

// BestCollegeSpellBonus returns the best college spell bonus for a spell.
func (e *Entity) BestCollegeSpellBonus(tags, colleges []string, tooltip *xio.ByteBuffer) fxp.Int {
	best := fxp.Min
	var bestTooltip string
	for _, college := range colleges {
		var buffer *xio.ByteBuffer
		if tooltip != nil {
			buffer = &xio.ByteBuffer{}
		}
		if pts := e.SpellBonusesFor(feature.SpellCollegeID, college, tags, buffer); best < pts {
			best = pts
			if buffer != nil {
				bestTooltip = buffer.String()
			}
		}
	}
	if tooltip != nil {
		tooltip.WriteString(bestTooltip)
	}
	if best == fxp.Min {
		best = 0
	}
	return best
}

// BonusFor returns the total bonus for the given ID.
func (e *Entity) BonusFor(featureID string, tooltip *xio.ByteBuffer) fxp.Int {
	var total fxp.Int
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

// CostReductionFor returns the total cost reduction for the given ID.
func (e *Entity) CostReductionFor(featureID string) fxp.Int {
	var total fxp.Int
	for _, f := range e.featureMap[strings.ToLower(featureID)] {
		if reduction, ok := f.(*feature.CostReduction); ok {
			total += reduction.Percentage
		}
	}
	if total > fxp.Eighty {
		total = fxp.Eighty
	}
	return total.Max(0)
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
				drMap[strings.ToLower(drBonus.Specialization)] += fxp.As[int](drBonus.AdjustedAmount())
				drBonus.AddToTooltip(tooltip)
			}
		}
	}
	return drMap
}

// BestSkillNamed returns the best skill that matches.
func (e *Entity) BestSkillNamed(name, specialization string, requirePoints bool, excludes map[string]bool) *Skill {
	var best *Skill
	level := fxp.Min
	for _, sk := range e.SkillNamed(name, specialization, requirePoints, excludes) {
		skillLevel := sk.CalculateLevel().Level
		if best == nil || level < skillLevel {
			best = sk
			level = skillLevel
		}
	}
	return best
}

// BaseSkill returns the best skill for the given default, or nil.
func (e *Entity) BaseSkill(def *SkillDefault, requirePoints bool) *Skill {
	if e == nil || def == nil || !def.SkillBased() {
		return nil
	}
	return e.BestSkillNamed(def.Name, def.Specialization, requirePoints, nil)
}

// SkillNamed returns a list of skills that match.
func (e *Entity) SkillNamed(name, specialization string, requirePoints bool, excludes map[string]bool) []*Skill {
	var list []*Skill
	TraverseSkills(func(sk *Skill) bool {
		if !sk.Container() && !excludes[sk.String()] {
			if !requirePoints || sk.Type == gid.Technique || sk.AdjustedPoints(nil) > 0 {
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
func (e *Entity) SkillComparedBonusFor(featureID, name, specialization string, tags []string, tooltip *xio.ByteBuffer) fxp.Int {
	var total fxp.Int
	for _, f := range e.featureMap[strings.ToLower(featureID)] {
		if bonus, ok := f.(*feature.SkillBonus); ok &&
			bonus.NameCriteria.Matches(name) &&
			bonus.SpecializationCriteria.Matches(specialization) &&
			bonus.TagsCriteria.Matches(tags...) {
			total += bonus.AdjustedAmount()
			bonus.AddToTooltip(tooltip)
		}
	}
	return total
}

// SkillPointComparedBonusFor returns the total bonus for the matching skill point bonuses.
func (e *Entity) SkillPointComparedBonusFor(featureID, name, specialization string, tags []string, tooltip *xio.ByteBuffer) fxp.Int {
	var total fxp.Int
	for _, f := range e.featureMap[strings.ToLower(featureID)] {
		if bonus, ok := f.(*feature.SkillPointBonus); ok &&
			bonus.NameCriteria.Matches(name) &&
			bonus.SpecializationCriteria.Matches(specialization) &&
			bonus.TagsCriteria.Matches(tags...) {
			total += bonus.AdjustedAmount()
			bonus.AddToTooltip(tooltip)
		}
	}
	return total
}

// SpellPointBonusesFor returns the total point bonus for the matching spell.
func (e *Entity) SpellPointBonusesFor(featureID, qualifier string, tags []string, tooltip *xio.ByteBuffer) fxp.Int {
	level := e.BonusFor(featureID, tooltip)
	level += e.BonusFor(featureID+"/"+strings.ToLower(qualifier), tooltip)
	level += e.SpellPointComparedBonusFor(featureID+"*", qualifier, tags, tooltip)
	return level
}

// SpellComparedBonusFor returns the total bonus for the matching spell bonuses.
func (e *Entity) SpellComparedBonusFor(featureID, name string, tags []string, tooltip *xio.ByteBuffer) fxp.Int {
	var total fxp.Int
	for _, f := range e.featureMap[strings.ToLower(featureID)] {
		if bonus, ok := f.(*feature.SpellBonus); ok &&
			bonus.NameCriteria.Matches(name) &&
			bonus.TagsCriteria.Matches(tags...) {
			total += bonus.AdjustedAmount()
			bonus.AddToTooltip(tooltip)
		}
	}
	return total
}

// SpellPointComparedBonusFor returns the total bonus for the matching spell point bonuses.
func (e *Entity) SpellPointComparedBonusFor(featureID, qualifier string, tags []string, tooltip *xio.ByteBuffer) fxp.Int {
	var total fxp.Int
	for _, f := range e.featureMap[strings.ToLower(featureID)] {
		if bonus, ok := f.(*feature.SpellPointBonus); ok &&
			bonus.NameCriteria.Matches(qualifier) &&
			bonus.TagsCriteria.Matches(tags...) {
			total += bonus.AdjustedAmount()
			bonus.AddToTooltip(tooltip)
		}
	}
	return total
}

// AddNamedWeaponDamageBonusesFor adds the bonuses for matching weapons that match to the map. If 'm' is nil, it will
// be created. The provided map (or the newly created one) will be returned.
func (e *Entity) AddNamedWeaponDamageBonusesFor(featureID, nameQualifier, usageQualifier string, tagsQualifier []string, dieCount int, tooltip *xio.ByteBuffer, m map[*feature.WeaponDamageBonus]bool) map[*feature.WeaponDamageBonus]bool {
	if m == nil {
		m = make(map[*feature.WeaponDamageBonus]bool)
	}
	for _, one := range e.NamedWeaponDamageBonusesFor(featureID, nameQualifier, usageQualifier, tagsQualifier, dieCount, tooltip) {
		m[one] = true
	}
	return m
}

// NamedWeaponDamageBonusesFor returns the bonuses for matching weapons.
func (e *Entity) NamedWeaponDamageBonusesFor(featureID, nameQualifier, usageQualifier string, tagsQualifier []string, dieCount int, tooltip *xio.ByteBuffer) []*feature.WeaponDamageBonus {
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
			bonus.TagsCriteria.Matches(tagsQualifier...) {
			bonuses = append(bonuses, bonus)
			level := bonus.LeveledAmount.Level
			bonus.LeveledAmount.Level = fxp.From(dieCount)
			bonus.AddToTooltip(tooltip)
			bonus.LeveledAmount.Level = level
		}
	}
	return bonuses
}

// NamedWeaponSkillBonusesFor returns the bonuses for matching weapons.
func (e *Entity) NamedWeaponSkillBonusesFor(featureID, nameQualifier, usageQualifier string, tagsQualifier []string, tooltip *xio.ByteBuffer) []*feature.SkillBonus {
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
			bonus.TagsCriteria.Matches(tagsQualifier...) {
			bonuses = append(bonuses, bonus)
			bonus.AddToTooltip(tooltip)
		}
	}
	return bonuses
}

// Move returns the current Move value for the given Encumbrance.
func (e *Entity) Move(enc datafile.Encumbrance) int {
	initialMove := e.ResolveAttributeCurrent(gid.BasicMove).Max(0)
	divisor := 2 * xmath.Min(CountThresholdOpMet(attribute.HalveMove, e.Attributes), 2)
	if divisor > 0 {
		initialMove = initialMove.Div(fxp.From(divisor)).Ceil()
	}
	move := initialMove.Mul(fxp.Ten + fxp.Two.Mul(enc.Penalty())).Div(fxp.Ten).Trunc()
	if move < fxp.One {
		if initialMove > 0 {
			return 1
		}
		return 0
	}
	return fxp.As[int](move)
}

// Dodge returns the current Dodge value for the given Encumbrance.
func (e *Entity) Dodge(enc datafile.Encumbrance) int {
	dodge := fxp.Three + e.DodgeBonus + e.ResolveAttributeCurrent(gid.BasicSpeed).Max(0)
	divisor := 2 * xmath.Min(CountThresholdOpMet(attribute.HalveDodge, e.Attributes), 2)
	if divisor > 0 {
		dodge = dodge.Div(fxp.From(divisor)).Ceil()
	}
	return fxp.As[int]((dodge + enc.Penalty()).Max(fxp.One))
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
	return measure.Weight(fxp.Int(e.BasicLift()).Mul(encumbrance.WeightMultiplier()))
}

// OneHandedLift returns the one-handed lift value.
func (e *Entity) OneHandedLift() measure.Weight {
	return measure.Weight(fxp.Int(e.BasicLift()).Mul(fxp.Two))
}

// TwoHandedLift returns the two-handed lift value.
func (e *Entity) TwoHandedLift() measure.Weight {
	return measure.Weight(fxp.Int(e.BasicLift()).Mul(fxp.Eight))
}

// ShoveAndKnockOver returns the shove & knock over value.
func (e *Entity) ShoveAndKnockOver() measure.Weight {
	return measure.Weight(fxp.Int(e.BasicLift()).Mul(fxp.Twelve))
}

// RunningShoveAndKnockOver returns the running shove & knock over value.
func (e *Entity) RunningShoveAndKnockOver() measure.Weight {
	return measure.Weight(fxp.Int(e.BasicLift()).Mul(fxp.TwentyFour))
}

// CarryOnBack returns the carry on back value.
func (e *Entity) CarryOnBack() measure.Weight {
	return measure.Weight(fxp.Int(e.BasicLift()).Mul(fxp.Fifteen))
}

// ShiftSlightly returns the shift slightly value.
func (e *Entity) ShiftSlightly() measure.Weight {
	return measure.Weight(fxp.Int(e.BasicLift()).Mul(fxp.Fifty))
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
	var v fxp.Int
	if e.SheetSettings.DamageProgression == attribute.KnowingYourOwnStrength {
		var diff fxp.Int
		if st > fxp.Nineteen {
			diff = st.Div(fxp.Ten).Trunc() - fxp.One
			st -= diff.Mul(fxp.Ten)
		}
		v = fxp.From(math.Pow(10, fxp.As[float64](st)/10)).Mul(fxp.Two)
		if st <= fxp.Six {
			v = v.Mul(fxp.Ten).Round().Div(fxp.Ten)
		} else {
			v = v.Round()
		}
		v = v.Mul(fxp.From(math.Pow(10, fxp.As[float64](diff))))
	} else {
		v = st.Mul(st).Div(fxp.Five)
	}
	if v >= fxp.Ten {
		v = v.Round()
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
		return strconv.Itoa(e.Profile.AdjustedSizeModifier())
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

// ResolveAttributeCurrent resolves the given attribute ID to its current value, or fxp.Min.
func (e *Entity) ResolveAttributeCurrent(attrID string) fxp.Int {
	if e != nil && e.Type == datafile.PC {
		return e.Attributes.Current(attrID)
	}
	return fxp.Min
}

// PreservesUserDesc returns true if the user description widget should be preserved when written to disk. Normally, only
// character sheets should return true for this.
func (e *Entity) PreservesUserDesc() bool {
	return e.Type == datafile.PC
}

// Ancestry returns the current Ancestry.
func (e *Entity) Ancestry() *ancestry.Ancestry {
	var anc *ancestry.Ancestry
	TraverseTraits(func(t *Trait) bool {
		if t.Container() && t.ContainerType == trait.Race {
			if anc = ancestry.Lookup(t.Ancestry, SettingsProvider.Libraries()); anc != nil {
				return true
			}
		}
		return false
	}, true, e.Traits...)
	if anc == nil {
		if anc = ancestry.Lookup(ancestry.Default, SettingsProvider.Libraries()); anc == nil {
			jot.Fatal(1, "unable to load default ancestry (Human)")
		}
	}
	return anc
}

// EquippedWeapons returns a sorted list of equipped weapons.
func (e *Entity) EquippedWeapons(weaponType weapon.Type) []*Weapon {
	m := make(map[uint32]*Weapon)
	TraverseTraits(func(a *Trait) bool {
		if !a.Container() {
			for _, w := range a.Weapons {
				if w.Type == weaponType {
					m[w.HashCode()] = w
				}
			}
		}
		return false
	}, true, e.Traits...)
	TraverseEquipment(func(eqp *Equipment) bool {
		if eqp.Equipped {
			for _, w := range eqp.Weapons {
				if w.Type == weaponType {
					m[w.HashCode()] = w
				}
			}
		}
		return false
	}, e.CarriedEquipment...)
	TraverseSkills(func(s *Skill) bool {
		if !s.Container() {
			for _, w := range s.Weapons {
				if w.Type == weaponType {
					m[w.HashCode()] = w
				}
			}
		}
		return false
	}, e.Skills...)
	TraverseSpells(func(s *Spell) bool {
		if !s.Container() {
			for _, w := range s.Weapons {
				if w.Type == weaponType {
					m[w.HashCode()] = w
				}
			}
		}
		return false
	}, e.Spells...)
	list := make([]*Weapon, 0, len(m))
	for _, v := range m {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Less(list[j]) })
	return list
}

// Reactions returns the current set of reactions.
func (e *Entity) Reactions() []*ConditionalModifier {
	m := make(map[string]*ConditionalModifier)
	TraverseTraits(func(a *Trait) bool {
		source := i18n.Text("from trait ") + a.String()
		if !a.Container() {
			e.reactionsFromFeatureList(source, a.Features, m)
		}
		for _, mod := range a.Modifiers {
			if !mod.Disabled {
				e.reactionsFromFeatureList(source, mod.Features, m)
			}
		}
		if a.CR != trait.None && a.CRAdj == ReactionPenalty {
			amt := fxp.From(ReactionPenalty.Adjustment(a.CR))
			situation := fmt.Sprintf(i18n.Text("from others when %s is triggered"), a.String())
			if r, exists := m[situation]; exists {
				r.Add(source, amt)
			} else {
				m[situation] = NewReaction(source, situation, amt)
			}
		}
		return false
	}, true, e.Traits...)
	TraverseEquipment(func(eqp *Equipment) bool {
		if eqp.Equipped && eqp.Quantity > 0 {
			source := i18n.Text("from equipment ") + eqp.Name
			e.reactionsFromFeatureList(source, eqp.Features, m)
			for _, mod := range eqp.Modifiers {
				if !mod.Disabled {
					e.reactionsFromFeatureList(source, mod.Features, m)
				}
			}
		}
		return false
	}, e.CarriedEquipment...)
	list := make([]*ConditionalModifier, 0, len(m))
	for _, v := range m {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Less(list[j]) })
	return list
}

func (e *Entity) reactionsFromFeatureList(source string, features feature.Features, m map[string]*ConditionalModifier) {
	for _, f := range features {
		if bonus, ok := f.(*feature.ReactionBonus); ok {
			amt := bonus.AdjustedAmount()
			if r, exists := m[bonus.Situation]; exists {
				r.Add(source, amt)
			} else {
				m[bonus.Situation] = NewReaction(source, bonus.Situation, amt)
			}
		}
	}
}

// ConditionalModifiers returns the current set of conditional modifiers.
func (e *Entity) ConditionalModifiers() []*ConditionalModifier {
	m := make(map[string]*ConditionalModifier)
	TraverseTraits(func(a *Trait) bool {
		source := i18n.Text("from trait ") + a.String()
		if !a.Container() {
			e.conditionalModifiersFromFeatureList(source, a.Features, m)
		}
		for _, mod := range a.Modifiers {
			if !mod.Disabled {
				e.conditionalModifiersFromFeatureList(source, mod.Features, m)
			}
		}
		return false
	}, true, e.Traits...)
	TraverseEquipment(func(eqp *Equipment) bool {
		if eqp.Equipped && eqp.Quantity > 0 {
			source := i18n.Text("from equipment ") + eqp.Name
			e.conditionalModifiersFromFeatureList(source, eqp.Features, m)
			for _, mod := range eqp.Modifiers {
				if !mod.Disabled {
					e.conditionalModifiersFromFeatureList(source, mod.Features, m)
				}
			}
		}
		return false
	}, e.CarriedEquipment...)
	list := make([]*ConditionalModifier, 0, len(m))
	for _, v := range m {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Less(list[j]) })
	return list
}

func (e *Entity) conditionalModifiersFromFeatureList(source string, features feature.Features, m map[string]*ConditionalModifier) {
	for _, f := range features {
		if bonus, ok := f.(*feature.ConditionalModifier); ok {
			amt := bonus.AdjustedAmount()
			if r, exists := m[bonus.Situation]; exists {
				r.Add(source, amt)
			} else {
				m[bonus.Situation] = NewReaction(source, bonus.Situation, amt)
			}
		}
	}
}

// TraitList implements ListProvider
func (e *Entity) TraitList() []*Trait {
	return e.Traits
}

// SetTraitList implements ListProvider
func (e *Entity) SetTraitList(list []*Trait) {
	e.Traits = list
}

// CarriedEquipmentList implements ListProvider
func (e *Entity) CarriedEquipmentList() []*Equipment {
	return e.CarriedEquipment
}

// SetCarriedEquipmentList implements ListProvider
func (e *Entity) SetCarriedEquipmentList(list []*Equipment) {
	e.CarriedEquipment = list
}

// OtherEquipmentList implements ListProvider
func (e *Entity) OtherEquipmentList() []*Equipment {
	return e.OtherEquipment
}

// SetOtherEquipmentList implements ListProvider
func (e *Entity) SetOtherEquipmentList(list []*Equipment) {
	e.OtherEquipment = list
}

// SkillList implements ListProvider
func (e *Entity) SkillList() []*Skill {
	return e.Skills
}

// SetSkillList implements ListProvider
func (e *Entity) SetSkillList(list []*Skill) {
	e.Skills = list
}

// SpellList implements ListProvider
func (e *Entity) SpellList() []*Spell {
	return e.Spells
}

// SetSpellList implements ListProvider
func (e *Entity) SetSpellList(list []*Spell) {
	e.Spells = list
}

// NoteList implements ListProvider
func (e *Entity) NoteList() []*Note {
	return e.Notes
}

// SetNoteList implements ListProvider
func (e *Entity) SetNoteList(list []*Note) {
	e.Notes = list
}

// CRC64 computes a CRC-64 value for the canonical disk format of the data. The ModifiedOn field is ignored for this
// calculation.
func (e *Entity) CRC64() uint64 {
	var buffer bytes.Buffer
	saved := e.ModifiedOn
	e.ModifiedOn = jio.Time{}
	defer func() { e.ModifiedOn = saved }()
	if err := jio.Save(context.Background(), &buffer, e); err != nil {
		return 0
	}
	return crc.Bytes(0, buffer.Bytes())
}
