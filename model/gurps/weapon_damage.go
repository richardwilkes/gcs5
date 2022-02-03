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

	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
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
		w.ArmorDivisor = f64d4.One
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (w *WeaponDamage) MarshalJSON() ([]byte, error) {
	// An armor divisor of 0 is not valid and 1 is very common, so suppress its output when 1.
	armorDivisor := w.ArmorDivisor
	if armorDivisor == f64d4.One {
		w.ArmorDivisor = 0
	}
	data, err := json.Marshal(&w.WeaponDamageData)
	w.ArmorDivisor = armorDivisor
	return data, err
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
	// TODO: Implement
	return ""
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
	if w.ArmorDivisor != f64d4.One {
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
			if w.FragmentationArmorDivisor != f64d4.One {
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
