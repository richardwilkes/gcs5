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

package workspace

import (
	"io/fs"

	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/ui/icons"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type menuKeySettingsDockable struct {
	SettingsDockable
	content *unison.Panel
}

// ShowMenuKeySettings shows the Menu Key settings.
func ShowMenuKeySettings() {
	ws, dc, found := Activate(func(d unison.Dockable) bool {
		_, ok := d.(*menuKeySettingsDockable)
		return ok
	})
	if !found && ws != nil {
		d := &menuKeySettingsDockable{}
		d.Self = d
		d.TabTitle = i18n.Text("Menu Keys")
		d.Extension = ".keys"
		d.Loader = d.load
		d.Saver = d.save
		d.Resetter = d.reset
		d.Setup(ws, dc, nil, nil, d.initContent)
	}
}

func (d *menuKeySettingsDockable) initContent(content *unison.Panel) {
	d.content = content
	d.content.SetLayout(&unison.FlexLayout{
		Columns:  7,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	d.fill()
}

func (d *menuKeySettingsDockable) reset() {
	g := settings.Global()
	g.KeyBindings.Reset()
	g.KeyBindings.MakeCurrent()
	d.sync()
}

func (d *menuKeySettingsDockable) sync() {
	d.content.RemoveAllChildren()
	d.fill()
	d.MarkForRedraw()
}

func (d *menuKeySettingsDockable) fill() {
	// TODO: Implement
}

func (d *menuKeySettingsDockable) createResetField(index int) {
	b := unison.NewSVGButton(icons.ResetSVG())
	b.Tooltip = unison.NewTooltipWithText("Reset this key binding")
	b.ClickCallback = func() {
		// TODO: Implement
	}
	b.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.MiddleAlignment,
		VAlign: unison.MiddleAlignment,
	})
	d.content.AddChild(b)
}

func (d *menuKeySettingsDockable) load(fileSystem fs.FS, filePath string) error {
	b, err := settings.NewKeyBindingsFromFS(fileSystem, filePath)
	if err != nil {
		return err
	}
	g := settings.Global()
	g.KeyBindings = b
	g.KeyBindings.MakeCurrent()
	d.sync()
	return nil
}

func (d *menuKeySettingsDockable) save(filePath string) error {
	return settings.Global().KeyBindings.Save(filePath)
}
