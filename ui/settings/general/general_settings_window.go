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

package general

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/richardwilkes/gcs/model/gurps"
	"github.com/richardwilkes/gcs/model/gurps/library"
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

type wndData struct {
	wnd                                 *unison.Window
	resetButton                         *unison.Button
	menuButton                          *unison.Button
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

var data *wndData

// Show the General Settings window.
func Show() {
	if data == nil {
		wnd, err := unison.NewWindow(i18n.Text("General Settings"))
		if err != nil {
			unison.ErrorDialogWithMessage(i18n.Text("Unable to open General Settings"), err.Error())
			return
		}
		wnd.WillCloseCallback = func() { data = nil }
		data = &wndData{
			wnd: wnd,
		}
		content := data.wnd.Content()
		content.SetLayout(&unison.FlexLayout{Columns: 1})
		content.AddChild(data.createToolbar())
		content.AddChild(data.createContent())
		data.wnd.Pack()
		data.nameField.RequestFocus()
	}
	data.wnd.ToFront()
}

func (d *wndData) createToolbar() *unison.Panel {
	toolbar := unison.NewPanel()
	toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, geom32.Insets{Bottom: 1}, false),
		unison.NewEmptyBorder(geom32.Insets{
			Top:    unison.StdVSpacing,
			Left:   unison.StdHSpacing,
			Bottom: unison.StdVSpacing,
			Right:  unison.StdHSpacing,
		})))
	spacer := unison.NewPanel()
	spacer.SetLayoutData(&unison.FlexLayoutData{HGrab: true})
	toolbar.AddChild(spacer)
	d.resetButton = unison.NewSVGButton(icons.ResetSVG())
	d.resetButton.ClickCallback = d.reset
	toolbar.AddChild(d.resetButton)
	d.menuButton = unison.NewSVGButton(icons.MenuSVG())
	d.menuButton.ClickCallback = d.showMenu
	toolbar.AddChild(d.menuButton)
	toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})
	return toolbar
}

