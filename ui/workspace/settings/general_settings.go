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
	"strconv"

	"github.com/richardwilkes/gcs/model/gurps"
	gsettings "github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/ui/icons"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/desktop"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

type generalSettingsDockable struct {
	Dockable
	nameField                           *unison.Field
	autoFillProfileCheckbox             *unison.CheckBox
	pointsField                         *widget.NumericField
	includeUnspentPointsInTotalCheckbox *unison.CheckBox
	techLevelField                      *unison.Field
	calendarPopup                       *unison.PopupMenu
	initialScaleField                   *widget.PercentageField
	exportResolutionField               *widget.IntegerField
	tooltipDelayField                   *widget.NumericField
	tooltipDismissalField               *widget.NumericField
	gCalcKeyField                       *unison.Field
}

// ShowGeneralSettings the General Settings window.
func ShowGeneralSettings() {
	ws, dc, found := Activate(func(d unison.Dockable) bool {
		_, ok := d.(*generalSettingsDockable)
		return ok
	})
	if !found && ws != nil {
		d := &generalSettingsDockable{}
		d.Self = d
		d.TabTitle = i18n.Text("General Settings")
		d.Extension = ".general"
		d.Loader = d.load
		d.Saver = d.save
		d.Resetter = d.reset
		d.Setup(ws, dc, nil, nil, d.initContent)
		d.nameField.RequestFocus()
	}
}

func (d *generalSettingsDockable) initContent(content *unison.Panel) {
	content.SetLayout(&unison.FlexLayout{
		Columns:  3,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	d.createPlayerAndDescFields(content)
	d.createInitialPointsFields(content)
	d.createTechLevelField(content)
	d.createCalendarPopup(content)
	d.createScaleField(content)
	d.createImageResolutionField(content)
	d.createTooltipDelayField(content)
	d.createTooltipDismissalField(content)
	d.createGCalcKeyField(content)
}

func (d *generalSettingsDockable) createPlayerAndDescFields(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Default Player Name")))
	d.nameField = widget.NewStringField(settings.Global().General.DefaultPlayerName, func(s string) {
		settings.Global().General.DefaultPlayerName = s
	})
	content.AddChild(d.nameField)
	d.autoFillProfileCheckbox = widget.NewCheckBox(i18n.Text("Fill in initial description"),
		settings.Global().General.AutoFillProfile, func(checked bool) { settings.Global().General.AutoFillProfile = checked })
	content.AddChild(d.autoFillProfileCheckbox)
}

func (d *generalSettingsDockable) createInitialPointsFields(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Initial Points")))
	d.pointsField = widget.NewNumericField(settings.Global().General.InitialPoints, gsettings.InitialPointsMin,
		gsettings.InitialPointsMax, func(v fixed.F64d4) { settings.Global().General.InitialPoints = v })
	content.AddChild(d.pointsField)
	d.includeUnspentPointsInTotalCheckbox = widget.NewCheckBox(i18n.Text("Include unspent points in total"),
		settings.Global().General.IncludeUnspentPointsInTotal,
		func(checked bool) { settings.Global().General.IncludeUnspentPointsInTotal = checked })
	content.AddChild(d.includeUnspentPointsInTotalCheckbox)
}

func (d *generalSettingsDockable) createTechLevelField(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Default Tech Level")))
	d.techLevelField = widget.NewStringField(settings.Global().General.DefaultTechLevel,
		func(s string) { settings.Global().General.DefaultTechLevel = s })
	d.techLevelField.Tooltip = unison.NewTooltipWithText(gurps.TechLevelInfo)
	content.AddChild(d.techLevelField)
	content.AddChild(unison.NewPanel())
}

func (d *generalSettingsDockable) createCalendarPopup(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Calendar")))
	d.calendarPopup = unison.NewPopupMenu()
	libraries := settings.Global().Libraries()
	for _, lib := range gsettings.AvailableCalendarRefs(libraries) {
		d.calendarPopup.AddDisabledItem(lib.Name)
		for _, one := range lib.List {
			d.calendarPopup.AddItem(one.Name)
		}
	}
	d.calendarPopup.Select(settings.Global().General.CalendarRef(libraries).Name)
	d.calendarPopup.SetLayoutData(&unison.FlexLayoutData{HSpan: 2})
	d.calendarPopup.SelectionCallback = func() {
		if name, ok := d.calendarPopup.Selected().(string); ok {
			settings.Global().General.CalendarName = name
		}
	}
	content.AddChild(d.calendarPopup)
}

func (d *generalSettingsDockable) createScaleField(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Initial Scale")))
	d.initialScaleField = widget.NewPercentageField(settings.Global().General.InitialUIScale,
		gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax,
		func(v fixed.F64d4) { settings.Global().General.InitialUIScale = v })
	content.AddChild(widget.WrapWithSpan(2, d.initialScaleField))
}

func (d *generalSettingsDockable) createImageResolutionField(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Image Export Resolution")))
	d.exportResolutionField = widget.NewIntegerField(settings.Global().General.ImageResolution,
		gsettings.ImageResolutionMin, gsettings.ImageResolutionMax,
		func(v int) { settings.Global().General.ImageResolution = v })
	content.AddChild(widget.WrapWithSpan(2, d.exportResolutionField, widget.NewFieldTrailingLabel(i18n.Text("ppi"))))
}

func (d *generalSettingsDockable) createTooltipDelayField(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Tooltip Delay")))
	d.tooltipDelayField = widget.NewNumericField(settings.Global().General.TooltipDelay, gsettings.TooltipDelayMin,
		gsettings.TooltipDelayMax, func(v fixed.F64d4) {
			s := settings.Global().General
			s.TooltipDelay = v
			s.UpdateToolTipTiming()
		})
	content.AddChild(widget.WrapWithSpan(2, d.tooltipDelayField, widget.NewFieldTrailingLabel(i18n.Text("seconds"))))
}

func (d *generalSettingsDockable) createTooltipDismissalField(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Tooltip Dismissal")))
	d.tooltipDismissalField = widget.NewNumericField(settings.Global().General.TooltipDismissal,
		gsettings.TooltipDismissalMin, gsettings.TooltipDismissalMax, func(v fixed.F64d4) {
			s := settings.Global().General
			s.TooltipDismissal = v
			s.UpdateToolTipTiming()
		})
	content.AddChild(widget.WrapWithSpan(2, d.tooltipDismissalField, widget.NewFieldTrailingLabel(i18n.Text("seconds"))))
}

func (d *generalSettingsDockable) createGCalcKeyField(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("GURPS Calculator Key")))
	button := unison.NewButton()
	button.HideBase = true
	baseline := button.Font.Baseline()
	button.Drawable = &unison.DrawableSVG{
		SVG:  icons.SearchSVG(),
		Size: geom32.NewSize(baseline, baseline),
	}
	button.ClickCallback = d.findGCalcKey
	d.gCalcKeyField = widget.NewStringField(settings.Global().General.GCalcKey, func(s string) {
		settings.Global().General.GCalcKey = s
	})
	content.AddChild(widget.WrapWithSpan(2, d.gCalcKeyField, button))
}

