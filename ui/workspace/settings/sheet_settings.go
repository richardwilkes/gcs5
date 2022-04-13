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

package settings

import (
	"io/fs"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/attribute"
	"github.com/richardwilkes/gcs/model/gurps/measure"
	"github.com/richardwilkes/gcs/model/paper"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/model/settings/display"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type sheetSettingsDockable struct {
	Dockable
	entity                             *gurps.Entity
	damageProgressionPopup             *unison.PopupMenu[attribute.DamageProgression]
	showAdvantageModifier              *unison.CheckBox
	showEquipmentModifier              *unison.CheckBox
	showSpellAdjustments               *unison.CheckBox
	showTitleInsteadOfNameInPageFooter *unison.CheckBox
	useMultiplicativeModifiers         *unison.CheckBox
	useModifyDicePlusAdds              *unison.CheckBox
	lengthUnitsPopup                   *unison.PopupMenu[measure.LengthUnits]
	weightUnitsPopup                   *unison.PopupMenu[measure.WeightUnits]
	userDescDisplayPopup               *unison.PopupMenu[display.Option]
	modifiersDisplayPopup              *unison.PopupMenu[display.Option]
	notesDisplayPopup                  *unison.PopupMenu[display.Option]
	skillLevelAdjDisplayPopup          *unison.PopupMenu[display.Option]
	paperSizePopup                     *unison.PopupMenu[paper.Size]
	orientationPopup                   *unison.PopupMenu[paper.Orientation]
}

// ShowSheetSettings the Sheet Settings window.
func ShowSheetSettings(entity *gurps.Entity) {
	ws, dc, found := Activate(func(d unison.Dockable) bool {
		if s, ok := d.(*sheetSettingsDockable); ok && s.entity == entity {
			return true
		}
		return false
	})
	if !found && ws != nil {
		d := &sheetSettingsDockable{entity: entity}
		d.Self = d
		if entity != nil {
			d.TabTitle = i18n.Text("Sheet Settings: " + entity.Profile.Name)
		} else {
			d.TabTitle = i18n.Text("Default Sheet Settings")
		}
		d.Extension = ".sheet"
		d.Loader = d.load
		d.Saver = d.save
		d.Resetter = d.reset
		d.Setup(ws, dc, nil, nil, d.initContent)
	}
}

func (d *sheetSettingsDockable) CloseWithEntity(entity *gurps.Entity) bool {
	return entity != nil && entity == d.entity
}

func (d *sheetSettingsDockable) settings() *gurps.SheetSettings {
	if d.entity != nil {
		return d.entity.SheetSettings
	}
	return settings.Global().Sheet
}

func (d *sheetSettingsDockable) initContent(content *unison.Panel) {
	content.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.DefaultLabelTheme.Font.LineHeight(),
	})
	d.createDamageProgression(content)
	d.createOptions(content)
	d.createUnitsOfMeasurement(content)
	d.createWhereToDisplay(content)
	d.createPageSettings(content)
}

func (d *sheetSettingsDockable) createDamageProgression(content *unison.Panel) {
	s := d.settings()
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	d.damageProgressionPopup = createSettingPopup(d, panel, i18n.Text("Damage Progression"),
		attribute.AllDamageProgression, s.DamageProgression,
		func(item attribute.DamageProgression) {
			d.damageProgressionPopup.Tooltip = unison.NewTooltipWithText(item.Tooltip())
			d.settings().DamageProgression = item
		})
	content.AddChild(panel)
}

func (d *sheetSettingsDockable) createOptions(content *unison.Panel) {
	s := d.settings()
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	d.showAdvantageModifier = d.addCheckBox(panel, i18n.Text("Show advantage modifier cost adjustments"),
		s.ShowAdvantageModifierAdj, func() {
			d.settings().ShowAdvantageModifierAdj = d.showAdvantageModifier.State == unison.OnCheckState
			d.syncSheet()
		})
	d.showEquipmentModifier = d.addCheckBox(panel, i18n.Text("Show equipment modifier cost & weight adjustments"),
		s.ShowEquipmentModifierAdj, func() {
			d.settings().ShowEquipmentModifierAdj = d.showEquipmentModifier.State == unison.OnCheckState
			d.syncSheet()
		})
	d.showSpellAdjustments = d.addCheckBox(panel, i18n.Text("Show spell ritual, cost & time adjustments"),
		s.ShowSpellAdj, func() {
			d.settings().ShowSpellAdj = d.showSpellAdjustments.State == unison.OnCheckState
			d.syncSheet()
		})
	d.showTitleInsteadOfNameInPageFooter = d.addCheckBox(panel,
		i18n.Text("Show the title instead of the name in the footer"), s.UseTitleInFooter, func() {
			d.settings().UseTitleInFooter = d.showTitleInsteadOfNameInPageFooter.State == unison.OnCheckState
			d.syncSheet()
		})
	d.useMultiplicativeModifiers = d.addCheckBox(panel,
		i18n.Text("Use Multiplicative Modifiers (PW102; changes point value)"), s.UseMultiplicativeModifiers, func() {
			d.settings().UseMultiplicativeModifiers = d.useMultiplicativeModifiers.State == unison.OnCheckState
			d.syncSheet()
		})
	d.useModifyDicePlusAdds = d.addCheckBox(panel, i18n.Text("Use Modifying Dice + Adds (B269)"),
		s.UseModifyingDicePlusAdds, func() {
			d.settings().UseModifyingDicePlusAdds = d.useModifyDicePlusAdds.State == unison.OnCheckState
			d.syncSheet()
		})
	content.AddChild(panel)
}

