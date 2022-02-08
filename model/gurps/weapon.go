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
	"bufio"
	"fmt"
	"strings"

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/datafile"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// WeaponOwner defines the methods required of a Weapon owner.
type WeaponOwner interface {
	fmt.Stringer
	OwningEntity() *Entity
	Description() string
	Notes() string
	FeatureList() feature.Features
	CategoryList() []string
}

// WeaponData holds the Weapon data that is written to disk.
type WeaponData struct {
	Type            weapon.Type     `json:"type"`
	Damage          *WeaponDamage   `json:"damage,omitempty"`
	MinimumStrength string          `json:"strength,omitempty"`
	Usage           string          `json:"usage,omitempty"`
	UsageNotes      string          `json:"usage_notes,omitempty"`
	Reach           string          `json:"reach,omitempty"`
	Parry           string          `json:"parry,omitempty"`
	Block           string          `json:"block,omitempty"`
	Accuracy        string          `json:"accuracy,omitempty"`
	Range           string          `json:"range,omitempty"`
	RateOfFire      string          `json:"rate_of_fire,omitempty"`
	Shots           string          `json:"shots,omitempty"`
	Bulk            string          `json:"bulk,omitempty"`
	Recoil          string          `json:"recoil,omitempty"`
	Defaults        []*SkillDefault `json:"defaults,omitempty"`
}

// Weapon holds the stats for a weapon.
type Weapon struct {
	WeaponData
	Owner WeaponOwner
}

// MarshalJSON implements json.Marshaler.
func (w *Weapon) MarshalJSON() ([]byte, error) {
	type calc struct {
		Level  fixed.F64d4 `json:"level,omitempty"`
		Parry  string      `json:"parry,omitempty"`
		Block  string      `json:"block,omitempty"`
		Range  string      `json:"range,omitempty"`
		Damage string      `json:"damage,omitempty"`
	}
	data := struct {
		WeaponData
		Calc calc `json:"calc"`
	}{
		WeaponData: w.WeaponData,
		Calc: calc{
			Level:  w.SkillLevel(nil).Max(0),
			Damage: w.Damage.ResolvedDamage(nil),
		},
	}
	if w.Type == weapon.Melee {
		data.Calc.Parry = w.ResolvedParry(nil)
		data.Calc.Block = w.ResolvedBlock(nil)
	} else {
		data.Calc.Range = w.ResolvedRange()
	}
	return json.Marshal(&data)
}

func (w *Weapon) String() string {
	if w.Owner == nil {
		return ""
	}
	return w.Owner.Description()
}

// SetOwner sets the owner and ensures sub-components have their owners set.
func (w *Weapon) SetOwner(owner WeaponOwner) {
	w.Owner = owner
	w.Damage.Owner = w
}

// Entity returns the owning entity, if any.
func (w *Weapon) Entity() *Entity {
	if w.Owner == nil {
		return nil
	}
	entity := w.Owner.OwningEntity()
	if entity == nil {
		return nil
	}
	return entity
}

// PC returns the owning PC, if any.
func (w *Weapon) PC() *Entity {
	if entity := w.Entity(); entity != nil && entity.Type == datafile.PC {
		return entity
	}
	return nil
}

// SkillLevel returns the resolved skill level.
func (w *Weapon) SkillLevel(tooltip *xio.ByteBuffer) fixed.F64d4 {
	pc := w.PC()
	if pc == nil {
		return 0
	}
	var primaryTooltip *xio.ByteBuffer
	if tooltip != nil {
		primaryTooltip = &xio.ByteBuffer{}
	}
	adj := w.skillLevelBaseAdjustment(pc, primaryTooltip) + w.skillLevelPostAdjustment(pc, primaryTooltip)
	best := fixed.F64d4Min
	for _, def := range w.Defaults {
		if level := def.SkillLevelFast(pc, false, nil, true); level != fixed.F64d4Min {
			level += adj
			if best < level {
				best = level
			}
		}
	}
	if best == fixed.F64d4Min {
		return 0
	}
	if tooltip != nil && primaryTooltip.Len() != 0 {
		if tooltip.Len() != 0 {
			tooltip.WriteByte('\n')
		}
		tooltip.WriteString(primaryTooltip.String())
	}
	if best < 0 {
		best = 0
	}
	return best
}

