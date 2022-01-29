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
	"github.com/richardwilkes/gcs/model/encoding"
	"github.com/richardwilkes/gcs/model/f64d4"
	"github.com/richardwilkes/gcs/model/gurps/ancestry"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// Masks for the various TypeBits.
const (
	MentalTypeMask = 1 << iota
	PhysicalTypeMask
	SocialTypeMask
	ExoticTypeMask
	SupernaturalTypeMask
)

const (
	advantageAncestryKey       = "ancestry"
	advantageBasePointsKey     = "base_points"
	advantageContainerTypeKey  = "container_type"
	advantageCRAdjKey          = "cr_adj"
	advantageCRKey             = "cr"
	advantageExoticKey         = "exotic"
	advantageLevelsKey         = "levels"
	advantageMentalKey         = "mental"
	advantagePhysicalKey       = "physical"
	advantagePointsPerLevelKey = "points_per_level"
	advantagePrereqsKey        = "prereqs"
	advantageRoundCostDownKey  = "round_down"
	advantageSocialKey         = "social"
	advantageSupernaturalKey   = "supernatural"
	advantageTypeKey           = "advantage"
	advantageUserDescKey       = "userdesc"
	advantageCalcPointsKey     = "points"
)

// Advantage holds an advantage, disadvantage, quirk, or perk.
type Advantage struct {
	Common
	Parent            *Advantage
	Levels            fixed.F64d4
	BasePoints        fixed.F64d4
	PointsPerLevel    fixed.F64d4
	Features          []*Feature
	Modifiers         []*AdvantageModifier
	Prereq            *Prereq
	Defaults          []*SkillDefault
	Weapons           []*Weapon
	Ancestry          *ancestry.Ancestry
	Children          []*Advantage
	UnsatisfiedReason string
	UserDesc          string
	Categories        []string
	ContainerType     AdvantageContainerType // TODO: Consider merging Container & ContainerType
	TypeBits          uint8
	CR                SelfControlRoll
	CRAdj             SelfControlRollAdj
	Satisfied         bool
	RoundCostDown     bool
	Enabled           bool
}

// NewAdvantageFromJSON creates a new Advantage from a JSON object.
func NewAdvantageFromJSON(parent *Advantage, data map[string]interface{}) *Advantage {
	a := &Advantage{Parent: parent}
	a.Common.FromJSON(advantageTypeKey, data)
	if a.Container {
		a.ContainerType = AdvantageContainerTypeFromKey(encoding.String(data[advantageContainerTypeKey]))
		array := encoding.Array(data[commonChildrenKey])
		if len(array) != 0 {
			a.Children = make([]*Advantage, len(array))
			for i, one := range array {
				a.Children[i] = NewAdvantageFromJSON(a, encoding.Object(one))
			}
		}
	} else {
		a.Enabled = !encoding.Bool(data[commonDisabledKey])
		a.RoundCostDown = encoding.Bool(data[advantageRoundCostDownKey])
		if encoding.Bool(data[advantageMentalKey]) {
			a.TypeBits |= MentalTypeMask
		}
		if encoding.Bool(data[advantagePhysicalKey]) {
			a.TypeBits |= PhysicalTypeMask
		}
		if encoding.Bool(data[advantageSocialKey]) {
			a.TypeBits |= SocialTypeMask
		}
		if encoding.Bool(data[advantageExoticKey]) {
			a.TypeBits |= ExoticTypeMask
		}
		if encoding.Bool(data[advantageSupernaturalKey]) {
			a.TypeBits |= SupernaturalTypeMask
		}
		if v, exists := data[advantageLevelsKey]; exists {
			a.Levels = encoding.Number(v)
		} else {
			a.Levels = f64d4.NegOne
		}
		a.BasePoints = encoding.Number(data[advantageBasePointsKey])
		a.PointsPerLevel = encoding.Number(data[advantagePointsPerLevelKey])
		a.Prereq = NewPrereqFromJSON(encoding.Object(data[advantagePrereqsKey]), nil)
		a.Defaults = SkillDefaultsListFromJSON(data)
		a.Weapons = WeaponsListFromJSON(data)
		a.Features = FeaturesListFromJSON(data)
	}
	a.CR = SelfControlRollFromJSON(advantageCRKey, data)
	a.CRAdj = SelfControlRollAdjFromKey(encoding.String(data[advantageCRAdjKey]))
	a.Modifiers = AdvantageModifiersListFromJSON(commonModifiersKey, data)
	a.UserDesc = encoding.String(data[advantageUserDescKey])
	a.Categories = StringListFromJSON(commonCategoriesKey, true, data)
	return a
}