func (d *sheetSettingsDockable) addCheckBox(panel *unison.Panel, title string, checked bool, onClick func()) *unison.CheckBox {
	checkbox := unison.NewCheckBox()
	checkbox.Text = title
	checkbox.State = unison.CheckStateFromBool(checked)
	checkbox.ClickCallback = onClick
	panel.AddChild(checkbox)
	return checkbox
}

func (d *sheetSettingsDockable) createUnitsOfMeasurement(content *unison.Panel) {
	s := d.settings()
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	panel.SetLayoutData(&unison.FlexLayoutData{HAlign: unison.FillAlignment})
	d.createHeader(panel, i18n.Text("Units of Measurement"))
	d.lengthUnitsPopup = createSettingPopup(d, panel, i18n.Text("Length Units"), measure.AllLengthUnits,
		s.DefaultLengthUnits, func(item measure.LengthUnits) { d.settings().DefaultLengthUnits = item })
	d.weightUnitsPopup = createSettingPopup(d, panel, i18n.Text("Length Units"), measure.AllWeightUnits,
		s.DefaultWeightUnits, func(item measure.WeightUnits) { d.settings().DefaultWeightUnits = item })
	content.AddChild(panel)
}

func (d *sheetSettingsDockable) createWhereToDisplay(content *unison.Panel) {
	s := d.settings()
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	panel.SetLayoutData(&unison.FlexLayoutData{HAlign: unison.FillAlignment})
	d.createHeader(panel, i18n.Text("Where to display…"))
	d.userDescDisplayPopup = createSettingPopup(d, panel, i18n.Text("User Description"), display.AllOption,
		s.UserDescriptionDisplay, func(option display.Option) { d.settings().UserDescriptionDisplay = option })
	d.modifiersDisplayPopup = createSettingPopup(d, panel, i18n.Text("Modifiers"), display.AllOption,
		s.ModifiersDisplay, func(option display.Option) { d.settings().ModifiersDisplay = option })
	d.notesDisplayPopup = createSettingPopup(d, panel, i18n.Text("Notes"), display.AllOption, s.NotesDisplay,
		func(option display.Option) { d.settings().NotesDisplay = option })
	d.skillLevelAdjDisplayPopup = createSettingPopup(d, panel, i18n.Text("Skill Level Adjustments"), display.AllOption,
		s.SkillLevelAdjDisplay, func(option display.Option) { d.settings().SkillLevelAdjDisplay = option })
	content.AddChild(panel)
}