func (w *Weapon) skillLevelBaseAdjustment(entity *Entity, tooltip *xio.ByteBuffer) fixed.F64d4 {
	var adj fixed.F64d4
	if minST := w.ResolvedMinimumStrength() - (entity.StrengthOrZero() + entity.StrikingStrengthBonus); minST > 0 {
		adj -= minST
	}
	nameQualifier := w.String()
	for _, bonus := range entity.NamedWeaponSkillBonusesFor(feature.WeaponNamedIDPrefix+"*", nameQualifier, w.Usage,
		w.Owner.CategoryList(), tooltip) {
		adj += bonus.AdjustedAmount()
	}
	for _, bonus := range entity.NamedWeaponSkillBonusesFor(feature.WeaponNamedIDPrefix+"/"+nameQualifier,
		nameQualifier, w.Usage, w.Owner.CategoryList(), tooltip) {
		adj += bonus.AdjustedAmount()
	}
	for _, f := range w.Owner.FeatureList() {
		adj += w.extractSkillBonus(f, tooltip)
	}
	if adq, ok := w.Owner.(*Advantage); ok {
		for _, mod := range adq.Modifiers {
			if !mod.Disabled {
				for _, f := range mod.Features {
					adj += w.extractSkillBonus(f, tooltip)
				}
			}
		}
	}
	if eqp, ok := w.Owner.(*Equipment); ok {
		for _, mod := range eqp.Modifiers {
			if !mod.Disabled {
				for _, f := range mod.Features {
					adj += w.extractSkillBonus(f, tooltip)
				}
			}
		}
	}
	return adj
}

func (w *Weapon) skillLevelPostAdjustment(entity *Entity, tooltip *xio.ByteBuffer) fixed.F64d4 {
	if w.Type.EnsureValid() == weapon.Melee && strings.Contains(w.Parry, "F") {
		return w.EncumbrancePenalty(entity, tooltip)
	}
	return 0
}

// EncumbrancePenalty returns the current encumbrance penalty.
func (w *Weapon) EncumbrancePenalty(entity *Entity, tooltip *xio.ByteBuffer) fixed.F64d4 {
	if entity == nil {
		return 0
	}
	penalty := entity.EncumbranceLevel(true).Penalty()
	if penalty != 0 && tooltip != nil {
		tooltip.WriteByte('\n')
		tooltip.WriteString(i18n.Text("Encumbrance"))
		tooltip.WriteString(" [")
		tooltip.WriteString(penalty.StringWithSign())
		tooltip.WriteByte(']')
	}
	return penalty
}

func (w *Weapon) extractSkillBonus(f feature.Feature, tooltip *xio.ByteBuffer) fixed.F64d4 {
	if sb, ok := f.(*feature.SkillBonus); ok {
		switch sb.SelectionType.EnsureValid() {
		case skill.SkillsWithName:
		case skill.ThisWeapon:
			if sb.SpecializationCriteria.Matches(w.Usage) {
				sb.AddToTooltip(tooltip)
				return sb.AdjustedAmount()
			}
		case skill.WeaponsWithName:
			if w.Owner != nil && sb.NameCriteria.Matches(w.Owner.String()) &&
				sb.SpecializationCriteria.Matches(w.Usage) && sb.CategoryCriteria.Matches(w.Owner.CategoryList()...) {
				sb.AddToTooltip(tooltip)
				return sb.AdjustedAmount()
			}
		default:
			jot.Fatal(1, "unhandled selection type: "+string(sb.SelectionType))
		}
	}
	return 0
}

// ResolvedParry returns the resolved parry level.
func (w *Weapon) ResolvedParry(tooltip *xio.ByteBuffer) string {
	return w.resolvedValue(w.Parry, gid.Parry, tooltip)
}

// ResolvedBlock returns the resolved block level.
func (w *Weapon) ResolvedBlock(tooltip *xio.ByteBuffer) string {
	return w.resolvedValue(w.Block, gid.Block, tooltip)
}

// ResolvedRange returns the range, fully resolved for the user's ST, if possible.
func (w *Weapon) ResolvedRange() string {
	//nolint:ifshort // No, pc isn't just used on the next line...
	pc := w.PC()
	if pc == nil {
		return w.Range
	}
	st := (pc.StrengthOrZero() + pc.ThrowingStrengthBonus).Trunc()
	var savedRange string
	calcRange := w.Range
	for calcRange != savedRange {
		calcRange = w.resolveRange(calcRange, st)
		savedRange = calcRange
	}
	return calcRange
}

