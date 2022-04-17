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
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace"
	"github.com/richardwilkes/toolbox/desktop"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed/f64d4"
	"github.com/richardwilkes/unison"
)

type generalSettingsDockable struct {
	Dockable
	nameField                           *widget.StringField
	autoFillProfileCheckbox             *unison.CheckBox
	pointsField                         *widget.NumericField
	includeUnspentPointsInTotalCheckbox *unison.CheckBox
	techLevelField                      *widget.StringField
	calendarPopup                       *unison.PopupMenu[string]
	initialListScaleField               *widget.PercentageField
	initialSheetScaleField              *widget.PercentageField
	exportResolutionField               *widget.IntegerField
	tooltipDelayField                   *widget.NumericField
	tooltipDismissalField               *widget.NumericField
	gCalcKeyField                       *widget.StringField
}

// ShowGeneralSettings the General Settings window.
func ShowGeneralSettings() {
	ws, dc, found := workspace.Activate(func(d unison.Dockable) bool {
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
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Initial List Scale")))
	d.initialListScaleField = widget.NewPercentageField(
		func() int { return settings.Global().General.InitialListUIScale },
		func(v int) { settings.Global().General.InitialListUIScale = v },
		gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax)
	content.AddChild(widget.WrapWithSpan(2, d.initialListScaleField))
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Initial Sheet Scale")))
	d.initialSheetScaleField = widget.NewPercentageField(
		func() int { return settings.Global().General.InitialSheetUIScale },
		func(v int) { settings.Global().General.InitialSheetUIScale = v },
		gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax)
	content.AddChild(widget.WrapWithSpan(2, d.initialSheetScaleField))
	d.createImageResolutionField(content)
	d.createTooltipDelayField(content)
	d.createTooltipDismissalField(content)
	d.createGCalcKeyField(content)
}

func (d *generalSettingsDockable) createPlayerAndDescFields(content *unison.Panel) {
	title := i18n.Text("Default Player Name")
	content.AddChild(widget.NewFieldLeadingLabel(title))
	d.nameField = widget.NewStringField(title,
		func() string { return settings.Global().General.DefaultPlayerName },
		func(s string) { settings.Global().General.DefaultPlayerName = s })
	content.AddChild(d.nameField)
	d.autoFillProfileCheckbox = widget.NewCheckBox(i18n.Text("Fill in initial description"),
		settings.Global().General.AutoFillProfile,
		func(checked bool) { settings.Global().General.AutoFillProfile = checked })
	content.AddChild(d.autoFillProfileCheckbox)
}

func (d *generalSettingsDockable) createInitialPointsFields(content *unison.Panel) {
	title := i18n.Text("Initial Points")
	content.AddChild(widget.NewFieldLeadingLabel(title))
	d.pointsField = widget.NewNumericField(title,
		func() f64d4.Int { return settings.Global().General.InitialPoints },
		func(v f64d4.Int) { settings.Global().General.InitialPoints = v }, gsettings.InitialPointsMin,
		gsettings.InitialPointsMax, false)
	content.AddChild(d.pointsField)
	d.includeUnspentPointsInTotalCheckbox = widget.NewCheckBox(i18n.Text("Include unspent points in total"),
		settings.Global().General.IncludeUnspentPointsInTotal,
		func(checked bool) { settings.Global().General.IncludeUnspentPointsInTotal = checked })
	content.AddChild(d.includeUnspentPointsInTotalCheckbox)
}

func (d *generalSettingsDockable) createTechLevelField(content *unison.Panel) {
	title := i18n.Text("Default Tech Level")
	content.AddChild(widget.NewFieldLeadingLabel(title))
	d.techLevelField = widget.NewStringField(title,
		func() string { return settings.Global().General.DefaultTechLevel },
		func(s string) { settings.Global().General.DefaultTechLevel = s })
	d.techLevelField.Tooltip = unison.NewTooltipWithText(gurps.TechLevelInfo)
	content.AddChild(d.techLevelField)
	content.AddChild(unison.NewPanel())
}

func (d *generalSettingsDockable) createCalendarPopup(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Calendar")))
	d.calendarPopup = unison.NewPopupMenu[string]()
	libraries := settings.Global().Libraries()
	for _, lib := range gsettings.AvailableCalendarRefs(libraries) {
		d.calendarPopup.AddDisabledItem(lib.Name)
		for _, one := range lib.List {
			d.calendarPopup.AddItem(one.Name)
		}
	}
	d.calendarPopup.Select(settings.Global().General.CalendarRef(libraries).Name)
	d.calendarPopup.SetLayoutData(&unison.FlexLayoutData{HSpan: 2})
	d.calendarPopup.SelectionCallback = func(_ int, item string) {
		settings.Global().General.CalendarName = item
	}
	content.AddChild(d.calendarPopup)
}

func (d *generalSettingsDockable) createImageResolutionField(content *unison.Panel) {
	title := i18n.Text("Image Export Resolution")
	content.AddChild(widget.NewFieldLeadingLabel(title))
	d.exportResolutionField = widget.NewIntegerField(title,
		func() int { return settings.Global().General.ImageResolution },
		func(v int) { settings.Global().General.ImageResolution = v },
		gsettings.ImageResolutionMin, gsettings.ImageResolutionMax, false)
	content.AddChild(widget.WrapWithSpan(2, d.exportResolutionField, widget.NewFieldTrailingLabel(i18n.Text("ppi"))))
}

func (d *generalSettingsDockable) createTooltipDelayField(content *unison.Panel) {
	title := i18n.Text("Tooltip Delay")
	content.AddChild(widget.NewFieldLeadingLabel(title))
	d.tooltipDelayField = widget.NewNumericField(title, func() f64d4.Int { return settings.Global().General.TooltipDelay },
		func(v f64d4.Int) {
			general := settings.Global().General
			general.TooltipDelay = v
			general.UpdateToolTipTiming()
		}, gsettings.TooltipDelayMin, gsettings.TooltipDelayMax, false)
	content.AddChild(widget.WrapWithSpan(2, d.tooltipDelayField, widget.NewFieldTrailingLabel(i18n.Text("seconds"))))
}

func (d *generalSettingsDockable) createTooltipDismissalField(content *unison.Panel) {
	title := i18n.Text("Tooltip Dismissal")
	content.AddChild(widget.NewFieldLeadingLabel(title))
	d.tooltipDismissalField = widget.NewNumericField(title, func() f64d4.Int {
		return settings.Global().General.TooltipDismissal
	}, func(v f64d4.Int) {
		general := settings.Global().General
		general.TooltipDismissal = v
		general.UpdateToolTipTiming()
	}, gsettings.TooltipDismissalMin, gsettings.TooltipDismissalMax, false)
	content.AddChild(widget.WrapWithSpan(2, d.tooltipDismissalField, widget.NewFieldTrailingLabel(i18n.Text("seconds"))))
}

func (d *generalSettingsDockable) createGCalcKeyField(content *unison.Panel) {
	title := i18n.Text("GURPS Calculator Key")
	content.AddChild(widget.NewFieldLeadingLabel(title))
	button := unison.NewButton()
	button.HideBase = true
	baseline := button.Font.Baseline()
	button.Drawable = &unison.DrawableSVG{
		SVG:  res.SearchSVG,
		Size: unison.NewSize(baseline, baseline),
	}
	button.ClickCallback = d.findGCalcKey
	d.gCalcKeyField = widget.NewStringField(title,
		func() string { return settings.Global().General.GCalcKey },
		func(s string) { settings.Global().General.GCalcKey = s })
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
	d.initialListScaleField.Set(s.InitialListUIScale)
	d.initialSheetScaleField.Set(s.InitialSheetUIScale)
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