func (d *sheetSettingsDockable) createPageSettings(content *unison.Panel) {
	s := d.settings()
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlexLayout{
		Columns:  4,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	panel.SetLayoutData(&unison.FlexLayoutData{HAlign: unison.FillAlignment})
	d.createHeader(panel, i18n.Text("Page Settings"))
	d.paperSizePopup = createSettingPopup(d, panel, i18n.Text("Paper Size"), paper.AllSize,
		s.Page.Size, func(option paper.Size) { d.settings().Page.Size = option })
	d.orientationPopup = createSettingPopup(d, panel, i18n.Text("Orientation"), paper.AllOrientation,
		s.Page.Orientation, func(option paper.Orientation) { d.settings().Page.Orientation = option })
	d.createPaperMarginField(panel, i18n.Text("Top Margin"), s.Page.TopMargin,
		func(value paper.Length) { d.settings().Page.TopMargin = value })
	d.createPaperMarginField(panel, i18n.Text("Bottom Margin"), s.Page.BottomMargin,
		func(value paper.Length) { d.settings().Page.BottomMargin = value })
	d.createPaperMarginField(panel, i18n.Text("Left Margin"), s.Page.LeftMargin,
		func(value paper.Length) { d.settings().Page.LeftMargin = value })
	d.createPaperMarginField(panel, i18n.Text("Right Margin"), s.Page.RightMargin,
		func(value paper.Length) { d.settings().Page.RightMargin = value })
	content.AddChild(panel)
}

func (d *sheetSettingsDockable) createPaperMarginField(panel *unison.Panel, title string, current paper.Length,
	set func(value paper.Length)) *unison.Field {
	panel.AddChild(widget.NewFieldLeadingLabel(title))
	field := unison.NewField()
	field.SetText(current.String())
	field.ValidateCallback = func() bool {
		_, err := paper.ParseLengthFromString(field.Text())
		return err == nil
	}
	field.ModifiedCallback = func() {
		if value, err := paper.ParseLengthFromString(field.Text()); err == nil {
			set(value)
			d.syncSheet()
		}
	}
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	panel.AddChild(field)
	return field
}

func createSettingPopup[T comparable](d *sheetSettingsDockable, panel *unison.Panel, title string, choices []T,
	current T, set func(option T)) *unison.PopupMenu[T] {
	panel.AddChild(widget.NewFieldLeadingLabel(title))
	popup := unison.NewPopupMenu[T]()
	for _, one := range choices {
		popup.AddItem(one)
	}
	popup.Select(current)
	popup.SelectionCallback = func(_ int, item T) {
		set(item)
		d.syncSheet()
	}
	panel.AddChild(popup)
	return popup
}

func (d *sheetSettingsDockable) createHeader(panel *unison.Panel, title string) {
	label := unison.NewLabel()
	label.Text = title
	desc := label.Font.Descriptor()
	desc.Weight = unison.BoldFontWeight
	label.Font = desc.Font()
	label.SetLayoutData(&unison.FlexLayoutData{HSpan: 2})
	panel.AddChild(label)
	sep := unison.NewSeparator()
	sep.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  2,
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	panel.AddChild(sep)
}

func (d *sheetSettingsDockable) reset() {
	if d.entity != nil {
		d.entity.SheetSettings = settings.Global().Sheet.Clone(d.entity)
	} else {
		settings.Global().Sheet = gurps.FactorySheetSettings()
	}
	d.sync()
}

func (d *sheetSettingsDockable) sync() {
	s := d.settings()
	d.damageProgressionPopup.Select(s.DamageProgression)
	d.showAdvantageModifier.State = unison.CheckStateFromBool(s.ShowAdvantageModifierAdj)
	d.showEquipmentModifier.State = unison.CheckStateFromBool(s.ShowEquipmentModifierAdj)
	d.showSpellAdjustments.State = unison.CheckStateFromBool(s.ShowSpellAdj)
	d.showTitleInsteadOfNameInPageFooter.State = unison.CheckStateFromBool(s.UseTitleInFooter)
	d.useMultiplicativeModifiers.State = unison.CheckStateFromBool(s.UseMultiplicativeModifiers)
	d.useModifyDicePlusAdds.State = unison.CheckStateFromBool(s.UseModifyingDicePlusAdds)
	d.lengthUnitsPopup.Select(s.DefaultLengthUnits)
	d.weightUnitsPopup.Select(s.DefaultWeightUnits)
	d.userDescDisplayPopup.Select(s.UserDescriptionDisplay)
	d.modifiersDisplayPopup.Select(s.ModifiersDisplay)
	d.notesDisplayPopup.Select(s.NotesDisplay)
	d.skillLevelAdjDisplayPopup.Select(s.SkillLevelAdjDisplay)
	d.MarkForRedraw()
}

func (d *sheetSettingsDockable) syncSheet() {
	if d.entity != nil {
		for _, wnd := range unison.Windows() {
			if ws := workspace.FromWindow(wnd); ws != nil {
				ws.DocumentDock.RootDockLayout().ForEachDockContainer(func(dc *unison.DockContainer) bool {
					for _, one := range dc.Dockables() {
						type signature interface {
							Entity() *gurps.Entity
							MarkForRebuild()
						}
						if s, ok := one.(signature); ok && s.Entity() == d.entity {
							s.MarkForRebuild()
						}
					}
					return false
				})
			}
		}
	}
}

func (d *sheetSettingsDockable) load(fileSystem fs.FS, filePath string) error {
	s, err := gurps.NewSheetSettingsFromFile(fileSystem, filePath)
	if err != nil {
		return err
	}
	if d.entity != nil {
		d.entity.SheetSettings = s
		s.SetOwningEntity(d.entity)
	} else {
		settings.Global().Sheet = s
	}
	d.sync()
	return nil
}

func (d *sheetSettingsDockable) save(filePath string) error {
	if d.entity != nil {
		return d.entity.SheetSettings.Save(filePath)
	}
	return settings.Global().Sheet.Save(filePath)
}