func (d *generalSettingsDockable) findGCalcKey() {
	if err := desktop.OpenBrowser("http://www.gurpscalculator.com/Character/ImportGCS"); err != nil {
		unison.ErrorDialogWithError(i18n.Text("Unable to open browser to determine GURPS Calculator Key"), err)
	}
}

func (d *generalSettingsDockable) reset() {
	*settings.Global().General = *gsettings.NewGeneral()
	d.sync()
}

func (d *generalSettingsDockable) sync() {
	s := settings.Global().General
	d.nameField.SetText(s.DefaultPlayerName)
	widget.SetCheckBoxState(d.autoFillProfileCheckbox, s.AutoFillProfile)
	d.pointsField.SetText(s.InitialPoints.String())
	widget.SetCheckBoxState(d.includeUnspentPointsInTotalCheckbox, s.IncludeUnspentPointsInTotal)
	d.techLevelField.SetText(s.DefaultTechLevel)
	d.calendarPopup.Select(s.CalendarRef(settings.Global().Libraries()).Name)
	d.initialScaleField.SetText(s.InitialUIScale.String() + "%")
	d.exportResolutionField.SetText(strconv.Itoa(s.ImageResolution))
	d.tooltipDelayField.SetText(s.TooltipDelay.String())
	d.tooltipDismissalField.SetText(s.TooltipDismissal.String())
	d.gCalcKeyField.SetText(s.GCalcKey)
	d.MarkForRedraw()
}

func (d *generalSettingsDockable) load(fileSystem fs.FS, filePath string) error {
	s, err := gsettings.NewGeneralFromFile(fileSystem, filePath)
	if err != nil {
		return err
	}
	*settings.Global().General = *s
	d.sync()
	return nil
}

func (d *generalSettingsDockable) save(filePath string) error {
	return settings.Global().General.Save(filePath)
}
