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

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
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
	entity        *gurps.Entity
	featureParent fmt.Stringer
	features      *feature.Features
}

func newFeaturesPanel(entity *gurps.Entity, featureParent fmt.Stringer, features *feature.Features) *featuresPanel {
	p := &featuresPanel{
		entity:        entity,
		featureParent: featureParent,
		features:      features,
	}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  2,
		HAlign: unison.FillAlignment,
		HGrab:  true,
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
		if created := p.createFeatureForType(lastFeatureTypeUsed); created != nil {
			*features = slices.Insert(*features, 0, created)
			p.insertFeaturePanel(1, created)
			unison.DockContainerFor(p).MarkForLayoutRecursively()
			widget.MarkModified(p)
		}
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
	p.addLeveledModifierLine(panel, f, &f.LeveledAmount)
	panel.AddChild(unison.NewPanel())
	wrapper := unison.NewPanel()
	var popup *unison.PopupMenu[attribute.BonusLimitation]
	attrChoicePopup := addAttributeChoicePopup(wrapper, p.entity, i18n.Text("to"), &f.Attribute,
		gurps.SizeFlag|gurps.DodgeFlag|gurps.ParryFlag|gurps.BlockFlag)
	callback := attrChoicePopup.SelectionCallback
	attrChoicePopup.SelectionCallback = func(index int, item *gurps.AttributeChoice) {
		callback(index, item)
		adjustPopupBlank(popup, f.Attribute != gid.Strength)
	}
	popup = addPopup(wrapper, attribute.AllBonusLimitation, &f.Limitation)
	adjustPopupBlank(popup, f.Attribute != gid.Strength)
	wrapper.SetLayout(&unison.FlexLayout{
		Columns:  len(wrapper.Children()),
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	wrapper.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	panel.AddChild(wrapper)
	return panel
}

func (p *featuresPanel) createConditionalModifierPanel(f *feature.ConditionalModifier) *unison.Panel {
	panel := p.createBasePanel(f)
	p.addLeveledModifierLine(panel, f, &f.LeveledAmount)
	panel.AddChild(unison.NewPanel())
	watermark := i18n.Text("Triggering Condition")
	field := widget.NewMultiLineStringField(watermark, func() string { return f.Situation },
		func(value string) {
			f.Situation = value
			panel.MarkForLayoutAndRedraw()
			widget.MarkModified(panel)
		})
	field.Watermark = watermark
	field.AutoScroll = false
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	panel.AddChild(field)
	return panel
}

func (p *featuresPanel) createDRBonusPanel(f *feature.DRBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	p.addLeveledModifierLine(panel, f, &f.LeveledAmount)

	panel.AddChild(unison.NewPanel())
	addHitLocationChoicePopup(panel, p.entity, i18n.Text("to the"), &f.Location)

	panel.AddChild(unison.NewPanel())
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{
		Columns:  3,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	wrapper.AddChild(widget.NewFieldLeadingLabel(i18n.Text("against")))
	field := widget.NewStringField(i18n.Text("Specialization"), func() string { return f.Specialization },
		func(value string) {
			f.Specialization = value
			f.Normalize()
			widget.MarkModified(wrapper)
		})
	field.Watermark = gid.All
	field.SetMinimumTextWidthUsing("Specialization")
	wrapper.AddChild(field)
	wrapper.AddChild(widget.NewFieldTrailingLabel(i18n.Text("attacks")))
	panel.AddChild(wrapper)
	return panel
}

func (p *featuresPanel) createReactionBonusPanel(f *feature.ReactionBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	p.addLeveledModifierLine(panel, f, &f.LeveledAmount)
	panel.AddChild(unison.NewPanel())
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createSkillBonusPanel(f *feature.SkillBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	p.addLeveledModifierLine(panel, f, &f.LeveledAmount)
	panel.AddChild(unison.NewPanel())
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createSkillPointBonusPanel(f *feature.SkillPointBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	p.addLeveledModifierLine(panel, f, &f.LeveledAmount)
	panel.AddChild(unison.NewPanel())
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createSpellBonusPanel(f *feature.SpellBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	p.addLeveledModifierLine(panel, f, &f.LeveledAmount)
	panel.AddChild(unison.NewPanel())
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createSpellPointBonusPanel(f *feature.SpellPointBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	p.addLeveledModifierLine(panel, f, &f.LeveledAmount)
	// TODO: Implement
	return panel
}

func (p *featuresPanel) createWeaponDamageBonusPanel(f *feature.WeaponDamageBonus) *unison.Panel {
	panel := p.createBasePanel(f)
	p.addLeveledModifierLine(panel, f, &f.LeveledAmount)
	panel.AddChild(unison.NewPanel())
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

func (p *featuresPanel) addLeveledModifierLine(parent *unison.Panel, f feature.Feature, amount *feature.LeveledAmount) {
	panel := unison.NewPanel()
	p.addTypeSwitcher(panel, f)
	addLeveledAmountPanel(panel, amount)
	panel.SetLayout(&unison.FlexLayout{
		Columns:  len(panel.Children()),
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	panel.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	parent.AddChild(panel)
}

func (p *featuresPanel) featureTypesList() []feature.Type {
	if e, ok := p.featureParent.(*gurps.Equipment); ok && e.Container() {
		return feature.AllType
	}
	return feature.AllWithoutContainedWeightType
}

func (p *featuresPanel) addTypeSwitcher(parent *unison.Panel, f feature.Feature) {
	currentType := f.FeatureType()
	popup := addPopup(parent, p.featureTypesList(), &currentType)
	popup.SelectionCallback = func(_ int, item feature.Type) {
		if newFeature := p.createFeatureForType(item); newFeature != nil {
			lastFeatureTypeUsed = item
			parent.Parent().RemoveFromParent()
			list := *p.features
			i := slices.IndexFunc(list, func(one feature.Feature) bool { return one == f })
			list[i] = newFeature
			p.insertFeaturePanel(i+1, newFeature)
			unison.DockContainerFor(p).MarkForLayoutRecursively()
			widget.MarkModified(p)
		}
	}
}

func (p *featuresPanel) createFeatureForType(featureType feature.Type) feature.Feature {
	switch featureType {
	case feature.AttributeBonusType:
		one := feature.NewAttributeBonus(lastAttributeIDUsed)
		one.Parent = p.featureParent
		return one
	case feature.ConditionalModifierType:
		one := feature.NewConditionalModifierBonus()
		one.Parent = p.featureParent
		return one
	case feature.ContainedWeightReductionType:
		return feature.NewContainedWeightReduction()
	case feature.CostReductionType:
		return feature.NewCostReduction(lastAttributeIDUsed)
	case feature.DRBonusType:
		one := feature.NewDRBonus()
		one.Parent = p.featureParent
		return one
	case feature.ReactionBonusType:
		one := feature.NewReactionBonus()
		one.Parent = p.featureParent
		return one
	case feature.SkillBonusType:
		one := feature.NewSkillBonus()
		one.Parent = p.featureParent
		return one
	case feature.SkillPointBonusType:
		one := feature.NewSkillPointBonus()
		one.Parent = p.featureParent
		return one
	case feature.SpellBonusType:
		one := feature.NewSpellBonus()
		one.Parent = p.featureParent
		return one
	case feature.SpellPointBonusType:
		one := feature.NewSpellPointBonus()
		one.Parent = p.featureParent
		return one
	case feature.WeaponBonusType:
		one := feature.NewWeaponDamageBonus()
		one.Parent = p.featureParent
		return one
	default:
		jot.Warn(errs.Newf("unknown feature type: %s", featureType.Key()))
		return nil
	}
}