// ToJSON emits this object as JSON.
func (a *Advantage) ToJSON(encoder *encoding.JSONEncoder, entity *Entity) {
	encoder.StartObject()
	a.Common.ToInlineJSON(advantageTypeKey, encoder)
	if a.Container {
		if a.ContainerType != Group {
			encoder.KeyedString(advantageContainerTypeKey, a.ContainerType.Key(), false, false)
			if a.ContainerType == Race {
				encoder.KeyedString(advantageAncestryKey, a.Ancestry.Name, true, true)
			}
		}
		if len(a.Children) != 0 {
			encoder.Key(commonChildrenKey)
			encoder.StartArray()
			for _, one := range a.Children {
				one.ToJSON(encoder, entity)
			}
			encoder.EndArray()
		}
	} else {
		encoder.KeyedBool(commonDisabledKey, !a.Enabled, true)
		encoder.KeyedBool(advantageRoundCostDownKey, a.RoundCostDown, true)
		encoder.KeyedBool(advantageMentalKey, a.TypeBits&MentalTypeMask != 0, true)
		encoder.KeyedBool(advantagePhysicalKey, a.TypeBits&PhysicalTypeMask != 0, true)
		encoder.KeyedBool(advantageSocialKey, a.TypeBits&SocialTypeMask != 0, true)
		encoder.KeyedBool(advantageExoticKey, a.TypeBits&ExoticTypeMask != 0, true)
		encoder.KeyedBool(advantageSupernaturalKey, a.TypeBits&SupernaturalTypeMask != 0, true)
		if a.Levels >= 0 {
			encoder.KeyedNumber(advantageLevelsKey, a.Levels, false)
		}
		encoder.KeyedNumber(advantageBasePointsKey, a.BasePoints, true)
		encoder.KeyedNumber(advantagePointsPerLevelKey, a.PointsPerLevel, true)
		encoding.ToKeyedJSON(a.Prereq, advantagePrereqsKey, encoder)
		SkillDefaultsListToJSON(a.Defaults, encoder)
		WeaponsListToJSON(a.Weapons, encoder)
		FeaturesListToJSON(a.Features, encoder)
	}
	a.CR.ToKeyedJSON(advantageCRKey, encoder)
	if a.CR != NoneRequired {
		a.CRAdj.ToKeyedJSON(advantageCRAdjKey, encoder)
	}
	AdvantageModifiersListToJSON(commonModifiersKey, a.Modifiers, encoder)
	if entity != nil && entity.PreservesUserDesc() {
		encoder.KeyedString(advantageUserDescKey, a.UserDesc, true, true)
	}
	StringListToJSON(commonCategoriesKey, a.Categories, encoder)
	// Emit the calculated values for third parties
	encoder.Key(calcKey)
	encoder.StartObject()
	encoder.KeyedNumber(advantageCalcPointsKey, a.AdjustedPoints(entity), false)
	encoder.EndObject()
	encoder.EndObject()
}

