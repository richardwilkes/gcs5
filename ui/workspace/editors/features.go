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

package editors

import (
	"fmt"
	"reflect"

	"github.com/richardwilkes/gcs/model/gurps/feature"
	"github.com/richardwilkes/gcs/model/gurps/gid"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
	"golang.org/x/exp/slices"
)

var (
	lastFeatureTypeUsed = feature.AttributeBonusType
	lastAttributeIDUsed = gid.Strength
)

type featuresPanel struct {
	unison.Panel
	featureParent fmt.Stringer
	features      *feature.Features
}

func newFeaturesPanel(featureParent fmt.Stringer, features *feature.Features) *featuresPanel {
	p := &featuresPanel{
		featureParent: featureParent,
		features:      features,
	}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{Columns: 1})
	p.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  2,
		HAlign: unison.FillAlignment,
	})
	p.SetBorder(unison.NewCompoundBorder(
		&widget.TitledBorder{
			Title: i18n.Text("Features"),
			Font:  unison.LabelFont,
		},
		unison.NewEmptyBorder(unison.NewUniformInsets(2))))
	p.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	}
	addButton := unison.NewSVGButton(res.CircledAddSVG)
	addButton.ClickCallback = func() {
		var created feature.Feature
		switch lastFeatureTypeUsed {
		case feature.AttributeBonusType:
			one := feature.NewAttributeBonus(lastAttributeIDUsed)
			one.Parent = featureParent
			created = one
		case feature.ConditionalModifierType:
			one := feature.NewConditionalModifierBonus()
			one.Parent = featureParent
			created = one
		case feature.ContainedWeightReductionType:
			created = feature.NewContainedWeightReduction()
		case feature.CostReductionType:
			created = feature.NewCostReduction(lastAttributeIDUsed)
		case feature.DRBonusType:
			one := feature.NewDRBonus()
			one.Parent = featureParent
			created = one
		case feature.ReactionBonusType:
			one := feature.NewReactionBonus()
			one.Parent = featureParent
			created = one
		case feature.SkillBonusType:
			one := feature.NewSkillBonus()
			one.Parent = featureParent
			created = one
		case feature.SkillPointBonusType:
			one := feature.NewSkillPointBonus()
			one.Parent = featureParent
			created = one
		case feature.SpellBonusType:
			one := feature.NewSpellBonus()
			one.Parent = featureParent
			created = one
		case feature.SpellPointBonusType:
			one := feature.NewSpellPointBonus()
			one.Parent = featureParent
			created = one
		case feature.WeaponBonusType:
			one := feature.NewWeaponDamageBonus()
			one.Parent = featureParent
			created = one
		default:
			jot.Warn(errs.Newf("unknown feature type: %s", lastFeatureTypeUsed.String()))
			return
		}
		*features = slices.Insert(*features, 0, created)
		p.insertFeaturePanel(1, created)
		unison.DockContainerFor(p).MarkForLayoutRecursively()
		widget.MarkModified(p)
	}
	p.AddChild(addButton)
	for i, one := range *features {
		p.insertFeaturePanel(i+1, one)
	}
	return p
}

func (p *featuresPanel) insertFeaturePanel(index int, f feature.Feature) {
	var panel *unison.Panel
	switch one := f.(type) {
	case *feature.AttributeBonus:
		panel = p.createAttributeBonusPanel(one)
	case *feature.ConditionalModifier:
		panel = p.createConditionalModifierPanel(one)
	case *feature.ContainedWeightReduction:
		panel = p.createContainedWeightReductionPanel(one)
	case *feature.CostReduction:
		panel = p.createCostReductionPanel(one)
	case *feature.DRBonus:
		panel = p.createDRBonusPanel(one)
	case *feature.ReactionBonus:
		panel = p.createReactionBonusPanel(one)
	case *feature.SkillBonus:
		panel = p.createSkillBonusPanel(one)
	case *feature.SkillPointBonus:
		panel = p.createSkillPointBonusPanel(one)
	case *feature.SpellBonus:
		panel = p.createSpellBonusPanel(one)
	case *feature.SpellPointBonus:
		panel = p.createSpellPointBonusPanel(one)
	case *feature.WeaponDamageBonus:
		panel = p.createWeaponDamageBonusPanel(one)
	default:
		jot.Warn(errs.Newf("unknown feature type: %s", reflect.TypeOf(f).String()))
		return
	}
	if panel != nil {
		panel.SetLayoutData(&unison.FlexLayoutData{
			HAlign: unison.FillAlignment,
			HGrab:  true,
		})
		p.AddChildAtIndex(panel, index)
	}
}

func (p *featuresPanel) createBasePanel(f feature.Feature) *unison.Panel {
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HAlign:   unison.FillAlignment,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	deleteButton := unison.NewSVGButton(res.TrashSVG)
	deleteButton.ClickCallback = func() {
		if i := slices.IndexFunc(*p.features, func(elem feature.Feature) bool { return elem == f }); i != -1 {
			*p.features = slices.Delete(*p.features, i, i+1)
		}
		panel.RemoveFromParent()
		unison.DockContainerFor(p).MarkForLayoutRecursively()
		widget.MarkModified(p)
	}
	panel.AddChild(deleteButton)
	return panel
}

func (p *featuresPanel) createAttributeBonusPanel(f *feature.AttributeBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createConditionalModifierPanel(f *feature.ConditionalModifier) *unison.Panel {
	panel := p.createBasePanel(f)
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createContainedWeightReductionPanel(f *feature.ContainedWeightReduction) *unison.Panel {
	panel := p.createBasePanel(f)
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createCostReductionPanel(f *feature.CostReduction) *unison.Panel {
	panel := p.createBasePanel(f)
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createDRBonusPanel(f *feature.DRBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createReactionBonusPanel(f *feature.ReactionBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createSkillBonusPanel(f *feature.SkillBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createSkillPointBonusPanel(f *feature.SkillPointBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createSpellBonusPanel(f *feature.SpellBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createSpellPointBonusPanel(f *feature.SpellPointBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createWeaponDamageBonusPanel(f *feature.WeaponDamageBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	// TODO: Implement
	return panel
}
