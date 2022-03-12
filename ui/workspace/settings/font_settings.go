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
	"fmt"
	"io/fs"

	"github.com/richardwilkes/gcs/model/settings"
	"github.com/richardwilkes/gcs/model/theme"
	"github.com/richardwilkes/gcs/res"
	"github.com/richardwilkes/gcs/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xmath/fixed"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

type fontSettingsDockable struct {
	Dockable
	content  *unison.Panel
	noUpdate bool
}

// ShowFontSettings shows the Font settings.
func ShowFontSettings() {
	ws, dc, found := Activate(func(d unison.Dockable) bool {
		_, ok := d.(*fontSettingsDockable)
		return ok
	})
	if !found && ws != nil {
		d := &fontSettingsDockable{}
		d.Self = d
		d.TabTitle = i18n.Text("Fonts")
		d.Extension = ".fonts"
		d.Loader = d.load
		d.Saver = d.save
		d.Resetter = d.reset
		d.Setup(ws, dc, nil, nil, d.initContent)
	}
}

func (d *fontSettingsDockable) initContent(content *unison.Panel) {
	d.content = content
	d.content.SetLayout(&unison.FlexLayout{
		Columns:  7,
		HSpacing: unison.StdHSpacing,
		VSpacing: unison.StdVSpacing,
	})
	d.fill()
}

func (d *fontSettingsDockable) reset() {
	g := settings.Global()
	g.Fonts.Reset()
	g.Fonts.MakeCurrent()
	d.sync()
}

func (d *fontSettingsDockable) sync() {
	d.content.RemoveAllChildren()
	d.fill()
	d.MarkForRedraw()
}

func (d *fontSettingsDockable) fill() {
	for i, one := range theme.CurrentFonts {
		if i%2 == 0 {
			d.content.AddChild(widget.NewFieldLeadingLabel(one.Title))
		} else {
			d.content.AddChild(widget.NewFieldInteriorLeadingLabel(one.Title))
		}
		d.createFamilyField(i)
		d.createSizeField(i)
		d.createWeightField(i)
		d.createSpacingField(i)
		d.createSlantField(i)
		d.createResetField(i)
	}
	notice := unison.NewLabel()
	notice.Text = "Changing fonts usually requires restarting the app to see content laid out correctly."
	notice.Font = unison.SystemFont
	notice.SetBorder(unison.NewEmptyBorder(geom32.Insets{Top: unison.StdVSpacing * 2}))
	notice.SetLayoutData(&unison.FlexLayoutData{
		HSpan:  7,
		VSpan:  1,
		HAlign: unison.MiddleAlignment,
		VAlign: unison.MiddleAlignment,
	})
	d.content.AddChild(notice)
}

func (d *fontSettingsDockable) createFamilyField(index int) {
	p := unison.NewPopupMenu()
	for _, fam := range unison.FontFamilies() {
		p.AddItem(fam)
	}
	p.Select(theme.CurrentFonts[index].Font.Descriptor().Family)
	p.SelectionCallback = func() {
		if d.noUpdate {
			return
		}
		fd := theme.CurrentFonts[index].Font.Descriptor()
		if s, ok := p.Selected().(string); ok {
			fd.Family = s
			d.applyFont(index, fd)
		}
	}
	d.content.AddChild(p)
}

func (d *fontSettingsDockable) createSizeField(index int) {
	field := widget.NewNumericField(fixed.F64d4FromFloat32(theme.CurrentFonts[index].Font.Size()), fixed.F64d4One,
		fixed.F64d4FromInt(999), false, func(v fixed.F64d4) {
			if d.noUpdate {
				return
			}
			fd := theme.CurrentFonts[index].Font.Descriptor()
			fd.Size = v.AsFloat32()
			d.applyFont(index, fd)
		})
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	d.content.AddChild(field)
}

func (d *fontSettingsDockable) createWeightField(index int) {
	p := unison.NewPopupMenu()
	for _, s := range unison.FontWeights {
		p.AddItem(s)
	}
	p.Select(theme.CurrentFonts[index].Font.Descriptor().Weight)
	p.SelectionCallback = func() {
		if d.noUpdate {
			return
		}
		fd := theme.CurrentFonts[index].Font.Descriptor()
		if w, ok := p.Selected().(unison.FontWeight); ok {
			fd.Weight = w
			d.applyFont(index, fd)
		}
	}
	d.content.AddChild(p)
}

func (d *fontSettingsDockable) createSpacingField(index int) {
	p := unison.NewPopupMenu()
	for _, s := range unison.Spacings {
		p.AddItem(s)
	}
	p.Select(theme.CurrentFonts[index].Font.Descriptor().Spacing)
	p.SelectionCallback = func() {
		if d.noUpdate {
			return
		}
		fd := theme.CurrentFonts[index].Font.Descriptor()
		if s, ok := p.Selected().(unison.FontSpacing); ok {
			fd.Spacing = s
			d.applyFont(index, fd)
		}
	}
	d.content.AddChild(p)
}

func (d *fontSettingsDockable) createSlantField(index int) {
	p := unison.NewPopupMenu()
	for _, s := range unison.Slants {
		p.AddItem(s)
	}
	p.Select(theme.CurrentFonts[index].Font.Descriptor().Slant)
	p.SelectionCallback = func() {
		if d.noUpdate {
			return
		}
		fd := theme.CurrentFonts[index].Font.Descriptor()
		if s, ok := p.Selected().(unison.FontSlant); ok {
			fd.Slant = s
			d.applyFont(index, fd)
		}
	}
	d.content.AddChild(p)
}

func (d *fontSettingsDockable) createResetField(index int) {
	b := unison.NewSVGButton(res.ResetSVG)
	b.Tooltip = unison.NewTooltipWithText("Reset this font")
	b.ClickCallback = func() {
		if unison.QuestionDialog(fmt.Sprintf(i18n.Text("Are you sure you want to reset %s?"),
			theme.CurrentFonts[index].Title), "") == unison.ModalResponseOK {
			for _, v := range theme.FactoryFonts {
				if v.ID != theme.CurrentFonts[index].ID {
					continue
				}
				d.applyFont(index, v.Font.Descriptor())
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

func (d *fontSettingsDockable) applyFont(index int, fd unison.FontDescriptor) {
	theme.CurrentFonts[index].Font.Font = fd.Font()
	children := d.content.Children()
	i := index * 7
	fd = theme.CurrentFonts[index].Font.Descriptor()
	d.noUpdate = true
	if p, ok := children[i+1].Self.(*unison.PopupMenu); ok {
		p.Select(fd.Family)
	}
	if nf, ok := children[i+2].Self.(*widget.NumericField); ok {
		nf.SetText(fixed.F64d4FromFloat32(fd.Size).String())
	}
	if p, ok := children[i+3].Self.(*unison.PopupMenu); ok {
		p.Select(fd.Weight)
	}
	if p, ok := children[i+4].Self.(*unison.PopupMenu); ok {
		p.Select(fd.Spacing)
	}
	if p, ok := children[i+5].Self.(*unison.PopupMenu); ok {
		p.Select(fd.Slant)
	}
	d.noUpdate = false
	unison.ThemeChanged()
}

func (d *fontSettingsDockable) load(fileSystem fs.FS, filePath string) error {
	s, err := theme.NewFontsFromFS(fileSystem, filePath)
	if err != nil {
		return err
	}
	g := settings.Global()
	g.Fonts = *s
	g.Fonts.MakeCurrent()
	d.sync()
	return nil
}

func (d *fontSettingsDockable) save(filePath string) error {
	return settings.Global().Fonts.Save(filePath)
}