func (w *Weapon) resolvedValue(input, baseDefaultType string, tooltip *xio.ByteBuffer) string {
	pc := w.PC()
	if pc == nil {
		return input
	}
	var buffer strings.Builder
	skillLevel := fixed.F64d4Max
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if buffer.Len() != 0 {
			buffer.WriteByte('\n')
		}
		if line != "" {
			max := len(line)
			i := 0
			for i < max && line[i] == ' ' {
				i++
			}
			if i < max {
				ch := line[i]
				neg := false
				modifier := 0
				found := false
				if ch == '-' || ch == '+' {
					neg = ch == '-'
					i++
					if i < max {
						ch = line[i]
					}
				}
				for i < max && ch >= '0' && ch <= '9' {
					found = true
					modifier *= 10
					modifier += int(ch - '0')
					i++
					if i < max {
						ch = line[i]
					}
				}
				if found {
					if skillLevel == fixed.F64d4Max {
						var primaryTooltip, secondaryTooltip *xio.ByteBuffer
						if tooltip != nil {
							primaryTooltip = &xio.ByteBuffer{}
						}
						preAdj := w.skillLevelBaseAdjustment(pc, primaryTooltip)
						postAdj := w.skillLevelPostAdjustment(pc, primaryTooltip)
						adj := fxp.Three
						if baseDefaultType == gid.Parry {
							adj += pc.ParryBonus
						} else {
							adj += pc.BlockBonus
						}
						best := fixed.F64d4Min
						for _, def := range w.Defaults {
							level := def.SkillLevelFast(pc, false, nil, true)
							if level == fixed.F64d4Min {
								continue
							}
							level += preAdj
							if baseDefaultType != def.Type() {
								level = (level.Div(fxp.Two) + adj).Trunc()
							}
							level += postAdj
							var possibleTooltip *xio.ByteBuffer
							if def.Type() == gid.Skill && def.Name == "Karate" {
								if tooltip != nil {
									possibleTooltip = &xio.ByteBuffer{}
								}
								level += w.EncumbrancePenalty(pc, possibleTooltip)
							}
							if best < level {
								best = level
								secondaryTooltip = possibleTooltip
							}
						}
						if best != fixed.F64d4Min && tooltip != nil {
							if primaryTooltip.Len() != 0 {
								if tooltip.Len() != 0 {
									tooltip.WriteByte('\n')
								}
								tooltip.WriteString(primaryTooltip.String())
							}
							if secondaryTooltip != nil && secondaryTooltip.Len() != 0 {
								if tooltip.Len() != 0 {
									tooltip.WriteByte('\n')
								}
								tooltip.WriteString(secondaryTooltip.String())
							}
						}
						skillLevel = best.Max(0)
					}
					if neg {
						modifier = -modifier
					}
					num := (skillLevel + fixed.F64d4FromInt(modifier)).Trunc().String()
					if i < max {
						buffer.WriteString(num)
						line = line[i:]
					} else {
						line = num
					}
				}
			}
		}
		buffer.WriteString(line)
	}
	return buffer.String()
}

func (w *Weapon) resolveRange(inRange string, st fixed.F64d4) string {
	where := strings.IndexByte(inRange, 'x')
	if where == -1 {
		return inRange
	}
	last := where + 1
	max := len(inRange)
	if last < max && inRange[last] == ' ' {
		last++
	}
	if last >= max {
		return inRange
	}
	ch := inRange[last]
	found := false
	decimal := false
	started := last
	for (!decimal && ch == '.') || (ch >= '0' && ch <= '9') {
		found = true
		if ch == '.' {
			decimal = true
		}
		last++
		if last >= max {
			break
		}
		ch = inRange[last]
	}
	if !found {
		return inRange
	}
	value, err := fixed.F64d4FromString(inRange[started:last])
	if err != nil {
		return inRange
	}
	var buffer strings.Builder
	if where > 0 {
		buffer.WriteString(inRange[:where])
	}
	buffer.WriteString(value.Mul(st).Trunc().String())
	if last < max {
		buffer.WriteString(inRange[last:])
	}
	return inRange
}

// ResolvedMinimumStrength returns the resolved minimum strength required to use this weapon, or 0 if there is none.
func (w *Weapon) ResolvedMinimumStrength() fixed.F64d4 {
	started := false
	value := 0
	for _, ch := range w.MinimumStrength {
		if ch >= '0' && ch <= '9' {
			value *= 10
			value += int(ch - '0')
			started = true
		} else if started {
			break
		}
	}
	return fixed.F64d4FromInt(value)
}

// FillWithNameableKeys adds any nameable keys found in this Weapon to the provided map.
func (w *Weapon) FillWithNameableKeys(m map[string]string) {
	for _, one := range w.Defaults {
		one.FillWithNameableKeys(m)
	}
}

// ApplyNameableKeys replaces any nameable keys found in this Weapon with the corresponding values in the provided map.
func (w *Weapon) ApplyNameableKeys(m map[string]string) {
	for _, one := range w.Defaults {
		one.ApplyNameableKeys(m)
	}
}
