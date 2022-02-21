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

package colors

import (
	"fmt"
	"io/fs"

	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/ui/icons"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/gcs/ui/workspace/dsettings"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

type dockable struct {
	dsettings.Dockable
	content *unison.Panel
}

// Show the Color settings.
func Show() {
	ws, dc, found := dsettings.Activate(func(d unison.Dockable) bool {
		_, ok := d.(*dockable)
		return ok
	})
	if !found && ws != nil {
		d := &dockable{}
		d.Self = d
		d.TabTitle = i18n.Text("Colors")
		d.Extension = ".colors"
		d.Loader = d.load
		d.Saver = d.save
		d.Resetter = d.reset
		d.Setup(ws, dc, nil, nil, d.initContent)
	}
}

func (d *dockable) initContent(content *unison.Panel) {
	d.content = content
	d.content.SetLayout(&unison.FlexLayout{
		Columns:  8,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	d.fill()
}

func (d *dockable) reset() {
	g := settings.Global()
	g.Colors.Reset()
	g.Colors.MakeCurrent()
	d.sync()
}

func (d *dockable) sync() {
	d.content.RemoveAllChildren()
	d.fill()
	d.MarkForRedraw()
}

func (d *dockable) fill() {
	for i, one := range theme.CurrentColors {
		if i%2 == 0 {
			d.content.AddChild(widget.NewFieldLeadingLabel(one.Title))
		} else {
			d.content.AddChild(widget.NewFieldInteriorLeadingLabel(one.Title))
		}
		d.createColorWellField(one, true)
		d.createColorWellField(one, false)
		d.createResetField(one)
	}
}

func (d *dockable) createColorWellField(c *theme.ThemedColor, light bool) {
	w := unison.NewWell()
	w.Mask = unison.ColorWellMask
	if light {
		w.SetInk(c.Color.Light)
		w.Tooltip = unison.NewTooltipWithText(i18n.Text("The color to use when light mode is enabled"))
		w.InkChangedCallback = func() {
			if clr, ok := w.Ink().(unison.Color); ok {
				c.Color.Light = clr
				unison.ThemeChanged()
			}
		}
	} else {
		w.SetInk(c.Color.Dark)
		w.Tooltip = unison.NewTooltipWithText(i18n.Text("The color to use when dark mode is enabled"))
		w.InkChangedCallback = func() {
			if clr, ok := w.Ink().(unison.Color); ok {
				c.Color.Dark = clr
				unison.ThemeChanged()
			}
		}
	}
	d.content.AddChild(w)
}

func (d *dockable) createResetField(c *theme.ThemedColor) {
	b := unison.NewSVGButton(icons.ResetSVG())
	b.Tooltip = unison.NewTooltipWithText("Reset this color")
	b.ClickCallback = func() {
		if unison.QuestionDialog(fmt.Sprintf(i18n.Text("Are you sure you want to reset %s?"), c.Title), "") == unison.ModalResponseOK {
			for _, v := range theme.FactoryColors {
				if v.ID != c.ID {
					continue
				}
				*c.Color = *v.Color
				i := b.Parent().IndexOfChild(b)
				children := b.Parent().Children()
				if w, ok := children[i-2].Self.(*unison.Well); ok {
					w.SetInk(c.Color.Light)
				}
				if w, ok := children[i-1].Self.(*unison.Well); ok {
					w.SetInk(c.Color.Dark)
				}
				unison.ThemeChanged()
				break
			}
		}
	}
	b.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.MiddleAlignment,
		VAlign: unison.MiddleAlignment,
	})
	d.content.AddChild(b)
}

func (d *dockable) load(fileSystem fs.FS, filePath string) error {
	s, err := theme.NewColorsFromFS(fileSystem, filePath)
	if err != nil {
		return err
	}
	g := settings.Global()
	g.Colors = *s
	g.Colors.MakeCurrent()
	d.sync()
	return nil
}

func (d *dockable) save(filePath string) error {
	return settings.Global().Colors.Save(filePath)
}
