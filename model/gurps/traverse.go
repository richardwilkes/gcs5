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

// TraverseTraits calls the function 'f' for each Trait and its children in the input list. Return true from the
// function to abort early.
func TraverseTraits(f func(*Trait) bool, onlyEnabled bool, in ...*Trait) {
	type trackingInfo struct {
		list  []*Trait
		index int
	}
	tracking := []*trackingInfo{
		{
			list:  in,
			index: 0,
		},
	}
	for {
		if len(tracking) == 0 {
			return
		}
		current := tracking[len(tracking)-1]
		if current.index >= len(current.list) {
			tracking = tracking[:len(tracking)-1]
		} else {
			one := current.list[current.index]
			current.index++
			if onlyEnabled && !one.Enabled() {
				continue
			}
			if f(one) {
				return
			}
			if one.Container() && len(one.Children) != 0 {
				tracking = append(tracking, &trackingInfo{list: one.Children})
			}
		}
	}
}

// TraverseTraitModifiers calls the function 'f' for each TraitModifier and its children in the input list.
// Return true from the function to abort early.
func TraverseTraitModifiers(f func(*TraitModifier) bool, onlyEnabled bool, in ...*TraitModifier) {
	type trackingInfo struct {
		list  []*TraitModifier
		index int
	}
	tracking := []*trackingInfo{
		{
			list:  in,
			index: 0,
		},
	}
	for {
		if len(tracking) == 0 {
			return
		}
		current := tracking[len(tracking)-1]
		if current.index >= len(current.list) {
			tracking = tracking[:len(tracking)-1]
		} else {
			one := current.list[current.index]
			current.index++
			if onlyEnabled && one.Disabled {
				continue
			}
			if f(one) {
				return
			}
			if one.Container() && len(one.Children) != 0 {
				tracking = append(tracking, &trackingInfo{list: one.Children})
			}
		}
	}
}

// TraverseEquipment calls the function 'f' for each Equipment and its children in the input list. Return true from the
// function to abort early.
func TraverseEquipment(f func(*Equipment) bool, in ...*Equipment) {
	type trackingInfo struct {
		list  []*Equipment
		index int
	}
	tracking := []*trackingInfo{
		{
			list:  in,
			index: 0,
		},
	}
	for {
		if len(tracking) == 0 {
			return
		}
		current := tracking[len(tracking)-1]
		if current.index >= len(current.list) {
			tracking = tracking[:len(tracking)-1]
		} else {
			one := current.list[current.index]
			if f(one) {
				return
			}
			current.index++
			if one.Container() && len(one.Children) != 0 {
				tracking = append(tracking, &trackingInfo{list: one.Children})
			}
		}
	}
}

// TraverseEquipmentModifiers calls the function 'f' for each EquipmentModifier and its children in the input list.
// Return true from the function to abort early.
func TraverseEquipmentModifiers(f func(*EquipmentModifier) bool, onlyEnabled bool, in ...*EquipmentModifier) {
	type trackingInfo struct {
		list  []*EquipmentModifier
		index int
	}
	tracking := []*trackingInfo{
		{
			list:  in,
			index: 0,
		},
	}
	for {
		if len(tracking) == 0 {
			return
		}
		current := tracking[len(tracking)-1]
		if current.index >= len(current.list) {
			tracking = tracking[:len(tracking)-1]
		} else {
			one := current.list[current.index]
			current.index++
			if onlyEnabled && one.Disabled {
				continue
			}
			if f(one) {
				return
			}
			if one.Container() && len(one.Children) != 0 {
				tracking = append(tracking, &trackingInfo{list: one.Children})
			}
		}
	}
}

// TraverseSkills calls the function 'f' for each Skill and its children in the input list. Return true from the function
// to abort early.
func TraverseSkills(f func(*Skill) bool, in ...*Skill) {
	type trackingInfo struct {
		list  []*Skill
		index int
	}
	tracking := []*trackingInfo{
		{
			list:  in,
			index: 0,
		},
	}
	for {
		if len(tracking) == 0 {
			return
		}
		current := tracking[len(tracking)-1]
		if current.index >= len(current.list) {
			tracking = tracking[:len(tracking)-1]
		} else {
			one := current.list[current.index]
			if f(one) {
				return
			}
			current.index++
			if one.Container() && len(one.Children) != 0 {
				tracking = append(tracking, &trackingInfo{list: one.Children})
			}
		}
	}
}

// TraverseSpells calls the function 'f' for each Spell and its children in the input list. Return true from the function
// to abort early.
func TraverseSpells(f func(*Spell) bool, in ...*Spell) {
	type trackingInfo struct {
		list  []*Spell
		index int
	}
	tracking := []*trackingInfo{
		{
			list:  in,
			index: 0,
		},
	}
	for {
		if len(tracking) == 0 {
			return
		}
		current := tracking[len(tracking)-1]
		if current.index >= len(current.list) {
			tracking = tracking[:len(tracking)-1]
		} else {
			one := current.list[current.index]
			if f(one) {
				return
			}
			current.index++
			if one.Container() && len(one.Children) != 0 {
				tracking = append(tracking, &trackingInfo{list: one.Children})
			}
		}
	}
}

// TraverseNotes calls the function 'f' for each Note and its children in the input list. Return true from the function
// to abort early.
func TraverseNotes(f func(*Note) bool, in ...*Note) {
	type trackingInfo struct {
		list  []*Note
		index int
	}
	tracking := []*trackingInfo{
		{
			list:  in,
			index: 0,
		},
	}
	for {
		if len(tracking) == 0 {
			return
		}
		current := tracking[len(tracking)-1]
		if current.index >= len(current.list) {
			tracking = tracking[:len(tracking)-1]
		} else {
			one := current.list[current.index]
			if f(one) {
				return
			}
			current.index++
			if one.Container() && len(one.Children) != 0 {
				tracking = append(tracking, &trackingInfo{list: one.Children})
			}
		}
	}
}
