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
	"github.com/richardwilkes/gcs/model/gurps/weapon"
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// WeaponDamage holds the damage information for a weapon.
type WeaponDamage struct {
	Owner                     *Weapon               `json:"-"`
	Type                      string                `json:"type"`
	StrengthType              weapon.StrengthDamage `json:"st,omitempty"`
	Base                      *dice.Dice            `json:"base,omitempty"`
	ArmorDivisor              fixed.F64d4           `json:"armor_divisor,omitempty"`
	Fragmentation             *dice.Dice            `json:"fragmentation,omitempty"`
	FragmentationArmorDivisor fixed.F64d4           `json:"fragmentation_armor_divisor,omitempty"`
	FragmentationType         string                `json:"fragmentation_type,omitempty"`
	ModifierPerDie            fixed.F64d4           `json:"modifier_per_die,omitempty"`
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
	// TODO: Implement
	return ""
}
