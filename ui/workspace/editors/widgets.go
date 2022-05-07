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

	"github.com/richardwilkes/gcs/model/fxp"
	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/gurps/skill"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/tbl"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

func addNameLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndStringField(parent, i18n.Text("Name"), "", fieldData)
}

func addSpecializationLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndStringField(parent, i18n.Text("Specialization"), "", fieldData)
}

func addPageRefLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndStringField(parent, i18n.Text("Page Reference"), tbl.PageRefTooltipText, fieldData)
}

func addNotesLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndMultiLineStringField(parent, i18n.Text("Notes"), "", fieldData)
}

func addVTTNotesLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndMultiLineStringField(parent, i18n.Text("VTT Notes"),
		i18n.Text("Any notes for VTT use; see the instructions for your VVT to determine if/how these can be used"),
		fieldData)
}

func addUserDescLabelAndField(parent *unison.Panel, fieldData *string) {
	addLabelAndMultiLineStringField(parent, i18n.Text("User Description"),
		i18n.Text("Additional notes for your own reference. These only exist in character sheets and will be removed if transferred to a data list or template"),
		fieldData)
}

func addTechLevelRequired(parent *unison.Panel, fieldData **string, includeField bool) {
	tl := i18n.Text("Tech Level")
	var field *widget.StringField
	if includeField {
		wrapper := addFlowWrapper(parent, tl, 2)
		field = widget.NewStringField(tl, func() string {
			if *fieldData == nil {
				return ""
			}
			return **fieldData
		}, func(value string) {
			if *fieldData == nil {
				return
			}
			**fieldData = value
			widget.MarkModified(parent)
		})
		if *fieldData == nil {
			field.SetEnabled(false)
		}
		insets := field.Border().Insets()
		field.MinimumTextWidth = field.Font.SimpleWidth("12^") + insets.Width()
		wrapper.AddChild(field)
		parent = wrapper
	} else {
		parent.AddChild(widget.NewFieldLeadingLabel(tl))
	}
	last := *fieldData
	required := last != nil
	parent.AddChild(widget.NewCheckBox(i18n.Text("Required"), required, func(b bool) {
		required = b
		if b {
			if last == nil {
				var data string
				last = &data
			}
			*fieldData = last
			if field != nil {
				field.SetEnabled(true)
			}
		} else {
			last = *fieldData
			*fieldData = nil
			if field != nil {
				field.SetEnabled(false)
			}
		}
		widget.MarkModified(parent)
	}))
}

func addDifficultyLabelAndFields(parent *unison.Panel, entity *gurps.Entity, difficulty *gurps.AttributeDifficulty) {
	wrapper := addFlowWrapper(parent, i18n.Text("Difficulty"), 3)
	current := -1
	choices := gurps.AttributeChoices(entity, false)
	for i, one := range choices {
		if one.Key == difficulty.Attribute {
			current = i
			break
		}
	}
	if current == -1 {
		current = len(choices)
		choices = append(choices, &gurps.AttributeChoice{
			Key:   difficulty.Attribute,
			Title: difficulty.Attribute,
		})
	}
	attrChoice := choices[current]
	attrChoicePopup := addPopup(wrapper, choices, &attrChoice)
	attrChoicePopup.SelectionCallback = func(_ int, item *gurps.AttributeChoice) {
		difficulty.Attribute = item.Key
		widget.MarkModified(parent)
	}
	wrapper.AddChild(widget.NewFieldTrailingLabel("/"))
	addPopup(wrapper, skill.AllDifficulty, &difficulty.Difficulty)
}

func addTagsLabelAndField(parent *unison.Panel, fieldData *[]string) {
	addLabelAndListField(parent, i18n.Text("Tags"), i18n.Text("tags"), fieldData)
}

func addLabelAndListField(parent *unison.Panel, labelText, pluralForTooltip string, fieldData *[]string) {
	tooltip := fmt.Sprintf(i18n.Text("Separate multiple %s with commas"), pluralForTooltip)
	label := widget.NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	field := widget.NewMultiLineStringField(labelText, func() string { return gurps.CombineTags(*fieldData) },
		func(value string) {
			*fieldData = gurps.ExtractTags(value)
			parent.MarkForLayoutAndRedraw()
			widget.MarkModified(parent)
		})
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	field.AutoScroll = false
	parent.AddChild(field)
}

func addLabelAndStringField(parent *unison.Panel, labelText, tooltip string, fieldData *string) {
	label := widget.NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	addStringField(parent, labelText, tooltip, fieldData)
}

func addStringField(parent *unison.Panel, labelText, tooltip string, fieldData *string) *widget.StringField {
	field := widget.NewStringField(labelText, func() string { return *fieldData },
		func(value string) {
			*fieldData = value
			widget.MarkModified(parent)
		})
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
	return field
}

