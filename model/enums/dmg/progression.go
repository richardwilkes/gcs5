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

package dmg

import (
	"strings"

	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible Progression values.
const (
	BasicSet Progression = iota
	KnowingYourOwnStrength
	NoSchoolGrognardDamage
	ThrustEqualsSwingMinus2
	SwingEqualsThrustPlus2
	PhoenixFlameD3
)

// Progression controls how Thrust and Swing are calculated.
type Progression uint8

// ProgressionFromString extracts a Progression from a string.
func ProgressionFromString(str string) Progression {
	for p := BasicSet; p <= PhoenixFlameD3; p++ {
		if strings.EqualFold(p.Key(), str) {
			return p
		}
	}
	return BasicSet
}

// Key returns the key used to represent this Progression.
func (p Progression) Key() string {
	switch p {
	case KnowingYourOwnStrength:
		return "knowing_your_own_strength"
	case NoSchoolGrognardDamage:
		return "no_school_grognard_damage"
	case ThrustEqualsSwingMinus2:
		return "thrust_equals_swing_minus_2"
	case SwingEqualsThrustPlus2:
		return "swing_equals_thrust_plus_2"
	case PhoenixFlameD3:
		return "phoenix_flame_d3"
	default: // BasicSet
		return "basic_set"
	}
}

// String implements fmt.Stringer.
func (p Progression) String() string {
	switch p {
	case KnowingYourOwnStrength:
		return i18n.Text("Knowing Your Own Strength")
	case NoSchoolGrognardDamage:
		return i18n.Text("No School Grognard")
	case ThrustEqualsSwingMinus2:
		return i18n.Text("Thrust = Swing-2")
	case SwingEqualsThrustPlus2:
		return i18n.Text("Swing = Thrust+2")
	case PhoenixFlameD3:
		return i18n.Text("PhoenixFlame d3")
	default: // BasicSet
		return i18n.Text("Basic Set")
	}
}

// Tooltip returns the tooltip for the Progression.
func (p Progression) Tooltip() string {
	tooltip := i18n.Text("Determines the method used to calculate thrust and swing damage")
	if footnote := p.Footnote(); footnote != "" {
		return tooltip + ".\n" + footnote
	}
	return tooltip
}

// Footnote returns a footnote for the Progression, if any.
func (p Progression) Footnote() string {
	switch p {
	case KnowingYourOwnStrength:
		return i18n.Text("Pyramid 3-83, pages 16-19")
	case NoSchoolGrognardDamage:
		return i18n.Text("https://noschoolgrognard.blogspot.com/2013/04/adjusting-swing-damage-in-dungeon.html")
	case ThrustEqualsSwingMinus2:
		return i18n.Text("https://github.com/richardwilkes/gcs/issues/97")
	case SwingEqualsThrustPlus2:
		return i18n.Text("Houserule originating with Kevin Smyth. See https://gamingballistic.com/2020/12/04/df-eastmarch-boss-fight-and-house-rules/")
	case PhoenixFlameD3:
		return i18n.Text("Houserule that use d3s instead of d6s for Damage. See: https://github.com/richardwilkes/gcs/pull/393")
	default: // BasicSet
		return ""
	}
}

// Thrust returns the thrust damage for the given strength.
func (p Progression) Thrust(strength int) *dice.Dice {
	switch p {
	case KnowingYourOwnStrength:
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
	case NoSchoolGrognardDamage:
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
	case ThrustEqualsSwingMinus2:
		d := p.Swing(strength)
		d.Modifier -= 2
		return d
	case PhoenixFlameD3:
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
	default: // BasicSet, SwingEqualsThrustPlus2
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
}

// Swing returns the swing damage for the given strength.
func (p Progression) Swing(strength int) *dice.Dice {
	switch p {
	case KnowingYourOwnStrength:
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
	case NoSchoolGrognardDamage:
		return p.Thrust(strength + 3)
	case SwingEqualsThrustPlus2:
		d := p.Thrust(strength)
		d.Modifier += 2
		return d
	case PhoenixFlameD3:
		return p.Thrust(strength)
	default: // BasicSet, ThrustEqualsSwingMinus2
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
}
