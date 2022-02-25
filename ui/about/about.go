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

package about

import (
	_ "embed"
	"runtime"

	"github.com/richardwilkes/gcs/ui/trampolines"
	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/cmdline"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom32"
	"github.com/richardwilkes/unison"
)

var (
	//go:embed "about.png"
	aboutImageData []byte
	aboutWnd       = &aboutWindow{}
)

type aboutWindow struct {
	*unison.Window
	img *unison.Image
}

// Show the about box.
func Show(_ unison.MenuItem) {
	if aboutWnd.Window == nil {
		if err := aboutWnd.prepare(); err != nil {
			jot.Error(err)
			return
		}
	}
	aboutWnd.ToFront()
}

func (w *aboutWindow) prepare() error {
	var err error
	if w.img == nil {
		if w.img, err = unison.NewImageFromBytes(aboutImageData, 0.5); err != nil {
			return errs.NewWithCause("unable to load about image", err)
		}
	}
	if w.Window, err = unison.NewWindow(i18n.Text("About GCS"), unison.NotResizableWindowOption()); err != nil {
		return errs.NewWithCause("unable to create about window", err)
	}
	trampolines.MenuSetup(w.Window)
	content := w.Content()
	content.DrawCallback = w.drawContentBackground
	content.SetLayout(w)
	w.Pack()
	r := w.ContentRect()
	usable := unison.PrimaryDisplay().Usable
	r.X = usable.X + (usable.Width-r.Width)/2
	r.Y = usable.Y + (usable.Height-r.Height)/2
	r.Point.Align()
	w.SetContentRect(r)
	return nil
}

func (w *aboutWindow) LayoutSizes(_ *unison.Panel, _ geom32.Size) (min, pref, max geom32.Size) {
	pref = w.img.LogicalSize()
	return pref, pref, pref
}

func (w *aboutWindow) PerformLayout(target *unison.Panel) {
	target.SetFrameRect(geom32.Rect{Size: w.img.LogicalSize()})
}

func (w *aboutWindow) drawContentBackground(gc *unison.Canvas, _ geom32.Rect) {
	r := w.Content().ContentRect(true)
	gc.DrawImageInRect(w.img, r, nil, nil)
	gc.DrawRect(r, unison.NewEvenlySpacedGradient(geom32.Point{Y: 0.25}, geom32.Point{Y: 1}, 0, 0,
		unison.Transparent, unison.Black).Paint(gc, r, unison.Fill))

	face := unison.MatchFontFace(unison.DefaultSystemFamilyName, unison.NormalFontWeight, unison.StandardSpacing, unison.NoSlant)
	paint := unison.RGB(128, 128, 128).Paint(gc, geom32.Rect{}, unison.Fill)
	text := unison.NewText(i18n.Text("This product includes copyrighted material from the GURPS game, which is used by permission of Steve Jackson Games."),
		&unison.TextDecoration{
			Font:  face.Font(7),
			Paint: paint,
		})
	const aboutMargin = 10
	y := r.Height - aboutMargin
	text.Draw(gc, (r.Width-text.Width())/2, y)
	y -= text.Height()
	font := face.Font(8)
	text = unison.NewText(i18n.Text("GURPS is a trademark of Steve Jackson Games, used by permission. All rights reserved."),
		&unison.TextDecoration{
			Font:  font,
			Paint: paint,
		})
	text.Draw(gc, (r.Width-text.Width())/2, y)
	lineHeight := text.Height()
	y -= lineHeight * 1.5
	yr := y

	paint = unison.RGB(204, 204, 204).Paint(gc, geom32.Rect{}, unison.Fill)
	unison.NewText(cmdline.Copyright(), &unison.TextDecoration{
		Font:  font,
		Paint: paint,
	}).Draw(gc, aboutMargin, y)
	y -= lineHeight
	if cmdline.BuildNumber != "" {
		unison.NewText(i18n.Text("Build ")+cmdline.BuildNumber, &unison.TextDecoration{
			Font:  font,
			Paint: paint,
		}).Draw(gc, aboutMargin, y)
		y -= lineHeight
	}
	var t string
	if cmdline.AppVersion != "" {
		t = i18n.Text("Version ") + cmdline.AppVersion
	} else {
		t = i18n.Text("Development Version")
	}
	unison.NewText(t, &unison.TextDecoration{
		Font: unison.MatchFontFace(unison.DefaultSystemFamilyName, unison.MediumFontWeight,
			unison.StandardSpacing, unison.NoSlant).Font(10),
		Paint: unison.White.Paint(gc, geom32.Rect{}, unison.Fill),
	}).Draw(gc, aboutMargin, y)

	right := r.Width - aboutMargin
	text = unison.NewText(runtime.GOARCH, &unison.TextDecoration{
		Font:  font,
		Paint: paint,
	})
	text.Draw(gc, right-text.Width(), yr)
	yr -= lineHeight
	switch runtime.GOOS {
	case toolbox.MacOS:
		t = "macOS"
	case toolbox.LinuxOS:
		t = "Linux"
	case toolbox.WindowsOS:
		t = "Windows"
	default:
		t = runtime.GOOS
	}
	text = unison.NewText(t, &unison.TextDecoration{
		Font:  font,
		Paint: paint,
	})
	text.Draw(gc, right-text.Width(), yr)
}
