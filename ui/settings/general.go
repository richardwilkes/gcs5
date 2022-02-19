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
	"github.com/richardwilkes/gcs/model/fxp"
	gsettings "github.com/richardwilkes/gcs/model/gurps/settings"
	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/ui/icons"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var generalSettingsWindow *unison.Window

// ShowGeneralSettings shows the General Settings window.
func ShowGeneralSettings() {
	if generalSettingsWindow == nil {
		var err error
		if generalSettingsWindow, err = unison.NewWindow(i18n.Text("General Settings")); err != nil {
			unison.ErrorDialogWithMessage(i18n.Text("Unable to open General Settings"), err.Error())
			return
		}
		generalSettingsWindow.WillCloseCallback = func() { generalSettingsWindow = nil }
		content := generalSettingsWindow.Content()
		content.SetLayout(&unison.FlexLayout{
			Columns:  3,
			HSpacing: unison.StdHSpacing,
			VSpacing: unison.StdVSpacing,
		})
		content.SetBorder(unison.NewEmptyBorder(geom32.NewUniformInsets(unison.StdHSpacing * 2)))

		generalSettings := settings.Global().General

		content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Default Player Name")))
		content.AddChild(widget.NewStringField(generalSettings.DefaultPlayerName, func(s string) {
			generalSettings.DefaultPlayerName = s
		}))
		checkbox := unison.NewCheckBox()
		checkbox.Text = i18n.Text("Fill in initial description")
		if generalSettings.AutoFillProfile {
			checkbox.State = unison.OnCheckState
		}
		checkbox.ClickCallback = func() {
			generalSettings.AutoFillProfile = checkbox.State == unison.OnCheckState
		}
		content.AddChild(checkbox)

		content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Initial Points")))
		content.AddChild(widget.NewNumericField(generalSettings.InitialPoints, 0, fixed.F64d4FromInt(9999999),
			func(v fixed.F64d4) { generalSettings.InitialPoints = v }))
		checkbox = unison.NewCheckBox()
		checkbox.Text = i18n.Text("Include unspent points in total")
		if generalSettings.IncludeUnspentPointsInTotal {
			checkbox.State = unison.OnCheckState
		}
		checkbox.ClickCallback = func() {
			generalSettings.IncludeUnspentPointsInTotal = checkbox.State == unison.OnCheckState
		}
		content.AddChild(checkbox)

		content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Default Tech Level")))
		content.AddChild(widget.NewStringField(generalSettings.DefaultTechLevel,
			func(s string) { generalSettings.DefaultTechLevel = s }))
		content.AddChild(unison.NewPanel())

		content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Calendar")))
		popup := unison.NewPopupMenu()
		libraries := settings.Global().Libraries()
		for _, lib := range gsettings.AvailableCalendarRefs(libraries) {
			popup.AddDisabledItem(lib.Name)
			for _, one := range lib.List {
				popup.AddItem(one.Name)
			}
		}
		popup.Select(generalSettings.CalendarRef(libraries).Name)
		popup.SetLayoutData(&unison.FlexLayoutData{HSpan: 2})
		popup.SelectionCallback = func() {
			if name, ok := popup.Selected().(string); ok {
				generalSettings.CalendarName = name
			}
		}
		content.AddChild(popup)

		content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Initial Scale")))
		percentField := widget.NewPercentageField(generalSettings.InitialUIScale, fixed.F64d4FromInt(25),
			fixed.F64d4FromInt(999), func(v fixed.F64d4) { generalSettings.InitialUIScale = v })
		content.AddChild(percentField)
		content.AddChild(unison.NewPanel())

		content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Image Export Resolution")))
		content.AddChild(widget.WrapWithSpan(2, widget.NewIntegerField(generalSettings.ImageResolution, 50, 300,
			func(v int) { generalSettings.ImageResolution = v }), widget.NewFieldTrailingLabel(i18n.Text("ppi"))))

		content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Tooltip Delay")))
		content.AddChild(widget.WrapWithSpan(2, widget.NewNumericField(generalSettings.ToolTipDelay, 0, fxp.Thirty,
			func(v fixed.F64d4) {
				generalSettings.ToolTipDelay = v
				generalSettings.UpdateToolTipTiming()
			}), widget.NewFieldTrailingLabel(i18n.Text("seconds"))))

		content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Tooltip Dismissal")))
		content.AddChild(widget.WrapWithSpan(2, widget.NewNumericField(generalSettings.ToolTipDismissal, 0,
			fixed.F64d4FromInt(3600), func(v fixed.F64d4) {
				generalSettings.ToolTipDismissal = v
				generalSettings.UpdateToolTipTiming()
			}), widget.NewFieldTrailingLabel(i18n.Text("seconds"))))

		content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("GURPS Calculator Key")))
		button := unison.NewButton()
		button.HideBase = true
		baseline := button.Font.Baseline()
		button.Drawable = &unison.DrawableSVG{
			SVG:  icons.SearchSVG(),
			Size: geom32.NewSize(baseline, baseline),
		}
		button.ClickCallback = func() {
			// TODO: Implement
			jot.Info("handle click callback for GURPS Calculator Key lookup")
		}
		content.AddChild(widget.WrapWithSpan(2, widget.NewStringField(generalSettings.GCalcKey, func(s string) { generalSettings.GCalcKey = s }), button))

		generalSettingsWindow.Pack()
	}
	generalSettingsWindow.ToFront()
}