// AdjustedPoints returns the total points, taking levels and modifiers into account. 'entity' may be nil.
func (a *Advantage) AdjustedPoints(entity *Entity) fixed.F64d4 {
	if !a.Enabled {
		return 0
	}
	if !a.Container {
		return AdjustedPoints(entity, a.BasePoints, a.Levels, a.PointsPerLevel, a.CR, a.AllModifiers(), a.RoundCostDown)
	}
	var points fixed.F64d4
	if a.ContainerType == AlternativeAbilities {
		values := make([]fixed.F64d4, len(a.Children))
		for i, one := range a.Children {
			values[i] = one.AdjustedPoints(entity)
			if values[i] > points {
				points = values[i]
			}
		}
		max := points
		found := false
		for _, v := range values {
			if !found && max == v {
				found = true
			} else {
				points += f64d4.ApplyRounding(calculateModifierPoints(v, f64d4.Twenty), a.RoundCostDown)
			}
		}
	} else {
		for _, one := range a.Children {
			points += one.AdjustedPoints(entity)
		}
	}
	return points
}

// AllModifiers returns the modifiers plus any inherited from parents.
func (a *Advantage) AllModifiers() []*AdvantageModifier {
	all := make([]*AdvantageModifier, len(a.Modifiers))
	copy(all, a.Modifiers)
	p := a.Parent
	for p != nil {
		all = append(all, p.Modifiers...)
		p = p.Parent
	}
	return all
}

// AdjustedPoints returns the total points, taking levels and modifiers into account. 'entity' may be nil.
func AdjustedPoints(entity *Entity, basePoints, levels, pointsPerLevel fixed.F64d4, cr SelfControlRoll, modifiers []*AdvantageModifier, roundCostDown bool) fixed.F64d4 {
	var baseEnh, levelEnh, baseLim, levelLim fixed.F64d4
	multiplier := cr.Multiplier()
	for _, one := range modifiers {
		if one.Enabled {
			modifier := one.CostModifier()
			switch one.CostType {
			case Percentage:
				switch one.Affects {
				case Total:
					if modifier < 0 {
						baseLim += modifier
						levelLim += modifier
					} else {
						baseEnh += modifier
						levelEnh += modifier
					}
				case BaseOnly:
					if modifier < 0 {
						baseLim += modifier
					} else {
						baseEnh += modifier
					}
				case LevelsOnly:
					if modifier < 0 {
						levelLim += modifier
					} else {
						levelEnh += modifier
					}
				}
			case Points:
				if one.Affects == LevelsOnly {
					pointsPerLevel += modifier
				} else {
					basePoints += modifier
				}
			case Multiplier:
				multiplier = multiplier.Mul(modifier)
			}
		}
	}
	modifiedBasePoints := basePoints
	leveledPoints := pointsPerLevel.Mul(levels)
	if baseEnh != 0 || baseLim != 0 || levelEnh != 0 || levelLim != 0 {
		if SheetSettingsFor(entity).UseMultiplicativeModifiers {
			if baseEnh == levelEnh && baseLim == levelLim {
				modifiedBasePoints = modifyPoints(modifyPoints(modifiedBasePoints+leveledPoints, baseEnh), f64d4.Max(baseLim, f64d4.NegEighty))
			} else {
				modifiedBasePoints = modifyPoints(modifyPoints(modifiedBasePoints, baseEnh), f64d4.Max(baseLim, f64d4.NegEighty)) +
					modifyPoints(modifyPoints(leveledPoints, levelEnh), f64d4.Max(levelLim, f64d4.NegEighty))
			}
		} else {
			baseMod := f64d4.Max(baseEnh+baseLim, f64d4.NegEighty)
			levelMod := f64d4.Max(levelEnh+levelLim, f64d4.NegEighty)
			if baseMod == levelMod {
				modifiedBasePoints = modifyPoints(modifiedBasePoints+leveledPoints, baseMod)
			} else {
				modifiedBasePoints = modifyPoints(modifiedBasePoints, baseMod) + modifyPoints(leveledPoints, levelMod)
			}
		}
	} else {
		modifiedBasePoints += leveledPoints
	}
	return f64d4.ApplyRounding(modifiedBasePoints.Mul(multiplier), roundCostDown)
}

func modifyPoints(points, modifier fixed.F64d4) fixed.F64d4 {
	return points + calculateModifierPoints(points, modifier)
}

func calculateModifierPoints(points, modifier fixed.F64d4) fixed.F64d4 {
	return points.Mul(modifier).Div(f64d4.Hundred)
}