func (d *wndData) createContent() *unison.Panel {
	content := unison.NewPanel()
	content.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	content.SetLayout(&unison.FlexLayout{
		Columns:  3,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	content.SetBorder(unison.NewEmptyBorder(geom32.NewUniformInsets(unison.StdHSpacing * 2)))
	d.createPlayerAndDescFields(content)
	d.createInitialPointsFields(content)
	d.createTechLevelField(content)
	d.createCalendarPopup(content)
	d.createScaleField(content)
	d.createImageResolutionField(content)
	d.createTooltipDelayField(content)
	d.createTooltipDismissalField(content)
	d.createGCalcKeyField(content)
	return content
}

func (d *wndData) createPlayerAndDescFields(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Default Player Name")))
	d.nameField = widget.NewStringField(settings.Global().General.DefaultPlayerName, func(s string) {
		settings.Global().General.DefaultPlayerName = s
	})
	content.AddChild(d.nameField)
	d.autoFillProfileCheckbox = widget.NewCheckBox(i18n.Text("Fill in initial description"),
		settings.Global().General.AutoFillProfile, func(checked bool) { settings.Global().General.AutoFillProfile = checked })
	content.AddChild(d.autoFillProfileCheckbox)
}

func (d *wndData) createInitialPointsFields(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Initial Points")))
	d.pointsField = widget.NewNumericField(settings.Global().General.InitialPoints, gsettings.InitialPointsMin,
		gsettings.InitialPointsMax, func(v fixed.F64d4) { settings.Global().General.InitialPoints = v })
	content.AddChild(d.pointsField)
	d.includeUnspentPointsInTotalCheckbox = widget.NewCheckBox(i18n.Text("Include unspent points in total"),
		settings.Global().General.IncludeUnspentPointsInTotal,
		func(checked bool) { settings.Global().General.IncludeUnspentPointsInTotal = checked })
	content.AddChild(d.includeUnspentPointsInTotalCheckbox)
}

func (d *wndData) createTechLevelField(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Default Tech Level")))
	d.techLevelField = widget.NewStringField(settings.Global().General.DefaultTechLevel,
		func(s string) { settings.Global().General.DefaultTechLevel = s })
	d.techLevelField.Tooltip = unison.NewTooltipWithText(gurps.TechLevelInfo)
	content.AddChild(d.techLevelField)
	content.AddChild(unison.NewPanel())
}

func (d *wndData) createCalendarPopup(content *unison.Panel) {
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

func (d *wndData) createScaleField(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Initial Scale")))
	d.initialScaleField = widget.NewPercentageField(settings.Global().General.InitialUIScale,
		gsettings.InitialUIScaleMin, gsettings.InitialUIScaleMax,
		func(v fixed.F64d4) { settings.Global().General.InitialUIScale = v })
	content.AddChild(widget.WrapWithSpan(2, d.initialScaleField))
}

func (d *wndData) createImageResolutionField(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Image Export Resolution")))
	d.exportResolutionField = widget.NewIntegerField(settings.Global().General.ImageResolution,
		gsettings.ImageResolutionMin, gsettings.ImageResolutionMax,
		func(v int) { settings.Global().General.ImageResolution = v })
	content.AddChild(widget.WrapWithSpan(2, d.exportResolutionField, widget.NewFieldTrailingLabel(i18n.Text("ppi"))))
}

func (d *wndData) createTooltipDelayField(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Tooltip Delay")))
	d.tooltipDelayField = widget.NewNumericField(settings.Global().General.TooltipDelay, gsettings.TooltipDelayMin,
		gsettings.TooltipDelayMax, func(v fixed.F64d4) {
			s := settings.Global().General
			s.TooltipDelay = v
			s.UpdateToolTipTiming()
		})
	content.AddChild(widget.WrapWithSpan(2, d.tooltipDelayField, widget.NewFieldTrailingLabel(i18n.Text("seconds"))))
}

func (d *wndData) createTooltipDismissalField(content *unison.Panel) {
	content.AddChild(widget.NewFieldLeadingLabel(i18n.Text("Tooltip Dismissal")))
	d.tooltipDismissalField = widget.NewNumericField(settings.Global().General.TooltipDismissal,
		gsettings.TooltipDismissalMin, gsettings.TooltipDismissalMax, func(v fixed.F64d4) {
			s := settings.Global().General
			s.TooltipDismissal = v
			s.UpdateToolTipTiming()
		})
	content.AddChild(widget.WrapWithSpan(2, d.tooltipDismissalField, widget.NewFieldTrailingLabel(i18n.Text("seconds"))))
}

func (d *wndData) createGCalcKeyField(content *unison.Panel) {
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

func (d *wndData) findGCalcKey() {
	if err := desktop.OpenBrowser("http://www.gurpscalculator.com/Character/ImportGCS"); err != nil {
		unison.ErrorDialogWithMessage(i18n.Text("Unable to open browser to determine GURPS Calculator Key"), err.Error())
	}
}

func (d *wndData) reset() {
	*settings.Global().General = *gsettings.NewGeneral()
	d.sync()
}

func (d *wndData) sync() {
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
	d.wnd.MarkForRedraw()
}

func (d *wndData) showMenu() {
	f := unison.DefaultMenuFactory()
	id := unison.ContextMenuIDFlag
	m := f.NewMenu(id, "", nil)
	id++
	m.InsertItem(-1, f.NewItem(id, i18n.Text("Import…"), 0, 0, nil, d.handleImport))
	id++
	m.InsertItem(-1, f.NewItem(id, i18n.Text("Export…"), 0, 0, nil, d.handleExport))
	id++
	libraries := settings.Global().Libraries()
	sets := library.ScanForNamedFileSets(nil, "", ".general", false, libraries)
	if len(sets) != 0 {
		m.InsertSeparator(-1, false)
		for _, lib := range sets {
			m.InsertItem(-1, f.NewItem(id, lib.Name, 0, 0, func(_ unison.MenuItem) bool { return false }, nil))
			id++
			for _, one := range lib.List {
				d.insertFileToLoad(m, id, one)
				id++
			}
		}
	}
	m.Popup(d.menuButton.RectToRoot(d.menuButton.ContentRect(true)), 0)
}

func (d *wndData) insertFileToLoad(m unison.Menu, id int, ref *library.NamedFileRef) {
	m.InsertItem(-1, m.Factory().NewItem(id, ref.Name, 0, 0, nil, func(_ unison.MenuItem) {
		s, err := gsettings.NewGeneralFromFile(ref.FileSystem, ref.FilePath)
		if err != nil {
			unison.ErrorDialogWithMessage(i18n.Text("Unable to load general settings"), err.Error())
			return
		}
		*settings.Global().General = *s
		d.sync()
	}))
}

func (d *wndData) handleImport(_ unison.MenuItem) {
	dialog := unison.NewOpenDialog()
	dialog.SetResolvesAliases(true)
	dialog.SetAllowedExtensions(".general")
	dialog.SetAllowsMultipleSelection(false)
	dialog.SetCanChooseDirectories(false)
	dialog.SetCanChooseFiles(true)
	if dialog.RunModal() {
		p := dialog.Path()
		s, err := gsettings.NewGeneralFromFile(os.DirFS(filepath.Dir(p)), filepath.Base(p))
		if err != nil {
			unison.ErrorDialogWithMessage(i18n.Text("Unable to load general settings"), err.Error())
			return
		}
		*settings.Global().General = *s
		d.sync()
	}
}

func (d *wndData) handleExport(_ unison.MenuItem) {
	dialog := unison.NewSaveDialog()
	dialog.SetAllowedExtensions(".general")
	if dialog.RunModal() {
		if err := settings.Global().General.Save(dialog.Path()); err != nil {
			unison.ErrorDialogWithMessage(i18n.Text("Unable to save general settings"), err.Error())
		}
	}
}