func addLabelAndMultiLineStringField(parent *unison.Panel, labelText, tooltip string, fieldData *string) {
	label := widget.NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	field := widget.NewMultiLineStringField(labelText, func() string { return *fieldData },
		func(value string) {
			*fieldData = value
			parent.MarkForLayoutAndRedraw()
			widget.MarkModified(parent)
		})
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	field.AutoScroll = false
	parent.AddChild(field)
}

func addIntegerField(parent *unison.Panel, labelText, tooltip string, fieldData *int, min, max int) *widget.IntegerField {
	field := widget.NewIntegerField(labelText, func() int { return *fieldData },
		func(value int) {
			*fieldData = value
			widget.MarkModified(parent)
		}, min, max, false)
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
	return field
}

func addLabelAndNumericField(parent *unison.Panel, labelText, tooltip string, fieldData *fxp.Int, min, max fxp.Int) *widget.NumericField {
	label := widget.NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	field := widget.NewNumericField(labelText, func() fxp.Int { return *fieldData },
		func(value fxp.Int) {
			*fieldData = value
			widget.MarkModified(parent)
		}, min, max, false)
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
	return field
}

func addNumericField(parent *unison.Panel, labelText, tooltip string, fieldData *fxp.Int, min, max fxp.Int) *widget.NumericField {
	field := widget.NewNumericField(labelText, func() fxp.Int { return *fieldData },
		func(value fxp.Int) {
			*fieldData = value
			widget.MarkModified(parent)
		}, min, max, false)
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
	return field
}

func addWeightField(parent *unison.Panel, labelText, tooltip string, entity *gurps.Entity, fieldData *measure.Weight) *widget.WeightField {
	field := widget.NewWeightField(labelText, entity, func() measure.Weight { return *fieldData },
		func(value measure.Weight) {
			*fieldData = value
			widget.MarkModified(parent)
		}, 0, measure.Weight(fxp.Max))
	if tooltip != "" {
		field.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(field)
	return field
}

func addCheckBox(parent *unison.Panel, labelText string, fieldData *bool) {
	parent.AddChild(widget.NewCheckBox(labelText, *fieldData, func(b bool) {
		*fieldData = b
		widget.MarkModified(parent)
	}))
}

func addInvertedCheckBox(parent *unison.Panel, labelText string, fieldData *bool) {
	parent.AddChild(widget.NewCheckBox(labelText, !*fieldData, func(b bool) {
		*fieldData = !b
		widget.MarkModified(parent)
	}))
}

func addFlowWrapper(parent *unison.Panel, labelText string, count int) *unison.Panel {
	parent.AddChild(widget.NewFieldLeadingLabel(labelText))
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{
		Columns:  count,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
		VAlign:   unison.MiddleAlignment,
	})
	parent.AddChild(wrapper)
	return wrapper
}

func addLabelAndPopup[T comparable](parent *unison.Panel, labelText, tooltip string, choices []T, fieldData *T) *unison.PopupMenu[T] {
	label := widget.NewFieldLeadingLabel(labelText)
	if tooltip != "" {
		label.Tooltip = unison.NewTooltipWithText(tooltip)
	}
	parent.AddChild(label)
	return addPopup[T](parent, choices, fieldData)
}

func addPopup[T comparable](parent *unison.Panel, choices []T, fieldData *T) *unison.PopupMenu[T] {
	popup := unison.NewPopupMenu[T]()
	for _, one := range choices {
		popup.AddItem(one)
	}
	popup.Select(*fieldData)
	popup.SelectionCallback = func(_ int, item T) {
		*fieldData = item
		widget.MarkModified(parent)
	}
	parent.AddChild(popup)
	return popup
}

func disableAndBlankField(field unison.Paneler) {
	panel := field.AsPanel()
	panel.SetEnabled(false)
	panel.DrawOverCallback = func(gc *unison.Canvas, rect unison.Rect) {
		rect = panel.ContentRect(false)
		var ink unison.Ink
		if f, ok := panel.Self.(*unison.Field); ok {
			ink = f.BackgroundInk
		} else {
			ink = unison.DefaultFieldTheme.BackgroundInk
		}
		gc.DrawRect(rect, ink.Paint(gc, rect, unison.Fill))
	}
}

func enableAndUnblankField(field unison.Paneler) {
	panel := field.AsPanel()
	panel.SetEnabled(true)
	panel.DrawOverCallback = nil
}

func disableAndBlankPopup[T comparable](popup *unison.PopupMenu[T]) {
	popup.SetEnabled(false)
	popup.DrawOverCallback = func(gc *unison.Canvas, rect unison.Rect) {
		rect = popup.ContentRect(false)
		unison.DrawRoundedRectBase(gc, rect, popup.CornerRadius, 1, popup.BackgroundInk, popup.EdgeInk)
	}
}

func enableAndUnblankPopup[T comparable](popup *unison.PopupMenu[T]) {
	popup.SetEnabled(true)
	popup.DrawOverCallback = nil
}
