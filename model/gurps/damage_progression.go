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

	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible DamageProgression values.
const (
	BasicSet DamageProgression = iota
	KnowingYourOwnStrength
	NoSchoolGrognardDamage
	ThrustEqualsSwingMinus2
	SwingEqualsThrustPlus2
	PhoenixFlameD3
)

type damageProgressionData struct {
	Key      string
	String   string
	Footnote string
	Thrust   func(strength int) *dice.Dice
	Swing    func(strength int) *dice.Dice
}

// DamageProgression controls how Thrust and Swing are calculated.
type DamageProgression uint8

var damageProgressionValues = []*damageProgressionData{
	{
		Key:    "basic_set",
		String: i18n.Text("Basic Set"),
		Thrust: basicSetThrust,
		Swing:  basicSetSwing,
	},
	{
		Key:      "knowing_your_own_strength",
		String:   i18n.Text("Knowing Your Own Strength"),
		Footnote: i18n.Text("Pyramid 3-83, pages 16-19"),
		Thrust:   kyosThrust,
		Swing:    kyosSwing,
	},
	{
		Key:      "no_school_grognard_damage",
		String:   i18n.Text("No School Grognard"),
		Footnote: i18n.Text("https://noschoolgrognard.blogspot.com/2013/04/adjusting-swing-damage-in-dungeon.html"),
		Thrust:   noSchoolGrognardThrust,
		Swing:    noSchoolGrognardSwing,
	},
	{
		Key:      "thrust_equals_swing_minus_2",
		String:   i18n.Text("Thrust = Swing-2"),
		Footnote: i18n.Text("https://github.com/richardwilkes/gcs/issues/97"),
		Thrust:   basicSetSwingMinus2,
		Swing:    basicSetSwing,
	},
	{
		Key:      "swing_equals_thrust_plus_2",
		String:   i18n.Text("Swing = Thrust+2"),
		Footnote: i18n.Text("Houserule originating with Kevin Smyth. See https://gamingballistic.com/2020/12/04/df-eastmarch-boss-fight-and-house-rules/"),
		Thrust:   basicSetThrust,
		Swing:    basicSetThrustPlus2,
	},
	{
		Key:      "phoenix_flame_d3",
		String:   i18n.Text("PhoenixFlame d3"),
		Footnote: i18n.Text("Houserule that use d3s instead of d6s for Damage. See: https://github.com/richardwilkes/gcs/pull/393"),
		Thrust:   phoenixFlameD3,
		Swing:    phoenixFlameD3,
	},
}

// DamageProgressionFromString extracts a DamageProgression from a key.
func DamageProgressionFromString(key string) DamageProgression {
	for i, one := range damageProgressionValues {
		if strings.EqualFold(key, one.Key) {
			return DamageProgression(i)
		}
	}
	return 0
}

// EnsureValid returns the first DamageProgression if this DamageProgression is not a known value.
func (d DamageProgression) EnsureValid() DamageProgression {
	if int(d) < len(damageProgressionValues) {
		return d
	}
	return 0
}

// Key returns the key used to represent this DamageProgression.
func (d DamageProgression) Key() string {
	return damageProgressionValues[d.EnsureValid()].Key
}

// String implements fmt.Stringer.
func (d DamageProgression) String() string {
	return damageProgressionValues[d.EnsureValid()].String
}

// Tooltip returns the tooltip for the DamageProgression.
func (d DamageProgression) Tooltip() string {
	tooltip := i18n.Text("Determines the method used to calculate thrust and swing damage")
	if footnote := d.Footnote(); footnote != "" {
		return tooltip + ".\n" + footnote
	}
	return tooltip
}

// Footnote returns a footnote for the DamageProgression, if any.
func (d DamageProgression) Footnote() string {
	return damageProgressionValues[d.EnsureValid()].Footnote
}

// Thrust returns the thrust damage for the given strength.
func (d DamageProgression) Thrust(strength int) *dice.Dice {
	return damageProgressionValues[d.EnsureValid()].Thrust(strength)
}

