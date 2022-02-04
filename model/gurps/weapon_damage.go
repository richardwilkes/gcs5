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

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// WeaponDamageData holds the WeaponDamage data that is written to disk.
type WeaponDamageData struct {
	Type                      string                `json:"type"`
	StrengthType              weapon.StrengthDamage `json:"st,omitempty"`
	Base                      *dice.Dice            `json:"base,omitempty"`
	ArmorDivisor              fixed.F64d4           `json:"armor_divisor,omitempty"`
	Fragmentation             *dice.Dice            `json:"fragmentation,omitempty"`
	FragmentationArmorDivisor fixed.F64d4           `json:"fragmentation_armor_divisor,omitempty"`
	FragmentationType         string                `json:"fragmentation_type,omitempty"`
	ModifierPerDie            fixed.F64d4           `json:"modifier_per_die,omitempty"`
}

// WeaponDamage holds the damage information for a weapon.
type WeaponDamage struct {
	WeaponDamageData
	Owner *Weapon
}

// UnmarshalJSON implements json.Unmarshaler.
func (w *WeaponDamage) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &w.WeaponDamageData); err != nil {
		return err
	}
	if w.ArmorDivisor == 0 {
		w.ArmorDivisor = fxp.One
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (w *WeaponDamage) MarshalJSON() ([]byte, error) {
	// An armor divisor of 0 is not valid and 1 is very common, so suppress its output when 1.
	armorDivisor := w.ArmorDivisor
	if armorDivisor == fxp.One {
		w.ArmorDivisor = 0
	}
	data, err := json.Marshal(&w.WeaponDamageData)
	w.ArmorDivisor = armorDivisor
	return data, err
}

func (w *WeaponDamage) String() string {
	var buffer strings.Builder
	if w.StrengthType != weapon.None {
		buffer.WriteString(w.StrengthType.String())
	}
	convertMods := false
	if w.Owner != nil {
		convertMods = SheetSettingsFor(w.Owner.Entity()).UseModifyingDicePlusAdds
	}
	if w.Base != nil {
		if base := w.Base.StringExtra(convertMods); base != "0" {
			if buffer.Len() != 0 && base[0] != '+' && base[0] != '-' {
				buffer.WriteByte('+')
			}
			buffer.WriteString(base)
		}
	}
	if w.ArmorDivisor != fxp.One {
		buffer.WriteByte('(')
		buffer.WriteString(w.ArmorDivisor.String())
		buffer.WriteByte(')')
	}
	if w.ModifierPerDie != 0 {
		if buffer.Len() != 0 {
			buffer.WriteByte(' ')
		}
		buffer.WriteByte('(')
		buffer.WriteString(w.ModifierPerDie.StringWithSign())
		buffer.WriteString(i18n.Text(" per die)"))
	}
	if t := strings.TrimSpace(w.Type); t != "" {
		buffer.WriteByte(' ')
		buffer.WriteString(t)
	}
	if w.Fragmentation != nil {
		if frag := w.Fragmentation.StringExtra(convertMods); frag != "0" {
			buffer.WriteString(" [")
			buffer.WriteString(frag)
			if w.FragmentationArmorDivisor != fxp.One {
				buffer.WriteByte('(')
				buffer.WriteString(w.FragmentationArmorDivisor.String())
				buffer.WriteByte(')')
			}
			buffer.WriteByte(' ')
			buffer.WriteString(w.FragmentationType)
			buffer.WriteByte(']')
		}
	}
	return strings.TrimSpace(buffer.String())
}

// DamageTooltip returns a formatted tooltip for the damage.
func (w *WeaponDamage) DamageTooltip() string {
	var tooltip xio.ByteBuffer
	w.ResolvedDamage(&tooltip)
	if tooltip.Len() == 0 {
		return i18n.Text("No additional modifiers")
	}
	return i18n.Text("Includes modifiers from") + tooltip.String()
}

// ResolvedDamage returns the damage, fully resolved for the user's sw or thr, if possible.
func (w *WeaponDamage) ResolvedDamage(tooltip *xio.ByteBuffer) string {
	if w.Owner == nil {
		return w.String()
	}
	pc := w.Owner.PC()
	if pc == nil {
		return w.String()
	}
	maxST := w.Owner.ResolvedMinimumStrength().Mul(fxp.Three)
	st := pc.StrengthOrZero() + pc.StrikingStrengthBonus
	if maxST > 0 && maxST < st {
		st = maxST
	}
	base := &dice.Dice{
		Sides:      6,
		Multiplier: 1,
	}
	if w.Base != nil {
		*base = *w.Base
	}
	adq, adqOK := w.Owner.Owner.(*Advantage)
	if adqOK && adq.IsLeveled() {
		multiplyDice(int(adq.Levels.AsInt64()), base)
	}
	intST := int(st.AsInt64())
	switch w.StrengthType {
	case weapon.Thrust:
		base = addDice(base, pc.ThrustFor(intST))
	case weapon.LeveledThrust:
		thrust := pc.ThrustFor(intST)
		if adqOK && adq.IsLeveled() {
			multiplyDice(int(adq.Levels.AsInt64()), thrust)
		}
		base = addDice(base, thrust)
	case weapon.Swing:
		base = addDice(base, pc.SwingFor(intST))
	case weapon.LeveledSwing:
		thrust := pc.SwingFor(intST)
		if adqOK && adq.IsLeveled() {
			multiplyDice(int(adq.Levels.AsInt64()), thrust)
		}
		base = addDice(base, thrust)
	}
	var bestDefault *SkillDefault
	best := fixed.F64d4Min
	for _, one := range w.Owner.Defaults {
		if one.SkillBased() {
			if level := one.SkillLevelFast(pc, false, nil, true); best < level {
				best = level
				bestDefault = one
			}
		}
	}
	bonusSet := make(map[*feature.WeaponDamageBonus]bool)
	categories := w.Owner.Owner.CategoryList()
	if bestDefault != nil {
		pc.AddWeaponComparedDamageBonusesFor(feature.SkillNameID+"*", bestDefault.Name, bestDefault.Specialization,
			categories, base.Count, tooltip, bonusSet)
		pc.AddWeaponComparedDamageBonusesFor(feature.SkillNameID+"/"+bestDefault.Name, bestDefault.Name,
			bestDefault.Specialization, categories, base.Count, tooltip, bonusSet)
	}
	nameQualifier := w.Owner.String()
	pc.AddNamedWeaponDamageBonusesFor(feature.WeaponNamedIDPrefix+"*", nameQualifier, w.Owner.Usage, categories,
		base.Count, tooltip, bonusSet)
	pc.AddNamedWeaponDamageBonusesFor(feature.WeaponNamedIDPrefix+"/"+nameQualifier, nameQualifier, w.Owner.Usage,
		categories, base.Count, tooltip, bonusSet)
	for _, f := range w.Owner.Owner.FeatureList() {
		w.extractWeaponDamageBonus(f, bonusSet, base.Count, tooltip)
	}
	if adqOK {
		for _, mod := range adq.Modifiers {
			if !mod.Disabled {
				for _, f := range mod.Features {
					w.extractWeaponDamageBonus(f, bonusSet, base.Count, tooltip)
				}
			}
		}
	}
	if eqp, ok := w.Owner.Owner.(*Equipment); ok {
		for _, mod := range eqp.Modifiers {
			if !mod.Disabled {
				for _, f := range mod.Features {
					w.extractWeaponDamageBonus(f, bonusSet, base.Count, tooltip)
				}
			}
		}
	}
	adjustForPhoenixFlame := pc.SheetSettings.DamageProgression == attribute.PhoenixFlameD3 && base.Sides == 3
	var percent fixed.F64d4
	for bonus := range bonusSet {
		if bonus.Percent {
			percent += bonus.Amount
		} else {
			amt := bonus.Amount
			if bonus.PerLevel {
				amt = amt.Mul(fixed.F64d4FromInt64(int64(base.Count)))
				if adjustForPhoenixFlame {
					amt = amt.Div(fxp.Two)
				}
			}
			base.Modifier += int(amt.AsInt64())
		}
	}
	if w.ModifierPerDie != 0 {
		amt := w.ModifierPerDie.Mul(fixed.F64d4FromInt64(int64(base.Count)))
		if adjustForPhoenixFlame {
			amt = amt.Div(fxp.Two)
		}
		base.Modifier += int(amt.AsInt64())
	}
	if percent != 0 {
		base = adjustDiceForPercentBonus(base, percent)
	}
	var buffer strings.Builder
	if base.Count != 0 || base.Modifier != 0 {
		buffer.WriteString(base.StringExtra(pc.SheetSettings.UseModifyingDicePlusAdds))
	}
	if w.ArmorDivisor != fxp.One {
		buffer.WriteByte('(')
		buffer.WriteString(w.ArmorDivisor.String())
		buffer.WriteByte(')')
	}
	if strings.TrimSpace(w.Type) != "" {
		if buffer.Len() != 0 {
			buffer.WriteByte(' ')
		}
		buffer.WriteString(w.Type)
	}
	if w.Fragmentation != nil {
		if frag := w.Fragmentation.StringExtra(pc.SheetSettings.UseModifyingDicePlusAdds); frag != "0" {
			if buffer.Len() != 0 {
				buffer.WriteByte(' ')
			}
			buffer.WriteByte('[')
			buffer.WriteString(frag)
			if w.FragmentationArmorDivisor != 1 {
				buffer.WriteByte('(')
				buffer.WriteString(w.FragmentationArmorDivisor.String())
				buffer.WriteByte(')')
			}
			buffer.WriteByte(' ')
			buffer.WriteString(w.FragmentationType)
			buffer.WriteByte(']')
		}
	}
	return buffer.String()
}

func (w *WeaponDamage) extractWeaponDamageBonus(f feature.Feature, set map[*feature.WeaponDamageBonus]bool, dieCount int, tooltip *xio.ByteBuffer) {
	if bonus, ok := f.(*feature.WeaponDamageBonus); ok {
		level := bonus.LeveledAmount.Level
		bonus.LeveledAmount.Level = fixed.F64d4FromInt64(int64(dieCount))
		switch bonus.SelectionType {
		case weapon.WithRequiredSkill:
		case weapon.ThisWeapon:
			if bonus.SpecializationCriteria.Matches(w.Owner.Usage) {
				if _, exists := set[bonus]; !exists {
					set[bonus] = true
					bonus.AddToTooltip(tooltip)
				}
			}
		case weapon.WithName:
			if bonus.NameCriteria.Matches(w.Owner.String()) && bonus.SpecializationCriteria.Matches(w.Owner.Usage) &&
				bonus.CategoryCriteria.Matches(w.Owner.Owner.CategoryList()...) {
				if _, exists := set[bonus]; !exists {
					set[bonus] = true
					bonus.AddToTooltip(tooltip)
				}
			}
		default:
			jot.Fatal(1, "unknown selection type: "+string(bonus.SelectionType))
		}
		bonus.LeveledAmount.Level = level
	}
}

func multiplyDice(multiplier int, d *dice.Dice) {
	d.Count *= multiplier
	d.Modifier *= multiplier
	if d.Multiplier != 1 {
		d.Multiplier *= multiplier
	}
}

func addDice(left, right *dice.Dice) *dice.Dice {
	if left.Sides > 1 && right.Sides > 1 && left.Sides != right.Sides {
		sides := xmath.MinInt(left.Sides, right.Sides)
		average := fixed.F64d4FromInt64(int64(sides + 1)).Div(fxp.Two)
		averageLeft := fixed.F64d4FromInt64(int64(left.Count * (left.Sides + 1))).Div(fxp.Two).Mul(fixed.F64d4FromInt64(int64(left.Multiplier)))
		averageRight := fixed.F64d4FromInt64(int64(right.Count * (right.Sides + 1))).Div(fxp.Two).Mul(fixed.F64d4FromInt64(int64(right.Multiplier)))
		averageBoth := averageLeft + averageRight
		return &dice.Dice{
			Count:      int(averageBoth.Div(average).AsInt64()),
			Sides:      sides,
			Modifier:   int(fxp.Round(fxp.Mod(averageBoth, average))) + left.Modifier + right.Modifier,
			Multiplier: 1,
		}
	}
	return &dice.Dice{
		Count:      left.Count + right.Count,
		Sides:      xmath.MaxInt(left.Sides, right.Sides),
		Modifier:   left.Modifier + right.Modifier,
		Multiplier: left.Multiplier + right.Multiplier - 1,
	}
}

func adjustDiceForPercentBonus(d *dice.Dice, percent fixed.F64d4) *dice.Dice {
	count := fixed.F64d4FromInt64(int64(d.Count))
	modifier := fixed.F64d4FromInt64(int64(d.Modifier))
	averagePerDie := fixed.F64d4FromInt64(int64(d.Sides + 1)).Div(fxp.Two)
	average := averagePerDie.Mul(count) + modifier
	modifier = modifier.Mul(fxp.Hundred + percent).Div(fxp.Hundred)
	if average < 0 {
		count = count.Mul(fxp.Hundred + percent).Div(fxp.Hundred).Max(0)
	} else {
		average = average.Mul(fxp.Hundred+percent).Div(fxp.Hundred) - modifier
		count = average.Div(averagePerDie).Trunc().Max(0)
		modifier += fxp.Round(average - count.Mul(averagePerDie))
	}
	return &dice.Dice{
		Count:      int(count.AsInt64()),
		Sides:      d.Sides,
		Modifier:   int(modifier.AsInt64()),
		Multiplier: d.Multiplier,
	}
}
