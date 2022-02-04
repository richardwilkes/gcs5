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

package attribute

import (
	"github.com/richardwilkes/rpgtools/dice"
	"github.com/richardwilkes/toolbox/i18n"
)

// Possible DamageProgression values.
const (
	BasicSet                = DamageProgression("basic_set")
	KnowingYourOwnStrength  = DamageProgression("knowing_your_own_strength")
	NoSchoolGrognardDamage  = DamageProgression("no_school_grognard_damage")
	ThrustEqualsSwingMinus2 = DamageProgression("thrust_equals_swing_minus_2")
	SwingEqualsThrustPlus2  = DamageProgression("swing_equals_thrust_plus_2")
	PhoenixFlameD3          = DamageProgression("phoenix_flame_d3")
)

// AllDamageProgressions is the complete set of DamageProgression values.
var AllDamageProgressions = []DamageProgression{
	BasicSet,
	KnowingYourOwnStrength,
	NoSchoolGrognardDamage,
	ThrustEqualsSwingMinus2,
	SwingEqualsThrustPlus2,
	PhoenixFlameD3,
}

// DamageProgression controls how Thrust and Swing are calculated.
type DamageProgression string

// EnsureValid ensures this is of a known value.
func (d DamageProgression) EnsureValid() DamageProgression {
	for _, one := range AllDamageProgressions {
		if one == d {
			return d
		}
	}
	return AllDamageProgressions[0]
}

// String implements fmt.Stringer.
func (d DamageProgression) String() string {
	switch d {
	case BasicSet:
		return i18n.Text("Basic Set")
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
	default:
		return BasicSet.String()
	}
}

// Footnote returns a footnote for the DamageProgression, if any.
func (d DamageProgression) Footnote() string {
	switch d {
	case BasicSet:
		return ""
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
	default:
		return BasicSet.String()
	}
}

// Tooltip returns the tooltip for the DamageProgression.
func (d DamageProgression) Tooltip() string {
	tooltip := i18n.Text("Determines the method used to calculate thrust and swing damage")
	if footnote := d.Footnote(); footnote != "" {
		return tooltip + ".\n" + footnote
	}
	return tooltip
}

// Thrust returns the thrust damage for the given strength.
func (d DamageProgression) Thrust(strength int) *dice.Dice {
	switch d {
	case BasicSet:
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
		thr := BasicSet.Swing(strength)
		thr.Modifier -= 2
		return thr
	case SwingEqualsThrustPlus2:
		return BasicSet.Thrust(strength)
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
	default:
		return BasicSet.Thrust(strength)
	}
}

// Swing returns the swing damage for the given strength.
func (d DamageProgression) Swing(strength int) *dice.Dice {
	switch d {
	case BasicSet:
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
		return NoSchoolGrognardDamage.Thrust(strength + 3)
	case ThrustEqualsSwingMinus2:
		return BasicSet.Swing(strength)
	case SwingEqualsThrustPlus2:
		sw := BasicSet.Thrust(strength)
		sw.Modifier += 2
		return sw
	case PhoenixFlameD3:
		return PhoenixFlameD3.Thrust(strength)
	default:
		return BasicSet.Swing(strength)
	}
}