// Swing returns the swing damage for the given strength.
func (d DamageProgression) Swing(strength int) *dice.Dice {
	return damageProgressionValues[d.EnsureValid()].Swing(strength)
}

func basicSetThrust(strength int) *dice.Dice {
	if strength < 19 {
		return &dice.Dice{
			Count:      1,
			Sides:      6,
			Modifier:   -(6 - (strength-1)/2),
			Multiplier: 1,
		}
	}
	value := strength - 11
	if strength > 50 {
		value--
		if strength > 79 {
			value -= 1 + (strength-80)/5
		}
	}
	return &dice.Dice{
		Count:      value/8 + 1,
		Sides:      6,
		Modifier:   value%8/2 - 1,
		Multiplier: 1,
	}
}

func basicSetSwing(strength int) *dice.Dice {
	if strength < 10 {
		return &dice.Dice{
			Count:      1,
			Sides:      6,
			Modifier:   -(5 - (strength-1)/2),
			Multiplier: 1,
		}
	}
	if strength < 28 {
		strength -= 9
		return &dice.Dice{
			Count:      strength/4 + 1,
			Sides:      6,
			Modifier:   strength%4 - 1,
			Multiplier: 1,
		}
	}
	value := strength
	if strength > 40 {
		value -= (strength - 40) / 5
	}
	if strength > 59 {
		value++
	}
	value += 9
	return &dice.Dice{
		Count:      value/8 + 1,
		Sides:      6,
		Modifier:   value%8/2 - 1,
		Multiplier: 1,
	}
}

func basicSetThrustPlus2(strength int) *dice.Dice {
	d := basicSetThrust(strength)
	d.Modifier += 2
	return d
}

func basicSetSwingMinus2(strength int) *dice.Dice {
	d := basicSetSwing(strength)
	d.Modifier -= 2
	return d
}

func kyosThrust(strength int) *dice.Dice {
	if strength < 12 {
		return &dice.Dice{
			Count:      1,
			Sides:      6,
			Modifier:   strength - 12,
			Multiplier: 1,
		}
	}
	return &dice.Dice{
		Count:      (strength - 7) / 4,
		Sides:      6,
		Modifier:   (strength+1)%4 - 1,
		Multiplier: 1,
	}
}

func kyosSwing(strength int) *dice.Dice {
	if strength < 10 {
		return &dice.Dice{
			Count:      1,
			Sides:      6,
			Modifier:   strength - 10,
			Multiplier: 1,
		}
	}
	return &dice.Dice{
		Count:      (strength - 5) / 4,
		Sides:      6,
		Modifier:   (strength-1)%4 - 1,
		Multiplier: 1,
	}
}

func noSchoolGrognardThrust(strength int) *dice.Dice {
	if strength < 11 {
		return &dice.Dice{
			Count:      1,
			Sides:      6,
			Modifier:   -(14 - strength) / 2,
			Multiplier: 1,
		}
	}
	strength -= 11
	return &dice.Dice{
		Count:      strength/8 + 1,
		Sides:      6,
		Modifier:   (strength%8)/2 - 1,
		Multiplier: 1,
	}
}

func noSchoolGrognardSwing(strength int) *dice.Dice {
	return noSchoolGrognardThrust(strength + 3)
}

func phoenixFlameD3(strength int) *dice.Dice {
	if strength < 7 {
		if strength < 1 {
			strength = 1
		}
		return &dice.Dice{
			Count:      1,
			Sides:      6,
			Modifier:   ((strength + 1) / 2) - 7,
			Multiplier: 1,
		}
	}
	if strength < 10 {
		return &dice.Dice{
			Count:      1,
			Sides:      3,
			Modifier:   ((strength + 1) / 2) - 5,
			Multiplier: 1,
		}
	}
	strength -= 8
	return &dice.Dice{
		Count:      strength / 2,
		Sides:      3,
		Modifier:   strength % 2,
		Multiplier: 1,
	}
}
